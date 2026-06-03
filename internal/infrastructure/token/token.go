package token

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"

	"github.com/AdityaTaggar05/Purgatorio/internal/config"
	"github.com/AdityaTaggar05/Purgatorio/internal/domain/model"
)

func LoadSigningKey(cfg config.JWTConfig) (*model.SigningKey, error) {
	keyData, err := os.ReadFile(cfg.PrivateKeyPath)

	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode([]byte(keyData))
	if block == nil {
		return nil, errors.New("invalid private key")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	privateKey := key.(*rsa.PrivateKey)

	return &model.SigningKey{
		ID:         "purg-key-2025-12",
		Issuer:     cfg.Issuer,
		PrivateKey: privateKey,
		PublicKey:  &privateKey.PublicKey,
	}, nil
}
