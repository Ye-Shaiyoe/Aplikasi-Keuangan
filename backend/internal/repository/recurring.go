package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/akrom/finance-backend/internal/database"
	"github.com/akrom/finance-backend/internal/model"
)

const recurringSelectCols = `r.id, r.user_id, r.category_id, r.amount, r.description, r.type, r.frequency,
	r.start_date::text, COALESCE(r.end_date::text, '') AS end_date, r.next_date::text, COALESCE(r.last_processed::text, '') AS last_processed, r.is_active,
	r.created_at, r.updated_at,
	c.name AS category_name, c.icon AS category_icon, c.color AS category_color`

func scanRecurring(row interface {
	Scan(dest ...interface{}) error
}, r *model.RecurringTransaction) error {
	return row.Scan(
		&r.ID, &r.UserID, &r.CategoryID, &r.Amount, &r.Description, &r.Type, &r.Frequency,
		&r.StartDate, &r.EndDate, &r.NextDate, &r.LastProcessed, &r.IsActive,
		&r.CreatedAt, &r.UpdatedAt,
		&r.CategoryName, &r.CategoryIcon, &r.CategoryColor,
	)
}

func GetRecurring(userID int) ([]model.RecurringTransaction, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM recurring_transactions r
		JOIN categories c ON c.id = r.category_id
		WHERE r.user_id = $1
		ORDER BY r.is_active DESC, r.next_date ASC, r.id DESC
	`, recurringSelectCols)

	rows, err := database.Pool.Query(context.Background(), query, userID)
	if err != nil {
		return nil, fmt.Errorf("query recurring: %w", err)
	}
	defer rows.Close()

	var items []model.RecurringTransaction
	for rows.Next() {
		var r model.RecurringTransaction
		if err := scanRecurring(rows, &r); err != nil {
			return nil, err
		}
		items = append(items, r)
	}
	return items, nil
}

func GetRecurringByID(id, userID int) (*model.RecurringTransaction, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM recurring_transactions r
		JOIN categories c ON c.id = r.category_id
		WHERE r.id = $1 AND r.user_id = $2
	`, recurringSelectCols)

	var r model.RecurringTransaction
	if err := scanRecurring(database.Pool.QueryRow(context.Background(), query, id, userID), &r); err != nil {
		return nil, fmt.Errorf("get recurring %d: %w", id, err)
	}
	return &r, nil
}

func CreateRecurring(userID int, req model.RecurringRequest) (*model.RecurringTransaction, error) {
	active := true
	if req.IsActive != nil {
		active = *req.IsActive
	}

	var endDate interface{}
	if req.EndDate != "" {
		endDate = req.EndDate
	}

	query := `
		INSERT INTO recurring_transactions (user_id, category_id, amount, description, type, frequency, start_date, end_date, next_date, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $7, $9)
		RETURNING id, user_id, category_id, amount, description, type, frequency,
			start_date::text, COALESCE(end_date::text, '') AS end_date, next_date::text, COALESCE(last_processed::text, '') AS last_processed, is_active,
			created_at, updated_at
	`
	var r model.RecurringTransaction
	err := database.Pool.QueryRow(context.Background(), query,
		userID, req.CategoryID, req.Amount, req.Description, req.Type, req.Frequency,
		req.StartDate, endDate, active).Scan(
		&r.ID, &r.UserID, &r.CategoryID, &r.Amount, &r.Description, &r.Type, &r.Frequency,
		&r.StartDate, &r.EndDate, &r.NextDate, &r.LastProcessed, &r.IsActive,
		&r.CreatedAt, &r.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("create recurring: %w", err)
	}
	r.CategoryID = req.CategoryID
	return &r, nil
}

func UpdateRecurring(id, userID int, req model.RecurringRequest) (*model.RecurringTransaction, error) {
	active := true
	if req.IsActive != nil {
		active = *req.IsActive
	}

	var endDate interface{}
	if req.EndDate != "" {
		endDate = req.EndDate
	}

	// If start_date changes, reset next_date to new start_date. Otherwise, keep next_date unless start_date moved past it.
	query := `
		UPDATE recurring_transactions SET
			category_id = $1, amount = $2, description = $3, type = $4, frequency = $5,
			start_date = $6, end_date = $7, is_active = $8,
			next_date = CASE 
				WHEN start_date <> $6::date THEN $6::date
				WHEN $6::date > next_date THEN $6::date 
				ELSE next_date 
			END,
			updated_at = NOW()
		WHERE id = $9 AND user_id = $10
		RETURNING id, user_id, category_id, amount, description, type, frequency,
			start_date::text, COALESCE(end_date::text, '') AS end_date, next_date::text, COALESCE(last_processed::text, '') AS last_processed, is_active,
			created_at, updated_at
	`
	var r model.RecurringTransaction
	err := database.Pool.QueryRow(context.Background(), query,
		req.CategoryID, req.Amount, req.Description, req.Type, req.Frequency,
		req.StartDate, endDate, active, id, userID).Scan(
		&r.ID, &r.UserID, &r.CategoryID, &r.Amount, &r.Description, &r.Type, &r.Frequency,
		&r.StartDate, &r.EndDate, &r.NextDate, &r.LastProcessed, &r.IsActive,
		&r.CreatedAt, &r.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("update recurring %d: %w", id, err)
	}
	return &r, nil
}

func DeleteRecurring(id, userID int) error {
	_, err := database.Pool.Exec(context.Background(),
		"DELETE FROM recurring_transactions WHERE id = $1 AND user_id = $2", id, userID)
	if err != nil {
		return fmt.Errorf("delete recurring %d: %w", id, err)
	}
	return nil
}

// GetDueRecurring returns active recurring transactions whose next_date is due.
// If userID > 0, only that user's items are returned; otherwise all users.
func GetDueRecurring(userID int) ([]model.RecurringTransaction, error) {
	query := `
		SELECT id, user_id, category_id, amount, description, type, frequency,
			start_date::text, COALESCE(end_date::text, '') AS end_date, next_date::text, COALESCE(last_processed::text, '') AS last_processed, is_active,
			created_at, updated_at
		FROM recurring_transactions
		WHERE is_active = TRUE AND next_date <= CURRENT_DATE
	`
	args := []interface{}{}
	if userID > 0 {
		query += "AND user_id = $1 "
		args = append(args, userID)
	}
	query += "ORDER BY next_date ASC"

	rows, err := database.Pool.Query(context.Background(), query, args...)
	if err != nil {
		return nil, fmt.Errorf("query due recurring: %w", err)
	}
	defer rows.Close()

	var items []model.RecurringTransaction
	for rows.Next() {
		var r model.RecurringTransaction
		if err := rows.Scan(
			&r.ID, &r.UserID, &r.CategoryID, &r.Amount, &r.Description, &r.Type, &r.Frequency,
			&r.StartDate, &r.EndDate, &r.NextDate, &r.LastProcessed, &r.IsActive,
			&r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

// AdvanceRecurring moves next_date forward by one frequency period and sets
// last_processed to the old next_date. If an end_date is set and the new
// next_date exceeds it, the recurring is deactivated.
func AdvanceRecurring(id int, nextDate time.Time, frequency, endDateStr string) error {
	var newDate time.Time
	switch frequency {
	case "daily":
		newDate = nextDate.AddDate(0, 0, 1)
	case "weekly":
		newDate = nextDate.AddDate(0, 0, 7)
	case "monthly":
		newDate = nextDate.AddDate(0, 1, 0)
	case "yearly":
		newDate = nextDate.AddDate(1, 0, 0)
	default:
		newDate = nextDate.AddDate(0, 1, 0)
	}

	active := true
	if endDateStr != "" {
		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err == nil && newDate.After(endDate) {
			active = false
		}
	}

	// last_processed = the next_date that was just processed
	_, err := database.Pool.Exec(context.Background(), `
		UPDATE recurring_transactions
		SET last_processed = $1, next_date = $2, is_active = $3, updated_at = NOW()
		WHERE id = $4
	`, nextDate.Format("2006-01-02"), newDate.Format("2006-01-02"), active, id)
	if err != nil {
		return fmt.Errorf("advance recurring %d: %w", id, err)
	}
	return nil
}
