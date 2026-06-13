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

type MoveBuildingRequestDTO struct {
	BuildingID string `json:"building_id" validate:"required"`
	FromX      int    `json:"from_x" validate:"min=0,max=29"`
	FromY      int    `json:"from_y" validate:"min=0,max=29"`
	ToX        int    `json:"to_x" validate:"min=0,max=29"`
	ToY        int    `json:"to_y" validate:"min=0,max=29"`
}

func (h *BaseHandler) HandleMoveBuilding(w http.ResponseWriter, r *http.Request) {
	var req MoveBuildingRequestDTO

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(r.Context(), w, fmt.Errorf("invalid request JSON"))
		return
	}

	if err := h.Validator.Struct(req); err != nil {
		response.ValidationFailed(r.Context(), w, err)
		return
	}

	userID := r.Context().Value(ctxkeys.UserID).(uuid.UUID)

	if err := h.Service.MoveBuilding(r.Context(), userID, req.BuildingID, req.FromX, req.FromY, req.ToX, req.ToY); err == nil {
		response.Success(w, nil, "building moved successfully")
	} else {
		switch {
		case errors.Is(err, service.ErrBuildingNotPlaced):
			response.BadRequest(r.Context(), w, err)
		case errors.Is(err, service.ErrPositionOccupied),
			errors.Is(err, service.ErrPositionOutOfBounds):
			response.BadRequest(r.Context(), w, err)
		default:
			response.InternalServerError(r.Context(), w, err)
		}
	}
}
