<script setup lang="ts">
import { ref, watch, computed } from 'vue';
import { 
  X, ImageIcon, Calendar, Clock, GraduationCap, 
  CheckCircle2, AlertCircle 
} from 'lucide-vue-next';
import { getNewsPostById } from '@/services/api/news_post';
import type { NewsPostDetailResponse } from '@/interfaces/news_post';

const props = defineProps<{
  isOpen: boolean;
  id: number;
}>();

const emit = defineEmits(['close']);
const API_BASE_URL = 'http://localhost:8080';

// --- State ---
const loading = ref(false);
const error = ref('');
const news = ref<any>(null);
const features = ref<any[]>([]);

// --- Helper: Date Formatter ---
const formatDate = (dateString: string | undefined) => {
    if (!dateString) return 'ไม่ระบุ';
    const date = new Date(dateString);
    if (isNaN(date.getTime()) || date.getFullYear() < 2000) return 'ไม่ระบุ';
    return date.toLocaleDateString('th-TH', { year: 'numeric', month: 'short', day: 'numeric' });
};

// --- Helper: Timeline Logic (Robust Version) ---
const getProgressPercent = (startStr: string, endStr: string) => {
    if (!startStr || !endStr) return 0;
    const now = new Date().getTime();
    const start = new Date(startStr).getTime();
    const end = new Date(endStr).getTime();
    
    if (now < start) return 0;
    if (now > end) return 100;

    const totalDuration = end - start;
    if (totalDuration <= 0) return 100; 

    const elapsed = now - start;
    return Math.min(100, Math.max(0, (elapsed / totalDuration) * 100));
};

const getDaysRemaining = (endStr: string) => {
    if (!endStr) return 'ไม่ระบุ';
    const now = new Date().getTime();
    const end = new Date(endStr).getTime();
    const diff = end - now;
    
    if (diff < 0) return 'หมดเขตแล้ว';
    const diffDays = Math.ceil(diff / (1000 * 60 * 60 * 24));
    
    if (diffDays === 0) return 'วันนี้วันสุดท้าย!'; 
    return `เหลืออีก ${diffDays} วัน`;
};

const getScholarshipStatus = (openDateStr: string | undefined, closeDateStr: string | undefined) => {
    if (!openDateStr || !closeDateStr) return { label: 'ไม่ระบุ', color: 'gray', icon: AlertCircle };
    const now = new Date();
    const openDate = new Date(openDateStr);
    const closeDate = new Date(closeDateStr);
    
    // Reset time for status check
    now.setHours(0,0,0,0); openDate.setHours(0,0,0,0); closeDate.setHours(0,0,0,0);

    if (now < openDate) return { label: 'เร็วๆ นี้', color: 'amber', icon: Clock };
    else if (now > closeDate) return { label: 'ปิดรับสมัคร', color: 'red', icon: CheckCircle2 };
    else return { label: 'กำลังเปิดรับ', color: 'green', icon: CheckCircle2 };
};

// --- Helper: Image & Features ---
const fileUrl = computed(() => {
  if (news.value?.file_path) {
    const cleanPath = news.value.file_path.startsWith('/') ? news.value.file_path.substring(1) : news.value.file_path;
    return `${API_BASE_URL}/${cleanPath}`;
  }
  return null;
});

const handleImageError = (event: Event) => { (event.target as HTMLImageElement).style.display = 'none'; };

const formatCondition = (feature: any) => {
    if (!feature) return "-";
    const op = feature.operator; const val = feature.value;
    const typeName = feature.Typefeature?.type_feature_name || feature.type_feature?.type_feature_name || ''; 
    const nameVal = parseFloat(val);
    const displayVal = !isNaN(nameVal) && nameVal > 1000 ? nameVal.toLocaleString() : val;
    
    let text = "";
    switch (op) {
        case '>=': text = `ไม่น้อยกว่า ${displayVal}`; break;
        case '<=': text = `ไม่เกิน ${displayVal}`; break;
        case '>': text = `มากกว่า ${displayVal}`; break;
        case '<': text = `น้อยกว่า ${displayVal}`; break;
        case '==': text = `เท่ากับ ${displayVal}`; break;
        case '!=': text = `ไม่เท่ากับ ${displayVal}`; break;
        default: text = `${op || ''} ${displayVal}`;
    }
    
    let unit = " ";
    if (typeName.includes("รายได้") || typeName.includes("เงินเดือน")) unit = " บาท";
    if (typeName) return `${typeName} ${text}${unit}`;
    return `${text}${unit}`;
};

// --- Fetch Data ---
const fetchData = async () => {
  if (!props.id) return;
  loading.value = true;
  error.value = '';
  try {
    const response = await getNewsPostById(props.id) as NewsPostDetailResponse; 
    // Handle both structure types (just in case)
    const data = { ...response, features_scholarship: response.features_scholarship || response.features || [] };
    news.value = data.news_post;
    features.value = data.features_scholarship;
  } catch (err) { 
      console.error(err); 
      error.value = "เกิดข้อผิดพลาดในการโหลดข้อมูล"; 
  } finally { 
      loading.value = false; 
  }
};

// Watch for ID changes or Open state
watch(() => [props.id, props.isOpen], ([newId, newIsOpen]) => {
    if (newIsOpen && newId) fetchData();
}, { immediate: true });

</script>

<template>
  <div 
    v-if="isOpen" 
    class="fixed inset-0 z-[1000] flex items-center justify-center bg-black/60 backdrop-blur-sm p-4 transition-opacity animate-fade-in font-prompt"
    @click.self="emit('close')"
  >
      
      <div class="bg-[#F8FAFC] w-full max-w-5xl h-[90vh] rounded-[32px] shadow-2xl overflow-hidden flex flex-col relative animate-scale-in">
          
          <button 
            @click="emit('close')" 
            class="absolute top-4 right-4 z-50 w-10 h-10 rounded-full bg-white/80 backdrop-blur shadow-md flex items-center justify-center text-slate-700 hover:bg-red-500 hover:text-white transition-all transform hover:rotate-90 duration-300"
          >
            <X :size="24" />
          </button>

          <div class="flex-1 overflow-y-auto custom-scrollbar bg-[#F8FAFC]">
            
            <div v-if="loading" class="h-full flex flex-col items-center justify-center gap-4">
                <span class="loading loading-spinner loading-lg text-[#1e3a8a]"></span>
            </div>

            <div v-else-if="error" class="h-full flex flex-col items-center justify-center text-center text-slate-500">
                <AlertCircle :size="48" class="text-red-400 mb-2" />
                <p>{{ error }}</p>
                <button @click="emit('close')" class="mt-4 text-[#1e3a8a] underline">ปิดหน้าต่าง</button>
            </div>

            <div v-else-if="news" class="pb-10">
                
                <div class="relative w-full h-[35vh] md:h-[45vh]">
                    <img v-if="fileUrl" :src="fileUrl" @error="handleImageError" class="w-full h-full object-cover" />
                    <div v-else class="w-full h-full bg-gradient-to-br from-slate-800 to-[#1e3a8a] flex items-center justify-center text-white/20"><ImageIcon :size="64" /></div>
                    <div class="absolute inset-0 bg-gradient-to-t from-[#F8FAFC] via-transparent to-transparent"></div>
                    
                    <div class="absolute bottom-0 left-0 right-0 p-8 pt-20">
                        <div class="max-w-4xl mx-auto">
                            <span class="px-3 py-1 bg-white/90 backdrop-blur rounded-full text-xs font-bold text-[#1e3a8a] shadow-sm mb-3 inline-block">
                                <Clock class="inline w-3 h-3 mr-1" /> {{ formatDate(news.UpdatedAt) }}
                            </span>
                            <h1 class="text-2xl md:text-3xl font-bold text-slate-800 leading-tight drop-shadow-sm">{{ news.title }}</h1>
                        </div>
                    </div>
                </div>

                <div class="max-w-4xl mx-auto px-6 md:px-8">
                    <div class="prose prose-slate max-w-none mb-10 text-slate-600 leading-relaxed whitespace-pre-line">
                        {{ news.post_detail }}
                    </div>

                    <div v-if="news.scholarship" class="bg-white rounded-3xl shadow-xl shadow-slate-200/50 overflow-hidden border border-slate-100">
                        
                        <div class="bg-[#1e3a8a] text-white p-6 md:p-8 relative overflow-hidden">
                            <div class="absolute top-0 right-0 p-6 opacity-10"><GraduationCap :size="120" /></div>
                            <div class="relative z-10 flex flex-col md:flex-row md:items-start md:justify-between gap-4">
                                <div>
                                    <p class="text-blue-200 text-xs font-medium uppercase tracking-wider mb-1">Scholarship Name</p>
                                    <h2 class="text-xl font-bold leading-snug">{{ news.scholarship.scholarship_name }}</h2>
                                    <span class="inline-block mt-3 px-3 py-1 bg-white/10 backdrop-blur rounded-lg text-sm border border-white/20">
                                        {{ news.scholarship.typescholarship?.type_name || 'ทุนทั่วไป' }}
                                    </span>
                                </div>
                                <div class="flex-shrink-0">
                                    <div 
                                        class="px-4 py-2 rounded-xl flex items-center gap-2 font-bold shadow-lg"
                                        :class="{
                                            'bg-green-500 text-white': getScholarshipStatus(news.scholarship.open_date, news.scholarship.close_date).color === 'green',
                                            'bg-red-500 text-white': getScholarshipStatus(news.scholarship.open_date, news.scholarship.close_date).color === 'red',
                                            'bg-amber-400 text-slate-900': getScholarshipStatus(news.scholarship.open_date, news.scholarship.close_date).color === 'amber',
                                            'bg-gray-400 text-white': getScholarshipStatus(news.scholarship.open_date, news.scholarship.close_date).color === 'gray'
                                        }"
                                    >
                                        <component :is="getScholarshipStatus(news.scholarship.open_date, news.scholarship.close_date).icon" :size="18" />
                                        {{ getScholarshipStatus(news.scholarship.open_date, news.scholarship.close_date).label }}
                                    </div>
                                </div>
                            </div>
                        </div>

                        <div class="p-6 md:p-8">
                            
                            <div class="bg-slate-50 rounded-2xl p-6 border border-slate-100 mb-8">
                                <h4 class="text-sm font-bold text-slate-400 mb-6 flex items-center justify-between">
                                    <span class="flex items-center gap-2"><Calendar :size="16" /> ช่วงเวลารับสมัคร</span>
                                    <span class="text-xs font-bold px-2 py-1 rounded bg-blue-100 text-[#1e3a8a]">
                                        {{ getDaysRemaining(news.scholarship.close_date) }}
                                    </span>
                                </h4>
                                <div class="relative pt-2 pb-6 px-2">
                                    <div class="absolute -top-1 left-0 transform -translate-x-1/4">
                                        <span class="text-xs font-bold text-slate-700 block text-center w-24">{{ formatDate(news.scholarship.open_date) }}</span>
                                    </div>
                                    <div class="absolute -top-1 right-0 transform translate-x-1/4">
                                        <span class="text-xs font-bold text-slate-700 block text-center w-24">{{ formatDate(news.scholarship.close_date) }}</span>
                                    </div>
                                    
                                    <div class="h-3 bg-slate-200 rounded-full w-full relative overflow-hidden mt-6 shadow-inner">
                                        <div 
                                            class="absolute top-0 left-0 h-full rounded-full transition-all duration-1000 ease-out shadow-sm"
                                            :class="{
                                                'bg-green-500': getProgressPercent(news.scholarship.open_date, news.scholarship.close_date) < 80,
                                                'bg-amber-500': getProgressPercent(news.scholarship.open_date, news.scholarship.close_date) >= 80 && getProgressPercent(news.scholarship.open_date, news.scholarship.close_date) < 100,
                                                'bg-red-500': getProgressPercent(news.scholarship.open_date, news.scholarship.close_date) >= 100
                                            }"
                                            :style="{ width: getProgressPercent(news.scholarship.open_date, news.scholarship.close_date) + '%' }"
                                        ></div>
                                    </div>

                                    <div 
                                        class="absolute top-[26px] transform -translate-x-1/2 transition-all duration-1000 z-10"
                                        :style="{ left: getProgressPercent(news.scholarship.open_date, news.scholarship.close_date) + '%' }"
                                        v-if="getProgressPercent(news.scholarship.open_date, news.scholarship.close_date) > 0 && getProgressPercent(news.scholarship.open_date, news.scholarship.close_date) < 100"
                                    >
                                        <div class="w-4 h-4 bg-white border-4 border-[#1e3a8a] rounded-full shadow-md"></div>
                                        <span class="absolute top-6 left-1/2 -translate-x-1/2 text-[10px] font-bold text-[#1e3a8a] whitespace-nowrap bg-blue-50 px-1 rounded">วันนี้</span>
                                    </div>
                                </div>
                            </div>

                            <div v-if="news.scholarship.description" class="mb-8">
                                <h4 class="font-bold text-slate-800 mb-2">รายละเอียดเพิ่มเติม</h4>
                                <p class="text-slate-600 text-sm leading-relaxed">{{ news.scholarship.description }}</p>
                            </div>

                            <div v-if="features.length > 0">
                                <div class="divider my-6"></div>
                                <h4 class="font-bold text-slate-800 mb-4 flex items-center gap-2">
                                    <CheckCircle2 class="text-[#1e3a8a]" :size="20" /> คุณสมบัติผู้สมัคร
                                </h4>
                                <ul class="grid grid-cols-1 gap-3">
                                    <li v-for="(feat, idx) in features" :key="idx" 
                                        class="flex items-start gap-3 p-3 bg-white border border-slate-100 rounded-xl shadow-sm"
                                    >
                                        <span class="flex-shrink-0 w-6 h-6 rounded-full bg-blue-50 text-[#1e3a8a] font-bold text-xs flex items-center justify-center mt-0.5">{{ idx + 1 }}</span>
                                        <span class="text-slate-700 text-sm font-medium">{{ formatCondition(feat) }}</span>
                                    </li>
                                </ul>
                            </div>
                        </div>
                    </div>

                    <div class="h-10"></div>
                </div>
            </div>
          </div>
      </div>
  </div>
</template>

<style scoped>
/* Animations */
.animate-fade-in { animation: fadeIn 0.2s ease-out; }
.animate-scale-in { animation: scaleIn 0.3s cubic-bezier(0.16, 1, 0.3, 1); }

@keyframes fadeIn { from { opacity: 0; } to { opacity: 1; } }
@keyframes scaleIn { 
  from { transform: scale(0.95) translateY(20px); opacity: 0; } 
  to { transform: scale(1) translateY(0); opacity: 1; } 
}

/* Custom Scrollbar */
.custom-scrollbar::-webkit-scrollbar { width: 6px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background-color: #cbd5e1; border-radius: 20px; }
</style>