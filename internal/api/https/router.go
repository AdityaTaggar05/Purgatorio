package https

import (
	"net/http"

	"github.com/AdityaTaggar05/Purgatorio/internal/api/https/auth"
	"github.com/go-chi/chi"
)

func NewRouter(authHandler *auth.AuthHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/auth/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Service is up and running!"))
		w.WriteHeader(http.StatusOK)
	})

	return r
}
