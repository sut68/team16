<script setup lang="ts">
import { ref, watch } from 'vue';
import { TypefeatureAPI } from '@/services/api/typefeature';
import type { TypeFeatureResponse, TypeFeatureCreate } from '@/interfaces';

const props = defineProps<{ isOpen: boolean }>();
const emit = defineEmits<{ (e: 'close'): void; (e: 'saved'): void }>();

// State
const typeFeatures = ref<TypeFeatureResponse[]>([]);
const form = ref<{ type_feature_name: string }>({ type_feature_name: '' });
const editingFeature = ref<TypeFeatureResponse | null>(null);

// Fetch all type features
const fetchTypeFeatures = async () => {
  typeFeatures.value = await TypefeatureAPI.getAll();
};

// Reset form
const resetForm = () => {
  form.value = { type_feature_name: '' };
  editingFeature.value = null;
};

// Watch modal open
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

// Save (add / update)
const handleSave = async () => {
  if (!form.value.type_feature_name.trim()) {
    alert('กรุณากรอกชื่อประเภทคุณสมบัติ');
    return;
  }

  if (editingFeature.value) {
    await TypefeatureAPI.update(editingFeature.value.ID, {
      ID: editingFeature.value.ID,
      type_feature_name: form.value.type_feature_name
    });
  } else {
    const createData: TypeFeatureCreate = {
      type_feature_name: form.value.type_feature_name
    };
    await TypefeatureAPI.create(createData);
  }

  await fetchTypeFeatures();
  resetForm();
  emit('saved');
};

// Edit
const handleEdit = (f: TypeFeatureResponse) => {
  editingFeature.value = f;
  form.value.type_feature_name = f.type_feature_name;
};

// Delete
const handleDelete = async (f: TypeFeatureResponse) => {
  if (confirm('ต้องการลบประเภทคุณสมบัตินี้?')) {
    await TypefeatureAPI.delete(f.ID);
    await fetchTypeFeatures();
    if (editingFeature.value?.ID === f.ID) resetForm();
  }
};
</script>

<template>
  <div v-if="isOpen" class="modal modal-open">
    <div class="modal-box max-w-md">
      <h3 class="font-bold text-lg mb-4">
        {{ editingFeature ? 'แก้ไขประเภทคุณสมบัติ' : 'เพิ่มประเภทคุณสมบัติ' }}
      </h3>

      <div class="space-y-3">
        <input
          v-model="form.type_feature_name"
          placeholder="ชื่อประเภทคุณสมบัติ"
          class="input input-bordered w-full"
        />

        <div class="flex gap-2">
          <button @click="handleSave" class="btn btn-primary flex-1">
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
              <button @click="handleEdit(f)" class="btn btn-xs bg-yellow-400 text-white">
                แก้ไข
              </button>
              <button @click="handleDelete(f)" class="btn btn-xs bg-red-500 text-white">
                ลบ
              </button>
            </div>
          </li>
          <li v-if="typeFeatures.length === 0" class="text-gray-400 text-sm py-2">
            ยังไม่มีประเภท
          </li>
        </ul>
      </div>

      <div class="modal-action">
        <button class="btn" @click="$emit('close')">ปิด</button>
      </div>
    </div>
  </div>
</template>