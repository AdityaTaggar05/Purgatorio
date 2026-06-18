package user

import (
	"net/http"

	"github.com/AdityaTaggar05/Purgatorio/pkg/ctxkeys"
	"github.com/AdityaTaggar05/Purgatorio/pkg/response"
	"github.com/google/uuid"
)

type EconomyResponseDTO struct {
	Penitence              int `json:"penitence"`
	Grace                  int `json:"grace"`
	MaxPenitence           int `json:"max_penitence"`
	OverflowPenitence      int `json:"overflow_penitence"`
}

func (h *UserHandler) HandleGetEconomy(w http.ResponseWriter, r *http.Request) {
	if data, err := h.Service.GetEconomy(r.Context(), r.Context().Value(ctxkeys.UserID).(uuid.UUID)); err == nil {
		response.JSON(w, http.StatusOK, EconomyResponseDTO{
			Penitence:         data.Penitence,
			Grace:             data.Grace,
			MaxPenitence:      data.MaxPenitence,
			OverflowPenitence: data.CollectorPendingPenitence,
		})
	} else {
		response.InternalServerError(r.Context(), w, err)
	}
}
