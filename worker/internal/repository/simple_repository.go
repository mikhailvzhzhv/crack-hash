package repository

import (
	"sync"

	"github.com/google/uuid"
	"github.com/mikhailvzhzhv/crack-hash/worker/internal/dto"
)

type SimpleRepository struct {
	mutex sync.Mutex
	store map[uuid.UUID]*dto.TaskWithContext
}

func NewSimpleRepository() *SimpleRepository {
	return &SimpleRepository{
		mutex: sync.Mutex{},
		store: make(map[uuid.UUID]*dto.TaskWithContext),
	}
}

func (s *SimpleRepository) SaveTaskWithContext(task *dto.TaskWithContext) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.store[task.RequestID] = task
}

func (s *SimpleRepository) GetTaskWithContext(requestID uuid.UUID) *dto.TaskWithContext {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.store[requestID]
}
