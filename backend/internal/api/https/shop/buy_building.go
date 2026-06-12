package shop

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/AdityaTaggar05/Purgatorio/internal/domain/service"
	"github.com/AdityaTaggar05/Purgatorio/pkg/ctxkeys"
	"github.com/AdityaTaggar05/Purgatorio/pkg/response"
	"github.com/google/uuid"
)

type BuyBuildingRequestDTO struct {
	BuildingID string `json:"building_id" validate:"required"`
}

func (h *ShopHandler) HandleBuyBuilding(w http.ResponseWriter, r *http.Request) {
	var req BuyBuildingRequestDTO

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(r.Context(), w, err)
		return
	}

	if err := h.Validator.Struct(req); err != nil {
		response.ValidationFailed(r.Context(), w, err)
		return
	}

	userID := r.Context().Value(ctxkeys.UserID).(uuid.UUID)

	if err := h.Service.BuyBuilding(r.Context(), userID, req.BuildingID); err == nil {
		response.Success(w, nil, "building purchased successfully")
	} else {
		switch {
		case errors.Is(err, service.ErrBuildingNotFound):
			response.NotFound(r.Context(), w, err)
		case errors.Is(err, service.ErrBuildingLimitReached),
			errors.Is(err, service.ErrInsufficientResources):
			response.BadRequest(r.Context(), w, err)
		default:
			response.InternalServerError(r.Context(), w, err)
		}
	}
}
