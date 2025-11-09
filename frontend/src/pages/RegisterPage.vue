<template>
  <div class="min-h-screen bg-base-200 flex items-center justify-center">
    <div class="card w-full max-w-sm bg-base-100 shadow-xl">
      <div class="card-body">
        <h2 class="card-title text-2xl justify-center">Register</h2>
        <form @submit.prevent="handleRegister">
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
              <span v-else>Register</span>
            </button>
          </div>

          <div v-if="message" role="alert" :class="['alert mt-4', successful ? 'alert-success' : 'alert-error']">
            <svg v-if="successful" xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
            <svg v-else xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
            <span>{{ message }}</span>
          </div>

          <p v-if="successful" class="text-center mt-4">
            You can now <router-link to="/login" class="link link-primary">login</router-link>.
          </p>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import AuthService from '../services/auth.service';

const user = ref({ username: '', password: '' });
const loading = ref(false);
const message = ref('');
const successful = ref(false);

const handleRegister = () => {
  loading.value = true;
  message.value = '';
  successful.value = false;

  AuthService.register(user.value).then(
    (response) => {
      loading.value = false;
      successful.value = true;
      message.value = response.data.message;
    },
    (error) => {
      loading.value = false;
      successful.value = false;
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