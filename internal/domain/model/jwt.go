package model

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"time"
)

type SigningKey struct {
	ID         string
	Issuer     string
	PublicKey  *rsa.PublicKey
	PrivateKey *rsa.PrivateKey
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshToken struct {
	UserID    string
	Token     string
	Revoked   bool
	ExpiresAt time.Time
}

func GenerateRefreshToken(userID string, ttl time.Duration) (RefreshToken, error) {
	b := make([]byte, 32)

	_, err := rand.Read(b)
	if err != nil {
		return RefreshToken{}, err
	}

	token := base64.URLEncoding.EncodeToString(b)

	return RefreshToken{
		UserID:    userID,
		Token:     token,
		Revoked:   false,
		ExpiresAt: time.Now().Add(ttl),
	}, nil
}
