package model

import (
	"time"
)

type Bet struct {
	ID           int       `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	GameId       int       `gorm:"column:game_id;foreignKey:ID" json:"game_id"` // 게임의 ID
	TeamId       int       `gorm:"column:team_id;foreignKey:ID" json:"team_id"` // 베팅하는 팀의 ID
	TargetTeamId int       `gorm:"column:target_team_id" json:"target_team_id"` // 베팅 대상 팀의 ID
	Status       string    `gorm:"column:status" json:"status"`
	UseYn        string    `gorm:"column:use_yn" json:"use_yn"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func NewBet(gameId int, teamId int, targetTeamId int) *Bet {

	return &Bet{
		GameId:       gameId,
		TeamId:       teamId,
		TargetTeamId: targetTeamId,
		Status:       "P",
		UseYn:        "Y",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}
