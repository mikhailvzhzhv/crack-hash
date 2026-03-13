package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/mikhailvzhzhv/crack-hash/manager/internal/handler"
)

const (
	hashGroupPath       = "/api/hash"
	hashGroupCrackPath  = "/crack"
	hashGroupStatusPath = "/status"
)

func NewHashRouter(router *gin.Engine, taskHandler *handler.TaskHandler) {
	hashApi := router.Group(hashGroupPath)

	hashApi.POST(hashGroupCrackPath, taskHandler.CrackHash)
	hashApi.GET("", taskHandler.Hash)
	hashApi.GET(hashGroupStatusPath, taskHandler.Status)
	hashApi.DELETE(hashGroupCrackPath, taskHandler.CancelRequest)
}

func NewInternalRouter(router *gin.Engine, resultHandler *handler.ResultHandler) {
	router.PATCH("/internal/api/manager/hash/crack/request", resultHandler.HandleResult)
}
