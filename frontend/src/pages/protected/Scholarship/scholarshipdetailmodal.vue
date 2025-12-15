<script setup lang="ts">
import { ref, watch } from 'vue';
import type {
  ScholarshipResponse,
  FeatureScholarshipResponse,
} from '@/interfaces';
import { FeatureScholarshipAPI } from '@/services/api/featurescholarship';

interface Props {
  isOpen: boolean;
  scholarship: ScholarshipResponse | null;
}

const props = defineProps<Props>();
const emit = defineEmits<{ (e: 'close'): void }>();

const features = ref<FeatureScholarshipResponse[]>([]);
const isLoading = ref(false);

// โหลดคุณสมบัติของทุน (ใช้ logic เดิมที่ใช้งานได้)
const loadFeatures = async () => {
  if (!props.scholarship) return;

  isLoading.value = true;
  try {
    const all = await FeatureScholarshipAPI.getAll();
    features.value = (all as FeatureScholarshipResponse[]).filter(
      f => f.scholarship_id === props.scholarship!.ID
    );
  } catch (err) {
    console.error(err);
  } finally {
    isLoading.value = false;
  }
};

// ✅ watch เหมือนไฟล์ต้นฉบับ
watch(
  () => props.isOpen,
  (open) => {
    if (open) loadFeatures();
  }
);
</script>

<template>
  <div v-if="isOpen && scholarship" class="modal modal-open">
    <div class="modal-box max-w-2xl">
      <h2 class="font-bold text-xl text-blue-800 mb-2">
        {{ scholarship.scholarship_name }}
      </h2>

      <p class="text-gray-600 mb-3 whitespace-pre-line">
        {{ scholarship.description }}
      </p>

      <!-- ข้อมูลทั่วไป -->
      <div class="grid grid-cols-2 gap-2 text-sm text-gray-700 mb-4">
        <div>
          <b>ประเภท:</b> {{ scholarship.typescholarship?.type_name }}
        </div>
        <div>
          <b>สถานะ:</b> {{ scholarship.statusscholarship?.status_name }}
        </div>
        <div>
          <b>บริษัท:</b> {{ scholarship.sponsor?.company_name }}
        </div>

        <!-- ✅ แสดง semaster แทน semester -->
        <div>
          <b>ภาคการศึกษา:</b>
          <span v-if="scholarship.semaster">
            ปี {{ scholarship.semaster.academic_year }}
            / เทอม {{ scholarship.semaster.term }}
            / รอบ {{ scholarship.semaster.round }}
          </span>
          <span v-else class="text-gray-400">ไม่ระบุ</span>
        </div>

        <div class="col-span-2">
          <b>ระยะเวลา:</b>
          {{ scholarship.open_date }} – {{ scholarship.close_date }}
        </div>
      </div>

      <!-- คุณสมบัติทุน -->
      <div class="border-t pt-3">
        <h3 class="font-bold mb-2">คุณสมบัติทุน</h3>

        <div v-if="isLoading" class="text-gray-400 text-sm">
          กำลังโหลดข้อมูล...
        </div>

        <ul v-else>
          <li
            v-for="f in features"
            :key="f.ID"
            class="border-b py-1 text-sm"
          >
            {{ f.typefeature?.type_feature_name }}
            {{ f.operator }}
            {{ f.value }}
          </li>

          <li v-if="features.length === 0" class="text-gray-400 text-sm py-2">
            ไม่มีคุณสมบัติกำหนดไว้
          </li>
        </ul>
      </div>

      <div class="modal-action">
        <button class="btn" @click="emit('close')">ปิด</button>
      </div>
    </div>
  </div>
</template>
