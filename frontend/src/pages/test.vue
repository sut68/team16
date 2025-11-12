<template>
  <div class="flex h-screen w-full bg-white">
    
    <aside
      class="h-screen shrink-0 bg-[#8B0015] text-white
             transition-all duration-300 ease-in-out"
      :class="[isSidebarOpen ? 'w-72' : 'w-20']"
      role="navigation"
    >
      <div class="pt-6 pb-4" :class="[isSidebarOpen ? 'px-6' : 'px-4']">
        <div 
          class="flex items-center gap-3" 
          :class="[isSidebarOpen ? 'flex-row' : 'flex-col items-center']"
        >
          <svg class="w-10 h-10 shrink-0" viewBox="0 0 24 24" fill="currentColor">
            <path d="M10.325 4.317a1 1 0 0 1 .894-.553h1.562a1 1 0 0 1 .894.553l.48.96a1 1 0 0 0 .708.536l1.053.21a1 1 0 0 1 .72.6l.593 1.48a1 1 0 0 1-.235 1.07l-.757.757a1 1 0 0 0 0 1.414l.757.757a1 1 0 0 1 .235 1.07l-.593 1.48a1 1 0 0 1-.72.6l-1.053.21a1 1 0 0 0-.708.536l-.48.96a1 1 0 0 1-.894.553h-1.562a1 1 0 0 1-.894-.553l-.48-.96a1 1 0 0 0-.708-.536l-1.053-.21a1 1 0 0 1-.72-.6l-.593-1.48a1 1 0 0 1 .235-1.07l.757-.757a1 1 0 0 0 0-1.414l-.757-.757a1 1 0 0 1-.235-1.07l.593-1.48a1 1 0 0 1 .72-.6l1.053-.21a1 1 0 0 0 .708.536l.48-.96z"/>
            <circle cx="12" cy="12" r="3.25" fill="#fff"/>
          </svg>
          <div class="leading-tight text-center" v-show="isSidebarOpen">
            <div class="text-sm font-semibold">ทุนการศึกษา</div>
            <div class="text-xs opacity-90">สำนักวิชาวิศวกรรมศาสตร์</div>
          </div>
        </div>
      </div>

      <div class="pr-0" :class="[isSidebarOpen ? 'pl-4' : 'pl-2']">
        <ul class="space-y-2">
          
          <li>
            <div v-if="currentPage === 'dashboard'">
              <div 
                class="flex items-center bg-white text-[#8B0015]
                       rounded-l-xl
                       -mr-4
                       rounded-tr-3xl rounded-br-3xl" :class="[
                  isSidebarOpen ? '' : 'justify-center',
                ]"
              >
                <span 
                  class="btn btn-ghost gap-3 flex-1 text-[#8B0015] normal-case"
                  :class="[isSidebarOpen ? 'justify-start' : 'justify-center']"
                >
                  <svg class="w-5 h-5 shrink-0" viewBox="0 0 24 24" fill="currentColor"><path d="M3 3h8v8H3V3zm10 0h8v6h-8V3zM3 13h6v8H3v-8zm8 4h10v4H11v-4z"/></svg>
                  <span v-show="isSidebarOpen">แดชบอร์ด</span>
                </span>
                <button 
                  class="btn btn-circle btn-outline btn-sm border-2" 
                  :class="[isSidebarOpen ? 'mr-2' : '']"
                  @click="toggleSidebar"
                >
                  <svg class="w-4 h-4 transition-transform duration-300" :class="[isSidebarOpen ? '' : 'rotate-180']" viewBox="0 0 24 24" fill="currentColor"><path d="M15 19l-7-7 7-7"/></svg>
                </button>
              </div>
            </div>
            
            <button
              v-else
              class="btn btn-ghost gap-3 w-full text-white normal-case opacity-90 hover:opacity-100"
              :class="[
                isSidebarOpen ? 'justify-start' : 'justify-center',
                currentPage === 'applications' ? 'rounded-br-3xl' : '' // S-Curve (Inactive)
              ]"
              @click="currentPage = 'dashboard'"
            >
              <svg class="w-5 h-5 shrink-0" viewBox="0 0 24 24" fill="currentColor"><path d="M3 3h8v8H3V3zm10 0h8v6h-8V3zM3 13h6v8H3v-8zm8 4h10v4H11v-4z"/></svg>
              <span v-show="isSidebarOpen">แดชบอร์ด</span>
            </button>
          </li>

          <li>
            <div v-if="currentPage === 'applications'">
              <div 
                class="flex items-center bg-white text-[#8B0015]
                       rounded-l-xl
                       -mr-4
                       rounded-tr-3xl rounded-br-3xl" :class="[
                  isSidebarOpen ? '' : 'justify-center',
                ]"
              >
                <span 
                  class="btn btn-ghost gap-3 flex-1 text-[#8B0015] normal-case"
                  :class="[isSidebarOpen ? 'justify-start' : 'justify-center']"
                >
                  <svg class="w-5 h-5 shrink-0" viewBox="0 0 24 24" fill="currentColor"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><path d="M14 2v6h6" fill="#ffffff55"/></svg>
                  <span v-show="isSidebarOpen">เอกสารสมัคร</span>
                </span>
                <button 
                  class="btn btn-circle btn-outline btn-sm border-2" 
                  :class="[isSidebarOpen ? 'mr-2' : '']"
                  @click="toggleSidebar"
                >
                  <svg class="w-4 h-4 transition-transform duration-300" :class="[isSidebarOpen ? '' : 'rotate-180']" viewBox="0 0 24 24" fill="currentColor"><path d="M15 19l-7-7 7-7"/></svg>
                </button>
              </div>
            </div>
            
            <button
              v-else
              class="btn btn-ghost gap-3 w-full text-white normal-case opacity-90 hover:opacity-100"
              :class="[
                isSidebarOpen ? 'justify-start' : 'justify-center',
                currentPage === 'dashboard' ? 'rounded-tr-3xl' : '' // S-Curve (Inactive)
              ]"
              @click="currentPage = 'applications'"
            >
              <svg class="w-5 h-5 shrink-0" viewBox="0 0 24 24" fill="currentColor"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><path d="M14 2v6h6" fill="#ffffff55"/></svg>
              <span v-show="isSidebarOpen">เอกสารสมัคร</span>
            </button>
          </li>

        </ul>
      </div>
    </aside>

    <main class="flex-1 p-8 overflow-y-auto">
      <h1 class="text-2xl font-bold">
        หน้า {{ currentPage === 'dashboard' ? 'แดชบอร์ด' : 'เอกสารสมัคร' }}
      </h1>
      <p class="mt-2 text-gray-600">
        สถานะ Sidebar: <strong>{{ isSidebarOpen ? 'เปิด' : 'ปิด' }}</strong>
      </p>

      <div class="mt-4 p-4 bg-gray-50 rounded-lg border border-gray-200">
        <p>นี่คือเนื้อหาของหน้า {{ currentPage }}</p>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';

// --- State ---
const isSidebarOpen = ref(true);
const currentPage = ref('applications'); // (ตั้งค่าเริ่มต้นให้ 'applications' active ตามภาพ)

const toggleSidebar = () => {
  isSidebarOpen.value = !isSidebarOpen.value;
};
</script>