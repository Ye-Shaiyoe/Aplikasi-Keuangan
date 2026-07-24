package repository

import (
	"context"
	"time"

	"github.com/akrom/finance-backend/internal/database"
	"github.com/akrom/finance-backend/internal/model"
)

func GetBudgets(userID, month, year int) ([]model.Budget, error) {
	rows, err := database.Pool.Query(context.Background(),
		`SELECT id, user_id, category_id, month, year, amount, created_at, updated_at
		 FROM budgets WHERE user_id = $1 AND month = $2 AND year = $3
		 ORDER BY category_id`, userID, month, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var budgets []model.Budget
	for rows.Next() {
		var b model.Budget
		if err := rows.Scan(&b.ID, &b.UserID, &b.CategoryID, &b.Month, &b.Year, &b.Amount, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, err
		}
		budgets = append(budgets, b)
	}
	return budgets, nil
}

func UpsertBudget(userID int, req model.BudgetRequest) (*model.Budget, error) {
	var b model.Budget
	err := database.Pool.QueryRow(context.Background(),
		`INSERT INTO budgets (user_id, category_id, month, year, amount)
		 VALUES ($1, $2, $3, $4, $5)
		 ON CONFLICT (user_id, category_id, month, year) 
		 DO UPDATE SET amount = $5, updated_at = NOW()
		 RETURNING id, user_id, category_id, month, year, amount, created_at, updated_at`,
		userID, req.CategoryID, req.Month, req.Year, req.Amount).
		Scan(&b.ID, &b.UserID, &b.CategoryID, &b.Month, &b.Year, &b.Amount, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func DeleteBudget(userID, id int) error {
	_, err := database.Pool.Exec(context.Background(),
		"DELETE FROM budgets WHERE id=$1 AND user_id=$2", id, userID)
	return err
}

func GetBudgetSummary(userID, month, year int) (*model.BudgetSummaryResponse, error) {
	rows, err := database.Pool.Query(context.Background(),
		`SELECT 
			b.id, b.user_id, b.category_id, b.month, b.year, b.amount,
			b.created_at, b.updated_at,
			COALESCE(c.name, '') as category_name,
			COALESCE(c.color, '#6b7280') as category_color,
			COALESCE(SUM(t.amount) FILTER (WHERE t.type = 'expense' AND t.user_id = $1), 0) as spent
		 FROM budgets b
		 LEFT JOIN categories c ON c.id = b.category_id
		 LEFT JOIN transactions t ON t.category_id = b.category_id 
			AND t.user_id = $1 
			AND EXTRACT(MONTH FROM t.date) = $2
			AND EXTRACT(YEAR FROM t.date) = $3
			AND t.type = 'expense'
		 WHERE b.user_id = $1 AND b.month = $2 AND b.year = $3
		 GROUP BY b.id, b.user_id, b.category_id, b.month, b.year, b.amount,
			b.created_at, b.updated_at, c.name, c.color
		 ORDER BY c.name`,
		userID, month, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var budgets []model.BudgetWithSpending
	var totalBudget, totalSpent int64

	for rows.Next() {
		var bs model.BudgetWithSpending
		if err := rows.Scan(
			&bs.ID, &bs.UserID, &bs.CategoryID, &bs.Month, &bs.Year, &bs.Amount,
			&bs.CreatedAt, &bs.UpdatedAt,
			&bs.CategoryName, &bs.CategoryColor, &bs.Spent,
		); err != nil {
			return nil, err
		}
		bs.Remaining = bs.Amount - bs.Spent
		bs.PercentUsed = calcPercent(bs.Spent, bs.Amount)
		if bs.CreatedAt.IsZero() {
			bs.CreatedAt = time.Now()
		}
		if bs.UpdatedAt.IsZero() {
			bs.UpdatedAt = time.Now()
		}
		totalBudget += bs.Amount
		totalSpent += bs.Spent
		budgets = append(budgets, bs)
	}

	if budgets == nil {
		budgets = []model.BudgetWithSpending{}
	}

	return &model.BudgetSummaryResponse{
		Month:       month,
		Year:        year,
		TotalBudget: totalBudget,
		TotalSpent:  totalSpent,
		Budgets:     budgets,
	}, nil
}

func calcPercent(spent, budget int64) int {
	if budget == 0 {
		return 0
	}
	pct := int((spent * 100) / budget)
	if pct > 100 {
		pct = 100
	}
	return pct
}

// GetBudgetsWithSpending returns all budgets for a month together with actual
// spending figures. Used by the alert service.
func GetBudgetsWithSpending(userID, month, year int) ([]model.BudgetWithSpending, error) {
	rows, err := database.Pool.Query(context.Background(),
		`SELECT
			b.id, b.user_id, b.category_id, b.month, b.year, b.amount,
			b.created_at, b.updated_at,
			COALESCE(c.name, '') AS category_name,
			COALESCE(c.color, '#6b7280') AS category_color,
			COALESCE(SUM(t.amount) FILTER (WHERE t.type = 'expense'), 0) AS spent
		 FROM budgets b
		 LEFT JOIN categories c ON c.id = b.category_id
		 LEFT JOIN transactions t ON t.category_id = b.category_id
			AND t.user_id = $1
			AND EXTRACT(MONTH FROM t.date) = $2
			AND EXTRACT(YEAR FROM t.date) = $3
		 WHERE b.user_id = $1 AND b.month = $2 AND b.year = $3
		 GROUP BY b.id, b.user_id, b.category_id, b.month, b.year, b.amount,
		          b.created_at, b.updated_at, c.name, c.color
		 ORDER BY c.name`,
		userID, month, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []model.BudgetWithSpending
	for rows.Next() {
		var bs model.BudgetWithSpending
		if err := rows.Scan(
			&bs.ID, &bs.UserID, &bs.CategoryID, &bs.Month, &bs.Year, &bs.Amount,
			&bs.CreatedAt, &bs.UpdatedAt,
			&bs.CategoryName, &bs.CategoryColor, &bs.Spent,
		); err != nil {
			return nil, err
		}
		bs.Remaining = bs.Amount - bs.Spent
		bs.PercentUsed = calcPercent(bs.Spent, bs.Amount)
		results = append(results, bs)
	}
	return results, nil
}

// CopyBudgetsFromMonth copies all budgets from a source month/year to a
// destination month/year, skipping categories that already have a budget there.
// Returns the number of rows inserted.
func CopyBudgetsFromMonth(userID, fromMonth, fromYear, toMonth, toYear int) (int, error) {
	tag, err := database.Pool.Exec(context.Background(),
		`INSERT INTO budgets (user_id, category_id, month, year, amount)
		 SELECT $1, category_id, $4, $5, amount
		 FROM budgets
		 WHERE user_id = $1 AND month = $2 AND year = $3
		 ON CONFLICT (user_id, category_id, month, year) DO NOTHING`,
		userID, fromMonth, fromYear, toMonth, toYear)
	if err != nil {
		return 0, err
	}
	return int(tag.RowsAffected()), nil
}
