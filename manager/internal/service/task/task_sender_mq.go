package task

import (
	"github.com/mikhailvzhzhv/crack-hash/shared/v2/models"
	shared "github.com/mikhailvzhzhv/crack-hash/shared/v2/util"
	amqp "github.com/rabbitmq/amqp091-go"
)

type TaskSenderMq struct {
	channel *amqp.Channel
}

func NewTaskSenderMq() *TaskSenderMq {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	shared.FailOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	shared.FailOnError(err, "Failed to open a channel")

	_, err = ch.QueueDeclare(
		"task", // name
		false,  // durable
		false,  // delete when unused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
	)
	shared.FailOnError(err, "Failed to declare a queue")

	return &TaskSenderMq{
		channel: ch,
	}
}

func (t *TaskSenderMq) SendTasks(tasks []*models.Task) {
	for _, task := range tasks {
		t.channel.Publish(
			"",
			"task",
			false,
			false,
			amqp.Publishing{Body: shared.StructToJSON(task)},
		)
	}
}
