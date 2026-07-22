package model

// MLPredictCategoryRequest is sent from frontend to Go, then proxied to Python ML service
type MLPredictCategoryRequest struct {
	Description string `json:"description" binding:"required"`
}

// MLPredictCategoryResponse is what the Python ML service returns
type MLPredictCategoryResponse struct {
	PredictedCategoryID   int                 `json:"predicted_category_id"`
	PredictedCategoryName string              `json:"predicted_category_name"`
	Confidence            float64             `json:"confidence"`
	TopPredictions        []MLTopPrediction   `json:"top_predictions"`
}

type MLTopPrediction struct {
	CategoryID   int     `json:"category_id"`
	CategoryName string  `json:"category_name"`
	Confidence   float64 `json:"confidence"`
}

// MLForecastRequest is sent to Python ML service with monthly data
type MLForecastRequest struct {
	MonthlyData []MLMonthlyData `json:"monthly_data"`
}

type MLMonthlyData struct {
	Month   int   `json:"month"`
	Year    int   `json:"year"`
	Income  int64 `json:"income"`
	Expense int64 `json:"expense"`
}

// MLForecastResponse is what the Python ML service returns
type MLForecastResponse struct {
	ForecastExpense int64              `json:"forecast_expense"`
	ForecastIncome  int64              `json:"forecast_income"`
	ForecastBalance int64              `json:"forecast_balance"`
	Confidence      MLConfidence       `json:"confidence"`
	Projections     []MLProjection     `json:"projections"`
	DataPoints      int                `json:"data_points"`
	ModelType       string             `json:"model_type"`
}

type MLConfidence struct {
	ExpenseR2 float64 `json:"expense_r2"`
	IncomeR2  float64 `json:"income_r2"`
	Level     string  `json:"level"`
}

type MLProjection struct {
	MonthOffset     int   `json:"month_offset"`
	ForecastExpense int64 `json:"forecast_expense"`
	ForecastIncome  int64 `json:"forecast_income"`
	ForecastBalance int64 `json:"forecast_balance"`
}

// MLTrainRequest is sent to Python ML service to train category classifier
type MLTrainRequest struct {
	Transactions []MLTrainTransaction `json:"transactions"`
	Categories   []MLTrainCategory    `json:"categories"`
}

type MLTrainTransaction struct {
	Description string `json:"description"`
	CategoryID  int    `json:"category_id"`
	Type        string `json:"type"`
}

type MLTrainCategory struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// MLTrainResponse is what the Python ML service returns after training
type MLTrainResponse struct {
	Status  string                 `json:"status"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details"`
}

// MLHealthResponse from the Python ML service health check
type MLHealthResponse struct {
	Status  string                 `json:"status"`
	Service string                 `json:"service"`
	Models  map[string]string      `json:"models"`
}
