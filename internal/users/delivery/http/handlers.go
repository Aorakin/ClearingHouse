package http

import (
	"net/http"
	"strings"

	"github.com/ClearingHouse/internal/users/interfaces"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type UsersHandlers struct {
	usersUsecase interfaces.UsersUsecase
}

func NewUsersHandler(usersUsecase interfaces.UsersUsecase) interfaces.UsersHandlers {
	return &UsersHandlers{
		usersUsecase: usersUsecase,
	}
}

func (h *UsersHandlers) LoginWithGoogle() gin.HandlerFunc {
	return func(c *gin.Context) {
		portal := c.DefaultQuery("portal", "portal")
		state := "login:" + portal
		url := h.usersUsecase.GenerateLoginURL(state, portal)
		c.Redirect(http.StatusTemporaryRedirect, url)
	}
}

func (h *UsersHandlers) Callback() gin.HandlerFunc {
	return func(c *gin.Context) {
		code := c.Query("code")
		state := c.Query("state")
		if code == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No code provided"})
			return
		}

		// state format: "{type}:{portal}" e.g. "login:admin"
		parts := strings.SplitN(state, ":", 2)
		portal := "portal"
		if len(parts) == 2 {
			portal = parts[1]
		}

		userInfo, err := h.usersUsecase.HandleGoogleCallback(code, portal, c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to handle callback"})
			return
		}

		c.JSON(http.StatusOK, userInfo)
	}
}

func (h *UsersHandlers) Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear()
		session.Save()
		c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
	}
}

func (h *UsersHandlers) TestSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		username := session.Get("userID")
		if username == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Welcome, " + username.(string)})
	}
}
