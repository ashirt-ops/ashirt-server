package middleware

import (
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// RateLimiter provides IP-based rate limiting middleware
type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
}

// NewRateLimiter creates a new rate limiter with the specified requests per second and burst size
// Example: NewRateLimiter(5, 10) allows 5 requests per second with a burst of 10
func NewRateLimiter(requestsPerSecond float64, burst int) *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     rate.Limit(requestsPerSecond),
		burst:    burst,
	}
}

// getLimiter returns the rate limiter for a given IP address
func (rl *RateLimiter) getLimiter(ip string) *rate.Limiter {
	rl.mu.RLock()
	limiter, exists := rl.limiters[ip]
	rl.mu.RUnlock()

	if !exists {
		rl.mu.Lock()
		limiter = rate.NewLimiter(rl.rate, rl.burst)
		rl.limiters[ip] = limiter
		rl.mu.Unlock()
	}

	return limiter
}

// getIP extracts the client IP from the request, checking X-Forwarded-For first
func getIP(r *http.Request) string {
	// Check X-Forwarded-For header (set by reverse proxies)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// X-Forwarded-For can contain multiple IPs, use the first one
		return xff
	}
	// Fall back to RemoteAddr
	return r.RemoteAddr
}

// Limit is the middleware handler that enforces rate limiting
func (rl *RateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := getIP(r)
		limiter := rl.getLimiter(ip)

		if !limiter.Allow() {
			http.Error(w, "Rate limit exceeded. Please try again later.", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// cleanup removes limiters that haven't been used recently (helps prevent memory leaks)
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		// Clear all limiters periodically - they'll be recreated on next request
		rl.limiters = make(map[string]*rate.Limiter)
		rl.mu.Unlock()
	}
}

// StartCleanup begins the cleanup goroutine to prevent memory leaks
func (rl *RateLimiter) StartCleanup() {
	go rl.cleanup()
}
