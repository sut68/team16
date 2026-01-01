<script setup lang="ts">
import { ref, computed, reactive, watch, onMounted, nextTick } from 'vue';
import { InterviewAPI, type InterviewRoundCreate, type InterviewRoundUpdate } from '@/services/api'; 
import { ScholarshipAPI } from '@/services/api/scholarship'; 
import { LocationAPI } from '@/services/api/location';
import type { InterviewRound, Interviewer, Location, InterviewMode, Slot } from '@/interfaces/interview';
import type { ScholarshipResponse } from '@/interfaces/scholarship';
import Swal from 'sweetalert2';
import { Plus } from 'lucide-vue-next';
// --- Interfaces ---
interface StudentProfile {
    ID: number;
    first_name_th: string;
    last_name_th: string;
    student_id: string;
    email: string;
    gpax: number;
    major?: {
        major_name: string;
    } | null;
}

// --- 1. State Management ---
const isLoading = ref(true);
const isModalOpen = ref(false);
const isFilterOpen = ref(false);
const isStudentDetailModalOpen = ref(false);
const modalMode = ref<'create' | 'view' | 'edit'>('create');
const activeTab = ref<'active' | 'history'>('active');
const editingRoundId = ref<number | null>(null);

// Data from API
const allRounds = ref<InterviewRound[]>([]);
const scholarships = ref<ScholarshipResponse[]>([]);
const interviewersList = ref<Interviewer[]>([]);
const locations = ref<Location[]>([]);
const interviewModes = ref<InterviewMode[]>([]);
const selectedRoundDetails = ref<InterviewRound | null>(null);
const selectedStudent = ref<StudentProfile | null>(null);


// Filter & Sort State
const searchQuery = ref('');
const filterStatus = ref('All');
const filterScholarship = ref('All');
const sortBy = ref('DateAsc');

// Form Data
const formData = reactive({
    name: '',
    scholarship_id: null as number | null,
    description: '',
    date: '',
    end_date: '',
    start_time: '09:00',
    end_time: '16:00',
    slot_duration: 30,
    interviewer_ids: [] as number[],
    interview_mode_id: null as number | null,
    location_id: null as number | null,
    meeting_link: '',
});

// เก็บค่าตั้งต้นตอนโหลดโหมด Edit เพื่อเช็คว่ามีการเปลี่ยนเวลา/วันที่ไหม
const initialFormDataForEdit = reactive({
    date: '',
    end_date: '',
    start_time: '',
    end_time: '',
    slot_duration: 30,
});

// ตรวจสอบว่ามีการแก้ไข "โครงสร้างเวลา" หรือไม่ (ถ้าแก้จะไปโชว์ Preview แทน Slots เดิม)
const isSlotConfigDirty = computed(() => {
    if (modalMode.value !== 'edit') return false;
    return (
        formData.date !== initialFormDataForEdit.date ||
        formData.end_date !== initialFormDataForEdit.end_date ||
        formData.start_time !== initialFormDataForEdit.start_time ||
        formData.end_time !== initialFormDataForEdit.end_time ||
        formData.slot_duration !== initialFormDataForEdit.slot_duration
    );
});

// --- 2. Data Fetching ---
onMounted(async () => {
    await fetchData();
});

const fetchData = async () => {
    isLoading.value = true;
    try {
        const [roundsRes, scholarsRes, interviewersRes, locationsRes, modesRes] = await Promise.allSettled([
            InterviewAPI.getAllRounds(),
            ScholarshipAPI.getAll(),
            InterviewAPI.getAllInterviewers(),
            LocationAPI.getAllLocations(),
            InterviewAPI.getAllModes(),
        ]);

        if (roundsRes.status === 'fulfilled') allRounds.value = roundsRes.value || [];
        if (scholarsRes.status === 'fulfilled') scholarships.value = scholarsRes.value || [];
        if (interviewersRes.status === 'fulfilled') interviewersList.value = interviewersRes.value || [];
        if (locationsRes.status === 'fulfilled') locations.value = locationsRes.value || [];
        if (modesRes.status === 'fulfilled') interviewModes.value = modesRes.value || [];

    } catch (error) {
        console.error(error);
        Swal.fire('เกิดข้อผิดพลาด', 'ไม่สามารถดึงข้อมูลได้', 'error');
    } finally {
        isLoading.value = false;
    }
};

// --- 3. Computed Logic ---
const getModeName = (round: InterviewRound) => {
    if ((round as any).interview_mode?.name) return (round as any).interview_mode.name;
    const mode = interviewModes.value.find(m => m.ID === round.interview_mode_id);
    return mode ? mode.name : 'Unknown';
};

const getLocationLabel = (round: InterviewRound) => {
    const modeName = getModeName(round);
    if (modeName === 'Onsite') {
        if ((round as any).location?.name) return (round as any).location.name;
        const loc = locations.value.find(l => l.ID === round.location_id);
        return loc ? loc.name : '-';
    }
    return 'Online Meeting';
};

const getRoundStatus = (round: InterviewRound): 'Closed' | 'Full' | 'Open' => {
    const roundEnd = new Date(round.end_date_time);
    const today = new Date();
    today.setHours(0, 0, 0, 0);

    if (roundEnd.getTime() < today.getTime()) return 'Closed';
    if (!round.slots || round.slots.length === 0) return 'Open';
    const isFull = round.slots.every(slot => slot.is_booked);
    return isFull ? 'Full' : 'Open';
}

const getBookedCount = (round: InterviewRound) => {
    if (!round.slots) return 0;
    return round.slots.filter(s => s.is_booked).length;
}

const stats = computed(() => {
    if (activeTab.value === 'active') {
        const activeRounds = allRounds.value.filter(r => getRoundStatus(r) !== 'Closed');
        const totalRounds = activeRounds.length;
        const totalCapacity = activeRounds.reduce((acc, r) => acc + (r.slots?.length || 0), 0);
        const totalBooked = activeRounds.reduce((acc, r) => acc + getBookedCount(r), 0);

        return {
            title1: 'รอบที่เปิดอยู่',
            value1: totalRounds,
            desc1: 'กำลังดำเนินการ',
            
            title2: 'นศ. ที่จองแล้ว',
            value2: totalBooked,
            desc2: 'คนของที่นั่งทั้งหมด',
            
            title3: 'ที่นั่งว่างเหลือ',
            value3: totalCapacity - totalBooked,
            desc3: 'ที่นั่ง'
        };
    } else { // history tab
        const closedRounds = allRounds.value.filter(r => getRoundStatus(r) === 'Closed');
        const totalRounds = closedRounds.length;
        const totalBooked = closedRounds.reduce((acc, r) => acc + getBookedCount(r), 0);
        const totalCapacity = closedRounds.reduce((acc, r) => acc + (r.slots?.length || 0), 0);
        
        return {
            title1: 'รอบที่เสร็จสิ้น',
            value1: totalRounds,
            desc1: 'ในประวัติทั้งหมด',
            
            title2: 'ผู้เข้าร่วมทั้งหมด',
            value2: totalBooked,
            desc2: `จาก ${totalCapacity} ที่นั่ง`,
            
            title3: 'ทุนที่เกี่ยวข้อง',
            value3: new Set(closedRounds.map(r => r.scholarship_id)).size,
            desc3: 'โครงการ'
        };
    }
});

const filteredRounds = computed(() => {
    let result = allRounds.value;

    if (activeTab.value === 'active') {
        result = result.filter(r => getRoundStatus(r) !== 'Closed');
    } else {
        result = result.filter(r => getRoundStatus(r) === 'Closed');
    }

    if (searchQuery.value) {
        const q = searchQuery.value.toLowerCase();
        result = result.filter(r =>
            r.name.toLowerCase().includes(q) ||
            (r.scholarship && r.scholarship.scholarship_name.toLowerCase().includes(q))
        );
    }

    if (filterStatus.value !== 'All') {
        result = result.filter(r => getRoundStatus(r) === filterStatus.value);
    }
    if (filterScholarship.value !== 'All') {
        result = result.filter(r => r.scholarship?.typescholarship?.type_name === filterScholarship.value);
    }

    result = [...result].sort((a, b) => {
        if (sortBy.value === 'DateAsc') return new Date(a.start_date_time).getTime() - new Date(b.start_date_time).getTime();
        if (sortBy.value === 'DateDesc') return new Date(b.start_date_time).getTime() - new Date(a.start_date_time).getTime();
        if (sortBy.value === 'BookedDesc') return getBookedCount(b) - getBookedCount(a);
        return 0;
    });

    return result;
});

const previewSlots = computed(() => {
    const slots: { time: string }[] = [];
    if (!formData.date || !formData.end_date || !formData.start_time || !formData.end_time || formData.slot_duration <= 0) return slots;

    const startDateTime = new Date(`${formData.date}T${formData.start_time}`);
    const endDateTime = new Date(`${formData.end_date}T${formData.end_time}`);

    let current = startDateTime.getTime();
    const end = endDateTime.getTime();

    while (current < end) {
        const d = new Date(current);
        const timeStr = `${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`;
        slots.push({ time: timeStr });
        current += formData.slot_duration * 60 * 1000;
    }
    return slots;
});

const getInitials = (interviewer: Interviewer) => {
    if (!interviewer) return '';
    return `${interviewer.interviewer_firstname?.charAt(0) ?? ''}${interviewer.interviewer_lastname?.charAt(0) ?? ''}`;
}

const getInterviewersForRound = (round: InterviewRound): Interviewer[] => {
    if (!round.slots || round.slots.length === 0) return [];
    const interviewers = new Map<number, Interviewer>();
    round.slots.forEach(slot => {
        slot.interviewer_slots?.forEach(is => {
            if (is.interviewer && !interviewers.has(is.interviewer.ID)) {
                interviewers.set(is.interviewer.ID, is.interviewer);
            }
        });
    });
    return Array.from(interviewers.values());
}

const selectedModeName = computed(() => {
    const selectedMode = interviewModes.value.find(m => m.ID === formData.interview_mode_id);
    return selectedMode ? selectedMode.name : '';
});

// --- 4. Watchers ---
watch(activeTab, () => {
    searchQuery.value = '';
    filterStatus.value = 'All';
    filterScholarship.value = 'All';
    isFilterOpen.value = false;
    sortBy.value = activeTab.value === 'active' ? 'DateAsc' : 'DateDesc';
});

// --- 5. Methods ---

const openStudentDetailModal = (slot: Slot) => {
    if (!slot.is_booked) return;

    const booking = slot.interviewe_bookings?.[0];

    if (booking && booking.application_scholarship?.application?.student_profile) {
        selectedStudent.value = booking.application_scholarship.application.student_profile;
        isStudentDetailModalOpen.value = true;
    } else {
        Swal.fire('ไม่พบข้อมูล', 'ไม่พบข้อมูลการจองหรือข้อมูลนักศึกษาสำหรับช่องเวลานี้', 'info');
    }
};

const openCreateModal = () => {
    modalMode.value = 'create';
    editingRoundId.value = null;
    selectedRoundDetails.value = null;
    Object.assign(formData, {
        name: '', scholarship_id: null, description: '',
        date: '', end_date: '', start_time: '09:00', end_time: '16:00',
        slot_duration: 30, interviewer_ids: [], interview_mode_id: null,
        location_id: null, meeting_link: ''
    });
    isModalOpen.value = true;
};

const populateFormWithRoundDetails = (details: InterviewRound) => {
    const start = new Date(details.start_date_time);
    const end = new Date(details.end_date_time);

    Object.assign(formData, {
        name: details.name,
        scholarship_id: details.scholarship_id,
        description: details.description,
        date: start.toISOString().split('T')[0],
        end_date: end.toISOString().split('T')[0],
        start_time: start.toLocaleTimeString('en-GB', { hour: '2-digit', minute: '2-digit' }),
        end_time: end.toLocaleTimeString('en-GB', { hour: '2-digit', minute: '2-digit' }),
        slot_duration: details.slot_duration,
        interviewer_ids: getInterviewersForRound(details).map(i => i.ID),
        interview_mode_id: details.interview_mode_id,
        location_id: details.location_id,
        meeting_link: details.meeting_link,
    });
};

const openViewModal = async (round: InterviewRound) => {
    modalMode.value = 'view';
    isModalOpen.value = true;
    isLoading.value = true;
    selectedRoundDetails.value = null; // ล้างค่าเก่า

    try {
        const details = await InterviewAPI.getRoundById(round.ID);
        selectedRoundDetails.value = details;
        populateFormWithRoundDetails(details);
    } catch (e) {
        Swal.fire('เกิดข้อผิดพลาด', 'ไม่สามารถดึงรายละเอียดรอบได้', 'error');
        isModalOpen.value = false;
    } finally {
        isLoading.value = false;
    }
};

const openEditModal = async (round: InterviewRound) => {
    modalMode.value = 'edit';
    editingRoundId.value = round.ID;
    isModalOpen.value = true;
    isLoading.value = true;
    selectedRoundDetails.value = null; // ล้างค่าเก่า

    try {
        const details = await InterviewAPI.getRoundById(round.ID);
        selectedRoundDetails.value = details;

        // 1. ใส่ข้อมูลลง Form
        populateFormWithRoundDetails(details);
        
        // 2. รอให้ Vue อัปเดต state ของ formData แป๊บนึง แล้วค่อยเซ็ตค่าเริ่มต้นสำหรับเปรียบเทียบ
        await nextTick();
        Object.assign(initialFormDataForEdit, {
            date: formData.date,
            end_date: formData.end_date,
            start_time: formData.start_time,
            end_time: formData.end_time,
            slot_duration: formData.slot_duration,
        });

    } catch (e) {
        Swal.fire('เกิดข้อผิดพลาด', 'ไม่สามารถดึงรายละเอียดรอบได้', 'error');
        isModalOpen.value = false;
    } finally {
        isLoading.value = false;
    }
};

const handleDelete = async (id: number) => {
    const result = await Swal.fire({
        title: 'คุณแน่ใจหรือไม่?',
        text: "คุณจะไม่สามารถย้อนกลับการกระทำนี้ได้!",
        icon: 'warning',
        showCancelButton: true,
        confirmButtonColor: '#d33',
        cancelButtonColor: '#3085d6',
        confirmButtonText: 'ใช่, ลบเลย!',
        cancelButtonText: 'ยกเลิก'
    });

    if (result.isConfirmed) {
        try {
            await InterviewAPI.deleteRound(id);
            Swal.fire('ลบสำเร็จ!', 'รอบการสัมภาษณ์ถูกลบเรียบร้อยแล้ว', 'success');
            await fetchData();
        } catch (e) {
            Swal.fire('เกิดข้อผิดพลาด', 'ไม่สามารถลบรอบได้', 'error');
        }
    }
};

const handleUpdate = async () => {
    if (!editingRoundId.value) return;

    if (!formData.scholarship_id || !formData.date || !formData.start_time || !formData.end_time || !formData.end_date || !formData.interview_mode_id) {
        Swal.fire('ข้อมูลไม่ถูกต้อง', 'กรุณากรอกข้อมูลที่จำเป็นให้ครบถ้วน', 'warning');
        return;
    }

    const startDateTime = new Date(`${formData.date}T${formData.start_time}`).toISOString();
    const endDateTime = new Date(`${formData.end_date}T${formData.end_time}`).toISOString();
    
    const payload: InterviewRoundUpdate = {
        name: formData.name,
        description: formData.description,
        start_date_time: startDateTime,
        end_date_time: endDateTime,
        scholarship_id: formData.scholarship_id,
        interview_mode_id: formData.interview_mode_id,
        location_id: formData.location_id ? Number(formData.location_id) : null,
        meeting_link: formData.meeting_link,
        interviewer_ids: formData.interviewer_ids,
        slot_duration: formData.slot_duration,
    };

    isLoading.value = true;
    try {
        await InterviewAPI.updateRound(editingRoundId.value, payload);
        Swal.fire('สำเร็จ', 'อัปเดตรอบสัมภาษณ์เรียบร้อยแล้ว!', 'success');
        isModalOpen.value = false;
        await fetchData();
    } catch (error: any) {
        Swal.fire('เกิดข้อผิดพลาด', error.response?.data?.error || 'ไม่สามารถอัปเดตรอบสัมภาษณ์ได้', 'error');
    } finally {
        isLoading.value = false;
    }
};

const saveRound = async () => {
    if (!formData.scholarship_id || !formData.date || !formData.start_time || !formData.end_time || !formData.end_date || !formData.interview_mode_id) {
        Swal.fire('ข้อมูลไม่ถูกต้อง', 'กรุณากรอกข้อมูลที่จำเป็นให้ครบถ้วน', 'warning');
        return;
    }

    const startDateTime = new Date(`${formData.date}T${formData.start_time}`).toISOString();
    const endDateTime = new Date(`${formData.end_date}T${formData.end_time}`).toISOString();
    const adminProfileId = 1;

    const payload: InterviewRoundCreate = {
        name: formData.name, description: formData.description,
        start_date_time: startDateTime, end_date_time: endDateTime,
        slot_duration: formData.slot_duration, scholarship_id: formData.scholarship_id,
        admin_profile_id: adminProfileId, interviewer_ids: formData.interviewer_ids,
        interview_mode_id: formData.interview_mode_id,
        location_id: formData.location_id ? Number(formData.location_id) : null,
        meeting_link: formData.meeting_link
    };

    isLoading.value = true;
    try {
        await InterviewAPI.createRound(payload);
        Swal.fire('สำเร็จ', 'สร้างรอบสัมภาษณ์เรียบร้อยแล้ว!', 'success');
        isModalOpen.value = false;
        await fetchData();
    } catch (error: any) {
        Swal.fire('เกิดข้อผิดพลาด', error.response?.data?.error || 'ไม่สามารถสร้างรอบสัมภาษณ์ได้', 'error');
    } finally {
        isLoading.value = false;
    }
};
</script>

<template>
    <div class="w-full mx-auto flex flex-col h-full p-6 bg-white rounded-tl-[30px] shadow overflow-visible" data-theme="light">

        <div class="flex flex-col md:flex-row items-start md:items-center justify-between gap-4 mb-8">
            <h1 class="text-2xl font-bold text-slate-800">จัดการรอบสัมภาษณ์</h1>
            <button @click="openCreateModal"
                class="btn btn-sm bg-white border border-gray-300 text-gray-700 hover:bg-gray-100 flex items-center gap-2 rounded-full px-5 h-10 shadow-sm transition-all font-medium">
                <Plus class="w-4 h-4"/>
                สร้างรอบสัมภาษณ์
            </button>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-3 bg-white shadow rounded-2xl border border-gray-100 w-full mb-8 divide-y md:divide-y-0 md:divide-x divide-gray-100">
            
            <div class="p-3 flex flex-row items-center justify-between">
                <div>
                    <div class="text-slate-500 text-sm mb-1">{{ stats.title1 }}</div>
                    <div class="text-[#1e3a8a] text-3xl font-bold">{{ stats.value1 }}</div>
                    <div class="text-xs text-info mt-1">{{ stats.desc1 }}</div>
                </div>
                <div class="text-[#1e3a8a] bg-blue-50 p-3 rounded-full">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="inline-block w-8 h-8 stroke-current">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"></path>
                    </svg>
                </div>
            </div>

            <div class="p-3 flex flex-row items-center justify-between">
                <div>
                    <div class="text-slate-500 text-sm mb-1">{{ stats.title2 }}</div>
                    <div class="text-emerald-700 text-3xl font-bold">{{ stats.value2 }}</div>
                    <div class="text-success font-medium text-xs mt-1">{{ stats.desc2 }}</div>
                </div>
                <div class="text-emerald-700 bg-green-50 p-3 rounded-full">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="inline-block w-8 h-8 stroke-current">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z"></path>
                    </svg>
                </div>
            </div>

            <div class="p-3 flex flex-row items-center justify-between">
                <div>
                    <div class="text-slate-500 text-sm mb-1">{{ stats.title3 }}</div>
                    <div class="text-orange-500 text-3xl font-bold">{{ stats.value3 }}</div>
                    <div class="text-xs text-warning mt-1">{{ stats.desc3 }}</div>
                </div>
                <div class="text-orange-500 bg-orange-50 p-3 rounded-full">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="inline-block w-8 h-8 stroke-current">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                    </svg>
                </div>
            </div>
            
        </div>

        <div class="relative z-50 flex flex-col xl:flex-row items-end xl:items-center justify-between gap-4 mb-6 border-b border-gray-200">
            <div class="flex gap-8 -mb-[1px]">
                <a @click="activeTab = 'active'" class="pb-3 px-1 text-base font-medium cursor-pointer transition-all border-b-[3px]"
                    :class="activeTab === 'active' ? 'text-[#1e3a8a] border-[#1e3a8a]' : 'text-slate-500 border-transparent hover:text-slate-700'">
                    รอบปัจจุบัน
                </a>
                <a @click="activeTab = 'history'" class="pb-3 px-1 text-base font-medium cursor-pointer transition-all border-b-[3px]"
                    :class="activeTab === 'history' ? 'text-[#1e3a8a] border-[#1e3a8a]' : 'text-slate-500 border-transparent hover:text-slate-700'">
                    ประวัติย้อนหลัง
                </a>
            </div>

            <div class="flex flex-col md:flex-row items-center gap-2 pb-4 xl:pb-2">
                <button @click="isFilterOpen = !isFilterOpen"
                    class="btn btn-sm btn-outline bg-white border-gray-300 text-gray-600 hover:bg-gray-50 h-10 rounded-full px-5 shadow-sm">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-2" viewBox="0 0 20 20" fill="currentColor">
                        <path fill-rule="evenodd" d="M3 3a1 1 0 011-1h12a1 1 0 011 1v3a1 1 0 01-.293.707L12 11.414V15a1 1 0 01-.293.707l-2 2A1 1 0 018 17v-5.586L3.293 6.707A1 1 0 013 6V3z" clip-rule="evenodd" />
                    </svg>
                    ตัวกรอง
                </button>

                <div v-if="isFilterOpen" class="fixed inset-0 z-50" @click="isFilterOpen = false"></div>
                <div v-if="isFilterOpen" class="absolute right-0 mt-2 w-72 bg-white rounded-xl shadow-2xl z-[100] border p-4 animate-pop-in">
                    <p class="font-bold text-base mb-3 text-slate-700">ตัวกรองข้อมูล</p>
                    <div class="space-y-3">
                        <div>
                            <label class="text-xs font-medium text-gray-500">สถานะ</label>
                            <select v-model="filterStatus" class="select select-bordered select-sm w-full h-10 bg-white">
                                <option value="All">ทั้งหมด</option>
                                <option value="Open">เปิดให้จอง</option>
                                <option value="Full">เต็มแล้ว</option>
                            </select>
                        </div>
                        <div>
                            <label class="text-xs font-medium text-gray-500">ประเภททุน</label>
                            <select v-model="filterScholarship" class="select select-bordered select-sm w-full h-10 bg-white">
                                <option value="All">ทั้งหมด</option>
                                <option v-for="t in [...new Set(scholarships.map(s => s.typescholarship?.type_name))].filter(Boolean)" :key="t" :value="t">{{ t }}</option>
                            </select>
                        </div>
                    </div>
                </div>

                <select v-model="sortBy" class="select select-bordered select-sm rounded-full h-10 bg-white border-gray-300 w-full md:w-44 shadow-sm px-4">
                    <option value="DateAsc">วันที่ (เก่าสุดก่อน)</option>
                    <option value="DateDesc">วันที่ (ใหม่สุดก่อน)</option>
                    <option value="BookedDesc">ยอดจองมากที่สุด</option>
                </select>

                <div class="relative w-full md:w-64">
                    <span class="absolute inset-y-0 left-0 flex items-center pl-3 text-gray-400">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                        </svg>
                    </span>
                    <input type="text" v-model="searchQuery" placeholder="ค้นหาชื่อรอบ, ทุน..."
                        class="input input-bordered input-sm h-10 w-full rounded-full pl-10 bg-white border-gray-300 shadow-sm" />
                </div>
            </div>
        </div>

        <div class="bg-slate-50/30 rounded-3xl shadow-sm border border-slate-100 p-6 min-h-[500px] flex-1">
            
            <div v-if="activeTab === 'active'" class="animate-fade-in">
                <div v-if="filteredRounds.length > 0" class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6">
                    <div v-for="round in filteredRounds" :key="round.ID"
                        class="group card bg-white border border-gray-200 hover:border-blue-400 hover:shadow-xl transition-all rounded-2xl overflow-visible">
                        
                        <div class="px-5 py-4 border-b border-gray-50 bg-slate-50/50 flex justify-between items-start rounded-t-2xl">
                            <div>
                                <div class="flex gap-2 mb-2">
                                    <div class="badge badge-sm font-medium border-none text-white shadow-sm" :class="{
                                        'badge-success': getRoundStatus(round) === 'Open',
                                        'badge-error': getRoundStatus(round) === 'Full',
                                        'badge-ghost text-gray-500': getRoundStatus(round) === 'Closed'
                                    }">
                                        {{ getRoundStatus(round) === 'Open' ? 'เปิดให้จอง' : (getRoundStatus(round) === 'Full' ? 'เต็มแล้ว' : 'ปิด') }}
                                    </div>
                                    <div v-if="round.scholarship" class="badge badge-sm badge-outline text-gray-500">
                                        {{ round.scholarship.scholarship_name }}
                                    </div>
                                </div>
                                <h3 class="font-bold text-[#1e3a8a] text-lg leading-tight">{{ round.name }}</h3>
                            </div>
                            <div class="dropdown dropdown-end">
                                <label tabindex="0" class="btn btn-circle btn-ghost btn-sm text-gray-400"><svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="w-5 h-5 stroke-current"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 5v.01M12 12v.01M12 19v.01M12 6a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2z"></path></svg></label>
                                <ul tabindex="0" class="dropdown-content menu p-2 shadow-lg bg-white rounded-box w-40 z-[10] border">
                                    <li><a @click="openViewModal(round)">รายละเอียด</a></li>
                                    <li><a @click="openEditModal(round)">แก้ไข</a></li>
                                    <div class="divider my-0"></div>
                                    <li><a @click="handleDelete(round.ID)" class="text-error">ลบรายการ</a></li>
                                </ul>
                            </div>
                        </div>

                        <div class="p-5 space-y-4 cursor-pointer" @click="openViewModal(round)">
                            <div class="flex items-center gap-3 text-sm">
                                <div class="w-8 h-8 rounded-full bg-blue-50 flex items-center justify-center text-blue-600">
                                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" /></svg>
                                </div>
                                <div>
                                    <p class="font-semibold text-slate-700">{{ new Date(round.start_date_time).toLocaleDateString('th-TH', { day: 'numeric', month: 'short', year: '2-digit'}) }}</p>
                                    <p class="text-xs text-gray-500">{{ new Date(round.start_date_time).toLocaleTimeString('th-TH', {hour: '2-digit', minute: '2-digit'}) }} - {{ new Date(round.end_date_time).toLocaleTimeString('th-TH', {hour: '2-digit', minute: '2-digit'}) }}</p>
                                </div>
                            </div>

                            <div class="flex items-center gap-3 text-sm">
                                <div class="w-8 h-8 rounded-full flex items-center justify-center" :class="getModeName(round) === 'Onsite' ? 'bg-orange-50 text-orange-600' : 'bg-purple-50 text-purple-600'">
                                    <svg v-if="getModeName(round) === 'Onsite'" xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M5.05 4.05a7 7 0 119.9 9.9L10 18.9l-4.95-4.95a7 7 0 010-9.9zM10 11a2 2 0 100-4 2 2 0 000 4z" clip-rule="evenodd" /></svg>
                                    <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor"><path d="M2 6a2 2 0 012-2h6a2 2 0 012 2v8a2 2 0 01-2 2H4a2 2 0 01-2-2V6zM14.553 7.106A1 1 0 0014 8v4a1 1 0 00.553.894l2 1A1 1 0 0018 13V7a1 1 0 00-1.447-.894l-2 1z" /></svg>
                                </div>
                                <div>
                                    <p class="font-semibold text-slate-700">{{ getModeName(round) }}</p>
                                    <p class="text-xs text-gray-500 truncate max-w-[150px]">{{ getLocationLabel(round) }}</p>
                                </div>
                            </div>

                            <div class="divider my-1"></div>

                            <div class="flex justify-between items-end">
                                <div class="flex-1 mr-4">
                                    <div class="flex justify-between text-xs mb-1">
                                        <span class="text-gray-500">การจอง ({{ round.slots?.length > 0 ? Math.round((getBookedCount(round)/round.slots.length)*100) : 0 }}%)</span>
                                        <span class="font-bold text-slate-700">{{ getBookedCount(round) }}/{{ round.slots?.length || 0 }}</span>
                                    </div>
                                    <progress class="progress w-full h-2" :class="getBookedCount(round) === (round.slots?.length || 0) ? 'progress-error' : 'progress-primary'" :value="getBookedCount(round)" :max="round.slots?.length || 1"></progress>
                                </div>
                                <div class="avatar-group -space-x-3">
                                    <div v-for="interviewer in getInterviewersForRound(round).slice(0, 3)" :key="interviewer.ID" class="avatar placeholder border-white">
                                        <div class="bg-neutral-focus text-neutral-content w-8 h-8 text-[10px]"><span>{{ getInitials(interviewer) }}</span></div>
                                    </div>
                                    <div v-if="getInterviewersForRound(round).length > 3" class="avatar placeholder border-white">
                                        <div class="bg-gray-200 text-gray-600 w-8 h-8 text-[10px]"><span>+{{ getInterviewersForRound(round).length - 3 }}</span></div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                <div v-else class="flex flex-col items-center justify-center py-20 text-center animate-fade-in">
                    <div class="bg-gray-50 p-6 rounded-full mb-4">
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-10 w-10 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
                    </div>
                    <h3 class="text-lg font-bold text-slate-600">ไม่มีรอบสัมภาษณ์ที่เปิดอยู่</h3>
                    <p class="text-gray-500">สร้างรอบใหม่ได้ที่ปุ่มมุมขวาบน</p>
                </div>
            </div>

            <div v-else class="animate-fade-in">
                <div class="overflow-x-auto bg-white rounded-2xl border border-gray-100">
                    <table class="table w-full">
                        <thead class="bg-slate-50 text-slate-600 uppercase text-xs">
                            <tr>
                                <th class="py-4">วันที่</th>
                                <th>ชื่อรอบสัมภาษณ์</th>
                                <th>ประเภททุน</th>
                                <th>สถิติการจอง</th>
                                <th>สถานะ</th>
                                <th class="text-right">จัดการ</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="round in filteredRounds" :key="round.ID" class="hover:bg-slate-50 transition-colors">
                                <td class="font-medium text-slate-700">
                                    <div class="flex flex-col">
                                        <span>{{ new Date(round.start_date_time).toLocaleDateString('th-TH', { dateStyle: 'medium' }) }}</span>
                                        <span class="text-xs text-gray-400 font-light">{{ new Date(round.start_date_time).toLocaleTimeString('th-TH', {hour: '2-digit', minute: '2-digit'}) }} - {{ new Date(round.end_date_time).toLocaleTimeString('th-TH', {hour: '2-digit', minute: '2-digit'}) }}</span>
                                    </div>
                                </td>
                                <td><div class="font-bold text-[#1e3a8a]">{{ round.name }}</div></td>
                                <td><div v-if="round.scholarship" class="badge badge-outline text-xs">{{ round.scholarship.scholarship_name }}</div></td>
                                <td>
                                    <div class="flex items-center gap-2">
                                        <span class="text-sm font-semibold">{{ getBookedCount(round) }}/{{ round.slots?.length || 0 }}</span>
                                        <progress class="progress progress-primary w-16 h-1.5" :value="getBookedCount(round)" :max="round.slots?.length || 1"></progress>
                                    </div>
                                </td>
                                <td><div class="badge bg-gray-100 text-gray-500 border-none">เสร็จสิ้น</div></td>
                                <td class="text-right"><button @click="openViewModal(round)" class="btn btn-ghost btn-xs text-blue-600">ดูข้อมูล</button></td>
                            </tr>
                            <tr v-if="filteredRounds.length === 0">
                                <td colspan="6" class="text-center py-12 text-gray-400">ไม่พบประวัติย้อนหลัง</td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>
        </div>

        <div v-if="isModalOpen" class="fixed inset-0 z-[100] flex items-center justify-center bg-black/60 backdrop-blur-sm p-4">
            <div class="bg-white w-full max-w-5xl h-[90vh] rounded-2xl shadow-2xl flex flex-col overflow-hidden animate-pop-in">
                
                <div class="px-6 py-4 border-b flex items-center justify-between bg-slate-50">
                    <div>
                        <h2 class="text-xl font-bold text-[#1e3a8a]">
                            {{ modalMode === 'create' ? 'สร้างรอบสัมภาษณ์ใหม่' : (modalMode === 'edit' ? 'แก้ไขรอบสัมภาษณ์' : 'รายละเอียดรอบสัมภาษณ์') }}
                        </h2>
                        <p class="text-xs text-gray-500 mt-1">
                            {{ modalMode === 'create' ? 'กรอกข้อมูลเพื่อสร้าง Slot อัตโนมัติ' : (modalMode === 'edit' ? 'การแก้ไขเวลาจะมีผลต่อ Slot เดิม' : 'ตรวจสอบรายชื่อผู้จองและจัดการข้อมูล') }}
                        </p>
                    </div>
                    <button @click="isModalOpen = false" class="btn btn-circle btn-ghost btn-sm text-gray-500">✕</button>
                </div>

                <div class="flex-1 overflow-y-auto p-6 bg-[#f8fafc]">
                    
                    <div v-if="isLoading && (modalMode === 'view' || modalMode === 'edit')" class="flex flex-col justify-center items-center h-full gap-4">
                        <span class="loading loading-spinner loading-lg text-primary"></span>
                        <p class="text-gray-500 animate-pulse">กำลังโหลดข้อมูลรอบสัมภาษณ์...</p>
                    </div>

                    <div v-else class="grid grid-cols-1 lg:grid-cols-3 gap-6">
                        <div class="lg:col-span-1 space-y-6">
                            <div class="card bg-white shadow-sm border border-gray-100 p-5 space-y-4">
                                <h3 class="font-bold text-slate-700 border-b pb-2">ข้อมูลทั่วไป</h3>
                                <div class="form-control w-full">
                                    <label class="label"><span class="label-text font-medium">ชื่อรอบสัมภาษณ์</span></label>
                                    <input v-model="formData.name" type="text" class="input input-bordered input-sm w-full" :disabled="modalMode === 'view'" />
                                </div>
                                <div class="form-control w-full">
                                    <label class="label"><span class="label-text font-medium">สำหรับทุน</span></label>
                                    <select v-model="formData.scholarship_id" class="select select-bordered select-sm w-full" :disabled="modalMode === 'view'">
                                        <option disabled :value="null">เลือกทุน...</option>
                                        <option v-for="s in scholarships" :key="s.ID" :value="s.ID">{{ s.scholarship_name }}</option>
                                    </select>
                                </div>
                                <div class="form-control w-full">
                                    <label class="label"><span class="label-text font-medium">รายละเอียดเพิ่มเติม</span></label>
                                    <textarea v-model="formData.description" class="textarea textarea-bordered text-sm h-20" :disabled="modalMode === 'view'"></textarea>
                                </div>
                            </div>

                            <div class="card bg-white shadow-sm border border-gray-100 p-5 space-y-4">
                                <h3 class="font-bold text-slate-700 border-b pb-2">เวลา & สถานที่</h3>
                                <div class="form-control w-full">
                                    <label class="label"><span class="label-text font-medium">วันที่</span></label>
                                    <div class="grid grid-cols-2 gap-2">
                                        <input v-model="formData.date" type="date" class="input input-bordered input-sm w-full" :disabled="modalMode === 'view'" />
                                        <input v-model="formData.end_date" type="date" class="input input-bordered input-sm w-full" :disabled="modalMode === 'view'" />
                                    </div>
                                </div>
                                <div class="grid grid-cols-2 gap-2">
                                    <div class="form-control">
                                        <label class="label"><span class="label-text font-medium">เริ่ม</span></label>
                                        <input v-model="formData.start_time" type="time" class="input input-bordered input-sm w-full" :disabled="modalMode === 'view'" />
                                    </div>
                                    <div class="form-control">
                                        <label class="label"><span class="label-text font-medium">สิ้นสุด</span></label>
                                        <input v-model="formData.end_time" type="time" class="input input-bordered input-sm w-full" :disabled="modalMode === 'view'" />
                                    </div>
                                </div>
                                <div class="form-control w-full">
                                    <label class="label"><span class="label-text font-medium">ระยะเวลา/คน(นาที)</span></label>
                                    <select v-model="formData.slot_duration" class="select select-bordered select-sm w-full" :disabled="modalMode === 'view'">
                                        <option :value="10">10 นาที</option><option :value="15">15 นาที</option><option :value="20">20 นาที</option><option :value="30">30 นาที</option><option :value="60">1 ชั่วโมง</option>
                                    </select>
                                </div>
                                <div class="divider my-0"></div>
                                <div class="form-control w-full">
                                    <label class="label"><span class="label-text font-medium">รูปแบบการสัมภาษณ์</span></label>
                                    <select v-model="formData.interview_mode_id" class="select select-bordered select-sm w-full" :disabled="modalMode === 'view'">
                                        <option disabled :value="null">เลือกรูปแบบ...</option>
                                        <option v-for="mode in interviewModes" :key="mode.ID" :value="mode.ID">{{ mode.name }}</option>
                                    </select>
                                </div>
                                <div v-if="selectedModeName === 'Onsite'" class="form-control w-full animate-fade-in">
                                    <label class="label"><span class="label-text font-medium">สถานที่</span></label>
                                    <select v-model="formData.location_id" class="select select-bordered select-sm w-full" :disabled="modalMode === 'view'">
                                        <option v-for="loc in locations" :key="loc.ID" :value="loc.ID">{{ loc.name }}</option>
                                    </select>
                                </div>
                                <div v-if="selectedModeName === 'Online'" class="form-control w-full animate-fade-in">
                                    <label class="label"><span class="label-text font-medium">Meeting Link</span></label>
                                    <input v-model="formData.meeting_link" type="text" class="input input-bordered input-sm w-full" :disabled="modalMode === 'view'" />
                                </div>
                            </div>
                        </div>

                        <div class="lg:col-span-2 space-y-6">
                            <div class="card bg-white shadow-sm border border-gray-100 p-5">
                                <h3 class="font-bold text-slate-700 mb-3">กรรมการสัมภาษณ์</h3>
                                <div class="flex flex-wrap gap-2">
                                    <label v-for="interviewer in interviewersList" :key="interviewer.ID"
                                        class="cursor-pointer border rounded-lg px-3 py-2 flex items-center gap-2 transition-all"
                                        :class="formData.interviewer_ids.includes(interviewer.ID) ? 'border-blue-500 bg-blue-50 ring-1 ring-blue-500' : 'border-gray-200'">
                                        <input type="checkbox" :value="interviewer.ID" v-model="formData.interviewer_ids" class="checkbox checkbox-primary checkbox-xs" :disabled="modalMode === 'view'" />
                                        <span class="text-sm">{{ interviewer.interviewer_firstname }} {{ interviewer.interviewer_lastname }}</span>
                                    </label>
                                </div>
                            </div>

                            <div class="card bg-white shadow-sm border border-gray-100 p-5 flex-1 min-h-[400px]">
                                <div class="flex items-center justify-between mb-4">
                                    <h3 class="font-bold text-slate-700 flex items-center gap-2">
                                        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-gray-500" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
                                        {{ (modalMode === 'create' || isSlotConfigDirty) ? 'พรีวิวตารางเวลาที่จะถูกสร้างใหม่' : 'ตารางการจองปัจจุบัน' }}
                                    </h3>
                                    <div class="text-xs text-gray-500">
                                        ทั้งหมด: {{ (modalMode === 'create' || isSlotConfigDirty) ? previewSlots.length : selectedRoundDetails?.slots?.length || 0 }} สล็อต
                                    </div>
                                </div>

                                <div v-if="modalMode === 'create' || isSlotConfigDirty">
                                    <div v-if="isSlotConfigDirty" class="p-3 bg-amber-50 text-amber-800 text-xs rounded-lg mb-4 border border-amber-200">
                                        ⚠️ <strong>คำเตือน:</strong> คุณกำลังแก้ไขข้อมูลเวลา/วันที่ ระบบจะลบ Slots เดิมทั้งหมด และสร้างใหม่ตามพรีวิวด้านล่างนี้เมื่อกดบันทึก
                                    </div>
                                    <div v-if="previewSlots.length > 0" class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-3">
                                        <div v-for="(slot, idx) in previewSlots" :key="idx" class="p-3 rounded-xl border text-center bg-white border-green-200">
                                            <span class="font-bold text-lg text-slate-700">{{ slot.time }}</span>
                                            <div class="badge badge-success badge-xs text-white block mx-auto mt-1">ว่าง</div>
                                        </div>
                                    </div>
                                    <div v-else class="flex flex-col items-center justify-center h-40 text-gray-400 border-2 border-dashed rounded-xl">
                                        <p>ระบุวันและเวลา เพื่อดูพรีวิวตาราง</p>
                                    </div>
                                </div>

                                <div v-else>
                                    <div v-if="selectedRoundDetails?.slots && selectedRoundDetails.slots.length > 0" class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-3">
                                        <div v-for="slot in selectedRoundDetails.slots" :key="slot.ID">
                                            <button @click="openStudentDetailModal(slot)" 
                                                :disabled="!slot.is_booked"
                                                class="w-full p-3 rounded-xl border text-center flex flex-col items-center justify-center transition-all h-full"
                                                :class="{
                                                    'bg-blue-50 border-blue-300 cursor-pointer hover:bg-blue-100 hover:border-blue-400': slot.is_booked,
                                                    'bg-white border-green-200 cursor-not-allowed': !slot.is_booked
                                                }">
                                                <span class="font-bold text-lg" :class="slot.is_booked ? 'text-slate-700' : 'text-slate-700'">
                                                    {{ new Date(slot.start_time).toLocaleTimeString('th-TH', { hour: '2-digit', minute: '2-digit' }) }}
                                                </span>
                                                
                                                <div v-if="slot.is_booked && slot.interviewe_bookings && slot.interviewe_bookings[0]?.application_scholarship?.application?.student_profile" 
                                                     class="mt-1 text-xs text-blue-700 font-semibold truncate w-full px-1">
                                                    {{ slot.interviewe_bookings[0].application_scholarship.application.student_profile.first_name_th }}
                                                </div>
                                                <div v-else class="badge badge-xs mt-1" :class="slot.is_booked ? 'badge-ghost text-gray-400' : 'badge-success text-white'">
                                                    {{ slot.is_booked ? 'จองแล้ว' : 'ว่าง' }}
                                                </div>
                                            </button>
                                        </div>
                                    </div>
                                    <div v-else class="flex flex-col items-center justify-center h-40 text-gray-400 border-2 border-dashed rounded-xl">
                                        <p>ไม่มีข้อมูลช่องเวลาสำหรับรอบนี้</p>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>

                <div class="p-4 border-t bg-white flex justify-end gap-2">
                    <button @click="isModalOpen = false" class="btn btn-ghost text-gray-500">ยกเลิก</button>
                    <button v-if="modalMode === 'create'" @click="saveRound" class="btn bg-[#1e3a8a] text-white hover:bg-[#152c6f]">สร้างรอบสัมภาษณ์</button>
                    <button v-if="modalMode === 'edit'" @click="handleUpdate" class="btn bg-blue-600 text-white hover:bg-blue-700">บันทึกการแก้ไข</button>
                </div>
            </div>
        </div>

        <!-- Student Detail Modal -->
        <div v-if="isStudentDetailModalOpen && selectedStudent" class="fixed inset-0 z-[101] flex items-center justify-center bg-black/60 backdrop-blur-sm p-4">
            <div class="bg-white w-full max-w-md rounded-2xl shadow-2xl flex flex-col overflow-hidden animate-pop-in">
                <div class="px-6 py-4 border-b flex items-center justify-between bg-slate-50">
                    <h2 class="text-lg font-bold text-[#1e3a8a]">รายละเอียดผู้สมัคร</h2>
                    <button @click="isStudentDetailModalOpen = false" class="btn btn-circle btn-ghost btn-sm text-gray-500">✕</button>
                </div>
                <div class="p-6 space-y-4">
                    <div class="flex items-center space-x-4">
                        <div class="avatar placeholder">
                            <div class="bg-neutral-focus text-neutral-content rounded-full w-16">
                                <span class="text-xl">{{ selectedStudent.first_name_th?.charAt(0) }}{{ selectedStudent.last_name_th?.charAt(0) }}</span>
                            </div>
                        </div>
                        <div>
                            <h3 class="font-bold text-xl text-slate-800">{{ selectedStudent.first_name_th }} {{ selectedStudent.last_name_th }}</h3>
                            <p class="text-sm text-gray-500">{{ selectedStudent.student_id }}</p>
                        </div>
                    </div>
                    <div class="divider my-2"></div>
                    <div class="grid grid-cols-1 sm:grid-cols-2 gap-x-4 gap-y-3 text-sm">
                        <div>
                            <label class="font-semibold text-gray-500 text-xs">Email</label>
                            <p class="text-gray-800">{{ selectedStudent.email }}</p>
                        </div>
                        <div>
                            <label class="font-semibold text-gray-500 text-xs">Major</label>
                            <p class="text-gray-800">{{ selectedStudent.major?.major_name || '-' }}</p>
                        </div>
                        <div>
                            <label class="font-semibold text-gray-500 text-xs">GPAX</label>
                            <p class="text-gray-800">{{ selectedStudent.gpax?.toFixed(2) || '-' }}</p>
                        </div>
                    </div>
                </div>
                <div class="p-4 border-t bg-gray-50/70 flex justify-end">
                    <button @click="isStudentDetailModalOpen = false" class="btn btn-ghost">ปิด</button>
                </div>
            </div>
        </div>

    </div>
</template>

<style scoped>
.animate-pop-in { animation: pop-in 0.2s cubic-bezier(0.16, 1, 0.3, 1) forwards; }
.animate-fade-in { animation: fade 0.3s ease forwards; }
@keyframes pop-in { 0% { opacity: 0; transform: scale(0.95); } 100% { opacity: 1; transform: scale(1); } }
@keyframes fade { from { opacity: 0; transform: translateY(-5px); } to { opacity: 1; transform: translateY(0); } }
::-webkit-scrollbar { width: 6px; }
::-webkit-scrollbar-track { background: transparent; }
::-webkit-scrollbar-thumb { background: #cbd5e1; border-radius: 10px; }
</style>