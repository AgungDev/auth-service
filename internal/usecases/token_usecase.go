package usecases

import (
	"auth_service/internal/domain"
	"auth_service/internal/repositories"
	"errors"
	"time"
)

type tokenUsecase struct {
	tokenRepo repositories.TokenRepositoryInterface
}

// NewTokenUsecase creates a new token usecase instance
func NewTokenUsecase(tokenRepo repositories.TokenRepositoryInterface) TokenUsecaseInterface {
	return &tokenUsecase{tokenRepo: tokenRepo}
}

// Create creates a new token
func (uc *tokenUsecase) Create(token *domain.Token) error {
	return uc.tokenRepo.Create(token)
}

// GetByID retrieves a token by ID
func (uc *tokenUsecase) GetByID(id uint) (*domain.Token, error) {
	return uc.tokenRepo.FindByID(id)
}

// GetByAccessToken retrieves a token by access token
func (uc *tokenUsecase) GetByAccessToken(accessToken string) (*domain.Token, error) {
	token, err := uc.tokenRepo.FindByAccessToken(accessToken)
	if err != nil {
		return nil, err
	}
	if token != nil && uc.IsTokenExpired(token) {
		return nil, errors.New("token has expired")
	}
	return token, nil
}

// GetByUserID retrieves a token by user ID
func (uc *tokenUsecase) GetByUserID(userID uint) (*domain.Token, error) {
	return uc.tokenRepo.FindByUserID(userID)
}

// Update updates token information
func (uc *tokenUsecase) Update(id uint, token *domain.Token) error {
	return uc.tokenRepo.Update(id, token)
}

// Delete removes a token
func (uc *tokenUsecase) Delete(id uint) error {
	return uc.tokenRepo.Delete(id)
}

// IsTokenExpired checks if a token has expired
func (uc *tokenUsecase) IsTokenExpired(token *domain.Token) bool {
	if token.ExpiresAt == 0 {
		return false
	}
	return time.Now().Unix() > token.ExpiresAt
}
