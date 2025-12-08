<script setup lang="ts">
import { ref, watch, computed } from 'vue';
import { makeApprovalDecision } from '@/services/api/approval';
// Import Interfaces
import type { ApprovalTaskResponse } from '@/interfaces';

interface TimelineEvent {
    id: number | string;
    title: string;
    date: string;
    description: string;
    actor: string;
    status: string; // current, past, past-approved, past-rejected, past-request-change, past-submitted
    type: 'task' | 'decision';
    timestamp: number;
}

const props = defineProps<{
    isOpen: boolean;
    documentData: ApprovalTaskResponse | null;
}>();

const emit = defineEmits(['close', 'action-completed']);

// State
const comment = ref('');
const actionType = ref<'approve' | 'reject' | 'request-change' | null>(null);
const isSubmitting = ref(false);
const submissionError = ref<string | null>(null);

// ✅ Helper: แปลงวันที่จาก Go Time String เป็น JS Date Object
const parseGoDate = (dateString: string | undefined): Date | null => {
    if (!dateString) return null;
    
    let dateStr: string = String(dateString);
    
    // 1. ตัดส่วน Monotonic clock (m=+...) ออก
    if (dateStr.includes(' m=')) {
        const parts: string[] = dateStr.split(' m=');
        dateStr = parts[0] as string;
    }
    
    // 2. แปลงรูปแบบให้เป็น ISO 8601
    // ตัวอย่าง: "2025-12-05 10:47:51.344406723 +0000 UTC" -> "2025-12-05T10:47:51.344Z"
    if (dateStr.includes('+0000 UTC')) {
        dateStr = dateStr.replace(' +0000 UTC', 'Z').replace(' ', 'T');
    }
    
    // 3. ตัด Nanoseconds ให้เหลือ Milliseconds (3 หลัก)
    dateStr = dateStr.replace(/(\.\d{3})\d+/, '$1');

    const date = new Date(dateStr);
    return isNaN(date.getTime()) ? null : date;
};

// Helper: Format Date (Thai)
const formatDate = (dateString: string | undefined | number) => {
    if (!dateString) return 'เมื่อสักครู่';
    
    // ถ้าเป็น number (timestamp) ให้แปลงเลย
    if (typeof dateString === 'number') {
        return new Date(dateString).toLocaleDateString('th-TH', {
            year: 'numeric',
            month: 'short',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit',
        });
    }

    // ถ้าเป็น string ให้ผ่าน parseGoDate ก่อน
    const date = parseGoDate(dateString);
    if (!date) return 'เมื่อสักครู่'; 
    
    return date.toLocaleDateString('th-TH', {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit',
    });
};

// Check if user can perform action
const canAction = computed(() => {
    if (!props.documentData?.status) return false;
    const s = props.documentData.status.toLowerCase();
    return s === 'pending' || s === 'request-change';
});

// ✅ Timeline Logic
const processedTimelineEvents = computed(() => {
    if (!props.documentData) return [];

    const events: TimelineEvent[] = [];
    const doc = props.documentData;
    
    // Determine the actor's name once with explicit checks for TypeScript
    const adminActorName = (doc.admin_profile && doc.admin_profile.admin_firstname) ? doc.admin_profile.admin_firstname : 'ระบบ';
    
    // Parse timestamp สำหรับการยื่นใบสมัคร
    const createdDate = parseGoDate(doc.CreatedAt);
    // ถ้า parse ไม่ได้ ให้ใช้เวลาปัจจุบันลบไปนิดหน่อย เพื่อให้มันไปอยู่ล่างสุด (เก่าสุด)
    let createdTs = createdDate ? createdDate.getTime() : Date.now() - 100000;

    // 1. Base Event: ยื่นใบสมัคร (เสมอ)
    events.push({
        id: 'created',
        title: 'ยื่นใบสมัครแล้ว',
        date: formatDate(doc.CreatedAt),
        description: 'ส่งเอกสารเข้าสู่ระบบเรียบร้อยแล้ว',
        actor: 'ผู้สมัคร',
        status: 'past-submitted', 
        type: 'task',
        timestamp: createdTs,
    });

    // 2. Decision Events: ประวัติการตัดสินใจ
    if (doc.approval_decisions && doc.approval_decisions.length > 0) {
        doc.approval_decisions.forEach(decision => {
            let title = '';
            let status = '';
            let description = decision.comment || '';
            const dType = decision.decision?.toLowerCase() || '';

            // ✅ กรอง "ดำเนินการ" ออก: ถ้าไม่ใช่สถานะหลัก (approve, reject, request-change) ให้ข้ามไปเลย
            if (!['approve', 'reject', 'request-change'].includes(dType)) {
                return;
            }

            switch (dType) {
                case 'approve':
                    title = 'อนุมัติแล้ว';
                    status = 'past-approved';
                    description = description || 'เอกสารได้รับการอนุมัติ';
                    break;
                case 'reject':
                    title = 'ปฏิเสธคำขอ';
                    status = 'past-rejected';
                    break;
                case 'request-change':
                    title = 'ขอข้อมูลเพิ่มเติม/แก้ไข';
                    status = 'past-request-change';
                    break;
            }

            // Parse timestamp ของการตัดสินใจ
            const decisionDate = parseGoDate(decision.decision_at);
            const ts = decisionDate ? decisionDate.getTime() : Date.now();
            
            events.push({
                id: decision.ID,
                title: title,
                date: formatDate(decision.decision_at),
                description: description,
                actor: doc.admin_profile ? doc.admin_profile.admin_firstname : 'เจ้าหน้าที่',
                status: status,
                type: 'decision',
                timestamp: ts, 
            });
        });
    }

    // 3. Handle Current Status & Missing Events
    const currentStatus = doc.status?.toLowerCase() || '';
    const nowTs = Date.now();

    // กรณีรอดำเนินการ (Pending)
    if (currentStatus === 'pending') {
        // หาวันที่อัปเดตล่าสุดที่แท้จริง
        const updateDate = parseGoDate(doc.UpdatedAt);
        const updateTs = updateDate ? updateDate.getTime() : nowTs;

        events.push({
            id: 'reviewing',
            title: 'เจ้าหน้าที่กำลังตรวจสอบ',
            date: formatDate(doc.UpdatedAt || doc.CreatedAt), 
            description: 'อยู่ในระหว่างการพิจารณาตรวจสอบความถูกต้อง',
            actor: adminActorName, 
            status: 'current',
            type: 'task',
            timestamp: updateTs > createdTs ? updateTs : createdTs + 1000, // ให้มั่นใจว่าอยู่หลัง CreatedAt
        });
    }
    // กรณีอนุมัติ (Approved) แต่ไม่มี Event ใน Timeline
    else if (currentStatus === 'approved') {
        const hasApprovedEvent = events.some(e => e.status === 'past-approved');
        if (!hasApprovedEvent) {
            events.push({
                id: 'auto-approved',
                title: 'อนุมัติแล้ว',
                date: formatDate(doc.UpdatedAt || doc.CreatedAt),
                description: 'เอกสารได้รับการอนุมัติ (สิ้นสุดกระบวนการ)',
                actor: 'เจ้าหน้าที่',
                status: 'past-approved',
                type: 'task',
                timestamp: nowTs,
            });
        }
    }
    // กรณีปฏิเสธ (Rejected) แต่ไม่มี Event ใน Timeline
    else if (currentStatus === 'rejected') {
        const hasRejectedEvent = events.some(e => e.status === 'past-rejected');
        if (!hasRejectedEvent) {
            events.push({
                id: 'auto-rejected',
                title: 'ปฏิเสธคำขอ',
                date: formatDate(doc.UpdatedAt || doc.CreatedAt),
                description: 'คำขอถูกปฏิเสธ',
                actor: 'เจ้าหน้าที่',
                status: 'past-rejected',
                type: 'task',
                timestamp: nowTs,
            });
        }
    }

    // Sort Descending (Newest First) - ล่าสุดอยู่บน
    events.sort((a, b) => b.timestamp - a.timestamp);

    // 4. Finalize Statuses
    if (events.length > 0) {
        const firstEvent = events[0];
        if (firstEvent && firstEvent.status === 'past') {
            firstEvent.status = 'current';
        }
        for(let i=1; i<events.length; i++) {
            const event = events[i];
            if (event && event.status === 'current') {
                event.status = 'past';
            }
        }
    }

    return events;
});

// Reset state
watch(() => props.isOpen, (newValue) => {
    if (newValue) {
        comment.value = '';
        actionType.value = null;
        isSubmitting.value = false;
        submissionError.value = null;
    }
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

const openDocument = () => {
    if (props.documentData?.application_document?.file_path) {
        const backendBaseUrl = 'http://localhost:8080'; 
        const fileUrl = `${backendBaseUrl}/uploads/${props.documentData.application_document.file_path}`;
        window.open(fileUrl, '_blank');
    } else {
        alert('ไม่พบไฟล์เอกสาร');
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
                                        <button @click="openDocument"
                                            class="btn btn-sm btn-ghost text-[#1e3a8a] hover:bg-blue-50">ดูไฟล์</button>
                                    </div>
                                    <div v-else class="text-center text-gray-400 py-4">ไม่พบเอกสารแนบ</div>
                                </div>
                            </div>
                        </div>

                        <div class="card bg-white shadow-sm border border-gray-100">
                            <div class="card-body p-5">
                                <h3 class="font-bold text-lg text-slate-700 mb-4 border-b pb-2">
                                    ข้อกำหนดทั้งหมดของทุนนี้
                                </h3>
                                
                                <div v-if="documentData.approval_requirement?.scholarship?.approval_requirements?.length">
                                    <ul class="space-y-3 list-disc list-inside text-slate-600 text-sm">
                                        <li v-for="req in documentData.approval_requirement.scholarship.approval_requirements" :key="req.ID">
                                            {{ req.description }}
                                        </li>
                                    </ul>
                                </div>
                                
                                <div v-else class="text-center text-gray-400 py-4">
                                    ไม่พบข้อมูลข้อกำหนด
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

                        <div class="card bg-white shadow-sm border border-gray-100">
                            <div class="card-body p-5">
                                <h3 class="font-bold text-lg text-slate-700 mb-4">ประวัติการดำเนินการ</h3>

                                <ul class="timeline timeline-vertical timeline-compact -ml-4">
                                    <li v-for="(event, index) in processedTimelineEvents" :key="event.id">
                                        
                                        <hr v-if="index > 0" class="bg-gray-200" />

                                        <div class="timeline-middle">
                                            <div v-if="event.status === 'current'"
                                                class="relative flex items-center justify-center w-6 h-6">
                                                <span class="absolute inline-flex h-full w-full rounded-full bg-blue-400 opacity-75 animate-ping"></span>
                                                <span class="relative inline-flex rounded-full h-6 w-6 bg-[#1e3a8a] border-4 border-blue-100"></span>
                                            </div>
                                            
                                            <div v-else-if="event.status === 'past-approved' || event.status === 'past-submitted'"
                                                class="w-6 h-6 rounded-full bg-green-500 text-white flex items-center justify-center">
                                                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20"
                                                    fill="currentColor" class="w-4 h-4">
                                                    <path fill-rule="evenodd"
                                                        d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                                                        clip-rule="evenodd" />
                                                </svg>
                                            </div>

                                            <div v-else-if="event.status === 'past-rejected'"
                                                class="w-6 h-6 rounded-full bg-red-500 text-white flex items-center justify-center">
                                                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20"
                                                    fill="currentColor" class="w-4 h-4">
                                                    <path fill-rule="evenodd"
                                                        d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                                                        clip-rule="evenodd" />
                                                </svg>
                                            </div>

                                            <div v-else-if="event.status === 'past-request-change'"
                                                class="w-6 h-6 rounded-full bg-orange-500 text-white flex items-center justify-center">
                                                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20"
                                                    fill="currentColor" class="w-4 h-4">
                                                    <path fill-rule="evenodd"
                                                        d="M13.586 3.586a2 2 0 112.828 2.828l-.793.793-2.828-2.828.793-.793zM11.379 5.793L3 14.172V17h2.828l8.38-8.379-2.83-2.828z"
                                                        clip-rule="evenodd" />
                                                </svg>
                                            </div>

                                            <div v-else
                                                class="w-6 h-6 rounded-full bg-gray-300 text-white flex items-center justify-center">
                                                <div class="w-2 h-2 bg-white rounded-full"></div>
                                            </div>
                                        </div>

                                        <div class="timeline-end timeline-box w-full border-none shadow-none p-0 pl-2 mb-6">
                                            <div class="font-bold text-slate-800 text-sm"
                                                :class="{'text-[#1e3a8a]': event.status === 'current'}">
                                                {{ event.title }}
                                            </div>
                                            <div class="text-xs text-gray-500 mb-1">
                                                {{ event.date }} • โดย {{ event.actor }}
                                            </div>
                                            <div v-if="(event.status === 'past-rejected' || event.status === 'past-request-change') && event.description" 
                                                 class="mt-2 p-2 rounded-lg text-xs border"
                                                 :class="event.status === 'past-rejected' ? 'bg-red-50 text-red-700 border-red-200' : 'bg-orange-50 text-orange-700 border-orange-200'">
                                                {{ event.description }}
                                            </div>
                                            <div v-else class="text-xs text-gray-600 break-words">
                                                {{ event.description }}
                                            </div>
                                        </div>
                                        
                                        <hr v-if="index < processedTimelineEvents.length - 1" class="bg-gray-200" />
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