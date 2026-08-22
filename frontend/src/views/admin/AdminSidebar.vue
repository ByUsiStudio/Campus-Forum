<template>
  <div class="page-container">
    <!-- 添加新侧边栏项表单 -->
    <el-card class="mb-4" shadow="never">
      <template #header>
        <span class="card-title">
          <el-icon :size="20" class="mr-1"><Plus /></el-icon>
          添加新侧边栏项
        </span>
      </template>
      <el-row :gutter="16">
        <el-col :xs="24" :sm="10">
          <el-input
            v-model="newItem.title"
            placeholder="例如：联系我们"
            clearable
          />
        </el-col>
        <el-col :xs="24" :sm="10">
          <el-input
            v-model="newItem.url"
            placeholder="例如：https://example.com"
            clearable
          />
        </el-col>
        <el-col :xs="24" :sm="4">
          <el-button type="primary" style="width: 100%" @click="addSidebarItem">
            <el-icon class="mr-1"><Plus /></el-icon>
            添加
          </el-button>
        </el-col>
      </el-row>
    </el-card>

    <!-- 侧边栏项列表 -->
    <el-card shadow="never">
      <template #header>
        <span class="card-title">
          <el-icon :size="20" class="mr-1"><Monitor /></el-icon>
          侧边栏配置
        </span>
      </template>

      <el-empty
        v-if="sidebarItems.length === 0"
        description="暂无侧边栏项"
        :image-size="64"
      />

      <el-table
        v-else
        :data="sidebarItems"
        row-key="id"
        style="width: 100%"
      >
        <el-table-column label="标题" min-width="160">
          <template #default="{ row }">
            <span class="font-weight-medium">{{ row.title }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="url" label="链接地址" min-width="220" show-overflow-tooltip />
        <el-table-column label="操作" width="100" align="center">
          <template #default="{ row, $index }">
            <el-button
              type="danger"
              size="small"
              link
              @click="removeSidebarItem($index)"
            >
              <el-icon><Delete /></el-icon>
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <template v-if="sidebarItems.length > 0" #footer>
        <div class="card-footer">
          <el-button
            type="primary"
            :loading="saving"
            @click="saveSidebarConfig"
          >
            <el-icon class="mr-1"><Check /></el-icon>
            保存配置
          </el-button>
        </div>
      </template>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Plus, Monitor, Delete, Check } from '@element-plus/icons-vue'
import { adminSidebarApi } from '../../api/admin'
import { success, error } from '@/utils/message'

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
.mb-4 {
  margin-bottom: 16px;
}
.mr-1 {
  margin-right: 4px;
}
.page-container {
  padding: 16px;
}
.card-title {
  display: inline-flex;
  align-items: center;
  font-weight: 600;
}
.font-weight-medium {
  font-weight: 500;
}
.card-footer {
  display: flex;
  justify-content: flex-end;
  padding-top: 16px;
  border-top: 1px solid var(--el-border-color-lighter);
}
</style>
