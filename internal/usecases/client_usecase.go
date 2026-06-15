package usecases

import (
	"context"
	"errors"
	"fmt"

	"auth_service/internal/domain"
	"auth_service/internal/domain/dto"
	"auth_service/internal/repositories"
	"auth_service/internal/security"
	loggerpkg "auth_service/pkg/logger"
	"github.com/google/uuid"
)

type clientUsecase struct {
	clientRepo repositories.ClientRepositoryInterface
	security   security.SecurityService
	logger     loggerpkg.Logger
}

// NewClientUsecase creates a new client usecase instance
func NewClientUsecase(clientRepo repositories.ClientRepositoryInterface, securityService security.SecurityService, logger loggerpkg.Logger) ClientUsecaseInterface {
	return &clientUsecase{clientRepo: clientRepo, security: securityService, logger: logger}
}

// Create creates a new client with hashed secret
func (uc *clientUsecase) Create(ctx context.Context, req dto.ClientRequest) (*domain.Client, error) {
	uc.logger.Info("", "", "", fmt.Sprintf("create client request received: client_id=%s, name=%s", req.ClientID, req.Name))

	// Check if client already exists
	existingClient, err := uc.clientRepo.FindByClientID(ctx, req.ClientID)
	if err != nil {
		uc.logger.Error("", "", "", fmt.Sprintf("failed to find client by client_id=%s: %v", req.ClientID, err))
		return nil, err
	}
	if existingClient != nil {
		uc.logger.Warn("", "", "", fmt.Sprintf("client already exists: client_id=%s", req.ClientID))
		return nil, errors.New("client already exists")
	}

	// Hash client secret
	hashedSecret, err := uc.security.HashPassword(req.ClientSecret)
	if err != nil {
		uc.logger.Error("", "", "", fmt.Sprintf("failed to hash client secret for client_id=%s: %v", req.ClientID, err))
		return nil, errors.New("failed to hash client secret")
	}

	// build child records for normalized tables
	var redirectURIs []domain.ClientRedirectURI
	for _, r := range req.RedirectURIs {
		redirectURIs = append(redirectURIs, domain.ClientRedirectURI{RedirectURI: r})
	}
	var grants []domain.ClientGrant
	for _, g := range req.Grants {
		grants = append(grants, domain.ClientGrant{GrantType: g})
	}

	client := &domain.Client{
		ClientID:         req.ClientID,
		Name:             req.Name,
		ClientSecretHash: hashedSecret,
		RedirectURIs:     redirectURIs,
		Grants:           grants,
		IsConfidential:   req.IsConfidential,
	}

	if err := uc.clientRepo.Create(ctx, client); err != nil {
		uc.logger.Error("", "", "", fmt.Sprintf("failed to create client client_id=%s: %v", req.ClientID, err))
		return nil, err
	}

	uc.logger.Info("", "", "", fmt.Sprintf("client created: client_id=%s, id=%s", client.ClientID, client.ID.String()))
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
	uc.logger.Info("", "", "", fmt.Sprintf("update client request received: id=%s, client_id=%s", id.String(), req.ClientID))

	hashedSecret, err := uc.security.HashPassword(req.ClientSecret)
	if err != nil {
		uc.logger.Error("", "", "", fmt.Sprintf("failed to hash client secret for client_id=%s: %v", req.ClientID, err))
		return errors.New("failed to hash client secret")
	}

	var redirectURIs []domain.ClientRedirectURI
	for _, r := range req.RedirectURIs {
		redirectURIs = append(redirectURIs, domain.ClientRedirectURI{RedirectURI: r})
	}
	var grants []domain.ClientGrant
	for _, g := range req.Grants {
		grants = append(grants, domain.ClientGrant{GrantType: g})
	}

	client := &domain.Client{
		ClientID:         req.ClientID,
		Name:             req.Name,
		ClientSecretHash: hashedSecret,
		RedirectURIs:     redirectURIs,
		Grants:           grants,
		IsConfidential:   req.IsConfidential,
	}

	if err := uc.clientRepo.Update(ctx, id, client); err != nil {
		uc.logger.Error("", "", "", fmt.Sprintf("failed to update client id=%s: %v", id.String(), err))
		return err
	}

	uc.logger.Info("", "", "", fmt.Sprintf("client updated: id=%s, client_id=%s", id.String(), req.ClientID))
	return nil
}

// Delete removes a client
func (uc *clientUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	uc.logger.Info("", "", "", fmt.Sprintf("delete client request received: id=%s", id.String()))
	if err := uc.clientRepo.Delete(ctx, id); err != nil {
		uc.logger.Error("", "", "", fmt.Sprintf("failed to delete client id=%s: %v", id.String(), err))
		return err
	}
	uc.logger.Info("", "", "", fmt.Sprintf("client deleted: id=%s", id.String()))
	return nil
}
