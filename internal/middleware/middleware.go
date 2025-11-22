package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ClearingHouse/internal/auth"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 1. Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "missing Authorization header"})
			c.Abort()
			return
		}

		// 2. Must be "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid Authorization header format"})
			c.Abort()
			return
		}

		accessToken := parts[1]

		// 3. Parse JWT using your RSA public key
		token, err := jwt.Parse(accessToken, func(t *jwt.Token) (interface{}, error) {
			// Ensure token is RS512
			if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("invalid signing method")
			}
			return auth.PublicKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid or expired access token"})
			c.Abort()
			return
		}

		// 4. Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid token claims"})
			c.Abort()
			return
		}

		// 5. Extract user info
		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid user_id in token"})
			c.Abort()
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid user_id format"})
			c.Abort()
			return
		}

		firstname, _ := claims["first_name"].(string)
		lastname, _ := claims["last_name"].(string)
		email, _ := claims["email"].(string)

		if userID == uuid.Nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "missing user_id in token"})
			c.Abort()
			return
		}

		// 6. Store into Gin context
		c.Set("userID", userID)
		c.Set("firstname", firstname)
		c.Set("lastname", lastname)
		c.Set("email", email)

		c.Next()
	}
}
