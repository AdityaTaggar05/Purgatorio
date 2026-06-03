package auth

import (
	"log/slog"

	"github.com/AdityaTaggar05/Purgatorio/internal/domain/service"
	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	Logger *slog.Logger
	Service *service.AuthService
	Validator *validator.Validate
}

func NewHandler(logger *slog.Logger, service *service.AuthService) *AuthHandler {
	return &AuthHandler{
		Logger: logger,
		Service: service,
		Validator: validator.New(validator.WithRequiredStructEnabled()),
	}
}
