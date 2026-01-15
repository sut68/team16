<script setup lang="ts">
import { ref, watch, nextTick, computed, onMounted, onBeforeUnmount } from 'vue';
import type { PropType } from 'vue';
import type { SponsorResponse, SponsorPayload } from '@/interfaces/sponsor';
import { SponsorService } from '@/services/sponsor/sponsor';
import { validateSponsorForm } from '@/validators/sponsor_validator';
import Swal from 'sweetalert2';

// Hooks
import { useModalFocusTrap } from '@/hooks/sponsor/useFocusTrap';
import { useIndustries } from '@/hooks/sponsor/useIndustries';

// Props & emits
const props = defineProps({
  isOpen: { type: Boolean as PropType<boolean>, default: false },
  loading: { type: Boolean as PropType<boolean>, default: false },
  initialData: { type: Object as PropType<SponsorResponse | null>, default: null },
});

const emit = defineEmits<{
  (e: 'update:isOpen', v: boolean): void;
  (e: 'close'): void;
  (e: 'updated', sponsor: SponsorResponse): void;
}>();

// Form State
const form = ref<Partial<SponsorPayload>>({
  company_name: '',
  website: null,
  industry_id: null,
  status: 'active',
  description: null,
});
const saving = ref(false);
const errors = ref<Record<string, any>>({});

// Use Hooks
const { industries, loading: industryLoading, error: industryError } = useIndustries({ autoLoad: true });

const { dialogId, focusFirstElement, onBackdropClick, close } = useModalFocusTrap(props, emit);

// Industry search state
const industrySearch = ref('');
const industryDropdownOpen = ref(false);
const industryDropdownRef = ref<HTMLElement | null>(null);

// Filtered industries based on search
const filteredIndustries = computed(() => {
  if (!industrySearch.value.trim()) {
    return industries.value;
  }
  const query = industrySearch.value.toLowerCase();
  return industries.value.filter(i => 
    i.name.toLowerCase().includes(query)
  );
});

// Select industry
function selectIndustry(id: number, name: string) {
  form.value.industry_id = id;
  industrySearch.value = name;
  industryDropdownOpen.value = false;
}

// Clear industry selection
function clearIndustry() {
  form.value.industry_id = null;
  industrySearch.value = '';
}

// Close dropdown when clicking outside
function onClickOutsideIndustry(e: MouseEvent) {
  if (!industryDropdownRef.value?.contains(e.target as Node)) {
    industryDropdownOpen.value = false;
  }
}

onMounted(() => {
  document.addEventListener('click', onClickOutsideIndustry);
});

onBeforeUnmount(() => {
  document.removeEventListener('click', onClickOutsideIndustry);
});

// Sync initialData -> form
watch(
  () => props.initialData,
  (v) => {
    form.value = {
      company_name: v?.company_name ?? '',
      website: v?.website ?? null,
      industry_id: v?.industry_id ?? null,
      status: v?.status ?? 'active',
      description: v?.description ?? null,
    };
    errors.value = {};
  },
  { immediate: true }
);

// Watch industries to set industrySearch after data loads
watch(
  () => industries.value,
  (list) => {
    if (list.length > 0 && props.initialData?.industry_id && !industrySearch.value) {
      const found = list.find(i => i.ID === props.initialData?.industry_id);
      if (found) {
        industrySearch.value = found.name;
      }
    }
  }
);

// When modal opens
watch(
  () => props.isOpen,
  (open) => {
    if (open) {
      nextTick(() => focusFirstElement());
      // Set industry search if data already loaded
      if (props.initialData?.industry_id && industries.value.length > 0) {
        const found = industries.value.find(i => i.ID === props.initialData?.industry_id);
        industrySearch.value = found?.name || '';
      }
    }
  }
);

// Helper: build partial diff (only include changed fields)
function buildPartialPayload(initial: SponsorResponse | null, current: Partial<SponsorPayload>) {
  const payload: Partial<SponsorPayload> = {};
  if (!initial) {
    Object.assign(payload, current);
    return payload;
  }

  if ((current.company_name ?? '') !== (initial.company_name ?? '')) {
    payload.company_name = current.company_name ?? '';
  }

  const curWebsite = current.website ?? null;
  const initWebsite = initial.website ?? null;
  if (String(curWebsite ?? '') !== String(initWebsite ?? '')) {
    payload.website = curWebsite;
  }

  if ((current.industry_id ?? null) !== (initial.industry_id ?? null)) {
    payload.industry_id = current.industry_id ?? null;
  }

  if ((current.status ?? '') !== (initial.status ?? '')) {
    payload.status = current.status ?? '';
  }

  const curDesc = current.description ?? null;
  const initDesc = initial.description ?? null;
  if (String(curDesc ?? '') !== String(initDesc ?? '')) {
    payload.description = curDesc;
  }

  return payload;
}

async function submit() {
  errors.value = {};
  const validateObj: SponsorPayload = {
    company_name: String(form.value.company_name ?? ''),
    website: form.value.website ?? null,
    industry_id: form.value.industry_id ?? null,
    status: String(form.value.status ?? 'active'),
    description: form.value.description ?? null,
    contacts: [],
  };

  const { valid, errors: vErrors } = validateSponsorForm(validateObj);
  const cleanedErrors = { ...(vErrors as Record<string, any>) };
  delete cleanedErrors.contacts;
  errors.value = cleanedErrors;

  if (!valid && Object.keys(cleanedErrors).length > 0) {
    await Swal.fire({ icon: 'warning', title: 'กรุณาตรวจสอบข้อมูล', text: 'มีข้อมูลที่ไม่ถูกต้อง' });
    await nextTick();
    focusFirstElement();
    return;
  }

  const partial = buildPartialPayload(props.initialData, form.value);

  if (Object.keys(partial).length === 0) {
    emit('update:isOpen', false);
    return;
  }

  saving.value = true;
  try {
    const updated = await SponsorService.update((props.initialData as SponsorResponse).ID, partial);
    emit('updated', updated);
    emit('update:isOpen', false);
    await Swal.fire({ icon: 'success', title: 'อัปเดตบริษัทเรียบร้อย' });
  } catch (err: any) {
    console.error('Sponsor update failed', err);
    const resp = err?.response?.data;
    if (resp?.errors) {
      errors.value = resp.errors;
      await Swal.fire({ icon: 'warning', title: 'ข้อมูลไม่ผ่านการตรวจสอบ', text: resp.message ?? 'โปรดตรวจสอบฟอร์ม' });
    } else {
      await Swal.fire({ icon: 'error', title: 'อัปเดตไม่สำเร็จ', text: resp?.message ?? err?.message ?? 'เกิดข้อผิดพลาด' });
    }
  } finally {
    saving.value = false;
  }
}
</script>

<template>
  <div>
    <teleport to="body">
    <div
      v-if="props.isOpen"
      class="fixed inset-0 z-[250] flex items-center justify-center bg-black/50 backdrop-blur-sm p-4"
      @click="onBackdropClick"
      data-theme="light"
    >
      <div
        ref="modalRef"
        role="dialog"
        :aria-modal="true"
        :aria-labelledby="`${dialogId}-title`"
        class="bg-white w-full max-w-3xl rounded-2xl shadow-2xl flex flex-col overflow-hidden animate-pop-in"
      >
        <!-- Header -->
        <div class="px-6 py-4 border-b flex items-center justify-between bg-slate-50">
          <div>
            <h2 :id="`${dialogId}-title`" class="text-xl font-bold text-[#1e3a8a]">แก้ไขข้อมูลบริษัท</h2>
          </div>

          <button class="btn btn-circle btn-ghost btn-sm" @click="close" aria-label="ปิด">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24"
              stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <!-- body -->
        <div class="p-4 overflow-y-auto">
          <div class="grid grid-cols-1 gap-4">
            <div>
              <label class="block text-sm text-gray-600 mb-1">ชื่อบริษัท <span class="text-red-500">*</span></label>
              <input v-model="form.company_name" class="input input-bordered w-full" />
              <p v-if="errors.company_name" class="text-xs text-red-500 mt-1">{{ errors.company_name }}</p>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label class="block text-sm text-gray-600 mb-1">เว็บไซต์</label>
                <input v-model="form.website" placeholder="https://example.com" class="input input-bordered w-full" />
                <p v-if="errors.website" class="text-xs text-red-500 mt-1">{{ errors.website }}</p>
              </div>

              <div>
                <label class="block text-sm text-gray-600 mb-1">สถานะ</label>
                <select v-model="form.status" class="select select-bordered w-full">
                  <option value="active">Active</option>
                  <option value="inactive">Inactive</option>
                </select>
                <p v-if="errors.status" class="text-xs text-red-500 mt-1">{{ errors.status }}</p>
              </div>
            </div>

            <div>
              <label class="block text-sm text-gray-600 mb-1">คำอธิบาย</label>
              <textarea v-model="form.description" rows="3" class="textarea textarea-bordered w-full"></textarea>
              <p v-if="errors.description" class="text-xs text-red-500 mt-1">{{ errors.description }}</p>
            </div>

            <!-- Industry (Searchable) -->
            <div ref="industryDropdownRef" class="relative">
              <label class="block text-sm text-gray-600 mb-1">อุตสาหกรรม</label>
              
              <div class="relative">
                <input
                  v-model="industrySearch"
                  type="text"
                  class="input input-bordered w-full pr-16"
                  placeholder="ค้นหาหรือเลือกอุตสาหกรรม..."
                  :disabled="industryLoading"
                  @focus="industryDropdownOpen = true"
                  @input="industryDropdownOpen = true; form.industry_id = null"
                />
                
                <!-- Clear button -->
                <button
                  v-if="industrySearch || form.industry_id"
                  type="button"
                  class="absolute right-10 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600"
                  @click.stop="clearIndustry"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
                
                <!-- Dropdown arrow -->
                <button
                  type="button"
                  class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400"
                  @click.stop="industryDropdownOpen = !industryDropdownOpen"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 transition-transform" :class="{ 'rotate-180': industryDropdownOpen }" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                  </svg>
                </button>
              </div>
              
              <!-- Dropdown list -->
              <div
                v-if="industryDropdownOpen && !industryLoading"
                class="absolute z-50 w-full mt-1 bg-white border border-gray-200 rounded-lg shadow-lg max-h-48 overflow-y-auto"
              >
                <div v-if="filteredIndustries.length === 0" class="px-4 py-3 text-sm text-gray-500 text-center">
                  ไม่พบอุตสาหกรรมที่ค้นหา
                </div>
                <button
                  v-for="i in filteredIndustries"
                  :key="i.ID"
                  type="button"
                  class="w-full px-4 py-2.5 text-left text-sm hover:bg-blue-50 transition-colors flex items-center justify-between"
                  :class="{ 'bg-blue-50 text-blue-700': form.industry_id === i.ID }"
                  @click="selectIndustry(i.ID, i.name)"
                >
                  <span>{{ i.name }}</span>
                  <svg v-if="form.industry_id === i.ID" xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-blue-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                  </svg>
                </button>
              </div>

              <p v-if="industryLoading" class="text-xs text-slate-400 mt-1">กำลังโหลดรายการ...</p>
              <p v-if="industryError" class="text-xs text-red-500 mt-1">{{ industryError }}</p>
            </div>
          </div>
        </div>

        <!-- footer -->
        <div class="px-4 py-3 border-t bg-slate-50 flex items-center justify-end gap-3">
          <button 
            class="btn btn-ghost hover:bg-gray-200 transition-all duration-200" 
            @click="close" 
            type="button"
          >
            ยกเลิก
          </button>
          <button 
            class="btn bg-[#1e3a8a] hover:bg-[#152c6f] text-white border-none
                   shadow-md hover:shadow-lg hover:-translate-y-0.5 
                   disabled:opacity-50 disabled:cursor-not-allowed disabled:transform-none
                   transition-all duration-200" 
            :disabled="saving" 
            @click="submit" 
            type="button"
          >
            <span v-if="saving" class="loading loading-spinner" aria-hidden="true"></span>
            <span v-else>บันทึกการเปลี่ยนแปลง</span>
          </button>
        </div>
      </div>
    </div>
    </teleport>
  </div>
</template>

<style scoped>
@keyframes pop-in { from { opacity: 0; transform: translateY(6px) scale(.995); } to { opacity: 1; transform: translateY(0) scale(1); } }
.animate-pop-in { animation: pop-in .12s cubic-bezier(.2,.9,.3,1); }
</style>
