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

const deploymentTimeSecs = 60

type wsClientMsg struct {
	Type   string                   `json:"type"`
	Troops []engine.TroopDeployment `json:"troops,omitempty"`
}

type wsServerMsg struct {
	Type        string               `json:"type"`
	TimeLeft    int                  `json:"time_left,omitempty"`
	Ticks       []engine.TickResult  `json:"ticks,omitempty"`
	BatchStart  int                  `json:"batch_start,omitempty"`
	Outcome     engine.BattleOutcome `json:"outcome,omitempty"`
	Destruction float64              `json:"destruction,omitempty"`
	Loot        int                  `json:"loot,omitempty"`
	SinMeter    int                  `json:"sin_meter,omitempty"`
	Duration    int                  `json:"duration_ticks,omitempty"`
	Message     string               `json:"message,omitempty"`
}

func (h *BattleHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	battleID, err := uuid.Parse(chi.URLParam(r, "battle_id"))
	if err != nil {
		response.BadRequest(r.Context(), w, fmt.Errorf("invalid battle id"))
		return
	}

	userID, ok := r.Context().Value(ctxkeys.UserID).(uuid.UUID)
	if !ok {
		response.Error(r.Context(), w, http.StatusUnauthorized, fmt.Errorf("missing user id"))
		return
	}

	conn, err := h.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.Logger.Error("ws upgrade failed", "error", err, "battle_id", battleID)
		response.Error(r.Context(), w, http.StatusBadRequest, fmt.Errorf("websocket upgrade failed: %w", err))
		return
	}
	defer conn.Close()

	var deployments []engine.TroopDeployment

	conn.WriteJSON(wsServerMsg{Type: "deployment_start", TimeLeft: deploymentTimeSecs})

	conn.SetReadDeadline(time.Now().Add(time.Duration(deploymentTimeSecs+10) * time.Second))

	for {
		var msg wsClientMsg
		if err := conn.ReadJSON(&msg); err != nil {
			break
		}

		if msg.Type == "done" {
			break
		}

		if msg.Type == "deploy" {
			deployments = append(deployments, msg.Troops...)
			conn.WriteJSON(wsServerMsg{Type: "deploy_ack", Message: "ok"})
		}
	}

	if len(deployments) == 0 {
		conn.WriteJSON(wsServerMsg{Type: "error", Message: "no troops deployed"})
		return
	}

	if err := h.Service.ValidateFullDeployment(r.Context(), userID, deployments); err != nil {
		switch {
		case errors.Is(err, service.ErrInsufficientArmyTroops):
			conn.WriteJSON(wsServerMsg{Type: "error", Message: "insufficient troops"})
		default:
			conn.WriteJSON(wsServerMsg{Type: "error", Message: "validation failed"})
		}
		return
	}

	sim, err := h.Service.PrepareSimulation(r.Context(), battleID, userID, deployments)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrBattleNotFound):
			conn.WriteJSON(wsServerMsg{Type: "error", Message: "battle not found"})
		case errors.Is(err, service.ErrBattleNotPending):
			conn.WriteJSON(wsServerMsg{Type: "error", Message: "battle is not pending"})
		default:
			conn.WriteJSON(wsServerMsg{Type: "error", Message: "failed to start battle"})
		}
		return
	}

	var allTicks []engine.TickResult
	for !sim.IsDone() {
		allTicks = append(allTicks, sim.NextTick())
	}

	tickBatches := make([][]engine.TickResult, 0)
	batchSize := 10
	for i := 0; i < len(allTicks); i += batchSize {
		end := min(i + batchSize, len(allTicks))
		tickBatches = append(tickBatches, allTicks[i:end])
	}

	conn.SetReadDeadline(time.Time{})

	for batchIdx, batch := range tickBatches {
		batchStart := batchIdx * batchSize
		err = conn.WriteJSON(wsServerMsg{
			Type:       "tick_batch",
			Ticks:      batch,
			BatchStart: batchStart,
		})
		if err != nil {
			h.Logger.Error("failed to write tick batch", "error", err)
			return
		}
		
		time.Sleep(500 * time.Millisecond)
	}

	outcome, err := h.Service.ResolveAndStore(r.Context(), battleID, sim, deployments)
	if err != nil {
		conn.WriteJSON(wsServerMsg{Type: "error", Message: "failed to store battle result"})
		return
	}

	conn.WriteJSON(wsServerMsg{
		Type:        "battle_end",
		Outcome:     outcome.Outcome,
		Destruction: outcome.Destruction,
		Loot:        outcome.Loot,
		SinMeter:    outcome.SinMeter,
		Duration:    outcome.Duration,
	})
}
