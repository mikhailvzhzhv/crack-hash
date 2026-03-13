package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/mikhailvzhzhv/crack-hash/worker/internal/handler"
)

func NewRouter(router *gin.Engine, taskHandler *handler.TaskHandler) {
	router.POST("/internal/api/worker/hash/crack/task", taskHandler.HandleTask)
	router.DELETE("/internal/api/worker/hash/crack/task", taskHandler.CancelTask)
}
