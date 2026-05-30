package http

import (
	"auth_service/internal/domain/dto"
	"auth_service/internal/usecases"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type userController struct {
	userUsecase usecases.UserUsecaseInterface
}

// NewUserController creates a new user controller instance
func NewUserController(userUsecase usecases.UserUsecaseInterface) *userController {
	return &userController{userUsecase: userUsecase}
}

// Route registers user-related routes
func (c *userController) Route(rg *gin.RouterGroup) {
	userGroup := rg.Group("/users")
	{
		userGroup.POST("/register", c.Register)
		userGroup.POST("/login", c.Login)
		userGroup.GET("/:id", c.GetByID)
		userGroup.GET("", c.GetAll)
		userGroup.PUT("/:id", c.Update)
		userGroup.DELETE("/:id", c.Delete)
	}
}

// Register handles user registration
func (c *userController) Register(ctx *gin.Context) {
	var req dto.UserRegisterRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := c.userUsecase.Register(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp := dto.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		FullName: user.FullName,
		IsActive: user.IsActive,
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"data":    resp,
	})
}

// Login handles user login
func (c *userController) Login(ctx *gin.Context) {
	var req dto.UserLoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := c.userUsecase.Login(req)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	resp := dto.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		FullName: user.FullName,
		IsActive: user.IsActive,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"data":    resp,
	})
}

// GetByID retrieves a user by ID
func (c *userController) GetByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	user, err := c.userUsecase.GetByID(uint(id))
	if err != nil || user == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	resp := dto.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		FullName: user.FullName,
		IsActive: user.IsActive,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User retrieved successfully",
		"data":    resp,
	})
}

// GetAll retrieves all users
func (c *userController) GetAll(ctx *gin.Context) {
	users, err := c.userUsecase.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var responses []dto.UserResponse
	for _, user := range users {
		responses = append(responses, dto.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			FullName: user.FullName,
			IsActive: user.IsActive,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Users retrieved successfully",
		"data":    responses,
	})
}

// Update updates user information
func (c *userController) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var req dto.UserUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.userUsecase.Update(uint(id), req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
	})
}

// Delete removes a user
func (c *userController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	if err := c.userUsecase.Delete(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}
