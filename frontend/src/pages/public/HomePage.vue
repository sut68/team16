<script setup lang="ts">
import { ref, computed } from 'vue';
import SuryaGraphicImport from '@/assets/brand/Surya_graphic.png';

// 2. ประกาศตัวแปรเพื่อรับค่า (แก้ปัญหา Type Error)
const SuryaGraphic = SuryaGraphicImport;

interface TabItem {
  id: string;
  label: string;
}

interface Scholarship {
  id: number;
  title: string;
  description: string;
  amount: string;
  category: 'bachelor' | 'master' | 'doctoral' | 'inter';
}

const activeTab = ref<string>('bachelor');

const tabs: TabItem[] = [
  { id: 'bachelor', label: 'ปริญญาตรี' },
  { id: 'master', label: 'ปริญญาโท' },
  { id: 'doctoral', label: 'ปริญญาเอก' },
  { id: 'inter', label: 'International' },
];

const scholarships: Scholarship[] = [
  {
    id: 1,
    title: 'ทุนกล้าจน (SUT Low Income)',
    description: 'ทุนสนับสนุนค่าหน่วยกิต 100% สำหรับนักศึกษาที่มีความประพฤติดีแต่ขาดแคลนทุนทรัพย์',
    amount: 'ทุน 100%',
    category: 'bachelor',
  },
  {
    id: 2,
    title: 'ทุนศักยภาพทางวิชาการ',
    description: 'สำหรับผู้มีผลการเรียนดีเยี่ยม (GPAX 3.50+) รับทุนตลอดหลักสูตร',
    amount: 'ทุนเต็มจำนวน',
    category: 'bachelor',
  },
  {
    id: 3,
    title: 'ทุนวิจัยระดับบัณฑิตศึกษา',
    description: 'สนับสนุนงานวิจัยและค่าเล่าเรียนสำหรับปริญญาโท',
    amount: 'ค่าเล่าเรียน + เบี้ยเลี้ยง',
    category: 'master',
  },
  {
    id: 4,
    title: 'Royal Golden Jubilee (RGJ)',
    description: 'โครงการปริญญาเอกกาญจนาภิเษก เพื่อสร้างนักวิจัยระดับสูง',
    amount: 'ทุนวิจัยเต็มจำนวน',
    category: 'doctoral',
  },
  {
    id: 5,
    title: 'ASEAN Scholarship',
    description: 'Full scholarship for students from ASEAN member countries.',
    amount: 'Full Tuition',
    category: 'inter',
  },
];

const filteredScholarships = computed(() => {
  return scholarships.filter((s) => s.category === activeTab.value);
});

const setActiveTab = (id: string) => {
  activeTab.value = id;
};
</script>

<template>
  <div class="min-h-screen bg-[#EEEEEE] font-sans" data-testid="homepage-container">

    <header class="bg-white shadow-sm sticky top-0 z-50" data-testid="homepage-header">
      <div class="container mx-auto px-4 h-16 flex items-center justify-between">
        <div class="flex items-center gap-2">
          <div
            class="w-10 h-10 bg-[#F26522] text-white flex items-center justify-center font-bold text-xl rounded shadow-sm border border-[#8B001D]">
            SUT</div>
          <div class="flex flex-col">
            <span class="font-bold text-lg text-[#253C90] leading-none">SURANAREE</span>
            <span class="text-[10px] text-gray-500 tracking-wider">UNIVERSITY OF TECHNOLOGY</span>
          </div>
        </div>

        <nav class="hidden md:flex gap-6 text-sm font-medium text-gray-600" data-testid="homepage-nav">
          <a href="#" class="hover:text-[#F26522] transition">สมัครเรียน</a>
          <a href="#" class="hover:text-[#F26522] transition">คณะและหลักสูตร</a>
          <a href="#" class="hover:text-[#F26522] transition">ชีวิตในมหาลัย</a>
          <a href="#" class="hover:text-[#F26522] transition">International</a>
        </nav>

        <router-link to="/login"
          class="bg-[#F26522] hover:bg-[#6d0016] text-white px-5 py-2 rounded-md font-bold text-sm transition shadow-lg flex items-center gap-2"
          data-testid="homepage-login-button">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
          </svg>
          ล็อคอิน
        </router-link>
      </div>
    </header>

    <section class="relative bg-[#8B001D] text-white py-16 md:py-24 overflow-hidden" data-testid="homepage-hero-section">
      <div class="absolute inset-0 z-0 bg-cover bg-left-top bg-no-repeat opacity-20"
        :style="{ backgroundImage: `url(${SuryaGraphic})`, backgroundPosition: '-50px -50px' }"></div>
      <div
        class="absolute top-0 right-0 w-80 h-80 bg-[#F26522] opacity-10 rounded-full blur-3xl transform translate-x-1/2 -translate-y-1/2">
      </div>
      <div
        class="absolute bottom-0 left-0 w-64 h-64 bg-[#8B001D] opacity-20 rounded-full blur-2xl transform -translate-x-1/2 translate-y-1/2">
      </div>

      <div class="container mx-auto px-4 relative z-10 flex flex-col md:flex-row items-center">
        <div class="md:w-1/2 text-center md:text-left mb-8 md:mb-0 pl-0 md:pl-8">
          <span class="inline-block bg-[#F26522] text-white text-xm font-bold px-3 py-1 rounded-full mb-4 shadow-md">
            เปิดรับสมัครแล้ว!!
          </span>
          <h1 class="text-4xl md:text-7xl font-extrabold mb-4 leading-tight" data-testid="homepage-hero-title">
            #TCAS69 <br />
            <span
              class="text-[#F26522] bg-white px-4 rounded-lg transform -skew-x-6 inline-block mt-2 font-black text-gray-500 italic">SUT
              TEAM</span>
          </h1>
          <p class="text-xl md:text-2xl font-light mb-8">
            ทุนเรียนดี SUT ตัวจริง รับทุนสูงสุด <span
              class="font-bold text-white border-b-4 border-[#F26522]">100%*</span>
          </p>
          <button
            class="bg-white text-[#253C90] px-8 py-3 rounded-full font-bold text-lg hover:shadow-xl hover:scale-105 transition transform border-2 border-transparent hover:border-[#F26522]"
            data-testid="homepage-apply-button">
            ยื่นสมัครเลย
          </button>
          <p class="mt-4 text-xm opacity-70">*เงื่อนไขเป็นไปตามที่มหาวิทยาลัยกำหนด</p>
        </div>

        <div class="md:w-1/2 flex justify-center relative">
          <div
            class="w-72 h-80 md:w-96 md:h-96 bg-gradient-to-tr from-[#EEEEEE] to-white/10 backdrop-blur-sm rounded-2xl border-2 border-white/20 flex items-center justify-center relative shadow-2xl">
            <span class="text-white/40 font-bold text-2xl">Presenter Image</span>

            <div
              class="absolute -top-4 -right-4 bg-[#253C90] text-white px-4 py-3 rounded-lg shadow-lg font-bold animate-bounce border border-white/30">
              ทุน 100%
            </div>
          </div>
        </div>
      </div>
    </section>

    <section class="container mx-auto px-4 -mt-8 relative z-20" data-testid="homepage-tabs-section">
      <div class="bg-white rounded-xl shadow-xl p-2 md:p-4 max-w-4xl mx-auto border-b-4 border-[#F26522]">
        <div class="flex flex-wrap justify-center gap-2 md:gap-4">
          <button v-for="tab in tabs" :key="tab.id" @click="setActiveTab(tab.id)"
            class="px-6 py-3 rounded-md font-bold text-sm md:text-base transition-all duration-300 w-full md:w-auto"
            :data-testid="`homepage-tab-${tab.id}`"
            :class="[
              activeTab === tab.id
                ? 'bg-[#F26522] text-white shadow-md scale-105'
                : 'bg-[#EEEEEE] text-gray-500 hover:bg-gray-200'
            ]">
            {{ tab.label }}
          </button>
        </div>
      </div>
    </section>

    <section class="container mx-auto px-4 py-16 max-w-5xl" data-testid="homepage-scholarships-section">
      <div class="mb-8 border-b pb-4 border-gray-300">
        <h2 class="text-3xl font-bold text-[#253C90]">
          ทุนการศึกษาที่ได้รับ:
          <span class="text-[#F26522] text-xl font-normal ml-2">({{ tabs.find((t: { id: any; }) => t.id ===
            activeTab)?.label }})</span>
        </h2>
      </div>

      <div v-if="filteredScholarships.length > 0" class="grid grid-cols-1 md:grid-cols-2 gap-6" data-testid="homepage-scholarships-grid">
        <div v-for="scholarship in filteredScholarships" :key="scholarship.id"
          class="bg-white border border-gray-200 rounded-lg p-6 hover:shadow-xl transition group cursor-pointer relative overflow-hidden"
          :data-testid="`homepage-scholarship-card-${scholarship.id}`">
          <div class="absolute top-0 left-0 w-2 h-full bg-[#F26522] group-hover:w-3 transition-all"></div>

          <h3 class="text-xl font-bold text-[#253C90] mb-2 group-hover:text-[#F26522] transition">{{ scholarship.title
            }}</h3>
          <p class="text-gray-600 mb-4 font-light">{{ scholarship.description }}</p>

          <div class="flex items-center justify-between mt-4">
            <span class="text-[#8B001D] font-bold bg-[#F26522]/10 px-3 py-1 rounded-full text-sm">
              {{ scholarship.amount }}
            </span>
            <button class="text-[#253C90] font-semibold text-sm hover:text-[#F26522] hover:underline">
              ดูรายละเอียด >
            </button>
          </div>
        </div>
      </div>

      <div v-else class="text-center py-12 text-gray-400" data-testid="homepage-scholarships-empty">
        <p>ไม่มีรายการทุนในหมวดหมู่นี้ขณะนี้</p>
      </div>
    </section>

    <div class="fixed right-0 top-1/2 transform -translate-y-1/2 z-50 flex flex-col gap-1 hidden md:flex" data-testid="homepage-social-links">
      <a href="#"
        class="w-10 h-10 bg-[#1877F2] text-white flex items-center justify-center hover:w-12 transition-all shadow-md">F</a>
      <a href="#"
        class="w-10 h-10 bg-[#FF0000] text-white flex items-center justify-center hover:w-12 transition-all shadow-md">Y</a>
      <a href="#"
        class="w-10 h-10 bg-[#06C755] text-white flex items-center justify-center hover:w-12 transition-all shadow-md">L</a>
    </div>

  </div>
</template>

<style scoped></style>