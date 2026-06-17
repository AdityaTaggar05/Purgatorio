package middleware

import (
	"context"
	"crypto/rsa"
	"fmt"
	"net/http"
	"strings"

	"github.com/AdityaTaggar05/Purgatorio/pkg/ctxkeys"
	"github.com/AdityaTaggar05/Purgatorio/pkg/purgerr"
	"github.com/AdityaTaggar05/Purgatorio/pkg/response"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func RequestAuthenticator(publicKey *rsa.PublicKey) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := extractToken(r)
			userID, err := parseAndVerify(token, publicKey)
			if err != nil {
				response.Error(r.Context(), w, http.StatusUnauthorized, err)
				return
			}
			ctx := context.WithValue(r.Context(), ctxkeys.UserID, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func parseAndVerify(tokenString string, publicKey *rsa.PublicKey) (uuid.UUID, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		return publicKey, nil
	})
	if err != nil || !token.Valid {
		return uuid.Nil, purgerr.Wrap(fmt.Errorf("invalid token"), err)
	}
	claims := token.Claims.(jwt.MapClaims)
	sub, _ := claims["sub"].(string)

	return uuid.Parse(sub)
}

func extractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	tokenParam := r.URL.Query().Get("token")

	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			return ""
		}
		return parts[1]
	} else if tokenParam != "" {
		return tokenParam
	} else {
		return ""
	}
}
