package engine

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type TroopDeployment struct {
	TroopType string `json:"troop_type"`
	Position  Point  `json:"position"`
	Count     int    `json:"count"`
}

type BuildingSnapshot struct {
	ID           string  `json:"id"`
	BuildingType string  `json:"building_type"`
	Category     string  `json:"category"`
	Level        int     `json:"level"`
	HP           int     `json:"hp"`
	MaxHP        int     `json:"max_hp"`
	DPS          int     `json:"dps"`
	Range        float64 `json:"range"`
	Position     Point   `json:"position"`
	Size         int     `json:"size"`
}

type BattleInput struct {
	Seed        int64
	Deployments []TroopDeployment
	Buildings   []BuildingSnapshot
	Catalog     TroopCatalog
	TicksPerSec int
	MaxDuration int
}

type PositionChange struct {
	EntityID string `json:"entity_id"`
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
}

type TickResult struct {
	Tick      int              `json:"tick"`
	HPChanges []HPChange       `json:"hp_changes"`
	Positions []PositionChange `json:"positions,omitempty"`
	Done      bool             `json:"done,omitempty"`
}

type HPChange struct {
	EntityID   string `json:"entity_id"`
	EntityType string `json:"entity_type"`
	NewHP      int    `json:"new_hp"`
	Delta      int    `json:"delta"`
}

type BattleOutcome string

const (
	Victory         BattleOutcome = "victory"
	Defeat          BattleOutcome = "defeat"
	ThresholdFailed BattleOutcome = "threshold_failed"
)

type BattleResult struct {
	Outcome     BattleOutcome `json:"outcome"`
	Destruction float64       `json:"destruction"`
	Loot        int           `json:"loot"`
	Duration    int           `json:"duration"`
}

type TroopCatalog map[string]TroopStats

type TroopStats struct {
	ID              string
	HP              int
	DPS             int
	Speed           float64
	Range           float64
	PreferredTarget string
}

type troopState struct {
	id       string
	troopType string
	hp       float64
	maxHP    int
	dps      int
	speed    float64
	range_   float64
	pos      Point
	targetID string
	alive    bool
}

type buildingState struct {
	id           string
	buildingType string
	category     string
	hp           float64
	maxHP        int
	dps          int
	range_       float64
	pos          Point
	size         int
	alive        bool
}
