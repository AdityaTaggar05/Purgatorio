package army

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

type DetrainTroopsRequestDTO struct {
	TroopID string `json:"troop_id" validate:"required"`
	Count   int    `json:"count" validate:"required,min=1"`
}

func (h *ArmyHandler) HandleDetrainTroops(w http.ResponseWriter, r *http.Request) {
	var req DetrainTroopsRequestDTO

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(r.Context(), w, fmt.Errorf("invalid request JSON"))
		return
	}

	if err := h.Validator.Struct(req); err != nil {
		response.ValidationFailed(r.Context(), w, err)
		return
	}

	userID := r.Context().Value(ctxkeys.UserID).(uuid.UUID)

	if err := h.Service.DetrainTroops(r.Context(), userID, req.TroopID, req.Count); err == nil {
		response.Created(w, nil, "troops detrained successfully")
	} else {
		switch {
		case errors.Is(err, service.ErrTroopNotFound):
			response.NotFound(r.Context(), w, err)
		case errors.Is(err, service.ErrInsufficientTroops):
			response.BadRequest(r.Context(), w, err)
		case errors.Is(err, service.ErrInsufficientResources):
			response.BadRequest(r.Context(), w, err)
		default:
			response.InternalServerError(r.Context(), w, err)
		}
	}
}
