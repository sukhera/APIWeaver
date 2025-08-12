# APIWeaver Errors Package

The `pkg/errors` package provides a comprehensive error handling system for the APIWeaver project. It offers structured error types, error collection, and formatting utilities that can be used across all packages in the project.

## Features

- **Structured Error Types**: Different error types for various scenarios (syntax, validation, configuration, etc.)
- **Error Severity Levels**: Info, warning, error, and fatal severity levels
- **Error Collection**: Collect and manage multiple errors with context
- **Error Formatting**: Rich error messages with line numbers, context, and suggestions
- **Builder Pattern**: Fluent API for creating complex errors
- **Configuration Errors**: Specialized error types for configuration issues
- **Timeout Errors**: Specialized error types for timeout scenarios

## Quick Start

```go
import "github.com/sukhera/APIWeaver/pkg/errors"

// Create a simple error
err := errors.NewError(errors.ErrorTypeValidation, "Invalid HTTP method").Build()

// Create a complex error with context
err = errors.NewError(errors.ErrorTypeSyntax, "Missing closing brace").
    WithCode("SYNTAX_001").
    AtLine(10).
    AtColumn(5).
    WithContext("function definition").
    WithSuggestion("Add closing brace '}'").
    InSource("parser").
    Build()

// Create a configuration error
configErr := errors.NewConfigError("Invalid timeout value")

// Create a timeout error
timeoutErr := errors.NewTimeoutError("database query", "30s")
```

## Error Types

### ParseError

The main error type for parsing-related issues:

```go
type ParseError struct {
    Type       ErrorType `json:"type"`
    Code       string    `json:"code,omitempty"`
    Message    string    `json:"message"`
    LineNumber int       `json:"line_number"`
    Column     int       `json:"column,omitempty"`
    Context    string    `json:"context,omitempty"`
    Suggestion string    `json:"suggestion,omitempty"`
    Source     string    `json:"source,omitempty"`
    Severity   Severity  `json:"severity"`
}
```

### ErrorType Constants

- `ErrorTypeSyntax` - Syntax errors
- `ErrorTypeValidation` - Validation errors
- `ErrorTypeConfig` - Configuration errors
- `ErrorTypeTimeout` - Timeout errors
- `ErrorTypeSchema` - Schema-related errors
- `ErrorTypeTable` - Table parsing errors
- `ErrorTypeFrontmatter` - Frontmatter errors
- `ErrorTypeEndpoint` - Endpoint parsing errors
- `ErrorTypeReference` - Reference errors

### Severity Levels

- `SeverityInfo` - Informational messages
- `SeverityWarning` - Warning messages
- `SeverityError` - Error messages
- `SeverityFatal` - Fatal errors

## Error Builder

Use the builder pattern to create complex errors:

```go
err := errors.NewError(errors.ErrorTypeValidation, "Invalid parameter").
    WithCode("INVALID_PARAM").
    AtLine(15).
    AtColumn(10).
    WithContext("GET /users endpoint").
    WithSuggestion("Use a valid parameter name").
    InSource("endpoint").
    Build()
```

### Builder Methods

- `WithCode(code string)` - Set error code
- `AtLine(line int)` - Set line number
- `AtColumn(column int)` - Set column number
- `WithContext(context string)` - Set error context
- `WithSuggestion(suggestion string)` - Set suggestion for fixing
- `InSource(source string)` - Set error source
- `Build()` - Create the final error

## Error Collector

Collect and manage multiple errors:

```go
collector := errors.NewErrorCollector(100) // max 100 errors

// Add errors
collector.Add(errors.NewError(errors.ErrorTypeSyntax, "Error 1").Build())
collector.Add(errors.NewError(errors.ErrorTypeValidation, "Error 2").Build())

// Check for errors
if collector.HasErrors() {
    fmt.Printf("Found %d errors\n", len(collector.GetErrors()))
}

// Get all errors
allErrors := collector.GetErrors()

// Clear all errors
collector.Clear()
```

### Collector Methods

- `Add(err *ParseError)` - Add a single error
- `AddMultiple(errors []*ParseError)` - Add multiple errors
- `GetErrors()` - Get all errors
- `GetWarnings()` - Get all warnings
- `GetAll()` - Get all errors and warnings
- `HasErrors()` - Check if there are any errors
- `HasWarnings()` - Check if there are any warnings
- `HasFatalErrors()` - Check if there are any fatal errors
- `Clear()` - Clear all errors and warnings
- `ToError()` - Convert to standard error

## Specialized Error Types

### Configuration Errors

```go
err := errors.NewConfigError("Invalid configuration value")
// Output: "configuration error: Invalid configuration value"
```

### Timeout Errors

```go
err := errors.NewTimeoutError("database query", "30s")
// Output: "operation 'database query' timed out after 30s"
```

## Error Formatting

Format multiple errors into a readable string:

```go
errors := []*ParseError{
    errors.NewError(errors.ErrorTypeValidation, "Error 1").AtLine(1).Build(),
    errors.NewError(errors.ErrorTypeSyntax, "Error 2").AtLine(2).Build(),
}

formatted := errors.FormatErrors(errors)
fmt.Println(formatted)
```

## Error Methods

### ParseError Methods

- `Error()` - Return formatted error string
- `IsError()` - Check if severity is error or fatal
- `IsWarning()` - Check if severity is warning
- `IsFatal()` - Check if severity is fatal

## Usage Examples

### In Configuration Validation

```go
func (c *Config) Validate() error {
    if c.Timeout <= 0 {
        return errors.NewConfigError("timeout must be positive")
    }
    
    if c.MaxRetries < 0 {
        return errors.NewConfigError("max_retries must be non-negative")
    }
    
    return nil
}
```

### In Parser

```go
func (p *Parser) Parse(content string) (*Document, error) {
    collector := errors.NewErrorCollector(50)
    
    // Parse content and collect errors
    if invalidMethod {
        collector.Add(errors.NewError(errors.ErrorTypeValidation, "Invalid HTTP method").
            AtLine(lineNumber).
            WithSuggestion("Use GET, POST, PUT, DELETE, etc.").
            Build())
    }
    
    // Return document with errors
    return &Document{
        Errors: collector.GetErrors(),
    }, collector.ToError()
}
```

### In API Handler

```go
func (h *Handler) ProcessRequest(req *Request) error {
    if req.Timeout > maxTimeout {
        return errors.NewTimeoutError("request processing", fmt.Sprintf("%ds", req.Timeout))
    }
    
    return nil
}
```

## Best Practices

1. **Use Appropriate Error Types**: Choose the right error type for your scenario
2. **Provide Context**: Always include relevant context in error messages
3. **Add Suggestions**: Include helpful suggestions for fixing errors
4. **Use Error Collector**: For operations that might generate multiple errors
5. **Set Line Numbers**: Include line numbers for parsing errors
6. **Use Error Codes**: Include error codes for programmatic error handling
7. **Check Severity**: Use severity levels appropriately (info, warning, error, fatal)

## Integration

This errors package is designed to be used across all packages in the APIWeaver project:

- `internal/config` - Configuration validation errors
- `internal/domain/parser` - Parsing and validation errors
- `internal/domain/builder` - AST building errors
- `cmd/apiweaver` - Application-level errors

The package provides a consistent error handling experience throughout the application.
