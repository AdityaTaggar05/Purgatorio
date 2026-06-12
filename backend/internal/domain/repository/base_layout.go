package repository

import (
	"context"
	"time"

	"github.com/AdityaTaggar05/Purgatorio/internal/domain/model"
	"github.com/google/uuid"
)

type BaseLayoutRepository interface {
	GetLayout(ctx context.Context, userID uuid.UUID) ([]model.PlacedBuilding, error)
	PlaceBuilding(ctx context.Context, pb model.PlacedBuilding) error
	RemoveBuilding(ctx context.Context, userID uuid.UUID, buildingID string, x, y int) error
	MoveBuilding(ctx context.Context, userID uuid.UUID, buildingID string, fromX, fromY, toX, toY int) error
	GetBuildingAtPosition(ctx context.Context, userID uuid.UUID, x, y int) (*model.PlacedBuilding, error)
	GetBuildingLevelStats(ctx context.Context, buildingID string, level int) (*model.BuildingLevel, error)
	StartUpgrade(ctx context.Context, userID uuid.UUID, buildingID string, x, y int, upgradeEndsAt time.Time) error
	GetReadyUpgrades(ctx context.Context, userID uuid.UUID) ([]model.PlacedBuilding, error)
	CompleteUpgrade(ctx context.Context, userID uuid.UUID, buildingID string, x, y int, newLevel int) error
	BumpLevel(ctx context.Context, userID uuid.UUID, buildingID string, x, y int, newLevel int) error
}
