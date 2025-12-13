<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import * as UserAPI from '@/services/api/user';
import Swal from 'sweetalert2';

const users = ref<any[]>([]);
const roles = ref<any[]>([]);
const majors = ref<any[]>([]);
const showModal = ref(false);
const isEditing = ref(false);
const editingUserId = ref<number | null>(null);

const initialForm = {
  username: '', password: '', role_id: 0,
  // Student
  student_id: '', national_id: '', first_name_th: '', last_name_th: '', 
  major_id: 0, gpax: 0, advisor_name: '',
  // Admin
  admin_firstname: '', admin_lastname: '', position: '',
  // Common
  email: '', phone: ''
};
const form = ref({ ...initialForm });

const fetchData = async () => {
  try {
    const [u, r, m] = await Promise.all([UserAPI.getUsers(), UserAPI.getRoles(), UserAPI.getMajors()]);
    users.value = u || [];
    roles.value = r || [];
    majors.value = m || [];
  } catch (e) { console.error(e); }
};

onMounted(fetchData);

const isStudentRole = computed(() => {
  const r = roles.value.find(role => role.ID === form.value.role_id);
  return r?.name?.toLowerCase() === 'student';
});

const openCreateModal = () => {
  isEditing.value = false;
  editingUserId.value = null;
  form.value = { ...initialForm };
  if (roles.value.length > 0) form.value.role_id = roles.value[0].ID;
  showModal.value = true;
};

const openEditModal = (user: any) => {
  isEditing.value = true;
  editingUserId.value = user.ID;
  form.value = { ...initialForm };
  form.value.username = user.username;
  form.value.role_id = user.RoleID || user.role?.ID;

  if (user.student_profile?.[0]) {
    const s = user.student_profile[0];
    form.value.student_id = s.student_id;
    form.value.national_id = s.national_id;
    form.value.first_name_th = s.first_name_th;
    form.value.last_name_th = s.last_name_th;
    form.value.major_id = s.major_id;
    form.value.gpax = s.gpax;
    form.value.advisor_name = s.advisor_name;
    form.value.email = s.email;
    form.value.phone = s.phone;
  } else if (user.admin_profile?.[0]) {
    const a = user.admin_profile[0];
    // Map ให้ตรงกับ Key ที่ Backend ส่งมา
    form.value.admin_firstname = a.admin_firstname; 
    form.value.admin_lastname = a.admin_lastname;
    form.value.position = a.position;
    form.value.email = a.email;
    form.value.phone = a.phone;
  }
  showModal.value = true;
};

const submitUser = async () => {
  let res;
  if (isEditing.value && editingUserId.value) {
    res = await UserAPI.updateUser(editingUserId.value, form.value);
  } else {
    res = await UserAPI.createUser(form.value);
  }

  // **สำคัญ: แก้ไข Logic การเช็คผลลัพธ์**
  // ถ้าระบบ `https.ts` ของคุณ return `res.data` ในกรณีสำเร็จ
  // เราจะเช็คว่า `res` มีค่า และ `message` มีค่าหรือไม่
  if (res && res.message) {
    Swal.fire({
      icon: 'success',
      title: isEditing.value ? 'แก้ไขสำเร็จ' : 'สร้างสำเร็จ',
      timer: 1500,
      showConfirmButton: false
    });
    showModal.value = false;
    await fetchData(); // รีเฟรชตารางทันที
  } else {
    // กรณี Error
    const errMsg = res?.data?.error || 'เกิดข้อผิดพลาด';
    Swal.fire('ล้มเหลว', errMsg, 'error');
  }
};

const removeUser = async (id: number) => {
  const conf = await Swal.fire({ title: 'ยืนยันการลบ?', showCancelButton: true, confirmButtonColor: '#d33' });
  if (conf.isConfirmed) {
    const res = await UserAPI.deleteUser(id);
    if (res && res.message) {
      Swal.fire('ลบสำเร็จ', '', 'success');
      fetchData();
    } else {
      Swal.fire('ผิดพลาด', res?.data?.error, 'error');
    }
  }
};
</script>

<template>
  <div class="p-6">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold text-gray-800">จัดการข้อมูลผู้ใช้งาน</h1>
      <button @click="openCreateModal" class="btn bg-[#1e3a8a] text-white hover:bg-blue-800 shadow-md">
        + เพิ่มผู้ใช้งาน
      </button>
    </div>

    <div class="bg-white shadow-lg rounded-xl overflow-hidden">
      <div class="overflow-x-auto">
        <table class="table w-full">
          <thead class="bg-gray-100 text-gray-600">
            <tr>
              <th class="py-3 px-6">ID</th>
              <th class="py-3 px-6">Username</th>
              <th class="py-3 px-6">Role</th>
              <th class="py-3 px-6">Name</th>
              <th class="py-3 px-6 text-center">Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="u in users" :key="u.ID" class="hover:bg-gray-50 border-b">
              <td class="py-3 px-6">{{ u.ID }}</td>
              <td class="py-3 px-6 font-bold">{{ u.username }}</td>
              <td class="py-3 px-6"><span class="badge badge-sm uppercase">{{ u.role?.name }}</span></td>
              <td class="py-3 px-6">
                <span v-if="u.student_profile?.length">{{ u.student_profile[0].first_name_th }} {{ u.student_profile[0].last_name_th }}</span>
                <span v-else-if="u.admin_profile?.length">{{ u.admin_profile[0].admin_firstname }} {{ u.admin_profile[0].admin_lastname }}</span>
              </td>
              <td class="py-3 px-6 text-center">
                <button @click="openEditModal(u)" class="btn btn-xs btn-warning text-white mr-2">แก้ไข</button>
                <button @click="removeUser(u.ID)" class="btn btn-xs btn-error text-white">ลบ</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm p-4">
      <div class="bg-white rounded-lg w-full max-w-3xl max-h-[90vh] overflow-y-auto flex flex-col">
        <div class="p-6 border-b sticky top-0 bg-white z-10 flex justify-between items-center">
          <h3 class="text-xl font-bold text-gray-800">{{ isEditing ? 'แก้ไขข้อมูล' : 'เพิ่มผู้ใช้งานใหม่' }}</h3>
          <button @click="showModal = false" class="btn btn-circle btn-ghost btn-sm">✕</button>
        </div>
        
        <div class="p-6 space-y-4">
          <div class="grid grid-cols-2 gap-4">
            <div class="form-control"><label class="label font-bold">Username</label><input v-model="form.username" class="input input-bordered" :disabled="isEditing"/></div>
            <div class="form-control"><label class="label font-bold">Password</label><input v-model="form.password" type="password" class="input input-bordered" :placeholder="isEditing ? 'เว้นว่างถ้าไม่เปลี่ยน' : ''"/></div>
            <div class="form-control col-span-2">
              <label class="label font-bold">Role</label>
              <select v-model.number="form.role_id" class="select select-bordered" :disabled="isEditing">
                <option v-for="r in roles" :key="r.ID" :value="r.ID">{{ r.name }}</option>
              </select>
            </div>
          </div>

          <div class="divider text-gray-400 text-sm">ข้อมูลส่วนตัว</div>

          <div v-if="isStudentRole" class="grid grid-cols-2 gap-4 animate-fade-in">
            <div class="form-control"><label class="label">รหัสนักศึกษา</label><input v-model="form.student_id" class="input input-bordered"/></div>
            <div class="form-control"><label class="label">เลขบัตรประชาชน</label><input v-model="form.national_id" class="input input-bordered"/></div>
            <div class="form-control"><label class="label">ชื่อ (ไทย)</label><input v-model="form.first_name_th" class="input input-bordered"/></div>
            <div class="form-control"><label class="label">นามสกุล (ไทย)</label><input v-model="form.last_name_th" class="input input-bordered"/></div>
            <div class="form-control"><label class="label">สาขาวิชา</label>
               <select v-model.number="form.major_id" class="select select-bordered">
                 <option v-for="m in majors" :key="m.ID" :value="m.ID">{{ m.major_name }}</option>
               </select>
            </div>
            <div class="form-control"><label class="label">เกรดเฉลี่ย</label><input v-model.number="form.gpax" type="number" step="0.01" class="input input-bordered"/></div>
            <div class="form-control col-span-2"><label class="label">อาจารย์ที่ปรึกษา</label><input v-model="form.advisor_name" class="input input-bordered"/></div>
          </div>

          <div v-else class="grid grid-cols-2 gap-4 animate-fade-in">
            <div class="form-control"><label class="label">ชื่อจริง</label><input v-model="form.admin_firstname" class="input input-bordered"/></div>
            <div class="form-control"><label class="label">นามสกุล</label><input v-model="form.admin_lastname" class="input input-bordered"/></div>
            <div class="form-control col-span-2"><label class="label">ตำแหน่ง</label><input v-model="form.position" class="input input-bordered"/></div>
          </div>

          <div class="grid grid-cols-2 gap-4 mt-2">
             <div class="form-control"><label class="label">อีเมล</label><input v-model="form.email" class="input input-bordered"/></div>
             <div class="form-control"><label class="label">เบอร์โทรศัพท์</label><input v-model="form.phone" class="input input-bordered"/></div>
          </div>
        </div>

        <div class="p-6 border-t bg-gray-50 flex justify-end gap-3 sticky bottom-0">
          <button @click="showModal = false" class="btn btn-ghost">ยกเลิก</button>
          <button @click="submitUser" class="btn bg-[#1e3a8a] text-white hover:bg-blue-900 px-8">บันทึก</button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.animate-fade-in { animation: fadeIn 0.3s ease-in-out; }
@keyframes fadeIn { from { opacity: 0; transform: translateY(5px); } to { opacity: 1; transform: translateY(0); } }
</style>