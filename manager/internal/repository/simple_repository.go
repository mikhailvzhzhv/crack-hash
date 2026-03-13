package repository

import (
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/mikhailvzhzhv/crack-hash/manager/internal/service/dto"
)

type SimpleRepository struct {
	mu               sync.Mutex
	store            map[uuid.UUID]*dto.Result
	hashRequestStore map[[16]byte]*dto.ResultAdditionalContext
}

func NewSimpleRepository() *SimpleRepository {
	return &SimpleRepository{
		store:            make(map[uuid.UUID]*dto.Result),
		mu:               sync.Mutex{},
		hashRequestStore: make(map[[16]byte]*dto.ResultAdditionalContext),
	}
}

func (r *SimpleRepository) Get(requestID uuid.UUID) *dto.Result {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.store[requestID]
}

func (r *SimpleRepository) Save(id uuid.UUID, modifier func(*dto.Result) *dto.Result) (*dto.Result, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	result, exists := r.store[id]
	if !exists {
		saved := modifier(nil)
		r.store[id] = saved

		return copyResult(saved), nil
	}

	updatedResult := modifier(result)
	if updatedResult == nil {
		return nil, errors.New("modifier returned nil")
	}

	r.store[id] = updatedResult

	return copyResult(updatedResult), nil
}

func copyResult(r *dto.Result) *dto.Result {
	if r == nil {
		return nil
	}

	dataCopy := make([]string, len(r.Data))
	copy(dataCopy, r.Data)

	return &dto.Result{
		RequestID:      r.RequestID,
		Status:         r.Status,
		Data:           dataCopy,
		CompletedTasks: r.CompletedTasks,
		TaskCount:      r.TaskCount,
		Request:        r.Request,
	}
}

// true если успешно сохранен, false если такой hashRequest уже существует
func (r *SimpleRepository) SaveHashRequest(hashRequest [16]byte, additionalCtx *dto.ResultAdditionalContext) (bool, uuid.UUID) {
	r.mu.Lock()
	defer r.mu.Unlock()

	savedAdditionalCtx := r.hashRequestStore[hashRequest]
	if savedAdditionalCtx != nil {
		return false, savedAdditionalCtx.RequestID
	}

	r.hashRequestStore[hashRequest] = additionalCtx

	return true, uuid.Nil
}

func (r *SimpleRepository) GetAddtitonalContext(hashRequest [16]byte) *dto.ResultAdditionalContext {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.hashRequestStore[hashRequest]
}

func (r *SimpleRepository) DeleteHashRequest(hashRequest [16]byte) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.hashRequestStore[hashRequest] = nil
}
