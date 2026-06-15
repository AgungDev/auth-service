package repositories

import (
	"auth_service/internal/domain"
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type clientRepository struct {
	db *gorm.DB
}

// NewClientRepository creates a new client repository instance
func NewClientRepository(db *gorm.DB) ClientRepositoryInterface {
	return &clientRepository{db: db}
}

// Create inserts a new client into the database
func (r *clientRepository) Create(ctx context.Context, client *domain.Client) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(client).Error; err != nil {
			return err
		}
		return nil
	})
}

// FindByID finds a client by ID
func (r *clientRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Client, error) {
	var client domain.Client
	if err := r.db.WithContext(ctx).Preload("Grants").Preload("RedirectURIs").Where("id = ?", id).First(&client).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &client, nil
}

// FindByClientID finds a client by client_id
func (r *clientRepository) FindByClientID(ctx context.Context, clientID string) (*domain.Client, error) {
	var client domain.Client
	if err := r.db.WithContext(ctx).Preload("Grants").Preload("RedirectURIs").Where("client_id = ?", clientID).First(&client).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &client, nil
}

// FindAll retrieves all clients
func (r *clientRepository) FindAll(ctx context.Context) ([]domain.Client, error) {
	var clients []domain.Client
	if err := r.db.WithContext(ctx).Preload("Grants").Preload("RedirectURIs").Find(&clients).Error; err != nil {
		return nil, err
	}
	return clients, nil
}

// Update updates an existing client
func (r *clientRepository) Update(ctx context.Context, id uuid.UUID, client *domain.Client) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&domain.Client{}).Where("id = ?", id).Updates(client).Error; err != nil {
			return err
		}
		return nil
	})
}

// Delete deletes a client by ID
func (r *clientRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&domain.Client{}, id).Error; err != nil {
			return err
		}
		return nil
	})
}
