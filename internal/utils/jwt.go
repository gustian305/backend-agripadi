package utils

import (
	"backend/config"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(config.Cfg.JWT)

type JWTClaims struct {
	UserID    string `json:"user_id"`
	NoTelepon string `json:"no_telepon"`
	jwt.RegisteredClaims
}

func getJWTSecret() []byte {
	return []byte(config.Cfg.JWT)
}

func GenerateJWT(userID string, noTelepon string) (string, int64, error) {

	expiration := time.Now().Add(24 * time.Hour)

	claims := JWTClaims{
		UserID: userID,
		NoTelepon: noTelepon,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "ai-padi-app",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(getJWTSecret())
	if err != nil {
		return "", 0, err
	}

	return signedToken, expiration.Unix(), nil
}

func ValidateJWT(tokenString string) (*JWTClaims, error) {

	token, err := jwt.ParseWithClaims(
		tokenString,
		&JWTClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return getJWTSecret(), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}