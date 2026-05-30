package repositories

import (
	"auth_service/internal/domain"
	"errors"

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
func (r *clientRepository) Create(client *domain.Client) error {
	return r.db.Create(client).Error
}

// FindByID finds a client by ID
func (r *clientRepository) FindByID(id uint) (*domain.Client, error) {
	var client domain.Client
	if err := r.db.Where("id = ?", id).First(&client).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &client, nil
}

// FindByClientID finds a client by client_id
func (r *clientRepository) FindByClientID(clientID string) (*domain.Client, error) {
	var client domain.Client
	if err := r.db.Where("client_id = ?", clientID).First(&client).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &client, nil
}

// FindAll retrieves all clients
func (r *clientRepository) FindAll() ([]domain.Client, error) {
	var clients []domain.Client
	if err := r.db.Find(&clients).Error; err != nil {
		return nil, err
	}
	return clients, nil
}

// Update updates an existing client
func (r *clientRepository) Update(id uint, client *domain.Client) error {
	return r.db.Model(&domain.Client{}).Where("id = ?", id).Updates(client).Error
}

// Delete deletes a client by ID
func (r *clientRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Client{}, id).Error
}
