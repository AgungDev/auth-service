package usecases

import (
	"context"
	"errors"

	"auth_service/internal/domain"
	"auth_service/internal/domain/dto"
	"auth_service/internal/repositories"
	"github.com/google/uuid"
)

type permissionUsecase struct {
	permissionRepo repositories.PermissionRepositoryInterface
}

// NewPermissionUsecase creates a new permission usecase instance
func NewPermissionUsecase(permissionRepo repositories.PermissionRepositoryInterface) PermissionUsecaseInterface {
	return &permissionUsecase{permissionRepo: permissionRepo}
}

// Create creates a new permission
func (uc *permissionUsecase) Create(ctx context.Context, req dto.PermissionRequest) (*domain.Permission, error) {
	// Check if permission already exists
	existingPermission, err := uc.permissionRepo.FindByName(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	if existingPermission != nil {
		return nil, errors.New("permission already exists")
	}

	// Create permission
	permission := &domain.Permission{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
	}

	if err := uc.permissionRepo.Create(ctx, permission); err != nil {
		return nil, err
	}

	return permission, nil
}

// GetByID retrieves a permission by ID
func (uc *permissionUsecase) GetByID(ctx context.Context, id uuid.UUID) (*domain.Permission, error) {
	return uc.permissionRepo.FindByID(ctx, id)
}

// GetByName retrieves a permission by name
func (uc *permissionUsecase) GetByName(ctx context.Context, name string) (*domain.Permission, error) {
	return uc.permissionRepo.FindByName(ctx, name)
}

// GetByUserID retrieves permissions linked to the user's roles
func (uc *permissionUsecase) GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Permission, error) {
	return uc.permissionRepo.FindByUserID(ctx, userID)
}

// GetAll retrieves all permissions
func (uc *permissionUsecase) GetAll(ctx context.Context) ([]domain.Permission, error) {
	return uc.permissionRepo.FindAll(ctx)
}

// Update updates permission information
func (uc *permissionUsecase) Update(ctx context.Context, id uuid.UUID, req dto.PermissionRequest) error {
	permission := &domain.Permission{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
	}
	return uc.permissionRepo.Update(ctx, id, permission)
}

// Delete removes a permission
func (uc *permissionUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	return uc.permissionRepo.Delete(ctx, id)
}
