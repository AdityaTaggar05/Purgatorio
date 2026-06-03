package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AdityaTaggar05/Purgatorio/internal/domain/service"
	"github.com/AdityaTaggar05/Purgatorio/pkg/response"
)

type RefreshRequestDTO struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

func (h *AuthHandler) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequestDTO

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(r.Context(), w, fmt.Errorf("invalid request JSON"))
		return
	}

	if err := h.Validator.Struct(req); err != nil {
		response.Error(r.Context(), w, http.StatusUnprocessableEntity, err)
		return
	}

	tokens, err := h.Service.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		switch err {
		case service.ErrInvalidRefreshTokenFormat, service.ErrInvalidRefreshToken:
			response.BadRequest(r.Context(), w, err)
		default:
			response.InternalServerError(r.Context(), w, err)
		}
		return
	}

	response.JSON(w, http.StatusOK, tokens)
}
