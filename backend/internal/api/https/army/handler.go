package army

import (
	"log/slog"

	"github.com/AdityaTaggar05/Purgatorio/internal/domain/service"
	"github.com/go-playground/validator/v10"
)

type ArmyHandler struct {
	Logger    *slog.Logger
	Service   *service.ArmyService
	Validator *validator.Validate
}

func NewHandler(logger *slog.Logger, svc *service.ArmyService) *ArmyHandler {
	return &ArmyHandler{
		Logger:    logger,
		Service:   svc,
		Validator: validator.New(validator.WithRequiredStructEnabled()),
	}
}
