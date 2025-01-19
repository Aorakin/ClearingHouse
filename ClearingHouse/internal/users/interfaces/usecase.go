package interfaces

import "github.com/gin-gonic/gin"

type UsersUsecase interface {
	GenerateLoginURL(string) string
	HandleGoogleCallback(string, *gin.Context) (map[string]interface{}, error)
}
