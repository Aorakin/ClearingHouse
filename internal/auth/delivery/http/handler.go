package http

import (
	"net/http"
	"strings"

	"github.com/ClearingHouse/internal/auth/interfaces"
	apiError "github.com/ClearingHouse/pkg/api_error"
	"github.com/ClearingHouse/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	authUsecase interfaces.AuthUsecase
}

func NewAuthHandler(authUsecase interfaces.AuthUsecase) interfaces.AuthHandler {
	return &AuthHandler{
		authUsecase: authUsecase,
	}
}

func (h *AuthHandler) GoogleLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		url := h.authUsecase.GenerateGoogleLoginURL("login")
		c.Redirect(http.StatusTemporaryRedirect, url)
	}
}

func (h *AuthHandler) GoogleRegister() gin.HandlerFunc {
	return func(c *gin.Context) {
		url := h.authUsecase.GenerateGoogleRegisterURL("register")
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

		// Handle registration flow
		if state == "register" {
			user, err := h.authUsecase.HandleGoogleRegisterCallback(code, c)
			if err != nil {
				c.JSON(response.ErrorResponseBuilder(apiError.NewInternalServerError(err)))
				return
			}
			accessToken, refreshToken, err := h.authUsecase.GenerateTokens(user)
			if err != nil {
				c.JSON(response.ErrorResponseBuilder(apiError.NewInternalServerError(err)))
				return
			}
			c.SetCookie("access_token", accessToken, 3600, "/", "localhost", false, true)
			c.SetCookie("refresh_token", refreshToken, 7*24*3600, "/", "localhost", false, true)
			c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "access_token": accessToken, "refresh_token": refreshToken})
			return
		}

		// Handle login flow
		user, err := h.authUsecase.HandleGoogleCallback(code, c)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewInternalServerError(err)))
			return
		}

		accessToken, refreshToken, err := h.authUsecase.GenerateTokens(user)
		if err != nil {
			c.JSON(response.ErrorResponseBuilder(apiError.NewInternalServerError(err)))
			return
		}
		c.SetCookie("access_token", accessToken, 3600, "/", ".localhost", true, true)
		c.SetCookie("refresh_token", refreshToken, 7*24*3600, "/", ".localhost", true, true)

		c.JSON(http.StatusOK, gin.H{"access_token": accessToken, "refresh_token": refreshToken})
	}
}

func (h *AuthHandler) Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Clear the access token cookie
		c.SetCookie("access_token", "", -1, "/", ".localhost", true, true)
		// Clear the refresh token cookie
		c.SetCookie("refresh_token", "", -1, "/", ".localhost", true, true)
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
