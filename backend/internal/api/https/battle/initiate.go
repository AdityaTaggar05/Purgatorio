package battle

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/AdityaTaggar05/Purgatorio/internal/domain/service"
	"github.com/AdityaTaggar05/Purgatorio/pkg/ctxkeys"
	"github.com/AdityaTaggar05/Purgatorio/pkg/response"
	"github.com/google/uuid"
)

type initiateRequestDTO struct {
	DefenderID string `json:"defender_id" validate:"required,uuid"`
}

type initiateResponseDTO struct {
	BattleID     uuid.UUID `json:"battle_id"`
	DefenderName string    `json:"defender_name"`
}

func (h *BattleHandler) HandleInitiate(w http.ResponseWriter, r *http.Request) {
	var req initiateRequestDTO

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(r.Context(), w, fmt.Errorf("invalid request JSON"))
		return
	}

	if err := h.Validator.Struct(req); err != nil {
		response.ValidationFailed(r.Context(), w, err)
		return
	}

	attackerID := r.Context().Value(ctxkeys.UserID).(uuid.UUID)
	defenderID, err := uuid.Parse(req.DefenderID)
	if err != nil {
		response.BadRequest(r.Context(), w, fmt.Errorf("invalid defender id"))
		return
	}

	battleID, defenderName, err := h.Service.InitiateBattle(r.Context(), attackerID, defenderID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrCannotAttackSelf):
			response.BadRequest(r.Context(), w, err)
		case errors.Is(err, service.ErrDefenderNotFound):
			response.NotFound(r.Context(), w, err)
		case errors.Is(err, service.ErrTerraceLevelMismatch):
			response.BadRequest(r.Context(), w, err)
		case errors.Is(err, service.ErrDefenderShieldActive):
			response.BadRequest(r.Context(), w, err)
		default:
			response.InternalServerError(r.Context(), w, err)
		}
		return
	}

	response.Created(w, initiateResponseDTO{
		BattleID:     battleID,
		DefenderName: defenderName,
	}, "battle initiated")
}
