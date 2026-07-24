package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akrom/finance-backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func TestRequestIDMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.RequestID())
	r.GET("/test", func(c *gin.Context) {
		reqID := middleware.GetRequestID(c)
		c.String(http.StatusOK, reqID)
	})

	// Case 1: auto-generate UUID
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", w.Code)
	}
	respID := w.Header().Get(middleware.RequestIDHeader)
	if respID == "" {
		t.Error("X-Request-ID header is missing in response")
	}
	if len(respID) != 36 { // 8-4-4-4-12 UUID format
		t.Errorf("request ID length = %d, want 36 (UUID v4)", len(respID))
	}

	// Case 2: preserve incoming ID
	customID := "custom-req-12345"
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/test", nil)
	req2.Header.Set(middleware.RequestIDHeader, customID)
	r.ServeHTTP(w2, req2)

	if w2.Header().Get(middleware.RequestIDHeader) != customID {
		t.Errorf("custom X-Request-ID not preserved, got %s", w2.Header().Get(middleware.RequestIDHeader))
	}
}
