package auth

import (
	"encoding/json"
	"net/http"

	"github.com/AdityaTaggar05/Purgatorio/pkg/response"
)

type RegisterRequestDTO struct {
	Email string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *AuthHandler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequestDTO

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(r.Context(), w, err)
		return
	}
}
