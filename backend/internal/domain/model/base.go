package model

import (
	"time"

	"github.com/google/uuid"
)

type BuildingCategory string

const (
	BuildingDefense BuildingCategory = "defense"
	BuildingArmy BuildingCategory = "army"
	BuildingResource BuildingCategory = "resource"
	BuildingOther BuildingCategory = "other"
)

func (b BuildingCategory) String() string {
	return string(b)
}

type Building struct {
	ID       string           `json:"id"`
	Name     string           `json:"name"`
	Size     int              `json:"size"`
	Price    int              `json:"price"`
	Currency Currency         `json:"currency"`
	Category BuildingCategory `json:"category"`
	Metadata BuildingMetadata `json:"metadata"`
}

type BuildingMetadata struct {
	UpgradeEndsAt *time.Time `json:"upgrade_ends_at"`
}

type BaseLayout struct {
	UserID    uuid.UUID  `json:"user_id"`
	Buildings []Building `json:"buildings"`
}
