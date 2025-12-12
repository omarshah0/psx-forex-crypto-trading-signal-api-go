package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/omarshah0/rest-api-with-social-auth/internal/database"
	"github.com/omarshah0/rest-api-with-social-auth/internal/utils"
)

type RateLimitMiddleware struct {
	redisDB *database.RedisDB
}

func NewRateLimitMiddleware(redisDB *database.RedisDB) *RateLimitMiddleware {
	return &RateLimitMiddleware{redisDB: redisDB}
}

// RateLimit middleware implements rate limiting per IP address
func (m *RateLimitMiddleware) RateLimit(requestsPerMinute int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get client IP
			ip := r.RemoteAddr
			if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
				ip = forwarded
			}

			// Rate limit key
			key := fmt.Sprintf("rate_limit:%s", ip)
			ctx := context.Background()

			// Check current count
			count, err := m.redisDB.Client.Incr(ctx, key).Result()
			if err != nil {
				// On error, allow the request but log it
				next.ServeHTTP(w, r)
				return
			}

			// Set expiry on first request
			if count == 1 {
				m.redisDB.Client.Expire(ctx, key, time.Minute)
			}

			// Check if limit exceeded
			if count > int64(requestsPerMinute) {
				utils.SendError(w, http.StatusTooManyRequests, "rate_limit_exceeded", "Too many requests. Please try again later.")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

