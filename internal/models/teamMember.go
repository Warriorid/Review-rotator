package models

type TeamMember struct {
    UserID   string `json:"user_id" binding:"required"`
    Username string `json:"username" binding:"required"`
    TeamID   int    `json:"team_id" binding:"required"`
    IsActive bool   `json:"is_active" binding:"required"`
}