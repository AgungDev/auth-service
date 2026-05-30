package http

import (
	"auth_service/internal/domain/dto"
	"auth_service/internal/usecases"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type clientController struct {
	clientUsecase usecases.ClientUsecaseInterface
}

// NewClientController creates a new client controller instance
func NewClientController(clientUsecase usecases.ClientUsecaseInterface) *clientController {
	return &clientController{clientUsecase: clientUsecase}
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
	var req dto.ClientRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client, err := c.clientUsecase.Create(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp := dto.ClientResponse{
		ID:             client.ID,
		ClientID:       client.ClientID,
		RedirectURIs:   client.RedirectURIs,
		Grants:         client.Grants,
		IsConfidential: client.IsConfidential,
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Client created successfully",
		"data":    resp,
	})
}

// GetByID retrieves a client by ID
func (c *clientController) GetByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid client id"})
		return
	}

	client, err := c.clientUsecase.GetByID(uint(id))
	if err != nil || client == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "client not found"})
		return
	}

	resp := dto.ClientResponse{
		ID:             client.ID,
		ClientID:       client.ClientID,
		RedirectURIs:   client.RedirectURIs,
		Grants:         client.Grants,
		IsConfidential: client.IsConfidential,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Client retrieved successfully",
		"data":    resp,
	})
}

// GetAll retrieves all clients
func (c *clientController) GetAll(ctx *gin.Context) {
	clients, err := c.clientUsecase.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var responses []dto.ClientResponse
	for _, client := range clients {
		responses = append(responses, dto.ClientResponse{
			ID:             client.ID,
			ClientID:       client.ClientID,
			RedirectURIs:   client.RedirectURIs,
			Grants:         client.Grants,
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
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid client id"})
		return
	}

	var req dto.ClientRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.clientUsecase.Update(uint(id), req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Client updated successfully",
	})
}

// Delete removes a client
func (c *clientController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid client id"})
		return
	}

	if err := c.clientUsecase.Delete(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Client deleted successfully",
	})
}
