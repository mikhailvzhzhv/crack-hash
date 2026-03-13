package task

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mikhailvzhzhv/crack-hash/manager/internal/handler/models"
	shared "github.com/mikhailvzhzhv/crack-hash/shared/v2/models"
	sharedUtil "github.com/mikhailvzhzhv/crack-hash/shared/v2/util"
)

type TaskSenderRest struct {
	client   *http.Client
	maxRetry int
}

func NewTaskSenderRest() *TaskSenderRest {
	return &TaskSenderRest{
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
		maxRetry: 3,
	}
}

func (t *TaskSenderRest) SendTasks(tasks []*shared.Task, cancelChan chan models.Status) {
Loop:
	for _, task := range tasks {
		select {
		case status := <-cancelChan:
			if status == models.StatusCancelled {
				break Loop
			}
		default:
		}

		time.Sleep(1 * time.Second)

		var err error
		for i := 0; i < t.maxRetry; i++ {
			err = t.SendTask(task)
			if err == nil {
				break
			}

			log.Println("Send task error:", err)

			time.Sleep(1 * time.Second)
		}

		if err != nil {
			log.Printf("Failed to send task %s after %d retries: %v", task.RequestID, t.maxRetry, err)
		}

	}
}

func (t *TaskSenderRest) SendTask(task *shared.Task) error {
	taskJSON := sharedUtil.StructToJSON(task)
	if taskJSON == nil {
		return fmt.Errorf("Failed to marshal task %+v", task)
	}

	resp, err := t.client.Post(
		"http://nginx/internal/api/worker/hash/crack/task",
		"application/json",
		bytes.NewBuffer(taskJSON),
	)

	if err != nil {
		return fmt.Errorf("Failed to send task %+v: %v", task, err)
	}

	resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("Received non-success status %d for task %+v", resp.StatusCode, task)
	}

	return nil
}

func (t *TaskSenderRest) SendCancel(requestID uuid.UUID, workerCount int) {
	for i := 0; i < workerCount; i++ {
		url := fmt.Sprintf(
			"http://nginx/internal/api/worker/hash/crack/task?requestId=%s",
			requestID.String(),
		)
		req, err := http.NewRequest(http.MethodDelete, url, nil)

		if err != nil {
			log.Printf("Cannot create request: %v", err)
			return
		}

		resp, err := t.client.Do(req)
		if err != nil {
			log.Printf("Error while send request: %v", err)
			return
		}

		resp.Body.Close()

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			log.Printf("Received non-success status %d for result %+v", resp.StatusCode, resp)
		}
	}
}
