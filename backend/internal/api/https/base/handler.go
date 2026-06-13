package base

import (
	"log/slog"

	"github.com/AdityaTaggar05/Purgatorio/internal/domain/service"
	"github.com/go-playground/validator/v10"
)

type BaseHandler struct {
	Logger    *slog.Logger
	Service   *service.BaseService
	Validator *validator.Validate
}

func NewHandler(logger *slog.Logger, svc *service.BaseService) *BaseHandler {
	return &BaseHandler{
		Logger:    logger,
		Service:   svc,
		Validator: validator.New(validator.WithRequiredStructEnabled()),
	}
}
