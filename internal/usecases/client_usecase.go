package usecases

import (
	"context"
	"errors"

	"auth_service/internal/domain"
	"auth_service/internal/domain/dto"
	"auth_service/internal/repositories"
	"auth_service/internal/security"
	"github.com/google/uuid"
)

type clientUsecase struct {
	clientRepo  repositories.ClientRepositoryInterface
	security    security.SecurityService
}

// NewClientUsecase creates a new client usecase instance
func NewClientUsecase(clientRepo repositories.ClientRepositoryInterface, securityService security.SecurityService) ClientUsecaseInterface {
	return &clientUsecase{clientRepo: clientRepo, security: securityService}
}

// Create creates a new client with hashed secret
func (uc *clientUsecase) Create(ctx context.Context, req dto.ClientRequest) (*domain.Client, error) {
	// Check if client already exists
	existingClient, err := uc.clientRepo.FindByClientID(ctx, req.ClientID)
	if err != nil {
		return nil, err
	}
	if existingClient != nil {
		return nil, errors.New("client already exists")
	}

	// Hash client secret
	hashedSecret, err := uc.security.HashPassword(req.ClientSecret)
	if err != nil {
		return nil, errors.New("failed to hash client secret")
	}

	// Create client
	client := &domain.Client{
		ClientID:         req.ClientID,
		Name:             req.Name,
		ClientSecretHash: hashedSecret,
		RedirectURIs:     req.RedirectURIs,
		Grants:           req.Grants,
		IsConfidential:   req.IsConfidential,
	}

	if err := uc.clientRepo.Create(ctx, client); err != nil {
		return nil, err
	}

	return client, nil
}

// GetByID retrieves a client by ID
func (uc *clientUsecase) GetByID(ctx context.Context, id uuid.UUID) (*domain.Client, error) {
	return uc.clientRepo.FindByID(ctx, id)
}

// GetByClientID retrieves a client by client_id
func (uc *clientUsecase) GetByClientID(ctx context.Context, clientID string) (*domain.Client, error) {
	return uc.clientRepo.FindByClientID(ctx, clientID)
}

// GetAll retrieves all clients
func (uc *clientUsecase) GetAll(ctx context.Context) ([]domain.Client, error) {
	return uc.clientRepo.FindAll(ctx)
}

// Update updates client information
func (uc *clientUsecase) Update(ctx context.Context, id uuid.UUID, req dto.ClientRequest) error {
	hashedSecret, err := uc.security.HashPassword(req.ClientSecret)
	if err != nil {
		return errors.New("failed to hash client secret")
	}

	client := &domain.Client{
		ClientID:         req.ClientID,
		Name:             req.Name,
		ClientSecretHash: hashedSecret,
		RedirectURIs:     req.RedirectURIs,
		Grants:           req.Grants,
		IsConfidential:   req.IsConfidential,
	}
	return uc.clientRepo.Update(ctx, id, client)
}

// Delete removes a client
func (uc *clientUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	return uc.clientRepo.Delete(ctx, id)
}
