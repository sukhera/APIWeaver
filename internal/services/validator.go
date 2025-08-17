package services

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/sukhera/APIWeaver/internal/config"
	"github.com/sukhera/APIWeaver/internal/domain/parser"
	"github.com/sukhera/APIWeaver/internal/domain/validator"
)

// ValidationResult represents the result of validation
type ValidationResult struct {
	Valid       bool               `json:"valid"`
	Errors      []string           `json:"errors,omitempty"`
	Warnings    []string           `json:"warnings,omitempty"`
	Suggestions []string           `json:"suggestions,omitempty"`
	Metadata    ValidationMetadata `json:"metadata"`
}

// ValidationMetadata contains metadata about the validation process
type ValidationMetadata struct {
	ProcessingTimeMs int    `json:"processing_time_ms"`
	InputSizeBytes   int    `json:"input_size_bytes"`
	InputType        string `json:"input_type"`
	ValidatorVersion string `json:"validator_version"`
}

// Validator service handles validation of Markdown and OpenAPI specifications
type Validator struct {
	config           *config.ExtendedConfig
	logger           *slog.Logger
	parser           *parser.Parser
	openapiValidator *validator.OpenAPIValidator
}

// NewValidator creates a new Validator service
func NewValidator(cfg *config.ExtendedConfig, logger *slog.Logger) *Validator {
	// Create parser with configuration
	parserInstance := parser.New(
		parser.WithStrictMode(cfg.StrictMode),
		parser.WithRecovery(cfg.EnableRecovery, cfg.MaxRecoveryAttempts),
		parser.WithTimeout(cfg.ParserTimeout),
		parser.WithAllowedMethods(cfg.AllowedMethods),
		parser.WithValidationLevel(cfg.ValidationLevel),
		parser.WithRequireExamples(cfg.RequireExamples),
		parser.WithMaxNestingDepth(cfg.MaxNestingDepth),
		parser.WithInitialSliceCapacity(cfg.InitialSliceCapacity),
	)

	// Create OpenAPI validator
	openapiValidator := validator.NewOpenAPIValidator(validator.Config{
		StrictMode:         cfg.StrictMode,
		ValidateExamples:   true,
		CheckBestPractices: true,
		AllowExtensions:    true,
	})

	return &Validator{
		config:           cfg,
		logger:           logger,
		parser:           parserInstance,
		openapiValidator: openapiValidator,
	}
}

// Validate validates content based on its type (markdown or openapi)
func (v *Validator) Validate(ctx context.Context, content, inputType string) (*ValidationResult, error) {
	startTime := time.Now()

	v.logger.InfoContext(ctx, "Starting validation",
		"input_size", len(content),
		"input_type", inputType,
	)

	var result *ValidationResult
	var err error

	switch inputType {
	case "markdown":
		result, err = v.validateMarkdown(ctx, content)
	case "openapi":
		result, err = v.validateOpenAPI(ctx, content)
	default:
		return nil, fmt.Errorf("unsupported input type: %s", inputType)
	}

	if err != nil {
		v.logger.ErrorContext(ctx, "Validation failed", "error", err)
		return nil, err
	}

	// Update metadata
	result.Metadata.ProcessingTimeMs = int(time.Since(startTime).Milliseconds())
	result.Metadata.InputSizeBytes = len(content)
	result.Metadata.InputType = inputType
	result.Metadata.ValidatorVersion = "1.0.0" // TODO: Get from build info

	v.logger.InfoContext(ctx, "Validation completed",
		"processing_time_ms", result.Metadata.ProcessingTimeMs,
		"input_type", inputType,
		"valid", result.Valid,
		"errors", len(result.Errors),
		"warnings", len(result.Warnings),
		"suggestions", len(result.Suggestions),
	)

	return result, nil
}

// validateMarkdown validates Markdown content for APIWeaver format compliance
func (v *Validator) validateMarkdown(ctx context.Context, content string) (*ValidationResult, error) {
	// Parse the markdown content
	doc, err := v.parser.ParseWithContext(ctx, content)
	if err != nil {
		return &ValidationResult{
			Valid:    false,
			Errors:   []string{err.Error()},
			Metadata: ValidationMetadata{},
		}, nil
	}

	var errors []string
	var warnings []string
	var suggestions []string

	// Collect parse errors and warnings
	for _, parseErr := range doc.Errors {
		if parseErr.IsError() {
			errors = append(errors, parseErr.Error())
		} else if parseErr.IsWarning() {
			warnings = append(warnings, parseErr.Error())
		}
	}

	// Additional validation rules
	if len(doc.Endpoints) == 0 {
		warnings = append(warnings, "No endpoints found in the document")
	}

	if doc.Frontmatter == nil {
		suggestions = append(suggestions, "Consider adding YAML frontmatter with API metadata")
	}

	// Check for missing descriptions
	for _, endpoint := range doc.Endpoints {
		if endpoint.Description == "" {
			suggestions = append(suggestions, fmt.Sprintf("Endpoint %s %s is missing a description", endpoint.Method, endpoint.Path))
		}
	}

	result := &ValidationResult{
		Valid:       len(errors) == 0,
		Errors:      errors,
		Warnings:    warnings,
		Suggestions: suggestions,
		Metadata:    ValidationMetadata{},
	}

	return result, nil
}

// validateOpenAPI validates OpenAPI specification content
func (v *Validator) validateOpenAPI(ctx context.Context, content string) (*ValidationResult, error) {
	// Use the OpenAPI validator
	validationResult, err := v.openapiValidator.Validate(ctx, content)
	if err != nil {
		return &ValidationResult{
			Valid:    false,
			Errors:   []string{err.Error()},
			Metadata: ValidationMetadata{},
		}, nil
	}

	// Convert validator result to service result
	result := &ValidationResult{
		Valid:       validationResult.Valid,
		Errors:      validationResult.Errors,
		Warnings:    validationResult.Warnings,
		Suggestions: validationResult.Suggestions,
		Metadata:    ValidationMetadata{},
	}

	return result, nil
}

// ValidateFile validates a file based on its extension
func (v *Validator) ValidateFile(ctx context.Context, filename string) (*ValidationResult, error) {
	v.logger.InfoContext(ctx, "Validating file", "filename", filename)

	// This would read the file and determine type based on extension
	// Implementation would use the file utilities from common package
	return nil, fmt.Errorf("not implemented")
}
