package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"zen-gosumer/consumer"
	"zen-gosumer/producer"

	"github.com/segmentio/kafka-go"
	"github.com/streadway/amqp"
)

func main() {

	// Kafka reader using environment variables.
	kafkaURL := os.Getenv("KAFKA_URL")
	topic := os.Getenv("KAFKA_WRITE_TOPIC")

	rabbitmqURL := os.Getenv("RABBITMQ_URL")
	port := os.Getenv("RABBITMQ_PORT")
	user := os.Getenv("RABBITMQ_USER")
	pass := os.Getenv("RABBITMQ_PASS")
	resender := os.Getenv("RABBITMQ_RESENDER_QUEUE")
	exchange := os.Getenv("RABBITMQ_EXCHANGE")
	binding := os.Getenv("RABBITMQ_BINDING_KEY")

	producer := producer.NewWriter(kafkaURL, topic)

	consumer := consumer.NewReader(rabbitmqURL, user, pass, port)
	consumer.CreateExchange(exchange)
	consumer.CreateResenderQueue(resender, exchange, binding)

	deliveries, err := consumer.Consume(
		resender, // name
		binding,  // consumerTag
	)
	if err != nil {
		log.Printf("Queue Declare: %s", err)
		panic(err)
	}
	go handle(deliveries, consumer.Done, producer)

	select {}
}

func handle(deliveries <-chan amqp.Delivery, done chan error, prod *producer.KafkaWriter) {
	it := 0
	for d := range deliveries {

		log.Printf(
			"got %dB delivery: [%v] %q",
			len(d.Body),
			d.DeliveryTag,
			d.Body,
		)

		msg := kafka.Message{
			Key:   []byte(fmt.Sprintf("Key-%d", it)),
			Value: d.Body,
		}
		prod.Write(context.Background(), msg)
		it++
		d.Ack(false)
	}
	log.Printf("handle: deliveries channel closed")
	done <- nil
}
