import { ref, onMounted, onUnmounted, type Ref } from 'vue';
import { approvalWs } from '@/services/api/approvalWebSocket';
import type { ApprovalTaskResponse } from '@/interfaces/approval';

/**
 * Vue Composable สำหรับ Student - รับ notification เมื่อ Admin อัพเดตสถานะเอกสาร
 * 
 * Flow:
 * 1. Admin กด Approve/Reject/Request-Change
 * 2. Backend broadcast 'approval_task_updated' event
 * 3. Hook นี้รับ event และอัปเดต reactive ref
 * 4. Component ใช้ watch เพื่อ react ต่อการเปลี่ยนแปลง
 */
export function useStudentApprovalWebSocket() {
  const isConnected = ref(false);
  const lastUpdate: Ref<ApprovalTaskResponse | null> = ref(null);
  const updateCount = ref(0);

  const unsubscribers: (() => void)[] = [];

  const handleConnectionChange = (connected: boolean) => {
    isConnected.value = connected;
    // console.log(`🔌 [Student] Approval WebSocket ${connected ? 'connected ✅' : 'disconnected ❌'}`);
  };

  const handleTaskUpdated = (updatedTask: ApprovalTaskResponse) => {
    // console.log('📢 [Student] Received approval update:', updatedTask);
    lastUpdate.value = updatedTask;
    updateCount.value++;
  };

  onMounted(() => {
    // console.log('🚀 [Student] Mounting approval WebSocket hook...');

    // Connect to WebSocket (จะ reuse connection เดิมถ้ามีอยู่แล้ว)
    approvalWs.connect();

    // Subscribe to approval task updates (when admin approves/rejects documents)
    unsubscribers.push(approvalWs.onConnectionChange(handleConnectionChange));
    unsubscribers.push(approvalWs.onTaskUpdated(handleTaskUpdated));

    // Subscribe to application status updates (when evaluation is completed)
    unsubscribers.push(approvalWs.on('application_status_updated', (data: any) => {
      // console.log('📢 [Student] Received application status update (evaluation):', data);
      lastUpdate.value = data;
      updateCount.value++;
    }));

    // Subscribe to screening updates (when admin screens applicant)
    unsubscribers.push(approvalWs.on('screening_status_updated', (data: any) => {
      // console.log('📢 [Student] Received screening status update:', data);
      lastUpdate.value = data;
      updateCount.value++;
    }));

    // Subscribe to batch screening updates
    unsubscribers.push(approvalWs.on('screening_batch_updated', (data: any) => {
      // console.log('📢 [Student] Received batch screening update:', data);
      lastUpdate.value = data;
      updateCount.value++;
    }));

    // Subscribe to interview booking updates (when student/admin books interview)
    unsubscribers.push(approvalWs.on('interview_booking_created', (data: any) => {
      // console.log('📢 [Student] Received interview booking created:', data);
      lastUpdate.value = data;
      updateCount.value++;
    }));

    // Subscribe to interview booking deletion (when booking is cancelled)
    unsubscribers.push(approvalWs.on('interview_booking_deleted', (data: any) => {
      // console.log('📢 [Student] Received interview booking deleted:', data);
      lastUpdate.value = data;
      updateCount.value++;
    }));

    // Subscribe to evaluation created (when admin starts evaluation)
    unsubscribers.push(approvalWs.on('evaluation_created', (data: any) => {
      // console.log('📢 [Student] Received evaluation created:', data);
      lastUpdate.value = data;
      updateCount.value++;
    }));

    // Subscribe to evaluation updated
    unsubscribers.push(approvalWs.on('evaluation_updated', (data: any) => {
      // console.log('📢 [Student] Received evaluation updated:', data);
      lastUpdate.value = data;
      updateCount.value++;
    }));

    // Subscribe to evaluation completed
    unsubscribers.push(approvalWs.on('evaluation_completed', (data: any) => {
      // console.log('📢 [Student] Received evaluation completed:', data);
      lastUpdate.value = data;
      updateCount.value++;
    }));

    // Check current connection status
    isConnected.value = approvalWs.isConnected();
    // console.log(`🔌 [Student] Initial connection status: ${isConnected.value ? 'connected' : 'disconnected'}`);
  });

  onUnmounted(() => {
    // console.log('🧹 [Student] Unmounting approval WebSocket hook...');
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
