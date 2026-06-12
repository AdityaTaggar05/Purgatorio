package repository

import (
	"context"

	"github.com/AdityaTaggar05/Purgatorio/internal/domain/model"
	"github.com/google/uuid"
)

type BaseRepository interface {
	GetResourceGenerationInfo(ctx context.Context, userID uuid.UUID) ([]model.ResourceGenerationInfo, error)
	RemoveUpgradeInfo(ctx context.Context, userID uuid.UUID) error
}
