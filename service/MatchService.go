package service

import (
	"PointCalculator/model"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type MatchService struct {
	db *gorm.DB
}

func NewMatchService(db *gorm.DB) *MatchService {
	return &MatchService{db: db}
}

// 매치 목록 조회
func (s *MatchService) GetMatchList() ([]model.Match, error) {
	var matches []model.Match

	// 매치 정보와 함께 게임과 팀 정보를 조인하여 조회
	query := s.db.Table("matches").
		Select("matches.*, games.name as game_name, t1.name as team1_name, t2.name as team2_name").
		Joins("LEFT JOIN games ON matches.game_id = games.id AND games.use_yn = 'Y'").
		Joins("LEFT JOIN teams t1 ON matches.player_team_id = t1.id AND t1.use_yn = 'Y'").
		Joins("LEFT JOIN teams t2 ON matches.opponent_team_id = t2.id AND t2.use_yn = 'Y'").
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
	if err := s.db.Where("id = ? and use_yn = ?", id, "Y").First(&match).Error; err != nil {
		return model.Match{}, errors.New("match not found")
	}
	return match, nil
}

// 매치 생성
func (s *MatchService) CreateMatch(inputMatch model.Match) (model.Match, error) {
	fmt.Printf("CreateMatch 입력 데이터: %+v\n", inputMatch)

	// 동일한 매치가 이미 존재하는지 확인
	if s.db.Where("game_id = ? AND player_team_id = ? AND opponent_team_id = ? AND use_yn = ?",
		inputMatch.GameId, inputMatch.PlayerTeamId, inputMatch.OpponentTeamId, "Y").First(&model.Match{}).Error == nil {
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
