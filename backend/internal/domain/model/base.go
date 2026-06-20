package model

import (
	"time"

	"github.com/google/uuid"
)

type BuildingCategory string

const (
	BuildingDefense  BuildingCategory = "defense"
	BuildingArmy     BuildingCategory = "army"
	BuildingResource BuildingCategory = "resource"
	BuildingOther    BuildingCategory = "other"
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
}

type BuildingMetadata struct {
	UpgradeEndsAt *time.Time `json:"upgrade_ends_at,omitempty"`
}

type PlacedBuilding struct {
	UserID     uuid.UUID
	BuildingID string
	X          int
	Y          int
	Level      int
	Metadata   *BuildingMetadata
}

type PlacedBuildingResponse struct {
	BuildingID      string           `json:"building_id"`
	Name            string           `json:"name"`
	Category        BuildingCategory `json:"category"`
	Level           int              `json:"level"`
	X               int              `json:"x"`
	Y               int              `json:"y"`
	Size            int              `json:"size"`
	HP              *int             `json:"hp,omitempty"`
	DPS             *int             `json:"dps,omitempty"`
	AttackRange     *float64         `json:"attack_range,omitempty"`
	ProductionRate  *int             `json:"production_rate,omitempty"`
	StorageCapacity *int             `json:"storage_capacity,omitempty"`
	UpgradeCost     *int             `json:"upgrade_cost,omitempty"`
	UpgradeTime     *int             `json:"upgrade_time,omitempty"`
	Metadata        *BuildingMetadata `json:"metadata,omitempty"`
}

type BaseLayoutResponse struct {
	Buildings []PlacedBuildingResponse `json:"buildings"`
	GridW     int                      `json:"grid_w"`
	GridH     int                      `json:"grid_h"`
}

type CheckInResult struct {
	CompletedUpgrades []CheckInUpgrade `json:"completed_upgrades"`
}

type CheckInUpgrade struct {
	BuildingID string `json:"building_id"`
	X          int    `json:"x"`
	Y          int    `json:"y"`
	FromLevel  int    `json:"from_level"`
	ToLevel    int    `json:"to_level"`
}
