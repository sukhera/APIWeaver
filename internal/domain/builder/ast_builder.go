package builder

import (
	"time"

	"github.com/sukhera/APIWeaver/internal/domain/parser"
	"github.com/sukhera/APIWeaver/pkg/errors"
)

// DocumentBuilder builds Document instances
type DocumentBuilder struct {
	document *parser.Document
	errors   []*errors.ParseError
}

// NewDocumentBuilder creates a new document builder
func NewDocumentBuilder() *DocumentBuilder {
	return &DocumentBuilder{
		document: &parser.Document{
			ParsedAt:  time.Now(),
			Endpoints: []*parser.Endpoint{},
			Errors:    []*errors.ParseError{},
		},
		errors: []*errors.ParseError{},
	}
}

// WithFrontmatter sets the frontmatter
func (b *DocumentBuilder) WithFrontmatter(frontmatter *parser.Frontmatter) *DocumentBuilder {
	b.document.Frontmatter = frontmatter
	return b
}

// AddEndpoint adds an endpoint to the document
func (b *DocumentBuilder) AddEndpoint(endpoint *parser.Endpoint) *DocumentBuilder {
	if endpoint != nil {
		b.document.Endpoints = append(b.document.Endpoints, endpoint)
	}
	return b
}

// AddEndpoints adds multiple endpoints
func (b *DocumentBuilder) AddEndpoints(endpoints []*parser.Endpoint) *DocumentBuilder {
	for _, endpoint := range endpoints {
		b.AddEndpoint(endpoint)
	}
	return b
}

// AddComponent adds a reusable component
func (b *DocumentBuilder) AddComponent(component *parser.Component) *DocumentBuilder {
	if component != nil {
		b.document.Components = append(b.document.Components, component)
	}
	return b
}

// AddComponents adds multiple components
func (b *DocumentBuilder) AddComponents(components []*parser.Component) *DocumentBuilder {
	for _, component := range components {
		b.AddComponent(component)
	}
	return b
}

// AddError adds a parse error
func (b *DocumentBuilder) AddError(err *errors.ParseError) *DocumentBuilder {
	if err != nil {
		b.document.Errors = append(b.document.Errors, err)
		b.errors = append(b.errors, err)
	}
	return b
}

// AddErrors adds multiple parse errors
func (b *DocumentBuilder) AddErrors(errs []*errors.ParseError) *DocumentBuilder {
	for _, err := range errs {
		b.AddError(err)
	}
	return b
}

// Build constructs the final Document
func (b *DocumentBuilder) Build() *parser.Document {
	// Ensure components slice is not nil
	if b.document.Components == nil {
		b.document.Components = []*parser.Component{}
	}

	return b.document
}

// HasErrors returns true if the builder has accumulated errors
func (b *DocumentBuilder) HasErrors() bool {
	return len(b.errors) > 0
}

// HasFatalErrors returns true if there are any fatal errors
func (b *DocumentBuilder) HasFatalErrors() bool {
	for _, err := range b.errors {
		if err.IsFatal() {
			return true
		}
	}
	return false
}

// EndpointBuilder builds Endpoint instances
type EndpointBuilder struct {
	endpoint *parser.Endpoint
}

// NewEndpointBuilder creates a new endpoint builder
func NewEndpointBuilder(method, path string, lineNumber int) *EndpointBuilder {
	return &EndpointBuilder{
		endpoint: &parser.Endpoint{
			Method:     method,
			Path:       path,
			LineNumber: lineNumber,
			Parameters: []*parser.Parameter{},
			Responses:  []*parser.Response{},
			Tags:       []string{},
		},
	}
}

// WithSummary sets the endpoint summary
func (b *EndpointBuilder) WithSummary(summary string) *EndpointBuilder {
	b.endpoint.Summary = summary
	return b
}

// WithDescription sets the endpoint description
func (b *EndpointBuilder) WithDescription(description string) *EndpointBuilder {
	b.endpoint.Description = description
	return b
}

// AddParameter adds a parameter to the endpoint
func (b *EndpointBuilder) AddParameter(param *parser.Parameter) *EndpointBuilder {
	if param != nil {
		b.endpoint.Parameters = append(b.endpoint.Parameters, param)
	}
	return b
}

// AddParameters adds multiple parameters
func (b *EndpointBuilder) AddParameters(params []*parser.Parameter) *EndpointBuilder {
	for _, param := range params {
		b.AddParameter(param)
	}
	return b
}

// WithRequestBody sets the request body
func (b *EndpointBuilder) WithRequestBody(requestBody *parser.RequestBody) *EndpointBuilder {
	b.endpoint.RequestBody = requestBody
	return b
}

// AddResponse adds a response to the endpoint
func (b *EndpointBuilder) AddResponse(response *parser.Response) *EndpointBuilder {
	if response != nil {
		b.endpoint.Responses = append(b.endpoint.Responses, response)
	}
	return b
}

// AddResponses adds multiple responses
func (b *EndpointBuilder) AddResponses(responses []*parser.Response) *EndpointBuilder {
	for _, response := range responses {
		b.AddResponse(response)
	}
	return b
}

// AddTag adds a tag to the endpoint
func (b *EndpointBuilder) AddTag(tag string) *EndpointBuilder {
	if tag != "" {
		b.endpoint.Tags = append(b.endpoint.Tags, tag)
	}
	return b
}

// AddTags adds multiple tags
func (b *EndpointBuilder) AddTags(tags []string) *EndpointBuilder {
	for _, tag := range tags {
		b.AddTag(tag)
	}
	return b
}

// Build constructs the final Endpoint
func (b *EndpointBuilder) Build() *parser.Endpoint {
	return b.endpoint
}

// ParameterBuilder builds Parameter instances
type ParameterBuilder struct {
	parameter *parser.Parameter
}

// NewParameterBuilder creates a new parameter builder
func NewParameterBuilder(name, location string, lineNumber int) *ParameterBuilder {
	return &ParameterBuilder{
		parameter: &parser.Parameter{
			Name:       name,
			In:         location,
			LineNumber: lineNumber,
			Type:       "string", // default type
		},
	}
}

// WithType sets the parameter type
func (b *ParameterBuilder) WithType(paramType string) *ParameterBuilder {
	b.parameter.Type = paramType
	return b
}

// Required marks the parameter as required
func (b *ParameterBuilder) Required() *ParameterBuilder {
	b.parameter.Required = true
	return b
}

// Optional marks the parameter as optional
func (b *ParameterBuilder) Optional() *ParameterBuilder {
	b.parameter.Required = false
	return b
}

// WithDescription sets the parameter description
func (b *ParameterBuilder) WithDescription(description string) *ParameterBuilder {
	b.parameter.Description = description
	return b
}

// WithExample sets the parameter example
func (b *ParameterBuilder) WithExample(example interface{}) *ParameterBuilder {
	b.parameter.Example = example
	return b
}

// WithSchema sets the parameter schema
func (b *ParameterBuilder) WithSchema(schema *parser.Schema) *ParameterBuilder {
	b.parameter.Schema = schema
	return b
}

// Build constructs the final Parameter
func (b *ParameterBuilder) Build() *parser.Parameter {
	return b.parameter
}
