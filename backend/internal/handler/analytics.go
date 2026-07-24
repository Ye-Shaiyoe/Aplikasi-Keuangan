package handler

import (
	"net/http"
	"strconv"

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

func GetExpenseHeatmap(c *gin.Context) {
	userID := middleware.GetUserID(c)
	year, _ := strconv.Atoi(c.Query("year"))

	heatmap, err := service.GetExpenseHeatmap(userID, year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, heatmap)
}

func GetNetWorthTimeline(c *gin.Context) {
	userID := middleware.GetUserID(c)
	months, _ := strconv.Atoi(c.Query("months"))

	timeline, err := service.GetNetWorthTimeline(userID, months)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, timeline)
}
