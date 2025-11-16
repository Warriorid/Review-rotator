package models

type Team struct {
    TeamID   int          `json:"-" db:"team_id"`
    TeamName string       `json:"team_name" db:"team_name" binding:"required"`
    Members  []TeamMember `json:"members" binding:"required"`
}