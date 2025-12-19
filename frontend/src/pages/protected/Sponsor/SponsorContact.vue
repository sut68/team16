<script setup lang="ts">
import { ref, watch, computed, nextTick } from 'vue';
import type { PropType } from 'vue';
import type { ContactPayload, ContactResponse } from '../../../interfaces/sponsor';
import { SponsorService } from '../../../services/sponsor/sponsor';
import Swal from 'sweetalert2';
import { useFocusTrap } from '../../../hooks/sponsor/useFocusTrap'; 
import { validateContacts, buildContactsBatch } from '@/validators/sponsorValidator';

// Props + Emits
const props = defineProps({
  isOpen: { type: Boolean as PropType<boolean>, default: false },
  sponsorId: { type: Number as PropType<number>, required: true },
  initialContacts: { type: Array as PropType<ContactResponse[]>, default: () => [] },
  disableBackdropClose: { type: Boolean as PropType<boolean>, default: false },
});

const emit = defineEmits<{
  (e: "update:isOpen", v: boolean): void;
  (e: "saved", contacts: ContactResponse[]): void;
  (e: "close"): void;
}>();

// local state
const localContacts = ref<ContactPayload[]>(
  (props.initialContacts ?? []).map(c => ({
    ID: c.ID,
    name: c.name,
    email: c.email,
    phone: c.phone,
    position: c.position ?? null,
  }))
);

// removed unused `loading`
const saving = ref(false);
// make errors non-nullable for simpler usage
const errors = ref<Record<number, Record<string, string>>>({});

// keep localContacts in sync when initialContacts changes (e.g. parent refetch)
watch(() => props.initialContacts, (v) => {
  localContacts.value = (v ?? []).map(c => ({
    ID: c.ID,
    name: c.name,
    email: c.email,
    phone: c.phone,
    position: c.position ?? null,
  }));
});

// useFocusTrap
const isOpenRef = computed(() => props.isOpen);
const { dialogId, focusFirstElement, onBackdropClick } = useFocusTrap(isOpenRef, {
  onClose: () => {
    emit('update:isOpen', false);
    emit('close');
  },
  disableBackdropClose: props.disableBackdropClose,
});

// ----- helpers -----
function addContact() {
  if (saving.value) return;
  localContacts.value.push({ name: '', email: '', phone: '', position: '' });
  // focus next tick to the first focusable (new row)
  nextTick(() => focusFirstElement());
}

async function removeContact(idx: number) {
  if (saving.value) return;

  const c = localContacts.value[idx] as any;
  if (c?.ID != null) {
    // use Swal for consistent UI
    const answer = await Swal.fire({
      title: 'ลบผู้ติดต่อ',
      text: 'คุณต้องการลบผู้ติดต่อรายการนี้จริงหรือไม่?',
      icon: 'warning',
      showCancelButton: true,
      confirmButtonText: 'ลบ',
      cancelButtonText: 'ยกเลิก',
    });
    if (!answer.isConfirmed) return;
  }
  localContacts.value.splice(idx, 1);
}

// helper: normalize/trim contacts (so validation receives clean data)
function normalizeContacts(list: ContactPayload[]) {
  return list.map(c => ({
    ...(c as any),
    name: (c.name ?? '').toString().trim(),
    email: (c.email ?? '').toString().trim(),
    phone: (c.phone ?? '').toString().trim(),
    position: (c.position ?? '') === null ? null : (c.position ?? '').toString().trim(),
    // preserve ID if exists
    ID: (c as any).ID ?? undefined,
  }));
}

// submit contacts only
async function saveContacts() {
  // clear previous errors
  errors.value = {};

  // normalize first
  const contacts = normalizeContacts(localContacts.value ?? []);

  // local validation
  const { valid, errors: vErrors } = validateContacts(contacts);
  if (!valid) {
    errors.value = vErrors ?? {};
    await Swal.fire({ icon: 'warning', title: 'โปรดตรวจสอบผู้ติดต่อ', text: 'มีข้อมูลบางรายการไม่ถูกต้อง' });
    // focus first error
    nextTick(() => focusFirstElement(true));
    return;
  }

  const batch = buildContactsBatch(props.initialContacts ?? [], contacts);
  // buildContactsBatch should return shape { upsert?:..., delete_ids?:... }
  if ((!batch.upsert || batch.upsert.length === 0) && (!batch.delete_ids || batch.delete_ids.length === 0)) {
    await Swal.fire({ icon: 'info', title: 'ไม่มีการเปลี่ยนแปลง', text: 'ไม่พบการเปลี่ยนแปลงผู้ติดต่อ' });
    return;
  }

  saving.value = true;
  try {
    const res = await SponsorService.updateContacts(props.sponsorId, batch as any);
    // res expected { contacts: [...] }
    if (res?.contacts) {
      const contactsResp: ContactResponse[] = (res.contacts as any[]).map((c) => ({
        ID: Number(c.ID), // coerce to number
        name: String(c.name ?? ''),
        email: String(c.email ?? ''),
        phone: String(c.phone ?? ''),
        position: c.position ?? null,
      }))

      const validContacts = contactsResp.filter(c => Number.isFinite(c.ID));
      // emit saved with fresh contacts
      if (validContacts.length === 0) {
        // fallback: emit the raw contacts as best-effort (cast)
        emit('saved', res.contacts as unknown as ContactResponse[]);
      } else {
        emit('saved', validContacts);
      }
      
      emit('update:isOpen', false);
      await Swal.fire({ icon: 'success', title: 'บันทึกผู้ติดต่อเรียบร้อย' });
    } else {
      // fallback: success without contacts list
      emit('update:isOpen', false);
      await Swal.fire({ icon: 'success', title: 'บันทึกสำเร็จ' });
      // convert payload back to ContactResponse-ish for emit (best-effort)
      const fallback = contacts.map(c => ({
        ID: (c as any).ID ?? undefined,
        name: c.name,
        email: c.email,
        phone: c.phone,
        position: c.position ?? null,
      })) as ContactResponse[];
      emit('saved', fallback);
    }
  } catch (err: any) {
    console.error('updateContacts failed', err);
    const resp = err?.response?.data;
    // backend might return errors mapping for contacts: { "0": { name: "..." } }
    if (resp?.errors && resp.errors.contacts) {
      errors.value = resp.errors.contacts;
      await Swal.fire({ icon: 'warning', title: 'ข้อมูลไม่ผ่านการตรวจสอบ', text: resp.message ?? 'โปรดตรวจสอบข้อผิดพลาด' });
      nextTick(() => focusFirstElement(true));
    } else {
      await Swal.fire({ icon: 'error', title: 'เกิดข้อผิดพลาด', text: resp?.message ?? err?.message ?? 'ไม่สามารถบันทึกได้' });
    }
  } finally {
    saving.value = false;
  }
}

// watch open -> focus
watch(isOpenRef, async (open) => {
  if (open) {
    await nextTick();
    focusFirstElement();
  }
});

function handleClose() {
  emit('update:isOpen', false);
  emit('close');
}
</script>

<template>
  <teleport to="body">
    <div
      v-if="props.isOpen"
      class="fixed inset-0 z-[300] flex items-center justify-center bg-black/50 backdrop-blur-sm p-4"
      @click="onBackdropClick"
      data-theme="light"
    >
      <div
        ref="modalRef"
        role="dialog"
        :aria-modal="true"
        :aria-labelledby="`${dialogId}-title`"
        class="bg-white w-full max-w-3xl max-h-[85vh] rounded-2xl shadow-2xl flex flex-col overflow-hidden animate-pop-in"
      >
        <div class="px-4 py-3 border-b flex items-center justify-between bg-slate-50">
          <h3 :id="`${dialogId}-title`" class="text-lg font-semibold">จัดการผู้ติดต่อ</h3>
          <div class="flex items-center gap-2">
            <button class="btn btn-ghost btn-sm" @click="handleClose" aria-label="ปิด">ปิด</button>
          </div>
        </div>

        <div class="p-4 overflow-y-auto flex-1">
          <div class="mb-3 flex items-center justify-between">
            <div class="text-sm text-slate-600">Sponsor ID: <strong>{{ props.sponsorId }}</strong></div>
            <button class="btn btn-sm btn-outline" @click="addContact" :disabled="saving">เพิ่มผู้ติดต่อ</button>
          </div>

          <div v-if="localContacts.length === 0" class="flex flex-col items-center gap-2 py-8 text-slate-500">
            <div class="text-sm">ยังไม่มีผู้ติดต่อ</div>
          </div>

          <div class="space-y-3">
            <template v-for="(c, idx) in localContacts" :key="(c as any).ID ?? idx">
              <div class="grid grid-cols-12 gap-2 items-start">
                <input
                  v-model="c.name"
                  class="input input-sm input-bordered col-span-12 md:col-span-4"
                  placeholder="ชื่อ"
                  :aria-invalid="errors[idx] && errors[idx].name ? 'true' : 'false'"
                />
                <input
                  v-model="c.email"
                  class="input input-sm input-bordered col-span-12 md:col-span-3"
                  placeholder="email"
                  :aria-invalid="errors[idx] && errors[idx].email ? 'true' : 'false'"
                />
                <input
                  v-model="c.phone"
                  class="input input-sm input-bordered col-span-12 md:col-span-3"
                  placeholder="โทร."
                  :aria-invalid="errors[idx] && errors[idx].phone ? 'true' : 'false'"
                />
                <input
                  v-model="c.position"
                  class="input input-sm input-bordered col-span-9 md:col-span-1"
                  placeholder="ตำแหน่ง"
                  :aria-invalid="errors[idx] && errors[idx].position ? 'true' : 'false'"
                />
                <button
                  class="btn btn-sm btn-error col-span-3 md:col-span-1"
                  @click.prevent="removeContact(idx)"
                  :disabled="saving"
                  aria-label="ลบผู้ติดต่อ"
                >
                  ลบ
                </button>

                <p v-if="errors[idx]" class="text-xs text-red-500 col-span-12">
                  {{ Object.values(errors[idx])[0] }}
                </p>
              </div>
            </template>
          </div>
        </div>

        <div class="px-4 py-3 border-t bg-slate-50 flex items-center justify-end gap-2">
          <button class="btn btn-ghost" @click="handleClose" type="button">ยกเลิก</button>
          <button class="btn btn-primary" :disabled="saving" @click="saveContacts" type="button">
            <span v-if="saving" class="loading loading-spinner" aria-hidden="true"></span>
            <span v-else>บันทึกผู้ติดต่อ</span>
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
