package http

import (
	"auth_service/internal/config"
	"auth_service/internal/repository"
	"auth_service/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TokenRequest struct {
	GrantType    string `json:"grant_type"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

var (
	userRepo *repository.UserRepository
	cfg      *config.Config
)

func SetUserRepo(repo *repository.UserRepository) {
	userRepo = repo
}

func SetConfig(c *config.Config) {
	cfg = c
}

func RegisterRoutes(r *gin.Engine) {
	r.POST("/oauth/token", TokenHandler)
}

func TokenHandler(c *gin.Context) {
	var req TokenRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if req.GrantType != "password" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported grant_type"})
		return
	}
	user, err := userRepo.FindByUsername(req.Username)
	if err != nil || !user.IsActive {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	if !pkg.CheckPasswordHash(req.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	roles := []interface{}{
		map[string]interface{}{
			"role":        "student", // TODO: ambil dari DB
			"permissions": []string{"courses:read", "profile:read"},
		},
	}
	accessToken, err := pkg.GenerateJWTRS256(
		"user:"+req.Username,
		cfg.AppName,
		roles,
		cfg.JWTPrivateKeyPath,
		15,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token generation failed"})
		return
	}
	resp := TokenResponse{
		AccessToken:  accessToken,
		TokenType:    "Bearer",
		ExpiresIn:    900,
		RefreshToken: "dummy-refresh-token",
	}
	c.JSON(http.StatusOK, resp)
}
