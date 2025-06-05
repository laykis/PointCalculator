package service

import (
	"PointCalculator/model"
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MatchService struct {
	db          *gorm.DB
	gameService *GameService
	teamService *TeamService
}

func NewMatchService(db *gorm.DB) *MatchService {
	return &MatchService{
		db:          db,
		gameService: NewGameService(db),
		teamService: NewTeamService(db),
	}
}

// 매치 목록 조회
func (s *MatchService) GetMatchList() ([]model.Match, error) {
	var matches []model.Match

	// 매치 정보와 함께 게임과 팀 정보를 조인하여 조회
	query := s.db.Table("matches").
		Select("matches.*, games.name as game_name, t1.name as team1_name, t2.name as team2_name").
		Joins("LEFT JOIN games ON matches.game_id = games.id").
		Joins("LEFT JOIN teams t1 ON matches.player_team_id = t1.id").
		Joins("LEFT JOIN teams t2 ON matches.opponent_team_id = t2.id").
		Where("matches.use_yn = ?", "Y")

	// SQL 쿼리 출력
	sql := query.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Find(&matches)
	})
	fmt.Printf("실행된 SQL 쿼리: %s\n", sql)

	if err := query.Scan(&matches).Error; err != nil {
		return nil, err
	}

	if len(matches) > 0 {
		fmt.Printf("첫 번째 매치 정보: GameName=%s, Team1Name=%s, Team2Name=%s\n",
			matches[0].GameName, matches[0].Team1Name, matches[0].Team2Name)
	} else {
		fmt.Println("매치 목록이 비어있습니다.")
	}

	return matches, nil
}

// 매치 조회
func (s *MatchService) GetMatch(id int) (model.Match, error) {
	var match model.Match

	query := s.db.Table("matches").
		Select("matches.*, games.name as game_name, t1.name as team1_name, t2.name as team2_name").
		Joins("LEFT JOIN games ON matches.game_id = games.id AND games.use_yn = 'Y'").
		Joins("LEFT JOIN teams t1 ON matches.player_team_id = t1.id AND t1.use_yn = 'Y'").
		Joins("LEFT JOIN teams t2 ON matches.opponent_team_id = t2.id AND t2.use_yn = 'Y'").
		Where("matches.id = ? AND matches.use_yn = ?", id, "Y")

	if err := query.First(&match).Error; err != nil {
		return model.Match{}, errors.New("match not found")
	}
	return match, nil
}

// 매치 생성
func (s *MatchService) CreateMatch(inputMatch model.Match) (model.Match, error) {
	fmt.Printf("CreateMatch 입력 데이터: %+v\n", inputMatch)

	// 동일한 매치가 이미 존재하는지 확인
	if s.db.Where("game_id = ? AND player_team_id = ? AND opponent_team_id = ? AND use_yn = ? AND status != ?",
		inputMatch.GameId, inputMatch.PlayerTeamId, inputMatch.OpponentTeamId, "Y", "C").First(&model.Match{}).Error == nil {
		return model.Match{}, errors.New("match already exists")
	}

	match := model.NewMatch(inputMatch.GameId, inputMatch.PlayerTeamId, inputMatch.OpponentTeamId)
	fmt.Printf("생성된 매치 데이터: %+v\n", match)

	if err := s.db.Create(match).Error; err != nil {
		fmt.Printf("매치 DB 생성 실패: %v\n", err)
		return model.Match{}, err
	}
	return *match, nil
}

// 매치 수정
func (s *MatchService) UpdateMatch(inputMatch model.Match) (model.Match, error) {
	match, err := s.GetMatch(inputMatch.ID)
	if err != nil {
		return model.Match{}, err
	}

	match.GameId = inputMatch.GameId
	match.PlayerTeamId = inputMatch.PlayerTeamId
	match.OpponentTeamId = inputMatch.OpponentTeamId
	match.Status = inputMatch.Status

	if err := s.db.Save(&match).Error; err != nil {
		return model.Match{}, err
	}
	return match, nil
}

// 매치 삭제
func (s *MatchService) DeleteMatch(inputMatch model.Match) (model.Match, error) {
	fmt.Printf("삭제할 매치 ID: %d\n", inputMatch.ID)

	var match model.Match
	if err := s.db.Where("id = ? and use_yn = ?", inputMatch.ID, "Y").First(&match).Error; err != nil {
		fmt.Printf("매치 조회 실패: %v\n", err)
		return model.Match{}, errors.New("match not found")
	}

	fmt.Printf("삭제할 매치 정보: %+v\n", match)

	err := s.db.Model(&match).Update("use_yn", "N").Error
	if err != nil {
		fmt.Printf("매치 삭제 실패: %v\n", err)
		return model.Match{}, err
	}
	return match, nil
}

// 매치 결과 처리 및 승리 팀 포인트 지급
func (s *MatchService) ProcessMatchResult(matchId int, winnerTeamId int) error {
	// 트랜잭션 시작
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// 매치 정보 조회
	var match model.Match
	if err := tx.Where("id = ? and use_yn = ?", matchId, "Y").First(&match).Error; err != nil {
		tx.Rollback()
		return errors.New("match not found")
	}

	// 이미 완료된 매치인지 확인
	if match.Status == "C" {
		tx.Rollback()
		return errors.New("match already completed")
	}

	// 승리 팀 포인트 업데이트 (기본 승리 포인트 5점)
	if err := tx.Model(&model.Team{}).
		Where("id = ?", winnerTeamId).
		Update("point", gorm.Expr("point + ?", 5)).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 매치 상태 업데이트
	if err := tx.Model(&match).
		Updates(map[string]interface{}{
			"status":     "C",
			"updated_at": time.Now(),
		}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 베팅 결과 처리
	var bets []model.Bet
	if err := tx.Where("match_id = ? AND use_yn = ?", matchId, "Y").Find(&bets).Error; err != nil {
		tx.Rollback()
		return err
	}

	for _, bet := range bets {
		// 베팅 결과 확인
		isWin := false
		if bet.BetType == "W" && bet.TargetTeamId == winnerTeamId {
			// 승리 베팅이 맞은 경우
			isWin = true
		} else if bet.BetType == "L" && bet.TargetTeamId != winnerTeamId {
			// 패배 베팅이 맞은 경우
			isWin = true
		}

		// 베팅 결과에 따른 포인트 처리
		if isWin {
			var addPoint int
			if bet.BettingPoint == 0 {
				// 베팅 포인트가 0점이면 1점만 지급
				addPoint = 1
			} else {
				// 베팅 성공 시 베팅 포인트의 2배를 지급
				addPoint = bet.BettingPoint * 2
			}

			if err := tx.Model(&model.Team{}).
				Where("id = ?", bet.TeamId).
				Update("point", gorm.Expr("point + ?", addPoint)).Error; err != nil {
				tx.Rollback()
				return err
			}
		}

		// 베팅 상태 업데이트
		if err := tx.Model(&bet).
			Updates(map[string]interface{}{
				"status":     "C",
				"updated_at": time.Now(),
			}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// 트랜잭션 커밋
	return tx.Commit().Error
}

// 진행중인 매치 목록 조회
func (s *MatchService) GetActiveMatches() ([]gin.H, error) {
	var matches []model.Match
	if err := s.db.Where("status = ? AND use_yn = ?", "P", "Y").Find(&matches).Error; err != nil {
		return nil, err
	}

	var result []gin.H
	for _, match := range matches {
		// 게임 정보 조회
		game, err := s.gameService.GetGame(match.GameId)
		if err != nil {
			continue
		}

		// 팀1 정보 조회
		team1, err := s.teamService.GetTeam(match.PlayerTeamId)
		if err != nil {
			continue
		}

		// 팀2 정보 조회
		team2, err := s.teamService.GetTeam(match.OpponentTeamId)
		if err != nil {
			continue
		}

		matchInfo := gin.H{
			"ID":             match.ID,
			"GameName":       game.Name,
			"Team1Name":      team1.Name,
			"Team2Name":      team2.Name,
			"PlayerTeamId":   match.PlayerTeamId,
			"OpponentTeamId": match.OpponentTeamId,
			"Status":         match.Status,
			"BetCount":       0, // 기본값 설정 (실제 값은 호출하는 쪽에서 설정)
		}
		result = append(result, matchInfo)
	}

	return result, nil
}
