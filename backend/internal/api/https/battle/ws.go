package battle

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/AdityaTaggar05/Purgatorio/internal/domain/service"
	"github.com/AdityaTaggar05/Purgatorio/internal/engine"
	"github.com/AdityaTaggar05/Purgatorio/pkg/ctxkeys"
	"github.com/AdityaTaggar05/Purgatorio/pkg/response"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

type wsDeployMessage struct {
	Type   string                   `json:"type"`
	Troops []engine.TroopDeployment `json:"troops"`
}

type wsTickBatch struct {
	Type       string              `json:"type"`
	Ticks      []engine.TickResult `json:"ticks"`
	BatchStart int                 `json:"batch_start"`
}

type wsBattleEnd struct {
	Type        string               `json:"type"`
	Outcome     engine.BattleOutcome `json:"outcome"`
	Destruction float64              `json:"destruction"`
	Loot        int                  `json:"loot"`
	SinMeter    int                  `json:"sin_meter"`
	Duration    int                  `json:"duration_ticks"`
}

type wsError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func (h *BattleHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	battleID, err := uuid.Parse(chi.URLParam(r, "battle_id"))
	if err != nil {
		response.BadRequest(r.Context(), w, fmt.Errorf("invalid battle id"))
		return
	}

	userID := r.Context().Value(ctxkeys.UserID).(uuid.UUID)

	conn, err := h.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(30 * time.Second))

	var deployMsg wsDeployMessage
	if err := conn.ReadJSON(&deployMsg); err != nil {
		conn.WriteJSON(wsError{Type: "error", Message: "failed to read deployment"})
		return
	}

	if deployMsg.Type != "deploy" {
		conn.WriteJSON(wsError{Type: "error", Message: "expected deploy message"})
		return
	}

	sim, err := h.Service.StartSimulation(r.Context(), battleID, userID, deployMsg.Troops)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrBattleNotFound):
			conn.WriteJSON(wsError{Type: "error", Message: "battle not found"})
		case errors.Is(err, service.ErrBattleNotPending):
			conn.WriteJSON(wsError{Type: "error", Message: "battle is not pending"})
		case errors.Is(err, service.ErrInsufficientArmyTroops):
			conn.WriteJSON(wsError{Type: "error", Message: "insufficient troops"})
		default:
			conn.WriteJSON(wsError{Type: "error", Message: "failed to start battle"})
		}
		return
	}

	conn.SetReadDeadline(time.Time{})

	batchStart := 0
	for !sim.IsDone() {
		batchSize := 10
		ticks := make([]engine.TickResult, 0, batchSize)
		for i := 0; i < batchSize && !sim.IsDone(); i++ {
			ticks = append(ticks, sim.NextTick())
		}

		conn.WriteJSON(wsTickBatch{
			Type:       "tick_batch",
			Ticks:      ticks,
			BatchStart: batchStart,
		})
		batchStart += len(ticks)
	}

	result, err := h.Service.ResolveAndStore(r.Context(), battleID, sim, deployMsg.Troops)
	if err != nil {
		conn.WriteJSON(wsError{Type: "error", Message: "failed to store battle result"})
		return
	}

	conn.WriteJSON(wsBattleEnd{
		Type:        "battle_end",
		Outcome:     result.Outcome,
		Destruction: result.Destruction,
		Loot:        result.Loot,
		SinMeter:    result.SinMeter,
		Duration:    result.Duration,
	})
}
