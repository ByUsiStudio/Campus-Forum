<script setup>import { ref, inject, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { adminApi } from '../api';
const router = useRouter();
const user = inject('user');
const activeTab = ref('users');
const users = ref([]);
const articles = ref([]);
const comments = ref([]);
const reports = ref([]);
const isLoading = ref(false);
const loadUsers = async () => {
 isLoading.value = true;
 try {
 const response = await adminApi.getUsers();
 users.value = response.data.users;
 }
 catch (error) {
 console.error('加载用户失败:', error);
 }
 finally {
 isLoading.value = false;
 }
};
const loadArticles = async () => {
 isLoading.value = true;
 try {
 const response = await adminApi.getArticles();
 articles.value = response.data.articles;
 }
 catch (error) {
 console.error('加载文章失败:', error);
 }
 finally {
 isLoading.value = false;
 }
};
const loadComments = async () => {
 isLoading.value = true;
 try {
 const response = await adminApi.getComments();
 comments.value = response.data.comments;
 }
 catch (error) {
 console.error('加载评论失败:', error);
 }
 finally {
 isLoading.value = false;
 }
};
const loadReports = async () => {
 isLoading.value = true;
 try {
 const response = await adminApi.getReports();
 reports.value = response.data.reports;
 }
 catch (error) {
 console.error('加载举报失败:', error);
 }
 finally {
 isLoading.value = false;
 }
};
const handleBanUser = async (userId) => {
 try {
 await adminApi.banUser(userId, { reason: '违规操作' });
 users.value = users.value.map(u => u.id === userId ? { ...u, is_banned: true } : u);
 }
 catch (error) {
 console.error('封禁用户失败:', error);
 }
};
const handleDeleteArticle = async (articleId) => {
 try {
 await adminApi.deleteArticle(articleId);
 articles.value = articles.value.filter(a => a.id !== articleId);
 }
 catch (error) {
 console.error('删除文章失败:', error);
 }
};
const handleDeleteComment = async (commentId) => {
 try {
 await adminApi.deleteComment(commentId);
 comments.value = comments.value.filter(c => c.id !== commentId);
 }
 catch (error) {
 console.error('删除评论失败:', error);
 }
};
const formatTime = (timeStr) => {
 const date = new Date(timeStr);
 return date.toLocaleString('zh-CN');
};
onMounted(() => {
 if (!user.value || user.value.role !== 'admin') {
 router.push('/');
 return;
 }
});
</script>

<template>
  <v-container class="max-w-6xl mx-auto py-8">
    <v-card rounded="xl" elevation="4">
      <v-card-title class="gradient-purple text-white">
        <v-icon class="mr-2">mdi-settings</v-icon>
        <span class="font-bold text-xl">管理后台</span>
      </v-card-title>
      
      <v-tabs v-model="activeTab" background-color="transparent" class="border-b">
        <v-tab value="users" @click="loadUsers">
          <v-icon class="mr-2">mdi-users</v-icon>
          用户管理
        </v-tab>
        <v-tab value="articles" @click="loadArticles">
          <v-icon class="mr-2">mdi-file</v-icon>
          文章管理
        </v-tab>
        <v-tab value="comments" @click="loadComments">
          <v-icon class="mr-2">mdi-comment</v-icon>
          评论管理
        </v-tab>
        <v-tab value="reports" @click="loadReports">
          <v-icon class="mr-2">mdi-flag</v-icon>
          举报管理
        </v-tab>
      </v-tabs>
      
      <v-tabs-items v-model="activeTab">
        <!-- 用户管理 -->
        <v-tab-item value="users">
          <div v-if="isLoading" class="text-center py-8">
            <v-progress-circular indeterminate color="primary" />
          </div>
          <v-data-table
            v-else
            :headers="[{ text: 'ID', value: 'id' }, { text: '用户名', value: 'username' }, { text: '昵称', value: 'display_name' }, { text: 'QQ号', value: 'qq_number' }, { text: '角色', value: 'role' }, { text: '状态', value: 'is_banned' }, { text: '注册时间', value: 'created_at' }, { text: '操作', value: 'actions' }]"
            :items="users"
            class="elevation-0"
          >
            <template v-slot:item.role="{ item }">
              <v-chip :color="item.role === 'admin' ? 'primary' : 'gray'">
                {{ item.role === 'admin' ? '管理员' : '普通用户' }}
              </v-chip>
            </template>
            <template v-slot:item.is_banned="{ item }">
              <span :class="item.is_banned ? 'text-error' : 'text-success'">
                {{ item.is_banned ? '已封禁' : '正常' }}
              </span>
            </template>
            <template v-slot:item.actions="{ item }">
              <v-btn 
                text 
                :color="item.is_banned ? 'success' : 'error'"
                @click="handleBanUser(item.id)"
              >
                {{ item.is_banned ? '解封' : '封禁' }}
              </v-btn>
            </template>
          </v-data-table>
        </v-tab-item>
        
        <!-- 文章管理 -->
        <v-tab-item value="articles">
          <div v-if="isLoading" class="text-center py-8">
            <v-progress-circular indeterminate color="primary" />
          </div>
          <v-data-table
            v-else
            :headers="[{ text: 'ID', value: 'id' }, { text: '标题', value: 'title' }, { text: '作者', value: 'user_name' }, { text: '分类', value: 'category_name' }, { text: '浏览量', value: 'view_count' }, { text: '创建时间', value: 'created_at' }, { text: '操作', value: 'actions' }]"
            :items="articles.map(a => ({ ...a, user_name: a.user?.display_name || a.user?.username, category_name: a.category?.name }))"
            class="elevation-0"
          >
            <template v-slot:item.actions="{ item }">
              <v-btn text color="error" @click="handleDeleteArticle(item.id)">
                删除
              </v-btn>
            </template>
          </v-data-table>
        </v-tab-item>
        
        <!-- 评论管理 -->
        <v-tab-item value="comments">
          <div v-if="isLoading" class="text-center py-8">
            <v-progress-circular indeterminate color="primary" />
          </div>
          <v-data-table
            v-else
            :headers="[{ text: 'ID', value: 'id' }, { text: '内容', value: 'content' }, { text: '作者', value: 'user_name' }, { text: '文章', value: 'article_title' }, { text: '创建时间', value: 'created_at' }, { text: '操作', value: 'actions' }]"
            :items="comments.map(c => ({ ...c, user_name: c.user?.display_name || c.user?.username, article_title: c.article?.title }))"
            class="elevation-0"
          >
            <template v-slot:item.actions="{ item }">
              <v-btn text color="error" @click="handleDeleteComment(item.id)">
                删除
              </v-btn>
            </template>
          </v-data-table>
        </v-tab-item>
        
        <!-- 举报管理 -->
        <v-tab-item value="reports">
          <div v-if="isLoading" class="text-center py-8">
            <v-progress-circular indeterminate color="primary" />
          </div>
          <v-data-table
            v-else
            :headers="[{ text: 'ID', value: 'id' }, { text: '举报类型', value: 'type' }, { text: '举报内容', value: 'content' }, { text: '举报者', value: 'reporter_name' }, { text: '状态', value: 'status' }, { text: '创建时间', value: 'created_at' }, { text: '操作', value: 'actions' }]"
            :items="reports.map(r => ({ ...r, reporter_name: r.reporter?.display_name || r.reporter?.username }))"
            class="elevation-0"
          >
            <template v-slot:item.status="{ item }">
              <v-chip :color="item.status === 'pending' ? 'warning' : 'success'">
                {{ item.status === 'pending' ? '待处理' : '已处理' }}
              </v-chip>
            </template>
            <template v-slot:item.actions="{ item }">
              <v-btn 
                text 
                color="primary" 
                @click="handleReport(item.id)"
                v-if="item.status === 'pending'"
              >
                处理
              </v-btn>
              <span v-else class="text-gray-400">已处理</span>
            </template>
          </v-data-table>
        </v-tab-item>
      </v-tabs-items>
    </v-card>
  </v-container>
</template>