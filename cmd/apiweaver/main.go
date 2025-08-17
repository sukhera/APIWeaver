package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/sukhera/APIWeaver/cmd/apiweaver/commands"
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
	Short: "APIWeaver - Markdown to API Specification Generator",
	Long: `APIWeaver is a powerful tool for converting markdown API documentation 
into structured OpenAPI 3.1 specifications. It supports generation, validation, 
and amendment of API specifications with comprehensive error reporting.`,
	Version: fmt.Sprintf("%s (commit: %s, built: %s)", version, commitSHA, buildTime),
}

func main() {
	// Add all commands
	rootCmd.AddCommand(commands.NewGenerateCmd())
	rootCmd.AddCommand(commands.NewAmendCmd())
	rootCmd.AddCommand(commands.NewValidateCmd())
	rootCmd.AddCommand(commands.NewServeCmd())

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}