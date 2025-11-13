package models

type Team struct {
	TeamName string       `json:"team_name" binding:"required"`
	Members  []TeamMember `json:"members" binding:"required"`
}