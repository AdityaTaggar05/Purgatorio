package response

import (
	"context"
	"encoding/json"
	"net/http"
	"slices"

	"github.com/AdityaTaggar05/Purgatorio/internal/api/https/middleware"
)

type Response struct {
	Success bool       `json:"success"`
	Data    any        `json:"data,omitempty"`
	Error   *ErrorData `json:"error,omitempty"`
	Message string     `json:"message,omitempty"`
}

type ErrorData struct {
	Code    string       `json:"code"`
	Message string       `json:"message"`
	Details []FieldError `json:"details,omitempty"`
}

func JSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := Response{
		Success: statusCode >= 200 && statusCode < 300,
		Data:    data,
	}

	json.NewEncoder(w).Encode(response)
}

func Success(w http.ResponseWriter, data any, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := Response{
		Success: true,
		Data:    data,
		Message: message,
	}

	json.NewEncoder(w).Encode(response)
}

func Created(w http.ResponseWriter, data any, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response := Response{
		Success: true,
		Data:    data,
		Message: message,
	}

	json.NewEncoder(w).Encode(response)
}

func Error(ctx context.Context, w http.ResponseWriter, statusCode int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// Logging error with full details
	if logCtx, ok := ctx.Value(middleware.LogContextKey).(*middleware.LogContext); ok {
		logCtx.Error = err
	}

	message := err.Error()

	// Redacting the error details for specific status codes
	if slices.Contains([]int{401, 403, 500}, statusCode) {
		message = http.StatusText(statusCode)
	}

	response := Response{
		Success: false,
		Error: &ErrorData{
			Code:    http.StatusText(statusCode),
			Message: message,
			// TODO: implement Details for api response as well
		},
	}

	json.NewEncoder(w).Encode(response)
}

func BadRequest(ctx context.Context, w http.ResponseWriter, err error) {
	Error(ctx, w, http.StatusBadRequest, err)
}

func NotFound(ctx context.Context, w http.ResponseWriter, err error) {
	Error(ctx, w, http.StatusNotFound, err)
}

func InternalServerError(ctx context.Context, w http.ResponseWriter, err error) {
	Error(ctx, w, http.StatusInternalServerError, err)
}

func Unauthorized(ctx context.Context, w http.ResponseWriter, err error) {
	Error(ctx, w, http.StatusUnauthorized, err)
}

func Forbidden(ctx context.Context, w http.ResponseWriter, err error) {
	Error(ctx, w, http.StatusForbidden, err)
}
