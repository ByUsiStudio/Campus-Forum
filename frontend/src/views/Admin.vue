<template>
  <div v-if="isInitialized && isAdmin">
    <v-card class="pa-6">
      <v-card-title class="text-h5 mb-4" style="color: rgb(var(--v-theme-primary));">
        管理后台
      </v-card-title>
      
      <v-tabs v-model="activeTab" color="primary" class="mb-4">
        <v-tab value="overview">概览</v-tab>
        <v-tab value="users">用户管理</v-tab>
        <v-tab value="articles">文章管理</v-tab>
        <v-tab value="comments">评论管理</v-tab>
        <v-tab value="categories">分区管理</v-tab>
        <v-tab value="sidebar">侧边栏配置</v-tab>
        <v-tab value="deletions">
          <v-badge :content="deletionRequests.length" color="error" :model-value="deletionRequests.length > 0">
            删除审核
          </v-badge>
        </v-tab>
        <v-tab value="announcement">公告管理</v-tab>
        <v-tab value="notifications">通知管理</v-tab>
      </v-tabs>
      
      <v-window v-model="activeTab">
        <!-- 概览 -->
        <v-window-item value="overview">
          <v-row>
            <v-col cols="12" sm="6" md="3">
              <v-card color="primary" class="pa-4">
                <div class="text-h3 text-white">{{ statistics.user_count }}</div>
                <div class="text-body-1 text-white opacity-80">用户总数</div>
              </v-card>
            </v-col>
            <v-col cols="12" sm="6" md="3">
              <v-card color="success" class="pa-4">
                <div class="text-h3 text-white">{{ statistics.article_count }}</div>
                <div class="text-body-1 text-white opacity-80">文章总数</div>
              </v-card>
            </v-col>
            <v-col cols="12" sm="6" md="3">
              <v-card color="info" class="pa-4">
                <div class="text-h3 text-white">{{ statistics.comment_count }}</div>
                <div class="text-body-1 text-white opacity-80">评论总数</div>
              </v-card>
            </v-col>
            <v-col cols="12" sm="6" md="3">
              <v-card color="warning" class="pa-4">
                <div class="text-h3 text-white">{{ statistics.view_count }}</div>
                <div class="text-body-1 text-white opacity-80">总浏览量</div>
              </v-card>
            </v-col>
          </v-row>
        </v-window-item>
        
        <!-- 用户管理 -->
        <v-window-item value="users">
          <v-card variant="outlined" class="pa-4 mb-4">
            <div class="text-body-2 text-medium-emphasis mb-2">用户总数：{{ users.length }}</div>
          </v-card>
          
          <v-table>
            <thead>
              <tr>
                <th>ID</th>
                <th>用户名</th>
                <th>显示名称</th>
                <th>QQ号码</th>
                <th>角色</th>
                <th>状态</th>
                <th>注册时间</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="user in users" :key="user.id">
                <td>{{ user.id }}</td>
                <td>{{ user.username }}</td>
                <td>{{ user.display_name }}</td>
                <td>{{ user.qq_number }}</td>
                <td>
                  <v-chip size="small" :color="user.role === 'admin' ? 'error' : 'default'">
                    {{ user.role === 'admin' ? '管理员' : '用户' }}
                  </v-chip>
                </td>
                <td>
                  <v-chip size="small" :color="user.status === 'banned' ? 'error' : 'success'">
                    {{ user.status === 'banned' ? '已封禁' : '正常' }}
                  </v-chip>
                </td>
                <td>{{ formatDate(user.created_at) }}</td>
                <td>
                  <v-btn variant="text" size="small" color="primary" @click="showEditRoleDialog(user)" v-if="currentUserId && user.id !== currentUserId">
                    修改角色
                  </v-btn>
                  <v-btn variant="text" size="small" color="warning" @click="showBanDialog(user)" v-if="currentUserId && user.id !== currentUserId && user.status !== 'banned'">
                    封禁
                  </v-btn>
                  <v-btn variant="text" size="small" color="success" @click="handleUnban(user)" v-if="currentUserId && user.id !== currentUserId && user.status === 'banned'">
                    解封
                  </v-btn>
                  <v-btn variant="text" size="small" color="error" @click="handleDeleteUser(user)" v-if="currentUserId && user.id !== currentUserId">
                    删除
                  </v-btn>
                </td>
              </tr>
            </tbody>
          </v-table>
        </v-window-item>
        
        <!-- 文章管理 -->
        <v-window-item value="articles">
          <v-card variant="outlined" class="pa-4 mb-4">
            <v-row align="center">
              <v-col cols="12" md="4">
                <v-select
                  v-model="articleFilter"
                  :items="articleStatusOptions"
                  label="筛选状态"
                  variant="outlined"
                  density="compact"
                  hide-details
                ></v-select>
              </v-col>
              <v-col cols="12" md="8" class="text-right">
                <v-pagination
                  v-model="articlePage"
                  :length="articleTotalPages"
                  :total-visible="5"
                  density="compact"
                ></v-pagination>
              </v-col>
            </v-row>
          </v-card>
          
          <v-table>
            <thead>
              <tr>
                <th>ID</th>
                <th>标题</th>
                <th>作者</th>
                <th>分区</th>
                <th>点赞</th>
                <th>浏览</th>
                <th>状态</th>
                <th>发布时间</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="article in articles" :key="article.id">
                <td>{{ article.id }}</td>
                <td class="text-truncate" style="max-width: 200px;">{{ article.title }}</td>
                <td>{{ article.User?.display_name }}</td>
                <td>{{ article.Category?.name }}</td>
                <td>{{ article.like_count }}</td>
                <td>{{ article.view_count }}</td>
                <td>
                  <v-chip size="small" :color="getStatusColor(article.status)">
                    {{ getStatusText(article.status) }}
                  </v-chip>
                </td>
                <td>{{ formatDate(article.created_at) }}</td>
                <td>
                  <v-btn variant="text" size="small" color="primary" @click="showStatusDialog(article)">
                    修改状态
                  </v-btn>
                </td>
              </tr>
            </tbody>
          </v-table>
        </v-window-item>
        
        <!-- 评论管理 -->
        <v-window-item value="comments">
          <v-card variant="outlined" class="pa-4 mb-4">
            <v-row align="center">
              <v-col cols="12" md="6">
                <div class="text-body-2">评论总数：{{ commentTotal }}</div>
              </v-col>
              <v-col cols="12" md="6" class="text-right">
                <v-pagination
                  v-model="commentPage"
                  :length="commentTotalPages"
                  :total-visible="5"
                  density="compact"
                ></v-pagination>
              </v-col>
            </v-row>
          </v-card>
          
          <v-list lines="three">
            <v-list-item v-for="comment in allComments" :key="comment.id" class="px-0">
              <template v-slot:prepend>
                <v-avatar color="primary" size="40">
                  <v-img :src="comment.User?.avatar"></v-img>
                </v-avatar>
              </template>
              
              <v-list-item-title>
                {{ comment.User?.display_name }}
                <span class="text-caption text-medium-emphasis ml-2">回复文章：{{ comment.Article?.title }}</span>
              </v-list-item-title>
              <v-list-item-subtitle class="mt-1">
                {{ comment.content }}
              </v-list-item-subtitle>
              <v-list-item-subtitle class="mt-1">
                {{ formatDate(comment.created_at) }}
              </v-list-item-subtitle>
              
              <template v-slot:append>
                <v-btn variant="text" size="small" color="error" @click="handleDeleteComment(comment.id)">
                  删除
                </v-btn>
              </template>
            </v-list-item>
          </v-list>
        </v-window-item>
        
        <!-- 分区管理 -->
        <v-window-item value="categories">
          <v-card variant="outlined" class="pa-4 mb-4">
            <v-card-title class="text-subtitle-1 pa-0 mb-4">添加新分区</v-card-title>
            <v-form @submit.prevent="addCategory">
              <v-row>
                <v-col cols="12" md="4">
                  <v-text-field
                    v-model="categoryForm.name"
                    label="分区名称"
                    variant="outlined"
                    density="compact"
                    hide-details
                  ></v-text-field>
                </v-col>
                <v-col cols="12" md="4">
                  <v-text-field
                    v-model="categoryForm.description"
                    label="描述"
                    variant="outlined"
                    density="compact"
                    hide-details
                  ></v-text-field>
                </v-col>
                <v-col cols="12" md="2">
                  <v-btn type="submit" color="primary" block>添加</v-btn>
                </v-col>
              </v-row>
            </v-form>
          </v-card>
          
          <v-list>
            <v-list-item v-for="cat in categories" :key="cat.id">
              <template v-slot:prepend>
                <v-avatar color="primary" size="40">
                  <span>{{ cat.sort_order || 0 }}</span>
                </v-avatar>
              </template>
              
              <v-list-item-title class="font-weight-bold">{{ cat.name }}</v-list-item-title>
              <v-list-item-subtitle>{{ cat.description || '无描述' }}</v-list-item-subtitle>
              
              <template v-slot:append>
                <v-btn variant="text" size="small" color="primary" @click="showEditCategoryDialog(cat)">
                  编辑
                </v-btn>
                <v-btn variant="text" size="small" color="error" @click="handleDeleteCategory(cat.id)">
                  删除
                </v-btn>
              </template>
            </v-list-item>
          </v-list>
        </v-window-item>
        
        <!-- 侧边栏配置 -->
        <v-window-item value="sidebar">
          <v-card-text class="pa-0">
            <p class="text-body-2 text-medium-emphasis mb-4">
              配置侧边栏链接列表
            </p>
            
            <div v-for="(item, index) in sidebarItems" :key="index" class="d-flex gap-2 mb-3 align-center">
              <v-text-field
                v-model="item.title"
                label="标题"
                variant="outlined"
                density="compact"
                hide-details
                class="flex-grow-1"
              ></v-text-field>
              <v-text-field
                v-model="item.link"
                label="链接"
                variant="outlined"
                density="compact"
                hide-details
                class="flex-grow-1"
              ></v-text-field>
              <v-text-field
                v-model="item.icon"
                label="图标"
                variant="outlined"
                density="compact"
                hide-details
                style="max-width: 120px;"
              ></v-text-field>
              <v-btn icon variant="text" color="error" @click="removeSidebarItem(index)">
                <v-icon>mdi-delete</v-icon>
              </v-btn>
            </div>
            
            <div class="d-flex gap-2 mt-4">
              <v-btn variant="outlined" @click="addSidebarItem">
                <v-icon start>mdi-plus</v-icon>
                添加链接
              </v-btn>
              <v-btn color="primary" @click="saveSidebarConfig">
                <v-icon start>mdi-content-save</v-icon>
                保存配置
              </v-btn>
            </div>
          </v-card-text>
        </v-window-item>
        
        <!-- 删除审核 -->
        <v-window-item value="deletions">
          <div v-if="deletionRequests.length === 0" class="text-center pa-8 text-medium-emphasis">
            暂无待审核申请
          </div>
          
          <v-card v-for="req in deletionRequests" :key="req.id" class="mb-4 pa-4" variant="outlined">
            <v-card-text>
              <div class="d-flex justify-space-between align-start">
                <div>
                  <div class="text-h6 mb-2">{{ req.article?.title || '文章已删除' }}</div>
                  <v-list-item-subtitle class="mb-1">
                    <v-icon size="small">mdi-account</v-icon>
                    申请人：{{ req.user?.display_name }}
                  </v-list-item-subtitle>
                  <v-list-item-subtitle class="mb-1">
                    <v-icon size="small">mdi-delete</v-icon>
                    删除原因：{{ req.reason }}
                  </v-list-item-subtitle>
                  <v-list-item-subtitle>
                    <v-icon size="small">mdi-clock</v-icon>
                    申请时间：{{ formatDate(req.created_at) }}
                  </v-list-item-subtitle>
                </div>
                <div class="d-flex gap-2">
                  <v-btn color="primary" variant="flat" @click="approveDeletion(req.id)">
                    批准删除
                  </v-btn>
                  <v-btn color="secondary" variant="outlined" @click="rejectDeletion(req.id)">
                    拒绝
                  </v-btn>
                </div>
              </div>
            </v-card-text>
          </v-card>
        </v-window-item>
        
        <!-- 公告管理 -->
        <v-window-item value="announcement">
          <v-card-text class="pa-0">
            <v-textarea
              v-model="announcementContent"
              label="公告内容（支持Markdown）"
              variant="outlined"
              rows="10"
              placeholder="输入公告内容..."
              class="mb-4"
            ></v-textarea>
            
            <v-btn color="primary" @click="saveAnnouncement" class="mb-4">
              <v-icon start>mdi-content-save</v-icon>
              保存公告
            </v-btn>
          </v-card-text>
        </v-window-item>
        
        <!-- 通知管理 -->
        <v-window-item value="notifications">
          <v-card variant="outlined" class="pa-4 mb-4">
            <v-card-title class="text-subtitle-1 pa-0 mb-4">发送系统通知</v-card-title>
            <v-form @submit.prevent="sendNotification">
              <v-row>
                <v-col cols="12" md="4">
                  <v-select
                    v-model="notificationForm.type"
                    :items="notificationTypes"
                    label="通知类型"
                    variant="outlined"
                    density="compact"
                    hide-details
                  ></v-select>
                </v-col>
                <v-col cols="12" md="4">
                  <v-text-field
                    v-model="notificationForm.title"
                    label="通知标题"
                    variant="outlined"
                    density="compact"
                    hide-details
                  ></v-text-field>
                </v-col>
                <v-col cols="12" md="4">
                  <v-select
                    v-model="notificationForm.target"
                    :items="notificationTargets"
                    label="发送对象"
                    variant="outlined"
                    density="compact"
                    hide-details
                  ></v-select>
                </v-col>
                <v-col cols="12">
                  <v-textarea
                    v-model="notificationForm.content"
                    label="通知内容"
                    variant="outlined"
                    rows="4"
                    hide-details
                  ></v-textarea>
                </v-col>
                <v-col cols="12">
                  <v-btn type="submit" color="primary">发送通知</v-btn>
                </v-col>
              </v-row>
            </v-form>
          </v-card>
          
          <v-card variant="outlined" class="pa-4">
            <v-card-title class="text-subtitle-1 pa-0 mb-4">已发送通知</v-card-title>
            <v-list lines="three">
              <v-list-item v-for="notification in notifications" :key="notification.id" class="px-0">
                <v-list-item-title>
                  <v-chip size="small" :color="getNotificationColor(notification.type)" class="mr-2">
                    {{ getNotificationTypeText(notification.type) }}
                  </v-chip>
                  {{ notification.title }}
                </v-list-item-title>
                <v-list-item-subtitle class="mt-1">{{ notification.content }}</v-list-item-subtitle>
                <v-list-item-subtitle class="mt-1">
                  发送时间：{{ formatDate(notification.created_at) }}
                </v-list-item-subtitle>
                <template v-slot:append>
                  <v-btn variant="text" size="small" color="error" @click="handleDeleteNotification(notification.id)">
                    删除
                  </v-btn>
                </template>
              </v-list-item>
            </v-list>
          </v-card>
        </v-window-item>
      </v-window>
    </v-card>
    
    <!-- 角色修改对话框 -->
    <v-dialog v-model="editRoleDialog.show" max-width="400">
      <v-card>
        <v-card-title>修改用户角色</v-card-title>
        <v-card-text>
          <div class="mb-4">用户：{{ editRoleDialog.user?.display_name }}</div>
          <v-radio-group v-model="editRoleDialog.role" inline>
            <v-radio label="管理员" value="admin"></v-radio>
            <v-radio label="普通用户" value="user"></v-radio>
          </v-radio-group>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn variant="text" @click="editRoleDialog.show = false">取消</v-btn>
          <v-btn color="primary" @click="confirmEditRole">确定</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    
    <!-- 文章状态修改对话框 -->
    <v-dialog v-model="statusDialog.show" max-width="400">
      <v-card>
        <v-card-title>修改文章状态</v-card-title>
        <v-card-text>
          <div class="mb-4">文章：{{ statusDialog.article?.title }}</div>
          <v-radio-group v-model="statusDialog.status" inline>
            <v-radio label="已发布" value="published"></v-radio>
            <v-radio label="待审核" value="pending"></v-radio>
            <v-radio label="已删除" value="deleted"></v-radio>
          </v-radio-group>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn variant="text" @click="statusDialog.show = false">取消</v-btn>
          <v-btn color="primary" @click="confirmChangeStatus">确定</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    
    <!-- 分区编辑对话框 -->
    <v-dialog v-model="editCategoryDialog.show" max-width="500">
      <v-card>
        <v-card-title>编辑分区</v-card-title>
        <v-card-text>
          <v-text-field v-model="editCategoryDialog.name" label="分区名称" class="mb-3"></v-text-field>
          <v-text-field v-model="editCategoryDialog.description" label="描述" class="mb-3"></v-text-field>
          <v-text-field v-model="editCategoryDialog.sort_order" label="排序" type="number"></v-text-field>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn variant="text" @click="editCategoryDialog.show = false">取消</v-btn>
          <v-btn color="primary" @click="confirmEditCategory">确定</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    
    <!-- 封禁用户对话框 -->
    <v-dialog v-model="banDialog.show" max-width="500">
      <v-card>
        <v-card-title>封禁用户</v-card-title>
        <v-card-text>
          <div class="mb-4">用户：{{ banDialog.user?.display_name }}</div>
          <v-textarea
            v-model="banDialog.reason"
            label="封禁原因"
            variant="outlined"
            rows="3"
            hint="请输入封禁原因"
            persistent-hint
          ></v-textarea>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn variant="text" @click="banDialog.show = false">取消</v-btn>
          <v-btn color="error" @click="handleBan">确认封禁</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
  
  <v-alert v-else type="error" variant="tonal">
    无权限访问
  </v-alert>
</template>

<script>
import { ref, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import api from '../api'
import { confirm as showConfirm, success as showSuccess, error as showError } from '../utils/modal'

export default {
  name: 'Admin',
  setup() {
    const router = useRouter()
    const activeTab = ref('overview')
    const isAdmin = ref(false)
    const currentUserId = ref(null)
    const isInitialized = ref(false)
    
    // 统计数据
    const statistics = ref({
      user_count: 0,
      article_count: 0,
      comment_count: 0,
      view_count: 0
    })
    
    // 用户管理
    const users = ref([])
    const editRoleDialog = ref({
      show: false,
      user: null,
      role: 'user'
    })
    
    // 文章管理
    const articles = ref([])
    const articlePage = ref(1)
    const articleTotalPages = ref(1)
    const articleFilter = ref('')
    const articleStatusOptions = [
      { title: '全部', value: '' },
      { title: '已发布', value: 'published' },
      { title: '待审核', value: 'pending' },
      { title: '已删除', value: 'deleted' }
    ]
    const statusDialog = ref({
      show: false,
      article: null,
      status: 'published'
    })
    
    // 评论管理
    const allComments = ref([])
    const commentPage = ref(1)
    const commentTotal = ref(0)
    const commentTotalPages = ref(1)
    
    // 删除审核
    const deletionRequests = ref([])
    
    // 分区管理
    const categories = ref([])
    const categoryForm = ref({
      name: '',
      description: ''
    })
    const editCategoryDialog = ref({
      show: false,
      id: null,
      name: '',
      description: '',
      sort_order: 0
    })
    
    // 侧边栏配置
    const sidebarItems = ref([])
    
    // 公告
    const announcementContent = ref('')
    
    // 通知管理
    const notifications = ref([])
    const notificationForm = ref({
      type: 'system',
      title: '',
      content: '',
      target: 'all'
    })
    const notificationTypes = [
      { title: '系统通知', value: 'system' },
      { title: '活动公告', value: 'activity' },
      { title: '更新通知', value: 'update' },
      { title: '警告通知', value: 'warning' }
    ]
    const notificationTargets = [
      { title: '所有用户', value: 'all' },
      { title: '仅管理员', value: 'admin' }
    ]
    
    // 封禁对话框
    const banDialog = ref({
      show: false,
      user: null,
      reason: ''
    })
    
    const checkAdmin = async () => {
      const storedUser = localStorage.getItem('user')
      
      if (!storedUser) {
        isAdmin.value = false
        isInitialized.value = true
        router.push('/login')
        return
      }
      
      let user
      try {
        user = JSON.parse(storedUser)
      } catch (e) {
        console.error('解析用户信息失败', e)
        isAdmin.value = false
        isInitialized.value = true
        router.push('/login')
        return
      }
      
      if (!user || !user.id) {
        isAdmin.value = false
        isInitialized.value = true
        router.push('/login')
        return
      }
      
      if (user.role === 'admin') {
        isAdmin.value = true
        currentUserId.value = user.id
        isInitialized.value = true
        return
      }
      
      try {
        const response = await api.get('/profile')
        const profile = response.data
        localStorage.setItem('user', JSON.stringify(profile))
        
        if (profile.role === 'admin') {
          isAdmin.value = true
          currentUserId.value = profile.id
        } else {
          isAdmin.value = false
          router.push('/')
        }
      } catch (error) {
        console.error('获取用户信息失败', error)
        isAdmin.value = false
        router.push('/login')
      }
      
      isInitialized.value = true
    }
    
    const loadStatistics = async () => {
      try {
        const response = await api.get('/admin/statistics')
        statistics.value = response.data
      } catch (error) {
        console.error('加载统计失败', error)
      }
    }
    
    const loadUsers = async () => {
      try {
        const response = await api.get('/admin/users')
        users.value = response.data.users
      } catch (error) {
        console.error('加载用户失败', error)
      }
    }
    
    const loadArticles = async () => {
      try {
        const params = { page: articlePage.value, page_size: 20 }
        if (articleFilter.value) {
          params.status = articleFilter.value
        }
        const response = await api.get('/admin/articles', { params })
        articles.value = response.data.articles
        articleTotalPages.value = response.data.total_pages
      } catch (error) {
        console.error('加载文章失败', error)
      }
    }
    
    const loadComments = async () => {
      try {
        const response = await api.get('/admin/comments', { params: { page: commentPage.value, page_size: 20 } })
        allComments.value = response.data.comments
        commentTotal.value = response.data.total
        commentTotalPages.value = response.data.total_pages
      } catch (error) {
        console.error('加载评论失败', error)
      }
    }
    
    const loadDeletionRequests = async () => {
      try {
        const response = await api.get('/deletion-requests')
        deletionRequests.value = response.data.requests || []
      } catch (error) {
        console.error('加载删除申请失败', error)
        deletionRequests.value = []
      }
    }
    
    const loadCategories = async () => {
      try {
        const response = await api.get('/categories')
        categories.value = response.data.categories || []
      } catch (error) {
        console.error('加载分区失败', error)
      }
    }
    
    const loadSidebarConfig = async () => {
      try {
        const response = await api.get('/sidebar-config')
        sidebarItems.value = response.data.items || []
      } catch (error) {
        console.error('加载侧边栏配置失败', error)
      }
    }
    
    const loadAnnouncement = async () => {
      try {
        const response = await api.get('/announcement')
        announcementContent.value = response.data.content || ''
      } catch (error) {
        console.error('加载公告失败', error)
      }
    }
    
    const loadNotifications = async () => {
      try {
        const response = await api.get('/notifications/admin')
        notifications.value = response.data.notifications || []
      } catch (error) {
        console.error('加载通知失败', error)
      }
    }
    
    // 用户操作
    const showEditRoleDialog = (user) => {
      editRoleDialog.value = {
        show: true,
        user: user,
        role: user.role
      }
    }
    
    const confirmEditRole = async () => {
      try {
        await api.put(`/admin/users/${editRoleDialog.value.user.id}/role`, { role: editRoleDialog.value.role })
        showSuccess('更新成功')
        editRoleDialog.value.show = false
        loadUsers()
      } catch (error) {
        console.error('更新角色失败', error)
        showError(error.response?.data?.error || '更新失败')
      }
    }
    
    const handleDeleteUser = async (user) => {
      const confirmed = await showConfirm(`确定要删除用户 "${user.display_name}" 吗？`, {
        title: '确认删除',
        icon: 'mdi-alert-circle',
        iconColor: 'error'
      })
      if (!confirmed) return
      
      try {
        await api.delete(`/admin/users/${user.id}`)
        showSuccess('删除成功')
        loadUsers()
      } catch (error) {
        console.error('删除用户失败', error)
        showError(error.response?.data?.error || '删除失败')
      }
    }
    
    const showBanDialog = (user) => {
      banDialog.value = {
        show: true,
        user: user,
        reason: ''
      }
    }
    
    const handleBan = async () => {
      if (!banDialog.value.reason.trim()) {
        showError('请输入封禁原因')
        return
      }
      
      try {
        await api.post(`/admin/users/${banDialog.value.user.id}/ban`, { reason: banDialog.value.reason })
        showSuccess('封禁成功')
        banDialog.value.show = false
        loadUsers()
      } catch (error) {
        console.error('封禁用户失败', error)
        showError(error.response?.data?.error || '封禁失败')
      }
    }
    
    const handleUnban = async (user) => {
      const confirmed = await showConfirm(`确定要解封用户 "${user.display_name}" 吗？`)
      if (!confirmed) return
      
      try {
        await api.post(`/admin/users/${user.id}/unban`)
        showSuccess('解封成功')
        loadUsers()
      } catch (error) {
        console.error('解封用户失败', error)
        showError(error.response?.data?.error || '解封失败')
      }
    }
    
    // 文章操作
    const getStatusColor = (status) => {
      const colors = { published: 'success', pending: 'warning', deleted: 'error' }
      return colors[status] || 'default'
    }
    
    const getStatusText = (status) => {
      const texts = { published: '已发布', pending: '待审核', deleted: '已删除' }
      return texts[status] || status
    }
    
    const showStatusDialog = (article) => {
      statusDialog.value = {
        show: true,
        article: article,
        status: article.status
      }
    }
    
    const confirmChangeStatus = async () => {
      try {
        await api.put(`/admin/articles/${statusDialog.value.article.id}/status`, { status: statusDialog.value.status })
        showSuccess('更新成功')
        statusDialog.value.show = false
        loadArticles()
      } catch (error) {
        console.error('更新状态失败', error)
        showError(error.response?.data?.error || '更新失败')
      }
    }
    
    // 评论操作
    const handleDeleteComment = async (commentId) => {
      const confirmed = await showConfirm('确定要删除这条评论吗？', {
        title: '确认删除',
        icon: 'mdi-alert-circle',
        iconColor: 'error'
      })
      if (!confirmed) return
      
      try {
        await api.delete(`/admin/comments/${commentId}`)
        showSuccess('删除成功')
        loadComments()
      } catch (error) {
        console.error('删除评论失败', error)
        showError('删除失败')
      }
    }
    
    // 删除审核操作
    const approveDeletion = async (id) => {
      const confirmed = await showConfirm('确定要批准删除这篇文章吗？', {
        title: '确认批准',
        icon: 'mdi-check-circle',
        iconColor: 'success'
      })
      if (!confirmed) return
      
      try {
        await api.post(`/deletion-requests/${id}/approve`)
        showSuccess('已批准删除')
        loadDeletionRequests()
      } catch (error) {
        console.error('批准删除失败', error)
        showError('操作失败')
      }
    }
    
    const rejectDeletion = async (id) => {
      const confirmed = await showConfirm('确定要拒绝这个删除申请吗？', {
        title: '确认拒绝',
        icon: 'mdi-close-circle',
        iconColor: 'warning'
      })
      if (!confirmed) return
      
      try {
        await api.post(`/deletion-requests/${id}/reject`)
        showSuccess('已拒绝删除')
        loadDeletionRequests()
      } catch (error) {
        console.error('拒绝删除失败', error)
        showError('操作失败')
      }
    }
    
    // 分区操作
    const addCategory = async () => {
      if (!categoryForm.value.name) {
        showError('请输入分区名称')
        return
      }
      
      try {
        await api.post('/categories', categoryForm.value)
        showSuccess('添加成功')
        categoryForm.value = { name: '', description: '' }
        loadCategories()
      } catch (error) {
        console.error('添加分区失败', error)
        showError('添加失败')
      }
    }
    
    const showEditCategoryDialog = (cat) => {
      editCategoryDialog.value = {
        show: true,
        id: cat.id,
        name: cat.name,
        description: cat.description || '',
        sort_order: cat.sort_order || 0
      }
    }
    
    const confirmEditCategory = async () => {
      try {
        await api.put(`/categories/${editCategoryDialog.value.id}`, {
          name: editCategoryDialog.value.name,
          description: editCategoryDialog.value.description,
          sort_order: parseInt(editCategoryDialog.value.sort_order) || 0
        })
        showSuccess('更新成功')
        editCategoryDialog.value.show = false
        loadCategories()
      } catch (error) {
        console.error('更新失败', error)
        showError('更新失败')
      }
    }
    
    const handleDeleteCategory = async (id) => {
      const confirmed = await showConfirm('确定要删除这个分区吗？', {
        title: '确认删除',
        icon: 'mdi-alert-circle',
        iconColor: 'error'
      })
      if (!confirmed) return
      
      try {
        await api.delete(`/categories/${id}`)
        showSuccess('删除成功')
        loadCategories()
      } catch (error) {
        console.error('删除分区失败', error)
        showError(error.response?.data?.error || '删除失败')
      }
    }
    
    // 侧边栏操作
    const addSidebarItem = () => {
      sidebarItems.value.push({ title: '', link: '', icon: '' })
    }
    
    const removeSidebarItem = (index) => {
      sidebarItems.value.splice(index, 1)
    }
    
    const saveSidebarConfig = async () => {
      try {
        await api.put('/sidebar-config', { items: sidebarItems.value })
        showSuccess('保存成功')
      } catch (error) {
        console.error('保存失败', error)
        showError('保存失败')
      }
    }
    
    // 公告操作
    const saveAnnouncement = async () => {
      try {
        await api.put('/announcement', { content: announcementContent.value })
        showSuccess('保存成功')
      } catch (error) {
        console.error('保存公告失败', error)
        showError('保存失败')
      }
    }
    
    // 通知操作
    const getNotificationColor = (type) => {
      const colors = { system: 'primary', activity: 'success', update: 'info', warning: 'warning' }
      return colors[type] || 'default'
    }
    
    const getNotificationTypeText = (type) => {
      const texts = { system: '系统通知', activity: '活动公告', update: '更新通知', warning: '警告通知' }
      return texts[type] || type
    }
    
    const sendNotification = async () => {
      if (!notificationForm.value.title || !notificationForm.value.content) {
        showError('请填写通知标题和内容')
        return
      }
      
      try {
        await api.post('/notifications', notificationForm.value)
        showSuccess('通知发送成功')
        notificationForm.value = { type: 'system', title: '', content: '', target: 'all' }
        loadNotifications()
      } catch (error) {
        console.error('发送通知失败', error)
        showError('发送失败')
      }
    }
    
    const handleDeleteNotification = async (id) => {
      const confirmed = await showConfirm('确定要删除这条通知吗？', {
        title: '确认删除',
        icon: 'mdi-alert-circle',
        iconColor: 'error'
      })
      if (!confirmed) return
      
      try {
        await api.delete(`/notifications/${id}`)
        showSuccess('删除成功')
        loadNotifications()
      } catch (error) {
        console.error('删除通知失败', error)
        showError('删除失败')
      }
    }
    
    const formatDate = (date) => {
      if (!date) return ''
      return new Date(date).toLocaleString('zh-CN')
    }
    
    // 监听分页和筛选变化
    watch([articlePage, articleFilter], () => {
      if (activeTab.value === 'articles') {
        loadArticles()
      }
    })
    
    watch(commentPage, () => {
      if (activeTab.value === 'comments') {
        loadComments()
      }
    })
    
    watch(activeTab, (newTab) => {
      if (newTab === 'overview') loadStatistics()
      if (newTab === 'users') loadUsers()
      if (newTab === 'articles') loadArticles()
      if (newTab === 'comments') loadComments()
      if (newTab === 'deletions') loadDeletionRequests()
      if (newTab === 'categories') loadCategories()
      if (newTab === 'sidebar') loadSidebarConfig()
      if (newTab === 'announcement') loadAnnouncement()
      if (newTab === 'notifications') loadNotifications()
    })
    
    onMounted(async () => {
      await checkAdmin()
      if (isAdmin.value) {
        loadStatistics()
        loadDeletionRequests()
        loadCategories()
        loadSidebarConfig()
        loadAnnouncement()
      }
    })
    
    return {
      activeTab,
      isAdmin,
      isInitialized,
      currentUserId,
      statistics,
      users,
      editRoleDialog,
      articles,
      articlePage,
      articleTotalPages,
      articleFilter,
      articleStatusOptions,
      statusDialog,
      allComments,
      commentPage,
      commentTotal,
      commentTotalPages,
      deletionRequests,
      categories,
      categoryForm,
      editCategoryDialog,
      sidebarItems,
      announcementContent,
      notifications,
      notificationForm,
      notificationTypes,
      notificationTargets,
      banDialog,
      showEditRoleDialog,
      confirmEditRole,
      handleDeleteUser,
      showBanDialog,
      handleBan,
      handleUnban,
      getStatusColor,
      getStatusText,
      showStatusDialog,
      confirmChangeStatus,
      handleDeleteComment,
      approveDeletion,
      rejectDeletion,
      addCategory,
      showEditCategoryDialog,
      confirmEditCategory,
      handleDeleteCategory,
      addSidebarItem,
      removeSidebarItem,
      saveSidebarConfig,
      saveAnnouncement,
      getNotificationColor,
      getNotificationTypeText,
      sendNotification,
      handleDeleteNotification,
      formatDate
    }
  }
}
</script>

<style scoped>
</style>
