package config

import "time"

type Config struct {
	Server ServerConfig
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port: mustGetEnv("PORT"),
			ReadTimeout: getDuration("SERVER_READ_TIMEOUT", 5*time.Second),
			WriteTimeout: getDuration("SERVER_WRITE_TIMEOUT", 10*time.Second),
		},
	}
}
