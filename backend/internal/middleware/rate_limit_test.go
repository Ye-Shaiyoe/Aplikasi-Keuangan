package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akrom/finance-backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func TestRateLimiterMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	// Allow 2 requests burst, 1 req/sec
	r.Use(middleware.RateLimiter(1, 2))
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// 1st request — allowed
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/ping", nil)
	r.ServeHTTP(w1, req1)
	if w1.Code != http.StatusOK {
		t.Errorf("req 1 status = %d, want 200", w1.Code)
	}

	// 2nd request — allowed (burst=2)
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/ping", nil)
	r.ServeHTTP(w2, req2)
	if w2.Code != http.StatusOK {
		t.Errorf("req 2 status = %d, want 200", w2.Code)
	}

	// 3rd request — should be rate limited (429)
	w3 := httptest.NewRecorder()
	req3, _ := http.NewRequest("GET", "/ping", nil)
	r.ServeHTTP(w3, req3)
	if w3.Code != http.StatusTooManyRequests {
		t.Errorf("req 3 status = %d, want 429", w3.Code)
	}
}
