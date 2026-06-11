package response

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/AdityaTaggar05/Purgatorio/pkg/ctxkeys"
	"github.com/go-playground/validator/v10"
)

// FieldError is the consumer-facing shape for a single validation failure.
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationFailed writes a 422 with clean FieldErrors to the consumer.
func ValidationFailed(ctx context.Context, w http.ResponseWriter, err error) {
	// Preserve raw validator error for internal logging
	if logCtx, ok := ctx.Value(ctxkeys.Log).(*LogContext); ok {
		logCtx.Error = err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnprocessableEntity)

	var fields []FieldError
	if ve, ok := errors.AsType[validator.ValidationErrors](err); ok {
		fields = make([]FieldError, len(ve))
		for i, fe := range ve {
			fields[i] = FieldError{
				Field:   strings.ToLower(fe.Field()),
				Message: tagToMessage(fe),
			}
		}
	}

	json.NewEncoder(w).Encode(Response{
		Success: false,
		Error: &ErrorData{
			Code:    http.StatusText(http.StatusUnprocessableEntity),
			Message: "one or more fields failed validation",
			Details: fields,
		},
	})
}

func tagToMessage(fe validator.FieldError) string {
	field := strings.ToLower(fe.Field())
	switch fe.Tag() {
	case "required":
		return field + " is required"
	case "email":
		return "must be a valid email address"
	case "min":
		return fmt.Sprintf("must be at least %s characters", fe.Param())
	case "max":
		return fmt.Sprintf("must be at most %s characters", fe.Param())
	case "alphanum":
		return "must contain only letters and numbers"
	default:
		return fmt.Sprintf("failed %s validation", fe.Tag())
	}
}
