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

func GetTransactions(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var filter model.TransactionFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transactions, total, err := service.GetTransactions(userID, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  transactions,
		"total": total,
		"page":  filter.Page,
		"limit": filter.Limit,
	})
}

func GetTransaction(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	transaction, err := service.GetTransactionByID(id, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
		return
	}
	c.JSON(http.StatusOK, transaction)
}

func CreateTransaction(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req model.TransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := time.Parse("2006-01-02", req.Date); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format, use YYYY-MM-DD"})
		return
	}

	transaction, err := service.CreateTransaction(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, transaction)
}

func UpdateTransaction(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req model.TransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := time.Parse("2006-01-02", req.Date); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format, use YYYY-MM-DD"})
		return
	}

	transaction, err := service.UpdateTransaction(id, userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, transaction)
}

func DeleteTransaction(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := service.DeleteTransaction(id, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "transaction deleted"})
}

func GetSummary(c *gin.Context) {
	userID := middleware.GetUserID(c)
	now := time.Now()
	month := now.Month()
	year := now.Year()

	if m := c.Query("month"); m != "" {
		if val, err := strconv.Atoi(m); err == nil && val >= 1 && val <= 12 {
			month = time.Month(val)
		}
	}
	if y := c.Query("year"); y != "" {
		if val, err := strconv.Atoi(y); err == nil && val > 0 {
			year = val
		}
	}

	summary, err := service.GetSummary(userID, int(month), year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, summary)
}

func GetYearlyTrend(c *gin.Context) {
	userID := middleware.GetUserID(c)
	year := time.Now().Year()
	if y := c.Query("year"); y != "" {
		if val, err := strconv.Atoi(y); err == nil && val > 2000 {
			year = val
		}
	}

	trend, err := service.GetYearlyTrend(userID, year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, trend)
}

func GetCategoryTrend(c *gin.Context) {
	userID := middleware.GetUserID(c)
	now := time.Now()
	month := int(now.Month())
	year := now.Year()

	if m := c.Query("month"); m != "" {
		if val, err := strconv.Atoi(m); err == nil && val >= 1 && val <= 12 {
			month = val
		}
	}
	if y := c.Query("year"); y != "" {
		if val, err := strconv.Atoi(y); err == nil && val > 0 {
			year = val
		}
	}

	trend, err := service.GetCategoryTrend(userID, month, year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, trend)
}

func ExportTransactionsCSV(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var filter model.TransactionFilter
	_ = c.ShouldBindQuery(&filter)

	csvData, err := service.ExportTransactionsCSV(userID, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=transaksi.csv")
	c.Data(http.StatusOK, "text/csv; charset=utf-8", csvData)
}

func BulkDeleteTransactions(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req model.BulkDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	deleted, err := service.BulkDeleteTransactions(userID, req.IDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.BulkDeleteResponse{
		Deleted: deleted,
		Message: "transaksi berhasil dihapus",
	})
}

func GetSpendingStreak(c *gin.Context) {
	userID := middleware.GetUserID(c)
	streak, err := service.GetSpendingStreak(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, streak)
}
