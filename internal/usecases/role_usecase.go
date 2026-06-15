package usecases

import (
	"context"
	"errors"

	"auth_service/internal/domain"
	"auth_service/internal/domain/dto"
	"auth_service/internal/repositories"
	"github.com/google/uuid"
)

type roleUsecase struct {
	roleRepo repositories.RoleRepositoryInterface
}

// NewRoleUsecase creates a new role usecase instance
func NewRoleUsecase(roleRepo repositories.RoleRepositoryInterface) RoleUsecaseInterface {
	return &roleUsecase{roleRepo: roleRepo}
}

// Create creates a new role
func (uc *roleUsecase) Create(ctx context.Context, req dto.RoleRequest) (*domain.Role, error) {
	// Check if role already exists
	existingRole, err := uc.roleRepo.FindByName(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	if existingRole != nil {
		return nil, errors.New("role already exists")
	}

	// Create role
	role := &domain.Role{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
	}

	if err := uc.roleRepo.Create(ctx, role); err != nil {
		return nil, err
	}

	return role, nil
}

// GetByID retrieves a role by ID
func (uc *roleUsecase) GetByID(ctx context.Context, id uuid.UUID) (*domain.Role, error) {
	return uc.roleRepo.FindByID(ctx, id)
}

// GetByName retrieves a role by name
func (uc *roleUsecase) GetByName(ctx context.Context, name string) (*domain.Role, error) {
	return uc.roleRepo.FindByName(ctx, name)
}

// GetAll retrieves all roles
func (uc *roleUsecase) GetAll(ctx context.Context) ([]domain.Role, error) {
	return uc.roleRepo.FindAll(ctx)
}

// Update updates role information
func (uc *roleUsecase) Update(ctx context.Context, id uuid.UUID, req dto.RoleRequest) error {
	role := &domain.Role{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
	}
	return uc.roleRepo.Update(ctx, id, role)
}

// Delete removes a role
func (uc *roleUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	return uc.roleRepo.Delete(ctx, id)
}

// AssignPermissions assigns permissions to a role
func (uc *roleUsecase) AssignPermissions(ctx context.Context, roleID uuid.UUID, permissionIDs []uuid.UUID) error {
	return uc.roleRepo.AssignPermissions(ctx, roleID, permissionIDs)
}

// GetPermissionsByRoleID retrieves permissions assigned to a role
func (uc *roleUsecase) GetPermissionsByRoleID(ctx context.Context, roleID uuid.UUID) ([]domain.Permission, error) {
	return uc.roleRepo.GetPermissionsByRoleID(ctx, roleID)
}

// RemovePermission removes a permission from a role
func (uc *roleUsecase) RemovePermission(ctx context.Context, roleID uuid.UUID, permissionID uuid.UUID) error {
	return uc.roleRepo.RemovePermission(ctx, roleID, permissionID)
}
