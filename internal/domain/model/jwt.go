package model

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"math/big"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

func (s *SigningKey) PublicKeyToJWK() map[string]string {
	n := base64.RawURLEncoding.EncodeToString(s.PublicKey.N.Bytes())
    e := base64.RawURLEncoding.EncodeToString(big.NewInt(int64(s.PublicKey.E)).Bytes())

    return map[string]string{
        "kty": "RSA",
        "kid": s.ID,
        "use": "sig",
        "alg": "RS256",
        "n":   n,
        "e":   e,
    }
}

func GenerateJWT(user User, signingKey *SigningKey, ttl time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": jwt.NewNumericDate(time.Now().Add(ttl)),
		"iat": jwt.NewNumericDate(time.Now()),
		"iss": signingKey.Issuer,
	})

	token.Header["kid"] = signingKey.ID

	return token.SignedString(signingKey.PrivateKey)
}
