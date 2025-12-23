<script setup lang="ts">
import { ref, watch, computed } from 'vue';
import { makeApprovalDecision } from '@/services/api/approval';
import { getMyProfile } from '@/services/api/user';
import type { ApprovalTaskResponse, SemasterResponse } from '@/interfaces';

interface ApprovalTaskDisplay extends ApprovalTaskResponse {
    roundText?: string;
    submission_date?: string;
    semaster?: SemasterResponse;
}

interface TimelineEvent {
    id: number | string;
    title: string;
    date: string;
    description: string;
    actor: string;
    status: string;
    type: 'task' | 'decision';
    timestamp: number;
    filePath?: string;
    fileName?: string;
    isOldFile?: boolean;
}

const props = defineProps<{
    isOpen: boolean;
    documentData: ApprovalTaskDisplay | null;
}>();

const emit = defineEmits(['close', 'action-completed']);

const comment = ref('');
const actionType = ref<'approve' | 'reject' | 'request-change' | null>(null);
const isSubmitting = ref(false);
const submissionError = ref<string | null>(null);
const adminId = ref<number | null>(null);


const fetchAdminProfile = async () => {
    submissionError.value = null; 
    try {
        const profile = await getMyProfile();
        
        if (!profile) {
            submissionError.value = 'ไม่สามารถโหลดข้อมูลโปรไฟล์ได้ (Profile is null)';
            return;
        }
        if (!profile.data) {
            submissionError.value = 'ข้อมูลโปรไฟล์ไม่สมบูรณ์ (Data is null)';
            return;
        }
        if (!profile.role) {
            submissionError.value = 'ไม่สามารถระบุสิทธิ์ผู้ใช้ได้ (Role is null)';
            return;
        }

        const role = profile.role.toLowerCase();
        if (!['admin', 'super_admin', 'staff'].includes(role)) {
            submissionError.value = `สิทธิ์ผู้ใช้ของคุณ ("${profile.role}") ไม่ใช่ Admin`;
            return;
        }
        
        const id = (profile.data as any).ID || (profile.data as any).id;
        if (!id && id !== 0) {
            submissionError.value = 'ไม่พบ ID ของผู้ดูแลระบบในข้อมูลโปรไฟล์';
            return;
        }

        adminId.value = id;
    } catch (error) {
        submissionError.value = 'เกิดข้อผิดพลาดในการดึงข้อมูลโปรไฟล์';
    }
};


const semesterText = computed(() => {
    const semaster = props.documentData?.semaster;
    if (!semaster) return 'N/A';
    return `ปีการศึกษา ${semaster.academic_year} เทอม ${semaster.term} รอบ ${semaster.round}`;
});

const parseGoDate = (dateString: string | undefined): Date | null => {
    if (!dateString) return null;
    const date = new Date(dateString);
    return isNaN(date.getTime()) ? null : date;
};

const formatDate = (dateString: string | undefined | number) => {
    if (!dateString) return 'เมื่อสักครู่';
    if (typeof dateString === 'number') {
        return new Date(dateString).toLocaleDateString('th-TH', {
            year: 'numeric', month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit',
        });
    }
    const date = parseGoDate(dateString);
    if (!date) return 'เมื่อสักครู่';
    return date.toLocaleDateString('th-TH', {
        year: 'numeric', month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit',
    });
};

const canAction = computed(() => {
    if (!props.documentData?.status) return false;
    const s = props.documentData.status.toLowerCase();
    if (s === 'approved' || s === 'rejected') return false;
    return true;
});

const latestDocument = computed(() => {
    const currentDoc = props.documentData?.application_document;
    if (!currentDoc) return null;

    const appScholarship = currentDoc.application_scholarship;
    const allDocs = appScholarship?.application_documents || appScholarship?.application_documents;

    if (allDocs && Array.isArray(allDocs) && allDocs.length > 0) {
        const sortedDocs = [...allDocs].sort((a, b) =>
            (parseGoDate(b.CreatedAt)?.getTime() || 0) - (parseGoDate(a.CreatedAt)?.getTime() || 0)
        );
        return sortedDocs[0];
    }
    return currentDoc;
});

const isResubmitted = computed(() => {
    if (!props.documentData?.approval_decisions || !latestDocument.value) return false;

    const uploadTime = parseGoDate(latestDocument.value.CreatedAt)?.getTime() || 0;
    const decisions = [...props.documentData.approval_decisions];

    if (decisions.length === 0) return false;

    decisions.sort((a, b) => (Number(b.ID) || 0) - (Number(a.ID) || 0));
    const lastDecision = decisions[0];
    if (!lastDecision) return false;

    const decisionTime = parseGoDate(lastDecision.decision_at)?.getTime() || 0;
    return lastDecision.decision === 'request-change' && uploadTime > decisionTime;
});

const processedTimelineEvents = computed(() => {
    if (!props.documentData || !props.documentData.application_document) return [];

    const events: TimelineEvent[] = [];
    const doc = props.documentData;
    const currentDoc = doc.application_document;

    const appScholarship = currentDoc.application_scholarship;
    const application = appScholarship?.application;
    let studentName = 'ผู้สมัคร';
    if (application?.student_profile) {
        studentName = `${application.student_profile.first_name_th} ${application.student_profile.last_name_th}`;
    }
    const defaultAdminName = 'เจ้าหน้าที่';

    if (appScholarship) {
        const applyDateRaw = appScholarship.CreatedAt;
        const applyTs = parseGoDate(applyDateRaw)?.getTime();
        if (applyTs) {
            events.push({
                id: `app-scholarship-${appScholarship.ID}`,
                title: 'ยื่นใบสมัครทุน',
                date: formatDate(applyDateRaw),
                description: 'นักศึกษายื่นความประสงค์ขอรับทุนการศึกษา',
                actor: studentName,
                status: 'past-submitted',
                type: 'task',
                timestamp: applyTs,
            });
        }
    }

    const allDocsRaw = appScholarship?.application_documents || appScholarship?.application_documents;

    if (allDocsRaw && Array.isArray(allDocsRaw) && allDocsRaw.length > 0) {
        const sortedDocs = [...allDocsRaw].sort((a, b) =>
            (parseGoDate(a.CreatedAt)?.getTime() || 0) - (parseGoDate(b.CreatedAt)?.getTime() || 0)
        );

        sortedDocs.forEach((d, index) => {
            const uploadTs = parseGoDate(d.CreatedAt)?.getTime();
            if (uploadTs) {
                const isFirstUpload = index === 0;
                const isCurrentFile = d.ID === latestDocument.value?.ID;

                let titleLabel = isFirstUpload ? 'อัปโหลดเอกสาร (ครั้งแรก)' : 'อัปโหลดเอกสารแก้ไข';
                if (isFirstUpload && isResubmitted.value && isCurrentFile) {
                    titleLabel = 'อัปโหลดเอกสารแก้ไข';
                }

                events.push({
                    id: `doc-history-${d.ID}`,
                    title: titleLabel,
                    date: formatDate(d.CreatedAt),
                    description: '',
                    actor: studentName,
                    status: 'past-submitted',
                    type: 'task',
                    timestamp: uploadTs,
                    filePath: d.file_path,
                    fileName: d.file_name,
                    isOldFile: !isCurrentFile
                });
            }
        });
    } else {
        const uploadDateRaw = currentDoc.CreatedAt;
        const uploadTs = parseGoDate(uploadDateRaw)?.getTime();
        let titleLabel = isResubmitted.value ? 'อัปโหลดเอกสารแก้ไข' : 'อัปโหลดเอกสาร';

        if (uploadDateRaw) {
            events.push({
                id: `doc-${currentDoc.ID}`,
                title: titleLabel,
                date: formatDate(uploadDateRaw),
                description: '',
                actor: studentName,
                status: 'past-submitted',
                type: 'task',
                timestamp: uploadTs || Date.now(),
                filePath: currentDoc.file_path,
                fileName: currentDoc.file_name,
                isOldFile: false
            });
        }
    }

    if (doc.approval_decisions && doc.approval_decisions.length > 0) {
        doc.approval_decisions.forEach(decision => {
            const dateRaw = decision.decision_at;
            const decisionTs = parseGoDate(dateRaw)?.getTime();
            if (!decisionTs) return;

            let title = '';
            let status = '';
            let description = decision.comment || '';
            const type = decision.decision?.toLowerCase();
            const adminName = decision.admin_profile?.admin_firstname || defaultAdminName;


            switch (type) {
                case 'approve':
                    title = 'อนุมัติแล้ว';
                    status = 'past-approved';
                    description = description || 'เอกสารได้รับการอนุมัติ';
                    break;
                case 'reject':
                    title = 'ปฏิเสธคำขอ';
                    status = 'past-rejected';
                    description = `เหตุผล: ${description || 'ไม่ระบุ'}`;
                    break;
                case 'request-change':
                    title = isResubmitted.value ? 'ขอข้อมูลเพิ่มเติม (ส่งแก้ไขแล้ว)' : 'ขอข้อมูลเพิ่มเติม/แก้ไข';
                    status = 'past-request-change';
                    description = `สิ่งที่ต้องแก้ไข: ${description || 'ไม่ระบุ'}`;
                    break;
                default:
                    return;
            }

            events.push({
                id: decision.ID,
                title: title,
                date: formatDate(dateRaw),
                description: description,
                actor: adminName,
                status: status,
                type: 'decision',
                timestamp: decisionTs,
            });
        });
    }

    const currentStatus = doc.status?.toLowerCase();
    const hasDecisions = doc.approval_decisions && doc.approval_decisions.length > 0;

    if (isResubmitted.value) {
        events.push({
            id: 'reviewing-new',
            title: 'รอตรวจสอบการแก้ไข',
            date: 'กำลังดำเนินการ',
            description: 'ผู้สมัครได้ส่งเอกสารแก้ไขเข้ามาใหม่แล้ว รอเจ้าหน้าที่ตรวจสอบ',
            actor: defaultAdminName,
            status: 'current',
            type: 'task',
            timestamp: Date.now() + 10000,
        });
    }
    else if (currentStatus === 'pending' && !hasDecisions) {
        const reviewingTs = parseGoDate(doc.CreatedAt)?.getTime();
        if (reviewingTs) {
            events.push({
                id: 'reviewing',
                title: 'เจ้าหน้าที่กำลังตรวจสอบ',
                date: formatDate(doc.CreatedAt),
                description: 'อยู่ในระหว่างการพิจารณาตรวจสอบความถูกต้อง',
                actor: defaultAdminName,
                status: 'current',
                type: 'task',
                timestamp: reviewingTs,
            });
        }
    }

    events.sort((a, b) => b.timestamp - a.timestamp);

    if (events.length > 0) {
        let currentSet = false;

        for (const event of events) {
            if (event.id === 'reviewing-new') {
                event.status = 'current';
                currentSet = true;
                continue;
            }

            if (!currentSet) {
                const decisionType = event.type === 'decision' ? (doc.approval_decisions.find(d => d.ID === event.id)?.decision?.toLowerCase() || '') : '';

                if (event.id === 'reviewing') {
                    event.status = 'current';
                    currentSet = true;
                } else {
                    const isMatch = (currentStatus === decisionType) ||
                        (currentStatus === 'approved' && decisionType === 'approve') ||
                        (currentStatus === 'rejected' && decisionType === 'reject');

                    if (isMatch) {
                        if (decisionType === 'approve' || decisionType === 'approved') {
                            event.status = 'past-approved';
                        } else if (decisionType === 'reject' || decisionType === 'rejected') {
                            event.status = 'past-rejected';
                        } else {
                            if (isResubmitted.value) {
                                event.status = 'past-request-change';
                            } else {
                                event.status = 'current';
                            }
                        }
                        currentSet = true;
                    } else {
                        event.status = event.status.replace('current', 'past');
                    }
                }
            } else {
                event.status = event.status.replace('current', 'past');
            }
        }
    }

    return events;
});

watch(() => props.isOpen, (newValue) => {
    if (newValue) {
        comment.value = '';
        actionType.value = null;
        isSubmitting.value = false;
        submissionError.value = null;
        fetchAdminProfile();
    }
});

const closeModal = () => {
    if (isSubmitting.value) return;
    emit('close');
};

const submitAction = async (type: 'approve' | 'reject' | 'request-change') => {
    if (!props.documentData || !type) return;

    if (!adminId.value) {
        if (!submissionError.value) {
            submissionError.value = 'ไม่สามารถยืนยันตัวตนผู้ดำเนินการได้ กรุณาลองใหม่อีกครั้ง';
        }
        return;
    }

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
            comment: comment.value,
            admin_id: adminId.value,
        });
        emit('action-completed');
        closeModal();
    } catch (error) {
        submissionError.value = 'เกิดข้อผิดพลาดในการบันทึกข้อมูล กรุณาลองใหม่อีกครั้ง';
    } finally {
        isSubmitting.value = false;
    }
};

const openDocument = () => {
    const filePath = latestDocument.value?.file_path;
    if (filePath) {
        const backendBaseUrl = 'http://localhost:8080';
        const fileUrl = `${backendBaseUrl}/${filePath}`;
        window.open(fileUrl, '_blank');
    } else {
        alert('ไม่พบไฟล์เอกสาร');
    }
};

const openSpecificDocument = (path: string) => {
    if (path) {
        const backendBaseUrl = 'http://localhost:8080';
        const fileUrl = `${backendBaseUrl}/${path}`;
        window.open(fileUrl, '_blank');
    }
};
</script>

<template>
    <div v-if="isOpen && documentData"
        class="fixed inset-0 z-[200] flex items-center justify-center bg-black/50 backdrop-blur-sm p-4 transition-opacity"
        data-testid="approval-detail-modal">
        <div
            class="bg-white w-full max-w-6xl h-[90vh] rounded-2xl shadow-2xl flex flex-col overflow-hidden animate-pop-in">
            <div class="px-6 py-4 border-b flex items-center justify-between bg-slate-50">
                <div>
                    <h2 class="text-xl font-bold text-[#1e3a8a] flex items-center gap-2" data-testid="scholarship-name">
                        {{ documentData.application_document.application_scholarship?.scholarship?.scholarship_name ||
                        'รายละเอียดทุนการศึกษา' }}

                        <span class="badge badge-info text-white animate-pulse" v-if="isResubmitted"
                            data-testid="scholarship-status-badge">
                            มีการส่งแก้ไขใหม่ (รอตรวจสอบ)
                        </span>
                        <span class="badge badge-success text-white"
                            v-else-if="documentData.status?.toLowerCase() === 'approved'"
                            data-testid="scholarship-status-badge">
                            อนุมัติแล้ว
                        </span>
                        <span class="badge badge-error text-white"
                            v-else-if="documentData.status?.toLowerCase() === 'rejected'"
                            data-testid="scholarship-status-badge">
                            ปฏิเสธ
                        </span>
                        <span class="badge badge-warning text-white"
                            v-else-if="documentData.status?.toLowerCase() === 'request-change'"
                            data-testid="scholarship-status-badge">
                            รอผู้สมัครแก้ไข
                        </span>
                        <span class="badge badge-warning text-white" v-else data-testid="scholarship-status-badge">
                            รอตรวจสอบ
                        </span>
                    </h2>
                    <div class="flex items-center gap-2 text-sm text-gray-500 mt-1">
                        <span data-testid="task-id">Task ID: #{{ documentData.ID }}</span>
                        <span>•</span>
                        <span class="bg-blue-50 text-blue-700 px-2 rounded border border-blue-100 text-xs"
                            data-testid="semester-text">
                            {{ semesterText }}
                        </span>
                    </div>
                </div>
                <button @click="closeModal" class="btn btn-circle btn-ghost btn-sm text-gray-500 hover:bg-gray-200"
                    data-testid="modal-close-button">
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
                        <div class="card bg-white shadow-sm border border-gray-100" data-testid="applicant-info-card">
                            <div class="card-body p-5">
                                <h3 class="font-bold text-lg text-slate-700 mb-4 border-b pb-2">ข้อมูลผู้สมัคร</h3>
                                <div v-if="documentData.application_document?.application_scholarship?.application?.student_profile"
                                    class="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm">
                                    <div>
                                        <span class="block text-gray-400">ชื่อ-นามสกุล</span>
                                        <span class="font-semibold text-slate-800 text-lg"
                                            data-testid="student-name-value">
                                            {{
                                                documentData.application_document.application_scholarship.application.student_profile.first_name_th
                                            }}
                                            {{
                                                documentData.application_document.application_scholarship.application.student_profile.last_name_th
                                            }}
                                        </span>
                                    </div>
                                    <div>
                                        <span class="block text-gray-400">รหัสนักศึกษา</span>
                                        <span class="font-semibold text-slate-800" data-testid="student-id-value">
                                            {{
                                                documentData.application_document.application_scholarship.application.student_profile.student_id
                                            || '-' }}
                                        </span>
                                    </div>
                                    <div>
                                        <span class="block text-gray-400">สาขาวิชา</span>
                                        <span class="font-semibold text-slate-800" data-testid="student-major-value">
                                            {{
                                                documentData.application_document.application_scholarship.application.student_profile.major?.major_name
                                            || 'N/A' }}
                                        </span>
                                    </div>
                                    <div>
                                        <span class="block text-gray-400">เกรดเฉลี่ย (GPAX)</span>
                                        <span class="font-semibold text-slate-800" data-testid="student-gpax-value">
                                            {{
                                                documentData.application_document.application_scholarship.application.student_profile.gpax
                                                    ?
                                                    documentData.application_document.application_scholarship.application.student_profile.gpax.toFixed(2)
                                                    : '-' }}
                                        </span>
                                    </div>
                                </div>
                                <div v-else class="text-center text-gray-400 py-4" data-testid="no-applicant-info">
                                    ไม่พบข้อมูลผู้สมัคร</div>
                            </div>
                        </div>

                        <div class="card bg-white shadow-sm border border-gray-100" data-testid="attachment-info-card">
                            <div class="card-body p-5">
                                <h3 class="font-bold text-lg text-slate-700 mb-4 border-b pb-2">
                                    เอกสารแนบ (ฉบับล่าสุด)
                                    <span v-if="isResubmitted"
                                        class="ml-2 badge badge-info text-white text-xs">NEW</span>
                                </h3>
                                <div class="space-y-3">
                                    <div v-if="latestDocument"
                                        class="flex items-center justify-between p-3 bg-slate-50 rounded-lg border border-gray-200">
                                        <div class="flex items-center gap-3">
                                            <div
                                                class="w-10 h-10 bg-red-100 rounded-lg flex items-center justify-center text-red-500 font-bold text-xs">
                                                PDF
                                            </div>
                                            <div>
                                                <p class="font-medium text-slate-700 truncate max-w-[300px]"
                                                    data-testid="file-name">
                                                    {{ latestDocument.file_name }}
                                                </p>
                                                <p class="text-xs text-gray-400" data-testid="file-upload-date">
                                                    อัปโหลดเมื่อ {{ formatDate(latestDocument.CreatedAt) }}
                                                </p>
                                            </div>
                                        </div>
                                        <button @click="openDocument"
                                            class="btn btn-sm btn-ghost text-[#1e3a8a] hover:bg-blue-50"
                                            data-testid="view-file-button">ดูไฟล์</button>
                                    </div>
                                    <div v-else class="text-center text-gray-400 py-4" data-testid="no-attachment">
                                        ไม่พบเอกสารแนบ</div>
                                </div>
                            </div>
                        </div>

                        <div class="card bg-white shadow-sm border border-gray-100" data-testid="requirements-card">
                            <div class="card-body p-5">
                                <h3 class="font-bold text-lg text-slate-700 mb-4 border-b pb-2">ข้อกำหนดทั้งหมดของทุนนี้
                                </h3>
                                <div
                                    v-if="documentData.application_document?.application_scholarship?.scholarship?.approval_requirements?.length">
                                    <ul class="space-y-3 list-disc list-inside text-slate-600 text-sm"
                                        data-testid="requirement-list">
                                        <li v-for="req in documentData.application_document.application_scholarship.scholarship.approval_requirements"
                                            :key="req.ID" :data-testid="`requirement-item-${req.ID}`">
                                            {{ req.requirement.name }}
                                        </li>
                                    </ul>
                                </div>
                                <div v-else class="text-center text-gray-400 py-4" data-testid="no-requirements">
                                    ไม่พบข้อมูลข้อกำหนด
                                </div>
                            </div>
                        </div>
                    </div>

                    <div class="lg:col-span-1 space-y-6">

                        <div v-if="canAction"
                            class="card bg-white shadow-md border border-blue-100 ring-4 ring-blue-50/50"
                            data-testid="decision-card">
                            <div class="card-body p-5">
                                <h3 class="font-bold text-lg text-slate-700 mb-4">ผลการพิจารณา</h3>
                                <div class="flex flex-col gap-2">
                                    <button @click="submitAction('approve')" :disabled="isSubmitting"
                                        class="btn bg-[#1e3a8a] hover:bg-[#152c6f] text-white w-full border-none"
                                        data-testid="approve-button">
                                        <span v-if="isSubmitting" class="loading loading-spinner loading-sm"></span>
                                        อนุมัติเอกสาร (Approve)
                                    </button>
                                    <div class="grid grid-cols-2 gap-2 mt-2">
                                        <button @click="actionType = 'request-change'" :disabled="isSubmitting"
                                            class="btn btn-outline btn-warning btn-sm hover:text-white"
                                            data-testid="request-change-button">
                                            ขอแก้ไข
                                        </button>
                                        <button @click="actionType = 'reject'" :disabled="isSubmitting"
                                            class="btn btn-outline btn-error btn-sm hover:text-white"
                                            data-testid="reject-button">
                                            ปฏิเสธ
                                        </button>
                                    </div>
                                </div>

                                <div v-if="actionType" class="mt-4 pt-4 border-t animate-fade-in"
                                    data-testid="comment-section">
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
                                            class="btn btn-ghost btn-xs text-gray-500"
                                            data-testid="cancel-comment-button">ยกเลิก</button>
                                        <button @click="submitAction(actionType!)"
                                            class="btn btn-sm text-white border-none"
                                            :class="actionType === 'reject' ? 'btn-error' : 'btn-warning'"
                                            :disabled="isSubmitting" data-testid="confirm-decision-button">
                                            <span v-if="isSubmitting"
                                                class="loading loading-spinner loading-xs"></span>
                                            ยืนยัน
                                        </button>
                                    </div>
                                </div>

                                <div v-if="submissionError"
                                        class="text-error text-xs mt-2 font-medium bg-red-50 p-2 rounded"
                                        data-testid="submission-error-message">
                                        {{ submissionError }}
                                </div>
                            </div>
                        </div>

                        <div v-else class="card bg-white shadow-sm border border-gray-100"
                            data-testid="decision-card-inactive">
                            <div class="card-body p-5 items-center text-center">
                                <div class="badge badge-lg p-4 font-bold text-white mb-2"
                                    :class="documentData.status?.toLowerCase() === 'approved' ? 'badge-success' : 'badge-error'"
                                    data-testid="inactive-decision-badge">
                                    {{ documentData.status?.toLowerCase() === 'approved' ? 'อนุมัติเรียบร้อยแล้ว' :
                                    'สิ้นสุดการดำเนินการ' }}
                                </div>
                                <p class="text-xs text-gray-400" data-testid="inactive-decision-status">สถานะปัจจุบัน: {{
                                    documentData.status }}</p>
                            </div>
                        </div>

                        <div class="card bg-white shadow-sm border border-gray-100" data-testid="timeline-card">
                            <div class="card-body p-5">
                                <h3 class="font-bold text-lg text-slate-700 mb-4">ประวัติการดำเนินการ</h3>
                                <ul class="timeline timeline-vertical timeline-compact -ml-4">
                                    <li v-for="(event, index) in processedTimelineEvents" :key="event.id"
                                        :data-testid="`timeline-event-${index}`">
                                        <hr v-if="index > 0" class="bg-gray-200" />
                                        <div class="timeline-middle">
                                            <div v-if="event.status === 'current'"
                                                class="relative flex items-center justify-center w-6 h-6">
                                                <span
                                                    class="absolute inline-flex h-full w-full rounded-full bg-blue-400 opacity-75 animate-ping"></span>
                                                <span
                                                    class="relative inline-flex rounded-full h-6 w-6 bg-[#1e3a8a] border-4 border-blue-100"></span>
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

                                        <div
                                            class="timeline-end timeline-box w-full border-none shadow-none p-0 pl-2 mb-6">
                                            <div class="font-bold text-slate-800 text-sm"
                                                :class="{ 'text-[#1e3a8a]': event.status === 'current' }"
                                                :data-testid="`timeline-event-title-${index}`">{{ event.title
                                                }}</div>
                                            <div class="text-xs text-gray-500 mb-1"
                                                :data-testid="`timeline-event-meta-${index}`">{{ event.date }} • โดย {{
                                                event.actor }}</div>

                                            <div v-if="(event.status === 'past-rejected' || event.status === 'past-request-change') && event.description"
                                                class="mt-2 p-2 rounded-lg text-xs border"
                                                :class="event.status === 'past-rejected' ? 'bg-red-50 text-red-700 border-red-200' : 'bg-orange-50 text-orange-700 border-orange-200'"
                                                :data-testid="`timeline-event-description-${index}`">
                                                {{ event.description }}
                                            </div>

                                            <div v-else-if="event.filePath"
                                                class="bg-white border border-gray-200 rounded-lg p-2.5 mt-2 shadow-sm max-w-sm">
                                                <div class="flex items-center justify-between gap-3">
                                                    <div class="flex items-center gap-2 overflow-hidden">
                                                        <div
                                                            class="w-8 h-8 bg-gray-100 rounded flex items-center justify-center text-gray-500 shrink-0">
                                                            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4"
                                                                fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                                                <path stroke-linecap="round" stroke-linejoin="round"
                                                                    stroke-width="2"
                                                                    d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                                                            </svg>
                                                        </div>
                                                        <div class="truncate text-xs font-medium text-slate-700"
                                                            :title="event.fileName"
                                                            :data-testid="`timeline-event-filename-${index}`">{{
                                                            event.fileName }}</div>
                                                    </div>
                                                    <button @click.prevent="openSpecificDocument(event.filePath)"
                                                        class="btn btn-xs btn-outline bg-white text-gray-600 hover:bg-gray-50 shrink-0 font-normal"
                                                        :data-testid="`timeline-event-view-file-button-${index}`">
                                                        {{ event.isOldFile ? 'ดูไฟล์เก่า' : 'ดูไฟล์' }}
                                                    </button>
                                                </div>
                                            </div>

                                            <div v-else class="text-xs text-gray-600 break-words mt-1"
                                                :data-testid="`timeline-event-description-${index}`">{{
                                                event.description }}</div>
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