package routes

import (
	"backend/internal/handler"
	"backend/internal/middleware"
	"backend/internal/service"

	"github.com/gin-gonic/gin"
)

func LoginRoutes(router *gin.RouterGroup) {

	// user login
	userLoginService := service.NewUserLoginService()
	userLoginHandler := handler.NewUserLoginHandler(userLoginService)

	router.POST("/login", userLoginHandler.LoginUser)
	router.POST("/logout", middleware.AuthMiddleware(), userLoginHandler.LogoutUser)
}
