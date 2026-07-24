package validator_test

import (
	"testing"
	"time"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/validator"
)

func TestValidateTransactionRequest(t *testing.T) {
	today := time.Now().Format("2006-01-02")

	tests := []struct {
		name    string
		req     model.TransactionRequest
		wantErr bool
	}{
		{
			name: "valid income",
			req: model.TransactionRequest{
				CategoryID:  1,
				Amount:      1000,
				Description: "Salary",
				Date:        today,
				Type:        "income",
			},
			wantErr: false,
		},
		{
			name: "valid expense",
			req: model.TransactionRequest{
				CategoryID:  2,
				Amount:      500,
				Description: "Coffee",
				Date:        today,
				Type:        "expense",
			},
			wantErr: false,
		},
		{
			name: "negative amount",
			req: model.TransactionRequest{
				CategoryID: 1, Amount: -100, Date: today, Type: "income",
			},
			wantErr: true,
		},
		{
			name: "zero amount",
			req: model.TransactionRequest{
				CategoryID: 1, Amount: 0, Date: today, Type: "income",
			},
			wantErr: true,
		},
		{
			name: "invalid type",
			req: model.TransactionRequest{
				CategoryID: 1, Amount: 100, Date: today, Type: "investment",
			},
			wantErr: true,
		},
		{
			name: "zero category",
			req: model.TransactionRequest{
				CategoryID: 0, Amount: 100, Date: today, Type: "income",
			},
			wantErr: true,
		},
		{
			name: "invalid date format",
			req: model.TransactionRequest{
				CategoryID: 1, Amount: 100, Date: "15/01/2024", Type: "income",
			},
			wantErr: true,
		},
		{
			name: "date too far in past",
			req: model.TransactionRequest{
				CategoryID: 1, Amount: 100, Date: "2010-01-01", Type: "income",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateTransactionRequest(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTransactionRequest() err = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateBulkDeleteRequest(t *testing.T) {
	tests := []struct {
		name    string
		ids     []int
		wantErr bool
	}{
		{"valid ids", []int{1, 2, 3}, false},
		{"empty list", []int{}, true},
		{"invalid id zero", []int{1, 0, 3}, true},
		{"duplicate ids", []int{1, 2, 2}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateBulkDeleteRequest(tt.ids)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateBulkDeleteRequest() err = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}
