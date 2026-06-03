package https

import (
	"log/slog"
	"net/http"

	"github.com/AdityaTaggar05/Purgatorio/internal/api/https/auth"
	"github.com/AdityaTaggar05/Purgatorio/internal/api/https/middleware"
	"github.com/go-chi/chi"
)

func NewRouter(logger *slog.Logger, authHandler *auth.AuthHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestLogger(logger))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Service is up and running!"))
	})

	return r
}
