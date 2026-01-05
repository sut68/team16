<script setup lang="ts">
  import { ref, watch, nextTick } from 'vue';
  import type { PropType } from 'vue';
  import type { ContactPayload, ContactResponse } from '@/interfaces/sponsor';
  import { SponsorService } from '@/services/sponsor/sponsor';
  import Swal from 'sweetalert2';
  import { useModalFocusTrap } from '@/hooks/sponsor/useFocusTrap';
  import { validateContacts, buildContactsBatch } from '@/validators/sponsor_validator';

  // Props
  const props = defineProps({
    isOpen: { type: Boolean as PropType<boolean>, default: false },
    sponsorId: { type: Number as PropType<number>, required: true },
    sponsorName: { type: String as PropType<string>, default: '' },
    initialContacts: { type: Array as PropType<ContactResponse[]>, default: () => [] },
    disableBackdropClose: { type: Boolean as PropType<boolean>, default: false },
  });

  const emit = defineEmits<{
    (e: "update:isOpen", v: boolean): void;
    (e: "saved", contacts: ContactResponse[]): void;
    (e: "close"): void;
  }>();

  // State
  const localContacts = ref<ContactPayload[]>(
    (props.initialContacts ?? []).map(c => ({
      ID: c.ID,
      name: c.name,
      email: c.email,
      phone: c.phone,
      position: c.position ?? null,
    }))
  );
  const saving = ref(false);
  const errors = ref<Record<number, Record<string, string>>>({});

  // Use Modal Focus Trap Hook
  const { dialogId, focusFirstElement, onBackdropClick, close: handleClose } = useModalFocusTrap(
    props, 
    emit, 
    { disableBackdropClose: props.disableBackdropClose }
  );

  // Sync localContacts เมื่อ parent ส่ง initialContacts ใหม่
  watch(() => props.initialContacts, (v) => {
    localContacts.value = (v ?? []).map(c => ({
      ID: c.ID,
      name: c.name,
      email: c.email,
      phone: c.phone,
      position: c.position ?? null,
    }));
  });

  // ทำความสะอาดข้อมูล contact (trim whitespace)
  function normalizeContacts(list: ContactPayload[]) {
    return list.map(c => ({
      ...(c as any),
      name: (c.name ?? '').toString().trim(),
      email: (c.email ?? '').toString().trim(),
      phone: (c.phone ?? '').toString().trim(),
      position: (c.position ?? '') === null ? null : (c.position ?? '').toString().trim(),
      ID: (c as any).ID ?? undefined,
    }));
  }

  // เพิ่มผู้ติดต่อใหม่ (row ว่าง)
  function addContact() {
    if (saving.value) return;
    localContacts.value.push({ name: '', email: '', phone: '', position: '' });
    nextTick(() => focusFirstElement());
  }

  // ลบผู้ติดต่อ (ถ้ามี ID จะถาม confirm ก่อน)
  async function removeContact(idx: number) {
    if (saving.value) return;

    const c = localContacts.value[idx] as any;
    if (c?.ID != null) {
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

  // บันทึกผู้ติดต่อทั้งหมด (validate -> build batch -> API call)
  async function saveContacts() {
    // ล้าง error ก่อนหน้า
    errors.value = {};

    // ทำความสะอาดข้อมูล
    const contacts = normalizeContacts(localContacts.value ?? []);

    // ตรวจสอบความถูกต้อง
    const { valid, errors: vErrors } = validateContacts(contacts);
    if (!valid) {
      errors.value = vErrors ?? {};
      await Swal.fire({ icon: 'warning', title: 'โปรดตรวจสอบผู้ติดต่อ', text: 'มีข้อมูลบางรายการไม่ถูกต้อง' });
      nextTick(() => focusFirstElement(true));
      return;
    }

    // สร้าง batch สำหรับ API (upsert + delete)
    const batch = buildContactsBatch(props.initialContacts ?? [], contacts);
    if ((!batch.upsert || batch.upsert.length === 0) && (!batch.delete_ids || batch.delete_ids.length === 0)) {
      await Swal.fire({ icon: 'info', title: 'ไม่มีการเปลี่ยนแปลง', text: 'ไม่พบการเปลี่ยนแปลงผู้ติดต่อ' });
      return;
    }

    // เรียก API
    saving.value = true;
    try {
      const res = await SponsorService.updateContacts(props.sponsorId, batch as any);
      
      if (res?.contacts) {
        // แปลง response เป็น ContactResponse[]
        const contactsResp: ContactResponse[] = (res.contacts as any[]).map((c) => ({
          ID: Number(c.ID),
          name: String(c.name ?? ''),
          email: String(c.email ?? ''),
          phone: String(c.phone ?? ''),
          position: c.position ?? null,
        }));

        const validContacts = contactsResp.filter(c => Number.isFinite(c.ID));
        emit('saved', validContacts.length > 0 ? validContacts : res.contacts as unknown as ContactResponse[]);
        emit('update:isOpen', false);
        await Swal.fire({ icon: 'success', title: 'บันทึกผู้ติดต่อเรียบร้อย' });
      } else {
        // กรณี response ไม่มี contacts (fallback)
        emit('update:isOpen', false);
        await Swal.fire({ icon: 'success', title: 'บันทึกสำเร็จ' });
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
      
      // แสดง error จาก backend (ถ้ามี)
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

</script>

<template>
  <div>
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
        <!-- Header -->
        <div class="px-6 py-4 border-b flex items-center justify-between bg-slate-50">
          <div>
            <h2 :id="`${dialogId}-title`" class="text-xl font-bold text-[#1e3a8a]">จัดการผู้ติดต่อ</h2>
          </div>

          <button class="btn btn-circle btn-ghost btn-sm" @click="handleClose" aria-label="ปิด">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24"
              stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <div class="p-4 overflow-y-auto flex-1">
          <div class="mb-3 flex items-center justify-between">
            <div class="text-sm text-slate-600">Sponsor: <strong>{{ props.sponsorName }}</strong></div>
            <button 
              class="inline-flex items-center gap-2 px-4 py-2 text-sm font-medium
                    text-blue-600 bg-blue-50 border border-blue-200 rounded-lg 
                    hover:bg-blue-100 hover:border-blue-300
                    focus:ring-2 focus:ring-blue-500 focus:ring-offset-2
                    disabled:opacity-50 disabled:cursor-not-allowed
                    transition-all duration-200" 
              @click="addContact" 
              :disabled="saving"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z" />
              </svg>
              เพิ่มผู้ติดต่อ
            </button>
          </div>

          <div v-if="localContacts.length === 0" class="flex flex-col items-center gap-2 py-8 text-slate-500">
            <div class="text-sm">ยังไม่มีผู้ติดต่อ</div>
          </div>

          <div class="space-y-4">
            <template v-for="(c, idx) in localContacts" :key="(c as any).ID ?? idx">
              <div class="bg-gray-50 rounded-xl p-4 border border-gray-200 hover:border-gray-300 transition-colors">
                <!-- Card Header -->
                <div class="flex items-center justify-between mb-3 pb-2 border-b border-gray-200">
                  <span class="text-sm font-medium text-gray-700">
                    ผู้ติดต่อที่ {{ idx + 1 }}
                  </span>
                  <button
                    class="p-2 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-lg 
                           disabled:opacity-50 disabled:cursor-not-allowed
                           transition-all duration-200"
                    @click.prevent="removeContact(idx)"
                    :disabled="saving"
                    aria-label="ลบผู้ติดต่อ"
                    title="ลบผู้ติดต่อ"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                    </svg>
                  </button>
                </div>

                <!-- Form Fields -->
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <!-- ชื่อ -->
                  <div class="space-y-1">
                    <label class="text-sm font-medium text-gray-600">
                      ชื่อ <span class="text-red-500">*</span>
                    </label>
                    <input
                      v-model="c.name"
                      class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all"
                      :class="{ 'border-red-500 bg-red-50': errors[idx]?.name }"
                      placeholder="ชื่อผู้ติดต่อ"
                      :aria-invalid="errors[idx]?.name ? 'true' : 'false'"
                    />
                    <p v-if="errors[idx]?.name" class="text-xs text-red-500">{{ errors[idx].name }}</p>
                  </div>

                  <!-- อีเมล -->
                  <div class="space-y-1">
                    <label class="text-sm font-medium text-gray-600">
                      อีเมล <span class="text-red-500">*</span>
                    </label>
                    <input
                      v-model="c.email"
                      type="email"
                      class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all"
                      :class="{ 'border-red-500 bg-red-50': errors[idx]?.email }"
                      placeholder="example@email.com"
                      :aria-invalid="errors[idx]?.email ? 'true' : 'false'"
                    />
                    <p v-if="errors[idx]?.email" class="text-xs text-red-500">{{ errors[idx].email }}</p>
                  </div>

                  <!-- เบอร์โทร -->
                  <div class="space-y-1">
                    <label class="text-sm font-medium text-gray-600">
                      เบอร์โทรศัพท์ <span class="text-red-500">*</span>
                    </label>
                    <input
                      v-model="c.phone"
                      type="tel"
                      class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all"
                      :class="{ 'border-red-500 bg-red-50': errors[idx]?.phone }"
                      placeholder="0xx-xxx-xxxx"
                      :aria-invalid="errors[idx]?.phone ? 'true' : 'false'"
                    />
                    <p v-if="errors[idx]?.phone" class="text-xs text-red-500">{{ errors[idx].phone }}</p>
                  </div>

                  <!-- ตำแหน่ง -->
                  <div class="space-y-1">
                    <label class="text-sm font-medium text-gray-600">
                      ตำแหน่ง
                    </label>
                    <input
                      v-model="c.position"
                      class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all"
                      :class="{ 'border-red-500 bg-red-50': errors[idx]?.position }"
                      placeholder="เช่น ผู้จัดการ, HR"
                      :aria-invalid="errors[idx]?.position ? 'true' : 'false'"
                    />
                    <p v-if="errors[idx]?.position" class="text-xs text-red-500">{{ errors[idx].position }}</p>
                  </div>
                </div>
              </div>
            </template>
          </div>
        </div>

        <div class="px-4 py-3 border-t bg-slate-50 flex items-center justify-end gap-3">
          <button 
            class="btn btn-ghost hover:bg-gray-200 transition-all duration-200" 
            @click="handleClose" 
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
            @click="saveContacts" 
            type="button"
          >
            <span v-if="saving" class="loading loading-spinner" aria-hidden="true"></span>
            <span v-else>บันทึกผู้ติดต่อ</span>
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
