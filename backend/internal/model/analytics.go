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
