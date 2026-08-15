package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	fmt.Println("Starting Kafka Check Script...")

	topic := "parcel.events"
	broker := "localhost:9092"
	groupID := fmt.Sprintf("dev-tools-test-group-%d", time.Now().Unix()) // Unique group to ensure we get the message

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// 1. Start Consumer in Background
	// We use a channel to signal when we've received a message
	doneChan := make(chan struct{})

	go func() {
		fmt.Printf("Consumer: Starting to listen on topic '%s' (Group: %s)...\n", topic, groupID)
		reader := kafka.NewReader(kafka.ReaderConfig{
			Brokers:     []string{broker},
			Topic:       topic,
			GroupID:     groupID,
			MinBytes:    10e3,
			MaxBytes:    10e6,
			StartOffset: kafka.FirstOffset,
		})
		defer reader.Close()

		for {
			m, err := reader.ReadMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return // Context cancelled/timeout
				}
				fmt.Printf("Consumer Error: %v\n", err)
				continue
			}
			// Pretty print the value
			var prettyJSON map[string]interface{}
			if err := json.Unmarshal(m.Value, &prettyJSON); err != nil {
				fmt.Printf("\n[CONSUMER] RECEIVED MESSAGE (Raw):\n%s\n", string(m.Value))
			} else {
				prettyBytes, _ := json.MarshalIndent(prettyJSON, "", "  ")
				fmt.Printf("\n[CONSUMER] RECEIVED MESSAGE!\nTopic: %s\nPartition: %d\nOffset: %d\nKey: %s\nValue:\n%s\n\n",
					m.Topic, m.Partition, m.Offset, string(m.Key), string(prettyBytes))
			}
			close(doneChan)
			return // Exit after one message for this test
		}
	}()

	// Wait a moment for consumer group to join/initialize
	fmt.Println("Waiting 2 seconds for consumer initialization...")
	time.Sleep(2 * time.Second)

	// 2. Produce Message
	fmt.Printf("Producer: Connecting to %s...\n", broker)
	writer := &kafka.Writer{
		Addr:     kafka.TCP(broker),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	defer writer.Close()

	timestamp := time.Now().Unix()
	messageKey := fmt.Sprintf("test-key-%d", timestamp)
	messageValue := fmt.Sprintf(`{"parcel_id": "test-parcel-%d", "status": "picked_up", "current_location": "Test City"}`, timestamp)

	fmt.Printf("Producer: Sending message... Key=%s\n", messageKey)
	err := writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(messageKey),
		Value: []byte(messageValue),
	})

	if err != nil {
		log.Fatalf("Producer FAILED: %v\n", err)
	}
	fmt.Println("Producer: SUCCESS! Message sent.")

	// Wait for consumer to finish or timeout
	select {
	case <-doneChan:
		fmt.Println("Test finished SUCCESS: Consumer received the message.")
	case <-ctx.Done():
		fmt.Println("Test finished TIMEOUT: Consumer did not receive the message in time.")
	}
}
