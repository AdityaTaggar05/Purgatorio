package config

import "time"

type JWTConfig struct {
	PublicKeyPath string
	PrivateKeyPath string
	Issuer string
	AccessTTL time.Duration
	RefreshTTL time.Duration
}
