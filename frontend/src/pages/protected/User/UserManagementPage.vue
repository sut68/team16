<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import * as UserAPI from '@/services/api/user';
import Swal from 'sweetalert2';
import StatsGrid from '@/components/ui/StatsGrid.vue';
import type { StatItem } from '@/components/ui/StatsGrid.vue';

// State
const users = ref<any[]>([]);
const roles = ref<any[]>([]);
const majors = ref<any[]>([]);
const showModal = ref(false);
const isEditing = ref(false);
const editingUserId = ref<number | null>(null);

// Form Data
const initialForm = {
  username: '', password: '', role_id: 0,
  // Student Fields
  student_id: '', 
  national_id: '', 
  first_name_th: '', last_name_th: '', 
  first_name_en: '', last_name_en: '', 
  current_year: 1, 
  major_id: 0, 
  gpax: 0, 
  advisor_name: '',
  // Admin Fields
  admin_firstname: '', admin_lastname: '', position: '', phone: '', email: ''
};
const form = ref({ ...initialForm });

// API Fetching
const fetchData = async () => {
  try {
    const [u, r, m] = await Promise.all([UserAPI.getUsers(), UserAPI.getRoles(), UserAPI.getMajors()]);
    users.value = u || [];
    roles.value = r || [];
    majors.value = m || [];
  } catch (e) { console.error(e); }
};

onMounted(fetchData);

// Check Role
const isStudentRole = computed(() => {
  const r = roles.value.find(role => role.ID === form.value.role_id);
  return r?.name?.toLowerCase() === 'student';
});

// User Statistics for StatsGrid
const userStats = computed<StatItem[]>(() => {
  const total = users.value.length;
  const students = users.value.filter(u => u.role?.name?.toLowerCase() === 'student').length;
  const admins = total - students;
  
  return [
    {
      title: 'ผู้ใช้งานทั้งหมด',
      value: total,
      description: 'Total Users',
      icon: 'users',
      color: 'blue'
    },
    {
      title: 'นักศึกษา',
      value: students,
      description: 'Students',
      icon: 'user-round-check',
      color: 'green'
    },
    {
      title: 'ผู้ดูแลระบบ',
      value: admins,
      description: 'Admins / Staff',
      icon: 'user-round-cog',
      color: 'purple'
    }
  ];
});

// ฟังก์ชันคำนวณเลขบัตรประชาชน (สูตรมาตรฐาน)
const checkThaiID = (id: string): boolean => {
  if (!id || id.length !== 13 || !/^\d+$/.test(id)) return false;
  let sum = 0;
  for (let i = 0; i < 12; i++) {
    sum += parseInt(id.charAt(i)) * (13 - i);
  }
  const remainder = sum % 11;
  const checkDigit = (11 - remainder) % 10;
  return checkDigit === parseInt(id.charAt(12));
};

// 🟢 ฟังก์ชันตรวจสอบข้อมูล (Validation) - บังคับกรอกทุกฟิลด์
const validateForm = (): string | null => {
  // 1. Check Common Fields
  if (!form.value.username?.trim()) return 'กรุณาระบุ Username';
  
  // Password บังคับเฉพาะตอนสร้างใหม่ (Edit เว้นว่างได้ถ้าไม่เปลี่ยน)
  if (!isEditing.value && !form.value.password) return 'กรุณาระบุ Password';
  
  if (!form.value.role_id) return 'กรุณาระบุ Role';

  // 2. Check Student Fields
  if (isStudentRole.value) {
    // Regex รหัสนักศึกษา (ขึ้นต้น B,C,M,D ตามด้วยเลข 7 หลัก)
    const studentIdPattern = /^[BCMD]\d{7}$/;
    if (!form.value.student_id?.trim()) return 'กรุณาระบุรหัสนักศึกษา';
    if (!studentIdPattern.test(form.value.student_id)) return 'รหัสนักศึกษาไม่ถูกต้อง (ต้องขึ้นต้น B,C,M,D ตามด้วยเลข 7 หลัก)';

    // ชื่อ-สกุล ไทย
    if (!form.value.first_name_th?.trim()) return 'กรุณาระบุชื่อจริง (ไทย)';
    if (!form.value.last_name_th?.trim()) return 'กรุณาระบุนามสกุล (ไทย)';

    // ชื่อ-สกุล อังกฤษ (ต้องเป็นภาษาอังกฤษเท่านั้น)
    const enNamePattern = /^[a-zA-Z\s]+$/;
    if (!form.value.first_name_en?.trim()) return 'กรุณาระบุ First Name (Eng)';
    if (!enNamePattern.test(form.value.first_name_en)) return 'First Name (Eng) ต้องเป็นภาษาอังกฤษเท่านั้น';
    if (!form.value.last_name_en?.trim()) return 'กรุณาระบุ Last Name (Eng)';
    if (!enNamePattern.test(form.value.last_name_en)) return 'Last Name (Eng) ต้องเป็นภาษาอังกฤษเท่านั้น';

    // บัตรประชาชน (บังคับเช็คสูตร)
    if (!form.value.national_id?.trim()) return 'กรุณาระบุเลขบัตรประชาชน';
    if (!checkThaiID(form.value.national_id)) return 'เลขบัตรประชาชนไม่ถูกต้อง (Checksum ไม่ผ่าน)';

    // การศึกษา & เกรด
    if (!form.value.major_id) return 'กรุณาเลือกสาขาวิชา';
    if (form.value.gpax === null || form.value.gpax === undefined || String(form.value.gpax) === '') return 'กรุณาระบุเกรดเฉลี่ย (GPAX)';
    if (form.value.gpax < 0 || form.value.gpax > 4.00) return 'GPAX ต้องอยู่ระหว่าง 0.00 - 4.00';
    
    if (!form.value.current_year) return 'กรุณาระบุชั้นปี';
    if (form.value.current_year < 1 || form.value.current_year > 8) return 'ชั้นปีต้องอยู่ระหว่าง 1-8';

    if (!form.value.advisor_name?.trim()) return 'กรุณาระบุชื่ออาจารย์ที่ปรึกษา';
  } 
  
  // 3. Check Admin Fields
  else {
    if (!form.value.admin_firstname?.trim()) return 'กรุณาระบุชื่อจริง';
    if (!form.value.admin_lastname?.trim()) return 'กรุณาระบุนามสกุล';
    if (!form.value.position?.trim()) return 'กรุณาระบุตำแหน่ง';
    
    // Email Validation
    const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!form.value.email?.trim()) return 'กรุณาระบุอีเมล';
    if (!emailPattern.test(form.value.email)) return 'รูปแบบอีเมลไม่ถูกต้อง';
    
    // Phone Validation (10 digits)
    if (!form.value.phone?.trim()) return 'กรุณาระบุเบอร์โทรศัพท์';
    if (!/^\d{10}$/.test(form.value.phone)) return 'เบอร์โทรศัพท์ต้องเป็นตัวเลข 10 หลัก';
  }

  return null; // ผ่าน
};

// Open Create Modal
const openCreateModal = () => {
  isEditing.value = false;
  editingUserId.value = null;
  form.value = { ...initialForm };
  // Default to Student role
  const studentRole = roles.value.find(r => r.name.toLowerCase() === 'student');
  if (studentRole) form.value.role_id = studentRole.ID;
  showModal.value = true;
};

// Open Edit Modal
const openEditModal = (user: any) => {
  isEditing.value = true;
  editingUserId.value = user.ID;
  
  // Clone data to form
  form.value = { ...initialForm };
  form.value.username = user.username;
  form.value.role_id = user.RoleID || user.role?.ID;

  if (user.student_profile?.[0]) {
    const s = user.student_profile[0];
    form.value.student_id = s.student_id;
    form.value.national_id = s.national_id;
    form.value.first_name_th = s.first_name_th;
    form.value.last_name_th = s.last_name_th;
    form.value.first_name_en = s.first_name_en;
    form.value.last_name_en = s.last_name_en;
    form.value.current_year = s.current_year;
    form.value.major_id = s.major_id;
    form.value.gpax = s.gpax;
    form.value.advisor_name = s.advisor_name;
  } else if (user.admin_profile?.[0]) {
    const a = user.admin_profile[0];
    form.value.admin_firstname = a.admin_firstname; 
    form.value.admin_lastname = a.admin_lastname;
    form.value.position = a.position;
    form.value.email = a.email;
    form.value.phone = a.phone;
  }
  showModal.value = true;
};

// Submit Logic
const submitUser = async () => {
  const errorMsg = validateForm();
  if (errorMsg) {
    return Swal.fire({
      icon: 'warning',
      title: 'ข้อมูลไม่ถูกต้อง',
      text: errorMsg,
      confirmButtonText: 'ตกลง'
    });
  }

  try {
    let res;
    if (isEditing.value && editingUserId.value) {
      res = await UserAPI.updateUser(editingUserId.value, form.value);
    } else {
      res = await UserAPI.createUser(form.value);
    }

    if (res && res.message) {
      Swal.fire({
        icon: 'success',
        title: isEditing.value ? 'แก้ไขสำเร็จ' : 'สร้างสำเร็จ',
        timer: 1500,
        showConfirmButton: false
      });
      showModal.value = false;
      await fetchData();
    } else {
      throw new Error(res?.data?.error || 'เกิดข้อผิดพลาดจากเซิร์ฟเวอร์');
    }
  } catch (err: any) {
    Swal.fire('ล้มเหลว', err.message, 'error');
  }
};

const removeUser = async (id: number) => {
  const conf = await Swal.fire({ 
    title: 'ยืนยันการลบ?', 
    text: "ข้อมูลที่เกี่ยวข้องทั้งหมดจะหายไป",
    icon: 'warning',
    showCancelButton: true, 
    confirmButtonColor: '#d33',
    confirmButtonText: 'ลบ',
    cancelButtonText: 'ยกเลิก'
  });
  
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
  <div class="w-full mx-auto flex flex-col h-full p-6 bg-white rounded-tl-[30px] shadow" data-theme="light">
    
    <div class="flex items-center justify-between mb-4">
      <h1 class="text-2xl font-bold text-[#1e3a8a]">จัดการข้อมูลผู้ใช้งาน</h1>
      <button 
        @click="openCreateModal" 
        class="btn btn-sm bg-white border border-gray-300 text-gray-700 hover:bg-gray-100 hover:border-gray-400 flex items-center justify-center gap-2 rounded-full px-5 h-10 shadow-sm transition-all"
      >
        <span class="text-xl leading-none font-bold text-blue-600">+</span>
        <span class="font-medium">เพิ่มผู้ใช้งาน</span>
      </button>
    </div>

    <!-- Stats Section -->
    <StatsGrid :stats="userStats" class="mb-6" />

    <div class="flex-1 min-h-0 flex flex-col gap-4">
      <div class="table-scroll overflow-x-auto overflow-y-auto min-h-0 bg-white rounded flex-1">
        <table class="table w-full">
          <thead class="bg-white/95">
            <tr class="sticky top-0 z-10 shadow-sm">
              <th class="px-4 py-3 text-left font-semibold text-gray-600">ID</th>
              <th class="px-4 py-3 text-left font-semibold text-gray-600">Username</th>
              <th class="px-4 py-3 text-left font-semibold text-gray-600">Role</th>
              <th class="px-4 py-3 text-left font-semibold text-gray-600">Name (TH)</th>
              <th class="px-4 py-3 text-left font-semibold text-gray-600">Major / Year</th>
              <th class="px-4 py-3 text-center font-semibold text-gray-600">Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="u in users" :key="u.ID" class="odd:bg-gray-50 hover:bg-gray-100 transition-colors">
              <td class="px-4 py-3 text-sm text-gray-600">{{ u.ID }}</td>
              <td class="px-4 py-3 text-sm font-bold text-[#1e3a8a]">{{ u.username }}</td>
              <td class="px-4 py-3">
                <span 
                  class="px-2 py-1 text-xs font-medium rounded-full uppercase"
                  :class="u.role?.name?.toLowerCase() === 'student' 
                    ? 'bg-green-100 text-green-700' 
                    : 'bg-purple-100 text-purple-700'"
                >
                  {{ u.role?.name }}
                </span>
              </td>
              <td class="px-4 py-3 text-sm text-gray-700">
                <span v-if="u.student_profile?.length">
                  {{ u.student_profile[0].first_name_th }} {{ u.student_profile[0].last_name_th }}
                </span>
                <span v-else-if="u.admin_profile?.length">
                  {{ u.admin_profile[0].admin_firstname }} {{ u.admin_profile[0].admin_lastname }}
                </span>
                <span v-else class="text-gray-400">-</span>
              </td>
              <td class="px-4 py-3 text-sm text-gray-600">
                <div v-if="u.student_profile?.length">
                  <div>{{ u.student_profile[0].major?.major_name || '-' }}</div>
                  <div class="text-xs text-gray-400">ชั้นปีที่ {{ u.student_profile[0].current_year || '-' }}</div>
                </div>
                <div v-else class="text-gray-400">-</div>
              </td>
              <td class="px-4 py-3 text-center">
                <div class="flex items-center justify-center gap-2">
                  <button @click="openEditModal(u)" class="btn btn-xs btn-circle btn-ghost text-yellow-600 hover:bg-yellow-100" title="แก้ไข">
                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M17 3a2.828 2.828 0 1 1 4 4L7.5 20.5 2 22l1.5-5.5L17 3z"></path></svg>
                  </button>
                  <button @click="removeUser(u.ID)" class="btn btn-xs btn-circle btn-ghost text-red-600 hover:bg-red-100" title="ลบ">
                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 6h18"></path><path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"></path><path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"></path></svg>
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <teleport to="body">
      <div v-if="showModal" class="fixed inset-0 z-[200] flex items-center justify-center bg-black/50 backdrop-blur-sm p-4 animate-fade-in" @click.self="showModal = false">
        <div class="bg-white w-full max-w-4xl h-[90vh] rounded-2xl shadow-2xl flex flex-col overflow-hidden animate-pop-in">
          
          <div class="px-6 py-4 border-b flex items-center justify-between bg-slate-50">
            <h3 class="text-xl font-bold text-[#1e3a8a]">
              {{ isEditing ? 'แก้ไขข้อมูลผู้ใช้งาน' : 'เพิ่มผู้ใช้งานใหม่' }}
            </h3>
            <button @click="showModal = false" class="btn btn-circle btn-ghost btn-sm text-gray-500">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" /></svg>
            </button>
          </div>

          <div class="flex-1 overflow-y-auto p-6 bg-[#f8fafc]">
             <div class="card bg-white shadow-sm border border-gray-100">
                <div class="card-body p-5">
                   
                   <h4 class="font-bold text-gray-800 mb-3 flex items-center gap-2 text-sm uppercase tracking-wide">
                      <span class="w-2 h-2 rounded-full bg-[#1e3a8a]"></span> ข้อมูลบัญชี (Account)
                   </h4>
                   <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
                      <div class="form-control">
                         <label class="label text-gray-600 text-sm mb-1">Username <span class="text-red-500">*</span></label>
                         <input v-model="form.username" class="input input-bordered w-full" placeholder="เช่น รหัสนักศึกษา" />
                      </div>
                      <div class="form-control">
                         <label class="label text-gray-600 text-sm mb-1">Password</label>
                         <input v-model="form.password" type="password" class="input input-bordered w-full" :placeholder="isEditing ? 'เว้นว่างถ้าไม่ต้องการเปลี่ยน' : 'กำหนดรหัสผ่าน'" />
                      </div>
                      <div class="form-control md:col-span-2">
                         <label class="label text-gray-600 text-sm mb-1">Role <span class="text-red-500">*</span></label>
                         <select v-model.number="form.role_id" class="select select-bordered w-full bg-white" :disabled="isEditing">
                            <option v-for="r in roles" :key="r.ID" :value="r.ID">{{ r.name }}</option>
                         </select>
                      </div>
                   </div>

                   <div v-if="isStudentRole" class="animate-fade-in space-y-6">
                      
                      <div>
                        <h4 class="font-bold text-gray-800 mb-3 flex items-center gap-2 text-sm uppercase tracking-wide">
                           <span class="w-2 h-2 rounded-full bg-orange-400"></span> ข้อมูลระบุตัวตน
                        </h4>
                        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                           <div class="form-control"><label class="label text-gray-600 text-sm">รหัสนักศึกษา <span class="text-red-500">*</span></label><input v-model="form.student_id" class="input input-bordered w-full"/></div>
                           <div class="form-control"><label class="label text-gray-600 text-sm">เลขบัตรประชาชน <span class="text-red-500">*</span></label><input v-model="form.national_id" class="input input-bordered w-full" maxlength="13"/></div>
                        </div>
                      </div>

                      <div>
                        <h4 class="font-bold text-gray-800 mb-3 flex items-center gap-2 text-sm uppercase tracking-wide">
                           <span class="w-2 h-2 rounded-full bg-green-500"></span> ข้อมูลชื่อ-นามสกุล
                        </h4>
                        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                           <div class="form-control"><label class="label text-gray-600 text-sm">ชื่อ (ภาษาไทย) <span class="text-red-500">*</span></label><input v-model="form.first_name_th" class="input input-bordered w-full"/></div>
                           <div class="form-control"><label class="label text-gray-600 text-sm">นามสกุล (ภาษาไทย) <span class="text-red-500">*</span></label><input v-model="form.last_name_th" class="input input-bordered w-full"/></div>
                           
                           <div class="form-control"><label class="label text-gray-600 text-sm">First Name (Eng) <span class="text-red-500">*</span></label><input v-model="form.first_name_en" class="input input-bordered w-full"/></div>
                           <div class="form-control"><label class="label text-gray-600 text-sm">Last Name (Eng) <span class="text-red-500">*</span></label><input v-model="form.last_name_en" class="input input-bordered w-full"/></div>
                        </div>
                      </div>

                      <div>
                        <h4 class="font-bold text-gray-800 mb-3 flex items-center gap-2 text-sm uppercase tracking-wide">
                           <span class="w-2 h-2 rounded-full bg-blue-500"></span> ข้อมูลการศึกษา
                        </h4>
                        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                           <div class="form-control">
                              <label class="label text-gray-600 text-sm">สาขาวิชา <span class="text-red-500">*</span></label>
                              <select v-model.number="form.major_id" class="select select-bordered w-full bg-white">
                                 <option :value="0">-- เลือกสาขา --</option>
                                 <option v-for="m in majors" :key="m.ID" :value="m.ID">{{ m.major_name }}</option>
                              </select>
                           </div>
                           <div class="form-control">
                              <label class="label text-gray-600 text-sm">ชั้นปีที่ (Current Year) <span class="text-red-500">*</span></label>
                              <input v-model.number="form.current_year" type="number" min="1" max="8" class="input input-bordered w-full"/>
                           </div>
                           <div class="form-control"><label class="label text-gray-600 text-sm">เกรดเฉลี่ยสะสม (GPAX) <span class="text-red-500">*</span></label><input v-model.number="form.gpax" type="number" step="0.01" min="0" max="4" class="input input-bordered w-full"/></div>
                           <div class="form-control"><label class="label text-gray-600 text-sm">อาจารย์ที่ปรึกษา <span class="text-red-500">*</span></label><input v-model="form.advisor_name" class="input input-bordered w-full"/></div>
                        </div>
                      </div>
                   </div>

                   <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-4 animate-fade-in">
                      <div class="form-control"><label class="label text-gray-600 text-sm">ชื่อจริง <span class="text-red-500">*</span></label><input v-model="form.admin_firstname" class="input input-bordered w-full"/></div>
                      <div class="form-control"><label class="label text-gray-600 text-sm">นามสกุล <span class="text-red-500">*</span></label><input v-model="form.admin_lastname" class="input input-bordered w-full"/></div>
                      <div class="form-control md:col-span-2"><label class="label text-gray-600 text-sm">ตำแหน่ง <span class="text-red-500">*</span></label><input v-model="form.position" class="input input-bordered w-full"/></div>
                      
                      <div class="divider md:col-span-2 text-xs text-gray-400">ข้อมูลติดต่อ (สำหรับเจ้าหน้าที่)</div>
                      
                      <div class="form-control"><label class="label text-gray-600 text-sm">อีเมล <span class="text-red-500">*</span></label><input v-model="form.email" class="input input-bordered w-full"/></div>
                      <div class="form-control"><label class="label text-gray-600 text-sm">เบอร์โทรศัพท์ <span class="text-red-500">*</span></label><input v-model="form.phone" class="input input-bordered w-full" maxlength="10"/></div>
                   </div>

                </div>
             </div>
          </div>

          <div class="px-6 py-4 border-t bg-slate-50 flex items-center justify-end gap-2">
            <button @click="showModal = false" class="btn btn-sm btn-ghost text-gray-700 rounded-full px-6 hover:bg-gray-200">ยกเลิก</button>
            <button @click="submitUser" class="btn btn-sm bg-[#1e3a8a] border-none text-white rounded-full px-6 hover:bg-blue-900 shadow-sm">บันทึกข้อมูล</button>
          </div>
        </div>
      </div>
    </teleport>

  </div>
</template>

<style scoped>
.table-scroll { scrollbar-gutter: stable both-edges; }
.table-scroll::-webkit-scrollbar { width: 8px; height: 8px; }
.table-scroll::-webkit-scrollbar-thumb { background: rgba(0,0,0,0.15); border-radius: 999px; }

@keyframes pop-in {
  from { opacity: 0; transform: translateY(8px) scale(0.98); }
  to { opacity: 1; transform: translateY(0) scale(1); }
}
.animate-pop-in { animation: pop-in 0.15s cubic-bezier(0.2, 0.9, 0.3, 1); }

@keyframes fadeIn {
  from { opacity: 0; } to { opacity: 1; }
}
.animate-fade-in { animation: fadeIn 0.15s ease-out; }
</style>