package logger

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"
)

// HTTPRequestFormatter formats HTTP request logs
type HTTPRequestFormatter struct {
	logger *slog.Logger
}

// NewHTTPRequestFormatter creates a new HTTP request formatter
func NewHTTPRequestFormatter(logger *slog.Logger) *HTTPRequestFormatter {
	return &HTTPRequestFormatter{logger: logger}
}

// LogRequest logs an HTTP request with consistent formatting
func (f *HTTPRequestFormatter) LogRequest(ctx context.Context, method, path, remoteAddr, userAgent string, statusCode int, duration time.Duration, bodySize int64) {
	f.logger.InfoContext(ctx, "HTTP request",
		"method", method,
		"path", path,
		"status", statusCode,
		"duration_ms", duration.Milliseconds(),
		"remote_addr", remoteAddr,
		"user_agent", userAgent,
		"body_size", bodySize,
	)
}

// ErrorFormatter formats error logs consistently
type ErrorFormatter struct {
	logger *slog.Logger
}

// NewErrorFormatter creates a new error formatter
func NewErrorFormatter(logger *slog.Logger) *ErrorFormatter {
	return &ErrorFormatter{logger: logger}
}

// LogError logs an error with context and stack trace
func (f *ErrorFormatter) LogError(ctx context.Context, err error, operation string, metadata map[string]interface{}) {
	args := []interface{}{
		"error", err.Error(),
		"operation", operation,
	}

	// Add metadata
	for key, value := range metadata {
		args = append(args, key, value)
	}

	f.logger.ErrorContext(ctx, "Operation failed", args...)
}

// LogParseError logs parsing errors with line/column information
func (f *ErrorFormatter) LogParseError(ctx context.Context, err error, filename string, line, column int, context string) {
	f.logger.ErrorContext(ctx, "Parse error",
		"error", err.Error(),
		"file", filename,
		"line", line,
		"column", column,
		"context", context,
	)
}

// MetricsFormatter formats metrics logs
type MetricsFormatter struct {
	logger *slog.Logger
}

// NewMetricsFormatter creates a new metrics formatter
func NewMetricsFormatter(logger *slog.Logger) *MetricsFormatter {
	return &MetricsFormatter{logger: logger}
}

// LogPerformanceMetrics logs performance metrics
func (f *MetricsFormatter) LogPerformanceMetrics(ctx context.Context, operation string, duration time.Duration, metadata map[string]interface{}) {
	args := []interface{}{
		"operation", operation,
		"duration_ms", duration.Milliseconds(),
	}

	// Add metadata
	for key, value := range metadata {
		args = append(args, key, value)
	}

	f.logger.InfoContext(ctx, "Performance metrics", args...)
}

// SecurityFormatter formats security-related logs
type SecurityFormatter struct {
	logger *slog.Logger
}

// NewSecurityFormatter creates a new security formatter
func NewSecurityFormatter(logger *slog.Logger) *SecurityFormatter {
	return &SecurityFormatter{logger: logger}
}

// LogSecurityEvent logs security events
func (f *SecurityFormatter) LogSecurityEvent(ctx context.Context, eventType, description string, metadata map[string]interface{}) {
	args := []interface{}{
		"event_type", eventType,
		"description", description,
	}

	// Add metadata
	for key, value := range metadata {
		args = append(args, key, value)
	}

	f.logger.WarnContext(ctx, "Security event", args...)
}

// SanitizeUserInput sanitizes user input for logging to prevent log injection
func SanitizeUserInput(input string) string {
	// Remove or escape potentially dangerous characters
	input = strings.ReplaceAll(input, "\n", "\\n")
	input = strings.ReplaceAll(input, "\r", "\\r")
	input = strings.ReplaceAll(input, "\t", "\\t")
	
	// Limit length to prevent log flooding
	if len(input) > 1000 {
		input = input[:1000] + "...[truncated]"
	}
	
	return input
}

// FormatError formats an error with additional context
func FormatError(err error, operation string, metadata map[string]interface{}) string {
	var parts []string
	parts = append(parts, fmt.Sprintf("operation=%s", operation))
	parts = append(parts, fmt.Sprintf("error=%s", err.Error()))
	
	for key, value := range metadata {
		parts = append(parts, fmt.Sprintf("%s=%v", key, value))
	}
	
	return strings.Join(parts, " ")
}