package repository

import (
	"context"

	"github.com/AdityaTaggar05/Purgatorio/internal/domain/model"
	"github.com/google/uuid"
)

type BaseRepository interface {
	GetResourceGenerationInfo(ctx context.Context, id uuid.UUID) ([]model.ResourceGenerationInfo, error)
}
