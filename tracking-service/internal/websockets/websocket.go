package websockets

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	Hub *Hub

	// The websocket connection.
	Conn *websocket.Conn

	// Buffered channel of outbound messages.
	Send chan []byte

	// ParcelID this client is interested in.
	ParcelID string
}

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
	// Registered clients by ParcelID.
	Clients map[string]map[*Client]bool

	// Inbound messages from the servers.
	Broadcast chan []byte

	// Register requests from the clients.
	Register chan *Client

	// Unregister requests from clients.
	Unregister chan *Client

	mu sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[string]map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			if h.Clients[client.ParcelID] == nil {
				h.Clients[client.ParcelID] = make(map[*Client]bool)
			}
			h.Clients[client.ParcelID][client] = true
			h.mu.Unlock()
		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Clients[client.ParcelID][client]; ok {
				delete(h.Clients[client.ParcelID], client)
				close(client.Send)
				if len(h.Clients[client.ParcelID]) == 0 {
					delete(h.Clients, client.ParcelID)
				}
			}
			h.mu.Unlock()
		case message := <-h.Broadcast:
			// message is expected to be a JSON with parcel_id
			var data struct {
				ParcelID string `json:"parcel_id"`
			}
			if err := json.Unmarshal(message, &data); err != nil {
				log.Printf("error unmarshaling broadcast message: %v", err)
				continue
			}

			h.mu.Lock()
			if clients, ok := h.Clients[data.ParcelID]; ok {
				for client := range clients {
					select {
					case client.Send <- message:
					default:
						close(client.Send)
						delete(h.Clients[data.ParcelID], client)
					}
				}
				if len(h.Clients[data.ParcelID]) == 0 {
					delete(h.Clients, data.ParcelID)
				}
			}
			h.mu.Unlock()
		}
	}
}

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for now
	},
}

// WritePump pumps messages from the hub to the websocket connection.
func (c *Client) WritePump() {
	defer func() {
		c.Conn.Close()
	}()
	for message := range c.Send {
		w, err := c.Conn.NextWriter(websocket.TextMessage)
		if err != nil {
			return
		}
		w.Write(message)

		// Add queued chat messages to the current websocket message.
		n := len(c.Send)
		for i := 0; i < n; i++ {
			w.Write([]byte("\n"))
			w.Write(<-c.Send)
		}

		if err := w.Close(); err != nil {
			return
		}
	}
	// The hub closed the channel.
	c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
}

// ReadPump pumps messages from the websocket connection to the hub.
// For tracking, we mainly care about outgoing updates, but we might want to keep the connection alive.
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()
	for {
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
	}
}
