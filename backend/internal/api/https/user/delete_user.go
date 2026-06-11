package user

import (
	"errors"
	"net/http"

	"github.com/AdityaTaggar05/Purgatorio/internal/domain/service"
	"github.com/AdityaTaggar05/Purgatorio/pkg/ctxkeys"
	"github.com/AdityaTaggar05/Purgatorio/pkg/response"
	"github.com/google/uuid"
)

func (h *UserHandler) HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	if err := h.Service.DeleteUser(r.Context(), r.Context().Value(ctxkeys.UserID).(uuid.UUID)); err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			response.NotFound(r.Context(), w, err)
		default:
			response.InternalServerError(r.Context(), w, err)
		}
		return
	}

	response.Success(w, nil, "account deleted successfully")
}

