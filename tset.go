package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
	"github.com/akrom/finance-backend/internal/util"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	return repository.GetTransactions(userID, filter)
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	return repository.GetTransactionByID(id, userID)
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.CreateTransaction(userID, req)
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	return repository.UpdateTransaction(id, userID, req)
}

func DeleteTransaction(id, userID int) error {
	return repository.DeleteTransaction(id, userID)
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	return repository.GetSummary(userID, month, year)
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	return repository.GetYearlyTrend(userID, year)
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	return repository.GetCategoryTrend(userID, month, year)
}

// ExportTransactionsCSV returns all transactions matching the filter as a
// UTF-8 CSV byte slice (with BOM for Excel). The filter page/limit are ignored;
// all matching rows are returned.
func ExportTransactionsCSV(userID int, filter model.TransactionFilter) ([]byte, error) {
	// Fetch all by using a high limit
	filter.Page = 1
	filter.Limit = 10_000
	transactions, _, err := repository.GetTransactions(userID, filter)
	if err != nil {
		return nil, fmt.Errorf("export csv: fetch transactions: %w", err)
	}
	return util.TransactionsToCSV(transactions), nil
}

// BulkDeleteTransactions deletes multiple transactions in a single operation.
// Only transactions that belong to userID are deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	deleted, err := repository.BulkDeleteTransactions(userID, ids)
	if err != nil {
		return 0, fmt.Errorf("bulk delete: %w", err)
	}
	return deleted, nil
}

// GetSpendingStreak calculates the current and all-time longest consecutive-day
// streak during which the user has recorded at least one transaction.
func GetSpendingStreak(userID int) (*model.SpendingStreakResponse, error) {
	dates, err := repository.GetAllTransactionDates(userID)
	if err != nil {
		return nil, fmt.Errorf("streak: fetch dates: %w", err)
	}

	if len(dates) == 0 {
		return &model.SpendingStreakResponse{}, nil
	}

	// De-duplicate and sort dates (should already be sorted from DB)
	seen := make(map[string]bool)
	var uniqueDates []time.Time
	for _, ds := range dates {
		if seen[ds] {
			continue
		}
		seen[ds] = true
		t, err := time.Parse("2006-01-02", ds)
		if err != nil {
			continue
		}
		uniqueDates = append(uniqueDates, t)
	}
	sort.Slice(uniqueDates, func(i, j int) bool {
		return uniqueDates[i].Before(uniqueDates[j])
	})

	// Calculate streaks
	currentStreak := 1
	longestStreak := 1
	tempStreak := 1

	for i := 1; i < len(uniqueDates); i++ {
		diff := uniqueDates[i].Sub(uniqueDates[i-1])
		if diff == 24*time.Hour {
			tempStreak++
		} else {
			tempStreak = 1
		}
		if tempStreak > longestStreak {
			longestStreak = tempStreak
		}
	}

	// Current streak: count backwards from today
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastDate := uniqueDates[len(uniqueDates)-1].UTC().Truncate(24 * time.Hour)

	isActiveToday := lastDate.Equal(today)

	// If last activity was today or yesterday, count the current streak
	if lastDate.Equal(today) || today.Sub(lastDate) == 24*time.Hour {
		currentStreak = 1
		for i := len(uniqueDates) - 2; i >= 0; i-- {
			expected := uniqueDates[i+1].UTC().Truncate(24*time.Hour).AddDate(0, 0, -1)
			if uniqueDates[i].UTC().Truncate(24*time.Hour).Equal(expected) {
				currentStreak++
			} else {
				break
			}
		}
	} else {
		currentStreak = 0 // streak is broken
	}

	return &model.SpendingStreakResponse{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
		LastActiveDay: util.FormatDate(lastDate),
		IsActiveToday: isActiveToday,
	}, nil
}
