package service

import (
    "math/rand"
    "review-rotator/internal/models"
    "review-rotator/internal/repository"
)

type PullRequestService struct {
    repo *repository.Repository
}

func NewPullRequestService(repo *repository.Repository) *PullRequestService {
    return &PullRequestService{repo: repo}
}

func (s *PullRequestService) CreatePullRequest(input models.CreatePRRequest) (*models.PullRequest, error) {
    if err := s.validatePRCreation(input); err != nil {
        return nil, err
    }
    teamMembers, err := s.getAvailableReviewers(input.AuthorID)
    if err != nil {
        return nil, err
    }
    reviewers := selectRandomReviewers(teamMembers, 2)
    pr := models.PullRequest{
        PullRequestID:   input.PullRequestID,
        PullRequestName: input.PullRequestName,
        AuthorID:        input.AuthorID,
        Status:          "OPEN",
        AssignedReviewers: reviewers,
    }
    return s.repo.PullRequest.CreatePullRequest(pr, reviewers)
}

func (s *PullRequestService) MergePullRequest(pullRequestID string) (*models.PullRequest, error) {
    pr, err := s.repo.PullRequest.MergePullRequest(pullRequestID)
    if err != nil {
        return nil, err
    }
    return pr, nil
}

func (s *PullRequestService) validatePRCreation(input models.CreatePRRequest) error {
    exists, err := s.repo.PullRequest.PRExists(input.PullRequestID)
    if err != nil {
        return err
    }
    if exists {
        return models.ErrPRExists
    }
    if err := s.repo.User.GetUserByID(input.AuthorID); err != nil {
        return models.ErrNotFound
    }

    return nil
}

func (s *PullRequestService) getAvailableReviewers(authorID string) ([]models.User, error) {
    team, err := s.repo.Team.GetTeamByUserID(authorID)
    if err != nil {
        return nil, models.ErrNotFound
    }

    return s.repo.User.GetActiveTeamMembers(team.TeamID, authorID)
}

func selectRandomReviewers(users []models.User, maxReviewers int) []string {
    if len(users) == 0 {
        return []string{}
    }
    count := len(users)
    if count > maxReviewers {
        count = maxReviewers
    }
    availableUsers := make([]models.User, len(users))
    copy(availableUsers, users)
    
    reviewers := make([]string, count)
    for i := 0; i < count; i++ {
        randomIndex := rand.Intn(len(availableUsers))
        reviewers[i] = availableUsers[randomIndex].UserID
        availableUsers = append(availableUsers[:randomIndex], availableUsers[randomIndex+1:]...)
    }
    return reviewers
}