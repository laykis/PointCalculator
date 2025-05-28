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
	var teamPoint model.Team
	if err := s.db.Table("teams").Where("team_id = ?", inputBet.TeamId).First(&teamPoint).Error; err != nil {
		return false, err
	}

	point := teamPoint.Point

	if point < inputBet.BettingPoint {
		return false, errors.New("팀 포인트가 부족합니다.")
	}

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

// 매치 상태에 따른 베팅 상태 업데이트
func (s *BetService) UpdateBetsByMatchStatus(matchId int) error {
	// 매치가 완료되면 해당 매치의 모든 베팅을 완료 상태로 변경
	if err := s.db.Model(&model.Bet{}).
		Where("match_id = ? AND use_yn = ?", matchId, "Y").
		Update("status", "C").Error; err != nil {
		return err
	}
	return nil
}

// 매치 ID로 베팅 목록 조회
func (s *BetService) GetBetsByMatchId(matchId int) ([]model.Bet, error) {
	var bets []model.Bet
	if err := s.db.Where("match_id = ? AND use_yn = ?", matchId, "Y").Find(&bets).Error; err != nil {
		return nil, err
	}
	return bets, nil
}
