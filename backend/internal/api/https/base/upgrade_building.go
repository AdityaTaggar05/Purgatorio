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

type UpgradeBuildingRequestDTO struct {
	BuildingID string `json:"building_id" validate:"required"`
	X          int    `json:"x" validate:"min=0,max=29"`
	Y          int    `json:"y" validate:"min=0,max=29"`
}

func (h *BaseHandler) HandleUpgradeBuilding(w http.ResponseWriter, r *http.Request) {
	var req UpgradeBuildingRequestDTO

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(r.Context(), w, fmt.Errorf("invalid request JSON"))
		return
	}

	if err := h.Validator.Struct(req); err != nil {
		response.ValidationFailed(r.Context(), w, err)
		return
	}

	userID := r.Context().Value(ctxkeys.UserID).(uuid.UUID)

	if err := h.Service.UpgradeBuilding(r.Context(), userID, req.BuildingID, req.X, req.Y); err == nil {
		response.Created(w, nil, "upgrade started successfully")
	} else {
		switch {
		case errors.Is(err, service.ErrBuildingNotPlaced):
			response.BadRequest(r.Context(), w, err)
		case errors.Is(err, service.ErrUpgradeAlreadyActive):
			response.BadRequest(r.Context(), w, err)
		case errors.Is(err, service.ErrMaxLevelReached):
			response.BadRequest(r.Context(), w, err)
		case errors.Is(err, service.ErrInsufficientResources):
			response.BadRequest(r.Context(), w, err)
		default:
			response.InternalServerError(r.Context(), w, err)
		}
	}
}
