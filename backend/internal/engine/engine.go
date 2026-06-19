package engine

import (
	"fmt"
	"math"
	"math/rand"
)

type hpRecord struct {
	newHP float64
	delta float64
}

type Simulation struct {
	rng         *rand.Rand
	troops      []*troopState
	buildings   []*buildingState
	tick        int
	maxTick     int
	ticksPerSec int
	done        bool
	changes     map[string]hpRecord
	initHP      map[string]int
	seed        int64
	idSeq       int
}

func NewSimulation(input BattleInput) *Simulation {
	rng := rand.New(rand.NewSource(input.Seed))

	ticksPerSec := input.TicksPerSec
	if ticksPerSec <= 0 {
		ticksPerSec = 20
	}

	maxTick := input.MaxDuration
	if maxTick <= 0 {
		maxTick = 180 * ticksPerSec
	}

	sim := &Simulation{
		rng:         rng,
		maxTick:     input.MaxDuration,
		ticksPerSec: input.TicksPerSec,
		changes:     make(map[string]hpRecord),
		initHP:      make(map[string]int),
		seed:        input.Seed,
	}

	for _, dep := range input.Deployments {
		stats, ok := input.Catalog[dep.TroopType]
		if !ok {
			continue
		}
		for i := 0; i < dep.Count; i++ {
			sim.idSeq++
			id := fmt.Sprintf("troop_%s_%d", dep.TroopType, sim.idSeq)
			sim.troops = append(sim.troops, &troopState{
				id:        id,
				troopType: dep.TroopType,
				hp:        float64(stats.HP),
				maxHP:     stats.HP,
				dps:       stats.DPS,
				speed:     stats.Speed,
				range_:    stats.Range,
				pos:       dep.Position,
				alive:     true,
			})
			sim.initHP[id] = stats.HP
		}
	}

	for _, b := range input.Buildings {
		sim.buildings = append(sim.buildings, &buildingState{
			id:           b.ID,
			buildingType: b.BuildingType,
			category:     b.Category,
			hp:           float64(b.HP),
			maxHP:        b.MaxHP,
			dps:          b.DPS,
			range_:       b.Range,
			pos:          b.Position,
			size:         b.Size,
			alive:        true,
		})
		sim.initHP[b.ID] = b.HP
	}

	return sim
}

func (s *Simulation) NextTick() TickResult {
	s.tick++
	s.changes = make(map[string]hpRecord)

	s.retargetTroops()

	step := 1.0 / float64(s.ticksPerSec)
	s.moveTroops(step)
	s.buildingsAttack(step)
	s.troopsAttack(step)
	s.checkDeaths()

	hpChanges := s.collectHPChanges()
	done := s.checkDone()

	positions := s.collectPositions()

	return TickResult{
		Tick:      s.tick,
		HPChanges: hpChanges,
		Positions: positions,
		Done:      done,
	}
}

func (s *Simulation) retargetTroops() {
	for _, t := range s.troops {
		if !t.alive {
			continue
		}
		if t.targetID == "" || !s.buildingAliveByID(t.targetID) {
			if b := findNearestBuildingForTroop(t, s.buildings); b != nil {
				t.targetID = b.id
				t.path = nil // Clear old path for new target
			} else {
				t.targetID = ""
				t.path = nil
			}
		}
	}
}

func (s *Simulation) buildingAliveByID(id string) bool {
	for _, b := range s.buildings {
		if b.id == id && b.alive {
			return true
		}
	}
	return false
}

func (s *Simulation) moveTroops(step float64) {
	for _, t := range s.troops {
		if !t.alive || t.targetID == "" {
			continue
		}
		target := s.buildingByID(t.targetID)
		if target == nil {
			continue
		}
		center := buildingCenter(target.pos, target.size)
		dist := distance(t.pos, center)
		speedStep := t.speed * step

		if dist <= t.range_ {
			t.path = nil
			continue
		}

		if len(t.path) == 0 {
			t.path = s.findPathForTroop(t, target)
		}

		if len(t.path) > 0 {
			wp := t.path[0]
			wpDist := distance(t.pos, wp)
			if wpDist <= speedStep {
				t.pos = wp
				t.path = t.path[1:]
			} else {
				moveToward(t, wp, speedStep)
			}
			continue
		}

		if dist-speedStep <= t.range_ {
			moveToward(t, center, dist-t.range_)
		} else {
			moveToward(t, center, speedStep)
		}
	}
}

func (s *Simulation) buildingByID(id string) *buildingState {
	for _, b := range s.buildings {
		if b.id == id {
			return b
		}
	}
	return nil
}

func (s *Simulation) buildingsAttack(step float64) {
	for _, b := range s.buildings {
		if !b.alive || b.dps == 0 || b.range_ <= 0 {
			continue
		}
		target := findNearestBuilding(s.troops, b, b.range_)
		if target == nil {
			continue
		}
		damage := float64(b.dps) * step
		target.hp -= damage
		s.changes[target.id] = hpRecord{newHP: target.hp, delta: -damage}
	}
}

func (s *Simulation) troopsAttack(step float64) {
	for _, t := range s.troops {
		if !t.alive || t.targetID == "" {
			continue
		}
		target := s.buildingByID(t.targetID)
		if target == nil || !target.alive {
			continue
		}
		dist := distance(t.pos, buildingCenter(target.pos, target.size))
		if dist > t.range_ {
			continue
		}
		damage := float64(t.dps) * step
		target.hp -= damage
		s.changes[target.id] = hpRecord{newHP: target.hp, delta: -damage}
	}
}

func (s *Simulation) checkDeaths() {
	for _, t := range s.troops {
		if t.alive && t.hp <= 0 {
			t.alive = false
			t.hp = 0
		}
	}
	for _, b := range s.buildings {
		if b.alive && b.hp <= 0 {
			b.alive = false
			b.hp = 0
		}
	}
}

func (s *Simulation) collectHPChanges() []HPChange {
	changes := make([]HPChange, 0, len(s.changes))
	for id, rec := range s.changes {
		entityType := "building"
		for _, t := range s.troops {
			if t.id == id {
				entityType = "troop"
				break
			}
		}
		changes = append(changes, HPChange{
			EntityID:   id,
			EntityType: entityType,
			NewHP:      int(math.Max(0, rec.newHP)),
			Delta:      int(math.Floor(rec.delta)),
		})
	}
	return changes
}

func (s *Simulation) collectPositions() []PositionChange {
	positions := make([]PositionChange, 0, len(s.troops))
	for _, t := range s.troops {
		if t.alive {
			positions = append(positions, PositionChange{
				EntityID: t.id,
				X:        t.pos.X,
				Y:        t.pos.Y,
			})
		}
	}
	return positions
}

func (s *Simulation) checkDone() bool {
	if s.done {
		return true
	}

	allTroopsDead := true
	for _, t := range s.troops {
		if t.alive {
			allTroopsDead = false
			break
		}
	}
	if allTroopsDead {
		s.done = true
		return true
	}

	allBuildingsDead := true
	for _, b := range s.buildings {
		if b.alive {
			allBuildingsDead = false
			break
		}
	}
	if allBuildingsDead {
		s.done = true
		return true
	}

	if s.tick >= s.maxTick {
		s.done = true
		return true
	}

	return false
}

func (s *Simulation) IsDone() bool {
	return s.done
}

func (s *Simulation) Seed() int64 {
	return s.seed
}

func (s *Simulation) Result() BattleResult {
	allTroopsDead := true
	for _, t := range s.troops {
		if t.alive {
			allTroopsDead = false
			break
		}
	}

	destruction := s.destructionPercent()

	if allTroopsDead {
		return BattleResult{
			Outcome:     Defeat,
			Destruction: destruction,
			Loot:        0,
			Duration:    s.tick,
		}
	}

	return BattleResult{
		Outcome:     Victory,
		Destruction: destruction,
		Loot:        0,
		Duration:    s.tick,
	}
}

func (s *Simulation) destructionPercent() float64 {
	totalMaxHP := 0.0
	totalDamage := 0.0
	for _, b := range s.buildings {
		if b.buildingType == "bastion" {
			continue
		}
		totalMaxHP += float64(b.maxHP)
		totalDamage += float64(b.maxHP) - b.hp
	}
	if totalMaxHP == 0 {
		return 0
	}
	pct := (totalDamage / totalMaxHP) * 100.0
	return math.Round(pct*10) / 10
}
