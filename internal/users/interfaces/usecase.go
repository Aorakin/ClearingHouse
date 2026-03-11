package interfaces

import (
	"github.com/gin-gonic/gin"
)

type UsersUsecase interface {
	GenerateLoginURL(state, portal string) string
	HandleGoogleCallback(code, portal string, c *gin.Context) (map[string]interface{}, error)
}
