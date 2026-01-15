import { ref, onMounted, onUnmounted } from 'vue';
import { approvalWs } from '@/services/api/approvalWebSocket';
import { getApprovalTasks } from '@/services/api/approval';
import type { ApprovalTaskResponse } from '@/interfaces/approval';
import Swal from 'sweetalert2';

export function useApprovalWebSocket() {
    const tasks = ref<ApprovalTaskResponse[]>([]);
    const isConnected = ref(false);
    const error = ref<string | null>(null);

    const unsubscribers: (() => void)[] = [];

    const handleConnectionChange = (connected: boolean) => {
        isConnected.value = connected;
        if (connected) {
            fetchInitialData();
        }
    };
    
    const fetchInitialData = async () => {
        try {
            const initialTasks = await getApprovalTasks();
            tasks.value = initialTasks || [];
            console.log('Initial approval tasks loaded via hook.');
        } catch (e) {
            console.error('Failed to fetch initial approval tasks:', e);
            error.value = 'Failed to load initial approval tasks.';
            Swal.fire('Error', error.value, 'error');
        }
    };
    const handleTaskUpdated = (updatedTask: ApprovalTaskResponse) => {
        const index = tasks.value.findIndex(t => t.ID === updatedTask.ID);
        if (index !== -1) {
            tasks.value.splice(index, 1, updatedTask);
        } else {
            // If the task is new (e.g. from a new document upload), add it.
            tasks.value.unshift(updatedTask);
        }
    };

    onMounted(() => {
        approvalWs.connect();

        unsubscribers.push(approvalWs.onConnectionChange(handleConnectionChange));
        unsubscribers.push(approvalWs.onTaskUpdated(handleTaskUpdated));
        
        isConnected.value = approvalWs.isConnected();
        if (isConnected.value) {
            fetchInitialData();
        }
    });

    onUnmounted(() => {
        unsubscribers.forEach(unsub => unsub());
        // We might not want to disconnect globally if another component is using it.
        // For now, we'll leave the connection open.
        // approvalWs.disconnect(); 
    });

    return {
        tasks,
        isConnected,
        error,
    };
}
