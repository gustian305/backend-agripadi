package routes

import (
	"backend/config"
	"backend/internal/handler"
	"backend/internal/middleware"
	"backend/internal/service"

	"github.com/gin-gonic/gin"
)

func ProfileRoutes(route *gin.RouterGroup) {
	profileService := service.NewUserProfileService(config.DB)
	profileHandler := handler.NewProfileHandler(profileService)

	authMiddleware := middleware.AuthMiddleware()

	profile := route.Group("/profile")
	profile.Use(authMiddleware)
	{
		profile.PUT("", profileHandler.UpdateProfileHandler)
		profile.GET("", profileHandler.GetProfileHandler)
	}
}