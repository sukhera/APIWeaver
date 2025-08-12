package errors

// ErrorBuilder provides a fluent interface for building ParseErrors
type ErrorBuilder struct {
	error *ParseError
}

// NewError creates a new error builder
func NewError(errorType ErrorType, message string) *ErrorBuilder {
	return &ErrorBuilder{
		error: &ParseError{
			Type:     errorType,
			Message:  message,
			Severity: SeverityError,
		},
	}
}

// NewWarning creates a new warning builder
func NewWarning(errorType ErrorType, message string) *ErrorBuilder {
	return &ErrorBuilder{
		error: &ParseError{
			Type:     errorType,
			Message:  message,
			Severity: SeverityWarning,
		},
	}
}

// NewFatal creates a new fatal error builder
func NewFatal(errorType ErrorType, message string) *ErrorBuilder {
	return &ErrorBuilder{
		error: &ParseError{
			Type:     errorType,
			Message:  message,
			Severity: SeverityFatal,
		},
	}
}

// WithCode sets the error code
func (b *ErrorBuilder) WithCode(code string) *ErrorBuilder {
	b.error.Code = code
	return b
}

// AtLine sets the line number
func (b *ErrorBuilder) AtLine(line int) *ErrorBuilder {
	b.error.LineNumber = line
	return b
}

// AtColumn sets the column number
func (b *ErrorBuilder) AtColumn(column int) *ErrorBuilder {
	b.error.Column = column
	return b
}

// AtPosition sets both line and column
func (b *ErrorBuilder) AtPosition(line, column int) *ErrorBuilder {
	b.error.LineNumber = line
	b.error.Column = column
	return b
}

// WithContext sets the error context
func (b *ErrorBuilder) WithContext(context string) *ErrorBuilder {
	b.error.Context = context
	return b
}

// WithSuggestion sets a suggestion for fixing the error
func (b *ErrorBuilder) WithSuggestion(suggestion string) *ErrorBuilder {
	b.error.Suggestion = suggestion
	return b
}

// InSource sets the source component
func (b *ErrorBuilder) InSource(source string) *ErrorBuilder {
	b.error.Source = source
	return b
}

// Build creates the final ParseError
func (b *ErrorBuilder) Build() *ParseError {
	return b.error
}

// Predefined error constructors for common scenarios

// NewSyntaxError creates a syntax error
func NewSyntaxError(message string, line int) *ParseError {
	return NewError(ErrorTypeSyntax, message).
		AtLine(line).
		Build()
}

// NewValidationError creates a validation error
func NewValidationError(message string, line int) *ParseError {
	return NewError(ErrorTypeValidation, message).
		AtLine(line).
		Build()
}

// NewSchemaError creates a schema-related error
func NewSchemaError(message string, line int) *ParseError {
	return NewError(ErrorTypeSchema, message).
		AtLine(line).
		InSource("schema").
		Build()
}

// NewEndpointError creates an endpoint-related error
func NewEndpointError(message string, line int) *ParseError {
	return NewError(ErrorTypeEndpoint, message).
		AtLine(line).
		InSource("endpoint").
		Build()
}

// NewFrontmatterError creates a frontmatter-related error
func NewFrontmatterError(message string, line int) *ParseError {
	return NewError(ErrorTypeFrontmatter, message).
		AtLine(line).
		InSource("frontmatter").
		Build()
}

// NewTableError creates a table-related error
func NewTableError(message string, line int) *ParseError {
	return NewError(ErrorTypeTable, message).
		AtLine(line).
		InSource("table").
		Build()
}
