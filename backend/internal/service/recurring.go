package service

import (
	"fmt"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/repository"
)

func GetRecurring(userID int) ([]model.RecurringTransaction, error) {
	return repository.GetRecurring(userID)
}

func GetRecurringByID(id, userID int) (*model.RecurringTransaction, error) {
	return repository.GetRecurringByID(id, userID)
}

func CreateRecurring(userID int, req model.RecurringRequest) (*model.RecurringTransaction, error) {
	return repository.CreateRecurring(userID, req)
}

func UpdateRecurring(id, userID int, req model.RecurringRequest) (*model.RecurringTransaction, error) {
	return repository.UpdateRecurring(id, userID, req)
}

func DeleteRecurring(id, userID int) error {
	return repository.DeleteRecurring(id, userID)
}

// ProcessDueRecurring generates real transactions for every due recurring item.
// If userID is 0, all users are processed (used at server startup); otherwise
// only the given user's items are processed (manual trigger). It loops to catch
// up on missed periods, capped to avoid runaway generation.
func ProcessDueRecurring(userID int) (int, error) {
	processed := 0
	const maxIterations = 500
	failedIDs := make(map[int]bool)

	for iter := 0; iter < maxIterations; iter++ {
		due, err := repository.GetDueRecurring(userID)
		if err != nil {
			return processed, fmt.Errorf("get due recurring: %w", err)
		}

		var freshDue []model.RecurringTransaction
		for _, r := range due {
			if !failedIDs[r.ID] {
				freshDue = append(freshDue, r)
			}
		}

		if len(freshDue) == 0 {
			break
		}

		for _, r := range freshDue {
			req := model.TransactionRequest{
				CategoryID:  r.CategoryID,
				Amount:      r.Amount,
				Description: r.Description,
				Date:        r.NextDate,
				Type:        r.Type,
			}

			parsedNextDate, err := time.Parse("2006-01-02", r.NextDate)
			if err != nil {
				fmt.Printf("recurring: parse next_date %s failed for recurring %d: %v\n", r.NextDate, r.ID, err)
				failedIDs[r.ID] = true
				continue
			}

			if _, err := repository.CreateTransaction(r.UserID, req); err != nil {
				// Log and skip this one so one bad row doesn't block the rest.
				fmt.Printf("recurring: create transaction for recurring %d failed: %v\n", r.ID, err)
				failedIDs[r.ID] = true
				continue
			}
			if err := repository.AdvanceRecurring(r.ID, parsedNextDate, r.Frequency, r.EndDate); err != nil {
				fmt.Printf("recurring: advance recurring %d failed: %v\n", r.ID, err)
				failedIDs[r.ID] = true
				continue
			}
			processed++
		}
	}

	return processed, nil
}
