package ws

import (
	"encoding/json"
	"log"
	"sync"
)

// ApprovalMessage represents a message for approval broadcasts
type ApprovalMessage struct {
	Type string      `json:"type"` // e.g., "approval_task_updated"
	Data interface{} `json:"data"`
}

// ApprovalHub manages WebSocket clients for approval notifications
type ApprovalHub struct {
	// Registered clients
	Clients map[*ApprovalClient]bool

	// Register requests from clients
	Register chan *ApprovalClient

	// Unregister requests from clients
	Unregister chan *ApprovalClient

	// Broadcast channel for messages to all clients
	Broadcast chan ApprovalMessage

	// Mutex for thread-safe operations
	mu sync.RWMutex
}

// Global instance of ApprovalHub
var ApprovalHubInstance *ApprovalHub

// NewApprovalHub creates a new ApprovalHub instance
func NewApprovalHub() *ApprovalHub {
	return &ApprovalHub{
		Clients:    make(map[*ApprovalClient]bool),
		Register:   make(chan *ApprovalClient),
		Unregister: make(chan *ApprovalClient),
		Broadcast:  make(chan ApprovalMessage, 256), // buffered channel
	}
}

// InitApprovalHub initializes the global ApprovalHub instance
func InitApprovalHub() {
	ApprovalHubInstance = NewApprovalHub()
	go ApprovalHubInstance.Run()
	log.Println("✅ ApprovalHub initialized and running")
}

// Run starts the ApprovalHub's main loop
func (h *ApprovalHub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.Clients[client] = true
			h.mu.Unlock()
			log.Printf("📢 ApprovalHub: Client registered (UserID: %d, Total: %d)", client.UserID, len(h.Clients))

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}
			h.mu.Unlock()
			log.Printf("📤 ApprovalHub: Client unregistered (UserID: %d, Total: %d)", client.UserID, len(h.Clients))

		case msg := <-h.Broadcast:
			h.mu.RLock()
			data, err := json.Marshal(msg)
			if err != nil {
				log.Printf("❌ ApprovalHub: Failed to marshal message: %v", err)
				h.mu.RUnlock()
				continue
			}

			// Broadcast to all approval clients
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

// BroadcastUpdate sends a message to all connected approval clients.
func (h *ApprovalHub) BroadcastUpdate(msgType string, data interface{}) {
	message := ApprovalMessage{
		Type: msgType,
		Data: data,
	}
	h.Broadcast <- message
}

// GetClientCount returns the number of connected clients
func (h *ApprovalHub) GetClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.Clients)
}
