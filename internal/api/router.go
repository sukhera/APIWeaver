package api

import (
	"log/slog"
	"net/http"

	"github.com/sukhera/APIWeaver/internal/api/handlers"
	"github.com/sukhera/APIWeaver/internal/api/middleware"
	"github.com/sukhera/APIWeaver/internal/config"
	"github.com/sukhera/APIWeaver/internal/storage"
)

// Router handles HTTP routing and middleware
type Router struct {
	config   *config.ExtendedConfig
	logger   *slog.Logger
	storage  storage.Storage
	handlers *handlers.Handlers
	mux      *http.ServeMux
}

// NewRouter creates a new router instance
func NewRouter(cfg *config.ExtendedConfig, logger *slog.Logger, store storage.Storage) *Router {
	// Create handlers
	h := handlers.New(cfg, logger, store)

	router := &Router{
		config:   cfg,
		logger:   logger,
		storage:  store,
		handlers: h,
		mux:      http.NewServeMux(),
	}

	router.setupRoutes()
	return router
}

// Handler returns the HTTP handler with middleware applied
func (r *Router) Handler() http.Handler {
	handler := http.Handler(r.mux)

	// Apply middleware stack (in reverse order - last applied executes first)
	handler = middleware.Recovery(r.logger)(handler)
	handler = middleware.Logging(r.logger)(handler)
	handler = middleware.CORS(r.config.Server.CORS)(handler)
	handler = middleware.Security()(handler)

	return handler
}

// setupRoutes configures all HTTP routes
func (r *Router) setupRoutes() {
	// Health and info endpoints
	r.mux.HandleFunc("GET /api/v1/health", r.handlers.Health)
	r.mux.HandleFunc("GET /api/v1/version", r.handlers.Version)

	// Core conversion endpoints
	r.mux.HandleFunc("POST /api/v1/generate", r.handlers.Generate)
	r.mux.HandleFunc("POST /api/v1/amend", r.handlers.Amend)
	r.mux.HandleFunc("POST /api/v1/validate", r.handlers.Validate)

	// Utility endpoints
	r.mux.HandleFunc("GET /api/v1/examples", r.handlers.Examples)

	// Static files (embedded web UI) - placeholder
	r.mux.HandleFunc("GET /", r.handlers.StaticFiles)
}
