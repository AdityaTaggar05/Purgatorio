package user

import (
	"net/http"

	"github.com/AdityaTaggar05/Purgatorio/pkg/ctxkeys"
	"github.com/AdityaTaggar05/Purgatorio/pkg/response"
	"github.com/google/uuid"
)

type CombatResponseDTO struct {
	SinMeter int `json:"sin_meter"`
}

func (h *UserHandler) HandleGetCombat(w http.ResponseWriter, r *http.Request) {
	data, err := h.Service.GetCombat(r.Context(), r.Context().Value(ctxkeys.UserID).(uuid.UUID))
	if err != nil {
		response.InternalServerError(r.Context(), w, err)
		return
	}
	response.JSON(w, http.StatusOK, CombatResponseDTO{
		SinMeter: data.SinMeter,
	})
}
