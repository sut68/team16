<script setup lang="ts">
import { ref, computed } from 'vue';
import { useRouter } from 'vue-router';
import SuryaGraphicImport from '@/assets/brand/Surya_graphic.png';
import SUTLogoImport from '@/assets/logo/Institute_of_Engineering_SUT_Logo.svg';
import NewsCardList from '@/components/ui/NewsCardList.vue';
import StatisticsSection from '@/components/ui/StatisticsSection.vue';
import type { NewsPost } from '@/interfaces/news_post';
import FooterBgImport from '@/assets/Suranapa-Tower01.jpg';


const SuryaGraphic = SuryaGraphicImport;
const SUTLogo = SUTLogoImport;
const FooterBg = FooterBgImport;

// Presenter image from assets (works in both dev and production)
import SutHomePageImport from '@/assets/SutHomePage.jpg';
const presenterImageUrl = SutHomePageImport;

// Router for navigation
const router = useRouter();

// Timeline data
const timelineRef = ref<HTMLElement | null>(null);

const thaiMonths = [
  'ม.ค.', 'ก.พ.', 'มี.ค.', 'เม.ย.', 'พ.ค.', 'มิ.ย.',
  'ก.ค.', 'ส.ค.', 'ก.ย.', 'ต.ค.', 'พ.ย.', 'ธ.ค.'
];

// Sample scholarship counts per month (in real app, fetch from API)
const scholarshipCounts = [2, 3, 5, 2, 1, 0, 1, 2, 3, 4, 2, 1];

const timelineMonths = computed(() => {
  const now = new Date();
  const currentMonth = now.getMonth();
  const currentYear = now.getFullYear();
  
  // Generate all 12 months starting from current month
  return Array.from({ length: 12 }, (_, i) => {
    const monthIndex = (currentMonth + i) % 12;
    const year = currentYear + Math.floor((currentMonth + i) / 12);
    const thaiYear = year + 543;
    
    return {
      name: thaiMonths[monthIndex],
      year: thaiYear.toString(),
      count: scholarshipCounts[monthIndex] ?? 0,
      isCurrentMonth: i === 0
    };
  });
});

function scrollTimeline(direction: 'left' | 'right') {
  if (timelineRef.value) {
    const scrollAmount = 150;
    timelineRef.value.scrollBy({
      left: direction === 'left' ? -scrollAmount : scrollAmount,
      behavior: 'smooth'
    });
  }
}

// Handle news card click
function handleViewNewsDetail(_news : NewsPost) {
  // Navigate to login or detail page
  // For now, redirect to login as per original code
  router.push('/login');
}

function handleViewAllNews() {
  router.push('/login');
}

</script>

<template>
  <div class="min-h-screen font-sans bg-gray-50" data-testid="homepage-container">

    <header class="bg-white shadow-sm sticky top-0 z-50" data-testid="homepage-header">
      <div class="container mx-auto px-4 h-[60px] flex items-center justify-between">
        <div class="flex items-center gap-3">
          <img 
            :src="SUTLogo" 
            alt="SUT Logo" 
            class="h-10 object-contain"
          />
        </div>

        <nav class="hidden md:flex gap-6 text-sm font-medium text-gray-600" data-testid="homepage-nav">
          <a href="#" class="hover:text-[#F26522] transition">สมัครเรียน</a>
          <a href="#" class="hover:text-[#F26522] transition">คณะและหลักสูตร</a>
          <a href="#" class="hover:text-[#F26522] transition">ชีวิตในมหาลัย</a>
          <a href="#" class="hover:text-[#F26522] transition">International</a>
        </nav>

        <router-link to="/login"
          class="bg-[#F26522] hover:bg-[#6d0016] text-white px-4 py-2 rounded-lg font-bold text-sm transition shadow-lg flex items-center gap-1.5"
          data-testid="homepage-login-button">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="m10 17 5-5-5-5"/>
            <path d="M15 12H3"/>
            <path d="M15 3h4a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-4"/>
          </svg>
          ล็อคอิน
        </router-link>
      </div>
    </header>

    <!-- Hero Banner -->
    <section class="relative bg-gradient-to-r from-[#8B001D] to-[#f97316] text-white py-16 md:py-24 overflow-hidden" data-testid="homepage-hero-section">
      <div class="absolute inset-0 z-0 overflow-hidden">
        <!-- Container for centering -->
        <div 
          class="absolute bottom-[-25vw] left-1/2 -translate-x-1/2" 
          style="width: 70vw; min-width: 600px; max-width: 1200px;"
        >
          <img 
            :src="SuryaGraphic" 
            alt="" 
            class="surya-rotate opacity-20 w-full h-full object-contain"
          />
        </div>
      </div>
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
            <!-- <span class="text-white/40 font-bold text-2xl">Presenter Image</span> -->
            <img :src="presenterImageUrl" alt="Presenter" class="w-full h-full object-cover rounded-2xl" />

            <div
              class="absolute -top-4 -right-4 bg-[#253C90] text-white px-4 py-3 rounded-lg shadow-lg font-bold animate-bounce border border-white/30">
              ทุน 100%
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Scholarship Timeline Section (Compact) -->
    <section class="container mx-auto px-4 -mt-8 relative z-20" data-testid="homepage-timeline-section">
      <div class="bg-white rounded-xl shadow-xl p-2 md:p-4 max-w-4xl mx-auto border-b-4 border-[#F26522]">
        <div class="flex items-center gap-3">
          <!-- Header -->
          <div class="flex items-center gap-2 flex-shrink-0">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-[#F26522]" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <rect width="18" height="18" x="3" y="4" rx="2" ry="2"/>
              <line x1="16" x2="16" y1="2" y2="6"/>
              <line x1="8" x2="8" y1="2" y2="6"/>
              <line x1="3" x2="21" y1="10" y2="10"/>
            </svg>
            <span class="text-sm font-bold text-gray-700 hidden sm:inline">ปฏิทินทุน</span>
          </div>
          
          <!-- Scroll Left Button -->
          <button 
            class="hidden md:flex w-6 h-6 bg-gray-100 hover:bg-[#F26522] hover:text-white rounded-full items-center justify-center text-gray-400 transition-all flex-shrink-0"
            @click="scrollTimeline('left')"
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path d="m15 18-6-6 6-6"/>
            </svg>
          </button>
          
          <!-- Timeline Months -->
          <div 
            ref="timelineRef"
            class="flex gap-2 overflow-x-auto scroll-smooth flex-1 py-2 px-2"
            style="scrollbar-width: none; -ms-overflow-style: none;"
          >
            <div 
              v-for="(month, index) in timelineMonths" 
              :key="index"
              @click="router.push('/login')"
              class="flex-shrink-0 cursor-pointer"
            >
              <div 
                class="relative px-4 py-2.5 rounded-lg border transition-all duration-200 hover:shadow-md hover:scale-105 flex items-center gap-2"
                :class="[
                  month.isCurrentMonth 
                    ? 'bg-gradient-to-r from-[#F26522] to-[#ea580c] border-[#F26522] text-white' 
                    : 'bg-gray-50 border-gray-200 hover:border-[#F26522] text-gray-700'
                ]"
              >
                <!-- Current month indicator -->
                <div 
                  v-if="month.isCurrentMonth"
                  class="w-2 h-2 bg-green-400 rounded-full animate-pulse"
                ></div>
                
                <!-- Month name -->
                <span class="text-xs font-bold">{{ month.name }}</span>
                
                <!-- Count badge -->
                <span 
                  v-if="month.count > 0"
                  class="text-[10px] font-bold px-1.5 py-0.5 rounded-full"
                  :class="month.isCurrentMonth ? 'bg-white/30 text-white' : 'bg-[#F26522]/10 text-[#F26522]'"
                >{{ month.count }}</span>
              </div>
            </div>
          </div>
          
          <!-- Scroll Right Button -->
          <button 
            class="hidden md:flex w-6 h-6 bg-gray-100 hover:bg-[#F26522] hover:text-white rounded-full items-center justify-center text-gray-400 transition-all flex-shrink-0"
            @click="scrollTimeline('right')"
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path d="m9 18 6-6-6-6"/>
            </svg>
          </button>
        </div>
      </div>
    </section>

    <!-- Statistics Section -->
    <StatisticsSection />

    <!-- News Cards Section -->
    <section class="container mx-auto px-4 py-12" data-testid="homepage-news-section">
      <NewsCardList
        :limit="6"
        :only-public="true"
        title="ข่าวประชาสัมพันธ์"
        :show-view-all="true"
        @view-detail="handleViewNewsDetail"
        @view-all="handleViewAllNews"
      />
    </section>

    <div class="fixed right-0 top-1/2 transform -translate-y-1/2 z-50 flex flex-col gap-1 hidden md:flex" data-testid="homepage-social-links">
      <a href="#"
        class="w-10 h-10 bg-[#1877F2] text-white flex items-center justify-center hover:w-12 transition-all shadow-md">F</a>
      <a href="#"
        class="w-10 h-10 bg-[#FF0000] text-white flex items-center justify-center hover:w-12 transition-all shadow-md">Y</a>
      <a href="#"
        class="w-10 h-10 bg-[#06C755] text-white flex items-center justify-center hover:w-12 transition-all shadow-md">L</a>
    </div>

    <!-- Footer -->
    <footer class="relative mt-16">
      <img 
        :src="FooterBg" 
        alt="" 
        class="w-full object-cover"
      />
      
      <!-- gradient: 100% gray-50 (top) to 0% transparent (bottom) -->
      <div class="absolute inset-0 bg-gradient-to-b from-gray-50 via-gray-50/30 via-40% to-transparent"></div>
      
      <!-- Footer content -->
      <div class="absolute inset-0 z-10 flex items-end pb-8">
        <div class="container mx-auto px-4">
          <div class="grid grid-cols-1 md:grid-cols-4 gap-8 text-white mb-8">
            <div class="md:col-span-1">
              <div class="flex items-center gap-3 mb-4">
                <div class="w-12 h-12 bg-[#F26522] rounded-lg flex items-center justify-center font-bold text-white">
                  SUT
                </div>
                <div>
                  <h3 class="font-bold text-lg">SURANAREE</h3>
                  <p class="text-xs text-white/70">University of Technology</p>
                </div>
              </div>
              <p class="text-sm text-white/80 leading-relaxed">
                ระบบทุนการศึกษา มหาวิทยาลัยเทคโนโลยีสุรนารี
              </p>
            </div>

            <!-- Quick Links เดี๋ยวค่อยมาปรับอีกที-->
            <div>
              <h4 class="font-bold mb-4 text-[#F26522]">เมนูลัด</h4>
              <ul class="space-y-2 text-sm text-white/80">
                <li><a href="#" class="hover:text-[#F26522] transition">หน้าแรก</a></li>
                <li><a href="#" class="hover:text-[#F26522] transition">ทุนการศึกษา</a></li>
                <li><a href="#" class="hover:text-[#F26522] transition">สมัครทุน</a></li>
                <li><a href="#" class="hover:text-[#F26522] transition">ติดต่อเรา</a></li>
              </ul>
            </div>

            <!-- Contact -->
            <div>
              <h4 class="font-bold mb-4 text-[#F26522]">ติดต่อเรา</h4>
              <ul class="space-y-2 text-sm text-white/80">
                <li class="flex items-center gap-2">
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z"/>
                  </svg>
                  044-224-xxxx
                </li>
                <li class="flex items-center gap-2">
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"/>
                  </svg>
                  scholarship@sut.ac.th
                </li>
              </ul>
            </div>

            <!-- Social -->
            <div>
              <h4 class="font-bold mb-4 text-[#F26522]">ติดตามเรา</h4>
              <div class="flex gap-3">
                <a href="#" class="w-10 h-10 bg-white/10 hover:bg-[#F26522] rounded-lg flex items-center justify-center transition">
                  <span class="font-bold">F</span>
                </a>
                <a href="#" class="w-10 h-10 bg-white/10 hover:bg-[#F26522] rounded-lg flex items-center justify-center transition">
                  <span class="font-bold">Y</span>
                </a>
                <a href="#" class="w-10 h-10 bg-white/10 hover:bg-[#F26522] rounded-lg flex items-center justify-center transition">
                  <span class="font-bold">L</span>
                </a>
              </div>
            </div>
          </div>

          <div class="border-t border-white/20 pt-6 text-center text-sm text-white/60">
            <p>พัฒนาโดย Team 16 | SE Project</p>
          </div>
        </div>
      </div>
    </footer>

  </div>
</template>

<style scoped>
@keyframes slow-rotate {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

.surya-rotate {
  /* ปรับเวลาหมุน */
  animation: slow-rotate 60s linear infinite;
  transform-origin: center center;
}
</style>