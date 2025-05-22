package service

import (
	"PointCalculator/model"
	"errors"

	"gorm.io/gorm"
)

type TeamService struct {
	db *gorm.DB
}

func NewTeamService(db *gorm.DB) *TeamService {
	return &TeamService{db: db}
}

func (s *TeamService) GetTeamList() ([]model.Team, error) {
	var teams []model.Team
	if err := s.db.Where("use_yn = ?", "Y").Find(&teams).Error; err != nil {
		return nil, err
	}
	return teams, nil
}

func (s *TeamService) GetTeam(id int) (model.Team, error) {
	var team model.Team
	if err := s.db.Where("id = ? and use_yn = ?", id, "Y").First(&team).Error; err != nil {
		return model.Team{}, err
	}
	return team, nil
}

func (s *TeamService) CreateTeam(name string) (model.Team, error) {

	if s.db.Where("name = ? and use_yn = ?", name, "Y").First(&model.Team{}).Error == nil {
		return model.Team{}, errors.New("team already exists")
	}

	team := model.NewTeam(name)
	if err := s.db.Create(team).Error; err != nil {
		return model.Team{}, err
	}
	return *team, nil
}

func (s *TeamService) UpdateTeam(inputTeam model.Team) (model.Team, error) {
	team, err := s.GetTeam(inputTeam.ID)
	if err != nil {
		return model.Team{}, err
	}

	team.Name = inputTeam.Name
	if err := s.db.Save(&team).Error; err != nil {
		return model.Team{}, err
	}
	return team, nil
}

func (s *TeamService) DeleteTeam(inputTeam model.Team) (model.Team, error) {
	var team model.Team
	if err := s.db.Where("id = ? and use_yn = ?", inputTeam.ID, "Y").First(&team).Error; err != nil {
		return model.Team{}, errors.New("team not found")
	}

	err := s.db.Model(&team).Update("use_yn", "N").Error
	if err != nil {
		return model.Team{}, err
	}

	return team, nil
}
