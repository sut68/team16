<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick } from "vue"
import { ChatroomAPI } from "@/services/api/chatroom"
import { AssistanceAPI } from "@/services/api/assistance"
import { ChatService } from "@/services/api/websocket"
import type { ChatroomRespone, AssistanceRespone } from "@/interfaces"
import { getMyProfile } from "@/services/api/user"

// ---------------- state ----------------
const rooms = ref<ChatroomRespone[]>([])
const selectedRoom = ref<ChatroomRespone | null>(null)
const messages = ref<AssistanceRespone[]>([])
const input = ref("")
const myId = ref<number>(0)

const chatService = new ChatService()
const token = sessionStorage.getItem("token")
const chatBox = ref<HTMLElement | null>(null)

// ---------------- methods ----------------
const scrollToBottom = () => {
  nextTick(() => {
    if (chatBox.value) chatBox.value.scrollTop = chatBox.value.scrollHeight
  })
}

const selectRoom = async (room: ChatroomRespone) => {
  selectedRoom.value = room
  messages.value = []

  chatService.disconnect()

  try {
    // load history
    const historyRes = await AssistanceAPI.getAll()
    const history: AssistanceRespone[] = Array.isArray(historyRes)
      ? historyRes
      : Array.isArray(historyRes?.data)
        ? historyRes.data
        : []

    messages.value = history.filter(m => m.chatroom_id === room.ID)
    scrollToBottom()

    // connect websocket
    if (!token) return
    chatService.connect(token, room.ID)
    chatService.onMessage(msg => {
      messages.value.push(msg)
      scrollToBottom()
    })
  } catch (err) {
    console.error("Error loading room:", err)
  }
}

const send = () => {
  if (!input.value.trim()) return
  if (!chatService.isConnected()) return
  if (!selectedRoom.value) return

  chatService.sendMessage(input.value)
  input.value = ""
}



const closeRoom = async () => {
  if (!selectedRoom.value) return
  try {
    await ChatroomAPI.update(selectedRoom.value.ID, { status_chatroom: "close" })
    chatService.disconnect()
    selectedRoom.value = null
    messages.value = []
    await loadRooms()
  } catch (err) {
    console.error("Failed to close room:", err)
  }
}

const loadRooms = async () => {
  try {
    const profileRes = await getMyProfile()
    myId.value = profileRes.data.ID

    const res = await ChatroomAPI.getAllOpen()
    rooms.value = Array.isArray(res.data) ? res.data : []

    if (rooms.value.length > 0 && !selectedRoom.value) {
      selectRoom(rooms.value[0]!)
    }
  } catch (err) {
    console.error("Failed to load rooms:", err)
  }
}

// ---------------- lifecycle ----------------
onMounted(() => loadRooms())
onUnmounted(() => chatService.disconnect())
</script>

<template>
  <div class="flex h-[600px] bg-white rounded-xl shadow overflow-hidden">

    <!-- Sidebar -->
    <aside class="w-64 border-r bg-gray-50">
      <div class="p-4 font-semibold text-gray-700">
        📂 ห้องแชท
      </div>

      <ul class="overflow-y-auto h-full">
        <li
          v-for="r in rooms"
          :key="r.ID"
          @click="selectRoom(r)"
          class="px-4 py-3 cursor-pointer flex justify-between items-center
                 hover:bg-gray-100"
          :class="r.ID === selectedRoom?.ID ? 'bg-blue-50 border-l-4 border-blue-500' : ''"
        >
          <span class="text-sm font-medium text-gray-700">
            {{ r.user?.username || `User ID: ${r.user_id}` }}
          </span>
          <span>
            <span v-if="r.status_chatroom === 'open'" class="text-green-500">●</span>
            <span v-else class="text-red-500">●</span>
          </span>
        </li>
      </ul>
    </aside>

    <!-- Chat Panel -->
    <section class="flex-1 flex flex-col" v-if="selectedRoom">
      <!-- Header -->
      <div class="px-6 py-4 border-b flex justify-between items-center">
        <h3 class="font-semibold text-gray-800">
          💬 แชทกับ {{ selectedRoom.user?.username || `User ID: ${selectedRoom.user_id}` }}
        </h3>
        <button
          @click="closeRoom"
          class="px-4 py-2 text-sm rounded-lg bg-red-500 text-white hover:bg-red-600"
        >
          ปิดห้อง
        </button>
      </div>

      <!-- Messages -->
      <div
        ref="chatBox"
        class="flex-1 overflow-y-auto px-6 py-4 space-y-3 bg-gray-50"
      >
        <div
          v-for="msg in messages"
          :key="msg.ID"
          class="flex"
          :class="msg.sender_id === myId ? 'justify-end' : 'justify-start'"
        >
          <div
            class="max-w-[70%] px-4 py-2 rounded-2xl text-sm shadow"
            :class="msg.sender_id === myId
              ? 'bg-blue-500 text-white rounded-br-none'
              : 'bg-white text-gray-800 rounded-bl-none'"
          >
            <div class="text-xs opacity-70 mb-1">
              {{ msg.sender.username }}
            </div>
            {{ msg.massage }}
          </div>
        </div>
      </div>

      <!-- Input -->
      <div class="p-4 border-t flex gap-2">
        <input
          v-model="input"
          @keyup.enter="send"
          placeholder="พิมพ์ข้อความ..."
          class="flex-1 px-4 py-2 rounded-lg border focus:outline-none focus:ring focus:ring-blue-200"
        />
        <button
          @click="send"
          class="px-5 py-2 rounded-lg bg-blue-500 text-white hover:bg-blue-600"
        >
          ส่ง
        </button>
      </div>
    </section>

    <!-- Empty state -->
    <section v-else class="flex-1 flex items-center justify-center text-gray-400">
      เลือกห้องแชททางซ้าย
    </section>
  </div>
</template>



