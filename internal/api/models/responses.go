package models

import (
	"time"

	"github.com/sukhera/APIWeaver/internal/services"
)

// Base response structure
type BaseResponse struct {
	Success   bool      `json:"success"`
	Timestamp time.Time `json:"timestamp"`
}

// Error response
type ErrorResponse struct {
	Success   bool         `json:"success"`
	Error     ErrorDetails `json:"error"`
	Timestamp time.Time    `json:"timestamp"`
}

type ErrorDetails struct {
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
	Code    int    `json:"code"`
}

// Health check response
type HealthResponse struct {
	Status    string     `json:"status"`
	Timestamp time.Time  `json:"timestamp"`
	Version   string     `json:"version"`
	System    SystemInfo `json:"system"`
}

type SystemInfo struct {
	GoVersion string `json:"go_version"`
	OS        string `json:"os"`
	Arch      string `json:"arch"`
}

// Version response
type VersionResponse struct {
	Version   string `json:"version"`
	CommitSHA string `json:"commit_sha"`
	BuildTime string `json:"build_time"`
	GoVersion string `json:"go_version"`
}

// Generate response
type GenerateResponse struct {
	Success   bool         `json:"success"`
	Data      GenerateData `json:"data,omitempty"`
	Errors    []string     `json:"errors,omitempty"`
	Warnings  []string     `json:"warnings,omitempty"`
	Timestamp time.Time    `json:"timestamp"`
}

type GenerateData struct {
	OpenAPI  string                        `json:"openapi"`
	Format   string                        `json:"format"`
	Metadata services.GenerationMetadata   `json:"metadata"`
}

// Amend response
type AmendResponse struct {
	Success   bool      `json:"success"`
	Data      AmendData `json:"data,omitempty"`
	Errors    []string  `json:"errors,omitempty"`
	Warnings  []string  `json:"warnings,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

type AmendData struct {
	OpenAPI   string                      `json:"openapi"`
	Format    string                      `json:"format"`
	Changes   []string                    `json:"changes"`
	Conflicts []string                    `json:"conflicts,omitempty"`
	Metadata  services.AmendmentMetadata  `json:"metadata"`
}

// Validate response
type ValidateResponse struct {
	Success   bool         `json:"success"`
	Data      ValidateData `json:"data,omitempty"`
	Errors    []string     `json:"errors,omitempty"`
	Warnings  []string     `json:"warnings,omitempty"`
	Timestamp time.Time    `json:"timestamp"`
}

type ValidateData struct {
	Valid        bool                        `json:"valid"`
	ErrorCount   int                         `json:"error_count"`
	WarningCount int                         `json:"warning_count"`
	Metadata     services.ValidationMetadata `json:"metadata"`
}

// Examples response
type ExamplesResponse struct {
	Success   bool              `json:"success"`
	Examples  []ExampleTemplate `json:"examples"`
	Timestamp time.Time         `json:"timestamp"`
}

type ExampleTemplate struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Content     string   `json:"content"`
	Category    string   `json:"category"`
	Tags        []string `json:"tags"`
}

// GetDefaultExamples returns default example templates
func GetDefaultExamples() []ExampleTemplate {
	return []ExampleTemplate{
		{
			ID:          "simple-api",
			Name:        "Simple API",
			Description: "A basic REST API with CRUD operations",
			Category:    "basic",
			Tags:        []string{"rest", "crud"},
			Content: `---
title: "Simple Task API"
version: "1.0.0"
description: "A simple task management API"
servers:
  - url: "https://api.example.com/v1"
---

# Simple Task API

## GET /tasks
Retrieve all tasks.

**Parameters:**
| Name | Type | Required | Description |
|------|------|----------|-------------|
| status | string | No | Filter by task status (pending, completed) |
| limit | integer | No | Number of tasks to return (default: 10) |

**Response (200):**
` + "```json\n" + `{
  "tasks": [
    {
      "id": "1",
      "title": "Example task",
      "completed": false
    }
  ]
}
` + "```" + `

## POST /tasks
Create a new task.

**Request Body:**
` + "```json\n" + `{
  "title": "string",
  "description": "string"
}
` + "```",
		},
		{
			ID:          "auth-api",
			Name:        "Authentication API",
			Description: "API with authentication endpoints",
			Category:    "auth",
			Tags:        []string{"auth", "security"},
			Content: `---
title: "Auth API"
version: "1.0.0"
description: "Authentication and authorization API"
---

# Authentication API

## POST /auth/login
User login endpoint.

**Request Body:**
` + "```json\n" + `{
  "email": "user@example.com",
  "password": "password123"
}
` + "```" + `

**Response (200):**
` + "```json\n" + `{
  "token": "jwt_token_here",
  "expires_in": 3600
}
` + "```",
		},
	}
}