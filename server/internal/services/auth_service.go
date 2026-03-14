package services

import (
	"oauth/internal/models"
	"oauth/internal/repositories"
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
) (*models.User, error){
	user, err := s.UserRepo.FindByEmail(email)

	if err != nil {
		newUser := models.User{
			Email: email,
			Name: name,
			Avatar: avtar,
			Provider: provider,
			ProviderID: providerID,
		}

		err := s.UserRepo.Create(&newUser)

		if err != nil {
			return nil,err
		}
		return &newUser, nil
	}

	if user.ProviderID == ""{
		user.Provider = provider
		user.ProviderID = providerID
	}
	return user, nil
}