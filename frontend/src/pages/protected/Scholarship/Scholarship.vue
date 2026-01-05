<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { ScholarshipAPI } from '@/services/api/scholarship';
import { StatusscholarshipAPI } from '@/services/api/statusscholarship';
import { TypescholarshipAPI } from '@/services/api/typescholarship';

import type {
  ScholarshipResponse,
  ScholarshipStatusResponse,
  ScholarshipTypeResponse
} from '@/interfaces';

import ScholarshipModal from './ScholarshipModal.vue';
import FeatureScholarshipModal from './featurescholarshipmodal.vue';
import ScholarshipDetailModal from './scholarshipdetailmodal.vue';

// =====================
// State
// =====================
const scholarships = ref<ScholarshipResponse[]>([]);
const statuses = ref<ScholarshipStatusResponse[]>([]);
const types = ref<ScholarshipTypeResponse[]>([]);

const isLoading = ref(true);
const error = ref<string | null>(null);

const searchQuery = ref('');
const sortOrder = ref<'newest' | 'oldest'>('newest');
const filterStatus = ref('all');
const filterType = ref('all');

// เพิ่ม / แก้ไขทุน
const isModalOpen = ref(false);
const selectedScholarship = ref<ScholarshipResponse | null>(null);

// ดูรายละเอียดทุน
const isDetailModalOpen = ref(false);
const detailScholarship = ref<ScholarshipResponse | null>(null);

// เพิ่มคุณสมบัติทุน
const isFeatureModalOpen = ref(false);
const currentScholarshipId = ref<number | null>(null);

// =====================
// Fetch
// =====================
const fetchScholarships = async () => {
  isLoading.value = true;
  error.value = null;
  try {
    scholarships.value = await ScholarshipAPI.getAll();
  } catch (err) {
    error.value = 'โหลดข้อมูลทุนไม่สำเร็จ';
    console.error(err);
  } finally {
    isLoading.value = false;
  }
};

const fetchStatuses = async () => {
  try {
    statuses.value = await StatusscholarshipAPI.getAll();
  } catch (err) {
    console.error('โหลดสถานะไม่สำเร็จ', err);
  }
};

const fetchTypes = async () => {
  try {
    types.value = await TypescholarshipAPI.getAll();
  } catch (err) {
    console.error('โหลดประเภทไม่สำเร็จ', err);
  }
};

onMounted(() => {
  fetchScholarships();
  fetchStatuses();
  fetchTypes();
});

// =====================
// Filtering
// =====================
const filteredScholarships = computed(() => {
  let data = [...scholarships.value];

  if (filterStatus.value !== 'all') {
    data = data.filter(
      s => s.statusscholarship?.status_name === filterStatus.value
    );
  }

  if (filterType.value !== 'all') {
    data = data.filter(
      s => s.typescholarship?.type_name === filterType.value
    );
  }

  if (searchQuery.value.trim()) {
    const q = searchQuery.value.toLowerCase();
    data = data.filter(
      s =>
        s.scholarship_name.toLowerCase().includes(q) ||
        s.description.toLowerCase().includes(q)
    );
  }

  data.sort((a, b) =>
    sortOrder.value === 'newest'
      ? new Date(b.open_date).getTime() - new Date(a.open_date).getTime()
      : new Date(a.open_date).getTime() - new Date(b.open_date).getTime()
  );

  return data;
});

// =====================
// Modal Controls
// =====================
const openAddModal = () => {
  selectedScholarship.value = null;
  isModalOpen.value = true;
};

const openEditModal = (item: ScholarshipResponse) => {
  selectedScholarship.value = item;
  isModalOpen.value = true;
};

const openDetailModal = (item: ScholarshipResponse) => {
  detailScholarship.value = item;
  isDetailModalOpen.value = true;
};

const handleDelete = async (item: ScholarshipResponse) => {
  if (confirm(`ต้องการลบทุน "${item.scholarship_name}" ใช่หรือไม่?`)) {
    await ScholarshipAPI.delete(item.ID);
    fetchScholarships();
  }
};

const handleSaved = () => {
  isModalOpen.value = false;
  fetchScholarships();
};

// =====================
// Feature Scholarship
// =====================
const openAddFeatureModal = (scholarshipId: number | null) => {
  currentScholarshipId.value = scholarshipId;
  isFeatureModalOpen.value = true;
};

const handleFeatureSaved = () => {
  isFeatureModalOpen.value = false;
};
</script>

<template>
  <!-- PAGE ROOT -->
  <div class="h-screen bg-[#f0f2f5] p-6 text-slate-800 flex flex-col">

    <!-- ================= HEADER / FILTER ================= -->
    <div class="shrink-0 mb-6">
      <div class="flex flex-col xl:flex-row justify-between items-center gap-4">
        <div class="flex gap-2 overflow-x-auto pb-2">
          <input
            v-model="searchQuery"
            placeholder="ค้นหาทุน..."
            class="input input-sm input-bordered rounded-full pl-4 w-64"
          />

          <select v-model="sortOrder" class="select select-sm select-bordered rounded-full">
            <option value="newest">ใหม่ล่าสุด</option>
            <option value="oldest">เก่าสุด</option>
          </select>

          <select v-model="filterStatus" class="select select-sm select-bordered rounded-full">
            <option value="all">สถานะทั้งหมด</option>
            <option v-for="st in statuses" :key="st.ID" :value="st.status_name">
              {{ st.status_name }}
            </option>
          </select>

          <select v-model="filterType" class="select select-sm select-bordered rounded-full">
            <option value="all">ทุกประเภท</option>
            <option v-for="tp in types" :key="tp.ID" :value="tp.type_name">
              {{ tp.type_name }}
            </option>
          </select>
        </div>

        <div class="flex gap-2">
          <button
            @click="openAddModal"
            class="btn btn-sm rounded-full bg-blue-600 text-white px-5 shadow"
          >
            ➕ เพิ่มทุนใหม่
          </button>
          <button
            @click="openAddFeatureModal(null)"
            class="btn btn-sm rounded-full bg-green-500 text-white px-5 shadow"
          >
            ➕ เพิ่มคุณสมบัติทุน
          </button>
        </div>
      </div>
    </div>

    <!-- ================= LIST (SCROLLABLE) ================= -->
    <div class="flex-1 overflow-y-auto pr-2">
      <transition-group name="fade" tag="div" class="space-y-4 pb-10">

        <div
          v-for="item in filteredScholarships"
          :key="item.ID"
          @click="openDetailModal(item)"
          class="card bg-white shadow-sm border hover:border-blue-300 rounded-2xl p-4 cursor-pointer"
        >
          <div class="flex justify-between items-start">
            <div>
              <h3 class="font-bold text-blue-800 text-lg">
                {{ item.scholarship_name }}
              </h3>
              <p class="text-gray-500 text-sm mt-1 line-clamp-2">
                {{ item.description }}
              </p>
              <div class="mt-2 text-xs text-gray-400">
                ประเภท: {{ item.typescholarship?.type_name }} |
                สถานะ: {{ item.statusscholarship?.status_name }} |
                บริษัท: {{ item.sponsor?.company_name }}
              </div>
              <div class="mt-1 text-xs text-gray-400">
                เปิดรับ: {{ item.open_date }} – ปิดรับ: {{ item.close_date }}
              </div>
            </div>

            <div class="flex gap-2 flex-col">
              <button
                @click.stop="openEditModal(item)"
                class="btn btn-xs rounded-full bg-yellow-400 text-white"
              >
                แก้ไข
              </button>
              <button
                @click.stop="handleDelete(item)"
                class="btn btn-xs rounded-full bg-red-500 text-white"
              >
                ลบ
              </button>
            </div>
          </div>
        </div>

        <div
          v-if="filteredScholarships.length === 0"
          class="text-center text-gray-400 py-16"
        >
          ไม่พบข้อมูลทุน
        </div>

      </transition-group>
    </div>

    <!-- ================= MODALS ================= -->
    <ScholarshipModal
      v-if="isModalOpen"
      :isOpen="isModalOpen"
      :scholarship="selectedScholarship"
      @close="isModalOpen = false"
      @saved="handleSaved"
    />

    <FeatureScholarshipModal
      v-if="isFeatureModalOpen"
      :isOpen="isFeatureModalOpen"
      :scholarshipId="currentScholarshipId"
      @close="isFeatureModalOpen = false"
      @saved="handleFeatureSaved"
    />

    <ScholarshipDetailModal
      :isOpen="isDetailModalOpen"
      :scholarship="detailScholarship"
      @close="isDetailModalOpen = false"
    />

  </div>
</template>
