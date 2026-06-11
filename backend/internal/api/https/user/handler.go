package user

import (
	"log/slog"

	"github.com/AdityaTaggar05/Purgatorio/internal/domain/service"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	Logger *slog.Logger
	Service *service.UserService
	Validator *validator.Validate
}

func NewHandler(logger *slog.Logger, service *service.UserService) *UserHandler {
	return &UserHandler{
		Logger: logger,
		Service: service,
		Validator: validator.New(validator.WithRequiredStructEnabled()),
	}
}
