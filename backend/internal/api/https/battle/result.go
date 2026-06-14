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

func (h *BattleHandler) HandleResult(w http.ResponseWriter, r *http.Request) {
	battleID, err := uuid.Parse(chi.URLParam(r, "battle_id"))
	if err != nil {
		response.BadRequest(r.Context(), w, fmt.Errorf("invalid battle id"))
		return
	}

	userID := r.Context().Value(ctxkeys.UserID).(uuid.UUID)
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

	_ = userID

	response.Success(w, replay, "battle result retrieved")
}
