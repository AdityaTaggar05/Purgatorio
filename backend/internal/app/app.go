package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AdityaTaggar05/Purgatorio/internal/api/https"
	"github.com/AdityaTaggar05/Purgatorio/internal/api/https/army"
	"github.com/AdityaTaggar05/Purgatorio/internal/api/https/auth"
	"github.com/AdityaTaggar05/Purgatorio/internal/api/https/base"
	"github.com/AdityaTaggar05/Purgatorio/internal/api/https/shop"
	"github.com/AdityaTaggar05/Purgatorio/internal/api/https/user"
	"github.com/AdityaTaggar05/Purgatorio/internal/config"
	"github.com/AdityaTaggar05/Purgatorio/internal/domain/repository"
	"github.com/AdityaTaggar05/Purgatorio/internal/domain/service"
	"github.com/AdityaTaggar05/Purgatorio/internal/infrastructure/postgres"
	"github.com/AdityaTaggar05/Purgatorio/internal/infrastructure/token"
)

type App struct {
	Config *config.Config
	Server *http.Server
	Logger *slog.Logger
	close func() error
}

func New(cfg *config.Config) (*App, error) {
	// 1) Logger Setup
	logger, close, err := initializeLogger(os.Getenv("LOG_FILE"))

	if err != nil {
		return nil, err
	}

	// 2) Infrastructure Setup
	ctx := context.Background()
	db, err := postgres.NewPostgresDB(ctx, cfg.Postgres)

	if err != nil {
		return nil, err
	}
	logger.Debug("connected to database!")

	// 3) Repository Setup
	var userRepo repository.UserRepository = postgres.NewUserRepository(db)
	var baseRepo repository.BaseRepository = postgres.NewBaseRepository(db)
	var shopRepo repository.ShopRepository = postgres.NewShopRepository(db)
	var baseLayoutRepo repository.BaseLayoutRepository = postgres.NewBaseLayoutRepository(db)
	var armyRepo repository.ArmyRepository = postgres.NewArmyRepository(db)

	// 4) Service Setup
	signingKey, err := token.LoadSigningKey(cfg.JWT)
	if err != nil {
		return nil, err
	}

	authService := service.NewAuthService(cfg.JWT, signingKey, userRepo)
	userService := service.NewUserService(userRepo, baseRepo)
	shopService := service.NewShopService(shopRepo, userRepo)
	baseService := service.NewBaseService(baseLayoutRepo, shopRepo, userRepo)
	armyService := service.NewArmyService(armyRepo, userRepo, baseLayoutRepo)

	// 5) Handler Setup
	authHandler := auth.NewHandler(logger, authService)
	userHandler := user.NewHandler(logger, userService)
	shopHandler := shop.NewHandler(logger, shopService)
	baseHandler := base.NewHandler(logger, baseService)
	armyHandler := army.NewHandler(logger, armyService)

	// 6) Router Setup
	router := https.NewRouter(logger, signingKey.PublicKey, authHandler, userHandler, shopHandler, baseHandler, armyHandler)

	// 7) Server Setup
	return &App{
		Config: cfg,
		Server: &http.Server{
			Addr: ":" + cfg.Server.Port,
			ReadTimeout: cfg.Server.ReadTimeout,
			WriteTimeout: cfg.Server.WriteTimeout,
			Handler: router,
		},
		Logger: logger,
		close: close,
	}, nil
}

func (a *App) Start() error {
	defer func() {
		if err := a.close(); err != nil {
			fmt.Fprint(os.Stderr, err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		a.Logger.Debug("Purgatorio Server listening on " + a.Server.Addr)

		if err := a.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.Logger.Error("HTTP server error: %v", "error", err)
		}
	}()

	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := a.Server.Shutdown(ctx); err != nil {
		a.Logger.Error("Server shutdown failed: %v\n", "error", err)
		return err
	}

	a.Logger.Debug("Purgatorio Server stopped gracefully!")
	return nil
}
