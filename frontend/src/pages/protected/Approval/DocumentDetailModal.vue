<script setup lang="ts">
import { ref } from 'vue';

const props = defineProps<{
    isOpen: boolean;
    documentData: any;
}>();

const emit = defineEmits(['close', 'approve', 'reject', 'request-change']);
const comment = ref('');
const actionType = ref<'approve' | 'reject' | 'request-change' | null>(null);
const timelineEvents = ref([
    {
        id: 3,
        title: 'รอการตรวจสอบ (ฉบับปัจจุบัน)',
        date: '05 ธ.ค. 2568 - 09:00',
        description: 'ผู้สมัครส่งเอกสารฉบับแก้ไขเข้ามาใหม่',
        actor: 'ผู้สมัคร',
        status: 'current',
        type: 'user',
    },
    {
        id: 2,
        title: 'ส่งคืนแก้ไข (Request Change)',
        date: '04 ธ.ค. 2568 - 14:30',
        description: 'รูปถ่ายนักศึกษาไม่ชัดเจน และขาดสำเนาหน้าสมุดบัญชี',
        actor: 'Admin A',
        status: 'past',
        type: 'admin_warning'
    },
    {
        id: 1,
        title: 'ยื่นคำขอครั้งแรก',
        date: '03 ธ.ค. 2568 - 08:00',
        description: 'ยื่นเอกสารผ่านระบบออนไลน์',
        actor: 'ผู้สมัคร',
        status: 'past',
        type: 'user',
        attachment: {
            name: 'รูปถ่าย_v1.jpg',
            url: '#'
        }
    }
]);

const closeModal = () => {
    emit('close');
    comment.value = '';
    actionType.value = null;
};

const submitAction = (type: 'approve' | 'reject' | 'request-change') => {
    emit(type, {
        id: props.documentData?.id,
        comment: comment.value
    });
    closeModal();
};
</script>

<template>
    <div v-if="isOpen"
        class="fixed inset-0 z-[200] flex items-center justify-center bg-black/50 backdrop-blur-sm p-4 transition-opacity">

        <div
            class="bg-white w-full max-w-6xl h-[90vh] rounded-2xl shadow-2xl flex flex-col overflow-hidden animate-pop-in">

            <div class="px-6 py-4 border-b flex items-center justify-between bg-slate-50">
                <div>
                    <h2 class="text-xl font-bold text-[#1e3a8a] flex items-center gap-2">
                        {{ documentData?.title || 'รายละเอียด' }}
                        <span class="badge badge-warning text-white"
                            v-if="documentData?.status === 'pending'">รอตรวจสอบ</span>
                        <span class="badge badge-success text-white"
                            v-else-if="documentData?.status === 'approved'">อนุมัติแล้ว</span>
                    </h2>
                    <div class="flex items-center gap-2 text-sm text-gray-500 mt-1">
                        <span>รหัส: #REQ-{{ documentData?.id }}</span>
                        <span>•</span>
                        <span class="bg-blue-50 text-blue-700 px-2 rounded border border-blue-100 text-xs">รอบ: {{
                            documentData?.round }}</span>
                    </div>
                </div>
                <button @click="closeModal" class="btn btn-circle btn-ghost btn-sm">
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
                                <div class="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm">
                                    <div>
                                        <span class="block text-gray-400">ชื่อ-นามสกุล</span>
                                        <span class="font-semibold text-slate-800 text-lg">{{ documentData?.applicant ||
                                            '-' }}</span>
                                    </div>
                                    <div>
                                        <span class="block text-gray-400">รหัสนักศึกษา</span>
                                        <span class="font-semibold text-slate-800">B64XXXXX</span>
                                    </div>
                                    <div>
                                        <span class="block text-gray-400">สาขาวิชา</span>
                                        <span class="font-semibold text-slate-800">วิศวกรรมคอมพิวเตอร์</span>
                                    </div>
                                    <div>
                                        <span class="block text-gray-400">เกรดเฉลี่ย (GPAX)</span>
                                        <span class="font-semibold text-slate-800">3.50</span>
                                    </div>
                                </div>
                            </div>
                        </div>

                        <div class="card bg-white shadow-sm border border-gray-100">
                            <div class="card-body p-5">
                                <h3 class="font-bold text-lg text-slate-700 mb-4 border-b pb-2">เอกสารแนบ (ฉบับล่าสุด)
                                </h3>
                                <div class="space-y-3">
                                    <div
                                        class="flex items-center justify-between p-3 bg-slate-50 rounded-lg border border-gray-200">
                                        <div class="flex items-center gap-3">
                                            <div
                                                class="w-10 h-10 bg-red-100 rounded-lg flex items-center justify-center text-red-500 font-bold text-xs">
                                                PDF</div>
                                            <div>
                                                <p class="font-medium text-slate-700">หนังสือรับรองรายได้.pdf</p>
                                                <p class="text-xs text-gray-400">2.5 MB • อัปโหลดล่าสุด 05 ธ.ค. 68</p>
                                            </div>
                                        </div>
                                        <button class="btn btn-sm btn-ghost text-[#1e3a8a]">ดูไฟล์</button>
                                    </div>
                                    <div
                                        class="flex items-center justify-between p-3 bg-slate-50 rounded-lg border border-gray-200">
                                        <div class="flex items-center gap-3">
                                            <div
                                                class="w-10 h-10 bg-blue-100 rounded-lg flex items-center justify-center text-blue-500 font-bold text-xs">
                                                IMG</div>
                                            <div>
                                                <p class="font-medium text-slate-700">รูปถ่ายชุดนักศึกษา.jpg</p>
                                                <p class="text-xs text-gray-400">1.2 MB • อัปโหลดล่าสุด 05 ธ.ค. 68</p>
                                            </div>
                                        </div>
                                        <button class="btn btn-sm btn-ghost text-[#1e3a8a]">ดูไฟล์</button>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>

                    <div class="lg:col-span-1 space-y-6">

                        <div v-if="documentData?.status === 'pending'"
                            class="card bg-white shadow-md border border-blue-100">
                            <div class="card-body p-5">
                                <h3 class="font-bold text-lg text-slate-700 mb-4">ผลการพิจารณา</h3>

                                <div class="flex flex-col gap-2">
                                    <button @click="submitAction('approve')"
                                        class="btn bg-[#1e3a8a] hover:bg-[#152c6f] text-white w-full border-none">
                                        อนุมัติเอกสาร (Approve)
                                    </button>
                                    <div class="grid grid-cols-2 gap-2 mt-2">
                                        <button @click="actionType = 'request-change'"
                                            class="btn btn-outline btn-warning btn-sm">ขอแก้ไข</button>
                                        <button @click="actionType = 'reject'"
                                            class="btn btn-outline btn-error btn-sm">ปฏิเสธ</button>
                                    </div>
                                </div>

                                <div v-if="actionType" class="mt-4 pt-4 border-t animate-fade-in">
                                    <p class="text-sm font-bold mb-2"
                                        :class="actionType === 'reject' ? 'text-error' : 'text-warning'">
                                        {{ actionType === 'reject' ? 'ระบุเหตุผลการปฏิเสธ' : 'ระบุสิ่งที่ต้องแก้ไข' }}
                                    </p>
                                    <textarea v-model="comment" class="textarea textarea-bordered w-full h-24 text-sm"
                                        placeholder="พิมพ์รายละเอียด..."></textarea>
                                    <div class="flex justify-end gap-2 mt-2">
                                        <button @click="actionType = null" class="btn btn-ghost btn-xs">ยกเลิก</button>
                                        <button @click="submitAction(actionType!)"
                                            class="btn btn-sm btn-primary">ยืนยัน</button>
                                    </div>
                                </div>
                            </div>
                        </div>

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

                                            <div v-if="event.type === 'admin_warning'"
                                                class="bg-orange-50 border border-orange-200 text-orange-800 text-xs p-2 rounded mt-1">
                                                {{ event.description }}
                                            </div>
                                            <div v-else class="text-xs text-gray-600">{{ event.description }}</div>

                                            <div v-if="event.attachment"
                                                class="mt-2 p-2 bg-gray-100 rounded border border-gray-200 flex items-center justify-between gap-2 max-w-[90%]">
                                                <div class="flex items-center gap-2 overflow-hidden">
                                                    <svg xmlns="http://www.w3.org/2000/svg"
                                                        class="h-4 w-4 text-gray-400 shrink-0" fill="none"
                                                        viewBox="0 0 24 24" stroke="currentColor">
                                                        <path stroke-linecap="round" stroke-linejoin="round"
                                                            stroke-width="2"
                                                            d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                                                    </svg>
                                                    <span
                                                        class="text-xs text-gray-500 truncate line-through decoration-red-400">{{
                                                            event.attachment.name }}</span>
                                                </div>
                                                <a href="#"
                                                    class="text-[10px] bg-white border px-1.5 py-0.5 rounded text-gray-500 hover:text-blue-600 hover:border-blue-300">ดูไฟล์เก่า</a>
                                            </div>

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