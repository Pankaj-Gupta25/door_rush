package messaging

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sachinggsingh/PDTS/parcel-service/internal/utils"
	"github.com/segmentio/kafka-go"
)

type KafkaProducer interface {
	PublishParcelEvent(ctx context.Context, parcelID string, status string) error
	Close() error
}

type kafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(brokers []string, topic string) KafkaProducer {
	return &kafkaProducer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *kafkaProducer) PublishParcelEvent(ctx context.Context, parcelID string, status string) error {
	event := map[string]string{
		"parcel_id": parcelID,
		"status":    status,
	}

	payload, err := json.Marshal(event)
	if err != nil {
		utils.Log.Error("Failed to marshal event")
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	err = p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(parcelID),
		Value: payload,
	})
	if err != nil {
		utils.Log.Error("Failed to write message to kafka")
		return fmt.Errorf("failed to write message to kafka: %w", err)
	}
	return nil
}

func (p *kafkaProducer) Close() error {
	return p.writer.Close()
}

// Mock Kafka Producer for testing
type MockKafkaProducer struct {
	LastParcelID string
	LastStatus   string
	Err          error
}

func (m *MockKafkaProducer) PublishParcelEvent(ctx context.Context, parcelID string, status string) error {
	m.LastParcelID = parcelID
	m.LastStatus = status

	return m.Err
}

func (m *MockKafkaProducer) Close() error {
	return nil
}
