package http

import (
	"fmt"
	"net/http"

	"auth_service/internal/domain/dto"

	"github.com/gin-gonic/gin"
)

func JSONSuccess(ctx *gin.Context, status int, message string, data any) {
	ctx.JSON(status, dto.APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func JSONCreated(ctx *gin.Context, message string, data any) {
	JSONSuccess(ctx, http.StatusCreated, message, data)
}

func JSONOK(ctx *gin.Context, message string, data any) {
	JSONSuccess(ctx, http.StatusOK, message, data)
}

func JSONNoContent(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusNoContent, dto.APIResponse{
		Success: true,
		Message: message,
	})
}

func JSONError(ctx *gin.Context, status int, message string) {
	ctx.JSON(status, dto.APIResponse{
		Success: false,
		Message: message,
		Error:   message,
	})
}

func JSONErrorf(ctx *gin.Context, status int, format string, a ...any) {
	message := fmt.Sprintf(format, a...)
	JSONError(ctx, status, message)
}
