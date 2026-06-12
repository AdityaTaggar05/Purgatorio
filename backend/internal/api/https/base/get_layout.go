package base

import (
	"net/http"

	"github.com/AdityaTaggar05/Purgatorio/internal/domain/model"
	"github.com/AdityaTaggar05/Purgatorio/pkg/ctxkeys"
	"github.com/AdityaTaggar05/Purgatorio/pkg/response"
	"github.com/google/uuid"
)

type LayoutResponseDTO struct {
	Buildings []model.PlacedBuildingResponse `json:"buildings"`
	GridW     int                            `json:"grid_w"`
	GridH     int                            `json:"grid_h"`
}

func (h *BaseHandler) HandleGetLayout(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(ctxkeys.UserID).(uuid.UUID)

	layout, err := h.Service.GetLayout(r.Context(), userID)
	if err != nil {
		response.Error(r.Context(), w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, LayoutResponseDTO{
		Buildings: layout.Buildings,
		GridW:     layout.GridW,
		GridH:     layout.GridH,
	})
}
