package models

type PullRequest struct {
    PullRequestID    string   `json:"pull_request_id" binding:"required"`
	PullRequestName  string   `json:"pull_request_name" binding:"required"`
	AuthorID         string   `json:"author_id" binding:"required"`
	Status           string   `json:"status" binding:"required"`
	AssignedReviewers []string `json:"assigned_reviewers" binding:"required"`
    CreatedAt        *string  `json:"createdAt,omitempty"`
    MergedAt         *string  `json:"mergedAt,omitempty"`
}

type PullRequestShort struct {
	PullRequestID   string `json:"pull_request_id" binding:"required"`
	PullRequestName string `json:"pull_request_name" binding:"required"`
	AuthorID        string `json:"author_id" binding:"required"`
	Status          string `json:"status" binding:"required"`
}

type CreatePRRequest struct {
    PullRequestID   string `json:"pull_request_id" binding:"required"`
    PullRequestName string `json:"pull_request_name" binding:"required"`
    AuthorID        string `json:"author_id" binding:"required"`
}

type MergePRRequest struct {
    PullRequestID string `json:"pull_request_id" binding:"required"`
}

type ReassignReviewerRequest struct {
    PullRequestID string `json:"pull_request_id" binding:"required"`
    OldUserID     string `json:"old_user_id" binding:"required"`
}

type ReassignReviewerResponse struct {
    PR         *PullRequest `json:"pr" binding:"required"`
    ReplacedBy string       `json:"replaced_by" binding:"required"`
}