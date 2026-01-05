<script setup lang="ts">
import { ref, watch } from 'vue';
import { TypefeatureAPI } from '@/services/api/typefeature';
import type { TypeFeatureResponse, TypeFeatureCreate } from '@/interfaces';

/* ---------------- props / emits ---------------- */
const props = defineProps<{ isOpen: boolean }>();
const emit = defineEmits<{ (e: 'close'): void; (e: 'saved'): void }>();

/* ---------------- state ---------------- */
const typeFeatures = ref<TypeFeatureResponse[]>([]);
const form = ref<{ type_feature_name: string }>({
  type_feature_name: ''
});
const editingFeature = ref<TypeFeatureResponse | null>(null);

/* ---------------- fetch ---------------- */
const fetchTypeFeatures = async () => {
  try {
    typeFeatures.value = await TypefeatureAPI.getAll();
  } catch (e) {
    console.error(e);
    alert('โหลดรายการประเภทคุณสมบัติไม่สำเร็จ');
  }
};

/* ---------------- reset ---------------- */
const resetForm = () => {
  form.value = { type_feature_name: '' };
  editingFeature.value = null;
};

/* ---------------- watch modal ---------------- */
watch(
  () => props.isOpen,
  (open) => {
    if (open) {
      fetchTypeFeatures();
      resetForm();
    }
  },
  { immediate: true }
);

/* ---------------- save (add / edit) ---------------- */
const handleSave = async () => {
  if (!form.value.type_feature_name.trim()) {
    alert('กรุณากรอกชื่อประเภทคุณสมบัติ');
    return;
  }

  try {
    if (editingFeature.value) {
      // แก้ไข
      await TypefeatureAPI.update(editingFeature.value.ID, {
        ID: editingFeature.value.ID,
        type_feature_name: form.value.type_feature_name
      });

      alert('แก้ไขประเภทคุณสมบัติเรียบร้อยแล้ว ✅');
    } else {
      // เพิ่มใหม่
      const createData: TypeFeatureCreate = {
        type_feature_name: form.value.type_feature_name
      };

      await TypefeatureAPI.create(createData);

      alert('เพิ่มประเภทคุณสมบัติเรียบร้อยแล้ว ✅');
    }

    await fetchTypeFeatures();
    resetForm();

    // ❌ ไม่ปิด modal
    emit('saved');
  } catch (e) {
    console.error(e);
    alert('บันทึกข้อมูลไม่สำเร็จ');
  }
};

/* ---------------- edit ---------------- */
const handleEdit = (f: TypeFeatureResponse) => {
  editingFeature.value = f;
  form.value.type_feature_name = f.type_feature_name;
};

/* ---------------- delete ---------------- */
const handleDelete = async (f: TypeFeatureResponse) => {
  if (!confirm('ต้องการลบประเภทคุณสมบัตินี้?')) return;

  try {
    await TypefeatureAPI.delete(f.ID);
    await fetchTypeFeatures();

    if (editingFeature.value?.ID === f.ID) {
      resetForm();
    }

    alert('ลบประเภทคุณสมบัติเรียบร้อยแล้ว');
  } catch (e) {
    console.error(e);
    alert('ลบข้อมูลไม่สำเร็จ');
  }
};
</script>

<template>
  <div v-if="isOpen" class="modal modal-open">
    <div class="modal-box max-w-md">
      <h3 class="font-bold text-lg mb-4">
        {{ editingFeature ? 'แก้ไขประเภทคุณสมบัติ' : 'เพิ่มประเภทคุณสมบัติ' }}
      </h3>

      <!-- Form -->
      <div class="space-y-3">
        <input
          v-model="form.type_feature_name"
          placeholder="ชื่อประเภทคุณสมบัติ"
          class="input input-bordered w-full"
        />

        <div class="flex gap-2">
          <button
            @click="handleSave"
            class="btn btn-primary flex-1"
          >
            {{ editingFeature ? 'บันทึกการแก้ไข' : 'เพิ่มประเภท' }}
          </button>

          <button
            v-if="editingFeature"
            @click="resetForm"
            class="btn flex-1 bg-gray-400 text-white"
          >
            ยกเลิกแก้ไข
          </button>
        </div>
      </div>

      <!-- List -->
      <div class="mt-4 border-t pt-2">
        <h4 class="font-bold mb-2">รายการประเภทคุณสมบัติ</h4>

        <ul>
          <li
            v-for="f in typeFeatures"
            :key="f.ID"
            class="flex justify-between items-center border-b py-1"
          >
            <div>{{ f.type_feature_name }}</div>
            <div class="flex gap-1">
              <button
                @click="handleEdit(f)"
                class="btn btn-xs bg-yellow-400 text-white"
              >
                แก้ไข
              </button>
              <button
                @click="handleDelete(f)"
                class="btn btn-xs bg-red-500 text-white"
              >
                ลบ
              </button>
            </div>
          </li>

          <li
            v-if="typeFeatures.length === 0"
            class="text-gray-400 text-sm py-2"
          >
            ยังไม่มีประเภท
          </li>
        </ul>
      </div>

      <!-- Footer -->
      <div class="modal-action">
        <button class="btn" @click="emit('close')">ปิด</button>
      </div>
    </div>
  </div>
</template>
