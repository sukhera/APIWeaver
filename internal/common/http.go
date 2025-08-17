package common

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// HTTPError represents an HTTP error with status code and message
type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Error implements the error interface
func (e HTTPError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("HTTP %d: %s (%s)", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("HTTP %d: %s", e.Code, e.Message)
}

// NewHTTPError creates a new HTTP error
func NewHTTPError(code int, message string, details ...string) HTTPError {
	var detail string
	if len(details) > 0 {
		detail = details[0]
	}
	return HTTPError{
		Code:    code,
		Message: message,
		Details: detail,
	}
}

// Common HTTP errors
var (
	ErrBadRequest          = NewHTTPError(http.StatusBadRequest, "Bad Request")
	ErrUnauthorized        = NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	ErrForbidden           = NewHTTPError(http.StatusForbidden, "Forbidden")
	ErrNotFound            = NewHTTPError(http.StatusNotFound, "Not Found")
	ErrMethodNotAllowed    = NewHTTPError(http.StatusMethodNotAllowed, "Method Not Allowed")
	ErrRequestTimeout      = NewHTTPError(http.StatusRequestTimeout, "Request Timeout")
	ErrPayloadTooLarge     = NewHTTPError(http.StatusRequestEntityTooLarge, "Payload Too Large")
	ErrUnsupportedMediaType = NewHTTPError(http.StatusUnsupportedMediaType, "Unsupported Media Type")
	ErrTooManyRequests     = NewHTTPError(http.StatusTooManyRequests, "Too Many Requests")
	ErrInternalServer      = NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	ErrNotImplemented      = NewHTTPError(http.StatusNotImplemented, "Not Implemented")
	ErrBadGateway          = NewHTTPError(http.StatusBadGateway, "Bad Gateway")
	ErrServiceUnavailable  = NewHTTPError(http.StatusServiceUnavailable, "Service Unavailable")
	ErrGatewayTimeout      = NewHTTPError(http.StatusGatewayTimeout, "Gateway Timeout")
)

// GetClientIP extracts the client IP address from the request
func GetClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first (for proxies)
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		// Take the first IP if there are multiple
		ips := strings.Split(forwarded, ",")
		return strings.TrimSpace(ips[0])
	}

	// Check X-Real-IP header
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Fall back to RemoteAddr
	ip := r.RemoteAddr
	if colon := strings.LastIndex(ip, ":"); colon != -1 {
		ip = ip[:colon]
	}
	return ip
}

// GetUserAgent extracts and normalizes the User-Agent header
func GetUserAgent(r *http.Request) string {
	ua := r.Header.Get("User-Agent")
	if ua == "" {
		return "Unknown"
	}
	
	// Truncate very long user agents
	if len(ua) > 500 {
		ua = ua[:500] + "..."
	}
	
	return ua
}

// IsAjaxRequest checks if the request is an AJAX request
func IsAjaxRequest(r *http.Request) bool {
	return r.Header.Get("X-Requested-With") == "XMLHttpRequest"
}

// GetContentType extracts the content type from the request
func GetContentType(r *http.Request) string {
	contentType := r.Header.Get("Content-Type")
	if idx := strings.Index(contentType, ";"); idx != -1 {
		contentType = contentType[:idx]
	}
	return strings.TrimSpace(strings.ToLower(contentType))
}

// IsJSONRequest checks if the request has JSON content type
func IsJSONRequest(r *http.Request) bool {
	contentType := GetContentType(r)
	return contentType == "application/json"
}

// IsMultipartRequest checks if the request has multipart content type
func IsMultipartRequest(r *http.Request) bool {
	contentType := GetContentType(r)
	return strings.HasPrefix(contentType, "multipart/")
}

// GetRequestSize returns the size of the request body
func GetRequestSize(r *http.Request) int64 {
	if r.ContentLength >= 0 {
		return r.ContentLength
	}
	return 0
}

// SetCacheHeaders sets cache-related headers
func SetCacheHeaders(w http.ResponseWriter, maxAge time.Duration) {
	if maxAge > 0 {
		w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", int(maxAge.Seconds())))
	} else {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
	}
}

// SetSecurityHeaders sets security-related headers
func SetSecurityHeaders(w http.ResponseWriter) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
	w.Header().Set("Content-Security-Policy", "default-src 'self'")
}

// SetCORSHeaders sets CORS headers
func SetCORSHeaders(w http.ResponseWriter, allowedOrigins []string, allowedMethods []string, allowedHeaders []string) {
	if len(allowedOrigins) == 1 && allowedOrigins[0] == "*" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	} else {
		// In a real implementation, you'd check the Origin header against allowedOrigins
		w.Header().Set("Access-Control-Allow-Origin", strings.Join(allowedOrigins, ","))
	}
	
	if len(allowedMethods) > 0 {
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(allowedMethods, ", "))
	}
	
	if len(allowedHeaders) > 0 {
		w.Header().Set("Access-Control-Allow-Headers", strings.Join(allowedHeaders, ", "))
	}
	
	w.Header().Set("Access-Control-Max-Age", "86400")
}

// ParseRange parses HTTP Range header
func ParseRange(rangeHeader string, size int64) (start, end int64, err error) {
	if !strings.HasPrefix(rangeHeader, "bytes=") {
		return 0, 0, fmt.Errorf("invalid range header")
	}
	
	rangeSpec := rangeHeader[6:] // Remove "bytes="
	parts := strings.Split(rangeSpec, "-")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid range format")
	}
	
	if parts[0] == "" && parts[1] == "" {
		return 0, 0, fmt.Errorf("invalid range values")
	}
	
	if parts[0] == "" {
		// Suffix range (-500)
		suffix, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return 0, 0, err
		}
		start = size - suffix
		end = size - 1
	} else if parts[1] == "" {
		// Start range (500-)
		start, err = strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return 0, 0, err
		}
		end = size - 1
	} else {
		// Full range (500-999)
		start, err = strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return 0, 0, err
		}
		end, err = strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return 0, 0, err
		}
	}
	
	// Validate range
	if start < 0 || end >= size || start > end {
		return 0, 0, fmt.Errorf("invalid range values")
	}
	
	return start, end, nil
}

// WithTimeout adds a timeout to an HTTP request
func WithTimeout(r *http.Request, timeout time.Duration) (*http.Request, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(r.Context(), timeout)
	return r.WithContext(ctx), cancel
}

// ResponseWriter wraps http.ResponseWriter to capture status code and size
type ResponseWriter struct {
	http.ResponseWriter
	StatusCode int
	Size       int
}

// NewResponseWriter creates a new ResponseWriter
func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		StatusCode:     200, // Default status code
	}
}

// WriteHeader captures the status code
func (rw *ResponseWriter) WriteHeader(code int) {
	rw.StatusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Write captures the response size
func (rw *ResponseWriter) Write(data []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(data)
	rw.Size += size
	return size, err
}