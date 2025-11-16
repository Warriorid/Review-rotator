package models

type User struct {
    UserID   string `json:"user_id" db:"user_id" binding:"required"`
    Username string `json:"username" db:"username" binding:"required"`
    TeamID   int    `json:"-" db:"team_id" binding:"required"`
    TeamName string `json:"team_name" db:"team_name"`                    
    IsActive bool   `json:"is_active" db:"is_active" binding:"required"`
}

type SetUserActiveRequest struct {
	UserID   string `json:"user_id" binding:"required"`
	IsActive bool   `json:"is_active"`
}

type UserReviewsResponse struct {
    UserID        string             `json:"user_id" binding:"required"`
    PullRequests  []PullRequestShort `json:"pull_requests" binding:"required"`
}