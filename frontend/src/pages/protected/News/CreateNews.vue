<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { createNewsPost } from '@/services/api/news_post';
import { getScholarships } from '@/services/api/scholarship';
import type { ScholarshipResponse } from '@/interfaces/scholarship';

const emit = defineEmits(['close', 'success']);
const isSaving = ref(false);

// --- List ทุน ---
const scholarships = ref<{id: number, name: string}[]>([]);

// --- Form state ---
const form = ref({
  title: '',
  post_detail: '',
  file_path: '',
  admin_id: null,          // รับจาก props หรือ session
  scholarship_id: null,    // เลือกจาก dropdown
  status_news_id: null     // เลือกจาก select
});

const newImageFile = ref<File | null>(null);
const previewImage = ref<string | null>(null);

const displayImage = computed(() => {
  if (previewImage.value) return previewImage.value;
  if (form.value.file_path) return form.value.file_path;
  return null;
});

const handleFileUpload = (event: Event) => {
  const target = event.target as HTMLInputElement;
  if (target.files && target.files[0]) {
    const file = target.files[0];
    if (file.size > 5 * 1024 * 1024) {
      //alert('ขนาดไฟล์ต้องไม่เกิน 5MB');
      return;
    }
    newImageFile.value = file;
    previewImage.value = URL.createObjectURL(file);
  }
};

const handleSave = async () => {
  // 1. Validation พื้นฐาน
  if (!form.value.title || !form.value.post_detail) {
    //alert('กรุณากรอกหัวข้อและเนื้อหาข่าว');
    return;
  }

  // 2. ตรวจสอบค่า ID (กันส่ง null/undefined ไปให้ Server)
  // สมมติ: ถ้าไม่มี admin_id ให้ลองดึงจาก localStorage หรือใส่ค่า Default ชั่วคราวเพื่อ Test
  const safeAdminId = form.value.admin_id ? String(form.value.admin_id) : '1'; 
  const safeScholarshipId = form.value.scholarship_id ? String(form.value.scholarship_id) : '1'; 

  isSaving.value = true;

  try {
    const formData = new FormData();
    formData.append('title', form.value.title.trim());
    formData.append('post_detail', form.value.post_detail.trim());
    
    // ส่งค่า ID ที่เช็คแล้วว่าไม่ว่าง
    formData.append('admin_id', safeAdminId);
    formData.append('scholarship_id', safeScholarshipId);
    formData.append('status_news_id', String(form.value.status_news_id));
    
    // จัดการไฟล์
    if (newImageFile.value) {
      formData.append('file_path', newImageFile.value);
    }

    // *** Debug: ดูว่าส่งอะไรไปบ้างใน Console ***
    console.log('--- Sending FormData ---');
    for (const pair of formData.entries()) {
        console.log(`${pair[0]}: ${pair[1]}`);
    }

    // ยิง API
    await createNewsPost(formData);

    //alert('บันทึกข้อมูลสำเร็จ');
    emit('success');
    emit('close');

  } catch (error: any) {
    console.error('Save Failed:', error);
    
    // ดักจับ Error Response จาก Server มาแสดง
    if (error.response) {
        console.log('Server Response Data:', error.response.data);
        //alert(`เกิดข้อผิดพลาด: ${JSON.stringify(error.response.data)}`);
    } else {
        //alert('เกิดข้อผิดพลาดในการเชื่อมต่อ Server (500)');
    }
  } finally {
    isSaving.value = false;
  }
};

// ดึงทุนจาก API
onMounted(async () => {
  try {
    const res: ScholarshipResponse[] = await getScholarships();
    scholarships.value = res.map(s => ({
      id: s.ID,
      name: s.scholarship_name
    }));
  } catch (error) {
    console.error('ไม่สามารถดึงทุนได้', error);
  }
});


const handleCancel = () => emit('close');
</script>

<template>
  <div class="w-full font-sans text-slate-800 flex justify-center items-start">
    <div class="w-full max-w-6xl bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden">
      
      <!-- Header -->
      <div class="px-8 py-5 border-b border-gray-100 flex flex-row justify-between items-center bg-white">
        <div>
          <h1 class="text-xl font-bold text-[#1e3a8a] flex items-center gap-2">เพิ่มข่าวสารใหม่</h1>
        </div>
        <div class="flex gap-3">
          <button @click="handleCancel" class="btn btn-sm btn-ghost text-gray-500">ยกเลิก</button>
          <button @click="handleSave" class="btn btn-sm bg-[#1e3a8a] hover:bg-[#152c6f] text-white border-none px-6" :disabled="isSaving">
            <span v-if="isSaving" class="loading loading-spinner loading-xs"></span>
            บันทึก
          </button>
        </div>
      </div>

      <!-- Form -->
      <div class="p-8">
        <div class="grid grid-cols-1 lg:grid-cols-3 gap-8 lg:gap-12">
          
          <!-- Left: Title + Detail -->
          <div class="lg:col-span-2 space-y-6">
            <div class="form-control w-full">
              <label class="label pt-0"><span class="label-text font-bold text-slate-700 text-base">หัวข้อข่าว <span class="text-red-500">*</span></span></label>
              <input v-model="form.title" type="text" placeholder="ระบุหัวข้อข่าว..." class="input input-bordered w-full h-11 text-base bg-gray-50/50" />
            </div>

            <div class="form-control w-full">
              <label class="label"><span class="label-text font-bold text-slate-700 text-base">เนื้อหาข่าว <span class="text-red-500">*</span></span></label>
              <textarea v-model="form.post_detail" class="textarea textarea-bordered h-64 text-base leading-relaxed p-4 bg-gray-50/50 resize-none" placeholder="รายละเอียด..."></textarea>
            </div>
          </div>

          <!-- Right: Scholarship + Status + Image -->
          <div class="lg:col-span-1 space-y-6">

            <!-- Scholarship Dropdown -->
            <div class="form-control w-full">
              <label class="label pt-0 pb-2"><span class="label-text font-bold text-slate-700 text-base">ทุนการศึกษา <span class="text-red-500">*</span></span></label>
              <select v-model="form.scholarship_id" class="select select-bordered w-full text-base bg-white focus:border-[#1e3a8a] focus:ring-1 focus:ring-[#1e3a8a]">
                <option value="" disabled>เลือกทุนการศึกษา</option>
                <option v-for="scholarship in scholarships" :key="scholarship.id" :value="scholarship.id">
                  {{ scholarship.name }}
                </option>
              </select>
            </div>

            <!-- Status Dropdown -->
            <div class="form-control w-full">
              <label class="label pt-0 pb-2"><span class="label-text font-bold text-slate-700 text-base">สถานะการแสดงผล <span class="text-red-500">*</span></span></label>
              <select v-model="form.status_news_id" class="select select-bordered w-full text-base bg-white focus:border-[#1e3a8a] focus:ring-1 focus:ring-[#1e3a8a]">
                <option value="" disabled>เลือกสถานะข่าว</option>
                <option :value="1" class="text-green-600 font-bold">● เผยแพร่ (Published)</option>
                <option :value="2" class="text-orange-500 font-bold">● ฉบับร่าง (Draft)</option>
                <option :value="3" class="text-red-600 font-bold">● ซ่อน (Hidden)</option>
              </select>
              <div class="mt-2 text-xs text-gray-400">
                <span v-if="form.status_news_id === 1">ข่าวจะแสดงบนหน้าเว็บไซต์ทันที</span>
                <span v-else>ข่าวจะถูกซ่อนไว้เฉพาะผู้ดูแลระบบ</span>
              </div>
            </div>

            <!-- Cover Image Upload -->
            <div>
              <label class="label pt-0 pb-2"><span class="label-text font-bold text-slate-700 text-base">รูปภาพหน้าปก</span></label>
              <div class="relative group">
                <div class="w-full aspect-[4/3] rounded-xl border-2 border-dashed border-gray-300 bg-slate-50 overflow-hidden flex flex-col items-center justify-center text-center hover:border-[#1e3a8a] hover:bg-blue-50/30">
                  <img v-if="displayImage" :src="displayImage" class="w-full h-full object-cover" />
                  <div v-else class="p-6 flex flex-col items-center">
                    <span class="text-sm font-semibold text-gray-600">อัปโหลดรูปภาพ</span>
                  </div>
                  <input type="file" accept="image/*" @change="handleFileUpload" class="absolute inset-0 w-full h-full opacity-0 cursor-pointer" />
                </div>
              </div>
            </div>

          </div>

        </div>
      </div>

    </div>
  </div>
</template>
