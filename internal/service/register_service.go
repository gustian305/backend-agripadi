package service

import (
	"backend/config"
	"backend/internal/dto"
	"backend/internal/models"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRegisterInterface interface {
	RegisterUser(req dto.UserRegisterRequest) (*models.User, error)
}

type UserRegisterService struct {
	DB *gorm.DB
}

func NewUserRegisterService() *UserRegisterService {
	return &UserRegisterService{DB: config.DB}
}

func (s *UserRegisterService) RegisterUser(req dto.UserRegisterRequest) (*models.User, error) {

	if req.Password != req.ConfirmPassword {
		return nil, errors.New("password dan confirm password tidak cocok")
	}

	var user models.User

	err := s.DB.Transaction(func(tx *gorm.DB) error {
		var existing models.User 

		err := tx.Where("no_telepon = ?", req.NoTelepon). 
		First(&existing).Error

		if err == nil {
			return errors.New("nomor telepon sudah terdaftar")
		}

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		hashedPassword, err := bcrypt.GenerateFromPassword(
			[]byte(req.Password),
			bcrypt.DefaultCost,
		)

		if err != nil {
			return err
		}

		user = models.User{
			NamaLengkap: req.NamaLengkap,
			NoTelepon:   req.NoTelepon,
			Password: string(hashedPassword),
		}

		err = tx.Create(&user).Error
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &user, nil
}