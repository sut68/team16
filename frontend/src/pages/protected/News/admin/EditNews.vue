<script setup lang="ts">
import { ref, computed, watch } from 'vue';
import { getNewsPostById, updateNewsPost } from '@/services/api/news_post';
import { validateNewsPostForm } from '@/validators/newspost_validators';

const props = defineProps<{ id: number }>();
const emit = defineEmits(['close', 'success']);

const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8080"; 

const isLoading = ref(false);
const isSaving = ref(false);
const errors = ref<Record<string, string>>({});

const statusOptions = [
  { id: 1, label: 'เผยแพร่สาธารณะ (Published)', class: 'text-emerald-600 font-medium' },
  { id: 4, label: 'เฉพาะสมาชิก (Members Only)', class: 'text-blue-600 font-medium' },
  { id: 2, label: 'ฉบับร่าง (Draft)', class: 'text-orange-500 font-medium' },
  { id: 3, label: 'จัดเก็บ (Archived)', class: 'text-slate-500 font-medium' }
];

const touched = ref<Record<string, boolean>>({
  title: false,
  post_detail: false,
  file_path: false
});

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

// ฟังก์ชันตรวจสอบความถูกต้อง
const validate = () => {
  const payload = {
    title: form.value.title,
    post_detail: form.value.post_detail,
    status_news_id: form.value.status_news_id,
    file_path: newImageFile.value || (oldFilePath.value ? 'existing' : undefined)
  };

  const result = validateNewsPostForm(payload as any);
  errors.value = Object.fromEntries(
    Object.entries(result.errors).map(([k, v]) => [k, v || ''])
  );
  return result.valid;
};

// ตรวจสอบเมื่อมีการพิมพ์ (แต่จะแสดงสีแดงเฉพาะ field ที่เคย touched แล้ว)
watch([form, newImageFile], () => validate(), { deep: true });

const displayImage = computed(() => {
  if (previewImage.value) return previewImage.value;
  if (oldFilePath.value) {
    const cleanPath = oldFilePath.value.startsWith('/') ? oldFilePath.value.substring(1) : oldFilePath.value;
    return `${API_URL}/${cleanPath}?t=${new Date().getTime()}`; 
  }
  return null;
});

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
    validate();
  } catch (error) {
    console.error(error);
    emit('close');
  } finally { isLoading.value = false; }
};

watch(() => props.id, (newId) => { if (newId) fetchNewsDetail(); }, { immediate: true }); 

const handleFileUpload = (event: Event) => {
  const target = event.target as HTMLInputElement;
  if (target.files?.[0]) {
    newImageFile.value = target.files[0];
    previewImage.value = URL.createObjectURL(target.files[0]);
    touched.value.file_path = true;
  }
};

const handleSave = async () => {
  // เมื่อกดบันทึก ให้ Mark ทุกช่องว่าถูกแตะแล้วเพื่อแสดง Error สีแดง
  Object.keys(touched.value).forEach(key => touched.value[key] = true);
  
  if (!validate()) return; // ถ้าไม่ผ่าน ไม่ต้องทำต่อ

  isSaving.value = true;
  try {
    const fd = new FormData();
    fd.append("title", form.value.title);
    fd.append("post_detail", form.value.post_detail);
    if (form.value.admin_id) fd.append("admin_id", String(form.value.admin_id));
    if (form.value.scholarship_id) fd.append("scholarship_id", String(form.value.scholarship_id));
    fd.append("status_news_id", String(form.value.status_news_id));
    if (newImageFile.value) fd.append("file_path", newImageFile.value);

    await updateNewsPost(props.id, fd);
    emit('success');
  } catch (error) {
    console.error(error);
  } finally { isSaving.value = false; }
};
</script>
<template>
  <div class="w-full h-full flex flex-col bg-white overflow-hidden">
    
    <div class="px-6 py-4 border-b border-slate-100 flex justify-between items-center bg-white shrink-0">
        <div>
            <h1 class="text-xl font-extrabold text-slate-800 tracking-tight">แก้ไขข่าวสาร</h1>
            <p class="text-[11px] text-slate-400 font-mono">Ref ID: {{ id }}</p>
        </div>
        <div class="flex gap-2">
            <button @click="emit('close')" class="btn btn-sm btn-ghost text-slate-500 hover:bg-slate-50" :disabled="isSaving">ยกเลิก</button>
            <button @click="handleSave" class="btn btn-sm bg-[#1e3a8a] hover:bg-[#152c6f] text-white border-none px-6" :disabled="isSaving">
                <span v-if="isSaving" class="loading loading-spinner loading-xs"></span>
                บันทึกการแก้ไข
            </button>
        </div>
    </div>
    
    <div class="flex-1 overflow-y-auto p-6 custom-scrollbar">
        <div v-if="isLoading" class="flex flex-col items-center justify-center h-64">
            <span class="loading loading-dots loading-lg text-[#1e3a8a]"></span>
        </div>

        <div v-else class="grid grid-cols-1 lg:grid-cols-12 gap-8 max-w-[1600px] mx-auto w-full">
            
            <div class="lg:col-span-8 space-y-6">
                <div class="form-control w-full">
                    <label class="label pt-0 pb-1.5"><span class="label-text font-bold text-slate-700 text-sm">หัวข้อข่าว *</span></label>
                    <input 
                        v-model="form.title" 
                        @blur="touched.title = true"
                        type="text" 
                        :class="['input input-bordered w-full h-11 transition-all rounded-lg focus:ring-1 focus:ring-[#1e3a8a]', (touched.title && errors.title) ? 'border-red-500 bg-red-50' : 'bg-slate-50 border-slate-200']" 
                        placeholder="ระบุหัวข้อข่าว..."
                    />
                    <p v-if="touched.title && errors.title" class="text-red-500 text-[11px] mt-1 ml-1">{{ errors.title }}</p>
                </div>

                <div class="form-control w-full">
                    <label class="label pb-1.5"><span class="label-text font-bold text-slate-700 text-sm">เนื้อหาข่าว *</span></label>
                    <textarea 
                        v-model="form.post_detail" 
                        @blur="touched.post_detail = true"
                        :class="['textarea textarea-bordered min-h-[450px] lg:h-[calc(100vh-320px)] p-4 text-base leading-relaxed transition-all rounded-lg resize-none focus:ring-1 focus:ring-[#1e3a8a]', (touched.post_detail && errors.post_detail) ? 'border-red-500 bg-red-50' : 'bg-slate-50 border-slate-200']"
                        placeholder="เขียนรายละเอียดข่าวสาร..."
                    ></textarea>
                    <p v-if="touched.post_detail && errors.post_detail" class="text-red-500 text-[11px] mt-1 ml-1">{{ errors.post_detail }}</p>
                </div>
            </div>

            <div class="lg:col-span-4 space-y-6">
                <div class="space-y-6">
                    <div class="bg-white p-5 rounded-xl border border-slate-100 shadow-sm">
                        <label class="label pt-0 pb-2"><span class="label-text font-bold text-slate-700 text-sm">สถานะ</span></label>
                        <select v-model="form.status_news_id" class="select select-bordered select-sm w-full h-10 bg-white border-slate-200 rounded-lg">
                            <option v-for="option in statusOptions" :key="option.id" :value="option.id" :class="option.class">
                                {{ option.label }}
                            </option>
                        </select>
                        <div class="mt-3 p-3 bg-blue-50/50 rounded-lg border border-blue-100/50 text-[12px] text-slate-500">
                            <p v-if="form.status_news_id === 1">✓ เผยแพร่สู่สาธารณะ</p>
                            <p v-if="form.status_news_id === 4">🔒 เฉพาะสมาชิก</p>
                            <p v-if="form.status_news_id === 2">📝 ฉบับร่าง</p>
                            <p v-if="form.status_news_id === 3">📦 เก็บเข้าคลัง</p>
                        </div>
                    </div>

                    <div class="bg-white p-5 rounded-xl border border-slate-100 shadow-sm">
                        <label class="label pt-0 pb-2 font-bold text-slate-700 text-sm">รูปภาพหน้าปก</label>
                        <div 
                            class="relative group aspect-[16/9] rounded-lg border-2 border-dashed flex items-center justify-center overflow-hidden transition-all duration-300"
                            :class="(touched.file_path && errors.file_path) ? 'border-red-400 bg-red-50' : 'border-slate-300 bg-slate-50 hover:border-[#1e3a8a]'"
                        >
                            <img v-if="displayImage" :src="displayImage" class="w-full h-full object-cover transition-transform group-hover:scale-105" />
                            <div v-else class="flex flex-col items-center text-slate-400 text-xs">
                                <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6 mb-1" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" /></svg>
                                <span>คลิกเพื่อเปลี่ยนรูป</span>
                            </div>
                            <input type="file" accept="image/*" @change="handleFileUpload" class="absolute inset-0 w-full h-full opacity-0 cursor-pointer" />
                        </div>
                        <p v-if="touched.file_path && errors.file_path" class="text-[10px] text-red-500 mt-2 text-center">{{ errors.file_path }}</p>
                    </div>
                </div>
            </div>
        </div>
    </div>
  </div>
</template>

<style scoped>
/* ปรับแต่ง Scrollbar ให้ดูสะอาดตา (เลือกใส่หรือไม่ก็ได้) */
.custom-scrollbar::-webkit-scrollbar { width: 6px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #e2e8f0; border-radius: 10px; }
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background: #cbd5e1; }
</style>