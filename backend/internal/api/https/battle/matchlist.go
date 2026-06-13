package battle

import (
	"net/http"

	"github.com/AdityaTaggar05/Purgatorio/internal/domain/model"
	"github.com/AdityaTaggar05/Purgatorio/pkg/ctxkeys"
	"github.com/AdityaTaggar05/Purgatorio/pkg/response"
	"github.com/google/uuid"
)

func (h *BattleHandler) HandleMatchList(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(ctxkeys.UserID).(uuid.UUID)

	list, err := h.Service.GetMatchList(r.Context(), userID)
	if err != nil {
		response.InternalServerError(r.Context(), w, err)
		return
	}

	if list == nil {
		list = []model.MatchListEntry{}
	}

	response.Success(w, list, "match list retrieved")
}
