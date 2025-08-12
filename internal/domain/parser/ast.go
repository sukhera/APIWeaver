package parser

import (
	"time"

	"github.com/sukhera/APIWeaver/pkg/errors"
)

// Document represents the root of a parsed Markdown API specification
type Document struct {
	Frontmatter *Frontmatter         `json:"frontmatter,omitempty"`
	Endpoints   []*Endpoint          `json:"endpoints"`
	Components  []*Component         `json:"components,omitempty"`
	ParsedAt    time.Time            `json:"parsed_at"`
	Errors      []*errors.ParseError `json:"errors,omitempty"`
}

// Frontmatter represents the optional YAML frontmatter at the beginning of the document
type Frontmatter struct {
	Title       string            `json:"title,omitempty"`
	Version     string            `json:"version,omitempty"`
	Description string            `json:"description,omitempty"`
	Servers     []Server          `json:"servers,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	LineNumber  int               `json:"line_number"`
}

// Server represents server configuration from frontmatter
type Server struct {
	URL         string `json:"url"`
	Description string `json:"description,omitempty"`
}

// Endpoint represents a parsed API endpoint
type Endpoint struct {
	Method      string       `json:"method"`
	Path        string       `json:"path"`
	Summary     string       `json:"summary,omitempty"`
	Description string       `json:"description,omitempty"`
	Parameters  []*Parameter `json:"parameters,omitempty"`
	RequestBody *RequestBody `json:"request_body,omitempty"`
	Responses   []*Response  `json:"responses,omitempty"`
	Tags        []string     `json:"tags,omitempty"`
	LineNumber  int          `json:"line_number"`
}

// Parameter represents a request parameter
type Parameter struct {
	Name        string      `json:"name"`
	In          string      `json:"in"` // "query", "path", "header", "cookie"
	Type        string      `json:"type"`
	Required    bool        `json:"required"`
	Description string      `json:"description,omitempty"`
	Example     interface{} `json:"example,omitempty"`
	Schema      *Schema     `json:"schema,omitempty"`
	LineNumber  int         `json:"line_number"`
}

// RequestBody represents the request body specification
type RequestBody struct {
	Description string             `json:"description,omitempty"`
	Required    bool               `json:"required"`
	Content     map[string]*Schema `json:"content"` // key: media type
	LineNumber  int                `json:"line_number"`
}

// Response represents an API response
type Response struct {
	StatusCode  string             `json:"status_code"`
	Description string             `json:"description,omitempty"`
	Headers     map[string]*Header `json:"headers,omitempty"`
	Content     map[string]*Schema `json:"content,omitempty"` // key: media type
	LineNumber  int                `json:"line_number"`
}

// Header represents a response header
type Header struct {
	Type        string      `json:"type"`
	Description string      `json:"description,omitempty"`
	Example     interface{} `json:"example,omitempty"`
}

// Schema represents a JSON/YAML schema definition
type Schema struct {
	Type        string             `json:"type,omitempty"`
	Format      string             `json:"format,omitempty"`
	Properties  map[string]*Schema `json:"properties,omitempty"`
	Items       *Schema            `json:"items,omitempty"`
	Required    []string           `json:"required,omitempty"`
	Enum        []interface{}      `json:"enum,omitempty"`
	Example     interface{}        `json:"example,omitempty"`
	Description string             `json:"description,omitempty"`
	Ref         string             `json:"$ref,omitempty"`
	AllOf       []*Schema          `json:"allOf,omitempty"`
	OneOf       []*Schema          `json:"oneOf,omitempty"`
	AnyOf       []*Schema          `json:"anyOf,omitempty"`
	LineNumber  int                `json:"line_number"`
}

// Component represents a reusable component definition
type Component struct {
	Name       string  `json:"name"`
	Type       string  `json:"type"` // "schema", "parameter", "response", etc.
	Schema     *Schema `json:"schema,omitempty"`
	LineNumber int     `json:"line_number"`
}
