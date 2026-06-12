<script setup>import { ref, inject, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { authApi, articleApi, favoriteApi, signinApi } from '../api';
const router = useRouter();
const user = inject('user');
const profile = ref(null);
const articles = ref([]);
const favorites = ref([]);
const activeTab = ref('articles');
const isLoading = ref(false);
const signinStatus = ref({
 hasSignedIn: false,
 signInDays: 0,
 totalSignIns: 0
});
const loadProfile = async () => {
 isLoading.value = true;
 try {
 const response = await authApi.getProfile();
 profile.value = response.data;
 }
 catch (error) {
 console.error('加载个人资料失败:', error);
 }
 finally {
 isLoading.value = false;
 }
};
const loadArticles = async () => {
 try {
 const response = await articleApi.getArticles({ user_id: profile.value.id });
 articles.value = response.data.articles;
 }
 catch (error) {
 console.error('加载文章失败:', error);
 }
};
const loadFavorites = async () => {
 try {
 const response = await favoriteApi.getFavorites();
 favorites.value = response.data.favorites;
 }
 catch (error) {
 console.error('加载收藏失败:', error);
 }
};
const loadSigninStatus = async () => {
 if (!user.value)
 return;
 try {
 const response = await signinApi.getStatus();
 signinStatus.value = response.data;
 }
 catch (error) {
 console.error('加载签到状态失败:', error);
 }
};
const handleSignin = async () => {
 if (!user.value) {
 router.push('/login');
 return;
 }
 try {
 const response = await signinApi.signin();
 signinStatus.value = {
 hasSignedIn: true,
 signInDays: response.data.sign_in_days,
 totalSignIns: response.data.total_sign_ins
 };
 }
 catch (error) {
 console.error('签到失败:', error);
 }
};
const formatTime = (timeStr) => {
 const date = new Date(timeStr);
 return date.toLocaleDateString('zh-CN');
};
onMounted(() => {
 if (!user.value) {
 router.push('/login');
 return;
 }
 loadProfile();
 loadSigninStatus();
});
</script>

<template>
  <v-container class="max-w-4xl mx-auto py-8">
    <v-card rounded="xl" elevation="4" class="mb-6" v-if="profile">
      <div class="gradient-purple p-6 rounded-t-xl">
        <div class="flex items-center">
          <v-avatar size="120" color="white">
            <v-icon size="60" color="primary">mdi-account</v-icon>
          </v-avatar>
          <div class="ml-6 text-white">
            <h2 class="text-2xl font-bold">{{ profile.display_name || profile.username }}</h2>
            <p class="text-white/80">@{{ profile.username }}</p>
            <p class="text-white/60 text-sm mt-1">{{ profile.bio || '暂无简介' }}</p>
          </div>
        </div>
        
        <div class="flex items-center mt-6">
          <div class="text-center mr-8">
            <span class="text-2xl font-bold">{{ profile.article_count }}</span>
            <p class="text-white/60 text-sm">文章</p>
          </div>
          <div class="text-center mr-8">
            <span class="text-2xl font-bold">{{ profile.follower_count }}</span>
            <p class="text-white/60 text-sm">粉丝</p>
          </div>
          <div class="text-center mr-8">
            <span class="text-2xl font-bold">{{ profile.following_count }}</span>
            <p class="text-white/60 text-sm">关注</p>
          </div>
          <div class="text-center">
            <span class="text-2xl font-bold">{{ signinStatus.totalSignIns }}</span>
            <p class="text-white/60 text-sm">签到</p>
          </div>
        </div>
      </div>
      
      <v-card-actions class="justify-center py-4">
        <v-btn 
          v-if="!signinStatus.hasSignedIn"
          color="primary" 
          @click="handleSignin"
        >
          <v-icon class="mr-1">mdi-calendar-check</v-icon>
          立即签到
        </v-btn>
        <div v-else class="text-success">
          <v-icon class="mr-1">mdi-check-circle</v-icon>
          今日已签到 · 连续 {{ signinStatus.signInDays }} 天
        </div>
      </v-card-actions>
    </v-card>
    
    <!-- Tab 切换 -->
    <v-card rounded="xl" elevation="4">
      <v-tabs v-model="activeTab" background-color="transparent" class="border-b">
        <v-tab value="articles">
          <v-icon class="mr-2">mdi-file</v-icon>
          我的文章
        </v-tab>
        <v-tab value="favorites">
          <v-icon class="mr-2">mdi-bookmark</v-icon>
          我的收藏
        </v-tab>
        <v-tab value="settings">
          <v-icon class="mr-2">mdi-settings</v-icon>
          个人设置
        </v-tab>
      </v-tabs>
      
      <v-tabs-items v-model="activeTab">
        <!-- 我的文章 -->
        <v-tab-item value="articles">
          <v-list v-if="articles.length > 0">
            <v-list-item 
              v-for="article in articles" 
              :key="article.id"
              class="border-b border-gray-100 last:border-0 cursor-pointer"
              @click="router.push(`/article/${article.id}`)"
            >
              <v-list-item-content>
                <v-list-item-title class="font-medium">{{ article.title }}</v-list-item-title>
                <v-list-item-subtitle>{{ article.content.slice(0, 50) }}...</v-list-item-subtitle>
              </v-list-item-content>
              <v-list-item-action>
                <span class="text-sm text-gray-500">{{ formatTime(article.created_at) }}</span>
              </v-list-item-action>
            </v-list-item>
          </v-list>
          <div v-else class="text-center py-12 text-gray-500">
            <v-icon size="48" class="mx-auto mb-2">mdi-file-question</v-icon>
            <p>暂无文章</p>
          </div>
        </v-tab-item>
        
        <!-- 我的收藏 -->
        <v-tab-item value="favorites">
          <v-list v-if="favorites.length > 0">
            <v-list-item 
              v-for="item in favorites" 
              :key="item.article.id"
              class="border-b border-gray-100 last:border-0 cursor-pointer"
              @click="router.push(`/article/${item.article.id}`)"
            >
              <v-list-item-content>
                <v-list-item-title class="font-medium">{{ item.article.title }}</v-list-item-title>
                <v-list-item-subtitle>{{ item.article.user?.display_name || item.article.user?.username }}</v-list-item-subtitle>
              </v-list-item-content>
              <v-list-item-action>
                <span class="text-sm text-gray-500">{{ formatTime(item.created_at) }}</span>
              </v-list-item-action>
            </v-list-item>
          </v-list>
          <div v-else class="text-center py-12 text-gray-500">
            <v-icon size="48" class="mx-auto mb-2">mdi-bookmark-outline</v-icon>
            <p>暂无收藏</p>
          </div>
        </v-tab-item>
        
        <!-- 个人设置 -->
        <v-tab-item value="settings">
          <v-card-text class="px-8">
            <v-text-field
              v-model="profile.display_name"
              label="昵称"
              class="mb-4"
            />
            <v-textarea
              v-model="profile.bio"
              label="个人简介"
              rows="3"
              class="mb-4"
            />
            <v-btn color="primary" @click="handleUpdateProfile">
              保存设置
            </v-btn>
          </v-card-text>
        </v-tab-item>
      </v-tabs-items>
    </v-card>
  </v-container>
</template>