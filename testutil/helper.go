package testutil

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/sukhera/APIWeaver/internal/domain/parser"
	"github.com/sukhera/APIWeaver/pkg/errors"
)

// TestContext creates a test context with timeout
func TestContext(t *testing.T) context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(cancel)
	return ctx
}

// CreateTestDocument creates a test document with default values
func CreateTestDocument() *parser.Document {
	return &parser.Document{
		Frontmatter: &parser.Frontmatter{
			Title:       "Test API",
			Version:     "1.0.0",
			Description: "Test API documentation",
			LineNumber:  1,
		},
		Endpoints: []*parser.Endpoint{
			CreateTestEndpoint(),
		},
		Components: []*parser.Component{},
		ParsedAt:   time.Now(),
		Errors:     []*errors.ParseError{},
	}
}

// CreateTestEndpoint creates a test endpoint with default values
func CreateTestEndpoint() *parser.Endpoint {
	return &parser.Endpoint{
		Method:      "GET",
		Path:        "/api/test",
		Summary:     "Test endpoint",
		Description: "A test endpoint for testing",
		Parameters:  []*parser.Parameter{},
		RequestBody: nil,
		Responses:   []*parser.Response{},
		Tags:        []string{"test"},
		LineNumber:  10,
	}
}

// CreateTestParameter creates a test parameter with default values
func CreateTestParameter() *parser.Parameter {
	return &parser.Parameter{
		Name:        "test_param",
		In:          "query",
		Type:        "string",
		Required:    false,
		Description: "A test parameter",
		Example:     "test_value",
		Schema:      nil,
		LineNumber:  15,
	}
}

// CreateTestSchema creates a test schema with default values
func CreateTestSchema() *parser.Schema {
	return &parser.Schema{
		Type:        "object",
		Format:      "",
		Properties:  map[string]*parser.Schema{},
		Items:       nil,
		Required:    []string{},
		Enum:        []interface{}{},
		Example:     map[string]interface{}{"key": "value"},
		Description: "A test schema",
		Ref:         "",
		AllOf:       []*parser.Schema{},
		OneOf:       []*parser.Schema{},
		AnyOf:       []*parser.Schema{},
		LineNumber:  20,
	}
}

// CreateTestResponse creates a test response with default values
func CreateTestResponse() *parser.Response {
	return &parser.Response{
		StatusCode:  "200",
		Description: "Success response",
		Headers:     map[string]*parser.Header{},
		Content:     map[string]*parser.Schema{},
		LineNumber:  25,
	}
}

// CreateTestComponent creates a test component with default values
func CreateTestComponent() *parser.Component {
	return &parser.Component{
		Name:       "TestComponent",
		Type:       "schema",
		Schema:     CreateTestSchema(),
		LineNumber: 30,
	}
}

// CreateTestHTTPRequest creates a test HTTP request
func CreateTestHTTPRequest(method, path string, body io.Reader) *http.Request {
	req := httptest.NewRequest(method, path, body)
	req.Header.Set("Content-Type", "application/json")
	return req
}

// CreateTestHTTPResponse creates a test HTTP response recorder
func CreateTestHTTPResponse() *httptest.ResponseRecorder {
	return httptest.NewRecorder()
}

// AssertJSONResponse asserts that a response body matches expected JSON
func AssertJSONResponse(t *testing.T, responseBody []byte, expected any) {
	var actual map[string]interface{}
	err := json.Unmarshal(responseBody, &actual)
	require.NoError(t, err)

	expectedBytes, err := json.Marshal(expected)
	require.NoError(t, err)

	var expectedMap map[string]interface{}
	err = json.Unmarshal(expectedBytes, &expectedMap)
	require.NoError(t, err)

	assert.Equal(t, expectedMap, actual)
}

// CreateTestParseError creates a test parse error
func CreateTestParseError() *errors.ParseError {
	return errors.NewError(errors.ErrorTypeValidation, "Test error").
		AtLine(1).
		WithContext("test").
		Build()
}

// MockCleanup ensures mocks are properly cleaned up
func MockCleanup(t *testing.T, mocks ...interface{}) {
	t.Cleanup(func() {
		for _, mock := range mocks {
			if m, ok := mock.(interface{ AssertExpectations(testing.TB) }); ok {
				m.AssertExpectations(t)
			}
		}
	})
}

// CreateTestMarkdownContent creates test markdown content
func CreateTestMarkdownContent() string {
	return `---
title: Test API
version: 1.0.0
description: Test API documentation
---

# Test API

## GET /api/test

Test endpoint description

### Parameters

- **test_param** (query, string, optional) - A test parameter

### Responses

- **200** - Success response
  - Content-Type: application/json
  - Schema:
    ` + "```json" + `
    {
      "type": "object",
      "properties": {
        "message": {
          "type": "string"
        }
      }
    }
    ` + "```" + `

### Example

` + "```bash" + `
curl -X GET "https://api.example.com/api/test?test_param=value"
` + "```" + `
`
}

// CreateTestInvalidMarkdownContent creates invalid markdown content for testing
func CreateTestInvalidMarkdownContent() string {
	return `---
title: Test API
version: 1.0.0
description: Test API documentation
---

# Test API

## INVALID_METHOD /invalid/path

Invalid endpoint with invalid HTTP method

### Parameters

- **invalid_param** (invalid_location, invalid_type, required) - Invalid parameter
`
}
