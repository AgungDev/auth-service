package repositories

import (
	"auth_service/internal/domain"
	"errors"

	"gorm.io/gorm"
)

type roleRepository struct {
	db *gorm.DB
}

// NewRoleRepository creates a new role repository instance
func NewRoleRepository(db *gorm.DB) RoleRepositoryInterface {
	return &roleRepository{db: db}
}

// Create inserts a new role into the database
func (r *roleRepository) Create(role *domain.Role) error {
	return r.db.Create(role).Error
}

// FindByID finds a role by ID
func (r *roleRepository) FindByID(id uint) (*domain.Role, error) {
	var role domain.Role
	if err := r.db.Where("id = ?", id).First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

// FindByName finds a role by name
func (r *roleRepository) FindByName(name string) (*domain.Role, error) {
	var role domain.Role
	if err := r.db.Where("name = ?", name).First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

// FindAll retrieves all roles
func (r *roleRepository) FindAll() ([]domain.Role, error) {
	var roles []domain.Role
	if err := r.db.Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

// Update updates an existing role
func (r *roleRepository) Update(id uint, role *domain.Role) error {
	return r.db.Model(&domain.Role{}).Where("id = ?", id).Updates(role).Error
}

// Delete deletes a role by ID
func (r *roleRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Role{}, id).Error
}
