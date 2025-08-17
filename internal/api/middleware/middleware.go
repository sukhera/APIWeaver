package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/sukhera/APIWeaver/internal/common"
	"github.com/sukhera/APIWeaver/internal/config"
)

// Logging middleware logs HTTP requests and responses
func Logging(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Wrap response writer to capture status and size
			rw := common.NewResponseWriter(w)

			// Call the next handler
			next.ServeHTTP(rw, r)

			// Log the request
			duration := time.Since(start)
			logger.Info("HTTP request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", rw.StatusCode,
				"duration_ms", duration.Milliseconds(),
				"remote_addr", common.GetClientIP(r),
				"user_agent", common.GetUserAgent(r),
				"size", rw.Size,
			)
		})
	}
}

// CORS middleware adds CORS headers
func CORS(config config.CORSConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !config.Enabled {
				next.ServeHTTP(w, r)
				return
			}

			// Set CORS headers
			common.SetCORSHeaders(w, config.AllowedOrigins, config.AllowedMethods, config.AllowedHeaders)

			if config.AllowCredentials {
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			// Handle preflight requests
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Security middleware adds security headers
func Security() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Add security headers
			common.SetSecurityHeaders(w)
			next.ServeHTTP(w, r)
		})
	}
}

// Recovery middleware recovers from panics and returns a 500 error
func Recovery(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error("Panic recovered",
						"error", err,
						"method", r.Method,
						"path", r.URL.Path,
						"remote_addr", common.GetClientIP(r),
					)

					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					_, _ = w.Write([]byte(`{"success":false,"error":{"message":"Internal server error","code":500}}`))
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

// RequestID middleware adds a unique request ID to each request
func RequestID() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Generate request ID (simplified)
			requestID := fmt.Sprintf("%d", time.Now().UnixNano())
			
			// Add to response headers
			w.Header().Set("X-Request-ID", requestID)
			
			// Add to request context for logging
			ctx := r.Context()
			// In a real implementation, you'd add the request ID to context
			
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RateLimit middleware provides basic rate limiting
func RateLimit(config config.RateLimitConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !config.Enabled {
				next.ServeHTTP(w, r)
				return
			}

			// Simplified rate limiting - in real implementation would use
			// a proper rate limiter with Redis or in-memory store
			// For MVP, we'll just pass through
			next.ServeHTTP(w, r)
		})
	}
}