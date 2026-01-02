<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue';
import { createNewsPost } from '@/services/api/news_post';
import { getScholarships } from '@/services/api/scholarship';
import { Get } from '@/services/api/https';
import { validateNewsPostForm } from '@/validators/newspost_validators';
import type { CreateNewsPostPayload } from '@/interfaces/news_post';
import { Bold, Italic, Underline, Palette, List, ListOrdered, AlignLeft, AlignCenter, AlignRight, AlignJustify } from 'lucide-vue-next';

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

const editorRef = ref<HTMLElement | null>(null);

const applyFormat = (command: string, value: string | undefined = undefined) => {
  document.execCommand(command, false, value);
  updateToolbarState();
};

const lastRange = ref<Range | null>(null);

const applyFontSize = (size: string) => {
  const selection = window.getSelection();
  let range: Range | null = null;
  
  if (selection && selection.rangeCount > 0) {
    const currentRange = selection.getRangeAt(0);
    if (editorRef.value && editorRef.value.contains(currentRange.commonAncestorContainer)) {
      range = currentRange;
    }
  }
  
  if (!range && lastRange.value) {
    range = lastRange.value;
    selection?.removeAllRanges();
    selection?.addRange(range);
  }
  
  if (!range) return;

  const span = document.createElement('span');
  span.style.fontSize = size + 'px';

  if (range.collapsed) {
     // กรณีไม่ได้เลือกข้อความ ให้ใส่ zero-width space เพื่อให้พิมพ์ต่อได้ในขนาดใหม่
     span.innerHTML = '&#8203;';
     range.insertNode(span);
     
     // ย้าย Cursor ไปข้างใน span
     range.selectNodeContents(span);
     range.collapse(false);
     selection?.removeAllRanges();
     selection?.addRange(range);
  } else {
     // กรณีเลือกข้อความ
     const content = range.extractContents();
     span.appendChild(content);
     range.insertNode(span);
     
     // Select กลับไปที่ข้อความเดิม
     range.selectNodeContents(span);
     selection?.removeAllRanges();
     selection?.addRange(range);
  }

  if (editorRef.value) {
      form.value.post_detail = editorRef.value.innerHTML;
      touched.value.post_detail = true;
  }
  updateToolbarState();
};

const activeFormats = ref({
  bold: false,
  italic: false,
  underline: false,
  insertUnorderedList: false,
  insertOrderedList: false,
  justifyLeft: false,
  justifyCenter: false,
  justifyRight: false,
  justifyFull: false,
  heading: 'p',
  fontSize: '16',
});

const updateToolbarState = () => {
  activeFormats.value = {
    bold: document.queryCommandState('bold'),
    italic: document.queryCommandState('italic'),
    underline: document.queryCommandState('underline'),
    insertUnorderedList: document.queryCommandState('insertUnorderedList'),
    insertOrderedList: document.queryCommandState('insertOrderedList'),
    justifyLeft: document.queryCommandState('justifyLeft'),
    justifyCenter: document.queryCommandState('justifyCenter'),
    justifyRight: document.queryCommandState('justifyRight'),
    justifyFull: document.queryCommandState('justifyFull'),
    heading: document.queryCommandValue('formatBlock') || 'p',
    fontSize: '16',
  };

  const selection = window.getSelection();
  if (selection && selection.rangeCount > 0) {
      const range = selection.getRangeAt(0);
      
      // Save range if we are in editor
      if (editorRef.value && editorRef.value.contains(range.commonAncestorContainer)) {
          lastRange.value = range.cloneRange();
          
          const parent = range.commonAncestorContainer.nodeType === 3 
            ? range.commonAncestorContainer.parentElement 
            : range.commonAncestorContainer as HTMLElement;
            
          if (parent) {
              const computedSize = window.getComputedStyle(parent).fontSize;
              if (computedSize) {
                  activeFormats.value.fontSize = parseFloat(computedSize).toString();
              }
          }
      }
  }
};

const onEditorInput = (event: Event) => {
  const target = event.target as HTMLElement;
  form.value.post_detail = target.innerHTML;
  touched.value.post_detail = true;
  updateToolbarState();
};

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
            <div :class="['border rounded-lg overflow-hidden bg-slate-50', (touched.post_detail && errors.post_detail) ? 'border-red-500' : 'border-gray-300 focus-within:border-[#1e3a8a]']">
              <!-- Toolbar -->
              <!-- Toolbar -->
              <div class="flex gap-1 p-2 bg-gray-100 border-b border-gray-200 flex-wrap items-center">

                <select 
                  @change="(e: Event) => applyFontSize((e.target as HTMLSelectElement).value)" 
                  .value="activeFormats.fontSize"
                  class="select select-bordered select-xs w-24 bg-white"
                  title="ขนาดตัวอักษร"
                >
                  <option value="12">12 px</option>
                  <option value="14">14 px</option>
                  <option value="16">16 px</option>
                  <option value="18">18 px</option>
                  <option value="20">20 px</option>
                  <option value="22">22 px</option>
                  <option value="24">24 px</option>
                  <option value="28">28 px</option>
                  <option value="32">32 px</option>
                  <option value="36">36 px</option>
                </select>

                <div class="w-px h-6 bg-gray-300 mx-1 self-center"></div>

                <button 
                  type="button" 
                  @mousedown.prevent="applyFormat('bold')" 
                  :class="['p-1.5 rounded transition-colors', activeFormats.bold ? 'bg-slate-300 text-slate-900' : 'hover:bg-gray-200 text-slate-700']" 
                  title="ตัวหนา"
                >
                  <Bold :size="18" />
                </button>
                <button 
                  type="button" 
                  @mousedown.prevent="applyFormat('italic')" 
                  :class="['p-1.5 rounded transition-colors', activeFormats.italic ? 'bg-slate-300 text-slate-900' : 'hover:bg-gray-200 text-slate-700']" 
                  title="ตัวเอียง"
                >
                  <Italic :size="18" />
                </button>
                <button 
                  type="button" 
                  @mousedown.prevent="applyFormat('underline')" 
                  :class="['p-1.5 rounded transition-colors', activeFormats.underline ? 'bg-slate-300 text-slate-900' : 'hover:bg-gray-200 text-slate-700']" 
                  title="ขีดเส้นใต้"
                >
                  <Underline :size="18" />
                </button>
                <div class="w-px h-6 bg-gray-300 mx-1 self-center"></div>
                <button 
                  type="button" 
                  @mousedown.prevent="applyFormat('insertUnorderedList')" 
                  :class="['p-1.5 rounded transition-colors', activeFormats.insertUnorderedList ? 'bg-slate-300 text-slate-900' : 'hover:bg-gray-200 text-slate-700']" 
                  title="รายการแบบจุด"
                >
                  <List :size="18" />
                </button>
                <button 
                  type="button" 
                  @mousedown.prevent="applyFormat('insertOrderedList')" 
                  :class="['p-1.5 rounded transition-colors', activeFormats.insertOrderedList ? 'bg-slate-300 text-slate-900' : 'hover:bg-gray-200 text-slate-700']" 
                  title="รายการแบบตัวเลข"
                >
                  <ListOrdered :size="18" />
                </button>
                <div class="w-px h-6 bg-gray-300 mx-1 self-center"></div>
                
                <button 
                  type="button" 
                  @mousedown.prevent="applyFormat('justifyLeft')" 
                  :class="['p-1.5 rounded transition-colors', activeFormats.justifyLeft ? 'bg-slate-300 text-slate-900' : 'hover:bg-gray-200 text-slate-700']" 
                  title="ชิดซ้าย"
                >
                  <AlignLeft :size="18" />
                </button>
                <button 
                  type="button" 
                  @mousedown.prevent="applyFormat('justifyCenter')" 
                  :class="['p-1.5 rounded transition-colors', activeFormats.justifyCenter ? 'bg-slate-300 text-slate-900' : 'hover:bg-gray-200 text-slate-700']" 
                  title="กึ่งกลาง"
                >
                  <AlignCenter :size="18" />
                </button>
                <button 
                  type="button" 
                  @mousedown.prevent="applyFormat('justifyRight')" 
                  :class="['p-1.5 rounded transition-colors', activeFormats.justifyRight ? 'bg-slate-300 text-slate-900' : 'hover:bg-gray-200 text-slate-700']" 
                  title="ชิดขวา"
                >
                  <AlignRight :size="18" />
                </button>
                <button 
                  type="button" 
                  @mousedown.prevent="applyFormat('justifyFull')" 
                  :class="['p-1.5 rounded transition-colors', activeFormats.justifyFull ? 'bg-slate-300 text-slate-900' : 'hover:bg-gray-200 text-slate-700']" 
                  title="กระจายแบบเต็ม"
                >
                  <AlignJustify :size="18" />
                </button>

                <div class="w-px h-6 bg-gray-300 mx-1 self-center"></div>

                <div class="flex items-center gap-1 p-1.5 hover:bg-gray-200 rounded cursor-pointer relative group" title="สีตัวอักษร">
                  <Palette :size="18" class="text-slate-700"/>
                  <input type="color" @input="(e: Event) => applyFormat('foreColor', (e.target as HTMLInputElement).value)" class="absolute inset-0 opacity-0 cursor-pointer w-full h-full" />
                </div>
              </div>
              <!-- Editable Area -->
              <div 
                ref="editorRef"
                contenteditable="true"
                class="rich-text-content p-4 min-h-[400px] outline-none max-h-[600px] overflow-y-auto prose max-w-none"
                @input="onEditorInput"
                @keyup="updateToolbarState"
                @mouseup="updateToolbarState"
                @click="updateToolbarState"
                @blur="touched.post_detail = true"
              ></div>
            </div>
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

<style scoped>
.rich-text-content :deep(ul) {
  list-style-type: disc !important;
  padding-left: 1.5rem !important;
  margin: 1em 0 !important;
}

.rich-text-content :deep(ol) {
  list-style-type: decimal !important;
  padding-left: 1.5rem !important;
  margin: 1em 0 !important;
}

.rich-text-content :deep(li) {
  margin: 0.5em 0;
  display: list-item !important; /* Force display list-item just in case */
}

/* Heading Styles */
.rich-text-content :deep(h1) {
  font-size: 2em;
  font-weight: bold;
  margin: 0.67em 0;
}

.rich-text-content :deep(h2) {
  font-size: 1.5em;
  font-weight: bold;
  margin: 0.75em 0;
}

.rich-text-content :deep(h3) {
  font-size: 1.17em;
  font-weight: bold;
  margin: 0.83em 0;
}

.rich-text-content :deep(h4) {
  font-size: 1em;
  font-weight: bold;
  margin: 1.12em 0;
}

.rich-text-content :deep(h5) {
  font-size: 0.83em;
  font-weight: bold;
  margin: 1.5em 0;
}

.rich-text-content :deep(h6) {
  font-size: 0.67em;
  font-weight: bold;
  margin: 1.67em 0;
}
</style>