<template>
  <div class="admin-container" v-if="isAdmin">
    <h2>管理后台</h2>
    
    <div class="admin-tabs">
      <button @click="activeTab = 'deletions'" :class="{ active: activeTab === 'deletions' }">
        删除审核 ({{ deletionRequests.length }})
      </button>
      <button @click="activeTab = 'categories'" :class="{ active: activeTab === 'categories' }">
        分区管理
      </button>
      <button @click="activeTab = 'sidebar'" :class="{ active: activeTab === 'sidebar' }">
        侧边栏配置
      </button>
      <button @click="activeTab = 'announcement'" :class="{ active: activeTab === 'announcement' }">
        公告管理
      </button>
    </div>
    
    <!-- 删除审核 -->
    <div v-if="activeTab === 'deletions'" class="tab-content">
      <h3>文章删除申请</h3>
      <div v-if="deletionRequests.length === 0" class="empty">暂无待审核申请</div>
      <div v-for="req in deletionRequests" :key="req.id" class="deletion-item">
        <div class="deletion-info">
          <h4>{{ req.article.title }}</h4>
          <p>申请人：{{ req.user.display_name }}</p>
          <p>删除原因：{{ req.reason }}</p>
          <p>申请时间：{{ formatDate(req.created_at) }}</p>
        </div>
        <div class="deletion-actions">
          <button @click="approveDeletion(req.id)" class="btn btn-primary">批准删除</button>
          <button @click="rejectDeletion(req.id)" class="btn btn-secondary">拒绝</button>
        </div>
      </div>
    </div>
    
    <!-- 分区管理 -->
    <div v-if="activeTab === 'categories'" class="tab-content">
      <h3>分区管理</h3>
      <div class="category-form">
        <input type="text" v-model="categoryForm.name" placeholder="分区名称">
        <input type="text" v-model="categoryForm.description" placeholder="描述">
        <input type="number" v-model="categoryForm.sort_order" placeholder="排序">
        <button @click="addCategory" class="btn btn-primary">添加分区</button>
      </div>
      
      <div class="category-list">
        <div v-for="cat in categories" :key="cat.id" class="category-item">
          <div class="category-info">
            <strong>{{ cat.name }}</strong>
            <span>{{ cat.description }}</span>
            <span>排序：{{ cat.sort_order }}</span>
          </div>
          <div class="category-actions">
            <button @click="editCategory(cat)" class="btn btn-secondary">编辑</button>
            <button @click="deleteCategory(cat.id)" class="btn btn-danger">删除</button>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 侧边栏配置 -->
    <div v-if="activeTab === 'sidebar'" class="tab-content">
      <h3>侧边栏配置</h3>
      <p class="info">配置侧边栏链接列表（支持JSON格式）</p>
      
      <div class="sidebar-items">
        <div v-for="(item, index) in sidebarItems" :key="index" class="sidebar-item">
          <input v-model="item.title" placeholder="标题">
          <input v-model="item.link" placeholder="链接">
          <input v-model="item.icon" placeholder="图标(emoji)">
          <button @click="removeSidebarItem(index)" class="btn-danger">删除</button>
        </div>
        <button @click="addSidebarItem" class="btn btn-secondary">添加链接</button>
        <button @click="saveSidebarConfig" class="btn btn-primary" style="margin-top: 20px;">保存配置</button>
      </div>
    </div>
    
    <!-- 公告管理 -->
    <div v-if="activeTab === 'announcement'" class="tab-content">
      <h3>公告管理</h3>
      <div class="form-group">
        <label>公告内容（支持Markdown）</label>
        <textarea v-model="announcementContent" rows="10" placeholder="输入公告内容..."></textarea>
      </div>
      <button @click="saveAnnouncement" class="btn btn-primary">保存公告</button>
      
      <div class="preview" v-if="announcementContent">
        <h4>预览</h4>
        <div class="markdown-body" v-html="previewHtml"></div>
      </div>
    </div>
  </div>
  
  <div v-else class="error-message">无权限访问</div>
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

<style scoped>
.admin-container {
  background: white;
  border-radius: 12px;
  padding: 30px;
}

.admin-tabs {
  display: flex;
  gap: 10px;
  border-bottom: 1px solid #e5e7eb;
  margin-bottom: 30px;
  padding-bottom: 10px;
}

.admin-tabs button {
  padding: 8px 20px;
  background: none;
  border: none;
  cursor: pointer;
  font-size: 14px;
  color: #6b7280;
}

.admin-tabs button.active {
  color: #10b981;
  border-bottom: 2px solid #10b981;
}

.tab-content {
  padding: 20px 0;
}

.deletion-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  margin-bottom: 15px;
}

.deletion-info h4 {
  margin-bottom: 10px;
}

.deletion-info p {
  margin-bottom: 5px;
  color: #6b7280;
  font-size: 14px;
}

.deletion-actions {
  display: flex;
  gap: 10px;
}

.category-form {
  display: flex;
  gap: 10px;
  margin-bottom: 20px;
  flex-wrap: wrap;
}

.category-form input {
  flex: 1;
  padding: 8px;
  border: 1px solid #d1d5db;
  border-radius: 4px;
}

.category-form input[type="number"] {
  width: 80px;
}

.category-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px;
  border-bottom: 1px solid #e5e7eb;
}

.category-info {
  display: flex;
  gap: 20px;
  align-items: center;
}

.category-actions {
  display: flex;
  gap: 10px;
}

.sidebar-items {
  max-width: 600px;
}

.sidebar-item {
  display: flex;
  gap: 10px;
  margin-bottom: 10px;
  align-items: center;
}

.sidebar-item input {
  flex: 1;
  padding: 8px;
  border: 1px solid #d1d5db;
  border-radius: 4px;
}

.preview {
  margin-top: 30px;
  padding: 20px;
  background: #f9fafb;
  border-radius: 8px;
}

.info {
  color: #6b7280;
  margin-bottom: 15px;
  font-size: 14px;
}

.empty {
  text-align: center;
  padding: 40px;
  color: #6b7280;
}
</style>