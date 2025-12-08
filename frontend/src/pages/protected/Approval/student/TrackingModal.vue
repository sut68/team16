<script setup lang="ts">
import { ref, computed } from 'vue';
import type { ApplicationScholarshipResponse } from '@/interfaces';

// --- Props & Emits ---
const props = defineProps<{
    isOpen: boolean;
    applicationData: ApplicationScholarshipResponse | null;
}>();

const emit = defineEmits(['close', 'upload-file', 'book-interview']);

// --- State ---
const fileInput = ref<HTMLInputElement | HTMLInputElement[] | null>(null);
const isUploading = ref(false);


// --- Logic: Stages Definition ---
const stages = [
    { id: 1, title: 'ยื่นใบสมัคร', description: 'ส่งข้อมูลเข้าสู่ระบบเรียบร้อยแล้ว' },
    { id: 2, title: 'ตรวจสอบคุณสมบัติ', description: 'เจ้าหน้าที่ตรวจสอบคุณสมบัติเบื้องต้น' },
    { id: 3, title: 'อัปโหลดเอกสาร', description: 'ส่งเอกสารประกอบการพิจารณา' },
    { id: 4, title: 'ตรวจสอบเอกสาร', description: 'เจ้าหน้าที่ตรวจสอบความถูกต้องของเอกสาร' },
    { id: 5, title: 'จองคิวสัมภาษณ์', description: 'เลือกเวลาเพื่อเข้าสัมภาษณ์' }
];

// --- Logic: Determine Current State ---
const processState = computed(() => {
    if (!props.applicationData) return { currentStep: 1, status: 'normal' };
    
    // The status on ApplicationScholarship is now the source of truth
    const status = props.applicationData.status?.toLowerCase();
    
    // Mapping Status Backend -> Step ID
    if (status === 'new') return { currentStep: 2, status: 'process' }; // This status is hypothetical now
    if (status === 'qualified') return { currentStep: 3, status: 'action' }; // Application started, ready for doc upload
    if (status === 'pending') return { currentStep: 4, status: 'process' }; // Docs uploaded, pending review
    if (status === 'request-change') return { currentStep: 4, status: 'warning' };
    if (status === 'approved') return { currentStep: 5, status: 'completed' }; // Final approval
    if (status === 'rejected') return { currentStep: 0, status: 'error' };

    return { currentStep: 1, status: 'normal' };
});

const getStepStatus = (stepId: number) => {
    const { currentStep, status } = processState.value;
    
    if (stepId < currentStep) return 'completed';
    
    if (stepId === currentStep) {
        if (status === 'error') return 'rejected';
        if (status === 'warning') return 'request-change';
        if (status === 'action') return 'action-needed';
        if (status === 'completed') return 'completed';
        return 'current';
    }
    
    return 'pending';
};

// Helper to get the most recent comment
const latestComment = computed(() => {
    if (!props.applicationData?.application_documents) return '';

    // Flatten all tasks from all documents, then all decisions from all tasks
    const allTasks = props.applicationData.application_documents.flatMap(doc => doc.approval_tasks || []);
    const allDecisions = allTasks.flatMap(task => task.approval_decisions || []);
    
    if (allDecisions.length === 0) return '';

    // Safely access the last element
    const latestDecision = allDecisions[allDecisions.length - 1];
    if (latestDecision) {
      return latestDecision.comment;
    }

    return ''; // Fallback if the last element is somehow undefined
});


// --- Actions ---
const triggerUpload = () => {
  // Defensive code to handle if fileInput.value becomes an array
  if (Array.isArray(fileInput.value)) {
    fileInput.value[0]?.click();
  } else {
    fileInput.value?.click();
  }
};

const handleFileChange = (e: Event) => {
    const target = e.target as HTMLInputElement;
    if (target.files?.length) {
        isUploading.value = true;
        // The parent component will handle the actual upload and state change.
        // We just show the loading state and emit the file.
        setTimeout(() => {
            emit('upload-file', target.files![0]);
            isUploading.value = false;
        }, 1000);
    }
};
</script>

<template>
    <div v-if="isOpen && applicationData" 
        class="fixed inset-0 z-[200] flex items-center justify-center bg-black/60 backdrop-blur-sm p-4 animate-fade">
        
        <div class="bg-white w-full max-w-3xl max-h-[90vh] rounded-xl shadow-2xl flex flex-col overflow-hidden animate-pop-in">
            
            <div class="bg-[#1e3a8a] px-6 py-4 shrink-0 flex justify-between items-start">
                <div>
                    <h2 class="text-xl font-bold text-white flex items-center gap-2">
                        ติดตามสถานะ
                        <span class="bg-white/20 text-xs px-2 py-0.5 rounded text-blue-50 font-normal border border-white/10">
                            #APP-{{ applicationData.ID }}
                        </span>
                    </h2>
                    <p class="text-blue-200 text-sm mt-1">{{ applicationData.scholarship?.scholarship_name || 'รายละเอียดทุน' }}</p>
                </div>
                <button @click="$emit('close')" class="btn btn-circle btn-ghost btn-sm text-white/70 hover:bg-white/20 hover:text-white">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" /></svg>
                </button>
            </div>

            <div class="flex-1 overflow-y-auto bg-white p-6 sm:p-8">
                <h3 class="font-bold text-slate-700 mb-6 text-lg border-b pb-2">ขั้นตอนการดำเนินการ</h3>

                <div class="pl-2">
                    <div v-for="(stage, index) in stages" :key="stage.id" class="relative pb-10 last:pb-0">
                        
                        <div v-if="index !== stages.length - 1" 
                             class="absolute left-[19px] top-10 bottom-0 w-1 transition-colors duration-500 z-0"
                             :class="getStepStatus(stage.id) === 'completed' ? 'bg-green-500' : 'bg-gray-300'">
                        </div>

                        <div class="relative flex items-start gap-4 z-10">
                            <div class="relative flex-shrink-0 w-10 h-10 rounded-full flex items-center justify-center border-2 font-bold transition-all bg-white shadow-sm"
                                :class="{
                                    'border-green-500 text-green-600': getStepStatus(stage.id) === 'completed',
                                    'border-[#1e3a8a] text-[#1e3a8a] ring-4 ring-blue-50': getStepStatus(stage.id) === 'current' || getStepStatus(stage.id) === 'action-needed',
                                    'border-orange-400 text-orange-500 ring-4 ring-orange-50': getStepStatus(stage.id) === 'request-change',
                                    'border-red-500 text-red-500': getStepStatus(stage.id) === 'rejected',
                                    'border-gray-300 text-gray-400': getStepStatus(stage.id) === 'pending'
                                }">
                                <svg v-if="getStepStatus(stage.id) === 'completed'" xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" /></svg>
                                <span v-else>{{ stage.id }}</span>
                            </div>

                            <div class="flex-1 pt-1">
                                <div class="flex flex-wrap justify-between items-start gap-2">
                                    <div>
                                        <h4 class="font-bold text-base"
                                            :class="{
                                                'text-green-700': getStepStatus(stage.id) === 'completed',
                                                'text-[#1e3a8a]': ['current', 'action-needed'].includes(getStepStatus(stage.id)),
                                                'text-orange-600': getStepStatus(stage.id) === 'request-change',
                                                'text-red-600': getStepStatus(stage.id) === 'rejected',
                                                'text-gray-400': getStepStatus(stage.id) === 'pending'
                                            }">
                                            {{ stage.title }}
                                        </h4>
                                        <p class="text-sm text-gray-500 mt-0.5">{{ stage.description }}</p>
                                    </div>
                                    
                                    <div class="shrink-0">
                                        <span v-if="getStepStatus(stage.id) === 'completed'" class="badge badge-success text-white badge-sm">เสร็จสิ้น</span>
                                        <span v-else-if="['current', 'process'].includes(getStepStatus(stage.id))" class="badge badge-info text-white badge-sm animate-pulse">กำลังดำเนินการ</span>
                                        <span v-else-if="getStepStatus(stage.id) === 'action-needed'" class="badge badge-primary badge-outline badge-sm">รอคุณดำเนินการ</span>
                                        <span v-else-if="getStepStatus(stage.id) === 'request-change'" class="badge badge-warning text-white badge-sm">ต้องแก้ไข</span>
                                        <span v-else-if="getStepStatus(stage.id) === 'rejected'" class="badge badge-error text-white badge-sm">ไม่ผ่าน</span>
                                    </div>
                                </div>

                                <div class="mt-3 animate-fade-in">
                                    <!-- Hidden file input moved here to avoid duplicate refs -->
                                    <input type="file" ref="fileInput" class="hidden" @change="handleFileChange" accept=".pdf,.jpg,.png">
                                    
                                    <div v-if="stage.id === 3 && getStepStatus(3) === 'action-needed'" 
                                            class="bg-blue-50 border border-blue-100 rounded-lg p-4 mt-2">
                                        <p class="text-sm text-blue-800 mb-3 font-medium flex items-center gap-2">
                                            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
                                            คุณสมบัติผ่านเกณฑ์แล้ว
                                        </p>
                                        <button @click="triggerUpload" :disabled="isUploading" class="btn bg-[#1e3a8a] hover:bg-blue-800 text-white btn-sm border-none w-full sm:w-auto">
                                            <span v-if="isUploading" class="loading loading-spinner loading-xs"></span>
                                            <span v-else class="flex items-center gap-2">
                                                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12" /></svg>
                                                อัปโหลดเอกสาร
                                            </span>
                                        </button>
                                    </div>

                                    <div v-if="stage.id === 4 && getStepStatus(4) === 'request-change'" 
                                            class="bg-red-50 border border-red-100 rounded-lg p-4 mt-2">
                                        <div class="flex items-start gap-3 mb-3">
                                            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-red-600 mt-0.5 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" /></svg>
                                            <div>
                                                <p class="text-sm font-bold text-red-800">แจ้งเตือน: เอกสารต้องได้รับการแก้ไข</p>
                                                <div class="text-sm text-red-700 mt-1 bg-white/50 p-2 rounded border border-red-100">
                                                    "{{ latestComment || 'กรุณาตรวจสอบความถูกต้องของเอกสารแล้วอัปโหลดใหม่อีกครั้ง' }}"
                                                </div>
                                            </div>
                                        </div>
                                        <button @click="triggerUpload" :disabled="isUploading" class="btn bg-red-600 hover:bg-red-700 text-white btn-sm border-none w-full sm:w-auto">
                                            <span v-if="isUploading" class="loading loading-spinner loading-xs"></span>
                                            <span v-else>อัปโหลดเอกสารแก้ไข</span>
                                        </button>
                                    </div>

                                    <div v-if="stage.id === 5 && getStepStatus(5) === 'completed'" 
                                            class="bg-green-50 border border-green-100 rounded-lg p-4 mt-2">
                                        <div class="flex items-center gap-3 mb-3">
                                            <div class="w-10 h-10 bg-green-100 rounded-full flex items-center justify-center text-green-600 shrink-0">
                                                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
                                            </div>
                                            <div>
                                                <p class="font-bold text-green-800">ยินดีด้วย! เอกสารผ่านการอนุมัติ</p>
                                                <p class="text-xs text-green-700">กรุณาจองคิวสัมภาษณ์เพื่อดำเนินการต่อ</p>
                                            </div>
                                        </div>
                                        <button @click="$emit('book-interview')" class="btn bg-[#1e3a8a] hover:bg-blue-900 text-white btn-sm border-none shadow-lg shadow-blue-200 w-full">
                                            จองคิวสัมภาษณ์ทันที
                                        </button>
                                    </div>

                                    <div v-if="getStepStatus(stage.id) === 'rejected'" class="bg-red-50 border border-red-100 rounded-lg p-3 text-sm text-red-700 mt-2">
                                        <span class="font-bold">เหตุผลที่ไม่ผ่าน:</span> {{ latestComment || 'ไม่ผ่านการพิจารณาในขั้นตอนนี้' }}
                                    </div>

                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <div class="p-4 bg-gray-50 border-t text-center shrink-0">
                 <button @click="$emit('close')" class="btn btn-ghost w-full max-w-xs text-gray-500 font-normal">ปิดหน้าต่าง</button>
            </div>
        </div>
    </div>
</template>

<style scoped>
@keyframes pop-in { 0% { opacity: 0; transform: scale(0.95) translateY(10px); } 100% { opacity: 1; transform: scale(1) translateY(0); } }
.animate-pop-in { animation: pop-in 0.2s cubic-bezier(0.16, 1, 0.3, 1) forwards; }

@keyframes fade { from { opacity: 0; } to { opacity: 1; } }
.animate-fade { animation: fade 0.2s ease-out forwards; }

@keyframes fade-in { from { opacity: 0; height: 0; transform: translateY(-5px); } to { opacity: 1; height: auto; transform: translateY(0); } }
.animate-fade-in { animation: fade-in 0.3s ease forwards; }
</style>
