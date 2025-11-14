package repository

import (
	"review-rotator/internal/models"
	"time"

	"github.com/jmoiron/sqlx"
)

type PullRequestPostgres struct {
    db *sqlx.DB
}

func NewPullRequestPostgres(db *sqlx.DB) *PullRequestPostgres {
    return &PullRequestPostgres{db: db}
}

func (r *PullRequestPostgres) PRExists(pullRequestID string) (bool, error) {
    var exists bool
    query := `SELECT EXISTS(SELECT 1 FROM pull_requests WHERE pull_request_id = $1)`
    err := r.db.QueryRow(query, pullRequestID).Scan(&exists)
    return exists, err
}

func (r *PullRequestPostgres) CreatePullRequest(pr models.PullRequest, reviewers []string) (*models.PullRequest, error) {
    tx, err := r.db.Beginx()
    if err != nil {
        return nil, err
    }
    defer tx.Rollback()

    query := `
        INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id, status, created_at)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING pull_request_id, pull_request_name, author_id, status, created_at
    `
    
    now := time.Now().Format(time.RFC3339)
    err = tx.QueryRow(
        query,
        pr.PullRequestID,
        pr.PullRequestName,
        pr.AuthorID,
        pr.Status,
        now,
    ).Scan(
        &pr.PullRequestID,
        &pr.PullRequestName,
        &pr.AuthorID,
        &pr.Status,
        &pr.CreatedAt,
    )
    if err != nil {
        return nil, err
    }

    for _, reviewerID := range reviewers {
        _, err = tx.Exec(
            "INSERT INTO pr_reviewers (pull_request_id, reviewer_id) VALUES ($1, $2)",
            pr.PullRequestID,
            reviewerID,
        )
        if err != nil {
            return nil, err
        }
    }

    pr.AssignedReviewers = reviewers

    if err := tx.Commit(); err != nil {
        return nil, err
    }

    return &pr, nil
}