package army

import "github.com/go-chi/chi"

func (h *ArmyHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/troops", h.HandleGetTroops)
	r.Get("/my-troops", h.HandleGetMyTroops)
	r.Post("/train", h.HandleTrainTroops)

	return r
}
