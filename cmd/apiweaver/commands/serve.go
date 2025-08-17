package commands

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/sukhera/APIWeaver/internal/api"
	"github.com/sukhera/APIWeaver/internal/config"
	"github.com/sukhera/APIWeaver/internal/logger"
	"github.com/sukhera/APIWeaver/internal/storage"
	"github.com/sukhera/APIWeaver/internal/storage/mongodb"
)

// NewServeCmd creates the serve command
func NewServeCmd() *cobra.Command {
	var (
		port       int
		host       string
		configFile string
		verbose    bool
		devMode    bool
	)

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the HTTP API server",
		Long: `Start the APIWeaver HTTP API server to provide web-based access
to the markdown parsing and OpenAPI generation functionality.

The server provides both a REST API and serves the embedded web UI.`,
		Example: `  apiweaver serve
  apiweaver serve --port 8080 --host 0.0.0.0
  apiweaver serve --config server.yaml --verbose --dev`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runServe(cmd.Context(), port, host, configFile, verbose, devMode)
		},
	}

	cmd.Flags().IntVarP(&port, "port", "p", 8080, "Server port")
	cmd.Flags().StringVarP(&host, "host", "H", "0.0.0.0", "Server host")
	cmd.Flags().StringVarP(&configFile, "config", "c", "", "Configuration file path")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
	cmd.Flags().BoolVar(&devMode, "dev", false, "Enable development mode")

	return cmd
}

func runServe(ctx context.Context, port int, host, configFile string, verbose, devMode bool) error {
	// Setup context with cancellation for graceful shutdown
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Load configuration with fallback to defaults
	var cfg *config.ExtendedConfig
	var err error

	if configFile != "" {
		cfg, err = config.Load(configFile)
		if err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}
	} else {
		// Use defaults when no config file specified
		cfg = &config.ExtendedConfig{
			Config: config.Default(),
			Server: config.ServerConfig{
				Port:    8080,
				Host:    "0.0.0.0",
				DevMode: false,
				CORS: config.CORSConfig{
					Enabled:        true,
					AllowedOrigins: []string{"*"},
					AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
					AllowedHeaders: []string{"Content-Type", "Authorization"},
				},
			},
			Logger: logger.Config{
				Level:  "info",
				Format: "json",
				Output: "stdout",
			},
			Storage: config.StorageConfig{
				MongoDB: config.MongoDBConfig{
					Enabled:  false,
					URI:      "mongodb://localhost:27017",
					Database: "apiweaver",
				},
			},
		}
	}

	// Override with command line flags
	if verbose {
		cfg.Verbose = true
	}
	cfg.Server.Port = port
	cfg.Server.Host = host
	cfg.Server.DevMode = devMode

	// Setup logger
	log, err := logger.New(cfg.Logger)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}

	log.Info("Starting APIWeaver server",
		"version", "dev", // TODO: Get from build info
		"port", port,
		"host", host,
		"dev_mode", devMode,
	)

	// Initialize storage (if configured)
	var store storage.Storage
	if cfg.Storage.MongoDB.Enabled {
		log.Info("Initializing MongoDB storage", "uri", cfg.Storage.MongoDB.URI)
		store, err = mongodb.NewMongoDB(cfg.Storage.MongoDB)
		if err != nil {
			log.Warn("Failed to initialize MongoDB storage, continuing without persistence", "error", err)
		} else {
			defer func() {
				if err := store.Close(); err != nil {
					log.Error("Failed to close storage", "error", err)
				}
			}()
		}
	}

	// Create and start HTTP server
	server, err := api.NewServer(cfg, log, store)
	if err != nil {
		return fmt.Errorf("failed to create server: %w", err)
	}

	// Start server in a goroutine
	serverErrChan := make(chan error, 1)
	go func() {
		if err := server.Start(ctx); err != nil {
			serverErrChan <- err
		}
	}()

	// Wait for shutdown signal or server error
	select {
	case sig := <-sigChan:
		log.Info("Received shutdown signal", "signal", sig)
		cancel()
		return server.Shutdown(context.Background())
	case err := <-serverErrChan:
		log.Error("Server error", "error", err)
		return err
	}
}
