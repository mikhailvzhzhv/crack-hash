package repository

import (
	"errors"
	"sync"

	"github.com/mikhailvzhzhv/crack-hash/manager/internal/handler/models"
)

type StatisticRepository struct {
	mu        sync.Mutex
	statistic *models.Statistic
}

func NewStatisticRepository() *StatisticRepository {
	return &StatisticRepository{
		mu:        sync.Mutex{},
		statistic: &models.Statistic{},
	}
}

func (s *StatisticRepository) Get() *models.Statistic {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.statistic
}

func (s *StatisticRepository) Update(modifier func(*models.Statistic) *models.Statistic) (*models.Statistic, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	updated := modifier(s.statistic)
	if updated == nil {
		return nil, errors.New("modifier returned nil")
	}

	return copyStatistic(updated), nil
}

func copyStatistic(s *models.Statistic) *models.Statistic {
	if s == nil {
		return nil
	}

	return &models.Statistic{
		TotalTasks:       s.TotalTasks,
		ActiveTasks:      s.ActiveTasks,
		CompletedTasks:   s.CompletedTasks,
		AvgExecutionTime: s.AvgExecutionTime,
		TasksParts:       s.TasksParts,
	}
}
