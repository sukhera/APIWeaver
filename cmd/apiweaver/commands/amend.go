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

// NewAmendCmd creates the amend command
func NewAmendCmd() *cobra.Command {
	var (
		changesFile  string
		outputFile   string
		outputFormat string
		configFile   string
		verbose      bool
		dryRun       bool
	)

	cmd := &cobra.Command{
		Use:   "amend [existing-spec-file]",
		Short: "Amend existing OpenAPI specification",
		Long: `Amend an existing OpenAPI specification with changes described in Markdown format.
The changes file should contain descriptions of modifications to apply to the existing spec.`,
		Args: cobra.ExactArgs(1),
		Example: `  apiweaver amend openapi.yaml --changes changes.md
  apiweaver amend api.json --changes updates.md --output updated-api.yaml
  apiweaver amend spec.yaml --changes mods.md --dry-run --verbose`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAmend(cmd.Context(), args[0], changesFile, outputFile, outputFormat, configFile, verbose, dryRun)
		},
	}

	cmd.Flags().StringVarP(&changesFile, "changes", "c", "", "Markdown file describing changes to apply (required)")
	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file for amended spec (defaults to overwrite input)")
	cmd.Flags().StringVarP(&outputFormat, "format", "f", "", "Output format (yaml, json) (auto-detected if not specified)")
	cmd.Flags().StringVar(&configFile, "config", "", "Configuration file path")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be changed without applying")

	if err := cmd.MarkFlagRequired("changes"); err != nil {
		// This should never fail for a valid flag name
		panic(fmt.Sprintf("failed to mark flag as required: %v", err))
	}

	return cmd
}

func runAmend(ctx context.Context, specFile, changesFile, outputFile, outputFormat, configFile string, verbose, dryRun bool) error {
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

	log.Info("Starting OpenAPI amendment",
		"spec_file", specFile,
		"changes_file", changesFile,
		"dry_run", dryRun,
	)

	// Clean and validate file paths
	specFile = filepath.Clean(specFile)
	changesFile = filepath.Clean(changesFile)

	// Read existing spec
	specContent, err := os.ReadFile(specFile) // #nosec G304 - file path is from CLI argument
	if err != nil {
		return fmt.Errorf("failed to read spec file %s: %w", specFile, err)
	}

	// Read changes
	changesContent, err := os.ReadFile(changesFile) // #nosec G304 - file path is from CLI argument
	if err != nil {
		return fmt.Errorf("failed to read changes file %s: %w", changesFile, err)
	}

	// Auto-detect output format if not specified
	if outputFormat == "" {
		outputFormat = detectFormat(specFile)
	}

	// Create amender service
	amenderService := services.NewAmender(cfg, log)

	// Apply amendments
	result, err := amenderService.Amend(ctx, string(specContent), string(changesContent), outputFormat, dryRun)
	if err != nil {
		log.Error("Amendment failed", "error", err)
		return fmt.Errorf("failed to amend OpenAPI spec: %w", err)
	}

	// Handle dry run
	if dryRun {
		fmt.Printf("Dry run - Changes that would be applied:\n\n")
		for i, change := range result.Changes {
			fmt.Printf("%d. %s\n", i+1, change)
		}
		if len(result.Conflicts) > 0 {
			fmt.Printf("\nConflicts that need resolution:\n\n")
			for i, conflict := range result.Conflicts {
				fmt.Printf("%d. %s\n", i+1, conflict)
			}
		}
		return nil
	}

	// Determine output file
	if outputFile == "" {
		outputFile = specFile // Overwrite original
	}

	// Write result
	if err := os.WriteFile(outputFile, []byte(result.Content), 0600); err != nil {
		return fmt.Errorf("failed to write output file %s: %w", outputFile, err)
	}

	log.Info("OpenAPI specification amended successfully", "output_file", outputFile)

	// Print summary
	if verbose {
		fmt.Fprintf(os.Stderr, "\nAmendment Summary:\n")
		fmt.Fprintf(os.Stderr, "  Changes applied: %d\n", len(result.Changes))
		fmt.Fprintf(os.Stderr, "  Processing time: %dms\n", result.Metadata.ProcessingTimeMs)
		if len(result.Warnings) > 0 {
			fmt.Fprintf(os.Stderr, "  Warnings: %d\n", len(result.Warnings))
		}
		if len(result.Conflicts) > 0 {
			fmt.Fprintf(os.Stderr, "  Conflicts resolved: %d\n", len(result.Conflicts))
		}
	}

	return nil
}

func detectFormat(filename string) string {
	if len(filename) > 5 && filename[len(filename)-5:] == ".json" {
		return "json"
	}
	return "yaml" // Default to YAML
}
