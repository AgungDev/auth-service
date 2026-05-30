package usecases

import (
	"auth_service/internal/domain"
	"auth_service/internal/domain/dto"
	"auth_service/internal/repositories"
	"auth_service/pkg"
	"errors"
)

type clientUsecase struct {
	clientRepo repositories.ClientRepositoryInterface
}

// NewClientUsecase creates a new client usecase instance
func NewClientUsecase(clientRepo repositories.ClientRepositoryInterface) ClientUsecaseInterface {
	return &clientUsecase{clientRepo: clientRepo}
}

// Create creates a new client with hashed secret
func (uc *clientUsecase) Create(req dto.ClientRequest) (*domain.Client, error) {
	// Check if client already exists
	existingClient, err := uc.clientRepo.FindByClientID(req.ClientID)
	if err != nil {
		return nil, err
	}
	if existingClient != nil {
		return nil, errors.New("client already exists")
	}

	// Hash client secret
	hashedSecret, err := pkg.HashPassword(req.ClientSecret)
	if err != nil {
		return nil, errors.New("failed to hash client secret")
	}

	// Create client
	client := &domain.Client{
		ClientID:         req.ClientID,
		ClientSecretHash: hashedSecret,
		RedirectURIs:     req.RedirectURIs,
		Grants:           req.Grants,
		IsConfidential:   req.IsConfidential,
	}

	if err := uc.clientRepo.Create(client); err != nil {
		return nil, err
	}

	return client, nil
}

// GetByID retrieves a client by ID
func (uc *clientUsecase) GetByID(id uint) (*domain.Client, error) {
	return uc.clientRepo.FindByID(id)
}

// GetByClientID retrieves a client by client_id
func (uc *clientUsecase) GetByClientID(clientID string) (*domain.Client, error) {
	return uc.clientRepo.FindByClientID(clientID)
}

// GetAll retrieves all clients
func (uc *clientUsecase) GetAll() ([]domain.Client, error) {
	return uc.clientRepo.FindAll()
}

// Update updates client information
func (uc *clientUsecase) Update(id uint, req dto.ClientRequest) error {
	hashedSecret, err := pkg.HashPassword(req.ClientSecret)
	if err != nil {
		return errors.New("failed to hash client secret")
	}

	client := &domain.Client{
		ClientID:         req.ClientID,
		ClientSecretHash: hashedSecret,
		RedirectURIs:     req.RedirectURIs,
		Grants:           req.Grants,
		IsConfidential:   req.IsConfidential,
	}
	return uc.clientRepo.Update(id, client)
}

// Delete removes a client
func (uc *clientUsecase) Delete(id uint) error {
	return uc.clientRepo.Delete(id)
}
