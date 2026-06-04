package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/AdityaTaggar05/Purgatorio/internal/domain/model"
	"github.com/AdityaTaggar05/Purgatorio/internal/domain/service"
	"github.com/AdityaTaggar05/Purgatorio/pkg/response"
)

type LoginRequestDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponseDTO struct {
	User   model.User      `json:"user"`
	Tokens model.TokenPair `json:"tokens"`
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
		response.JSON(w, http.StatusOK, LoginResponseDTO{
			User:   user,
			Tokens: tokens,
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
