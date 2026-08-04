package websockets

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Since we are testing Hub which uses real channels and goroutines,
// we'll use a simple mock-like approach.

func TestHub_RegisterUnregister(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	client := &Client{
		Hub:      hub,
		ParcelID: "parcel-1",
		Send:     make(chan []byte, 256),
	}

	hub.Register <- client

	// Give it a moment to register
	time.Sleep(10 * time.Millisecond)

	hub.mu.Lock()
	assert.True(t, hub.Clients["parcel-1"][client])
	hub.mu.Unlock()

	hub.Unregister <- client
	time.Sleep(10 * time.Millisecond)

	hub.mu.Lock()
	assert.Nil(t, hub.Clients["parcel-1"])
	hub.mu.Unlock()
}

func TestHub_Broadcast(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	client := &Client{
		Hub:      hub,
		ParcelID: "parcel-1",
		Send:     make(chan []byte, 256),
	}

	hub.Register <- client
	time.Sleep(10 * time.Millisecond)

	update := map[string]string{
		"parcel_id": "parcel-1",
		"status":    "IN_TRANSIT",
	}
	msg, _ := json.Marshal(update)

	hub.Broadcast <- msg

	select {
	case received := <-client.Send:
		var data map[string]string
		json.Unmarshal(received, &data)
		assert.Equal(t, "parcel-1", data["parcel_id"])
		assert.Equal(t, "IN_TRANSIT", data["status"])
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timed out waiting for broadcast")
	}
}
