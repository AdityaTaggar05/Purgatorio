package battle

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/AdityaTaggar05/Purgatorio/internal/domain/service"
	"github.com/AdityaTaggar05/Purgatorio/internal/engine"
	"github.com/AdityaTaggar05/Purgatorio/pkg/response"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

type resultResponseDTO struct {
	Ticks      []engine.TickResult      `json:"ticks"`
	Result     engine.BattleResult      `json:"result"`
	Deployment []engine.TroopDeployment `json:"deployment"`
}

func (h *BattleHandler) HandleResult(w http.ResponseWriter, r *http.Request) {
	battleID, err := uuid.Parse(chi.URLParam(r, "battle_id"))
	if err != nil {
		response.BadRequest(r.Context(), w, fmt.Errorf("invalid battle id"))
		return
	}

	simResult, err := h.Service.GetReplay(r.Context(), battleID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrBattleNotFound) || errors.Is(err, service.ErrReplayNotFound):
			response.NotFound(r.Context(), w, err)
		default:
			response.InternalServerError(r.Context(), w, err)
		}
		return
	}

	response.Success(w, resultResponseDTO{
		Ticks:      simResult.Ticks,
		Result:     simResult.Result,
		Deployment: simResult.Deployment,
	}, "battle result retrieved")
}
