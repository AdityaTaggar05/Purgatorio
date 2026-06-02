package config

import "time"

type ServerConfig struct {
	Port int
	ReadTimeout time.Duration
	WriteTimeout time.Duration
}
