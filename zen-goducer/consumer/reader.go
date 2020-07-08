package consumer

import (
	"context"
	"strings"

	"github.com/segmentio/kafka-go"
)

// KafkaReader struct
type KafkaReader struct {
	reader *kafka.Reader
}

// NewReader method - instantiate new consumer Callback Broker
func NewReader(uri, topic, group string) *KafkaReader {
	brokers := strings.Split(uri, ",")
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  group,
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	return &KafkaReader{
		reader: r,
	}
}

// Close method messages
func (r KafkaReader) Close() error {
	return r.reader.Close()
}

// Read method messages
func (r KafkaReader) Read(ctx context.Context) (kafka.Message, error) {
	return r.reader.ReadMessage(ctx)
}
