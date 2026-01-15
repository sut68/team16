package ws

import (
	"encoding/json"
	"log"
	"sync"
)

// InterviewMessage represents a message for interview broadcasts
type InterviewMessage struct {
	Type string      `json:"type"` // e.g., "interview_slot_updated", "interview_booking_created"
	Data interface{} `json:"data"`
}

// InterviewHub manages WebSocket clients for interview notifications
type InterviewHub struct {
	// Registered clients
	Clients map[*InterviewClient]bool

	// Register requests from clients
	Register chan *InterviewClient

	// Unregister requests from clients
	Unregister chan *InterviewClient

	// Broadcast channel for messages to all clients
	Broadcast chan InterviewMessage

	// Mutex for thread-safe operations
	mu sync.RWMutex
}

// Global instance of InterviewHub
var InterviewHubInstance *InterviewHub

// NewInterviewHub creates a new InterviewHub instance
func NewInterviewHub() *InterviewHub {
	return &InterviewHub{
		Clients:    make(map[*InterviewClient]bool),
		Register:   make(chan *InterviewClient),
		Unregister: make(chan *InterviewClient),
		Broadcast:  make(chan InterviewMessage, 256), // buffered channel
	}
}

// InitInterviewHub initializes the global InterviewHub instance
func InitInterviewHub() {
	InterviewHubInstance = NewInterviewHub()
	go InterviewHubInstance.Run()
	log.Println("✅ InterviewHub initialized and running")
}

// Run starts the InterviewHub's main loop
func (h *InterviewHub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.Clients[client] = true
			h.mu.Unlock()
			log.Printf("📢 InterviewHub: Client registered (UserID: %d, Role: %s, Total: %d)", client.UserID, client.Role, len(h.Clients))

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}
			h.mu.Unlock()
			log.Printf("📤 InterviewHub: Client unregistered (UserID: %d, Total: %d)", client.UserID, len(h.Clients))

		case msg := <-h.Broadcast:
			h.mu.RLock()
			data, err := json.Marshal(msg)
			if err != nil {
				log.Printf("❌ InterviewHub: Failed to marshal message: %v", err)
				h.mu.RUnlock()
				continue
			}

			// Broadcast to all interview clients
			for client := range h.Clients {
				select {
				case client.Send <- data:
					// Message sent successfully
				default:
					// Client's send buffer is full, remove it
					h.mu.RUnlock()
					h.mu.Lock()
					close(client.Send)
					delete(h.Clients, client)
					h.mu.Unlock()
					h.mu.RLock()
				}
			}
			h.mu.RUnlock()
		}
	}
}

// BroadcastUpdate sends a message to all connected interview clients.
func (h *InterviewHub) BroadcastUpdate(msgType string, data interface{}) {
	if h == nil {
		log.Println("⚠️ InterviewHubInstance is nil. Cannot broadcast.")
		return
	}
	message := InterviewMessage{
		Type: msgType,
		Data: data,
	}
	h.Broadcast <- message
}

// GetClientCount returns the number of connected clients
func (h *InterviewHub) GetClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.Clients)
}
