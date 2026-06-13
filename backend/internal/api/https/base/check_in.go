package base

import (
	"net/http"

	"github.com/AdityaTaggar05/Purgatorio/internal/domain/model"
	"github.com/AdityaTaggar05/Purgatorio/pkg/ctxkeys"
	"github.com/AdityaTaggar05/Purgatorio/pkg/response"
	"github.com/google/uuid"
)

type CheckInResponseDTO struct {
	CompletedUpgrades []model.CheckInUpgrade `json:"completed_upgrades"`
}

func (h *BaseHandler) HandleCheckIn(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(ctxkeys.UserID).(uuid.UUID)

	result, err := h.Service.CheckIn(r.Context(), userID)
	if err != nil {
		response.Error(r.Context(), w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, CheckInResponseDTO{
		CompletedUpgrades: result.CompletedUpgrades,
	})
}
