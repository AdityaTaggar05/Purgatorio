package auth

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/AdityaTaggar05/Purgatorio/internal/domain/service"
	"github.com/AdityaTaggar05/Purgatorio/pkg/response"
)

func (h *AuthHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		response.Unauthorized(r.Context(), w, fmt.Errorf("No refresh token"))
		return
	}

	if err := h.Service.Logout(r.Context(), cookie.Value); err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidRefreshTokenFormat):
			response.BadRequest(r.Context(), w, err)
		default:
			response.InternalServerError(r.Context(), w, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
