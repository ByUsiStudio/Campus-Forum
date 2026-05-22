<template>
  <div v-if="isAdmin">
    <v-card class="pa-6">
      <v-card-title class="text-h5 mb-4" style="color: rgb(var(--v-theme-primary));">
        管理后台
      </v-card-title>
      
      <v-tabs v-model="activeTab" color="primary" class="mb-4">
        <v-tab value="deletions">
          <v-badge :content="deletionRequests.length" color="error" :model-value="deletionRequests.length > 0">
            删除审核
          </v-badge>
        </v-tab>
        <v-tab value="categories">分区管理</v-tab>
        <v-tab value="sidebar">侧边栏配置</v-tab>
        <v-tab value="announcement">公告管理</v-tab>
      </v-tabs>
      
      <v-window v-model="activeTab">
        <!-- 删除审核 -->
        <v-window-item value="deletions">
          <div v-if="deletionRequests.length === 0" class="text-center pa-8 text-medium-emphasis">
            暂无待审核申请
          </div>
          
          <v-card v-for="req in deletionRequests" :key="req.id" class="mb-4 pa-4" variant="outlined">
            <v-card-text>
              <div class="d-flex justify-space-between align-start">
                <div>
                  <div class="text-h6 mb-2">{{ req.article.title }}</div>
                  <v-list-item-subtitle class="mb-1">
                    <v-icon size="small">mdi-account</v-icon>
                    申请人：{{ req.user.display_name }}
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
                  <v-text-field
                    v-model="categoryForm.sort_order"
                    label="排序"
                    type="number"
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
                  <span>{{ cat.sort_order }}</span>
                </v-avatar>
              </template>
              
              <v-list-item-title class="font-weight-bold">{{ cat.name }}</v-list-item-title>
              <v-list-item-subtitle>{{ cat.description || '无描述' }}</v-list-item-subtitle>
              
              <template v-slot:append>
                <v-btn variant="text" size="small" color="primary" @click="editCategory(cat)">
                  编辑
                </v-btn>
                <v-btn variant="text" size="small" color="error" @click="deleteCategory(cat.id)">
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
            
            <v-card v-if="announcementContent" variant="outlined" class="pa-4">
              <v-card-title class="text-subtitle-1 pa-0 mb-2">预览</v-card-title>
              <div class="markdown-body" v-html="previewHtml"></div>
            </v-card>
          </v-card-text>
        </v-window-item>
      </v-window>
    </v-card>
  </div>
  
  <v-alert v-else type="error" variant="tonal">
    无权限访问
  </v-alert>
</template>

<script>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import api from '../api'
import MarkdownIt from 'markdown-it'

const md = new MarkdownIt({
  html: true,
  linkify: true
})

export default {
  name: 'Admin',
  setup() {
    const router = useRouter()
    const activeTab = ref('deletions')
    const deletionRequests = ref([])
    const categories = ref([])
    const sidebarItems = ref([])
    const announcementContent = ref('')
    const isAdmin = ref(false)
    
    const categoryForm = ref({
      name: '',
      description: '',
      sort_order: 0
    })
    
    const previewHtml = computed(() => {
      return md.render(announcementContent.value)
    })
    
    const checkAdmin = () => {
      const user = JSON.parse(localStorage.getItem('user') || '{}')
      if (user.role !== 'admin') {
        isAdmin.value = false
        router.push('/')
      } else {
        isAdmin.value = true
      }
    }
    
    const loadDeletionRequests = async () => {
      try {
        const response = await api.get('/deletion-requests')
        deletionRequests.value = response.data.requests
      } catch (error) {
        console.error('加载删除申请失败', error)
      }
    }
    
    const loadCategories = async () => {
      try {
        const response = await api.get('/categories')
        categories.value = response.data.categories
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
    
    const approveDeletion = async (id) => {
      try {
        await api.post(`/deletion-requests/${id}/approve`)
        alert('已批准删除')
        loadDeletionRequests()
      } catch (error) {
        console.error('批准删除失败', error)
        alert('操作失败')
      }
    }
    
    const rejectDeletion = async (id) => {
      try {
        await api.post(`/deletion-requests/${id}/reject`)
        alert('已拒绝删除')
        loadDeletionRequests()
      } catch (error) {
        console.error('拒绝删除失败', error)
        alert('操作失败')
      }
    }
    
    const addCategory = async () => {
      if (!categoryForm.value.name) {
        alert('请输入分区名称')
        return
      }
      
      try {
        await api.post('/categories', categoryForm.value)
        categoryForm.value = { name: '', description: '', sort_order: 0 }
        loadCategories()
        alert('添加成功')
      } catch (error) {
        console.error('添加分区失败', error)
        alert('添加失败')
      }
    }
    
    const editCategory = (cat) => {
      const newName = prompt('输入新名称', cat.name)
      const newDesc = prompt('输入新描述', cat.description)
      const newOrder = prompt('输入排序', cat.sort_order)
      
      if (newName) {
        api.put(`/categories/${cat.id}`, {
          name: newName,
          description: newDesc,
          sort_order: parseInt(newOrder) || 0
        }).then(() => {
          loadCategories()
          alert('更新成功')
        }).catch(err => {
          console.error('更新失败', err)
          alert('更新失败')
        })
      }
    }
    
    const deleteCategory = async (id) => {
      if (!confirm('确定要删除这个分区吗？')) return
      
      try {
        await api.delete(`/categories/${id}`)
        loadCategories()
        alert('删除成功')
      } catch (error) {
        console.error('删除分区失败', error)
        alert(error.response?.data?.error || '删除失败')
      }
    }
    
    const addSidebarItem = () => {
      sidebarItems.value.push({ title: '', link: '', icon: '' })
    }
    
    const removeSidebarItem = (index) => {
      sidebarItems.value.splice(index, 1)
    }
    
    const saveSidebarConfig = async () => {
      try {
        await api.put('/sidebar-config', { items: sidebarItems.value })
        alert('保存成功')
      } catch (error) {
        console.error('保存失败', error)
        alert('保存失败')
      }
    }
    
    const saveAnnouncement = async () => {
      try {
        await api.put('/announcement', { content: announcementContent.value })
        alert('保存成功')
      } catch (error) {
        console.error('保存公告失败', error)
        alert('保存失败')
      }
    }
    
    const formatDate = (date) => {
      return new Date(date).toLocaleString('zh-CN')
    }
    
    onMounted(() => {
      checkAdmin()
      loadDeletionRequests()
      loadCategories()
      loadSidebarConfig()
      loadAnnouncement()
    })
    
    return {
      activeTab,
      deletionRequests,
      categories,
      sidebarItems,
      announcementContent,
      isAdmin,
      categoryForm,
      previewHtml,
      approveDeletion,
      rejectDeletion,
      addCategory,
      editCategory,
      deleteCategory,
      addSidebarItem,
      removeSidebarItem,
      saveSidebarConfig,
      saveAnnouncement,
      formatDate
    }
  }
}
</script>