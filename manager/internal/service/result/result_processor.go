package result

import (
	"github.com/google/uuid"
	"github.com/mikhailvzhzhv/crack-hash/manager/internal/handler/models"
	"github.com/mikhailvzhzhv/crack-hash/manager/internal/repository"
	"github.com/mikhailvzhzhv/crack-hash/manager/internal/service/dto"
	shared "github.com/mikhailvzhzhv/crack-hash/shared/v2/models"
)

type ResultProcessor struct {
	repo          *repository.SimpleRepository
	statisticRepo *repository.StatisticRepository
}

func NewResultProcessor(repo *repository.SimpleRepository, statisticRepo *repository.StatisticRepository) *ResultProcessor {
	return &ResultProcessor{
		repo:          repo,
		statisticRepo: statisticRepo,
	}
}

func (r *ResultProcessor) Save(resultResponse *shared.Result) {
	r.repo.Save(resultResponse.RequestID, func(r *dto.Result) *dto.Result {
		if r.Status != models.StatusInProgress {
			return r
		}

		r.Data = append(r.Data, resultResponse.Words...)
		r.CompletedTasks++

		if r.CompletedTasks == r.TaskCount {
			r.Status = models.StatusReady
		}

		return r
	})

	r.statisticRepo.Update(func(s *models.Statistic) *models.Statistic {
		result := r.repo.Get(resultResponse.RequestID)
		if result.Status == models.StatusReady {
			s.ActiveTasks--
			s.CompletedTasks++
		}

		if result.Status == models.StatusCancelled {
			s.ActiveTasks--
		}

		avg := (s.TasksParts*s.AvgExecutionTime + resultResponse.AvgExecutionTime) / (s.TasksParts + 1)
		s.AvgExecutionTime = avg
		s.TasksParts++

		return s
	})
}

func (r *ResultProcessor) GetResult(requestID uuid.UUID) *models.ResultResponse {
	result := r.repo.Get(requestID)

	return &models.ResultResponse{
		Status: result.Status,
		Data:   result.Data,
	}
}
