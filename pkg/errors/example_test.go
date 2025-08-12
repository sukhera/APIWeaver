package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Example demonstrates how to use the errors package
func Example() {
	// Create a new parse error
	err := NewError(ErrorTypeValidation, "Invalid HTTP method").
		WithCode("INVALID_METHOD").
		AtLine(10).
		WithContext("GET /users").
		WithSuggestion("Use a valid HTTP method like GET, POST, PUT, etc.").
		InSource("endpoint").
		Build()

	fmt.Println(err.Error())
	// Output: [INVALID_METHOD] line 10 in endpoint Invalid HTTP method (suggestion: Use a valid HTTP method like GET, POST, PUT, etc.)
}

// ExampleNewConfigError demonstrates creating configuration errors
func ExampleNewConfigError() {
	err := NewConfigError("Invalid configuration value")
	fmt.Println(err.Error())
	// Output: configuration error: Invalid configuration value
}

// ExampleNewTimeoutError demonstrates creating timeout errors
func ExampleNewTimeoutError() {
	err := NewTimeoutError("parsing", "30s")
	fmt.Println(err.Error())
	// Output: operation 'parsing' timed out after 30s
}

// ExampleErrorCollector demonstrates using the error collector
func ExampleErrorCollector() {
	collector := NewErrorCollector(100) // max 100 errors

	// Add some errors
	collector.Add(NewError(ErrorTypeSyntax, "Missing closing brace").AtLine(5).Build())
	collector.Add(NewError(ErrorTypeValidation, "Invalid parameter").AtLine(10).Build())

	// Get all errors
	errors := collector.GetErrors()
	fmt.Printf("Found %d errors\n", len(errors))
	// Output: Found 2 errors
}

// TestErrorUsage demonstrates various ways to use the errors package
func TestErrorUsage(t *testing.T) {
	// Test creating a basic error
	err := NewError(ErrorTypeValidation, "Test error").Build()
	assert.NotNil(t, err)
	assert.Equal(t, ErrorTypeValidation, err.Type)
	assert.Equal(t, "Test error", err.Message)

	// Test creating an error with all fields
	err = NewError(ErrorTypeSyntax, "Complex error").
		WithCode("TEST_001").
		AtLine(42).
		AtColumn(10).
		WithContext("test context").
		WithSuggestion("Try fixing it").
		InSource("test").
		Build()

	assert.NotNil(t, err)
	assert.Equal(t, "TEST_001", err.Code)
	assert.Equal(t, 42, err.LineNumber)
	assert.Equal(t, 10, err.Column)
	assert.Equal(t, "test context", err.Context)
	assert.Equal(t, "Try fixing it", err.Suggestion)
	assert.Equal(t, "test", err.Source)

	// Test error formatting
	expected := "[TEST_001] line 42:10 in test Complex error (suggestion: Try fixing it)"
	assert.Equal(t, expected, err.Error())

	// Test severity methods
	assert.True(t, err.IsError())
	assert.False(t, err.IsWarning())
	assert.False(t, err.IsFatal())
}

// TestErrorCollectorUsage demonstrates using the error collector
func TestErrorCollectorUsage(t *testing.T) {
	collector := NewErrorCollector(100) // max 100 errors

	// Test adding errors
	err1 := NewError(ErrorTypeValidation, "Error 1").Build()
	err2 := NewError(ErrorTypeSyntax, "Error 2").Build()

	collector.Add(err1)
	collector.Add(err2)

	errors := collector.GetErrors()
	assert.Len(t, errors, 2)

	// Test error checking
	assert.True(t, collector.HasErrors())
	assert.False(t, collector.HasWarnings())

	// Test clearing errors
	collector.Clear()
	assert.Len(t, collector.GetErrors(), 0)
}

// TestConfigErrorUsage demonstrates configuration error usage
func TestConfigErrorUsage(t *testing.T) {
	err := NewConfigError("Invalid setting")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "configuration error")
	assert.Contains(t, err.Error(), "Invalid setting")
}

// TestTimeoutErrorUsage demonstrates timeout error usage
func TestTimeoutErrorUsage(t *testing.T) {
	err := NewTimeoutError("database query", "5s")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "database query")
	assert.Contains(t, err.Error(), "5s")
}

// TestErrorFormatting demonstrates error formatting utilities
func TestErrorFormatting(t *testing.T) {
	errors := []*ParseError{
		NewError(ErrorTypeValidation, "Error 1").AtLine(1).Build(),
		NewError(ErrorTypeSyntax, "Error 2").AtLine(2).Build(),
	}

	formatted := FormatErrors(errors)
	assert.Contains(t, formatted, "Error 1")
	assert.Contains(t, formatted, "Error 2")
	assert.Contains(t, formatted, "line 1")
	assert.Contains(t, formatted, "line 2")
}
