package producer

import (
	"context"
	"strings"

	"github.com/segmentio/kafka-go"
)

// KafkaWriter struct
type KafkaWriter struct {
	writer *kafka.Writer
}

// NewReader method - instantiate new consumer Callback Broker
func NewWriter(uri, topic string) *KafkaWriter {
	brokers := strings.Split(uri, ",")
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  brokers,
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})

	return &KafkaWriter{
		writer: w,
	}
}

// Close method messages
func (w *KafkaWriter) Close() error {
	return w.writer.Close()
}

// Write messages to Kafka
func (w *KafkaWriter) Write(ctx context.Context, msg kafka.Message) error {
	return w.writer.WriteMessages(ctx, msg)
}
