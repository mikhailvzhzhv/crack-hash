package result

import (
	"log"

	shared "github.com/mikhailvzhzhv/crack-hash/shared/v2/util"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ResultAggregatorMq struct {
	msgs            <-chan amqp.Delivery
	resultProcessor *ResultProcessor
}

func NewResultAggregator(resultProcessor *ResultProcessor) *ResultAggregatorMq {

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	shared.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	shared.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"result_queue", // name
		true,           // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
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

	return &ResultAggregatorMq{
		msgs:            msgs,
		resultProcessor: resultProcessor,
	}
}

func (r *ResultAggregatorMq) Aggregate() {
	go func() {
		for d := range r.msgs {
			log.Printf("Received a message: %s", d.Body)

			r.resultProcessor.Save(shared.JSONToResult(d.Body))

			log.Printf("Done")
			d.Ack(false)
		}
	}()
}
