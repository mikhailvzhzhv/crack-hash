package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikhailvzhzhv/crack-hash/manager/internal/handler/models"
	"github.com/mikhailvzhzhv/crack-hash/manager/internal/service/result"
	shared "github.com/mikhailvzhzhv/crack-hash/shared/v2/models"
)

type ResultHandler struct {
	rp *result.ResultProcessor
}

func NewResultHandler(rp *result.ResultProcessor) *ResultHandler {
	return &ResultHandler{
		rp: rp,
	}
}

func (r *ResultHandler) HandleResult(c *gin.Context) {
	var result shared.Result
	if err := c.ShouldBindJSON(&result); err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(err))
		return
	}

	log.Printf("Received a message: %s", result)
	r.rp.Save(&result)

	log.Printf("Done")
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
