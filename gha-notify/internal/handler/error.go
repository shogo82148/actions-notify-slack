package handler

import (
	"context"
	"log/slog"
	"net/http"
)

type validationError struct {
	err error
}

func newValidationError(err error) error {
	return &validationError{err: err}
}

func (e *validationError) Error() string {
	return "handler: validation error: " + e.err.Error()
}

func (e *validationError) Unwrap() error {
	return e.err
}

func handleError(ctx context.Context, w http.ResponseWriter, err error) {
	switch err.(type) {
	case *validationError:
		slog.WarnContext(ctx, "validation error", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		slog.ErrorContext(ctx, "unexpected error", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
