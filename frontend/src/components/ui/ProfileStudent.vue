<script setup lang="ts">
import { ref } from 'vue';
import * as UserAPI from '@/services/api/user';
import Swal from 'sweetalert2';

const props = defineProps<{ initialData: any }>();
const emit = defineEmits(['refresh']);

const isEditing = ref(false);
const activeTab = ref('personal');
const form = ref({ ...props.initialData });

// Init family info ถ้าไม่มี
if (!form.value.family_info) {
  form.value.family_info = {
    father_name: '', father_occupation: '', father_income: 0,
    mother_name: '', mother_occupation: '', mother_income: 0,
    guardian_name: '', guardian_relation: '', guardian_occupation: '', guardian_income: 0
  };
}

const save = async () => {
  // Format Date
  if (form.value.birth_date) form.value.birth_date = new Date(form.value.birth_date).toISOString();
  
  const payload = { ...form.value, family_info: form.value.family_info };
  const res = await UserAPI.updateMyProfile(payload);
  
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
      <h2 class="text-2xl font-bold text-gray-800">ข้อมูลส่วนตัวนิสิต</h2>
      <div>
        <button v-if="!isEditing" @click="isEditing = true" class="btn btn-warning text-white">แก้ไขข้อมูล</button>
        <div v-else class="space-x-2">
          <button @click="isEditing = false; emit('refresh')" class="btn btn-ghost">ยกเลิก</button>
          <button @click="save" class="btn bg-[#1e3a8a] text-white">บันทึก</button>
        </div>
      </div>
    </div>

    <div class="tabs tabs-boxed mb-6 bg-gray-100 p-1 rounded-lg inline-flex">
      <a class="tab px-6" :class="{ 'bg-white font-bold shadow-sm text-primary': activeTab === 'personal' }" @click="activeTab = 'personal'">ข้อมูลส่วนตัว</a>
      <a class="tab px-6" :class="{ 'bg-white font-bold shadow-sm text-primary': activeTab === 'family' }" @click="activeTab = 'family'">ครอบครัว</a>
    </div>

    <div v-if="activeTab === 'personal'" class="grid grid-cols-1 md:grid-cols-2 gap-6">
      <div class="alert alert-info col-span-1 md:col-span-2 text-sm text-white">
        ข้อมูลวิชาการและชื่อภาษาไทย ไม่สามารถแก้ไขได้ หากต้องการแก้ไขโปรดติดต่อเจ้าหน้าที่
      </div>

      <div class="form-control"><label class="label text-gray-500">รหัสนักศึกษา</label><input :value="form.student_id" disabled class="input input-bordered bg-gray-100"/></div>
      <div class="form-control"><label class="label text-gray-500">สาขาวิชา</label><input :value="form.major?.major_name" disabled class="input input-bordered bg-gray-100"/></div>
      <div class="form-control"><label class="label text-gray-500">ชื่อ-นามสกุล (ไทย)</label><input :value="form.first_name_th + ' ' + form.last_name_th" disabled class="input input-bordered bg-gray-100"/></div>
      <div class="form-control"><label class="label text-gray-500">เลขบัตรประชาชน</label><input :value="form.national_id" disabled class="input input-bordered bg-gray-100"/></div>
      <div class="form-control"><label class="label text-gray-500">เกรดเฉลี่ย (GPAX)</label><input :value="form.gpax" disabled class="input input-bordered bg-gray-100"/></div>
      <div class="form-control"><label class="label text-gray-500">อาจารย์ที่ปรึกษา</label><input :value="form.advisor_name" disabled class="input input-bordered bg-gray-100"/></div>

      <div class="divider col-span-1 md:col-span-2 text-gray-400">แก้ไขข้อมูลติดต่อได้</div>

      <div class="form-control"><label class="label font-bold">ชื่อ (EN)</label><input v-model="form.first_name_en" :disabled="!isEditing" class="input input-bordered"/></div>
      <div class="form-control"><label class="label font-bold">นามสกุล (EN)</label><input v-model="form.last_name_en" :disabled="!isEditing" class="input input-bordered"/></div>
      <div class="form-control"><label class="label font-bold">วันเกิด</label><input v-model="form.birth_date" type="date" :disabled="!isEditing" class="input input-bordered"/></div>
      <div class="form-control"><label class="label font-bold">อีเมล</label><input v-model="form.email" :disabled="!isEditing" class="input input-bordered"/></div>
      <div class="form-control"><label class="label font-bold">เบอร์โทรศัพท์</label><input v-model="form.phone" :disabled="!isEditing" class="input input-bordered"/></div>
      <div class="form-control md:col-span-2"><label class="label font-bold">ที่อยู่</label><textarea v-model="form.permanent_address" :disabled="!isEditing" class="textarea textarea-bordered"></textarea></div>
    </div>

    <div v-if="activeTab === 'family'" class="space-y-6">
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <h3 class="md:col-span-3 font-bold border-b pb-2">ข้อมูลบิดา</h3>
        <div class="form-control"><label class="label">ชื่อ-นามสกุล</label><input v-model="form.family_info.father_name" :disabled="!isEditing" class="input input-bordered"/></div>
        <div class="form-control"><label class="label">อาชีพ</label><input v-model="form.family_info.father_occupation" :disabled="!isEditing" class="input input-bordered"/></div>
        <div class="form-control"><label class="label">รายได้</label><input v-model.number="form.family_info.father_income" type="number" :disabled="!isEditing" class="input input-bordered"/></div>
      </div>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <h3 class="md:col-span-3 font-bold border-b pb-2">ข้อมูลมารดา</h3>
        <div class="form-control"><label class="label">ชื่อ-นามสกุล</label><input v-model="form.family_info.mother_name" :disabled="!isEditing" class="input input-bordered"/></div>
        <div class="form-control"><label class="label">อาชีพ</label><input v-model="form.family_info.mother_occupation" :disabled="!isEditing" class="input input-bordered"/></div>
        <div class="form-control"><label class="label">รายได้</label><input v-model.number="form.family_info.mother_income" type="number" :disabled="!isEditing" class="input input-bordered"/></div>
      </div>
    </div>
  </div>
</template>