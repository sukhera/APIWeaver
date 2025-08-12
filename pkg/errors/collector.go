package errors

import (
	"errors"
	"fmt"
	"strings"
)

// ErrorCollector collects and manages parsing errors
type ErrorCollector struct {
	errors    []*ParseError
	warnings  []*ParseError
	maxErrors int
	context   string
}

// NewErrorCollector creates a new error collector
func NewErrorCollector(maxErrors int) *ErrorCollector {
	return &ErrorCollector{
		errors:    []*ParseError{},
		warnings:  []*ParseError{},
		maxErrors: maxErrors,
	}
}

// SetContext sets the current parsing context
func (c *ErrorCollector) SetContext(context string) {
	c.context = context
}

// Add adds an error to the collection
func (c *ErrorCollector) Add(err *ParseError) {
	if err == nil {
		return
	}

	// Set context if not already set
	if err.Context == "" && c.context != "" {
		err.Context = c.context
	}

	if err.IsError() {
		c.errors = append(c.errors, err)

		// Stop collecting if we hit the max error limit
		if c.maxErrors > 0 && len(c.errors) >= c.maxErrors {
			fatalErr := NewFatal(ErrorTypeValidation,
				fmt.Sprintf("too many errors (limit: %d)", c.maxErrors)).Build()
			c.errors = append(c.errors, fatalErr)
		}
	} else if err.IsWarning() {
		c.warnings = append(c.warnings, err)
	}
}

// AddMultiple adds multiple errors
func (c *ErrorCollector) AddMultiple(errors []*ParseError) {
	for _, err := range errors {
		c.Add(err)
	}
}

// HasErrors returns true if there are any errors
func (c *ErrorCollector) HasErrors() bool {
	return len(c.errors) > 0
}

// HasWarnings returns true if there are any warnings
func (c *ErrorCollector) HasWarnings() bool {
	return len(c.warnings) > 0
}

// HasFatalErrors returns true if there are any fatal errors
func (c *ErrorCollector) HasFatalErrors() bool {
	for _, err := range c.errors {
		if err.IsFatal() {
			return true
		}
	}
	return false
}

// GetErrors returns all errors
func (c *ErrorCollector) GetErrors() []*ParseError {
	return c.errors
}

// GetWarnings returns all warnings
func (c *ErrorCollector) GetWarnings() []*ParseError {
	return c.warnings
}

// GetAll returns all errors and warnings
func (c *ErrorCollector) GetAll() []*ParseError {
	all := make([]*ParseError, 0, len(c.errors)+len(c.warnings))
	all = append(all, c.errors...)
	all = append(all, c.warnings...)
	return all
}

// Clear clears all collected errors and warnings
func (c *ErrorCollector) Clear() {
	c.errors = []*ParseError{}
	c.warnings = []*ParseError{}
}

// ToError converts the collected errors to a standard error
func (c *ErrorCollector) ToError() error {
	if !c.HasErrors() {
		return nil
	}

	if len(c.errors) == 1 {
		return c.errors[0]
	}

	// Create a multi-error
	messages := make([]string, len(c.errors))
	for i, err := range c.errors {
		messages[i] = err.Error()
	}

	return errors.New(strings.Join(messages, "; "))
}
