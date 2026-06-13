package model

type BuildingLevel struct {
	BuildingID      string   `json:"building_id"`
	Level           int      `json:"level"`
	HP              *int     `json:"hp,omitempty"`
	DamagePerSec    *int     `json:"damage_per_sec,omitempty"`
	ProductionRate  *int     `json:"production_rate,omitempty"`
	StorageCapacity *int     `json:"storage_capacity,omitempty"`
	AttackRange     *float64 `json:"attack_range,omitempty"`
	UpgradeCost     int      `json:"upgrade_cost"`
	UpgradeTime     int      `json:"upgrade_time"`
}

type ShopItem struct {
	Building     Building `json:"building"`
	CurrentOwned int      `json:"current_owned"`
	MaxAllowed   int      `json:"max_allowed"`
	CanBuy       bool     `json:"can_buy"`
}
