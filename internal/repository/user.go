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


func (r *UserPostgres) GetUserByID(userID string) error {
    var user models.User
    query := `
        SELECT u.user_id, u.username, u.team_id, u.is_active, t.team_name
        FROM users u
        JOIN teams t ON u.team_id = t.team_id
        WHERE u.user_id = $1
    `
    
    err := r.db.QueryRow(query, userID).Scan(
        &user.UserID,
        &user.Username,
        &user.TeamID,
        &user.IsActive,
        &user.TeamName,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return models.ErrNotFound
        }
        return err
    }
    
    return nil
}

func (r *UserPostgres) GetActiveTeamMembers(teamID int, excludeUserID string) ([]models.User, error) {
    var users []models.User
    query := `
        SELECT u.user_id, u.username, u.team_id, u.is_active, t.team_name
        FROM users u
        JOIN teams t ON u.team_id = t.team_id
        WHERE u.team_id = $1 AND u.is_active = true AND u.user_id != $2
        ORDER BY u.user_id
    `
    
    rows, err := r.db.Query(query, teamID, excludeUserID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var user models.User
        err := rows.Scan(
            &user.UserID,
            &user.Username,
            &user.TeamID,
            &user.IsActive,
            &user.TeamName,
        )
        if err != nil {
            return nil, err
        }
        users = append(users, user)
    }

    return users, nil
}