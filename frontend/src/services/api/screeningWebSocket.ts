// src/services/api/screeningWebSocket.ts
// WebSocket service สำหรับ Screening Notifications (Realtime)

export interface ScreeningResult {
  screening_id: number;
  application_id: number;
  scholarship_id: number;
  scholarship_name: string;
  student_name: string;
  passed: boolean;
  passed_criteria: number;
  total_criteria: number;
  failed_reasons?: string[];
  processed_by?: string;
}

export interface BatchProgress {
  scholarship_id: number;
  scholarship_name: string;
  total: number;
  processed: number;
  passed: number;
  failed: number;
}

export interface BatchComplete {
  scholarship_id: number;
  scholarship_name: string;
  total: number;
  passed: number;
  failed: number;
}

export interface ScreeningMessage {
  type: 'screening_result' | 'progress' | 'batch_complete';
  data: ScreeningResult | BatchProgress | BatchComplete;
}

type ScreeningResultCallback = (result: ScreeningResult) => void;
type ProgressCallback = (progress: BatchProgress) => void;
type BatchCompleteCallback = (complete: BatchComplete) => void;
type ConnectionCallback = (connected: boolean) => void;

export class ScreeningWebSocketService {
  private ws: WebSocket | null = null;
  private resultCallbacks: ScreeningResultCallback[] = [];
  private progressCallbacks: ProgressCallback[] = [];
  private batchCompleteCallbacks: BatchCompleteCallback[] = [];
  private connectionCallbacks: ConnectionCallback[] = [];
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 5;
  private reconnectDelay = 3000;

  /**
   * Connect to screening WebSocket
   * ใช้ credentials (cookies) สำหรับ authentication
   */
  connect(): void {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      return;
    }

    const apiUrl = import.meta.env.VITE_API_URL || '/api';
    let wsUrl: string;

    if (apiUrl.startsWith('/')) {
      // Relative path - สร้าง full URL จาก window.location
      const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
      wsUrl = `${protocol}//${window.location.host}${apiUrl}/ws/screening`;
    } else {
      // Absolute URL - แปลง http → ws
      wsUrl = apiUrl.replace('http', 'ws') + '/ws/screening';
    }

    this.ws = new WebSocket(wsUrl);

    this.ws.onopen = () => {
      this.reconnectAttempts = 0;
      this.connectionCallbacks.forEach(cb => cb(true));
    };

    this.ws.onclose = (event) => {
      this.ws = null;
      this.connectionCallbacks.forEach(cb => cb(false));

      // Auto reconnect if not intentionally closed
      if (event.code !== 1000 && this.reconnectAttempts < this.maxReconnectAttempts) {
        this.reconnectAttempts++;
        setTimeout(() => this.connect(), this.reconnectDelay);
      }
    };

    this.ws.onerror = (error) => {
      console.error('⚠️ Screening WebSocket error:', error);
    };

    this.ws.onmessage = (event) => {
      try {
        const message: ScreeningMessage = JSON.parse(event.data);
        this.handleMessage(message);
      } catch (e) {
        console.error('Invalid screening WebSocket message:', event.data);
      }
    };
  }

  /**
   * Handle incoming message and dispatch to callbacks
   */
  private handleMessage(message: ScreeningMessage): void {
    switch (message.type) {
      case 'screening_result':
        const result = message.data as ScreeningResult;
        this.resultCallbacks.forEach(cb => cb(result));
        break;

      case 'progress':
        const progress = message.data as BatchProgress;
        this.progressCallbacks.forEach(cb => cb(progress));
        break;

      case 'batch_complete':
        const complete = message.data as BatchComplete;
        this.batchCompleteCallbacks.forEach(cb => cb(complete));
        break;
    }
  }

  /**
   * Subscribe to screening results
   */
  onScreeningResult(callback: ScreeningResultCallback): () => void {
    this.resultCallbacks.push(callback);
    return () => {
      this.resultCallbacks = this.resultCallbacks.filter(cb => cb !== callback);
    };
  }

  /**
   * Subscribe to batch progress
   */
  onProgress(callback: ProgressCallback): () => void {
    this.progressCallbacks.push(callback);
    return () => {
      this.progressCallbacks = this.progressCallbacks.filter(cb => cb !== callback);
    };
  }

  /**
   * Subscribe to batch completion
   */
  onBatchComplete(callback: BatchCompleteCallback): () => void {
    this.batchCompleteCallbacks.push(callback);
    return () => {
      this.batchCompleteCallbacks = this.batchCompleteCallbacks.filter(cb => cb !== callback);
    };
  }

  /**
   * Subscribe to connection status changes
   */
  onConnectionChange(callback: ConnectionCallback): () => void {
    this.connectionCallbacks.push(callback);
    return () => {
      this.connectionCallbacks = this.connectionCallbacks.filter(cb => cb !== callback);
    };
  }

  /**
   * Check if connected
   */
  isConnected(): boolean {
    return this.ws !== null && this.ws.readyState === WebSocket.OPEN;
  }

  /**
   * Disconnect WebSocket
   */
  disconnect(): void {
    if (this.ws) {
      this.ws.close(1000, 'User disconnected');
      this.ws = null;
      this.resultCallbacks = [];
      this.progressCallbacks = [];
      this.batchCompleteCallbacks = [];
      this.connectionCallbacks = [];
    }
  }
}

// Singleton instance
export const screeningWs = new ScreeningWebSocketService();
