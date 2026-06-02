package app

import (
	"net/http"

	"github.com/AdityaTaggar05/Purgatorio/internal/config"
)

type App struct {
	Config *config.Config
	Server *http.Server
}

func New(cfg *config.Config) (*App, error) {
	return &App{
		Config: cfg,
		Server: &http.Server{
			Addr: ":" + cfg.Server.Port,
			ReadTimeout: cfg.Server.ReadTimeout,
			WriteTimeout: cfg.Server.WriteTimeout,
		},
	}, nil
}

func (a *App) Start() error {
	return nil
}
