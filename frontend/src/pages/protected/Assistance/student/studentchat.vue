<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick } from "vue"
import { ChatroomAPI } from "@/services/api/chatroom"
import { AssistanceAPI } from "@/services/api/assistance"
import { ChatService } from "@/services/api/websocket"
import type { AssistanceRespone, ChatroomRespone } from "@/interfaces"
import { getMyProfile } from "@/services/api/user"

// ---------------- state ----------------
const messages = ref<AssistanceRespone[]>([])
const input = ref("")
const room = ref<ChatroomRespone | null>(null)
const myId = ref<number>(0)

const chatService = new ChatService()
const token = sessionStorage.getItem("token")

const chatBox = ref<HTMLElement | null>(null)

// ---------------- methods ----------------
const send = () => {
  if (!input.value.trim()) return
  if (!chatService.isConnected()) return

  chatService.sendMessage(input.value)
  input.value = ""
}



const scrollToBottom = () => {
  nextTick(() => {
    if (chatBox.value) {
      chatBox.value.scrollTop = chatBox.value.scrollHeight
    }
  })
}

// ---------------- lifecycle ----------------
onMounted(async () => {
  try {
    // 1. profile
    const profileRes = await getMyProfile()
    myId.value = profileRes.data.ID

    // 2. get open room
    const roomRes = await ChatroomAPI.getMyOpen()
    if (!roomRes.data) {
      const createRes = await ChatroomAPI.create()
      room.value = createRes.data
    } else {
      room.value = roomRes.data
    }

    if (!room.value) return

    // 3. load history
    const historyRes = await AssistanceAPI.getAll()
    const history: AssistanceRespone[] = Array.isArray(historyRes)
      ? historyRes
      : Array.isArray(historyRes?.data)
        ? historyRes.data
        : []

    messages.value = history.filter(m => m.chatroom_id === room.value!.ID)
    scrollToBottom()

    // 4. websocket
    if (!token) {
      console.warn("No token found, cannot connect websocket")
      return
    }

    chatService.connect(token, room.value.ID)
    chatService.onMessage(msg => {
      messages.value.push(msg)
      scrollToBottom()
    })

  } catch (err) {
    console.error("StudentChat error:", err)
  }
})

onUnmounted(() => {
  chatService.disconnect()
})
</script>

<template>
  <div class="max-w-xl mx-auto bg-white rounded-xl shadow flex flex-col h-[600px]">

    <!-- Header -->
    <div class="px-6 py-4 border-b font-semibold text-gray-800">
      📨 ติดต่อและขอความช่วยเหลือ
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
          class="max-w-[75%] px-4 py-2 rounded-2xl text-sm shadow"
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
  </div>
</template>



