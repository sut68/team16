// src/services/api/approvalWebSocket.ts
// WebSocket service for Approval Notifications (Realtime)

import type { ApprovalTaskResponse } from '@/interfaces';

// Define the structure of messages coming from the approval websocket
export interface ApprovalMessage {
  type: 'approval_task_updated'; // Add other types here later
  data: ApprovalTaskResponse; // Using ApprovalTaskResponse for now, can be a union type later
}

// Define callback signatures
type ApprovalMessageCallback<T> = (data: T) => void;
type ConnectionCallback = (connected: boolean) => void;

export class ApprovalWebSocketService {
  private ws: WebSocket | null = null;
  private connectionCallbacks: ConnectionCallback[] = [];
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 5;
  private reconnectDelay = 3000;

  // A map to hold callbacks for specific message types
  private messageCallbacks: Map<string, ApprovalMessageCallback<any>[]> = new Map();

  /**
   * Connect to the approval WebSocket using credentials (cookies) for authentication.
   */
  connect(): void {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      console.log('🔌 Approval WebSocket already connected');
      return;
    }

    const apiUrl = import.meta.env.VITE_API_URL || '/api';
    let wsUrl: string;

    if (apiUrl.startsWith('/')) {
      const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
      wsUrl = `${protocol}//${window.location.host}${apiUrl}/ws/approval`;
    } else {
      wsUrl = apiUrl.replace('http', 'ws') + '/ws/approval';
    }

    console.log('🔌 Connecting to Approval WebSocket:', wsUrl);
    this.ws = new WebSocket(wsUrl);

    this.ws.onopen = () => {
      console.log('✅ Approval WebSocket connected');
      this.reconnectAttempts = 0;
      this.connectionCallbacks.forEach(cb => cb(true));
    };

    this.ws.onclose = (event) => {
      console.log('❌ Approval WebSocket disconnected:', event.code, event.reason);
      this.ws = null;
      this.connectionCallbacks.forEach(cb => cb(false));

      // Auto-reconnect if not intentionally closed
      if (event.code !== 1000 && this.reconnectAttempts < this.maxReconnectAttempts) {
        this.reconnectAttempts++;
        console.log(`🔄 Reconnecting approval ws in ${this.reconnectDelay}ms (attempt ${this.reconnectAttempts}/${this.maxReconnectAttempts})`);
        setTimeout(() => this.connect(), this.reconnectDelay);
      }
    };

    this.ws.onerror = (error) => {
      console.error('⚠️ Approval WebSocket error:', error);
    };

    this.ws.onmessage = (event) => {
      try {
        const message: ApprovalMessage = JSON.parse(event.data);
        this.handleMessage(message);
      } catch (e) {
        console.error('Invalid approval WebSocket message:', event.data);
      }
    };
  }

  /**
   * Handle incoming messages and dispatch to the correct callbacks based on message type.
   */
  private handleMessage(message: ApprovalMessage): void {
    const callbacks = this.messageCallbacks.get(message.type);
    if (callbacks) {
      console.log(`📢 Approval WS [${message.type}]:`, message.data);
      callbacks.forEach(cb => cb(message.data));
    } else {
      console.warn(`No handler for approval ws message type: ${message.type}`);
    }
  }

  /**
   * Subscribe to a specific type of message from the approval channel.
   * @param type The message type to subscribe to.
   * @param callback The function to call when a message of that type is received.
   * @returns An unsubscribe function.
   */
  on<T>(type: string, callback: ApprovalMessageCallback<T>): () => void {
    if (!this.messageCallbacks.has(type)) {
      this.messageCallbacks.set(type, []);
    }
    this.messageCallbacks.get(type)!.push(callback);

    // Return an unsubscribe function
    return () => {
      const callbacks = this.messageCallbacks.get(type);
      if (callbacks) {
        const index = callbacks.indexOf(callback);
        if (index > -1) {
          callbacks.splice(index, 1);
        }
      }
    };
  }

  /**
   * Subscribe to connection status changes.
   */
  onConnectionChange(callback: ConnectionCallback): () => void {
    this.connectionCallbacks.push(callback);
    return () => {
      this.connectionCallbacks = this.connectionCallbacks.filter(cb => cb !== callback);
    };
  }

  /**
   * Check if the WebSocket is currently connected.
   */
  isConnected(): boolean {
    return this.ws !== null && this.ws.readyState === WebSocket.OPEN;
  }

  /**
   * Disconnect the WebSocket and clear all callbacks.
   */
  disconnect(): void {
    if (this.ws) {
      console.log('🔌 Disconnecting Approval WebSocket');
      this.ws.close(1000, 'User disconnected');
      this.ws = null;
      this.messageCallbacks.clear();
      this.connectionCallbacks = [];
    }
  }
}

// Export a singleton instance of the service
export const approvalWs = new ApprovalWebSocketService();
