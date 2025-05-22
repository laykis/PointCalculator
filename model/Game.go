package model

import (
	"time"
)

type Game struct {
	ID        int       `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Name      string    `gorm:"column:name" json:"name"`
	UseYn     string    `gorm:"column:use_yn" json:"use_yn"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func NewGame(name string) *Game {
	return &Game{
		Name:      name,
		UseYn:     "Y",
		CreatedAt: time.Now(),
	}
}
