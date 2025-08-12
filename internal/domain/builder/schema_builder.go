package builder

import (
	"github.com/sukhera/APIWeaver/internal/domain/parser"
)

// SchemaBuilder builds Schema instances using a fluent interface
type SchemaBuilder struct {
	schema *parser.Schema
}

// NewSchemaBuilder creates a new schema builder
func NewSchemaBuilder(lineNumber int) *SchemaBuilder {
	return &SchemaBuilder{
		schema: &parser.Schema{
			LineNumber: lineNumber,
			Properties: make(map[string]*parser.Schema),
		},
	}
}

// WithType sets the schema type
func (b *SchemaBuilder) WithType(schemaType string) *SchemaBuilder {
	b.schema.Type = schemaType
	return b
}

// WithFormat sets the schema format
func (b *SchemaBuilder) WithFormat(format string) *SchemaBuilder {
	b.schema.Format = format
	return b
}

// WithDescription sets the schema description
func (b *SchemaBuilder) WithDescription(description string) *SchemaBuilder {
	b.schema.Description = description
	return b
}

// WithExample sets the schema example
func (b *SchemaBuilder) WithExample(example interface{}) *SchemaBuilder {
	b.schema.Example = example
	return b
}

// AddProperty adds a property to an object schema
func (b *SchemaBuilder) AddProperty(name string, property *parser.Schema) *SchemaBuilder {
	if name != "" && property != nil {
		b.schema.Properties[name] = property
	}
	return b
}

// WithItems sets the items schema for array types
func (b *SchemaBuilder) WithItems(items *parser.Schema) *SchemaBuilder {
	b.schema.Items = items
	return b
}

// AddRequired adds a required property name
func (b *SchemaBuilder) AddRequired(propertyName string) *SchemaBuilder {
	if propertyName != "" {
		b.schema.Required = append(b.schema.Required, propertyName)
	}
	return b
}

// WithEnum sets the enum values
func (b *SchemaBuilder) WithEnum(values []interface{}) *SchemaBuilder {
	b.schema.Enum = values
	return b
}

// WithRef sets a reference to another schema
func (b *SchemaBuilder) WithRef(ref string) *SchemaBuilder {
	b.schema.Ref = ref
	return b
}

// Build constructs the final Schema
func (b *SchemaBuilder) Build() *parser.Schema {
	return b.schema
}

// ResponseBuilder builds Response instances
type ResponseBuilder struct {
	response *parser.Response
}

// NewResponseBuilder creates a new response builder
func NewResponseBuilder(statusCode string, lineNumber int) *ResponseBuilder {
	return &ResponseBuilder{
		response: &parser.Response{
			StatusCode: statusCode,
			LineNumber: lineNumber,
			Headers:    make(map[string]*parser.Header),
			Content:    make(map[string]*parser.Schema),
		},
	}
}

// WithDescription sets the response description
func (b *ResponseBuilder) WithDescription(description string) *ResponseBuilder {
	b.response.Description = description
	return b
}

// AddHeader adds a response header
func (b *ResponseBuilder) AddHeader(name string, header *parser.Header) *ResponseBuilder {
	if name != "" && header != nil {
		b.response.Headers[name] = header
	}
	return b
}

// AddContent adds response content for a media type
func (b *ResponseBuilder) AddContent(mediaType string, schema *parser.Schema) *ResponseBuilder {
	if mediaType != "" && schema != nil {
		b.response.Content[mediaType] = schema
	}
	return b
}

// Build constructs the final Response
func (b *ResponseBuilder) Build() *parser.Response {
	return b.response
}

// RequestBodyBuilder builds RequestBody instances
type RequestBodyBuilder struct {
	requestBody *parser.RequestBody
}

// NewRequestBodyBuilder creates a new request body builder
func NewRequestBodyBuilder(lineNumber int) *RequestBodyBuilder {
	return &RequestBodyBuilder{
		requestBody: &parser.RequestBody{
			LineNumber: lineNumber,
			Content:    make(map[string]*parser.Schema),
		},
	}
}

// WithDescription sets the request body description
func (b *RequestBodyBuilder) WithDescription(description string) *RequestBodyBuilder {
	b.requestBody.Description = description
	return b
}

// Required marks the request body as required
func (b *RequestBodyBuilder) Required() *RequestBodyBuilder {
	b.requestBody.Required = true
	return b
}

// Optional marks the request body as optional
func (b *RequestBodyBuilder) Optional() *RequestBodyBuilder {
	b.requestBody.Required = false
	return b
}

// AddContent adds content for a media type
func (b *RequestBodyBuilder) AddContent(mediaType string, schema *parser.Schema) *RequestBodyBuilder {
	if mediaType != "" && schema != nil {
		b.requestBody.Content[mediaType] = schema
	}
	return b
}

// Build constructs the final RequestBody
func (b *RequestBodyBuilder) Build() *parser.RequestBody {
	return b.requestBody
}

// Convenience methods for common schema patterns

// BuildStringSchema creates a simple string schema
func BuildStringSchema(lineNumber int) *parser.Schema {
	return NewSchemaBuilder(lineNumber).
		WithType("string").
		Build()
}

// BuildIntegerSchema creates a simple integer schema
func BuildIntegerSchema(lineNumber int) *parser.Schema {
	return NewSchemaBuilder(lineNumber).
		WithType("integer").
		Build()
}

// BuildArraySchema creates an array schema with the given item schema
func BuildArraySchema(itemSchema *parser.Schema, lineNumber int) *parser.Schema {
	return NewSchemaBuilder(lineNumber).
		WithType("array").
		WithItems(itemSchema).
		Build()
}

// BuildObjectSchema creates an object schema with the given properties
func BuildObjectSchema(properties map[string]*parser.Schema, required []string, lineNumber int) *parser.Schema {
	builder := NewSchemaBuilder(lineNumber).WithType("object")

	for name, prop := range properties {
		builder.AddProperty(name, prop)
	}

	for _, req := range required {
		builder.AddRequired(req)
	}

	return builder.Build()
}
