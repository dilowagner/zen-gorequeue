package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"zen-goducer/consumer"
	"zen-goducer/producer"
)

func main() {
	// Kafka reader using environment variables.
	kafkaURL := os.Getenv("KAFKA_URL")
	topic := os.Getenv("KAFKA_READ_TOPIC")
	group := os.Getenv("KAFKA_GROUP_ID")

	// Rabbitmq writer using environment variables.
	rabbitmqURL := os.Getenv("RABBITMQ_URL")
	port := os.Getenv("RABBITMQ_PORT")
	user := os.Getenv("RABBITMQ_USER")
	pass := os.Getenv("RABBITMQ_PASS")
	receiver := os.Getenv("RABBITMQ_RECEIVER_QUEUE")
	exchange := os.Getenv("RABBITMQ_EXCHANGE")
	binding := os.Getenv("RABBITMQ_BINDING_KEY")

	// Kafka Consumer Client
	consumer := consumer.NewReader(kafkaURL, topic, group)
	defer consumer.Close()

	// Rabbitmq Producer Client
	producer := producer.NewWriter(rabbitmqURL, user, pass, port)
	producer.CreateReceiverQueue(receiver, exchange, binding)

	fmt.Println("start consuming ... !!")
	for {
		m, err := consumer.Read(context.Background())
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("message at topic:%v partition:%v offset:%v	%s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
		err = producer.Publish([]byte(`{"key":"teste outro teste"`), "20000")
		if err != nil {
			log.Fatalln(err)
		}
	}
}
