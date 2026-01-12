<script setup lang="ts">
  import { ref, onMounted, provide } from 'vue'
  import { ClipboardList, Settings2, Layers, UserPlus } from 'lucide-vue-next'
  import { getAllApplicationScholarships } from '@/services/api/application'

  // Tab Components
  import EvaluationList from './EvaluationList.vue'
  import EvaluationCriteriaList from './EvaluationCriteriaList.vue'
  import RoundCriteriaManager from './RoundCriteriaManager.vue'
  import QualifiedApplicantsTab from './QualifiedApplicantsTab.vue'

  // Qualified Applicants Count (for badge)
  const qualifiedCount = ref(0)

  // Fetch count on mount
  async function fetchQualifiedCount() {
    try {
      const [qualified, scheduled] = await Promise.all([
        getAllApplicationScholarships('qualified'),
        getAllApplicationScholarships('interview_scheduled')
      ])
      qualifiedCount.value = qualified.length + scheduled.length
    } catch (e) {
      console.error('Error fetching qualified count:', e)
    }
  }

  // Provide a refresh function for child components
  function refreshQualifiedCount() {
    fetchQualifiedCount()
  }
  provide('refreshQualifiedCount', refreshQualifiedCount)

  onMounted(fetchQualifiedCount)

  // Tab State
  type TabType = 'evaluations' | 'qualified' | 'criteria' | 'roundCriteria'
  const activeTab = ref<TabType>('evaluations')

  const tabs: { id: TabType; label: string; icon: any; description: string; badge?: boolean }[] = [
    { 
      id: 'evaluations', 
      label: 'การประเมิน', 
      icon: ClipboardList,
      description: 'ประเมินและตัดสินผู้สมัครทุน',
      badge: false
    },
    { 
      id: 'qualified', 
      label: 'ผู้สมัครพร้อมประเมิน', 
      icon: UserPlus,
      description: 'สร้างการประเมินสำหรับผู้สมัครใหม่',
      badge: true
    },
    { 
      id: 'criteria', 
      label: 'เกณฑ์การประเมิน', 
      icon: Settings2,
      description: 'จัดการเกณฑ์และน้ำหนักคะแนน',
      badge: false
    },
    { 
      id: 'roundCriteria', 
      label: 'เกณฑ์ประจำรอบ', 
      icon: Layers,
      description: 'กำหนดเกณฑ์สำหรับแต่ละรอบสัมภาษณ์',
      badge: false
    },
  ]
</script>

<template>
  <div class="flex flex-col h-full bg-white rounded-tl-[30px] shadow overflow-hidden">
    <!-- Page Header with Tabs -->
    <div class="px-6 pt-6">
      <!-- Title -->
      <div class="mb-6">
        <h1 class="text-2xl font-bold text-slate-800">พิจารณาผู้รับทุน</h1>
      </div>

      <!-- Tab Navigation (InterviewRoundManager style) -->
      <div class="flex gap-8 border-b border-gray-200">
        <a
          v-for="tab in tabs"
          :key="tab.id"
          @click="activeTab = tab.id"
          class="pb-3 px-1 text-base font-medium cursor-pointer transition-all border-b-[3px] -mb-[1px] flex items-center gap-2 relative"
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
          <!-- Badge for qualified tab -->
          <span 
            v-if="tab.badge && qualifiedCount > 0 && activeTab !== tab.id"
            class="ml-1 min-w-5 h-5 px-1.5 bg-red-500 text-white text-xs font-bold rounded-full flex items-center justify-center animate-pulse"
          >
            {{ qualifiedCount }}
          </span>
          <span 
            v-else-if="tab.badge"
            class="ml-1 text-xs opacity-60"
          >
            ({{ qualifiedCount }})
          </span>
        </a>
      </div>
    </div>

    <!-- Tab Content -->
    <div class="flex-1 min-h-0 overflow-hidden">
      <Transition name="fade" mode="out-in">
        <KeepAlive>
          <EvaluationList v-if="activeTab === 'evaluations'" :key="'evaluations'" />
          <QualifiedApplicantsTab v-else-if="activeTab === 'qualified'" :key="'qualified'" @created="fetchQualifiedCount" />
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
