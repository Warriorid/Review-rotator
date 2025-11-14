package repository

import (
	"database/sql"
	"review-rotator/internal/models"

	"github.com/jmoiron/sqlx"
)

type UserPostgres struct {
	db *sqlx.DB
}

func NewUserPostgres(db *sqlx.DB) *UserPostgres {
	return &UserPostgres{db: db}
}


func (r *UserPostgres) SetUserActive(userID string, isActive bool) (*models.User, error) {
	var user models.User
	
	query := `
		UPDATE users 
		SET is_active = $1 
		WHERE user_id = $2
		RETURNING user_id, username, team_id, is_active
	`
	
	err := r.db.QueryRow(query, isActive, userID).Scan(
		&user.UserID, 
		&user.Username, 
		&user.TeamID, 
		&user.IsActive,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	
	err = r.db.QueryRow(
		"SELECT team_name FROM teams WHERE team_id = $1", 
		user.TeamID,
	).Scan(&user.TeamName)
	if err != nil {
		return nil, err
	}
	
	return &user, nil
}