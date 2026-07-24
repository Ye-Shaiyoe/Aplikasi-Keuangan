package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// tokenBucket is a simple token-bucket rate limiter for a single client.
type tokenBucket struct {
	tokens     float64
	maxTokens  float64
	refillRate float64 // tokens per second
	lastRefill time.Time
	mu         sync.Mutex
}

// allow returns true and consumes one token if the bucket has capacity.
func (b *tokenBucket) allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(b.lastRefill).Seconds()
	b.tokens += elapsed * b.refillRate
	if b.tokens > b.maxTokens {
		b.tokens = b.maxTokens
	}
	b.lastRefill = now

	if b.tokens >= 1 {
		b.tokens--
		return true
	}
	return false
}

// rateLimiterStore manages per-IP token buckets.
type rateLimiterStore struct {
	mu      sync.RWMutex
	buckets sync.Map // key: IP string → *tokenBucket
	// Cleanup parameters
	cleanupInterval time.Duration
	lastCleanup     time.Time
	maxBuckets      int
	// Bucket parameters
	maxTokens  float64
	refillRate float64
}

// getBucket retrieves or creates a token bucket for the given key.
func (s *rateLimiterStore) getBucket(key string) *tokenBucket {
	if v, ok := s.buckets.Load(key); ok {
		return v.(*tokenBucket)
	}
	b := &tokenBucket{
		tokens:     s.maxTokens,
		maxTokens:  s.maxTokens,
		refillRate: s.refillRate,
		lastRefill: time.Now(),
	}
	actual, _ := s.buckets.LoadOrStore(key, b)
	return actual.(*tokenBucket)
}

// maybeCleanup removes buckets that have not been seen for a while.
func (s *rateLimiterStore) maybeCleanup() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if time.Since(s.lastCleanup) < s.cleanupInterval {
		return
	}
	s.lastCleanup = time.Now()
	s.buckets.Range(func(k, v interface{}) bool {
		b := v.(*tokenBucket)
		b.mu.Lock()
		idle := time.Since(b.lastRefill) > s.cleanupInterval
		b.mu.Unlock()
		if idle {
			s.buckets.Delete(k)
		}
		return true
	})
}

// RateLimiter returns a Gin middleware that enforces a per-IP token-bucket rate
// limit using only the Go standard library.
//
//   requestsPerSecond – sustained request rate allowed per IP
//   burstSize         – maximum burst capacity (token bucket ceiling)
func RateLimiter(requestsPerSecond float64, burstSize int) gin.HandlerFunc {
	store := &rateLimiterStore{
		maxTokens:       float64(burstSize),
		refillRate:      requestsPerSecond,
		cleanupInterval: 5 * time.Minute,
		lastCleanup:     time.Now(),
		maxBuckets:      10_000,
	}

	return func(c *gin.Context) {
		// Periodic cleanup to prevent unbounded memory growth
		store.maybeCleanup()

		ip := c.ClientIP()
		bucket := store.getBucket(ip)

		if !bucket.allow() {
			c.Header("Retry-After", "1")
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "terlalu banyak permintaan, coba lagi sebentar",
			})
			return
		}
		c.Next()
	}
}
