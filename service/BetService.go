package service

import (
	"PointCalculator/model"

	"gorm.io/gorm"
)

type BetService struct {
	db *gorm.DB
}

func NewBetService(db *gorm.DB) *BetService {
	return &BetService{db: db}
}

func (s *BetService) Betting(bettingPoint int, betWinOrLose bool, teamId int, targetTeamId int) bool {
	var team model.Team
	if err := s.db.Where("id = ? and use_yn = ?", teamId, "Y").First(&team).Error; err != nil {
		return false
	}

	if team.Point < bettingPoint {
		return false
	}
	if betWinOrLose {
		bettingPoint = DoublePoint(bettingPoint)
	}

	return true
}
