package repository

import (
	"context"
	"fmt"

	"github.com/akrom/finance-backend/internal/database"
	"github.com/akrom/finance-backend/internal/model"
)

// GetAllTransactionsForML fetches all transactions for a user (for ML training)
func GetAllTransactionsForML(userID int) ([]model.MLTrainTransaction, error) {
	query := `
		SELECT t.description, t.category_id, t.type
		FROM transactions t
		WHERE t.user_id = $1 AND t.description IS NOT NULL AND t.description != ''
		ORDER BY t.date ASC
	`

	rows, err := database.Pool.Query(context.Background(), query, userID)
	if err != nil {
		return nil, fmt.Errorf("ml get transactions: %w", err)
	}
	defer rows.Close()

	var transactions []model.MLTrainTransaction
	for rows.Next() {
		var t model.MLTrainTransaction
		if err := rows.Scan(&t.Description, &t.CategoryID, &t.Type); err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}

	return transactions, nil
}

// GetAllCategoriesForML fetches all categories for a user (for ML training)
func GetAllCategoriesForML(userID int) ([]model.MLTrainCategory, error) {
	query := `
		SELECT id, name, type
		FROM categories
		WHERE user_id = $1
		ORDER BY id
	`

	rows, err := database.Pool.Query(context.Background(), query, userID)
	if err != nil {
		return nil, fmt.Errorf("ml get categories: %w", err)
	}
	defer rows.Close()

	var categories []model.MLTrainCategory
	for rows.Next() {
		var c model.MLTrainCategory
		if err := rows.Scan(&c.ID, &c.Name, &c.Type); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	return categories, nil
}
