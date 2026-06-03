package auth

import (
	"log/slog"

	"github.com/AdityaTaggar05/Purgatorio/internal/domain/service"
)

type AuthHandler struct {
	Logger *slog.Logger
	Service *service.AuthService
}

func NewHandler(logger *slog.Logger, service *service.AuthService) *AuthHandler {
	return &AuthHandler{
		Logger: logger,
		Service: service,
	}
}
