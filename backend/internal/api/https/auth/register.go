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

type RegisterRequestDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,alphanum,min=3,max=14"`
	Password string `json:"password" validate:"required"`
}

type RegisterResponseDTO struct {
	User        model.User `json:"user"`
	AccessToken string     `json:"access_token"`
}

func (h *AuthHandler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequestDTO

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(r.Context(), w, fmt.Errorf("invalid request JSON"))
		return
	}

	if err := h.Validator.Struct(req); err != nil {
		response.ValidationFailed(r.Context(), w, err)
		return
	}

	if user, tokens, err := h.Service.Register(r.Context(), req.Email, req.Username, req.Password); err == nil {
		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    tokens.RefreshToken,
			Expires:  time.Now().Add(h.Service.Config.RefreshTTL),
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
			Path:     "/auth",
		})

		response.JSON(w, http.StatusCreated, RegisterResponseDTO{
			User:        user,
			AccessToken: tokens.AccessToken,
		})
	} else {
		switch {
		case errors.Is(err, service.ErrUserAlreadyExists):
			response.Error(r.Context(), w, http.StatusConflict, err)
		default:
			response.InternalServerError(r.Context(), w, err)
		}
	}
}
