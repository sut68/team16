<script setup lang="ts">
  import { ref, onMounted } from 'vue';
  import { InterviewAPI } from '@/services/api';
  import type { InterviewRound } from '@/interfaces/interview';
  import { InterviewRoundCriteriaService } from '@/services/evaluation/evaluation';
  import RoundCriteriaModal from './RoundCriteriaModal.vue';
  import { ChevronRight, Calendar, Users, CheckCircle } from 'lucide-vue-next';

  // State
  const isLoading = ref(false);
  const error = ref<string | null>(null);
  const interviewRounds = ref<InterviewRound[]>([]);
  const roundCriteriaCounts = ref<Record<number, number>>({});

  // Modal State
  const isCriteriaModalOpen = ref(false);
  const selectedRound = ref<{ id: number; name: string } | null>(null);

  // Fetch interview rounds
  async function fetchRounds() {
    isLoading.value = true;
    error.value = null;
    
    try {
      const rounds = await InterviewAPI.getAllRounds();
      interviewRounds.value = rounds || [];
      
      // Fetch criteria count for each round
      const counts: Record<number, number> = {};
      await Promise.all(
        interviewRounds.value.map(async (round) => {
          try {
            const criteria = await InterviewRoundCriteriaService.getByRoundId(round.ID);
            counts[round.ID] = criteria?.length || 0;
          } catch {
            counts[round.ID] = 0;
          }
        })
      );
      roundCriteriaCounts.value = counts;
    } catch (err: any) {
      console.error('Failed to fetch interview rounds:', err);
      error.value = err?.message || 'ไม่สามารถโหลดข้อมูลรอบสัมภาษณ์ได้';
    } finally {
      isLoading.value = false;
    }
  }

  // Open criteria modal
  function openCriteriaModal(round: InterviewRound) {
    selectedRound.value = {
      id: round.ID,
      name: round.name,
    };
    isCriteriaModalOpen.value = true;
  }

  // Close criteria modal
  function closeCriteriaModal() {
    isCriteriaModalOpen.value = false;
    selectedRound.value = null;
  }

  // Handle criteria updated
  async function handleCriteriaUpdated() {
    // Refresh criteria counts
    if (selectedRound.value) {
      try {
        const criteria = await InterviewRoundCriteriaService.getByRoundId(selectedRound.value.id);
        roundCriteriaCounts.value[selectedRound.value.id] = criteria?.length || 0;
      } catch {
        // Ignore error
      }
    }
  }

  // Format date
  function formatDate(dateString: string): string {
    if (!dateString) return '';
    const date = new Date(dateString);
    return date.toLocaleDateString('th-TH', {
      day: 'numeric',
      month: 'short',
      year: 'numeric',
    });
  }

  // Get round status
  function getRoundStatus(round: InterviewRound): 'active' | 'upcoming' | 'closed' {
    const now = new Date();
    const start = new Date(round.start_date_time);
    const end = new Date(round.end_date_time);
    
    if (now < start) return 'upcoming';
    if (now > end) return 'closed';
    return 'active';
  }

  function getStatusLabel(status: string): string {
    const labels: Record<string, string> = {
      active: 'กำลังดำเนินการ',
      upcoming: 'กำลังจะมาถึง',
      closed: 'ปิดแล้ว',
    };
    return labels[status] || status;
  }

  function getStatusClass(status: string): string {
    const classes: Record<string, string> = {
      active: 'bg-green-100 text-green-700 border-green-200',
      upcoming: 'bg-blue-100 text-blue-700 border-blue-200',
      closed: 'bg-gray-100 text-gray-500 border-gray-200',
    };
    return classes[status] || 'bg-gray-100 text-gray-500';
  }

  onMounted(fetchRounds);
</script>

<template>
  <div class="h-full flex flex-col p-6 overflow-auto">
    <!-- Header -->
    <div class="flex items-center justify-between mb-6">
      <div class="flex items-center gap-3">
        <h1 class="text-2xl font-bold text-[#1e3a8a]">เกณฑ์ประจำรอบสัมภาษณ์</h1>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="isLoading" class="flex-1 flex items-center justify-center">
      <span class="loading loading-spinner loading-lg text-purple-500"></span>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="flex-1 flex flex-col items-center justify-center">
      <div class="bg-red-50 rounded-xl p-6 text-center max-w-md">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12 text-red-400 mx-auto mb-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
        </svg>
        <p class="text-red-600 font-medium mb-2">{{ error }}</p>
        <button @click="fetchRounds" class="text-sm text-red-500 hover:text-red-700 underline">
          ลองใหม่อีกครั้ง
        </button>
      </div>
    </div>

    <!-- Empty State -->
    <div v-else-if="interviewRounds.length === 0" class="flex-1 flex flex-col items-center justify-center">
      <div class="bg-gray-50 rounded-xl p-8 text-center max-w-md">
        <Calendar class="h-16 w-16 text-gray-300 mx-auto mb-4" />
        <h3 class="text-lg font-medium text-gray-600 mb-2">ไม่มีรอบสัมภาษณ์</h3>
        <p class="text-gray-500 text-sm">ยังไม่มีรอบสัมภาษณ์ในระบบ กรุณาสร้างรอบสัมภาษณ์ก่อน</p>
      </div>
    </div>

    <!-- Rounds List -->
    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      <div
        v-for="round in interviewRounds"
        :key="round.ID"
        @click="openCriteriaModal(round)"
        class="group bg-white rounded-xl border border-gray-200 p-5 cursor-pointer
               transition-all duration-200 hover:shadow-lg hover:border-purple-300 hover:-translate-y-1"
      >
        <!-- Header -->
        <div class="flex items-start justify-between mb-3">
          <div class="flex-1 min-w-0">
            <h3 class="font-semibold text-gray-800 truncate group-hover:text-purple-600 transition-colors">
              {{ round.name }}
            </h3>
            <p v-if="round.scholarship" class="text-sm text-gray-500 truncate mt-1">
              {{ round.scholarship.scholarship_name }}
            </p>
          </div>
          <span 
            class="badge badge-sm border ml-2 flex-shrink-0"
            :class="getStatusClass(getRoundStatus(round))"
          >
            {{ getStatusLabel(getRoundStatus(round)) }}
          </span>
        </div>

        <!-- Info -->
        <div class="space-y-2 text-sm text-gray-600 mb-4">
          <div class="flex items-center gap-2">
            <Calendar class="h-4 w-4 text-gray-400" />
            <span>{{ formatDate(round.start_date_time) }} - {{ formatDate(round.end_date_time) }}</span>
          </div>
          <div class="flex items-center gap-2">
            <Users class="h-4 w-4 text-gray-400" />
            <span>{{ round.slots?.length || 0 }} สล็อต</span>
          </div>
        </div>

        <!-- Footer -->
        <div class="flex items-center justify-between pt-3 border-t border-gray-100">
          <div class="flex items-center gap-2">
            <CheckCircle class="h-4 w-4 text-purple-500" />
            <span class="text-sm font-medium text-purple-600">
              {{ roundCriteriaCounts[round.ID] || 0 }} เกณฑ์
            </span>
          </div>
          <ChevronRight class="h-5 w-5 text-gray-400 group-hover:text-purple-500 transition-colors" />
        </div>
      </div>
    </div>

    <!-- Criteria Modal -->
    <RoundCriteriaModal
      :is-open="isCriteriaModalOpen"
      :round-id="selectedRound?.id || 0"
      :round-name="selectedRound?.name || ''"
      @close="closeCriteriaModal"
      @updated="handleCriteriaUpdated"
    />
  </div>
</template>

<style scoped>
/* Custom scrollbar */
::-webkit-scrollbar {
  width: 6px;
}
::-webkit-scrollbar-track {
  background: transparent;
}
::-webkit-scrollbar-thumb {
  background: #cbd5e1;
  border-radius: 10px;
}
::-webkit-scrollbar-thumb:hover {
  background: #94a3b8;
}
</style>
