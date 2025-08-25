package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		id := session.Get("userID")
		if id == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			c.Abort()
			return
		}

		userUUID, err := uuid.Parse(id.(string))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid session UUID"})
			c.Abort()
			return
		}

		c.Set("userID", userUUID)
		c.Next()
	}
}
