package repository

import (
	"review-rotator/internal/models"

	"github.com/jmoiron/sqlx"
)

type StatisticsPostgres struct {
	db *sqlx.DB
}

func NewStatisticsPostgres (db *sqlx.DB) *StatisticsPostgres {
	return &StatisticsPostgres{db: db}
}

func (r *StatisticsPostgres) GetReviewStats() (*models.StatsResponse, error) {
    userStats, err := r.getUserReviewStats()
    if err != nil {
        return nil, err
    }
    totalPRs, err := r.getTotalPRs()
    if err != nil {
        return nil, err
    }
    return &models.StatsResponse{
        UserStats: userStats,
        TotalPRs:  totalPRs,
    }, nil
}

func (r *StatisticsPostgres) getUserReviewStats() ([]models.UserStat, error) {
    query := `
        SELECT u.user_id, u.username, COUNT(prr.pull_request_id) as review_count
        FROM users u
        LEFT JOIN pr_reviewers prr ON u.user_id = prr.reviewer_id
        WHERE u.is_active = true
        GROUP BY u.user_id, u.username
        ORDER BY review_count DESC, u.username
    `
    
    rows, err := r.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var userStats []models.UserStat
    for rows.Next() {
        var stat models.UserStat
        err := rows.Scan(
            &stat.UserID,
            &stat.Username,
            &stat.ReviewCount,
        )
        if err != nil {
            return nil, err
        }
        userStats = append(userStats, stat)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return userStats, nil
}

func (r *StatisticsPostgres) getTotalPRs() (int, error) {
    var total int
    query := `SELECT COUNT(*) FROM pull_requests`
    err := r.db.QueryRow(query).Scan(&total)
    if err != nil {
        return 0, err
    }
    return total, nil
}