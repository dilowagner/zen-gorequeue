package producer

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

// sourceExchange const
const (
	sourceExchange = "amq.direct"
)

// RabbitmqWriter struct
type RabbitmqWriter struct {
	channel *amqp.Channel
	binding string
}

// NewWriter return RabbitmqWriter struct
func NewWriter(host, user, pass, port string) *RabbitmqWriter {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s", user, pass, host, port)
	conn, err := amqp.Dial(url)
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	return &RabbitmqWriter{
		channel: ch,
	}
}

// CreateReceiverQueue - Create Queue Receiver queue send to Queue Resender
func (w *RabbitmqWriter) CreateReceiverQueue(queue, exchange, bindingKey string) {

	w.binding = bindingKey

	args := make(map[string]interface{})
	args["x-message-ttl"] = 60000
	args["x-dead-letter-exchange"] = exchange
	args["x-dead-letter-routing-key"] = bindingKey

	w.create(queue, args)
	w.bind(queue, sourceExchange, bindingKey)
}

func (w *RabbitmqWriter) create(queue string, args map[string]interface{}) {
	q, err := w.channel.QueueDeclare(
		queue, // name of the queue
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // noWait
		args,  // arguments
	)
	if err != nil {
		log.Printf("Queue Declare: %s", err)
		panic(err)
	}

	log.Printf("Declared Queue (%q %d messages, %d consumers)",
		q.Name, q.Messages, q.Consumers)
}

func (w *RabbitmqWriter) bind(queue, exchange, bindingKey string) {
	if err := w.channel.QueueBind(
		queue,      // name of the queue
		bindingKey, // bindingKey
		exchange,   // sourceExchange
		false,      // noWait
		nil,        // arguments
	); err != nil {
		log.Printf("Queue Bind: %s", err)
		panic(err)
	}
}

// Publish message to Rabbitmq
// body []byte - Body Message
// exp string - expiration time
func (w *RabbitmqWriter) Publish(body []byte, exp string) error {
	message := amqp.Publishing{
		Body:       body,
		Expiration: exp,
	}

	log.Printf("BINDING_KEY: %s", w.binding)
	return w.channel.Publish(sourceExchange, w.binding, false, false, message)
}
