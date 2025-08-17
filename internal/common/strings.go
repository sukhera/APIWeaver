package common

import (
	"regexp"
	"strings"
	"unicode"
)

// ToCamelCase converts a string to camelCase
func ToCamelCase(s string) string {
	if s == "" {
		return ""
	}
	
	words := splitWords(s)
	if len(words) == 0 {
		return ""
	}
	
	result := strings.ToLower(words[0])
	for i := 1; i < len(words); i++ {
		result += strings.Title(strings.ToLower(words[i]))
	}
	
	return result
}

// ToPascalCase converts a string to PascalCase
func ToPascalCase(s string) string {
	if s == "" {
		return ""
	}
	
	words := splitWords(s)
	var result strings.Builder
	
	for _, word := range words {
		if word != "" {
			result.WriteString(strings.Title(strings.ToLower(word)))
		}
	}
	
	return result.String()
}

// ToSnakeCase converts a string to snake_case
func ToSnakeCase(s string) string {
	if s == "" {
		return ""
	}
	
	words := splitWords(s)
	var result []string
	
	for _, word := range words {
		if word != "" {
			result = append(result, strings.ToLower(word))
		}
	}
	
	return strings.Join(result, "_")
}

// ToKebabCase converts a string to kebab-case
func ToKebabCase(s string) string {
	if s == "" {
		return ""
	}
	
	words := splitWords(s)
	var result []string
	
	for _, word := range words {
		if word != "" {
			result = append(result, strings.ToLower(word))
		}
	}
	
	return strings.Join(result, "-")
}

// splitWords splits a string into words based on various delimiters
func splitWords(s string) []string {
	// Handle camelCase and PascalCase
	re := regexp.MustCompile(`([a-z])([A-Z])`)
	s = re.ReplaceAllString(s, "${1} ${2}")
	
	// Split on common delimiters
	words := regexp.MustCompile(`[^a-zA-Z0-9]+`).Split(s, -1)
	
	// Filter empty strings
	var result []string
	for _, word := range words {
		if word != "" {
			result = append(result, word)
		}
	}
	
	return result
}

// Truncate truncates a string to the specified length with ellipsis
func Truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	
	if maxLen <= 3 {
		return s[:maxLen]
	}
	
	return s[:maxLen-3] + "..."
}

// IsEmpty checks if a string is empty or contains only whitespace
func IsEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

// IsValidIdentifier checks if a string is a valid identifier (alphanumeric + underscore, not starting with digit)
func IsValidIdentifier(s string) bool {
	if s == "" {
		return false
	}
	
	// Must start with letter or underscore
	if !unicode.IsLetter(rune(s[0])) && s[0] != '_' {
		return false
	}
	
	// Rest must be alphanumeric or underscore
	for _, r := range s[1:] {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			return false
		}
	}
	
	return true
}

// ContainsAny checks if string contains any of the given substrings
func ContainsAny(s string, substrings []string) bool {
	for _, sub := range substrings {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}

// RemoveEmptyLines removes empty lines from a multi-line string
func RemoveEmptyLines(s string) string {
	lines := strings.Split(s, "\n")
	var result []string
	
	for _, line := range lines {
		if !IsEmpty(line) {
			result = append(result, line)
		}
	}
	
	return strings.Join(result, "\n")
}

// IndentLines indents each line of a multi-line string
func IndentLines(s string, indent string) string {
	lines := strings.Split(s, "\n")
	var result []string
	
	for _, line := range lines {
		if !IsEmpty(line) {
			result = append(result, indent+line)
		} else {
			result = append(result, line)
		}
	}
	
	return strings.Join(result, "\n")
}

// NormalizeWhitespace normalizes whitespace in a string
func NormalizeWhitespace(s string) string {
	// Replace multiple whitespace characters with single space
	re := regexp.MustCompile(`\s+`)
	return strings.TrimSpace(re.ReplaceAllString(s, " "))
}

// ExtractLines extracts lines from a string within a given range
func ExtractLines(s string, start, end int) string {
	lines := strings.Split(s, "\n")
	
	if start < 0 {
		start = 0
	}
	if end > len(lines) {
		end = len(lines)
	}
	if start >= end {
		return ""
	}
	
	return strings.Join(lines[start:end], "\n")
}