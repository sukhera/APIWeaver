package services

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/sukhera/APIWeaver/internal/config"
	"github.com/sukhera/APIWeaver/internal/domain/generator"
	"github.com/sukhera/APIWeaver/internal/domain/parser"
)

// GenerationResult represents the result of OpenAPI generation
type GenerationResult struct {
	Content  string             `json:"content"`
	Format   string             `json:"format"`
	Metadata GenerationMetadata `json:"metadata"`
	Warnings []string           `json:"warnings,omitempty"`
	Errors   []string           `json:"errors,omitempty"`
}

// GenerationMetadata contains metadata about the generation process
type GenerationMetadata struct {
	ProcessingTimeMs int `json:"processing_time_ms"`
	InputSizeBytes   int `json:"input_size_bytes"`
	OutputSizeBytes  int `json:"output_size_bytes"`
	EndpointCount    int `json:"endpoint_count"`
	ComponentCount   int `json:"component_count"`
}

// Generator service handles OpenAPI spec generation
type Generator struct {
	config    *config.ExtendedConfig
	logger    *slog.Logger
	parser    *parser.Parser
	generator *generator.Generator
}

// NewGenerator creates a new Generator service
func NewGenerator(cfg *config.ExtendedConfig, logger *slog.Logger) *Generator {
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

	// Create generator
	generatorInstance := generator.New(generator.Config{
		Format:          cfg.OutputFormat,
		PrettyPrint:     cfg.PrettyPrint,
		IncludeExamples: true,
		ValidateOutput:  true,
		StrictMode:      cfg.StrictMode,
	})

	return &Generator{
		config:    cfg,
		logger:    logger,
		parser:    parserInstance,
		generator: generatorInstance,
	}
}

// Generate generates an OpenAPI specification from Markdown content
func (g *Generator) Generate(ctx context.Context, content string, format string) (*GenerationResult, error) {
	startTime := time.Now()

	g.logger.InfoContext(ctx, "Starting OpenAPI generation",
		"input_size", len(content),
		"format", format,
	)

	// Parse the markdown content
	doc, err := g.parser.ParseWithContext(ctx, content)
	if err != nil {
		g.logger.ErrorContext(ctx, "Failed to parse markdown", "error", err)
		return nil, fmt.Errorf("failed to parse markdown: %w", err)
	}

	// Check for parse errors
	var parseErrors []string
	var parseWarnings []string

	for _, parseErr := range doc.Errors {
		if parseErr.IsError() {
			parseErrors = append(parseErrors, parseErr.Error())
		} else if parseErr.IsWarning() {
			parseWarnings = append(parseWarnings, parseErr.Error())
		}
	}

	// If there are fatal errors and we're in strict mode, return early
	if len(parseErrors) > 0 && g.config.StrictMode {
		return &GenerationResult{
			Format:   format,
			Errors:   parseErrors,
			Warnings: parseWarnings,
			Metadata: GenerationMetadata{
				ProcessingTimeMs: int(time.Since(startTime).Milliseconds()),
				InputSizeBytes:   len(content),
				EndpointCount:    len(doc.Endpoints),
				ComponentCount:   len(doc.Components),
			},
		}, fmt.Errorf("parsing failed with %d errors", len(parseErrors))
	}

	// Generate OpenAPI specification
	spec, err := g.generator.Generate(ctx, doc, format)
	if err != nil {
		g.logger.ErrorContext(ctx, "Failed to generate OpenAPI spec", "error", err)
		return nil, fmt.Errorf("failed to generate OpenAPI spec: %w", err)
	}

	processingTime := time.Since(startTime)

	result := &GenerationResult{
		Content:  spec,
		Format:   format,
		Warnings: parseWarnings,
		Errors:   parseErrors,
		Metadata: GenerationMetadata{
			ProcessingTimeMs: int(processingTime.Milliseconds()),
			InputSizeBytes:   len(content),
			OutputSizeBytes:  len(spec),
			EndpointCount:    len(doc.Endpoints),
			ComponentCount:   len(doc.Components),
		},
	}

	g.logger.InfoContext(ctx, "OpenAPI generation completed",
		"processing_time_ms", result.Metadata.ProcessingTimeMs,
		"endpoint_count", result.Metadata.EndpointCount,
		"component_count", result.Metadata.ComponentCount,
		"output_size", result.Metadata.OutputSizeBytes,
		"warnings", len(parseWarnings),
		"errors", len(parseErrors),
	)

	return result, nil
}

// GenerateFromFile generates an OpenAPI specification from a Markdown file
func (g *Generator) GenerateFromFile(ctx context.Context, filename string, format string) (*GenerationResult, error) {
	g.logger.InfoContext(ctx, "Generating from file", "filename", filename)

	// This would read the file and call Generate
	// Implementation would use the file utilities from common package
	return nil, fmt.Errorf("not implemented")
}

// ValidateInput validates markdown input before generation
func (g *Generator) ValidateInput(ctx context.Context, content string) error {
	if content == "" {
		return fmt.Errorf("input content is empty")
	}

	// Additional validation logic would go here
	return nil
}
