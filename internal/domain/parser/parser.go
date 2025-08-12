package parser

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/sukhera/APIWeaver/pkg/errors"
)

// Parser represents the main markdown parser
type Parser struct {
	config *ParserConfig
}

// ParserConfig holds configuration for the parser
type ParserConfig struct {
	StrictMode           bool
	EnableRecovery       bool
	MaxRecoveryAttempts  int
	Timeout              time.Duration
	InitialSliceCapacity int
	AllowedMethods       []string
	ValidationLevel      string
	RequireExamples      bool
	MaxNestingDepth      int
}

// ParserOption is a functional option for configuring the parser
type ParserOption func(*ParserConfig)

// WithStrictMode enables strict parsing mode
func WithStrictMode(strict bool) ParserOption {
	return func(cfg *ParserConfig) {
		cfg.StrictMode = strict
	}
}

// WithRecovery enables error recovery with specified attempts
func WithRecovery(enable bool, maxAttempts int) ParserOption {
	return func(cfg *ParserConfig) {
		cfg.EnableRecovery = enable
		cfg.MaxRecoveryAttempts = maxAttempts
	}
}

// WithTimeout sets the parser timeout
func WithTimeout(timeout time.Duration) ParserOption {
	return func(cfg *ParserConfig) {
		cfg.Timeout = timeout
	}
}

// WithAllowedMethods sets the allowed HTTP methods
func WithAllowedMethods(methods []string) ParserOption {
	return func(cfg *ParserConfig) {
		cfg.AllowedMethods = methods
	}
}

// WithValidationLevel sets the validation level
func WithValidationLevel(level string) ParserOption {
	return func(cfg *ParserConfig) {
		cfg.ValidationLevel = level
	}
}

// WithRequireExamples requires examples in the specification
func WithRequireExamples(require bool) ParserOption {
	return func(cfg *ParserConfig) {
		cfg.RequireExamples = require
	}
}

// WithMaxNestingDepth sets the maximum schema nesting depth
func WithMaxNestingDepth(depth int) ParserOption {
	return func(cfg *ParserConfig) {
		cfg.MaxNestingDepth = depth
	}
}

// WithInitialSliceCapacity sets the initial slice capacity for better performance
func WithInitialSliceCapacity(capacity int) ParserOption {
	return func(cfg *ParserConfig) {
		cfg.InitialSliceCapacity = capacity
	}
}

// New creates a new parser with the given options
func New(options ...ParserOption) *Parser {
	config := defaultConfig()

	for _, option := range options {
		option(config)
	}

	return &Parser{
		config: config,
	}
}

// defaultConfig returns a default parser configuration
func defaultConfig() *ParserConfig {
	return &ParserConfig{
		StrictMode:           false,
		EnableRecovery:       true,
		MaxRecoveryAttempts:  3,
		Timeout:              30 * time.Second,
		InitialSliceCapacity: 100,
		AllowedMethods:       []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		ValidationLevel:      "basic",
		RequireExamples:      false,
		MaxNestingDepth:      10,
	}
}

// Parse parses markdown content and returns a Document
func (p *Parser) Parse(content string) (*Document, error) {
	// Create error collector for multiple errors
	collector := errors.NewErrorCollector(p.config.MaxRecoveryAttempts)

	// Create document
	doc := &Document{
		ParsedAt: time.Now(),
		Errors:   []*errors.ParseError{},
	}

	// Parse frontmatter
	frontmatter, remainingContent, err := p.parseFrontmatter(content)
	if err != nil {
		if parseErr, ok := err.(*errors.ParseError); ok {
			collector.Add(parseErr)
		} else {
			collector.Add(errors.NewError(errors.ErrorTypeFrontmatter, err.Error()).Build())
		}
	} else {
		doc.Frontmatter = frontmatter
	}

	// Parse endpoints
	endpoints, endpointErrors := p.parseEndpoints(remainingContent)
	doc.Endpoints = endpoints
	for _, err := range endpointErrors {
		collector.Add(err)
	}

	// Parse components
	components, componentErrors := p.parseComponents(remainingContent)
	doc.Components = components
	for _, err := range componentErrors {
		collector.Add(err)
	}

	// Validate document
	validationErrors := p.validateDocument(doc)
	for _, err := range validationErrors {
		collector.Add(err)
	}

	// Set errors from collector
	doc.Errors = collector.GetErrors()

	// Return error if in strict mode and there are errors
	if p.config.StrictMode && collector.HasErrors() {
		return doc, collector.ToError()
	}

	return doc, nil
}

// ParseWithContext parses content with a context for cancellation
func (p *Parser) ParseWithContext(ctx context.Context, content string) (*Document, error) {
	// Create a channel for the result
	resultChan := make(chan *Document, 1)
	errorChan := make(chan error, 1)

	go func() {
		doc, err := p.Parse(content)
		if err != nil {
			errorChan <- err
		} else {
			resultChan <- doc
		}
	}()

	select {
	case doc := <-resultChan:
		return doc, nil
	case err := <-errorChan:
		return nil, err
	case <-ctx.Done():
		return nil, errors.NewTimeoutError("parsing", p.config.Timeout.String())
	}
}

// parseFrontmatter parses YAML frontmatter from the content
func (p *Parser) parseFrontmatter(content string) (*Frontmatter, string, error) {
	// This is a placeholder implementation
	// In a real implementation, you would parse YAML frontmatter here
	return nil, content, nil
}

// parseEndpoints parses endpoints from the content
func (p *Parser) parseEndpoints(content string) ([]*Endpoint, []*errors.ParseError) {
	// This is a placeholder implementation
	// In a real implementation, you would parse endpoints here
	_ = content // TODO: Implement endpoint parsing from content
	return []*Endpoint{}, []*errors.ParseError{}
}

// parseComponents parses reusable components from the content
func (p *Parser) parseComponents(content string) ([]*Component, []*errors.ParseError) {
	// This is a placeholder implementation
	// In a real implementation, you would parse components here
	_ = content // TODO: Implement component parsing from content
	return []*Component{}, []*errors.ParseError{}
}

// validateDocument validates the parsed document
func (p *Parser) validateDocument(doc *Document) []*errors.ParseError {
	var parseErrors []*errors.ParseError

	// Validate endpoints
	for _, endpoint := range doc.Endpoints {
		if !p.isValidMethod(endpoint.Method) {
			parseErrors = append(parseErrors, errors.NewError(errors.ErrorTypeValidation,
				fmt.Sprintf("Invalid HTTP method: %s", endpoint.Method)).
				AtLine(endpoint.LineNumber).
				WithSuggestion("Use one of: "+strings.Join(p.config.AllowedMethods, ", ")).
				Build())
		}
	}

	return parseErrors
}

// isValidMethod checks if the HTTP method is valid
func (p *Parser) isValidMethod(method string) bool {
	for _, validMethod := range p.config.AllowedMethods {
		if method == validMethod {
			return true
		}
	}
	return false
}

// GetConfig returns a copy of the parser configuration
func (p *Parser) GetConfig() ParserConfig {
	return *p.config
}
