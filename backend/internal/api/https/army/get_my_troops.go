package army

import (
	"net/http"

	"github.com/AdityaTaggar05/Purgatorio/internal/domain/model"
	"github.com/AdityaTaggar05/Purgatorio/pkg/ctxkeys"
	"github.com/AdityaTaggar05/Purgatorio/pkg/response"
	"github.com/google/uuid"
)

func (h *ArmyHandler) HandleGetMyTroops(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(ctxkeys.UserID).(uuid.UUID)

	result, err := h.Service.GetMyTroops(r.Context(), userID)
	if err != nil {
		response.Error(r.Context(), w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, model.MyTroopsResponse{
		Troops:       result.Troops,
		UsedCapacity: result.UsedCapacity,
		MaxCapacity:  result.MaxCapacity,
	})
}
