package repository

import (
	"database/sql"
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
    `
    now := time.Now().Format(time.RFC3339)
    _, err = tx.Exec(query, 
        pr.PullRequestID,
        pr.PullRequestName,
        pr.AuthorID,
        pr.Status,
        now,
    )
    if err != nil {
        return nil, err
    }
    reviewerQuery := `INSERT INTO pr_reviewers (pull_request_id, reviewer_id) VALUES ($1, $2)`
    for _, reviewerID := range reviewers {
        _, err = tx.Exec(reviewerQuery, pr.PullRequestID, reviewerID)
        if err != nil {
            return nil, err
        }
    }
    if err := tx.Commit(); err != nil {
        return nil, err
    }
    pr.CreatedAt = &now
    pr.AssignedReviewers = reviewers
    return &pr, nil
}

func (r *PullRequestPostgres) MergePullRequest(pullRequestID string) (*models.PullRequest, error) {
    pr, err := r.GetPullRequestByID(pullRequestID)
    if err != nil {
        return nil, err
    }
    if pr.Status == "MERGED" {
        return pr, nil
    }
    now := time.Now().Format(time.RFC3339)
    _, err = r.db.Exec(
        "UPDATE pull_requests SET status = 'MERGED', merged_at = $1 WHERE pull_request_id = $2",
        now, pullRequestID,
    )
    if err != nil {
        return nil, err
    }
    return r.GetPullRequestByID(pullRequestID)
}

func (r *PullRequestPostgres) GetPullRequestByID(pullRequestID string) (*models.PullRequest, error) {
    var pr models.PullRequest
    query := `
        SELECT pull_request_id, pull_request_name, author_id, status, created_at, merged_at
        FROM pull_requests 
        WHERE pull_request_id = $1
    `
    err := r.db.QueryRow(query, pullRequestID).Scan(
        &pr.PullRequestID,
        &pr.PullRequestName,
        &pr.AuthorID,
        &pr.Status,
        &pr.CreatedAt,
        &pr.MergedAt,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, models.ErrNotFound
        }
        return nil, err
    }
    reviewersQuery := `
        SELECT reviewer_id 
        FROM pr_reviewers 
        WHERE pull_request_id = $1
    `
    rows, err := r.db.Query(reviewersQuery, pullRequestID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var reviewerID string
        if err := rows.Scan(&reviewerID); err != nil {
            return nil, err
        }
        pr.AssignedReviewers = append(pr.AssignedReviewers, reviewerID)
    }

    return &pr, nil
}

func (r *PullRequestPostgres) ReassignReviewer(pullRequestID string, oldReviewerID string, newReviewerID string) (*models.PullRequest, error) {
    tx, err := r.db.Beginx()
    if err != nil {
        return nil, err
    }
    defer tx.Rollback()
    var exists bool
    err = tx.QueryRow(
        "SELECT EXISTS(SELECT 1 FROM pr_reviewers WHERE pull_request_id = $1 AND reviewer_id = $2)",
        pullRequestID, oldReviewerID,
    ).Scan(&exists)
    if err != nil {
        return nil, err
    }
    if !exists {
        return nil, models.ErrNotAssigned
    }
    _, err = tx.Exec(
        "DELETE FROM pr_reviewers WHERE pull_request_id = $1 AND reviewer_id = $2",
        pullRequestID, oldReviewerID,
    )
    if err != nil {
        return nil, err
    }
    _, err = tx.Exec(
        "INSERT INTO pr_reviewers (pull_request_id, reviewer_id) VALUES ($1, $2)",
        pullRequestID, newReviewerID,
    )
    if err != nil {
        return nil, err
    }

    if err := tx.Commit(); err != nil {
        return nil, err
    }
    return r.GetPullRequestByID(pullRequestID)
}