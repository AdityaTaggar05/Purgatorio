package model

type User struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
	XP           int    `json:"xp"`
	Level        int    `json:"level"`
	TerraceLevel int    `json:"terrace_level"`
}
