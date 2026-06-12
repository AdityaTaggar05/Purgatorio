package base

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/AdityaTaggar05/Purgatorio/internal/domain/service"
	"github.com/AdityaTaggar05/Purgatorio/pkg/ctxkeys"
	"github.com/AdityaTaggar05/Purgatorio/pkg/response"
	"github.com/google/uuid"
)

type RemoveBuildingRequestDTO struct {
	BuildingID string `json:"building_id" validate:"required"`
	X          int    `json:"x" validate:"min=0,max=29"`
	Y          int    `json:"y" validate:"min=0,max=29"`
}

func (h *BaseHandler) HandleRemoveBuilding(w http.ResponseWriter, r *http.Request) {
	var req RemoveBuildingRequestDTO

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(r.Context(), w, fmt.Errorf("invalid request JSON"))
		return
	}

	if err := h.Validator.Struct(req); err != nil {
		response.ValidationFailed(r.Context(), w, err)
		return
	}

	userID := r.Context().Value(ctxkeys.UserID).(uuid.UUID)

	if err := h.Service.RemoveBuilding(r.Context(), userID, req.BuildingID, req.X, req.Y); err == nil {
		response.Success(w, nil, "building removed successfully")
	} else {
		switch {
		case errors.Is(err, service.ErrBuildingNotPlaced):
			response.BadRequest(r.Context(), w, err)
		default:
			response.InternalServerError(r.Context(), w, err)
		}
	}
}
