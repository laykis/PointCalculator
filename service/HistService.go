package service

import (
	"PointCalculator/model"
	"fmt"

	"gorm.io/gorm"
)

type HistService struct {
	db *gorm.DB
}

func NewHistService(db *gorm.DB) *HistService {
	return &HistService{db: db}
}

// 히스토리 생성
func (s *HistService) CreateHist(hist *model.Hist) error {
	return s.db.Create(hist).Error
}

// 히스토리 목록 조회 (최신순)
func (s *HistService) GetHistList() ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	// 히스토리와 팀 정보를 조인하여 조회
	query := s.db.Table("hists h").
		Select("h.*, t.name as team_name, m.game_id, g.name as game_name").
		Joins("LEFT JOIN teams t ON h.team_id = t.id").
		Joins("LEFT JOIN matches m ON h.match_id = m.id").
		Joins("LEFT JOIN games g ON m.game_id = g.id").
		Where("h.use_yn = ?", "Y").
		Order("h.created_at DESC")

	if err := query.Find(&results).Error; err != nil {
		return nil, err
	}

	// 히스토리 타입에 따른 내용 포맷팅
	for i := range results {
		switch results[i]["type"] {
		case model.HIST_TYPE_BET:
			results[i]["type_text"] = "베팅"
		case model.HIST_TYPE_BET_WIN:
			results[i]["type_text"] = "베팅 성공"
		case model.HIST_TYPE_BET_LOSE:
			results[i]["type_text"] = "베팅 실패"
		case model.HIST_TYPE_MATCH_WIN:
			results[i]["type_text"] = "매치 승리"
		case model.HIST_TYPE_DOUBLE:
			results[i]["type_text"] = "더블 찬스"
		case model.HIST_TYPE_TRIPLE:
			results[i]["type_text"] = "트리플 찬스"
		}

		// 포인트 변동 표시
		point := results[i]["point"].(int64)
		if point > 0 {
			results[i]["point_text"] = fmt.Sprintf("+%d", point)
		} else {
			results[i]["point_text"] = fmt.Sprintf("%d", point)
		}
	}

	return results, nil
}
