package repository

import (
	"review-rotator/internal/models"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	Team
}

type Team interface {
	AddTeam(team models.Team) (*models.Team, error)
	GetTeam(teamName string) (*models.Team, error)
}


func NewRepository(db *sqlx.DB) *Repository {
    return &Repository{
		Team: NewTeamPostgres(db),
    }
}