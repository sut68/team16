<script setup lang="ts">
import { ref, computed, reactive, watch, onMounted } from 'vue';
import { InterviewAPI, type InterviewRoundCreate } from '@/services/api'; // ตรวจสอบ path ให้ตรงกับโปรเจคจริง
import { ScholarshipAPI } from '@/services/api/scholarship'; // ตรวจสอบ path ให้ตรงกับโปรเจคจริง
import { LocationAPI } from '@/services/api/location';
import type { InterviewRound, Interviewer, Location, InterviewMode } from '@/interfaces/interview';
import type { ScholarshipResponse } from '@/interfaces/scholarship';
import Swal from 'sweetalert2';

// --- 1. State Management ---
const isLoading = ref(true);
const isModalOpen = ref(false);
const isFilterOpen = ref(false);
const modalMode = ref<'create' | 'view'>('create');
const activeTab = ref<'active' | 'history'>('active');

// Data from API
const allRounds = ref<InterviewRound[]>([]);
const scholarships = ref<ScholarshipResponse[]>([]);
const interviewersList = ref<Interviewer[]>([]);
const locations = ref<Location[]>([]);
const interviewModes = ref<InterviewMode[]>([]); // New state for modes
const selectedRoundDetails = ref<InterviewRound | null>(null);

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
    interview_mode_id: null as number | null, // Changed from mode
    location_id: null as number | null,
    meeting_link: '',
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
            InterviewAPI.getAllModes(), // Fetch modes
        ]);

        if (roundsRes.status === 'fulfilled') {
            allRounds.value = roundsRes.value || [];
        } else {
            console.error("Failed to fetch interview rounds:", roundsRes.reason);
            allRounds.value = [];
        }

        if (scholarsRes.status === 'fulfilled') {
            scholarships.value = scholarsRes.value || [];
        } else {
            console.error("Failed to fetch scholarships:", scholarsRes.reason);
            scholarships.value = [];
        }

        if (interviewersRes.status === 'fulfilled') {
            interviewersList.value = interviewersRes.value || [];
        } else {
            console.error("Failed to fetch interviewers:", interviewersRes.reason);
            interviewersList.value = [];
        }

        if (locationsRes.status === 'fulfilled') {
            locations.value = locationsRes.value || [];
        } else {
            console.error("Failed to fetch locations:", locationsRes.reason);
            locations.value = [];
        }
		
        if (modesRes.status === 'fulfilled') {
            interviewModes.value = modesRes.value || [];
        } else {
            console.error("Failed to fetch interview modes:", modesRes.reason);
            interviewModes.value = [];
        }

    } catch (error) {
        console.error(error);
        Swal.fire('Error', 'Could not fetch data', 'error');
    } finally {
        isLoading.value = false;
    }
};

// --- 3. Computed Logic ---
const getModeName = (round: InterviewRound) => {
    // กรณีมี object แนบมา (ขึ้นอยู่กับ backend response)
    if ((round as any).interview_mode?.name) return (round as any).interview_mode.name;
    // กรณีมีแค่ ID ให้ค้นจาก state
    const mode = interviewModes.value.find(m => m.ID === round.interview_mode_id);
    return mode ? mode.name : 'Unknown';
};

// ดึงชื่อสถานที่ หรือ ข้อความ Online Meeting
const getLocationLabel = (round: InterviewRound) => {
    const modeName = getModeName(round);
    if (modeName === 'Onsite') {
        // กรณีมี object location แนบมา
        if ((round as any).location?.name) return (round as any).location.name;
        // กรณีมีแค่ ID
        const loc = locations.value.find(l => l.ID === round.location_id);
        return loc ? loc.name : '-';
    }
    return 'Online Meeting';
};
const getRoundStatus = (round: InterviewRound): 'Closed' | 'Full' | 'Open' => {
    const roundEnd = new Date(round.end_date_time);
    
    // Set 'today' to the beginning of the current day (00:00:00)
    const today = new Date();
    today.setHours(0, 0, 0, 0);

    // A round is closed only if its end time is before the start of today
    if (roundEnd.getTime() < today.getTime()) {
        return 'Closed';
    }
    if (!round.slots || round.slots.length === 0) {
        return 'Open';
    }
    const isFull = round.slots.every(slot => slot.is_booked);
    if (isFull) {
        return 'Full';
    }
    return 'Open';
}

const getBookedCount = (round: InterviewRound) => {
    if (!round.slots) return 0;
    return round.slots.filter(s => s.is_booked).length;
}

const stats = computed(() => {
    const activeRounds = allRounds.value.filter(r => getRoundStatus(r) !== 'Closed');
    const totalRounds = activeRounds.length;
    const totalCapacity = activeRounds.reduce((acc, r) => acc + (r.slots?.length || 0), 0);
    const totalBooked = activeRounds.reduce((acc, r) => acc + getBookedCount(r), 0);
    const bookedPercent = totalCapacity > 0 ? Math.round((totalBooked / totalCapacity) * 100) : 0;
    return { totalRounds, totalCapacity, totalBooked, bookedPercent };
});

const filteredRounds = computed(() => {
    let result = allRounds.value;

    // Tab Logic
    if (activeTab.value === 'active') {
        result = result.filter(r => getRoundStatus(r) !== 'Closed');
    } else {
        result = result.filter(r => getRoundStatus(r) === 'Closed');
    }

    // Search
    if (searchQuery.value) {
        const q = searchQuery.value.toLowerCase();
        result = result.filter(r =>
            r.name.toLowerCase().includes(q) ||
            (r.scholarship && r.scholarship.scholarship_name.toLowerCase().includes(q))
        );
    }

    // Filters
    if (filterStatus.value !== 'All') {
        result = result.filter(r => getRoundStatus(r) === filterStatus.value);
    }
    if (filterScholarship.value !== 'All') {
        // เช็ค structure จริงๆ ว่า typescholarship อยู่ตรงไหน
        // สมมติว่าอยู่ที่ r.scholarship.typescholarship
        result = result.filter(r => r.scholarship?.typescholarship?.type_name === filterScholarship.value);
    }

    // Sort
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

    // สร้าง Date object แบบง่าย (ระวังเรื่อง Timezone หากใช้จริงจังควรใช้ library เช่น dayjs)
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
const openCreateModal = () => {
    modalMode.value = 'create';
    // Reset Form (แก้ไข: เพิ่มการ reset mode, location_id, meeting_link)
    Object.assign(formData, {
        name: '',
        scholarship_id: null,
        description: '',
        date: '',
        end_date: '',
        start_time: '09:00',
        end_time: '16:00',
        slot_duration: 30,
        interviewer_ids: [],
        interview_mode_id: null,
        location_id: null,
        meeting_link: ''
    });
    selectedRoundDetails.value = null;
    isModalOpen.value = true;
};

const openViewModal = async (round: InterviewRound) => {
    modalMode.value = 'view';
    isLoading.value = true;
    isModalOpen.value = true;
    try {
        const details = await InterviewAPI.getRoundById(round.ID);
        selectedRoundDetails.value = details;
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
            meeting_link: details.meeting_link
        });
    } catch (e) {
        Swal.fire('Error', 'Could not fetch round details', 'error');
        isModalOpen.value = false;
    } finally {
        isLoading.value = false;
    }
};

const handleDelete = async (id: number) => {
    const result = await Swal.fire({
        title: 'Are you sure?',
        text: "You won't be able to revert this!",
        icon: 'warning',
        showCancelButton: true,
        confirmButtonColor: '#d33',
        cancelButtonColor: '#3085d6',
        confirmButtonText: 'Yes, delete it!'
    });

    if (result.isConfirmed) {
        try {
            await InterviewAPI.deleteRound(id);
            Swal.fire('Deleted!', 'The interview round has been deleted.', 'success');
            await fetchData();
        } catch (e) {
            Swal.fire('Error', 'Failed to delete round.', 'error');
        }
    }
};

const saveRound = async () => {
    if (!formData.scholarship_id || !formData.date || !formData.start_time || !formData.end_time || !formData.end_date || !formData.interview_mode_id) {
        Swal.fire('Invalid Data', 'Please fill all required fields', 'warning');
        return;
    }

    const startDateTime = new Date(`${formData.date}T${formData.start_time}`).toISOString();
    const endDateTime = new Date(`${formData.end_date}T${formData.end_time}`).toISOString();

    // สมมติ Admin ID (ในระบบจริงควรดึงจาก Store)
    const adminProfileId = 1;

    const payload: InterviewRoundCreate = {
        name: formData.name,
        description: formData.description,
        start_date_time: startDateTime,
        end_date_time: endDateTime,
        slot_duration: formData.slot_duration,
        scholarship_id: formData.scholarship_id,
        admin_profile_id: adminProfileId,
        interviewer_ids: formData.interviewer_ids,
        interview_mode_id: formData.interview_mode_id,
        location_id: formData.location_id ? Number(formData.location_id) : null,
        meeting_link: formData.meeting_link
    };

    isLoading.value = true;
    try {
        await InterviewAPI.createRound(payload);
        Swal.fire('Success', 'Interview round created successfully!', 'success');
        isModalOpen.value = false;
        await fetchData();
    } catch (error: any) {
        console.error(error);
        Swal.fire('Error', error.response?.data?.error || 'Failed to create interview round', 'error');
    } finally {
        isLoading.value = false;
    }
};
</script>

<template>
    <div class="w-full mx-auto flex flex-col h-full p-6 bg-white rounded-tl-[30px] shadow overflow-visible"
        data-theme="light">

        <div class="flex flex-col md:flex-row items-start md:items-center justify-between gap-4 mb-8">
            <div>
            </div>
            <button @click="openCreateModal"
                class="btn btn-sm bg-white border border-gray-300 text-gray-700 hover:bg-gray-100 hover:border-gray-400 flex items-center justify-center gap-2 rounded-full px-5 h-10 w-full md:w-auto shadow-sm transition-all duration-150">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                    <path fill-rule="evenodd"
                        d="M10 3a1 1 0 011 1v5h5a1 1 0 110 2h-5v5a1 1 0 11-2 0v-5H4a1 1 0 110-2h5V4a1 1 0 011-1z"
                        clip-rule="evenodd" />
                </svg>
                สร้างรอบสัมภาษณ์
            </button>
        </div>

        <div
            class="stats stats-vertical lg:stats-horizontal shadow bg-white border border-gray-100 w-full mb-8 relative z-0">
            <div class="stat">
                <div class="stat-figure text-[#1e3a8a] bg-blue-50 p-3 rounded-full">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
                        class="inline-block w-8 h-8 stroke-current">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10">
                        </path>
                    </svg>
                </div>
                <div class="stat-title text-slate-500">รอบที่เปิดอยู่</div>
                <div class="stat-value text-[#1e3a8a]">{{ stats.totalRounds }}</div>
                <div class="stat-desc">กำลังดำเนินการ</div>
            </div>

            <div class="stat">
                <div class="stat-figure text-secondary bg-pink-50 p-3 rounded-full">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
                        class="inline-block w-8 h-8 stroke-current">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z">
                        </path>
                    </svg>
                </div>
                <div class="stat-title text-slate-500">นศ. ที่จองแล้ว</div>
                <div class="stat-value text-secondary">{{ stats.totalBooked }}</div>
                <div class="stat-desc text-success font-medium">↗︎ {{ stats.bookedPercent }}% ของที่นั่งทั้งหมด</div>
            </div>

            <div class="stat">
                <div class="stat-figure text-orange-500 bg-orange-50 p-3 rounded-full">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
                        class="inline-block w-8 h-8 stroke-current">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                    </svg>
                </div>
                <div class="stat-title text-slate-500">ที่นั่งว่างเหลือ</div>
                <div class="stat-value text-orange-500">{{ stats.totalCapacity - stats.totalBooked }}</div>
                <div class="stat-desc">ที่นั่ง</div>
            </div>
        </div>

        <div
            class="relative z-50 flex flex-col xl:flex-row items-end xl:items-center justify-between gap-4 mb-6 border-b border-gray-200">

            <div class="flex gap-8 -mb-[1px]">
                <a @click="activeTab = 'active'"
                    class="pb-3 px-1 text-base font-medium cursor-pointer transition-all duration-200 border-b-[3px]"
                    :class="activeTab === 'active'
                        ? 'text-[#1e3a8a] border-[#1e3a8a]'
                        : 'text-slate-500 border-transparent hover:text-slate-700 hover:border-slate-300'">
                    รอบปัจจุบัน
                </a>
                <a @click="activeTab = 'history'"
                    class="pb-3 px-1 text-base font-medium cursor-pointer transition-all duration-200 border-b-[3px]"
                    :class="activeTab === 'history'
                        ? 'text-[#1e3a8a] border-[#1e3a8a]'
                        : 'text-slate-500 border-transparent hover:text-slate-700 hover:border-slate-300'">
                    ประวัติย้อนหลัง
                </a>
            </div>

            <div class="flex flex-col md:flex-row items-center gap-2 w-full xl:w-auto pb-4 xl:pb-2">

                <div class="relative w-full md:w-auto">
                    <button @click="isFilterOpen = !isFilterOpen"
                        class="btn btn-sm btn-outline bg-white border-gray-300 text-gray-600 hover:bg-gray-50 hover:border-gray-400 gap-2 h-10 rounded-full font-normal px-5 w-full md:w-auto shadow-sm">
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
                            <path fill-rule="evenodd"
                                d="M3 3a1 1 0 011-1h12a1 1 0 011 1v3a1 1 0 01-.293.707L12 11.414V15a1 1 0 01-.293.707l-2 2A1 1 0 018 17v-5.586L3.293 6.707A1 1 0 013 6V3z"
                                clip-rule="evenodd" />
                        </svg>
                        ตัวกรอง
                    </button>

                    <div v-if="isFilterOpen" class="fixed inset-0 z-50" @click="isFilterOpen = false"></div>
                    <div v-if="isFilterOpen"
                        class="absolute right-0 mt-2 w-72 bg-white rounded-xl shadow-2xl z-[100] border p-4 animate-pop-in">

                        <div v-if="activeTab === 'active'">
                            <p class="font-bold text-base mb-3 text-slate-700">กรองรอบปัจจุบัน</p>
                            <div class="space-y-3">
                                <div>
                                    <label class="text-xs font-medium text-gray-500">สถานะ</label>
                                    <select v-model="filterStatus"
                                        class="select select-bordered select-sm w-full h-10 bg-white text-sm">
                                        <option value="All">ทั้งหมด</option>
                                        <option value="Open">เปิดให้จอง</option>
                                        <option value="Full">เต็มแล้ว</option>
                                    </select>
                                </div>
                                <div>
                                    <label class="text-xs font-medium text-gray-500">ประเภททุน</label>
                                    <select v-model="filterScholarship"
                                        class="select select-bordered select-sm w-full h-10 bg-white text-sm">
                                        <option value="All">ทั้งหมด</option>
                                        <option
                                            v-for="t in [...new Set(scholarships.map(s => s.typescholarship.type_name))]"
                                            :key="t" :value="t">{{ t }}</option>
                                    </select>
                                </div>
                            </div>
                        </div>

                        <div v-else>
                            <p class="font-bold text-base mb-3 text-slate-700">กรองประวัติย้อนหลัง</p>
                            <select v-model="filterScholarship"
                                class="select select-bordered select-sm w-full h-10 bg-white text-sm">
                                <option value="All">ทั้งหมด</option>
                                <option v-for="t in [...new Set(scholarships.map(s => s.typescholarship.type_name))]"
                                    :key="t" :value="t">{{ t }}</option>
                            </select>
                        </div>

                    </div>
                </div>

                <select v-model="sortBy"
                    class="select select-bordered select-sm rounded-full h-10 bg-white text-sm border-gray-300 focus:border-[#1e3a8a] focus:ring-[#1e3a8a] w-full md:w-40 shadow-sm px-6">
                    <option value="DateAsc">วันที่ (เก่าสุดก่อน)</option>
                    <option value="DateDesc">วันที่ (ใหม่สุดก่อน)</option>
                    <option value="BookedDesc">ยอดจองมากที่สุด</option>
                </select>

                <div class="relative w-full md:w-64">
                    <span class="absolute inset-y-0 left-0 flex items-center pl-3 text-gray-400">
                        <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24"
                            stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                        </svg>
                    </span>
                    <input type="text" v-model="searchQuery" placeholder="ค้นหาชื่อรอบ, ทุน..."
                        class="input input-bordered input-sm h-10 w-full rounded-full pl-10 bg-white border-gray-300 shadow-sm focus:border-[#1e3a8a] focus:ring-1 focus:ring-[#1e3a8a] text-sm" />
                </div>

            </div>
        </div>

        <div class="bg-slate-50/30 rounded-3xl shadow-sm border border-slate-100 p-6 min-h-[600px] flex-1 relative z-0">

            <div v-if="activeTab === 'active'" class="animate-fade-in">
                <div v-if="filteredRounds.length > 0" class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6">
    <div v-for="round in filteredRounds" :key="round.ID"
        class="group card bg-white border border-gray-200 hover:border-blue-400 hover:shadow-xl transition-all duration-300 rounded-2xl overflow-visible">

        <div class="px-5 py-4 border-b border-gray-50 bg-slate-50/50 flex justify-between items-start rounded-t-2xl relative">
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
                <h3 class="font-bold text-[#1e3a8a] text-lg leading-tight group-hover:text-blue-600 transition-colors">
                    {{ round.name }}
                </h3>
            </div>

            <div class="dropdown dropdown-end">
                <label tabindex="0" class="btn btn-circle btn-ghost btn-sm text-gray-400 hover:bg-white">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="w-5 h-5 stroke-current">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 5v.01M12 12v.01M12 19v.01M12 6a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2z"></path>
                    </svg>
                </label>
                <ul tabindex="0" class="dropdown-content menu p-2 shadow-lg bg-white rounded-box w-40 z-[10] border border-gray-100">
                    <li><a @click="openViewModal(round)" class="text-slate-600">รายละเอียด</a></li>
                    <li><a class="text-slate-600">แก้ไข</a></li>
                    <div class="divider my-0"></div>
                    <li><a @click="handleDelete(round.ID)" class="text-error hover:bg-red-50">ลบรายการ</a></li>
                </ul>
            </div>
        </div>

        <div class="p-5 space-y-4 cursor-pointer" @click="openViewModal(round)">
            
            <div class="flex items-center gap-3 text-sm text-gray-600">
                <div class="w-8 h-8 rounded-full bg-blue-50 flex items-center justify-center text-blue-600">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                    </svg>
                </div>
                <div>
                    <p class="font-semibold text-slate-700">
                        วันที่ {{ new Date(round.start_date_time).toLocaleDateString('th-TH', { day: 'numeric', month: 'short', year: '2-digit'}) }}
                    </p>
                    <p class="text-xs">
                        {{ new Date(round.start_date_time).toLocaleTimeString('th-TH', {hour: '2-digit', minute: '2-digit'}) }} - 
                        {{ new Date(round.end_date_time).toLocaleTimeString('th-TH', {hour: '2-digit', minute: '2-digit'}) }}
                    </p>
                </div>
            </div>

            <div class="flex items-center gap-3 text-sm text-gray-600">
                <div class="w-8 h-8 rounded-full flex items-center justify-center" 
                    :class="getModeName(round) === 'Onsite' ? 'bg-orange-50 text-orange-600' : 'bg-purple-50 text-purple-600'">
                    
                    <svg v-if="getModeName(round) === 'Onsite'" xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
                        <path fill-rule="evenodd" d="M5.05 4.05a7 7 0 119.9 9.9L10 18.9l-4.95-4.95a7 7 0 010-9.9zM10 11a2 2 0 100-4 2 2 0 000 4z" clip-rule="evenodd" />
                    </svg>
                    <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
                        <path d="M2 6a2 2 0 012-2h6a2 2 0 012 2v8a2 2 0 01-2 2H4a2 2 0 01-2-2V6zM14.553 7.106A1 1 0 0014 8v4a1 1 0 00.553.894l2 1A1 1 0 0018 13V7a1 1 0 00-1.447-.894l-2 1z" />
                    </svg>
                </div>
                <div>
                    <p class="font-semibold text-slate-700">{{ getModeName(round) }}</p>
                    <p class="text-xs truncate max-w-[180px]">{{ getLocationLabel(round) }}</p>
                </div>
            </div>

            <div class="divider my-1"></div>

            <div class="flex justify-between items-end">
                <div class="flex-1 mr-4">
                    <div class="flex justify-between text-xs mb-1">
                        <span class="text-gray-500">การจอง ({{ round.slots?.length > 0 ? Math.round((getBookedCount(round)/round.slots.length)*100) : 0 }}%)</span>
                        <span class="font-bold text-slate-700">{{ getBookedCount(round) }}/{{ round.slots?.length || 0 }}</span>
                    </div>
                    <progress class="progress w-full h-2"
                        :class="getBookedCount(round) === (round.slots?.length || 0) ? 'progress-error' : 'progress-primary'"
                        :value="getBookedCount(round)" :max="round.slots?.length || 1"></progress>
                </div>

                <div class="avatar-group -space-x-3">
                    <div v-for="interviewer in getInterviewersForRound(round).slice(0, 3)" :key="interviewer.ID" class="avatar placeholder border-white">
                        <div class="bg-neutral-focus text-neutral-content w-8 h-8 text-[10px]">
                            <span>{{ getInitials(interviewer) }}</span>
                        </div>
                    </div>
                    <div v-if="getInterviewersForRound(round).length > 3" class="avatar placeholder border-white">
                        <div class="bg-gray-200 text-gray-600 w-8 h-8 text-[10px]">
                            <span>+{{ getInterviewersForRound(round).length - 3 }}</span>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</div>

                <div v-else class="flex flex-col items-center justify-center py-20 text-center animate-fade-in">
                    <div class="bg-gray-50 p-6 rounded-full mb-4">
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-10 w-10 text-gray-400" fill="none"
                            viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                        </svg>
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
                            <tr v-for="round in filteredRounds" :key="round.ID"
                                class="hover:bg-slate-50 transition-colors">
                                <td class="font-medium text-slate-700">
                                    <div class="flex flex-col">
                                        <span>{{ new Date(round.start_date_time).toLocaleDateString('th-TH', {
                                            dateStyle: 'medium' }) }}</span>
                                        <span class="text-xs text-gray-400 font-light">{{ new
                                            Date(round.start_date_time).toLocaleTimeString('th-TH', {hour: '2-digit',
                                            minute: '2-digit'}) }} - {{ new
                                                Date(round.end_date_time).toLocaleTimeString('th-TH', {hour: '2-digit',
                                            minute: '2-digit'}) }}</span>
                                    </div>
                                </td>
                                <td>
                                    <div class="font-bold text-[#1e3a8a]">{{ round.name }}</div>
                                </td>
                                <td>
                                    <div v-if="round.scholarship" class="badge badge-outline text-xs">{{
                                        round.scholarship.scholarship_name }}</div>
                                </td>
                                <td>
                                    <div class="flex items-center gap-2">
                                        <span class="text-sm font-semibold">{{ getBookedCount(round) }}/{{
                                            round.slots?.length || 0 }}</span>
                                        <progress class="progress progress-primary w-16 h-1.5"
                                            :value="getBookedCount(round)" :max="round.slots?.length || 1"></progress>
                                    </div>
                                </td>
                                <td>
                                    <div class="badge bg-gray-100 text-gray-500 border-none">เสร็จสิ้น</div>
                                </td>
                                <td class="text-right">
                                    <button @click="openViewModal(round)"
                                        class="btn btn-ghost btn-xs text-blue-600">ดูข้อมูล</button>
                                </td>
                            </tr>
                            <tr v-if="filteredRounds.length === 0">
                                <td colspan="6" class="text-center py-12 text-gray-400">
                                    ไม่พบประวัติย้อนหลัง
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>

        </div>

        <div v-if="isModalOpen"
            class="fixed inset-0 z-[100] flex items-center justify-center bg-black/60 backdrop-blur-sm p-4">
            <div
                class="bg-white w-full max-w-5xl h-[90vh] rounded-2xl shadow-2xl flex flex-col overflow-hidden animate-pop-in">

                <div class="px-6 py-4 border-b flex items-center justify-between bg-slate-50">
                    <div>
                        <h2 class="text-xl font-bold text-[#1e3a8a]">
                            {{ modalMode === 'create' ? 'สร้างรอบสัมภาษณ์ใหม่' : 'รายละเอียดรอบสัมภาษณ์' }}
                        </h2>
                        <p class="text-xs text-gray-500 mt-1">
                            {{ modalMode === 'create' ? 'กรอกข้อมูลเพื่อสร้าง Slot อัตโนมัติ' :
                            'ตรวจสอบรายชื่อผู้จองและจัดการข้อมูล' }}
                        </p>
                    </div>
                    <button @click="isModalOpen = false"
                        class="btn btn-circle btn-ghost btn-sm text-gray-500">✕</button>
                </div>

                <div class="flex-1 overflow-y-auto p-6 bg-[#f8fafc]">
                    <div v-if="isLoading && modalMode === 'view'" class="flex justify-center items-center h-full">
                        <span class="loading loading-spinner loading-lg"></span>
                    </div>
                    <div v-else class="grid grid-cols-1 lg:grid-cols-3 gap-6">
                        <div class="lg:col-span-1 space-y-6">
                            <div class="card bg-white shadow-sm border border-gray-100 p-5 space-y-4">
                                <h3 class="font-bold text-slate-700 border-b pb-2">ข้อมูลทั่วไป</h3>
                                <div class="form-control w-full">
                                    <label class="label"><span
                                            class="label-text font-medium">ชื่อรอบสัมภาษณ์</span></label>
                                    <input v-model="formData.name" type="text" placeholder="เช่น รอบเรียนดี 1/67"
                                        class="input input-bordered input-sm w-full" :disabled="modalMode === 'view'" />
                                </div>
                                <div class="form-control w-full">
                                    <label class="label"><span class="label-text font-medium">สำหรับทุน</span></label>
                                    <select v-model="formData.scholarship_id"
                                        class="select select-bordered select-sm w-full" :disabled="modalMode === 'view'">
                                        <option disabled :value="null">เลือกทุน...</option>
                                        <option v-for="s in scholarships" :key="s.ID" :value="s.ID">{{
                                            s.scholarship_name }}</option>
                                    </select>
                                </div>
                                <div class="form-control w-full">
                                    <label class="label"><span
                                            class="label-text font-medium">รายละเอียดเพิ่มเติม</span></label>
                                    <textarea v-model="formData.description"
                                        class="textarea textarea-bordered text-sm h-20"
                                        placeholder="เช่น สิ่งที่ต้องเตรียมมา..."
                                        :disabled="modalMode === 'view'"></textarea>
                                </div>
                            </div>
                            <div class="card bg-white shadow-sm border border-gray-100 p-5 space-y-4">
                                <h3 class="font-bold text-slate-700 border-b pb-2">เวลา & สถานที่</h3>
                                <div class="form-control w-full">
                                    <label class="label"><span class="label-text font-medium">วันที่เริ่ม</span></label>
                                    <input v-model="formData.date" type="date"
                                        class="input input-bordered input-sm w-full" :disabled="modalMode === 'view'" />
                                </div>
                                <div class="form-control w-full">
                                    <label class="label"><span class="label-text font-medium">วันที่สิ้นสุด</span></label>
                                    <input v-model="formData.end_date" type="date"
                                        class="input input-bordered input-sm w-full" :disabled="modalMode === 'view'" />
                                </div>
                                <div class="grid grid-cols-2 gap-2">
                                    <div class="form-control">
                                        <label class="label"><span class="label-text font-medium">เวลาที่เริ่ม</span></label>
                                        <input v-model="formData.start_time" type="time"
                                            class="input input-bordered input-sm w-full"
                                            :disabled="modalMode === 'view'" />
                                    </div>
                                    <div class="form-control">
                                        <label class="label"><span class="label-text font-medium">เวลาที่สิ้นสุด</span></label>
                                        <input v-model="formData.end_time" type="time"
                                            class="input input-bordered input-sm w-full"
                                            :disabled="modalMode === 'view'" />
                                    </div>
                                </div>
                                <div class="form-control w-full">
                                    <label class="label"><span class="label-text font-medium">ระยะเวลา/คน(นาที)</span></label>
                                    <select v-model="formData.slot_duration"
                                        class="select select-bordered select-sm w-full" :disabled="modalMode === 'view'">
                                        <option :value="10">10 นาที</option>
                                        <option :value="15">15 นาที</option>
                                        <option :value="20">20 นาที</option>
                                        <option :value="30">30 นาที</option>
                                        <option :value="60">1 ชั่วโมง</option>
                                    </select>
                                    <div class="divider my-2"></div>

                                    <div class="form-control w-full">
                                        <label class="label"><span class="label-text font-medium">รูปแบบ</span></label>
                                        <select v-model="formData.interview_mode_id"
                                            class="select select-bordered select-sm w-full" :disabled="modalMode === 'view'">
                                            <option disabled :value="null">เลือกรูปแบบ...</option>
                                            <option v-for="mode in interviewModes" :key="mode.ID" :value="mode.ID">{{
                                                mode.name }}</option>
                                        </select>
                                    </div>

                                    <div v-if="selectedModeName === 'Onsite'" class="form-control w-full animate-fade-in">
                                        <label class="label"><span class="label-text font-medium">สถานที่</span></label>
                                        <select v-model="formData.location_id"
                                            class="select select-bordered select-sm w-full"
                                            :disabled="modalMode === 'view'">
                                            <option :value="null" disabled>เลือกห้องสัมภาษณ์...</option>
                                            <option v-for="loc in locations" :key="loc.ID" :value="loc.ID">{{ loc.name
                                                }}</option>
                                        </select>
                                    </div>

                                    <div v-if="selectedModeName === 'Online'" class="form-control w-full animate-fade-in">
                                        <label class="label"><span class="label-text font-medium">Meeting
                                                Link</span></label>
                                        <input v-model="formData.meeting_link" type="text"
                                            placeholder="https://zoom.us/..."
                                            class="input input-bordered input-sm w-full"
                                            :disabled="modalMode === 'view'" />
                                    </div>
                                </div>
                            </div>
                        </div>
                        <div class="lg:col-span-2 space-y-6">
                            <div class="card bg-white shadow-sm border border-gray-100 p-5">
                                <h3 class="font-bold text-slate-700 mb-3">กรรมการสัมภาษณ์</h3>
                                <div class="flex flex-wrap gap-2">
                                    <label v-for="interviewer in interviewersList" :key="interviewer.ID"
                                        class="cursor-pointer border rounded-lg px-3 py-2 flex items-center gap-2 transition-all hover:bg-gray-50"
                                        :class="formData.interviewer_ids.includes(interviewer.ID) ? 'border-blue-500 bg-blue-50 ring-1 ring-blue-500' : 'border-gray-200'">
                                        <input type="checkbox" :value="interviewer.ID"
                                            v-model="formData.interviewer_ids"
                                            class="checkbox checkbox-primary checkbox-xs"
                                            :disabled="modalMode === 'view'" />
                                        <span class="text-sm">{{ interviewer.interviewer_firstname }} {{
                                            interviewer.interviewer_lastname }}</span>
                                    </label>
                                </div>
                            </div>
                            <div class="card bg-white shadow-sm border border-gray-100 p-5 flex-1 min-h-[400px]">
                                <div class="flex items-center justify-between mb-4">
                                    <h3 class="font-bold text-slate-700 flex items-center gap-2">
                                        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-gray-500"
                                            fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                                d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                                        </svg>
                                        {{ modalMode === 'create' ? 'พรีวิวตารางเวลา' : 'ตารางการจอง' }}
                                    </h3>
                                    <div class="text-xs text-gray-500">
                                        ทั้งหมด: {{ modalMode === 'create' ? previewSlots.length :
                                        selectedRoundDetails?.slots?.length || 0 }} สล็อต
                                    </div>
                                </div>
                                <div v-if="modalMode === 'create'">
                                    <div v-if="previewSlots.length > 0"
                                        class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-3">
                                        <div v-for="(slot, idx) in previewSlots" :key="idx"
                                            class="relative p-3 rounded-xl border text-center bg-white border-green-200 cursor-default">
                                            <span class="font-bold text-lg text-slate-700">{{ slot.time }}</span>
                                            <div class="badge badge-success badge-xs text-white">ว่าง</div>
                                        </div>
                                    </div>
                                    <div v-else
                                        class="flex flex-col items-center justify-center h-40 text-gray-400 border-2 border-dashed border-gray-200 rounded-xl">
                                        <p>กรุณาระบุเวลาเริ่ม-จบ และระยะเวลา</p>
                                        <p class="text-xs">ระบบจะคำนวณรอบให้อัตโนมัติ</p>
                                    </div>
                                </div>
                                <div v-else>
                                    <div v-if="selectedRoundDetails && selectedRoundDetails.slots"
                                        class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-3">
                                        <div v-for="slot in selectedRoundDetails.slots" :key="slot.ID"
                                            class="relative p-3 rounded-xl border text-center transition-all duration-200 flex flex-col items-center justify-center gap-1 group"
                                            :class="[
                                                slot.is_booked
                                                    ? 'bg-gray-100 border-gray-200 cursor-help'
                                                    : 'bg-white border-green-200'
                                            ]">
                                            <span class="font-bold text-lg"
                                                :class="slot.is_booked ? 'text-gray-400' : 'text-slate-700'">
                                                {{ new Date(slot.start_time).toLocaleTimeString('th-TH', {
                                                    hour:
                                                '2-digit', minute: '2-digit' }) }}
                                            </span>
                                            <div v-if="slot.is_booked" class="badge badge-ghost badge-xs text-gray-400">
                                                จองแล้ว</div>
                                            <div v-else class="badge badge-success badge-xs text-white">ว่าง</div>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>

                <div class="p-4 border-t bg-white flex justify-end gap-2">
                    <button @click="isModalOpen = false" class="btn btn-ghost text-gray-500">ยกเลิก</button>
                    <button v-if="modalMode === 'create'" @click="saveRound"
                        class="btn bg-[#1e3a8a] text-white hover:bg-[#152c6f]">
                        ยืนยันการสร้าง
                    </button>
                </div>

            </div>
        </div>

    </div>
</template>

<style scoped>
.animate-pop-in {
    animation: pop-in 0.2s cubic-bezier(0.16, 1, 0.3, 1) forwards;
}

.animate-fade-in {
    animation: fade 0.3s ease forwards;
}

@keyframes pop-in {
    0% {
        opacity: 0;
        transform: scale(0.95);
    }

    100% {
        opacity: 1;
        transform: scale(1);
    }
}

@keyframes fade {
    from {
        opacity: 0;
        transform: translateY(-5px);
    }

    to {
        opacity: 1;
        transform: translateY(0);
    }
}

/* Custom Scrollbar */
::-webkit-scrollbar {
    width: 6px;
}

::-webkit-scrollbar-track {
    background: transparent;
}

::-webkit-scrollbar-thumb {
    background: #cbd5e1;
    border-radius: 10px;
}

::-webkit-scrollbar-thumb:hover {
    background: #94a3b8;
}
</style>