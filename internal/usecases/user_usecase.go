package usecases

import (
	"auth_service/internal/domain"
	"auth_service/internal/domain/dto"
	"auth_service/internal/repositories"
	"auth_service/pkg"
	"errors"
)

type userUsecase struct {
	userRepo repositories.UserRepositoryInterface
}

// NewUserUsecase creates a new user usecase instance
func NewUserUsecase(userRepo repositories.UserRepositoryInterface) UserUsecaseInterface {
	return &userUsecase{userRepo: userRepo}
}

// Register handles user registration with validation and password hashing
func (uc *userUsecase) Register(req dto.UserRegisterRequest) (*domain.User, error) {
	// Check if username already exists
	existingUser, err := uc.userRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	// Check if email already exists
	existingEmail, err := uc.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if existingEmail != nil {
		return nil, errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := pkg.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Create user
	user := &domain.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		FullName:     req.FullName,
		IsActive:     true,
	}

	if err := uc.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login validates user credentials
func (uc *userUsecase) Login(req dto.UserLoginRequest) (*domain.User, error) {
	// Find user by username
	user, err := uc.userRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, err
	}

	// Check if user exists
	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("user is not active")
	}

	// Verify password
	if !pkg.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

// GetByID retrieves a user by ID
func (uc *userUsecase) GetByID(id uint) (*domain.User, error) {
	return uc.userRepo.FindByID(id)
}

// GetByUsername retrieves a user by username
func (uc *userUsecase) GetByUsername(username string) (*domain.User, error) {
	return uc.userRepo.FindByUsername(username)
}

// GetAll retrieves all users
func (uc *userUsecase) GetAll() ([]domain.User, error) {
	return uc.userRepo.FindAll()
}

// Update updates user information
func (uc *userUsecase) Update(id uint, req dto.UserUpdateRequest) error {
	user := &domain.User{
		FullName: req.FullName,
		Email:    req.Email,
	}
	return uc.userRepo.Update(id, user)
}

// Delete removes a user
func (uc *userUsecase) Delete(id uint) error {
	return uc.userRepo.Delete(id)
}
