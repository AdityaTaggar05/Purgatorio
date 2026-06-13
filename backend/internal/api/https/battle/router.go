package battle

import "github.com/go-chi/chi"

func (h *BattleHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/matchlist", h.HandleMatchList)
	r.Post("/initiate", h.HandleInitiate)
	r.Get("/{battle_id}/ws", h.HandleWebSocket)
	r.Get("/{battle_id}/result", h.HandleResult)
	r.Get("/{battle_id}/replay", h.HandleReplay)
	r.Get("/attacks", h.HandleAttacks)
	r.Get("/defenses", h.HandleDefenses)

	return r
}
