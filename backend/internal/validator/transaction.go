package validator

import (
	"time"

	apperr "github.com/akrom/finance-backend/internal/errors"
	"github.com/akrom/finance-backend/internal/model"
)

const (
	maxYearsBack  = 5
	maxYearsAhead = 1
)

// ValidateTransactionRequest validates all business rules for creating or
// updating a transaction. It returns a *ValidationError on the first rule
// violation, or nil when the request is valid.
func ValidateTransactionRequest(req model.TransactionRequest) error {
	// Amount must be positive
	if req.Amount <= 0 {
		return apperr.NewValidation("amount", "jumlah harus lebih dari 0")
	}

	// Type must be income or expense
	if req.Type != "income" && req.Type != "expense" {
		return apperr.NewValidation("type", "jenis transaksi harus 'income' atau 'expense'")
	}

	// CategoryID must be set
	if req.CategoryID <= 0 {
		return apperr.NewValidation("category_id", "kategori harus dipilih")
	}

	// Date must be parseable
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return apperr.NewValidation("date", "format tanggal tidak valid, gunakan YYYY-MM-DD")
	}

	// Date range: not older than maxYearsBack years, not more than maxYearsAhead years ahead
	now := time.Now()
	earliest := now.AddDate(-maxYearsBack, 0, 0)
	latest := now.AddDate(maxYearsAhead, 0, 0)
	if date.Before(earliest) {
		return apperr.NewValidation("date", "tanggal terlalu jauh ke masa lalu (maksimal 5 tahun)")
	}
	if date.After(latest) {
		return apperr.NewValidation("date", "tanggal terlalu jauh ke masa depan (maksimal 1 tahun)")
	}

	return nil
}

// ValidateBulkDeleteRequest validates a bulk-delete payload.
func ValidateBulkDeleteRequest(ids []int) error {
	if len(ids) == 0 {
		return apperr.NewValidation("ids", "daftar id tidak boleh kosong")
	}
	if len(ids) > 100 {
		return apperr.NewValidation("ids", "maksimal 100 transaksi dapat dihapus sekaligus")
	}
	seen := make(map[int]bool, len(ids))
	for _, id := range ids {
		if id <= 0 {
			return apperr.NewValidation("ids", "setiap id harus bilangan positif")
		}
		if seen[id] {
			return apperr.NewValidation("ids", "terdapat id duplikat dalam daftar")
		}
		seen[id] = true
	}
	return nil
}
