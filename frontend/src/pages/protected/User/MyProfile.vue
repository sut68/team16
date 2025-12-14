<script setup lang="ts">
import { ref, onMounted } from 'vue';
import * as UserAPI from '@/services/api/user';
import Swal from 'sweetalert2';

const role = ref('');
const isEditing = ref(false);
const activeTab = ref('personal');

const profile = ref<any>({
  student_id: '', national_id: '', first_name_th: '', last_name_th: '', 
  first_name_en: '', last_name_en: '', major: { major_name: '' }, 
  gpax: 0, advisor_name: '', birth_date: '',
  permanent_address: '', current_address: '', province: '', siblings_count: 0,
  admin_firstname: '', admin_lastname: '', position: '',
  email: '', phone: '',
  family_info: {
    father_name: '', father_occupation: '', father_income: 0,
    mother_name: '', mother_occupation: '', mother_income: 0,
    guardian_name: '', guardian_relation: '', guardian_occupation: '', guardian_income: 0
  }
});

const loadProfile = async () => {
  const res = await UserAPI.getMyProfile();
  if (res) {
    role.value = res.role;
    profile.value = { ...profile.value, ...res.data };
    if (res.family) profile.value.family_info = res.family;
    if (profile.value.birth_date) profile.value.birth_date = profile.value.birth_date.split('T')[0];
  }
};

onMounted(loadProfile);

const saveProfile = async () => {
  if (profile.value.birth_date) profile.value.birth_date = new Date(profile.value.birth_date).toISOString();
  
  const payload = { ...profile.value, family_info: profile.value.family_info };
  const res = await UserAPI.updateMyProfile(payload);
  if (res && res.message) {
    Swal.fire('สำเร็จ', 'บันทึกข้อมูลเรียบร้อย', 'success');
    isEditing.value = false;
    loadProfile();
  } else {
    Swal.fire('ผิดพลาด', 'บันทึกไม่สำเร็จ', 'error');
  }
};
</script>

<template>
  <div class="h-screen overflow-y-auto bg-gray-50 pb-20">
    <div class="max-w-5xl mx-auto p-6">
      
      <div class="bg-white shadow rounded-lg p-6 mb-6">
        <div class="flex justify-between items-center border-b pb-4 mb-6">
          <h2 class="text-2xl font-bold text-gray-800">
            {{ role === 'student' ? 'ข้อมูลส่วนตัวนิสิต' : 'ข้อมูลผู้ดูแลระบบ' }}
          </h2>
          <div>
            <button v-if="!isEditing" @click="isEditing = true" class="btn btn-warning text-white">แก้ไขข้อมูล</button>
            <div v-else>
              <button @click="isEditing = false; loadProfile()" class="btn btn-ghost mr-2">ยกเลิก</button>
              <button @click="saveProfile" class="btn btn-success text-white">บันทึก</button>
            </div>
          </div>
        </div>

        <div v-if="role === 'admin'" class="grid grid-cols-2 gap-6">
          <div class="form-control"><label class="label font-bold">ชื่อจริง</label><input v-model="profile.admin_firstname" :disabled="!isEditing" class="input input-bordered"/></div>
          <div class="form-control"><label class="label font-bold">นามสกุล</label><input v-model="profile.admin_lastname" :disabled="!isEditing" class="input input-bordered"/></div>
          <div class="form-control"><label class="label font-bold">ตำแหน่ง</label><input v-model="profile.position" :disabled="!isEditing" class="input input-bordered"/></div>
          <div class="form-control"><label class="label font-bold">อีเมล</label><input v-model="profile.email" :disabled="!isEditing" class="input input-bordered"/></div>
          <div class="form-control"><label class="label font-bold">เบอร์โทรศัพท์</label><input v-model="profile.phone" :disabled="!isEditing" class="input input-bordered"/></div>
        </div>

        <div v-if="role === 'student'">
          <div class="tabs tabs-boxed mb-6 bg-gray-100 p-1 inline-flex rounded-lg">
            <a class="tab px-6" :class="{ 'bg-white font-bold shadow-sm text-primary': activeTab === 'personal' }" @click="activeTab = 'personal'">ข้อมูลส่วนตัว</a>
            <a class="tab px-6" :class="{ 'bg-white font-bold shadow-sm text-primary': activeTab === 'family' }" @click="activeTab = 'family'">ครอบครัว</a>
          </div>

          <div v-if="activeTab === 'personal'">
            <div class="alert alert-info mb-6 text-sm text-white shadow-sm">ข้อมูลในช่องสีเทา ไม่สามารถแก้ไขได้ หากต้องการแก้ไขโปรดติดต่อเจ้าหน้าที่</div>
            
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div class="form-control"><label class="label text-gray-500">รหัสนักศึกษา</label><input :value="profile.student_id" disabled class="input input-bordered bg-gray-100"/></div>
              <div class="form-control"><label class="label text-gray-500">เลขบัตรประชาชน</label><input :value="profile.national_id" disabled class="input input-bordered bg-gray-100"/></div>
              <div class="form-control"><label class="label text-gray-500">ชื่อ-นามสกุล (ไทย)</label><input :value="profile.first_name_th + ' ' + profile.last_name_th" disabled class="input input-bordered bg-gray-100"/></div>
              <div class="form-control"><label class="label text-gray-500">สาขาวิชา</label><input :value="profile.major?.major_name" disabled class="input input-bordered bg-gray-100"/></div>
              <div class="form-control"><label class="label text-gray-500">เกรดเฉลี่ย (GPAX)</label><input :value="profile.gpax" disabled class="input input-bordered bg-gray-100"/></div>
              <div class="form-control"><label class="label text-gray-500">อาจารย์ที่ปรึกษา</label><input :value="profile.advisor_name" disabled class="input input-bordered bg-gray-100"/></div>

              <div class="divider col-span-1 md:col-span-2 text-gray-400">แก้ไขข้อมูลได้</div>
              
              <div class="form-control"><label class="label font-bold">ชื่อ (อังกฤษ)</label><input v-model="profile.first_name_en" :disabled="!isEditing" class="input input-bordered focus:border-primary"/></div>
              <div class="form-control"><label class="label font-bold">นามสกุล (อังกฤษ)</label><input v-model="profile.last_name_en" :disabled="!isEditing" class="input input-bordered focus:border-primary"/></div>
              <div class="form-control"><label class="label font-bold">วันเกิด</label><input v-model="profile.birth_date" type="date" :disabled="!isEditing" class="input input-bordered focus:border-primary"/></div>
              <div class="form-control"><label class="label font-bold">อีเมล</label><input v-model="profile.email" :disabled="!isEditing" class="input input-bordered focus:border-primary"/></div>
              <div class="form-control"><label class="label font-bold">เบอร์โทรศัพท์</label><input v-model="profile.phone" :disabled="!isEditing" class="input input-bordered focus:border-primary"/></div>
              <div class="form-control"><label class="label font-bold">จังหวัด</label><input v-model="profile.province" :disabled="!isEditing" class="input input-bordered focus:border-primary"/></div>
              <div class="form-control"><label class="label font-bold">จำนวนพี่น้อง</label><input v-model.number="profile.siblings_count" type="number" :disabled="!isEditing" class="input input-bordered focus:border-primary"/></div>
              <div class="form-control md:col-span-2"><label class="label font-bold">ที่อยู่ตามทะเบียนบ้าน</label><textarea v-model="profile.permanent_address" :disabled="!isEditing" class="textarea textarea-bordered focus:border-primary"></textarea></div>
              <div class="form-control md:col-span-2"><label class="label font-bold">ที่อยู่ปัจจุบัน</label><textarea v-model="profile.current_address" :disabled="!isEditing" class="textarea textarea-bordered focus:border-primary"></textarea></div>
            </div>
          </div>

          <div v-if="activeTab === 'family'" class="space-y-8 animate-fade-in">
            <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
                <h3 class="md:col-span-3 font-bold text-lg text-primary border-b pb-2">ข้อมูลบิดา</h3>
                <div class="form-control"><label class="label font-bold">ชื่อ-นามสกุล</label><input v-model="profile.family_info.father_name" :disabled="!isEditing" class="input input-bordered"/></div>
                <div class="form-control"><label class="label font-bold">อาชีพ</label><input v-model="profile.family_info.father_occupation" :disabled="!isEditing" class="input input-bordered"/></div>
                <div class="form-control"><label class="label font-bold">รายได้</label><input v-model.number="profile.family_info.father_income" type="number" :disabled="!isEditing" class="input input-bordered"/></div>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
                <h3 class="md:col-span-3 font-bold text-lg text-primary border-b pb-2">ข้อมูลมารดา</h3>
                <div class="form-control"><label class="label font-bold">ชื่อ-นามสกุล</label><input v-model="profile.family_info.mother_name" :disabled="!isEditing" class="input input-bordered"/></div>
                <div class="form-control"><label class="label font-bold">อาชีพ</label><input v-model="profile.family_info.mother_occupation" :disabled="!isEditing" class="input input-bordered"/></div>
                <div class="form-control"><label class="label font-bold">รายได้</label><input v-model.number="profile.family_info.mother_income" type="number" :disabled="!isEditing" class="input input-bordered"/></div>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <h3 class="md:col-span-2 font-bold text-lg text-primary border-b pb-2">ผู้ปกครอง</h3>
                <div class="form-control"><label class="label font-bold">ชื่อ-นามสกุล</label><input v-model="profile.family_info.guardian_name" :disabled="!isEditing" class="input input-bordered"/></div>
                <div class="form-control"><label class="label font-bold">ความเกี่ยวข้อง</label><input v-model="profile.family_info.guardian_relation" :disabled="!isEditing" class="input input-bordered"/></div>
                <div class="form-control"><label class="label font-bold">อาชีพ</label><input v-model="profile.family_info.guardian_occupation" :disabled="!isEditing" class="input input-bordered"/></div>
                <div class="form-control"><label class="label font-bold">รายได้</label><input v-model.number="profile.family_info.guardian_income" type="number" :disabled="!isEditing" class="input input-bordered"/></div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.animate-fade-in { animation: fadeIn 0.4s ease-out; }
@keyframes fadeIn { from { opacity: 0; transform: translateY(10px); } to { opacity: 1; transform: translateY(0); } }
</style>