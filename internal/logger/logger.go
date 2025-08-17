package logger

import (
	"io"
	"log/slog"
	"os"
)

// Config represents logger configuration
type Config struct {
	Level      string `mapstructure:"level" json:"level"`
	Format     string `mapstructure:"format" json:"format"` // "json" or "text"
	Output     string `mapstructure:"output" json:"output"` // "stdout", "stderr", or file path
	AddSource  bool   `mapstructure:"add_source" json:"add_source"`
	TimeFormat string `mapstructure:"time_format" json:"time_format"`
}

// DefaultConfig returns default logger configuration
func DefaultConfig() Config {
	return Config{
		Level:      "info",
		Format:     "json",
		Output:     "stdout",
		AddSource:  false,
		TimeFormat: "2006-01-02T15:04:05.000Z07:00",
	}
}

// New creates a new structured logger based on configuration
func New(cfg Config) (*slog.Logger, error) {
	// Set default values if empty
	if cfg.Level == "" {
		cfg.Level = "info"
	}
	if cfg.Format == "" {
		cfg.Format = "json"
	}
	if cfg.Output == "" {
		cfg.Output = "stdout"
	}

	// Parse log level
	var level slog.Level
	switch cfg.Level {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn", "warning":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	// Determine output writer
	var writer io.Writer
	switch cfg.Output {
	case "stdout":
		writer = os.Stdout
	case "stderr":
		writer = os.Stderr
	default:
		// Assume it's a file path
		file, err := os.OpenFile(cfg.Output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, err
		}
		writer = file
	}

	// Create handler options
	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: cfg.AddSource,
	}

	// Create appropriate handler
	var handler slog.Handler
	switch cfg.Format {
	case "text":
		handler = slog.NewTextHandler(writer, opts)
	default: // "json"
		handler = slog.NewJSONHandler(writer, opts)
	}

	// Create and return logger
	logger := slog.New(handler)
	return logger, nil
}

// WithCorrelationID adds a correlation ID to all log entries
func WithCorrelationID(logger *slog.Logger, correlationID string) *slog.Logger {
	return logger.With("correlation_id", correlationID)
}

// WithComponent adds a component name to all log entries
func WithComponent(logger *slog.Logger, component string) *slog.Logger {
	return logger.With("component", component)
}