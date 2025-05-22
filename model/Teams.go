package model

import "time"

type Team struct {
	ID        int       `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Name      string    `gorm:"column:name" json:"name"`
	Point     int       `gorm:"column:point" json:"point"`
	UseYn     string    `gorm:"column:use_yn" json:"use_yn"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}

// NewTeam creates a new Team instance with default values
func NewTeam(name string) *Team {
	return &Team{
		Name:      name,
		Point:     0,
		UseYn:     "Y",
		CreatedAt: time.Now(),
	}
}
