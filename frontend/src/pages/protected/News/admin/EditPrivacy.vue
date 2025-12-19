<script setup lang="ts">
import { ref, watch } from 'vue';
import { getNewsPostById, updateNewsPost } from '@/services/api/news_post';

const props = defineProps<{
  isOpen: boolean;
  newsId: number;
}>();

const emit = defineEmits(['close', 'save']);

const isSaving = ref(false);
const isLoading = ref(false);

// --- 1. นิยาม Status Options (เพิ่ม field: borderColor, radioInner) ---
// แก้ไข: ใส่ชื่อ Class เต็มๆ เพื่อให้ Tailwind ทำงานได้ถูกต้อง
const statusOptions = [
  { 
    id: 1, 
    label: 'เผยแพร่สาธารณะ (Public)', 
    desc: 'แสดงผลทันที บุคคลทั่วไปสามารถเห็นได้',
    // Card Style
    style: 'border-emerald-500 bg-emerald-50/50 hover:bg-emerald-50 ring-emerald-500', 
    // Radio Border Color
    borderColor: 'border-emerald-500',
    // Inner Dot Color
    dotColor: 'bg-emerald-500',
    icon: 'M12 21a9.004 9.004 0 008.716-6.747M12 21a9.004 9.004 0 01-8.716-6.747M12 21c2.485 0 4.5-4.03 4.5-9S14.485 3 12 3m0 18c-2.485 0-4.5-4.03-4.5-9S9.515 3 12 3m0 0a8.997 8.997 0 017.843 4.582M12 3a8.997 8.997 0 00-7.843 4.582m15.686 0A11.953 11.953 0 0112 10.5c-2.998 0-5.74-1.1-7.843-2.918m15.686 0A8.959 8.959 0 0121 12c0 .778-.099 1.533-.284 2.253m0 0A17.919 17.919 0 0112 16.5c-3.162 0-6.133-.815-8.716-2.247m0 0A9.015 9.015 0 013 12c0-1.605.42-3.113 1.157-4.418'
  },
  { 
    id: 4, 
    label: 'เฉพาะสมาชิก (Members Only)', 
    desc: 'แสดงผลทันที แต่ต้องล็อกอินก่อนถึงจะเห็น',
    style: 'border-indigo-500 bg-indigo-50/50 hover:bg-indigo-50 ring-indigo-500',
    borderColor: 'border-indigo-500',
    dotColor: 'bg-indigo-500',
    icon: 'M16.5 10.5V6.75a4.5 4.5 0 10-9 0v3.75m-.75 11.25h10.5a2.25 2.25 0 002.25-2.25v-6.75a2.25 2.25 0 00-2.25-2.25H6.75a2.25 2.25 0 00-2.25 2.25v6.75a2.25 2.25 0 002.25 2.25z'
  },
  { 
    id: 2, 
    label: 'ฉบับร่าง (Draft)', 
    desc: 'ซ่อนไว้แก้ไข ยังไม่แสดงบนหน้าเว็บ',
    style: 'border-orange-400 bg-orange-50/50 hover:bg-orange-50 ring-orange-400',
    borderColor: 'border-orange-400',
    dotColor: 'bg-orange-400',
    icon: 'M16.862 4.487l1.687-1.688a1.875 1.875 0 112.652 2.652L10.582 16.07a4.5 4.5 0 01-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 011.13-1.897l8.932-8.931zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0115.75 21H5.25A2.25 2.25 0 013 18.75V8.25A2.25 2.25 0 015.25 6H10'
  },
  { 
    id: 3, 
    label: 'จัดเก็บ (Archived)', 
    desc: 'ข่าวหมดอายุ/เลิกใช้งาน (เก็บเป็นประวัติ)',
    style: 'border-slate-400 bg-slate-50/50 hover:bg-slate-50 ring-slate-400',
    borderColor: 'border-slate-400',
    dotColor: 'bg-slate-400',
    icon: 'M20.25 7.5l-.625 10.632a2.25 2.25 0 01-2.247 2.118H6.622a2.25 2.25 0 01-2.247-2.118L3.75 7.5M10 11.25h4M3.375 7.5h17.25c.621 0 1.125-.504 1.125-1.125v-1.5c0-.621-.504-1.125-1.125-1.125H3.375c-.621 0-1.125.504-1.125 1.125v1.5c0 .621.504 1.125 1.125 1.125z'
  },
];

// --- 2. Form Setup ---
const form = ref({
  title: '',       
  post_detail: '', 
  admin_id: null as number | null,
  scholarship_id: null as number | null,
  status_news_id: 1, 
});

// --- 3. ดึงข้อมูลเมื่อเปิด Modal ---
watch(() => props.isOpen, async (newVal) => {
  if (newVal && props.newsId) {
    isLoading.value = true;
    try {
      const response: any = await getNewsPostById(props.newsId);
      const data = response.news_post; 

      form.value = {
        title: data.title || '',
        post_detail: data.post_detail || '',
        admin_id: data.admin_id,       
        scholarship_id: data.scholarship_id, 
        status_news_id: data.status_news_id || 1, 
      };

    } catch (error) {
      console.error("Failed to fetch settings:", error);
    } finally {
      isLoading.value = false;
    }
  }
});

// --- 4. บันทึกข้อมูล ---
const handleSave = async () => {
  isSaving.value = true;
  try {
    const formData = new FormData();
    formData.append('title', form.value.title);
    formData.append('post_detail', form.value.post_detail);
    
    if (form.value.admin_id) formData.append('admin_id', String(form.value.admin_id));
    if (form.value.scholarship_id) formData.append('scholarship_id', String(form.value.scholarship_id));
    
    formData.append('status_news_id', String(form.value.status_news_id));

    await updateNewsPost(props.newsId, formData);
    
    emit('save'); 
    emit('close'); 
  } catch (error) {
    console.error(error);
    alert('เกิดข้อผิดพลาดในการบันทึก');
  } finally {
    isSaving.value = false;
  }
};
</script>

<template>
  <div v-if="isOpen" class="fixed inset-0 z-[1000] flex items-center justify-center bg-slate-900/60 backdrop-blur-sm p-4 transition-opacity">
    <div class="bg-white w-full max-w-lg rounded-2xl shadow-2xl overflow-hidden transform transition-all scale-100 flex flex-col max-h-[90vh] animate-bounce-in border border-white/20">
      
      <div class="bg-white px-8 py-5 border-b border-gray-100 flex justify-between items-center shrink-0">
        <div>
            <h3 class="text-xl font-extrabold text-slate-800 tracking-tight">ตั้งค่าความเป็นส่วนตัว</h3>
            <p class="text-xs text-slate-400 mt-0.5 font-mono">News ID: #{{ newsId }}</p>
        </div>
        <button @click="$emit('close')" class="btn btn-sm btn-circle btn-ghost text-slate-400 hover:bg-slate-100 transition-colors">
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" class="w-5 h-5"><path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" /></svg>
        </button>
      </div>

      <div v-if="isLoading" class="flex-1 flex flex-col items-center justify-center p-12 text-slate-400">
        <span class="loading loading-dots loading-lg text-[#1e3a8a]"></span>
        <p class="mt-4 text-sm font-medium animate-pulse">กำลังโหลดข้อมูล...</p>
      </div>

      <div v-else class="p-8 overflow-y-auto custom-scrollbar">
        <div class="form-control">
            <label class="label pb-4 pt-0">
                <span class="label-text font-bold text-slate-700 text-base">การมองเห็น</span>
            </label>
            
            <div class="space-y-4">
                <label 
                    v-for="option in statusOptions" 
                    :key="option.id"
                    class="group flex items-start gap-4 p-4 border rounded-xl cursor-pointer transition-all duration-200 relative overflow-hidden"
                    :class="[
                        form.status_news_id === option.id 
                            ? option.style + ' shadow-md ring-1 ring-offset-0 bg-opacity-100' 
                            : 'border-slate-200 bg-white hover:border-slate-300 hover:bg-slate-50'
                    ]"
                >
                    <input 
                        type="radio" 
                        :value="option.id" 
                        v-model="form.status_news_id" 
                        class="sr-only"
                    />

                    <div class="mt-1">
                        <div class="w-5 h-5 rounded-full border-2 flex items-center justify-center transition-all duration-300"
                            :class="[
                                form.status_news_id === option.id 
                                    ? option.borderColor 
                                    : 'border-slate-300 group-hover:border-slate-400'
                            ]"
                        >
                            <div class="w-2.5 h-2.5 rounded-full transition-all duration-300"
                                :class="[
                                    option.dotColor, 
                                    form.status_news_id === option.id ? 'scale-100 opacity-100' : 'scale-0 opacity-0'
                                ]"
                            ></div>
                        </div>
                    </div>
                    
                    <div class="flex-1 z-10">
                        <div class="flex items-center gap-2 mb-1.5">
                            <span class="font-bold text-slate-800 text-sm group-hover:text-[#1e3a8a] transition-colors">
                                {{ option.label }}
                            </span>
                        </div>
                        <p class="text-xs text-slate-500 font-medium leading-relaxed group-hover:text-slate-600">
                            {{ option.desc }}
                        </p>
                    </div>

                    <div class="absolute top-4 right-4 text-slate-400/20 group-hover:text-slate-400/40 transition-colors">
                        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" class="w-8 h-8">
                            <path stroke-linecap="round" stroke-linejoin="round" :d="option.icon" />
                        </svg>
                    </div>
                </label>
            </div>
        </div>
      </div>

      <div class="bg-gray-50 px-8 py-5 flex justify-end gap-3 border-t border-gray-100 shrink-0">
        <button @click="$emit('close')" class="btn btn-sm h-10 px-5 rounded-lg btn-ghost text-slate-500 hover:bg-slate-200 font-normal">ยกเลิก</button>
        <button 
            @click="handleSave" 
            class="btn btn-sm h-10 px-6 rounded-lg bg-[#1e3a8a] hover:bg-[#152c6f] text-white border-none shadow-md shadow-blue-900/10"
            :disabled="isSaving || isLoading"
        >
            <span v-if="isSaving" class="loading loading-spinner loading-xs"></span>
            {{ isSaving ? 'กำลังบันทึก...' : 'บันทึกการเปลี่ยนแปลง' }}
        </button>
      </div>

    </div>
  </div>
</template>

<style scoped>
@keyframes bounceIn {
  0% { opacity: 0; transform: scale(0.95); }
  100% { opacity: 1; transform: scale(1); }
}
.animate-bounce-in {
  animation: bounceIn 0.2s cubic-bezier(0.175, 0.885, 0.32, 1.275);
}

.custom-scrollbar::-webkit-scrollbar { width: 5px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #cbd5e1; border-radius: 10px; }
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background: #94a3b8; }
</style>