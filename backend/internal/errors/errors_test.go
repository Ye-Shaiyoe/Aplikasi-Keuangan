package errors_test

import (
	"fmt"
	"testing"

	apperr "github.com/akrom/finance-backend/internal/errors"
)

func TestErrorConstructorsAndCheckers(t *testing.T) {
	nf := apperr.NewNotFound("Transaction", 42)
	if !apperr.IsNotFound(nf) {
		t.Error("IsNotFound should be true for NotFoundError")
	}
	if nf.Error() != "Transaction with id 42 not found" {
		t.Errorf("unexpected error string: %s", nf.Error())
	}

	ve := apperr.NewValidation("amount", "must be positive")
	if !apperr.IsValidation(ve) {
		t.Error("IsValidation should be true for ValidationError")
	}

	ce := apperr.NewConflict("email already registered")
	if !apperr.IsConflict(ce) {
		t.Error("IsConflict should be true for ConflictError")
	}

	ue := apperr.NewUnauthorized("invalid token")
	if !apperr.IsUnauthorized(ue) {
		t.Error("IsUnauthorized should be true for UnauthorizedError")
	}

	ie := apperr.NewInternal("db query", fmt.Errorf("connection refused"))
	if ie.Unwrap() == nil {
		t.Error("InternalError Unwrap() should return inner error")
	}
}
