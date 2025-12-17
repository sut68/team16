<script setup lang="ts">
import { ref, watch } from 'vue';
import * as UserAPI from '@/services/api/user';
import Swal from 'sweetalert2';

const props = defineProps<{ initialData: any }>();
const emit = defineEmits(['refresh']);

const isEditing = ref(false);
const form = ref<any>({});

// เมื่อ props เปลี่ยน ให้ update form
watch(() => props.initialData, (newVal) => {
  if (newVal) form.value = JSON.parse(JSON.stringify(newVal));
}, { immediate: true });

// 🟢 ฟังก์ชันตรวจสอบข้อมูล (Validation)
const validateForm = (): string | null => {
  const f = form.value;

  if (!f.admin_firstname?.trim()) return 'กรุณาระบุชื่อจริง';
  if (!f.admin_lastname?.trim()) return 'กรุณาระบุนามสกุล';
  if (!f.position?.trim()) return 'กรุณาระบุตำแหน่ง';

  // ตรวจสอบ Email
  const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  if (!f.email) return 'กรุณาระบุอีเมล';
  if (!emailPattern.test(f.email)) return 'รูปแบบอีเมลไม่ถูกต้อง';

  // ตรวจสอบเบอร์โทร (10 หลัก)
  const phonePattern = /^\d{10}$/;
  if (!f.phone) return 'กรุณาระบุเบอร์โทรศัพท์';
  if (!phonePattern.test(f.phone)) return 'เบอร์โทรศัพท์ต้องเป็นตัวเลข 10 หลัก';

  return null; // ผ่าน
};

const save = async () => {
  // 1. เรียกใช้ Validation ก่อน
  const errorMsg = validateForm();
  if (errorMsg) {
    return Swal.fire({
      icon: 'warning',
      title: 'ข้อมูลไม่ถูกต้อง',
      text: errorMsg,
      confirmButtonText: 'ตกลง'
    });
  }

  // 2. ถ้าผ่าน ค่อยยิง API
  try {
    const res = await UserAPI.updateMyProfile(form.value);
    if (res && (res.status === 200 || res.message)) {
      Swal.fire({ 
        icon: 'success', 
        title: 'บันทึกสำเร็จ', 
        timer: 1500, 
        showConfirmButton: false 
      });
      isEditing.value = false;
      emit('refresh'); // แจ้ง parent ให้โหลดข้อมูลใหม่
    } else {
      throw new Error('บันทึกไม่สำเร็จ');
    }
  } catch (err: any) {
    Swal.fire('ผิดพลาด', err?.response?.data?.error || err.message, 'error');
  }
};
</script>

<template>
  <div class="bg-white p-6 rounded-2xl border border-gray-100 shadow-sm animate-fade-in">
    <div class="flex justify-between items-center mb-6 pb-4 border-b">
      <div>
        <h2 class="text-xl font-bold text-[#1e3a8a]">ข้อมูลผู้ดูแลระบบ</h2>
        <p class="text-xs text-gray-500">จัดการข้อมูลส่วนตัว</p>
      </div>
      <div>
        <button v-if="!isEditing" @click="isEditing = true" class="btn btn-sm btn-warning text-white rounded-full px-6 shadow-sm hover:shadow-md transition-all">
          <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="mr-1"><path d="M17 3a2.828 2.828 0 1 1 4 4L7.5 20.5 2 22l1.5-5.5L17 3z"></path></svg>
          แก้ไข
        </button>
        <div v-else class="flex gap-2">
          <button @click="isEditing = false; emit('refresh')" class="btn btn-sm btn-ghost rounded-full text-gray-500">ยกเลิก</button>
          <button @click="save" class="btn btn-sm bg-[#1e3a8a] text-white rounded-full border-none shadow-md hover:bg-blue-900 transition-all">บันทึก</button>
        </div>
      </div>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
      <div class="form-control">
        <label class="label text-sm text-gray-600 font-medium">ชื่อจริง <span v-if="isEditing" class="text-red-500 ml-1">*</span></label>
        <input v-model="form.admin_firstname" :disabled="!isEditing" class="input input-bordered w-full focus:border-[#1e3a8a] focus:ring-1 focus:ring-[#1e3a8a]" placeholder="ระบุชื่อจริง"/>
      </div>
      <div class="form-control">
        <label class="label text-sm text-gray-600 font-medium">นามสกุล <span v-if="isEditing" class="text-red-500 ml-1">*</span></label>
        <input v-model="form.admin_lastname" :disabled="!isEditing" class="input input-bordered w-full focus:border-[#1e3a8a] focus:ring-1 focus:ring-[#1e3a8a]" placeholder="ระบุนามสกุล"/>
      </div>
      <div class="form-control md:col-span-2">
        <label class="label text-sm text-gray-600 font-medium">ตำแหน่ง <span v-if="isEditing" class="text-red-500 ml-1">*</span></label>
        <input v-model="form.position" :disabled="!isEditing" class="input input-bordered w-full focus:border-[#1e3a8a] focus:ring-1 focus:ring-[#1e3a8a]" placeholder="เช่น เจ้าหน้าที่ทะเบียน"/>
      </div>
      <div class="form-control">
        <label class="label text-sm text-gray-600 font-medium">อีเมล <span v-if="isEditing" class="text-red-500 ml-1">*</span></label>
        <input v-model="form.email" :disabled="!isEditing" type="email" class="input input-bordered w-full focus:border-[#1e3a8a] focus:ring-1 focus:ring-[#1e3a8a]" placeholder="example@sut.ac.th"/>
      </div>
      <div class="form-control">
        <label class="label text-sm text-gray-600 font-medium">เบอร์โทรศัพท์ <span v-if="isEditing" class="text-red-500 ml-1">*</span></label>
        <input v-model="form.phone" :disabled="!isEditing" type="tel" maxlength="10" class="input input-bordered w-full focus:border-[#1e3a8a] focus:ring-1 focus:ring-[#1e3a8a]" placeholder="08xxxxxxxx"/>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* เพิ่ม Animation เล็กน้อย */
.animate-fade-in {
  animation: fadeIn 0.3s ease-in-out;
}
@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>