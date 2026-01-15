// src/hooks/useScreeningSocket.ts
// Vue Composable สำหรับ Screening WebSocket

import { ref, onMounted, onUnmounted } from 'vue';
import {
  screeningWs,
  type ScreeningResult,
  type BatchProgress,
  type BatchComplete
} from '@/services/api/screeningWebSocket';

export function useScreeningSocket() {
  const isConnected = ref(false);
  const lastResult = ref<ScreeningResult | null>(null);
  const currentProgress = ref<BatchProgress | null>(null);
  const results = ref<ScreeningResult[]>([]);

  // Unsubscribe functions
  let unsubResult: (() => void) | null = null;
  let unsubProgress: (() => void) | null = null;
  let unsubComplete: (() => void) | null = null;
  let unsubConnection: (() => void) | null = null;

  // Connection status callback
  const handleConnectionChange = (connected: boolean) => {
    isConnected.value = connected;
  };

  // Screening result callback
  const handleResult = (result: ScreeningResult) => {
    lastResult.value = result;
    results.value.unshift(result); // Add to beginning

    // Keep only last 50 results
    if (results.value.length > 50) {
      results.value = results.value.slice(0, 50);
    }
  };

  // Progress callback
  const handleProgress = (progress: BatchProgress) => {
    currentProgress.value = progress;
  };

  // Batch complete callback
  const handleBatchComplete = (_complete: BatchComplete) => {
    currentProgress.value = null; // Clear progress when done
  };

  // Connect on mount
  onMounted(() => {
    screeningWs.connect();

    unsubConnection = screeningWs.onConnectionChange(handleConnectionChange);
    unsubResult = screeningWs.onScreeningResult(handleResult);
    unsubProgress = screeningWs.onProgress(handleProgress);
    unsubComplete = screeningWs.onBatchComplete(handleBatchComplete);

    // Set initial connection status
    isConnected.value = screeningWs.isConnected();
  });

  // Cleanup on unmount
  onUnmounted(() => {
    unsubConnection?.();
    unsubResult?.();
    unsubProgress?.();
    unsubComplete?.();
  });

  // Manual connect/disconnect functions
  const connect = () => {
    screeningWs.connect();
  };

  const disconnect = () => {
    screeningWs.disconnect();
    isConnected.value = false;
  };

  // Clear results
  const clearResults = () => {
    results.value = [];
    lastResult.value = null;
  };

  return {
    // State
    isConnected,
    lastResult,
    currentProgress,
    results,

    // Actions
    connect,
    disconnect,
    clearResults
  };
}
