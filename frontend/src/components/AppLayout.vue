<script setup>import { ref, inject, computed, onMounted, onUnmounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import UserAvatar from './UserAvatar.vue';
const router = useRouter();
const route = useRoute();
const user = inject('user');
const siteConfig = inject('siteConfig');
const isDrawerOpen = ref(false);
const isMobile = ref(false);
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
const bottomNavItems = [
 { title: '首页', icon: 'mdi-home', path: '/' },
 { title: '搜索', icon: 'mdi-magnify', path: '/search' },
 { title: '发布', icon: 'mdi-plus-circle', path: '/create' },
 { title: '通知', icon: 'mdi-bell', path: '/notifications' },
 { title: '我的', icon: 'mdi-account', path: '/profile' }
];
const drawerItems = [
 { title: '好友', icon: 'mdi-users', path: '/friends' },
 { title: '聊天', icon: 'mdi-chat', path: '/chat' },
 { title: '收藏', icon: 'mdi-heart', path: '/profile' },
 { title: '设置', icon: 'mdi-settings', path: '/profile' },
 { title: '管理员', icon: 'mdi-shield', path: '/admin' }
];
const topNavItems = [
 { title: '首页', path: '/' },
 { title: '分类', path: '/category/1' },
 { title: '好友', path: '/friends' },
 { title: '聊天', path: '/chat' },
 { title: '通知', path: '/notifications' }
];
const logout = () => {
 localStorage.removeItem('token');
 localStorage.removeItem('user');
 router.push('/login');
};
</script>

<template>
  <v-app>
    <!-- 移动端底部栏 -->
    <v-bottom-navigation
      v-if="isMobile"
      v-model="currentPath"
      class="mdc-bottom-navigation--shifting"
      fixed
      app
    >
      <v-btn
        v-for="item in bottomNavItems"
        :key="item.path"
        :value="item.path"
        :to="item.path"
        active-class="primary"
      >
        <template v-slot:icon>
          <v-icon>{{ item.icon }}</v-icon>
        </template>
        <template v-slot:text>
          {{ item.title }}
        </template>
      </v-btn>
    </v-bottom-navigation>

    <!-- 移动端功能抽屉按钮 -->
    <v-btn
      v-if="isMobile"
      fab
      bottom
      right
      color="primary"
      class="mr-4 mb-24"
      @click="toggleDrawer"
    >
      <v-icon>mdi-apps</v-icon>
    </v-btn>

    <!-- 功能抽屉 -->
    <v-bottom-sheet
      v-model="isDrawerOpen"
      inset
    >
      <v-card class="pa-4">
        <v-card-title class="text-h6">功能菜单</v-card-title>
        <v-list>
          <v-list-item
            v-for="item in drawerItems"
            :key="item.path"
            @click="goTo(item.path)"
            class="cursor-pointer"
          >
            <v-list-item-icon>
              <v-icon>{{ item.icon }}</v-icon>
            </v-list-item-icon>
            <v-list-item-title>{{ item.title }}</v-list-item-title>
          </v-list-item>
          <v-divider class="my-2"></v-divider>
          <v-list-item
            v-if="user"
            @click="logout"
            class="cursor-pointer"
          >
            <v-list-item-icon>
              <v-icon color="error">mdi-logout</v-icon>
            </v-list-item-icon>
            <v-list-item-title class="text-error">退出登录</v-list-item-title>
          </v-list-item>
        </v-list>
      </v-card>
    </v-bottom-sheet>

    <!-- PC端顶部栏 -->
    <v-app-bar
      v-if="!isMobile"
      app
      color="primary"
      dark
    >
      <v-toolbar-title>
        <v-btn to="/" text class="text-white">
          {{ siteConfig?.site_title || '校园论坛' }}
        </v-btn>
      </v-toolbar-title>

      <v-spacer></v-spacer>

      <!-- 导航菜单 -->
      <v-btn
        v-for="item in topNavItems"
        :key="item.path"
        :to="item.path"
        text
        :class="{ 'font-weight-bold': currentPath === item.path }"
      >
        {{ item.title }}
      </v-btn>

      <v-spacer></v-spacer>

      <!-- 搜索框 -->
      <v-text-field
        v-model="searchQuery"
        placeholder="搜索..."
        prepend-inner-icon="mdi-magnify"
        class="hidden-sm-and-down"
        style="width: 300px;"
        @keyup.enter="router.push('/search?query=' + searchQuery)"
      ></v-text-field>

      <!-- 用户菜单 -->
      <template v-if="user">
        <v-btn icon to="/notifications">
          <v-icon>mdi-bell</v-icon>
        </v-btn>
        <v-menu offset-y>
          <template v-slot:activator="{ props }">
            <v-btn
              icon
              v-bind="props"
            >
              <UserAvatar :user="user" :size="40" />
            </v-btn>
          </template>
          <v-list>
            <v-list-item to="/profile">
              <v-list-item-icon>
                <v-icon>mdi-account</v-icon>
              </v-list-item-icon>
              <v-list-item-title>个人中心</v-list-item-title>
            </v-list-item>
            <v-list-item v-if="user.role === 'admin'" to="/admin">
              <v-list-item-icon>
                <v-icon>mdi-shield</v-icon>
              </v-list-item-icon>
              <v-list-item-title>管理后台</v-list-item-title>
            </v-list-item>
            <v-divider></v-divider>
            <v-list-item @click="logout">
              <v-list-item-icon>
                <v-icon color="error">mdi-logout</v-icon>
              </v-list-item-icon>
              <v-list-item-title class="text-error">退出登录</v-list-item-title>
            </v-list-item>
          </v-list>
        </v-menu>
      </template>

      <template v-else>
        <v-btn to="/login" text>登录</v-btn>
        <v-btn to="/register" text>注册</v-btn>
      </template>
    </v-app-bar>

    <!-- 主内容区 -->
    <v-main :class="{ 'pb-20': isMobile }">
      <router-view />
    </v-main>
  </v-app>
</template>

<script>
export default {
  name: 'AppLayout'
}
</script>
