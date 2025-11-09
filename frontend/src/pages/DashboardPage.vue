<template>
  <div class="min-h-screen bg-base-200 flex items-center justify-center">
    <div class="text-center">
      <div class="hero-content">
        <div class="max-w-md">
          <h1 class="text-5xl font-bold">Welcome!</h1>
          <p class="py-6">You are successfully logged in to your dashboard.</p>
          <button @click="handleLogout" class="btn btn-error">Logout</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router';
import AuthService from '../services/auth.service';

const router = useRouter();

const handleLogout = () => {
  AuthService.logout();
  router.push('/login');
};
</script>