package repository

import (
	"review-rotator/internal/models"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	Team
	User
}

type Team interface {
	AddTeam(team models.Team) (*models.Team, error)
	GetTeam(teamName string) (*models.Team, error)
}
type User interface {
	SetUserActive(userID string, isActive bool) (*models.User, error)
}


func NewRepository(db *sqlx.DB) *Repository {
    return &Repository{
		Team: NewTeamPostgres(db),
		User: NewUserPostgres(db),
    }
}