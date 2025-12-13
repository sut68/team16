<script setup lang="ts">
import { ref } from 'vue';
import * as UserAPI from '@/services/api/user';
import Swal from 'sweetalert2';

const props = defineProps<{ initialData: any }>();
const emit = defineEmits(['refresh']);

const isEditing = ref(false);
const form = ref({ ...props.initialData });

const save = async () => {
  const res = await UserAPI.updateMyProfile(form.value);
  if (res && res.status === 200) {
    Swal.fire('สำเร็จ', 'บันทึกข้อมูลเรียบร้อย', 'success');
    isEditing.value = false;
    emit('refresh');
  } else {
    Swal.fire('ผิดพลาด', res?.data?.error || 'บันทึกไม่สำเร็จ', 'error');
  }
};
</script>

<template>
  <div class="bg-white rounded-xl shadow-sm border border-gray-100 p-6 animate-fade-in">
    <div class="flex justify-between items-center mb-6 pb-4 border-b">
      <h2 class="text-2xl font-bold text-gray-800">ข้อมูลผู้ดูแลระบบ (Admin)</h2>
      <div>
        <button v-if="!isEditing" @click="isEditing = true" class="btn btn-warning text-white">แก้ไขข้อมูล</button>
        <div v-else class="space-x-2">
          <button @click="isEditing = false; emit('refresh')" class="btn btn-ghost">ยกเลิก</button>
          <button @click="save" class="btn bg-[#1e3a8a] text-white">บันทึก</button>
        </div>
      </div>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
      <div class="form-control">
        <label class="label font-bold">ชื่อจริง</label>
        <input v-model="form.admin_first_name" :disabled="!isEditing" class="input input-bordered focus:border-primary"/>
      </div>
      <div class="form-control">
        <label class="label font-bold">นามสกุล</label>
        <input v-model="form.admin_last_name" :disabled="!isEditing" class="input input-bordered focus:border-primary"/>
      </div>
      <div class="form-control">
        <label class="label font-bold">ตำแหน่ง</label>
        <input v-model="form.position" :disabled="!isEditing" class="input input-bordered focus:border-primary"/>
      </div>
      <div class="form-control">
        <label class="label font-bold">อีเมล</label>
        <input v-model="form.email" :disabled="!isEditing" class="input input-bordered focus:border-primary"/>
      </div>
      <div class="form-control">
        <label class="label font-bold">เบอร์โทรศัพท์</label>
        <input v-model="form.phone" :disabled="!isEditing" class="input input-bordered focus:border-primary"/>
      </div>
    </div>
  </div>
</template>