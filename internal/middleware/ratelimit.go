package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter represents a simple in-memory rate limiter
type RateLimiter struct {
	visitors map[string]*Visitor
	mu       sync.RWMutex
	rate     int           // requests per window
	window   time.Duration // time window
}

// Visitor represents a client making requests
type Visitor struct {
	limiter  *TokenBucket
	lastSeen time.Time
}

// TokenBucket implements token bucket algorithm
type TokenBucket struct {
	tokens    int
	capacity  int
	refillRate int
	lastRefill time.Time
	mu        sync.Mutex
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(rate int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*Visitor),
		rate:     rate,
		window:   window,
	}

	// Cleanup expired visitors every 5 minutes
	go rl.cleanupVisitors()

	return rl
}

// Allow checks if a request should be allowed
func (rl *RateLimiter) Allow(clientIP string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	visitor, exists := rl.visitors[clientIP]
	if !exists {
		visitor = &Visitor{
			limiter:  NewTokenBucket(rl.rate, rl.rate),
			lastSeen: time.Now(),
		}
		rl.visitors[clientIP] = visitor
	}

	visitor.lastSeen = time.Now()
	return visitor.limiter.Allow()
}

// NewTokenBucket creates a new token bucket
func NewTokenBucket(capacity, refillRate int) *TokenBucket {
	return &TokenBucket{
		tokens:     capacity,
		capacity:   capacity,
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

// Allow checks if a token can be consumed
func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.refill()

	if tb.tokens > 0 {
		tb.tokens--
		return true
	}

	return false
}

// refill adds tokens based on elapsed time
func (tb *TokenBucket) refill() {
	now := time.Now()
	elapsed := now.Sub(tb.lastRefill)
	tokensToAdd := int(elapsed.Seconds()) * tb.refillRate / 60 // refill per minute

	if tokensToAdd > 0 {
		tb.tokens += tokensToAdd
		if tb.tokens > tb.capacity {
			tb.tokens = tb.capacity
		}
		tb.lastRefill = now
	}
}

// cleanupVisitors removes old visitors
func (rl *RateLimiter) cleanupVisitors() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		for ip, visitor := range rl.visitors {
			if time.Since(visitor.lastSeen) > 10*time.Minute {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// RateLimitMiddleware creates a rate limiting middleware
func RateLimitMiddleware(rate int, window time.Duration) gin.HandlerFunc {
	limiter := NewRateLimiter(rate, window)

	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		if !limiter.Allow(clientIP) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "Rate limit exceeded",
				"message": "Too many requests. Please try again later.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GlobalRateLimitMiddleware applies a global rate limit
func GlobalRateLimitMiddleware() gin.HandlerFunc {
	return RateLimitMiddleware(100, time.Minute) // 100 requests per minute
}

// AuthRateLimitMiddleware applies stricter rate limiting for auth endpoints
func AuthRateLimitMiddleware() gin.HandlerFunc {
	return RateLimitMiddleware(5, time.Minute) // 5 requests per minute for auth
}

// ApiRateLimitMiddleware applies standard rate limiting for API endpoints
func ApiRateLimitMiddleware() gin.HandlerFunc {
	return RateLimitMiddleware(60, time.Minute) // 60 requests per minute for API
}
