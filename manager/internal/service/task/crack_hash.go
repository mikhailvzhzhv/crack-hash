package task

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/mikhailvzhzhv/crack-hash/manager/internal/handler/models"
	"github.com/mikhailvzhzhv/crack-hash/manager/internal/repository"
	"github.com/mikhailvzhzhv/crack-hash/manager/internal/service/dto"
	shared "github.com/mikhailvzhzhv/crack-hash/shared/v2/models"
)

type CrackHashService struct {
	taskSender    TaskSender
	batchSize     int
	repository    Repository
	workerCount   int
	statisticRepo *repository.StatisticRepository
}

func NewCrackHashService(taskSender TaskSender, repository Repository, statisticRepo *repository.StatisticRepository) (*CrackHashService, error) {
	batchSizeStr := os.Getenv("BATCH_SIZE")
	if len(batchSizeStr) == 0 {
		return nil, errors.New("env variable BATCH_SIZE is absent")
	}

	batchSize, err := strconv.Atoi(batchSizeStr)
	if err != nil {
		return nil, err
	}

	workerCountStr := os.Getenv("WORKER_COUNT")
	if len(workerCountStr) == 0 {
		return nil, errors.New("env variable WORKER_COUNT is absent")
	}

	workerCount, err := strconv.Atoi(workerCountStr)
	if err != nil {
		return nil, err
	}

	return &CrackHashService{
		taskSender:    taskSender,
		batchSize:     batchSize,
		repository:    repository,
		workerCount:   workerCount,
		statisticRepo: statisticRepo,
	}, nil
}

func (c *CrackHashService) Hash(request *models.HashRequest) string {
	hash := md5.Sum([]byte(request.Word))

	return hex.EncodeToString(hash[:])
}

func (c *CrackHashService) CrackHash(request *models.CrackHashRequest) (*models.CrackHashResponse, error) {
	ok, requestID, cancelChan := c.saveTaskRequest(request)
	if !ok {
		estimatedCombinations := c.calculateEstimatedCombinations(request)

		return &models.CrackHashResponse{
			RequestId:             requestID,
			EstimatedCombinations: estimatedCombinations,
		}, nil
	}

	tasks, estimatedCombinations, err := c.createTasks(requestID, request)
	if err != nil {
		return nil, err
	}

	go c.taskSender.SendTasks(tasks, cancelChan)

	c.updateStatistic()

	return &models.CrackHashResponse{
		RequestId:             requestID,
		EstimatedCombinations: estimatedCombinations,
	}, nil
}

func (c *CrackHashService) createTasks(requestID uuid.UUID, request *models.CrackHashRequest) ([]*shared.Task, int, error) {
	tasks := make([]*shared.Task, 0)

	partCount, estimatedCombinations := c.calculateBatchCount(request)
	for i := range partCount {
		task, err := c.createTask(requestID, request, i, partCount)
		if err != nil {
			return nil, 0, err
		}

		tasks = append(tasks, task)
	}

	return tasks, estimatedCombinations, nil
}

func (c *CrackHashService) createTask(requestID uuid.UUID, request *models.CrackHashRequest, partNumber int, partCount int) (*shared.Task, error) {
	hash, err := hex.DecodeString(request.Hash)
	if err != nil {
		return nil, err
	}

	return &shared.Task{
		RequestID:  requestID,
		PartNumber: partNumber,
		PartCount:  partCount,
		Alphabet:   request.Alphabet,
		MaxWordLen: request.MaxLength,
		TargetHash: [16]byte(hash),
		Algorithm:  shared.Algorithm(request.Algorithm),
	}, nil
}

func (c *CrackHashService) calculateBatchCount(request *models.CrackHashRequest) (int, int) {
	estimatedCombinations := float64(c.calculateEstimatedCombinations(request))

	return int(math.Ceil(estimatedCombinations / float64(c.batchSize))), int(estimatedCombinations)
}

func (c *CrackHashService) calculateEstimatedCombinations(request *models.CrackHashRequest) int {
	estimatedCombinations := float64(0)
	for i := 1; i <= request.MaxLength; i++ {
		estimatedCombinations += math.Pow(float64(len(request.Alphabet)), float64(i))
	}

	return int(estimatedCombinations)
}

func (c *CrackHashService) saveTaskRequest(request *models.CrackHashRequest) (bool, uuid.UUID, chan models.Status) {
	hashRequest := c.getHashRequet(request)

	requestID := uuid.New()
	cancelChan := make(chan models.Status)
	additionalCtx := &dto.ResultAdditionalContext{
		RequestID:  requestID,
		CancelChan: cancelChan,
	}

	if ok, savedRequestID := c.repository.SaveHashRequest(hashRequest, additionalCtx); !ok {
		return false, savedRequestID, nil
	}

	partCount, _ := c.calculateBatchCount(request)
	log.Printf("partCount: %d", partCount)

	c.repository.Save(requestID, func(r *dto.Result) *dto.Result {
		return dto.NewResult(requestID, partCount, request)
	})

	return true, requestID, cancelChan
}

func (c *CrackHashService) getHashRequet(request *models.CrackHashRequest) [16]byte {
	stringData := strings.Join([]string{request.Alphabet, request.Hash, strconv.Itoa(request.MaxLength)}, "_")
	byteData := []byte(stringData)

	return md5.Sum(byteData)
}

func (c *CrackHashService) CancelRequest(requestID uuid.UUID) error {
	result := c.repository.Get(requestID)
	if result == nil {
		return fmt.Errorf("cannot find request with id: %s", requestID)
	}

	if result.Status != models.StatusInProgress {
		return fmt.Errorf("task in inactive status: %s", result)

	}

	hashRequest := c.getHashRequet(result.Request)

	c.repository.Save(requestID, func(r *dto.Result) *dto.Result {
		r.Status = models.StatusCancelled

		return r
	})

	addtitonalCtx := c.repository.GetAddtitonalContext(hashRequest)
	if addtitonalCtx == nil {
		return fmt.Errorf("cannot find addtitonalCtx with hashRequest: %s", hashRequest)
	}

	addtitonalCtx.CancelChan <- models.StatusCancelled

	c.repository.DeleteHashRequest(hashRequest)

	c.statisticRepo.Update(func(s *models.Statistic) *models.Statistic {
		s.ActiveTasks--

		return s
	})

	go c.taskSender.SendCancel(requestID, c.workerCount)

	return nil
}

func (c *CrackHashService) updateStatistic() {
	c.statisticRepo.Update(func(s *models.Statistic) *models.Statistic {
		s.TotalTasks++
		s.ActiveTasks++

		return s
	})
}
