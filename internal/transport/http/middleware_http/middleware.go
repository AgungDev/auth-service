package middleware_http

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"auth_service/internal/security"
	loggerpkg "auth_service/pkg/logger"
	"github.com/gin-gonic/gin"
)

const ContextKeyCorrelationID = "correlation_id"

func RequestLoggerMiddleware(log loggerpkg.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		correlationID := ctx.GetHeader("X-Correlation-ID")
		if correlationID == "" {
			correlationID = log.NewCorrelationID()
		}
		ctx.Set(ContextKeyCorrelationID, correlationID)
		ctx.Writer.Header().Set("X-Correlation-ID", correlationID)
		ctx.Request = ctx.Request.WithContext(
			context.WithValue(ctx.Request.Context(), ContextKeyCorrelationID, correlationID),
		)

		start := time.Now()
		ctx.Next()
		duration := time.Since(start)

		log.Info(
			correlationID,
			"",
			"",
			"request completed",
			slog.String("method", ctx.Request.Method),
			slog.String("path", ctx.FullPath()),
			slog.Int("status", ctx.Writer.Status()),
			slog.String("duration", duration.String()),
		)
	}
}

func GetCorrelationID(ctx *gin.Context) string {
	if cid, ok := ctx.Get(ContextKeyCorrelationID); ok {
		if str, ok := cid.(string); ok {
			return str
		}
		return fmt.Sprint(cid)
	}
	return ctx.GetHeader("X-Correlation-ID")
}

func AuthMiddleware(securityService security.SecurityService, logger loggerpkg.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		correlationID := GetCorrelationID(ctx)
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			logger.Warn(correlationID, "", "anonymous", "missing authorization header")
			ctx.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "authorization header is required", "error": "authorization header is required"})
			return
		}

		parts := strings.Fields(authHeader)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			logger.Warn(correlationID, "", "anonymous", "invalid authorization header format")
			ctx.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "authorization header must be Bearer token", "error": "invalid authorization format"})
			return
		}

		claims, err := securityService.ValidateJWT(parts[1])
		if err != nil {
			logger.Warn(correlationID, "", "anonymous", "invalid bearer token: "+err.Error())
			ctx.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "invalid or expired token", "error": "invalid or expired token"})
			return
		}

		ctx.Set("jwt_claims", claims)
		ctx.Set("user_subject", claims.Sub)
		ctx.Next()
	}
}
