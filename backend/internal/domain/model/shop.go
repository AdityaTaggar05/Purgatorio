package model

type BuildingLevel struct {
	BuildingID      string   `json:"building_id"`
	Level           int      `json:"level"`
	HP              *int     `json:"hp"`
	DamagePerSec    *int     `json:"damage_per_sec"`
	ProductionRate  *int     `json:"production_rate"`
	StorageCapacity *int     `json:"storage_capacity"`
	AttackRange     *float64 `json:"attack_range"`
	UpgradeCost     int      `json:"upgrade_cost"`
	UpgradeTime     int      `json:"upgrade_time"`
}

type ShopItem struct {
	Building     Building        `json:"building"`
	CurrentOwned int             `json:"current_owned"`
	MaxAllowed   int             `json:"max_allowed"`
	CanBuy       bool            `json:"can_buy"`
	Levels       []BuildingLevel `json:"levels"`
}
