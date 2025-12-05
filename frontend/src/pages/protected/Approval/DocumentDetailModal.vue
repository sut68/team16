<script setup lang="ts">
import { ref, watch, computed } from 'vue';
import { makeApprovalDecision } from '@/services/api/approval';
import type { ApprovalTaskResponse } from '@/interfaces';
const props = defineProps<{
    isOpen: boolean;
    documentData: ApprovalTaskResponse | null;
}>();
const emit = defineEmits(['close', 'action-completed']);
const comment = ref('');
const actionType = ref<'approve' | 'reject' | 'request-change' | null>(null);
const isSubmitting = ref(false);
const submissionError = ref<string | null>(null);
const timelineEvents = ref([
    {
        id: 3,
        title: 'รอการตรวจสอบ (ฉบับปัจจุบัน)',
        date: new Date().toLocaleDateString('th-TH'),
        description: 'เอกสารอยู่ในสถานะรอพิจารณา',
        actor: 'ระบบ',
        status: 'current',
        type: 'user',
    }
]);
watch(() => props.isOpen, (newValue) => {
    if (newValue) {
        comment.value = '';
        actionType.value = null;
        isSubmitting.value = false;
        submissionError.value = null;
    }
});
const canAction = computed(() => {
    if (!props.documentData?.status) return false;
    const s = props.documentData.status.toLowerCase();
    return s === 'pending' || s === 'request-change';
});

const closeModal = () => {
    if (isSubmitting.value) return;
    emit('close');
};

const submitAction = async (type: 'approve' | 'reject' | 'request-change') => {
    if (!props.documentData || !type) return;

    if ((type === 'reject' || type === 'request-change') && !comment.value.trim()) {
        submissionError.value = 'กรุณาระบุเหตุผลประกอบการพิจารณา';
        return;
    }

    isSubmitting.value = true;
    submissionError.value = null;

    try {
        await makeApprovalDecision({
            task_id: props.documentData.ID,
            decision: type,
            comment: comment.value
        });

        emit('action-completed');
        closeModal();
    } catch (error) {
        console.error('Failed to submit action:', error);
        submissionError.value = 'เกิดข้อผิดพลาดในการบันทึกข้อมูล กรุณาลองใหม่อีกครั้ง';
    } finally {
        isSubmitting.value = false;
    }
};
</script>

<template>
    <div v-if="isOpen && documentData"
        class="fixed inset-0 z-[200] flex items-center justify-center bg-black/50 backdrop-blur-sm p-4 transition-opacity">
        <div
            class="bg-white w-full max-w-6xl h-[90vh] rounded-2xl shadow-2xl flex flex-col overflow-hidden animate-pop-in">
            <div class="px-6 py-4 border-b flex items-center justify-between bg-slate-50">
                <div>
                    <h2 class="text-xl font-bold text-[#1e3a8a] flex items-center gap-2">
                        {{ documentData.approval_requirement?.scholarship?.scholarship_name || 'รายละเอียดทุนการศึกษา'
                        }}

                        <span class="badge badge-warning text-white"
                            v-if="documentData.status?.toLowerCase() === 'pending'">รอตรวจสอบ</span>
                        <span class="badge badge-warning text-white"
                            v-else-if="documentData.status?.toLowerCase() === 'request-change'">รอแก้ไข</span>
                        <span class="badge badge-success text-white"
                            v-else-if="documentData.status?.toLowerCase() === 'approved'">อนุมัติแล้ว</span>
                        <span class="badge badge-error text-white"
                            v-else-if="documentData.status?.toLowerCase() === 'rejected'">ปฏิเสธ</span>
                    </h2>
                    <div class="flex items-center gap-2 text-sm text-gray-500 mt-1">
                        <span>Task ID: #{{ documentData.ID }}</span>
                        <span>•</span>
                        <span class="bg-blue-50 text-blue-700 px-2 rounded border border-blue-100 text-xs">
                            รอบ: {{ documentData.round || '1/2568' }}
                        </span>
                    </div>
                </div>
                <button @click="closeModal" class="btn btn-circle btn-ghost btn-sm text-gray-500 hover:bg-gray-200">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24"
                        stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M6 18L18 6M6 6l12 12" />
                    </svg>
                </button>
            </div>
            <div class="flex-1 overflow-y-auto p-6 bg-[#f8fafc]">
                <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
                    <div class="lg:col-span-2 space-y-6">

                        <div class="card bg-white shadow-sm border border-gray-100">
                            <div class="card-body p-5">
                                <h3 class="font-bold text-lg text-slate-700 mb-4 border-b pb-2">ข้อมูลผู้สมัคร</h3>

                                <div v-if="documentData.application?.student_profile"
                                    class="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm">
                                    <div>
                                        <span class="block text-gray-400">ชื่อ-นามสกุล</span>
                                        <span class="font-semibold text-slate-800 text-lg">
                                            {{ documentData.application.student_profile.first_name_th }}
                                            {{ documentData.application.student_profile.last_name_th }}
                                        </span>
                                    </div>
                                    <div>
                                        <span class="block text-gray-400">รหัสนักศึกษา</span>
                                        <span class="font-semibold text-slate-800">
                                            {{ documentData.application.student_profile.student_id || '-' }}
                                        </span>
                                    </div>
                                    <div>
                                        <span class="block text-gray-400">สาขาวิชา</span>
                                        <span class="font-semibold text-slate-800">
                                            {{ documentData.application.student_profile.major?.major_name || 'N/A' }}
                                        </span>
                                    </div>
                                    <div>
                                        <span class="block text-gray-400">เกรดเฉลี่ย (GPAX)</span>
                                        <span class="font-semibold text-slate-800">
                                            {{ documentData.application.student_profile.gpax ?
                                                documentData.application.student_profile.gpax.toFixed(2) : '-' }}
                                        </span>
                                    </div>
                                </div>
                                <div v-else class="text-center text-gray-400 py-4">ไม่พบข้อมูลผู้สมัคร</div>
                            </div>
                        </div>

                        <div class="card bg-white shadow-sm border border-gray-100">
                            <div class="card-body p-5">
                                <h3 class="font-bold text-lg text-slate-700 mb-4 border-b pb-2">เอกสารแนบ (ฉบับล่าสุด)
                                </h3>
                                <div class="space-y-3">
                                    <div v-if="documentData.application_document"
                                        class="flex items-center justify-between p-3 bg-slate-50 rounded-lg border border-gray-200">
                                        <div class="flex items-center gap-3">
                                            <div
                                                class="w-10 h-10 bg-red-100 rounded-lg flex items-center justify-center text-red-500 font-bold text-xs">
                                                PDF
                                            </div>
                                            <div>
                                                <p class="font-medium text-slate-700 truncate max-w-[300px]">
                                                    {{ documentData.application_document.file_name }}
                                                </p>
                                                <p class="text-xs text-gray-400">
                                                    อัปโหลดโดย {{ documentData.application_document.uploaded_by }}
                                                </p>
                                            </div>
                                        </div>
                                        <button
                                            class="btn btn-sm btn-ghost text-[#1e3a8a] hover:bg-blue-50">ดูไฟล์</button>
                                    </div>
                                    <div v-else class="text-center text-gray-400 py-4">ไม่พบเอกสารแนบ</div>
                                </div>
                            </div>
                        </div>
                    </div>
                    <div class="lg:col-span-1 space-y-6">
                        <div v-if="canAction"
                            class="card bg-white shadow-md border border-blue-100 ring-4 ring-blue-50/50">
                            <div class="card-body p-5">
                                <h3 class="font-bold text-lg text-slate-700 mb-4">ผลการพิจารณา</h3>

                                <div class="flex flex-col gap-2">
                                    <button @click="submitAction('approve')" :disabled="isSubmitting"
                                        class="btn bg-[#1e3a8a] hover:bg-[#152c6f] text-white w-full border-none">
                                        <span v-if="isSubmitting" class="loading loading-spinner loading-sm"></span>
                                        อนุมัติเอกสาร (Approve)
                                    </button>

                                    <div class="grid grid-cols-2 gap-2 mt-2">
                                        <button @click="actionType = 'request-change'" :disabled="isSubmitting"
                                            class="btn btn-outline btn-warning btn-sm hover:text-white">
                                            ขอแก้ไข
                                        </button>
                                        <button @click="actionType = 'reject'" :disabled="isSubmitting"
                                            class="btn btn-outline btn-error btn-sm hover:text-white">
                                            ปฏิเสธ
                                        </button>
                                    </div>
                                </div>

                                <div v-if="actionType" class="mt-4 pt-4 border-t animate-fade-in">
                                    <p class="text-sm font-bold mb-2 flex items-center gap-1"
                                        :class="actionType === 'reject' ? 'text-error' : 'text-warning'">
                                        {{ actionType === 'reject' ? 'ระบุเหตุผลการปฏิเสธ' : 'ระบุสิ่งที่ต้องแก้ไข' }}
                                    </p>

                                    <textarea v-model="comment"
                                        class="textarea textarea-bordered w-full h-24 text-sm focus:ring-2 focus:ring-opacity-50"
                                        :class="actionType === 'reject' ? 'focus:border-error focus:ring-error' : 'focus:border-warning focus:ring-warning'"
                                        placeholder="พิมพ์รายละเอียด..."></textarea>

                                    <div v-if="submissionError"
                                        class="text-error text-xs mt-2 font-medium bg-red-50 p-2 rounded">
                                        {{ submissionError }}
                                    </div>

                                    <div class="flex justify-end gap-2 mt-3">
                                        <button @click="actionType = null"
                                            class="btn btn-ghost btn-xs text-gray-500">ยกเลิก</button>
                                        <button @click="submitAction(actionType!)"
                                            class="btn btn-sm text-white border-none"
                                            :class="actionType === 'reject' ? 'btn-error' : 'btn-warning'"
                                            :disabled="isSubmitting">
                                            <span v-if="isSubmitting" class="loading loading-spinner loading-xs"></span>
                                            ยืนยัน
                                        </button>
                                    </div>
                                </div>
                            </div>
                        </div>
                        <div v-else class="card bg-white shadow-sm border border-gray-100">
                            <div class="card-body p-5 items-center text-center">
                                <div class="badge badge-lg p-4 font-bold text-white mb-2"
                                    :class="documentData.status?.toLowerCase() === 'approved' ? 'badge-success' : 'badge-error'">
                                    {{ documentData.status?.toLowerCase() === 'approved' ? 'อนุมัติเรียบร้อยแล้ว' :
                                    'สิ้นสุดการดำเนินการ' }}
                                </div>
                                <p class="text-xs text-gray-400">สถานะปัจจุบัน: {{ documentData.status }}</p>
                            </div>
                        </div>

                        <!-- Timeline Section -->
                        <div class="card bg-white shadow-sm border border-gray-100">
                            <div class="card-body p-5">
                                <h3 class="font-bold text-lg text-slate-700 mb-4">ประวัติการดำเนินการ</h3>

                                <ul class="timeline timeline-vertical timeline-compact -ml-4">
                                    <li v-for="(event, index) in timelineEvents" :key="event.id">
                                        <hr v-if="index !== timelineEvents.length - 1" class="bg-gray-200" />
                                        <div class="timeline-middle">
                                            <div v-if="event.status === 'past'"
                                                class="w-6 h-6 rounded-full bg-blue-100 text-blue-600 flex items-center justify-center">
                                                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20"
                                                    fill="currentColor" class="w-4 h-4">
                                                    <path fill-rule="evenodd"
                                                        d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                                                        clip-rule="evenodd" />
                                                </svg>
                                            </div>
                                            <div v-else
                                                class="w-6 h-6 rounded-full bg-[#1e3a8a] border-4 border-blue-100">
                                            </div>
                                        </div>

                                        <div
                                            class="timeline-end timeline-box w-full border-none shadow-none p-0 pl-2 mb-6">
                                            <div class="font-bold text-slate-800 text-sm">{{ event.title }}</div>
                                            <div class="text-xs text-gray-500 mb-1">{{ event.date }} • โดย {{
                                                event.actor }}</div>
                                            <div class="text-xs text-gray-600">{{ event.description }}</div>
                                        </div>
                                        <hr v-if="index !== 0" class="bg-gray-200" />
                                    </li>
                                </ul>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
@keyframes pop-in {
    0% {
        opacity: 0;
        transform: scale(0.95) translateY(10px);
    }

    100% {
        opacity: 1;
        transform: scale(1) translateY(0);
    }
}

.animate-pop-in {
    animation: pop-in 0.2s cubic-bezier(0.16, 1, 0.3, 1) forwards;
}

.animate-fade-in {
    animation: fade 0.3s ease forwards;
}

@keyframes fade {
    from {
        opacity: 0;
        height: 0;
    }

    to {
        opacity: 1;
        height: auto;
    }
}
</style>