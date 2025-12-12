package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/omarshah0/rest-api-with-social-auth/internal/config"
	"github.com/omarshah0/rest-api-with-social-auth/internal/repositories"
	"github.com/omarshah0/rest-api-with-social-auth/internal/utils"
)

type LoggingMiddleware struct {
	logRepo       *repositories.LogRepository
	sensitiveKeys []string
	ignoredKeys   []string
}

func NewLoggingMiddleware(logRepo *repositories.LogRepository, loggingConfig config.LoggingConfig) *LoggingMiddleware {
	return &LoggingMiddleware{
		logRepo:       logRepo,
		sensitiveKeys: loggingConfig.SensitiveKeys,
		ignoredKeys:   loggingConfig.IgnoredKeys,
	}
}

// responseWriter wraps http.ResponseWriter to capture response data
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
		body:           &bytes.Buffer{},
	}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.body.Write(b)
	return rw.ResponseWriter.Write(b)
}

// LogRequests middleware logs all requests to MongoDB
func (m *LoggingMiddleware) LogRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		// Read request body
		var requestBody interface{}
		if r.Body != nil && (r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH") {
			bodyBytes, err := io.ReadAll(r.Body)
			if err == nil && len(bodyBytes) > 0 {
				// Restore the body for handlers
				r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

				// Parse JSON body
				var bodyMap map[string]interface{}
				if err := json.Unmarshal(bodyBytes, &bodyMap); err == nil {
					requestBody = utils.MaskSensitiveData(bodyMap, m.sensitiveKeys, m.ignoredKeys)
				}
			}
		}

		// Capture response
		rw := newResponseWriter(w)

		// Call next handler
		next.ServeHTTP(rw, r)

		// Calculate duration
		duration := time.Since(startTime)

		// Get user ID from context if available
		var userID *int64
		if uid, ok := GetUserIDFromContext(r.Context()); ok {
			userID = &uid
		}

		// Parse response body
		var responseBody interface{}
		if rw.body.Len() > 0 {
			var bodyMap map[string]interface{}
			if err := json.Unmarshal(rw.body.Bytes(), &bodyMap); err == nil {
				responseBody = utils.MaskSensitiveData(bodyMap, m.sensitiveKeys, m.ignoredKeys)
			}
		}

		// Get IP address
		ip := r.RemoteAddr
		if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
			ip = forwarded
		}

		// Create request log
		requestLog := &repositories.RequestLog{
			Timestamp:   startTime,
			Method:      r.Method,
			Path:        r.URL.Path,
			Headers:     m.maskHeaders(r.Header),
			Body:        requestBody,
			QueryParams: m.convertURLValues(r.URL.Query()),
			Response: repositories.ResponseLog{
				Status:  rw.statusCode,
				Headers: m.maskHeaders(rw.Header()),
				Body:    responseBody,
			},
			UserID:    userID,
			IPAddress: ip,
			Duration:  duration.Milliseconds(),
		}

		// Log to MongoDB asynchronously with background context
		go func() {
			// Use background context instead of request context to avoid cancellation
			ctx := context.Background()
			if err := m.logRepo.Create(ctx, requestLog); err != nil {
				log.Printf("Failed to log request: %v", err)
			}
		}()
	})
}

// maskHeaders masks sensitive headers
func (m *LoggingMiddleware) maskHeaders(headers http.Header) map[string]interface{} {
	headerMap := make(map[string]interface{})
	for key, values := range headers {
		if len(values) == 1 {
			headerMap[key] = values[0]
		} else {
			headerMap[key] = values
		}
	}

	masked := utils.MaskSensitiveData(headerMap, m.sensitiveKeys, m.ignoredKeys)
	if maskedMap, ok := masked.(map[string]interface{}); ok {
		return maskedMap
	}
	return headerMap
}

// convertURLValues converts url.Values to map[string]interface{}
func (m *LoggingMiddleware) convertURLValues(values map[string][]string) map[string]interface{} {
	result := make(map[string]interface{})
	for key, vals := range values {
		if len(vals) == 1 {
			result[key] = vals[0]
		} else {
			result[key] = vals
		}
	}
	return result
}
