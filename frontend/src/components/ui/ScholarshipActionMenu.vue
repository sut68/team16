<script setup lang="ts">
import { ref, watch, onUnmounted, nextTick } from 'vue'
import { MoreVertical, Eye, Pencil, Settings2, Trash2 } from 'lucide-vue-next'

// Props
const props = defineProps<{
  scholarshipId: number
  scholarshipName?: string
}>()

// Emits
const emit = defineEmits<{
  view: [id: number]
  edit: [id: number]
  manageFeatures: [id: number]
  delete: [id: number]
  menuOpen: []
  menuClose: []
}>()

// State
const isOpen = ref(false)
const menuRef = ref<HTMLElement | null>(null)

// Close menu
function close() {
  if (isOpen.value) {
    isOpen.value = false
    emit('menuClose')
  }
}

// Toggle menu
function toggle(event: Event) {
  event.stopPropagation()
  
  if (isOpen.value) {
    close()
  } else {
    isOpen.value = true
    emit('menuOpen')
    // Add click listener after a tick to avoid immediate close
    nextTick(() => {
      document.addEventListener('click', handleClickOutside, true)
    })
  }
}

// Handle click outside
function handleClickOutside(event: MouseEvent) {
  if (menuRef.value && !menuRef.value.contains(event.target as Node)) {
    close()
    document.removeEventListener('click', handleClickOutside, true)
  }
}

// Action handlers
function handleView(event: Event) {
  event.stopPropagation()
  emit('view', props.scholarshipId)
  close()
  document.removeEventListener('click', handleClickOutside, true)
}

function handleEdit(event: Event) {
  event.stopPropagation()
  emit('edit', props.scholarshipId)
  close()
  document.removeEventListener('click', handleClickOutside, true)
}

function handleManageFeatures(event: Event) {
  event.stopPropagation()
  emit('manageFeatures', props.scholarshipId)
  close()
  document.removeEventListener('click', handleClickOutside, true)
}

function handleDelete(event: Event) {
  event.stopPropagation()
  emit('delete', props.scholarshipId)
  close()
  document.removeEventListener('click', handleClickOutside, true)
}

// Cleanup on unmount
onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside, true)
})
</script>

<template>
  <div class="relative" ref="menuRef">
    <!-- Trigger Button -->
    <button
      @click="toggle"
      class="p-2 text-gray-400 hover:text-gray-600 hover:bg-gray-100 rounded-lg transition-colors"
      title="เมนู"
    >
      <MoreVertical class="w-5 h-5" />
    </button>

    <!-- Dropdown Menu -->
    <Transition
      enter-active-class="transition ease-out duration-100"
      enter-from-class="transform opacity-0 scale-95"
      enter-to-class="transform opacity-100 scale-100"
      leave-active-class="transition ease-in duration-75"
      leave-from-class="transform opacity-100 scale-100"
      leave-to-class="transform opacity-0 scale-95"
    >
      <div
        v-if="isOpen"
        class="absolute right-0 mt-2 w-48 bg-white rounded-xl shadow-lg border border-gray-200 py-1 z-50"
      >
        <!-- View -->
        <button
          @click.stop="handleView"
          class="w-full px-4 py-2.5 text-left text-sm text-gray-700 hover:bg-gray-50 flex items-center gap-3 transition-colors"
        >
          <Eye class="w-4 h-4 text-gray-400" />
          ดูรายละเอียด
        </button>

        <!-- Edit -->
        <button
          @click.stop="handleEdit"
          class="w-full px-4 py-2.5 text-left text-sm text-gray-700 hover:bg-gray-50 flex items-center gap-3 transition-colors"
        >
          <Pencil class="w-4 h-4 text-gray-400" />
          แก้ไข
        </button>

        <!-- Manage Features -->
        <button
          @click.stop="handleManageFeatures"
          class="w-full px-4 py-2.5 text-left text-sm text-gray-700 hover:bg-gray-50 flex items-center gap-3 transition-colors"
        >
          <Settings2 class="w-4 h-4 text-gray-400" />
          จัดการคุณสมบัติ
        </button>

        <!-- Divider -->
        <div class="border-t border-gray-100 my-1"></div>

        <!-- Delete -->
        <button
          @click.stop="handleDelete"
          class="w-full px-4 py-2.5 text-left text-sm text-red-600 hover:bg-red-50 flex items-center gap-3 transition-colors"
        >
          <Trash2 class="w-4 h-4" />
          ลบ
        </button>
      </div>
    </Transition>
  </div>
</template>
