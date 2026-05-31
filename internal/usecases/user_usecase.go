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

type userUsecase struct {
	userRepo  repositories.UserRepositoryInterface
	security security.SecurityService
}

// NewUserUsecase creates a new user usecase instance
func NewUserUsecase(userRepo repositories.UserRepositoryInterface, securityService security.SecurityService) UserUsecaseInterface {
	return &userUsecase{userRepo: userRepo, security: securityService}
}

// Register handles user registration with validation and password hashing
func (uc *userUsecase) Register(ctx context.Context, req dto.UserRegisterRequest) (*domain.User, error) {
	// Check if username already exists
	existingUser, err := uc.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	// Check if email already exists
	existingEmail, err := uc.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existingEmail != nil {
		return nil, errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := uc.security.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Create user
	user := &domain.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		FullName:     req.FullName,
		Status:       "active",
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login validates user credentials
func (uc *userUsecase) Login(ctx context.Context, req dto.UserLoginRequest) (*domain.User, error) {
	// Find user by username
	user, err := uc.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	// Check if user exists
	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	// Check if user is active
	if user.Status != "active" {
		return nil, errors.New("user is not active")
	}

	// Verify password
	if !uc.security.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

// GetByID retrieves a user by ID
func (uc *userUsecase) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return uc.userRepo.FindByID(ctx, id)
}

// GetByUsername retrieves a user by username
func (uc *userUsecase) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	return uc.userRepo.FindByUsername(ctx, username)
}

// GetAll retrieves all users
func (uc *userUsecase) GetAll(ctx context.Context) ([]domain.User, error) {
	return uc.userRepo.FindAll(ctx)
}

// Update updates user information
func (uc *userUsecase) Update(ctx context.Context, id uuid.UUID, req dto.UserUpdateRequest) error {
	user := &domain.User{
		FullName: req.FullName,
		Email:    req.Email,
	}
	return uc.userRepo.Update(ctx, id, user)
}

// Delete removes a user
func (uc *userUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	return uc.userRepo.Delete(ctx, id)
}
