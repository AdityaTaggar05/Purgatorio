package engine

import "math"

const pfPad = 6

type cell struct {
	x, y int
}

func (s *Simulation) aStar(start cell, goals []cell) []cell {
	// Determine search bounds: bounding box of all buildings + padding
	minX, minY := math.MaxInt32, math.MaxInt32
	maxX, maxY := math.MinInt32, math.MinInt32
	for _, b := range s.buildings {
		bx := int(b.pos.X)
		by := int(b.pos.Y)
		be := bx + b.size
		bt := by + b.size
		if bx < minX {
			minX = bx
		}
		if by < minY {
			minY = by
		}
		if be > maxX {
			maxX = be
		}
		if bt > maxY {
			maxY = bt
		}
	}
	if minX == math.MaxInt32 {
		return nil
	}
	minX -= pfPad
	minY -= pfPad
	maxX += pfPad
	maxY += pfPad

	// Find the goal nearest the start for heuristic
	bestGoal := goals[0]
	bestDist := manhattan(start, bestGoal)
	for _, g := range goals[1:] {
		d := manhattan(start, g)
		if d < bestDist {
			bestDist = d
			bestGoal = g
		}
	}

	// A* with simple open set (slice-based for small grids)
	type node struct {
		cell
		g, f      float64
		px, py    int // parent cell
		hasParent bool
		closed    bool
	}

	open := []*node{{cell: start, g: 0, f: bestDist}}
	closed := make(map[cell]*node)
	closed[start] = open[0]

	var found *node

	for len(open) > 0 {
		bestIdx := 0
		for i, n := range open[1:] {
			if n.f < open[bestIdx].f {
				bestIdx = i + 1
			}
		}
		cur := open[bestIdx]
		open[bestIdx] = open[len(open)-1]
		open = open[:len(open)-1]

		if cur.closed {
			continue
		}
		cur.closed = true

		// Check if we've reached any goal
		for _, g := range goals {
			if cur.x == g.x && cur.y == g.y {
				found = cur
				goto reconstruct
			}
		}

		// 4-directional neighbors
		dirs := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
		for _, d := range dirs {
			nx, ny := cur.x+d[0], cur.y+d[1]
			if nx < minX || nx >= maxX || ny < minY || ny >= maxY {
				continue
			}
			nc := cell{nx, ny}
			if s.isBlockedCell(nc) {
				continue
			}
			if existing, ok := closed[nc]; ok && existing.closed {
				continue
			}
			ng := cur.g + 1.0
			nh := manhattan(nc, bestGoal)
			nf := ng + nh

			if existing, ok := closed[nc]; ok {
				if ng >= existing.g {
					continue
				}
				existing.g = ng
				existing.f = nf
				existing.px, existing.py = cur.x, cur.y
				existing.hasParent = true
				existing.closed = false
				open = append(open, existing)
			} else {
				nn := &node{cell: nc, g: ng, f: nf, px: cur.x, py: cur.y, hasParent: true}
				closed[nc] = nn
				open = append(open, nn)
			}
		}
	}

reconstruct:
	if found == nil {
		return nil
	}

	var path []cell
	for n := found; n.hasParent; n = closed[cell{n.px, n.py}] {
		path = append(path, n.cell)
	}
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path
}

func manhattan(a, b cell) float64 {
	return math.Abs(float64(a.x-b.x)) + math.Abs(float64(a.y-b.y))
}

func (s *Simulation) isBlockedCell(c cell) bool {
	for _, b := range s.buildings {
		if !b.alive {
			continue
		}
		bx := int(b.pos.X)
		by := int(b.pos.Y)
		if c.x >= bx && c.x < bx+b.size && c.y >= by && c.y < by+b.size {
			return true
		}
	}
	return false
}

// adjacentCells returns unblocked cells that are adjacent to (touching) a building's footprint.
func (s *Simulation) adjacentCells(b *buildingState) []cell {
	bx := int(b.pos.X)
	by := int(b.pos.Y)
	var cells []cell
	for dx := -1; dx <= b.size; dx++ {
		for dy := -1; dy <= b.size; dy++ {
			// Skip cells inside the building footprint
			if dx >= 0 && dx < b.size && dy >= 0 && dy < b.size {
				continue
			}
			c := cell{bx + dx, by + dy}
			if !s.isBlockedCell(c) {
				cells = append(cells, c)
			}
		}
	}
	return cells
}

// findPathForTroop computes an A* path for a troop to reach its target building.
// Returns waypoints as continuous Point positions (cell centers).
func (s *Simulation) findPathForTroop(t *troopState, target *buildingState) []Point {
	goals := s.adjacentCells(target)
	if len(goals) == 0 {
		return nil
	}

	start := cell{x: int(math.Floor(t.pos.X)), y: int(math.Floor(t.pos.Y))}

	// If start is already adjacent to the target and in range, no path needed
	for _, g := range goals {
		if start.x == g.x && start.y == g.y {
			return nil
		}
	}

	path := s.aStar(start, goals)
	if path == nil {
		return nil
	}

	// Convert cell path to float Points (cell centers)
	ptPath := make([]Point, len(path))
	for i, c := range path {
		ptPath[i] = Point{X: float64(c.x) + 0.5, Y: float64(c.y) + 0.5}
	}
	return ptPath
}
