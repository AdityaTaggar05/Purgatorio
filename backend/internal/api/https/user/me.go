package user

import (
	"errors"
	"net/http"

	"github.com/AdityaTaggar05/Purgatorio/internal/domain/service"
	"github.com/AdityaTaggar05/Purgatorio/pkg/ctxkeys"
	"github.com/AdityaTaggar05/Purgatorio/pkg/response"
)

func (h *UserHandler) HandleMe(w http.ResponseWriter, r *http.Request) {
	if me, err := h.Service.GetUserByID(r.Context(), r.Context().Value(ctxkeys.UserID).(string)); err == nil {
		response.JSON(w, http.StatusOK, me)
	} else {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			response.Error(r.Context(), w, http.StatusNotFound, err)
		default:
			response.InternalServerError(r.Context(), w, err)
		}
	}
}
