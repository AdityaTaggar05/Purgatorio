package model

import (
	"time"

	"github.com/google/uuid"
)

type AuthAndUser struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	XP           int       `json:"xp"`
	Level        int       `json:"level"`
	TerraceLevel int       `json:"terrace_level"`
}

type User struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	XP           int       `json:"xp"`
	Level        int       `json:"level"`
	TerraceLevel int       `json:"terrace_level"`
}

type UserEconomy struct {
	ID                        uuid.UUID
	Penitence                 int
	Grace                     int
	MaxPenitence              int
	CollectorPendingPenitence int
	CollectorResetAt          time.Time
}
