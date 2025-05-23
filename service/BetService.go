package service

import (
	"PointCalculator/model"
	"errors"
	"gorm.io/gorm"
)

type BetService struct {
	db *gorm.DB
}

func NewBetService(db *gorm.DB) *BetService {
	return &BetService{db: db}
}

func (s *BetService) CreateBet(inputBet model.Bet) (bool, error) {

	teamPoint := s.db.Where("team_id = ?", inputBet.TeamId).First(&model.Tean{}).Error; err != nil {
		return false, errors.New("팀 정보가 없습니다.")
	}

	if teamPoint < inputBet.BettingPoint {
	

	bet := model.NewBet(inputBet.MatchID, inputBet.TeamId, inputBet.TargetTeamId, inputBet.BettingPoint)
	if err := s.db.Create(bet).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (s *BetService) GetBet(id int) (*model.Bet, error) {
	var bet model.Bet
	if err := s.db.Where("id = ?", id).First(&bet).Error; err != nil {
		return nil, err
	}
	return &bet, nil
}

func (s *BetService) GetBetList() ([]model.Bet, error) {
	var bets []model.Bet
	if err := s.db.Find(&bets).Error; err != nil {
		return nil, err
	}
	return bets, nil
}
