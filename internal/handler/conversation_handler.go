package handler

import (
	"backend/internal/dto"
	"backend/internal/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ConversationHandler struct {
	conversationService service.ConversationInterface
}

func NewConversationHandler(conversationService service.ConversationInterface) *ConversationHandler {
	return &ConversationHandler{conversationService: conversationService}
}

func (h *ConversationHandler) CreateConversationHandler(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		log.Println("user not found")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	userID := userIDRaw.(uuid.UUID)

	var req dto.CreateConversationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("failed to bind json")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.conversationService.CreateConversation(userID, req)
	if err != nil {
		log.Println("failed to create conversation")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    res,
	})
}

func (h *ConversationHandler) GetUserConversationHandler(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		log.Println("user not found")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	userID := userIDRaw.(uuid.UUID)

	res, err := h.conversationService.GetUserConversation(userID)
	if err != nil {
		log.Println("failed to list conversation")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    res,
	})
}

func (h *ConversationHandler) GetConversationDetailHandler(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		log.Println("user not found")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	userID := userIDRaw.(uuid.UUID)

	conversationIDParam := c.Param("conversation_id")

	conversationID, err := uuid.Parse(conversationIDParam)
	if err != nil {
		log.Println("invalid conversation id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid conversation ID"})
		return
	}

	res, err := h.conversationService.GetConversationDetail(userID, conversationID)
	if err != nil {
		log.Println("failed to get conversation detail")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    res,
	})
}

func (h *ConversationHandler) DeleteConversationHandler(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		log.Println("user not found")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	userID := userIDRaw.(uuid.UUID)


	var req dto.DeleteConversationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("failed to bind json")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.conversationService.DeleteConversation(userID, req)
	if err != nil {
		log.Println("failed to delete conversation")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}