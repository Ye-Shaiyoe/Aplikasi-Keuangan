package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger is a Gin middleware that writes structured request logs to stdout.
// Each log line is a pipe-delimited string containing:
//
//	timestamp | request_id | method | path | status | latency | client_ip | error (if any)
//
// It is intentionally dependency-free — no external logging library required.
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Build log entry after handler chain completes
		latency := time.Since(start)
		status := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		reqID := GetRequestID(c)
		errMsg := c.Errors.ByType(gin.ErrorTypePrivate).String()

		if query != "" {
			path = path + "?" + query
		}

		// Colour the status code for terminal readability
		statusColor := statusColorCode(status)
		resetColor := "\033[0m"

		line := fmt.Sprintf(
			"[GIN] %s | %s%3d%s | %13v | %-15s | %s %-7s %s",
			time.Now().Format("2006/01/02 - 15:04:05"),
			statusColor, status, resetColor,
			latency,
			clientIP,
			reqID,
			method,
			path,
		)
		if errMsg != "" {
			line += " | ERR: " + errMsg
		}
		fmt.Println(line)
	}
}

// statusColorCode returns an ANSI escape code for the given HTTP status.
func statusColorCode(status int) string {
	switch {
	case status >= 200 && status < 300:
		return "\033[32m" // green
	case status >= 300 && status < 400:
		return "\033[36m" // cyan
	case status >= 400 && status < 500:
		return "\033[33m" // yellow
	default:
		return "\033[31m" // red
	}
}
