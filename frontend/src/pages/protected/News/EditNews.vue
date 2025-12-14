<script setup lang="ts">
import { ref, computed, watch } from 'vue';
import { getNewsPostById, updateNewsPost } from '@/services/api/news_post';
// import type { UpdateNewsPost } from '@/interfaces/news_post'; // เปิดใช้ถ้ามี Type นี้

const props = defineProps<{ id: number }>();
const emit = defineEmits(['close', 'success']);

const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8080"; 

const isLoading = ref(false);
const isSaving = ref(false);

// --- 1. นิยาม Status Options (ให้ตรงกับ Seed และหน้า Modal) ---
const statusOptions = [
  { id: 1, label: '● เผยแพร่สาธารณะ (Published)', class: 'text-green-600 font-bold' },
  { id: 5, label: '● เฉพาะสมาชิก (Members Only)', class: 'text-indigo-600 font-bold' }, // เพิ่ม ID 5
  { id: 2, label: '● ฉบับร่าง (Draft)', class: 'text-orange-500 font-bold' },
  { id: 3, label: '● จัดเก็บ (Archived)', class: 'text-slate-500 font-bold' }
];

const form = ref({
  title: '',
  post_detail: '',
  admin_id: null as number | null,       // เริ่มต้นเป็น null รอ Fetch
  scholarship_id: null as number | null, // เริ่มต้นเป็น null รอ Fetch
  status_news_id: 1
});

const oldFilePath = ref('');
const newImageFile = ref<File | undefined>(undefined);
const previewImage = ref<string | null>(null);

// Computed: จัดการแสดงผลรูปภาพ
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

// --- Custom Popup System Definition ---
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

/**
 * ฟังก์ชันสำหรับเรียกแสดงหน้าต่างแจ้งเตือน (Modal Dialog)
 * รองรับการทำงานแบบ Await เพื่อรอผลลัพธ์จากผู้ใช้
 */
const showPopup = (type: 'success' | 'error' | 'confirm', title: string, message: string) => {
  return new Promise<boolean>((resolve) => {
    popup.value = {
      isOpen: true,
      type,
      title,
      message,
      resolve
    };
  });
};


const fetchNewsDetail = async () => {
  if (!props.id) return;

  isLoading.value = true;
  try {
    const data = await getNewsPostById(props.id);
    
    // --- 2. Map ข้อมูลจริงเข้า Form (ไม่ Hardcode เลข 1) ---
    form.value = {
      title: data.title || '',
      post_detail: data.post_detail || '',
      admin_id: data.admin_id,             // ใช้ค่าเดิมจาก DB
      scholarship_id: data.scholarship_id, // ใช้ค่าเดิมจาก DB
      status_news_id: data.status_news_id ?? 1
    };
    
    oldFilePath.value = typeof data.file_path === 'string' ? data.file_path : '';
    newImageFile.value = undefined;
    previewImage.value = null;
    
  } catch (error) {
    console.error("Error fetching news:", error);
    showPopup('error', 'เกิดข้อผิดพลาด', 'ไม่สามารถโหลดข้อมูลข่าวได้');
    emit('close');
  } finally {
    isLoading.value = false;
  }
};

watch(() => props.id, (newId) => {
    if (newId) {
        fetchNewsDetail();
    }
}, { immediate: true }); 

const handleFileUpload = (event: Event) => {
  const target = event.target as HTMLInputElement;
  if (target.files && target.files[0]) {
    const file = target.files[0];
    if (file.size > 5 * 1024 * 1024) {
        showPopup('error', 'ข้อผิดพลาด', 'ขนาดไฟล์ต้องไม่เกิน 5MB');
        return;
    }
    newImageFile.value = file;
    previewImage.value = URL.createObjectURL(file);
  }
};

const handleSave = async () => {
  // Validation
  if (!form.value.title?.trim() || !form.value.post_detail?.trim()) {
      showPopup('error', 'ข้อผิดพลาด', 'กรุณากรอกหัวข้อและเนื้อหาข่าวให้ครบถ้วน');
      return;
  }
  
  // เช็ค ID สำคัญ
  if (!form.value.admin_id || !form.value.scholarship_id) {
      showPopup('error', 'ข้อผิดพลาด', 'ข้อมูลผู้ดูแลระบบหรือทุนการศึกษาไม่ถูกต้อง');
      return;
  }

  isSaving.value = true;
  try {
    const fd = new FormData();

    fd.append("title", form.value.title);
    fd.append("post_detail", form.value.post_detail);
    
    // ส่งค่าเดิมกลับไป (หรือจะเปลี่ยน admin_id เป็นคนปัจจุบันที่ login ก็ได้ แล้วแต่ Logic)
    fd.append("admin_id", String(form.value.admin_id));
    fd.append("scholarship_id", String(form.value.scholarship_id));
    
    // ส่ง Status ตามที่เลือก
    fd.append("status_news_id", String(form.value.status_news_id));

    if (newImageFile.value) {
      fd.append("file_path", newImageFile.value);
    }

    await updateNewsPost(props.id, fd);

    showPopup('success', 'ดำเนินการสำเร็จ', 'ข้อมูลข่าวได้รับการอัปเดตเรียบร้อยแล้ว');
    emit('success');

  } catch (error) {
    console.error("Error updating news:", error);
    showPopup('error', 'ข้อผิดพลาด', 'เกิดข้อผิดพลาดในการบันทึกข้อมูลข่าว');
  } finally {
    isSaving.value = false;
  }
};

</script>

<template>
  <div class="w-full font-sans text-slate-800 flex justify-center items-start">
    <div class="w-full max-w-6xl bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden">
        
        <div class="px-8 py-5 border-b border-gray-100 flex flex-row justify-between items-center bg-white">
            <div>
                <h1 class="text-xl font-bold text-[#1e3a8a]">แก้ไขข่าวสาร</h1>
                <span class="text-xs text-gray-400">Ref ID: {{ id }}</span>
            </div>
            <div class="flex gap-3">
                <button @click="emit('close')" class="btn btn-sm btn-ghost text-gray-500" :disabled="isSaving">ยกเลิก</button>
                <button @click="handleSave" class="btn btn-sm bg-[#1e3a8a] text-white border-none hover:bg-blue-900" :disabled="isSaving">
                    <span v-if="isSaving" class="loading loading-spinner loading-xs"></span>
                    {{ isSaving ? 'กำลังบันทึก...' : 'บันทึกการแก้ไข' }}
                </button>
            </div>
        </div>
        
        <div v-if="isLoading" class="flex justify-center py-20">
            <span class="loading loading-spinner loading-lg text-[#1e3a8a]"></span>
        </div>

        <div v-else class="p-8">
            <div class="grid grid-cols-1 lg:grid-cols-3 gap-8 lg:gap-12">
                <div class="lg:col-span-2 space-y-6">
                    <div class="form-control w-full">
                        <label class="label pt-0"><span class="label-text font-bold text-slate-700">หัวข้อข่าว <span class="text-red-500">*</span></span></label>
                        <input v-model="form.title" type="text" class="input input-bordered w-full bg-gray-50/50 focus:border-[#1e3a8a]" />
                    </div>
                    <div class="form-control w-full">
                        <label class="label"><span class="label-text font-bold text-slate-700">เนื้อหาข่าว <span class="text-red-500">*</span></span></label>
                        <textarea v-model="form.post_detail" class="textarea textarea-bordered h-64 p-4 bg-gray-50/50 focus:border-[#1e3a8a] text-base leading-relaxed"></textarea>
                    </div>
                </div>

                <div class="lg:col-span-1 space-y-6">
                    
                    <div class="form-control w-full">
                        <label class="label pt-0 pb-2"><span class="label-text font-bold text-slate-700">สถานะการแสดงผล</span></label>
                        <select v-model="form.status_news_id" class="select select-bordered w-full bg-white text-base focus:border-[#1e3a8a]">
                            <option 
                                v-for="option in statusOptions" 
                                :key="option.id" 
                                :value="option.id"
                                :class="option.class"
                            >
                                {{ option.label }}
                            </option>
                        </select>
                        <div class="mt-2 text-xs text-gray-400 px-1">
                            <span v-if="form.status_news_id === 1">ข่าวจะแสดงให้ทุกคนเห็นทันที</span>
                            <span v-else-if="form.status_news_id === 5">ต้อง Login ก่อนจึงจะเห็นข่าวนี้</span>
                            <span v-else-if="form.status_news_id === 2">ข่าวถูกซ่อนไว้ (เฉพาะ Admin เห็น)</span>
                            <span v-else-if="form.status_news_id === 3">ข่าวนี้ถูกเก็บเข้า Archive แล้ว</span>
                        </div>
                    </div>

                    <div>
                        <label class="label pt-0 pb-2 flex justify-between">
                            <span class="label-text font-bold text-slate-700">รูปภาพหน้าปก</span>
                            <span v-if="newImageFile" class="text-xs text-green-600 font-medium">รอการอัปโหลด...</span>
                        </label>
                        <div class="relative group cursor-pointer">
                            <div class="w-full aspect-[4/3] rounded-xl border-2 border-dashed border-gray-300 bg-slate-50 flex flex-col items-center justify-center overflow-hidden hover:border-[#1e3a8a] hover:bg-blue-50/30 transition-all relative">
                                <img v-if="displayImage" :src="displayImage" class="w-full h-full object-cover" />
                                <div v-else class="flex flex-col items-center justify-center text-gray-400">
                                    <span>คลิกเพื่ออัปโหลด</span>
                                </div>
                                <div class="absolute inset-0 bg-black/40 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity">
                                    <span class="text-white font-medium bg-black/50 px-3 py-1 rounded-full text-sm">เปลี่ยนรูปภาพ</span>
                                </div>
                                <input type="file" accept="image/*" @change="handleFileUpload" class="absolute inset-0 w-full h-full opacity-0 cursor-pointer" />
                            </div>
                            <div class="mt-2 text-xs text-gray-400 text-center">รองรับไฟล์ JPG, PNG ขนาดไม่เกิน 5MB</div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
  </div>
</template>