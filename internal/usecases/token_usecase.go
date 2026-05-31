package usecases

import (
	"context"
	"errors"
	"time"

	"auth_service/internal/domain"
	"auth_service/internal/repositories"
	"github.com/google/uuid"
)

type tokenUsecase struct {
	tokenRepo repositories.TokenRepositoryInterface
}

// NewTokenUsecase creates a new token usecase instance
func NewTokenUsecase(tokenRepo repositories.TokenRepositoryInterface) TokenUsecaseInterface {
	return &tokenUsecase{tokenRepo: tokenRepo}
}

// Create creates a new token
func (uc *tokenUsecase) Create(ctx context.Context, token *domain.Token) error {
	return uc.tokenRepo.Create(ctx, token)
}

// GetByID retrieves a token by ID
func (uc *tokenUsecase) GetByID(ctx context.Context, id uuid.UUID) (*domain.Token, error) {
	return uc.tokenRepo.FindByID(ctx, id)
}

// GetByRefreshTokenHash retrieves a token by refresh token hash
func (uc *tokenUsecase) GetByRefreshTokenHash(ctx context.Context, refreshTokenHash string) (*domain.Token, error) {
	token, err := uc.tokenRepo.FindByRefreshTokenHash(ctx, refreshTokenHash)
	if err != nil {
		return nil, err
	}
	if token != nil && uc.IsTokenExpired(token) {
		return nil, errors.New("token has expired")
	}
	return token, nil
}

// GetByUserID retrieves a token by user ID
func (uc *tokenUsecase) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Token, error) {
	return uc.tokenRepo.FindByUserID(ctx, userID)
}

// Update updates token information
func (uc *tokenUsecase) Update(ctx context.Context, id uuid.UUID, token *domain.Token) error {
	return uc.tokenRepo.Update(ctx, id, token)
}

// Delete removes a token
func (uc *tokenUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	return uc.tokenRepo.Delete(ctx, id)
}

// DeleteByUserID removes tokens for a user
func (uc *tokenUsecase) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	return uc.tokenRepo.DeleteByUserID(ctx, userID)
}

// IsTokenExpired checks if a token has expired
func (uc *tokenUsecase) IsTokenExpired(token *domain.Token) bool {
	if token == nil {
		return true
	}
	if token.ExpiresAt.IsZero() {
		return false
	}
	return time.Now().After(token.ExpiresAt)
}
