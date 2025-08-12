package errors

import (
	"fmt"
	"strings"
)

// FormatErrors formats errors for display
func FormatErrors(errors []*ParseError) string {
	if len(errors) == 0 {
		return "No errors"
	}

	var builder strings.Builder

	// Group by severity
	var fatal, errs, warnings, infos []*ParseError
	for _, err := range errors {
		switch err.Severity {
		case SeverityFatal:
			fatal = append(fatal, err)
		case SeverityError:
			errs = append(errs, err)
		case SeverityWarning:
			warnings = append(warnings, err)
		case SeverityInfo:
			infos = append(infos, err)
		}
	}

	// Format each group
	if len(fatal) > 0 {
		builder.WriteString("FATAL ERRORS:\n")
		for _, err := range fatal {
			builder.WriteString(fmt.Sprintf("  %s\n", err.Error()))
		}
		builder.WriteString("\n")
	}

	if len(errs) > 0 {
		builder.WriteString("ERRORS:\n")
		for _, err := range errs {
			builder.WriteString(fmt.Sprintf("  %s\n", err.Error()))
		}
		builder.WriteString("\n")
	}

	if len(warnings) > 0 {
		builder.WriteString("WARNINGS:\n")
		for _, err := range warnings {
			builder.WriteString(fmt.Sprintf("  %s\n", err.Error()))
		}
		builder.WriteString("\n")
	}

	if len(infos) > 0 {
		builder.WriteString("INFO:\n")
		for _, err := range infos {
			builder.WriteString(fmt.Sprintf("  %s\n", err.Error()))
		}
	}

	return strings.TrimSpace(builder.String())
}

// FilterErrors filters errors by type or severity
func FilterErrors(errors []*ParseError, filter func(*ParseError) bool) []*ParseError {
	var filtered []*ParseError
	for _, err := range errors {
		if filter(err) {
			filtered = append(filtered, err)
		}
	}
	return filtered
}

// FilterBySeverity filters errors by severity
func FilterBySeverity(errors []*ParseError, severity Severity) []*ParseError {
	return FilterErrors(errors, func(err *ParseError) bool {
		return err.Severity == severity
	})
}

// FilterByType filters errors by type
func FilterByType(errors []*ParseError, errorType ErrorType) []*ParseError {
	return FilterErrors(errors, func(err *ParseError) bool {
		return err.Type == errorType
	})
}
