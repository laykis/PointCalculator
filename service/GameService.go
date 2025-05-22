package service

import (
	"PointCalculator/model"

	"gorm.io/gorm"
)

type GameService struct {
	db *gorm.DB
}

func NewGameService(db *gorm.DB) *GameService {
	return &GameService{db: db}
}

func (s *GameService) CreateGame(gameName string) (model.Game, error) {
	game := model.NewGame(gameName)
	if err := s.db.Create(game).Error; err != nil {
		return model.Game{}, err
	}
	return *game, nil
}

func (s *GameService) GetGameList() ([]model.Game, error) {
	var games []model.Game
	if err := s.db.Where("use_yn = ?", "Y").Find(&games).Error; err != nil {
		return nil, err
	}
	return games, nil
}
