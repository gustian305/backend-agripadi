package models

import (
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`

	ConversationID uuid.UUID    `gorm:"type:uuid;index" json:"conversation_id"`
	Conversation   Conversation `gorm:"foreignKey:ConversationID;OnDelete:CASCADE" json:"conversation"`

	SenderType string `gorm:"type:varchar(20)" json:"sender_type"` // user / ai
	ChatType string `gorm:"type:varchar(20)" json:"chat_type"` // text / image
	Content    string `gorm:"type:text" json:"content"`

	AIResult     string  `gorm:"type:varchar(255)" json:"ai_result"`
	AIConfidence float64 `gorm:"type:decimal(5,4)" json:"ai_confidence"`

	Status string `gorm:"type:varchar(20)" json:"status"` // pending / success / failed

	CreatedAt time.Time `json:"created_at"`


	Attachment []ChatAttachment `gorm:"foreignKey:ChatID" json:"attachment"`
}

type ChatAttachment struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`

	ChatID uuid.UUID `gorm:"type:uuid;index" json:"chat_id"`
	Chat   Chat   `gorm:"foreignKey:ChatID" json:"chat"`

	FileURL  string `gorm:"type:varchar(255)" json:"file_url"`
	FileName string `gorm:"type:varchar(255)" json:"file_name"`
	FileSize int64  `gorm:"type:bigint" json:"file_size"`
	MimeType string `gorm:"type:varchar(255)" json:"mime_type"`

	CreatedAt time.Time `json:"created_at"`
}

type AIResponse struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`

	ChatID uuid.UUID `gorm:"type:uuid;index" json:"chat_id"`
	Chat   Chat   `gorm:"foreignKey:ChatID" json:"chat"`

	ModelName string `gorm:"type:varchar(20)" json:"model_name"`
	Token     int64  `gorm:"type:bigint" json:"token"`
	LatencyMs int64  `gorm:"type:bigint" json:"latency_ms"`

	CreatedAt time.Time `json:"created_at"`
}
