<script setup>import { ref, inject } from 'vue';
import { useRouter } from 'vue-router';
import { authApi } from '../api';
const router = useRouter();
const setUser = inject('setUser');
const form = ref({
 username: '',
 password: ''
});
const isLoading = ref(false);
const error = ref('');
const handleLogin = async () => {
 if (!form.value.username || !form.value.password) {
 error.value = '请填写用户名和密码';
 return;
 }
 isLoading.value = true;
 try {
 const response = await authApi.login(form.value);
 localStorage.setItem('token', response.data.token);
 setUser(response.data.user);
 router.push('/');
 }
 catch (err) {
 error.value = err.response?.data?.error || '登录失败';
 }
 finally {
 isLoading.value = false;
 }
};
</script>

<template>
  <div class="min-h-screen gradient-purple-light d-flex align-items-center justify-center p-4">
    <v-card 
      class="w-full max-w-md" 
      elevation="8"
      rounded="xl"
    >
      <v-card-title class="text-center py-6">
        <v-icon 
          size="48" 
          color="primary"
          class="mx-auto mb-3"
        >
          mdi-forum
        </v-icon>
        <h2 class="text-2xl font-bold text-primary">校园论坛</h2>
        <p class="text-gray-500 mt-2">欢迎回来</p>
      </v-card-title>
      
      <v-card-text>
        <v-form @submit.prevent="handleLogin">
          <v-text-field
            v-model="form.username"
            label="用户名"
            placeholder="请输入用户名"
            prepend-icon="mdi-account"
            class="mb-4"
            :error-messages="error ? ['请输入用户名'] : []"
          />
          
          <v-text-field
            v-model="form.password"
            label="密码"
            type="password"
            placeholder="请输入密码"
            prepend-icon="mdi-lock"
            class="mb-6"
            :error-messages="error ? ['请输入密码'] : []"
          />
          
          <v-btn
            type="submit"
            color="primary"
            class="w-full mb-4"
            :loading="isLoading"
            block
            size="large"
          >
            <span v-if="!isLoading">登录</span>
            <span v-else>登录中...</span>
          </v-btn>
        </v-form>
        
        <div v-if="error" class="text-center mb-4">
          <v-alert 
            type="error" 
            border="bottom" 
            class="text-sm"
          >
            {{ error }}
          </v-alert>
        </div>
        
        <div class="text-center">
          <span class="text-gray-500">还没有账号？</span>
          <v-btn 
            text 
            color="primary" 
            @click="router.push('/register')"
          >
            立即注册
          </v-btn>
        </div>
      </v-card-text>
    </v-card>
  </div>
</template>