package common

import (
	"context"
)

// ErrorHandler handles an error.
type ErrorHandler interface {
	Handle(err error)
	HandleContext(ctx context.Context, err error)
}

// NoopErrorHandler is an error handler that discards every error.
type NoopErrorHandler struct{}

func (NoopErrorHandler) Handle(_ error)                           {}
func (NoopErrorHandler) HandleContext(_ context.Context, _ error) {}

type ValidationError struct {
	Violates map[string][]string
}

func (ValidationError) Error() string {
	return "invalid input"
}

func (e ValidationError) Violations() map[string][]string {
	return e.Violates
}

func (ValidationError) Validation() bool {
	return true
}

func (ValidationError) ServiceError() bool {
	return true
}
