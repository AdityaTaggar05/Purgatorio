package engine

import "math"

func moveToward(t *troopState, target Point, step float64) {
	dx := target.X - t.pos.X
	dy := target.Y - t.pos.Y
	dist := math.Sqrt(dx*dx + dy*dy)

	if dist <= step {
		t.pos = target
		return
	}

	t.pos.X += (dx / dist) * step
	t.pos.Y += (dy / dist) * step
}
