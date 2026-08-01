package auth

import (
	"log/slog"
	"regexp"

	"github.com/AdityaTaggar05/Purgatorio/internal/domain/service"
	"github.com/go-playground/validator/v10"
)

var (
	hasUpper   = regexp.MustCompile(`[A-Z]`)
	hasLower   = regexp.MustCompile(`[a-z]`)
	hasNumber  = regexp.MustCompile(`[0-9]`)
	hasSpecial = regexp.MustCompile(`[!@#~$%^&*(),.?":{}|<>]`)
)

func validatePasswordComplexity(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < 8 {
		return false
	}

	return hasUpper.MatchString(password) && hasLower.MatchString(password) && hasNumber.MatchString(password) && hasSpecial.MatchString(password)
}

type AuthHandler struct {
	Logger    *slog.Logger
	Service   *service.AuthService
	Validator *validator.Validate
}

func NewHandler(logger *slog.Logger, service *service.AuthService) *AuthHandler {
	validate := validator.New(validator.WithRequiredStructEnabled())

	validate.RegisterValidation("password", validatePasswordComplexity)

	return &AuthHandler{
		Logger:    logger,
		Service:   service,
		Validator: validate,
	}
}
