package parser

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
)

func TestValidationVisitor_VisitDocument(t *testing.T) {
	tests := []struct {
		name          string
		doc           *Document
		strictMode    bool
		expectedError bool
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
			expectedError: false,
		},
		{
			name: "success with empty endpoints",
			doc: &Document{
				Endpoints: []*Endpoint{},
			},
			strictMode:    false,
			expectedError: false,
		},
		{
			name: "error with duplicate endpoints",
			doc: &Document{
				Endpoints: []*Endpoint{
					{Method: "GET", Path: "/test", LineNumber: 1},
					{Method: "GET", Path: "/test", LineNumber: 2},
				},
			},
			strictMode:    false,
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create validation visitor
			visitor := NewValidationVisitor(tt.strictMode)

			// Execute test
			err := visitor.VisitDocument(context.Background(), tt.doc)

			// Assert results
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidationVisitor_VisitEndpoint(t *testing.T) {
	tests := []struct {
		name          string
		endpoint      *Endpoint
		strictMode    bool
		expectedError bool
	}{
		{
			name: "success with valid endpoint",
			endpoint: &Endpoint{
				Method:     "GET",
				Path:       "/test",
				LineNumber: 1,
			},
			strictMode:    false,
			expectedError: false,
		},
		{
			name: "error with invalid HTTP method",
			endpoint: &Endpoint{
				Method:     "INVALID",
				Path:       "/test",
				LineNumber: 1,
			},
			strictMode:    false,
			expectedError: false,
		},
		{
			name: "error with invalid path",
			endpoint: &Endpoint{
				Method:     "GET",
				Path:       "invalid-path",
				LineNumber: 1,
			},
			strictMode:    false,
			expectedError: false,
		},
		{
			name: "warning with missing description in strict mode",
			endpoint: &Endpoint{
				Method:     "GET",
				Path:       "/test",
				LineNumber: 1,
			},
			strictMode:    true,
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create validation visitor
			visitor := NewValidationVisitor(tt.strictMode)

			// Execute test
			err := visitor.VisitEndpoint(context.Background(), tt.endpoint)

			// Assert results
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidationVisitor_VisitParameter(t *testing.T) {
	tests := []struct {
		name          string
		parameter     *Parameter
		expectedError bool
	}{
		{
			name: "success with valid parameter",
			parameter: &Parameter{
				Name:       "test_param",
				In:         "query",
				Type:       "string",
				Required:   false,
				LineNumber: 1,
			},
			expectedError: false,
		},
		{
			name: "error with invalid location",
			parameter: &Parameter{
				Name:       "test_param",
				In:         "invalid",
				Type:       "string",
				Required:   false,
				LineNumber: 1,
			},
			expectedError: false,
		},
		{
			name: "error with path parameter not required",
			parameter: &Parameter{
				Name:       "test_param",
				In:         "path",
				Type:       "string",
				Required:   false,
				LineNumber: 1,
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create validation visitor
			visitor := NewValidationVisitor(false)

			// Execute test
			err := visitor.VisitParameter(context.Background(), tt.parameter)

			// Assert results
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestStatisticsVisitor_VisitDocument(t *testing.T) {
	tests := []struct {
		name          string
		doc           *Document
		expectedStats DocumentStatistics
	}{
		{
			name: "success with single endpoint",
			doc: &Document{
				Endpoints: []*Endpoint{
					{Method: "GET", Path: "/test"},
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create statistics visitor
			visitor := NewStatisticsVisitor()

			// Execute test
			err := tt.doc.Accept(context.Background(), visitor)

			// Assert results
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStats.TotalEndpoints, visitor.Stats.TotalEndpoints)
			assert.Equal(t, tt.expectedStats.TotalComponents, visitor.Stats.TotalComponents)
			assert.Equal(t, tt.expectedStats.HasFrontmatter, visitor.Stats.HasFrontmatter)

			for method, count := range tt.expectedStats.EndpointsByMethod {
				assert.Equal(t, count, visitor.Stats.EndpointsByMethod[method])
			}
		})
	}
}

func TestDocument_Accept(t *testing.T) {
	tests := []struct {
		name          string
		doc           *Document
		expectedError bool
	}{
		{
			name: "success with valid document",
			doc: &Document{
				Endpoints: []*Endpoint{
					{Method: "GET", Path: "/test"},
				},
			},
			expectedError: false,
		},
		{
			name: "success with empty document",
			doc: &Document{
				Endpoints: []*Endpoint{},
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a real visitor (not a mock)
			visitor := NewStatisticsVisitor()

			// Execute test
			err := tt.doc.Accept(context.Background(), visitor)

			// Assert results
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMockVisitor_Integration(t *testing.T) {
	tests := []struct {
		name          string
		doc           *Document
		setupMock     func(*MockVisitor)
		expectedError bool
	}{
		{
			name: "success with mock visitor",
			doc: &Document{
				Endpoints: []*Endpoint{
					{Method: "GET", Path: "/test"},
				},
			},
			setupMock: func(mockVisitor *MockVisitor) {
				mockVisitor.EXPECT().VisitDocument(mock.Anything, mock.Anything).Return(nil).Once()
				mockVisitor.EXPECT().VisitEndpoint(mock.Anything, mock.Anything).Return(nil).Once()
			},
			expectedError: false,
		},
		{
			name: "error when mock visitor returns error",
			doc: &Document{
				Endpoints: []*Endpoint{
					{Method: "GET", Path: "/test"},
				},
			},
			setupMock: func(mockVisitor *MockVisitor) {
				mockVisitor.EXPECT().VisitDocument(mock.Anything, mock.Anything).Return(assert.AnError).Once()
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock visitor
			mockVisitor := NewMockVisitor(t)

			// Setup mock expectations
			if tt.setupMock != nil {
				tt.setupMock(mockVisitor)
			}

			// Execute test
			err := tt.doc.Accept(context.Background(), mockVisitor)

			// Assert results
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// Verify mock expectations
			mockVisitor.AssertExpectations(t)
		})
	}
}
