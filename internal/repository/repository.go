package repository

import (
	"review-rotator/internal/models"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	Team
	User
	PullRequest
	Statistics
}

type Team interface {
	AddTeam(team models.Team) (*models.Team, error)
	GetTeam(teamName string) (*models.Team, error)
	GetTeamByUserID(userID string) (*models.Team, error)
}
type User interface {
	SetUserActive(userID string, isActive bool) (*models.User, error)
	GetActiveTeamMembers(teamID int, excludeUserID string) ([]models.User, error)
    GetUserByID(userID string) error
	GetUserTeamID(userID string) (int, error)
	GetUserReviewRequests(userID string) ([]models.PullRequestShort, error)
}

type PullRequest interface {
    CreatePullRequest(pr models.PullRequest, reviewers []string) (*models.PullRequest, error)
	PRExists(pullRequestID string) (bool, error)
	GetPullRequestByID(pullRequestID string) (*models.PullRequest, error)
    MergePullRequest(pullRequestID string) (*models.PullRequest, error)
	ReassignReviewer(pullRequestID string, oldReviewerID string, newReviewerID string) (*models.PullRequest, error)
}

type Statistics interface{
	GetReviewStats() (*models.StatsResponse, error)
}


func NewRepository(db *sqlx.DB) *Repository {
    return &Repository{
		Team: NewTeamPostgres(db),
		User: NewUserPostgres(db),
		PullRequest: NewPullRequestPostgres(db),
		Statistics: NewStatisticsPostgres(db),
    }
}