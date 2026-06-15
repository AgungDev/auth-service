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

type roleController struct {
	roleUsecase usecases.RoleUsecaseInterface
	logger      loggerpkg.Logger
}

// NewRoleController creates a new role controller instance
func NewRoleController(roleUsecase usecases.RoleUsecaseInterface, logger loggerpkg.Logger) *roleController {
	return &roleController{roleUsecase: roleUsecase, logger: logger}
}

// Route registers role-related routes
func (c *roleController) Route(rg *gin.RouterGroup) {
	roleGroup := rg.Group("/roles")
	{
		roleGroup.POST("", c.Create)
		roleGroup.GET("/:id", c.GetByID)
		roleGroup.GET("", c.GetAll)
		roleGroup.PUT("/:id", c.Update)
		roleGroup.DELETE("/:id", c.Delete)
		roleGroup.POST("/:id/permissions", c.AssignPermissions)
		roleGroup.GET("/:id/permissions", c.GetPermissions)
		roleGroup.DELETE("/:id/permissions/:permissionId", c.RemovePermission)
	}
}

// Create creates a new role
func (c *roleController) Create(ctx *gin.Context) {
	correlationID := middleware.GetCorrelationID(ctx)
	c.logger.Info(correlationID, "", "anonymous", "role create request received")
	var req dto.RoleRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.Warn(correlationID, "", "anonymous", "role create validation failed: "+err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role, err := c.roleUsecase.Create(ctx.Request.Context(), req)
	if err != nil {
		c.logger.Error(correlationID, "", "anonymous", "role create failed: "+err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp := dto.RoleResponse{
		ID:          role.ID.String(),
		Code:        role.Code,
		Name:        role.Name,
		Description: role.Description,
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Role created successfully",
		"data":    resp,
	})
}

// GetByID retrieves a role by ID
func (c *roleController) GetByID(ctx *gin.Context) {
	correlationID := middleware.GetCorrelationID(ctx)
	c.logger.Info(correlationID, "", "anonymous", "role get request received")
	roleID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid role id"})
		return
	}

	role, err := c.roleUsecase.GetByID(ctx.Request.Context(), roleID)
	if err != nil || role == nil {
		c.logger.Warn(correlationID, "", "anonymous", "role not found")
		ctx.JSON(http.StatusNotFound, gin.H{"error": "role not found"})
		return
	}

	resp := dto.RoleResponse{
		ID:          role.ID.String(),
		Code:        role.Code,
		Name:        role.Name,
		Description: role.Description,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Role retrieved successfully",
		"data":    resp,
	})
}

// GetAll retrieves all roles
func (c *roleController) GetAll(ctx *gin.Context) {
	correlationID := middleware.GetCorrelationID(ctx)
	c.logger.Info(correlationID, "", "anonymous", "role list request received")
	roles, err := c.roleUsecase.GetAll(ctx.Request.Context())
	if err != nil {
		c.logger.Error(correlationID, "", "anonymous", "role list failed: "+err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var responses []dto.RoleResponse
	for _, role := range roles {
		responses = append(responses, dto.RoleResponse{
			ID:          role.ID.String(),
			Code:        role.Code,
			Name:        role.Name,
			Description: role.Description,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Roles retrieved successfully",
		"data":    responses,
	})
}

// Update updates a role
func (c *roleController) Update(ctx *gin.Context) {
	correlationID := middleware.GetCorrelationID(ctx)
	c.logger.Info(correlationID, "", "anonymous", "role update request received")
	roleID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid role id"})
		return
	}

	var req dto.RoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.roleUsecase.Update(ctx.Request.Context(), roleID, req); err != nil {
		c.logger.Error(correlationID, "", "anonymous", "role update failed: "+err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Role updated successfully",
	})
}

// Delete removes a role
func (c *roleController) Delete(ctx *gin.Context) {
	correlationID := middleware.GetCorrelationID(ctx)
	c.logger.Info(correlationID, "", "anonymous", "role delete request received")
	roleID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid role id"})
		return
	}

	if err := c.roleUsecase.Delete(ctx.Request.Context(), roleID); err != nil {
		c.logger.Error(correlationID, "", "anonymous", "role delete failed: "+err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Role deleted successfully",
	})
}

func (c *roleController) AssignPermissions(ctx *gin.Context) {
	correlationID := middleware.GetCorrelationID(ctx)
	c.logger.Info(correlationID, "", "anonymous", "role assign permissions request received")

	roleID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid role id"})
		return
	}

	var req dto.PermissionAssignmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	permissionIDs := make([]uuid.UUID, 0, len(req.PermissionIDs))
	for _, permissionID := range req.PermissionIDs {
		parsed, err := uuid.Parse(permissionID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid permission id"})
			return
		}
		permissionIDs = append(permissionIDs, parsed)
	}

	if err := c.roleUsecase.AssignPermissions(ctx.Request.Context(), roleID, permissionIDs); err != nil {
		c.logger.Error(correlationID, "", "anonymous", "assign role permissions failed: "+err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Role permissions assigned successfully"})
}

func (c *roleController) GetPermissions(ctx *gin.Context) {
	correlationID := middleware.GetCorrelationID(ctx)
	c.logger.Info(correlationID, "", "anonymous", "role permissions request received")

	roleID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid role id"})
		return
	}

	permissions, err := c.roleUsecase.GetPermissionsByRoleID(ctx.Request.Context(), roleID)
	if err != nil {
		c.logger.Error(correlationID, "", "anonymous", "get role permissions failed: "+err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

	ctx.JSON(http.StatusOK, gin.H{"message": "Role permissions retrieved successfully", "data": resp})
}

func (c *roleController) RemovePermission(ctx *gin.Context) {
	correlationID := middleware.GetCorrelationID(ctx)
	c.logger.Info(correlationID, "", "anonymous", "role remove permission request received")

	roleID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid role id"})
		return
	}

	permissionID, err := uuid.Parse(ctx.Param("permissionId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid permission id"})
		return
	}

	if err := c.roleUsecase.RemovePermission(ctx.Request.Context(), roleID, permissionID); err != nil {
		c.logger.Error(correlationID, "", "anonymous", "remove role permission failed: "+err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Role permission removed successfully"})
}
