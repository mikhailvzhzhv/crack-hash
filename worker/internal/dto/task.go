package dto

import (
	"context"

	"github.com/google/uuid"
)

type TaskWithContext struct {
	RequestID  uuid.UUID
	Ctx        context.Context
	CancelFunc context.CancelFunc
}
