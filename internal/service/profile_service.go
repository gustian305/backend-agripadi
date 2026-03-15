package service

import (
	"backend/internal/dto"
	"backend/internal/models"
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserProfileInterface interface {
	UpdateProfile(userID uuid.UUID, req dto.UpdateProfileRequest) (*dto.UpdateProfileResponse, error)
	GetProfile(userID uuid.UUID) (*dto.UpdateProfileResponse, error)
}

type UserProfileService struct {
	DB *gorm.DB
}

func NewUserProfileService(db *gorm.DB) *UserProfileService {
	return &UserProfileService{DB: db}
}

func (s* UserProfileService) UpdateProfile(userID uuid.UUID, req dto.UpdateProfileRequest) (*dto.UpdateProfileResponse, error) {
	var user models.User

	err := s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&user, "id = ?", userID).Error; err != nil {
			return errors.New("user tidak ditemukan")
		}

		if req.NamaLengkap != "" {
			user.NamaLengkap = req.NamaLengkap
		}

		if req.NoTelepon != "" && req.NoTelepon != user.NoTelepon {
			
			var existing models.User
			
			err := tx.Where("no_telepon = ? AND id != ?", req.NoTelepon, userID).First(&existing).Error

			if err == nil {
				return errors.New("nomor telepon sudah terdaftar")
			}

			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}

			user.NoTelepon = req.NoTelepon
		}

		if req.Password != "" {
			if req.Password != req.ConfirmPassword {
				return errors.New("password dan confirm password tidak cocok")
			}

			hashedPassword, err := bcrypt.GenerateFromPassword(
				[]byte(req.Password),
				bcrypt.DefaultCost,
			)

			if err != nil {
				return err
			}

			user.Password = string(hashedPassword)
		}

		if err := tx.Save(&user).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &dto.UpdateProfileResponse{
		ID:          user.ID.String(),
		NamaLengkap: user.NamaLengkap,
		NoTelepon:   user.NoTelepon,
	}, nil
}

func (s *UserProfileService) GetProfile(userID uuid.UUID) (*dto.UpdateProfileResponse, error) {
	var user models.User

	err := s.DB.First(&user, "id = ?", userID).Error
	if err != nil {
		return nil, errors.New("user tidak ditemukan")
	}

	return &dto.UpdateProfileResponse{
		ID:          user.ID.String(),
		NamaLengkap: user.NamaLengkap,
		NoTelepon:   user.NoTelepon,
	}, nil
}