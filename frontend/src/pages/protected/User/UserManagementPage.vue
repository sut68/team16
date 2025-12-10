<template>
  <div class="min-h-screen bg-[#f0f2f5] p-6 font-sans text-slate-800">
    <div class="max-w-5xl mx-auto">
      <!-- Header -->
      <div class="flex items-start justify-between mb-6">
        <div>
          <h1 class="text-2xl font-bold text-[#1e3a8a]">User Management</h1>
          <p class="text-sm text-slate-500 mt-1">จัดการบัญชีผู้ใช้ — แสดงชื่อ, username, รหัสนักศึกษา, role</p>
        </div>

        <div class="flex items-center gap-3">
          <input
            v-model="q"
            @input="onSearch"
            type="search"
            placeholder="ค้นหา username หรือ role..."
            class="input input-bordered h-10 w-72 bg-white"
          />
          <button
            class="rounded-md px-4 py-2 text-white font-medium shadow-lg flex items-center gap-2"
            :style="createButtonStyle"
            @click="openCreate"
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
            </svg>
            Create user
          </button>
        </div>
      </div>

      <!-- Card Table -->
      <div class="bg-white rounded-2xl shadow-sm overflow-hidden">
        <!-- top info -->
        <div class="px-6 py-3 border-b text-sm text-slate-600 flex items-center justify-between">
          <div>Showing {{ filtered.length }} users</div>
          <div class="flex items-center gap-3">
            <label class="text-xs text-slate-500 mr-2">Per page</label>
            <select v-model.number="perPage" class="select select-sm select-bordered">
              <option :value="5">5</option>
              <option :value="10">10</option>
              <option :value="20">20</option>
            </select>
          </div>
        </div>

        <!-- headers -->
        <div class="grid grid-cols-12 gap-4 px-6 py-3 bg-[#fafafa] text-xs uppercase tracking-wide text-slate-600">
          <div class="col-span-5">ชื่อ - นามสกุล</div>
          <div class="col-span-2">USERNAME</div>
          <div class="col-span-2">รหัสนักศึกษา</div>
          <div class="col-span-2">ROLE</div>
          <div class="col-span-1 text-right">ACTIONS</div>
        </div>

        <!-- rows -->
        <div v-if="loading" class="p-6 text-center text-slate-500">
          <span class="loading loading-spinner"></span> กำลังโหลด...
        </div>

        <div v-else>
          <div v-if="!filtered.length" class="p-6 text-center text-slate-500">ไม่มีผู้ใช้</div>

          <div v-for="u in paged" :key="u.ID" class="grid grid-cols-12 gap-4 px-6 py-4 items-center border-b hover:bg-slate-50">
            <!-- name -->
            <div class="col-span-5 flex items-center gap-4">
              <div class="h-10 w-10 rounded-full bg-[#e8eefb] text-[#1e3a8a] flex items-center justify-center font-semibold">
                {{ initials(u.username) }}
              </div>
              <div>
                <div class="text-sm font-medium text-slate-900">{{ profileName(u) }}</div>
                <div class="text-xs text-slate-400">ID: {{ u.ID }}</div>
              </div>
            </div>

            <!-- username -->
            <div class="col-span-2 text-sm text-slate-700">{{ u.username }}</div>

            <!-- student id -->
            <div class="col-span-2 text-sm text-slate-700">
              <span v-if="studentIdOf(u)">{{ studentIdOf(u) }}</span>
              <span v-else class="text-slate-400">—</span>
            </div>

            <!-- role -->
            <div class="col-span-2">
              <span :class="roleBadgeClass(u)" class="px-3 py-1 rounded-full text-xs font-semibold">
                {{ (u.role?.name || guessRole(u.role_id)).toUpperCase() }}
              </span>
            </div>

            <!-- actions -->
            <div class="col-span-1 text-right">
              <button class="text-indigo-600 hover:underline mr-3" @click="onEdit(u)" :disabled="!canEdit(u)">แก้ไข</button>
              <button class="text-red-500 hover:underline" @click="confirmDelete(u)" :disabled="!canDelete(u)">ลบ</button>
            </div>
          </div>
        </div>

        <!-- pagination footer -->
        <div class="px-6 py-3 flex items-center justify-between text-sm text-slate-600">
          <div>Page {{ page }} / {{ totalPages }}</div>
          <div class="flex gap-2">
            <button class="btn btn-ghost btn-sm" :disabled="page <= 1" @click="page--">Prev</button>
            <button class="btn btn-ghost btn-sm" :disabled="page >= totalPages" @click="page++">Next</button>
          </div>
        </div>
      </div>
    </div>

    <!-- Create Modal -->
    <teleport to="body">
      <div v-if="showCreate" class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4" @click.self="closeCreate">
        <div class="w-full max-w-lg bg-white rounded-2xl shadow-lg overflow-hidden">
          <div class="px-6 py-4 border-b flex items-center justify-between bg-[#fbfbff]">
            <h3 class="text-lg font-semibold text-[#1e3a8a]">Create User</h3>
            <button class="btn btn-ghost" @click="closeCreate">✕</button>
          </div>

          <form @submit.prevent="submitCreate" class="p-6 grid gap-4">
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
              <div>
                <label class="block text-sm text-slate-700 mb-1">Username *</label>
                <input v-model="form.username" class="input input-bordered w-full" />
                <p v-if="errors.username" class="text-xs text-red-500 mt-1">{{ errors.username }}</p>
              </div>

              <div>
                <label class="block text-sm text-slate-700 mb-1">Password *</label>
                <input v-model="form.password" type="password" class="input input-bordered w-full" />
                <p v-if="errors.password" class="text-xs text-red-500 mt-1">{{ errors.password }}</p>
              </div>
            </div>

            <div>
              <label class="block text-sm text-slate-700 mb-1">Role</label>
              <select v-model="form.role" class="select select-bordered w-full">
                <option value="student">Student</option>
                <option value="admin">Admin</option>
                <option value="admin" v-if="myID === 1">Admin</option>
              </select>
              <p v-if="errors.role" class="text-xs text-red-500 mt-1">{{ errors.role }}</p>
            </div>

            <div class="flex justify-end gap-3 mt-4">
              <button type="button" class="btn btn-ghost" @click="closeCreate">Cancel</button>
              <button type="submit" class="btn btn-primary" :disabled="loadingCreate">
                <span v-if="loadingCreate" class="loading loading-spinner mr-2"></span>
                Create
              </button>
            </div>
          </form>
        </div>
      </div>
    </teleport>

    <!-- Delete Confirmation -->
    <teleport to="body">
      <div v-if="deleting" class="fixed inset-0 z-60 flex items-center justify-center p-4 bg-black/30">
        <div class="bg-white rounded-lg p-6 w-full max-w-md shadow-lg">
          <h3 class="text-lg font-semibold mb-2">ยืนยันการลบ</h3>
          <p class="text-sm text-slate-600 mb-4">ต้องการลบผู้ใช้ <strong>{{ deleting.username }}</strong> หรือไม่? การลบไม่สามารถกู้คืนได้</p>
          <div class="flex justify-end gap-2">
            <button class="btn btn-ghost" @click="deleting = null">ยกเลิก</button>
            <button class="btn btn-error" @click="doDelete" :disabled="loadingDelete">
              <span v-if="loadingDelete" class="loading loading-spinner mr-2"></span>ลบ
            </button>
          </div>
        </div>
      </div>
    </teleport>

    <!-- toast -->
    <div v-if="toast" class="fixed right-6 bottom-6 z-60">
      <div :class="['px-4 py-2 rounded-md shadow-lg text-white', toast.type === 'error' ? 'bg-red-600' : 'bg-green-600']">
        {{ toast.msg }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, nextTick } from 'vue'
import type { UserResponse, StudentProfileResponse, AdminProfileResponse } from '@/interfaces/user'
import * as UserAPI from '@/services/api/user'
import axios from 'axios'

// state
const users = ref<UserResponse[]>([])
const studentProfiles = ref<Record<number, StudentProfileResponse | undefined>>({})
const adminProfiles = ref<Record<number, AdminProfileResponse | undefined>>({})
const loading = ref(false)

const q = ref('')
const page = ref(1)
const perPage = ref(10)

const showCreate = ref(false)
const loadingCreate = ref(false)
const loadingDelete = ref(false)
const deleting = ref<UserResponse | null>(null)

const form = ref({ username: '', password: '', role: 'student' })
const errors = ref<Record<string,string>>({})
const toast = ref<{ msg: string; type?: 'ok'|'error' } | null>(null)

// me from localStorage
const me = JSON.parse(localStorage.getItem('me') || 'null')
const myID = me?.ID ?? me?.id ?? me?.user_id
const myRole = (me?.role?.name || '').toLowerCase()

const createButtonStyle = "background: linear-gradient(90deg,#1e3a8a,#3b82f6); padding: 10px 14px; border-radius: 8px;"

// helpers
function initials(name?: string) { if (!name) return 'U'; return name.slice(0,2).toUpperCase() }
function guessRole(role_id?: number | null) { if (role_id === 1) return 'admin'; if (role_id === 2) return 'student'; return 'user' }
function roleBadgeClass(u: UserResponse) {
  const r = (u.role?.name || guessRole(u.role_id)).toLowerCase()
  if (r === 'admin') return 'bg-red-100 text-red-700'
  if (r === 'student') return 'bg-green-100 text-green-700'
  return 'bg-slate-100 text-slate-700'
}
function profileName(u: UserResponse) {
  const st = studentProfiles.value[u.ID as number]
  if (st) return `${st.first_name_th} ${st.last_name_th}`
  const ad = adminProfiles.value[u.ID as number]
  if (ad) return `${ad.admin_firstname} ${ad.admin_lastname}`
  return u.username
}
function studentIdOf(u: UserResponse) {
  const st = studentProfiles.value[u.ID as number]
  if (st) return st.student_id
  return null
}

// load data
async function load() {
  loading.value = true
  try {
    users.value = await UserAPI.listUsers()

    // Try bulk fetch profiles if backend supports it; otherwise it's ok to skip.
    try {
      const res = await axios.get<StudentProfileResponse[]>('/api/profile/all-student')
      res.data.forEach(s => studentProfiles.value[s.user_id] = s)
    } catch (e) {
      // ignore if endpoint not present
    }

    try {
      const res2 = await axios.get<AdminProfileResponse[]>('/api/profile/all-admin')
      res2.data.forEach(a => adminProfiles.value[a.user_id] = a)
    } catch (e) {
      // ignore
    }

  } catch (e:any) {
    showToast('ไม่สามารถโหลดผู้ใช้ได้', 'error')
    console.error(e)
  } finally {
    loading.value = false
  }
}
onMounted(load)

// search & pagination
const filtered = computed(() => {
  const term = q.value.trim().toLowerCase()
  if (!term) return users.value
  return users.value.filter(u => {
    const uname = (u.username||'').toLowerCase()
    const role = (u.role?.name||'').toLowerCase()
    const pname = (profileName(u)||'').toLowerCase()
    return uname.includes(term) || role.includes(term) || pname.includes(term)
  })
})
const totalPages = computed(() => Math.max(1, Math.ceil(filtered.value.length / perPage.value)))
const paged = computed(() => {
  const start = (page.value - 1) * perPage.value
  return filtered.value.slice(start, start + perPage.value)
})
function onSearch() { page.value = 1 }

// create
function openCreate() { showCreate.value = true; nextTick(()=> { /* autofocus if needed */ }) }
function closeCreate() { showCreate.value = false; form.value = { username:'', password:'', role:'student' }; errors.value = {} }
function validateCreate() {
  errors.value = {}
  if (!form.value.username || form.value.username.trim().length < 3) errors.value.username = 'โปรดระบุ username อย่างน้อย 3 ตัว'
  if (!form.value.password || form.value.password.length < 4) errors.value.password = 'รหัสผ่านต้องอย่างน้อย 4 ตัว'
  if (form.value.role === 'admin' && myID !== 1) errors.value.role = 'เฉพาะ admin (ID=1) สร้าง Admin ได้'
  return Object.keys(errors.value).length === 0
}
async function submitCreate() {
  if (!validateCreate()) return
  loadingCreate.value = true
  try {
    await UserAPI.createUser({ username: form.value.username.trim(), password: form.value.password, role: (form.value.role||'student').toLowerCase() })
    showToast('สร้างผู้ใช้เรียบร้อย','ok')
    await load()
    closeCreate()
  } catch (e:any) {
    showToast(e?.response?.data?.error || e?.message || 'สร้างผิดพลาด','error')
    console.error(e)
  } finally { loadingCreate.value = false }
}

// delete
function confirmDelete(u: UserResponse) { deleting.value = u }
async function doDelete() {
  if (!deleting.value) return
  loadingDelete.value = true
  try {
    await UserAPI.deleteUser(deleting.value.ID!)
    showToast('ลบผู้ใช้เรียบร้อย','ok')
    await load()
    deleting.value = null
  } catch (e:any) {
    showToast(e?.response?.data?.error || e?.message || 'ลบผิดพลาด','error')
    console.error(e)
  } finally { loadingDelete.value = false }
}

// actions / permissions
function canDelete(u: UserResponse) {
  if (!me) return false
  const myid = myID
  if (!myid) return false
  if (myid === u.ID) return false
  const targetRole = (u.role?.name || guessRole(u.role_id)).toLowerCase()
  if (targetRole === 'admin' && myid !== 1) return false
  if (!(myid === 1 || myRole === 'admin')) return false
  return true
}
function canEdit(u: UserResponse) {
  return canDelete(u) || (myID === u.ID)
}
function onEdit(u: UserResponse) {
  // TODO: implement edit flow (open modal / navigate)
  alert(`Implement edit for ${u.username}`)
}

// toast
function showToast(msg: string, type: 'ok'|'error' = 'ok') {
  toast.value = { msg, type }
  setTimeout(()=> (toast.value = null), 3200)
}
</script>

<style scoped>
.input, .select { background-color: white; }
.btn-primary { background-image: linear-gradient(90deg,#1e3a8a,#3b82f6); border: none; color: white; }
.btn-error { background-image: linear-gradient(90deg,#dc2626,#ef4444); border: none; color: white; }
.loading-spinner { border-width: 3px; height: 1rem; width: 1rem; border-radius: 9999px; border-top-color: transparent; border-right-color: transparent; animation: spin 1s linear infinite; display:inline-block; vertical-align:middle; }
@keyframes spin { to { transform: rotate(360deg); } }
</style>
