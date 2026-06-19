package engine

import "math"

func findNearestBuilding(troops []*troopState, building *buildingState, maxRange float64) *troopState {
	var nearest *troopState
	var nearestDist float64

	for _, t := range troops {
		if !t.alive {
			continue
		}
		dist := distance(t.pos, buildingCenter(building.pos, building.size))
		if dist <= maxRange {
			if nearest == nil || dist < nearestDist {
				nearest = t
				nearestDist = dist
			}
		}
	}
	return nearest
}

func findNearestBuildingForTroop(t *troopState, buildings []*buildingState) *buildingState {
	var nearest *buildingState
	var nearestBastion *buildingState
	var nearestDist float64
	var nearestBastionDist float64

	for _, b := range buildings {
		if !b.alive {
			continue
		}
		dist := distance(t.pos, buildingCenter(b.pos, b.size))
		if b.buildingType == "bastion" {
			if nearestBastion == nil || dist < nearestBastionDist {
				nearestBastion = b
				nearestBastionDist = dist
			}
			continue
		}
		if nearest == nil || dist < nearestDist {
			nearest = b
			nearestDist = dist
		}
	}

	// Prefer non-bastion buildings. Only target a bastion if no other building exists.
	if nearest != nil {
		return nearest
	}
	return nearestBastion
}

func distance(a, b Point) float64 {
	dx := a.X - b.X
	dy := a.Y - b.Y
	return math.Sqrt(dx*dx + dy*dy)
}

func buildingCenter(pos Point, size int) Point {
	return Point{X: pos.X + float64(size)/2.0, Y: pos.Y + float64(size)/2.0}
}
