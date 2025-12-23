<script setup lang="ts">
import { ref, computed } from "vue";
import { useRouter } from "vue-router";
import type { ApplicationScholarshipResponse } from "@/interfaces";
import Swal from "sweetalert2";

const props = defineProps<{
  isOpen: boolean;
  applicationData: ApplicationScholarshipResponse | null;
}>();

const emit = defineEmits(["close", "upload-file", "book-interview"]);

const router = useRouter();
const fileInput = ref<HTMLInputElement | HTMLInputElement[] | null>(null);
const isUploading = ref(false);

// Navigate to interview booking page
const goToInterviewBooking = () => {
  emit("close");
  router.push("/dashboard/schedule");
};

const stages = [
  {
    id: 1,
    title: "ยื่นใบสมัคร",
    description: "ส่งข้อมูลเข้าสู่ระบบเรียบร้อยแล้ว",
  },
  {
    id: 2,
    title: "ตรวจสอบคุณสมบัติ",
    description: "เจ้าหน้าที่ตรวจสอบคุณสมบัติเบื้องต้น",
  },
  { id: 3, title: "อัปโหลดเอกสาร", description: "ส่งเอกสารประกอบการพิจารณา" },
  {
    id: 4,
    title: "ตรวจสอบเอกสาร",
    description: "เจ้าหน้าที่ตรวจสอบความถูกต้องของเอกสาร",
  },
  { id: 5, title: "จองคิวสัมภาษณ์", description: "เลือกเวลาเพื่อเข้าสัมภาษณ์" },
  {
    id: 6,
    title: "ผลการสัมภาษณ์",
    description: "ประกาศผลการคัดเลือกและสัมภาษณ์",
  },
];

const parseGoDate = (dateString: string | undefined): Date | null => {
  if (!dateString) return null;

  const dateStr = String(dateString).split(" m=")[0];
  if (dateStr && dateStr.includes("+0000 UTC")) {
    const formattedStr = dateStr.replace(" +0000 UTC", "Z").replace(" ", "T");
    const date = new Date(formattedStr);
    return isNaN(date.getTime()) ? null : date;
  }
  if (!dateStr) return null;
  const date = new Date(dateStr);
  return isNaN(date.getTime()) ? null : date;
};

const processState = computed(() => {
  if (!props.applicationData)
    return { currentStep: 1, status: "normal", message: "" };

  const data = props.applicationData;
  const mainStatus = data.status?.toLowerCase();

  const screening = data.screening;

  // Check screening status first
  if (screening) {
    const screeningStatusId = screening.status_screening_id;
    const screeningStatusName = screening.status_screening?.name?.toLowerCase();

    // If screening rejected (ID: 3 or has rejection reason)
    if (
      screeningStatusId === 3 ||
      screening.rejection_reason ||
      screeningStatusName === "rejected" ||
      screeningStatusName === "not pass" ||
      screeningStatusName?.includes("ไม่ผ่าน")
    ) {
      return {
        currentStep: 2,
        status: "error",
        message:
          screening.rejection_reason || "คุณสมบัติไม่ผ่านเกณฑ์การคัดกรอง",
      };
    }

    // If screening is still pending (ID: 1 = รอตรวจสอบ) - stay at step 2
    if (screeningStatusId === 1) {
      return {
        currentStep: 2,
        status: "process",
        message: "รอเจ้าหน้าที่ตรวจสอบคุณสมบัติ",
      };
    }

    // If screening passed (ID: 2) - check if documents are uploaded
    if (screeningStatusId === 2) {
      const docs = data.application_documents;
      if (!docs || docs.length === 0) {
        // Screening passed but no documents uploaded - go to step 3
        return {
          currentStep: 3,
          status: "action",
          message: "ผ่านการคัดกรองแล้ว กรุณาอัปโหลดเอกสาร",
        };
      }
      // If documents exist, continue to check document approval status below
    }
  } else if (mainStatus === "new") {
    // No screening record yet but status is new - should wait for screening
    return {
      currentStep: 2,
      status: "process",
      message: "รอเจ้าหน้าที่ตรวจสอบคุณสมบัติ",
    };
  }

  const docs = data.application_documents;
  let latestDecisionType = "";
  let latestDecisionComment = "";
  let latestDocTime = 0;
  let latestDecisionTime = 0;

  if (docs && docs.length > 0) {
    const sortedDocs = [...docs].sort(
      (a, b) =>
        (parseGoDate(b.CreatedAt)?.getTime() || 0) -
        (parseGoDate(a.CreatedAt)?.getTime() || 0)
    );
    latestDocTime = parseGoDate(sortedDocs[0]?.CreatedAt)?.getTime() || 0;

    const allDecisions: any[] = [];
    docs.forEach((doc) => {
      if (doc.approval_tasks) {
        doc.approval_tasks.forEach((task) => {
          if (task.approval_decisions) {
            task.approval_decisions.forEach((decision) => {
              allDecisions.push(decision);
            });
          }
        });
      }
    });

    if (allDecisions.length > 0) {
      allDecisions.sort((a, b) => (Number(b.ID) || 0) - (Number(a.ID) || 0));
      const latest = allDecisions[0];
      latestDecisionType = latest.decision?.toLowerCase();
      latestDecisionComment = latest.comment;
      latestDecisionTime = parseGoDate(latest.decision_at)?.getTime() || 0;
    }
  }

  if (latestDecisionType === "rejected" || latestDecisionType === "reject") {
    return {
      currentStep: 4,
      status: "error",
      message: latestDecisionComment,
    };
  }

  if (latestDecisionType === "request-change") {
    if (latestDocTime > latestDecisionTime) {
      return {
        currentStep: 4,
        status: "process",
        message: "รอตรวจสอบเอกสารแก้ไข",
      };
    } else {
      return {
        currentStep: 3,
        status: "warning",
        message: latestDecisionComment,
      };
    }
  }

  // Check for 'qualified' status - this means all documents are approved
  // and student can proceed to interview booking
  if (mainStatus === "qualified") {
    return {
      currentStep: 5,
      status: "action",
      message: "เอกสารผ่านการอนุมัติแล้ว กรุณาจองคิวสัมภาษณ์",
    };
  }

  if (["approve", "approved"].includes(latestDecisionType)) {
    if (!["completed", "final-approved"].includes(mainStatus || "")) {
      return {
        currentStep: 5,
        status: "action",
        message: "เอกสารผ่านการอนุมัติแล้ว",
      };
    }
  }

  if (mainStatus === "new")
    return { currentStep: 2, status: "process", message: "" };

  // Handle 'pending' status - check if documents have been uploaded
  if (mainStatus === "pending") {
    if (docs && docs.length > 0)
      return { currentStep: 4, status: "process", message: "" };
    return { currentStep: 3, status: "action", message: "" };
  }

  if (["approved", "approve"].includes(mainStatus || "")) {
    return { currentStep: 5, status: "action", message: "" };
  }

  if (["completed", "final-approved"].includes(mainStatus || "")) {
    return { currentStep: 6, status: "completed", message: "" };
  }

  if (["rejected", "reject"].includes(mainStatus || "")) {
    if (!docs || docs.length === 0)
      return { currentStep: 2, status: "error", message: "ไม่ผ่านการพิจารณา" };
    return { currentStep: 4, status: "error", message: "ไม่ผ่านการพิจารณา" };
  }

  // Handle empty status (after screening passed) - check if documents exist
  if (mainStatus === "" || !mainStatus) {
    if (docs && docs.length > 0) {
      return { currentStep: 4, status: "process", message: "รอเจ้าหน้าที่ตรวจสอบเอกสาร" };
    }
    return { currentStep: 3, status: "action", message: "กรุณาอัปโหลดเอกสาร" };
  }

  return { currentStep: 1, status: "normal", message: "" };
});

const getStepStatus = (stepId: number) => {
  const { currentStep, status } = processState.value;

  if (stepId < currentStep) return "completed";

  if (stepId === currentStep) {
    if (status === "error") return "rejected";
    if (status === "warning") return "request-change";
    if (status === "action") return "action-needed";
    if (status === "completed") return "completed";
    return "current";
  }

  return "pending";
};

const latestComment = computed(() => {
  const docs =
    props.applicationData?.application_documents ||
    props.applicationData?.application_documents;
  if (!docs) return "";

  const allDecisions: any[] = [];

  docs.forEach((doc) => {
    if (doc.approval_tasks) {
      doc.approval_tasks.forEach((task) => {
        if (task.approval_decisions) {
          task.approval_decisions.forEach((decision) => {
            allDecisions.push(decision);
          });
        }
      });
    }
  });

  if (allDecisions.length === 0) return "";
  allDecisions.sort((a, b) => (Number(b.ID) || 0) - (Number(a.ID) || 0));
  return allDecisions[0]?.comment || "";
});

const triggerUpload = () => {
  if (Array.isArray(fileInput.value)) {
    fileInput.value[0]?.click();
  } else {
    fileInput.value?.click();
  }
};

const handleFileChange = async (e: Event) => {
  const target = e.target as HTMLInputElement;
  const file = target.files?.[0];

  if (file) {
    const result = await Swal.fire({
      title: "ยืนยันการส่งเอกสาร",
      text: `คุณต้องการส่งไฟล์ "${file.name}" ใช่หรือไม่?`,
      icon: "question",
      showCancelButton: true,
      confirmButtonText: "ยืนยันส่งข้อมูล",
      cancelButtonText: "ยกเลิก",
      confirmButtonColor: "#1e3a8a",
      cancelButtonColor: "#d33",
      reverseButtons: true,
      width: "400px",
      customClass: {
        title: "text-lg font-bold",
        htmlContainer: "text-sm",
        popup: "rounded-xl",
      },
    });

    if (result.isConfirmed) {
      isUploading.value = true;

      setTimeout(() => {
        emit("upload-file", file);
        isUploading.value = false;

        Swal.fire({
          title: "อัปโหลดสำเร็จ!",
          text: "ระบบได้รับเอกสารของคุณเรียบร้อยแล้ว",
          icon: "success",
          confirmButtonText: "ตกลง",
          confirmButtonColor: "#1e3a8a",
          timer: 2000,
          timerProgressBar: true,
          width: "400px",
          customClass: {
            title: "text-lg font-bold",
            htmlContainer: "text-sm",
            popup: "rounded-xl",
          },
        });

        target.value = "";
      }, 1000);
    } else {
      target.value = "";
    }
  }
};
</script>

<template>
  <div
    v-if="isOpen && applicationData"
    class="fixed inset-0 z-[200] flex items-center justify-center bg-black/60 backdrop-blur-sm p-4 animate-fade"
  >
    <div
      class="bg-white w-full max-w-3xl max-h-[90vh] rounded-xl shadow-2xl flex flex-col overflow-hidden animate-pop-in"
    >
      <div
        class="bg-[#1e3a8a] px-6 py-4 shrink-0 flex justify-between items-start"
      >
        <div>
          <h2 class="text-xl font-bold text-white flex items-center gap-2">
            ติดตามสถานะ
            <span
              class="bg-white/20 text-xs px-2 py-0.5 rounded text-blue-50 font-normal border border-white/10"
            >
              #APP-{{ applicationData.ID }}
            </span>
          </h2>
          <p class="text-blue-200 text-sm mt-1">
            {{
              applicationData.scholarship?.scholarship_name || "รายละเอียดทุน"
            }}
          </p>
        </div>
        <button
          @click="$emit('close')"
          class="btn btn-circle btn-ghost btn-sm text-white/70 hover:bg-white/20 hover:text-white"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            class="h-6 w-6"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M6 18L18 6M6 6l12 12"
            />
          </svg>
        </button>
      </div>

      <div class="flex-1 overflow-y-auto bg-white p-6 sm:p-8">
        <h3 class="font-bold text-slate-700 mb-6 text-lg border-b pb-2">
          ขั้นตอนการดำเนินการ
        </h3>

        <div class="pl-2">
          <div
            v-for="(stage, index) in stages"
            :key="stage.id"
            class="relative pb-10 last:pb-0"
          >
            <div
              v-if="index !== stages.length - 1"
              class="absolute left-[19px] top-10 bottom-0 w-1 transition-colors duration-500 z-0"
              :class="
                getStepStatus(stage.id) === 'completed'
                  ? 'bg-green-500'
                  : 'bg-gray-300'
              "
            ></div>

            <div class="relative flex items-start gap-4 z-10">
              <div
                class="relative flex-shrink-0 w-10 h-10 rounded-full flex items-center justify-center border-2 font-bold transition-all bg-white shadow-sm"
                :class="{
                  'border-green-500 text-green-600':
                    getStepStatus(stage.id) === 'completed',
                  'border-[#1e3a8a] text-[#1e3a8a] ring-4 ring-blue-50':
                    getStepStatus(stage.id) === 'current' ||
                    getStepStatus(stage.id) === 'action-needed',
                  'border-orange-400 text-orange-500 ring-4 ring-orange-50':
                    getStepStatus(stage.id) === 'request-change',
                  'border-red-500 text-red-500':
                    getStepStatus(stage.id) === 'rejected',
                  'border-gray-300 text-gray-400':
                    getStepStatus(stage.id) === 'pending',
                }"
              >
                <svg
                  v-if="getStepStatus(stage.id) === 'completed'"
                  xmlns="http://www.w3.org/2000/svg"
                  class="h-6 w-6"
                  viewBox="0 0 20 20"
                  fill="currentColor"
                >
                  <path
                    fill-rule="evenodd"
                    d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                    clip-rule="evenodd"
                  />
                </svg>
                <span v-else>{{ stage.id }}</span>
              </div>

              <div class="flex-1 pt-1">
                <div class="flex flex-wrap justify-between items-start gap-2">
                  <div>
                    <h4
                      class="font-bold text-base"
                      :class="{
                        'text-green-700':
                          getStepStatus(stage.id) === 'completed',
                        'text-[#1e3a8a]': ['current', 'action-needed'].includes(
                          getStepStatus(stage.id)
                        ),
                        'text-orange-600':
                          getStepStatus(stage.id) === 'request-change',
                        'text-red-600': getStepStatus(stage.id) === 'rejected',
                        'text-gray-400': getStepStatus(stage.id) === 'pending',
                      }"
                    >
                      {{ stage.title }}
                    </h4>
                    <p class="text-sm text-gray-500 mt-0.5">
                      {{ stage.description }}
                    </p>
                  </div>

                  <div class="shrink-0">
                    <span
                      v-if="getStepStatus(stage.id) === 'completed'"
                      class="badge badge-success text-white badge-sm"
                      >เสร็จสิ้น</span
                    >
                    <span
                      v-else-if="
                        ['current', 'process'].includes(getStepStatus(stage.id))
                      "
                      class="badge badge-info text-white badge-sm animate-pulse"
                      >กำลังดำเนินการ</span
                    >
                    <span
                      v-else-if="getStepStatus(stage.id) === 'action-needed'"
                      class="badge badge-primary badge-outline badge-sm"
                      >รอคุณดำเนินการ</span
                    >
                    <span
                      v-else-if="getStepStatus(stage.id) === 'request-change'"
                      class="badge badge-warning text-white badge-sm"
                      >ต้องแก้ไข</span
                    >
                    <span
                      v-else-if="getStepStatus(stage.id) === 'rejected'"
                      class="badge badge-error text-white badge-sm"
                      >ไม่ผ่าน</span
                    >
                  </div>
                </div>

                <div class="mt-3 animate-fade-in">
                  <input
                    type="file"
                    ref="fileInput"
                    class="hidden"
                    @change="handleFileChange"
                    accept=".pdf,.jpg,.png"
                  />

                  <div v-if="stage.id === 3">
                    <div
                      v-if="getStepStatus(3) === 'action-needed'"
                      class="bg-blue-50 border border-blue-100 rounded-lg p-4 mt-2"
                    >
                      <p
                        class="text-sm text-blue-800 mb-3 font-medium flex items-center gap-2"
                      >
                        <svg
                          xmlns="http://www.w3.org/2000/svg"
                          class="h-4 w-4"
                          fill="none"
                          viewBox="0 0 24 24"
                          stroke="currentColor"
                        >
                          <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                          />
                        </svg>
                        คุณสมบัติผ่านเกณฑ์แล้ว
                      </p>
                      <button
                        @click="triggerUpload"
                        :disabled="isUploading"
                        class="btn bg-[#1e3a8a] hover:bg-blue-800 text-white btn-sm border-none w-full sm:w-auto"
                      >
                        <span
                          v-if="isUploading"
                          class="loading loading-spinner loading-xs"
                        ></span>
                        <span v-else class="flex items-center gap-2">
                          <svg
                            xmlns="http://www.w3.org/2000/svg"
                            class="h-4 w-4"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke="currentColor"
                          >
                            <path
                              stroke-linecap="round"
                              stroke-linejoin="round"
                              stroke-width="2"
                              d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12"
                            />
                          </svg>
                          อัปโหลดเอกสาร
                        </span>
                      </button>
                    </div>

                    <div
                      v-if="getStepStatus(3) === 'request-change'"
                      class="bg-orange-50 border border-orange-200 rounded-lg p-4 mt-3 animate-fade-in"
                    >
                      <div class="flex items-start gap-3 mb-3">
                        <div class="bg-orange-100 p-2 rounded-full shrink-0">
                          <svg
                            xmlns="http://www.w3.org/2000/svg"
                            class="h-5 w-5 text-orange-600"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke="currentColor"
                          >
                            <path
                              stroke-linecap="round"
                              stroke-linejoin="round"
                              stroke-width="2"
                              d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
                            />
                          </svg>
                        </div>
                        <div class="w-full">
                          <p class="text-sm font-bold text-orange-900">
                            เอกสารต้องได้รับการแก้ไข
                          </p>
                          <div
                            class="mt-2 bg-white p-3 rounded border border-orange-100 text-sm text-slate-700 shadow-sm"
                          >
                            <span
                              class="font-semibold text-orange-800 text-xs block mb-1"
                              >สิ่งที่ต้องแก้ไข:</span
                            >
                            "{{
                              processState.message ||
                              "กรุณาตรวจสอบความถูกต้องของเอกสารแล้วอัปโหลดใหม่อีกครั้ง"
                            }}"
                          </div>
                        </div>
                      </div>
                      <button
                        @click="triggerUpload"
                        :disabled="isUploading"
                        class="btn bg-orange-500 hover:bg-orange-600 text-white btn-sm border-none w-full sm:w-auto shadow-md"
                      >
                        <span
                          v-if="isUploading"
                          class="loading loading-spinner loading-xs"
                        ></span>
                        <span v-else class="flex items-center gap-2">
                          <svg
                            xmlns="http://www.w3.org/2000/svg"
                            class="h-4 w-4"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke="currentColor"
                          >
                            <path
                              stroke-linecap="round"
                              stroke-linejoin="round"
                              stroke-width="2"
                              d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12"
                            />
                          </svg>
                          อัปโหลดเอกสารใหม่
                        </span>
                      </button>
                    </div>
                  </div>

                  <div
                    v-if="
                      stage.id === 5 &&
                      ['action-needed', 'current'].includes(getStepStatus(5))
                    "
                    class="bg-green-50 border border-green-100 rounded-lg p-4 mt-2"
                  >
                    <div class="flex items-center gap-3 mb-3">
                      <div
                        class="w-10 h-10 bg-green-100 rounded-full flex items-center justify-center text-green-600 shrink-0"
                      >
                        <svg
                          xmlns="http://www.w3.org/2000/svg"
                          class="h-6 w-6"
                          fill="none"
                          viewBox="0 0 24 24"
                          stroke="currentColor"
                        >
                          <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                          />
                        </svg>
                      </div>
                      <div>
                        <p class="font-bold text-green-800">
                          ยินดีด้วย! เอกสารผ่านการอนุมัติ
                        </p>
                        <p class="text-xs text-green-700">
                          กรุณาจองคิวสัมภาษณ์เพื่อดำเนินการต่อ
                        </p>
                      </div>
                    </div>

                    <button
                      @click="goToInterviewBooking"
                      class="btn bg-[#1e3a8a] hover:bg-blue-900 text-white btn-sm border-none shadow-lg shadow-blue-200 w-full"
                    >
                      จองคิวสัมภาษณ์ทันที
                    </button>
                  </div>

                  <div
                    v-if="getStepStatus(stage.id) === 'rejected'"
                    class="bg-red-50 border border-red-100 rounded-lg p-3 text-sm text-red-700 mt-2"
                  >
                    <span class="font-bold">เหตุผลที่ไม่ผ่าน:</span>
                    {{ latestComment || "ไม่ผ่านการพิจารณาในขั้นตอนนี้" }}
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="p-4 bg-gray-50 border-t text-center shrink-0">
        <button
          @click="$emit('close')"
          class="btn btn-ghost w-full max-w-xs text-gray-500 font-normal"
        >
          ปิดหน้าต่าง
        </button>
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

@keyframes fade {
  from {
    opacity: 0;
  }

  to {
    opacity: 1;
  }
}

.animate-fade {
  animation: fade 0.2s ease-out forwards;
}

@keyframes fade-in {
  from {
    opacity: 0;
    height: 0;
    transform: translateY(-5px);
  }

  to {
    opacity: 1;
    height: auto;
    transform: translateY(0);
  }
}

.animate-fade-in {
  animation: fade-in 0.3s ease forwards;
}
</style>