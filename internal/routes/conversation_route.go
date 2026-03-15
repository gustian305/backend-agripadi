package routes

import (
	"backend/config"
	"backend/internal/handler"
	"backend/internal/middleware"
	"backend/internal/service"

	"github.com/gin-gonic/gin"
)

func ConversationRoutes(router *gin.RouterGroup) {

	conversationService := service.NewConversationService(config.DB)
	conversationHandler := handler.NewConversationHandler(conversationService)

	conversation := router.Group("/conversations")

	conversation.POST("", middleware.AuthMiddleware(), conversationHandler.CreateConversationHandler)
	conversation.GET("", middleware.AuthMiddleware(), conversationHandler.GetUserConversationHandler)
	conversation.GET("/:conversation_id", middleware.AuthMiddleware(), conversationHandler.GetConversationDetailHandler)
	conversation.DELETE("", middleware.AuthMiddleware(), conversationHandler.DeleteConversationHandler)
}
