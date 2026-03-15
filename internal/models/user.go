package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`

	NamaLengkap string `gorm:"type:varchar(150);not null" json:"nama_lengkap"`
	NoTelepon   string `gorm:"type:varchar(20)" json:"no_telepon"`
	Password    string `gorm:"type:varchar(255);not null" json:"password"`
	Validasi    bool   `gorm:"type:boolean;default:false" json:"validasi"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Conversation []Conversation `gorm:"foreignKey:UserID" json:"conversation"`
}
