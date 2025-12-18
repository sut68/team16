<script setup lang="ts">
import { ref, computed, watch } from 'vue';
import { getNewsPostById, updateNewsPost } from '@/services/api/news_post';

const props = defineProps<{ id: number }>();
const emit = defineEmits(['close', 'success']);

const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8080"; 

const isLoading = ref(false);
const isSaving = ref(false);

const statusOptions = [
  { id: 1, label: 'เผยแพร่สาธารณะ (Public)', class: 'text-emerald-600 font-medium' },
  { id: 5, label: 'เฉพาะสมาชิก (Members Only)', class: 'text-indigo-600 font-medium' },
  { id: 2, label: 'ฉบับร่าง (Draft)', class: 'text-orange-500 font-medium' },
  { id: 3, label: 'จัดเก็บ (Archived)', class: 'text-slate-500 font-medium' }
];

const form = ref({
  title: '',
  post_detail: '',
  admin_id: null as number | null,
  scholarship_id: null as number | null,
  status_news_id: 1
});

const oldFilePath = ref('');
const newImageFile = ref<File | undefined>(undefined);
const previewImage = ref<string | null>(null);

const displayImage = computed(() => {
  if (previewImage.value) return previewImage.value;
  if (oldFilePath.value) {
    const cleanPath = oldFilePath.value.startsWith('/') 
        ? oldFilePath.value.substring(1) 
        : oldFilePath.value;
    return `${API_URL}/${cleanPath}?t=${new Date().getTime()}`; 
  }
  return null;
});

// --- Popup System ---
interface PopupState {
  isOpen: boolean;
  type: 'success' | 'error' | 'confirm';
  title: string;
  message: string;
  resolve: ((value: boolean) => void) | null;
}

const popup = ref<PopupState>({
  isOpen: false,
  type: 'success',
  title: '',
  message: '',
  resolve: null
});

const showPopup = (type: 'success' | 'error' | 'confirm', title: string, message: string) => {
  return new Promise<boolean>((resolve) => {
    popup.value = { isOpen: true, type, title, message, resolve };
  });
};

const closePopup = (result: boolean) => {
  if (popup.value.resolve) popup.value.resolve(result);
  popup.value.isOpen = false;
};

// --- Fetch Data ---
const fetchNewsDetail = async () => {
  if (!props.id) return;

  isLoading.value = true;
  try {
    const response = await getNewsPostById(props.id);
    const data = response.news_post as any;

    form.value = {
      title: data.title || '',
      post_detail: data.post_detail || '',
      admin_id: data.admin_id,
      scholarship_id: data.scholarship_id,
      status_news_id: data.status_news_id ?? 1
    };
    
    oldFilePath.value = typeof data.file_path === 'string' ? data.file_path : '';
    newImageFile.value = undefined;
    previewImage.value = null;
    
  } catch (error) {
    console.error("Error:", error);
    showPopup('error', 'ผิดพลาด', 'ไม่สามารถโหลดข้อมูลข่าวได้');
    emit('close');
  } finally {
    isLoading.value = false;
  }
};

watch(() => props.id, (newId) => {
    if (newId) fetchNewsDetail();
}, { immediate: true }); 

const handleFileUpload = (event: Event) => {
  const target = event.target as HTMLInputElement;
  if (target.files && target.files[0]) {
    const file = target.files[0];
    if (file.size > 5 * 1024 * 1024) {
        showPopup('error', 'ไฟล์ใหญ่เกินไป', 'ขนาดไฟล์ต้องไม่เกิน 5MB');
        return;
    }
    newImageFile.value = file;
    previewImage.value = URL.createObjectURL(file);
  }
};

const handleSave = async () => {
  if (!form.value.title?.trim() || !form.value.post_detail?.trim()) {
      showPopup('error', 'ข้อมูลไม่ครบ', 'กรุณากรอกหัวข้อและเนื้อหาข่าว');
      return;
  }
  
  isSaving.value = true;
  try {
    const fd = new FormData();
    fd.append("title", form.value.title);
    fd.append("post_detail", form.value.post_detail);
    if (form.value.admin_id) fd.append("admin_id", String(form.value.admin_id));
    if (form.value.scholarship_id) fd.append("scholarship_id", String(form.value.scholarship_id));
    fd.append("status_news_id", String(form.value.status_news_id));

    if (newImageFile.value) {
      fd.append("file_path", newImageFile.value);
    }

    await updateNewsPost(props.id, fd);
    await showPopup('success', 'สำเร็จ', 'บันทึกการแก้ไขเรียบร้อยแล้ว');
    emit('success');

  } catch (error) {
    console.error("Update Error:", error);
    showPopup('error', 'ผิดพลาด', 'บันทึกข้อมูลไม่สำเร็จ');
  } finally {
    isSaving.value = false;
  }
};
</script>

<template>
  <div class="w-full h-full p-4 md:p-8 flex justify-center bg-[#f8fafc]">
    
    <div class="w-full max-w-5xl bg-white rounded-2xl shadow-xl shadow-slate-200/60 overflow-hidden border border-slate-100 flex flex-col h-fit animate-fade-in-up">
        
        <div class="px-8 py-6 border-b border-slate-100 flex flex-row justify-between items-center bg-white sticky top-0 z-10">
            <div>
                <h1 class="text-2xl font-extrabold text-slate-800 tracking-tight">แก้ไขข่าวสาร</h1>
                <p class="text-sm text-slate-400 mt-1 font-mono">Reference ID: {{ id }}</p>
            </div>
            <div class="flex gap-3">
                <button @click="emit('close')" class="btn btn-ghost text-slate-500 hover:bg-slate-100" :disabled="isSaving">
                    ยกเลิก
                </button>
                <button @click="handleSave" class="btn bg-[#1e3a8a] hover:bg-[#152c6f] text-white border-none shadow-lg shadow-blue-900/20 px-6 gap-2" :disabled="isSaving">
                    <span v-if="isSaving" class="loading loading-spinner loading-xs"></span>
                    {{ isSaving ? 'กำลังบันทึก...' : 'บันทึกการแก้ไข' }}
                </button>
            </div>
        </div>
        
        <div v-if="isLoading" class="flex flex-col items-center justify-center py-32 text-slate-400">
            <span class="loading loading-dots loading-lg text-[#1e3a8a]"></span>
            <p class="mt-4 text-sm font-medium animate-pulse">กำลังดึงข้อมูล...</p>
        </div>

        <div v-else class="p-8 md:p-10">
            <div class="grid grid-cols-1 lg:grid-cols-12 gap-8 lg:gap-12">
                
                <div class="lg:col-span-8 space-y-8">
                    <div class="form-control w-full">
                        <label class="label pt-0 pb-2">
                            <span class="label-text font-bold text-slate-700 text-base">หัวข้อข่าว <span class="text-red-500">*</span></span>
                        </label>
                        <input 
                            v-model="form.title" 
                            type="text" 
                            class="input input-bordered w-full h-12 text-lg px-4 bg-slate-50 border-slate-200 focus:border-[#1e3a8a] focus:ring-1 focus:ring-[#1e3a8a] focus:bg-white transition-all rounded-xl placeholder:text-slate-400" 
                            placeholder="ระบุหัวข้อข่าวที่น่าสนใจ..." 
                        />
                    </div>

                    <div class="form-control w-full">
                        <label class="label pb-2">
                            <span class="label-text font-bold text-slate-700 text-base">เนื้อหาข่าว <span class="text-red-500">*</span></span>
                        </label>
                        <textarea 
                            v-model="form.post_detail" 
                            class="textarea textarea-bordered min-h-[400px] p-5 text-base leading-7 bg-slate-50 border-slate-200 focus:border-[#1e3a8a] focus:ring-1 focus:ring-[#1e3a8a] focus:bg-white transition-all rounded-xl placeholder:text-slate-400 resize-none scrollbar-thin scrollbar-thumb-slate-200 scrollbar-track-transparent"
                            placeholder="เขียนรายละเอียดข่าวสารประชาสัมพันธ์ที่นี่..."
                        ></textarea>
                    </div>
                </div>

                <div class="lg:col-span-4 space-y-8">
                    
                    <div class="bg-white p-6 rounded-2xl border border-slate-100 shadow-sm">
                        <div class="form-control w-full">
                            <label class="label pt-0 pb-3">
                                <span class="label-text font-bold text-slate-700 flex items-center gap-2">
                                    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="w-4 h-4 text-slate-400"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm.75-13a.75.75 0 00-1.5 0v5c0 .414.336.75.75.75h4a.75.75 0 000-1.5h-3.25V5z" clip-rule="evenodd" /></svg>
                                    สถานะการแสดงผล
                                </span>
                            </label>
                            <select v-model="form.status_news_id" class="select select-bordered w-full h-12 bg-white text-base border-slate-200 focus:border-[#1e3a8a] rounded-xl font-medium">
                                <option v-for="option in statusOptions" :key="option.id" :value="option.id" :class="option.class">
                                    {{ option.label }}
                                </option>
                            </select>
                            
                            <div class="mt-4 p-3 bg-blue-50/50 rounded-lg border border-blue-100/50">
                                <div class="flex gap-2">
                                    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="w-5 h-5 text-blue-500 shrink-0 mt-0.5"><path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a.75.75 0 000 1.5h.253a.25.25 0 01.244.304l-.459 2.066A1.75 1.75 0 0010.747 15H11a.75.75 0 000-1.5h-.253a.25.25 0 01-.244-.304l.459-2.066A1.75 1.75 0 009.253 9H9z" clip-rule="evenodd" /></svg>
                                    <p class="text-xs text-slate-500 leading-relaxed">
                                        <span v-if="form.status_news_id === 1">ข่าวนี้จะปรากฏบนหน้าเว็บไซต์ทันทีที่บันทึก ทุกคนสามารถเข้าถึงได้</span>
                                        <span v-else-if="form.status_news_id === 5">เฉพาะผู้ใช้งานที่เข้าสู่ระบบเท่านั้นจึงจะมองเห็นข่าวนี้</span>
                                        <span v-else-if="form.status_news_id === 2">ข่าวถูกซ่อนไว้ชั่วคราว มีเพียง Admin เท่านั้นที่มองเห็น</span>
                                        <span v-else-if="form.status_news_id === 3">ข่าวเก่าที่ถูกเก็บเข้ากรุ ไม่แสดงหน้าแรกแต่ยังค้นหาได้</span>
                                    </p>
                                </div>
                            </div>
                        </div>
                    </div>

                    <div class="bg-white p-6 rounded-2xl border border-slate-100 shadow-sm">
                        <label class="label pt-0 pb-3 flex justify-between items-center">
                            <span class="label-text font-bold text-slate-700 flex items-center gap-2">
                                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="w-4 h-4 text-slate-400"><path fill-rule="evenodd" d="M1 5.25A2.25 2.25 0 013.25 3h13.5A2.25 2.25 0 0119 5.25v9.5A2.25 2.25 0 0116.75 17H3.25A2.25 2.25 0 011 14.75v-9.5zm1.5 5.81v3.69c0 .414.336.75.75.75h13.5a.75.75 0 00.75-.75v-2.69l-2.22-2.219a.75.75 0 00-1.06 0l-1.91 1.909.47.47a.75.75 0 11-1.06 1.06L6.53 8.091a.75.75 0 00-1.06 0l-2.97 2.97zM12 7a1 1 0 11-2 0 1 1 0 012 0z" clip-rule="evenodd" /></svg>
                                รูปภาพหน้าปก
                            </span>
                            <span v-if="newImageFile" class="badge badge-success badge-sm text-white font-normal gap-1">
                                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="w-3 h-3"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.857-9.809a.75.75 0 00-1.214-.882l-3.483 4.79-1.88-1.88a.75.75 0 10-1.06 1.061l2.5 2.5a.75.75 0 001.137-.089l4-5.5z" clip-rule="evenodd" /></svg>
                                เลือกแล้ว
                            </span>
                        </label>
                        
                        <div class="relative group cursor-pointer w-full">
                            <div class="w-full aspect-[16/9] rounded-xl border-2 border-dashed border-slate-300 bg-slate-50 flex flex-col items-center justify-center overflow-hidden hover:border-[#1e3a8a] hover:bg-blue-50/50 transition-all duration-300 relative">
                                
                                <img v-if="displayImage" :src="displayImage" class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-105" />
                                
                                <div v-else class="flex flex-col items-center justify-center text-slate-400 gap-3">
                                    <div class="w-12 h-12 rounded-full bg-white shadow-sm flex items-center justify-center">
                                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6"><path stroke-linecap="round" stroke-linejoin="round" d="M3 16.5v2.25A2.25 2.25 0 005.25 21h13.5A2.25 2.25 0 0021 18.75V16.5m-13.5-9L12 3m0 0l4.5 4.5M12 3v13.5" /></svg>
                                    </div>
                                    <span class="text-sm font-medium">คลิกเพื่ออัปโหลด</span>
                                </div>

                                <div class="absolute inset-0 bg-slate-900/40 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity backdrop-blur-[2px]">
                                    <span class="btn btn-sm bg-white border-none text-slate-800 hover:bg-slate-100 shadow-md">
                                        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="w-4 h-4 mr-1"><path d="M5.433 13.917l1.262-3.155A4 4 0 017.58 9.42l6.92-6.918a2.121 2.121 0 013 3l-6.92 6.918c-.383.383-.84.685-1.343.886l-3.154 1.262a.5.5 0 01-.65-.65z" /><path d="M3.5 5.75c0-.69.56-1.25 1.25-1.25H10A.75.75 0 0010 3H4.75A2.75 2.75 0 002 5.75v9.5A2.75 2.75 0 004.75 18h9.5A2.75 2.75 0 0017 15.25V10a.75.75 0 00-1.5 0v5.25c0 .69-.56 1.25-1.25 1.25h-9.5c-.69 0-1.25-.56-1.25-1.25v-9.5z" /></svg>
                                        เปลี่ยนรูปภาพ
                                    </span>
                                </div>
                                
                                <input type="file" accept="image/*" @change="handleFileUpload" class="absolute inset-0 w-full h-full opacity-0 cursor-pointer" />
                            </div>
                            <div class="mt-3 text-xs text-slate-400 text-center flex items-center justify-center gap-1">
                                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="w-3 h-3"><path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a.75.75 0 000 1.5h.253a.25.25 0 01.244.304l-.459 2.066A1.75 1.75 0 0010.747 15H11a.75.75 0 000-1.5h-.253a.25.25 0 01-.244-.304l.459-2.066A1.75 1.75 0 009.253 9H9z" clip-rule="evenodd" /></svg>
                                รองรับไฟล์ JPG, PNG ขนาดไม่เกิน 5MB
                            </div>
                        </div>
                    </div>

                </div>
            </div>
        </div>

        <div v-if="popup.isOpen" class="fixed inset-0 z-[2000] flex items-center justify-center bg-slate-900/60 backdrop-blur-sm p-4">
            <div class="bg-white w-full max-w-sm rounded-2xl shadow-2xl p-6 text-center animate-bounce-in border border-white/20">
                <div class="mx-auto mb-4 w-14 h-14 flex items-center justify-center rounded-full shadow-inner"
                    :class="{
                        'bg-emerald-50 text-emerald-500': popup.type === 'success',
                        'bg-rose-50 text-rose-500': popup.type === 'error',
                        'bg-amber-50 text-amber-500': popup.type === 'confirm'
                    }">
                    <svg v-if="popup.type === 'success'" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor" class="w-7 h-7"><path stroke-linecap="round" stroke-linejoin="round" d="m4.5 12.75 6 6 9-13.5" /></svg>
                    <svg v-else-if="popup.type === 'error'" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor" class="w-7 h-7"><path stroke-linecap="round" stroke-linejoin="round" d="M6 18 18 6M6 6l12 12" /></svg>
                    <svg v-else xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor" class="w-7 h-7"><path stroke-linecap="round" stroke-linejoin="round" d="M12 9v3.75m-9.303 3.376c-.866 1.5.217 3.374 1.948 3.374h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 3.378c-.866-1.5-3.032-1.5-3.898 0L2.697 16.126ZM12 15.75h.007v.008H12v-.008Z" /></svg>
                </div>
                <h3 class="text-lg font-bold text-slate-800 mb-2">{{ popup.title }}</h3>
                <p class="text-slate-500 mb-6 text-sm leading-relaxed">{{ popup.message }}</p>
                <div class="flex gap-2 justify-center">
                    <button v-if="popup.type === 'confirm'" @click="closePopup(false)" class="btn btn-sm h-10 bg-white border-slate-200 text-slate-600 hover:bg-slate-50 flex-1 rounded-lg">ยกเลิก</button>
                    <button @click="closePopup(true)" class="btn btn-sm h-10 border-none text-white flex-1 rounded-lg shadow-md"
                        :class="{
                            'bg-emerald-500 hover:bg-emerald-600 shadow-emerald-200': popup.type === 'success', 
                            'bg-rose-500 hover:bg-rose-600 shadow-rose-200': popup.type === 'error' || popup.type === 'confirm'
                        }">
                        {{ popup.type === 'confirm' ? 'ยืนยัน' : 'ตกลง' }}
                    </button>
                </div>
            </div>
        </div>

    </div>
  </div>
</template>

<style scoped>
.animate-fade-in-up { animation: fadeInUp 0.4s ease-out; }
@keyframes fadeInUp { from { opacity: 0; transform: translateY(20px); } to { opacity: 1; transform: translateY(0); } }
@keyframes bounceIn { 0% { opacity: 0; transform: scale(0.95); } 100% { opacity: 1; transform: scale(1); } }
.animate-bounce-in { animation: bounceIn 0.2s cubic-bezier(0.175, 0.885, 0.32, 1.275); }
</style>