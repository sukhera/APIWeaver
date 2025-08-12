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

func TestParser_ValidateDocument(t *testing.T) {
	tests := []struct {
		name          string
		doc           *Document
		strictMode    bool
		expectedCount int
	}{
		{
			name: "success with valid document",
			doc: &Document{
				Endpoints: []*Endpoint{
					{
						Method:     "GET",
						Path:       "/test",
						LineNumber: 1,
					},
				},
			},
			strictMode:    false,
			expectedCount: 0,
		},
		{
			name: "error with invalid HTTP method",
			doc: &Document{
				Endpoints: []*Endpoint{
					{
						Method:     "INVALID",
						Path:       "/test",
						LineNumber: 1,
					},
				},
			},
			strictMode:    false,
			expectedCount: 1,
		},
		{
			name: "error with invalid path",
			doc: &Document{
				Endpoints: []*Endpoint{
					{
						Method:     "GET",
						Path:       "invalid-path",
						LineNumber: 1,
					},
				},
			},
			strictMode:    false,
			expectedCount: 1,
		},
		{
			name: "warning with missing description in strict mode",
			doc: &Document{
				Endpoints: []*Endpoint{
					{
						Method:     "GET",
						Path:       "/test",
						LineNumber: 1,
					},
				},
			},
			strictMode:    true,
			expectedCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := ValidateDocument(context.Background(), tt.doc, tt.strictMode)
			assert.Len(t, errors, tt.expectedCount)
		})
	}
}

func TestParser_GetDocumentStatistics(t *testing.T) {
	tests := []struct {
		name          string
		doc           *Document
		expectedStats DocumentStatistics
	}{
		{
			name: "success with single endpoint",
			doc: &Document{
				Endpoints: []*Endpoint{
					{
						Method: "GET",
						Path:   "/test",
					},
				},
				Components: []*Component{},
			},
			expectedStats: DocumentStatistics{
				TotalEndpoints:    1,
				EndpointsByMethod: map[string]int{"GET": 1},
				TotalComponents:   0,
				HasFrontmatter:    false,
			},
		},
		{
			name: "success with multiple endpoints",
			doc: &Document{
				Endpoints: []*Endpoint{
					{Method: "GET", Path: "/test1"},
					{Method: "POST", Path: "/test2"},
					{Method: "GET", Path: "/test3"},
				},
				Components: []*Component{
					{Name: "TestComponent", Type: "schema"},
				},
			},
			expectedStats: DocumentStatistics{
				TotalEndpoints:    3,
				EndpointsByMethod: map[string]int{"GET": 2, "POST": 1},
				TotalComponents:   1,
				HasFrontmatter:    false,
			},
		},
		{
			name: "success with frontmatter",
			doc: &Document{
				Frontmatter: &Frontmatter{Title: "Test API"},
				Endpoints:   []*Endpoint{},
				Components:  []*Component{},
			},
			expectedStats: DocumentStatistics{
				TotalEndpoints:  0,
				TotalComponents: 0,
				HasFrontmatter:  true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stats := GetDocumentStatistics(context.Background(), tt.doc)

			assert.Equal(t, tt.expectedStats.TotalEndpoints, stats.TotalEndpoints)
			assert.Equal(t, tt.expectedStats.TotalComponents, stats.TotalComponents)
			assert.Equal(t, tt.expectedStats.HasFrontmatter, stats.HasFrontmatter)

			for method, count := range tt.expectedStats.EndpointsByMethod {
				assert.Equal(t, count, stats.EndpointsByMethod[method])
			}
		})
	}
}

func TestParser_TransformDocument(t *testing.T) {
	tests := []struct {
		name          string
		doc           *Document
		transforms    []func(interface{}) interface{}
		expectedError bool
	}{
		{
			name: "success with no transformations",
			doc: &Document{
				Endpoints: []*Endpoint{
					{Method: "GET", Path: "/test"},
				},
			},
			transforms:    []func(interface{}) interface{}{},
			expectedError: false,
		},
		{
			name: "success with path transformation",
			doc: &Document{
				Endpoints: []*Endpoint{
					{Method: "GET", Path: "/test"},
				},
			},
			transforms: []func(interface{}) interface{}{
				func(obj interface{}) interface{} {
					if endpoint, ok := obj.(*Endpoint); ok {
						endpoint.Path = "/transformed" + endpoint.Path
					}
					return obj
				},
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := TransformDocument(context.Background(), tt.doc, tt.transforms...)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
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
