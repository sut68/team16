package ws

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer
	screeningWriteWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	screeningPongWait = 60 * time.Second

	// Send pings to peer with this period (must be less than pongWait)
	screeningPingPeriod = (screeningPongWait * 9) / 10

	// Maximum message size allowed from peer
	screeningMaxMessageSize = 512
)

// ScreeningClient represents a WebSocket client for screening notifications
type ScreeningClient struct {
	Hub     *ScreeningHub
	Conn    *websocket.Conn
	Send    chan []byte
	UserID  uint
	IsAdmin bool
}

// NewScreeningClient creates a new ScreeningClient
func NewScreeningClient(hub *ScreeningHub, conn *websocket.Conn, userID uint, isAdmin bool) *ScreeningClient {
	return &ScreeningClient{
		Hub:     hub,
		Conn:    conn,
		Send:    make(chan []byte, 256),
		UserID:  userID,
		IsAdmin: isAdmin,
	}
}

// ReadPump pumps messages from the WebSocket connection to the hub
// In screening context, we mainly listen for client pings/pongs
func (c *ScreeningClient) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(screeningMaxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(screeningPongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(screeningPongWait))
		return nil
	})

	for {
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("ScreeningClient ReadPump error: %v", err)
			}
			break
		}
		// For screening, we don't process incoming messages from clients
		// The hub is primarily for broadcasting results TO clients
	}
}

// WritePump pumps messages from the hub to the WebSocket connection
func (c *ScreeningClient) WritePump() {
	ticker := time.NewTicker(screeningPingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(screeningWriteWait))
			if !ok {
				// The hub closed the channel
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages to the current WebSocket message
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(screeningWriteWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
