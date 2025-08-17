package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/sukhera/APIWeaver/internal/config"
	"github.com/sukhera/APIWeaver/internal/storage"
)

// Server represents the HTTP API server
type Server struct {
	config  *config.ExtendedConfig
	logger  *slog.Logger
	storage storage.Storage
	server  *http.Server
	router  *Router
}

// NewServer creates a new API server instance
func NewServer(cfg *config.ExtendedConfig, logger *slog.Logger, store storage.Storage) (*Server, error) {
	// Create router with dependencies
	router := NewRouter(cfg, logger, store)

	return &Server{
		config:  cfg,
		logger:  logger,
		storage: store,
		router:  router,
	}, nil
}

// Start starts the HTTP server
func (s *Server) Start(ctx context.Context) error {
	addr := fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port)
	
	s.server = &http.Server{
		Addr:         addr,
		Handler:      s.router.Handler(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	s.logger.Info("Starting HTTP server", "addr", addr)

	// Start server in a goroutine
	errChan := make(chan error, 1)
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	// Wait for context cancellation or server error
	select {
	case err := <-errChan:
		return fmt.Errorf("server failed to start: %w", err)
	case <-ctx.Done():
		return s.Shutdown(context.Background())
	}
}

// Shutdown gracefully shuts down the HTTP server
func (s *Server) Shutdown(ctx context.Context) error {
	if s.server == nil {
		return nil
	}

	s.logger.Info("Shutting down HTTP server")
	
	// Set a timeout for graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if err := s.server.Shutdown(shutdownCtx); err != nil {
		s.logger.Error("Error during server shutdown", "error", err)
		return err
	}

	s.logger.Info("HTTP server shutdown complete")
	return nil
}