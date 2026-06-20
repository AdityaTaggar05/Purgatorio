package service

import (
	"context"
	"fmt"
	"time"

	"github.com/AdityaTaggar05/Purgatorio/internal/domain/model"
	"github.com/AdityaTaggar05/Purgatorio/internal/domain/repository"
	"github.com/AdityaTaggar05/Purgatorio/pkg/purgerr"
	"github.com/google/uuid"
)

type BaseService struct {
	BaseLayoutRepo repository.BaseLayoutRepository
	ShopRepo       repository.ShopRepository
	UserRepo       repository.UserRepository
}

func NewBaseService(baseLayoutRepo repository.BaseLayoutRepository, shopRepo repository.ShopRepository, userRepo repository.UserRepository) *BaseService {
	return &BaseService{
		BaseLayoutRepo: baseLayoutRepo,
		ShopRepo:       shopRepo,
		UserRepo:       userRepo,
	}
}

func (s *BaseService) GetLayout(ctx context.Context, userID uuid.UUID) (*model.BaseLayoutResponse, error) {
	placed, err := s.BaseLayoutRepo.GetLayout(ctx, userID)
	if err != nil {
		return nil, purgerr.Wrap(fmt.Errorf("failed to get layout"), err)
	}

	response := &model.BaseLayoutResponse{
		Buildings: make([]model.PlacedBuildingResponse, 0, len(placed)),
		GridW:     30,
		GridH:     30,
	}

	for _, pb := range placed {
		building, err := s.ShopRepo.GetBuildingByID(ctx, pb.BuildingID)
		if err != nil {
			continue
		}

		stats, err := s.BaseLayoutRepo.GetBuildingLevelStats(ctx, pb.BuildingID, pb.Level)
		if err != nil {
			continue
		}

		var upgradeCost, upgradeTime *int
		if _, err := s.BaseLayoutRepo.GetBuildingLevelStats(ctx, pb.BuildingID, pb.Level+1); err == nil {
			upgradeCost = &stats.UpgradeCost
			upgradeTime = &stats.UpgradeTime
		}

		response.Buildings = append(response.Buildings, model.PlacedBuildingResponse{
			BuildingID:      pb.BuildingID,
			Name:            building.Name,
			Category:        building.Category,
			Level:           pb.Level,
			X:               pb.X,
			Y:               pb.Y,
			Size:            building.Size,
			HP:              stats.HP,
			DPS:             stats.DamagePerSec,
			AttackRange:     stats.AttackRange,
			ProductionRate:  stats.ProductionRate,
			StorageCapacity: stats.StorageCapacity,
			UpgradeCost:     upgradeCost,
			UpgradeTime:     upgradeTime,
			Metadata:        pb.Metadata,
		})
	}

	return response, nil
}

func (s *BaseService) PlaceBuilding(ctx context.Context, userID uuid.UUID, buildingID string, x, y int) error {
	building, err := s.ShopRepo.GetBuildingByID(ctx, buildingID)
	if err != nil {
		return purgerr.Wrap(ErrBuildingNotFound, err)
	}

	if x < 0 || y < 0 || x+building.Size > 30 || y+building.Size > 30 {
		return purgerr.Wrap(ErrPositionOutOfBounds, ErrPositionOutOfBounds)
	}

	counts, err := s.ShopRepo.GetUserBuildingCounts(ctx, userID)
	if err != nil {
		return purgerr.Wrap(fmt.Errorf("failed to get building counts"), err)
	}

	placed, err := s.BaseLayoutRepo.GetLayout(ctx, userID)
	if err != nil {
		return purgerr.Wrap(fmt.Errorf("failed to get layout"), err)
	}

	placedCount := 0
	for _, pb := range placed {
		if pb.BuildingID == buildingID {
			placedCount++
		}
	}

	owned := counts[buildingID]
	if placedCount >= owned {
		return purgerr.Wrap(ErrNotEnoughBuildingsInInventory, ErrNotEnoughBuildingsInInventory)
	}

	sizes, err := s.buildingSizeMap(ctx)
	if err != nil {
		return err
	}

	if err := checkOverlap(placed, sizes, x, y, building.Size); err != nil {
		return err
	}

	pb := model.PlacedBuilding{
		UserID:     userID,
		BuildingID: buildingID,
		X:          x,
		Y:          y,
		Level:      1,
		Metadata:   &model.BuildingMetadata{},
	}

	if err := s.BaseLayoutRepo.PlaceBuilding(ctx, pb); err != nil {
		return purgerr.Wrap(fmt.Errorf("failed to place building"), err)
	}

	return nil
}

func (s *BaseService) RemoveBuilding(ctx context.Context, userID uuid.UUID, buildingID string, x, y int) error {
	pb, err := s.BaseLayoutRepo.GetBuildingAtPosition(ctx, userID, x, y)
	if err != nil {
		return purgerr.Wrap(ErrBuildingNotPlaced, err)
	}

	if pb.BuildingID != buildingID {
		return purgerr.Wrap(ErrBuildingNotPlaced, ErrBuildingNotPlaced)
	}

	if err := s.BaseLayoutRepo.RemoveBuilding(ctx, userID, buildingID, x, y); err != nil {
		return purgerr.Wrap(fmt.Errorf("failed to remove building"), err)
	}

	return nil
}

func (s *BaseService) MoveBuilding(ctx context.Context, userID uuid.UUID, buildingID string, fromX, fromY, toX, toY int) error {
	pb, err := s.BaseLayoutRepo.GetBuildingAtPosition(ctx, userID, fromX, fromY)
	if err != nil {
		return purgerr.Wrap(ErrBuildingNotPlaced, err)
	}

	if pb.BuildingID != buildingID {
		return purgerr.Wrap(ErrBuildingNotPlaced, ErrBuildingNotPlaced)
	}

	building, err := s.ShopRepo.GetBuildingByID(ctx, buildingID)
	if err != nil {
		return purgerr.Wrap(ErrBuildingNotFound, err)
	}

	if toX < 0 || toY < 0 || toX+building.Size > 30 || toY+building.Size > 30 {
		return purgerr.Wrap(ErrPositionOutOfBounds, ErrPositionOutOfBounds)
	}

	if fromX == toX && fromY == toY {
		return nil
	}

	placed, err := s.BaseLayoutRepo.GetLayout(ctx, userID)
	if err != nil {
		return purgerr.Wrap(fmt.Errorf("failed to get layout"), err)
	}

	sizes, err := s.buildingSizeMap(ctx)
	if err != nil {
		return err
	}

	filtered := make([]model.PlacedBuilding, 0, len(placed))
	for _, pb := range placed {
		if pb.BuildingID == buildingID && pb.X == fromX && pb.Y == fromY {
			continue
		}
		filtered = append(filtered, pb)
	}

	if err := checkOverlap(filtered, sizes, toX, toY, building.Size); err != nil {
		return err
	}

	if err := s.BaseLayoutRepo.MoveBuilding(ctx, userID, buildingID, fromX, fromY, toX, toY); err != nil {
		return purgerr.Wrap(fmt.Errorf("failed to move building"), err)
	}

	return nil
}

func (s *BaseService) buildingSizeMap(ctx context.Context) (map[string]int, error) {
	buildings, err := s.ShopRepo.GetAllBuildings(ctx)
	if err != nil {
		return nil, purgerr.Wrap(fmt.Errorf("failed to resolve building sizes"), err)
	}

	sizes := make(map[string]int, len(buildings))
	for _, b := range buildings {
		sizes[b.ID] = b.Size
	}
	return sizes, nil
}

func checkOverlap(placed []model.PlacedBuilding, sizes map[string]int, x, y, size int) error {
	for _, pb := range placed {
		pbSize, ok := sizes[pb.BuildingID]
		if !ok {
			continue
		}

		if overlap(pb.X, pb.Y, pbSize, x, y, size) {
			return purgerr.Wrap(ErrPositionOccupied, ErrPositionOccupied)
		}
	}

	return nil
}

func overlap(x1, y1, s1, x2, y2, s2 int) bool {
	return x1 < x2+s2 && x1+s1 > x2 && y1 < y2+s2 && y1+s1 > y2
}

func (s *BaseService) UpgradeBuilding(ctx context.Context, userID uuid.UUID, buildingID string, x, y int) error {
	pb, err := s.BaseLayoutRepo.GetBuildingAtPosition(ctx, userID, x, y)
	if err != nil {
		return purgerr.Wrap(ErrBuildingNotPlaced, err)
	}

	if pb.BuildingID != buildingID {
		return purgerr.Wrap(ErrBuildingNotPlaced, ErrBuildingNotPlaced)
	}

	if pb.Metadata != nil && pb.Metadata.UpgradeEndsAt != nil {
		return purgerr.Wrap(ErrUpgradeAlreadyActive, ErrUpgradeAlreadyActive)
	}

	currentStats, err := s.BaseLayoutRepo.GetBuildingLevelStats(ctx, buildingID, pb.Level)
	if err != nil {
		return purgerr.Wrap(ErrBuildingNotFound, err)
	}

	nextLevel := pb.Level + 1
	if _, err := s.BaseLayoutRepo.GetBuildingLevelStats(ctx, buildingID, nextLevel); err != nil {
		return purgerr.Wrap(ErrMaxLevelReached, err)
	}

	eco, err := s.UserRepo.GetEconomy(ctx, userID)
	if err != nil {
		return purgerr.Wrap(ErrUserNotFound, err)
	}

	if eco.Penitence < currentStats.UpgradeCost {
		return purgerr.Wrap(ErrInsufficientResources, ErrInsufficientResources)
	}

	eco.Penitence -= currentStats.UpgradeCost

	if err := s.UserRepo.UpdateEconomy(ctx, eco); err != nil {
		return purgerr.Wrap(fmt.Errorf("failed to deduct upgrade cost"), err)
	}

	upgradeEndsAt := time.Now().Add(time.Duration(currentStats.UpgradeTime) * time.Second)
	if err := s.BaseLayoutRepo.StartUpgrade(ctx, userID, buildingID, x, y, upgradeEndsAt); err != nil {
		return purgerr.Wrap(fmt.Errorf("failed to start upgrade"), err)
	}

	return nil
}

func (s *BaseService) CheckIn(ctx context.Context, userID uuid.UUID) (*model.CheckInResult, error) {
	ready, err := s.BaseLayoutRepo.GetReadyUpgrades(ctx, userID)
	if err != nil {
		return nil, purgerr.Wrap(fmt.Errorf("failed to get ready upgrades"), err)
	}

	result := &model.CheckInResult{
		CompletedUpgrades: make([]model.CheckInUpgrade, 0, len(ready)),
	}

	for _, pb := range ready {
		building, err := s.ShopRepo.GetBuildingByID(ctx, pb.BuildingID)
		if err != nil {
			continue
		}

		newLevel := pb.Level + 1

		if building.Category == model.BuildingResource {
			if err := s.BaseLayoutRepo.BumpLevel(ctx, userID, pb.BuildingID, pb.X, pb.Y, newLevel); err != nil {
				continue
			}
		} else {
			if err := s.BaseLayoutRepo.CompleteUpgrade(ctx, userID, pb.BuildingID, pb.X, pb.Y, newLevel); err != nil {
				continue
			}
		}

		result.CompletedUpgrades = append(result.CompletedUpgrades, model.CheckInUpgrade{
			BuildingID: pb.BuildingID,
			X:          pb.X,
			Y:          pb.Y,
			FromLevel:  pb.Level,
			ToLevel:    newLevel,
		})

		if pb.BuildingID == "sanctum" {
			user, err := s.UserRepo.GetUserByID(ctx, userID)
			if err != nil {
				continue
			}
			if err := s.UserRepo.UpdateTerraceLevel(ctx, userID, user.TerraceLevel+1); err != nil {
				continue
			}
		}
	}

	return result, nil
}
