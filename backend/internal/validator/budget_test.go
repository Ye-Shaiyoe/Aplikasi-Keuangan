package validator_test

import (
	"testing"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/validator"
)

func TestValidateBudgetRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     model.BudgetRequest
		wantErr bool
	}{
		{"valid budget", model.BudgetRequest{CategoryID: 1, Month: 6, Year: 2024, Amount: 1000000}, false},
		{"zero category", model.BudgetRequest{CategoryID: 0, Month: 6, Year: 2024, Amount: 1000000}, true},
		{"invalid month 0", model.BudgetRequest{CategoryID: 1, Month: 0, Year: 2024, Amount: 1000000}, true},
		{"invalid month 13", model.BudgetRequest{CategoryID: 1, Month: 13, Year: 2024, Amount: 1000000}, true},
		{"invalid year 2010", model.BudgetRequest{CategoryID: 1, Month: 6, Year: 2010, Amount: 1000000}, true},
		{"negative amount", model.BudgetRequest{CategoryID: 1, Month: 6, Year: 2024, Amount: -500}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateBudgetRequest(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateBudgetRequest() err = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateCopyBudgetRequest(t *testing.T) {
	if err := validator.ValidateCopyBudgetRequest(1, 2024, 2, 2024); err != nil {
		t.Errorf("expected valid copy request, got %v", err)
	}
	if err := validator.ValidateCopyBudgetRequest(5, 2024, 4, 2024); err == nil {
		t.Error("expected error when copying backwards, got nil")
	}
}

func TestValidateAlertThreshold(t *testing.T) {
	if err := validator.ValidateAlertThreshold(80); err != nil {
		t.Errorf("threshold 80 should be valid, got %v", err)
	}
	if err := validator.ValidateAlertThreshold(0); err == nil {
		t.Error("threshold 0 should fail")
	}
	if err := validator.ValidateAlertThreshold(105); err == nil {
		t.Error("threshold 105 should fail")
	}
}
