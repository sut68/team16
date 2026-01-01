<script setup lang="ts">
  import { ref } from 'vue'
  import { ClipboardList, Settings2, ListChecks } from 'lucide-vue-next'

  // Tab Components
  import EvaluationList from './EvaluationList.vue'
  import EvaluationCriteriaList from './EvaluationCriteriaList.vue'

  // Tab State
  type TabType = 'evaluations' | 'criteria'
  const activeTab = ref<TabType>('evaluations')

  const tabs = [
    { 
      id: 'evaluations', 
      label: 'รายการประเมิน', 
      icon: ClipboardList,
      description: 'ประเมินและตัดสินผู้สมัครทุน'
    },
    { 
      id: 'criteria', 
      label: 'เกณฑ์การประเมิน', 
      icon: Settings2,
      description: 'จัดการเกณฑ์และน้ำหนักคะแนน'
    },
  ] as const
</script>

<template>
  <div class="flex flex-col h-full bg-white rounded-tl-[30px] shadow overflow-hidden">
    <!-- Page Header with Tabs -->
    <div class="bg-white border-b border-gray-200 shadow-sm rounded-tl-[30px]">
      <div class="px-6 pt-6">
        <!-- Title -->
        <div class="flex items-center gap-3 mb-4">
          <div class="w-10 h-10 bg-gradient-to-br from-indigo-500 to-blue-600 rounded-xl flex items-center justify-center shadow-lg">
            <ListChecks class="w-5 h-5 text-white" />
          </div>
          <div>
            <h1 class="text-2xl font-bold text-gray-900">ระบบพิจารณาผู้รับทุน</h1>
            <p class="text-sm text-gray-500">ประเมินและคัดเลือกผู้ได้รับทุนการศึกษา</p>
          </div>
        </div>

        <!-- Tab Navigation -->
        <nav class="flex gap-1" aria-label="Tabs">
          <button
            v-for="tab in tabs"
            :key="tab.id"
            @click="activeTab = tab.id"
            class="group relative px-5 py-3 rounded-t-xl font-medium text-sm transition-all duration-200 flex items-center gap-2"
            :class="[
              activeTab === tab.id
                ? 'bg-white text-indigo-600 shadow-sm border border-b-0 border-gray-200'
                : 'text-gray-500 hover:text-gray-700 hover:bg-gray-50'
            ]"
          >
            <component 
              :is="tab.icon" 
              class="w-4 h-4 transition-colors"
              :class="activeTab === tab.id ? 'text-indigo-600' : 'text-gray-400 group-hover:text-gray-500'"
            />
            {{ tab.label }}
            
            <!-- Active Indicator -->
            <span
              v-if="activeTab === tab.id"
              class="absolute bottom-0 left-0 right-0 h-0.5 bg-indigo-600 rounded-t"
            ></span>
          </button>
        </nav>
      </div>
    </div>

    <!-- Tab Content -->
    <div class="flex-1 min-h-0 overflow-hidden">
      <Transition name="fade" mode="out-in">
        <KeepAlive>
          <component 
            :is="activeTab === 'evaluations' ? EvaluationList : EvaluationCriteriaList"
            :key="activeTab"
          />
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
