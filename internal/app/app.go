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
	"github.com/AdityaTaggar05/Purgatorio/internal/api/https/auth"
	"github.com/AdityaTaggar05/Purgatorio/internal/config"
	"github.com/AdityaTaggar05/Purgatorio/internal/infrastructure/postgres"
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
	_ = postgres.NewPostgresDB(logger, ctx, cfg.Postgres)

	// 3) Repository Setup

	// 4) Service Setup

	// 5) Handler Setup
	authHandler := auth.NewHandler()

	// 6) Router Setup
	router := https.NewRouter(logger, authHandler)

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
			a.Logger.Error("HTTP server error: %v", err)
		}
	}()

	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := a.Server.Shutdown(ctx); err != nil {
		a.Logger.Error("Server shutdown failed: %v\n", err)
		return err
	}

	a.Logger.Debug("Purgatorio Server stopped gracefully!")
	return nil
}
