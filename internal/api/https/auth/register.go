package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AdityaTaggar05/Purgatorio/internal/domain/model"
	"github.com/AdityaTaggar05/Purgatorio/pkg/response"
)

type RegisterRequestDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,alphanum,min=3,max=14"`
	Password string `json:"password" validate:"required"`
}

type RegsiterResponseDTO struct {
	User   model.User      `json:"user"`
	Tokens model.TokenPair `json:"tokens"`
}

func (h *AuthHandler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequestDTO

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(r.Context(), w, fmt.Errorf("invalid request JSON"))
		return
	}

	if err := h.Validator.Struct(req); err != nil {
		response.Error(r.Context(), w, http.StatusUnprocessableEntity, err)
		return
	}

	user, tokens, err := h.Service.Register(r.Context(), req.Email, req.Username, req.Password)
	if err != nil {

	}

	response.JSON(w, http.StatusCreated, RegsiterResponseDTO{
		User: user,
		Tokens: tokens,
	})
}
