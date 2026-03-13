package task

import (
	"github.com/google/uuid"
	"github.com/mikhailvzhzhv/crack-hash/manager/internal/handler/models"
	"github.com/mikhailvzhzhv/crack-hash/manager/internal/service/dto"
	shared "github.com/mikhailvzhzhv/crack-hash/shared/v2/models"
)

type TaskSender interface {
	SendTasks(tasks []*shared.Task, cancelChan chan models.Status)
	SendCancel(requestID uuid.UUID, workerCount int)
}

type Repository interface {
	Get(requestID uuid.UUID) *dto.Result
	Save(id uuid.UUID, modifier func(*dto.Result) *dto.Result) (*dto.Result, error)
	SaveHashRequest(hashRequest [16]byte, additionalCtx *dto.ResultAdditionalContext) (bool, uuid.UUID)
	GetAddtitonalContext(hashRequest [16]byte) *dto.ResultAdditionalContext
	DeleteHashRequest(hashRequest [16]byte)
}
