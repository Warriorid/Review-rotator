package service

import (
	"review-rotator/internal/models"
	"review-rotator/internal/repository"
)

type UserService struct {
	repo repository.User
}

func NewUserService(repo repository.User) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) SetUserActive(userID string, isActive bool) (*models.User, error) {
	user, err := s.repo.SetUserActive(userID, isActive)
	if err != nil {
		return nil, err
	}
	return user, nil
}