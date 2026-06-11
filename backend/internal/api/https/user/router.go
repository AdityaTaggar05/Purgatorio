package user

import "github.com/go-chi/chi"

func (h *UserHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/me", h.HandleMe)
	r.Delete("/me", h.HandleDeleteUser)

	r.Get("/:id", h.HandleGetUser)

	return r
}
