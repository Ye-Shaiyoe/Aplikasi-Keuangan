package service

import (
	"fmt"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
)

func GetSavingsGoals(userID int) ([]model.SavingsGoal, error) {
	return repository.GetSavingsGoals(userID)
}

func GetSavingsGoalByID(userID, id int) (*model.SavingsGoal, error) {
	return repository.GetSavingsGoalByID(userID, id)
}

func CreateSavingsGoal(userID int, req model.SavingsGoalRequest) (*model.SavingsGoal, error) {
	return repository.CreateSavingsGoal(userID, req)
}

func UpdateSavingsGoal(userID, id int, req model.SavingsGoalRequest) (*model.SavingsGoal, error) {
	return repository.UpdateSavingsGoal(userID, id, req)
}

func DeleteSavingsGoal(userID, id int) error {
	return repository.DeleteSavingsGoal(userID, id)
}

func DepositToSavingsGoal(userID, id int, amount int64) (*model.SavingsGoal, error) {
	return repository.DepositToSavingsGoal(userID, id, amount)
}

func WithdrawFromSavingsGoal(userID, id int, amount int64) (*model.SavingsGoal, error) {
	return repository.WithdrawFromSavingsGoal(userID, id, amount)
}

// GetSavingsProjection estimates when the goal will be reached based on the
// average monthly deposit over the life of the goal. If current_amount is
// already >= target_amount the goal is marked as completed.
func GetSavingsProjection(userID, goalID int) (*model.SavingsProjectionResponse, error) {
	goal, err := repository.GetSavingsGoalByID(userID, goalID)
	if err != nil {
		return nil, fmt.Errorf("projection: get goal: %w", err)
	}

	remaining := goal.TargetAmount - goal.CurrentAmount
	progress := float64(0)
	if goal.TargetAmount > 0 {
		progress = float64(goal.CurrentAmount) / float64(goal.TargetAmount) * 100
	}

	if remaining <= 0 {
		return &model.SavingsProjectionResponse{
			GoalID:              goal.ID,
			GoalName:            goal.Name,
			TargetAmount:        goal.TargetAmount,
			CurrentAmount:       goal.CurrentAmount,
			RemainingAmount:     0,
			ProgressPercent:     100,
			AvgMonthlyDeposit:   0,
			EstimatedMonths:     0,
			EstimatedCompletion: "Tercapai",
			IsOnTrack:           true,
		}, nil
	}

	// Estimate average monthly deposit based on how long the goal has existed
	// and how much has been saved so far.
	var avgMonthly int64
	monthsSinceCreation := 1
	if !goal.CreatedAt.IsZero() {
		years := time.Now().Year() - goal.CreatedAt.Year()
		months := int(time.Now().Month()) - int(goal.CreatedAt.Month())
		total := years*12 + months
		if total > 0 {
			monthsSinceCreation = total
		}
	}
	if goal.CurrentAmount > 0 {
		avgMonthly = goal.CurrentAmount / int64(monthsSinceCreation)
	}

	estimatedCompletion := "tidak dapat diprediksi"
	estimatedMonths := 0
	isOnTrack := false

	if avgMonthly > 0 {
		estimatedMonths = int((remaining + avgMonthly - 1) / avgMonthly) // ceiling division
		completion := time.Now().AddDate(0, estimatedMonths, 0)
		estimatedCompletion = completion.Format("2006-01-02")

		// Check if on track relative to deadline
		if goal.Deadline != nil && *goal.Deadline != "" {
			deadline, err := time.Parse("2006-01-02", *goal.Deadline)
			if err == nil {
				isOnTrack = !completion.After(deadline)
			}
		} else {
			isOnTrack = true // no deadline, always on track
		}
	}

	return &model.SavingsProjectionResponse{
		GoalID:              goal.ID,
		GoalName:            goal.Name,
		TargetAmount:        goal.TargetAmount,
		CurrentAmount:       goal.CurrentAmount,
		RemainingAmount:     remaining,
		ProgressPercent:     progress,
		AvgMonthlyDeposit:   avgMonthly,
		EstimatedMonths:     estimatedMonths,
		EstimatedCompletion: estimatedCompletion,
		IsOnTrack:           isOnTrack,
	}, nil
}

// AutoAllocateSavings distributes an amount proportionally across all active
// savings goals based on each goal's remaining deficit. Goals that are already
// complete are skipped.
func AutoAllocateSavings(userID int, amount int64) (*model.AutoAllocateResponse, error) {
	goals, err := repository.GetSavingsGoals(userID)
	if err != nil {
		return nil, fmt.Errorf("auto allocate: get goals: %w", err)
	}

	// Filter to goals that still need funding
	var active []model.SavingsGoal
	var totalDeficit int64
	for _, g := range goals {
		if g.CurrentAmount < g.TargetAmount {
			active = append(active, g)
			totalDeficit += g.TargetAmount - g.CurrentAmount
		}
	}

	if len(active) == 0 || totalDeficit == 0 {
		return &model.AutoAllocateResponse{
			TotalAllocated: 0,
			Allocations:    []model.AllocationResult{},
			Message:        "Semua target tabungan sudah tercapai",
		}, nil
	}

	var allocations []model.AllocationResult
	var totalAllocated int64

	for i, g := range active {
		deficit := g.TargetAmount - g.CurrentAmount
		var alloc int64
		if i == len(active)-1 {
			// Last goal gets the remainder to avoid rounding errors
			alloc = amount - totalAllocated
		} else {
			alloc = amount * deficit / totalDeficit
		}
		if alloc <= 0 {
			continue
		}
		// Cap at remaining deficit
		if alloc > deficit {
			alloc = deficit
		}
		if _, err := repository.DepositToSavingsGoal(userID, g.ID, alloc); err != nil {
			return nil, fmt.Errorf("auto allocate: deposit to goal %d: %w", g.ID, err)
		}
		allocations = append(allocations, model.AllocationResult{
			GoalID:    g.ID,
			GoalName:  g.Name,
			Allocated: alloc,
		})
		totalAllocated += alloc
	}

	if allocations == nil {
		allocations = []model.AllocationResult{}
	}

	return &model.AutoAllocateResponse{
		TotalAllocated: totalAllocated,
		Allocations:    allocations,
		Message:        fmt.Sprintf("Berhasil mengalokasikan Rp %d ke %d target tabungan", totalAllocated, len(allocations)),
	}, nil
}
