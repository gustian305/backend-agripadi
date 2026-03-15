package service

import (
	"backend/internal/dto"
	"backend/internal/models"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChatInterface interface {
	SendUserChat(userID uuid.UUID, req dto.SendChatRequest) (*dto.ChatResponse, error)
}

type ChatService struct {
	DB *gorm.DB

	conversationService ConversationInterface
	expertSystemService *ExpertSystemService
	llmService          *LLMService
	cnnService          *CNNService
	intentService       *IntentService
}

func NewChatService(db *gorm.DB, conversationService ConversationInterface, expertSystemService *ExpertSystemService, llmService *LLMService, cnnService *CNNService, intentService *IntentService) ChatInterface {
	return &ChatService{
		DB:                  db,
		conversationService: conversationService,
		expertSystemService: expertSystemService,
		llmService:          llmService,
		cnnService:          cnnService,
		intentService:       intentService,
	}
}

func (s *ChatService) SendUserChat(userID uuid.UUID, req dto.SendChatRequest) (*dto.ChatResponse, error) {
	conversationID, err := uuid.Parse(req.ConversationID)
	if err != nil {
		return nil, err
	}

	err = s.conversationService.ValidateConversationOwnership(userID, conversationID)
	if err != nil {
		return nil, err
	}

	if req.ChatType != "text" && req.ChatType != "image" {
		return nil, fmt.Errorf("invalid chat type")
	}

	var response *dto.ChatResponse

	err = s.DB.Transaction(func(tx *gorm.DB) error {

		userChat := models.Chat{
			ID:             uuid.New(),
			ConversationID: conversationID,
			SenderType:     "user",
			ChatType:       req.ChatType,
			Content:        req.Content,
			Status:         "success",
		}

		if err := tx.Create(&userChat).Error; err != nil {
			return err
		}

		if req.ChatType == "image" {
			resp, err := s.handleImageChat(tx, conversationID, &userChat)
			if err != nil {
				return err
			}
			response = resp
			return nil
		}

		resp, err := s.handleTextChat(tx, conversationID, &userChat)
		if err != nil {
			return err
		}
		response = resp

		return nil
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *ChatService) handleImageChat(tx *gorm.DB, conversationID uuid.UUID, imageChat *models.Chat) (*dto.ChatResponse, error) {

	result, err := s.cnnService.DetectHama(filepath.Join(".", imageChat.Content))
	if err != nil {
		return nil, err
	}

	imageChat.AIResult = result.Prediction
	imageChat.AIConfidence = result.Confidence

	if err := tx.Save(imageChat).Error; err != nil {
		return nil, err
	}

	err = tx.Model(&models.Conversation{}).
		Where("id = ?", conversationID).
		Updates(map[string]interface{}{
			"ai_state":  dto.AIStateAskLand,
			"last_hama": result.Prediction,
		}).Error

	if err != nil {
		return nil, err
	}

	text := fmt.Sprintf(
		"Saya mendeteksi hama %s dengan tingkat kepercayaan %.2f%%.\n\nUntuk memberikan rekomendasi pestisida yang tepat, berapa luas lahan padi Anda (dalam hektar)?",
		result.Prediction,
		result.Confidence*100,
	)

	return s.saveAIResponse(tx, conversationID, text, result)
}

func (s *ChatService) handleTextChat(tx *gorm.DB, conversationID uuid.UUID, userMessage *models.Chat) (*dto.ChatResponse, error) {

	var conversation models.Conversation

	if err := tx.First(&conversation, "id = ?", conversationID).Error; err != nil {
		return nil, err
	}

	switch conversation.AIState {

	case dto.AIStateAskLand:
		return s.processLandArea(tx, &conversation, userMessage)

	case dto.AIStateAskPhase:
		return s.processRicePhase(tx, &conversation, userMessage)

	default:
		return s.handleEducationChat(tx, conversationID, userMessage)
	}
}

func (s *ChatService) processLandArea(tx *gorm.DB, conversation *models.Conversation, userMessage *models.Chat) (*dto.ChatResponse, error) {

	land, err := strconv.ParseFloat(userMessage.Content, 64)
	if err != nil {

		text := "Mohon masukkan luas lahan dalam angka. Contoh: 1.5"

		return s.saveAIResponse(tx, conversation.ID, text, nil)
	}

	if land <= 0 || land > 100 {

		text := "Luas lahan tidak valid. Masukkan angka antara 0.1 sampai 100 hektar."

		return s.saveAIResponse(tx, conversation.ID, text, nil)
	}

	conversation.LandArea = land
	conversation.AIState = dto.AIStateAskPhase

	if err := tx.Save(conversation).Error; err != nil {
		return nil, err
	}

	text := "Apa fase pertumbuhan padi saat ini? (vegetatif / generatif)"

	return s.saveAIResponse(tx, conversation.ID, text, nil)
}

func (s *ChatService) processRicePhase(tx *gorm.DB, conversation *models.Conversation, userMessage *models.Chat) (*dto.ChatResponse, error) {

	phase := strings.ToLower(userMessage.Content)

	if phase != "vegetatif" && phase != "generatif" {

		text := "Fase tidak valid. Silakan pilih: vegetatif atau generatif."

		return s.saveAIResponse(tx, conversation.ID, text, nil)
	}

	conversation.RicePhase = phase
	conversation.AIState = dto.AIStateIdle

	if err := tx.Save(conversation).Error; err != nil {
		return nil, err
	}

	pestisida, err := s.expertSystemService.GetPestisidaByHama(conversation.LastHama)
	if err != nil {
		return nil, err
	}

	recommendation, err := s.expertSystemService.CalculateDosage(
		pestisida,
		conversation.LandArea,
	)
	if err != nil {
		return nil, err
	}

	var names []string

	for _, p := range recommendation {
		names = append(names, p.NamaPestisida)
	}

	text, err := s.llmService.GeneratePesticideExplanation(
		conversation.LastHama,
		conversation.LandArea,
		phase,
		names,
	)

	if err != nil {
		return nil, err
	}

	if err := tx.Save(conversation).Error; err != nil {
		return nil, err
	}

	return s.saveAIResponse(tx, conversation.ID, text, nil)
}

func (s *ChatService) handleEducationChat(tx *gorm.DB, conversationID uuid.UUID, userMessage *models.Chat) (*dto.ChatResponse, error) {

	intent := s.intentService.DetectIntent(userMessage.Content)

	switch intent {

	case IntentGreeting:
		return s.saveAIResponse(tx, conversationID, greetingMessage(), nil)

	case IntentHelp:
		return s.saveAIResponse(tx, conversationID, helpMessage(), nil)

	case IntentThanks:
		return s.saveAIResponse(tx, conversationID, thanksMessage(), nil)

	case IntentRice:
		text, err := s.llmService.GenerateRiceKnowledge(userMessage.Content)
		if err != nil {
			return nil, err
		}
		return s.saveAIResponse(tx, conversationID, text, nil)

	default:
		text := "Maaf, saya hanya dapat membantu pertanyaan terkait budidaya padi."
		return s.saveAIResponse(tx, conversationID, text, nil)
	}

}

func (s *ChatService) saveAIResponse(tx *gorm.DB, conversationID uuid.UUID, text string, detection *dto.CNNPredictionResponse) (*dto.ChatResponse, error) {

	msg := models.Chat{
		ID:             uuid.New(),
		ConversationID: conversationID,
		SenderType:     "ai",
		ChatType:       "text",
		Content:        text,
		Status:         "success",
		CreatedAt:      time.Now(),
	}

	if detection != nil {
		msg.AIResult = detection.Prediction
		msg.AIConfidence = detection.Confidence
	}

	if err := tx.Create(&msg).Error; err != nil {
		return nil, err
	}

	return &dto.ChatResponse{
		ID:             msg.ID.String(),
		ConversationID: msg.ConversationID.String(),
		SenderType:     msg.SenderType,
		ChatType:       msg.ChatType,
		Content:        msg.Content,
		Status:         msg.Status,
		AIResult:       msg.AIResult,
		AIConfidence:   msg.AIConfidence,
		CreatedAt:      msg.CreatedAt.Format(time.RFC3339),
	}, nil
}
