package repository

import (
	"context"

	"github.com/AdityaTaggar05/Purgatorio/internal/domain/model"
	"github.com/google/uuid"
)

type ShopRepository interface {
	GetAllBuildings(ctx context.Context) ([]model.Building, error)
	GetBuildingByID(ctx context.Context, id string) (model.Building, error)
	GetUserBuildingCounts(ctx context.Context, userID uuid.UUID) (map[string]int, error)
	GetLimitsByTerrace(ctx context.Context, terraceLevel int) (map[string]int, error)
	GetBuildingLevels(ctx context.Context, buildingID string) ([]model.BuildingLevel, error)
	PurchaseBuilding(ctx context.Context, userID uuid.UUID, buildingID string, price int, currency model.Currency) error
}
