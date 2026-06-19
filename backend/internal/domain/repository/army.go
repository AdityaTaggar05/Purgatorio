package repository

import (
	"context"

	"github.com/AdityaTaggar05/Purgatorio/internal/domain/model"
	"github.com/google/uuid"
)

type ArmyRepository interface {
	GetAllTroops(ctx context.Context) ([]model.Troop, error)
	GetTroopByID(ctx context.Context, troopID string) (*model.Troop, error)
	GetUserArmy(ctx context.Context, userID uuid.UUID) (*model.UserArmy, error)
	AddTroops(ctx context.Context, userID uuid.UUID, troopID string, count, usedCapacity int) error
	RemoveTroops(ctx context.Context, userID uuid.UUID, troopID string, count, usedCapacity int) error
}
