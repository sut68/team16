<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue';
import { updateScreeningStatus } from '@/services/api/screening'; 
import { Get } from '@/services/api/https'; 

const props = defineProps<{
    isOpen: boolean;
    documentData: any;
}>();

const emit = defineEmits(['close', 'action-completed']);

// --- Local State ---
const comment = ref('');
const actionType = ref<'approve' | 'reject' | null>(null);
const isSubmitting = ref(false);
const currentUser = ref('');

onMounted(async () => {
    try {
        const res = await Get('/profile/me');
        if (res && res.data) {
             const data = res.data;
             if (data.admin_firstname) {
                 currentUser.value = data.admin_firstname;
             }
        }
    } catch (error) {
        console.error('Error fetching current user:', error);
    }
});

// --- Helper Functions ---
const getRootData = () => {
    if (props.documentData && props.documentData.data) {
        return props.documentData.data;
    }
    return props.documentData || {};
};

const parseGoDate = (dateString: string | undefined): Date | null => {
    if (!dateString || dateString === '0001-01-01T00:00:00Z') return null;
    let dateStr = String(dateString);
    if (dateStr.includes(' m=')) {
        dateStr = dateStr.split(' m=')[0] as string;
    }
    if (dateStr.includes('+0000 UTC')) {
        dateStr = dateStr.replace(' +0000 UTC', 'Z').replace(' ', 'T');
    }
    dateStr = dateStr.replace(/(\.\d{3})\d+/, '$1');
    const date = new Date(dateStr);
    return isNaN(date.getTime()) ? null : date;
};

const formatDate = (dateString: string | undefined | number) => {
    if (!dateString) return 'เมื่อสักครู่';

    if (typeof dateString === 'number') {
        return new Date(dateString).toLocaleDateString('th-TH', {
            year: 'numeric', month: 'short', day: 'numeric',
            hour: '2-digit', minute: '2-digit',
        });
    }

    const date = parseGoDate(dateString);
    if (!date) return 'เมื่อสักครู่'; 
    
    return date.toLocaleDateString('th-TH', {
        year: 'numeric', month: 'short', day: 'numeric',
        hour: '2-digit', minute: '2-digit',
    });
};

const formatNumber = (num: number) => num.toLocaleString('th-TH', { maximumFractionDigits: 2 });
const mapOperatorToText = (op: string) => {
    switch (op) {
        case '>=': return 'ไม่ต่ำกว่า';
        case '<=': return 'ไม่เกิน';
        case '>': return 'มากกว่า';
        case '<': return 'น้อยกว่า';
        case '=':
        case '==': return '';
        default: return op;
    }
};

// --- Computed: Header Info ---
const headerInfo = computed(() => {
    const root = getRootData();
    // ดึงข้อมูลจาก application_scholarship ตาม API response
    const appSch = root.application_scholarship || root.ApplicationScholarship || {};
    const scholarship = appSch.scholarship || appSch.Scholarship || {};
    const application = appSch.application || appSch.Application || {};
    const student = application.student_profile || application.StudentProfile || {};
    
    const statusId = root.StatusScreeningID || root.status_screening_id;
    let statusStr = root.status || 'pending';
    if(statusId === 2) statusStr = 'approved';
    if(statusId === 3) statusStr = 'rejected';

    // ดึงข้อมูล round จาก semaster
    const semaster = scholarship.semaster || scholarship.Semaster || {};
    const roundText = semaster.term && semaster.academic_year 
        ? `${semaster.term}/${semaster.academic_year}` 
        : (semaster.round || '1/2568');

    return {
        id: root.ID || '-',
        title: scholarship.scholarship_name || scholarship.ScholarshipName || 'ไม่ระบุชื่อทุน',
        applicant: `${student.first_name_th || student.FirstNameTH || ''} ${student.last_name_th || student.LastNameTH || ''}`.trim() || 'ไม่ระบุชื่อ',
        status: statusStr,
        round: roundText
    };
});

// --- Computed: Timeline Logic ---
const timelineEvents = computed(() => {
    const root = getRootData();
    const events: any[] = []; // Explicit type as any[] or define interface
    
    // 1. Created Event
    const createdDate = parseGoDate(root.CreatedAt);
    const createdTs = createdDate ? createdDate.getTime() : Date.now() - 100000;

    events.push({
        id: 'created',
        title: 'ยื่นใบสมัครแล้ว',
        date: formatDate(root.CreatedAt),
        description: 'ส่งเอกสารเข้าสู่ระบบเรียบร้อยแล้ว',
        actor: 'ผู้สมัคร',
        status: 'past-submitted',
        timestamp: createdTs
    });

    // 2. Determine Current State
    const status = headerInfo.value.status?.toLowerCase();
    const nowTs = Date.now();
    const updateDate = parseGoDate(root.UpdatedAt);
    const updateTs = updateDate ? updateDate.getTime() : nowTs;
    const adminName = root.admin_profile?.admin_firstname || root.AdminProfile?.AdminFirstname || 'เจ้าหน้าที่';

    if (status === 'pending' || status.includes('รอ')) {
        events.push({
            id: 'reviewing',
            title: 'เจ้าหน้าที่กำลังตรวจสอบ',
            date: formatDate(root.UpdatedAt || root.CreatedAt),
            description: 'อยู่ในระหว่างการพิจารณาคุณสมบัติ',
            actor: currentUser.value || adminName,
            status: 'current',
            timestamp: updateTs > createdTs ? updateTs : createdTs + 1000
        });
    } else if (status === 'approved' || status?.includes('ผ่าน')) {
        events.push({
            id: 'auto-approved',
            title: 'อนุมัติแล้ว',
            date: formatDate(root.UpdatedAt),
            description: 'ผ่านเกณฑ์การคัดกรองเบื้องต้น',
            actor: adminName,
            status: 'past-approved',
            timestamp: updateTs
        });
    } else if (status === 'rejected' || status?.includes('ไม่ผ่าน')) {
        events.push({
            id: 'auto-rejected',
            title: 'ปฏิเสธคำขอ',
            date: formatDate(root.UpdatedAt),
            description: root.rejection_reason || root.RejectionReason || 'ไม่ผ่านเกณฑ์ที่กำหนด',
            actor: adminName,
            status: 'past-rejected',
            timestamp: updateTs
        });
    }

    // Sort
    events.sort((a, b) => b.timestamp - a.timestamp);
    
    // Fix: Handle possible undefined array elements
    if (events.length > 0) {
        const first = events[0];
        // ตรวจสอบ first ก่อนใช้งาน
        if (first && first.status !== 'past-approved' && first.status !== 'past-rejected') {
             if (first.status.startsWith('past')) first.status = 'current';
        }
        
        // Loop เริ่มที่ 1
        for(let i=1; i<events.length; i++) {
             const currentEvent = events[i];
             // ตรวจสอบ currentEvent ก่อนใช้งาน
             if (currentEvent && currentEvent.status === 'current') {
                currentEvent.status = 'past';
             }
        }
    }

    return events;
});

// --- Computed: Criteria Logic ---
const screeningCriteria = computed(() => {
    const root = getRootData();
    // ดึงข้อมูลจาก application_scholarship.scholarship ตาม API response
    const appSch = root.application_scholarship || root.ApplicationScholarship || {};
    const scholarship = appSch.scholarship || appSch.Scholarship || root.scholarship || root.Scholarship || {};
    const rawFeatures = scholarship.featurescholarships || scholarship.FeatureScholarships || scholarship.feature_scholarships || []; 

    if (!Array.isArray(rawFeatures) || rawFeatures.length === 0) return [];

    // ดึงข้อมูล student จาก application_scholarship.application.student_profile
    const application = appSch.application || appSch.Application || root.application || {};
    const student = application.student_profile || application.StudentProfile || {};
    const family = student.family_info || student.FamilyInfo || {};

    return rawFeatures.map((item: any, index: number) => {
        const typeFeature = item.Typefeature || item.TypeFeature || {}; 
        const label = item.feature_scholarship_name || item.FeatureScholarshipName || typeFeature.type_feature_name || `เกณฑ์ข้อที่ ${index + 1}`;
        const fullTextToCheck = `${label} ${typeFeature.type_feature_name || ''}`.toLowerCase();
        const operator = item.operator || item.Operator || '>=';        
        const requiredValue = parseFloat(item.value || item.Value || '0');
        
        let studentValueNum = 0;
        let studentValueStr = '-';
        let isPassed = false;
        let unit = '';

        if (fullTextToCheck.includes('เกรด') || fullTextToCheck.includes('gpax')) {
            studentValueNum = parseFloat(student.gpax || '0');
            studentValueStr = studentValueNum.toFixed(2);
        } else if (fullTextToCheck.includes('รายได้') || fullTextToCheck.includes('income')) {
            const fatherInc = Number(family.father_income || 0);
            const motherInc = Number(family.mother_income || 0);
            const guardianInc = Number(family.guardian_income || 0);
            const guardianIsParent = family.guardian_is_parent || '';
            
            // Only count guardian income if guardian is "other" (not father or mother)
            // This prevents counting the same income twice
            let totalFamilyIncome = fatherInc + motherInc;
            if (guardianIsParent === 'other' || guardianIsParent === '') {
              // Only add guardian income if it's a different person
              // Check if guardian name is different from father/mother name
              const guardianName = family.guardian_name || '';
              const fatherName = family.father_name || '';
              const motherName = family.mother_name || '';
              if (guardianName && guardianName !== fatherName && guardianName !== motherName) {
                totalFamilyIncome += guardianInc;
              }
            }

            if (fullTextToCheck.includes('ต่อคน') || fullTextToCheck.includes('เฉลี่ย') || fullTextToCheck.includes('สมาชิก')) {
                let parentCount = 0;
                if (fatherInc > 0 || family.father_name) parentCount++;
                if (motherInc > 0 || family.mother_name) parentCount++;
                // Only count guardian as additional person if they are "other"
                if (guardianIsParent === 'other' && (guardianInc > 0 || family.guardian_name)) {
                  const guardianName = family.guardian_name || '';
                  const fatherName = family.father_name || '';
                  const motherName = family.mother_name || '';
                  if (guardianName && guardianName !== fatherName && guardianName !== motherName) {
                    parentCount++;
                  }
                }
                if (parentCount === 0) parentCount = 1; 

                const totalMembers = 1 + Number(student.siblings_count || 0) + parentCount;
                studentValueNum = totalFamilyIncome / (totalMembers > 0 ? totalMembers : 1);
                studentValueStr = formatNumber(studentValueNum) + ' บ./คน';
                unit = 'บาท/คน';
            } else {
                studentValueNum = totalFamilyIncome;
                studentValueStr = formatNumber(studentValueNum) + ' บ.';
                unit = 'บาท';
            }
        } else if (fullTextToCheck.includes('ระยะเวลา') || fullTextToCheck.includes('duration') || fullTextToCheck.includes('ชั้นปี')) {
            studentValueNum = parseInt(student.current_year || student.CurrentYear || '0');
            studentValueStr = `ชั้นปีที่ ${studentValueNum}`;
            unit = 'ปี';
        } else if (fullTextToCheck.includes('พี่น้อง')) {
            studentValueNum = parseInt(student.siblings_count || '0');
            studentValueStr = `${studentValueNum} คน`;
            unit = 'คน';
        }
        
        switch (operator) {
            case '>=': isPassed = studentValueNum >= requiredValue; break;
            case '<=': isPassed = studentValueNum <= requiredValue; break;
            case '>':  isPassed = studentValueNum > requiredValue; break;
            case '<':  isPassed = studentValueNum < requiredValue; break;
            case '=': case '==': isPassed = studentValueNum == requiredValue; break;
            default:   isPassed = studentValueNum >= requiredValue;
        }

        const requirementText = `${mapOperatorToText(operator)} ${formatNumber(requiredValue)} ${unit}`;

        return { id: item.ID || index, label, requirement: requirementText, studentValue: studentValueStr, isPassed };
    });
});

const passCount = computed(() => screeningCriteria.value.filter(c => c.isPassed).length);
const totalCriteria = computed(() => screeningCriteria.value.length);
const isAllPassed = computed(() => totalCriteria.value === 0 || passCount.value === totalCriteria.value);

const canTakeAction = computed(() => {
  const status = headerInfo.value.status?.toLowerCase() || '';
  return status === 'pending' || status.includes('รอ');
});

// --- Watch & Actions ---
watch(() => props.isOpen, (val) => {
    if (val) {
        comment.value = '';
        actionType.value = null;
        isSubmitting.value = false;
    }
});

const closeModal = () => {
    if (isSubmitting.value) return;
    emit('close');
};

const submitAction = async (type: 'approve' | 'reject') => {
    const root = getRootData();
    if (!root.ID) return;

    let statusId = type === 'approve' ? 2 : 3;

    isSubmitting.value = true;
    try {
        await updateScreeningStatus(root.ID, {
            status_screening_id: statusId,
            rejection_reason: type === 'reject' ? comment.value.trim() : null
        });

        emit('action-completed');
        closeModal();

    } catch (error) {
        console.error(error);
        alert('เกิดข้อผิดพลาดในการบันทึก');
    } finally {
        isSubmitting.value = false;
    }
};
</script>

<template>
    <div v-if="isOpen" class="fixed inset-0 z-[200] flex items-center justify-center bg-black/60 backdrop-blur-sm p-4 transition-all">
        
        <div class="bg-white w-full max-w-6xl h-[90vh] rounded-2xl shadow-2xl flex flex-col overflow-hidden animate-pop-in">
            
            <div class="px-6 py-4 border-b flex items-center justify-between bg-slate-50 sticky top-0 z-10">
                <div>
                    <h2 class="text-xl font-bold text-[#1e3a8a] flex items-center gap-2">
                        {{ headerInfo.title }}
                        <span class="badge text-white border-none"
                            :class="{
                                'badge-warning': headerInfo.status.includes('รอ') || headerInfo.status === 'pending',
                                'badge-success': headerInfo.status.includes('ผ่าน') || headerInfo.status === 'approved',
                                'badge-error': headerInfo.status.includes('ไม่ผ่าน') || headerInfo.status === 'rejected'
                            }">
                            {{ headerInfo.status === 'pending' ? 'รอตรวจสอบ' : (headerInfo.status === 'approved' ? 'ผ่านการคัดกรอง' : 'ไม่ผ่าน') }}
                        </span>
                    </h2>
                    <div class="flex items-center gap-2 text-sm text-gray-500 mt-1">
                        <span>Task ID: #{{ headerInfo.id }}</span>
                        <span>•</span>
                        <span class="bg-blue-50 text-blue-700 px-2 rounded border border-blue-100 text-xs">
                            รอบ: {{ headerInfo.round }}
                        </span>
                        <span>•</span>
                        <span class="font-medium text-slate-700">{{ headerInfo.applicant }}</span>
                    </div>
                </div>
                <button @click="closeModal" class="btn btn-circle btn-ghost btn-sm hover:bg-slate-200">✕</button>
            </div>
    
            <div class="flex-1 overflow-y-auto p-6 bg-[#f8fafc]">
                <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
                    
                    <div class="lg:col-span-2 flex flex-col gap-6">
                        <div class="card bg-white shadow-sm border border-gray-200 overflow-hidden">
                            <div class="px-5 py-4 border-b border-gray-100 bg-white flex justify-between items-center">
                                <h3 class="font-bold text-lg text-slate-700">ตรวจสอบคุณสมบัติ (Automated Check)</h3>
                                <div class="text-sm font-medium px-3 py-1 rounded-full border" 
                                    :class="isAllPassed ? 'bg-green-50 text-green-700 border-green-200' : 'bg-orange-50 text-orange-700 border-orange-200'">
                                    ผ่านเกณฑ์ {{ passCount }}/{{ totalCriteria }} ข้อ
                                </div>
                            </div>
    
                            <div class="overflow-x-auto">
                                <table class="w-full text-sm text-left">
                                    <thead class="text-xs text-gray-500 uppercase bg-gray-50 border-b">
                                        <tr>
                                            <th class="px-6 py-3">เกณฑ์</th>
                                            <th class="px-6 py-3">เงื่อนไขที่กำหนด</th>
                                            <th class="px-6 py-3">ข้อมูลผู้สมัคร</th>
                                            <th class="px-6 py-3 text-center">ผลลัพธ์</th>
                                        </tr>
                                    </thead>
                                    <tbody class="divide-y divide-gray-100">
                                        <tr v-if="screeningCriteria.length === 0">
                                            <td colspan="4" class="px-6 py-8 text-center text-gray-400">ไม่พบข้อมูลเกณฑ์ หรือ ข้อมูลยังไม่โหลด</td>
                                        </tr>
                                        <tr v-for="item in screeningCriteria" :key="item.id" class="hover:bg-slate-50">
                                            <td class="px-6 py-4 font-medium text-slate-700">{{ item.label }}</td>
                                            <td class="px-6 py-4 text-gray-500 bg-gray-50/50 font-mono text-xs">{{ item.requirement }}</td>
                                            <td class="px-6 py-4" :class="{'text-red-600 font-bold': !item.isPassed}">{{ item.studentValue }}</td>
                                            <td class="px-6 py-4 text-center">
                                                <span v-if="item.isPassed" class="badge badge-success text-white badge-sm">ผ่าน</span>
                                                <span v-else class="badge badge-error text-white badge-sm">ไม่ผ่าน</span>
                                            </td>
                                        </tr>
                                    </tbody>
                                </table>
                            </div>
                        </div>
                    </div>
    
                    <div class="lg:col-span-1 space-y-6">
                        
                        <div v-if="canTakeAction" class="card bg-white shadow-md border border-blue-100 ring-4 ring-blue-50/50">
                            <div class="card-body p-5">
                                <h3 class="font-bold text-lg text-slate-700 mb-4">ผลการพิจารณา</h3>

                                <div class="flex flex-col gap-3">
                                    <button @click="submitAction('approve')" 
                                        :disabled="isSubmitting || !isAllPassed"
                                        class="btn bg-[#1e3a8a] hover:bg-[#152c6f] text-white w-full border-none shadow-sm disabled:bg-slate-200 disabled:text-slate-400">
                                        <span v-if="isSubmitting" class="loading loading-spinner loading-xs"></span>
                                        <span v-else-if="!isAllPassed">คุณสมบัติไม่ครบถ้วน</span>
                                        <span v-else>ผ่านการคัดกรอง</span>
                                    </button>
                                    
                                    <button @click="actionType = 'reject'" :disabled="isSubmitting"
                                        class="btn btn-outline btn-error w-full hover:text-white">
                                        ไม่ผ่านการคัดกรอง
                                    </button>
                                </div>

                                <div v-if="actionType === 'reject'" class="mt-4 pt-4 border-t animate-fade-in">
                                    <p class="text-sm font-bold mb-2 text-error flex items-center gap-1">
                                        ระบุเหตุผลที่ปฏิเสธ:
                                    </p>
                                    <textarea v-model="comment" 
                                        class="textarea textarea-bordered w-full h-24 text-sm border-error focus:ring-error"
                                        placeholder="เช่น เกรดเฉลี่ยไม่ถึงเกณฑ์, รายได้เกินกำหนด..."></textarea>
                                    <div class="flex justify-end gap-2 mt-3">
                                        <button @click="actionType = null" class="btn btn-ghost btn-xs text-gray-500">ยกเลิก</button>
                                        <button @click="submitAction('reject')" 
                                            class="btn btn-sm btn-error text-white border-none">
                                            ยืนยันผล
                                        </button>
                                    </div>
                                </div>
                            </div>
                        </div>

                        <div v-else class="card bg-white shadow-sm border border-gray-100">
                             <div class="card-body p-5 items-center text-center">
                                <div class="badge badge-lg p-4 font-bold text-white mb-2"
                                    :class="headerInfo.status.includes('ผ่าน') || headerInfo.status === 'approved' ? 'badge-success' : 'badge-error'">
                                    {{ headerInfo.status === 'approved' ? 'ผ่านการคัดกรองแล้ว' : 'ไม่ผ่านการคัดกรอง' }}
                                </div>
                                <p class="text-xs text-gray-400">การดำเนินการเสร็จสิ้น</p>
                             </div>
                        </div>

                        <div class="card bg-white shadow-sm border border-gray-100">
                            <div class="card-body p-5">
                                <h3 class="font-bold text-lg text-slate-700 mb-4 border-b pb-2">ประวัติการดำเนินการ</h3>
                                <ul class="timeline timeline-vertical timeline-compact -ml-4">
                                    <li v-for="(event, index) in timelineEvents" :key="event.id">
                                        
                                        <hr v-if="index > 0" class="bg-gray-200" />
                                        
                                        <div class="timeline-middle">
                                            <div v-if="event.status === 'current'" class="relative flex items-center justify-center w-6 h-6">
                                                <span class="absolute inline-flex h-full w-full rounded-full bg-blue-400 opacity-75 animate-ping"></span>
                                                <span class="relative inline-flex rounded-full h-4 w-4 bg-[#1e3a8a]"></span>
                                            </div>
                                            <div v-else-if="event.status.includes('approved') || event.status.includes('submitted')" 
                                                 class="w-6 h-6 rounded-full bg-green-500 text-white flex items-center justify-center scale-75">
                                                 <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="w-4 h-4">
                                                    <path fill-rule="evenodd" d="M16.704 4.153a.75.75 0 01.143 1.052l-8 10.5a.75.75 0 01-1.127.075l-4.5-4.5a.75.75 0 011.06-1.06l3.894 3.893 7.48-9.817a.75.75 0 011.05-.143z" clip-rule="evenodd" />
                                                  </svg>
                                            </div>
                                            <div v-else-if="event.status.includes('rejected')" 
                                                 class="w-6 h-6 rounded-full bg-red-500 text-white flex items-center justify-center scale-75">
                                                 <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="w-4 h-4">
                                                    <path d="M6.28 5.22a.75.75 0 00-1.06 1.06L8.94 10l-3.72 3.72a.75.75 0 101.06 1.06L10 11.06l3.72 3.72a.75.75 0 101.06-1.06L11.06 10l3.72-3.72a.75.75 0 00-1.06-1.06L10 8.94 6.28 5.22z" />
                                                  </svg>
                                            </div>
                                            <div v-else class="w-6 h-6 rounded-full bg-gray-300 scale-75"></div>
                                        </div>
                                        
                                        <div class="timeline-end timeline-box w-full border-none shadow-none p-0 pl-2 mb-6">
                                            <div class="font-bold text-slate-800 text-sm" :class="{'text-[#1e3a8a]': event.status === 'current'}">
                                                {{ event.title }}
                                            </div>
                                            <div class="text-xs text-gray-500 mb-2">
                                                {{ event.date }} • โดย {{ event.actor }}
                                            </div>
                    
                                            <div v-if="event.status.includes('rejected')" 
                                                 class="bg-red-50 border border-red-100 rounded-xl p-3 text-sm text-red-800 w-full break-words">
                                                {{ event.description }}
                                            </div>
                                            
                                            <div v-else class="text-xs text-gray-600 break-words">
                                                {{ event.description }}
                                            </div>
                                        </div>

                                        <hr v-if="index < timelineEvents.length - 1" class="bg-gray-200" />
                                    </li>
                                </ul>
                            </div>
                        </div>

                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
@keyframes pop-in {
    0% { opacity: 0; transform: scale(0.95) translateY(10px); }
    100% { opacity: 1; transform: scale(1) translateY(0); }
}
.animate-pop-in { animation: pop-in 0.2s cubic-bezier(0.16, 1, 0.3, 1) forwards; }
.animate-fade-in { animation: fade 0.3s ease forwards; }
@keyframes fade {
    from { opacity: 0; height: 0; transform: translateY(-5px); }
    to { opacity: 1; height: auto; transform: translateY(0); }
}
</style>