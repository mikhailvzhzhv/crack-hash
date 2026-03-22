package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mikhailvzhzhv/crack-hash/worker/internal/handler"
	"github.com/mikhailvzhzhv/crack-hash/worker/internal/repository"
	"github.com/mikhailvzhzhv/crack-hash/worker/internal/routers"
	"github.com/mikhailvzhzhv/crack-hash/worker/internal/service"
)

func main() {
	repo := repository.NewSimpleRepository()
	resultSender := service.NewResultSenderMq()
	taskProcessor, err := service.NewTaskProcessor(resultSender, repo)
	if err != nil {
		log.Fatalf("error while create TaskProcessor: %v", err)
	}

	taskReceiver := service.NewTaskReceiver(taskProcessor, resultSender)
	taskReceiver.Receive()

	taskHandler := handler.NewTaskHandler(taskProcessor)

	router := gin.Default()

	routers.NewRouter(router, taskHandler)

	router.Run(":9000")
}
