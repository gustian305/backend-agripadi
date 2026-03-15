package routes

import (
	"backend/config"
	"backend/internal/handler"
	"backend/internal/middleware"
	"backend/internal/service"

	"github.com/gin-gonic/gin"
)

func ChatRoutes(route *gin.RouterGroup) {

	conversationService := service.NewConversationService(config.DB)

	cnnService := service.NewCNNService()

	expertSystemService := service.NewExpertSystemService(config.DB)

	llmService := service.NewLLMService()

	intentService := service.NewIntentService(expertSystemService)

	chatService := service.NewChatService(
		config.DB,
		conversationService,
		expertSystemService,
		llmService,
		cnnService,
		intentService,
	)

	chatHandler := handler.NewChatHandler(chatService)

	chat := route.Group("/chat")

	chat.POST("/send", middleware.AuthMiddleware(), chatHandler.SendUserChatHandler)
	chat.POST("/upload-image", middleware.AuthMiddleware(), chatHandler.UploadUserImageChatHandler)
}