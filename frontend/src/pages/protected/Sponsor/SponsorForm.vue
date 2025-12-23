<script setup lang="ts">
import { ref, watch, onMounted, nextTick, computed } from 'vue';
import type { PropType } from 'vue';
import type { IndustryResponse } from '@/interfaces/sponsor';
import type { SponsorPayload, ContactPayload } from '@/interfaces/sponsor';
import { IndustryService } from '@/services/sponsor/industry';
import { validateSponsorForm } from '@/validators/sponsorValidator';
import { useFocusTrap } from '../../../hooks/sponsor/useFocusTrap';

// ---------------- Props / Emits ----------------
const props = defineProps({
  isOpen: { type: Boolean as PropType<boolean>, default: false },
  loading: { type: Boolean as PropType<boolean>, default: false },
  initialData: { type: Object as PropType<Record<string, any> | null>, default: null },
});
const emit = defineEmits(['update:isOpen', 'close', 'create']);

// ---------------- Local state ----------------
const industries = ref<IndustryResponse[]>([]);
const industryLoading = ref<boolean>(true);

const form = ref<SponsorPayload>({
  company_name: '',
  website: null,
  industry_id: null,
  status: 'active',
  description: null,
  contacts: [] as ContactPayload[],
});

const errors = ref<Record<string, any>>({});

// ---------------- useFocusTrap ----------------
const isOpenRef = computed(() => props.isOpen);
const { dialogId, focusFirstElement, onBackdropClick } = useFocusTrap(isOpenRef, {
  onClose: () => {
    emit('update:isOpen', false);
    emit('close');
  },
});

// ---------------- Fetch industries ----------------
onMounted(async () => {
  industryLoading.value = true;
  try {
    const res = await IndustryService.getAll();
    industries.value = res ?? [];
  } catch (err) {
    console.error('โหลดอุตสาหกรรมผิดพลาด:', err);
    industries.value = [];
  } finally {
    industryLoading.value = false;
  }
});

// ---------------- Watch isOpen (reset form when opened) ----------------
watch(
  () => props.isOpen,
  (open) => {
    if (open) {
      resetForm();
      // composable handles focus & body scroll
      nextTick(() => {
        // ถ้ามี error จาก server หรือ validation ก่อนหน้า ให้โฟกัส field แรกที่มี error
        if (Object.keys(errors.value || {}).length > 0) {
          focusFirstElement(true);
        }
      });
    } else {
      // closed -> nothing extra here
    }
  }
);

// ---------------- Helpers ----------------
function resetForm() {
  form.value = {
    company_name: props.initialData?.company_name ?? '',
    website: props.initialData?.website ?? null,
    industry_id: props.initialData?.industry_id ?? null,
    status: props.initialData?.status ?? 'active',
    description: props.initialData?.description ?? null,
    contacts: (props.initialData?.contacts ?? []) as ContactPayload[],
  };
  errors.value = {};
}

function close() {
  emit('update:isOpen', false);
  emit('close');
}

function addContact() {
  form.value.contacts?.push({ name: '', email: '', phone: '', position: '' } as ContactPayload);
}

function removeContact(idx: number) {
  form.value.contacts?.splice(idx, 1);
}

function submit() {
  const { valid, errors: vErrors } = validateSponsorForm(form.value);
  errors.value = vErrors;
  if (!valid) {
    nextTick(() => focusFirstElement(true));
    return;
  }

  const payload: SponsorPayload = {
    company_name: String(form.value.company_name).trim(),
    website: form.value.website ? String(form.value.website).trim() : null,
    industry_id: form.value.industry_id ?? null,
    status: String(form.value.status ?? 'active'),
    description: form.value.description ? String(form.value.description).trim() : null,
    contacts: (form.value.contacts ?? []).map(c => ({
      name: String((c as ContactPayload).name ?? '').trim(),
      email: String((c as ContactPayload).email ?? '').trim(),
      phone: String((c as ContactPayload).phone ?? '').trim(),
      position: String((c as ContactPayload).position ?? '').trim(),
    })) as ContactPayload[],
  };

  emit('create', payload);
}
</script>

<template>
  <div>
    <teleport to="body">
    <div
      v-if="props.isOpen"
      class="fixed inset-0 z-[200] flex items-center justify-center bg-black/50 backdrop-blur-sm p-4"
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
            <h2 :id="`${dialogId}-title`" class="text-xl font-bold text-[#1e3a8a]">เพิ่มบริษัทใหม่</h2>
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

                <!-- Industry -->
                <div>
                  <label class="block text-sm text-gray-600 mb-1">อุตสาหกรรม</label>

                  <select
                    v-model="form.industry_id"
                    class="select select-bordered w-full bg-white"
                    :disabled="industryLoading"
                  >
                    <option :value="null">-- เลือกอุตสาหกรรม --</option>
                    <option v-for="i in industries" :key="i.ID" :value="i.ID">{{ i.name }}</option>
                  </select>

                  <p v-if="industryLoading" class="text-xs text-slate-400 mt-1">กำลังโหลดรายการ...</p>
                </div>

                <!-- description -->
                <div>
                  <label class="block text-sm text-gray-600 mb-1">คำอธิบาย</label>
                  <textarea v-model="form.description" rows="3" class="textarea textarea-bordered w-full"></textarea>
                </div>

                <!-- Contact -->
                <div class="border rounded p-4 bg-white">
                  <div class="flex items-center justify-between mb-3">
                    <h3 class="text-sm font-medium">ผู้ติดต่อ</h3>

                    <button 
                      type="button" 
                      class="btn btn-sm bg-white border border-gray-300 text-gray-700 hover:bg-gray-100 hover:border-gray-400 flex items-center justify-center gap-2 rounded-full px-5 h-10 w-full md:w-auto shadow-sm transition-all duration-150" 
                      @click="addContact"
                    >
                      เพิ่มผู้ติดต่อ
                    </button>
                  </div>
                  
                  <div v-if="form.contacts?.length === 0" 
                    class="flex flex-col items-center gap-2"
                  >
                    <video 
                      src="../../../assets/animation/addContact.webm"
                      class="w-20 h-20 object-contain"
                      autoplay
                      muted
                      loop
                      playsinline
                    ></video>

                    <span class="text-xs text-gray-500">ยังไม่มีผู้ติดต่อ</span>
                  </div>

                  <div class="space-y-4">
                    <template v-for="(c, idx) in form.contacts" :key="idx">
                      <div class="grid grid-cols-12 gap-2 items-start">

                        <input
                          v-model="c.name"
                          class="input input-sm input-bordered col-span-12 md:col-span-4"
                          placeholder="ชื่อ"
                          :data-has-error="errors.contacts?.[idx]?.name ? 'true' : null"
                        />

                        <input
                          v-model="c.email"
                          class="input input-sm input-bordered col-span-12 md:col-span-3"
                          placeholder="email"
                          :data-has-error="errors.contacts?.[idx]?.email ? 'true' : null"
                        />

                        <input
                          v-model="c.position"
                          class="input input-sm input-bordered col-span-12 md:col-span-3"
                          placeholder="ตำแหน่ง"
                          :data-has-error="errors.contacts?.[idx]?.position ? 'true' : null"
                        />

                        <input
                          v-model="c.phone"
                          class="input input-sm input-bordered col-span-11 md:col-span-1"
                          placeholder="โทร."
                          :data-has-error="errors.contacts?.[idx]?.phone ? 'true' : null"
                        />

                        <button
                          type="button"
                          class="btn btn-sm btn-error btn-outline col-span-1"
                          @click="removeContact(idx)"
                        >
                          ลบ
                        </button>

                        <!-- Contact Errors -->
                        <p
                          v-if="errors.contacts?.[idx]"
                          class="text-xs text-red-500 col-span-12"
                        >
                          {{ Object.values(errors.contacts[idx])[0] }}
                        </p>

                      </div>
                    </template>
                  </div>
                </div>
              </form>
            </div>
          </div>
        </div>

        <!-- Footer -->
        <div class="px-6 py-4 border-t bg-slate-50 flex items-center justify-end gap-2">
          <button 
            class="btn btn-sm btn-ghost text-gray-700
                hover:bg-red-100 hover:border-red-400 hover:text-red-700
                flex items-center justify-center gap-2
                rounded-full px-5 h-10 w-full md:w-auto
                transition-all duration-150" 
            @click="close" 
            type="button"
          >
            ยกเลิก
          </button>

          <button 
            class="btn btn-sm bg-white border border-blue-300 text-gray-700
                hover:bg-blue-100 hover:border-blue-400 hover:text-blue-700
                flex items-center justify-center gap-2
                rounded-full px-5 h-10 w-full md:w-auto shadow-sm
                transition-all duration-150" 
            @click="submit" 
            :disabled="props.loading"
          >
            <span v-if="props.loading" class="loading loading-spinner" aria-hidden="true"></span>
            <span v-else>สร้างบริษัท</span>
          </button>
        </div>
      </div>
    </div>
    </teleport>
  </div>
</template>

<style scoped>
/* small pop-in animation */
@keyframes pop-in {
  from { opacity: 0; transform: translateY(6px) scale(.995); }
  to { opacity: 1; transform: translateY(0) scale(1); }
}
.animate-pop-in {
  animation: pop-in .12s cubic-bezier(.2,.9,.3,1);
}

/* sr-only for accessibility if needed */
.sr-only {
  position: absolute;
  width: 1px;
  height: 1px;
  padding: 0;
  margin: -1px;
  overflow: hidden;
  clip: rect(0,0,0,0);
  white-space: nowrap;
  border: 0;
}
</style>
