package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"backend/config"
	"backend/entity"
	"backend/services"
	"backend/ws"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var hub = ws.NewHub()

func init() {
	go hub.Run()
}

type IncomingMessage struct {
	Massage string `json:"massage"`
}

func WebSocketHandler(c *gin.Context) {
	token := c.Query("token")

	// If empty, try Cookie "access_token"
	if token == "" {
		if cookie, err := c.Request.Cookie("access_token"); err == nil {
			token = cookie.Value
		}
	}

	chatroomIDStr := c.Query("chatroom_id")
	chatroomID, _ := strconv.Atoi(chatroomIDStr)

	if token == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	claims, err := services.ValidateJWT(token)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// ตรวจ chatroom
	var room entity.Chatroom
	if err := config.DB.First(&room, chatroomID).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if room.Statuschatroom != "open" {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := &ws.Client{
		Conn:       conn,
		Send:       make(chan []byte),
		UserID:     claims.UserID,
		ChatroomID: uint(chatroomID),
	}

	hub.Register <- client

	go writePump(client)
	readPump(client)
}

func readPump(client *ws.Client) {
	defer func() {
		hub.Unregister <- client
		client.Conn.Close()
	}()

	for {
		_, msgBytes, err := client.Conn.ReadMessage()
		if err != nil {
			break
		}

		var incoming IncomingMessage
		if err := json.Unmarshal(msgBytes, &incoming); err != nil {
			continue
		}

		// ... (inside readPump loop)
		assist := entity.Assistance{
			Massage:    incoming.Massage,
			SenderID:   client.UserID,
			ChatroomID: client.ChatroomID,
		}

		if err := config.DB.Create(&assist).Error; err != nil {
			// Log error
			// log.Printf("Error saving message: %v", err)

			// Optional: Send error back to client
			errResp := map[string]string{"error": "Failed to save message"}
			errBytes, _ := json.Marshal(errResp)
			client.Send <- errBytes
			continue
		}

		if err := config.DB.Preload("Sender").Preload("Chatroom").First(&assist, assist.ID).Error; err != nil {
			// log.Printf("Error preloading message: %v", err)
			continue
		}

		resp, _ := json.Marshal(assist)

		hub.Broadcast <- ws.Message{
			ChatroomID: client.ChatroomID,
			Data:       resp,
		}

		// Also broadcast to admins (Notification System) to update sidebar preview
		if ws.ApprovalHubInstance != nil {
			ws.ApprovalHubInstance.BroadcastUpdate("new_chat_message", assist)
		}
	}
}

func writePump(client *ws.Client) {
	defer client.Conn.Close()

	for msg := range client.Send {
		err := client.Conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			break
		}
	}
}
