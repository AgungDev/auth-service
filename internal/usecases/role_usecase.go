package usecases

import (
	"auth_service/internal/domain"
	"auth_service/internal/domain/dto"
	"auth_service/internal/repositories"
	"errors"
)

type roleUsecase struct {
	roleRepo repositories.RoleRepositoryInterface
}

// NewRoleUsecase creates a new role usecase instance
func NewRoleUsecase(roleRepo repositories.RoleRepositoryInterface) RoleUsecaseInterface {
	return &roleUsecase{roleRepo: roleRepo}
}

// Create creates a new role
func (uc *roleUsecase) Create(req dto.RoleRequest) (*domain.Role, error) {
	// Check if role already exists
	existingRole, err := uc.roleRepo.FindByName(req.Name)
	if err != nil {
		return nil, err
	}
	if existingRole != nil {
		return nil, errors.New("role already exists")
	}

	// Create role
	role := &domain.Role{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := uc.roleRepo.Create(role); err != nil {
		return nil, err
	}

	return role, nil
}

// GetByID retrieves a role by ID
func (uc *roleUsecase) GetByID(id uint) (*domain.Role, error) {
	return uc.roleRepo.FindByID(id)
}

// GetByName retrieves a role by name
func (uc *roleUsecase) GetByName(name string) (*domain.Role, error) {
	return uc.roleRepo.FindByName(name)
}

// GetAll retrieves all roles
func (uc *roleUsecase) GetAll() ([]domain.Role, error) {
	return uc.roleRepo.FindAll()
}

// Update updates role information
func (uc *roleUsecase) Update(id uint, req dto.RoleRequest) error {
	role := &domain.Role{
		Name:        req.Name,
		Description: req.Description,
	}
	return uc.roleRepo.Update(id, role)
}

// Delete removes a role
func (uc *roleUsecase) Delete(id uint) error {
	return uc.roleRepo.Delete(id)
}
