package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/akrom/finance-backend/internal/database"
	"github.com/akrom/finance-backend/internal/model"
)

func GetTransactions(userID int, filter model.TransactionFilter) ([]model.Transaction, int, error) {
	where := []string{"t.user_id = $1"}
	args := []interface{}{userID}
	argIdx := 2

	if filter.Month > 0 && filter.Year > 0 {
		where = append(where, fmt.Sprintf("EXTRACT(MONTH FROM t.date) = $%d", argIdx))
		args = append(args, filter.Month)
		argIdx++
		where = append(where, fmt.Sprintf("EXTRACT(YEAR FROM t.date) = $%d", argIdx))
		args = append(args, filter.Year)
		argIdx++
	} else if filter.Month > 0 {
		where = append(where, fmt.Sprintf("EXTRACT(MONTH FROM t.date) = $%d", argIdx))
		args = append(args, filter.Month)
		argIdx++
	} else if filter.Year > 0 {
		where = append(where, fmt.Sprintf("EXTRACT(YEAR FROM t.date) = $%d", argIdx))
		args = append(args, filter.Year)
		argIdx++
	}

	if filter.Type != "" {
		where = append(where, fmt.Sprintf("t.type = $%d", argIdx))
		args = append(args, filter.Type)
		argIdx++
	}

	if filter.CategoryID > 0 {
		where = append(where, fmt.Sprintf("t.category_id = $%d", argIdx))
		args = append(args, filter.CategoryID)
		argIdx++
	}

	whereClause := strings.Join(where, " AND ")

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM transactions t WHERE %s", whereClause)
	var total int
	err := database.Pool.QueryRow(context.Background(), countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count transactions: %w", err)
	}

	page := filter.Page
	if page < 1 {
		page = 1
	}
	limit := filter.Limit
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	query := fmt.Sprintf(`
		SELECT t.id, t.category_id, t.amount, t.description, t.date::text, t.type, t.created_at, t.updated_at,
			   c.name AS category_name, c.icon AS category_icon, c.color AS category_color
		FROM transactions t
		JOIN categories c ON c.id = t.category_id
		WHERE %s
		ORDER BY t.date DESC, t.id DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := database.Pool.Query(context.Background(), query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("query transactions: %w", err)
	}
	defer rows.Close()

	var transactions []model.Transaction
	for rows.Next() {
		var t model.Transaction
		if err := rows.Scan(&t.ID, &t.CategoryID, &t.Amount, &t.Description, &t.Date, &t.Type,
			&t.CreatedAt, &t.UpdatedAt, &t.CategoryName, &t.CategoryIcon, &t.CategoryColor); err != nil {
			return nil, 0, err
		}
		transactions = append(transactions, t)
	}

	return transactions, total, nil
}

func GetTransactionByID(id, userID int) (*model.Transaction, error) {
	query := `
		SELECT t.id, t.category_id, t.amount, t.description, t.date::text, t.type, t.created_at, t.updated_at,
			   c.name AS category_name, c.icon AS category_icon, c.color AS category_color
		FROM transactions t
		JOIN categories c ON c.id = t.category_id
		WHERE t.id = $1 AND t.user_id = $2
	`
	var t model.Transaction
	err := database.Pool.QueryRow(context.Background(), query, id, userID).Scan(
		&t.ID, &t.CategoryID, &t.Amount, &t.Description, &t.Date, &t.Type,
		&t.CreatedAt, &t.UpdatedAt, &t.CategoryName, &t.CategoryIcon, &t.CategoryColor)
	if err != nil {
		return nil, fmt.Errorf("get transaction %d: %w", id, err)
	}
	return &t, nil
}

func CreateTransaction(userID int, req model.TransactionRequest) (*model.Transaction, error) {
	query := `
		INSERT INTO transactions (user_id, category_id, amount, description, date, type)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, category_id, amount, description, date::text, type, created_at, updated_at
	`
	var t model.Transaction
	err := database.Pool.QueryRow(context.Background(), query,
		userID, req.CategoryID, req.Amount, req.Description, req.Date, req.Type).Scan(
		&t.ID, &t.CategoryID, &t.Amount, &t.Description, &t.Date, &t.Type,
		&t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("create transaction: %w", err)
	}
	return &t, nil
}

func UpdateTransaction(id, userID int, req model.TransactionRequest) (*model.Transaction, error) {
	query := `
		UPDATE transactions SET category_id=$1, amount=$2, description=$3, date=$4, type=$5, updated_at=NOW()
		WHERE id=$6 AND user_id=$7
		RETURNING id, category_id, amount, description, date::text, type, created_at, updated_at
	`
	var t model.Transaction
	err := database.Pool.QueryRow(context.Background(), query,
		req.CategoryID, req.Amount, req.Description, req.Date, req.Type, id, userID).Scan(
		&t.ID, &t.CategoryID, &t.Amount, &t.Description, &t.Date, &t.Type,
		&t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("update transaction %d: %w", id, err)
	}
	return &t, nil
}

func DeleteTransaction(id, userID int) error {
	_, err := database.Pool.Exec(context.Background(), "DELETE FROM transactions WHERE id = $1 AND user_id = $2", id, userID)
	if err != nil {
		return fmt.Errorf("delete transaction %d: %w", id, err)
	}
	return nil
}

// GetAllTransactionDates returns all distinct transaction dates for a user sorted
// ascending. Used to calculate consecutive-day streaks.
func GetAllTransactionDates(userID int) ([]string, error) {
	rows, err := database.Pool.Query(context.Background(),
		`SELECT DISTINCT date::text FROM transactions WHERE user_id = $1 ORDER BY date ASC`,
		userID)
	if err != nil {
		return nil, fmt.Errorf("get transaction dates: %w", err)
	}
	defer rows.Close()

	var dates []string
	for rows.Next() {
		var d string
		if err := rows.Scan(&d); err != nil {
			return nil, err
		}
		dates = append(dates, d)
	}
	return dates, nil
}

// BulkDeleteTransactions deletes transactions by IDs that belong to userID.
// It returns the number of rows actually deleted.
func BulkDeleteTransactions(userID int, ids []int) (int, error) {
	if len(ids) == 0 {
		return 0, nil
	}

	// Build parameterised IN clause: ($2,$3,...)
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids)+1)
	args[0] = userID
	for i, id := range ids {
		placeholders[i] = fmt.Sprintf("$%d", i+2)
		args[i+1] = id
	}

	query := fmt.Sprintf(
		"DELETE FROM transactions WHERE user_id = $1 AND id IN (%s)",
		strings.Join(placeholders, ","),
	)

	tag, err := database.Pool.Exec(context.Background(), query, args...)
	if err != nil {
		return 0, fmt.Errorf("bulk delete transactions: %w", err)
	}
	return int(tag.RowsAffected()), nil
}

func GetSummary(userID, month, year int) (*model.SummaryResponse, error) {
	query := `
		SELECT
			COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0) AS total_income,
			COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0) AS total_expense
		FROM transactions
		WHERE user_id = $1 AND EXTRACT(MONTH FROM date) = $2 AND EXTRACT(YEAR FROM date) = $3
	`
	var summary model.SummaryResponse
	err := database.Pool.QueryRow(context.Background(), query, userID, month, year).Scan(
		&summary.TotalIncome, &summary.TotalExpense)
	if err != nil {
		return nil, fmt.Errorf("get summary: %w", err)
	}
	summary.Balance = summary.TotalIncome - summary.TotalExpense

	catQuery := `
		SELECT c.id, c.name, c.icon, c.color,
			   COALESCE(SUM(t.amount), 0) AS total,
			   COUNT(t.id) AS count
		FROM categories c
		LEFT JOIN transactions t ON t.category_id = c.id
			AND t.user_id = $1
			AND EXTRACT(MONTH FROM t.date) = $2 AND EXTRACT(YEAR FROM t.date) = $3
		WHERE c.user_id = $1
		GROUP BY c.id, c.name, c.icon, c.color
		ORDER BY total DESC
	`
	rows, err := database.Pool.Query(context.Background(), catQuery, userID, month, year)
	if err != nil {
		return nil, fmt.Errorf("get summary by category: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var cs model.CategorySummary
		if err := rows.Scan(&cs.CategoryID, &cs.CategoryName, &cs.CategoryIcon,
			&cs.CategoryColor, &cs.Total, &cs.Count); err != nil {
			return nil, err
		}
		summary.ByCategory = append(summary.ByCategory, cs)
	}

	return &summary, nil
}

func GetYearlyTrend(userID, year int) (*model.YearlyTrendResponse, error) {
	query := `
		SELECT
			EXTRACT(MONTH FROM date) AS month,
			COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0) AS income,
			COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0) AS expense
		FROM transactions
		WHERE user_id = $1 AND EXTRACT(YEAR FROM date) = $2
		GROUP BY EXTRACT(MONTH FROM date)
		ORDER BY month
	`
	rows, err := database.Pool.Query(context.Background(), query, userID, year)
	if err != nil {
		return nil, fmt.Errorf("yearly trend: %w", err)
	}
	defer rows.Close()

	monthMap := make(map[int]model.MonthlyTrend)
	var totalIncome, totalExpense int64

	for rows.Next() {
		var mt model.MonthlyTrend
		if err := rows.Scan(&mt.Month, &mt.Income, &mt.Expense); err != nil {
			return nil, err
		}
		mt.Balance = mt.Income - mt.Expense
		monthMap[mt.Month] = mt
		totalIncome += mt.Income
		totalExpense += mt.Expense
	}

	// Build full 12-month array
	months := make([]model.MonthlyTrend, 12)
	for i := 1; i <= 12; i++ {
		if mt, ok := monthMap[i]; ok {
			months[i-1] = mt
		} else {
			months[i-1] = model.MonthlyTrend{Month: i, Income: 0, Expense: 0, Balance: 0}
		}
	}

	return &model.YearlyTrendResponse{
		Year:         year,
		TotalIncome:  totalIncome,
		TotalExpense: totalExpense,
		Balance:      totalIncome - totalExpense,
		Months:       months,
	}, nil
}

func GetCategoryTrend(userID, month, year int) (*model.CategoryTrendResponse, error) {
	// Get all expense categories
	catQuery := `
		SELECT c.id, c.name, c.color
		FROM categories c
		WHERE c.user_id = $1 AND c.type = 'expense'
		ORDER BY c.name
	`
	catRows, err := database.Pool.Query(context.Background(), catQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("category trend cats: %w", err)
	}
	defer catRows.Close()

	type catInfo struct {
		id    int
		name  string
		color string
	}
	var cats []catInfo
	for catRows.Next() {
		var c catInfo
		if err := catRows.Scan(&c.id, &c.name, &c.color); err != nil {
			return nil, err
		}
		cats = append(cats, c)
	}
	catRows.Close()

	// Get weekly expense totals per category
	weekQuery := `
		SELECT
			category_id,
			CEIL(EXTRACT(DAY FROM date) / 7.0)::int AS week,
			COALESCE(SUM(amount), 0) AS total
		FROM transactions
		WHERE user_id = $1
			AND EXTRACT(MONTH FROM date) = $2
			AND EXTRACT(YEAR FROM date) = $3
			AND type = 'expense'
		GROUP BY category_id, CEIL(EXTRACT(DAY FROM date) / 7.0)::int
		ORDER BY category_id, week
	`
	weekRows, err := database.Pool.Query(context.Background(), weekQuery, userID, month, year)
	if err != nil {
		return nil, fmt.Errorf("category trend weeks: %w", err)
	}
	defer weekRows.Close()

	weekMap := make(map[int]map[int]int64) // catId -> week -> amount
	for weekRows.Next() {
		var catID, week int
		var total int64
		if err := weekRows.Scan(&catID, &week, &total); err != nil {
			return nil, err
		}
		if _, ok := weekMap[catID]; !ok {
			weekMap[catID] = make(map[int]int64)
		}
		weekMap[catID][week] = total
	}

	// Get total and count per category
	totalQuery := `
		SELECT category_id, COALESCE(SUM(amount), 0) AS total, COUNT(*) AS count
		FROM transactions
		WHERE user_id = $1
			AND EXTRACT(MONTH FROM date) = $2
			AND EXTRACT(YEAR FROM date) = $3
			AND type = 'expense'
		GROUP BY category_id
	`
	totalRows, err := database.Pool.Query(context.Background(), totalQuery, userID, month, year)
	if err != nil {
		return nil, fmt.Errorf("category trend totals: %w", err)
	}
	defer totalRows.Close()

	totals := make(map[int]struct {
		total int64
		count int
	})
	for totalRows.Next() {
		var catID int
		var total int64
		var count int
		if err := totalRows.Scan(&catID, &total, &count); err != nil {
			return nil, err
		}
		totals[catID] = struct {
			total int64
			count int
		}{total, count}
	}

	// Build response
	var categories []model.CategoryWeeklyTrend
	for _, c := range cats {
		weekly := make([]model.WeeklyTrend, 4)
		for w := 1; w <= 4; w++ {
			weekly[w-1] = model.WeeklyTrend{Week: w, Expense: weekMap[c.id][w]}
		}
		t := totals[c.id]
		categories = append(categories, model.CategoryWeeklyTrend{
			CategoryID:    c.id,
			CategoryName:  c.name,
			CategoryColor: c.color,
			Total:         t.total,
			Count:         t.count,
			Weekly:        weekly,
		})
	}

	if categories == nil {
		categories = []model.CategoryWeeklyTrend{}
	}

	return &model.CategoryTrendResponse{
		Month:      month,
		Year:       year,
		Categories: categories,
	}, nil
}
