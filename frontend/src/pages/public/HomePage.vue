<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import SuryaGraphicImport from '@/assets/brand/Surya_graphic.png';
import SUTLogoImport from '@/assets/logo/Institute_of_Engineering_SUT_Logo.svg';
import { Get } from '@/services/api/https';
import ScholarshipCardList from '@/components/ui/ScholarshipCardList.vue';
import StatisticsSection from '@/components/ui/StatisticsSection.vue';
import type { ScholarshipResponse } from '@/interfaces/scholarship';
import FooterBgImport from '@/assets/Suranapa-Tower01.jpg';

// 2. ประกาศตัวแปรเพื่อรับค่า (แก้ปัญหา Type Error)
const SuryaGraphic = SuryaGraphicImport;
const SUTLogo = SUTLogoImport;
const FooterBg = FooterBgImport;

interface TabItem {
  id: string;
  label: string;
}

interface NewsPost {
  ID: number;
  title: string;
  post_detail: string;
  scholarship: {
    sponsor: {
      company_name: string; // Corrected from sponsor_name
    };
    typescholarship: {
      type_name: string;
    };
  };
}

const activeTab = ref<string>('all');
const newsPosts = ref<NewsPost[]>([]);

const tabs: TabItem[] = [
  { id: 'all', label: 'ทั้งหมด' },
  { id: 'bachelor', label: 'ปริญญาตรี' },
  { id: 'master', label: 'ปริญญาโท' },
  { id: 'doctoral', label: 'ปริญญาเอก' },
  { id: 'inter', label: 'International' },
];

const fetchNews = async () => {
  try {
    // Use Get with requireAuth = false to avoid token warning and headers
    const response = await Get('/newsposts', false);
    if (response) {
      newsPosts.value = response;
    }
  } catch (error) {
    console.error('Error fetching news:', error);
  }
};

onMounted(() => {
  fetchNews();
});

// const filteredScholarships = computed(() => {
//   // Debug: Log data to see what we got
//   console.log('News Data:', newsPosts.value);
  
//   // Return all news for now because backend types (Full/Partial) don't match tabs (Bachelor/Master)
//   return newsPosts.value; 

//   /* 
//   const currentTabLabel = tabs.find(t => t.id === activeTab.value)?.label;
//   if (!currentTabLabel) return [];
  
//   return newsPosts.value.filter((n) => {
//     // Optional chaining in case of missing data
//     return n.scholarship?.typescholarship?.type_name === currentTabLabel;
//   });
//   */
// });

const setActiveTab = (id: string) => {
  activeTab.value = id;
};

// Router for navigation
const router = useRouter();

// Handle scholarship card click
function handleViewScholarshipDetail(_scholarship: ScholarshipResponse) {
  // Navigate to login or detail page
  router.push('/login');
}

function handleViewAllScholarships() {
  router.push('/login');
}
</script>

<template>
  <div class="min-h-screen font-sans bg-gray-50" data-testid="homepage-container">

    <header class="bg-white shadow-sm sticky top-0 z-50" data-testid="homepage-header">
      <div class="container mx-auto px-4 h-20 flex items-center justify-between">
        <div class="flex items-center gap-3">
          <img 
            :src="SUTLogo" 
            alt="SUT Logo" 
            class="h-14 object-contain"
          />
        </div>

        <nav class="hidden md:flex gap-8 text-base font-medium text-gray-600" data-testid="homepage-nav">
          <a href="#" class="hover:text-[#F26522] transition">สมัครเรียน</a>
          <a href="#" class="hover:text-[#F26522] transition">คณะและหลักสูตร</a>
          <a href="#" class="hover:text-[#F26522] transition">ชีวิตในมหาลัย</a>
          <a href="#" class="hover:text-[#F26522] transition">International</a>
        </nav>

        <router-link to="/login"
          class="bg-[#F26522] hover:bg-[#6d0016] text-white px-6 py-2.5 rounded-lg font-bold text-base transition shadow-lg flex items-center gap-2"
          data-testid="homepage-login-button">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
          </svg>
          ล็อคอิน
        </router-link>
      </div>
    </header>

    <!-- Hero Banner -->
    <section class="relative bg-[#8B001D] text-white py-16 md:py-24 overflow-hidden" data-testid="homepage-hero-section">
      <div class="absolute inset-0 z-0 overflow-hidden opacity-20">
        <img 
          :src="SuryaGraphic" 
          alt="" 
          class="surya-rotate"
          style="position: absolute; top: -50px; left: -50px;"
        />
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
            <img src="http://localhost:8080/uploads/news/Bodyslam1859.jpg" alt="Presenter" class="w-full h-full object-cover rounded-2xl" />

            <div
              class="absolute -top-4 -right-4 bg-[#253C90] text-white px-4 py-3 rounded-lg shadow-lg font-bold animate-bounce border border-white/30">
              ทุน 100%
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Tabs -->
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

    <!-- <section class="container mx-auto px-4 py-16 max-w-5xl" data-testid="homepage-scholarships-section">
      <div class="mb-8 border-b pb-4 border-gray-300">
        <h2 class="text-3xl font-bold text-[#253C90]">
          ทุนการศึกษาที่ได้รับ:
          <span class="text-[#F26522] text-xl font-normal ml-2">({{ tabs.find((t: { id: any; }) => t.id ===
            activeTab)?.label }})</span>
        </h2>
      </div>

      <div v-if="filteredScholarships.length > 0" class="grid grid-cols-1 md:grid-cols-2 gap-6" data-testid="homepage-scholarships-grid">
        <div v-for="news in filteredScholarships" :key="news.ID"
          class="bg-white border border-gray-200 rounded-lg p-6 hover:shadow-xl transition group cursor-pointer relative overflow-hidden"
          :data-testid="`homepage-scholarship-card-${news.ID}`">
          <div class="absolute top-0 left-0 w-2 h-full bg-[#F26522] group-hover:w-3 transition-all"></div>

          <h3 class="text-xl font-bold text-[#253C90] mb-2 group-hover:text-[#F26522] transition">{{ news.title
            }}</h3>
          <p class="text-gray-600 mb-4 font-light line-clamp-2">{{ news.post_detail }}</p>

          <div class="flex items-center justify-between mt-4">
            <span class="text-[#8B001D] font-bold bg-[#F26522]/10 px-3 py-1 rounded-full text-sm">
              {{ news.scholarship?.sponsor?.company_name || 'ทุนการศึกษา' }}
            </span>
            <router-link to="/login" class="text-[#253C90] font-semibold text-sm hover:text-[#F26522] hover:underline">
              ดูรายละเอียด >
            </router-link>
          </div>
        </div>
      </div>

      <div v-else class="text-center py-12 text-gray-400" data-testid="homepage-scholarships-empty">
        <p>ไม่มีรายการทุนในหมวดหมู่นี้ขณะนี้</p>
      </div>
    </section> -->

    <!-- Statistics Section สถิติเดี๋ยวไปดึงข้อมูลมาอีกที -->
    <StatisticsSection />

    <!-- Scholarship Cards Section ส่วนแสดงทุนการศึกษาที่เปิดรับสมัคร -->
    <section class="container mx-auto px-4 py-12" data-testid="homepage-scholarship-cards">
      <ScholarshipCardList
        :limit="6"
        :only-open="true"
        title="ทุนการศึกษาที่เปิดรับสมัคร"
        :show-view-all="true"
        @view-detail="handleViewScholarshipDetail"
        @view-all="handleViewAllScholarships"
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