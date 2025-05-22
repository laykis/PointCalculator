package model

import "time"

type Hist struct {
	ID        int       `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	TeamID    int       `gorm:"column:team_id;foreignKey:ID" json:"team_id"`
	Content   string    `gorm:"column:content" json:"content"`
	PrevPoint int       `gorm:"column:prev_point" json:"prev_point"`
	Point     int       `gorm:"column:point" json:"point"`
	UseYn     string    `gorm:"column:use_yn" json:"use_yn"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}
