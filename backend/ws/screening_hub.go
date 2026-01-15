package ws

import (
	"encoding/json"
	"log"
	"sync"
)

// ScreeningMessage represents a message for screening broadcasts
type ScreeningMessage struct {
	Type string      `json:"type"` // "screening_result", "progress", "batch_complete"
	Data interface{} `json:"data"`
}

// ScreeningResult represents the result of a single screening
type ScreeningResult struct {
	ScreeningID     uint     `json:"screening_id"`
	ApplicationID   uint     `json:"application_id"`
	ScholarshipID   uint     `json:"scholarship_id"`
	ScholarshipName string   `json:"scholarship_name"`
	StudentName     string   `json:"student_name"`
	Passed          bool     `json:"passed"`
	PassedCriteria  int      `json:"passed_criteria"`
	TotalCriteria   int      `json:"total_criteria"`
	FailedReasons   []string `json:"failed_reasons,omitempty"`
	ProcessedBy     string   `json:"processed_by,omitempty"`
}

// BatchProgress represents progress of a batch screening operation
type BatchProgress struct {
	ScholarshipID   uint   `json:"scholarship_id"`
	ScholarshipName string `json:"scholarship_name"`
	Total           int    `json:"total"`
	Processed       int    `json:"processed"`
	Passed          int    `json:"passed"`
	Failed          int    `json:"failed"`
}

// ScreeningHub manages WebSocket clients for screening notifications
type ScreeningHub struct {
	// Registered clients
	Clients map[*ScreeningClient]bool

	// Register requests from clients
	Register chan *ScreeningClient

	// Unregister requests from clients
	Unregister chan *ScreeningClient

	// Broadcast channel for messages to all clients
	Broadcast chan ScreeningMessage

	// Mutex for thread-safe operations
	mu sync.RWMutex
}

// Global instance of ScreeningHub
var ScreeningHubInstance *ScreeningHub

// NewScreeningHub creates a new ScreeningHub instance
func NewScreeningHub() *ScreeningHub {
	return &ScreeningHub{
		Clients:    make(map[*ScreeningClient]bool),
		Register:   make(chan *ScreeningClient),
		Unregister: make(chan *ScreeningClient),
		Broadcast:  make(chan ScreeningMessage, 256), // buffered channel
	}
}

// InitScreeningHub initializes the global ScreeningHub instance
func InitScreeningHub() {
	ScreeningHubInstance = NewScreeningHub()
	go ScreeningHubInstance.Run()
	log.Println("✅ ScreeningHub initialized and running")
}

// Run starts the ScreeningHub's main loop
func (h *ScreeningHub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.Clients[client] = true
			h.mu.Unlock()
			log.Printf("📢 ScreeningHub: Client registered (UserID: %d, Total: %d)", client.UserID, len(h.Clients))

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}
			h.mu.Unlock()
			log.Printf("📤 ScreeningHub: Client unregistered (UserID: %d, Total: %d)", client.UserID, len(h.Clients))

		case msg := <-h.Broadcast:
			h.mu.RLock()
			data, err := json.Marshal(msg)
			if err != nil {
				log.Printf("❌ ScreeningHub: Failed to marshal message: %v", err)
				h.mu.RUnlock()
				continue
			}

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

// BroadcastScreeningResult broadcasts a screening result to all clients
func (h *ScreeningHub) BroadcastScreeningResult(result ScreeningResult) {
	msg := ScreeningMessage{
		Type: "screening_result",
		Data: result,
	}
	h.Broadcast <- msg
}

// BroadcastProgress broadcasts batch progress to all clients
func (h *ScreeningHub) BroadcastProgress(progress BatchProgress) {
	msg := ScreeningMessage{
		Type: "progress",
		Data: progress,
	}
	h.Broadcast <- msg
}

// BroadcastBatchComplete broadcasts batch completion to all clients
func (h *ScreeningHub) BroadcastBatchComplete(scholarshipID uint, scholarshipName string, total, passed, failed int) {
	msg := ScreeningMessage{
		Type: "batch_complete",
		Data: map[string]interface{}{
			"scholarship_id":   scholarshipID,
			"scholarship_name": scholarshipName,
			"total":            total,
			"passed":           passed,
			"failed":           failed,
		},
	}
	h.Broadcast <- msg
}

// GetClientCount returns the number of connected clients
func (h *ScreeningHub) GetClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.Clients)
}
