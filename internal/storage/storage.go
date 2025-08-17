package storage

import (
	"context"
	"time"
)

// Storage defines the interface for data persistence
type Storage interface {
	// Specs management
	SaveSpec(ctx context.Context, spec *Spec) error
	GetSpec(ctx context.Context, id string) (*Spec, error)
	ListSpecs(ctx context.Context, filters SpecFilters) ([]*Spec, error)
	DeleteSpec(ctx context.Context, id string) error

	// Conversion history
	SaveConversion(ctx context.Context, conversion *Conversion) error
	GetConversion(ctx context.Context, id string) (*Conversion, error)
	ListConversions(ctx context.Context, filters ConversionFilters) ([]*Conversion, error)

	// Examples management
	SaveExample(ctx context.Context, example *Example) error
	GetExample(ctx context.Context, id string) (*Example, error)
	ListExamples(ctx context.Context, filters ExampleFilters) ([]*Example, error)

	// Health check
	Health(ctx context.Context) error
	Close() error
}

// Spec represents a stored OpenAPI specification
type Spec struct {
	ID        string            `json:"id" bson:"_id"`
	Title     string            `json:"title" bson:"title"`
	Version   string            `json:"version" bson:"version"`
	Content   string            `json:"content" bson:"content"`
	Format    string            `json:"format" bson:"format"` // yaml, json
	Metadata  map[string]string `json:"metadata" bson:"metadata"`
	CreatedAt time.Time         `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time         `json:"updated_at" bson:"updated_at"`
}

// Conversion represents a conversion history record
type Conversion struct {
	ID             string            `json:"id" bson:"_id"`
	InputContent   string            `json:"input_content" bson:"input_content"`
	OutputContent  string            `json:"output_content" bson:"output_content"`
	InputFormat    string            `json:"input_format" bson:"input_format"`       // markdown
	OutputFormat   string            `json:"output_format" bson:"output_format"`     // yaml, json
	ProcessingTime int               `json:"processing_time" bson:"processing_time"` // milliseconds
	Success        bool              `json:"success" bson:"success"`
	Errors         []string          `json:"errors" bson:"errors"`
	Warnings       []string          `json:"warnings" bson:"warnings"`
	Metadata       map[string]string `json:"metadata" bson:"metadata"`
	CreatedAt      time.Time         `json:"created_at" bson:"created_at"`
}

// Example represents a template example
type Example struct {
	ID          string            `json:"id" bson:"_id"`
	Name        string            `json:"name" bson:"name"`
	Description string            `json:"description" bson:"description"`
	Content     string            `json:"content" bson:"content"`
	Category    string            `json:"category" bson:"category"`
	Tags        []string          `json:"tags" bson:"tags"`
	Metadata    map[string]string `json:"metadata" bson:"metadata"`
	CreatedAt   time.Time         `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at" bson:"updated_at"`
}

// SpecFilters defines filters for spec queries
type SpecFilters struct {
	Title    string
	Version  string
	Format   string
	Limit    int
	Offset   int
	SortBy   string
	SortDesc bool
}

// ConversionFilters defines filters for conversion queries
type ConversionFilters struct {
	Success  *bool
	Format   string
	DateFrom *time.Time
	DateTo   *time.Time
	Limit    int
	Offset   int
	SortBy   string
	SortDesc bool
}

// ExampleFilters defines filters for example queries
type ExampleFilters struct {
	Category string
	Tags     []string
	Limit    int
	Offset   int
	SortBy   string
	SortDesc bool
}
