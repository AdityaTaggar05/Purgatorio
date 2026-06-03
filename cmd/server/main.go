package main

import (
	"fmt"
	"os"

	"github.com/AdityaTaggar05/Purgatorio/internal/app"
	"github.com/AdityaTaggar05/Purgatorio/internal/config"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()
	server, err := app.New(cfg)

	if err != nil {
		server.Logger.Error(fmt.Sprintf("unable to initialize the server: %v", err))
		os.Exit(1)
	}

	if err := server.Start(); err != nil {
		server.Logger.Error(fmt.Sprintf("unable to start the server: %v", err))
		os.Exit(1)
	}
}
