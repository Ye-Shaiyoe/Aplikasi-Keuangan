package model

import "time"

type Budget struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	CategoryID int       `json:"category_id"`
	Month      int       `json:"month"`
	Year       int       `json:"year"`
	Amount     int64     `json:"amount"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type BudgetRequest struct {
	CategoryID int   `json:"category_id" binding:"required"`
	Month      int   `json:"month" binding:"required,min=1,max=12"`
	Year       int   `json:"year" binding:"required,min=2020,max=2100"`
	Amount     int64 `json:"amount" binding:"required,min=0"`
}

// BudgetWithSpending joins budget with actual expense total
type BudgetWithSpending struct {
	Budget
	CategoryName  string `json:"category_name"`
	CategoryColor string `json:"category_color"`
	Spent         int64  `json:"spent"`
	Remaining     int64  `json:"remaining"`
	PercentUsed   int    `json:"percent_used"`
}

type BudgetSummaryResponse struct {
	Month       int                  `json:"month"`
	Year        int                  `json:"year"`
	TotalBudget int64                `json:"total_budget"`
	TotalSpent  int64                `json:"total_spent"`
	Budgets     []BudgetWithSpending `json:"budgets"`
}

// BudgetAlert is generated when actual spending exceeds a configured threshold
// percentage of the budget for a category.
type BudgetAlert struct {
	CategoryID    int     `json:"category_id"`
	CategoryName  string  `json:"category_name"`
	CategoryColor string  `json:"category_color"`
	BudgetAmount  int64   `json:"budget_amount"`
	SpentAmount   int64   `json:"spent_amount"`
	PercentUsed   float64 `json:"percent_used"`
	Severity      string  `json:"severity"` // "warning" (>=80%) or "critical" (>=100%)
	Message       string  `json:"message"`
}

// BudgetAlertResponse wraps a list of budget alerts with summary counts.
type BudgetAlertResponse struct {
	Month    int           `json:"month"`
	Year     int           `json:"year"`
	Alerts   []BudgetAlert `json:"alerts"`
	Warnings int           `json:"warnings"` // count of categories >= warningThreshold but < 100%
	Critical int           `json:"critical"` // count of categories >= 100%
}

// CopyBudgetResponse reports the result of a copy-budget-from-last-month operation.
type CopyBudgetResponse struct {
	Copied  int    `json:"copied"`
	Message string `json:"message"`
}
