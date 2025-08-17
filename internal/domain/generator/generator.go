package generator

import (
	"context"
	"fmt"

	"github.com/sukhera/APIWeaver/internal/domain/parser"
)

// Config holds generator configuration
type Config struct {
	Format          string
	PrettyPrint     bool
	IncludeExamples bool
	ValidateOutput  bool
	StrictMode      bool
}

// Generator generates OpenAPI specifications from parsed documents
type Generator struct {
	config Config
}

// New creates a new Generator instance
func New(config Config) *Generator {
	return &Generator{
		config: config,
	}
}

// Generate generates an OpenAPI specification from a parsed document
func (g *Generator) Generate(ctx context.Context, doc *parser.Document, format string) (string, error) {
	if doc == nil {
		return "", fmt.Errorf("document is nil")
	}

	// For MVP, return a mock OpenAPI spec based on the document
	switch format {
	case "json":
		return g.generateJSON(ctx, doc)
	case "yaml":
		return g.generateYAML(ctx, doc)
	default:
		return g.generateYAML(ctx, doc) // Default to YAML
	}
}

// generateYAML generates YAML format OpenAPI spec
func (g *Generator) generateYAML(ctx context.Context, doc *parser.Document) (string, error) {
	// Mock implementation - in real implementation this would use the AST
	spec := `openapi: 3.1.0
info:
  title: Generated API
  version: 1.0.0
  description: API generated from markdown`

	if doc.Frontmatter != nil {
		if doc.Frontmatter.Title != "" {
			spec = `openapi: 3.1.0
info:
  title: ` + doc.Frontmatter.Title + `
  version: ` + getVersionOrDefault(doc.Frontmatter.Version) + `
  description: ` + getDescriptionOrDefault(doc.Frontmatter.Description)
		}
	}

	spec += `
paths:`

	// Add endpoints
	if len(doc.Endpoints) == 0 {
		spec += `
  /example:
    get:
      summary: Example endpoint
      responses:
        '200':
          description: Success
          content:
            application/json:
              schema:
                type: object
                properties:
                  message:
                    type: string
                    example: "Hello, World!"`
	} else {
		for _, endpoint := range doc.Endpoints {
			spec += fmt.Sprintf(`
  %s:
    %s:
      summary: %s
      responses:
        '200':
          description: Success`, 
				endpoint.Path, 
				endpoint.Method, 
				getEndpointSummary(endpoint))
		}
	}

	spec += `
components:
  schemas:
    Error:
      type: object
      properties:
        message:
          type: string
        code:
          type: integer`

	return spec, nil
}

// generateJSON generates JSON format OpenAPI spec
func (g *Generator) generateJSON(ctx context.Context, doc *parser.Document) (string, error) {
	// Mock implementation - in real implementation this would build proper JSON
	return `{
  "openapi": "3.1.0",
  "info": {
    "title": "Generated API",
    "version": "1.0.0",
    "description": "API generated from markdown"
  },
  "paths": {
    "/example": {
      "get": {
        "summary": "Example endpoint",
        "responses": {
          "200": {
            "description": "Success",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "message": {
                      "type": "string",
                      "example": "Hello, World!"
                    }
                  }
                }
              }
            }
          }
        }
      }
    }
  },
  "components": {
    "schemas": {
      "Error": {
        "type": "object",
        "properties": {
          "message": {
            "type": "string"
          },
          "code": {
            "type": "integer"
          }
        }
      }
    }
  }
}`, nil
}

// Helper functions
func getVersionOrDefault(version string) string {
	if version == "" {
		return "1.0.0"
	}
	return version
}

func getDescriptionOrDefault(description string) string {
	if description == "" {
		return "API generated from markdown"
	}
	return description
}

func getEndpointSummary(endpoint *parser.Endpoint) string {
	if endpoint.Summary != "" {
		return endpoint.Summary
	}
	if endpoint.Description != "" {
		return endpoint.Description
	}
	return fmt.Sprintf("%s %s", endpoint.Method, endpoint.Path)
}