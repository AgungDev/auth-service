package repositories

import (
	"auth_service/internal/domain"
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type tokenRepository struct {
	db *gorm.DB
}

// NewTokenRepository creates a new token repository instance
func NewTokenRepository(db *gorm.DB) TokenRepositoryInterface {
	return &tokenRepository{db: db}
}

// Create inserts a new token into the database
func (r *tokenRepository) Create(ctx context.Context, token *domain.Token) error {
	return r.db.WithContext(ctx).Create(token).Error
}

// FindByID finds a token by ID
func (r *tokenRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Token, error) {
	var token domain.Token
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&token).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &token, nil
}

// FindByRefreshTokenHash finds a token by refresh token hash
func (r *tokenRepository) FindByRefreshTokenHash(ctx context.Context, refreshTokenHash string) (*domain.Token, error) {
	var token domain.Token
	if err := r.db.WithContext(ctx).Where("refresh_token_hash = ?", refreshTokenHash).First(&token).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &token, nil
}

// FindByUserID finds tokens by user ID
func (r *tokenRepository) FindByUserID(ctx context.Context, userID uuid.UUID) (*domain.Token, error) {
	var token domain.Token
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&token).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &token, nil
}

// Update updates an existing token
func (r *tokenRepository) Update(ctx context.Context, id uuid.UUID, token *domain.Token) error {
	return r.db.WithContext(ctx).Model(&domain.Token{}).Where("id = ?", id).Updates(token).Error
}

// Delete deletes a token by ID
func (r *tokenRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.Token{}, id).Error
}

// DeleteByUserID deletes all tokens for the specified user
func (r *tokenRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&domain.Token{}).Error
}
