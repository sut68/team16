package controllers

import (
	"log"
	"net/http"

	"backend/ws"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var approvalUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins for now, can be restricted in production
		return true
	},
}

// ApprovalWebSocketHandler handles WebSocket connections for approval notifications
func ApprovalWebSocketHandler(c *gin.Context) {
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
        c.JSON(http.StatusForbidden, gin.H{"error": "User role not found"})
		return
    }

    role, ok := roleInterface.(string)
    if !ok || role != "admin" {
        c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions. Admin role required."})
		return
    }


	// Check if ApprovalHub is initialized
	if ws.ApprovalHubInstance == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Approval hub not initialized"})
		return
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := approvalUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade WebSocket connection for approval: %v", err)
		return
	}

	// Create new client and register with hub
	client := ws.NewApprovalClient(ws.ApprovalHubInstance, conn, userID)
	ws.ApprovalHubInstance.Register <- client

	log.Printf("🔌 New approval WebSocket client connected (UserID: %d)", userID)

	// Start goroutines for reading and writing
	go client.WritePump()
	go client.ReadPump()
}
