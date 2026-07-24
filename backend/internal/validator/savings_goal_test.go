package validator_test

import (
	"testing"

	"github.com/akrom/finance-backend/internal/model"
	"github.com/akrom/finance-backend/internal/validator"
)

func TestValidateSavingsGoalRequest(t *testing.T) {
	futureDate := "2030-12-31"

	tests := []struct {
		name    string
		req     model.SavingsGoalRequest
		wantErr bool
	}{
		{"valid goal", model.SavingsGoalRequest{Name: "Beli Laptop", TargetAmount: 15000000, Deadline: &futureDate}, false},
		{"empty name", model.SavingsGoalRequest{Name: "", TargetAmount: 5000000}, true},
		{"single char name", model.SavingsGoalRequest{Name: "A", TargetAmount: 5000000}, true},
		{"zero target amount", model.SavingsGoalRequest{Name: "Dana Darurat", TargetAmount: 0}, true},
		{"negative target amount", model.SavingsGoalRequest{Name: "Dana Darurat", TargetAmount: -1000}, true},
		{"invalid deadline format", model.SavingsGoalRequest{Name: "Liburan", TargetAmount: 10000000, Deadline: strPtr("31-12-2030")}, true},
		{"past deadline", model.SavingsGoalRequest{Name: "Liburan", TargetAmount: 10000000, Deadline: strPtr("2020-01-01")}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateSavingsGoalRequest(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSavingsGoalRequest() err = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func strPtr(s string) *string {
	return &s
}

func TestValidateDepositWithdraw(t *testing.T) {
	if err := validator.ValidateDepositWithdraw(500000, 1000000, 5000000, false); err != nil {
		t.Errorf("valid deposit failed: %v", err)
	}
	if err := validator.ValidateDepositWithdraw(-100, 1000000, 5000000, false); err == nil {
		t.Error("negative deposit should fail")
	}
	if err := validator.ValidateDepositWithdraw(2000000, 1000000, 5000000, true); err == nil {
		t.Error("withdrawal exceeding current amount should fail")
	}
}
