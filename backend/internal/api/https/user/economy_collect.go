package user

import (
	"net/http"

	"github.com/AdityaTaggar05/Purgatorio/pkg/ctxkeys"
	"github.com/AdityaTaggar05/Purgatorio/pkg/response"
	"github.com/google/uuid"
)

func (h *UserHandler) HandleEconomyCollect(w http.ResponseWriter, r *http.Request) {
	if data, err := h.Service.EconomyCollect(r.Context(), r.Context().Value(ctxkeys.UserID).(uuid.UUID)); err == nil {
		response.JSON(w, http.StatusOK, EconomyResponseDTO{
			Penitence:         data.Penitence,
			Grace:             data.Grace,
			MaxPenitence:      data.MaxPenitence,
			OverflowPenitence: data.CollectorPendingPenitence,
		})
	} else {
		response.Error(r.Context(), w, http.StatusInternalServerError, err)
	}
}
