package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/mikhailvzhzhv/crack-hash/manager/internal/handler"
)

const (
	statisticGroupPath = "/api/statistic"
)

func NewStatisticRouter(router *gin.Engine, statisticHandler *handler.StatisticHandler) {
	hashApi := router.Group(statisticGroupPath)

	hashApi.GET("", statisticHandler.Statistic)
}
