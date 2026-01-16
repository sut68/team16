import { ref, onMounted, onUnmounted, type Ref } from 'vue';
import { approvalWs } from '@/services/api/approvalWebSocket';

/**
 * Vue Composable สำหรับ EvaluationList - รับ notification เมื่อมีการอัปเดต evaluation
 * 
 * Flow:
 * 1. Admin บันทึก evaluation หรือตัดสินผล
 * 2. Backend broadcast 'evaluation_updated' / 'application_status_updated' event
 * 3. Hook นี้รับ event และอัปเดต reactive ref
 * 4. Component ใช้ watch เพื่อ react ต่อการเปลี่ยนแปลง
 */
export function useEvaluationWebSocket() {
  const isConnected = ref(false);
  const lastUpdate: Ref<any> = ref(null);
  const updateCount = ref(0);

  const unsubscribers: (() => void)[] = [];

  const handleConnectionChange = (connected: boolean) => {
    isConnected.value = connected;
    // console.log(`🔌 [Evaluation] WebSocket ${connected ? 'connected ✅' : 'disconnected ❌'}`);
  };

  const handleUpdate = (data: any) => {
    // console.log('📢 [Evaluation] Received update:', data);
    lastUpdate.value = data;
    updateCount.value++;
  };

  onMounted(() => {
    // console.log('🚀 [Evaluation] Mounting WebSocket hook...');

    // Connect to WebSocket (จะ reuse connection เดิมถ้ามีอยู่แล้ว)
    approvalWs.connect();

    // Subscribe to events
    unsubscribers.push(approvalWs.onConnectionChange(handleConnectionChange));

    // Listen for evaluation-related events
    unsubscribers.push(approvalWs.on('evaluation_updated', handleUpdate));
    unsubscribers.push(approvalWs.on('evaluation_created', handleUpdate));
    unsubscribers.push(approvalWs.on('evaluation_completed', handleUpdate));
    unsubscribers.push(approvalWs.on('application_status_updated', handleUpdate));

    // Check current connection status
    isConnected.value = approvalWs.isConnected();
    // console.log(`🔌 [Evaluation] Initial connection status: ${isConnected.value ? 'connected' : 'disconnected'}`);
  });

  onUnmounted(() => {
    // console.log('🧹 [Evaluation] Unmounting WebSocket hook...');
    // Unsubscribe from events
    unsubscribers.forEach(unsub => unsub());

    // Don't disconnect - other components might be using it
  });

  return {
    isConnected,
    lastUpdate,
    updateCount,
  };
}
