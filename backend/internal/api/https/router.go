package https

import (
	"crypto/rsa"
	"log/slog"
	"net/http"

	"github.com/AdityaTaggar05/Purgatorio/internal/api/https/army"
	"github.com/AdityaTaggar05/Purgatorio/internal/api/https/auth"
	"github.com/AdityaTaggar05/Purgatorio/internal/api/https/base"
	"github.com/AdityaTaggar05/Purgatorio/internal/api/https/battle"
	"github.com/AdityaTaggar05/Purgatorio/internal/api/https/middleware"
	"github.com/AdityaTaggar05/Purgatorio/internal/api/https/shop"
	"github.com/AdityaTaggar05/Purgatorio/internal/api/https/user"
	"github.com/AdityaTaggar05/Purgatorio/pkg/response"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

func NewRouter(logger *slog.Logger, publicKey *rsa.PublicKey, authHandler *auth.AuthHandler, userHandler *user.UserHandler, shopHandler *shop.ShopHandler, baseHandler *base.BaseHandler, armyHandler *army.ArmyHandler, battleHandler *battle.BattleHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestLogger(logger))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:5173",
		},

		AllowedMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
		},

		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
		},

		AllowCredentials: true,
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		response.Success(w, nil, "Service is up and running!")
	})

	r.Mount("/auth", authHandler.Routes())
	r.Get("/.well-known/jwks.json", authHandler.HandleJWKS)

	protected := chi.NewMux()
	protected.Use(middleware.RequestAuthenticator(publicKey))
	protected.Mount("/user", userHandler.Routes())
	protected.Mount("/shop", shopHandler.Routes())
	protected.Mount("/base", baseHandler.Routes())
	protected.Mount("/army", armyHandler.Routes())
	protected.Mount("/battle", battleHandler.Routes())

	r.Mount("/", protected)

	return r
}
