<script setup lang="ts">
import { ref, computed } from 'vue';

// --- 1. Interfaces ---
interface InterviewRound {
  ID: number;
  name: string;
  scholarship_name: string;
  date: string;
  start_time: string;
  end_time: string;
  slot_duration: number;
  mode: 'Online' | 'Onsite';
  location_name?: string;
  meeting_link?: string;
  interviewers: string[];
  total_slots: number;
  booked_slots: number;
  status: 'Open' | 'Closed' | 'Full';
}

interface TimeSlot {
  time: string;
  isBooked: boolean;
  isSelected: boolean;
}

interface StudentBooking {
  round_id: number;
  round_name: string;
  slot_time: string;
  location: string;
  date: string;
  mode: 'Online' | 'Onsite';
}

// --- 2. State & Mock Data ---

// Mock User Profile (นักศึกษาที่ผ่านเกณฑ์ทุนเรียนดีเท่านั้น)
const currentUser = ref({
  id: 'B6400001',
  name: 'นายสมชาย รักเรียน',
  passed_scholarships: ['ทุนเรียนดี'] 
});

const isBookingModalOpen = ref(false);
const selectedRound = ref<InterviewRound | null>(null);
const selectedTimeSlot = ref<string | null>(null);
const myBooking = ref<StudentBooking | null>(null);

// Mock Rounds Data (รวมทุกทุน)
const rounds = ref<InterviewRound[]>([
  {
    ID: 1,
    name: 'สัมภาษณ์ทุนเรียนดี รุ่นที่ 1/2567',
    scholarship_name: 'ทุนเรียนดี',
    date: '2024-12-25',
    start_time: '09:00',
    end_time: '12:00',
    slot_duration: 30,
    mode: 'Onsite',
    location_name: 'อาคารเรียนรวม 2 ห้อง 201',
    interviewers: ['ดร.สมชาย'],
    total_slots: 6,
    booked_slots: 4,
    status: 'Open'
  },
  {
    ID: 2,
    name: 'สัมภาษณ์ทุนขาดแคลน (รอบเก็บตก)',
    scholarship_name: 'ทุนขาดแคลนทุนทรัพย์',
    date: '2024-12-28',
    start_time: '13:00',
    end_time: '16:00',
    slot_duration: 20,
    mode: 'Online',
    interviewers: ['อ.วิจัย'],
    total_slots: 9,
    booked_slots: 9,
    status: 'Full'
  },
  {
    ID: 3,
    name: 'สัมภาษณ์ทุนเรียนดี (เพิ่มเติม)', // อันนี้ทุนเรียนดีอีกอัน
    scholarship_name: 'ทุนเรียนดี',
    date: '2024-12-26',
    start_time: '09:00',
    end_time: '12:00',
    slot_duration: 30,
    mode: 'Onsite',
    location_name: 'อาคารเรียนรวม 2 ห้อง 202',
    interviewers: ['อ.ใจดี'],
    total_slots: 6,
    booked_slots: 0,
    status: 'Open'
  }
]);

// --- 3. Computed Logic (Simplified) ---

const displayedRounds = computed(() => {
  // กรอง 2 ชั้น:
  // 1. ต้องไม่ปิด (Closed)
  // 2. ต้องเป็นทุนที่ user มีสิทธิ์ (passed_scholarships)
  return rounds.value.filter(r => 
    r.status !== 'Closed' && 
    currentUser.value.passed_scholarships.includes(r.scholarship_name)
  );
});

const availableSlots = computed<TimeSlot[]>(() => {
  if (!selectedRound.value) return [];
  const slots: TimeSlot[] = [];
  const { start_time, end_time, slot_duration } = selectedRound.value;
  const [startH = 0, startM = 0] = start_time.split(':').map(Number);
  const [endH = 0, endM = 0] = end_time.split(':').map(Number);
  
  let currentMin = startH * 60 + startM;
  const endMin = endH * 60 + endM;
  let mockSeed = selectedRound.value.ID; 

  while (currentMin + slot_duration <= endMin) {
    const h = Math.floor(currentMin / 60);
    const m = currentMin % 60;
    const timeStr = `${String(h).padStart(2, '0')}:${String(m).padStart(2, '0')}`;
    mockSeed = (mockSeed * 9301 + 49297) % 233280;
    const isBookedMock = selectedRound.value.status === 'Full' ? true : (mockSeed % 10 < 4);

    slots.push({
      time: timeStr,
      isBooked: isBookedMock,
      isSelected: selectedTimeSlot.value === timeStr
    });
    currentMin += slot_duration;
  }
  return slots;
});

// --- 4. Methods ---
const openBookingModal = (round: InterviewRound) => {
  if (myBooking.value) {
    alert('กรุณายกเลิกการจองเดิมก่อน');
    return;
  }
  selectedRound.value = round;
  selectedTimeSlot.value = null;
  isBookingModalOpen.value = true;
};

const selectSlot = (slot: TimeSlot) => {
  if (!slot.isBooked) selectedTimeSlot.value = slot.time;
};

const confirmBooking = () => {
  if (!selectedRound.value || !selectedTimeSlot.value) return;
  if (confirm(`ยืนยันการจองเวลา ${selectedTimeSlot.value}?`)) {
    myBooking.value = {
      round_id: selectedRound.value.ID,
      round_name: selectedRound.value.name,
      slot_time: selectedTimeSlot.value,
      location: selectedRound.value.mode === 'Onsite' ? selectedRound.value.location_name! : 'Online',
      date: selectedRound.value.date,
      mode: selectedRound.value.mode
    };
    const idx = rounds.value.findIndex(r => r.ID === selectedRound.value!.ID);
    const targetRound = rounds.value[idx];
    if (targetRound) targetRound.booked_slots++;
    isBookingModalOpen.value = false;
  }
};

const cancelBooking = () => {
  if(confirm('ต้องการยกเลิกการจอง?')) {
    if (myBooking.value) {
        const idx = rounds.value.findIndex(r => r.ID === myBooking.value?.round_id);
        const targetRound = rounds.value[idx];
        if (targetRound) targetRound.booked_slots--;
    }
    myBooking.value = null;
  }
};
</script>

<template>
  <div class="w-full mx-auto flex flex-col h-full p-6 bg-white rounded-tl-[30px] shadow overflow-visible" data-theme="light">
    
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-[#1e3a8a] mb-1">จองรอบสัมภาษณ์</h1>
      <p class="text-gray-500 text-sm">ยินดีต้อนรับ, {{ currentUser.name }}</p>
    </div>

    <div class="mb-8 bg-blue-50/50 border border-blue-100 rounded-xl p-4 flex items-start gap-3">
      <div class="bg-blue-100 p-2 rounded-lg text-blue-600 mt-1">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
          <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd" />
        </svg>
      </div>
      <div>
        <h4 class="font-bold text-[#1e3a8a] text-sm">สิทธิ์การจองของคุณ</h4>
        <p class="text-sm text-slate-600 mt-1">
          ระบบแสดงเฉพาะรอบของทุนที่คุณ <u>ผ่านการพิจารณา</u> แล้ว:
        </p>
        <div class="flex flex-wrap gap-2 mt-2">
          <span v-for="sc in currentUser.passed_scholarships" :key="sc" 
            class="badge badge-primary badge-outline bg-white border-blue-200 text-blue-700">
            {{ sc }}
          </span>
        </div>
      </div>
    </div>

    <div v-if="myBooking" class="mb-8 animate-fade-in">
      <div class="card bg-gradient-to-r from-[#1e3a8a] to-[#2563eb] text-white shadow-lg relative overflow-hidden">
         <div class="card-body relative z-10 flex flex-col md:flex-row justify-between items-center gap-6 p-6">
            <div class="flex items-center gap-4">
                <div class="w-12 h-12 bg-white/20 rounded-full flex items-center justify-center backdrop-blur-sm">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" /></svg>
                </div>
                <div>
                   <h2 class="font-bold text-lg">จองสำเร็จ</h2>
                   <p class="text-blue-100 text-sm">นัดหมายการสัมภาษณ์ของคุณ</p>
                </div>
            </div>
            <div class="flex-1 bg-white/10 rounded-lg p-3 backdrop-blur-md border border-white/20 text-sm w-full">
               <div class="grid grid-cols-1 sm:grid-cols-2 gap-x-4 gap-y-1">
                  <p><span class="opacity-70 text-xs mr-2">รอบ:</span>{{ myBooking.round_name }}</p>
                  <p><span class="opacity-70 text-xs mr-2">วันที่:</span>{{ new Date(myBooking.date).toLocaleDateString('th-TH') }}</p>
                  <p><span class="opacity-70 text-xs mr-2">เวลา:</span>{{ myBooking.slot_time }} น.</p>
                  <p><span class="opacity-70 text-xs mr-2">ที่:</span>{{ myBooking.location }}</p>
               </div>
            </div>
            <button @click="cancelBooking" class="btn btn-sm btn-circle btn-ghost text-white hover:bg-white/20 tooltip tooltip-left" data-tip="ยกเลิก">✕</button>
         </div>
      </div>
    </div>

    <div class="flex-1">
        <h3 v-if="!myBooking" class="font-bold text-slate-700 mb-4 flex items-center gap-2">
           รอบที่เปิดให้จอง <span class="badge badge-ghost badge-sm">{{ displayedRounds.length }}</span>
        </h3>

        <div v-if="displayedRounds.length > 0" class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-5 animate-fade-in">
          <div v-for="round in displayedRounds" :key="round.ID" 
            class="card bg-white border border-gray-200 shadow-sm hover:shadow-md transition-all rounded-xl">
            
            <div class="p-5">
               <div class="flex justify-between items-start mb-3">
                  <span class="badge badge-sm badge-ghost text-gray-500">{{ round.scholarship_name }}</span>
                  <span class="text-xs font-semibold" :class="round.status === 'Open' ? 'text-green-600' : 'text-red-500'">
                      {{ round.status === 'Open' ? 'ว่าง' : 'เต็ม' }}
                  </span>
               </div>
               
               <h3 class="font-bold text-[#1e3a8a] mb-4 line-clamp-2 min-h-[3rem]">{{ round.name }}</h3>

               <div class="space-y-2 text-sm text-gray-600 mb-4">
                  <div class="flex items-center gap-2">
                     <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" /></svg>
                     <span>{{ new Date(round.date).toLocaleDateString('th-TH', { dateStyle: 'medium'}) }}</span>
                  </div>
                  <div class="flex items-center gap-2">
                     <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
                     <span>{{ round.start_time }} - {{ round.end_time }}</span>
                  </div>
               </div>

               <div class="w-full bg-gray-100 rounded-full h-1.5 mb-2">
                  <div class="bg-blue-600 h-1.5 rounded-full" :style="{ width: (round.booked_slots / round.total_slots * 100) + '%' }"></div>
               </div>
               <div class="text-xs text-right text-gray-400 mb-4">
                  ว่าง {{ round.total_slots - round.booked_slots }} ที่นั่ง
               </div>

               <button 
                  @click="openBookingModal(round)"
                  :disabled="round.status === 'Full' || (myBooking !== null)"
                  class="btn btn-sm w-full bg-[#1e3a8a] border-none text-white hover:bg-[#152c6f] disabled:bg-gray-100 disabled:text-gray-400">
                  {{ myBooking ? 'จองแล้ว' : (round.status === 'Full' ? 'ที่นั่งเต็ม' : 'เลือก') }}
               </button>
            </div>
          </div>
        </div>

        <div v-else class="flex flex-col items-center justify-center py-12 text-center border-2 border-dashed border-gray-200 rounded-2xl bg-gray-50/50">
           <p class="text-gray-500 font-medium">ไม่พบรอบสัมภาษณ์</p>
           <p class="text-xs text-gray-400 mt-1">คุณอาจยังไม่ได้รับสิทธิ์ หรือรอบยังไม่เปิด</p>
        </div>
    </div>

    <div v-if="isBookingModalOpen && selectedRound" class="fixed inset-0 z-[100] flex items-center justify-center bg-black/60 p-4 backdrop-blur-sm">
       <div class="bg-white w-full max-w-lg rounded-2xl shadow-xl overflow-hidden animate-pop-in flex flex-col max-h-[90vh]">
          <div class="p-4 border-b bg-slate-50 flex justify-between items-center">
             <span class="font-bold text-[#1e3a8a]">เลือกเวลา</span>
             <button @click="isBookingModalOpen = false" class="btn btn-xs btn-circle btn-ghost">✕</button>
          </div>
          <div class="p-4 overflow-y-auto flex-1 bg-white">
             <div class="grid grid-cols-4 gap-2">
                <button v-for="(slot, idx) in availableSlots" :key="idx"
                   @click="selectSlot(slot)"
                   :disabled="slot.isBooked"
                   class="btn btn-sm font-normal"
                   :class="slot.isSelected ? 'btn-primary bg-[#1e3a8a] text-white' : (slot.isBooked ? 'btn-disabled bg-gray-100' : 'btn-outline border-gray-200 text-gray-600')">
                   {{ slot.time }}
                </button>
             </div>
          </div>
          <div class="p-4 border-t bg-gray-50 flex justify-end gap-2">
             <button @click="isBookingModalOpen = false" class="btn btn-sm btn-ghost">ยกเลิก</button>
             <button @click="confirmBooking" :disabled="!selectedTimeSlot" class="btn btn-sm bg-[#1e3a8a] text-white">ยืนยัน</button>
          </div>
       </div>
    </div>

  </div>
</template>

<style scoped>
.animate-pop-in { animation: pop-in 0.2s cubic-bezier(0.16, 1, 0.3, 1) forwards; }
.animate-fade-in { animation: fade 0.4s ease forwards; }
@keyframes pop-in { 0% { opacity: 0; transform: scale(0.95); } 100% { opacity: 1; transform: scale(1); } }
@keyframes fade { from { opacity: 0; transform: translateY(10px); } to { opacity: 1; transform: translateY(0); } }
</style>