package middleware

import (
	"net/http"
	"sync"
	"time"
)

// RateLimiter implements a simple rate limiting middleware
type RateLimiter struct {
	requests map[string][]time.Time
	mu       sync.RWMutex
	limit    int
	window   time.Duration
}

// NewRateLimiter creates a new rate limiter
// limit: maximum number of requests
// window: time window for the limit (e.g., 1 minute)
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
	
	// Cleanup old entries periodically
	go rl.cleanup()
	
	return rl
}

// Middleware returns an HTTP middleware that enforces rate limiting
func (rl *RateLimiter) Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Use IP address as identifier
		ip := r.RemoteAddr
		
		if !rl.allow(ip) {
			http.Error(w, "Too many requests. Please try again later.", http.StatusTooManyRequests)
			return
		}
		
		next(w, r)
	}
}

// allow checks if a request from the given IP should be allowed
func (rl *RateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	
	now := time.Now()
	windowStart := now.Add(-rl.window)
	
	// Get existing requests for this IP
	requests := rl.requests[ip]
	
	// Filter out old requests outside the window
	validRequests := []time.Time{}
	for _, reqTime := range requests {
		if reqTime.After(windowStart) {
			validRequests = append(validRequests, reqTime)
		}
	}
	
	// Check if limit is exceeded
	if len(validRequests) >= rl.limit {
		return false
	}
	
	// Add current request
	validRequests = append(validRequests, now)
	rl.requests[ip] = validRequests
	
	return true
}

// cleanup periodically removes old entries to prevent memory leak
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.window)
	defer ticker.Stop()
	
	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		windowStart := now.Add(-rl.window)
		
		for ip, requests := range rl.requests {
			validRequests := []time.Time{}
			for _, reqTime := range requests {
				if reqTime.After(windowStart) {
					validRequests = append(validRequests, reqTime)
				}
			}
			
			if len(validRequests) == 0 {
				delete(rl.requests, ip)
			} else {
				rl.requests[ip] = validRequests
			}
		}
		rl.mu.Unlock()
	}
}
