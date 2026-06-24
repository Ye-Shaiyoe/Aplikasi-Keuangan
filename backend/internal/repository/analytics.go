package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/akrom/finance-backend/internal/database"
	"github.com/akrom/finance-backend/internal/model"
)

type AnalyticsRepository struct{}

func GetHistoricalMetrics(userID int, limitMonths int) ([]model.MonthlyMetric, error) {
	// Past limitMonths including current month
	query := fmt.Sprintf(`
		SELECT
			EXTRACT(MONTH FROM date) AS m,
			EXTRACT(YEAR FROM date) AS y,
			COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0) AS income,
			COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0) AS expense
		FROM transactions
		WHERE user_id = $1 AND date >= date_trunc('month', NOW() - INTERVAL '%d month')
		GROUP BY EXTRACT(YEAR FROM date), EXTRACT(MONTH FROM date)
		ORDER BY y ASC, m ASC
	`, limitMonths-1)

	rows, err := database.Pool.Query(context.Background(), query, userID)
	if err != nil {
		return nil, fmt.Errorf("historical metrics: %w", err)
	}
	defer rows.Close()

	monthNames := []string{"", "Jan", "Feb", "Mar", "Apr", "Mei", "Jun", "Jul", "Agu", "Sep", "Okt", "Nov", "Des"}
	var metrics []model.MonthlyMetric

	for rows.Next() {
		var m, y int
		var inc, exp int64
		if err := rows.Scan(&m, &y, &inc, &exp); err != nil {
			return nil, err
		}

		metrics = append(metrics, model.MonthlyMetric{
			MonthName: fmt.Sprintf("%s %d", monthNames[m], y),
			Month:     m,
			Year:      y,
			Income:    inc,
			Expense:   exp,
		})
	}

	// If no metrics found or we want to make sure the current month exists
	if len(metrics) == 0 {
		now := time.Now()
		metrics = append(metrics, model.MonthlyMetric{
			MonthName: fmt.Sprintf("%s %d", monthNames[int(now.Month())], now.Year()),
			Month:     int(now.Month()),
			Year:      now.Year(),
			Income:    0,
			Expense:   0,
		})
	}

	return metrics, nil
}

type BudgetCompare struct {
	CategoryID   int
	BudgetAmount int64
	ActualAmount int64
}

func GetBudgetComparison(userID, month, year int) ([]BudgetCompare, error) {
	query := `
		SELECT
			b.category_id,
			b.amount AS budget_amount,
			COALESCE(SUM(t.amount), 0) AS actual_amount
		FROM budgets b
		LEFT JOIN transactions t ON t.category_id = b.category_id
			AND t.user_id = b.user_id
			AND EXTRACT(MONTH FROM t.date) = b.month
			AND EXTRACT(YEAR FROM t.date) = b.year
			AND t.type = 'expense'
		WHERE b.user_id = $1 AND b.month = $2 AND b.year = $3
		GROUP BY b.category_id, b.amount
	`

	rows, err := database.Pool.Query(context.Background(), query, userID, month, year)
	if err != nil {
		return nil, fmt.Errorf("budget comparison: %w", err)
	}
	defer rows.Close()

	var comparisons []BudgetCompare
	for rows.Next() {
		var bc BudgetCompare
		if err := rows.Scan(&bc.CategoryID, &bc.BudgetAmount, &bc.ActualAmount); err != nil {
			return nil, err
		}
		comparisons = append(comparisons, bc)
	}

	return comparisons, nil
}

type CategorySpending struct {
	CategoryName string
	TotalAmount  int64
	Count        int
}

func GetTopExpenses(userID, month, year int) ([]CategorySpending, error) {
	query := `
		SELECT c.name, COALESCE(SUM(t.amount), 0) as total, COUNT(t.id) as cnt
		FROM categories c
		JOIN transactions t ON t.category_id = c.id
		WHERE t.user_id = $1 AND t.type = 'expense'
		  AND EXTRACT(MONTH FROM t.date) = $2
		  AND EXTRACT(YEAR FROM t.date) = $3
		GROUP BY c.name
		ORDER BY total DESC
		LIMIT 5
	`

	rows, err := database.Pool.Query(context.Background(), query, userID, month, year)
	if err != nil {
		return nil, fmt.Errorf("top expenses: %w", err)
	}
	defer rows.Close()

	var spendings []CategorySpending
	for rows.Next() {
		var cs CategorySpending
		if err := rows.Scan(&cs.CategoryName, &cs.TotalAmount, &cs.Count); err != nil {
			return nil, err
		}
		spendings = append(spendings, cs)
	}
	return spendings, nil
}
