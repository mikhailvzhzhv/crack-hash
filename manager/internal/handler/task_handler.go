package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mikhailvzhzhv/crack-hash/manager/internal/handler/models"
	"github.com/mikhailvzhzhv/crack-hash/manager/internal/service/result"
	"github.com/mikhailvzhzhv/crack-hash/manager/internal/service/task"
)

type TaskHandler struct {
	crackHashService *task.CrackHashService
	resultProcessor  *result.ResultProcessor
}

func NewTaskHandler(crackHashService *task.CrackHashService, resultProcessor *result.ResultProcessor) *TaskHandler {
	return &TaskHandler{
		crackHashService: crackHashService,
		resultProcessor:  resultProcessor,
	}
}

func (h *TaskHandler) CrackHash(ctx *gin.Context) {
	var request models.CrackHashRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, models.NewErrorResponse(err))
		return
	}

	response, err := h.crackHashService.CrackHash(&request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.NewErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (h *TaskHandler) Hash(ctx *gin.Context) {
	var request models.HashRequest
	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, models.NewErrorResponse(err))
		return
	}

	hash := h.crackHashService.Hash(&request)

	ctx.JSON(http.StatusOK, &models.HashResponse{Hash: hash})
}

func (h *TaskHandler) Status(ctx *gin.Context) {
	var requestID = ctx.Query("requestId")

	id, err := uuid.Parse(requestID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.NewErrorResponse(err))
		return
	}

	result := h.resultProcessor.GetResult(id)

	ctx.JSON(http.StatusOK, result)
}

func (h *TaskHandler) CancelRequest(ctx *gin.Context) {
	var requestID = ctx.Query("requestId")

	id, err := uuid.Parse(requestID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.NewErrorResponse(err))
		return
	}

	err = h.crackHashService.CancelRequest(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.NewErrorResponse(err))
	}

	ctx.Status(http.StatusOK)
}
