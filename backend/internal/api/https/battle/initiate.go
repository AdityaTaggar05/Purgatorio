package battle

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/AdityaTaggar05/Purgatorio/internal/domain/model"
	"github.com/AdityaTaggar05/Purgatorio/internal/domain/service"
	"github.com/AdityaTaggar05/Purgatorio/pkg/ctxkeys"
	"github.com/AdityaTaggar05/Purgatorio/pkg/response"
	"github.com/google/uuid"
)

func (h *BattleHandler) HandleInitiate(w http.ResponseWriter, r *http.Request) {
	var req model.InitiateRequest

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

	resp, err := h.Service.InitiateBattle(r.Context(), attackerID, defenderID)
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

	response.Created(w, resp, "battle initiated")
}
