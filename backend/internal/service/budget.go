package service

import (
	"fmt"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
)

func GetBudgets(userID, month, year int) ([]model.Budget, error) {
	return repository.GetBudgets(userID, month, year)
}

func UpsertBudget(userID int, req model.BudgetRequest) (*model.Budget, error) {
	return repository.UpsertBudget(userID, req)
}

func DeleteBudget(userID, id int) error {
	return repository.DeleteBudget(userID, id)
}

func GetBudgetSummary(userID, month, year int) (*model.BudgetSummaryResponse, error) {
	return repository.GetBudgetSummary(userID, month, year)
}

// AlertBudgetOverspend checks each category budget for the given month/year and
// returns alerts for categories whose spending has reached the threshold.
// threshold is a percentage in [1, 100]; typical values are 80 (warning) or 100 (critical).
func AlertBudgetOverspend(userID, month, year, threshold int) (*model.BudgetAlertResponse, error) {
	bsSlice, err := repository.GetBudgetsWithSpending(userID, month, year)
	if err != nil {
		return nil, fmt.Errorf("alert budget: get budgets: %w", err)
	}

	if threshold < 1 {
		threshold = 80 // sensible default
	}

	var alerts []model.BudgetAlert
	warnings := 0
	critical := 0

	for _, bs := range bsSlice {
		if bs.Amount == 0 {
			continue
		}
		pct := float64(bs.Spent) / float64(bs.Amount) * 100
		if pct < float64(threshold) {
			continue
		}

		severity := "warning"
		msg := fmt.Sprintf("Pengeluaran %.0f%% dari anggaran %s", pct, bs.CategoryName)
		if pct >= 100 {
			severity = "critical"
			msg = fmt.Sprintf("Anggaran %s sudah terlampaui (%.0f%%)", bs.CategoryName, pct)
			critical++
		} else {
			warnings++
		}

		alerts = append(alerts, model.BudgetAlert{
			CategoryID:    bs.CategoryID,
			CategoryName:  bs.CategoryName,
			CategoryColor: bs.CategoryColor,
			BudgetAmount:  bs.Amount,
			SpentAmount:   bs.Spent,
			PercentUsed:   pct,
			Severity:      severity,
			Message:       msg,
		})
	}

	if alerts == nil {
		alerts = []model.BudgetAlert{}
	}

	return &model.BudgetAlertResponse{
		Month:    month,
		Year:     year,
		Alerts:   alerts,
		Warnings: warnings,
		Critical: critical,
	}, nil
}

// CopyBudgetsFromLastMonth copies all budgets from the previous calendar month
// to the given month/year. Returns the number of budgets copied.
func CopyBudgetsFromLastMonth(userID, toMonth, toYear int) (int, error) {
	// Calculate the previous month
	ref := time.Date(toYear, time.Month(toMonth), 1, 0, 0, 0, 0, time.UTC).AddDate(0, -1, 0)
	fromMonth := int(ref.Month())
	fromYear := ref.Year()

	copied, err := repository.CopyBudgetsFromMonth(userID, fromMonth, fromYear, toMonth, toYear)
	if err != nil {
		return 0, fmt.Errorf("copy budgets: %w", err)
	}
	return copied, nil
}
