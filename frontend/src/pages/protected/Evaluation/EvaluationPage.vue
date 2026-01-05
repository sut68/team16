<script setup lang="ts">
  import { ref } from 'vue'
  import { ClipboardList, Settings2, Layers } from 'lucide-vue-next'

  // Tab Components
  import EvaluationList from './EvaluationList.vue'
  import EvaluationCriteriaList from './EvaluationCriteriaList.vue'
  import RoundCriteriaManager from './RoundCriteriaManager.vue'

  // Tab State
  type TabType = 'evaluations' | 'criteria' | 'roundCriteria'
  const activeTab = ref<TabType>('evaluations')

  const tabs = [
    { 
      id: 'evaluations', 
      label: 'การประเมิน', 
      icon: ClipboardList,
      description: 'ประเมินและตัดสินผู้สมัครทุน'
    },
    { 
      id: 'criteria', 
      label: 'เกณฑ์การประเมิน', 
      icon: Settings2,
      description: 'จัดการเกณฑ์และน้ำหนักคะแนน'
    },
    { 
      id: 'roundCriteria', 
      label: 'เกณฑ์ประจำรอบ', 
      icon: Layers,
      description: 'กำหนดเกณฑ์สำหรับแต่ละรอบสัมภาษณ์'
    },
  ] as const
</script>

<template>
  <div class="flex flex-col h-full bg-white rounded-tl-[30px] shadow overflow-hidden">
    <!-- Page Header with Tabs -->
    <div class="px-6 pt-6">
      <!-- Title -->
      <div class="mb-6">
        <h1 class="text-2xl font-bold text-gray-900">พิจารณาผู้รับทุน</h1>
      </div>

      <!-- Tab Navigation (InterviewRoundManager style) -->
      <div class="flex gap-8 border-b border-gray-200">
        <a
          v-for="tab in tabs"
          :key="tab.id"
          @click="activeTab = tab.id"
          class="pb-3 px-1 text-base font-medium cursor-pointer transition-all border-b-[3px] -mb-[1px] flex items-center gap-2"
          :class="activeTab === tab.id 
            ? 'text-[#1e3a8a] border-[#1e3a8a]' 
            : 'text-slate-500 border-transparent hover:text-slate-700'"
        >
          <component 
            :is="tab.icon" 
            class="w-4 h-4"
            :class="activeTab === tab.id ? 'text-[#1e3a8a]' : 'text-slate-400'"
          />
          {{ tab.label }}
        </a>
      </div>
    </div>

    <!-- Tab Content -->
    <div class="flex-1 min-h-0 overflow-hidden">
      <Transition name="fade" mode="out-in">
        <KeepAlive>
          <EvaluationList v-if="activeTab === 'evaluations'" :key="'evaluations'" />
          <EvaluationCriteriaList v-else-if="activeTab === 'criteria'" :key="'criteria'" />
          <RoundCriteriaManager v-else :key="'roundCriteria'" />
        </KeepAlive>
      </Transition>
    </div>
  </div>
</template>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.15s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
