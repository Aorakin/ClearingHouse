package interfaces

import "github.com/gin-gonic/gin"

type ProjectsHandlers interface {
	GetAll() gin.HandlerFunc
	Get() gin.HandlerFunc
	Add() gin.HandlerFunc
	Edit() gin.HandlerFunc
	Delete() gin.HandlerFunc
}
