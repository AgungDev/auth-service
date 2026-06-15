package http

import (
	"auth_service/internal/domain/dto"
	middleware "auth_service/internal/transport/http/middleware_http"
	"auth_service/internal/usecases"
	loggerpkg "auth_service/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type clientController struct {
	clientUsecase usecases.ClientUsecaseInterface
	logger        loggerpkg.Logger
}

// NewClientController creates a new client controller instance
func NewClientController(clientUsecase usecases.ClientUsecaseInterface, logger loggerpkg.Logger) *clientController {
	return &clientController{clientUsecase: clientUsecase, logger: logger}
}

// Route registers client-related routes
func (c *clientController) Route(rg *gin.RouterGroup) {
	clientGroup := rg.Group("/clients")
	{
		clientGroup.POST("", c.Create)
		clientGroup.GET("/:id", c.GetByID)
		clientGroup.GET("", c.GetAll)
		clientGroup.PUT("/:id", c.Update)
		clientGroup.DELETE("/:id", c.Delete)
	}
}

// Create creates a new client
func (c *clientController) Create(ctx *gin.Context) {
	correlationID := middleware.GetCorrelationID(ctx)
	c.logger.Info(correlationID, "", "anonymous", "client create request received")
	var req dto.ClientRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.Warn(correlationID, "", "anonymous", "client create validation failed: "+err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client, err := c.clientUsecase.Create(ctx.Request.Context(), req)
	if err != nil {
		c.logger.Error(correlationID, "", "anonymous", "client create failed: "+err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// convert related structs to string slices
	var redirectURIs []string
	for _, r := range client.RedirectURIs {
		redirectURIs = append(redirectURIs, r.RedirectURI)
	}
	var grants []string
	for _, g := range client.Grants {
		grants = append(grants, g.GrantType)
	}

	resp := dto.ClientResponse{
		ID:             client.ID.String(),
		ClientID:       client.ClientID,
		Name:           client.Name,
		RedirectURIs:   redirectURIs,
		Grants:         grants,
		IsConfidential: client.IsConfidential,
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Client created successfully",
		"data":    resp,
	})
}

// GetByID retrieves a client by ID
func (c *clientController) GetByID(ctx *gin.Context) {
	correlationID := middleware.GetCorrelationID(ctx)
	c.logger.Info(correlationID, "", "anonymous", "client get request received")
	clientID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid client id"})
		return
	}

	client, err := c.clientUsecase.GetByID(ctx.Request.Context(), clientID)
	if err != nil || client == nil {
		c.logger.Warn(correlationID, "", "anonymous", "client not found")
		ctx.JSON(http.StatusNotFound, gin.H{"error": "client not found"})
		return
	}

	var redirectURIs []string
	for _, r := range client.RedirectURIs {
		redirectURIs = append(redirectURIs, r.RedirectURI)
	}
	var grants []string
	for _, g := range client.Grants {
		grants = append(grants, g.GrantType)
	}

	resp := dto.ClientResponse{
		ID:             client.ID.String(),
		ClientID:       client.ClientID,
		Name:           client.Name,
		RedirectURIs:   redirectURIs,
		Grants:         grants,
		IsConfidential: client.IsConfidential,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Client retrieved successfully",
		"data":    resp,
	})
}

// GetAll retrieves all clients
func (c *clientController) GetAll(ctx *gin.Context) {
	correlationID := middleware.GetCorrelationID(ctx)
	c.logger.Info(correlationID, "", "anonymous", "client list request received")
	clients, err := c.clientUsecase.GetAll(ctx.Request.Context())
	if err != nil {
		c.logger.Error(correlationID, "", "anonymous", "client list failed: "+err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var responses []dto.ClientResponse
	for _, client := range clients {
		var redirectURIs []string
		for _, r := range client.RedirectURIs {
			redirectURIs = append(redirectURIs, r.RedirectURI)
		}
		var grants []string
		for _, g := range client.Grants {
			grants = append(grants, g.GrantType)
		}
		responses = append(responses, dto.ClientResponse{
			ID:             client.ID.String(),
			ClientID:       client.ClientID,
			Name:           client.Name,
			RedirectURIs:   redirectURIs,
			Grants:         grants,
			IsConfidential: client.IsConfidential,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Clients retrieved successfully",
		"data":    responses,
	})
}

// Update updates a client
func (c *clientController) Update(ctx *gin.Context) {
	correlationID := middleware.GetCorrelationID(ctx)
	c.logger.Info(correlationID, "", "anonymous", "client update request received")
	clientID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid client id"})
		return
	}

	var req dto.ClientRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.clientUsecase.Update(ctx.Request.Context(), clientID, req); err != nil {
		c.logger.Error(correlationID, "", "anonymous", "client update failed: "+err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Client updated successfully",
	})
}

// Delete removes a client
func (c *clientController) Delete(ctx *gin.Context) {
	correlationID := middleware.GetCorrelationID(ctx)
	c.logger.Info(correlationID, "", "anonymous", "client delete request received")
	clientID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid client id"})
		return
	}

	if err := c.clientUsecase.Delete(ctx.Request.Context(), clientID); err != nil {
		c.logger.Error(correlationID, "", "anonymous", "client delete failed: "+err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Client deleted successfully",
	})
}
