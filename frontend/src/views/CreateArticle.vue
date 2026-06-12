<script setup>import { ref, inject, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { articleApi, categoryApi } from '../api';
const router = useRouter();
const user = inject('user');
const form = ref({
 title: '',
 content: '',
 category_id: ''
});
const categories = ref([]);
const isLoading = ref(false);
const error = ref('');
const loadCategories = async () => {
 try {
 const response = await categoryApi.getCategories();
 categories.value = response.data.categories;
 }
 catch (error) {
 console.error('加载分类失败:', error);
 }
};
const handleSubmit = async () => {
 if (!form.value.title.trim()) {
 error.value = '请输入文章标题';
 return;
 }
 if (!form.value.content.trim()) {
 error.value = '请输入文章内容';
 return;
 }
 if (!form.value.category_id) {
 error.value = '请选择分类';
 return;
 }
 isLoading.value = true;
 try {
 await articleApi.createArticle(form.value);
 router.push('/');
 }
 catch (err) {
 error.value = err.response?.data?.error || '发布失败';
 }
 finally {
 isLoading.value = false;
 }
};
onMounted(() => {
 if (!user.value) {
 router.push('/login');
 return;
 }
 loadCategories();
});
</script>

<template>
  <v-container class="max-w-3xl mx-auto py-8">
    <v-card rounded="xl" elevation="4">
      <v-card-title class="gradient-purple text-white">
        <v-icon class="mr-2">mdi-file-plus</v-icon>
        <span class="font-bold">发布文章</span>
      </v-card-title>
      
      <v-card-text>
        <v-text-field
          v-model="form.title"
          label="文章标题"
          placeholder="请输入文章标题"
          class="mb-4"
          :error-messages="error ? ['请输入标题'] : []"
        />
        
        <v-select
          v-model="form.category_id"
          :items="categories"
          item-text="name"
          item-value="id"
          label="选择分类"
          placeholder="请选择文章分类"
          class="mb-4"
          :error-messages="error ? ['请选择分类'] : []"
        />
        
        <v-textarea
          v-model="form.content"
          label="文章内容"
          placeholder="写下你的想法..."
          rows="15"
          rounded="lg"
          class="mb-6"
          :error-messages="error ? ['请输入内容'] : []"
        />
        
        <div v-if="error" class="mb-4">
          <v-alert type="error" border="bottom" class="text-sm">
            {{ error }}
          </v-alert>
        </div>
        
        <v-btn
          color="primary"
          class="mr-4"
          :loading="isLoading"
          @click="handleSubmit"
        >
          <span v-if="!isLoading">发布文章</span>
          <span v-else>发布中...</span>
        </v-btn>
        
        <v-btn
          text
          color="gray"
          @click="router.push('/')"
        >
          取消
        </v-btn>
      </v-card-text>
    </v-card>
  </v-container>
</template>