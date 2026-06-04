package auth

import "github.com/go-chi/chi"

func (h *AuthHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/register", h.HandleRegister)
	r.Post("/login", h.HandleLogin)
	r.Post("/logout", h.HandleLogout)
	r.Post("/refresh", h.HandleRefresh)

	return r
}
