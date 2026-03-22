package service

import (
	"bytes"
	"context"
	"crypto/md5"
	"errors"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	shared "github.com/mikhailvzhzhv/crack-hash/shared/v2/models"
	"github.com/mikhailvzhzhv/crack-hash/worker/internal/dto"
	"github.com/mikhailvzhzhv/crack-hash/worker/internal/repository"
	"github.com/mikhailvzhzhv/crack-hash/worker/internal/util"
)

type TaskProcessor struct {
	batchSize    int
	resultSender ResultSender
	repo         *repository.SimpleRepository
}

func NewTaskProcessor(resultSender ResultSender, repo *repository.SimpleRepository) (*TaskProcessor, error) {
	batchSizeStr := os.Getenv("BATCH_SIZE")
	if len(batchSizeStr) == 0 {
		return nil, errors.New("env variable BATCH_SIZE is absent")
	}

	batchSize, err := strconv.Atoi(batchSizeStr)
	if err != nil {
		return nil, err
	}

	return &TaskProcessor{
		batchSize:    batchSize,
		resultSender: resultSender,
		repo:         repo,
	}, nil
}

func (s *TaskProcessor) ProcessTask(task *shared.Task) *shared.Result {
	if task.Algorithm != shared.MD5 {
		panic("unknown algorithm: " + task.Algorithm)
	}

	ctx, cancel := s.saveTaskWithContext(task)
	defer cancel()

	wordCount := 0
	targetWords := make([]string, 0)
	generator := util.NewWordGenerator(task, s.batchSize)
	timer := NewTimer()

	timer.Start()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		word, hasNext := generator.Next()
		if !hasNext {
			break
		}

		worldHash := md5.Sum([]byte(word))
		if bytes.Equal(worldHash[:], task.TargetHash[:]) {
			targetWords = append(targetWords, word)
		}

		wordCount++
	}

	timer.Stop()
	log.Printf("wordCount: %d; seconds: %f", wordCount, timer.GetSeconds())
	avgExecutionTime := int(float64(wordCount) / timer.GetSeconds())

	return &shared.Result{
		RequestID:        task.RequestID,
		Words:            targetWords,
		PartNumber:       task.PartNumber,
		PartCount:        task.PartCount,
		AvgExecutionTime: avgExecutionTime,
	}
}

func (s *TaskProcessor) saveTaskWithContext(task *shared.Task) (context.Context, context.CancelFunc) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	taskWithCtx := &dto.TaskWithContext{
		RequestID:  task.RequestID,
		Ctx:        ctx,
		CancelFunc: cancelFunc,
	}

	s.repo.SaveTaskWithContext(taskWithCtx)

	return ctx, cancelFunc
}

func (s *TaskProcessor) CancelRequest(id uuid.UUID) {
	taskCtx := s.repo.GetTaskWithContext(id)
	taskCtx.CancelFunc()
}
