package controllers

import (
	"log"
	"net/http"

	"backend/ws"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var interviewUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins for now, can be restricted in production
		return true
	},
}

// InterviewWebSocketHandler handles WebSocket connections for real-time interview updates
func InterviewWebSocketHandler(c *gin.Context) {
	// Get user ID and role from JWT token (set by JWTAuth middleware)
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userID, ok := userIDInterface.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	roleInterface, exists := c.Get("role")
	if !exists {
		// Default to "student" if role is not explicitly set (e.g., for general users)
		roleInterface = "student"
	}
	role, ok := roleInterface.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid role type"})
		return
	}


	// Check if InterviewHub is initialized
	if ws.InterviewHubInstance == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Interview hub not initialized"})
		return
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := interviewUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade WebSocket connection for interview: %v", err)
		return
	}

	// Create new client and register with hub
	client := ws.NewInterviewClient(ws.InterviewHubInstance, conn, userID, role)
	ws.InterviewHubInstance.Register <- client

	log.Printf("🔌 New interview WebSocket client connected (UserID: %d, Role: %s)", userID, role)

	// Start goroutines for reading and writing
	go client.WritePump()
	go client.ReadPump()
}
