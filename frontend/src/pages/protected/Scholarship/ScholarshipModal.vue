<script setup lang="ts">
import { ref, watch } from 'vue';
import type {
  ScholarshipCreate,
  ScholarshipUpdate,
  ScholarshipResponse,
  ScholarshipStatusResponse,
  ScholarshipTypeResponse,
  SponsorResponse,
  TypeFeatureResponse,
  FeatureScholarshipResponse,
  FeatureScholarshipCreate,
  SemasterResponse,
} from '@/interfaces';

import { ScholarshipAPI } from '@/services/api/scholarship';
import { StatusscholarshipAPI } from '@/services/api/statusscholarship';
import { TypescholarshipAPI } from '@/services/api/typescholarship';
import { SponsorService } from '@/services/sponsor/sponsor';
import { FeatureScholarshipAPI } from '@/services/api/featurescholarship';
import { TypefeatureAPI } from '@/services/api/typefeature';

/* ---------------- props / emits ---------------- */
interface Props {
  isOpen: boolean;
  scholarship: ScholarshipResponse | null;
}
const props = defineProps<Props>();
const emit = defineEmits(['close', 'saved']);

/* ---------------- dropdown data ---------------- */
const statusOptions = ref<ScholarshipStatusResponse[]>([]);
const typeOptions = ref<ScholarshipTypeResponse[]>([]);
const sponsorOptions = ref<SponsorResponse[]>([]);
const typeFeatureOptions = ref<TypeFeatureResponse[]>([]);
const semasterOptions = ref<SemasterResponse[]>([]);

/* ---------------- scholarship form ---------------- */
const form = ref<ScholarshipCreate>({
  scholarship_name: '',
  description: '',
  open_date: '',
  close_date: '',
  statusscholarship_id: 0,
  typescholarship_id: 0,
  sponsor_id: 0,
  semaster_id: 0,
});

/* ---------------- feature form ---------------- */
const featureForm = ref<{
  typefeature_id: number | null;
  operator: string;
  value: string;
}>({
  typefeature_id: null,
  operator: '',
  value: '',
});

type LocalFeature = FeatureScholarshipResponse & { tempId?: number };
const featureList = ref<LocalFeature[]>([]);
let nextTempId = 1;

/* ---------------- load data ---------------- */
const loadDropdownData = async () => {
  try {
    const [
      typesRes,
      statusesRes,
      sponsorsRes,
      typeFeaturesRes,
      featuresRes,
      scholarshipsRes,
    ] = await Promise.all([
      TypescholarshipAPI.getAll(),
      StatusscholarshipAPI.getAll(),
      SponsorService.getAll(),
      TypefeatureAPI.getAll(),
      FeatureScholarshipAPI.getAll(),
      ScholarshipAPI.getAll(),
    ]);

    typeOptions.value = typesRes ?? [];
    statusOptions.value = statusesRes ?? [];
    sponsorOptions.value = sponsorsRes ?? [];
    typeFeatureOptions.value = typeFeaturesRes ?? [];

    /* ดึง semaster จาก scholarship (ไม่ซ้ำ) */
    const semasterMap = new Map<number, SemasterResponse>();
    (scholarshipsRes as ScholarshipResponse[]).forEach(s => {
      if (s.semaster) {
        semasterMap.set(s.semaster.ID, s.semaster);
      }
    });
    semasterOptions.value = Array.from(semasterMap.values());

    if (props.scholarship) {
      form.value = {
        scholarship_name: props.scholarship.scholarship_name,
        description: props.scholarship.description,
        open_date: props.scholarship.open_date,
        close_date: props.scholarship.close_date,
        statusscholarship_id: props.scholarship.statusscholarship_id,
        typescholarship_id: props.scholarship.typescholarship_id,
        sponsor_id: props.scholarship.sponsor_id,
        semaster_id: props.scholarship.semaster.ID,
      };

      featureList.value = (featuresRes as FeatureScholarshipResponse[]).filter(
        f => f.scholarship_id === props.scholarship!.ID
      );
    } else {
      form.value = {
        scholarship_name: '',
        description: '',
        open_date: '',
        close_date: '',
        statusscholarship_id: statusOptions.value[0]?.ID ?? 0,
        typescholarship_id: typeOptions.value[0]?.ID ?? 0,
        sponsor_id: sponsorOptions.value[0]?.ID ?? 0,
        semaster_id: semasterOptions.value[0]?.ID ?? 0,
      };
      featureList.value = [];
    }
  } catch (e) {
    console.error(e);
    alert('โหลดข้อมูลไม่สำเร็จ');
  }
};

watch(
  () => props.isOpen,
  (open) => {
    if (open) loadDropdownData();
  },
  { immediate: true }
);

/* ---------------- submit ---------------- */
const handleSubmit = async () => {
  try {
    if (!form.value.scholarship_name) {
      return alert('กรุณากรอกชื่อทุน');
    }
    if (!form.value.description){
      return alert('กรุณากรอกรายละเอียด');
    }
    if (form.value.description.length < 5){
      return alert('กรุณากรอกรายละเอียดอย่างน้อย 5 ตัวอักษร')
    }
    if (!form.value.open_date){
      return alert('กรุณาเลือกวันเริ่มต้น')
    }
    if (!form.value.semaster_id) {
      return alert('กรุณาเลือกภาคการศึกษา');
    }

    let scholarshipId: number | undefined = props.scholarship?.ID;

    if (props.scholarship) {
      await ScholarshipAPI.update(
        props.scholarship.ID,
        form.value as ScholarshipUpdate
      );
      scholarshipId = props.scholarship.ID;
    } else {
      const created = await ScholarshipAPI.create(
        form.value as ScholarshipCreate
      );
      scholarshipId = (created as ScholarshipResponse).ID;
    }

    if (!scholarshipId) throw new Error('Missing scholarship ID');

    if (props.scholarship) {
      const allFeatures =
        (await FeatureScholarshipAPI.getAll()) as FeatureScholarshipResponse[];
      const toDelete = allFeatures.filter(
        f => f.scholarship_id === scholarshipId
      );
      for (const f of toDelete) {
        await FeatureScholarshipAPI.delete(f.ID);
      }
    }

    for (const f of featureList.value) {
      if (!f.typefeature_id || !f.operator) continue;
      const payload: FeatureScholarshipCreate = {
        scholarship_id: scholarshipId,
        typefeature_id: f.typefeature_id,
        operator: f.operator,
        value: f.value,
      };
      await FeatureScholarshipAPI.create(payload);
    }

    emit('saved');
  } catch (e) {
    console.error(e);
    alert('บันทึกข้อมูลไม่สำเร็จ');
  }
};

/* ---------------- feature actions ---------------- */
const addFeature = () => {
  if (!featureForm.value.typefeature_id) {
    return alert('เลือกประเภทคุณสมบัติ');
  }
  if (!featureForm.value.operator || !featureForm.value.value) {
    return alert('กรอกข้อมูลให้ครบ');
  }

  const tf =
    typeFeatureOptions.value.find(
      t => t.ID === featureForm.value.typefeature_id
    ) || null;

  featureList.value.push({
    ID: 0,
    scholarship_id: props.scholarship?.ID ?? 0,
    typefeature_id: featureForm.value.typefeature_id!,
    Typefeature: tf as any,
    operator: featureForm.value.operator,
    value: featureForm.value.value,
    tempId: nextTempId++,
  });

  featureForm.value = { typefeature_id: null, operator: '', value: '' };
};

const deleteFeature = (f: LocalFeature) => {
  if (!confirm('ลบคุณสมบัตินี้?')) return;
  featureList.value = featureList.value.filter(i =>
    i.ID > 0 ? i.ID !== f.ID : i.tempId !== f.tempId
  );
};
</script>

<template>
  <div v-if="isOpen" class="modal modal-open">
    <div class="modal-box max-w-xl">
      <h2 class="font-bold text-lg mb-3">
        {{ scholarship ? 'แก้ไขทุน' : 'เพิ่มทุนใหม่' }}
      </h2>

      <div class="space-y-3">
        <input
          v-model="form.scholarship_name"
          class="input input-bordered w-full"
          placeholder="ชื่อทุน"
        />

        <textarea
          v-model="form.description"
          class="textarea textarea-bordered w-full"
          rows="3"
          placeholder="รายละเอียดทุน"
        />

        <!-- เลือกภาคการศึกษา -->
        <select v-model.number="form.semaster_id" class="select select-bordered w-full">
          <option :value="0" disabled>-- เลือกภาคการศึกษา --</option>
          <option
            v-for="s in semasterOptions"
            :key="s.ID"
            :value="s.ID"
          >
            ปี {{ s.academic_year }} / เทอม {{ s.term }} / รอบ {{ s.round }}
          </option>
        </select>

        <select v-model.number="form.typescholarship_id" class="select select-bordered w-full">
          <option :value="0" disabled>-- เลือกประเภททุน --</option>
          <option v-for="t in typeOptions" :key="t.ID" :value="t.ID">
            {{ t.type_name }}
          </option>
        </select>

        <select v-model.number="form.statusscholarship_id" class="select select-bordered w-full">
          <option :value="0" disabled>-- เลือกสถานะทุน --</option>
          <option v-for="s in statusOptions" :key="s.ID" :value="s.ID">
            {{ s.status_name }}
          </option>
        </select>

        <select v-model.number="form.sponsor_id" class="select select-bordered w-full">
          <option :value="0" disabled>-- เลือกผู้สนับสนุน --</option>
          <option v-for="sp in sponsorOptions" :key="sp.ID" :value="sp.ID">
            {{ sp.company_name }}
          </option>
        </select>

        <div class="grid grid-cols-2 gap-2">
          <input v-model="form.open_date" type="date" class="input input-bordered w-full" />
          <input v-model="form.close_date" type="date" class="input input-bordered w-full" />
        </div>

        <!-- Feature -->
        <div class="border-t pt-3">
          <h3 class="font-bold mb-2">คุณสมบัติทุน</h3>

          <div class="flex gap-2 mb-2">
            <select v-model.number="featureForm.typefeature_id" class="select select-bordered flex-1">
              <option :value="null" disabled>-- ประเภท --</option>
              <option v-for="f in typeFeatureOptions" :key="f.ID" :value="f.ID">
                {{ f.type_feature_name }}
              </option>
            </select>

            <input v-model="featureForm.operator" class="input input-bordered w-24" placeholder="Op" />
            <input v-model="featureForm.value" class="input input-bordered flex-1" placeholder="ค่า" />
            <button class="btn btn-primary" @click="addFeature">เพิ่ม</button>
          </div>

          <ul>
            <li
              v-for="f in featureList"
              :key="f.ID > 0 ? f.ID : f.tempId"
              class="flex justify-between items-center border-b py-1"
            >
              <span>
                {{ f.Typefeature?.type_feature_name }}
                {{ f.operator }}
                {{ f.value }}
                <span v-if="f.ID === 0" class="text-xs text-gray-400 ml-1">
                  (ยังไม่บันทึก)
                </span>
              </span>
              <button class="btn btn-xs btn-error" @click="deleteFeature(f)">
                ลบ
              </button>
            </li>

            <li v-if="featureList.length === 0" class="text-gray-400 text-sm">
              ยังไม่มีคุณสมบัติ
            </li>
          </ul>
        </div>
      </div>

      <div class="modal-action">
        <button class="btn" @click="$emit('close')">ยกเลิก</button>
        <button class="btn btn-primary" @click="handleSubmit">บันทึก</button>
      </div>
    </div>
  </div>
</template>
