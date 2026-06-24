package handler

import (
	"net/http"

	"github.com/akrom/finance-backend/internal/middleware"
	"github.com/akrom/finance-backend/internal/service"
	"github.com/gin-gonic/gin"
)

func GetAdvancedAnalytics(c *gin.Context) {
	userID := middleware.GetUserID(c)

	analytics, err := service.GetAdvancedAnalytics(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, analytics)
}
