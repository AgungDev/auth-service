package usecases

import (
	"auth_service/internal/domain"
	"auth_service/internal/domain/dto"
	"auth_service/internal/repositories"
	"errors"
)

type permissionUsecase struct {
	permissionRepo repositories.PermissionRepositoryInterface
}

// NewPermissionUsecase creates a new permission usecase instance
func NewPermissionUsecase(permissionRepo repositories.PermissionRepositoryInterface) PermissionUsecaseInterface {
	return &permissionUsecase{permissionRepo: permissionRepo}
}

// Create creates a new permission
func (uc *permissionUsecase) Create(req dto.PermissionRequest) (*domain.Permission, error) {
	// Check if permission already exists
	existingPermission, err := uc.permissionRepo.FindByName(req.Name)
	if err != nil {
		return nil, err
	}
	if existingPermission != nil {
		return nil, errors.New("permission already exists")
	}

	// Create permission
	permission := &domain.Permission{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := uc.permissionRepo.Create(permission); err != nil {
		return nil, err
	}

	return permission, nil
}

// GetByID retrieves a permission by ID
func (uc *permissionUsecase) GetByID(id uint) (*domain.Permission, error) {
	return uc.permissionRepo.FindByID(id)
}

// GetByName retrieves a permission by name
func (uc *permissionUsecase) GetByName(name string) (*domain.Permission, error) {
	return uc.permissionRepo.FindByName(name)
}

// GetAll retrieves all permissions
func (uc *permissionUsecase) GetAll() ([]domain.Permission, error) {
	return uc.permissionRepo.FindAll()
}

// Update updates permission information
func (uc *permissionUsecase) Update(id uint, req dto.PermissionRequest) error {
	permission := &domain.Permission{
		Name:        req.Name,
		Description: req.Description,
	}
	return uc.permissionRepo.Update(id, permission)
}

// Delete removes a permission
func (uc *permissionUsecase) Delete(id uint) error {
	return uc.permissionRepo.Delete(id)
}
