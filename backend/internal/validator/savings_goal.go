package validator

import (
	"strings"
	"time"

	apperr "github.com/akrom/finance-backend/internal/errors"
	"github.com/akrom/finance-backend/internal/model"
)

// ValidateSavingsGoalRequest validates business rules for creating or updating
// a savings goal.
func ValidateSavingsGoalRequest(req model.SavingsGoalRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return apperr.NewValidation("name", "nama target tabungan tidak boleh kosong")
	}
	if len(req.Name) < 2 {
		return apperr.NewValidation("name", "nama target tabungan minimal 2 karakter")
	}
	if len(req.Name) > 100 {
		return apperr.NewValidation("name", "nama target tabungan maksimal 100 karakter")
	}

	if req.TargetAmount <= 0 {
		return apperr.NewValidation("target_amount", "jumlah target harus lebih dari 0")
	}

	// Validate deadline only if provided
	if req.Deadline != nil && *req.Deadline != "" {
		deadline, err := time.Parse("2006-01-02", *req.Deadline)
		if err != nil {
			return apperr.NewValidation("deadline", "format deadline tidak valid, gunakan YYYY-MM-DD")
		}
		today := time.Now().UTC().Truncate(24 * time.Hour)
		if !deadline.After(today) {
			return apperr.NewValidation("deadline", "deadline harus di masa depan")
		}
		// Deadline must not exceed 50 years from now
		maxDeadline := time.Now().AddDate(50, 0, 0)
		if deadline.After(maxDeadline) {
			return apperr.NewValidation("deadline", "deadline terlalu jauh ke masa depan (maksimal 50 tahun)")
		}
	}

	return nil
}

// ValidateDepositWithdraw validates that amount is positive for deposit/withdraw.
func ValidateDepositWithdraw(amount int64, currentAmount, targetAmount int64, isWithdraw bool) error {
	if amount <= 0 {
		return apperr.NewValidation("amount", "jumlah harus lebih dari 0")
	}
	if isWithdraw && amount > currentAmount {
		return apperr.NewValidation("amount", "jumlah penarikan melebihi saldo tabungan saat ini")
	}
	return nil
}

// ValidateAutoAllocate validates an auto-allocate request.
func ValidateAutoAllocate(amount int64) error {
	if amount <= 0 {
		return apperr.NewValidation("amount", "jumlah alokasi harus lebih dari 0")
	}
	return nil
}
