package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mikhailvzhzhv/crack-hash/shared/v2/models"
	"github.com/mikhailvzhzhv/crack-hash/worker/internal/service"
)

type TaskHandler struct {
	taskProcessor *service.TaskProcessor
}

func NewTaskHandler(taskProcessor *service.TaskProcessor) *TaskHandler {
	return &TaskHandler{
		taskProcessor: taskProcessor,
	}
}

func (t *TaskHandler) HandleTask(ctx *gin.Context) {
	var task models.Task
	if err := ctx.ShouldBindJSON(&task); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	go t.taskProcessor.ProcessTask(&task)

	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (t *TaskHandler) CancelTask(ctx *gin.Context) {
	var requestID = ctx.Query("requestId")

	id, err := uuid.Parse(requestID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	t.taskProcessor.CancelRequest(id)

	ctx.Status(http.StatusOK)
}
