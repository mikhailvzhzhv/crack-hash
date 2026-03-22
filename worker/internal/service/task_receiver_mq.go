package service

import (
	"log"

	shared "github.com/mikhailvzhzhv/crack-hash/shared/v2/util"
	amqp "github.com/rabbitmq/amqp091-go"
)

type TaskReceiverMq struct {
	msgs          <-chan amqp.Delivery
	taskProcessor *TaskProcessor
	resultSender  ResultSender
}

func NewTaskReceiver(taskProcessor *TaskProcessor, resultSender ResultSender) *TaskReceiverMq {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	shared.FailOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	shared.FailOnError(err, "Failed to open a channel")

	q, err := ch.QueueDeclare(
		"task", // name
		false,  // durable
		false,  // delete when unused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
	)
	shared.FailOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	shared.FailOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	shared.FailOnError(err, "Failed to register a consumer")

	return &TaskReceiverMq{
		msgs:          msgs,
		taskProcessor: taskProcessor,
		resultSender:  resultSender,
	}
}

func (t *TaskReceiverMq) Receive() {
	go func() {
		for d := range t.msgs {
			log.Printf("Received a message: %s", d.Body)

			task := shared.JSONToTask(d.Body)

			result := t.taskProcessor.ProcessTask(task)
			log.Printf("result: %s", result)
			t.resultSender.SendResult(result)

			log.Printf("Done")
			d.Ack(false)
		}
	}()
}
