package service

import (
	"math/rand"
	"review-rotator/internal/models"
)

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

func selectRandomReviewer(users []models.User) string {
    if len(users) == 0 {
        return ""
    }
    randomIndex := rand.Intn(len(users))
    return users[randomIndex].UserID
}

func (s *PullRequestService) validatePRForReassignment(pullRequestID string) (*models.PullRequest, error) {
    pr, err := s.repo.PullRequest.GetPullRequestByID(pullRequestID)
    if err != nil {
        return nil, err
    }
    
    if pr.Status == "MERGED" {
        return nil, models.ErrPRMerged
    }
    
    return pr, nil
}

func (s *PullRequestService) validateReviewerAssignment(pr *models.PullRequest, oldUserID string) error {
    for _, reviewer := range pr.AssignedReviewers {
        if reviewer == oldUserID {
            return nil
        }
    }
    return models.ErrNotAssigned
}

func (s *PullRequestService) findReplacementCandidate(pr *models.PullRequest, oldUserID string) (string, error) {
    oldReviewerTeamID, err := s.repo.User.GetUserTeamID(oldUserID)
    if err != nil {
        return "", models.ErrNotFound
    }
    availableUsers, err := s.repo.User.GetActiveTeamMembers(oldReviewerTeamID, oldUserID)
    if err != nil {
        return "", err
    }
    filteredUsers := s.filterReplacementCandidates(availableUsers, pr, oldUserID)
    if len(filteredUsers) == 0 {
        return "", models.ErrNoCandidate
    }
    return selectRandomReviewer(filteredUsers), nil
}

func (s *PullRequestService) filterReplacementCandidates(availableUsers []models.User, pr *models.PullRequest, oldUserID string) []models.User {
    filteredUsers := []models.User{}
    for _, user := range availableUsers {
        if s.isValidReplacementCandidate(user, pr, oldUserID) {
            filteredUsers = append(filteredUsers, user)
        }
    }
    return filteredUsers
}

func (s *PullRequestService) isValidReplacementCandidate(user models.User, pr *models.PullRequest, oldUserID string) bool {
    if user.UserID == pr.AuthorID {
        return false
    }
    for _, reviewer := range pr.AssignedReviewers {
        if reviewer == user.UserID && reviewer != oldUserID {
            return false
        }
    }
    
    return true
}