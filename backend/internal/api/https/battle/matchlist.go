package battle

import (
	"net/http"

	"github.com/AdityaTaggar05/Purgatorio/pkg/ctxkeys"
	"github.com/AdityaTaggar05/Purgatorio/pkg/response"
	"github.com/google/uuid"
)

type matchListEntryDTO struct {
	UserID       uuid.UUID `json:"user_id"`
	Username     string    `json:"username"`
	TerraceLevel int       `json:"terrace_level"`
}

func (h *BattleHandler) HandleMatchList(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(ctxkeys.UserID).(uuid.UUID)

	list, err := h.Service.GetMatchList(r.Context(), userID)
	if err != nil {
		response.InternalServerError(r.Context(), w, err)
		return
	}

	entries := make([]matchListEntryDTO, 0, len(list))
	for _, m := range list {
		entries = append(entries, matchListEntryDTO{
			UserID:       m.UserID,
			Username:     m.Username,
			TerraceLevel: m.TerraceLevel,
		})
	}

	response.Success(w, entries, "match list retrieved")
}
