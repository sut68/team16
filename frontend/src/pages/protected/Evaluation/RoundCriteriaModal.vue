<script setup lang="ts">
  import { ref, computed, watch } from 'vue';
  import type { EvaluationCriterionResponse, InterviewRoundCriterionResponse } from '@/interfaces/evaluation';
  import { EvaluationCriteriaService, InterviewRoundCriteriaService } from '@/services/evaluation/evaluation';
  import Swal from 'sweetalert2';
  import { Plus } from 'lucide-vue-next';

  // Props
  const props = defineProps<{
    isOpen: boolean;
    roundId: number;
    roundName: string;
  }>();

  // Emits
  const emit = defineEmits<{
    (e: 'close'): void;
    (e: 'updated'): void;
  }>();

  // State
  const isLoading = ref(false);
  const allCriteria = ref<EvaluationCriterionResponse[]>([]);
  const roundCriteria = ref<InterviewRoundCriterionResponse[]>([]);
  const showAddDropdown = ref(false);

  // Computed: เกณฑ์ที่ยังไม่ได้เพิ่มในรอบนี้
  const availableCriteria = computed(() => {
    const usedIds = roundCriteria.value.map(rc => rc.evaluation_criterion_id);
    return allCriteria.value.filter(c => c.is_active && !usedIds.includes(c.ID));
  });

  // Watch: โหลดข้อมูลเมื่อ modal เปิด
  watch(() => props.isOpen, async (isOpen) => {
    if (isOpen && props.roundId) {
      await fetchData();
    }
  });

  // Functions
  async function fetchData() {
    isLoading.value = true;
    try {
      const [criteriaRes, roundCriteriaRes] = await Promise.all([
        EvaluationCriteriaService.getAll(),
        InterviewRoundCriteriaService.getByRoundId(props.roundId),
      ]);
      allCriteria.value = criteriaRes || [];
      roundCriteria.value = roundCriteriaRes || [];
    } catch (error) {
      console.error('Failed to fetch criteria:', error);
      Swal.fire('Error', 'ไม่สามารถโหลดข้อมูลเกณฑ์ได้', 'error');
    } finally {
      isLoading.value = false;
    }
  }

  async function addCriterion(criterionId: number) {
    try {
      const newRoundCriterion = await InterviewRoundCriteriaService.addToRound(props.roundId, {
        evaluation_criterion_id: criterionId,
        weight: 1.0,
        is_enabled: true,
      });
      roundCriteria.value.push(newRoundCriterion);
      showAddDropdown.value = false;
      emit('updated');
      
      Swal.fire({
        toast: true,
        position: 'top-end',
        icon: 'success',
        title: 'เพิ่มเกณฑ์เรียบร้อย',
        showConfirmButton: false,
        timer: 1500,
      });
    } catch (error) {
      console.error('Failed to add criterion:', error);
      Swal.fire('Error', 'ไม่สามารถเพิ่มเกณฑ์ได้', 'error');
    }
  }

  async function removeCriterion(roundCriterionId: number) {
    const result = await Swal.fire({
      title: 'ยืนยันการลบ?',
      text: 'ต้องการลบเกณฑ์นี้ออกจากรอบสัมภาษณ์หรือไม่?',
      icon: 'warning',
      showCancelButton: true,
      confirmButtonColor: '#d33',
      cancelButtonColor: '#3085d6',
      confirmButtonText: 'ลบ',
      cancelButtonText: 'ยกเลิก',
    });

    if (result.isConfirmed) {
      try {
        await InterviewRoundCriteriaService.removeFromRound(roundCriterionId);
        roundCriteria.value = roundCriteria.value.filter(rc => rc.ID !== roundCriterionId);
        emit('updated');
        
        Swal.fire({
          toast: true,
          position: 'top-end',
          icon: 'success',
          title: 'ลบเกณฑ์เรียบร้อย',
          showConfirmButton: false,
          timer: 1500,
        });
      } catch (error) {
        console.error('Failed to remove criterion:', error);
        Swal.fire('Error', 'ไม่สามารถลบเกณฑ์ได้', 'error');
      }
    }
  }

  async function updateCriterion(roundCriterion: InterviewRoundCriterionResponse, updates: { weight?: number; is_enabled?: boolean }) {
    try {
      const updated = await InterviewRoundCriteriaService.update(roundCriterion.ID, updates);
      const index = roundCriteria.value.findIndex(rc => rc.ID === roundCriterion.ID);
      if (index !== -1) {
        roundCriteria.value[index] = updated;
      }
      emit('updated');
    } catch (error) {
      console.error('Failed to update criterion:', error);
      Swal.fire('Error', 'ไม่สามารถอัปเดตเกณฑ์ได้', 'error');
    }
  }

  function toggleEnabled(roundCriterion: InterviewRoundCriterionResponse) {
    updateCriterion(roundCriterion, { is_enabled: !roundCriterion.is_enabled });
  }

  function handleWeightChange(roundCriterion: InterviewRoundCriterionResponse, newWeight: number) {
    if (newWeight >= 0 && newWeight <= 10) {
      updateCriterion(roundCriterion, { weight: newWeight });
    }
  }

  function closeModal() {
    showAddDropdown.value = false;
    emit('close');
  }

  // Score type labels
  function getScoreTypeLabel(type: string): string {
    const labels: Record<string, string> = {
      numeric: 'คะแนน',
      grade: 'เกรด',
      pass_fail: 'ผ่าน/ไม่ผ่าน',
    };
    return labels[type] || type;
  }

  function getScoreTypeBadgeClass(type: string): string {
    const classes: Record<string, string> = {
      numeric: 'bg-blue-100 text-blue-700',
      grade: 'bg-purple-100 text-purple-700',
      pass_fail: 'bg-green-100 text-green-700',
    };
    return classes[type] || 'bg-gray-100 text-gray-700';
  }
</script>

<template>
  <!-- Modal Backdrop -->
  <div 
    v-if="isOpen"
    class="fixed inset-0 z-[110] flex items-center justify-center bg-black/60 backdrop-blur-sm p-4"
    @click.self="closeModal"
  >
    <!-- Modal Content -->
    <div class="bg-white w-full max-w-3xl max-h-[85vh] rounded-2xl shadow-2xl flex flex-col overflow-hidden animate-pop-in">
      
      <!-- Header -->
      <div class="px-6 py-4 border-b flex items-center justify-between bg-gradient-to-r from-[#1e3a8a] to-[#3b5998]">
        <div>
          <h2 class="text-xl font-bold text-white flex items-center gap-2">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01" />
            </svg>
            จัดการเกณฑ์การประเมิน
          </h2>
          <p class="text-white/70 text-sm mt-1">{{ roundName }}</p>
        </div>
        <button 
          @click="closeModal"
          class="btn btn-circle btn-ghost btn-sm text-white hover:bg-white/20"
        >
          ✕
        </button>
      </div>

      <!-- Content -->
      <div class="flex-1 overflow-y-auto p-6 bg-gray-50">
        
        <!-- Loading State -->
        <div v-if="isLoading" class="flex justify-center items-center h-40">
          <span class="loading loading-spinner loading-lg text-[#1e3a8a]"></span>
        </div>

        <div v-else>
          <!-- Add Criterion Section -->
          <div class="mb-6">
            <div class="flex items-center justify-between mb-3">
              <h3 class="font-semibold text-slate-700">เกณฑ์ที่ใช้ในรอบนี้</h3>
              
              <!-- Add Button with Dropdown -->
              <div class="relative">
                <button 
                  @click="showAddDropdown = !showAddDropdown"
                  class="btn-ghost-rounded"
                  :disabled="availableCriteria.length === 0"
                >
                  <Plus class="w-4 h-4" />
                  <span class="font-medium">เพิ่มเกณฑ์</span>
                </button>

                <!-- Dropdown -->
                <div 
                  v-if="showAddDropdown && availableCriteria.length > 0"
                  class="absolute right-0 mt-2 w-72 bg-white rounded-xl shadow-xl border z-10 max-h-60 overflow-y-auto"
                >
                  <div class="p-2">
                    <p class="text-xs text-gray-500 px-3 py-2 border-b">เลือกเกณฑ์ที่ต้องการเพิ่ม</p>
                    <div 
                      v-for="criterion in availableCriteria" 
                      :key="criterion.ID"
                      @click="addCriterion(criterion.ID)"
                      class="px-3 py-2 hover:bg-blue-50 cursor-pointer rounded-lg flex items-center justify-between gap-2 transition-colors"
                    >
                      <div>
                        <p class="font-medium text-slate-700 text-sm">{{ criterion.name }}</p>
                        <p class="text-xs text-gray-500">คะแนนเต็ม: {{ criterion.max_score }}</p>
                      </div>
                      <span 
                        class="badge badge-sm"
                        :class="getScoreTypeBadgeClass(criterion.score_type)"
                      >
                        {{ getScoreTypeLabel(criterion.score_type) }}
                      </span>
                    </div>
                  </div>
                </div>

                <!-- Backdrop for dropdown -->
                <div 
                  v-if="showAddDropdown" 
                  class="fixed inset-0 z-0" 
                  @click="showAddDropdown = false"
                ></div>
              </div>
            </div>

            <!-- Empty State -->
            <div 
              v-if="roundCriteria.length === 0"
              class="bg-white border-2 border-dashed border-gray-200 rounded-xl p-8 text-center"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12 mx-auto text-gray-300 mb-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
              </svg>
              <p class="text-gray-500">ยังไม่มีเกณฑ์การประเมินในรอบนี้</p>
              <p class="text-sm text-gray-400 mt-1">กดปุ่ม "เพิ่มเกณฑ์" เพื่อเริ่มต้น</p>
            </div>

            <!-- Criteria List -->
            <div v-else class="space-y-3">
              <div 
                v-for="rc in roundCriteria" 
                :key="rc.ID"
                class="bg-white rounded-xl border border-gray-200 p-4 transition-all hover:shadow-md"
                :class="{ 'opacity-50': !rc.is_enabled }"
              >
                <div class="flex items-center justify-between gap-4">
                  <!-- Left: Info -->
                  <div class="flex-1 min-w-0">
                    <div class="flex items-center gap-2 mb-1">
                      <h4 class="font-semibold text-slate-700 truncate">
                        {{ rc.evaluation_criterion?.name || 'Unknown' }}
                      </h4>
                      <span 
                        class="badge badge-sm"
                        :class="getScoreTypeBadgeClass(rc.evaluation_criterion?.score_type || '')"
                      >
                        {{ getScoreTypeLabel(rc.evaluation_criterion?.score_type || '') }}
                      </span>
                    </div>
                    <p class="text-sm text-gray-500 truncate">
                      {{ rc.evaluation_criterion?.description || 'ไม่มีคำอธิบาย' }}
                    </p>
                    <div class="text-xs text-gray-400 mt-1">
                      คะแนนเต็ม: {{ rc.evaluation_criterion?.max_score || 0 }}
                    </div>
                  </div>

                  <!-- Right: Controls -->
                  <div class="flex items-center gap-3">
                    <!-- Weight Input -->
                    <div class="flex flex-col items-center">
                      <label class="text-xs text-gray-500 mb-1">น้ำหนัก</label>
                      <input 
                        type="number"
                        :value="rc.weight"
                        @change="(e) => handleWeightChange(rc, parseFloat((e.target as HTMLInputElement).value))"
                        min="0"
                        max="10"
                        step="0.1"
                        class="input input-bordered input-sm w-16 text-center"
                      />
                    </div>

                    <!-- Toggle Enabled -->
                    <div class="flex flex-col items-center">
                      <label class="text-xs text-gray-500 mb-1">เปิดใช้</label>
                      <input 
                        type="checkbox"
                        :checked="rc.is_enabled"
                        @change="toggleEnabled(rc)"
                        class="toggle toggle-primary toggle-sm"
                      />
                    </div>

                    <!-- Delete Button -->
                    <button 
                      @click="removeCriterion(rc.ID)"
                      class="btn btn-ghost btn-sm text-red-500 hover:bg-red-50"
                    >
                      <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                      </svg>
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Info Box -->
          <div class="bg-blue-50 border border-blue-100 rounded-xl p-4 mt-4">
            <div class="flex gap-3">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-blue-500 flex-shrink-0 mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <div class="text-sm text-blue-700">
                <p class="font-medium mb-1">คำแนะนำ</p>
                <ul class="list-disc list-inside text-blue-600 space-y-1">
                  <li>เกณฑ์ที่เพิ่มจะปรากฏในแบบประเมินของผู้สมัครทุกคนในรอบนี้</li>
                  <li>น้ำหนัก (Weight) ใช้ในการคำนวณคะแนนรวม</li>
                  <li>สามารถปิดการใช้งานเกณฑ์ชั่วคราวได้โดยไม่ต้องลบ</li>
                </ul>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Footer -->
      <div class="p-4 border-t bg-white flex justify-end">
        <button 
          @click="closeModal"
          class="btn bg-[#1e3a8a] text-white hover:bg-[#152c6f]"
        >
          เสร็จสิ้น
        </button>
      </div>

    </div>
  </div>
</template>

<style scoped>
.animate-pop-in {
  animation: pop-in 0.2s cubic-bezier(0.16, 1, 0.3, 1) forwards;
}

@keyframes pop-in {
  0% {
    opacity: 0;
    transform: scale(0.95);
  }
  100% {
    opacity: 1;
    transform: scale(1);
  }
}

.btn-ghost-rounded {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  padding-left: 1rem;
  padding-right: 1rem;
  height: 2.5rem;
  font-size: 0.875rem;
  font-weight: 500;
  color: #374151;
  background-color: white;
  border: 1px solid #d1d5db;
  border-radius: 9999px;
  box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
  transition: all 150ms ease;
  cursor: pointer;
}

.btn-ghost-rounded:hover:not(:disabled) {
  background-color: #f3f4f6;
  border-color: #9ca3af;
}

.btn-ghost-rounded:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>
