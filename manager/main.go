package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mikhailvzhzhv/crack-hash/manager/internal/handler"
	"github.com/mikhailvzhzhv/crack-hash/manager/internal/repository"
	"github.com/mikhailvzhzhv/crack-hash/manager/internal/routers"
	"github.com/mikhailvzhzhv/crack-hash/manager/internal/service/result"
	"github.com/mikhailvzhzhv/crack-hash/manager/internal/service/statistic"
	"github.com/mikhailvzhzhv/crack-hash/manager/internal/service/task"
)

func main() {
	repo := repository.NewSimpleRepository()
	statisticRepo := repository.NewStatisticRepository()

	resultProcessor := result.NewResultProcessor(repo, statisticRepo)

	taskSender := task.NewTaskSenderRest()
	statisticService := statistic.NewStatisticService(statisticRepo)
	crackHashService, err := task.NewCrackHashService(taskSender, repo, statisticRepo)
	if err != nil {
		log.Fatalf("error while create CrackHashService: %v", err)
	}

	taskHandler := handler.NewTaskHandler(crackHashService, resultProcessor)
	resultHandler := handler.NewResultHandler(resultProcessor)
	statisticHandler := handler.NewStatisticHandler(statisticService)

	router := gin.Default()
	routers.NewHashRouter(router, taskHandler)
	routers.NewInternalRouter(router, resultHandler)
	routers.NewStatisticRouter(router, statisticHandler)

	router.Run(":8080")
}
