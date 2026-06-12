package shop

import (
	"log/slog"

	"github.com/AdityaTaggar05/Purgatorio/internal/domain/service"
	"github.com/go-playground/validator/v10"
)

type ShopHandler struct {
	Logger    *slog.Logger
	Service   *service.ShopService
	Validator *validator.Validate
}

func NewHandler(logger *slog.Logger, svc *service.ShopService) *ShopHandler {
	return &ShopHandler{
		Logger:    logger,
		Service:   svc,
		Validator: validator.New(validator.WithRequiredStructEnabled()),
	}
}
