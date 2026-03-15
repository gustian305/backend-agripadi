package routes

import (
	"backend/internal/handler"
	"backend/internal/service"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup) {

	// user register
	userRegisterService := service.NewUserRegisterService()
	userRegisterHandler := handler.NewUserRegisterHandler(userRegisterService)

	router.POST("/register", userRegisterHandler.RegisterUser)
}
