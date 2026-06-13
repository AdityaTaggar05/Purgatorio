package base

import "github.com/go-chi/chi"

func (h *BaseHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/layout", h.HandleGetLayout)
	r.Post("/layout", h.HandlePlaceBuilding)
	r.Delete("/layout", h.HandleRemoveBuilding)
	r.Put("/layout", h.HandleMoveBuilding)

	r.Post("/upgrade", h.HandleUpgradeBuilding)
	r.Post("/check-in", h.HandleCheckIn)

	return r
}
