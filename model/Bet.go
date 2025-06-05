package model

import (
	"time"
)

type Bet struct {
	ID           int       `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	MatchID      int       `gorm:"column:match_id;foreignKey:ID" json:"match_id"` // 게임의 ID
	TeamId       int       `gorm:"column:team_id;foreignKey:ID" json:"team_id"`   // 베팅하는 팀의 ID
	TargetTeamId int       `gorm:"column:target_team_id" json:"target_team_id"`   // 베팅 대상 팀의 ID
	BettingPoint int       `gorm:"column:betting_point" json:"betting_point"`     // 베팅 포인트
	BetType      string    `gorm:"column:bet_type" json:"bet_type"`               // 베팅 유형 (W: 승리, L: 패배)
	Status       string    `gorm:"column:status" json:"status"`                   // 베팅 상태 (P: 진행중, C: 완료)
	UseYn        string    `gorm:"column:use_yn" json:"use_yn"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func NewBet(matchId int, teamId int, targetTeamId int, bettingPoint int, betType string) *Bet {
	return &Bet{
		MatchID:      matchId,
		TeamId:       teamId,
		TargetTeamId: targetTeamId,
		BettingPoint: bettingPoint,
		BetType:      betType,
		Status:       "P",
		UseYn:        "Y",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}
