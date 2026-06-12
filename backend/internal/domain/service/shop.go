package service

import (
	"context"
	"fmt"

	"github.com/AdityaTaggar05/Purgatorio/internal/domain/model"
	"github.com/AdityaTaggar05/Purgatorio/internal/domain/repository"
	"github.com/AdityaTaggar05/Purgatorio/pkg/purgerr"
	"github.com/google/uuid"
)

type ShopService struct {
	ShopRepo repository.ShopRepository
	UserRepo repository.UserRepository
}

func NewShopService(shopRepo repository.ShopRepository, userRepo repository.UserRepository) *ShopService {
	return &ShopService{ShopRepo: shopRepo, UserRepo: userRepo}
}

func (s *ShopService) GetShop(ctx context.Context, userID uuid.UUID) ([]model.ShopItem, error) {
	user, err := s.UserRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, purgerr.Wrap(ErrUserNotFound, err)
	}

	buildings, err := s.ShopRepo.GetAllBuildings(ctx)
	if err != nil {
		return nil, purgerr.Wrap(fmt.Errorf("failed to get shop buildings"), err)
	}

	counts, err := s.ShopRepo.GetUserBuildingCounts(ctx, userID)
	if err != nil {
		return nil, purgerr.Wrap(fmt.Errorf("failed to get building counts"), err)
	}

	limits, err := s.ShopRepo.GetLimitsByTerrace(ctx, user.TerraceLevel)
	if err != nil {
		return nil, purgerr.Wrap(fmt.Errorf("failed to get building limits"), err)
	}

	items := make([]model.ShopItem, 0, len(buildings))
	for _, b := range buildings {
		owned := counts[b.ID]
		maxAllowed := limits[b.ID]

		items = append(items, model.ShopItem{
			Building:     b,
			CurrentOwned: owned,
			MaxAllowed:   maxAllowed,
			CanBuy:       owned < maxAllowed,
		})
	}

	return items, nil
}

func (s *ShopService) BuyBuilding(ctx context.Context, userID uuid.UUID, buildingID string) error {
	building, err := s.ShopRepo.GetBuildingByID(ctx, buildingID)
	if err != nil {
		return purgerr.Wrap(ErrBuildingNotFound, err)
	}

	eco, err := s.UserRepo.GetEconomy(ctx, userID)
	if err != nil {
		return purgerr.Wrap(ErrUserNotFound, err)
	}

	user, err := s.UserRepo.GetUserByID(ctx, userID)
	if err != nil {
		return purgerr.Wrap(ErrUserNotFound, err)
	}

	counts, err := s.ShopRepo.GetUserBuildingCounts(ctx, userID)
	if err != nil {
		return purgerr.Wrap(fmt.Errorf("failed to get building counts"), err)
	}

	limits, err := s.ShopRepo.GetLimitsByTerrace(ctx, user.TerraceLevel)
	if err != nil {
		return purgerr.Wrap(fmt.Errorf("failed to get building limits"), err)
	}

	owned := counts[buildingID]
	maxAllowed := limits[buildingID]

	if owned >= maxAllowed {
		return purgerr.Wrap(ErrBuildingLimitReached, ErrBuildingLimitReached)
	}

	balance := eco.Penitence
	if building.Currency == model.CurrencyGrace {
		balance = eco.Grace
	}

	if balance < building.Price {
		return purgerr.Wrap(ErrInsufficientResources, ErrInsufficientResources)
	}

	if err := s.ShopRepo.PurchaseBuilding(ctx, userID, buildingID, building.Price, building.Currency); err != nil {
		return purgerr.Wrap(fmt.Errorf("purchase failed"), err)
	}

	return nil
}
