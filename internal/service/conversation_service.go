package service

import (
	"backend/internal/dto"
	"backend/internal/models"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ConversationInterface interface {
	CreateConversation(userID uuid.UUID, req dto.CreateConversationRequest) (*dto.ConversationResponse, error)
	GetUserConversation(userID uuid.UUID) ([]dto.ConversationListResponse, error)
	GetConversationDetail(userID uuid.UUID, conversationID uuid.UUID) (*dto.ConversationDetailResponse, error)
	ValidateConversationOwnership(userID uuid.UUID, conversationID uuid.UUID) error
	DeleteConversation(userID uuid.UUID, req dto.DeleteConversationRequest) error
}

type ConversationService struct {
	DB *gorm.DB
}

func NewConversationService(db *gorm.DB) ConversationInterface {
	return &ConversationService{DB: db}
}

func (s *ConversationService) CreateConversation(userID uuid.UUID, req dto.CreateConversationRequest) (*dto.ConversationResponse, error) {

	conversation := models.Conversation{
		ID:        uuid.New(),
		UserID:    userID,
		IsBotChat: req.IsBotChat,
		AIState:   dto.AIStateIdle,
	}

	err := s.DB.Create(&conversation).Error
	if err != nil {
		return nil, err
	}

	res := dto.ConversationResponse{
		ID:        conversation.ID.String(),
		UserID:    conversation.UserID.String(),
		IsBotChat: conversation.IsBotChat,
		AIState:   conversation.AIState,
		CreatedAt: conversation.CreatedAt.Format(time.RFC3339),
	}

	return &res, nil
}

func (s *ConversationService) GetUserConversation(userID uuid.UUID) ([]dto.ConversationListResponse, error) {

	var conversations []models.Conversation

	err := s.DB.
		Preload("Chats").
		Where("user_id = ?", userID).
		Order("updated_at DESC").
		Find(&conversations).Error

	if err != nil {
		return nil, err
	}

	var res []dto.ConversationListResponse

	for _, conv := range conversations {

		lastMessage := ""
		count := len(conv.Chats)

		if count > 0 {
			lastMessage = conv.Chats[count-1].Content
		}

		res = append(res, dto.ConversationListResponse{
			ID:          conv.ID.String(),
			LastMessage: lastMessage,
			ChatCount:   count,
			CreatedAt:   conv.CreatedAt.Format(time.RFC3339),
		})
	}

	return res, nil
}

func (s *ConversationService) GetConversationDetail(userID uuid.UUID, conversationID uuid.UUID) (*dto.ConversationDetailResponse, error) {

	err := s.ValidateConversationOwnership(userID, conversationID)
	if err != nil {
		return nil, err
	}

	var conversation models.Conversation

	err = s.DB.
		Preload("Chats.Attachment").
		Where("id = ?", conversationID).
		First(&conversation).Error

	if err != nil {
		return nil, err
	}

	// penting: jangan nil slice
	messages := make([]dto.ChatResponse, 0)

	for _, m := range conversation.Chats {

		msg := dto.ChatResponse{
			ID:             m.ID.String(),
			ConversationID: m.ConversationID.String(),
			SenderType:     m.SenderType,
			ChatType:       m.ChatType,
			Content:        m.Content,
			Status:         m.Status,
			AIResult:       m.AIResult,
			AIConfidence:   m.AIConfidence,
			CreatedAt:      m.CreatedAt.Format(time.RFC3339),
		}

		messages = append(messages, msg)
	}

	res := dto.ConversationDetailResponse{
		ID:    conversation.ID.String(),
		Chats: messages,
	}

	return &res, nil
}

func (s *ConversationService) ValidateConversationOwnership(userID uuid.UUID, conversationID uuid.UUID) error {
	var count int64
	err := s.DB.Model(&models.Conversation{}).
		Where("id = ? AND user_id = ?", conversationID, userID).
		Count(&count).Error
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("conversation not found or access denied")
	}

	return nil
}

func (s *ConversationService) DeleteConversation(userID uuid.UUID, req dto.DeleteConversationRequest) error {
	if len(req.ConversationID) == 0 {
		return errors.New("conversation id is required")
	}

	conversationID, err := uuid.Parse(req.ConversationID[0])
	if err != nil {
		return errors.New("invalid conversation id")
	}

	err = s.ValidateConversationOwnership(userID, conversationID)
	if err != nil {
		return err
	}

	var ids []uuid.UUID

	for _, id := range req.ConversationID {
		conversationId, err := uuid.Parse(id)
		if err != nil {
			return errors.New("invalid conversation id")
		}
		ids = append(ids, conversationId)
	}

	err = s.DB.Where("id IN ? AND user_id = ?", ids, userID).Delete(&models.Conversation{}).Error
	if err != nil {
		return err
	}


	return nil
}