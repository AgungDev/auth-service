package logger

import (
	"crypto/rand"
	"fmt"
	"log/slog"
	"os"
)

type Logger interface {
	Info(correlationID string, userID any, role string, message string, attrs ...slog.Attr)
	Warn(correlationID string, userID any, role string, message string, attrs ...slog.Attr)
	Error(correlationID string, userID any, role string, message string, attrs ...slog.Attr)
	NewCorrelationID() string
}

type loggerImpl struct {
	logger         *slog.Logger
	serviceName    string
	serviceVersion string
	environment    string
}

func NewLogger(serviceName, serviceVersion, environment string) Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: false})
	return &loggerImpl{
		logger:         slog.New(handler),
		serviceName:    serviceName,
		serviceVersion: serviceVersion,
		environment:    environment,
	}
}

func (l *loggerImpl) attrs(correlationID string, userID any, role string, attrs ...slog.Attr) []any {
	fields := []any{
		slog.String("service", l.serviceName),
	}
	if l.serviceVersion != "" {
		fields = append(fields, slog.String("version", l.serviceVersion))
	}
	if l.environment != "" {
		fields = append(fields, slog.String("environment", l.environment))
	}
	if correlationID != "" {
		fields = append(fields, slog.String("correlation_id", correlationID))
	}
	if userID != nil && userID != "" {
		fields = append(fields, slog.Any("user_id", userID))
	}
	if role != "" {
		fields = append(fields, slog.String("role", role))
	}
	for _, attr := range attrs {
		fields = append(fields, attr)
	}
	return fields
}

func (l *loggerImpl) Info(correlationID string, userID any, role string, message string, attrs ...slog.Attr) {
	l.logger.With(l.attrs(correlationID, userID, role, attrs...)...).Info(message)
}

func (l *loggerImpl) Warn(correlationID string, userID any, role string, message string, attrs ...slog.Attr) {
	l.logger.With(l.attrs(correlationID, userID, role, attrs...)...).Warn(message)
}

func (l *loggerImpl) Error(correlationID string, userID any, role string, message string, attrs ...slog.Attr) {
	l.logger.With(l.attrs(correlationID, userID, role, attrs...)...).Error(message)
}

func (l *loggerImpl) NewCorrelationID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("%d", b)
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
