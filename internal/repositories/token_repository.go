package repositories

import (
	"auth_service/internal/domain"
	"errors"

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
func (r *tokenRepository) Create(token *domain.Token) error {
	return r.db.Create(token).Error
}

// FindByID finds a token by ID
func (r *tokenRepository) FindByID(id uint) (*domain.Token, error) {
	var token domain.Token
	if err := r.db.Where("id = ?", id).First(&token).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &token, nil
}

// FindByAccessToken finds a token by access token
func (r *tokenRepository) FindByAccessToken(accessToken string) (*domain.Token, error) {
	var token domain.Token
	if err := r.db.Where("access_token = ?", accessToken).First(&token).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &token, nil
}

// FindByUserID finds tokens by user ID
func (r *tokenRepository) FindByUserID(userID uint) (*domain.Token, error) {
	var token domain.Token
	if err := r.db.Where("user_id = ?", userID).First(&token).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &token, nil
}

// Update updates an existing token
func (r *tokenRepository) Update(id uint, token *domain.Token) error {
	return r.db.Model(&domain.Token{}).Where("id = ?", id).Updates(token).Error
}

// Delete deletes a token by ID
func (r *tokenRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Token{}, id).Error
}
