package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"oauth/internal/database"
	"oauth/internal/models"
	"oauth/internal/repositories"
	"oauth/internal/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthService struct {
	UserRepo repositories.UserRepository
}

func (s *AuthService) HandleOAuthLogin(
	email string,
	name string,
	avtar string,
	provider models.Provider,
	providerID string,
) (*models.User, error) {
	user, err := s.UserRepo.FindByEmail(email)

	if err != nil {
		newUser := models.User{
			Email:      email,
			Name:       name,
			Avatar:     avtar,
			Provider:   provider,
			ProviderID: providerID,
		}

		err := s.UserRepo.Create(&newUser)

		if err != nil {
			return nil, err
		}
		return &newUser, nil
	}

	if user.ProviderID == "" {
		user.Provider = provider
		user.ProviderID = providerID
	}
	return user, nil
}

func (s *AuthService) CreateSession(userID primitive.ObjectID) (*models.Session, error) {
	refreshToken, err := utils.GenerateRefreshToken(userID.Hex())
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	key := refreshTokenKey(refreshToken)

	// Store only a hashed token key in Redis to avoid exposing raw refresh tokens at rest.
	err = database.RedisClient.Set(ctx, key, userID.Hex(), utils.RefreshTokenTTL).Err()
	if err != nil {
		return nil, err
	}

	session := models.Session{
		UserID:       userID,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(utils.RefreshTokenTTL),
		CreatedAt:    time.Now(),
	}

	return &session, nil
}

func (s *AuthService) RefreshSession(oldRefreshToken string) (string, string, error) {
	ctx := context.Background()

	oldKey := refreshTokenKey(oldRefreshToken)

	// GETDEL atomically fetches and invalidates the current token to block replay.
	userIDHex, err := database.RedisClient.GetDel(ctx, oldKey).Result()
	if err != nil {
		if err == database.RedisNil {
			return "", "", fmt.Errorf("invalid session")
		}
		return "", "", err
	}

	accessToken, err := utils.GenerateAccessToken(userIDHex)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := utils.GenerateRefreshToken(userIDHex)
	if err != nil {
		return "", "", err
	}

	newKey := refreshTokenKey(newRefreshToken)
	err = database.RedisClient.Set(ctx, newKey, userIDHex, utils.RefreshTokenTTL).Err()
	if err != nil {
		return "", "", err
	}

	return accessToken, newRefreshToken, nil
}

func refreshTokenKey(refreshToken string) string {
	hash := sha256.Sum256([]byte(refreshToken))
	return "session:refresh:" + hex.EncodeToString(hash[:])
}
