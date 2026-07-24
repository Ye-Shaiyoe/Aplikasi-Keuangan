package validator

import (
	"time"

	apperr "github.com/akrom/finance-backend/internal/errors"
	"github.com/akrom/finance-backend/internal/model"
)

const (
	minYear = 2020
	maxYear = 2100
)

// ValidateBudgetRequest validates business rules for creating or updating a budget.
func ValidateBudgetRequest(req model.BudgetRequest) error {
	if req.CategoryID <= 0 {
		return apperr.NewValidation("category_id", "kategori harus dipilih")
	}

	if req.Month < 1 || req.Month > 12 {
		return apperr.NewValidation("month", "bulan harus antara 1 dan 12")
	}

	if req.Year < minYear || req.Year > maxYear {
		return apperr.NewValidation("year", "tahun harus antara 2020 dan 2100")
	}

	if req.Amount < 0 {
		return apperr.NewValidation("amount", "jumlah anggaran tidak boleh negatif")
	}

	return nil
}

// ValidateCopyBudgetRequest validates the parameters for copying budgets from one
// month to another.
func ValidateCopyBudgetRequest(fromMonth, fromYear, toMonth, toYear int) error {
	if fromMonth < 1 || fromMonth > 12 {
		return apperr.NewValidation("from_month", "bulan sumber harus antara 1 dan 12")
	}
	if toMonth < 1 || toMonth > 12 {
		return apperr.NewValidation("to_month", "bulan tujuan harus antara 1 dan 12")
	}
	if fromYear < minYear || fromYear > maxYear {
		return apperr.NewValidation("from_year", "tahun sumber tidak valid")
	}
	if toYear < minYear || toYear > maxYear {
		return apperr.NewValidation("to_year", "tahun tujuan tidak valid")
	}

	from := time.Date(fromYear, time.Month(fromMonth), 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(toYear, time.Month(toMonth), 1, 0, 0, 0, 0, time.UTC)
	if !to.After(from) {
		return apperr.NewValidation("to_month", "bulan tujuan harus setelah bulan sumber")
	}
	return nil
}

// ValidateAlertThreshold validates that a budget-alert threshold is in [1, 100].
func ValidateAlertThreshold(threshold int) error {
	if threshold < 1 || threshold > 100 {
		return apperr.NewValidation("threshold", "threshold harus antara 1 dan 100")
	}
	return nil
}
