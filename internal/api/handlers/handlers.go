package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"runtime"
	"time"

	"github.com/sukhera/APIWeaver/internal/api/models"
	"github.com/sukhera/APIWeaver/internal/config"
	"github.com/sukhera/APIWeaver/internal/services"
	"github.com/sukhera/APIWeaver/internal/storage"
)

// Handlers contains all HTTP request handlers
type Handlers struct {
	config    *config.ExtendedConfig
	logger    *slog.Logger
	storage   storage.Storage
	generator *services.Generator
	amender   *services.Amender
	validator *services.Validator
}

// New creates a new handlers instance
func New(cfg *config.ExtendedConfig, logger *slog.Logger, store storage.Storage) *Handlers {
	return &Handlers{
		config:    cfg,
		logger:    logger,
		storage:   store,
		generator: services.NewGenerator(cfg, logger),
		amender:   services.NewAmender(cfg, logger),
		validator: services.NewValidator(cfg, logger),
	}
}

// Health handles GET /api/v1/health
func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	response := models.HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Version:   "dev", // TODO: Get from build info
		System: models.SystemInfo{
			GoVersion: runtime.Version(),
			OS:        runtime.GOOS,
			Arch:      runtime.GOARCH,
		},
	}

	// Check storage health if available
	if h.storage != nil {
		if err := h.storage.Health(r.Context()); err != nil {
			response.Status = "degraded"
			h.logger.Warn("Storage health check failed", "error", err)
		}
	}

	h.writeJSONResponse(w, http.StatusOK, response)
}

// Version handles GET /api/v1/version
func (h *Handlers) Version(w http.ResponseWriter, r *http.Request) {
	response := models.VersionResponse{
		Version:   "dev", // TODO: Get from build info
		CommitSHA: "unknown",
		BuildTime: "unknown",
		GoVersion: runtime.Version(),
	}

	h.writeJSONResponse(w, http.StatusOK, response)
}

// Generate handles POST /api/v1/generate
func (h *Handlers) Generate(w http.ResponseWriter, r *http.Request) {
	req, err := models.ParseGenerateRequest(r)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	h.logger.Info("Processing generate request",
		"content_length", len(req.Content),
		"format", req.Format,
	)

	// Generate OpenAPI spec
	result, err := h.generator.Generate(r.Context(), req.Content, req.Format)
	if err != nil {
		h.logger.Error("Generation failed", "error", err)
		h.writeErrorResponse(w, http.StatusBadRequest, "Generation failed", err.Error())
		return
	}

	response := models.GenerateResponse{
		Success: true,
		Data: models.GenerateData{
			OpenAPI:  result.Content,
			Format:   result.Format,
			Metadata: result.Metadata,
		},
		Errors:    result.Errors,
		Warnings:  result.Warnings,
		Timestamp: time.Now(),
	}

	h.writeJSONResponse(w, http.StatusOK, response)
}

// Amend handles POST /api/v1/amend
func (h *Handlers) Amend(w http.ResponseWriter, r *http.Request) {
	req, err := models.ParseAmendRequest(r)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	h.logger.Info("Processing amend request",
		"spec_length", len(req.ExistingSpec),
		"changes_length", len(req.Changes),
		"format", req.Format,
	)

	// Apply amendments
	result, err := h.amender.Amend(r.Context(), req.ExistingSpec, req.Changes, req.Format, false)
	if err != nil {
		h.logger.Error("Amendment failed", "error", err)
		h.writeErrorResponse(w, http.StatusBadRequest, "Amendment failed", err.Error())
		return
	}

	response := models.AmendResponse{
		Success: true,
		Data: models.AmendData{
			OpenAPI:   result.Content,
			Format:    result.Format,
			Changes:   result.Changes,
			Conflicts: result.Conflicts,
			Metadata:  result.Metadata,
		},
		Errors:    result.Errors,
		Warnings:  result.Warnings,
		Timestamp: time.Now(),
	}

	h.writeJSONResponse(w, http.StatusOK, response)
}

// Validate handles POST /api/v1/validate
func (h *Handlers) Validate(w http.ResponseWriter, r *http.Request) {
	req, err := models.ParseValidateRequest(r)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	h.logger.Info("Processing validate request",
		"content_length", len(req.Content),
		"type", req.Type,
	)

	// Validate content
	result, err := h.validator.Validate(r.Context(), req.Content, req.Type)
	if err != nil {
		h.logger.Error("Validation failed", "error", err)
		h.writeErrorResponse(w, http.StatusBadRequest, "Validation failed", err.Error())
		return
	}

	response := models.ValidateResponse{
		Success: result.Valid,
		Data: models.ValidateData{
			Valid:        result.Valid,
			ErrorCount:   len(result.Errors),
			WarningCount: len(result.Warnings),
			Metadata:     result.Metadata,
		},
		Errors:    result.Errors,
		Warnings:  result.Warnings,
		Timestamp: time.Now(),
	}

	h.writeJSONResponse(w, http.StatusOK, response)
}

// Examples handles GET /api/v1/examples
func (h *Handlers) Examples(w http.ResponseWriter, r *http.Request) {
	// Get examples from storage if available
	var examples []models.ExampleTemplate

	if h.storage != nil {
		storageExamples, err := h.storage.ListExamples(r.Context(), storage.ExampleFilters{})
		if err != nil {
			h.logger.Warn("Failed to get examples from storage", "error", err)
		} else {
			// Convert storage examples to API models
			for _, ex := range storageExamples {
				examples = append(examples, models.ExampleTemplate{
					ID:          ex.ID,
					Name:        ex.Name,
					Description: ex.Description,
					Content:     ex.Content,
					Category:    ex.Category,
					Tags:        ex.Tags,
				})
			}
		}
	}

	// Fallback to hardcoded examples if no storage
	if len(examples) == 0 {
		examples = models.GetDefaultExamples()
	}

	response := models.ExamplesResponse{
		Success:   true,
		Examples:  examples,
		Timestamp: time.Now(),
	}

	h.writeJSONResponse(w, http.StatusOK, response)
}

// StaticFiles handles static file serving (placeholder for embedded web UI)
func (h *Handlers) StaticFiles(w http.ResponseWriter, r *http.Request) {
	// For MVP, return a simple HTML page
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`<!DOCTYPE html>
<html>
<head>
    <title>APIWeaver</title>
</head>
<body>
    <h1>APIWeaver API Server</h1>
    <p>The APIWeaver backend is running. Use the API endpoints:</p>
    <ul>
        <li><a href="/api/v1/health">Health Check</a></li>
        <li><a href="/api/v1/version">Version Info</a></li>
        <li><a href="/api/v1/examples">Example Templates</a></li>
    </ul>
</body>
</html>`))
}

// Helper methods

func (h *Handlers) writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("Failed to encode JSON response", "error", err)
	}
}

func (h *Handlers) writeErrorResponse(w http.ResponseWriter, statusCode int, message, details string) {
	response := models.ErrorResponse{
		Success: false,
		Error: models.ErrorDetails{
			Message: message,
			Details: details,
			Code:    statusCode,
		},
		Timestamp: time.Now(),
	}

	h.writeJSONResponse(w, statusCode, response)
}
