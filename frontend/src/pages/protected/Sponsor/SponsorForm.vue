<script setup lang="ts">
import { ref, watch, onMounted, onBeforeUnmount, nextTick } from 'vue';
import type { PropType } from 'vue';
import type { IndustryResponse } from '../../../interfaces/sponsor';
import { IndustryService } from "../../../services/sponsor/industry"

// -------------------- Props --------------------
const props = defineProps({
  isOpen: { type: Boolean as PropType<boolean>, default: false },
  loading: { type: Boolean as PropType<boolean>, default: false },
  initialData: { type: Object as PropType<Record<string, any> | null>, default: null },
});

// -------------------- Emits --------------------
const emit = defineEmits(['update:isOpen', 'close', 'create']);

// -------------------- Form State --------------------
const industries = ref<IndustryResponse[]>([]);

const form = ref({
  company_name: '',
  website: '',
  industry_id: null as number | null,
});

onMounted(async () => {
  try {
    const res = await IndustryService.getAll();
    industries.value = res;
  } catch (err) {
    console.error("โหลดอุตสาหกรรมผิดพลาด:", err);
  }
})

const errors = ref<{ [k: string]: string }>({});

// ref สำหรับ modal container
const modalRef = ref<HTMLElement | null>(null);

// เก็บ element ที่ focus ก่อนเปิด popup เพื่อคืนค่าเมื่อปิด
let previouslyFocusedElement: HTMLElement | null = null;

// id สำหรับ aria-labelledby
const dialogId = `sponsor-form-modal-${Math.random().toString(36).slice(2, 9)}`;

// -------------------- Watch: เปิด/ปิด Popup --------------------
watch(
  () => props.isOpen,
  async (open) => {
    if (open) {
      resetForm();
      document.body.style.overflow = 'hidden';
      previouslyFocusedElement = document.activeElement as HTMLElement | null;
      await nextTick();
      focusFirstElement();
    } else {
      document.body.style.overflow = '';
      if (previouslyFocusedElement?.focus) {
        previouslyFocusedElement.focus();
      }
      previouslyFocusedElement = null;
    }
  }
);

// -------------------- Reset Form --------------------
function resetForm() {
  form.value = {
    company_name: props.initialData?.company_name ?? '',
    website: props.initialData?.website ?? '',
    industry_id: props.initialData?.industry_id ?? null,
  };
  errors.value = {};
}

// -------------------- Validation --------------------
function validate() {
  errors.value = {};

  if (!form.value.company_name || form.value.company_name.trim().length < 2) {
    errors.value.company_name = 'โปรดระบุชื่อบริษัท (อย่างน้อย 2 ตัวอักษร)';
  }

  if (form.value.website && !/^https?:\/\//i.test(form.value.website)) {
    errors.value.website = 'ใส่ URL ให้มี http:// หรือ https://';
  }

  return Object.keys(errors.value).length === 0;
}

// -------------------- Close & Submit --------------------
function close() {
  emit('update:isOpen', false);
  emit('close');
}

function submit() {
  if (!validate()) {
    focusFirstElement(true);
    return;
  }

  emit('create', {
    company_name: form.value.company_name.trim(),
    website: form.value.website?.trim() || null,
    industry_id: form.value.industry_id ?? null,
  });
}

// -------------------- Focus / Keyboard handling --------------------
// Selector สำหรับ element ที่สามารถ focus ได้
const FOCUSABLE_SELECTOR = [
  'a[href]',
  'button:not([disabled])',
  'input:not([disabled]):not([type="hidden"])',
  'textarea:not([disabled])',
  'select:not([disabled])',
  '[tabindex]:not([tabindex="-1"])'
].join(',');

// Helper: ตรวจ element ที่เห็นได้ (ไม่ hidden)
function isVisible(el: HTMLElement) {
  // offsetParent เป็นวิธีที่ค่อนข้าง reliable ในกรณีทั่วไป
  return !!(el.offsetParent || el.getClientRects().length);
}

// focus ช่องแรก (หรือช่องที่มี error)
function focusFirstElement(focusError = false) {
  const container = modalRef.value;
  if (!container) return;

  if (focusError) {
    const errorField = container.querySelector<HTMLElement>('[data-has-error="true"]');
    if (errorField && typeof errorField.focus === 'function') {
      errorField.focus();
      return;
    }
  }

  const el = container.querySelector<HTMLElement>(FOCUSABLE_SELECTOR);
  if (el && typeof el.focus === 'function') el.focus();
}

// onKeydown: Escape ปิด, Tab/Shift+Tab ทำ focus-trap
function onKeydown(e: KeyboardEvent) {
  if (!props.isOpen) return;

  if (e.key === 'Escape') {
    e.preventDefault();
    close();
    return;
  }

  if (e.key !== 'Tab') return;

  const container = modalRef.value;
  if (!container) return;

  const nodeList = container.querySelectorAll<HTMLElement>(FOCUSABLE_SELECTOR);
  const focusable = Array.from(nodeList).filter((el) => {
    if (el.hasAttribute('disabled')) return false;
    if (el.getAttribute('tabindex') === '-1') return false;
    return isVisible(el);
  });

  // ถ้าไม่มี element ให้เลิกทำ
  if (focusable.length === 0) return;

  // ---------- type-safe narrow + assertion ----------
  // ภายในบล็อกนี้ focusable.length > 0 แล้ว ดังนั้นการใช้ [0] และ [length-1] ปลอดภัย
  const first = focusable[0] as HTMLElement;
  const last = focusable[focusable.length - 1] as HTMLElement;
  // --------------------------------------------------

  if (e.shiftKey) {
    if (document.activeElement === first) {
      e.preventDefault();
      last.focus();
    }
    return;
  }

  if (document.activeElement === last) {
    e.preventDefault();
    first.focus();
  }
}

// Backdrop click -> ปิด เมื่อคลิกนอก dialog card
function onBackdropClick(e: MouseEvent) {
  if (e.target === e.currentTarget) {
    close();
  }
}

// -------------------- Lifecycle --------------------
onMounted(() => window.addEventListener('keydown', onKeydown));
onBeforeUnmount(() => window.removeEventListener('keydown', onKeydown));
</script>

<template>
  <teleport to="body">
    <div
      v-if="isOpen"
      class="fixed inset-0 z-[200] flex items-center justify-center bg-black/50 backdrop-blur-sm p-4 transition-opacity"
      @click="onBackdropClick"
      data-theme="light"
    >
      <div
        ref="modalRef"
        role="dialog"
        :aria-modal="true"
        :aria-labelledby="`${dialogId}-title`"
        class="bg-white w-full max-w-6xl h-[90vh] rounded-2xl shadow-2xl flex flex-col overflow-hidden animate-pop-in"
      >
        <!-- Header -->
        <div class="px-6 py-4 border-b flex items-center justify-between bg-slate-50">
          <div>
            <h2 :id="`${dialogId}-title`" class="text-xl font-bold text-[#1e3a8a] flex item-center gap-2">เพิ่มบริษัทใหม่</h2>
          </div>

          <button class="btn btn-circle btn-ghost btn-sm" @click="close" aria-label="ปิด">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24"
              stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <!-- Body -->
        <div class="flex-1 overflow-y-auto p-6 bg-[#f8fafc]">
          <div class="card bg-white shadow-sm border border-gray-100">
            <div class="card-body p-5">

              <form @submit.prevent="submit" class="grid grid-cols-1 gap-4">
                <!-- Company Name -->
                <div>
                  <label class="block text-sm text-gray-600 mb-1">
                    ชื่อบริษัท <span class="text-red-500">*</span>
                  </label>
    
                  <input
                    v-model="form.company_name"
                    type="text"
                    class="input input-bordered w-full"
                    placeholder="ชื่อบริษัท"
                    :aria-invalid="errors.company_name ? 'true' : 'false'"
                    :data-has-error="errors.company_name ? 'true' : null"
                  />
    
                  <p v-if="errors.company_name" class="text-xs text-red-500 mt-1">
                    {{ errors.company_name }}
                  </p>
                </div>
    
                <!-- Website -->
                <div>
                  <label class="block text-sm text-gray-600 mb-1">เว็บไซต์</label>
                  <div class="relative">
                    <!-- Icon -->
                     <svg class="absolute left-3 top-1/2 -translate-y-1/2 h-[1em] opacity-50" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
                      <g stroke-linejoin="round" stroke-linecap="round" stroke-width="2.5" fill="none" stroke="currentColor">
                        <path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71"></path>
                        <path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71"></path>
                      </g>
                    </svg>

                    <input
                      v-model="form.website"
                      type="url"
                      placeholder="https://example.com"
                      class="input input-bordered w-full placeholder:text-gray-400 pl-10"
                      :aria-invalid="errors.website ? 'true' : 'false'"
                      :data-has-error="errors.website ? 'true' : null"
                    />

                    <p v-if="errors.website" class="text-xs text-red-500 mt-1">
                      {{ errors.website }}
                    </p>
                  </div>
                </div>
    
                <!-- Industry ID -->
                <div>
                  <label class="block text-sm text-gray-600 mb-1">อุตสาหกรรม (ID)</label>
    
                  <select
                    v-model="form.industry_id"
                    class="select select-bordered w-full bg-white"
                  >
                    <option :value="null">-- เลือกอุตสาหกรรม --</option>
                    <option v-for="i in industries" :key="i.ID" :value="i.ID">{{ i.name }}</option>
                  </select>
                  
                </div>
              </form>
            </div>
            
          </div>
        </div>

        <!-- Footer -->
        <div class="px-6 py-4 border-t bg-slate-50 flex items-center justify-end gap-2">
          <button class="btn btn-ghost" @click="close" type="button">ยกเลิก</button>

          <button class="btn btn-primary" @click="submit" :disabled="loading">
            <span v-if="loading" class="loading loading-spinner" aria-hidden="true"></span>
            <span v-else>สร้างบริษัท</span>
          </button>
        </div>
      </div>
    </div>
  </teleport>
</template>

<style scoped>
/* sr-only สำหรับ accessibility (พร้อมใช้ถ้าต้องการ) */
.sr-only {
  position: absolute;
  width: 1px;
  height: 1px;
  padding: 0;
  margin: -1px;
  overflow: hidden;
  clip: rect(0, 0, 0, 0);
  white-space: nowrap;
  border: 0;
}
</style>
