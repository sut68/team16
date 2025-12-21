<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { jwtDecode } from 'jwt-decode';
import { Eye, EyeOff } from 'lucide-vue-next';

import { authAPI } from '../../services/api';
import type { LoginUserRequest } from '../../interfaces/user';

import SuryaGraphic from '@/assets/brand/Surya_graphic.png';
import SUTLogo from '@/assets/logo/SUT_logo_orange.png';

const user = ref<LoginUserRequest>({
  username: '',
  password: '',
});

const loading = ref(false);
const message = ref('');
const showPassword = ref(false);

const router = useRouter();

const handleLogin = async () => {
  loading.value = true;
  message.value = '';
  try {
    const data = await authAPI.login(user.value);

    if (data?.token) {
      sessionStorage.setItem('token', data.token);
      if (data.token_type) {
        sessionStorage.setItem('token_type', data.token_type);
      }

      const decodedToken: { role: string } = jwtDecode(data.token);

      if (decodedToken.role === 'admin') {
        await router.push('/admin');
      } else if (['user', 'student'].includes(decodedToken.role)) {
        await router.push('/dashboard');
      } else {
        await router.push('/');
      }

      window.location.reload();
      return;
    }

    message.value = data?.error ?? data?.message ?? 'ข้อมูลไม่ถูกต้อง';
  } catch {
    message.value = 'เกิดข้อผิดพลาดในการเชื่อมต่อ';
  } finally {
    loading.value = false;
  }
};
</script>

<template>
  <div
    class="h-screen bg-[#8B001D] flex items-center justify-center
           p-4 md:p-8 overflow-hidden relative"
  >
    <!-- background graphic -->
    <div
      class="absolute inset-0 z-0 opacity-30 bg-no-repeat"
      :style="{
        backgroundImage: `url(${SuryaGraphic})`,
        backgroundPosition: 'right center',
        backgroundSize: '150%',
      }"
    ></div>

    <!-- glow -->
    <div
      class="absolute top-0 right-0 w-96 h-96 bg-[#F26522]/20
             rounded-full blur-[120px] translate-x-1/2 -translate-y-1/2"
    ></div>
    <div
      class="absolute bottom-0 left-0 w-64 h-64 bg-[#F26522]/10
             rounded-full blur-[100px] -translate-x-1/2 translate-y-1/2"
    ></div>

    <!-- CARD -->
    <div
      class="bg-white w-full max-w-6xl h-full max-h-[680px]
             flex flex-col lg:flex-row rounded-[2.5rem]
             shadow-[0_50px_100px_-20px_rgba(0,0,0,0.5)]
             overflow-hidden border border-white/10
             z-10 animate__animated animate__zoomIn"
    >
      <!-- LEFT -->
      <div
        class="w-full lg:w-[45%] p-8 md:p-12
               flex flex-col justify-center bg-white
               relative overflow-y-auto"
      >
        <!-- logo -->
        <div class="flex items-center gap-3 mb-6">
          <img
            :src="SUTLogo"
            alt="SUT Logo"
            class="w-16 h-auto drop-shadow-md object-contain"
          />
          <div class="flex flex-col">
            <span class="font-black text-[#253C90] text-lg uppercase">
              Scholarship
            </span>
            <span
              class="text-[9px] text-gray-400 font-bold tracking-widest uppercase"
            >
              University Portal
            </span>
          </div>
        </div>

        <!-- title -->
        <h1
          class="text-3xl md:text-4xl font-black text-gray-800
                 mb-4 leading-tight"
        >
          Start your <br />
          <span
            class="text-transparent bg-clip-text
                   bg-gradient-to-r from-[#F26522] to-[#8B001D]"
          >
            success story
          </span>
        </h1>

        <!-- form -->
        <form @submit.prevent="handleLogin" class="space-y-6">
          <!-- username -->
          <div class="form-control">
            <label
              class="text-[10px] font-bold text-[#253C90]
                     uppercase tracking-widest mb-1"
            >
              Username / ID
            </label>
            <input
              type="text"
              v-model="user.username"
              placeholder="username"
              required
              class="w-full py-3 bg-transparent border-b-2
                     border-gray-100 outline-none
                     focus:border-[#F26522] transition-all
                     text-gray-700 placeholder:text-gray-200"
            />
          </div>

          <!-- password + eye -->
          <div class="form-control">
            <label
              class="text-[10px] font-bold text-[#253C90]
                     uppercase tracking-widest mb-1"
            >
              Password
            </label>

            <div class="relative">
              <input
                :type="showPassword ? 'text' : 'password'"
                v-model="user.password"
                placeholder="password"
                required
                class="w-full py-3 bg-transparent border-b-2
                       border-gray-100 outline-none
                       focus:border-[#F26522] transition-all
                       text-gray-700 placeholder:text-gray-200 pr-10"
              />

              <button
                type="button"
                @click="showPassword = !showPassword"
                class="absolute right-1 top-1/2 -translate-y-1/2
                       text-gray-400 hover:text-[#F26522]
                       transition-colors"
              >
                <Eye v-if="!showPassword" class="w-5 h-5" />
                <EyeOff v-else class="w-5 h-5" />
              </button>
            </div>
          </div>

          <!-- submit -->
          <button
            type="submit"
            :disabled="loading"
            class="w-full lg:w-auto px-16 py-4
                   bg-[#F26522] text-white font-black
                   rounded-full transition-all
                   hover:bg-[#8B001D] hover:shadow-2xl
                   active:scale-95 disabled:opacity-50"
          >
            <span v-if="loading" class="loading loading-spinner loading-sm"></span>
            <span v-else class="text-sm tracking-widest uppercase">
              Sign In
            </span>
          </button>
        </form>

        <!-- error -->
        <div
          v-if="message"
          class="mt-4 p-3 bg-red-50 text-[#8B001D]
                 rounded-xl text-xs border-l-4 border-[#8B001D]"
        >
          {{ message }}
        </div>
      </div>

      <!-- RIGHT -->
      <div
        class="hidden lg:flex lg:w-[55%]
               bg-gradient-to-br from-[#8B001D] to-[#6d0016]
               relative items-center justify-center
               p-12 overflow-hidden border-l border-gray-50"
      >
        <div
          class="relative z-10 w-full max-w-sm
                 bg-white/5 backdrop-blur-2xl
                 border border-white/10
                 p-12 rounded-[3rem]
                 shadow-2xl text-center"
        >
          <div class="floating-element mb-8">
            <img
              src="https://illustrations.popsy.co/white/studying.svg"
              alt="Scholarship"
              class="w-64 h-auto mx-auto
                     drop-shadow-[0_35px_35px_rgba(0,0,0,0.4)]"
            />
          </div>

          <p class="text-white/60 font-light text-sm leading-relaxed">
            "ทุนเรียนดี SUT ตัวจริง <br />
            รับทุนสูงสุด 100%*"
          </p>
        </div>

        <div
          class="absolute bottom-[-10%] right-[-10%]
                 w-64 h-64 bg-[#F26522]/20
                 blur-[80px] rounded-full"
        ></div>
      </div>
    </div>
  </div>
</template>

<style scoped>
@import "https://cdnjs.cloudflare.com/ajax/libs/animate.css/4.1.1/animate.min.css";

.floating-element {
  animation: float 6s ease-in-out infinite;
}

@keyframes float {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-20px); }
}

.overflow-y-auto::-webkit-scrollbar {
  width: 0;
}
</style>
