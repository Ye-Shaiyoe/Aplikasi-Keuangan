package errors

import "fmt"

// NotFoundError is returned when a requested resource does not exist or does not
// belong to the requesting user.
type NotFoundError struct {
	Resource string
	ID       int
}

func (e *NotFoundError) Error() string {
	if e.ID != 0 {
		return fmt.Sprintf("%s with id %d not found", e.Resource, e.ID)
	}
	return fmt.Sprintf("%s not found", e.Resource)
}

// ValidationError is returned when input data fails business-rule validation.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
	}
	return fmt.Sprintf("validation error: %s", e.Message)
}

// ConflictError is returned when an operation would violate a uniqueness constraint.
type ConflictError struct {
	Message string
}

func (e *ConflictError) Error() string {
	return fmt.Sprintf("conflict: %s", e.Message)
}

// UnauthorizedError is returned when a user attempts to access a resource they
// do not own or do not have permission to access.
type UnauthorizedError struct {
	Message string
}

func (e *UnauthorizedError) Error() string {
	return fmt.Sprintf("unauthorized: %s", e.Message)
}

// InternalError wraps an underlying error that is not safe to expose to clients.
type InternalError struct {
	Op  string // operation that failed
	Err error
}

func (e *InternalError) Error() string {
	return fmt.Sprintf("internal error during %s: %v", e.Op, e.Err)
}

func (e *InternalError) Unwrap() error { return e.Err }

// Constructors — convenience functions for common error creation.

func NewNotFound(resource string, id int) *NotFoundError {
	return &NotFoundError{Resource: resource, ID: id}
}

func NewValidation(field, message string) *ValidationError {
	return &ValidationError{Field: field, Message: message}
}

func NewConflict(message string) *ConflictError {
	return &ConflictError{Message: message}
}

func NewUnauthorized(message string) *UnauthorizedError {
	return &UnauthorizedError{Message: message}
}

func NewInternal(op string, err error) *InternalError {
	return &InternalError{Op: op, Err: err}
}

// Type-assertion helpers.

func IsNotFound(err error) bool {
	_, ok := err.(*NotFoundError)
	return ok
}

func IsValidation(err error) bool {
	_, ok := err.(*ValidationError)
	return ok
}

func IsConflict(err error) bool {
	_, ok := err.(*ConflictError)
	return ok
}

func IsUnauthorized(err error) bool {
	_, ok := err.(*UnauthorizedError)
	return ok
}
