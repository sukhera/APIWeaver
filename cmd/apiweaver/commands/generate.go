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

// NewGenerateCmd creates the generate command
func NewGenerateCmd() *cobra.Command {
	var (
		outputFile   string
		outputFormat string
		configFile   string
		verbose      bool
	)

	cmd := &cobra.Command{
		Use:   "generate [input-file]",
		Short: "Generate OpenAPI specification from Markdown",
		Long: `Generate a complete OpenAPI 3.1 specification from a structured Markdown file.
The input file should contain API documentation in APIWeaver's Markdown format
with endpoints, parameters, and response definitions.`,
		Args: cobra.ExactArgs(1),
		Example: `  apiweaver generate api-docs.md
  apiweaver generate docs.md --output openapi.yaml --format yaml
  apiweaver generate example.md --config config.yaml --verbose`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGenerate(cmd.Context(), args[0], outputFile, outputFormat, configFile, verbose)
		},
	}

	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file for generated OpenAPI spec")
	cmd.Flags().StringVarP(&outputFormat, "format", "f", "yaml", "Output format (yaml, json)")
	cmd.Flags().StringVarP(&configFile, "config", "c", "", "Configuration file path")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")

	return cmd
}

func runGenerate(ctx context.Context, inputFile, outputFile, outputFormat, configFile string, verbose bool) error {
	// Load configuration
	cfg, err := config.Load(configFile)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Override with command line flags
	if verbose {
		cfg.Verbose = true
	}

	// Setup logger
	log, err := logger.New(cfg.Logger)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}

	log.Info("Starting OpenAPI generation",
		"input_file", inputFile,
		"output_file", outputFile,
		"format", outputFormat,
	)

	// Clean and validate input file path
	inputFile = filepath.Clean(inputFile)

	// Read input file
	content, err := os.ReadFile(inputFile) // #nosec G304 - file path is from CLI argument
	if err != nil {
		return fmt.Errorf("failed to read input file %s: %w", inputFile, err)
	}

	// Create generator service
	generatorService := services.NewGenerator(cfg, log)

	// Generate OpenAPI spec
	spec, err := generatorService.Generate(ctx, string(content), outputFormat)
	if err != nil {
		log.Error("Generation failed", "error", err)
		return fmt.Errorf("failed to generate OpenAPI spec: %w", err)
	}

	// Output result
	if outputFile != "" {
		if err := os.WriteFile(outputFile, []byte(spec.Content), 0600); err != nil {
			return fmt.Errorf("failed to write output file %s: %w", outputFile, err)
		}
		log.Info("OpenAPI specification generated successfully", "output_file", outputFile)
	} else {
		fmt.Print(spec.Content)
	}

	// Print summary
	if verbose {
		fmt.Fprintf(os.Stderr, "\nGeneration Summary:\n")
		fmt.Fprintf(os.Stderr, "  Endpoints: %d\n", spec.Metadata.EndpointCount)
		fmt.Fprintf(os.Stderr, "  Components: %d\n", spec.Metadata.ComponentCount)
		fmt.Fprintf(os.Stderr, "  Processing time: %dms\n", spec.Metadata.ProcessingTimeMs)
		if len(spec.Warnings) > 0 {
			fmt.Fprintf(os.Stderr, "  Warnings: %d\n", len(spec.Warnings))
		}
	}

	return nil
}
