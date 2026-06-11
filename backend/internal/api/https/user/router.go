package user

import "github.com/go-chi/chi"

func (h *UserHandler) Routes() chi.Router {
	r := chi.NewRouter()

	return r
}
