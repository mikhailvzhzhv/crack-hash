package statistic

import (
	"github.com/mikhailvzhzhv/crack-hash/manager/internal/handler/models"
	"github.com/mikhailvzhzhv/crack-hash/manager/internal/repository"
)

type StatisticService struct {
	statisticRepo *repository.StatisticRepository
}

func NewStatisticService(statisticRepo *repository.StatisticRepository) *StatisticService {
	return &StatisticService{
		statisticRepo: statisticRepo,
	}
}

func (s *StatisticService) GetStatistic() *models.Statistic {
	return s.statisticRepo.Get()
}
