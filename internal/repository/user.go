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

func (r *UserPostgres) GetUserTeamID(userID string) (int, error) {
    var teamID int
    query := `SELECT team_id FROM users WHERE user_id = $1`
    
    err := r.db.QueryRow(query, userID).Scan(&teamID)
    if err != nil {
        if err == sql.ErrNoRows {
            return 0, models.ErrNotFound
        }
        return 0, err
    }
    
    return teamID, nil
}

func (r *UserPostgres) GetUserReviewRequests(userID string) ([]models.PullRequestShort, error) {
    query := `
        SELECT pr.pull_request_id, pr.pull_request_name, pr.author_id, pr.status
        FROM pull_requests pr
        JOIN pr_reviewers prr ON pr.pull_request_id = prr.pull_request_id
        WHERE prr.reviewer_id = $1
        ORDER BY pr.created_at DESC
    `
    
    rows, err := r.db.Query(query, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var pullRequests []models.PullRequestShort
    for rows.Next() {
        var pr models.PullRequestShort
        err := rows.Scan(
            &pr.PullRequestID,
            &pr.PullRequestName,
            &pr.AuthorID,
            &pr.Status,
        )
        if err != nil {
            return nil, err
        }
        pullRequests = append(pullRequests, pr)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return pullRequests, nil
}