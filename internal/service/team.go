package service

import (
	"database/sql"
	"review-rotator/internal/models"
	"review-rotator/internal/repository"
)

type TeamService struct {
	repo repository.Team
}

func NewTeamService(repo repository.Team) *TeamService {
	return &TeamService{repo: repo}
}

func (s *TeamService) AddTeam(team models.Team)(*models.Team, error){
	_, err := s.repo.GetTeam(team.TeamName)
	if err == nil {
		return nil, models.ErrTeamExists
	}
	if err != sql.ErrNoRows {
		return nil, err
	}
	var newTeam *models.Team
	newTeam, err = s.repo.AddTeam(team)
	if err != nil {
		return nil, err
	}

	return newTeam, nil
}

func (s *TeamService) GetTeam(teamName string) (*models.Team, error) {
	team, err := s.repo.GetTeam(teamName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return team, nil
}