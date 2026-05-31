package usecases

import (
	"context"
	"errors"
	"fmt"
	"time"

	"auth_service/internal/domain"
	"auth_service/internal/domain/dto"
	"auth_service/internal/security"
	"github.com/google/uuid"
)

var (
	ErrUserIDMismatch = errors.New("user id does not match token subject")
	ErrNotAuthorized  = errors.New("not authorized")
)

type authUsecase struct {
	userUsecase       UserUsecaseInterface
	clientUsecase     ClientUsecaseInterface
	permissionUsecase PermissionUsecaseInterface
	tokenUsecase      TokenUsecaseInterface
	security          security.SecurityService
	appName           string
	accessTokenTTLMin int
	refreshTokenTTLHr int
}

func NewAuthUsecase(
	userUsecase UserUsecaseInterface,
	clientUsecase ClientUsecaseInterface,
	permissionUsecase PermissionUsecaseInterface,
	tokenUsecase TokenUsecaseInterface,
	securityService security.SecurityService,
	appName string,
	accessTokenTTLMinutes int,
	refreshTokenTTLHours int,
) AuthUsecaseInterface {
	return &authUsecase{
		userUsecase:       userUsecase,
		clientUsecase:     clientUsecase,
		permissionUsecase: permissionUsecase,
		tokenUsecase:      tokenUsecase,
		security:          securityService,
		appName:           appName,
		accessTokenTTLMin: accessTokenTTLMinutes,
		refreshTokenTTLHr: refreshTokenTTLHours,
	}
}

func (uc *authUsecase) Register(ctx context.Context, req dto.UserRegisterRequest) (*domain.User, error) {
	return uc.userUsecase.Register(ctx, req)
}

func (uc *authUsecase) Login(ctx context.Context, req dto.UserLoginRequest) (*dto.AuthLoginResponse, error) {
	user, err := uc.userUsecase.Login(ctx, req)
	if err != nil {
		return nil, err
	}

	var clientID *uuid.UUID
	if req.ClientID != "" {
		client, err := uc.clientUsecase.GetByClientID(ctx, req.ClientID)
		if err != nil {
			return nil, err
		}
		if client == nil {
			return nil, errors.New("invalid client")
		}
		clientID = &client.ID
	}

	roles := []interface{}{
		map[string]interface{}{
			"role":        "user",
			"permissions": []string{"profile:read"},
		},
	}

	accessToken, err := uc.security.GenerateJWT(user.ID.String(), uc.appName, roles, uc.accessTokenTTLMin)
	if err != nil {
		return nil, err
	}

	refreshToken := uuid.NewString()
	tokenRecord := &domain.Token{
		UserID:           user.ID,
		ClientID:         clientID,
		RefreshTokenHash: refreshToken,
		ExpiresAt:        time.Now().Add(time.Duration(uc.refreshTokenTTLHr) * time.Hour),
	}

	if err := uc.tokenUsecase.Create(ctx, tokenRecord); err != nil {
		return nil, err
	}

	return &dto.AuthLoginResponse{
		AccessToken:  accessToken,
		TokenType:    "Bearer",
		ExpiresIn:    uc.accessTokenTTLMin * 60,
		RefreshToken: refreshToken,
		User: dto.UserResponse{
			ID:       user.ID.String(),
			Username: user.Username,
			Email:    user.Email,
			FullName: user.FullName,
			Status:   user.Status,
		},
	}, nil
}

func (uc *authUsecase) Refresh(ctx context.Context, req dto.RefreshTokenRequest) (*dto.TokenResponse, error) {
	tokenRecord, err := uc.tokenUsecase.GetByRefreshTokenHash(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}
	if tokenRecord == nil {
		return nil, errors.New("invalid refresh token")
	}

	user, err := uc.userUsecase.GetByID(ctx, tokenRecord.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid refresh token")
	}

	roles := []interface{}{
		map[string]interface{}{
			"role":        "user",
			"permissions": []string{"profile:read"},
		},
	}

	accessToken, err := uc.security.GenerateJWT(user.ID.String(), uc.appName, roles, uc.accessTokenTTLMin)
	if err != nil {
		return nil, err
	}

	newRefreshToken := uuid.NewString()
	tokenRecord.RefreshTokenHash = newRefreshToken
	tokenRecord.ExpiresAt = time.Now().Add(time.Duration(uc.refreshTokenTTLHr) * time.Hour)

	if err := uc.tokenUsecase.Update(ctx, tokenRecord.ID, tokenRecord); err != nil {
		return nil, err
	}

	return &dto.TokenResponse{
		AccessToken:  accessToken,
		TokenType:    "Bearer",
		ExpiresIn:    uc.accessTokenTTLMin * 60,
		RefreshToken: newRefreshToken,
	}, nil
}

func (uc *authUsecase) Introspect(ctx context.Context, token string) (*dto.IntrospectResponse, error) {
	claims, err := uc.security.ValidateJWT(token)
	if err != nil {
		return nil, err
	}

	roles := make([]string, 0, len(claims.Roles))
	for _, rawRole := range claims.Roles {
		switch value := rawRole.(type) {
		case string:
			roles = append(roles, value)
		case map[string]interface{}:
			if roleName, ok := value["role"].(string); ok {
				roles = append(roles, roleName)
			}
		default:
			roles = append(roles, fmt.Sprintf("%v", value))
		}
	}

	return &dto.IntrospectResponse{
		Active:      true,
		User:        map[string]string{"sub": claims.Sub, "iss": claims.Iss},
		Roles:       roles,
		Permissions: []string{},
	}, nil
}

func (uc *authUsecase) Authorize(ctx context.Context, subject string, roles []interface{}, req dto.AuthorizeRequest) (*dto.AuthorizeResponse, error) {
	if req.UserID != subject {
		return nil, ErrUserIDMismatch
	}

	if !hasPermission(roles, req.Permission) {
		return nil, ErrNotAuthorized
	}

	return &dto.AuthorizeResponse{Authorized: true}, nil
}

func (uc *authUsecase) GetPermissionsByUserID(ctx context.Context, userID uuid.UUID) ([]string, error) {
	permissions, err := uc.permissionUsecase.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	permissionCodes := make([]string, 0, len(permissions))
	for _, permission := range permissions {
		permissionCodes = append(permissionCodes, permission.Code)
	}
	return permissionCodes, nil
}

func (uc *authUsecase) Logout(ctx context.Context, userID uuid.UUID) error {
	return uc.tokenUsecase.DeleteByUserID(ctx, userID)
}

func hasPermission(roles []interface{}, permission string) bool {
	for _, rawRole := range roles {
		if roleMap, ok := rawRole.(map[string]interface{}); ok {
			if roleName, ok := roleMap["role"].(string); ok && roleName == "admin" {
				return true
			}

			if permissions, ok := roleMap["permissions"].([]interface{}); ok {
				for _, rawPermission := range permissions {
					if fmt.Sprintf("%v", rawPermission) == permission {
						return true
					}
				}
			}
		}
	}
	return false
}
