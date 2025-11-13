package repository

import (
	"review-rotator/internal/models"
	"github.com/jmoiron/sqlx"
)

type TeamPostgres struct {
	db *sqlx.DB
}

func NewTeamPostgres(db *sqlx.DB) *TeamPostgres{
	return &TeamPostgres{db: db}
}

func (r *TeamPostgres) AddTeam(team models.Team) (*models.Team, error) {
    tx, err := r.db.Beginx()
    if err != nil {
        return nil, err
    }

    defer tx.Rollback()
    var teamID int
    err = tx.QueryRow(
        "INSERT INTO teams (team_name) VALUES ($1) RETURNING team_id", 
        team.TeamName,
    ).Scan(&teamID)
    if err != nil {
        return nil, err
    }

    for _, member := range team.Members {
        _, err = tx.Exec(`
            INSERT INTO users (user_id, username, team_id, is_active) 
            VALUES ($1, $2, $3, $4)
            ON CONFLICT (user_id) 
            DO UPDATE SET username = $2, team_id = $3, is_active = $4
        `, member.UserID, member.Username, teamID, member.IsActive)
        if err != nil {
            return nil, err
        }
    }

    if err := tx.Commit(); err != nil {
        return nil, err
    }
    team.TeamID = teamID
    return &team, nil
}

func (r *TeamPostgres) GetTeam(teamName string) (*models.Team, error) {
    var team models.Team
    err := r.db.QueryRow(
        "SELECT team_id, team_name FROM teams WHERE team_name = $1", 
        teamName,
    ).Scan(&team.TeamID, &team.TeamName)
    if err != nil {
        return nil, err
    }
    query := `
        SELECT u.user_id, u.username, u.is_active, t.team_name
        FROM users u
        JOIN teams t ON u.team_id = t.team_id
        WHERE u.team_id = $1
    `
    rows, err := r.db.Query(query, team.TeamID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var member models.TeamMember
        var teamName string
        err := rows.Scan(&member.UserID, &member.Username, &member.IsActive, &teamName)
        if err != nil {
            return nil, err
        }
        team.Members = append(team.Members, member)
    }

    return &team, nil
}