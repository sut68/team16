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

    // อ่าน WebSocket URL จาก ENV หรือ derive จาก API URL
    const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080/api';

    // สร้าง WebSocket URL
    let wsUrl: string;
    if (import.meta.env.VITE_WS_URL) {
      // ถ้ามีการตั้งค่า VITE_WS_URL ใช้ค่านั้น
      wsUrl = import.meta.env.VITE_WS_URL;
    } else if (apiUrl.startsWith('/')) {
      // ถ้าเป็น relative path (เช่น /api) ให้สร้าง full URL จาก window.location
      const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
      wsUrl = `${protocol}//${window.location.host}${apiUrl}`;
    } else {
      // ถ้าเป็น absolute URL ให้แปลง http → ws
      wsUrl = apiUrl.replace('http', 'ws');
    }

    const url = `${wsUrl}/ws?token=${token}&chatroom_id=${chatroomId}`;
    this.ws = new WebSocket(url);

    this.ws.onopen = () => {
      // console.log("✅ WebSocket connected:", chatroomId);
    };

    this.ws.onclose = () => {
      // console.log("❌ WebSocket disconnected");
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
