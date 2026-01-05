<script setup lang="ts">
import { ref, watch } from 'vue';
import * as UserAPI from '@/services/api/user';
import Swal from 'sweetalert2';

const props = defineProps<{ initialData: any }>();
const emit = defineEmits(['refresh']);

const isEditing = ref(false);
const activeTab = ref('personal');
const form = ref<any>({});
const isSameAddress = ref(false);
const thaiProvinces = [
  "กรุงเทพมหานคร", "กระบี่", "กาญจนบุรี", "กาฬสินธุ์", "กำแพงเพชร", "ขอนแก่น", "จันทบุรี", "ฉะเชิงเทรา", "ชลบุรี", "ชัยนาท",
  "ชัยภูมิ", "ชุมพร", "เชียงราย", "เชียงใหม่", "ตรัง", "ตราด", "ตาก", "นครนายก", "นครปฐม", "นครพนม",
  "นครราชสีมา", "นครศรีธรรมราช", "นครสวรรค์", "นนทบุรี", "นราธิวาส", "น่าน", "บึงกาฬ", "บุรีรัมย์", "ปทุมธานี", "ประจวบคีรีขันธ์",
  "ปราจีนบุรี", "ปัตตานี", "พระนครศรีอยุธยา", "พะเยา", "พังงา", "พัทลุง", "พิจิตร", "พิษณุโลก", "เพชรบุรี", "เพชรบูรณ์",
  "แพร่", "ภูเก็ต", "มหาสารคาม", "มุกดาหาร", "แม่ฮ่องสอน", "ยโสธร", "ยะลา", "ร้อยเอ็ด", "ระนอง", "ระยอง",
  "ราชบุรี", "ลพบุรี", "ลำปาง", "ลำพูน", "เลย", "ศรีสะเกษ", "สกลนคร", "สงขลา", "สตูล", "สมุทรปราการ",
  "สมุทรสงคราม", "สมุทรสาคร", "สระแก้ว", "สระบุรี", "สิงห์บุรี", "สุโขทัย", "สุพรรณบุรี", "สุราษฎร์ธานี", "สุรินทร์", "หนองคาย",
  "หนองบัวลำภู", "อ่างทอง", "อำนาจเจริญ", "อุดรธานี", "อุตรดิตถ์", "อุทัยธานี", "อุบลราชธานี"
];
// 1. เมื่อกด Checkbox "เหมือนที่อยู่ตามทะเบียนบ้าน"
watch(isSameAddress, (checked) => {
  if (checked) {
    form.value.current_address = form.value.permanent_address;
  }
});

// 2. เมื่อพิมพ์แก้ "ที่อยู่ตามทะเบียนบ้าน" (ให้ Sync ไปที่อยู่ปัจจุบันด้วย ถ้าติ๊กอยู่)
watch(() => form.value.permanent_address, (newVal) => {
  if (isSameAddress.value) {
    form.value.current_address = newVal;
  }
});

watch(() => props.initialData, (newVal) => {
  if (newVal) {
    // Deep clone to avoid mutating props directly
    form.value = JSON.parse(JSON.stringify(newVal));

    // Format date for input type="date"
    if (form.value.birth_date) {
      const dateOnly = form.value.birth_date.split('T')[0];

      // เช็คว่าถ้าเป็นปี 0001 (ค่า Default ของ Go) ให้เคลียร์ค่าเป็นว่าง
      if (dateOnly === '0001-01-01') {
        form.value.birth_date = '';
      } else {
        form.value.birth_date = dateOnly;
      }
    }

    // Check Address Consistency on Load
    // ถ้ามีข้อมูลและข้อมูลเหมือนกัน ให้ติ๊กถูกอัตโนมัติ
    if (form.value.permanent_address &&
      form.value.permanent_address === form.value.current_address) {
      isSameAddress.value = true;
    } else {
      isSameAddress.value = false;
    }

    initFamily();
  }
}, { immediate: true });

function initFamily() {
  if (!form.value.family_info) {
    form.value.family_info = {
      father_name: '', father_occupation: '', father_income: 0,
      mother_name: '', mother_occupation: '', mother_income: 0,
      guardian_name: '', guardian_relation: '', guardian_occupation: '', guardian_income: 0,
      guardian_is_parent: '' // "father", "mother", or "other"
    };
  }
  // Ensure guardian_is_parent exists for existing records
  if (form.value.family_info && form.value.family_info.guardian_is_parent === undefined) {
    form.value.family_info.guardian_is_parent = '';
  }
}

// Handle guardian type change - copy data from father/mother if selected
function onGuardianTypeChange() {
  const guardianType = form.value.family_info.guardian_is_parent;
  
  if (guardianType === 'father') {
    // Copy father data to guardian fields
    form.value.family_info.guardian_name = form.value.family_info.father_name;
    form.value.family_info.guardian_relation = 'บิดา';
    form.value.family_info.guardian_occupation = form.value.family_info.father_occupation;
    form.value.family_info.guardian_income = form.value.family_info.father_income;
  } else if (guardianType === 'mother') {
    // Copy mother data to guardian fields
    form.value.family_info.guardian_name = form.value.family_info.mother_name;
    form.value.family_info.guardian_relation = 'มารดา';
    form.value.family_info.guardian_occupation = form.value.family_info.mother_occupation;
    form.value.family_info.guardian_income = form.value.family_info.mother_income;
  } else if (guardianType === 'other') {
    // Clear guardian fields for new input
    form.value.family_info.guardian_name = '';
    form.value.family_info.guardian_relation = '';
    form.value.family_info.guardian_occupation = '';
    form.value.family_info.guardian_income = 0;
  }
}


// STRICT VALIDATION FUNCTION (Checks EVERY field)
const validateForm = (): string | null => {
  const f = form.value;
  const fam = f.family_info || {};

  // --- 1. Personal Information Validation ---
  if (!f.first_name_th?.trim()) return 'กรุณาระบุชื่อ (ภาษาไทย)';
  if (!f.last_name_th?.trim()) return 'กรุณาระบุนามสกุล (ภาษาไทย)';

  // English Name Check (A-Z only)
  const enPattern = /^[a-zA-Z\s]+$/;
  if (!f.first_name_en?.trim()) return 'กรุณาระบุ First Name';
  if (!enPattern.test(f.first_name_en)) return 'First Name ต้องเป็นภาษาอังกฤษเท่านั้น';
  if (!f.last_name_en?.trim()) return 'กรุณาระบุ Last Name';
  if (!enPattern.test(f.last_name_en)) return 'Last Name ต้องเป็นภาษาอังกฤษเท่านั้น';

  // Birth Date
  if (!f.birth_date) return 'กรุณาระบุวันเกิด';

  // Contact Info
  const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  if (!f.email?.trim()) return 'กรุณาระบุอีเมล';
  if (!emailPattern.test(f.email)) return 'รูปแบบอีเมลไม่ถูกต้อง';

  const phonePattern = /^\d{10}$/;
  if (!f.phone) return 'กรุณาระบุเบอร์โทรศัพท์';
  if (!phonePattern.test(f.phone)) return 'เบอร์โทรศัพท์ต้องเป็นตัวเลข 10 หลัก';

  // Address Info
  if (!f.province?.trim()) return 'กรุณาระบุจังหวัดภูมิลำเนา';
  if (!f.permanent_address?.trim()) return 'กรุณาระบุที่อยู่ตามทะเบียนบ้าน';
  if (!f.current_address?.trim()) return 'กรุณาระบุที่อยู่ปัจจุบัน';

  if (f.siblings_count === undefined || f.siblings_count === null || f.siblings_count < 0) return 'กรุณาระบุจำนวนพี่น้อง (ใส่ 0 หากไม่มี)';

  // --- 2. Family Information Validation (Mandatory) ---

  // Father
  if (!fam.father_name?.trim()) return 'กรุณาระบุชื่อ-สกุล บิดา';
  if (!fam.father_occupation?.trim()) return 'กรุณาระบุอาชีพ บิดา';
  if (fam.father_income === undefined || fam.father_income === null || fam.father_income < 0) return 'กรุณาระบุรายได้บิดา (ใส่ 0 หากไม่มี)';

  // Mother
  if (!fam.mother_name?.trim()) return 'กรุณาระบุชื่อ-สกุล มารดา';
  if (!fam.mother_occupation?.trim()) return 'กรุณาระบุอาชีพ มารดา';
  if (fam.mother_income === undefined || fam.mother_income === null || fam.mother_income < 0) return 'กรุณาระบุรายได้มารดา (ใส่ 0 หากไม่มี)';

  // Guardian - first check if guardian type is selected
  if (!fam.guardian_is_parent) return 'กรุณาเลือกว่าผู้ปกครองคือใคร (บิดา/มารดา/บุคคลอื่น)';
  
  // Only validate other guardian fields if user selected "other"
  if (fam.guardian_is_parent === 'other') {
    if (!fam.guardian_name?.trim()) return 'กรุณาระบุชื่อ-สกุล ผู้ปกครอง';
    if (!fam.guardian_relation?.trim()) return 'กรุณาระบุความเกี่ยวข้องกับผู้ปกครอง';
    if (!fam.guardian_occupation?.trim()) return 'กรุณาระบุอาชีพ ผู้ปกครอง';
    if (fam.guardian_income === undefined || fam.guardian_income === null || fam.guardian_income < 0) return 'กรุณาระบุรายได้ผู้ปกครอง (ใส่ 0 หากไม่มี)';
  }

  return null; // Passed validation
};

const save = async () => {
  const errorMsg = validateForm();

  if (errorMsg) {
    return Swal.fire({
      icon: 'warning',
      title: 'ข้อมูลไม่ครบถ้วน',
      text: errorMsg,
      confirmButtonText: 'ตกลง',
      confirmButtonColor: '#1e3a8a'
    });
  }

  // Prepare payload
  if (form.value.birth_date) {
    // Ensure date is in ISO format if needed by backend, or keep YYYY-MM-DD
    form.value.birth_date = new Date(form.value.birth_date).toISOString();
  }

  const payload = { ...form.value, family_info: form.value.family_info };

  try {
    const res = await UserAPI.updateMyProfile(payload);
    if (res && (res.status === 200 || res.message)) {
      Swal.fire({
        icon: 'success',
        title: 'บันทึกสำเร็จ',
        timer: 1500,
        showConfirmButton: false
      });
      isEditing.value = false;
      emit('refresh');
    } else {
      throw new Error('บันทึกไม่สำเร็จ');
    }
  } catch (err: any) {
    Swal.fire('ผิดพลาด', err?.response?.data?.error || err.message, 'error');
  }
};
</script>

<template>
  <div class="bg-white p-2 animate-fade-in relative">

    <div class="flex flex-col md:flex-row justify-between items-start md:items-center mb-6 pb-4 border-b gap-4">
      <div>
        <h2 class="text-xl font-bold text-[#1e3a8a]">ข้อมูลส่วนตัวนักศึกษา</h2>
        <p class="text-xs text-gray-500 mt-1">จัดการข้อมูลประวัติส่วนตัวและครอบครัว</p>
      </div>
      <div>
        <button v-if="!isEditing" @click="isEditing = true"
          class="btn btn-sm btn-warning text-white rounded-full px-6 shadow-sm hover:bg-yellow-500 transition-colors">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24"
            stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
          </svg>
          แก้ไขข้อมูล
        </button>
        <div v-else class="flex gap-2">
          <button @click="isEditing = false; emit('refresh')"
            class="btn btn-sm btn-ghost rounded-full text-gray-500 hover:bg-gray-100">ยกเลิก</button>
          <button @click="save"
            class="btn btn-sm bg-[#1e3a8a] border-none text-white rounded-full px-6 shadow-sm hover:bg-blue-900 transition-colors flex items-center gap-1">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24"
              stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
            บันทึก
          </button>
        </div>
      </div>
    </div>

    <div class="flex items-center gap-4 mb-6 border-b border-gray-100">
      <button class="pb-2 text-sm font-medium border-b-2 transition-all duration-200"
        :class="activeTab === 'personal' ? 'border-[#1e3a8a] text-[#1e3a8a]' : 'border-transparent text-gray-500 hover:text-gray-700'"
        @click="activeTab = 'personal'">
        ข้อมูลส่วนตัว
      </button>
      <button class="pb-2 text-sm font-medium border-b-2 transition-all duration-200"
        :class="activeTab === 'family' ? 'border-[#1e3a8a] text-[#1e3a8a]' : 'border-transparent text-gray-500 hover:text-gray-700'"
        @click="activeTab = 'family'">
        ข้อมูลครอบครัว
      </button>
    </div>

    <div v-if="activeTab === 'personal'" class="space-y-6 animate-fade-in">

      <div class="bg-blue-50/50 p-6 rounded-2xl border border-blue-100 relative">
        <h3 class="text-sm font-bold text-[#1e3a8a] mb-4 flex items-center gap-2">
          <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 24 24">
            <path d="M12 3L1 9l11 6 9-4.91V17h2V9M5 13.18v4L12 21l7-3.82v-4L12 17l-7-3.82z" />
          </svg>
          ข้อมูลทางการศึกษา (แก้ไขไม่ได้)
        </h3>

        <div class="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
          <div class="form-control">
            <span class="text-[10px] font-bold text-gray-400 uppercase">รหัสนักศึกษา</span>
            <span class="font-semibold text-gray-800">{{ form.student_id || '-' }}</span>
          </div>
          <div class="form-control">
            <span class="text-[10px] font-bold text-gray-400 uppercase">ชั้นปีที่ (Year)</span>
            <span class="font-semibold text-gray-800">{{ form.current_year || '-' }}</span>
          </div>
          <div class="form-control md:col-span-2">
            <span class="text-[10px] font-bold text-gray-400 uppercase">สาขาวิชา</span>
            <span class="font-semibold text-gray-800">{{ form.major?.major_name || '-' }}</span>
          </div>
          <div class="form-control">
            <span class="text-[10px] font-bold text-gray-400 uppercase">GPAX</span>
            <span class="font-semibold text-gray-800">{{ form.gpax || '0.00' }}</span>
          </div>
          <div class="form-control">
            <span class="text-[10px] font-bold text-gray-400 uppercase">เลขบัตรประชาชน</span>
            <span class="font-semibold text-gray-800">{{ form.national_id || '-' }}</span>
          </div>
          <div class="form-control md:col-span-2">
            <span class="text-[10px] font-bold text-gray-400 uppercase">อาจารย์ที่ปรึกษา</span>
            <span class="font-semibold text-gray-800">{{ form.advisor_name || '-' }}</span>
          </div>
        </div>
      </div>

      <div>
        <h3 class="text-sm font-bold text-gray-800 mb-4 pl-2 border-l-4 border-yellow-400">ข้อมูลทั่วไป (แก้ไขได้)</h3>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div class="form-control">
            <label class="label text-sm text-gray-600">ชื่อ (ภาษาไทย) <span v-if="isEditing"
                class="text-red-500">*</span></label>
            <input v-model="form.first_name_th" :disabled="!isEditing" class="input input-bordered w-full" />
          </div>
          <div class="form-control">
            <label class="label text-sm text-gray-600">นามสกุล (ภาษาไทย) <span v-if="isEditing"
                class="text-red-500">*</span></label>
            <input v-model="form.last_name_th" :disabled="!isEditing" class="input input-bordered w-full" />
          </div>

          <div class="form-control">
            <label class="label text-sm text-gray-600">First Name (EN) <span v-if="isEditing"
                class="text-red-500">*</span></label>
            <input v-model="form.first_name_en" :disabled="!isEditing" class="input input-bordered w-full"
              placeholder="Ex. Somchai" />
          </div>
          <div class="form-control">
            <label class="label text-sm text-gray-600">Last Name (EN) <span v-if="isEditing"
                class="text-red-500">*</span></label>
            <input v-model="form.last_name_en" :disabled="!isEditing" class="input input-bordered w-full"
              placeholder="Ex. Jaidee" />
          </div>

          <div class="form-control">
            <label class="label text-sm text-gray-600">วันเกิด <span v-if="isEditing"
                class="text-red-500">*</span></label>
            <input v-model="form.birth_date" type="date" :disabled="!isEditing" class="input input-bordered w-full" />
          </div>
          <div class="form-control">
            <label class="label text-sm text-gray-600">
              จังหวัดภูมิลำเนา <span v-if="isEditing" class="text-red-500">*</span>
            </label>
            <select v-model="form.province" :disabled="!isEditing"
              class="select select-bordered w-full font-normal text-base"
              :class="{ 'text-gray-400': !form.province && isEditing }">
              <option disabled value="">กรุณาเลือกจังหวัด</option>
              <option v-for="province in thaiProvinces" :key="province" :value="province">
                {{ province }}
              </option>
            </select>
          </div>
          <div class="form-control">
            <label class="label text-sm text-gray-600">อีเมล <span v-if="isEditing"
                class="text-red-500">*</span></label>
            <input v-model="form.email" :disabled="!isEditing" class="input input-bordered w-full" />
          </div>
          <div class="form-control">
            <label class="label text-sm text-gray-600">เบอร์โทรศัพท์ <span v-if="isEditing"
                class="text-red-500">*</span></label>
            <input v-model="form.phone" :disabled="!isEditing" maxlength="10" class="input input-bordered w-full"
              placeholder="08xxxxxxxx" />
          </div>
          <div class="form-control">
            <label class="label text-sm text-gray-600">จำนวนพี่น้อง (คน) <span v-if="isEditing"
                class="text-red-500">*</span></label>
            <input v-model.number="form.siblings_count" type="number" min="0" :disabled="!isEditing"
              class="input input-bordered w-full" />
          </div>

          <div class="hidden md:block"></div>

          <div class="form-control md:col-span-2">
            <label class="label text-sm text-gray-600">ที่อยู่ตามทะเบียนบ้าน <span v-if="isEditing"
                class="text-red-500">*</span></label>
            <textarea v-model="form.permanent_address" :disabled="!isEditing" class="textarea textarea-bordered w-full"
              rows="2"></textarea>
          </div>

          <div class="md:col-span-2" v-if="isEditing">
            <div class="flex items-center gap-2 mt-[-10px] mb-2 select-none">
              <input type="checkbox" id="sameAddress" v-model="isSameAddress"
                class="checkbox checkbox-sm checkbox-primary" />
              <label for="sameAddress" class="text-sm text-gray-600 cursor-pointer">
                ที่อยู่ปัจจุบันเหมือนกับที่อยู่ตามทะเบียนบ้าน
              </label>
            </div>
          </div>

          <div class="form-control md:col-span-2">
            <label class="label text-sm text-gray-600">ที่อยู่ปัจจุบัน <span v-if="isEditing"
                class="text-red-500">*</span></label>
            <textarea v-model="form.current_address" :disabled="!isEditing || isSameAddress"
              class="textarea textarea-bordered w-full" :class="{ 'bg-gray-100': !isEditing || isSameAddress }"
              rows="2"></textarea>
          </div>
        </div>
      </div>
    </div>

    <div v-if="activeTab === 'family'" class="space-y-6 animate-fade-in">
      <div class="alert text-xs text-blue-700 bg-blue-50 border border-blue-200 flex items-start">
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
          class="stroke-current shrink-0 w-5 h-5 mr-2 mt-0.5">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
        </svg>
        <span>กรุณากรอกข้อมูลให้ครบถ้วนทุกช่องเพื่อประโยชน์สูงสุดในการพิจารณาทุนการศึกษา (หากไม่มีรายได้ให้ใส่ 0)</span>
      </div>

      <div class="p-5 border rounded-2xl bg-white hover:border-blue-200 transition-colors shadow-sm">
        <h4 class="font-bold text-[#1e3a8a] border-b pb-2 mb-3 text-sm flex justify-between">
          ข้อมูลบิดา
        </h4>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div class="form-control">
            <label class="label text-xs">ชื่อ-นามสกุล <span v-if="isEditing" class="text-red-500">*</span></label>
            <input v-model="form.family_info.father_name" :disabled="!isEditing"
              class="input input-sm input-bordered w-full" />
          </div>
          <div class="form-control">
            <label class="label text-xs">อาชีพ <span v-if="isEditing" class="text-red-500">*</span></label>
            <input v-model="form.family_info.father_occupation" :disabled="!isEditing"
              class="input input-sm input-bordered w-full" />
          </div>
          <div class="form-control">
            <label class="label text-xs">รายได้/เดือน (บาท) <span v-if="isEditing" class="text-red-500">*</span></label>
            <input v-model.number="form.family_info.father_income" type="number" min="0" :disabled="!isEditing"
              class="input input-sm input-bordered w-full" />
          </div>
        </div>
      </div>

      <div class="p-5 border rounded-2xl bg-white hover:border-blue-200 transition-colors shadow-sm">
        <h4 class="font-bold text-[#1e3a8a] border-b pb-2 mb-3 text-sm">ข้อมูลมารดา</h4>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div class="form-control">
            <label class="label text-xs">ชื่อ-นามสกุล <span v-if="isEditing" class="text-red-500">*</span></label>
            <input v-model="form.family_info.mother_name" :disabled="!isEditing"
              class="input input-sm input-bordered w-full" />
          </div>
          <div class="form-control">
            <label class="label text-xs">อาชีพ <span v-if="isEditing" class="text-red-500">*</span></label>
            <input v-model="form.family_info.mother_occupation" :disabled="!isEditing"
              class="input input-sm input-bordered w-full" />
          </div>
          <div class="form-control">
            <label class="label text-xs">รายได้/เดือน (บาท) <span v-if="isEditing" class="text-red-500">*</span></label>
            <input v-model.number="form.family_info.mother_income" type="number" min="0" :disabled="!isEditing"
              class="input input-sm input-bordered w-full" />
          </div>
        </div>
      </div>

      <div class="p-5 border rounded-2xl bg-white hover:border-blue-200 transition-colors shadow-sm">
        <h4 class="font-bold text-[#1e3a8a] border-b pb-2 mb-3 text-sm">ผู้ปกครอง</h4>
        
        <!-- Guardian Type Selector -->
        <div class="form-control mb-4">
          <label class="label text-xs font-semibold text-gray-700">ผู้ปกครองของนักศึกษา <span v-if="isEditing" class="text-red-500">*</span></label>
          <select 
            v-model="form.family_info.guardian_is_parent" 
            :disabled="!isEditing"
            class="select select-sm select-bordered w-full md:w-1/2"
            @change="onGuardianTypeChange"
          >
            <option value="">-- กรุณาเลือก --</option>
            <option value="father">บิดา (ใช้ข้อมูลจากบิดาข้างต้น)</option>
            <option value="mother">มารดา (ใช้ข้อมูลจากมารดาข้างต้น)</option>
            <option value="other">บุคคลอื่น (ระบุข้อมูลเพิ่มเติม)</option>
          </select>
        </div>

        <!-- Info box when father/mother is selected -->
        <div v-if="form.family_info.guardian_is_parent === 'father' || form.family_info.guardian_is_parent === 'mother'" 
             class="alert bg-green-50 border border-green-200 text-green-700 text-xs mb-4">
          <svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-5 w-5" fill="none" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <span>
            ระบบจะใช้ข้อมูลจาก<strong>{{ form.family_info.guardian_is_parent === 'father' ? 'บิดา' : 'มารดา' }}</strong>เป็นผู้ปกครอง
          </span>
        </div>

        <!-- Guardian fields - only show when 'other' is selected -->
        <div v-if="form.family_info.guardian_is_parent === 'other'" class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div class="form-control">
            <label class="label text-xs">ชื่อ-นามสกุลผู้ปกครอง <span v-if="isEditing" class="text-red-500">*</span></label>
            <input v-model="form.family_info.guardian_name" :disabled="!isEditing"
              class="input input-sm input-bordered w-full" placeholder="เช่น นายสมศักดิ์ ใจดี" />
          </div>
          <div class="form-control">
            <label class="label text-xs">ความเกี่ยวข้อง <span v-if="isEditing" class="text-red-500">*</span></label>
            <input v-model="form.family_info.guardian_relation" :disabled="!isEditing"
              class="input input-sm input-bordered w-full" placeholder="เช่น ปู่, ย่า, ลุง, ป้า" />
          </div>
          <div class="form-control">
            <label class="label text-xs">อาชีพ <span v-if="isEditing" class="text-red-500">*</span></label>
            <input v-model="form.family_info.guardian_occupation" :disabled="!isEditing"
              class="input input-sm input-bordered w-full" />
          </div>
          <div class="form-control">
            <label class="label text-xs">รายได้/เดือน (บาท) <span v-if="isEditing" class="text-red-500">*</span></label>
            <input v-model.number="form.family_info.guardian_income" type="number" min="0" :disabled="!isEditing"
              class="input input-sm input-bordered w-full" />
          </div>
        </div>

        <!-- Summary when father/mother is selected -->
        <div v-if="form.family_info.guardian_is_parent === 'father'" class="text-sm text-gray-600 bg-gray-50 p-3 rounded-lg">
          <p><strong>ผู้ปกครอง:</strong> {{ form.family_info.father_name || '-' }}</p>
          <p><strong>อาชีพ:</strong> {{ form.family_info.father_occupation || '-' }}</p>
          <p><strong>รายได้:</strong> {{ (form.family_info.father_income || 0).toLocaleString() }} บาท/เดือน</p>
        </div>
        <div v-if="form.family_info.guardian_is_parent === 'mother'" class="text-sm text-gray-600 bg-gray-50 p-3 rounded-lg">
          <p><strong>ผู้ปกครอง:</strong> {{ form.family_info.mother_name || '-' }}</p>
          <p><strong>อาชีพ:</strong> {{ form.family_info.mother_occupation || '-' }}</p>
          <p><strong>รายได้:</strong> {{ (form.family_info.mother_income || 0).toLocaleString() }} บาท/เดือน</p>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.animate-fade-in {
  animation: fadeIn 0.3s ease-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(5px);
  }

  to {
    opacity: 1;
    transform: translateY(0);
  }
}

input:disabled,
textarea:disabled {
  background-color: #f9fafb;
  border-color: #e5e7eb;
  color: #374151;
  cursor: not-allowed;
}
</style>