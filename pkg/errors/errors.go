package errors

import (
	"fmt"
	"strings"
)

// ParseError represents a parsing error with context
type ParseError struct {
	Type       ErrorType `json:"type"`
	Code       string    `json:"code,omitempty"`
	Message    string    `json:"message"`
	LineNumber int       `json:"line_number"`
	Column     int       `json:"column,omitempty"`
	Context    string    `json:"context,omitempty"`
	Suggestion string    `json:"suggestion,omitempty"`
	Source     string    `json:"source,omitempty"` // e.g., "frontmatter", "endpoint", "schema"
	Severity   Severity  `json:"severity"`
}

// ErrorType represents different categories of parsing errors
type ErrorType string

const (
	ErrorTypeSyntax      ErrorType = "syntax"
	ErrorTypeValidation  ErrorType = "validation"
	ErrorTypeConfig      ErrorType = "config"
	ErrorTypeTimeout     ErrorType = "timeout"
	ErrorTypeSchema      ErrorType = "schema"
	ErrorTypeTable       ErrorType = "table"
	ErrorTypeFrontmatter ErrorType = "frontmatter"
	ErrorTypeEndpoint    ErrorType = "endpoint"
	ErrorTypeReference   ErrorType = "reference"
)

// Severity represents the severity level of an error
type Severity string

const (
	SeverityInfo    Severity = "info"
	SeverityWarning Severity = "warning"
	SeverityError   Severity = "error"
	SeverityFatal   Severity = "fatal"
)

// Error implements the error interface
func (e *ParseError) Error() string {
	var parts []string

	if e.Code != "" {
		parts = append(parts, fmt.Sprintf("[%s]", e.Code))
	}

	if e.LineNumber > 0 {
		if e.Column > 0 {
			parts = append(parts, fmt.Sprintf("line %d:%d", e.LineNumber, e.Column))
		} else {
			parts = append(parts, fmt.Sprintf("line %d", e.LineNumber))
		}
	}

	if e.Source != "" {
		parts = append(parts, fmt.Sprintf("in %s", e.Source))
	}

	parts = append(parts, e.Message)

	result := strings.Join(parts, " ")

	if e.Suggestion != "" {
		result += fmt.Sprintf(" (suggestion: %s)", e.Suggestion)
	}

	return result
}

// IsError returns true if this is an error-level issue
func (e *ParseError) IsError() bool {
	return e.Severity == SeverityError || e.Severity == SeverityFatal
}

// IsWarning returns true if this is a warning-level issue
func (e *ParseError) IsWarning() bool {
	return e.Severity == SeverityWarning
}

// IsFatal returns true if this is a fatal error
func (e *ParseError) IsFatal() bool {
	return e.Severity == SeverityFatal
}

// ConfigError represents configuration-related errors
type ConfigError struct {
	Field   string
	Value   interface{}
	Message string
}

func (e *ConfigError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("configuration error in field '%s': %s", e.Field, e.Message)
	}
	return fmt.Sprintf("configuration error: %s", e.Message)
}

// NewConfigError creates a new configuration error
func NewConfigError(message string) *ConfigError {
	return &ConfigError{Message: message}
}

// TimeoutError represents timeout-related errors
type TimeoutError struct {
	Operation string
	Duration  string
}

func (e *TimeoutError) Error() string {
	return fmt.Sprintf("operation '%s' timed out after %s", e.Operation, e.Duration)
}

// NewTimeoutError creates a new timeout error
func NewTimeoutError(operation, duration string) *TimeoutError {
	return &TimeoutError{
		Operation: operation,
		Duration:  duration,
	}
}
