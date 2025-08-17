package parser

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParser_New(t *testing.T) {
	tests := []struct {
		name     string
		options  []ParserOption
		expected bool
	}{
		{
			name:     "success with default options",
			options:  []ParserOption{},
			expected: true,
		},
		{
			name: "success with strict mode",
			options: []ParserOption{
				WithStrictMode(true),
			},
			expected: true,
		},
		{
			name: "success with custom timeout",
			options: []ParserOption{
				WithTimeout(60),
			},
			expected: true,
		},
		{
			name: "success with multiple options",
			options: []ParserOption{
				WithStrictMode(true),
				WithTimeout(30),
				WithAllowedMethods([]string{"GET", "POST"}),
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := New(tt.options...)
			assert.Equal(t, tt.expected, parser != nil)
		})
	}
}

func TestParser_Parse(t *testing.T) {
	tests := []struct {
		name          string
		content       string
		setupParser   func(*Parser)
		expectedError bool
		expectedDoc   bool
	}{
		{
			name:    "success with valid markdown",
			content: "# Test API\n\n## GET /test\n\nTest endpoint",
			setupParser: func(p *Parser) {
				// Use default parser settings
			},
			expectedError: false,
			expectedDoc:   true,
		},
		{
			name:    "success with empty content",
			content: "",
			setupParser: func(p *Parser) {
				// Use default parser settings
			},
			expectedError: false,
			expectedDoc:   true,
		},
		{
			name:    "success with minimal content",
			content: "# Test API\n\n## GET /test\n\nTest endpoint",
			setupParser: func(p *Parser) {
				// Use default parser settings
			},
			expectedError: false,
			expectedDoc:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := New()
			if tt.setupParser != nil {
				tt.setupParser(parser)
			}

			doc, err := parser.Parse(tt.content)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, doc)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, doc)
				if tt.expectedDoc {
					assert.NotNil(t, doc)
					assert.NotZero(t, doc.ParsedAt)
				}
			}
		})
	}
}

func TestParser_ParseWithContext(t *testing.T) {
	tests := []struct {
		name          string
		content       string
		setupParser   func(*Parser)
		expectedError bool
		expectedDoc   bool
	}{
		{
			name:    "success with valid markdown",
			content: "# Test API\n\n## GET /test\n\nTest endpoint",
			setupParser: func(p *Parser) {
				// Use default parser settings
			},
			expectedError: false,
			expectedDoc:   true,
		},
		{
			name:    "success with empty content",
			content: "",
			setupParser: func(p *Parser) {
				// Use default parser settings
			},
			expectedError: false,
			expectedDoc:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			parser := New()
			if tt.setupParser != nil {
				tt.setupParser(parser)
			}

			doc, err := parser.ParseWithContext(ctx, tt.content)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, doc)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, doc)
				if tt.expectedDoc {
					assert.NotNil(t, doc)
					assert.NotZero(t, doc.ParsedAt)
				}
			}
		})
	}
}

func TestParser_GetConfig(t *testing.T) {
	tests := []struct {
		name     string
		options  []ParserOption
		expected ParserConfig
	}{
		{
			name:    "default config",
			options: []ParserOption{},
			expected: ParserConfig{
				StrictMode:     false,
				Timeout:        30 * time.Second,
				AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
			},
		},
		{
			name: "custom config",
			options: []ParserOption{
				WithStrictMode(true),
				WithTimeout(60 * time.Second),
				WithAllowedMethods([]string{"GET", "POST"}),
			},
			expected: ParserConfig{
				StrictMode:     true,
				Timeout:        60 * time.Second,
				AllowedMethods: []string{"GET", "POST"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := New(tt.options...)
			config := parser.GetConfig()

			assert.Equal(t, tt.expected.StrictMode, config.StrictMode)
			assert.Equal(t, tt.expected.Timeout, config.Timeout)
			assert.Equal(t, tt.expected.AllowedMethods, config.AllowedMethods)
		})
	}
}
