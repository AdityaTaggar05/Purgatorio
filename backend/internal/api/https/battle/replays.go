package battle

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/AdityaTaggar05/Purgatorio/internal/domain/service"
	"github.com/AdityaTaggar05/Purgatorio/pkg/ctxkeys"
	"github.com/AdityaTaggar05/Purgatorio/pkg/response"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (h *BattleHandler) HandleAttacks(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(ctxkeys.UserID).(uuid.UUID)

	battles, err := h.Service.GetRecentBattles(r.Context(), userID, 20)
	if err != nil {
		response.InternalServerError(r.Context(), w, err)
		return
	}

	response.Success(w, battles, "attack history retrieved")
}

func (h *BattleHandler) HandleDefenses(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(ctxkeys.UserID).(uuid.UUID)

	battles, err := h.Service.GetRecentBattles(r.Context(), userID, 20)
	if err != nil {
		response.InternalServerError(r.Context(), w, err)
		return
	}

	response.Success(w, battles, "defense history retrieved")
}

func (h *BattleHandler) HandleReplay(w http.ResponseWriter, r *http.Request) {
	battleID, err := uuid.Parse(chi.URLParam(r, "battle_id"))
	if err != nil {
		response.BadRequest(r.Context(), w, fmt.Errorf("invalid battle id"))
		return
	}

	replay, err := h.Service.GetReplay(r.Context(), battleID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrBattleNotFound):
			response.NotFound(r.Context(), w, err)
		default:
			response.InternalServerError(r.Context(), w, err)
		}
		return
	}

	response.Success(w, replay, "replay retrieved")
}
