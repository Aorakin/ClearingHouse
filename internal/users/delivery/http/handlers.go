package http

import (
	"net/http"

	"github.com/ClearingHouse/internal/users/dtos"
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
		url := h.usersUsecase.GenerateLoginURL("state-token")
		c.Redirect(http.StatusTemporaryRedirect, url)
	}
}

func (h *UsersHandlers) Callback() gin.HandlerFunc {
	return func(c *gin.Context) {
		code := c.Query("code")
		if code == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No code provided"})
			return
		}

		userInfo, err := h.usersUsecase.HandleGoogleCallback(code, c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to handle callback"})
			return
		}

		c.JSON(http.StatusOK, userInfo)
	}
}

func (h *UsersHandlers) Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		var registerInput dtos.RegisterInput
		if err := c.ShouldBindJSON(&registerInput); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		if err := h.usersUsecase.Register(registerInput); err != nil {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "registered successfully"})
	}
}

func (h *UsersHandlers) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		user, err := h.usersUsecase.Login(username, password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
			return
		}

		session := sessions.Default(c)
		session.Set("userID", user.ID)
		session.Save()

		c.JSON(http.StatusOK, gin.H{"message": "Login successful"})

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
