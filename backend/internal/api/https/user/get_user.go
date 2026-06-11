package user

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/AdityaTaggar05/Purgatorio/internal/domain/service"
	"github.com/AdityaTaggar05/Purgatorio/pkg/purgerr"
	"github.com/AdityaTaggar05/Purgatorio/pkg/response"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (h *UserHandler) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := uuid.Validate(id); err != nil {
		response.Error(r.Context(), w, http.StatusBadRequest, purgerr.Wrap(fmt.Errorf("invalid user id"), err))
		return
	}

	if user, err := h.Service.GetUserByID(r.Context(), uuid.MustParse(id)); err == nil {
		response.JSON(w, http.StatusOK, user)
	} else {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			response.NotFound(r.Context(), w, err)
		default:
			response.InternalServerError(r.Context(), w, err)
		}
	}
}

