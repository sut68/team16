<script setup lang="ts">
import { ref, watch, computed, onMounted } from 'vue';
import { getNewsPostById } from '@/services/api/news_post';

const props = defineProps<{
  isOpen: boolean;
  id: number;
}>();

const emit = defineEmits(['close']);

// --- ส่วน Helper Formatter (เหมือนเดิม) ---
const formatCondition = (feature: any) => {
    // ... (Code เดิมของคุณตรงนี้ไม่ต้องแก้) ...
    const op = feature.operator;
    const val = feature.value;
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
        default: text = `${op} ${displayVal}`;
    }
    let unit = " ";
    if (typeName.includes("รายได้") || typeName.includes("เงินเดือน")) unit = " บาท";
    if (typeName) return `${typeName} ${text}${unit}`;
    return `${text}${unit}`;
};

const newsData = ref<any>(null);
const isLoading = ref(false);

const fetchData = async () => {
  if (props.id) {
    isLoading.value = true;
    try {
      const data = await getNewsPostById(props.id);
      newsData.value = data;
    } catch (error) {
      console.error("Failed to load preview:", error);
    } finally {
      isLoading.value = false;
    }
  }
};

onMounted(() => { fetchData(); });
watch(
  () => [props.id, props.isOpen], // จับตาดูทั้ง ID และสถานะเปิด/ปิด
  ([newId, newIsOpen]) => {
    if (newIsOpen && newId) {
      fetchData(); // โหลดข้อมูลเมื่อเปิด Modal และมี ID
    }
  }
);

// Helper: สร้าง URL รูปภาพ
const fileUrl = computed(() => {
  if (newsData.value?.file_path) {
    // ตัด / ตัวหน้าออกถ้ามี เพื่อป้องกัน double slash
    const cleanPath = newsData.value.file_path.startsWith('/') 
        ? newsData.value.file_path.substring(1) 
        : newsData.value.file_path;
    return `http://localhost:8080/${cleanPath}`;
  }
  return null;
});

// Helper: จัดการวันที่
const formatDate = (dateString: string) => {
    if (!dateString) return '-';
    return new Date(dateString).toLocaleDateString('th-TH', {
        year: 'numeric', month: 'long', day: 'numeric'
    });
};

// Helper: ซ่อนรูปที่ Error (เช่น กรณีไฟล์จริงเป็น PDF หรือหาไฟล์ไม่เจอ)
const handleImageError = (event: Event) => {
    const target = event.target as HTMLImageElement;
    target.style.display = 'none'; // ซ่อน tag img ไปเลย
    // ซ่อนปุ่มดูรูปเต็มด้วย (โดยการหา parent แล้วซ่อน)
    const parent = target.parentElement;
    if (parent) parent.style.display = 'none';
};
</script>

<template>
  <div v-if="isOpen" class="fixed inset-0 z-[1000] flex items-center justify-center bg-black/60 backdrop-blur-sm p-4 transition-opacity animate-fade-in">
    
    <div class="bg-white w-full max-w-4xl h-[90vh] rounded-2xl shadow-2xl overflow-hidden flex flex-col relative animate-scale-in">
      
        <button @click="emit('close')" class="absolute top-4 right-4 z-10 btn btn-circle btn-sm bg-white/90 hover:bg-white border-none shadow-md text-slate-600">✕</button>

        <div v-if="isLoading" class="flex-1 flex items-center justify-center">
            <span class="loading loading-spinner loading-lg text-[#1e3a8a]"></span>
        </div>

        <div v-else-if="newsData" class="flex-1 overflow-y-auto bg-gray-50">
            
            <div class="w-full bg-slate-200 relative group min-h-[200px] flex items-center justify-center overflow-hidden">
                
                <div v-if="fileUrl" class="relative w-full flex justify-center bg-black/5">
                    <img 
                        :src="fileUrl" 
                        @error="handleImageError" 
                        class="w-auto h-auto max-h-[500px] object-contain shadow-sm" 
                    />
                    
                    <a :href="fileUrl" target="_blank" class="absolute bottom-4 right-4 btn btn-sm bg-black/50 border-none text-white hover:bg-black/70 backdrop-blur-md opacity-0 group-hover:opacity-100 transition-opacity">
                        🔍 ดูรูปขนาดเต็ม
                    </a>
                </div>

                <div v-else class="flex flex-col items-center justify-center text-gray-400 py-10">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-12 h-12 mb-2 opacity-50">
                        <path stroke-linecap="round" stroke-linejoin="round" d="m2.25 15.75 5.159-5.159a2.25 2.25 0 0 1 3.182 0l5.159 5.159m-1.5-1.5 1.409-1.409a2.25 2.25 0 0 1 3.182 0l2.909 2.909m-18 3.75h16.5a1.5 1.5 0 0 0 1.5-1.5V6a1.5 1.5 0 0 0-1.5-1.5H3.75A1.5 1.5 0 0 0 2.25 6v12a1.5 1.5 0 0 0 1.5 1.5Zm10.5-11.25h.008v.008h-.008V8.25Zm.375 0a.375.375 0 1 1-.75 0 .375.375 0 0 1 .75 0Z" />
                    </svg>
                    <span class="text-sm font-medium">ไม่มีภาพปก</span>
                </div>
            </div>

            <div class="p-8 md:p-10 max-w-3xl mx-auto bg-white -mt-6 rounded-t-3xl relative shadow-sm min-h-[50vh]">
                
                <h1 class="text-2xl md:text-3xl font-bold text-[#1e3a8a] mb-2 leading-tight">
                    {{ newsData.title }}
                </h1>
                <p class="text-sm text-gray-400 mb-6">ประกาศเมื่อ: {{ formatDate(newsData.UpdatedAt) }}</p>

                <div class="prose prose-lg text-slate-700 whitespace-pre-wrap leading-relaxed mb-8">
                    {{ newsData.post_detail }}
                </div>

                <div v-if="newsData.scholarship" class="bg-blue-50/50 border border-blue-100 rounded-xl p-6">
                    <h3 class="text-lg font-bold text-[#1e3a8a] mb-4 flex items-center gap-2">
                        🎓 รายละเอียดทุนการศึกษา
                    </h3>
                    
                    <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
                        <div>
                            <span class="text-xs text-gray-500 block">ชื่อทุน</span>
                            <span class="font-semibold text-slate-800">{{ newsData.scholarship.scholarship_name }}</span>
                        </div>
                        <div>
                             <span class="text-xs text-gray-500 block">ประเภท</span>
                            <span class="badge badge-primary badge-outline mt-1">{{ newsData.scholarship.typescholarship?.type_name || '-' }}</span>
                        </div>
                        <div>
                            <span class="text-xs text-gray-500 block">เปิดรับสมัคร</span>
                            <span class="font-medium text-green-700">{{ formatDate(newsData.scholarship.open_date) }}</span>
                        </div>
                        <div>
                            <span class="text-xs text-gray-500 block">ปิดรับสมัคร</span>
                            <span class="font-medium text-red-600">{{ formatDate(newsData.scholarship.close_date) }}</span>
                        </div>
                    </div>

                    <div v-if="newsData.scholarship.featurescholarships?.length > 0">
                        <div class="divider my-2"></div>
                        <span class="text-sm font-bold text-slate-700 mb-2 block">คุณสมบัติผู้สมัคร:</span>
                        <ul class="space-y-2 text-sm text-slate-600 bg-white p-3 rounded-lg border border-gray-100">
                            <li v-for="feature in newsData.scholarship.featurescholarships" :key="feature.ID" class="flex items-center gap-2">
                                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="w-4 h-4 text-green-600 flex-shrink-0">
                                    <path fill-rule="evenodd" d="M16.704 4.153a.75.75 0 0 1 .143 1.052l-8 10.5a.75.75 0 0 1-1.127.075l-4.5-4.5a.75.75 0 0 1 1.06-1.06l3.894 3.893 7.48-9.817a.75.75 0 0 1 1.05-.143Z" clip-rule="evenodd" />
                                </svg>
                                <span>{{ formatCondition(feature) }}</span>
                            </li>
                        </ul>
                    </div>
                </div>

            </div>

        </div>
    </div>
  </div>
</template>

<style scoped>
.animate-fade-in { animation: fadeIn 0.2s ease-out; }
.animate-scale-in { animation: scaleIn 0.2s ease-out; }

@keyframes fadeIn { from { opacity: 0; } to { opacity: 1; } }
@keyframes scaleIn { from { transform: scale(0.95); opacity: 0; } to { transform: scale(1); opacity: 1; } }
</style>