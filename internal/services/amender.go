package services

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/sukhera/APIWeaver/internal/config"
	"github.com/sukhera/APIWeaver/internal/domain/amender"
)

// AmendmentResult represents the result of OpenAPI amendment
type AmendmentResult struct {
	Content   string            `json:"content"`
	Format    string            `json:"format"`
	Changes   []string          `json:"changes"`
	Conflicts []string          `json:"conflicts,omitempty"`
	Metadata  AmendmentMetadata `json:"metadata"`
	Warnings  []string          `json:"warnings,omitempty"`
	Errors    []string          `json:"errors,omitempty"`
}

// AmendmentMetadata contains metadata about the amendment process
type AmendmentMetadata struct {
	ProcessingTimeMs  int `json:"processing_time_ms"`
	InputSizeBytes    int `json:"input_size_bytes"`
	OutputSizeBytes   int `json:"output_size_bytes"`
	ChangesApplied    int `json:"changes_applied"`
	ConflictsResolved int `json:"conflicts_resolved"`
}

// Amender service handles OpenAPI spec amendments
type Amender struct {
	config  *config.ExtendedConfig
	logger  *slog.Logger
	amender *amender.Amender
}

// NewAmender creates a new Amender service
func NewAmender(cfg *config.ExtendedConfig, logger *slog.Logger) *Amender {
	// Create amender with configuration
	amenderInstance := amender.New(amender.Config{
		StrictMode:           cfg.StrictMode,
		AllowBreakingChanges: !cfg.StrictMode,
		ConflictResolution:   amender.ConflictResolutionAuto,
		ValidateOutput:       true,
	})

	return &Amender{
		config:  cfg,
		logger:  logger,
		amender: amenderInstance,
	}
}

// Amend applies changes to an existing OpenAPI specification
func (a *Amender) Amend(ctx context.Context, existingSpec, changes, format string, dryRun bool) (*AmendmentResult, error) {
	startTime := time.Now()

	a.logger.InfoContext(ctx, "Starting OpenAPI amendment",
		"spec_size", len(existingSpec),
		"changes_size", len(changes),
		"format", format,
		"dry_run", dryRun,
	)

	// Parse the existing specification
	spec, err := a.amender.ParseSpec(ctx, existingSpec, format)
	if err != nil {
		a.logger.ErrorContext(ctx, "Failed to parse existing spec", "error", err)
		return nil, fmt.Errorf("failed to parse existing specification: %w", err)
	}

	// Parse the changes
	changeSet, err := a.amender.ParseChanges(ctx, changes)
	if err != nil {
		a.logger.ErrorContext(ctx, "Failed to parse changes", "error", err)
		return nil, fmt.Errorf("failed to parse changes: %w", err)
	}

	// Apply changes
	result, err := a.amender.ApplyChanges(ctx, spec, changeSet, dryRun)
	if err != nil {
		a.logger.ErrorContext(ctx, "Failed to apply changes", "error", err)
		return nil, fmt.Errorf("failed to apply changes: %w", err)
	}

	// Serialize result if not dry run
	var content string
	var outputSize int
	if !dryRun {
		content, err = a.amender.SerializeSpec(ctx, result.Spec, format)
		if err != nil {
			a.logger.ErrorContext(ctx, "Failed to serialize amended spec", "error", err)
			return nil, fmt.Errorf("failed to serialize amended specification: %w", err)
		}
		outputSize = len(content)
	}

	processingTime := time.Since(startTime)

	amendmentResult := &AmendmentResult{
		Content:   content,
		Format:    format,
		Changes:   result.Changes,
		Conflicts: result.Conflicts,
		Warnings:  result.Warnings,
		Errors:    result.Errors,
		Metadata: AmendmentMetadata{
			ProcessingTimeMs:  int(processingTime.Milliseconds()),
			InputSizeBytes:    len(existingSpec) + len(changes),
			OutputSizeBytes:   outputSize,
			ChangesApplied:    len(result.Changes),
			ConflictsResolved: len(result.Conflicts),
		},
	}

	a.logger.InfoContext(ctx, "OpenAPI amendment completed",
		"processing_time_ms", amendmentResult.Metadata.ProcessingTimeMs,
		"changes_applied", amendmentResult.Metadata.ChangesApplied,
		"conflicts_resolved", amendmentResult.Metadata.ConflictsResolved,
		"output_size", amendmentResult.Metadata.OutputSizeBytes,
		"warnings", len(result.Warnings),
		"errors", len(result.Errors),
		"dry_run", dryRun,
	)

	return amendmentResult, nil
}

// ValidateChanges validates changes before applying them
func (a *Amender) ValidateChanges(ctx context.Context, changes string) error {
	if changes == "" {
		return fmt.Errorf("changes content is empty")
	}

	// Parse and validate changes
	_, err := a.amender.ParseChanges(ctx, changes)
	if err != nil {
		return fmt.Errorf("invalid changes format: %w", err)
	}

	return nil
}

// PreviewChanges provides a preview of what changes would be applied
func (a *Amender) PreviewChanges(ctx context.Context, existingSpec, changes, format string) (*AmendmentResult, error) {
	return a.Amend(ctx, existingSpec, changes, format, true)
}
