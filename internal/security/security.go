package security

import (
	"auth_service/internal/config"
	"auth_service/pkg"
)

type SecurityService interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
	GenerateJWT(sub, iss string, roles []interface{}, expireMinutes int) (string, error)
	ValidateJWT(tokenString string) (*pkg.Claims, error)
}

type securityService struct {
	jwtConfig config.JWTConfig
}

func NewSecurityService(jwtConfig config.JWTConfig) SecurityService {
	return &securityService{jwtConfig: jwtConfig}
}

func (s *securityService) HashPassword(password string) (string, error) {
	return pkg.HashPassword(password)
}

func (s *securityService) CheckPasswordHash(password, hash string) bool {
	return pkg.CheckPasswordHash(password, hash)
}

func (s *securityService) GenerateJWT(sub, iss string, roles []interface{}, expireMinutes int) (string, error) {
	return pkg.GenerateJWTRS256(sub, iss, roles, s.jwtConfig.PrivateKeyPath, expireMinutes)
}

func (s *securityService) ValidateJWT(tokenString string) (*pkg.Claims, error) {
	return pkg.ValidateJWTRS256(tokenString, s.jwtConfig.PublicKeyPath)
}
