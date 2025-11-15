package models

type StatsResponse struct {
    UserStats []UserStat `json:"user_stats" binding:"required"`
    TotalPRs  int        `json:"total_prs" binding:"required"`
}

type UserStat struct {
    UserID    string `json:"user_id" binding:"required"`
    Username  string `json:"username" binding:"required"`
    ReviewCount int  `json:"review_count" binding:"required"`
}