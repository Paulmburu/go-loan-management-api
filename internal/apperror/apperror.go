package apperror

import "errors"

// Predefined application-level errors used across the app.
var (
	ErrNotFound   = errors.New("not found")
	ErrConflict   = errors.New("conflict")
	ErrValidation = errors.New("validation error")
	ErrInternal   = errors.New("internal error")
)
