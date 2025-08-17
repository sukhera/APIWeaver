package models

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// GenerateRequest represents a request to generate OpenAPI spec
type GenerateRequest struct {
	Content string `json:"content"`
	Format  string `json:"format"` // "yaml" or "json"
}

// AmendRequest represents a request to amend OpenAPI spec
type AmendRequest struct {
	ExistingSpec string `json:"existing_spec"`
	Changes      string `json:"changes"`
	Format       string `json:"format"`
}

// ValidateRequest represents a request to validate content
type ValidateRequest struct {
	Content string `json:"content"`
	Type    string `json:"type"` // "markdown" or "openapi"
}

// ParseGenerateRequest parses a generate request from HTTP request
func ParseGenerateRequest(r *http.Request) (*GenerateRequest, error) {
	contentType := r.Header.Get("Content-Type")
	
	if strings.HasPrefix(contentType, "multipart/form-data") {
		return parseMultipartGenerateRequest(r)
	} else if strings.HasPrefix(contentType, "application/json") {
		return parseJSONGenerateRequest(r)
	}
	
	return nil, fmt.Errorf("unsupported content type: %s", contentType)
}

// ParseAmendRequest parses an amend request from HTTP request
func ParseAmendRequest(r *http.Request) (*AmendRequest, error) {
	var req AmendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("failed to decode JSON request: %w", err)
	}

	if req.Format == "" {
		req.Format = "yaml" // default
	}

	if req.ExistingSpec == "" {
		return nil, fmt.Errorf("existing_spec is required")
	}

	if req.Changes == "" {
		return nil, fmt.Errorf("changes is required")
	}

	return &req, nil
}

// ParseValidateRequest parses a validate request from HTTP request
func ParseValidateRequest(r *http.Request) (*ValidateRequest, error) {
	contentType := r.Header.Get("Content-Type")
	
	if strings.HasPrefix(contentType, "multipart/form-data") {
		return parseMultipartValidateRequest(r)
	} else if strings.HasPrefix(contentType, "application/json") {
		return parseJSONValidateRequest(r)
	}
	
	return nil, fmt.Errorf("unsupported content type: %s", contentType)
}

// Helper functions

func parseJSONGenerateRequest(r *http.Request) (*GenerateRequest, error) {
	var req GenerateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("failed to decode JSON request: %w", err)
	}

	if req.Format == "" {
		req.Format = "yaml" // default
	}

	return &req, nil
}

func parseMultipartGenerateRequest(r *http.Request) (*GenerateRequest, error) {
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10MB max
		return nil, fmt.Errorf("failed to parse multipart form: %w", err)
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		return nil, fmt.Errorf("failed to get file from form: %w", err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file content: %w", err)
	}

	format := r.FormValue("format")
	if format == "" {
		format = "yaml" // default
	}

	return &GenerateRequest{
		Content: string(content),
		Format:  format,
	}, nil
}

func parseJSONValidateRequest(r *http.Request) (*ValidateRequest, error) {
	var req ValidateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("failed to decode JSON request: %w", err)
	}

	if req.Type == "" {
		req.Type = "markdown" // default
	}

	if req.Content == "" {
		return nil, fmt.Errorf("content is required")
	}

	return &req, nil
}

func parseMultipartValidateRequest(r *http.Request) (*ValidateRequest, error) {
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10MB max
		return nil, fmt.Errorf("failed to parse multipart form: %w", err)
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		return nil, fmt.Errorf("failed to get file from form: %w", err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file content: %w", err)
	}

	validateType := r.FormValue("type")
	if validateType == "" {
		validateType = "markdown" // default
	}

	return &ValidateRequest{
		Content: string(content),
		Type:    validateType,
	}, nil
}