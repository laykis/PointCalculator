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
	// 트랜잭션 시작
	tx := s.db.Begin()
	if tx.Error != nil {
		return false, tx.Error
	}

	// 팀 포인트 조회
	var teamPoint model.Team
	if err := tx.Table("teams").Where("id = ?", inputBet.TeamId).First(&teamPoint).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	point := teamPoint.Point

	if point < inputBet.BettingPoint {
		tx.Rollback()
		return false, errors.New("팀 포인트가 부족합니다.")
	}

	// 베팅 포인트 차감
	if err := tx.Model(&model.Team{}).
		Where("id = ?", inputBet.TeamId).
		Update("point", gorm.Expr("point - ?", inputBet.BettingPoint)).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	// 베팅 생성
	bet := model.NewBet(inputBet.MatchID, inputBet.TeamId, inputBet.TargetTeamId, inputBet.BettingPoint, inputBet.BetType)
	if err := tx.Create(bet).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	// 트랜잭션 커밋
	if err := tx.Commit().Error; err != nil {
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
	if err := s.db.Where("use_yn = ?", "Y").Find(&bets).Error; err != nil {
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

// 베팅 삭제
func (s *BetService) DeleteBet(id int) error {
	// 트랜잭션 시작
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// 베팅 정보 조회
	var bet model.Bet
	if err := tx.Where("id = ? AND use_yn = ?", id, "Y").First(&bet).Error; err != nil {
		tx.Rollback()
		return errors.New("bet not found")
	}

	// 베팅이 진행중인 상태인지 확인
	if bet.Status != "P" {
		tx.Rollback()
		return errors.New("completed bet cannot be deleted")
	}

	// 베팅 포인트 반환
	if err := tx.Model(&model.Team{}).
		Where("id = ?", bet.TeamId).
		Update("point", gorm.Expr("point + ?", bet.BettingPoint)).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 소프트 삭제 (use_yn = 'N'으로 변경)
	if err := tx.Model(&bet).Update("use_yn", "N").Error; err != nil {
		tx.Rollback()
		return err
	}

	// 트랜잭션 커밋
	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

// 팀별 진행중인 베팅 수 조회
func (s *BetService) GetActiveBetCountByTeamId(teamId int) (int, error) {
	var count int64
	if err := s.db.Model(&model.Bet{}).
		Where("team_id = ? AND status = ? AND use_yn = ?", teamId, "P", "Y").
		Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

// 매치별 베팅 수 조회
func (s *BetService) GetBetCountByMatchId(matchId int) (int, error) {
	var count int64
	if err := s.db.Model(&model.Bet{}).
		Where("match_id = ? AND status = ? AND use_yn = ?", matchId, "P", "Y").
		Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}
