package auth

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/AdityaTaggar05/Purgatorio/internal/domain/service"
	"github.com/AdityaTaggar05/Purgatorio/pkg/response"
)

func (h *AuthHandler) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		response.Unauthorized(r.Context(), w, fmt.Errorf("No refresh token"))
		return
	}

	tokens, err := h.Service.Refresh(r.Context(), cookie.Value)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidRefreshTokenFormat), errors.Is(err, service.ErrInvalidRefreshToken):
			response.BadRequest(r.Context(), w, err)
		default:
			response.InternalServerError(r.Context(), w, err)
		}
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		Expires:  time.Now().Add(h.Service.Config.RefreshTTL),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/auth",
	})

	response.JSON(w, http.StatusOK, map[string]any{"access_token": tokens.AccessToken})
}
