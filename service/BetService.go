package service

import (
	"PointCalculator/model"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type BetService struct {
	db          *gorm.DB
	histService *HistService
}

func NewBetService(db *gorm.DB) *BetService {
	return &BetService{
		db:          db,
		histService: NewHistService(db),
	}
}

func (s *BetService) CreateBet(inputBet model.Bet) (bool, error) {
	// 트랜잭션 시작
	tx := s.db.Begin()
	if tx.Error != nil {
		return false, tx.Error
	}

	// 동일한 매치에 대한 중복 베팅 체크
	var existingBet model.Bet
	if err := tx.Where("match_id = ? AND team_id = ? AND use_yn = ? AND status = ?",
		inputBet.MatchID, inputBet.TeamId, "Y", "P").First(&existingBet).Error; err == nil {
		tx.Rollback()
		return false, errors.New("이미 이 매치에 베팅한 기록이 있습니다")
	}

	// 팀 포인트 조회
	var team model.Team
	if err := tx.Table("teams").Where("id = ?", inputBet.TeamId).First(&team).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	if team.Point < inputBet.BettingPoint {
		tx.Rollback()
		return false, errors.New("팀 포인트가 부족합니다.")
	}

	// 베팅 포인트 차감
	newPoint := team.Point - inputBet.BettingPoint
	if err := tx.Model(&model.Team{}).
		Where("id = ?", inputBet.TeamId).
		Update("point", newPoint).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	// 베팅 생성
	bet := model.NewBet(
		inputBet.MatchID,
		inputBet.TeamId,
		inputBet.TargetTeamId,
		inputBet.BettingPoint,
		inputBet.BetType,
		inputBet.IsDouble,
		inputBet.IsTriple,
	)

	if err := tx.Create(bet).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	// 베팅 히스토리 생성
	hist := model.NewHist(
		inputBet.TeamId,
		model.HIST_TYPE_BET,
		fmt.Sprintf("매치 #%d에 %d 포인트 베팅", inputBet.MatchID, inputBet.BettingPoint),
		team.Point,
		-inputBet.BettingPoint,
		newPoint,
		&inputBet.MatchID,
		&bet.ID,
	)

	if err := tx.Create(hist).Error; err != nil {
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

	// 팀 정보 조회
	var team model.Team
	if err := tx.Where("id = ?", bet.TeamId).First(&team).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 베팅 포인트 반환
	newPoint := team.Point + bet.BettingPoint
	if err := tx.Model(&model.Team{}).
		Where("id = ?", bet.TeamId).
		Update("point", newPoint).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 베팅 취소 히스토리 생성
	hist := model.NewHist(
		bet.TeamId,
		model.HIST_TYPE_BET,
		fmt.Sprintf("매치 #%d 베팅 취소로 %d 포인트 반환", bet.MatchID, bet.BettingPoint),
		team.Point,
		bet.BettingPoint,
		newPoint,
		&bet.MatchID,
		&bet.ID,
	)

	if err := tx.Create(hist).Error; err != nil {
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

// UpdateBet 베팅 정보를 수정합니다.
func (s *BetService) UpdateBet(bet model.Bet) error {
	// 트랜잭션 시작
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 기존 베팅 정보 조회
	var currentBet model.Bet
	if err := tx.Where("id = ? AND use_yn = ?", bet.ID, "Y").First(&currentBet).Error; err != nil {
		tx.Rollback()
		return errors.New("bet not found")
	}

	// 베팅이 진행중인 상태인지 확인
	if currentBet.Status != "P" {
		tx.Rollback()
		return errors.New("completed bet cannot be modified")
	}

	// 기존 베팅 포인트와 새로운 베팅 포인트의 차이 계산
	pointDiff := currentBet.BettingPoint - bet.BettingPoint

	// 팀 정보 조회
	var team model.Team
	if err := tx.Where("id = ?", currentBet.TeamId).First(&team).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 포인트가 증가하는 경우 팀 포인트 잔액 확인
	if pointDiff < 0 && team.Point < -pointDiff {
		tx.Rollback()
		return errors.New("팀 포인트가 부족합니다")
	}

	// 베팅 정보 업데이트
	if err := tx.Model(&currentBet).Updates(map[string]interface{}{
		"bet_type":      bet.BetType,
		"betting_point": bet.BettingPoint,
		"is_double":     bet.IsDouble,
		"is_triple":     bet.IsTriple,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 팀 포인트 업데이트
	newPoint := team.Point + pointDiff
	if err := tx.Model(&model.Team{}).
		Where("id = ?", currentBet.TeamId).
		Update("point", newPoint).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 베팅 수정 히스토리 생성
	hist := model.NewHist(
		currentBet.TeamId,
		model.HIST_TYPE_BET,
		fmt.Sprintf("매치 #%d 베팅 수정으로 %d 포인트 변동", currentBet.MatchID, pointDiff),
		team.Point,
		pointDiff,
		newPoint,
		&currentBet.MatchID,
		&currentBet.ID,
	)

	if err := tx.Create(hist).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
