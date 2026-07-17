<template>
  <div class="admin-sidebar">
    <div class="layui-card mb-4">
      <div class="layui-card-header">
        <i class="fa-solid fa-circle-plus"></i>
        <span>添加新侧边栏项</span>
      </div>
      <div class="layui-card-body">
        <div class="form-row">
          <div class="form-item">
            <input
              v-model="newItem.title"
              type="text"
              placeholder="例如：联系我们"
              class="layui-input"
            />
          </div>
          <div class="form-item">
            <input
              v-model="newItem.url"
              type="text"
              placeholder="例如：https://example.com"
              class="layui-input"
            />
          </div>
          <div class="form-item btn-item">
            <button class="layui-btn" @click="addSidebarItem">
              <i class="fa-solid fa-plus"></i>添加
            </button>
          </div>
        </div>
      </div>
    </div>

    <div class="layui-card">
      <div class="layui-card-header">
        <i class="fa-solid fa-globe"></i>
        <span>侧边栏配置</span>
      </div>

      <div v-if="sidebarItems.length > 0" class="item-list">
        <div v-for="(item, index) in sidebarItems" :key="item.id || index" class="item-row">
          <div class="item-avatar">
            <i class="fa-solid fa-link"></i>
          </div>
          <div class="item-content">
            <div class="item-title">{{ item.title }}</div>
            <div class="item-url">{{ item.url }}</div>
          </div>
          <button class="delete-btn" @click="removeSidebarItem(index)">
            <i class="fa-solid fa-trash"></i>
          </button>
        </div>
      </div>

      <div v-else class="empty-state">
        <i class="fa-solid fa-inbox"></i>
        <div class="empty-text">暂无侧边栏项</div>
      </div>

      <div v-if="sidebarItems.length > 0" class="layui-card-footer">
        <button class="layui-btn" @click="saveSidebarConfig" :disabled="saving">
          <i class="fa-solid fa-check"></i>保存配置
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { adminSidebarApi } from '../../api/admin'
import { success, error } from '../../utils/modal'

const sidebarItems = ref([])
const newItem = ref({
  title: '',
  url: ''
})
const saving = ref(false)

const loadSidebarConfig = async () => {
  try {
    const response = await adminSidebarApi.getConfig()
    sidebarItems.value = response.data.items || []
  } catch (err) {
    console.error('加载侧边栏配置失败', err)
    error('加载侧边栏配置失败')
  }
}

const addSidebarItem = () => {
  if (!newItem.value.title || !newItem.value.url) {
    error('请填写标题和链接地址')
    return
  }
  
  sidebarItems.value.push({
    id: Date.now(),
    title: newItem.value.title,
    url: newItem.value.url
  })
  
  newItem.value = {
    title: '',
    url: ''
  }
}

const removeSidebarItem = (index) => {
  sidebarItems.value.splice(index, 1)
}

const saveSidebarConfig = async () => {
  saving.value = true
  try {
    await adminSidebarApi.updateConfig({ items: sidebarItems.value })
    success('保存成功')
  } catch (err) {
    console.error('保存侧边栏配置失败', err)
    error(err.response?.data?.error || '保存失败')
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  loadSidebarConfig()
})
</script>

<style scoped>
.admin-sidebar {
  padding: 20px;
}

.layui-card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 18px;
  font-weight: 600;
}

.form-row {
  display: flex;
  gap: 12px;
}

.form-item {
  flex: 1;

  &.btn-item {
    flex: none;
  }
}

.item-list {
  padding: 8px 0;
}

.item-row {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px 20px;
  border-bottom: 1px solid #f0f0f0;

  &:last-child {
    border-bottom: none;
  }
}

.item-avatar {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  background: rgba(30, 159, 255, 0.1);
  color: #1E9FFF;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  flex-shrink: 0;
}

.item-content {
  flex: 1;
  min-width: 0;
}

.item-title {
  font-size: 15px;
  font-weight: 500;
  color: #333;
  margin-bottom: 4px;
}

.item-url {
  font-size: 13px;
  color: #999;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.delete-btn {
  padding: 6px;
  background: transparent;
  border: none;
  color: #999;
  font-size: 16px;
  cursor: pointer;

  &:hover {
    color: #FF5722;
  }
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
}

.empty-state i {
  font-size: 48px;
  color: #e0e0e0;
  margin-bottom: 16px;
}

.empty-text {
  font-size: 15px;
  color: #999;
}

.layui-card-footer {
  display: flex;
  justify-content: flex-end;
  padding: 16px 20px;
}

@media (max-width: 768px) {
  .form-row {
    flex-direction: column;
  }
}
</style>