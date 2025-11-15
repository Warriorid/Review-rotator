package service

import (
	"review-rotator/internal/models"
	"review-rotator/internal/repository"
)

type Service struct {
	Team
	User
	PullRequest
}

type Team interface {
	AddTeam(team models.Team) (*models.Team, error)
	GetTeam(teamName string) (*models.Team, error)
}

type User interface {
	SetUserActive(userID string, isActive bool) (*models.User, error)
}
type PullRequest interface {
    CreatePullRequest(input models.CreatePRRequest) (*models.PullRequest, error)
	MergePullRequest(pullRequestID string) (*models.PullRequest, error)
	ReassignReviewer(input models.ReassignReviewerRequest) (*models.ReassignReviewerResponse, error)
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Team: NewTeamService(repos.Team),
		User: NewUserService(repos.User),
		PullRequest: NewPullRequestService(repos),
	}
}
