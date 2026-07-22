package handler

import (
	"net/http"

	"github.com/akrom/finance-backend/internal/middleware"
	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/service"
	"github.com/gin-gonic/gin"
)

// MLPredictCategory handles POST /api/ml/predict-category
func MLPredictCategory(c *gin.Context) {
	var req model.MLPredictCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := service.MLPredictCategory(req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// MLForecast handles POST /api/ml/forecast
func MLForecast(c *gin.Context) {
	userID := middleware.GetUserID(c)

	result, err := service.MLForecastForUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// MLTrainModel handles POST /api/ml/train
func MLTrainModel(c *gin.Context) {
	userID := middleware.GetUserID(c)

	result, err := service.MLTrainModel(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// MLHealth handles GET /api/ml/health
func MLHealth(c *gin.Context) {
	result, err := service.MLHealthCheck()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error":  "ML service is not available",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}
