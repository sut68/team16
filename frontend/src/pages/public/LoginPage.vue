<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { jwtDecode } from 'jwt-decode';
import { Eye, EyeOff, Sparkles, GraduationCap, ArrowRight } from 'lucide-vue-next';

import { authAPI } from '../../services/api';
import type { LoginUserRequest } from '../../interfaces/user';
import SUTLogo from '@/assets/logo/SUT_logo_orange.png'; // Make sure this path is correct

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
  
  // Fake delay for animation effect
  await new Promise(resolve => setTimeout(resolve, 800));

  try {
    const data = await authAPI.login(user.value);

    if (data?.token) {
      sessionStorage.setItem('token', data.token);
      if (data.token_type) {
        sessionStorage.setItem('token_type', data.token_type);
      }

      const decodedToken: { role: string } = jwtDecode(data.token);

      if (decodedToken.role === 'admin') {
        await router.push('/admin/news');
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
  <div class="min-h-screen relative w-full overflow-hidden flex items-center justify-center bg-[#8B001D]">
    
    <!-- Animated Background -->
    <div class="absolute inset-0 z-0">
      <div class="absolute w-[200%] h-[200%] top-[-50%] left-[-50%] animate-gradient-spin">
        <div class="w-full h-full bg-gradient-to-br from-[#5e0013] via-[#8B001D] to-[#F26522] opacity-80 blur-3xl"></div>
      </div>
      <div class="absolute inset-0 bg-[url('https://www.transparenttextures.com/patterns/cubes.png')] opacity-10 mix-blend-overlay"></div>
    </div>

    <!-- Floating Orbs -->
    <div class="absolute top-20 left-20 w-48 h-48 bg-[#FFD700] rounded-full blur-[80px] animate-float-slow opacity-40"></div>
    <div class="absolute bottom-20 right-20 w-64 h-64 bg-[#F26522] rounded-full blur-[100px] animate-float-delayed opacity-50"></div>

    <!-- Main Card Container -->
    <div class="relative z-10 w-full max-w-5xl h-[85vh] max-h-[700px] flex rounded-[3rem] overflow-hidden shadow-2xl animate-card-entrance backdrop-blur-xl border border-white/20 bg-white/10">
      
      <!-- LEFT SIDE: Form -->
      <div class="w-full lg:w-[45%] p-8 md:p-12 flex flex-col justify-center relative bg-white/10 backdrop-blur-md border-r border-white/10">
        
        <!-- Logo -->
        <div class="flex items-center gap-4 mb-4 animate-fade-in-down" style="animation-delay: 0.2s;">
          <div class="w-16 h-16 bg-white rounded-2xl flex items-center justify-center shadow-lg transform rotate-3 hover:rotate-0 transition-all duration-300">
             <img :src="SUTLogo" alt="SUT" class="w-12 h-auto" />
          </div>
          <div>
            <h2 class="text-white font-black text-2xl tracking-tighter leading-none shadow-black drop-shadow-md">SUT</h2>
            <span class="text-[#FFD700] text-xs font-bold tracking-[0.2em] uppercase drop-shadow">Scholarships</span>
          </div>
        </div>

        <!-- Title -->
        <div class="mb-10 animate-fade-in-up" style="animation-delay: 0.4s;">
          <h1 class="text-4xl md:text-5xl font-black text-white leading-tight drop-shadow-md">
            Unlock Your <br/>
            <span class="text-transparent bg-clip-text bg-gradient-to-r from-[#FFD700] to-[#fff] animate-pulse-slow">Potential</span>
          </h1>
          <p class="text-white/80 mt-4 font-light text-sm">เข้าสู่ระบบเพื่ออนาคตทางการศึกษาที่เหนือกว่า</p>
        </div>

        <!-- Form -->
        <form @submit.prevent="handleLogin" class="space-y-6 animate-fade-in-up" style="animation-delay: 0.6s;">
          
          <!-- Username Input -->
          <div class="group">
             <div class="relative transition-all duration-300 transform group-hover:-translate-y-1">
               <input 
                 type="text" 
                 v-model="user.username" 
                 required
                 class="peer w-full bg-white/20 text-white rounded-2xl px-6 py-4 outline-none border-2 border-transparent focus:border-[#FFD700] focus:bg-white/30 transition-all duration-300 placeholder-transparent shadow-inner"
                 placeholder="Username"
                 id="username"
               />
               <label for="username" class="absolute left-6 top-4 text-white/70 text-sm transition-all duration-300 peer-placeholder-shown:text-base peer-placeholder-shown:top-4 peer-focus:top-[-10px] peer-focus:text-xs peer-focus:text-[#FFD700] peer-focus:bg-[#8B001D] peer-focus:px-2 peer-focus:rounded peer-valid:top-[-10px] peer-valid:text-xs peer-valid:text-[#FFD700] peer-valid:bg-[#8B001D] peer-valid:px-2 peer-valid:rounded whitespace-nowrap">
                  <GraduationCap class="inline w-4 h-4 mr-1"/> รหัสนักศึกษา / Username
               </label>
             </div>
          </div>

          <!-- Password Input -->
          <div class="group">
             <div class="relative transition-all duration-300 transform group-hover:-translate-y-1">
               <input 
                 :type="showPassword ? 'text' : 'password'"
                 v-model="user.password" 
                 required
                 class="peer w-full bg-white/20 text-white rounded-2xl px-6 py-4 outline-none border-2 border-transparent focus:border-[#FFD700] focus:bg-white/30 transition-all duration-300 placeholder-transparent shadow-inner"
                 placeholder="Password"
                 id="password"
               />
               <label for="password" class="absolute left-6 top-4 text-white/70 text-sm transition-all duration-300 peer-placeholder-shown:text-base peer-placeholder-shown:top-4 peer-focus:top-[-10px] peer-focus:text-xs peer-focus:text-[#FFD700] peer-focus:bg-[#8B001D] peer-focus:px-2 peer-focus:rounded peer-valid:top-[-10px] peer-valid:text-xs peer-valid:text-[#FFD700] peer-valid:bg-[#8B001D] peer-valid:px-2 peer-valid:rounded">
                 Password
               </label>
               <button type="button" @click="showPassword = !showPassword" class="absolute right-4 top-1/2 -translate-y-1/2 text-white/50 hover:text-[#FFD700] transition-colors">
                  <Eye v-if="!showPassword" class="w-5 h-5"/>
                  <EyeOff v-else class="w-5 h-5"/>
               </button>
             </div>
          </div>

          <!-- Submit Button -->
          <button 
            type="submit" 
            :disabled="loading"
            class="px-8 py-4 w-full bg-gradient-to-r from-[#FFD700] to-[#F26522] text-[#8B001D] font-black rounded-2xl shadow-lg transform transition-all duration-300 hover:scale-105 hover:shadow-[#F26522]/50 active:scale-95 disabled:opacity-70 disabled:cursor-not-allowed group overflow-hidden relative"
          >
            <div class="absolute inset-0 bg-white/40 translate-x-[-100%] group-hover:translate-x-[100%] transition-transform duration-700 skew-x-12"></div>
            <div class="flex items-center justify-center gap-2 relative z-10">
              <span v-if="loading" class="animate-spin w-5 h-5 border-2 border-[#8B001D] border-t-transparent rounded-full"></span>
              <span v-else>SIGN IN</span>
              <ArrowRight v-if="!loading" class="w-5 h-5 group-hover:translate-x-1 transition-transform"/>
              
            </div>
          </button>

          <!-- Error Message -->
          <div v-if="message" class="mt-4 p-4 rounded-xl bg-white/90 border border-red-500 text-red-600 text-sm flex items-center gap-2 animate-shake font-bold shadow-md">
            <span class="w-2 h-2 rounded-full bg-red-600"></span>
            {{ message }}
          </div>
        </form>

        <div class="mt-8 text-center text-white/40 text-xs font-light animate-fade-in" style="animation-delay: 1s;">
          &copy; 2024 Suranaree University of Technology
        </div>
      </div>

      <!-- RIGHT SIDE: Visuals -->
      <div class="hidden lg:flex lg:w-[55%] relative items-center justify-center overflow-hidden bg-gradient-to-br from-[#F26522]/20 to-[#8B001D]/20 backdrop-blur-sm">
        
        <!-- Animated Elements -->
        <div class="absolute inset-0 z-0">
           <div class="absolute top-[20%] left-[20%] w-32 h-32 bg-white opacity-10 blur-2xl animate-pulse"></div>
           <div class="absolute bottom-[20%] right-[20%] w-64 h-64 bg-[#F26522] opacity-20 blur-3xl animate-pulse" style="animation-delay: 1s;"></div>
        </div>

        <div class="relative z-10 text-center animate-float-slow">
           <div class="w-80 h-80 bg-gradient-to-br from-white/20 to-transparent backdrop-blur-xl rounded-[4rem] border border-white/30 shadow-2xl flex items-center justify-center transform rotate-6 hover:rotate-0 transition-all duration-700 group">
              <div class="absolute inset-0 bg-white opacity-0 group-hover:opacity-20 transition-opacity duration-500 rounded-[4rem]"></div>
              
              <div class="text-center p-8">
                <div class="w-20 h-20 bg-white rounded-2xl mx-auto mb-6 flex items-center justify-center shadow-lg transform group-hover:scale-110 transition-transform duration-500">
                  <Sparkles class="w-10 h-10 text-[#F26522] animate-spin-slow" />
                </div>
                <h3 class="text-2xl font-bold text-white mb-2 drop-shadow-md">Be The Future</h3>
                <p class="text-white/80 text-sm leading-relaxed font-medium">
                  "โอกาสทางการศึกษา <br/> ที่ไร้ขีดจำกัดรอคุณอยู่"
                </p>
              </div>
           </div>
           
           <!-- Decorative cards behind -->
           <div class="absolute top-10 -right-10 w-40 h-40 bg-[#F26522]/30 backdrop-blur-md rounded-3xl border border-white/10 -z-10 transform rotate-12 animate-float-delayed"></div>
           <div class="absolute -bottom-10 -left-10 w-32 h-32 bg-[#8B001D]/40 backdrop-blur-md rounded-full -z-10 animate-pulse-slow"></div>
        </div>
      </div>

    </div>
  </div>
</template>

<style scoped>
/* Keyframe Animations */
@keyframes gradient-spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

@keyframes float-slow {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-20px); }
}

@keyframes float-delayed {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(15px); }
}

@keyframes card-entrance {
  0% { opacity: 0; transform: scale(0.9) translateY(20px); }
  100% { opacity: 1; transform: scale(1) translateY(0); }
}

@keyframes fade-in-down {
  0% { opacity: 0; transform: translateY(-20px); }
  100% { opacity: 1; transform: translateY(0); }
}

@keyframes fade-in-up {
  0% { opacity: 0; transform: translateY(20px); }
  100% { opacity: 1; transform: translateY(0); }
}

@keyframes pulse-slow {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.6; }
}

@keyframes shake {
  0%, 100% { transform: translateX(0); }
  25% { transform: translateX(-5px); }
  75% { transform: translateX(5px); }
}

@keyframes spin-slow {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

/* Utility Classes for Animation using Tailwind arbitrary values */
.animate-gradient-spin {
  animation: gradient-spin 20s linear infinite;
}

.animate-float-slow {
  animation: float-slow 6s ease-in-out infinite;
}

.animate-float-delayed {
  animation: float-delayed 7s ease-in-out infinite 1s;
}

.animate-card-entrance {
  animation: card-entrance 1s cubic-bezier(0.2, 0.8, 0.2, 1) forwards;
}

.animate-fade-in-down {
  animation: fade-in-down 0.8s ease-out forwards;
  opacity: 0; /* Star hidden */
}

.animate-fade-in-up {
  animation: fade-in-up 0.8s ease-out forwards;
  opacity: 0; /* Start hidden */
}

.animate-pulse-slow {
  animation: pulse-slow 4s ease-in-out infinite;
}

.animate-shake {
  animation: shake 0.4s ease-in-out;
}

.animate-spin-slow {
  animation: spin-slow 10s linear infinite;
}
</style>
