package config

import (
	"backend/internal/models"

	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {

	err := db.AutoMigrate(
		&models.User{},
		&models.Pestisida{},
		&models.Conversation{},
		&models.Chat{},
		&models.ChatAttachment{},
		&models.AIResponse{},
	)

	return err
}