package model

import "github.com/google/uuid"

type Troop struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	TrainingCost    int     `json:"training_cost"`
	Space           int     `json:"space"`
	HP              int     `json:"hp"`
	DPS             int     `json:"dps"`
	Speed           float64 `json:"speed"`
	AttackRange     float64 `json:"attack_range"`
	PreferredTarget string  `json:"preferred_target"`
}

type UserArmy struct {
	UserID       uuid.UUID
	Troops       map[string]int `json:"troops"`
	UsedCapacity int            `json:"used_capacity"`
}

type MyTroopsResponse struct {
	Troops       map[string]int `json:"troops"`
	UsedCapacity int            `json:"used_capacity"`
	MaxCapacity  int            `json:"max_capacity"`
}
