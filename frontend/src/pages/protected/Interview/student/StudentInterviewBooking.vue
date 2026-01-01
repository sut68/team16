<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { InterviewAPI } from '@/services/api';
import { LocationAPI } from '@/services/api/location'; 
import { getMyProfile } from '@/services/api/user';
import { getStudentApplications } from '@/services/api/application';
import type { InterviewRound, Slot, InterviewBooking, Location, InterviewMode } from '@/interfaces/interview'; 
import type { MyProfileResponse } from '@/interfaces/user';
import type { ApplicationScholarshipResponse } from '@/interfaces/application_scholarship';
import Swal from 'sweetalert2';

// Interface เสริมสำหรับแสดงผล
interface InterviewBookingWithDetails extends InterviewBooking {
    slot?: Slot;
    round?: InterviewRound;
}

// --- 1. State Management ---
const isLoading = ref(true);
const studentProfile = ref<MyProfileResponse | null>(null);
const allRounds = ref<InterviewRound[]>([]);
const myBookings = ref<InterviewBookingWithDetails[]>([]);
const studentApplications = ref<ApplicationScholarshipResponse[]>([]);
const locations = ref<Location[]>([]);
const interviewModes = ref<InterviewMode[]>([]);

const isBookingModalOpen = ref(false);
const selectedRound = ref<InterviewRound | null>(null);
const selectedSlot = ref<Slot | null>(null);

// --- 2. Data Fetching ---
onMounted(async () => {
    await fetchData();
});

const fetchData = async () => {
    isLoading.value = true;
    try {
        const profileRes = await getMyProfile();
        studentProfile.value = profileRes;

        if (profileRes && profileRes.role && profileRes.role === 'student') {
            const studentId = profileRes.data.ID;

            const [appsRes, roundsRes, bookingRes, locationsRes, modesRes] = await Promise.all([
                getStudentApplications(studentId),
                InterviewAPI.getAllRounds(),
                InterviewAPI.getStudentBookings(studentId),
                LocationAPI.getAllLocations(),
                InterviewAPI.getAllModes()
            ]);

            studentApplications.value = appsRes || [];
            allRounds.value = roundsRes || [];
            locations.value = locationsRes || [];
            interviewModes.value = modesRes || [];

            if (bookingRes && bookingRes.length > 0) {
                myBookings.value = bookingRes.map(booking => {
                    const slot = allRounds.value
                        .flatMap(r => r.slots || [])
                        .find(s => s.ID === booking.slot_id);
                    
                    const round = allRounds.value.find(r => r.slots?.some(s => s.ID === booking.slot_id));

                    return { ...booking, slot, round };
                }).filter(b => b.slot && b.round); // Ensure booking is valid
            } else {
                myBookings.value = [];
            }
        }
    } catch (error) {
        console.error("Failed to fetch initial data:", error);
        Swal.fire('ผิดพลาด', 'ไม่สามารถโหลดข้อมูลได้', 'error');
    } finally {
        isLoading.value = false;
    }
};

// --- 3. Computed Logic ---

const qualifiedScholarshipIds = computed(() => {
    return studentApplications.value
        .filter(app => app.status === 'qualified')
        .map(app => app.scholarship_id);
});

// Scholarship IDs for which an interview has already been booked.
const bookedScholarshipIds = computed(() => {
    return myBookings.value.map(booking => booking.round?.scholarship_id).filter(id => id !== undefined);
});


const displayedRounds = computed(() => {
    if (!allRounds.value) return [];
    // Show rounds for scholarships that are qualified AND not yet booked.
    return allRounds.value.filter(round =>
        qualifiedScholarshipIds.value.includes(round.scholarship_id) &&
        !bookedScholarshipIds.value.includes(round.scholarship_id)
    );
});

const availableSlots = computed(() => {
    if (!selectedRound.value || !selectedRound.value.slots) return [];
    return [...selectedRound.value.slots].sort((a, b) =>
        new Date(a.start_time).getTime() - new Date(b.start_time).getTime()
    );
});

const getBookedCount = (round: InterviewRound) => {
    if (!round.slots) return 0;
    return round.slots.filter(s => s.is_booked).length;
}

// --- Helper Functions for UI ---

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

const getRoundStatus = (round: InterviewRound) => {
    if (!round.slots || round.slots.length === 0) return 'Open';
    const isFull = round.slots.every(slot => slot.is_booked);
    return isFull ? 'Full' : 'Open';
};

// --- 4. Methods ---
const openBookingModal = async (round: InterviewRound) => {
    isLoading.value = true;
    try {
        const freshRound = await InterviewAPI.getRoundById(round.ID);
        selectedRound.value = freshRound;
        selectedSlot.value = null;
        isBookingModalOpen.value = true;
    } catch (error) {
        Swal.fire('ผิดพลาด', 'ไม่สามารถโหลดข้อมูลรอบสัมภาษณ์ได้', 'error');
    } finally {
        isLoading.value = false;
    }
};

const selectSlot = (slot: Slot) => {
    if (!slot.is_booked) {
        selectedSlot.value = slot;
    }
};

const confirmBooking = async () => {
    if (!selectedRound.value || !selectedSlot.value) return;

    const app = studentApplications.value.find(a => a.scholarship_id === selectedRound.value?.scholarship_id);

    if (!app) {
        Swal.fire('ผิดพลาด', 'ไม่พบใบสมัครสำหรับทุนนี้', 'error');
        return;
    }

    const timeDisplay = new Date(selectedSlot.value.start_time).toLocaleTimeString('th-TH', { hour: '2-digit', minute: '2-digit' });

    const result = await Swal.fire({
        title: 'ยืนยันการจอง',
        text: `คุณต้องการจองเวลา ${timeDisplay} ใช่หรือไม่?`,
        icon: 'question',
        showCancelButton: true,
        confirmButtonText: 'ยืนยัน',
        confirmButtonColor: '#1e3a8a',
        cancelButtonText: 'ยกเลิก',
    });

    if (result.isConfirmed) {
        isLoading.value = true;
        try {
            const payload = {
                slot_id: selectedSlot.value.ID,
                application_scholarship_id: app.ID,
                status: 'confirmed',
            };

            await InterviewAPI.createBooking(payload);

            await Swal.fire('สำเร็จ!', 'การจองของคุณถูกยืนยันแล้ว', 'success');

            isBookingModalOpen.value = false;
            await fetchData();
        } catch (error: any) {
            console.error("Booking failed:", error);
            const msg = error.response?.data?.error || 'การจองล้มเหลว';
            Swal.fire('ผิดพลาด', msg, 'error');
        } finally {
            isLoading.value = false;
        }
    }
};

const cancelBooking = async (bookingId: number) => {
    const result = await Swal.fire({
        title: 'ยืนยันการยกเลิก',
        text: 'เมื่อยกเลิกแล้ว สิทธิ์จะว่างให้ผู้อื่นจองทันที',
        icon: 'warning',
        showCancelButton: true,
        confirmButtonText: 'ใช่, ยกเลิก',
        confirmButtonColor: '#d33',
        cancelButtonText: 'ไม่',
    });

    if (result.isConfirmed) {
        isLoading.value = true;
        try {
            await InterviewAPI.deleteBooking(bookingId);
            Swal.fire('สำเร็จ', 'ยกเลิกการจองเรียบร้อยแล้ว', 'success');
            await fetchData(); // Refetch all data to update UI
        } catch (error) {
            console.error("Cancellation failed:", error);
            Swal.fire('ผิดพลาด', 'การยกเลิกล้มเหลว', 'error');
        } finally {
            isLoading.value = false;
        }
    }
};
</script>

<template>
    <div class="w-full mx-auto flex flex-col h-full p-6 bg-white rounded-tl-[30px] shadow overflow-visible"
        data-theme="light">

        <div v-if="isLoading" class="flex justify-center items-center h-full">
            <span class="loading loading-spinner loading-lg"></span>
        </div>

        <div v-else class="animate-fade-in">
            <div class="mb-6">
                <h1 class="text-2xl font-bold text-[#1e3a8a] mb-1">จองรอบสัมภาษณ์</h1>
                <p class="text-gray-500 text-sm" v-if="studentProfile && 'first_name_th' in studentProfile.data">{{
                    (studentProfile.data as any).first_name_th }} {{ (studentProfile.data as any).last_name_th }}</p>
            </div>

            <div class="mb-8 bg-blue-50/50 border border-blue-100 rounded-xl p-4 flex items-start gap-3">
                <div class="bg-blue-100 p-2 rounded-lg text-blue-600 mt-1">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                        <path fill-rule="evenodd"
                            d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z"
                            clip-rule="evenodd" />
                    </svg>
                </div>
                <div>
                    <h4 class="font-bold text-[#1e3a8a] text-sm">สิทธิ์การจองของคุณ</h4>
                    <p class="text-sm text-slate-600 mt-1">
                        ระบบแสดงเฉพาะรอบของทุนที่คุณ <u>ผ่านการพิจารณา</u> และได้รับสิทธิ์เข้าสัมภาษณ์แล้ว
                    </p>
                </div>
            </div>

            <div v-if="myBookings.length > 0" class="mb-8 animate-fade-in space-y-4">
                 <h3 class="font-bold text-slate-700 flex items-center gap-2">
                    นัดหมายของคุณ <span class="badge badge-accent badge-sm">{{ myBookings.length }}</span>
                </h3>
                <div v-for="booking in myBookings" :key="booking.ID"
                    class="card bg-gradient-to-r from-[#1e3a8a] to-[#2563eb] text-white shadow-lg relative overflow-hidden">
                    <div
                        class="card-body relative z-10 flex flex-col md:flex-row justify-between items-center gap-6 p-6">
                        <div class="flex items-center gap-4 flex-1">
                            <div
                                class="w-12 h-12 bg-white/20 rounded-full flex items-center justify-center backdrop-blur-sm shrink-0">
                                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-white" fill="none"
                                    viewBox="0 0 24 24" stroke="currentColor">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                        d="M5 13l4 4L19 7" />
                                </svg>
                            </div>
                            <div>
                                <h2 class="font-bold text-lg leading-tight">{{ booking.round?.scholarship?.scholarship_name }}</h2>
                                <p class="text-blue-100 text-sm">{{ booking.round?.name }}</p>
                            </div>
                        </div>
                        <div v-if="booking.slot"
                            class="bg-white/10 rounded-lg p-3 backdrop-blur-md border border-white/20 text-sm w-full md:w-auto md:min-w-[280px]">
                            <div class="grid grid-cols-1 sm:grid-cols-2 gap-x-4 gap-y-1">
                                <p><span class="opacity-70 text-xs mr-2">วันที่:</span>{{ new
                                    Date(booking.slot.start_time).toLocaleDateString('th-TH', { dateStyle: 'medium' })
                                }}</p>
                                <p><span class="opacity-70 text-xs mr-2">เวลา:</span>{{ new
                                    Date(booking.slot.start_time).toLocaleTimeString('th-TH', {
                                        hour: '2-digit',
                                        minute: '2-digit'
                                    }) }} น.</p>
                            </div>
                        </div>
                        <button @click="cancelBooking(booking.ID)"
                            class="btn btn-sm btn-circle btn-ghost text-white hover:bg-white/20 tooltip tooltip-left"
                            data-tip="ยกเลิก">✕</button>
                    </div>
                </div>
            </div>

            <div class="flex-1">
                <h3 v-if="displayedRounds.length > 0" class="font-bold text-slate-700 mb-4 flex items-center gap-2">
                    รอบที่เปิดให้จอง <span class="badge badge-ghost badge-sm">{{ displayedRounds.length }}</span>
                </h3>

                <div v-if="displayedRounds.length > 0"
                    class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6 animate-fade-in">
                    <div v-for="round in displayedRounds" :key="round.ID"
                        class="group card bg-white border border-gray-200 hover:border-blue-400 hover:shadow-xl transition-all duration-300 rounded-2xl overflow-visible">

                        <div
                            class="px-5 py-4 border-b border-gray-50 bg-slate-50/50 flex flex-col justify-start items-start rounded-t-2xl relative">
                            <div class="flex gap-2 mb-2">
                                <div class="badge badge-sm font-medium border-none text-white shadow-sm" :class="{
                                    'badge-success': getRoundStatus(round) === 'Open',
                                    'badge-error': getRoundStatus(round) === 'Full'
                                }">
                                    {{ getRoundStatus(round) === 'Open' ? 'เปิดให้จอง' : 'เต็มแล้ว' }}
                                </div>
                                <div v-if="round.scholarship" class="badge badge-sm badge-outline text-gray-500">
                                    {{ round.scholarship.scholarship_name }}
                                </div>
                            </div>
                            <h3
                                class="font-bold text-[#1e3a8a] text-lg leading-tight group-hover:text-blue-600 transition-colors line-clamp-2">
                                {{ round.name }}
                            </h3>
                        </div>

                        <div class="p-5 space-y-4">
                            <div class="flex items-center gap-3 text-sm text-gray-600">
                                <div
                                    class="w-8 h-8 rounded-full bg-blue-50 flex items-center justify-center text-blue-600 shrink-0">
                                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none"
                                        viewBox="0 0 24 24" stroke="currentColor">
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                            d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                                    </svg>
                                </div>
                                <div>
                                    <p class="font-semibold text-slate-700">
                                        วันที่ {{ new Date(round.start_date_time).toLocaleDateString('th-TH', {
                                            day:
                                                'numeric', month: 'short', year: '2-digit'
                                        }) }}
                                    </p>
                                    <p class="text-xs">
                                        {{ new Date(round.start_date_time).toLocaleTimeString('th-TH', {
                                            hour: '2-digit',
                                            minute: '2-digit'
                                        }) }} -
                                        {{ new Date(round.end_date_time).toLocaleTimeString('th-TH', {
                                            hour: '2-digit',
                                            minute: '2-digit'
                                        }) }}
                                    </p>
                                </div>
                            </div>

                            <div class="flex items-center gap-3 text-sm text-gray-600">
                                <div class="w-8 h-8 rounded-full flex items-center justify-center shrink-0"
                                    :class="getModeName(round) === 'Onsite' ? 'bg-orange-50 text-orange-600' : 'bg-purple-50 text-purple-600'">
                                    <svg v-if="getModeName(round) === 'Onsite'" xmlns="http://www.w3.org/2000/svg"
                                        class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
                                        <path fill-rule="evenodd"
                                            d="M5.05 4.05a7 7 0 119.9 9.9L10 18.9l-4.95-4.95a7 7 0 010-9.9zM10 11a2 2 0 100-4 2 2 0 000 4z"
                                            clip-rule="evenodd" />
                                    </svg>
                                    <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20"
                                        fill="currentColor">
                                        <path
                                            d="M2 6a2 2 0 012-2h6a2 2 0 012 2v8a2 2 0 01-2 2H4a2 2 0 01-2-2V6zM14.553 7.106A1 1 0 0014 8v4a1 1 0 00.553.894l2 1A1 1 0 0018 13V7a1 1 0 00-1.447-.894l-2 1z" />
                                    </svg>
                                </div>
                                <div class="overflow-hidden">
                                    <p class="font-semibold text-slate-700">{{ getModeName(round) }}</p>
                                    <p class="text-xs truncate w-full" :title="getLocationLabel(round)">{{
                                        getLocationLabel(round) }}</p>
                                </div>
                            </div>

                            <div class="divider my-1"></div>

                            <div class="flex justify-between items-end gap-3">
                                <div class="flex-1">
                                    <div class="flex justify-between text-xs mb-1">
                                        <span class="text-gray-500">ว่าง ({{ round.slots?.length - getBookedCount(round)
                                        }} ที่นั่ง)</span>
                                        <span class="font-bold text-slate-700">{{ getBookedCount(round) }}/{{
                                            round.slots?.length || 0 }}</span>
                                    </div>
                                    <progress class="progress w-full h-2"
                                        :class="getBookedCount(round) === (round.slots?.length || 0) ? 'progress-error' : 'progress-primary'"
                                        :value="getBookedCount(round)" :max="round.slots?.length || 1"></progress>
                                </div>

                                <button @click="openBookingModal(round)"
                                    :disabled="getRoundStatus(round) === 'Full'"
                                    class="btn btn-sm px-4 bg-[#1e3a8a] border-none text-white hover:bg-[#152c6f] disabled:bg-gray-100 disabled:text-gray-400 shadow-sm whitespace-nowrap">
                                    {{ getRoundStatus(round) === 'Full' ? 'ที่นั่งเต็ม' : 'เลือกเวลา' }}
                                </button>
                            </div>
                        </div>
                    </div>
                </div>

                <div v-if="displayedRounds.length === 0 && myBookings.length === 0"
                    class="flex flex-col items-center justify-center py-12 text-center border-2 border-dashed border-gray-200 rounded-2xl bg-gray-50/50">
                    <p class="text-gray-500 font-medium">ไม่พบรอบสัมภาษณ์ที่สามารถจองได้</p>
                    <p class="text-xs text-gray-400 mt-1">คุณอาจยังไม่ได้รับสิทธิ์ หรือรอบยังไม่เปิด</p>
                </div>
            </div>
        <Teleport to="body">
            <div v-if="isBookingModalOpen && selectedRound" 
                 class="fixed inset-0 z-[100] flex items-start justify-center pt-[50px]  bg-black/60 backdrop-blur-sm p-4 transition-all duration-300">
    
                <div class="bg-white w-full h-full md:max-w-6xl md:h-[80vh] md:rounded-2xl shadow-2xl overflow-hidden flex flex-col md:flex-row animate-pop-in relative">
                    
                    

                    <div class="w-full md:w-[350px] bg-slate-50 border-r border-slate-200 flex flex-col shrink-0">
                        
                        <div class="p-5 bg-[#1e3a8a] text-white">
                            <div class="badge bg-orange-500 text-white border-none mb-2 text-xs">
                                {{ selectedRound.scholarship?.scholarship_name || 'ทุนการศึกษา' }}
                            </div>
                            <h2 class="text-xl font-bold leading-tight mb-1">{{ selectedRound.name }}</h2>
                            <div class="flex items-center gap-2 text-blue-100 text-xs">
                                <svg xmlns="http://www.w3.org/2000/svg" class="h-3 w-3" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M6 2a1 1 0 00-1 1v1H4a2 2 0 00-2 2v10a2 2 0 002 2h12a2 2 0 002-2V6a2 2 0 00-2-2h-1V3a1 1 0 10-2 0v1H7V3a1 1 0 00-1-1zm0 5a1 1 0 000 2h8a1 1 0 100-2H6z" clip-rule="evenodd" /></svg>
                                {{ new Date(selectedRound.start_date_time).toLocaleDateString('th-TH', { weekday: 'short', day: 'numeric', month: 'short', year: '2-digit'}) }}
                            </div>
                        </div>

                        <div class="p-5 flex-1 overflow-y-auto space-y-4">
                            <div class="group">
                                <h3 class="text-[15px] font-bold text-slate-400 uppercase tracking-wider mb-1">สถานที่ / รูปแบบ</h3>
                                <div class="flex items-start gap-2">
                                    <div class="w-8 h-8 rounded-full bg-blue-100 text-blue-600 flex items-center justify-center shrink-0 mt-0.5">
                                        <svg v-if="getModeName(selectedRound) === 'Onsite'" xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M5.05 4.05a7 7 0 119.9 9.9L10 18.9l-4.95-4.95a7 7 0 010-9.9zM10 11a2 2 0 100-4 2 2 0 000 4z" clip-rule="evenodd" /></svg>
                                        <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor"><path d="M2 6a2 2 0 012-2h6a2 2 0 012 2v8a2 2 0 01-2 2H4a2 2 0 01-2-2V6zM14.553 7.106A1 1 0 0014 8v4a1 1 0 00.553.894l2 1A1 1 0 0018 13V7a1 1 0 00-1.447-.894l-2 1z" /></svg>
                                    </div>
                                    <div>
                                        <p class="font-bold text-slate-700 text-[14px]">{{ getModeName(selectedRound) }}</p>
                                        <p class="text-slate-500 text-sm leading-relaxed">{{ getLocationLabel(selectedRound) }}</p>
                                        <p v-if="getModeName(selectedRound) === 'Online' && (selectedRound as any).meeting_link" class="text-[10px] text-blue-500 mt-1 truncate max-w-[150px]">Link จะแสดงเมื่อยืนยัน</p>
                                    </div>
                                </div>
                            </div>

                            <div class="divider my-1"></div>

                            <div>
                                <h3 class="text-[15px] font-bold text-slate-400 uppercase tracking-wider mb-1">หมายเหตุ</h3>
                                <div class="bg-white p-3 rounded-lg border border-slate-200 text-sm text-slate-600 leading-relaxed shadow-sm">
                                    {{ selectedRound.description || 'ไม่มีรายละเอียดเพิ่มเติม' }}
                                </div>
                            </div>
                        </div>
                    </div>

                    <div class="flex-1 bg-white flex flex-col h-full min-h-0">
                        <div class="md:hidden p-3 border-b text-center font-bold text-slate-700 text-sm">เลือกเวลาสัมภาษณ์</div>

                        <div class="flex-1 overflow-y-auto p-4 custom-scrollbar">
                            <div class="flex justify-between items-center mb-4">
                                <div>
                                    <h3 class="text-lg font-bold text-slate-800">เลือกช่วงเวลาที่ต้องการจองสัมภาษณ์</h3>
                                    <p class="text-slate-500 text-sm">เลือกช่วงเวลาที่คุณสะดวก</p>
                                </div>
                                <div class="hidden sm:flex gap-2 text-[13px]">
                                    <div class="flex items-center gap-1"><span class="w-2 h-2 rounded-full bg-white border border-green-400"></span> ว่าง</div>
                                    <div class="flex items-center gap-1"><span class="w-2 h-2 rounded-full bg-blue-50 border border-blue-600"></span> ที่เลือก</div>
                                    <div class="flex items-center gap-1"><span class="w-2 h-2 rounded-full bg-gray-100 border border-gray-200"></span> ไม่ว่าง</div>
                                </div>
                            </div>
                            
                            <div v-if="isLoading" class="flex flex-col justify-center items-center h-32 gap-2">
                                <span class="loading loading-spinner loading-md text-blue-600"></span>
                            </div>

                            <div v-else class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-3">
                                <button v-for="slot in availableSlots" :key="slot.ID"
                                    @click="!slot.is_booked && selectSlot(slot)"
                                    :disabled="slot.is_booked"
                                    class="relative group flex flex-col items-center justify-center p-2 rounded-xl border transition-all duration-200 h-20"
                                    :class="[
                                        slot.is_booked 
                                            ? 'bg-gray-50 border-gray-100 opacity-60 cursor-not-allowed' 
                                            : (selectedSlot?.ID === slot.ID 
                                                ? 'bg-blue-50 border-blue-600 shadow-sm scale-[1.02] z-10' 
                                                : 'bg-white border-slate-100 hover:border-blue-300 hover:shadow-md hover:-translate-y-0.5')
                                    ]">
                                    
                                    <div v-if="selectedSlot?.ID === slot.ID" class="absolute top-1 right-1 text-blue-600 animate-bounce-short">
                                        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" /></svg>
                                    </div>

                                    <span class="text-lg font-bold tracking-tight mb-1" 
                                        :class="selectedSlot?.ID === slot.ID ? 'text-blue-700' : (slot.is_booked ? 'text-gray-400' : 'text-slate-700')">
                                        {{ new Date(slot.start_time).toLocaleTimeString('th-TH', { hour: '2-digit', minute: '2-digit' }) }}
                                    </span>
                                    
                                    <span class="text-[12px] px-2 py-0.5 rounded-full font-medium"
                                        :class="[
                                            slot.is_booked ? 'bg-gray-200 text-gray-500' :
                                            (selectedSlot?.ID === slot.ID ? 'bg-blue-200 text-blue-800' : 'bg-green-100 text-green-700')
                                        ]">
                                        {{ slot.is_booked ? 'ไม่ว่าง' : (selectedSlot?.ID === slot.ID ? 'เลือกแล้ว' : 'ว่าง') }}
                                    </span>
                                </button>
                            </div>
                        </div>

                        <div class="p-4 border-t border-slate-100 bg-white flex flex-col sm:flex-row justify-between items-center gap-3 shadow-[0_-4px_6px_-1px_rgba(0,0,0,0.05)] z-20">
                            <div class="flex items-center gap-2 w-full sm:w-auto">
                                <div class="w-8 h-8 rounded-full bg-slate-100 flex items-center justify-center text-slate-400">
                                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
                                </div>
                                <div>
                                    <p class="text-[10px] text-slate-500">เวลาที่เลือก</p>
                                    <p class="font-bold text-slate-800 text-sm">
                                        {{ selectedSlot ? new Date(selectedSlot.start_time).toLocaleTimeString('th-TH', { hour: '2-digit', minute: '2-digit' }) + ' น.' : '-' }}
                                    </p>
                                </div>
                            </div>

                            <div class="flex gap-2 w-full sm:w-auto text-sm">
                                <button @click="isBookingModalOpen = false" class="btn btn-sm btn-ghost text-slate-500 hover:bg-slate-100 flex-1 sm:flex-none">ยกเลิก</button>
                                <button @click="confirmBooking" 
                                    :disabled="!selectedSlot" 
                                    class="btn btn-sm bg-[#1e3a8a] text-white hover:bg-[#152c6f] border-none px-6 flex-1 sm:flex-none shadow-md shadow-blue-900/10 transition-all hover:scale-105 disabled:bg-slate-200 disabled:text-slate-400 disabled:shadow-none disabled:scale-100">
                                    ยืนยันการจอง
                                </button>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
            </Teleport>
        </div>
    </div>
</template>

<style scoped>
.animate-pop-in {
    animation: pop-in 0.2s cubic-bezier(0.16, 1, 0.3, 1) forwards;
}

.animate-fade-in {
    animation: fade 0.4s ease forwards;
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
        transform: translateY(10px);
    }

    to {
        opacity: 1;
        transform: translateY(0);
    }
}

/* Animation สำหรับ Modal */
.animate-pop-in {
    animation: pop-in 0.3s cubic-bezier(0.16, 1, 0.3, 1) forwards;
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

/* Animation เด้งดึ๋งเล็กน้อยตอนกดเลือก */
.animate-bounce-short {
    animation: bounce-short 0.4s ease-in-out;
}

@keyframes bounce-short {

    0%,
    100% {
        transform: translateY(0);
    }

    50% {
        transform: translateY(-20%);
    }
}

/* Scrollbar สวยๆ */
.custom-scrollbar::-webkit-scrollbar {
    width: 6px;
}

.custom-scrollbar::-webkit-scrollbar-track {
    background: transparent;
}

.custom-scrollbar::-webkit-scrollbar-thumb {
    background: #cbd5e1;
    border-radius: 10px;
}

.custom-scrollbar::-webkit-scrollbar-thumb:hover {
    background: #94a3b8;
}
</style>