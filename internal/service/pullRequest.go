package service

import (
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
func (s *PullRequestService) ReassignReviewer(input models.ReassignReviewerRequest) (*models.ReassignReviewerResponse, error) {
    pr, err := s.validatePRForReassignment(input.PullRequestID)
    if err != nil {
        return nil, err
    }
    if err := s.validateReviewerAssignment(pr, input.OldUserID); err != nil {
        return nil, err
    }
    newReviewer, err := s.findReplacementCandidate(pr, input.OldUserID)
    if err != nil {
        return nil, err
    }
    updatedPR, err := s.repo.PullRequest.ReassignReviewer(input.PullRequestID, input.OldUserID, newReviewer)
    if err != nil {
        return nil, err
    }
    return &models.ReassignReviewerResponse{
        PR:         updatedPR,
        ReplacedBy: newReviewer,
    }, nil
}