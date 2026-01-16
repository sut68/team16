// src/services/api/interviewWebSocket.ts
import type { InterviewRound, Slot, InterviewBooking, InterviewBookingDTOResponse } from '@/interfaces/interview';

// This combines all possible data types from the interview websocket
export type InterviewData = InterviewRound | Slot | InterviewBooking | { id: number };

export type InterviewMessageType =
  | 'INTERVIEW_ROUND_CREATED'
  | 'INTERVIEW_ROUND_UPDATED'
  | 'INTERVIEW_ROUND_DELETED'
  | 'INTERVIEW_SLOT_UPDATED'
  | 'INTERVIEW_BOOKING_CREATED'
  | 'INTERVIEW_BOOKING_UPDATED'
  | 'INTERVIEW_BOOKING_DELETED';

export interface InterviewMessage {
  type: InterviewMessageType;
  data: any;
}

type InterviewMessageCallback<T> = (data: T) => void;
type ConnectionCallback = (connected: boolean) => void;

export class InterviewWebSocketService {
  private ws: WebSocket | null = null;
  private connectionCallbacks: ConnectionCallback[] = [];
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 5;
  private reconnectDelay = 3000;

  // A map to hold arrays of callbacks for each message type
  private messageCallbacks: Map<InterviewMessageType, InterviewMessageCallback<any>[]> = new Map();

  private getWebSocketURL(): string {
    const apiUrl = import.meta.env.VITE_API_URL || '/api';
    let wsUrl: string;

    if (apiUrl.startsWith('/')) {
      const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
      wsUrl = `${protocol}//${window.location.host}${apiUrl}/ws/interview`;
    } else {
      wsUrl = apiUrl.replace(/https?:\/\//, 'ws://') + '/ws/interview';
    }
    return wsUrl;
  }

  connect(): void {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      return;
    }

    const wsUrl = this.getWebSocketURL();
    // console.log('🔌 Connecting to Interview WebSocket:', wsUrl);
    this.ws = new WebSocket(wsUrl);

    this.ws.onopen = () => {
      // console.log('✅ Interview WebSocket connected');
      this.reconnectAttempts = 0;
      this.connectionCallbacks.forEach(cb => cb(true));
    };

    this.ws.onclose = (event) => {
      // console.log('❌ Interview WebSocket disconnected:', event.code, event.reason);
      this.ws = null;
      this.connectionCallbacks.forEach(cb => cb(false));

      if (event.code !== 1000 && this.reconnectAttempts < this.maxReconnectAttempts) {
        this.reconnectAttempts++;
        setTimeout(() => this.connect(), this.reconnectDelay);
      }
    };

    this.ws.onerror = (error) => {
      console.error('⚠️ Interview WebSocket error:', error);
    };

    this.ws.onmessage = (event) => {
      try {
        const message: InterviewMessage = JSON.parse(event.data);
        this.handleMessage(message);
      } catch (e) {
        console.error('Invalid interview WebSocket message:', event.data);
      }
    };
  }

  private handleMessage(message: InterviewMessage): void {
    const callbacks = this.messageCallbacks.get(message.type);
    if (callbacks) {
      // console.log(`📢 Interview WS [${message.type}]:`, message.data);
      callbacks.forEach(cb => cb(message.data));
    } else {
      console.warn(`No handler for interview ws message type: ${message.type}`);
    }
  }

  // Generic subscription method
  private on<T>(type: InterviewMessageType, callback: InterviewMessageCallback<T>): () => void {
    if (!this.messageCallbacks.has(type)) {
      this.messageCallbacks.set(type, []);
    }
    const callbacks = this.messageCallbacks.get(type)!;
    callbacks.push(callback);

    // Return an unsubscribe function
    return () => {
      const index = callbacks.indexOf(callback);
      if (index > -1) {
        callbacks.splice(index, 1);
      }
    };
  }

  // Specific subscription methods
  onRoundCreated(callback: InterviewMessageCallback<InterviewRound>) { return this.on('INTERVIEW_ROUND_CREATED', callback); }
  onRoundUpdated(callback: InterviewMessageCallback<InterviewRound>) { return this.on('INTERVIEW_ROUND_UPDATED', callback); }
  onRoundDeleted(callback: InterviewMessageCallback<{ id: number }>) { return this.on('INTERVIEW_ROUND_DELETED', callback); }
  onSlotUpdated(callback: InterviewMessageCallback<Slot>) { return this.on('INTERVIEW_SLOT_UPDATED', callback); }
  onBookingCreated(callback: InterviewMessageCallback<InterviewBookingDTOResponse>) { return this.on('INTERVIEW_BOOKING_CREATED', callback); }
  onBookingUpdated(callback: InterviewMessageCallback<InterviewBookingDTOResponse>) { return this.on('INTERVIEW_BOOKING_UPDATED', callback); }
  onBookingDeleted(callback: InterviewMessageCallback<InterviewBookingDTOResponse>) { return this.on('INTERVIEW_BOOKING_DELETED', callback); }

  onConnectionChange(callback: ConnectionCallback): () => void {
    this.connectionCallbacks.push(callback);
    return () => {
      this.connectionCallbacks = this.connectionCallbacks.filter(cb => cb !== callback);
    };
  }

  isConnected(): boolean {
    return this.ws !== null && this.ws.readyState === WebSocket.OPEN;
  }

  disconnect(): void {
    if (this.ws) {
      this.ws.close(1000, 'User disconnected');
      this.ws = null;
      this.messageCallbacks.clear();
      this.connectionCallbacks = [];
    }
  }
}

export const interviewWs = new InterviewWebSocketService();
