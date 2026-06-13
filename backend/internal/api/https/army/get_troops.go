package army

import (
	"net/http"

	"github.com/AdityaTaggar05/Purgatorio/internal/domain/model"
	"github.com/AdityaTaggar05/Purgatorio/pkg/response"
)

type TroopsResponseDTO struct {
	Troops []model.Troop `json:"troops"`
}

func (h *ArmyHandler) HandleGetTroops(w http.ResponseWriter, r *http.Request) {
	troops, err := h.Service.GetTroops(r.Context())
	if err != nil {
		response.Error(r.Context(), w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, TroopsResponseDTO{Troops: troops})
}
