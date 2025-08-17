package validator

import (
	"context"
	"fmt"
	"strings"
)

// Config holds validator configuration
type Config struct {
	StrictMode         bool
	ValidateExamples   bool
	CheckBestPractices bool
	AllowExtensions    bool
}

// OpenAPIValidator validates OpenAPI specifications
type OpenAPIValidator struct {
	config Config
}

// NewOpenAPIValidator creates a new OpenAPI validator
func NewOpenAPIValidator(config Config) *OpenAPIValidator {
	return &OpenAPIValidator{
		config: config,
	}
}

// ValidationResult represents the result of validation
type ValidationResult struct {
	Valid       bool
	Errors      []string
	Warnings    []string
	Suggestions []string
}

// Validate validates an OpenAPI specification
func (v *OpenAPIValidator) Validate(ctx context.Context, content string) (*ValidationResult, error) {
	var errors []string
	var warnings []string
	var suggestions []string

	// Basic validation - check if it looks like OpenAPI
	if !strings.Contains(content, "openapi") {
		errors = append(errors, "Missing 'openapi' field")
	}

	if !strings.Contains(content, "info") {
		errors = append(errors, "Missing 'info' object")
	}

	if !strings.Contains(content, "paths") {
		warnings = append(warnings, "No 'paths' object found - API has no endpoints")
	}

	// Version validation
	if strings.Contains(content, "openapi: 2.") || strings.Contains(content, `"openapi": "2.`) {
		warnings = append(warnings, "OpenAPI 2.x (Swagger) detected - consider upgrading to OpenAPI 3.1")
	}

	// Best practices check
	if v.config.CheckBestPractices {
		if !strings.Contains(content, "description") {
			suggestions = append(suggestions, "Consider adding descriptions to improve API documentation")
		}

		if !strings.Contains(content, "examples") && !strings.Contains(content, "example") {
			suggestions = append(suggestions, "Consider adding examples to improve API usability")
		}

		if !strings.Contains(content, "components") {
			suggestions = append(suggestions, "Consider using components for reusable schemas")
		}
	}

	// Strict mode checks
	if v.config.StrictMode {
		if strings.Contains(content, "x-") && !v.config.AllowExtensions {
			warnings = append(warnings, "OpenAPI extensions (x-*) found in strict mode")
		}
	}

	// Example validation
	if v.config.ValidateExamples {
		// This would validate that examples match their schemas
		// Mock implementation
		if strings.Contains(content, "example") {
			suggestions = append(suggestions, "Examples found - ensure they match their schemas")
		}
	}

	result := &ValidationResult{
		Valid:       len(errors) == 0,
		Errors:      errors,
		Warnings:    warnings,
		Suggestions: suggestions,
	}

	return result, nil
}

// ValidateSchema validates a JSON schema
func (v *OpenAPIValidator) ValidateSchema(ctx context.Context, schema map[string]interface{}) error {
	// Mock implementation
	if schema == nil {
		return fmt.Errorf("schema is nil")
	}

	// Check required fields
	if _, ok := schema["type"]; !ok {
		return fmt.Errorf("schema missing 'type' field")
	}

	return nil
}

// ValidateExample validates an example against a schema
func (v *OpenAPIValidator) ValidateExample(ctx context.Context, example interface{}, schema map[string]interface{}) error {
	// Mock implementation
	if example == nil {
		return fmt.Errorf("example is nil")
	}

	if schema == nil {
		return fmt.Errorf("schema is nil")
	}

	// This would validate the example against the schema
	return nil
}