package service

import (
	"github.com/mikhailvzhzhv/crack-hash/shared/v2/models"
	shared "github.com/mikhailvzhzhv/crack-hash/shared/v2/util"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ResultSenderMq struct {
	channel *amqp.Channel
}

func (r *ResultSenderMq) SendResult(result *models.Result) {
	r.channel.Publish(
		"",
		"result",
		false,
		false,
		amqp.Publishing{Body: shared.StructToJSON(result)},
	)
}

func NewResultSenderMq() *ResultSenderMq {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	shared.FailOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	shared.FailOnError(err, "Failed to open a channel")

	_, err = ch.QueueDeclare(
		"result", // name
		false,    // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	shared.FailOnError(err, "Failed to declare a queue")

	return &ResultSenderMq{
		channel: ch,
	}
}
