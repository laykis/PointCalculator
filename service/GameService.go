package service

import (
	"PointCalculator/model"
	"errors"

	"gorm.io/gorm"
)

type GameService struct {
	db *gorm.DB
}

func NewGameService(db *gorm.DB) *GameService {
	return &GameService{db: db}
}

func (s *GameService) CreateGame(gameName string) (model.Game, error) {
	if s.db.Where("name = ? and use_yn = ?", gameName, "Y").First(&model.Game{}).Error == nil {
		return model.Game{}, errors.New("game already exists")
	}

	game := model.NewGame(gameName)
	if err := s.db.Create(game).Error; err != nil {
		return model.Game{}, err
	}
	return *game, nil
}

func (s *GameService) UpdateGame(inputGame model.Game) (model.Game, error) {
	game, err := s.GetGame(inputGame)
	if err != nil {
		return model.Game{}, err
	}

	game.Name = inputGame.Name
	if err := s.db.Save(&game).Error; err != nil {
		return model.Game{}, err
	}
	return game, nil
}

func (s *GameService) DeleteGame(inputGame model.Game) (model.Game, error) {
	var game model.Game
	if err := s.db.Where("id = ? and use_yn = ?", inputGame.ID, "Y").First(&game).Error; err != nil {
		return model.Game{}, errors.New("game not found")
	}

	err := s.db.Model(&game).Update("use_yn", "N").Error
	if err != nil {
		return model.Game{}, err
	}
	return game, nil
}

func (s *GameService) GetGame(inputGame model.Game) (model.Game, error) {
	var game model.Game
	if err := s.db.Where("id = ? and use_yn = ?", inputGame.ID, "Y").First(&game).Error; err != nil {
		return model.Game{}, errors.New("game not found")
	}
	return game, nil
}

func (s *GameService) GetGameList() ([]model.Game, error) {
	var games []model.Game
	if err := s.db.Where("use_yn = ?", "Y").Find(&games).Error; err != nil {
		return nil, err
	}
	return games, nil
}
