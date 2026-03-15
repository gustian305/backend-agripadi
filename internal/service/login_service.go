package service

import (
	"backend/config"
	"backend/internal/dto"
	"backend/internal/models"
	"backend/internal/utils"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserLoginInterface interface {
	LoginUser(req dto.UserLoginRequest) (*dto.UserLoginResponse, error)
	LogoutUser(token string) error
}

type UserLoginService struct {
	DB *gorm.DB
}

func NewUserLoginService() *UserLoginService {
	return &UserLoginService{DB: config.DB}
}

func (s *UserLoginService) LoginUser(req dto.UserLoginRequest) (*dto.UserLoginResponse, error) {

	var user models.User

	err := s.DB.
		Where("no_telepon = ?", req.NoTelepon).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("nomor telepon tidak ditemukan")
		}
		return nil, err
	}

	// cek password
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.Password),
	)

	if err != nil {
		return nil, errors.New("password salah")
	}

	// generate JWT
	token, exp, err := utils.GenerateJWT(
		user.ID.String(),
		user.NoTelepon,
	)

	if err != nil {
		return nil, err
	}

	resp := dto.UserLoginResponse{
		ID:          user.ID.String(),
		Token:       token,
		ExpiresIn:   exp,
		NamaLengkap: user.NamaLengkap,
		NoTelepon:   user.NoTelepon,
		Validasi:    user.Validasi,
	}

	return &resp, nil
}

func (s *UserLoginService) LogoutUser(token string) error {

	if token == "" {
		return errors.New("token tidak ditemukan")
	}

	return nil
}