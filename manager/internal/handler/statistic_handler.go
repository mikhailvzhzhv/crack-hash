package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikhailvzhzhv/crack-hash/manager/internal/service/statistic"
)

type StatisticHandler struct {
	statisticService *statistic.StatisticService
}

func NewStatisticHandler(statisticService *statistic.StatisticService) *StatisticHandler {
	return &StatisticHandler{
		statisticService: statisticService,
	}
}

func (s *StatisticHandler) Statistic(ctx *gin.Context) {
	statistic := s.statisticService.GetStatistic()

	ctx.JSON(http.StatusOK, statistic)
}
