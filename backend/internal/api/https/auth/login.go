package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/AdityaTaggar05/Purgatorio/internal/domain/model"
	"github.com/AdityaTaggar05/Purgatorio/internal/domain/service"
	"github.com/AdityaTaggar05/Purgatorio/pkg/response"
)

type LoginRequestDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponseDTO struct {
	User        model.User `json:"user"`
	AccessToken string     `json:"access_token"`
}

func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequestDTO

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(r.Context(), w, fmt.Errorf("invalid request JSON"))
		return
	}

	if err := h.Validator.Struct(req); err != nil {
		response.ValidationFailed(r.Context(), w, err)
		return
	}

	if user, tokens, err := h.Service.Login(r.Context(), req.Email, req.Password); err == nil {
		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    tokens.RefreshToken,
			Expires:  time.Now().Add(h.Service.Config.RefreshTTL),
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
			Path:     "/auth",
		})

		response.JSON(w, http.StatusOK, LoginResponseDTO{
			User:   user,
			AccessToken: tokens.AccessToken,
		})
	} else {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			response.Unauthorized(r.Context(), w, err)
		case errors.Is(err, service.ErrIncorrectPassword):
			response.Unauthorized(r.Context(), w, err)
		default:
			response.InternalServerError(r.Context(), w, err)
		}
	}
}
