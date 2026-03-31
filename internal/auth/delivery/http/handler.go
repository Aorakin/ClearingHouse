package http

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/ClearingHouse/internal/auth/dtos"
	"github.com/ClearingHouse/internal/auth/interfaces"
	apiError "github.com/ClearingHouse/pkg/api_error"
	"github.com/ClearingHouse/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	authUsecase interfaces.AuthUsecase
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func NewAuthHandler(authUsecase interfaces.AuthUsecase) interfaces.AuthHandler {
	return &AuthHandler{
		authUsecase: authUsecase,
	}
}

func (h *AuthHandler) GoogleLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		portal := c.DefaultQuery("portal", "portal")
		state := "login:" + portal
		url := h.authUsecase.GenerateGoogleLoginURL(state, portal)
		log.Printf("State: %s and URL: %s", state, url)
		c.Redirect(http.StatusTemporaryRedirect, url)
	}
}

func (h *AuthHandler) GoogleRegister() gin.HandlerFunc {
	return func(c *gin.Context) {
		portal := c.DefaultQuery("portal", "portal")
		state := "register:" + portal
		url := h.authUsecase.GenerateGoogleRegisterURL(state, portal)
		c.Redirect(http.StatusTemporaryRedirect, url)
	}
}

func (h *AuthHandler) GoogleCallback() gin.HandlerFunc {
	return func(c *gin.Context) {
		code := c.Query("code")
		state := c.Query("state")

		if code == "" {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError("no code provided")))
			return
		}

		// state format: "{type}:{portal}" e.g. "login:admin", "register:portal"
		stateType, portal := parseState(state)

		// Handle registration flow
		if stateType == "register" {
			user, err := h.authUsecase.HandleGoogleRegisterCallback(code, portal, c)
			if err != nil {
				c.JSON(response.ErrorResponseBuilder(apiError.NewInternalServerError(err)))
				return
			}
			accessToken, refreshToken, err := h.authUsecase.GenerateTokens(user)
			if err != nil {
				c.JSON(response.ErrorResponseBuilder(apiError.NewInternalServerError(err)))
				return
			}
			prodCookieDomain := getEnv("AUTH_COOKIE_DOMAIN", ".onepointfive.life")
			refreshCookiePath := getEnv("AUTH_REFRESH_COOKIE_PATH", "/users/auth")
			cookieSecure := getEnvBool("AUTH_COOKIE_SECURE", true)
			c.SetCookie("access_token", accessToken, 3600, "/", ".localhost", true, true)
			c.SetCookie("refresh_token", refreshToken, 7*24*3600, "/", ".localhost", true, true)
			c.SetCookie("access_token", accessToken, 7*24*3600, "/", prodCookieDomain, cookieSecure, true)
			c.SetCookie("refresh_token", refreshToken, 7*24*3600, refreshCookiePath, prodCookieDomain, cookieSecure, true)
			c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "access_token": accessToken, "refresh_token": refreshToken})
			return
		}

		// Handle login flow
		user, err := h.authUsecase.HandleGoogleCallback(code, portal, c)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewInternalServerError(err)))
			return
		}

		accessToken, refreshToken, err := h.authUsecase.GenerateTokens(user)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewInternalServerError(err)))
			return
		}
		prodCookieDomain := getEnv("AUTH_COOKIE_DOMAIN", ".onepointfive.life")
		refreshCookiePath := getEnv("AUTH_REFRESH_COOKIE_PATH", "/users/auth")
		cookieSecure := getEnvBool("AUTH_COOKIE_SECURE", true)
		c.SetCookie("access_token", accessToken, 3600, "/", ".localhost", true, true)
		c.SetCookie("refresh_token", refreshToken, 7*24*3600, "/", ".localhost", true, true)
		c.SetCookie("access_token", accessToken, 7*24*3600, "/", prodCookieDomain, cookieSecure, true)
		c.SetCookie("refresh_token", refreshToken, 7*24*3600, refreshCookiePath, prodCookieDomain, cookieSecure, true)

		c.JSON(http.StatusOK, gin.H{"access_token": accessToken, "refresh_token": refreshToken})
	}
}

// parseState splits a state string in "type:portal" format.
// Returns (stateType, portal), defaulting to ("login", "portal") on malformed input.
func parseState(state string) (stateType, portal string) {
	parts := strings.SplitN(state, ":", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return state, "portal"
}

// ManualRegister handles manual user registration for testing without OAuth
func (h *AuthHandler) ManualRegister() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dtos.ManualRegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err.Error())))
			return
		}

		user, err := h.authUsecase.ManualRegister(req.Email, req.FirstName, req.LastName)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError(err.Error())))
			return
		}

		accessToken, refreshToken, err := h.authUsecase.GenerateTokens(user)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewInternalServerError(err)))
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message":       "User registered successfully",
			"access_token":  accessToken,
			"refresh_token": refreshToken,
			"user":          user,
		})
	}
}

func (h *AuthHandler) Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get refresh token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
				refreshToken := parts[1]
				// Blacklist the refresh token
				if err := h.authUsecase.BlacklistToken(refreshToken); err != nil {
					// Log error but don't fail logout
					c.JSON(response.ErrorResponseBuilder(apiError.NewInternalServerError(err)))
					return
				}
			}
		}

		prodCookieDomain := getEnv("AUTH_COOKIE_DOMAIN", ".onepointfive.life")
		cookieSecure := getEnvBool("AUTH_COOKIE_SECURE", true)

		// Clear the access token cookie
		c.SetCookie("access_token", "", -1, "/", prodCookieDomain, cookieSecure, true)
		// Clear the refresh token cookie
		c.SetCookie("refresh_token", "", -1, "/", prodCookieDomain, cookieSecure, true)
		c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
	}
}

func (h *AuthHandler) RefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError("no Authorization header provided")))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.JSON(response.ErrorResponseBuilder(apiError.NewBadRequestError("invalid Authorization header format")))
			return
		}

		refreshToken := parts[1]

		newAccessToken, err := h.authUsecase.RefreshAccessToken(refreshToken)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewInternalServerError(err)))
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"access_token": newAccessToken,
		})
	}
}

func (h *AuthHandler) GetMe() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uuid.UUID)
		if userID == uuid.Nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewUnauthorizedError("unauthorized")))
			return
		}

		user, err := h.authUsecase.GetUserByID(userID)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewInternalServerError(err)))
			return
		}

		c.JSON(http.StatusOK, user)
	}
}
