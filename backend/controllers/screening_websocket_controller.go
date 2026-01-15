package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"backend/ws"
)

var screeningUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins for now, can be restricted in production
		return true
	},
}

// ScreeningWebSocketHandler handles WebSocket connections for screening notifications
func ScreeningWebSocketHandler(c *gin.Context) {
	// Get user ID from JWT token (set by JWTAuth middleware)
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

	// Check if user is admin (optional - for filtering notifications later)
	isAdmin := false
	if roleInterface, exists := c.Get("role"); exists {
		if role, ok := roleInterface.(string); ok {
			isAdmin = role == "admin"
		}
	}

	// Check if ScreeningHub is initialized
	if ws.ScreeningHubInstance == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Screening hub not initialized"})
		return
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := screeningUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade WebSocket connection: %v", err)
		return
	}

	// Create new client and register with hub
	client := ws.NewScreeningClient(ws.ScreeningHubInstance, conn, userID, isAdmin)
	ws.ScreeningHubInstance.Register <- client

	log.Printf("🔌 New screening WebSocket client connected (UserID: %d, IsAdmin: %v)", userID, isAdmin)

	// Start goroutines for reading and writing
	go client.WritePump()
	go client.ReadPump()
}

// GetScreeningHubStatus returns the current status of the screening hub (for debugging)
func GetScreeningHubStatus(c *gin.Context) {
	if ws.ScreeningHubInstance == nil {
		c.JSON(http.StatusOK, gin.H{
			"status":       "not_initialized",
			"client_count": 0,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":       "running",
		"client_count": ws.ScreeningHubInstance.GetClientCount(),
	})
}
