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

type userController struct {
	userUsecase usecases.UserUsecaseInterface
	logger      loggerpkg.Logger
}

// NewUserController creates a new user controller instance
func NewUserController(userUsecase usecases.UserUsecaseInterface, logger loggerpkg.Logger) *userController {
	return &userController{userUsecase: userUsecase, logger: logger}
}

// Route registers user-related routes
func (c *userController) Route(rg *gin.RouterGroup) {
	userGroup := rg.Group("/users")
	{
		userGroup.POST("", c.Create)
		userGroup.GET("/:id", c.GetByID)
		userGroup.GET("", c.GetAll)
		userGroup.PUT("/:id", c.Update)
		userGroup.DELETE("/:id", c.Delete)
		userGroup.POST("/:id/roles", c.AssignRoles)
		userGroup.GET("/:id/roles", c.GetRoles)
		userGroup.DELETE("/:id/roles/:roleId", c.RemoveRole)
		userGroup.POST("/:id/permissions", c.AssignPermissions)
		userGroup.GET("/:id/permissions", c.GetPermissions)
		userGroup.DELETE("/:id/permissions/:permissionId", c.RemovePermission)
	}
}

// GetByID retrieves a user by ID
func (c *userController) GetByID(ctx *gin.Context) {
	correlationID := middleware.GetCorrelationID(ctx)
	c.logger.Info(correlationID, "", "anonymous", "user get request received")
	userID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		JSONError(ctx, http.StatusBadRequest, "invalid user id")
		return
	}

	user, err := c.userUsecase.GetByID(ctx.Request.Context(), userID)
	if err != nil || user == nil {
		c.logger.Warn(correlationID, "", "anonymous", "user not found")
		JSONError(ctx, http.StatusNotFound, "user not found")
		return
	}

	resp := dto.UserResponse{
		ID:       user.ID.String(),
		Username: user.Username,
		Email:    user.Email,
		FullName: user.FullName,
		Status:   user.Status,
	}

	JSONOK(ctx, "user retrieved", resp)
}

// GetAll retrieves all users
func (c *userController) GetAll(ctx *gin.Context) {
	correlationID := middleware.GetCorrelationID(ctx)
	c.logger.Info(correlationID, "", "anonymous", "user list request received")
	users, err := c.userUsecase.GetAll(ctx.Request.Context())
	if err != nil {
		c.logger.Error(correlationID, "", "anonymous", "user list failed: "+err.Error())
		JSONError(ctx, http.StatusInternalServerError, "failed to retrieve users")
		return
	}

	var responses []dto.UserResponse
	for _, user := range users {
		responses = append(responses, dto.UserResponse{
			ID:       user.ID.String(),
			Username: user.Username,
			Email:    user.Email,
			FullName: user.FullName,
			Status:   user.Status,
		})
	}

	JSONOK(ctx, "users retrieved", responses)
}

// Update updates user information
func (c *userController) Update(ctx *gin.Context) {
	correlationID := middleware.GetCorrelationID(ctx)
	c.logger.Info(correlationID, "", "anonymous", "user update request received")
	userID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		JSONError(ctx, http.StatusBadRequest, "invalid user id")
		return
	}

	var req dto.UserUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.Warn(correlationID, "", "anonymous", "user update validation failed: "+err.Error())
		JSONError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := c.userUsecase.Update(ctx.Request.Context(), userID, req); err != nil {
		c.logger.Error(correlationID, "", "anonymous", "user update failed: "+err.Error())
		JSONError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	JSONOK(ctx, "user updated", nil)
}

// Delete removes a user
func (c *userController) Delete(ctx *gin.Context) {
	correlationID := middleware.GetCorrelationID(ctx)
	c.logger.Info(correlationID, "", "anonymous", "user delete request received")
	userID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		JSONError(ctx, http.StatusBadRequest, "invalid user id")
		return
	}

	if err := c.userUsecase.Delete(ctx.Request.Context(), userID); err != nil {
		c.logger.Error(correlationID, "", "anonymous", "user delete failed: "+err.Error())
		JSONError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	JSONOK(ctx, "user deleted", nil)
}

func (c *userController) Create(ctx *gin.Context) {
	correlationID := middleware.GetCorrelationID(ctx)
	c.logger.Info(correlationID, "", "anonymous", "user create request received")

	var req dto.UserRegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.Warn(correlationID, "", "anonymous", "user create validation failed: "+err.Error())
		JSONError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	user, err := c.userUsecase.Register(ctx.Request.Context(), req)
	if err != nil {
		c.logger.Error(correlationID, "", "anonymous", "user create failed: "+err.Error())
		JSONError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := dto.UserResponse{
		ID:       user.ID.String(),
		Username: user.Username,
		Email:    user.Email,
		FullName: user.FullName,
		Status:   user.Status,
	}

	JSONCreated(ctx, "user created", resp)
}

func (c *userController) AssignRoles(ctx *gin.Context) {
	correlationID := middleware.GetCorrelationID(ctx)
	c.logger.Info(correlationID, "", "anonymous", "user assign roles request received")

	userID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		JSONError(ctx, http.StatusBadRequest, "invalid user id")
		return
	}

	var req dto.RoleAssignmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.Warn(correlationID, "", "anonymous", "assign roles validation failed: "+err.Error())
		JSONError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	roleIDs := make([]uuid.UUID, 0, len(req.RoleIDs))
	for _, roleID := range req.RoleIDs {
		parsed, err := uuid.Parse(roleID)
		if err != nil {
			JSONError(ctx, http.StatusBadRequest, "invalid role id")
			return
		}
		roleIDs = append(roleIDs, parsed)
	}

	if err := c.userUsecase.AssignRoles(ctx.Request.Context(), userID, roleIDs); err != nil {
		c.logger.Error(correlationID, "", "anonymous", "assign roles failed: "+err.Error())
		JSONError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	JSONOK(ctx, "roles assigned", nil)
}

func (c *userController) GetRoles(ctx *gin.Context) {
	correlationID := middleware.GetCorrelationID(ctx)
	c.logger.Info(correlationID, "", "anonymous", "user roles request received")

	userID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		JSONError(ctx, http.StatusBadRequest, "invalid user id")
		return
	}

	roles, err := c.userUsecase.GetRolesByUserID(ctx.Request.Context(), userID)
	if err != nil {
		c.logger.Error(correlationID, "", "anonymous", "get user roles failed: "+err.Error())
		JSONError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	var resp []dto.RoleResponse
	for _, role := range roles {
		resp = append(resp, dto.RoleResponse{
			ID:          role.ID.String(),
			Code:        role.Code,
			Name:        role.Name,
			Description: role.Description,
		})
	}

	JSONOK(ctx, "user roles retrieved", resp)
}

func (c *userController) RemoveRole(ctx *gin.Context) {
	correlationID := middleware.GetCorrelationID(ctx)
	c.logger.Info(correlationID, "", "anonymous", "user remove role request received")

	userID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		JSONError(ctx, http.StatusBadRequest, "invalid user id")
		return
	}

	roleID, err := uuid.Parse(ctx.Param("roleId"))
	if err != nil {
		JSONError(ctx, http.StatusBadRequest, "invalid role id")
		return
	}

	if err := c.userUsecase.RemoveRole(ctx.Request.Context(), userID, roleID); err != nil {
		c.logger.Error(correlationID, "", "anonymous", "remove role failed: "+err.Error())
		JSONError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	JSONOK(ctx, "role removed", nil)
}

func (c *userController) AssignPermissions(ctx *gin.Context) {
	correlationID := middleware.GetCorrelationID(ctx)
	c.logger.Info(correlationID, "", "anonymous", "user assign permissions request received")

	userID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		JSONError(ctx, http.StatusBadRequest, "invalid user id")
		return
	}

	var req dto.PermissionAssignmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.Warn(correlationID, "", "anonymous", "assign permissions validation failed: "+err.Error())
		JSONError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	permissionIDs := make([]uuid.UUID, 0, len(req.PermissionIDs))
	for _, permissionID := range req.PermissionIDs {
		parsed, err := uuid.Parse(permissionID)
		if err != nil {
			JSONError(ctx, http.StatusBadRequest, "invalid permission id")
			return
		}
		permissionIDs = append(permissionIDs, parsed)
	}

	if err := c.userUsecase.AssignPermissions(ctx.Request.Context(), userID, permissionIDs); err != nil {
		c.logger.Error(correlationID, "", "anonymous", "assign permissions failed: "+err.Error())
		JSONError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	JSONOK(ctx, "permissions assigned", nil)
}

func (c *userController) GetPermissions(ctx *gin.Context) {
	correlationID := middleware.GetCorrelationID(ctx)
	c.logger.Info(correlationID, "", "anonymous", "user permissions request received")

	userID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		JSONError(ctx, http.StatusBadRequest, "invalid user id")
		return
	}

	permissions, err := c.userUsecase.GetPermissions(ctx.Request.Context(), userID)
	if err != nil {
		c.logger.Error(correlationID, "", "anonymous", "get user permissions failed: "+err.Error())
		JSONError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	var resp []dto.PermissionResponse
	for _, permission := range permissions {
		resp = append(resp, dto.PermissionResponse{
			ID:          permission.ID.String(),
			Code:        permission.Code,
			Name:        permission.Name,
			Description: permission.Description,
		})
	}

	JSONOK(ctx, "user permissions retrieved", resp)
}

func (c *userController) RemovePermission(ctx *gin.Context) {
	correlationID := middleware.GetCorrelationID(ctx)
	c.logger.Info(correlationID, "", "anonymous", "user remove permission request received")

	userID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		JSONError(ctx, http.StatusBadRequest, "invalid user id")
		return
	}

	permissionID, err := uuid.Parse(ctx.Param("permissionId"))
	if err != nil {
		JSONError(ctx, http.StatusBadRequest, "invalid permission id")
		return
	}

	if err := c.userUsecase.RemovePermission(ctx.Request.Context(), userID, permissionID); err != nil {
		c.logger.Error(correlationID, "", "anonymous", "remove permission failed: "+err.Error())
		JSONError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	JSONOK(ctx, "permission removed", nil)
}
