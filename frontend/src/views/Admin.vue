<script setup>
import { ref, inject, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { adminAuthApi, adminStatsApi } from '../api/admin'

const router = useRouter()
const user = inject('user')

const activeTab = ref('overview')
const isAdmin = ref(false)
const statistics = ref({
  user_count: 0,
  article_count: 0,
  comment_count: 0,
  view_count: 0
})
const isLoading = ref(true)

const menuItems = [
  { title: '概览', value: 'overview', icon: 'mdi-view-dashboard' },
  { title: '用户管理', value: 'users', icon: 'mdi-account-group' },
  { title: '文章管理', value: 'articles', icon: 'mdi-file-document' },
  { title: '评论管理', value: 'comments', icon: 'mdi-comment' },
  { title: '分区管理', value: 'categories', icon: 'mdi-folder' },
  { title: '头衔管理', value: 'titles', icon: 'mdi-medal' },
  { title: '举报管理', value: 'reports', icon: 'mdi-alert' },
  { title: '公告管理', value: 'announcement', icon: 'mdi-bullhorn' },
  { title: '网站配置', value: 'site-config', icon: 'mdi-cog' },
  { title: '侧边栏配置', value: 'sidebar', icon: 'mdi-menu' },
  { title: '用户通知', value: 'user-notifications', icon: 'mdi-bell' }
]

const checkAdmin = async () => {
  try {
    const response = await adminAuthApi.checkAdmin()
    isAdmin.value = response.data.is_admin
    if (!isAdmin.value) {
      router.push('/')
    }
  } catch (error) {
    console.error('检查管理员权限失败:', error)
    router.push('/')
  }
}

const loadStatistics = async () => {
  try {
    const response = await adminStatsApi.getStatistics()
    statistics.value = response.data
  } catch (error) {
    console.error('加载统计数据失败:', error)
  }
}

onMounted(async () => {
  if (!user.value) {
    router.push('/login')
    return
  }
  await checkAdmin()
  if (isAdmin.value) {
    await loadStatistics()
  }
  isLoading.value = false
})
</script>

<template>
  <v-app>
    <!-- 顶部导航 -->
    <v-app-bar app color="primary">
      <v-btn icon @click="router.push('/')">
        <v-icon>mdi-arrow-left</v-icon>
      </v-btn>
      <v-toolbar-title>管理后台</v-toolbar-title>
      <v-spacer></v-spacer>
      <v-chip color="white" text-color="primary">
        <v-icon start>mdi-account</v-icon>
        {{ user?.username }}
      </v-chip>
    </v-app-bar>
    
    <!-- 加载状态 -->
    <v-main v-if="isLoading" class="d-flex align-center justify-center">
      <v-progress-circular indeterminate color="primary" size="64" />
    </v-main>
    
    <!-- 主内容 -->
    <v-main v-else>
      <v-container fluid class="pa-0">
        <v-row no-gutters>
          <!-- 左侧导航 -->
          <v-col cols="12" md="3" lg="2">
            <v-card flat tile class="h-100">
              <v-list dense nav>
                <v-list-item
                  v-for="item in menuItems"
                  :key="item.value"
                  :value="item.value"
                  :active="activeTab === item.value"
                  @click="activeTab = item.value"
                >
                  <v-list-item-icon>
                    <v-icon>{{ item.icon }}</v-icon>
                  </v-list-item-icon>
                  <v-list-item-title>{{ item.title }}</v-list-item-title>
                </v-list-item>
              </v-list>
            </v-card>
          </v-col>
          
          <!-- 右侧内容 -->
          <v-col cols="12" md="9" lg="10">
            <v-container fluid class="pa-4">
              <!-- 概览 -->
              <template v-if="activeTab === 'overview'">
                <v-row>
                  <v-col cols="12" sm="6" md="3">
                    <v-card>
                      <v-card-text class="text-center">
                        <v-icon size="48" color="primary">mdi-account-group</v-icon>
                        <div class="text-h4 mt-2">{{ statistics.user_count }}</div>
                        <div class="text-grey">用户总数</div>
                      </v-card-text>
                    </v-card>
                  </v-col>
                  <v-col cols="12" sm="6" md="3">
                    <v-card>
                      <v-card-text class="text-center">
                        <v-icon size="48" color="success">mdi-file-document</v-icon>
                        <div class="text-h4 mt-2">{{ statistics.article_count }}</div>
                        <div class="text-grey">文章总数</div>
                      </v-card-text>
                    </v-card>
                  </v-col>
                  <v-col cols="12" sm="6" md="3">
                    <v-card>
                      <v-card-text class="text-center">
                        <v-icon size="48" color="warning">mdi-comment</v-icon>
                        <div class="text-h4 mt-2">{{ statistics.comment_count }}</div>
                        <div class="text-grey">评论总数</div>
                      </v-card-text>
                    </v-card>
                  </v-col>
                  <v-col cols="12" sm="6" md="3">
                    <v-card>
                      <v-card-text class="text-center">
                        <v-icon size="48" color="info">mdi-eye</v-icon>
                        <div class="text-h4 mt-2">{{ statistics.view_count }}</div>
                        <div class="text-grey">总浏览量</div>
                      </v-card-text>
                    </v-card>
                  </v-col>
                </v-row>
              </template>
              
              <!-- 用户管理 -->
              <template v-if="activeTab === 'users'">
                <AdminUsers />
              </template>
              
              <!-- 文章管理 -->
              <template v-if="activeTab === 'articles'">
                <AdminArticles />
              </template>
              
              <!-- 评论管理 -->
              <template v-if="activeTab === 'comments'">
                <AdminComments />
              </template>
              
              <!-- 分区管理 -->
              <template v-if="activeTab === 'categories'">
                <AdminCategories />
              </template>
              
              <!-- 头衔管理 -->
              <template v-if="activeTab === 'titles'">
                <AdminTitles />
              </template>
              
              <!-- 举报管理 -->
              <template v-if="activeTab === 'reports'">
                <AdminReports />
              </template>
              
              <!-- 公告管理 -->
              <template v-if="activeTab === 'announcement'">
                <AdminAnnouncement />
              </template>
              
              <!-- 网站配置 -->
              <template v-if="activeTab === 'site-config'">
                <AdminSiteConfig />
              </template>
              
              <!-- 侧边栏配置 -->
              <template v-if="activeTab === 'sidebar'">
                <AdminSidebar />
              </template>
              
              <!-- 用户通知 -->
              <template v-if="activeTab === 'user-notifications'">
                <AdminUserNotifications />
              </template>
            </v-container>
          </v-col>
        </v-row>
      </v-container>
    </v-main>
  </v-app>
</template>

<script>
import AdminUsers from './admin/AdminUsers.vue'
import AdminArticles from './admin/AdminArticles.vue'
import AdminComments from './admin/AdminComments.vue'
import AdminCategories from './admin/AdminCategories.vue'
import AdminTitles from './admin/AdminTitles.vue'
import AdminReports from './admin/AdminReports.vue'
import AdminAnnouncement from './admin/AdminAnnouncement.vue'
import AdminSiteConfig from './admin/AdminSiteConfig.vue'
import AdminSidebar from './admin/AdminSidebar.vue'
import AdminUserNotifications from './admin/AdminUserNotifications.vue'

export default {
  components: {
    AdminUsers,
    AdminArticles,
    AdminComments,
    AdminCategories,
    AdminTitles,
    AdminReports,
    AdminAnnouncement,
    AdminSiteConfig,
    AdminSidebar,
    AdminUserNotifications
  }
}
</script>