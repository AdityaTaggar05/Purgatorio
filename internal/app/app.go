package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AdityaTaggar05/Purgatorio/internal/api/https"
	"github.com/AdityaTaggar05/Purgatorio/internal/api/https/auth"
	"github.com/AdityaTaggar05/Purgatorio/internal/config"
)

type App struct {
	Config *config.Config
	Server *http.Server
}

func New(cfg *config.Config) (*App, error) {
	authHandler := auth.NewHandler()
	router := https.NewRouter(authHandler)

	return &App{
		Config: cfg,
		Server: &http.Server{
			Addr: ":" + cfg.Server.Port,
			ReadTimeout: cfg.Server.ReadTimeout,
			WriteTimeout: cfg.Server.WriteTimeout,
			Handler: router,
		},
	}, nil
}

func (a *App) Start() error {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		fmt.Println("Purgatorio Server listening on " + a.Server.Addr)

		if err := a.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := a.Server.Shutdown(ctx); err != nil {
		fmt.Printf("Server shutdown failed: %v\n", err)
		return err
	}

	fmt.Println("Purgatorio Server stopped gracefully!")
	return nil
}
