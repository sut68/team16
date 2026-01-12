<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick, computed } from "vue"
import { ChatroomAPI } from "@/services/api/chatroom"
import { AssistanceAPI } from "@/services/api/assistance"
import { ChatService } from "@/services/api/websocket"
import type { AssistanceRespone, ChatroomRespone } from "@/interfaces"
import { getMyProfile } from "@/services/api/user"
import { UserRoundCheck, UserRoundCog, Search, Send, MessageSquare, X } from "lucide-vue-next"

// ---------------- state ----------------
const messages = ref<AssistanceRespone[]>([])
const input = ref("")
const room = ref<ChatroomRespone | null>(null)
const myId = ref<number>(0)
const isLoading = ref(true)

const chatService = new ChatService()
const token = sessionStorage.getItem("token")

const chatBox = ref<HTMLElement | null>(null)
const isSearching = ref(false)
const messageSearchQuery = ref("")

// ---------------- helpers ----------------
const formatTime = (dateStr: string) => {
  if (!dateStr) return ""
  const date = new Date(dateStr)
  return date.toLocaleTimeString("th-TH", { hour: "2-digit", minute: "2-digit" })
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return ""
  const date = new Date(dateStr)
  const today = new Date()
  const yesterday = new Date(today)
  yesterday.setDate(yesterday.getDate() - 1)
  
  if (date.toDateString() === today.toDateString()) {
    return "วันนี้"
  } else if (date.toDateString() === yesterday.toDateString()) {
    return "เมื่อวาน"
  }
  return date.toLocaleDateString("th-TH", { day: "numeric", month: "short" })
}

// Filter messages based on search query
const filteredMessages = computed(() => {
  if (!messageSearchQuery.value.trim()) return messages.value
  const query = messageSearchQuery.value.toLowerCase()
  return messages.value.filter(m => m.massage.toLowerCase().includes(query))
})

// Group messages by date
const groupedMessages = computed(() => {
  const groups: { date: string; messages: AssistanceRespone[] }[] = []
  let currentDate = ""
  
  filteredMessages.value.forEach(msg => {
    const msgDate = new Date(msg.CreatedAt).toDateString()
    if (msgDate !== currentDate) {
      currentDate = msgDate
      groups.push({ date: msg.CreatedAt, messages: [msg] })
    } else {
      groups[groups.length - 1]!.messages.push(msg)
    }
  })
  
  return groups
})

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

const toggleSearch = () => {
  isSearching.value = !isSearching.value
  if (!isSearching.value) {
    messageSearchQuery.value = ""
  }
}

// ---------------- lifecycle ----------------
onMounted(async () => {
  isLoading.value = true
  try {
    // 1. profile
    const profileRes = await getMyProfile()
    myId.value = profileRes.data.ID

    // 2. get open room
    const roomRes = await ChatroomAPI.getMyOpen()
    if (!roomRes || !roomRes.data) {
      const createRes = await ChatroomAPI.create()
      room.value = createRes
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
  } finally {
    isLoading.value = false
  }
})

onUnmounted(() => {
  chatService.disconnect()
})
</script>

<template>
  <div class="h-full flex flex-col bg-white rounded-tl-[30px] overflow-hidden shadow" data-theme="light">
    
    <!-- Chat Header -->
    <div class="px-6 py-4 border-b bg-white flex items-center justify-between shrink-0 relative overflow-hidden">
      <!-- Default Header Content -->
      <div v-if="!isSearching" class="flex items-center gap-3">
        <div class="w-10 h-10 rounded-full bg-gradient-to-br from-green-400 to-green-600 flex items-center justify-center text-white shadow">
          <UserRoundCog :size="20" />
        </div>
        <div>
          <h3 class="font-semibold text-gray-800">ฝ่ายช่วยเหลือ</h3>
          <p class="text-xs text-green-500 flex items-center gap-1">
            <span class="w-2 h-2 rounded-full bg-green-400 animate-pulse"></span>
            พร้อมให้บริการ
          </p>
        </div>
      </div>
      
      <div v-if="!isSearching" class="flex items-center gap-2">
        <button 
          @click="toggleSearch"
          class="p-2 rounded-full hover:bg-gray-100 text-gray-400 transition-colors"
        >
          <Search :size="20" />
        </button>
      </div>

      <!-- Search Bar Content -->
      <div v-if="isSearching" class="flex-1 flex items-center gap-2 animate-slide-down">
        <div class="relative flex-1">
          <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
          <input 
            v-model="messageSearchQuery"
            type="text" 
            placeholder="ค้นหาข้อความในแชท..." 
            class="w-full pl-9 pr-4 py-2 bg-gray-100 rounded-lg text-sm border-none focus:ring-2 focus:ring-blue-300 focus:bg-white transition-all"
            autoFocus
          />
        </div>
        <button 
          @click="toggleSearch"
          class="p-2 rounded-full hover:bg-gray-100 text-gray-500 transition-colors"
        >
          <X :size="20" />
        </button>
      </div>
    </div>
    
    <!-- Messages Area -->
    <div 
      ref="chatBox"
      class="flex-1 overflow-y-auto p-6 bg-gradient-to-b from-gray-50 to-white"
    >
      <!-- Loading State -->
      <div v-if="isLoading" class="flex items-center justify-center h-full">
        <div class="text-center text-gray-400">
          <div class="loading loading-spinner loading-lg mb-2"></div>
          <p>กำลังโหลด...</p>
        </div>
      </div>
      
      <!-- Empty State -->
      <div v-else-if="messages.length === 0" class="flex flex-col items-center justify-center h-full text-gray-400">
        <div class="w-24 h-24 rounded-full bg-blue-50 flex items-center justify-center mb-4">
          <MessageSquare :size="48" class="text-blue-300" />
        </div>
        <h3 class="font-semibold text-gray-600 mb-1">เริ่มการสนทนา</h3>
        <p class="text-sm">พิมพ์ข้อความเพื่อติดต่อฝ่ายช่วยเหลือ</p>
      </div>
      
      <!-- Messages -->
      <div v-else class="space-y-6">
        <template v-for="group in groupedMessages" :key="group.date">
          <!-- Date Divider -->
          <div class="flex items-center justify-center">
            <span class="px-3 py-1 bg-gray-200 text-gray-500 text-xs rounded-full">
              {{ formatDate(group.date) }}
            </span>
          </div>
          
          <!-- Messages in Group -->
          <div 
            v-for="msg in group.messages" 
            :key="msg.ID"
            class="flex"
            :class="msg.sender_id === myId ? 'justify-end' : 'justify-start'"
          >
            <!-- Other's message (left) -->
            <div v-if="msg.sender_id !== myId" class="flex items-end gap-2 max-w-[70%]">
              <div class="w-8 h-8 rounded-full bg-gradient-to-br from-green-400 to-green-600 flex items-center justify-center text-white shrink-0 shadow">
                <UserRoundCog :size="16" />
              </div>
              <div>
                <div class="bg-white px-4 py-3 rounded-2xl rounded-bl-md shadow-sm border">
                  <p class="text-gray-800 text-sm whitespace-pre-wrap">{{ msg.massage }}</p>
                </div>
                <span class="text-xs text-gray-400 ml-2 mt-1 block">{{ formatTime(msg.CreatedAt) }}</span>
              </div>
            </div>
            
            <!-- My message (right) -->
            <div v-else class="flex items-end gap-2 max-w-[70%] flex-row-reverse">
              <div class="w-8 h-8 rounded-full bg-gradient-to-br from-blue-400 to-blue-600 flex items-center justify-center text-white shrink-0 shadow">
                <UserRoundCheck :size="16" />
              </div>
              <div>
                <div class="bg-blue-500 px-4 py-3 rounded-2xl rounded-br-md text-white shadow-md">
                  <p class="text-sm whitespace-pre-wrap">{{ msg.massage }}</p>
                </div>
                <span class="text-xs text-gray-400 mr-2 mt-1 block text-right">{{ formatTime(msg.CreatedAt) }}</span>
              </div>
            </div>
          </div>
        </template>
      </div>
    </div>
    
    <!-- Input Area -->
    <div class="p-4 border-t bg-white shrink-0">
      <div class="flex items-center gap-3">
        <!-- Input -->
        <input
          v-model="input"
          @keyup.enter="send"
          type="text"
          placeholder="Type a message..."
          class="flex-1 px-4 py-3 bg-gray-100 rounded-xl border-none focus:ring-2 focus:ring-blue-300 focus:bg-white transition-all"
        />
        
        <!-- Send Button -->
        <button
          @click="send"
          :disabled="!input.trim()"
          class="px-5 py-3 bg-blue-500 hover:bg-blue-600 disabled:bg-gray-300 disabled:cursor-not-allowed text-white rounded-xl font-medium flex items-center gap-2 transition-colors shadow-md hover:shadow-lg"
        >
          Send
          <Send :size="20" />
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Custom scrollbar */
::-webkit-scrollbar {
  width: 6px;
}

::-webkit-scrollbar-track {
  background: transparent;
}

::-webkit-scrollbar-thumb {
  background: rgba(0, 0, 0, 0.1);
  border-radius: 10px;
}

::-webkit-scrollbar-thumb:hover {
  background: rgba(0, 0, 0, 0.2);
}

.animate-pulse {
  animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

.animate-slide-down {
  animation: slideDown 0.2s ease-out;
}

@keyframes slideDown {
  from { transform: translateY(-100%); opacity: 0; }
  to { transform: translateY(0); opacity: 1; }
}
</style>
