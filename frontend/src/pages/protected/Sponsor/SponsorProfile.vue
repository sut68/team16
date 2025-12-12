<script setup lang="ts">
import { ref } from 'vue'

type Contact = {
  id: number
  name: string
  position?: string
  email?: string
  phone?: string
}

type Sponsor = {
  id: number | null
  name: string
  logo_url: string
  industry?: { id: number; name: string } | null
  email?: string
  phone?: string
  website?: string
  scholarship_count?: number
  contacts?: Contact[]
}

// --- Mock data (ให้เห็นหน้าได้เลย) ---
const sponsor = ref<Sponsor>({
  id: 1,
  name: 'ACME Scholarship Foundation',
  logo_url: 'https://via.placeholder.com/150?text=Logo',
  industry: { id: 10, name: 'Technology' },
  email: 'contact@acme.example',
  phone: '+66 2 123 4567',
  website: 'www.acme.example',
  scholarship_count: 5,
  contacts: []
})

// แยก contacts เพื่อให้ใช้งานง่าย (template ใช้ contacts)
const contacts = ref<Contact[]>([
  { id: 1, name: 'นางสาวปรียา ชาญชัย', position: 'HR Manager', email: 'priya@acme.example', phone: '+66 89 111 2222' },
  { id: 2, name: 'นายสมชาย ใจดี', position: 'PR Officer', email: 'somchai@acme.example', phone: '+66 85 333 4444' }
])

// ถ้าต้องการ ให้ sponsor.contacts ชี้ไปที่ contacts (ไม่บังคับ)
sponsor.value.contacts = contacts.value

// --- Simple handlers (local only, ให้ลองคลิกดูงาน) ---
function onAddContact() {
  // สร้าง dummy contact ใหม่แล้ว push เข้า list (เพื่อ preview หน้า)
  const nextId = (contacts.value.length ? Math.max(...contacts.value.map(c => c.id)) : 0) + 1
  const newC: Contact = {
    id: nextId,
    name: `New Contact ${nextId}`,
    position: 'Position',
    email: `new${nextId}@example.com`,
    phone: 'N/A'
  }
  contacts.value.push(newC)
  sponsor.value.contacts = contacts.value
  window.alert('เพิ่มผู้ติดต่อตัวอย่างแล้ว — เปลี่ยนข้อมูลได้เอง')
}

// function onEditContact(c: Contact) {
//   const newName = window.prompt('แก้ไขชื่อผู้ติดต่อ', c.name)
//   if (newName != null) {
//     const idx = contacts.value.findIndex(x => x.id === c.id)
//     if (idx !== -1) {
//       contacts.value[idx].name = newName   
//     }
//     sponsor.value.contacts = contacts.value 
//   }
// }

function onDeleteContact(c: Contact) {
  const ok = window.confirm(`ลบผู้ติดต่อ ${c.name} จริงหรือไม่?`)
  if (!ok) return
  contacts.value = contacts.value.filter(x => x.id !== c.id)
  sponsor.value.contacts = contacts.value
}

// เพื่อให้ template ใช้ชื่อเดียวกับเดิม (ปุ่มใน template ต้องเรียกชื่อเหล่านี้)
const addContact = onAddContact
// const editContact = onEditContact
const deleteContact = onDeleteContact
</script>


<template>
  <div class="space-y-6">
    <!-- ================= Row 1: Summary ================= -->
    <section class="bg-white p-6 rounded-lg shadow">
      <div class="flex items-center gap-6">
        <img
          :src="sponsor.logo_url"
          class="w-20 h-20 rounded object-cover"
          alt="logo"
        />

        <div class="flex-1">
          <h1 class="text-2xl font-semibold">{{ sponsor.name }}</h1>

          <p class="text-gray-500">
            Industry: {{ sponsor.industry?.name || '-' }}
          </p>

          <div class="flex gap-4 mt-2 text-sm text-gray-600">
            <span>📧 {{ sponsor.email }}</span>
            <span>📞 {{ sponsor.phone }}</span>
            <span>🌐 {{ sponsor.website }}</span>
          </div>
        </div>

        <div class="text-right">
          <div class="text-lg font-semibold">{{ sponsor.scholarship_count }}</div>
          <div class="text-xs text-gray-500">จำนวนทุน</div>
        </div>
      </div>
    </section>

    <!-- ================= Row 2: Main + Contact ================= -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      
      <!-- Main Content (ทุน + รายละเอียด) -->
      <section class="lg:col-span-2 bg-white p-6 rounded-lg shadow min-h-[300px]">
        <h2 class="text-xl font-semibold mb-4">ข้อมูลทุนทั้งหมด</h2>
        <div class="text-gray-500 text-sm">
          (เตรียมวางข้อมูลทุนในอนาคต)
        </div>
      </section>

      <!-- Contacts -->
      <section class="bg-white p-6 rounded-lg shadow">
        <div class="flex justify-between items-center mb-4">
          <h2 class="text-xl font-semibold">ผู้ติดต่อบริษัท</h2>
          <button class="px-3 py-1 text-sm bg-blue-600 text-white rounded">
            + เพิ่มผู้ติดต่อ
          </button>
        </div>

        <div v-if="contacts?.length === 0" class="text-gray-500 text-sm">
          ยังไม่มีผู้ติดต่อ
        </div>

        <div v-for="c in contacts" :key="c.id" class="border-b py-3">
          <div class="font-medium">{{ c.name }}</div>
          <div class="text-sm text-gray-500">{{ c.position }}</div>
          <div class="text-sm text-gray-500">📧 {{ c.email }}</div>
          <div class="text-sm text-gray-500">📞 {{ c.phone }}</div>

          <div class="flex gap-2 mt-2">
            <button class="text-blue-600 text-xs">แก้ไข</button>
            <button class="text-red-600 text-xs">ลบ</button>
          </div>
        </div>
      </section>

    </div>
  </div>
</template>
