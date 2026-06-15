package logger

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"strings"
	"time"

	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

const ContextKeyLogSource = "log_source"

type GormLogger struct {
	logger         *slog.Logger
	serviceName    string
	serviceVersion string
	environment    string
	logLevel       gormlogger.LogLevel
}

func NewGormLogger(serviceName, serviceVersion, environment string, level gormlogger.LogLevel) gormlogger.Interface {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: false})
	return &GormLogger{
		logger:         slog.New(handler),
		serviceName:    serviceName,
		serviceVersion: serviceVersion,
		environment:    environment,
		logLevel:       level,
	}
}

func (g *GormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	ng := *g
	ng.logLevel = level
	return &ng
}

func (g *GormLogger) Info(ctx context.Context, msg string, data ...any) {
	if g.logLevel < gormlogger.Info {
		return
	}
	g.logger.Info(msg, gormAttrs(g.serviceName, g.environment, ctx, data...)...)
}

func (g *GormLogger) Warn(ctx context.Context, msg string, data ...any) {
	if g.logLevel < gormlogger.Warn {
		return
	}
	g.logger.Warn(msg, gormAttrs(g.serviceName, g.environment, ctx, data...)...)
}

func (g *GormLogger) Error(ctx context.Context, msg string, data ...any) {
	if g.logLevel < gormlogger.Error {
		return
	}
	g.logger.Error(msg, gormAttrs(g.serviceName, g.environment, ctx, data...)...)
}

func (g *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if g.logLevel == gormlogger.Silent {
		return
	}

	sql, rows := fc()
	duration := time.Since(begin)
	operation, table := parseQueryMetadata(sql)

	attrs := []any{
		slog.String("service", g.serviceName),
	}
	if g.serviceVersion != "" {
		attrs = append(attrs, slog.String("version", g.serviceVersion))
	}
	if g.environment != "" {
		attrs = append(attrs, slog.String("environment", g.environment))
	}
	if source := getLogSource(ctx); source != "" {
		attrs = append(attrs, slog.String("source", source))
	}
	if cid := getCorrelationID(ctx); cid != "" {
		attrs = append(attrs, slog.String("correlation_id", cid))
	}
	if operation != "" {
		attrs = append(attrs, slog.String("operation", operation))
	}
	if table != "" {
		attrs = append(attrs, slog.String("table", table))
	}
	attrs = append(attrs,
		slog.Int64("rows", rows),
		slog.String("duration", duration.String()),
	)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if g.logLevel >= gormlogger.Warn {
				g.logger.Debug("record not found", attrs...)
			}
			return
		}
		if g.logLevel >= gormlogger.Error {
			g.logger.Error(err.Error(), attrs...)
		}
		return
	}

	if g.logLevel >= gormlogger.Info {
		g.logger.Info("gorm query", attrs...)
	}
}

func parseQueryMetadata(sql string) (string, string) {
	lower := strings.TrimSpace(strings.ToLower(sql))
	parts := strings.Fields(lower)
	if len(parts) == 0 {
		return "", ""
	}

	op := parts[0]
	table := ""

	switch op {
	case "select":
		for i, p := range parts {
			if p == "from" && i+1 < len(parts) {
				table = strings.Trim(parts[i+1], `"`)
				break
			}
		}
	case "insert":
		if len(parts) > 2 && parts[1] == "into" {
			table = strings.Trim(parts[2], `"`)
		}
	case "update":
		if len(parts) > 1 {
			table = strings.Trim(parts[1], `"`)
		}
	case "delete":
		for i, p := range parts {
			if p == "from" && i+1 < len(parts) {
				table = strings.Trim(parts[i+1], `"`)
				break
			}
		}
	}

	return op, table
}

func getCorrelationID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if val := ctx.Value("correlation_id"); val != nil {
		if s, ok := val.(string); ok {
			return s
		}
	}
	return ""
}

func getLogSource(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if val := ctx.Value(ContextKeyLogSource); val != nil {
		if s, ok := val.(string); ok {
			return s
		}
	}
	return ""
}

func gormAttrs(serviceName, environment string, ctx context.Context, data ...any) []any {
	attrs := []any{
		slog.String("service", serviceName),
	}
	if environment != "" {
		attrs = append(attrs, slog.String("environment", environment))
	}

	if cid := getCorrelationID(ctx); cid != "" {
		attrs = append(attrs, slog.String("correlation_id", cid))
	}

	if len(data) > 0 {
		attrs = append(attrs, slog.Any("data", data))
	}

	return attrs
}
