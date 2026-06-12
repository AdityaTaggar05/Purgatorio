package shop

import "github.com/go-chi/chi"

func (h *ShopHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.HandleGetShop)
	r.Post("/buy", h.HandleBuyBuilding)

	return r
}
