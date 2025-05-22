package model

import (
	"time"
)

type Match struct {
	ID             int       `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	GameId         int       `gorm:"column:game_id;foreignKey:ID" json:"game_id"`                   // 게임의 ID
	PlayerTeamId   int       `gorm:"column:player_team_id;foreignKey:ID" json:"player_team_id"`     // 플레이어 팀의 ID
	OpponentTeamId int       `gorm:"column:opponent_team_id;foreignKey:ID" json:"opponent_team_id"` // 상대 팀의 ID
	Status         string    `gorm:"column:status" json:"status"`                                   // 경기 상태
	UseYn          string    `gorm:"column:use_yn" json:"use_yn"`                                   // 사용 여부
	CreatedAt      time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func NewMatch(gameId int, playerTeamId int, opponentTeamId int) *Match {
	return &Match{
		GameId:         gameId,
		PlayerTeamId:   playerTeamId,
		OpponentTeamId: opponentTeamId,
		Status:         "P",
		UseYn:          "Y",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}
