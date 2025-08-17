package commands

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/sukhera/APIWeaver/internal/config"
	"github.com/sukhera/APIWeaver/internal/logger"
	"github.com/sukhera/APIWeaver/internal/services"
)

// NewValidateCmd creates the validate command
func NewValidateCmd() *cobra.Command {
	var (
		inputType    string
		configFile   string
		verbose      bool
		strict       bool
		outputFormat string
	)

	cmd := &cobra.Command{
		Use:   "validate [input-file]",
		Short: "Validate Markdown or OpenAPI specification",
		Long: `Validate either a Markdown file for APIWeaver format compliance
or an OpenAPI specification for standard compliance and best practices.`,
		Args: cobra.ExactArgs(1),
		Example: `  apiweaver validate api-docs.md --type markdown
  apiweaver validate openapi.yaml --type openapi --strict
  apiweaver validate spec.json --type openapi --format json --verbose`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runValidate(cmd.Context(), args[0], inputType, configFile, verbose, strict, outputFormat)
		},
	}

	cmd.Flags().StringVarP(&inputType, "type", "t", "", "Input type: markdown, openapi (auto-detected if not specified)")
	cmd.Flags().StringVarP(&configFile, "config", "c", "", "Configuration file path")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
	cmd.Flags().BoolVarP(&strict, "strict", "s", false, "Enable strict validation mode")
	cmd.Flags().StringVarP(&outputFormat, "format", "f", "text", "Output format (text, json)")

	return cmd
}

func runValidate(ctx context.Context, inputFile, inputType, configFile string, verbose, strict bool, outputFormat string) error {
	// Load configuration
	cfg, err := config.Load(configFile)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Override with command line flags
	if verbose {
		cfg.Verbose = true
	}
	if strict {
		cfg.StrictMode = true
	}

	// Setup logger
	log, err := logger.New(cfg.Logger)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}

	// Auto-detect input type if not specified
	if inputType == "" {
		inputType = detectInputType(inputFile)
	}

	log.Info("Starting validation",
		"input_file", inputFile,
		"input_type", inputType,
		"strict", strict,
	)

	// Clean and validate input file path
	inputFile = filepath.Clean(inputFile)

	// Read input file
	content, err := os.ReadFile(inputFile) // #nosec G304 - file path is from CLI argument
	if err != nil {
		return fmt.Errorf("failed to read input file %s: %w", inputFile, err)
	}

	// Create validator service
	validatorService := services.NewValidator(cfg, log)

	// Validate content
	result, err := validatorService.Validate(ctx, string(content), inputType)
	if err != nil {
		log.Error("Validation failed", "error", err)
		return fmt.Errorf("validation failed: %w", err)
	}

	// Output results based on format
	switch outputFormat {
	case "json":
		return outputValidationJSON(result)
	default:
		return outputValidationText(result, verbose)
	}
}

func detectInputType(filename string) string {
	// Check file extension
	if len(filename) > 3 && filename[len(filename)-3:] == ".md" {
		return "markdown"
	}
	if len(filename) > 5 && (filename[len(filename)-5:] == ".yaml" || filename[len(filename)-5:] == ".json") {
		return "openapi"
	}
	if len(filename) > 4 && filename[len(filename)-4:] == ".yml" {
		return "openapi"
	}

	// Default to markdown for unknown extensions
	return "markdown"
}

func outputValidationText(result *services.ValidationResult, verbose bool) error {
	// Print overall result
	if result.Valid {
		fmt.Printf("✅ Validation PASSED\n\n")
	} else {
		fmt.Printf("❌ Validation FAILED\n\n")
	}

	// Print errors
	if len(result.Errors) > 0 {
		fmt.Printf("Errors (%d):\n", len(result.Errors))
		for i, err := range result.Errors {
			fmt.Printf("  %d. %s\n", i+1, err)
		}
		fmt.Println()
	}

	// Print warnings
	if len(result.Warnings) > 0 {
		fmt.Printf("Warnings (%d):\n", len(result.Warnings))
		for i, warning := range result.Warnings {
			fmt.Printf("  %d. %s\n", i+1, warning)
		}
		fmt.Println()
	}

	// Print suggestions if verbose
	if verbose && len(result.Suggestions) > 0 {
		fmt.Printf("Suggestions (%d):\n", len(result.Suggestions))
		for i, suggestion := range result.Suggestions {
			fmt.Printf("  %d. %s\n", i+1, suggestion)
		}
		fmt.Println()
	}

	// Print summary
	fmt.Printf("Summary:\n")
	fmt.Printf("  Total issues: %d\n", len(result.Errors)+len(result.Warnings))
	fmt.Printf("  Processing time: %dms\n", result.Metadata.ProcessingTimeMs)

	if !result.Valid {
		os.Exit(1)
	}

	return nil
}

func outputValidationJSON(result *services.ValidationResult) error {
	// This would output the validation result as JSON
	// Implementation would marshal the result to JSON
	fmt.Printf("JSON output not yet implemented\n")
	return nil
}
