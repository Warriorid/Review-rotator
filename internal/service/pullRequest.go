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
    
	exists, err := s.repo.PullRequest.PRExists(input.PullRequestID)
    if err != nil {
        return nil, err
    }
    if exists {
        return nil, models.ErrPRExists	//409
    }

    err = s.repo.User.GetUserByID(input.AuthorID)
    if err != nil {
        return nil, models.ErrNotFound	//404
    }

    team, err := s.repo.Team.GetTeamByUserID(input.AuthorID)
    if err != nil {
        return nil, models.ErrNotFound	//404
    }

    teamMembers, err := s.repo.User.GetActiveTeamMembers(team.TeamID, input.AuthorID)
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
    createdPR, err := s.repo.PullRequest.CreatePullRequest(pr, reviewers)
    if err != nil {
        return nil, err
    }
    return createdPR, nil
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