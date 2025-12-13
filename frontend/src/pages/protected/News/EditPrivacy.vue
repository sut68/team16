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

// --- 1. นิยาม Status Options ---
const statusOptions = [
  { 
    id: 1, 
    name: 'Published', 
    label: 'เผยแพร่สาธารณะ (Public)', 
    desc: 'แสดงผลทันที บุคคลทั่วไปสามารถเห็นได้',
    style: 'border-green-500 bg-green-50/50 hover:bg-green-100',
    dot: 'bg-green-500'
  },
  { 
    id: 5, 
    name: 'Members Only', 
    label: 'เฉพาะสมาชิก (Members Only)', 
    desc: 'แสดงผลทันที แต่ต้องล็อกอินก่อนถึงจะเห็น',
    style: 'border-indigo-500 bg-indigo-50/50 hover:bg-indigo-100',
    dot: 'bg-indigo-500'
  },
  { 
    id: 2, 
    name: 'Draft', 
    label: 'ฉบับร่าง (Draft)', 
    desc: 'ซ่อนไว้แก้ไข ยังไม่แสดงบนหน้าเว็บ',
    style: 'border-orange-400 bg-orange-50/50 hover:bg-orange-100',
    dot: 'bg-orange-400'
  },
  { 
    id: 3, 
    name: 'Archived', 
    label: 'จัดเก็บ (Archived)', 
    desc: 'ข่าวหมดอายุ/เลิกใช้งาน (เก็บเป็นประวัติ)',
    style: 'border-slate-400 bg-slate-50/50 hover:bg-slate-100',
    dot: 'bg-slate-400'
  },
];

// --- 2. Form Setup (แก้ไขค่าเริ่มต้น) ---
const form = ref({
  title: '',       
  post_detail: '', 
  // ✅ แก้ไข: ไม่ Hardcode เลข 1 แล้ว ให้เป็น null รอรับค่าจาก API
  admin_id: null as number | null,
  scholarship_id: null as number | null,
  status_news_id: 1, // Default ไว้ที่ 1 ก่อน (เดี๋ยวจะถูกทับด้วยค่าจาก API)
});

// --- 3. ดึงข้อมูลเมื่อเปิด Modal ---
watch(() => props.isOpen, async (newVal) => {
  if (newVal && props.newsId) {
    isLoading.value = true;
    try {
      const data = await getNewsPostById(props.newsId);
      
      // Map ข้อมูลจริงจาก Database ใส่ Form
      form.value = {
        title: data.title,
        post_detail: data.post_detail,
        admin_id: data.admin_id,             // ค่าจริงจะถูกใส่ตรงนี้
        scholarship_id: data.scholarship_id, // ค่าจริงจะถูกใส่ตรงนี้
        status_news_id: data.status_news_id, 
      };

    } catch (error) {
      console.error("Failed to fetch settings:", error);
      // กรณีโหลดพลาด อาจจะ alert บอก user หรือปิด modal
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
    
    // แปลงค่าเป็น String (ถ้าเป็น null จะส่ง string ว่า "null" หรือ "0" แล้วแต่ backend จะรับไหวไหม 
    // แต่ปกติ watch ข้างบนจะทำงานก่อน ทำให้มีค่าเสมอ)
    formData.append('admin_id', String(form.value.admin_id || '')); 
    formData.append('scholarship_id', String(form.value.scholarship_id || ''));
    
    // ส่ง Status ID ที่เลือก
    formData.append('status_news_id', String(form.value.status_news_id));

    await updateNewsPost(props.newsId, formData);
    
    // ❌ เอา alert ออกตามที่คุยกันไว้
    // alert('อัปเดตสถานะเรียบร้อยแล้ว');
    
    emit('save'); 
    emit('close'); 
  } catch (error) {
    console.error(error);
    alert('เกิดข้อผิดพลาดในการบันทึก'); // ตรงนี้ไฟล์แม่ (NewsList) จะจัดการเองไม่ได้ถ้าเป็น error ภายในนี้ แต่ปล่อยไว้ก่อนได้
  } finally {
    isSaving.value = false;
  }
};
</script>

<template>
  <div v-if="isOpen" class="fixed inset-0 z-[1000] flex items-center justify-center bg-black/40 backdrop-blur-sm p-4 transition-opacity">
    <div class="bg-white w-full max-w-md rounded-xl shadow-2xl overflow-hidden transform transition-all scale-100 flex flex-col max-h-[90vh]">
      
      <div class="bg-white px-6 py-4 border-b border-gray-100 flex justify-between items-center shrink-0">
        <div>
            <h3 class="text-lg font-bold text-[#1e3a8a]">จัดการสถานะข่าวสาร</h3>
            <p class="text-xs text-gray-500">News ID: #{{ newsId }}</p>
        </div>
        <button @click="$emit('close')" class="btn btn-sm btn-circle btn-ghost text-gray-400 hover:bg-gray-100">✕</button>
      </div>

      <div v-if="isLoading" class="flex-1 flex items-center justify-center p-10">
        <span class="loading loading-spinner loading-md text-[#1e3a8a]"></span>
      </div>

      <div v-else class="p-6 overflow-y-auto custom-scrollbar">
        <div class="form-control">
            <label class="label pb-3">
                <span class="label-text font-bold text-slate-700 text-base">เลือกสถานะ (Select Status)</span>
            </label>
            
            <div class="space-y-3">
                <label 
                    v-for="option in statusOptions" 
                    :key="option.id"
                    class="flex items-start gap-3 p-4 border rounded-xl cursor-pointer transition-all duration-200 relative overflow-hidden"
                    :class="[
                        form.status_news_id === option.id ? option.style + ' shadow-md ring-1 ring-offset-1' : 'border-gray-200 bg-white hover:border-gray-300'
                    ]"
                >
                    <input 
                        type="radio" 
                        :value="option.id" 
                        v-model="form.status_news_id" 
                        class="radio radio-sm mt-1"
                        :class="{'radio-primary': form.status_news_id === option.id}" 
                    />
                    
                    <div class="flex-1 z-10">
                        <div class="flex items-center gap-2 mb-1">
                            <span class="w-2 h-2 rounded-full" :class="option.dot"></span>
                            <span class="font-bold text-slate-800 text-sm">{{ option.label }}</span>
                        </div>
                        <p class="text-xs text-gray-500 font-medium leading-relaxed">{{ option.desc }}</p>
                    </div>

                    <div v-if="form.status_news_id === option.id" class="absolute top-4 right-4 text-slate-600/20">
                        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" class="w-6 h-6">
                            <path fill-rule="evenodd" d="M2.25 12c0-5.385 4.365-9.75 9.75-9.75s9.75 4.365 9.75 9.75-4.365 9.75-9.75 9.75S2.25 17.385 2.25 12Zm13.36-1.814a.75.75 0 1 0-1.22-.872l-3.236 4.53L9.53 12.22a.75.75 0 0 0-1.06 1.06l2.25 2.25a.75.75 0 0 0 1.14-.094l3.75-5.25Z" clip-rule="evenodd" />
                        </svg>
                    </div>
                </label>
            </div>
        </div>
      </div>

      <div class="bg-gray-50 px-6 py-4 flex justify-end gap-3 border-t border-gray-100 shrink-0">
        <button @click="$emit('close')" class="btn btn-sm btn-ghost text-gray-500 font-normal">ยกเลิก</button>
        <button 
            @click="handleSave" 
            class="btn btn-sm bg-[#1e3a8a] hover:bg-[#152c6f] text-white border-none shadow-sm px-6"
            :disabled="isSaving || isLoading"
        >
            <span v-if="isSaving" class="loading loading-spinner loading-xs"></span>
            บันทึก
        </button>
      </div>

    </div>
  </div>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar { width: 6px; }
.custom-scrollbar::-webkit-scrollbar-track { background: #f1f1f1; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #d1d5db; border-radius: 10px; }
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background: #9ca3af; }
</style>