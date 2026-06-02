package https

import (
	"net/http"

	"github.com/AdityaTaggar05/Purgatorio/internal/api/https/auth"
	"github.com/go-chi/chi"
)

func NewRouter(authHandler *auth.AuthHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Service is up and running!"))
	})

	return r
}
