<script setup lang="ts">
import { ref, watch, computed, nextTick, onMounted } from 'vue';
import type { PropType } from 'vue';
import type { SponsorResponse, SponsorPayload } from '@/interfaces/sponsor';
import type { IndustryResponse } from '@/interfaces/sponsor'; // <-- ปรับ path ตามจริง
import { SponsorService } from '@/services/sponsor/sponsor';
import { IndustryService } from '@/services/sponsor/industry'; // <-- ตรวจสอบ path ให้ตรง
import { validateSponsorForm } from '@/validators/sponsorValidator';
import Swal from 'sweetalert2';
import { useFocusTrap } from './useFocusTrap';

// Props & emits
const props = defineProps({
  isOpen: { type: Boolean as PropType<boolean>, default: false },
  loading: { type: Boolean as PropType<boolean>, default: false },
  initialData: { type: Object as PropType<SponsorResponse | null>, default: null },
  industries: { type: Array as PropType<IndustryResponse[] | null>, default: null }, // optional pre-loaded list
});

const emit = defineEmits<{
  (e: 'update:isOpen', v: boolean): void;
  (e: 'close'): void;
  (e: 'updated', sponsor: SponsorResponse): void;
}>();

// Local state
const form = ref<Partial<SponsorPayload>>({
  company_name: '',
  website: null,
  industry_id: null,
  status: 'active',
  description: null,
});
const saving = ref(false);
const errors = ref<Record<string, any>>({});

// industries state (if parent didn't pass, we'll load)
const industries = ref<IndustryResponse[]>([]);
const industryLoading = ref<boolean>(false);
const industryError = ref<string | null>(null);

// focus trap
const isOpenRef = computed(() => props.isOpen);
const { modalRef, dialogId, focusFirstElement, onBackdropClick } = useFocusTrap(isOpenRef, {
  onClose: () => {
    emit('update:isOpen', false);
    emit('close');
  },
});

// sync initialData -> form when opens or initialData changes
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

// loadIndustries function (uses IndustryService.getAll() which returns res.data)
async function loadIndustries() {
  // if parent already provided industries, don't load
  if (props.industries && props.industries.length > 0) {
    industries.value = props.industries;
    return;
  }

  // avoid duplicate loads
  if (industryLoading.value || industries.value.length > 0) return;

  industryLoading.value = true;
  industryError.value = null;

  try {
    const res = await IndustryService.getAll();
    // per your note, IndustryService.getAll() returns res.data directly
    industries.value = res ?? [];
  } catch (err: any) {
    console.error('โหลดอุตสาหกรรมผิดพลาด:', err);
    industryError.value = err?.message ?? 'โหลดอุตสาหกรรมไม่สำเร็จ';
    industries.value = [];
  } finally {
    industryLoading.value = false;
  }
}

// When modal opens, load industries if needed and focus
watch(
  () => props.isOpen,
  (open) => {
    if (open) {
      nextTick(() => focusFirstElement());
      // Load if parent didn't pass industries
      if (!props.industries) {
        loadIndustries();
      } else {
        industries.value = props.industries;
      }
    }
  }
);

// Also load on mount if desirable (optional)
onMounted(() => {
  // only load if parent didn't pass industries and nothing loaded yet
  if (!props.industries && industries.value.length === 0) {
    loadIndustries();
  } else if (props.industries) {
    industries.value = props.industries;
  }
});

// helper: build partial diff (only include changed fields)
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

function close() {
  emit('update:isOpen', false);
  emit('close');
}
</script>

<template>
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
        <!-- header -->
        <div class="px-4 py-3 border-b flex items-center justify-between bg-slate-50">
          <h3 :id="`${dialogId}-title`" class="text-lg font-semibold">แก้ไขข้อมูลบริษัท</h3>
          <div>
            <button class="btn btn-ghost btn-sm" @click="close">ปิด</button>
          </div>
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

            <!-- Industry dropdown -->
            <div>
              <label class="block text-sm text-gray-600 mb-1">อุตสาหกรรม</label>

              <select
                v-model.number="form.industry_id"
                class="select select-bordered w-full bg-white"
                :disabled="industryLoading"
              >
                <option :value="null">-- เลือกอุตสาหกรรม --</option>
                <option v-for="i in industries" :key="i.ID" :value="i.ID">{{ i.name }}</option>
              </select>

              <p v-if="industryLoading" class="text-xs text-slate-400 mt-1">กำลังโหลดรายการ...</p>
              <p v-if="industryError" class="text-xs text-red-500 mt-1">{{ industryError }}</p>
            </div>
          </div>
        </div>

        <!-- footer -->
        <div class="px-4 py-3 border-t bg-slate-50 flex items-center justify-end gap-2">
          <button class="btn btn-ghost" @click="close" type="button">ยกเลิก</button>
          <button class="btn btn-primary" :disabled="saving" @click="submit" type="button">
            <span v-if="saving" class="loading loading-spinner" aria-hidden="true"></span>
            <span v-else>บันทึกการเปลี่ยนแปลง</span>
          </button>
        </div>
      </div>
    </div>
  </teleport>
</template>

<style scoped>
@keyframes pop-in { from { opacity: 0; transform: translateY(6px) scale(.995); } to { opacity: 1; transform: translateY(0) scale(1); } }
.animate-pop-in { animation: pop-in .12s cubic-bezier(.2,.9,.3,1); }
</style>
