<template>
  <v-container fluid class="pa-6">
    <v-card>
      <v-card-title>
        <v-icon class="mr-2">mdi-web</v-icon>
        侧边栏配置
      </v-card-title>
      <v-card-text>
        <div class="mb-4">
          <v-btn color="primary" variant="flat" @click="addSidebarItem">
            <v-icon class="mr-1">mdi-plus</v-icon>
            添加侧边栏项
          </v-btn>
        </div>
        
        <v-list v-if="sidebarItems.length > 0">
          <v-list-item
            v-for="(item, index) in sidebarItems"
            :key="item.id || index"
            class="align-center"
          >
            <v-list-item-content>
              <v-list-item-title>{{ item.title }}</v-list-item-title>
              <v-list-item-subtitle>{{ item.url }}</v-list-item-subtitle>
            </v-list-item-content>
            <v-list-item-action>
              <v-btn
                variant="text"
                color="error"
                size="small"
                @click="removeSidebarItem(index)"
              >
                <v-icon>mdi-delete</v-icon>
              </v-btn>
            </v-list-item-action>
          </v-list-item>
        </v-list>
        
        <div v-else class="text-center text-gray-400 py-8">
          <v-icon size="48" color="grey lighten-3">mdi-inbox</v-icon>
          <p class="mt-2">暂无侧边栏项</p>
        </div>
      </v-card-text>
      <v-card-actions>
        <v-btn color="primary" variant="flat" @click="saveSidebarConfig">
          保存配置
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-container>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { adminSidebarApi } from '../../api/admin'
import { success, error } from '../../utils/modal'

const sidebarItems = ref([])

const loadSidebarConfig = async () => {
  try {
    const response = await adminSidebarApi.getConfig()
    sidebarItems.value = response.data.items || []
  } catch (error) {
    console.error('加载侧边栏配置失败', error)
  }
}

const addSidebarItem = () => {
  const title = prompt('请输入侧边栏项标题：')
  if (!title) return
  
  const url = prompt('请输入链接地址：')
  if (!url) return
  
  sidebarItems.value.push({
    id: Date.now(),
    title,
    url
  })
}

const removeSidebarItem = (index) => {
  sidebarItems.value.splice(index, 1)
}

const saveSidebarConfig = async () => {
  try {
    await adminSidebarApi.updateConfig({ items: sidebarItems.value })
    success('保存成功')
  } catch (error) {
    console.error('保存侧边栏配置失败', error)
    error(error.response?.data?.error || '保存失败')
  }
}

onMounted(() => {
  loadSidebarConfig()
})
</script>