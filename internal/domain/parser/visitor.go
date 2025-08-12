package parser

import (
	"context"
	"strings"

	"github.com/sukhera/APIWeaver/pkg/errors"
)

// Visitor pattern for traversing and manipulating AST nodes
// This allows for separation of algorithms from the data structure

// Visitor interface defines methods for visiting different AST node types
type Visitor interface {
	VisitDocument(ctx context.Context, doc *Document) error
	VisitFrontmatter(ctx context.Context, frontmatter *Frontmatter) error
	VisitEndpoint(ctx context.Context, endpoint *Endpoint) error
	VisitParameter(ctx context.Context, parameter *Parameter) error
	VisitRequestBody(ctx context.Context, requestBody *RequestBody) error
	VisitResponse(ctx context.Context, response *Response) error
	VisitSchema(ctx context.Context, schema *Schema) error
	VisitComponent(ctx context.Context, component *Component) error
}

// Visitable interface for AST nodes that can accept visitors
type Visitable interface {
	Accept(ctx context.Context, visitor Visitor) error
}

// Make AST nodes implement Visitable interface

// Document Accept method
func (d *Document) Accept(ctx context.Context, visitor Visitor) error {
	if err := visitor.VisitDocument(ctx, d); err != nil {
		return err
	}

	// Visit frontmatter
	if d.Frontmatter != nil {
		if err := d.Frontmatter.Accept(ctx, visitor); err != nil {
			return err
		}
	}

	// Visit all endpoints
	for _, endpoint := range d.Endpoints {
		if err := endpoint.Accept(ctx, visitor); err != nil {
			return err
		}
	}

	// Visit all components
	for _, component := range d.Components {
		if err := component.Accept(ctx, visitor); err != nil {
			return err
		}
	}

	return nil
}

// Frontmatter Accept method
func (f *Frontmatter) Accept(ctx context.Context, visitor Visitor) error {
	return visitor.VisitFrontmatter(ctx, f)
}

// Endpoint Accept method
func (e *Endpoint) Accept(ctx context.Context, visitor Visitor) error {
	if err := visitor.VisitEndpoint(ctx, e); err != nil {
		return err
	}

	// Visit all parameters
	for _, param := range e.Parameters {
		if err := param.Accept(ctx, visitor); err != nil {
			return err
		}
	}

	// Visit request body
	if e.RequestBody != nil {
		if err := e.RequestBody.Accept(ctx, visitor); err != nil {
			return err
		}
	}

	// Visit all responses
	for _, response := range e.Responses {
		if err := response.Accept(ctx, visitor); err != nil {
			return err
		}
	}

	return nil
}

// Parameter Accept method
func (p *Parameter) Accept(ctx context.Context, visitor Visitor) error {
	if err := visitor.VisitParameter(ctx, p); err != nil {
		return err
	}

	// Visit schema if present
	if p.Schema != nil {
		if err := p.Schema.Accept(ctx, visitor); err != nil {
			return err
		}
	}

	return nil
}

// RequestBody Accept method
func (r *RequestBody) Accept(ctx context.Context, visitor Visitor) error {
	if err := visitor.VisitRequestBody(ctx, r); err != nil {
		return err
	}

	// Visit all content schemas
	for _, schema := range r.Content {
		if err := schema.Accept(ctx, visitor); err != nil {
			return err
		}
	}

	return nil
}

// Response Accept method
func (r *Response) Accept(ctx context.Context, visitor Visitor) error {
	if err := visitor.VisitResponse(ctx, r); err != nil {
		return err
	}

	// Visit all content schemas
	for _, schema := range r.Content {
		if err := schema.Accept(ctx, visitor); err != nil {
			return err
		}
	}

	return nil
}

// Schema Accept method
func (s *Schema) Accept(ctx context.Context, visitor Visitor) error {
	if err := visitor.VisitSchema(ctx, s); err != nil {
		return err
	}

	// Visit all properties
	for _, prop := range s.Properties {
		if err := prop.Accept(ctx, visitor); err != nil {
			return err
		}
	}

	// Visit items if present
	if s.Items != nil {
		if err := s.Items.Accept(ctx, visitor); err != nil {
			return err
		}
	}

	// Visit allOf schemas
	for _, subSchema := range s.AllOf {
		if err := subSchema.Accept(ctx, visitor); err != nil {
			return err
		}
	}

	// Visit anyOf schemas
	for _, subSchema := range s.AnyOf {
		if err := subSchema.Accept(ctx, visitor); err != nil {
			return err
		}
	}

	// Visit oneOf schemas
	for _, subSchema := range s.OneOf {
		if err := subSchema.Accept(ctx, visitor); err != nil {
			return err
		}
	}

	return nil
}

// Component Accept method
func (c *Component) Accept(ctx context.Context, visitor Visitor) error {
	if err := visitor.VisitComponent(ctx, c); err != nil {
		return err
	}

	// Visit schema if present
	if c.Schema != nil {
		if err := c.Schema.Accept(ctx, visitor); err != nil {
			return err
		}
	}

	return nil
}

// Base visitor implementation that provides default no-op behavior
type BaseVisitor struct{}

func (v *BaseVisitor) VisitDocument(ctx context.Context, doc *Document) error { return nil }
func (v *BaseVisitor) VisitFrontmatter(ctx context.Context, frontmatter *Frontmatter) error {
	return nil
}
func (v *BaseVisitor) VisitEndpoint(ctx context.Context, endpoint *Endpoint) error    { return nil }
func (v *BaseVisitor) VisitParameter(ctx context.Context, parameter *Parameter) error { return nil }
func (v *BaseVisitor) VisitRequestBody(ctx context.Context, requestBody *RequestBody) error {
	return nil
}
func (v *BaseVisitor) VisitResponse(ctx context.Context, response *Response) error    { return nil }
func (v *BaseVisitor) VisitSchema(ctx context.Context, schema *Schema) error          { return nil }
func (v *BaseVisitor) VisitComponent(ctx context.Context, component *Component) error { return nil }

// Concrete visitor implementations

// ValidationVisitor performs comprehensive validation during traversal
type ValidationVisitor struct {
	BaseVisitor
	errors      []*errors.ParseError
	strictMode  bool
	currentPath string
}

func NewValidationVisitor(strictMode bool) *ValidationVisitor {
	return &ValidationVisitor{
		errors:     []*errors.ParseError{},
		strictMode: strictMode,
	}
}

func (v *ValidationVisitor) VisitDocument(ctx context.Context, doc *Document) error {
	v.currentPath = "document"

	if len(doc.Endpoints) == 0 {
		v.addError("error", "document must contain at least one endpoint", 0)
	}

	// Check for duplicate endpoint paths
	paths := make(map[string]*Endpoint)
	for _, endpoint := range doc.Endpoints {
		key := endpoint.Method + " " + endpoint.Path
		if existing := paths[key]; existing != nil {
			v.addError("error", "duplicate endpoint: "+key, endpoint.LineNumber)
		}
		paths[key] = endpoint
	}

	return nil
}

func (v *ValidationVisitor) VisitEndpoint(ctx context.Context, endpoint *Endpoint) error {
	v.currentPath = "endpoint[" + endpoint.Method + " " + endpoint.Path + "]"

	// Validate HTTP method
	validMethods := []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	valid := false
	for _, method := range validMethods {
		if endpoint.Method == method {
			valid = true
			break
		}
	}
	if !valid {
		v.addError("error", "invalid HTTP method: "+endpoint.Method, endpoint.LineNumber)
	}

	// Validate path
	if !strings.HasPrefix(endpoint.Path, "/") {
		v.addError("error", "path must start with /", endpoint.LineNumber)
	}

	// Check for required descriptions in strict mode
	if v.strictMode && endpoint.Description == "" {
		v.addError("warning", "endpoint description is recommended", endpoint.LineNumber)
	}

	return nil
}

func (v *ValidationVisitor) VisitParameter(ctx context.Context, parameter *Parameter) error {
	v.currentPath += ".parameter[" + parameter.Name + "]"

	// Validate parameter location
	validLocations := []string{"query", "path", "header", "cookie"}
	valid := false
	for _, location := range validLocations {
		if parameter.In == location {
			valid = true
			break
		}
	}
	if !valid {
		v.addError("error", "invalid parameter location: "+parameter.In, parameter.LineNumber)
	}

	// Path parameters must be required
	if parameter.In == "path" && !parameter.Required {
		v.addError("error", "path parameters must be required", parameter.LineNumber)
	}

	return nil
}

func (v *ValidationVisitor) VisitSchema(ctx context.Context, schema *Schema) error {
	v.currentPath += ".schema"

	// Validate schema references
	if schema.Ref != "" && schema.Type != "" {
		v.addError("error", "schema cannot have both $ref and type", schema.LineNumber)
	}

	// Check for circular references (simplified check)
	if schema.Ref != "" && strings.Contains(schema.Ref, v.currentPath) {
		v.addError("warning", "potential circular reference detected", schema.LineNumber)
	}

	return nil
}

func (v *ValidationVisitor) addError(errorType, message string, lineNumber int) {
	_ = errorType // TODO: Use errorType to determine severity level
	v.errors = append(v.errors, errors.NewError(errors.ErrorTypeValidation, message).
		AtLine(lineNumber).
		WithContext(v.currentPath).
		Build())
}

func (v *ValidationVisitor) GetErrors() []*errors.ParseError {
	return v.errors
}

// StatisticsVisitor collects statistics about the AST
type StatisticsVisitor struct {
	BaseVisitor
	Stats DocumentStatistics
}

type DocumentStatistics struct {
	TotalEndpoints    int
	EndpointsByMethod map[string]int
	TotalParameters   int
	ParametersByType  map[string]int
	TotalSchemas      int
	SchemasByType     map[string]int
	MaxSchemaDepth    int
	HasFrontmatter    bool
	TotalComponents   int
	AveragePathLength float64
}

func NewStatisticsVisitor() *StatisticsVisitor {
	return &StatisticsVisitor{
		Stats: DocumentStatistics{
			EndpointsByMethod: make(map[string]int),
			ParametersByType:  make(map[string]int),
			SchemasByType:     make(map[string]int),
		},
	}
}

func (v *StatisticsVisitor) VisitDocument(ctx context.Context, doc *Document) error {
	v.Stats.TotalEndpoints = len(doc.Endpoints)
	v.Stats.TotalComponents = len(doc.Components)
	v.Stats.HasFrontmatter = doc.Frontmatter != nil

	// Calculate average path length
	if len(doc.Endpoints) > 0 {
		totalLength := 0
		for _, endpoint := range doc.Endpoints {
			totalLength += len(endpoint.Path)
		}
		v.Stats.AveragePathLength = float64(totalLength) / float64(len(doc.Endpoints))
	}

	return nil
}

func (v *StatisticsVisitor) VisitEndpoint(ctx context.Context, endpoint *Endpoint) error {
	v.Stats.EndpointsByMethod[endpoint.Method]++
	return nil
}

func (v *StatisticsVisitor) VisitParameter(ctx context.Context, parameter *Parameter) error {
	v.Stats.TotalParameters++
	v.Stats.ParametersByType[parameter.Type]++
	return nil
}

func (v *StatisticsVisitor) VisitSchema(ctx context.Context, schema *Schema) error {
	v.Stats.TotalSchemas++
	if schema.Type != "" {
		v.Stats.SchemasByType[schema.Type]++
	}

	// Calculate depth (simplified)
	depth := v.calculateSchemaDepth(schema, 0)
	if depth > v.Stats.MaxSchemaDepth {
		v.Stats.MaxSchemaDepth = depth
	}

	return nil
}

func (v *StatisticsVisitor) calculateSchemaDepth(schema *Schema, currentDepth int) int {
	maxDepth := currentDepth

	for _, prop := range schema.Properties {
		depth := v.calculateSchemaDepth(prop, currentDepth+1)
		if depth > maxDepth {
			maxDepth = depth
		}
	}

	if schema.Items != nil {
		depth := v.calculateSchemaDepth(schema.Items, currentDepth+1)
		if depth > maxDepth {
			maxDepth = depth
		}
	}

	return maxDepth
}

// TransformVisitor modifies the AST during traversal
type TransformVisitor struct {
	BaseVisitor
	transformations []func(interface{}) interface{}
}

func NewTransformVisitor() *TransformVisitor {
	return &TransformVisitor{
		transformations: []func(interface{}) interface{}{},
	}
}

func (v *TransformVisitor) AddTransformation(transform func(interface{}) interface{}) {
	v.transformations = append(v.transformations, transform)
}

func (v *TransformVisitor) VisitEndpoint(ctx context.Context, endpoint *Endpoint) error {
	// Example transformation: normalize paths
	endpoint.Path = strings.ToLower(endpoint.Path)

	// Apply custom transformations
	for _, transform := range v.transformations {
		if result := transform(endpoint); result != nil {
			if transformed, ok := result.(*Endpoint); ok {
				*endpoint = *transformed
			}
		}
	}

	return nil
}

// Helper functions for using visitors

// ValidateDocument validates a document using the validation visitor
func ValidateDocument(ctx context.Context, doc *Document, strictMode bool) []*errors.ParseError {
	visitor := NewValidationVisitor(strictMode)
	if err := doc.Accept(ctx, visitor); err != nil {
		// Add the error to the visitor's error collection
		visitor.errors = append(visitor.errors, errors.NewError(errors.ErrorTypeValidation, err.Error()).Build())
	}
	return visitor.GetErrors()
}

// GetDocumentStatistics collects statistics about a document
func GetDocumentStatistics(ctx context.Context, doc *Document) DocumentStatistics {
	visitor := NewStatisticsVisitor()
	_ = doc.Accept(ctx, visitor) // Ignore errors for statistics collection
	return visitor.Stats
}

// TransformDocument applies transformations to a document
func TransformDocument(ctx context.Context, doc *Document, transforms ...func(interface{}) interface{}) error {
	visitor := NewTransformVisitor()
	for _, transform := range transforms {
		visitor.AddTransformation(transform)
	}
	return doc.Accept(ctx, visitor)
}
