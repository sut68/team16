<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { createNewsPost } from '@/services/api/news_post';
import { getScholarships } from '@/services/api/scholarship';

const emit = defineEmits(['close', 'success']);
const isSaving = ref(false);
const scholarships = ref<{id: number, name: string}[]>([]);

const statusOptions = [
  { id: 1, label: 'สาธารณะ (Public)', class: 'text-emerald-600' },
  { id: 5, label: 'เฉพาะสมาชิก (Members)', class: 'text-indigo-600' },
  { id: 2, label: 'ฉบับร่าง (Draft)', class: 'text-orange-500' },
  { id: 3, label: 'จัดเก็บ (Archived)', class: 'text-slate-500' }
];

const form = ref({ title: '', post_detail: '', file_path: '', admin_id: 1, scholarship_id: null as number | null, status_news_id: 1 });
const newImageFile = ref<File | null>(null);
const previewImage = ref<string | null>(null);

const displayImage = computed(() => previewImage.value || form.value.file_path || null);

const handleFileUpload = (event: Event) => {
  const target = event.target as HTMLInputElement;
  if (target.files?.[0]) {
    const file = target.files[0];
    newImageFile.value = file;
    previewImage.value = URL.createObjectURL(file);
  }
};

const handleSave = async () => {
  if (!form.value.title || !form.value.post_detail) return alert('กรุณากรอกข้อมูลให้ครบ');
  isSaving.value = true;
  try {
    const fd = new FormData();
    fd.append('title', form.value.title);
    fd.append('post_detail', form.value.post_detail);
    fd.append('admin_id', String(form.value.admin_id));
    if (form.value.scholarship_id) fd.append('scholarship_id', String(form.value.scholarship_id));
    fd.append('status_news_id', String(form.value.status_news_id));
    if (newImageFile.value) fd.append('file_path', newImageFile.value);
    await createNewsPost(fd);
    emit('success');
  } catch (error) { console.error(error); alert('บันทึกไม่สำเร็จ'); } 
  finally { isSaving.value = false; }
};

onMounted(async () => {
  try {
    const res: any = await getScholarships();
    const data = Array.isArray(res) ? res : (res.data || []);
    scholarships.value = data.map((s: any) => ({ id: s.ID, name: s.scholarship_name }));
  } catch (e) { console.error(e); }
});
</script>

<template>
  <div class="h-full flex flex-col bg-white">
    <div class="px-6 py-4 border-b border-gray-200 flex justify-between items-center shrink-0">
      <div>
        <h1 class="text-xl font-bold text-[#1e3a8a]">เพิ่มข่าวสารใหม่</h1>
        <p class="text-xs text-gray-500">Create New Post</p>
      </div>
      <div class="flex gap-2">
        <button @click="emit('close')" class="btn btn-sm btn-ghost text-gray-500 font-normal">ยกเลิก</button>
        <button @click="handleSave" class="btn btn-sm bg-[#1e3a8a] text-white border-none hover:bg-[#152c6f]" :disabled="isSaving">
          <span v-if="isSaving" class="loading loading-spinner loading-xs"></span> บันทึก
        </button>
      </div>
    </div>

    <div class="flex-1 overflow-y-auto p-6 custom-scrollbar">
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
        
        <div class="lg:col-span-2 space-y-6">
          <div class="form-control w-full">
            <label class="label pt-0 pb-1"><span class="label-text font-bold text-slate-700">หัวข้อข่าว <span class="text-red-500">*</span></span></label>
            <input v-model="form.title" type="text" class="input input-bordered w-full bg-slate-50 focus:border-[#1e3a8a]" placeholder="ระบุหัวข้อ..." />
          </div>
          <div class="form-control w-full">
            <label class="label pb-1"><span class="label-text font-bold text-slate-700">เนื้อหาข่าว <span class="text-red-500">*</span></span></label>
            <textarea v-model="form.post_detail" class="textarea textarea-bordered min-h-[400px] bg-slate-50 focus:border-[#1e3a8a] text-base leading-relaxed" placeholder="รายละเอียด..."></textarea>
          </div>
        </div>

        <div class="lg:col-span-1 space-y-6">
          <div class="p-5 border border-gray-200 rounded-xl bg-white shadow-sm space-y-4">
            <div class="form-control w-full">
                <label class="label pt-0 pb-1"><span class="label-text font-bold text-slate-700">ทุนการศึกษา</span></label>
                <select v-model="form.scholarship_id" class="select select-bordered w-full bg-slate-50">
                    <option :value="null">เลือกทุนการศึกษา</option>
                    <option v-for="s in scholarships" :key="s.id" :value="s.id">{{ s.name }}</option>
                </select>
            </div>
            <div class="form-control w-full">
                <label class="label pt-0 pb-1"><span class="label-text font-bold text-slate-700">สถานะการเผยแพร่</span></label>
                <select v-model="form.status_news_id" class="select select-bordered w-full bg-slate-50">
                    <option v-for="opt in statusOptions" :key="opt.id" :value="opt.id" :class="opt.class">{{ opt.label }}</option>
                </select>
            </div>
          </div>

          <div class="p-5 border border-gray-200 rounded-xl bg-white shadow-sm">
            <label class="label pt-0 pb-2 flex justify-between">
                <span class="label-text font-bold text-slate-700">รูปภาพหน้าปก</span>
                <span v-if="newImageFile" class="text-xs text-green-600 font-medium">เลือกแล้ว</span>
            </label>
            <div class="relative group w-full aspect-[4/3] rounded-lg border-2 border-dashed border-gray-300 bg-slate-50 flex items-center justify-center overflow-hidden hover:border-[#1e3a8a] transition-colors cursor-pointer">
                <img v-if="displayImage" :src="displayImage" class="w-full h-full object-cover" />
                <div v-else class="text-center text-gray-400">
                    <span class="text-sm block">คลิกเพื่ออัปโหลด</span>
                </div>
                <input type="file" accept="image/*" @change="handleFileUpload" class="absolute inset-0 opacity-0 cursor-pointer" />
            </div>
            <p class="text-xs text-center text-gray-400 mt-2">รองรับไฟล์ JPG, PNG</p>
          </div>
        </div>

      </div>
    </div>
  </div>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar { width: 6px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #cbd5e1; border-radius: 20px; }
</style>