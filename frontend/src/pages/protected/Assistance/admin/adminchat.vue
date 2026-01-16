<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick, computed } from "vue"
import { ChatroomAPI } from "@/services/api/chatroom"
import { AssistanceAPI } from "@/services/api/assistance"
import { ChatService } from "@/services/api/websocket"
import type { ChatroomRespone, AssistanceRespone } from "@/interfaces"
import { getMyProfile } from "@/services/api/user"
import { UserRoundCheck, UserRoundCog, Search, Send, X } from "lucide-vue-next"

// ---------------- state ----------------
const rooms = ref<ChatroomRespone[]>([])
const selectedRoom = ref<ChatroomRespone | null>(null)
const messages = ref<AssistanceRespone[]>([])
const allMessages = ref<AssistanceRespone[]>([])
const input = ref("")
const searchQuery = ref("")
const myId = ref<number>(0)
const isLoading = ref(true)
const isLoadingMessages = ref(false)
const isSearchingMessages = ref(false)
const messageSearchQuery = ref("")

const chatService = new ChatService()
const token = sessionStorage.getItem("token")
const chatBox = ref<HTMLElement | null>(null)
const isConnected = ref(false)
const isConnecting = ref(false)
let connectionCheckInterval: number | null = null

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

// Filter rooms by search
const filteredRooms = computed(() => {
  if (!searchQuery.value.trim()) return rooms.value
  const query = searchQuery.value.toLowerCase()
  return rooms.value.filter(r => {
    const username = r.user?.username?.toLowerCase() || ""
    const studentId = r.user?.student_profile?.[0]?.student_id?.toLowerCase() || ""
    return username.includes(query) || studentId.includes(query)
  })
})

// Helper to get display name (Student ID or Username)
const getDisplayName = (user: any) => {
  return user?.student_profile?.[0]?.student_id || user?.username || "Unknown"
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

// Get last message for a room
const getLastMessage = (roomId: number) => {
  const roomMessages = allMessages.value.filter(m => m.chatroom_id === roomId)
  if (roomMessages.length === 0) return "เริ่มการสนทนา..."
  const lastMsg = roomMessages[roomMessages.length - 1]
  const msg = lastMsg?.massage || ""
  return msg.substring(0, 30) + (msg.length > 30 ? "..." : "")
}

const getLastMessageTime = (roomId: number) => {
  const roomMessages = allMessages.value.filter(m => m.chatroom_id === roomId)
  if (roomMessages.length === 0) return ""
  return formatTime(roomMessages[roomMessages.length - 1]?.CreatedAt || "")
}

// ---------------- methods ----------------
const scrollToBottom = () => {
  nextTick(() => {
    if (chatBox.value) {
      chatBox.value.scrollTop = chatBox.value.scrollHeight
      // Double check scroll after a short delay to ensure images/layout are stable
      setTimeout(() => {
        if (chatBox.value) chatBox.value.scrollTop = chatBox.value.scrollHeight
      }, 100)
    }
  })
}

const toggleMessageSearch = () => {
  isSearchingMessages.value = !isSearchingMessages.value
  if (!isSearchingMessages.value) {
    messageSearchQuery.value = ""
  }
}

const selectRoom = async (room: ChatroomRespone) => {
  // If clicking same room, do nothing
  if (selectedRoom.value?.ID === room.ID) return;

  selectedRoom.value = room
  messages.value = []
  isLoadingMessages.value = true
  isSearchingMessages.value = false
  messageSearchQuery.value = ""
  isConnecting.value = true

  chatService.disconnect()

  // Wait a bit to ensure socket allows new connection (debounce)
  setTimeout(async () => {
    try {
        // Filter messages for this room from cache
        messages.value = allMessages.value.filter(m => m.chatroom_id === room.ID)
        scrollToBottom()

        // connect websocket
        if (token) {
            chatService.connect(token, room.ID)
        }
        
        // Listen for messages
        chatService.onMessage(msg => {
            messages.value.push(msg)
            allMessages.value.push(msg)
            scrollToBottom()
        })
    } catch (err) {
        console.error("Error loading room:", err)
    } finally {
        isLoadingMessages.value = false
        isConnecting.value = false
    }
  }, 300);
}

const leaveRoom = () => {
  chatService.disconnect();
  selectedRoom.value = null;
  messages.value = [];
}

const send = () => {
  if (!input.value.trim()) return
  
  if (!chatService.isConnected()) {
    // Try to reconnect if room is selected
     if (selectedRoom.value && token) {
        console.log("Attempting to reconnect...");
        chatService.connect(token, selectedRoom.value.ID);
        // Wait a bit? Or just alert user to try again
     }
     alert('ไม่สามารถส่งข้อความได้ (Socket Disconnected) กรุณาลองกดส่งใหม่อีกครั้ง');
     return;
  }
  
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
  isLoading.value = true
  try {
    const profileRes = await getMyProfile()
    myId.value = profileRes.data.ID

    // Load all messages first
    const historyRes = await AssistanceAPI.getAll()
    allMessages.value = Array.isArray(historyRes)
      ? historyRes
      : Array.isArray(historyRes?.data)
        ? historyRes.data
        : []

    const res = await ChatroomAPI.getAllOpen()
    rooms.value = Array.isArray(res.data) ? res.data : []

    // Removed auto-select first room
    // if (rooms.value.length > 0 && !selectedRoom.value) {
    //   selectRoom(rooms.value[0]!)
    // }
  } catch (err) {
    console.error("Failed to load rooms:", err)
  } finally {
    isLoading.value = false
  }
}

// ---------------- lifecycle ----------------
// ---------------- lifecycle ----------------
onMounted(() => {
    loadRooms();
    // Start connection checker
    connectionCheckInterval = window.setInterval(() => {
        isConnected.value = chatService.isConnected();
    }, 2000);
})

onUnmounted(() => {
    chatService.disconnect();
    if (connectionCheckInterval) clearInterval(connectionCheckInterval);
})
</script>

<template>
  <div class="h-full flex bg-gray-100 rounded-tl-[30px] overflow-hidden" data-theme="light">
    
    <!-- Left Sidebar: Chat List -->
    <div class="w-80 bg-white border-r flex flex-col shrink-0">
      <!-- Header -->
      <div class="p-4 border-b">
        <h2 class="text-2xl font-bold text-slate-800 flex items-center gap-2">
          ห้องแชททั้งหมด
        </h2>
        <p class="text-gray-500 text-xs mt-1 ">{{ rooms.length }} ห้องที่เปิดอยู่</p>
      </div>
      
      <!-- Search Bar -->
      <div class="p-3 border-b">
        <div class="relative">
          <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
          <input 
            v-model="searchQuery"
            type="text" 
            placeholder="ค้นหาผู้ใช้ / รหัสนิสิต..." 
            class="w-full pl-9 pr-4 py-2 bg-gray-100 rounded-lg text-sm border-none focus:ring-2 focus:ring-blue-300 focus:bg-white transition-all"
          />
        </div>
      </div>
      
      <!-- Chat List -->
      <div class="flex-1 overflow-y-auto">
        <!-- Loading -->
        <div v-if="isLoading" class="flex items-center justify-center h-32 text-gray-400">
          <div class="loading loading-spinner loading-md"></div>
        </div>
        
        <!-- Empty -->
        <div v-else-if="filteredRooms.length === 0" class="flex flex-col items-center justify-center h-32 text-gray-400 p-4">
          <svg xmlns="http://www.w3.org/2000/svg" class="w-12 h-12 mb-2 opacity-40" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
          </svg>
          <p class="text-sm">ไม่มีห้องแชท</p>
        </div>
        
        <!-- Room List -->
        <div 
          v-else
          v-for="r in filteredRooms"
          :key="r.ID"
          @click="selectRoom(r)"
          class="flex items-center gap-3 p-4 cursor-pointer hover:bg-gray-50 transition-colors border-b border-gray-100"
          :class="r.ID === selectedRoom?.ID ? 'bg-blue-50 border-l-4 border-l-blue-500' : ''"
        >
          <!-- Avatar -->
          <div class="w-12 h-12 rounded-full bg-gradient-to-br from-blue-400 to-blue-600 flex items-center justify-center text-white shrink-0 shadow relative">
            <UserRoundCheck :size="24" />
            <!-- Online indicator -->
            <span 
              v-if="r.status_chatroom === 'open'"
              class="absolute bottom-0 right-0 w-3.5 h-3.5 bg-green-400 border-2 border-white rounded-full"
            ></span>
          </div>
          
          <!-- Info -->
          <div class="flex-1 min-w-0">
            <div class="flex items-center justify-between mb-0.5">
              <h4 class="font-semibold text-gray-800 truncate text-sm">
                {{ getDisplayName(r.user) }}
              </h4>
              <span class="text-xs text-gray-400 shrink-0">
                {{ getLastMessageTime(r.ID) }}
              </span>
            </div>
            <p class="text-xs text-gray-500 truncate">{{ getLastMessage(r.ID) }}</p>
          </div>
        </div>
      </div>
    </div>
    
    <!-- Right Side: Chat Window -->
    <div class="flex-1 flex flex-col bg-white min-w-0">
      
      <!-- Selected Room -->
      <template v-if="selectedRoom">
        <!-- Chat Header -->
        <div class="px-6 py-4 border-b bg-white flex items-center justify-between shrink-0 relative overflow-hidden h-[72px]">
          <!-- Default Header Content -->
          <div v-if="!isSearchingMessages" class="flex items-center gap-3">
            <div class="w-10 h-10 rounded-full bg-gradient-to-br from-blue-400 to-blue-600 flex items-center justify-center text-white shadow relative">
              <UserRoundCheck :size="20" />
              <!-- Status Dot on Avatar -->
              <span 
                :class="['absolute bottom-0 right-0 w-3 h-3 border-2 border-white rounded-full', isConnected ? 'bg-green-500' : 'bg-red-500']"
                :title="isConnected ? 'Online' : 'Offline'"
              ></span>
            </div>
            <div>
              <h3 class="font-semibold text-gray-800 flex items-center gap-2">
                {{ getDisplayName(selectedRoom.user) }}
                <!-- Status Badge Text -->
                <span v-if="!isConnected" class="text-[10px] px-1.5 py-0.5 bg-red-100 text-red-600 rounded-md font-medium">Offline</span>
              </h3>
              <p class="text-xs text-green-600 flex items-center gap-1" v-if="isConnected">
                กำลังสนทนา
              </p>
              <p class="text-xs text-slate-400" v-else>
                ไม่ได้เชื่อมต่อ
              </p>
            </div>
          </div>
          
          <div v-if="!isSearchingMessages" class="flex items-center gap-3">
            <!-- Search -->
            <button 
              @click="toggleMessageSearch"
              class="w-9 h-9 flex items-center justify-center rounded-full hover:bg-gray-100 text-gray-500 transition-colors tooltip"
              title="ค้นหาข้อความ"
            >
              <Search :size="18" />
            </button>

            <div class="h-6 w-px bg-gray-200 mx-1"></div>

            <!-- End Chat / Close Room -->
            <button 
              @click="closeRoom"
              class="px-4 py-2 text-xs rounded-lg bg-red-50 text-red-600 hover:bg-red-100 border border-red-200 font-medium transition-colors shadow-sm flex items-center gap-2"
              title="จบการสนทนา (ปิดเคส)"
            >
              <span>จบการสนทนา</span>
            </button>

            <!-- Leave / Minimize -->
             <button 
              @click="leaveRoom"
              class="w-9 h-9 flex items-center justify-center rounded-full hover:bg-gray-100 text-gray-400 hover:text-gray-600 transition-colors"
              title="พักแชท (ออกหน้านี้)"
            >
              <X :size="20" />
            </button>
          </div>

          <!-- Search Bar Content -->
          <div v-if="isSearchingMessages" class="flex-1 flex items-center gap-2 animate-slide-down">
            <div class="relative flex-1">
              <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
              <input 
                v-model="messageSearchQuery"
                type="text" 
                placeholder="ค้นหาข้อความในแชทนี้..." 
                class="w-full pl-9 pr-4 py-2 bg-gray-100 rounded-lg text-sm border-none focus:ring-2 focus:ring-blue-300 focus:bg-white transition-all"
                autoFocus
              />
            </div>
            <button 
              @click="toggleMessageSearch"
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
          <!-- Loading -->
          <div v-if="isLoadingMessages" class="flex items-center justify-center h-full">
            <div class="text-center text-gray-400">
              <div class="loading loading-spinner loading-lg mb-2"></div>
              <p>กำลังโหลดข้อความ...</p>
            </div>
          </div>
          
          <!-- Empty -->
          <div v-else-if="messages.length === 0" class="flex flex-col items-center justify-center h-full text-gray-400">
            <div class="w-20 h-20 rounded-full bg-gray-100 flex items-center justify-center mb-4">
              <svg xmlns="http://www.w3.org/2000/svg" class="w-10 h-10 text-gray-300" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
              </svg>
            </div>
            <p class="text-sm">ยังไม่มีข้อความในห้องนี้</p>
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
                <!-- User's message (left) -->
                <div v-if="msg.sender_id !== myId" class="flex items-end gap-2 max-w-[70%]">
                  <div class="w-8 h-8 rounded-full bg-gradient-to-br from-blue-400 to-blue-600 flex items-center justify-center text-white shrink-0 shadow">
                    <UserRoundCheck :size="16" />
                  </div>
                  <div>
                    <div class="bg-white px-4 py-3 rounded-2xl rounded-bl-md shadow-sm border">
                      <p class="text-gray-800 text-sm whitespace-pre-wrap">{{ msg.massage }}</p>
                    </div>
                    <span class="text-xs text-gray-400 ml-2 mt-1 block">{{ formatTime(msg.CreatedAt) }}</span>
                  </div>
                </div>
                
                <!-- Admin's message (right) -->
                <div v-else class="flex items-end gap-2 max-w-[70%] flex-row-reverse">
                  <div class="w-8 h-8 rounded-full bg-gradient-to-br from-green-400 to-green-600 flex items-center justify-center text-white shrink-0 shadow">
                    <UserRoundCog :size="16" />
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
            <input
              v-model="input"
              @keyup.enter="send"
              type="text"
              placeholder="Type a message..."
              class="flex-1 px-4 py-3 bg-gray-100 rounded-xl border-none focus:ring-2 focus:ring-blue-300 focus:bg-white transition-all"
            />
            
            <button
              @click="send"
              :disabled="!input.trim() || isConnecting"
              class="px-5 py-3 bg-blue-500 hover:bg-blue-600 disabled:bg-gray-300 disabled:cursor-not-allowed text-white rounded-xl font-medium flex items-center gap-2 transition-colors shadow-md hover:shadow-lg"
            >
              <span v-if="isConnecting" class="loading loading-spinner loading-xs"></span>
              <span v-else>Send</span>
              <Send :size="18" />
            </button>
          </div>
        </div>
      </template>
      
      <!-- No Room Selected -->
      <div v-else class="flex-1 flex flex-col items-center justify-center text-gray-400 bg-gradient-to-b from-gray-50 to-white">
        <div class="w-32 h-32 rounded-full bg-gray-100 flex items-center justify-center mb-6">
          <svg xmlns="http://www.w3.org/2000/svg" class="w-16 h-16 text-gray-300" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
          </svg>
        </div>
        <h3 class="text-lg font-medium text-gray-600 mb-2">เลือกห้องแชท</h3>
        <p class="text-sm">คลิกที่ห้องแชททางซ้ายเพื่อเริ่มการสนทนา</p>
      </div>
      
    </div>
  </div>
</template>

<style scoped>
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
