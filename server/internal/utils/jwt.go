package utils

import (
	"oauth/internal/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	AccessTokenTTL  = 15 * time.Minute
	RefreshTokenTTL = 7 * 24 * time.Hour
)

func GenerateJWT(userID string) (string, error) {
	return generateToken(userID, AccessTokenTTL)
}

func GenerateAccessToken(userID string) (string, error) {
	return generateToken(userID, AccessTokenTTL)
}

func GenerateRefreshToken(userID string) (string, error) {
	return generateToken(userID, RefreshTokenTTL)
}

func generateToken(userID string, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(ttl).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(config.Config.JWTSecret))
}
