import { ref, onMounted, onUnmounted } from 'vue';
import { interviewWs } from '@/services/api/interviewWebSocket';
import type { InterviewRound, Slot, InterviewBookingDTOResponse } from '@/interfaces/interview';
import { InterviewAPI } from '@/services/api';
import Swal from 'sweetalert2';


export function useInterviewWebSocket() {
    const rounds = ref<InterviewRound[]>([]);
    const isConnected = ref(false);
    const error = ref<string | null>(null);

    // Unsubscribe functions
    const unsubscribers: (() => void)[] = [];

    const handleConnectionChange = (connected: boolean) => {
        isConnected.value = connected;
        if (connected) {
            // Fetch initial data upon connection
            fetchInitialData();
        }
    };
    
    const fetchInitialData = async () => {
        try {
            const initialRounds = await InterviewAPI.getAllRounds();
            rounds.value = initialRounds || [];
            console.log('Initial interview rounds loaded via hook.');
        } catch (e) {
            console.error('Failed to fetch initial interview rounds:', e);
            error.value = 'Failed to load initial data.';
            Swal.fire('Error', 'Failed to load initial interview data.', 'error');
        }
    };

    const handleRoundCreated = (newRound: InterviewRound) => {
        rounds.value.push(newRound);
    };
    
    const handleRoundUpdated = (updatedRound: InterviewRound) => {
        const index = rounds.value.findIndex(r => r.ID === updatedRound.ID);
        if (index !== -1) {
            rounds.value.splice(index, 1, updatedRound);
        } else {
            rounds.value.push(updatedRound); // If not found, add it
        }
    };

    const handleRoundDeleted = (data: { id: number }) => {
        rounds.value = rounds.value.filter(r => r.ID !== data.id);
    };

    const handleSlotUpdated = (updatedSlot: Slot) => {
        for (const round of rounds.value) {
            const index = round.slots.findIndex((s: Slot) => s.ID === updatedSlot.ID);
            if (index !== -1) {
                round.slots[index] = updatedSlot;
                break;
            }
        }
    };

    const refetchRound = async (roundId: number) => {
        try {
            const updatedRound = await InterviewAPI.getRoundById(roundId);
            if (updatedRound) {
                const index = rounds.value.findIndex(r => r.ID === updatedRound.ID);
                if (index !== -1) {
                    rounds.value.splice(index, 1, updatedRound);
                }
            }
        } catch (e) {
            console.error(`Failed to refetch round ${roundId} after WebSocket update`, e);
        }
    };

    const handleBookingCreatedOrUpdated = (bookingDto: InterviewBookingDTOResponse) => {
        refetchRound(bookingDto.interview_round_id);
    };

    const handleBookingDeleted = (bookingDto: InterviewBookingDTOResponse) => {
        refetchRound(bookingDto.interview_round_id);
    };

    onMounted(() => {
        interviewWs.connect();

        unsubscribers.push(interviewWs.onConnectionChange(handleConnectionChange));
        unsubscribers.push(interviewWs.onRoundCreated(handleRoundCreated));
        unsubscribers.push(interviewWs.onRoundUpdated(handleRoundUpdated));
        unsubscribers.push(interviewWs.onRoundDeleted(handleRoundDeleted));
        unsubscribers.push(interviewWs.onSlotUpdated(handleSlotUpdated));
        unsubscribers.push(interviewWs.onBookingCreated(handleBookingCreatedOrUpdated));
        unsubscribers.push(interviewWs.onBookingUpdated(handleBookingCreatedOrUpdated));
        unsubscribers.push(interviewWs.onBookingDeleted(handleBookingDeleted));
        
        isConnected.value = interviewWs.isConnected();
        if (isConnected.value) {
             fetchInitialData();
        }
    });

    onUnmounted(() => {
        unsubscribers.forEach(unsub => unsub());
        // Optional: decide if the connection should be disconnected when no component is using it.
        // interviewWs.disconnect();
    });

    return {
        rounds,
        isConnected,
        error,
    };
}
