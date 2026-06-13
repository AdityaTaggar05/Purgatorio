package model

import (
	"time"

	"github.com/AdityaTaggar05/Purgatorio/internal/engine"
	"github.com/google/uuid"
)

type MatchListEntry struct {
	UserID       uuid.UUID `json:"user_id"`
	Username     string    `json:"username"`
	TerraceLevel int       `json:"terrace_level"`
}
