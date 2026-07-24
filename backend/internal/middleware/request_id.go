package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"
)

const RequestIDHeader = "X-Request-ID"
const requestIDKey = "request_id"

// generateRequestID creates a cryptographically random 16-byte hex string
// (32 hex characters) without any external UUID library.
func generateRequestID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		// Fallback: use a fixed-format placeholder (should never happen in practice)
		return "00000000000000000000000000000000"
	}
	// Format as UUID v4 (8-4-4-4-12)
	b[6] = (b[6] & 0x0f) | 0x40 // version 4
	b[8] = (b[8] & 0x3f) | 0x80 // variant bits
	return hex.EncodeToString(b[:4]) + "-" +
		hex.EncodeToString(b[4:6]) + "-" +
		hex.EncodeToString(b[6:8]) + "-" +
		hex.EncodeToString(b[8:10]) + "-" +
		hex.EncodeToString(b[10:])
}

// RequestID is a Gin middleware that injects a unique request ID into every
// request. If the incoming request already carries an X-Request-ID header that
// header value is reused; otherwise a new UUID v4 is generated.
// The ID is stored in the Gin context under the key "request_id" and echoed
// back in the response header.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader(RequestIDHeader)
		if id == "" {
			id = generateRequestID()
		}
		c.Set(requestIDKey, id)
		c.Header(RequestIDHeader, id)
		c.Header("X-Content-Type-Options", "nosniff")
		c.Next()
	}
}

// GetRequestID retrieves the request ID from the Gin context.
// Returns an empty string when no request ID has been set.
func GetRequestID(c *gin.Context) string {
	if v, exists := c.Get(requestIDKey); exists {
		if id, ok := v.(string); ok {
			return id
		}
	}
	return ""
}

// RequestIDFromHeader is a convenience handler that can be used in tests or
// simple HTTP servers without Gin to extract the request ID.
func RequestIDFromHeader(r *http.Request) string {
	return r.Header.Get(RequestIDHeader)
}
