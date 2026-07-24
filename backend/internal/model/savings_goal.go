package model

import "time"

type SavingsGoal struct {
	ID            int       `json:"id"`
	UserID        int       `json:"user_id"`
	Name          string    `json:"name"`
	TargetAmount  int64     `json:"target_amount"`
	CurrentAmount int64     `json:"current_amount"`
	Deadline      *string   `json:"deadline"` // optional, format: "2025-12-31"
	Color         string    `json:"color"`
	ImageURL      string    `json:"image_url"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type SavingsGoalRequest struct {
	Name         string  `json:"name" binding:"required,min=1"`
	TargetAmount int64   `json:"target_amount" binding:"required,min=1"`
	Deadline     *string `json:"deadline"`
	Color        string  `json:"color"`
	ImageURL     string  `json:"image_url"`
}

type SavingsGoalDeposit struct {
	Amount int64 `json:"amount" binding:"required,min=1"`
}

// SavingsProjectionResponse estimates when a savings goal will be reached based
// on the user's historical average deposit rate.
type SavingsProjectionResponse struct {
	GoalID              int     `json:"goal_id"`
	GoalName            string  `json:"goal_name"`
	TargetAmount        int64   `json:"target_amount"`
	CurrentAmount       int64   `json:"current_amount"`
	RemainingAmount     int64   `json:"remaining_amount"`
	ProgressPercent     float64 `json:"progress_percent"`
	AvgMonthlyDeposit   int64   `json:"avg_monthly_deposit"`   // Rp per month
	EstimatedMonths     int     `json:"estimated_months"`      // months to reach target
	EstimatedCompletion string  `json:"estimated_completion"` // YYYY-MM-DD or "tidak dapat diprediksi"
	IsOnTrack           bool    `json:"is_on_track"`           // true if deadline >= estimated completion
}

// AllocationResult describes how much was allocated to a single savings goal.
type AllocationResult struct {
	GoalID   int    `json:"goal_id"`
	GoalName string `json:"goal_name"`
	Allocated int64 `json:"allocated"`
}

// AutoAllocateResponse wraps the results of an auto-allocate operation.
type AutoAllocateResponse struct {
	TotalAllocated int64              `json:"total_allocated"`
	Allocations    []AllocationResult `json:"allocations"`
	Message        string             `json:"message"`
}
