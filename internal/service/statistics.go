package service

import (
	"review-rotator/internal/models"
	"review-rotator/internal/repository"
)

type StatisticsService struct {
	repo repository.Statistics
}

func NewStatisticsService (repo repository.Statistics) *StatisticsService {
	return &StatisticsService{repo: repo}
}

func (s *StatisticsService) GetReviewStats() (*models.StatsResponse, error) {
    stats, err := s.repo.GetReviewStats()
    if err != nil {
        return nil, err
    }
    return stats, nil
}