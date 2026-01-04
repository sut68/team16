// src/services/api/websocket.ts
import type { AssistanceRespone } from "@/interfaces/assistance";

type MessageCallback = (msg: AssistanceRespone) => void;

export class ChatService {
  private ws: WebSocket | null = null;
  private callbacks: MessageCallback[] = [];

  /**
   * connect websocket (ผูก token + chatroom)
   */
  connect(token: string, chatroomId: number) {
    // กัน connect ซ้ำ
    if (this.ws) return;

    const url = `ws://localhost:8080/api/ws?token=${token}&chatroom_id=${chatroomId}`;
    this.ws = new WebSocket(url);

    this.ws.onopen = () => {
      console.log("✅ WebSocket connected:", chatroomId);
    };

    this.ws.onclose = () => {
      console.log("❌ WebSocket disconnected");
      this.ws = null;
    };

    this.ws.onerror = (err) => {
      console.error("WebSocket error:", err);
    };

    this.ws.onmessage = (event) => {
      try {
        const data: AssistanceRespone = JSON.parse(event.data);
        this.callbacks.forEach(cb => cb(data));
      } catch (e) {
        console.error("Invalid websocket message:", event.data);
      }
    };
  }

  /**
   * ส่งข้อความ (ส่งเฉพาะข้อความ)
   */
  sendMessage(message: string) {
    if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
      console.warn("WebSocket not connected");
      return;
    }

    this.ws.send(JSON.stringify({
      massage: message,
    }));
  }

  /**
   * subscribe รับข้อความ
   */
  onMessage(callback: MessageCallback) {
    this.callbacks.push(callback);
  }
  
  isConnected(): boolean {
    return this.ws !== null && this.ws.readyState === WebSocket.OPEN
  }

  /**
   * disconnect websocket
   */
  disconnect() {
    if (this.ws) {
      this.ws.close();
      this.ws = null;
      this.callbacks = [];
    }
  }
}
