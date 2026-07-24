package model

type AdvancedAnalyticsResponse struct {
	HealthScore     int             `json:"health_score"`     // 0-100 score
	HealthRating    string          `json:"health_rating"`    // "Sangat Baik", "Baik", "Cukup", "Perlu Perhatian"
	SavingsRate     float64         `json:"savings_rate"`     // percentage: (income - expense) / income * 100
	BudgetAdherence float64         `json:"budget_adherence"` // percentage of categories within budget
	ForecastExpense int64           `json:"forecast_expense"` // predicted next-month expense
	ForecastIncome  int64           `json:"forecast_income"`  // predicted next-month income
	Insights        []string        `json:"insights"`         // personalized tips
	MonthlyMetrics  []MonthlyMetric `json:"monthly_metrics"`  // history for frontend chart
}

type MonthlyMetric struct {
	MonthName string `json:"month_name"`
	Month     int    `json:"month"`
	Year      int    `json:"year"`
	Income    int64  `json:"income"`
	Expense   int64  `json:"expense"`
}

// HeatmapDay holds aggregated spending data for a single calendar day.
type HeatmapDay struct {
	Date    string `json:"date"`   // YYYY-MM-DD
	Amount  int64  `json:"amount"` // total expense for the day
	Count   int    `json:"count"`  // number of expense transactions
	Level   int    `json:"level"`  // 0-4 intensity level for rendering
}

// ExpenseHeatmapResponse contains a full year of daily expense data.
type ExpenseHeatmapResponse struct {
	Year      int          `json:"year"`
	MaxAmount int64        `json:"max_amount"` // highest single-day total (for normalisation)
	Days      []HeatmapDay `json:"days"`
}

// NetWorthPoint is the cumulative balance at a single month boundary.
type NetWorthPoint struct {
	MonthName     string `json:"month_name"`
	Month         int    `json:"month"`
	Year          int    `json:"year"`
	MonthlyIncome int64  `json:"monthly_income"`
	MonthlyExpense int64 `json:"monthly_expense"`
	CumulativeBalance int64 `json:"cumulative_balance"`
}

// NetWorthTimelineResponse wraps a time-ordered list of net-worth snapshots.
type NetWorthTimelineResponse struct {
	Points          []NetWorthPoint `json:"points"`
	CurrentBalance  int64           `json:"current_balance"`
	TotalMonths     int             `json:"total_months"`
}
