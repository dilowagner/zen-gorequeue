package consumer

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

// RabbitmqReader struct
type RabbitmqReader struct {
	conn    *amqp.Connection
	channel *amqp.Channel

	Done chan error
}

// NewReader - create Rabbitmq reader struct
func NewReader(host, user, pass, port string) *RabbitmqReader {

	c := &RabbitmqReader{
		conn:    nil,
		channel: nil,
		Done:    make(chan error),
	}
	var err error

	url := fmt.Sprintf("amqp://%s:%s@%s:%s", user, pass, host, port)
	c.conn, err = amqp.Dial(url)
	if err != nil {
		panic(err)
	}

	go func() {
		log.Printf("closing: %s", <-c.conn.NotifyClose(make(chan *amqp.Error)))
	}()

	c.channel, err = c.conn.Channel()
	if err != nil {
		panic(err)
	}

	return c
}

// Shutdown
func (r *RabbitmqReader) Shutdown(binding string) error {
	// will close() the deliveries channel
	if err := r.channel.Cancel(binding, true); err != nil {
		log.Printf("Consumer cancel failed: %s", err)
		return err
	}

	if err := r.conn.Close(); err != nil {
		log.Printf("AMQP connection close error: %s", err)
		return err
	}

	defer log.Printf("AMQP shutdown OK")

	// wait for handle() to exit
	return <-r.Done
}

// ConfigureExchange method - create exchange by name
func (r *RabbitmqReader) CreateExchange(exchange string) {
	log.Printf("Got Channel, declaring Exchange (%q)", exchange)
	if err := r.channel.ExchangeDeclare(
		exchange, // name of the exchange
		"direct", // type
		true,     // durable
		false,    // delete when complete
		false,    // internal
		false,    // noWait
		nil,      // arguments
	); err != nil {
		log.Printf("Exchange Declare: %s", err)
	}
}

// CreateResenderQueue - Create Resender queue to repost after 60s
func (r *RabbitmqReader) CreateResenderQueue(queue, exchange, bindingKey string) {

	r.create(queue, nil)
	r.bind(queue, exchange, bindingKey)
}

// create - QueueDeclare
func (r *RabbitmqReader) create(queue string, args map[string]interface{}) {
	q, err := r.channel.QueueDeclare(
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

// bind - QueueBind
func (r *RabbitmqReader) bind(queue, exchange, bindingKey string) {
	if err := r.channel.QueueBind(
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

// Consume rabbitmq queue
func (r *RabbitmqReader) Consume(name, binding string) (<-chan amqp.Delivery, error) {
	return r.channel.Consume(
		name,    // name
		binding, // bindingKey,
		false,   // noAck
		false,   // exclusive
		false,   // noLocal
		false,   // noWait
		nil,     // arguments
	)
}
