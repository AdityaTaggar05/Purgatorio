package battle

import (
	"log/slog"
	"net/http"

	"github.com/AdityaTaggar05/Purgatorio/internal/domain/service"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
)

type BattleHandler struct {
	Logger    *slog.Logger
	Service   *service.BattleService
	Validator *validator.Validate
	Upgrader  websocket.Upgrader
}

func NewHandler(logger *slog.Logger, svc *service.BattleService) *BattleHandler {
	return &BattleHandler{
		Logger:    logger,
		Service:   svc,
		Validator: validator.New(validator.WithRequiredStructEnabled()),
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}
}
