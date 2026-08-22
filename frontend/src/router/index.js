import { createRouter, createWebHistory } from 'vue-router'

// 无需登录即可访问的浏览类页面 / 认证类页面（列表前缀匹配：文章、分类）
const publicPrefixes = [
  '/login', '/register', '/forgot-password', '/search',
  '/leaderboard', '/topics', '/category/', '/article/', '/video', '/403'
]

function isPublic(path) {
  return publicPrefixes.some(p => path === p || path.startsWith(p))
}

const routes = [
  {
    path: '/',
    name: 'Home',
    component: () => import('@/views/Home.vue'),
    meta: { public: true, title: '首页' }
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { public: true, title: '登录' }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/Register.vue'),
    meta: { public: true, title: '注册' }
  },
  {
    path: '/forgot-password',
    name: 'ForgotPassword',
    component: () => import('@/views/ForgotPassword.vue'),
    meta: { public: true, title: '找回密码' }
  },
  {
    path: '/article/:id',
    name: 'Article',
    component: () => import('@/views/Article.vue'),
    meta: { public: true, title: '文章详情' }
  },
  {
    path: '/create',
    name: 'CreateArticle',
    component: () => import('@/views/CreateArticle.vue'),
    meta: { title: '写文章' }
  },
  {
    path: '/profile',
    name: 'Profile',
    component: () => import('@/views/Profile.vue'),
    meta: { title: '个人中心' }
  },
  {
    path: '/admin',
    component: () => import('@/views/admin/Admin.vue'),
    meta: { requiresAdmin: true, title: '管理后台' },
    children: [
      { path: '', redirect: { name: 'AdminIndex' } },
      { path: 'index', name: 'AdminIndex', component: () => import('@/views/admin/AdminIndex.vue'), meta: { title: '数据概览' } },
      { path: 'statistics', name: 'AdminStatistics', component: () => import('@/views/admin/AdminStatistics.vue'), meta: { title: '数据统计' } },
      { path: 'users', name: 'AdminUsers', component: () => import('@/views/admin/AdminUsers.vue'), meta: { title: '用户管理' } },
      { path: 'articles', name: 'AdminArticles', component: () => import('@/views/admin/AdminArticles.vue'), meta: { title: '文章管理' } },
      { path: 'comments', name: 'AdminComments', component: () => import('@/views/admin/AdminComments.vue'), meta: { title: '评论管理' } },
      { path: 'categories', name: 'AdminCategories', component: () => import('@/views/admin/AdminCategories.vue'), meta: { title: '分区管理' } },
      { path: 'titles', name: 'AdminTitles', component: () => import('@/views/admin/AdminTitles.vue'), meta: { title: '头衔管理' } },
      { path: 'sidebar', name: 'AdminSidebar', component: () => import('@/views/admin/AdminSidebar.vue'), meta: { title: '侧边栏配置' } },
      { path: 'deletions', name: 'AdminDeletions', component: () => import('@/views/admin/AdminDeletions.vue'), meta: { title: '删除申请' } },
      { path: 'reports', name: 'AdminReports', component: () => import('@/views/admin/AdminReports.vue'), meta: { title: '举报管理' } },
      { path: 'announcement', name: 'AdminAnnouncement', component: () => import('@/views/admin/AdminAnnouncement.vue'), meta: { title: '公告管理' } },
      { path: 'notifications', name: 'AdminUserNotifications', component: () => import('@/views/admin/AdminUserNotifications.vue'), meta: { title: '用户通知与权限' } },
      { path: 'system-logs', name: 'AdminSystemLogs', component: () => import('@/views/admin/AdminSystemLogs.vue'), meta: { title: '系统日志' } },
      { path: 'siteconfig', name: 'AdminSiteConfig', component: () => import('@/views/admin/AdminSiteConfig.vue'), meta: { title: '网站配置' } },
      { path: 'smtpconfig', name: 'AdminSMTPConfig', component: () => import('@/views/admin/AdminSMTPConfig.vue'), meta: { title: '邮件配置' } },
      { path: 'admin-notifications', name: 'AdminNotifications', component: () => import('@/views/admin/AdminNotifications.vue'), meta: { title: '系统通知' } }
    ]
  },
  { path: '/category/:id', name: 'Category', component: () => import('@/views/Category.vue'), meta: { public: true, title: '分区' } },
  { path: '/video', name: 'VideoPlayer', component: () => import('@/views/VideoPlayerPage.vue'), meta: { public: true, title: '视频' } },
  { path: '/notifications', name: 'Notifications', component: () => import('@/views/Notifications.vue'), meta: { title: '我的通知' } },
  { path: '/follow-list', name: 'FollowList', component: () => import('@/views/FollowList.vue'), meta: { title: '关注/粉丝' } },
  { path: '/search', name: 'Search', component: () => import('@/views/Search.vue'), meta: { public: true, title: '搜索' } },
  { path: '/signin', name: 'SignIn', component: () => import('@/views/SignIn.vue'), meta: { title: '每日签到' } },
  { path: '/leaderboard', name: 'Leaderboard', component: () => import('@/views/Leaderboard.vue'), meta: { public: true, title: '排行榜' } },
  { path: '/topics', name: 'TopicList', component: () => import('@/views/TopicList.vue'), meta: { public: true, title: '话题' } },
  { path: '/collections', name: 'CollectionList', component: () => import('@/views/CollectionList.vue'), meta: { title: '我的收藏' } },
  { path: '/403', name: 'Forbidden', component: () => import('@/views/Forbidden.vue'), meta: { public: true, title: '无权访问' } },
  { path: '/:pathMatch(.*)*', name: 'NotFound', component: () => import('@/views/NotFound.vue'), meta: { public: true, title: '页面未找到' } }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior() {
    return { top: 0 }
  }
})

router.beforeEach((to) => {
  const loggedIn = !!localStorage.getItem('token')

  // 已登录访问认证页 → 回首页
  if (loggedIn && (to.path === '/login' || to.path === '/register')) {
    return { path: '/' }
  }

  const publicish = to.meta.public || isPublic(to.path)

  // 受限页面：未登录跳登录页，并携带目标以便登录后回跳
  if (!publicish && !loggedIn) {
    return { path: '/login', query: to.fullPath !== '/' ? { redirect: to.fullPath } : {} }
  }

  return true
})

export default router
