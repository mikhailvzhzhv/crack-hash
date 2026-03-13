package dto

import (
	"github.com/google/uuid"
	"github.com/mikhailvzhzhv/crack-hash/manager/internal/handler/models"
)

type Result struct {
	RequestID      uuid.UUID
	Status         models.Status
	Data           []string
	CompletedTasks int
	TaskCount      int
	Request        *models.CrackHashRequest
}

func NewResult(requestID uuid.UUID, partCount int, request *models.CrackHashRequest) *Result {
	return &Result{
		RequestID:      requestID,
		Status:         models.StatusInProgress,
		Data:           make([]string, 0),
		CompletedTasks: 0,
		TaskCount:      partCount,
		Request:        request,
	}
}

type ResultAdditionalContext struct {
	RequestID  uuid.UUID
	CancelChan chan models.Status
}
