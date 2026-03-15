package models

import (
	"time"

	"github.com/google/uuid"
)

type Conversation struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`

	UserID uuid.UUID `gorm:"type:uuid;index" json:"user_id"`
	User   User      `gorm:"foreignKey:UserID" json:"user"`

	IsBotChat bool `gorm:"type:boolean" json:"is_bot_chat"`

	AIState string `gorm:"type:varchar(50)" json:"ai_state"`

	LandArea float64 `gorm:"type:float" json:"land_area"`
	RicePhase string `gorm:"type:varchar(50)" json:"rice_phase"`
	LastHama string `gorm:"type:varchar(50)" json:"last_hama"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Chats []Chat `gorm:"foreignKey:ConversationID" json:"chats"`
}
