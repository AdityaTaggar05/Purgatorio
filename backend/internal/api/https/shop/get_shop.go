package shop

import (
	"net/http"

	"github.com/AdityaTaggar05/Purgatorio/internal/domain/model"
	"github.com/AdityaTaggar05/Purgatorio/pkg/ctxkeys"
	"github.com/AdityaTaggar05/Purgatorio/pkg/response"
	"github.com/google/uuid"
)

type ShopResponseDTO struct {
	Items []model.ShopItem `json:"items"`
}

func (h *ShopHandler) HandleGetShop(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(ctxkeys.UserID).(uuid.UUID)

	items, err := h.Service.GetShop(r.Context(), userID)
	if err != nil {
		response.Error(r.Context(), w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, ShopResponseDTO{Items: items})
}
