package config

import "time"

type Config struct {
	Server ServerConfig
	Postgres PostgresConfig
	JWT JWTConfig
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port: mustGetEnv("PORT"),
			ReadTimeout: getDuration("SERVER_READ_TIMEOUT", 5*time.Second),
			WriteTimeout: getDuration("SERVER_WRITE_TIMEOUT", 10*time.Second),
		},
		Postgres: PostgresConfig{
			URL: mustGetEnv("DB_URL"),
			MaxOpenConns: getInt("DB_MAX_OPEN_CONNS", 10),
		},
		JWT: JWTConfig{
			PrivateKeyPath: mustGetEnv("JWT_PRIVATE_KEY_PATH"),
			PublicKeyPath: mustGetEnv("JWT_PUBLIC_KEY_PATH"),
			Issuer: getEnv("JWT_ISSUER", "purgatorio"),
			AccessTTL: getDuration("JWT_ACCESS_TTL", 30*time.Minute),
			RefreshTTL: getDuration("JWT_REFRESH_TTL", 7*24*time.Hour),
		},
	}
}
