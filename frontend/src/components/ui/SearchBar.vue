<script setup lang="ts">
import { Search, X } from 'lucide-vue-next'

const model = defineModel<string>({ required: true })

defineProps<{
  placeholder?: string
}>()

defineEmits<{
  clear: []
}>()
</script>

<template>
  <div class="relative flex-1 md:flex-none w-full md:w-96" data-testid="search-container">
    <span class="absolute inset-y-0 left-0 flex items-center pl-3 text-gray-400 pointer-events-none">
      <Search class="w-5 h-5" />
    </span>

    <input
      type="text"
      v-model="model"
      :placeholder="placeholder ?? 'ค้นหา...'"
      class="input input-bordered input-sm h-10 w-full rounded-full pl-10 
             bg-white border-gray-300 shadow-sm 
             focus:border-[#1e3a8a] focus:ring-1 focus:ring-[#1e3a8a] text-sm"
      data-testid="search-input"
    />

    <button
      v-if="model"
      @click="model = ''; $emit('clear')"
      class="absolute right-2 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-700"
      data-testid="search-clear-btn"
    >
      <X class="w-4 h-4" />
    </button>
  </div>
</template>