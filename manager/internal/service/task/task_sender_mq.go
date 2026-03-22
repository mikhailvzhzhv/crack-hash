package task

import (
	"log"

	"github.com/google/uuid"
	"github.com/mikhailvzhzhv/crack-hash/manager/internal/handler/models"
	shared "github.com/mikhailvzhzhv/crack-hash/shared/v2/models"
	shared_util "github.com/mikhailvzhzhv/crack-hash/shared/v2/util"
	amqp "github.com/rabbitmq/amqp091-go"
)

type TaskSenderMq struct {
	channel *amqp.Channel
}

func NewTaskSenderMq() *TaskSenderMq {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	shared_util.FailOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	shared_util.FailOnError(err, "Failed to open a channel")

	_, err = ch.QueueDeclare(
		"task", // name
		false,  // durable
		false,  // delete when unused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
	)
	shared_util.FailOnError(err, "Failed to declare a queue")

	return &TaskSenderMq{
		channel: ch,
	}
}

func (t *TaskSenderMq) SendTasks(tasks []*shared.Task, cancelChan chan models.Status) {
Loop:
	for _, task := range tasks {
		select {
		case status := <-cancelChan:
			if status == models.StatusCancelled {
				break Loop
			}
		default:
		}

		log.Printf("send task: %s", task)

		t.channel.Publish(
			"",
			"task",
			false,
			false,
			amqp.Publishing{Body: shared_util.StructToJSON(task)},
		)
	}
}

func (t *TaskSenderMq) SendCancel(requestID uuid.UUID, workerCount int) {
	panic("Unimplemented")
}
