package model

import "time"

// 히스토리 타입 상수
const (
	HIST_TYPE_BET       = "BET"       // 베팅
	HIST_TYPE_BET_WIN   = "BET_WIN"   // 베팅 성공
	HIST_TYPE_BET_LOSE  = "BET_LOSE"  // 베팅 실패
	HIST_TYPE_MATCH_WIN = "MATCH_WIN" // 매치 승리
	HIST_TYPE_DOUBLE    = "DOUBLE"    // 더블 찬스 적용
	HIST_TYPE_TRIPLE    = "TRIPLE"    // 트리플 찬스 적용
)

type Hist struct {
	ID         int       `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	TeamID     int       `gorm:"column:team_id;foreignKey:ID" json:"team_id"`
	Type       string    `gorm:"column:type" json:"type"`               // 히스토리 타입
	Content    string    `gorm:"column:content" json:"content"`         // 내용 설명
	PrevPoint  int       `gorm:"column:prev_point" json:"prev_point"`   // 이전 포인트
	Point      int       `gorm:"column:point" json:"point"`             // 변동 포인트
	FinalPoint int       `gorm:"column:final_point" json:"final_point"` // 최종 포인트
	MatchID    *int      `gorm:"column:match_id" json:"match_id"`       // 관련 매치 ID (옵션)
	BetID      *int      `gorm:"column:bet_id" json:"bet_id"`           // 관련 베팅 ID (옵션)
	UseYn      string    `gorm:"column:use_yn" json:"use_yn"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
}

// 새 히스토리 생성
func NewHist(teamId int, histType string, content string, prevPoint int, point int, finalPoint int, matchId *int, betId *int) *Hist {
	return &Hist{
		TeamID:     teamId,
		Type:       histType,
		Content:    content,
		PrevPoint:  prevPoint,
		Point:      point,
		FinalPoint: finalPoint,
		MatchID:    matchId,
		BetID:      betId,
		UseYn:      "Y",
		CreatedAt:  time.Now(),
	}
}
