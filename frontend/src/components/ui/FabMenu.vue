<script setup lang="ts">
import { ref } from 'vue';
import { 
  MessageCircle, 
  Facebook, 
  Phone,
  AlertCircle
} from 'lucide-vue-next';
import EngiGearLogoImport from '@/assets/logo/ENGi-Gear-Logo-Official.png';

const EngiGearLogo = EngiGearLogoImport;
const showFacebookModal = ref(false);
const showPhoneModal = ref(false);

const handleFacebookClick = () => {
  showFacebookModal.value = true;
};

const handlePhoneClick = () => {
  showPhoneModal.value = true;
};

const closeFacebookModal = () => {
  showFacebookModal.value = false;
};

const closePhoneModal = () => {
  showPhoneModal.value = false;
};
</script>

<template>
  <!-- Pure CSS Material FAB Menu (Radial Layout) -->
  <div class="fixed bottom-8 right-8 z-[60] flex items-center justify-center w-14 h-14" data-testid="fab-menu">
    <!-- Checkbox Hack Trigger -->
    <input type="checkbox" id="fab-toggle" class="peer hidden" />
    
    <!-- Backdrop Overlay (Click outside to close) -->
    <label 
      for="fab-toggle" 
      class="fixed inset-0 z-[-1] hidden peer-checked:block cursor-default"
    ></label>
    
    <!-- Action 1: Contact/Assistance (Pops UP) -->
    <router-link 
      to="/dashboard/assistance" 
      class="group absolute flex items-center justify-center transition-all duration-500 opacity-0 scale-50 z-0
             peer-checked:opacity-100 peer-checked:scale-100 peer-checked:-translate-y-24 peer-checked:z-50"
    >
      <span class="absolute -top-10 opacity-0 group-hover:opacity-100 transition-all duration-300 transform translate-y-4 group-hover:translate-y-0 bg-white text-gray-800 text-xs font-bold px-3 py-1.5 rounded-lg shadow-lg border border-gray-100 whitespace-nowrap">
        ติดต่อช่วยเหลือ
      </span>
      <div class="w-12 h-12 bg-white text-[#06C755] rounded-full shadow-lg flex items-center justify-center border border-gray-100 transition-all duration-300 group-hover:shadow-[0_0_15px_rgba(6,199,85,0.4)] group-hover:scale-110">
        <MessageCircle class="w-6 h-6" />
      </div>
    </router-link>

    <!-- Action 2: Facebook (Pops DIAGONAL) -->
    <button 
      @click="handleFacebookClick"
      class="group absolute flex items-center justify-center transition-all duration-500 opacity-0 scale-50 z-0 delay-75
             peer-checked:opacity-100 peer-checked:scale-100 peer-checked:-translate-y-16 peer-checked:-translate-x-16 peer-checked:z-50"
    >
      <span class="absolute right-14 opacity-0 group-hover:opacity-100 transition-all duration-300 transform translate-x-4 group-hover:translate-x-0 bg-white text-gray-800 text-xs font-bold px-3 py-1.5 rounded-lg shadow-lg border border-gray-100 whitespace-nowrap">
        Facebook
      </span>
      <div class="w-12 h-12 bg-white text-[#1877F2] rounded-full shadow-lg flex items-center justify-center border border-gray-100 transition-all duration-300 group-hover:shadow-[0_0_15px_rgba(24,119,242,0.4)] group-hover:scale-110">
        <Facebook class="w-6 h-6" />
      </div>
    </button>

    <!-- Action 3: Phone (Pops LEFT) -->
    <button 
      @click="handlePhoneClick"
      class="group absolute flex items-center justify-center transition-all duration-500 opacity-0 scale-50 z-0 delay-150
             peer-checked:opacity-100 peer-checked:scale-100 peer-checked:-translate-x-24 peer-checked:z-50"
    >
      <span class="absolute right-14 opacity-0 group-hover:opacity-100 transition-all duration-300 transform translate-x-4 group-hover:translate-x-0 bg-white text-gray-800 text-xs font-bold px-3 py-1.5 rounded-lg shadow-lg border border-gray-100 whitespace-nowrap">
        โทรสอบถาม
      </span>
      <div class="w-12 h-12 bg-white text-red-500 rounded-full shadow-lg flex items-center justify-center border border-gray-100 transition-all duration-300 group-hover:shadow-[0_0_15px_rgba(239,68,68,0.4)] group-hover:scale-110">
        <Phone class="w-6 h-6" />
      </div>
    </button>

    <!-- Main Trigger Button -->
    <label 
      for="fab-toggle" 
      class="relative z-10 w-14 h-14 bg-white hover:bg-gray-50 text-white rounded-full shadow-[0_8px_25px_rgba(242,101,34,0.4)] flex items-center justify-center cursor-pointer transition-all duration-700 peer-checked:rotate-[360deg] active:scale-95 shadow-orange-500/20"
    >
      <img 
        :src="EngiGearLogo" 
        alt="Engineering Gear" 
        class="w-9 h-auto object-contain peer-checked:animate-none animate-[spin_10s_linear_infinite]" 
      />
    </label>
  </div>

  <!-- DaisyUI Modal for Facebook -->
  <div v-if="showFacebookModal" class="modal modal-open" @click.self="closeFacebookModal">
    <div class="modal-box bg-white shadow-2xl border border-gray-100">
      <div class="flex items-center gap-3 mb-4">
        <div class="w-12 h-12 bg-blue-100 rounded-full flex items-center justify-center">
          <AlertCircle class="w-6 h-6 text-[#1877F2]" />
        </div>
        <h3 class="font-bold text-lg text-gray-800">แจ้งเตือน</h3>
      </div>
      <p class="text-gray-600 text-center">ขออภัย! ขณะนี้ยังไม่มี Facebook Page</p>
      <p class="text-gray-500 text-sm mt-2 text-center">กรุณาติดต่อผ่านช่องทางอื่น เช่น ระบบช่วยเหลือ หรือโทรศัพท์</p>
      <div class="modal-action">
        <button @click="closeFacebookModal" class="btn btn-primary bg-[#F26522] border-none hover:bg-[#d95a1e] text-white">
          ปิด
        </button>
      </div>
    </div>
  </div>

  <!-- DaisyUI Modal for Phone -->
  <div v-if="showPhoneModal" class="modal modal-open" @click.self="closePhoneModal">
    <div class="modal-box bg-white shadow-2xl border border-gray-100">
      <div class="flex items-center gap-3 mb-4">
        <div class="w-12 h-12 bg-green-100 rounded-full flex items-center justify-center">
          <Phone class="w-6 h-6 text-green-600" />
        </div>
        <h3 class="font-bold text-lg text-gray-800">ติดต่อเรา</h3>
      </div>
      <div class="text-center py-4">
        <p class="text-gray-500 text-sm mb-2">สำนักงานทุนการศึกษา มทส.</p>
        <p class="text-2xl font-bold text-gray-800 mb-1">044-224-XXX</p>
        <p class="text-gray-400 text-xs">วันจันทร์ - ศุกร์ เวลา 08:30 - 16:30 น.</p>
      </div>
      <div class="modal-action">
        <button @click="closePhoneModal" class="btn btn-primary bg-[#F26522] border-none hover:bg-[#d95a1e] text-white">
          ปิด
        </button>
      </div>
    </div>
  </div>
</template>
