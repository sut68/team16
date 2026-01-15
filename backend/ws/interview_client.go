package ws

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	interviewWriteWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	interviewPongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	interviewPingPeriod = (interviewPongWait * 9) / 10
)

// InterviewClient is a middleman between the WebSocket connection and the hub.
type InterviewClient struct {
	Hub *InterviewHub

	// The WebSocket connection.
	Conn *websocket.Conn

	// Buffered channel of outbound messages.
	Send chan []byte

	// UserID of the connected user
	UserID uint

	// Role of the connected user ("student" or "admin")
	Role string
}

// NewInterviewClient creates a new InterviewClient
func NewInterviewClient(hub *InterviewHub, conn *websocket.Conn, userID uint, role string) *InterviewClient {
	return &InterviewClient{
		Hub:    hub,
		Conn:   conn,
		Send:   make(chan []byte, 256),
		UserID: userID,
		Role:   role,
	}
}

// ReadPump pumps messages from the WebSocket connection to the hub.
func (c *InterviewClient) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()
	c.Conn.SetReadDeadline(time.Now().Add(interviewPongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(interviewPongWait)); return nil })
	for {
		// This client is for server-sent events, so we don't expect messages from client.
		// This loop is mainly for handling pong and disconnect.
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
	}
}

// WritePump pumps messages from the hub to the WebSocket connection.
func (c *InterviewClient) WritePump() {
	ticker := time.NewTicker(interviewPingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(interviewWriteWait))
			if !ok {
				// The hub closed the channel.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(interviewWriteWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
