<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue';
import { createNewsPost } from '@/services/api/news_post';
import { getScholarships } from '@/services/api/scholarship';
import { Get } from '@/services/api/https';
import { validateNewsPostForm } from '@/validators/newspost_validators';
import type { CreateNewsPostPayload } from '@/interfaces/news_post';

const emit = defineEmits(['close', 'success']);
const isSaving = ref(false);
const scholarships = ref<{ id: number; name: string }[]>([]);

const touched = ref<Record<string, boolean>>({
  title: false,
  post_detail: false,
  scholarship_id: false,
  status_news_id: false,
  file_path: false
});

const form = ref({
  title: '',
  post_detail: '',
  file_path: '',
  admin_id: 0,
  scholarship_id: null as number | null,
  status_news_id: 1
});

const newImageFile = ref<File | null>(null);
const previewImage = ref<string | null>(null);
const errors = ref<Record<string, string>>({});
const isValid = ref(false);

const displayImage = computed(() => previewImage.value || form.value.file_path || null);

const performValidation = () => {
  const payload: CreateNewsPostPayload = {
    title: form.value.title,
    post_detail: form.value.post_detail,
    admin_id: form.value.admin_id,
    scholarship_id: form.value.scholarship_id,
    status_news_id: form.value.status_news_id,
    file_path: newImageFile.value || undefined,
  };

  const result = validateNewsPostForm(payload);
  isValid.value = result.valid;
  errors.value = Object.fromEntries(
    Object.entries(result.errors).map(([k, v]) => [k, v || ''])
  );
};

// แก้ไข warning: ลบ newValues, oldValues ออกเพราะไม่ได้ใช้
watch(
  [form, newImageFile],
  () => {
    if (form.value.title !== "") touched.value.title = true;
    if (form.value.post_detail !== "") touched.value.post_detail = true;
    if (form.value.scholarship_id !== null) touched.value.scholarship_id = true;
    if (newImageFile.value !== null) touched.value.file_path = true;
    
    performValidation();
  },
  { deep: true }
);

const handleFileUpload = (event: Event) => {
  const target = event.target as HTMLInputElement;
  if (target.files?.[0]) {
    const file = target.files[0];
    newImageFile.value = file;
    previewImage.value = URL.createObjectURL(file);
    touched.value.file_path = true;
  }
};

const handleSave = async () => {
  if (!isValid.value) return;
  isSaving.value = true;
  try {
    const payload: CreateNewsPostPayload = {
      title: form.value.title,
      post_detail: form.value.post_detail,
      admin_id: form.value.admin_id,
      scholarship_id: form.value.scholarship_id,
      status_news_id: form.value.status_news_id,
      file_path: newImageFile.value || undefined,
    };
    const fd = new FormData();
    Object.entries(payload).forEach(([key, value]) => {
      if (value !== undefined && value !== null) fd.append(key, value as any);
    });
    await createNewsPost(fd);
    emit('success');
  } catch (e) {
    console.error(e);
    alert('บันทึกไม่สำเร็จ');
  } finally {
    isSaving.value = false;
  }
};

onMounted(async () => {
  performValidation();
  
  // ดึง admin profile จาก API เพื่อให้ได้ ID ที่ถูกต้อง
  try {
    const profileRes: any = await Get('/profile/me');
    if (profileRes && profileRes.role === 'admin' && profileRes.data) {
      form.value.admin_id = profileRes.data.ID;
    }
  } catch (e) {
    console.error('Error fetching admin profile:', e);
  }
  
  try {
    const res: any = await getScholarships();
    const data = Array.isArray(res) ? res : (res.data || []);
    scholarships.value = data.map((s: any) => ({ id: s.ID, name: s.scholarship_name }));
  } catch (e) {
    console.error(e);
  }
});
</script>

<template>
  <div class="h-full flex flex-col bg-white">
    <div class="px-6 py-4 border-b border-gray-200 flex justify-between items-center shrink-0">
      <h1 class="text-xl font-bold text-[#1e3a8a]">เพิ่มข่าวสารใหม่</h1>
      <div class="flex gap-2">
        <button @click="emit('close')" class="btn btn-sm btn-ghost text-gray-500 font-normal">ยกเลิก</button>
        <button @click="handleSave" class="btn btn-sm bg-[#1e3a8a] text-white border-none hover:bg-[#152c6f]" :disabled="isSaving || !isValid">
          <span v-if="isSaving" class="loading loading-spinner loading-xs"></span> บันทึก
        </button>
      </div>
    </div>

    <div class="flex-1 overflow-y-auto p-6 custom-scrollbar">
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <div class="lg:col-span-2 space-y-6">
          <div class="form-control w-full">
            <label class="label pt-0 pb-1 font-bold text-slate-700">หัวข้อข่าว *</label>
            <input 
              v-model="form.title" 
              @blur="touched.title = true"
              type="text" 
              :class="['input input-bordered w-full bg-slate-50 focus:border-[#1e3a8a]', (touched.title && errors.title) ? 'border-red-500' : '']" 
              placeholder="ระบุหัวข้อ..." 
            />
            <p v-if="touched.title && errors.title" class="text-xs text-red-500 mt-1">{{ errors.title }}</p>
          </div>

          <div class="form-control w-full">
            <label class="label pb-1 font-bold text-slate-700">เนื้อหาข่าว *</label>
            <textarea 
              v-model="form.post_detail" 
              @blur="touched.post_detail = true"
              :class="['textarea textarea-bordered min-h-[400px] bg-slate-50 focus:border-[#1e3a8a]', (touched.post_detail && errors.post_detail) ? 'border-red-500' : '']" 
              placeholder="รายละเอียด..."
            ></textarea>
            <p v-if="touched.post_detail && errors.post_detail" class="text-xs text-red-500 mt-1">{{ errors.post_detail }}</p>
          </div>
        </div>

        <div class="lg:col-span-1 space-y-6">
          <div class="p-5 border border-gray-200 rounded-xl bg-white shadow-sm space-y-4">
            <div class="form-control w-full">
              <label class="label pt-0 pb-1 font-bold text-slate-700">ทุนการศึกษา</label>
              <select v-model.number="form.scholarship_id" @change="touched.scholarship_id = true" :class="['select select-bordered w-full', (touched.scholarship_id && errors.scholarship_id) ? 'border-red-500' : '']">
                <option :value="null">เลือกทุนการศึกษา</option>
                <option v-for="s in scholarships" :key="s.id" :value="s.id">{{ s.name }}</option>
              </select>
              <p v-if="touched.scholarship_id && errors.scholarship_id" class="text-xs text-red-500 mt-1">{{ errors.scholarship_id }}</p>
            </div>

            <div class="form-control w-full">
              <label class="label pt-0 pb-1 font-bold text-slate-700">รูปภาพหน้าปก</label>
              <div 
                class="relative group aspect-[4/3] rounded-lg border-2 border-dashed border-gray-300 flex items-center justify-center overflow-hidden hover:border-[#1e3a8a] cursor-pointer"
                :class="(touched.file_path && errors.file_path) ? 'border-red-400 bg-red-50' : ''"
              >
                <img v-if="displayImage" :src="displayImage" class="w-full h-full object-cover" />
                <span v-else class="text-sm text-gray-400">คลิกเพื่ออัปโหลด</span>
                <input type="file" accept="image/*" @change="handleFileUpload" class="absolute inset-0 opacity-0 cursor-pointer" />
              </div>
              <p v-if="touched.file_path && errors.file_path" class="text-xs text-red-500 mt-2 text-center">{{ errors.file_path }}</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>