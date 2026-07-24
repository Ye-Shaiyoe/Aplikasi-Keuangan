package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/akrom/finance-backend/internal/middleware"
	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/service"
	"github.com/gin-gonic/gin"
)

func GetBudgets(c *gin.Context) {
	userID := middleware.GetUserID(c)
	month, _ := strconv.Atoi(c.Query("month"))
	year, _ := strconv.Atoi(c.Query("year"))
	if month < 1 || month > 12 {
		month = int(time.Now().Month())
	}
	if year < 2000 {
		year = time.Now().Year()
	}

	budgets, err := service.GetBudgets(userID, month, year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, budgets)
}

func UpsertBudget(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req model.BudgetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	budget, err := service.UpsertBudget(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, budget)
}

func DeleteBudget(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := service.DeleteBudget(userID, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "budget deleted"})
}

func GetBudgetSummary(c *gin.Context) {
	userID := middleware.GetUserID(c)
	month, _ := strconv.Atoi(c.Query("month"))
	year, _ := strconv.Atoi(c.Query("year"))
	if month < 1 || month > 12 {
		month = int(time.Now().Month())
	}
	if year < 2000 {
		year = time.Now().Year()
	}

	summary, err := service.GetBudgetSummary(userID, month, year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, summary)
}

func GetBudgetAlerts(c *gin.Context) {
	userID := middleware.GetUserID(c)
	month, _ := strconv.Atoi(c.Query("month"))
	year, _ := strconv.Atoi(c.Query("year"))
	threshold, _ := strconv.Atoi(c.Query("threshold"))

	if month < 1 || month > 12 {
		month = int(time.Now().Month())
	}
	if year < 2000 {
		year = time.Now().Year()
	}
	if threshold < 1 {
		threshold = 80
	}

	alerts, err := service.AlertBudgetOverspend(userID, month, year, threshold)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, alerts)
}

func CopyBudgetsFromLastMonth(c *gin.Context) {
	userID := middleware.GetUserID(c)
	month, _ := strconv.Atoi(c.Query("month"))
	year, _ := strconv.Atoi(c.Query("year"))

	if month < 1 || month > 12 {
		month = int(time.Now().Month())
	}
	if year < 2000 {
		year = time.Now().Year()
	}

	copied, err := service.CopyBudgetsFromLastMonth(userID, month, year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.CopyBudgetResponse{
		Copied:  copied,
		Message: "anggaran bulan lalu berhasil disalin",
	})
}
