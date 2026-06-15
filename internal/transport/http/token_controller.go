package http

import (
	"auth_service/internal/domain"
	"auth_service/internal/domain/dto"
	"auth_service/internal/security"
	middleware "auth_service/internal/transport/http/middleware_http"
	"auth_service/internal/usecases"
	loggerpkg "auth_service/pkg/logger"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TokenControllerConfig struct {
	AppName               string
	AccessTokenTTLMinutes int
}

type tokenController struct {
	userUsecase   usecases.UserUsecaseInterface
	clientUsecase usecases.ClientUsecaseInterface
	tokenUsecase  usecases.TokenUsecaseInterface
	security      security.SecurityService
	config        TokenControllerConfig
	logger        loggerpkg.Logger
}

// NewTokenController creates a new token controller instance
func NewTokenController(
	userUsecase usecases.UserUsecaseInterface,
	clientUsecase usecases.ClientUsecaseInterface,
	tokenUsecase usecases.TokenUsecaseInterface,
	securityService security.SecurityService,
	config TokenControllerConfig,
	logger loggerpkg.Logger,
) *tokenController {
	return &tokenController{
		userUsecase:   userUsecase,
		clientUsecase: clientUsecase,
		tokenUsecase:  tokenUsecase,
		security:      securityService,
		config:        config,
		logger:        logger,
	}
}

// Route registers token-related routes
func (c *tokenController) Route(rg *gin.RouterGroup) {
	tokenGroup := rg.Group("/oauth")
	{
		tokenGroup.POST("/token", c.GetToken)
	}
}

// GetToken handles OAuth token generation
func (c *tokenController) GetToken(ctx *gin.Context) {
	correlationID := middleware.GetCorrelationID(ctx)
	c.logger.Info(correlationID, "", "anonymous", "token request received")
	var req dto.TokenRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.Warn(correlationID, "", "anonymous", "token request validation failed: "+err.Error())
		JSONError(ctx, http.StatusBadRequest, "invalid request")
		return
	}

	if req.GrantType != "password" && req.GrantType != "refresh_token" {
		c.logger.Warn(correlationID, "", "anonymous", "unsupported grant_type: "+req.GrantType)
		JSONError(ctx, http.StatusBadRequest, "unsupported grant_type")
		return
	}

	client, err := c.clientUsecase.GetByClientID(ctx.Request.Context(), req.ClientID)
	if err != nil || client == nil {
		c.logger.Warn(correlationID, "", "anonymous", "invalid client authentication")
		JSONError(ctx, http.StatusUnauthorized, "invalid client")
		return
	}

	if !c.security.CheckPasswordHash(req.ClientSecret, client.ClientSecretHash) {
		c.logger.Warn(correlationID, "", "anonymous", "invalid client credentials")
		JSONError(ctx, http.StatusUnauthorized, "invalid client credentials")
		return
	}

	if req.GrantType == "password" {
		c.handlePasswordGrant(ctx, req, client)
		return
	}

	c.handleRefreshTokenGrant(ctx, req, client)
}

// handlePasswordGrant handles the password grant flow
func (c *tokenController) handlePasswordGrant(ctx *gin.Context, req dto.TokenRequest, client *domain.Client) {
	correlationID := middleware.GetCorrelationID(ctx)
	loginReq := dto.UserLoginRequest{
		Username: req.Username,
		Password: req.Password,
	}

	user, err := c.userUsecase.Login(ctx.Request.Context(), loginReq)
	if err != nil {
		c.logger.Warn(correlationID, "", "anonymous", "token password grant login failed: "+err.Error())
		JSONError(ctx, http.StatusUnauthorized, "invalid credentials")
		return
	}

	roles := []interface{}{
		map[string]interface{}{
			"role":        "user",
			"permissions": []string{"profile:read"},
		},
	}

	accessToken, err := c.security.GenerateJWT(
		user.ID.String(),
		c.config.AppName,
		roles,
		c.config.AccessTokenTTLMinutes,
	)
	if err != nil {
		c.logger.Error(correlationID, "", "anonymous", "token generation failed: "+err.Error())
		JSONError(ctx, http.StatusInternalServerError, "token generation failed")
		return
	}

	refreshToken := uuid.NewString()
	clientID := client.ID
	tokenRecord := &domain.Token{
		UserID:           user.ID,
		ClientID:         &clientID,
		RefreshTokenHash: refreshToken,
		ExpiresAt:        time.Now().Add(24 * time.Hour),
	}

	if err := c.tokenUsecase.Create(ctx.Request.Context(), tokenRecord); err != nil {
		c.logger.Error(correlationID, "", "anonymous", "refresh token save failed: "+err.Error())
		JSONError(ctx, http.StatusInternalServerError, "unable to store refresh token")
		return
	}

	resp := dto.TokenResponse{
		AccessToken:  accessToken,
		TokenType:    "Bearer",
		ExpiresIn:    c.config.AccessTokenTTLMinutes * 60,
		RefreshToken: refreshToken,
	}

	JSONOK(ctx, "token granted", resp)
}

// handleRefreshTokenGrant handles the refresh token grant flow
func (c *tokenController) handleRefreshTokenGrant(ctx *gin.Context, req dto.TokenRequest, client *domain.Client) {
	correlationID := middleware.GetCorrelationID(ctx)
	c.logger.Info(correlationID, "", "anonymous", "refresh token grant request received")

	if req.RefreshToken == "" {
		JSONError(ctx, http.StatusBadRequest, "refresh_token is required")
		return
	}

	tokenRecord, err := c.tokenUsecase.GetByRefreshTokenHash(ctx.Request.Context(), req.RefreshToken)
	if err != nil || tokenRecord == nil {
		c.logger.Warn(correlationID, "", "anonymous", "refresh token not found")
		JSONError(ctx, http.StatusUnauthorized, "invalid refresh token")
		return
	}

	if c.tokenUsecase.IsTokenExpired(tokenRecord) {
		c.logger.Warn(correlationID, "", "anonymous", "refresh token expired")
		JSONError(ctx, http.StatusUnauthorized, "refresh token expired")
		return
	}

	user, err := c.userUsecase.GetByID(ctx.Request.Context(), tokenRecord.UserID)
	if err != nil || user == nil {
		c.logger.Warn(correlationID, "", "anonymous", "user not found for refresh token")
		JSONError(ctx, http.StatusUnauthorized, "invalid refresh token")
		return
	}

	roles := []interface{}{
		map[string]interface{}{
			"role":        "user",
			"permissions": []string{"profile:read"},
		},
	}

	accessToken, err := c.security.GenerateJWT(
		user.ID.String(),
		c.config.AppName,
		roles,
		c.config.AccessTokenTTLMinutes,
	)
	if err != nil {
		c.logger.Error(correlationID, "", "anonymous", "refresh token generation failed: "+err.Error())
		JSONError(ctx, http.StatusInternalServerError, "refresh token generation failed")
		return
	}

	newRefreshToken := uuid.NewString()
	tokenRecord.RefreshTokenHash = newRefreshToken
	tokenRecord.ExpiresAt = time.Now().Add(24 * time.Hour)
	if err := c.tokenUsecase.Update(ctx.Request.Context(), tokenRecord.ID, tokenRecord); err != nil {
		c.logger.Error(correlationID, "", "anonymous", "refresh token update failed: "+err.Error())
		JSONError(ctx, http.StatusInternalServerError, "unable to update refresh token")
		return
	}

	resp := dto.TokenResponse{
		AccessToken:  accessToken,
		TokenType:    "Bearer",
		ExpiresIn:    c.config.AccessTokenTTLMinutes * 60,
		RefreshToken: newRefreshToken,
	}

	JSONOK(ctx, "token refreshed", resp)
}
