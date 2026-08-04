package messaging

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sachinggsingh/PDTS-Go/tracking-service/internal/models"
	"github.com/sachinggsingh/PDTS-Go/tracking-service/internal/service"
	"github.com/sachinggsingh/PDTS-Go/tracking-service/internal/utils"
	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	reader  *kafka.Reader
	service service.TrackingService
	logger  *utils.Logger
}

func NewKafkaConsumer(brokers []string, topic string, groupID string, service service.TrackingService, logger *utils.Logger) *KafkaConsumer {
	return &KafkaConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			Topic:    topic,
			GroupID:  groupID,
			MinBytes: 10e3, // 10KB
			MaxBytes: 10e6, // 10MB
		}),
		service: service,
		logger:  logger,
	}
}

func (c *KafkaConsumer) Start(ctx context.Context) {
	c.logger.Info("Starting Kafka consumer...")
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			c.logger.Error(fmt.Sprintf("Failed to read message: %v", err))
			continue
		}

		c.logger.Info(fmt.Sprintf("Received message from Kafka: %s", string(m.Value)))

		var event struct {
			ParcelID        string `json:"parcel_id"`
			Status          string `json:"status"`
			CurrentLocation string `json:"current_location"`
		}

		if err := json.Unmarshal(m.Value, &event); err != nil {
			c.logger.Error(fmt.Sprintf("Failed to unmarshal event: %v", err))
			continue
		}

		// Initialize tracking
		tracking := &models.Tracking{
			ParcelID:        event.ParcelID,
			Status:          models.ParcelPickedUp,
			CurrentLocation: event.CurrentLocation,
		}

		_, err = c.service.UpsertTracking(tracking)
		if err != nil {
			c.logger.Error(fmt.Sprintf("Failed to initialize/update tracking for parcel %s: %v", event.ParcelID, err))
		} else {
			c.logger.Info(fmt.Sprintf("Successfully initialized/updated tracking for parcel %s", event.ParcelID))
		}
	}
}

func (c *KafkaConsumer) Close() error {
	return c.reader.Close()
}
