<template>
  <div class="min-h-screen bg-base-200 flex items-center justify-center">
    <div class="card w-full max-w-sm bg-base-100 shadow-xl">
      <div class="card-body">
        <h2 class="card-title text-2xl justify-center">Login</h2>
        <form @submit.prevent="handleLogin">
          <div class="form-control">
            <label class="label">
              <span class="label-text">Username</span>
            </label>
            <input type="text" v-model="user.username" required class="input input-bordered" />
          </div>
          <div class="form-control mt-4">
            <label class="label">
              <span class="label-text">Password</span>
            </label>
            <input type="password" v-model="user.password" required class="input input-bordered" />
          </div>
          <div class="form-control mt-6">
            <button type="submit" class="btn btn-primary" :disabled="loading">
              <span v-if="loading" class="loading loading-spinner"></span>
              <span v-else>Login</span>
            </button>
          </div>
          <div v-if="message" role="alert" class="alert alert-error mt-4">
            <svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
            <span>{{ message }}</span>
          </div>
          <p class="text-center mt-4">
            Don't have an account? <router-link to="/register" class="link link-primary">Register here</router-link>
          </p>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import AuthService from '../services/auth.service';

const user = ref({ username: '', password: '' });
const loading = ref(false);
const message = ref('');
const router = useRouter();

const handleLogin = () => {
  loading.value = true;
  message.value = '';
  AuthService.login(user.value).then(
    () => {
      router.push('/dashboard');
    },
    (error) => {
      loading.value = false;
      message.value =
        (error.response &&
          error.response.data &&
          error.response.data.error) ||
        error.message ||
        error.toString();
    }
  );
};
</script>