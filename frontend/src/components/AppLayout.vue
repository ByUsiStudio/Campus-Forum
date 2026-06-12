<script setup>
import { ref, inject, computed, onMounted, onUnmounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import UserAvatar from './UserAvatar.vue';

const router = useRouter();
const route = useRoute();
const user = inject('user');
const siteConfig = inject('siteConfig');

const isDrawerOpen = ref(false);
const isMobile = ref(false);
const searchQuery = ref('');
const currentPath = computed(() => route.path);

const handleResize = () => {
  isMobile.value = window.innerWidth < 1024;
};

onMounted(() => {
  handleResize();
  window.addEventListener('resize', handleResize);
});

onUnmounted(() => {
  window.removeEventListener('resize', handleResize);
});

const toggleDrawer = () => {
  isDrawerOpen.value = !isDrawerOpen.value;
};

const closeDrawer = () => {
  isDrawerOpen.value = false;
};

const goTo = (path) => {
  router.push(path);
  closeDrawer();
};

const handleSearch = () => {
  if (searchQuery.value.trim()) {
    router.push('/search?query=' + encodeURIComponent(searchQuery.value.trim()));
    searchQuery.value = '';
  }
};

const bottomNavItems = [
  { title: '首页', icon: 'mdi-home', path: '/' },
  { title: '搜索', icon: 'mdi-magnify', path: '/search' },
  { title: '发布', icon: 'mdi-plus-circle-outline', path: '/create' },
  { title: '通知', icon: 'mdi-bell-outline', path: '/notifications' },
  { title: '我的', icon: 'mdi-account-outline', path: '/profile' }
];

const drawerItems = [
  { title: '好友', icon: 'mdi-account-group-outline', path: '/friends' },
  { title: '聊天', icon: 'mdi-chat-outline', path: '/chat' },
  { title: '收藏', icon: 'mdi-heart-outline', path: '/profile' },
  { title: '设置', icon: 'mdi-cog-outline', path: '/profile' },
  { title: '管理员', icon: 'mdi-shield-outline', path: '/admin' }
];

const topNavItems = [
  { title: '首页', path: '/' },
  { title: '分类', path: '/category/1' },
  { title: '好友', path: '/friends' },
  { title: '聊天', path: '/chat' }
];

const logout = () => {
  localStorage.removeItem('token');
  localStorage.removeItem('user');
  router.push('/login');
};

const isActiveRoute = (path) => {
  if (path === '/') return currentPath.value === '/';
  return currentPath.value.startsWith(path);
};
</script>

<template>
  <v-app>
    <!-- 移动端底部导航 -->
    <v-bottom-navigation
      v-if="isMobile"
      v-model="currentPath"
      :bg-color="'white'"
      color="primary"
      grow
      class="mobile-bottom-nav"
    >
      <v-btn
        v-for="item in bottomNavItems"
        :key="item.path"
        :value="item.path"
        :to="item.path"
      >
        <v-icon>{{ item.icon }}</v-icon>
        <span class="text-caption">{{ item.title }}</span>
      </v-btn>
    </v-bottom-navigation>

    <!-- PC端顶部导航栏 -->
    <v-app-bar
      v-if="!isMobile"
      :elevation="1"
      color="white"
      class="app-header"
      fixed
    >
      <!-- Logo -->
      <v-app-bar-title class="ml-4">
        <router-link to="/" class="text-decoration-none">
          <div class="d-flex align-center">
            <v-icon color="primary" size="28" class="mr-2">mdi-forum</v-icon>
            <span class="text-h6 font-weight-bold text-grey-darken-3">
              {{ siteConfig?.site_title || '校园论坛' }}
            </span>
          </div>
        </router-link>
      </v-app-bar-title>

      <!-- 导航菜单 -->
      <v-tabs
        v-model="currentPath"
        color="primary"
        slider-size="2"
        class="ml-8"
      >
        <v-tab
          v-for="item in topNavItems"
          :key="item.path"
          :to="item.path"
          :value="item.path"
          class="text-body-1"
        >
          {{ item.title }}
        </v-tab>
      </v-tabs>

      <v-spacer></v-spacer>

      <!-- 搜索框 -->
      <v-text-field
        v-model="searchQuery"
        placeholder="搜索帖子、用户..."
        prepend-inner-icon="mdi-magnify"
        density="compact"
        variant="outlined"
        class="search-field mx-4"
        style="max-width: 280px;"
        hide-details
        @keyup.enter="handleSearch"
      />

      <!-- 通知按钮 -->
      <v-btn
        v-if="user"
        icon
        :to="'/notifications'"
        variant="text"
        color="grey-darken-1"
        class="mr-1"
      >
        <v-badge
          color="error"
          dot
        >
          <v-icon>mdi-bell-outline</v-icon>
        </v-badge>
      </v-btn>

      <!-- 用户菜单 -->
      <template v-if="user">
        <v-menu offset-y location="bottom end">
          <template #activator="{ props }">
            <v-btn
              icon
              variant="text"
              v-bind="props"
              class="mr-2"
            >
              <UserAvatar :user="user" :size="36" />
            </v-btn>
          </template>
          <v-card min-width="200" class="user-menu-card">
            <v-list density="compact">
              <v-list-item to="/profile" prepend-icon="mdi-account-outline">
                <v-list-item-title>个人中心</v-list-item-title>
              </v-list-item>
              <v-list-item v-if="user.role === 'admin'" to="/admin" prepend-icon="mdi-shield-outline">
                <v-list-item-title>管理后台</v-list-item-title>
              </v-list-item>
              <v-divider class="my-1"></v-divider>
              <v-list-item @click="logout" prepend-icon="mdi-logout" class="text-error">
                <v-list-item-title>退出登录</v-list-item-title>
              </v-list-item>
            </v-list>
          </v-card>
        </v-menu>
      </template>

      <template v-else>
        <v-btn
          to="/login"
          variant="text"
          color="primary"
          class="mr-2"
        >
          登录
        </v-btn>
        <v-btn
          to="/register"
          variant="flat"
          color="primary"
          class="mr-4"
        >
          注册
        </v-btn>
      </template>
    </v-app-bar>

    <!-- 移动端顶部栏 -->
    <v-app-bar
      v-if="isMobile"
      :elevation="0"
      color="white"
      class="mobile-app-bar"
      fixed
    >
      <v-app-bar-title class="text-center">
        <router-link to="/" class="text-decoration-none">
          <div class="d-flex align-center justify-center">
            <v-icon color="primary" size="24" class="mr-1">mdi-forum</v-icon>
            <span class="text-subtitle-1 font-weight-bold text-grey-darken-3">
              {{ siteConfig?.site_title || '论坛' }}
            </span>
          </div>
        </router-link>
      </v-app-bar-title>
    </v-app-bar>

    <!-- 主内容区 -->
    <v-main :class="{ 'pb-mobile-nav': isMobile, 'pt-app-bar': true }">
      <router-view v-slot="{ Component }">
        <transition name="page" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </v-main>

    <!-- 移动端功能菜单 -->
    <v-navigation-drawer
      v-if="isMobile"
      v-model="isDrawerOpen"
      location="bottom"
      temporary
      class="mobile-drawer"
    >
      <v-list density="compact">
        <v-list-item
          v-for="item in drawerItems"
          :key="item.path"
          @click="goTo(item.path)"
          :prepend-icon="item.icon"
        >
          <v-list-item-title>{{ item.title }}</v-list-item-title>
        </v-list-item>
        <v-divider class="my-2"></v-divider>
        <v-list-item
          v-if="user"
          @click="logout"
          prepend-icon="mdi-logout"
          class="text-error"
        >
          <v-list-item-title>退出登录</v-list-item-title>
        </v-list-item>
      </v-list>
    </v-navigation-drawer>

    <!-- 移动端菜单按钮 -->
    <v-btn
      v-if="isMobile && user"
      icon
      variant="flat"
      color="primary"
      class="menu-fab"
      @click="toggleDrawer"
    >
      <v-icon>mdi-menu</v-icon>
    </v-btn>
  </v-app>
</template>

<style scoped>
.app-header {
  border-bottom: 1px solid rgba(149, 117, 205, 0.1) !important;
}

.search-field :deep(.v-field) {
  border-radius: 20px;
  background: #FAFAFA;
}

.search-field :deep(.v-field__outline) {
  border-color: rgba(149, 117, 205, 0.2);
}

.search-field :deep(.v-field--focused .v-field__outline) {
  border-color: #9575CD;
}

.user-menu-card {
  border-radius: 12px !important;
}

.menu-fab {
  position: fixed !important;
  bottom: 72px;
  right: 16px;
  z-index: 1001;
}

/* 移动端顶部栏 */
.mobile-app-bar {
  border-bottom: 1px solid rgba(149, 117, 205, 0.1) !important;
}

/* 移动端底部导航栏 */
.mobile-bottom-nav {
  position: fixed !important;
  bottom: 0;
  left: 0;
  right: 0;
  z-index: 1000;
  box-shadow: 0 -2px 10px rgba(149, 117, 205, 0.15) !important;
  border-top: 1px solid rgba(149, 117, 205, 0.1);
}

/* 主内容区顶部 padding，确保不被导航栏遮挡 */
.pt-app-bar {
  padding-top: 64px !important;
}

/* 主内容区底部 padding，确保不被导航栏遮挡 */
.pb-mobile-nav {
  padding-bottom: 56px !important;
}
</style>
