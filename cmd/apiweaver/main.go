package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/sukhera/APIWeaver/internal/config"
	"github.com/sukhera/APIWeaver/internal/domain/parser"
	"github.com/sukhera/APIWeaver/pkg/errors"
)

// Version information - should be set during build
var (
	version   = "dev"
	commitSHA = "unknown"
	buildTime = "unknown"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "apiweaver",
	Short: "APIWeaver - Markdown to API Specification Parser",
	Long: `APIWeaver is a powerful tool for parsing markdown files and converting them 
into structured API specifications. It supports various output formats and 
provides comprehensive validation and error reporting.`,
	Version: fmt.Sprintf("%s (commit: %s, built: %s)", version, commitSHA, buildTime),
}

// parseCmd represents the parse command
var parseCmd = &cobra.Command{
	Use:   "parse [input-file]",
	Short: "Parse a markdown file into API specification",
	Long: `Parse a markdown file and convert it into a structured API specification.
The input file should contain API documentation in markdown format with
endpoints, parameters, and response definitions.`,
	Args: cobra.ExactArgs(1),
	Example: `  apiweaver parse example.md
  apiweaver parse api-docs.md --output result.json --strict
  apiweaver parse docs.md --config config.yaml --verbose`,
	RunE: runParse,
}

var (
	// Global flags
	configFile string
	verbose    bool
	strict     bool

	// Parse command flags
	outputFile   string
	outputFormat string
)

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "Configuration file path")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
	rootCmd.PersistentFlags().BoolVarP(&strict, "strict", "s", false, "Enable strict parsing mode")

	// Parse command flags
	parseCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file for parsed results")
	parseCmd.Flags().StringVarP(&outputFormat, "format", "f", "json", "Output format (json, yaml, text)")

	// Add parse command to root
	rootCmd.AddCommand(parseCmd)
}

// runParse executes the parse command
func runParse(cmd *cobra.Command, args []string) error {
	inputFile := args[0]

	// Load configuration
	cfg, err := loadConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Setup context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Run the parsing operation
	if err := run(ctx, cfg, inputFile, outputFile); err != nil {
		return fmt.Errorf("parsing failed: %w", err)
	}

	return nil
}

// loadConfig loads and validates the application configuration
func loadConfig() (*config.Config, error) {
	// Create a new Viper instance
	v := config.NewViperConfig()

	// Set defaults
	v.SetDefault("strict_mode", strict)
	v.SetDefault("verbose", verbose)
	v.SetDefault("output_format", outputFormat)

	// Read config file if specified
	if configFile != "" {
		v.SetConfigFile(configFile)
		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// Override with command line flags
	v.Set("strict_mode", strict)
	v.Set("verbose", verbose)
	v.Set("output_format", outputFormat)

	// Create config from Viper
	cfg := config.FromViper(v)

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// readFileWithContext reads a file with context support
func readFileWithContext(ctx context.Context, filename string) ([]byte, error) {
	// Validate file path for security
	if err := validateFilePath(filename); err != nil {
		return nil, fmt.Errorf("invalid file path: %w", err)
	}

	// Create a channel for the result
	resultChan := make(chan []byte, 1)
	errorChan := make(chan error, 1)

	go func() {
		// #nosec G304 -- File path is validated by validateFilePath function
		content, err := os.ReadFile(filename)
		if err != nil {
			errorChan <- err
		} else {
			resultChan <- content
		}
	}()

	select {
	case content := <-resultChan:
		return content, nil
	case err := <-errorChan:
		return nil, err
	case <-ctx.Done():
		return nil, errors.NewTimeoutError("file reading", "context cancelled")
	}
}

// parseWithContext parses content with context support
func parseWithContext(ctx context.Context, content string, cfg *config.Config) (*parser.Document, error) {
	// Create parser with configuration using functional options
	p := parser.New(
		parser.WithStrictMode(cfg.StrictMode),
		parser.WithRecovery(cfg.EnableRecovery, cfg.MaxRecoveryAttempts),
		parser.WithTimeout(cfg.ParserTimeout),
		parser.WithAllowedMethods(cfg.AllowedMethods),
		parser.WithValidationLevel(cfg.ValidationLevel),
		parser.WithRequireExamples(cfg.RequireExamples),
		parser.WithMaxNestingDepth(cfg.MaxNestingDepth),
		parser.WithInitialSliceCapacity(cfg.InitialSliceCapacity),
	)

	return p.ParseWithContext(ctx, content)
}

// writeOutputWithContext writes output with context support
func writeOutputWithContext(ctx context.Context, outputFile string, doc *parser.Document, cfg *config.Config) error {
	// Validate file path for security
	if err := validateFilePath(outputFile); err != nil {
		return fmt.Errorf("invalid output file path: %w", err)
	}

	// Create a channel for the result
	resultChan := make(chan error, 1)

	go func() {
		// #nosec G304 -- File path is validated by validateFilePath function
		file, err := os.Create(outputFile)
		if err != nil {
			resultChan <- err
			return
		}
		defer func() {
			if closeErr := file.Close(); closeErr != nil && err == nil {
				resultChan <- closeErr
			}
		}()

		// Write document summary based on output format
		switch cfg.OutputFormat {
		case "json":
			// Write JSON format
			encoder := json.NewEncoder(file)
			if cfg.PrettyPrint {
				encoder.SetIndent("", "  ")
			}
			result := map[string]interface{}{
				"metadata": map[string]interface{}{
					"parsed_at":  doc.ParsedAt.Format(time.RFC3339),
					"endpoints":  len(doc.Endpoints),
					"components": len(doc.Components),
					"errors":     len(doc.Errors),
				},
				"configuration": map[string]interface{}{
					"strict_mode":      cfg.StrictMode,
					"validation_level": cfg.ValidationLevel,
					"output_format":    cfg.OutputFormat,
				},
			}
			if err := encoder.Encode(result); err != nil {
				resultChan <- err
				return
			}
		default:
			// Write text format (default)
			lines := []string{
				"APIWeaver Parse Results",
				"======================",
				fmt.Sprintf("Parsed at: %s", doc.ParsedAt.Format(time.RFC3339)),
				fmt.Sprintf("Endpoints: %d", len(doc.Endpoints)),
				fmt.Sprintf("Components: %d", len(doc.Components)),
				fmt.Sprintf("Errors: %d", len(doc.Errors)),
			}

			// Add configuration details if verbose or pretty print is enabled
			if cfg.Verbose || cfg.PrettyPrint {
				lines = append(lines, "",
					"Configuration:",
					fmt.Sprintf("  Strict mode: %v", cfg.StrictMode),
					fmt.Sprintf("  Validation level: %s", cfg.ValidationLevel),
					fmt.Sprintf("  Output format: %s", cfg.OutputFormat),
				)
			}

			for _, line := range lines {
				if _, writeErr := fmt.Fprintln(file, line); writeErr != nil {
					resultChan <- writeErr
					return
				}
			}
		}

		resultChan <- nil
	}()

	select {
	case err := <-resultChan:
		return err
	case <-ctx.Done():
		return errors.NewTimeoutError("file writing", "context cancelled")
	}
}

// run executes the main application logic
func run(ctx context.Context, cfg *config.Config, inputFile, outputFile string) error {
	if cfg.Verbose {
		fmt.Printf("Starting APIWeaver v%s\n", version)
		fmt.Printf("Parsing file: %s\n", inputFile)
	}

	// Create a timeout context for the entire operation
	timeoutCtx, cancel := context.WithTimeout(ctx, cfg.ParserTimeout)
	defer cancel()

	// Read input file with context
	content, err := readFileWithContext(timeoutCtx, inputFile)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	if cfg.Verbose {
		fmt.Printf("File size: %d bytes\n", len(content))
		fmt.Printf("Configuration loaded successfully\n")
		fmt.Printf("Strict mode: %v\n", cfg.StrictMode)
		fmt.Printf("Verbose mode: %v\n", cfg.Verbose)
		fmt.Printf("Parser timeout: %v\n", cfg.ParserTimeout)
	}

	// Parse the content with context
	doc, err := parseWithContext(timeoutCtx, string(content), cfg)
	if err != nil {
		return fmt.Errorf("failed to parse content: %w", err)
	}

	if cfg.Verbose {
		fmt.Printf("Parsing completed successfully\n")
		fmt.Printf("Endpoints found: %d\n", len(doc.Endpoints))
		fmt.Printf("Components found: %d\n", len(doc.Components))
		fmt.Printf("Errors found: %d\n", len(doc.Errors))
	}

	// Output results with context
	if outputFile != "" {
		if err := writeOutputWithContext(timeoutCtx, outputFile, doc, cfg); err != nil {
			return fmt.Errorf("failed to write output: %w", err)
		}
		if cfg.Verbose {
			fmt.Printf("Results written to: %s\n", outputFile)
		}
	} else {
		// Print summary to stdout
		printSummary(doc, cfg)
	}

	return nil
}

// validateFilePath validates a file path for security
func validateFilePath(path string) error {
	// Check for path traversal attempts
	if strings.Contains(path, "..") {
		return fmt.Errorf("path traversal not allowed")
	}

	// Check for null bytes or other dangerous characters
	if strings.ContainsAny(path, "\x00") {
		return fmt.Errorf("null bytes not allowed in file path")
	}

	return nil
}

// printSummary prints a summary of the parsing results to stdout
func printSummary(doc *parser.Document, cfg *config.Config) {
	fmt.Printf("\n=== APIWeaver Parse Summary ===\n")
	fmt.Printf("Parsed at: %s\n", doc.ParsedAt.Format(time.RFC3339))
	fmt.Printf("Endpoints: %d\n", len(doc.Endpoints))
	fmt.Printf("Components: %d\n", len(doc.Components))
	fmt.Printf("Errors: %d\n", len(doc.Errors))

	if len(doc.Errors) > 0 {
		fmt.Printf("\nErrors found:\n")
		for i, err := range doc.Errors {
			fmt.Printf("  %d. %s\n", i+1, err.Error())
		}
	}

	if cfg.Verbose {
		fmt.Printf("\nConfiguration used:\n")
		fmt.Printf("  Strict mode: %v\n", cfg.StrictMode)
		fmt.Printf("  Parser timeout: %v\n", cfg.ParserTimeout)
		fmt.Printf("  Validation level: %s\n", cfg.ValidationLevel)
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		// Use our custom error types for better error reporting
		if cfgErr, ok := err.(*errors.ConfigError); ok {
			fmt.Fprintf(os.Stderr, "Configuration error: %v\n", cfgErr)
		} else if timeoutErr, ok := err.(*errors.TimeoutError); ok {
			fmt.Fprintf(os.Stderr, "Timeout error: %v\n", timeoutErr)
		} else {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		os.Exit(1)
	}
}
