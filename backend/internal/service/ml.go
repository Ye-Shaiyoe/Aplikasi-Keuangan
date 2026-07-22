package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
)

var mlServiceURL string

func init() {
	mlServiceURL = os.Getenv("ML_SERVICE_URL")
	if mlServiceURL == "" {
		mlServiceURL = "http://localhost:8000"
	}
}

var httpClient = &http.Client{
	Timeout: 30 * time.Second,
}

// MLPredictCategory sends a description to the Python ML service and returns prediction
func MLPredictCategory(description string) (*model.MLPredictCategoryResponse, error) {
	reqBody := model.MLPredictCategoryRequest{Description: description}
	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("ml service: marshal request: %w", err)
	}

	resp, err := httpClient.Post(
		mlServiceURL+"/api/ml/predict-category",
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf("ml service: request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ml service: read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ml service: status %d: %s", resp.StatusCode, string(respBody))
	}

	var result model.MLPredictCategoryResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("ml service: decode response: %w", err)
	}

	return &result, nil
}

// MLForecast sends monthly data to the Python ML service for forecasting
func MLForecast(monthlyData []model.MLMonthlyData) (*model.MLForecastResponse, error) {
	reqBody := model.MLForecastRequest{MonthlyData: monthlyData}
	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("ml service: marshal request: %w", err)
	}

	resp, err := httpClient.Post(
		mlServiceURL+"/api/ml/forecast",
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf("ml service: request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ml service: read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ml service: status %d: %s", resp.StatusCode, string(respBody))
	}

	var result model.MLForecastResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("ml service: decode response: %w", err)
	}

	return &result, nil
}

// MLTrainModel collects user's transactions and categories, then sends to Python ML for training
func MLTrainModel(userID int) (*model.MLTrainResponse, error) {
	// Fetch all transactions for this user
	transactions, err := repository.GetAllTransactionsForML(userID)
	if err != nil {
		return nil, fmt.Errorf("ml train: fetch transactions: %w", err)
	}

	// Fetch all categories for this user
	categories, err := repository.GetAllCategoriesForML(userID)
	if err != nil {
		return nil, fmt.Errorf("ml train: fetch categories: %w", err)
	}

	reqBody := model.MLTrainRequest{
		Transactions: transactions,
		Categories:   categories,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("ml service: marshal train request: %w", err)
	}

	resp, err := httpClient.Post(
		mlServiceURL+"/api/ml/train",
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf("ml service: train request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ml service: read train response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ml service: train status %d: %s", resp.StatusCode, string(respBody))
	}

	var result model.MLTrainResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("ml service: decode train response: %w", err)
	}

	return &result, nil
}

// MLHealthCheck checks if the Python ML service is healthy
func MLHealthCheck() (*model.MLHealthResponse, error) {
	resp, err := httpClient.Get(mlServiceURL + "/api/ml/health")
	if err != nil {
		return nil, fmt.Errorf("ml service: health check failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ml service: read health response: %w", err)
	}

	var result model.MLHealthResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("ml service: decode health response: %w", err)
	}

	return &result, nil
}

// MLForecastForUser fetches user's historical data and sends to ML service
func MLForecastForUser(userID int) (*model.MLForecastResponse, error) {
	// Get 12 months of historical data
	metrics, err := repository.GetHistoricalMetrics(userID, 12)
	if err != nil {
		return nil, fmt.Errorf("ml forecast: fetch metrics: %w", err)
	}

	monthlyData := make([]model.MLMonthlyData, len(metrics))
	for i, m := range metrics {
		monthlyData[i] = model.MLMonthlyData{
			Month:   m.Month,
			Year:    m.Year,
			Income:  m.Income,
			Expense: m.Expense,
		}
	}

	return MLForecast(monthlyData)
}
