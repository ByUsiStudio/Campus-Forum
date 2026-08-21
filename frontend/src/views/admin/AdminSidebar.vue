<template>
  <v-container fluid class="pa-4 pa-md-6">
    <!-- 添加新侧边栏项表单 -->
    <v-card class="mb-4" variant="flat" rounded="lg">
      <v-card-title class="pb-2">
        <v-icon start size="20">mdi-plus-circle</v-icon>
        添加新侧边栏项
      </v-card-title>
      <v-card-text>
        <v-row dense>
          <v-col cols="12" sm="5">
            <v-text-field
              v-model="newItem.title"
              label="标题"
              placeholder="例如：联系我们"
              variant="outlined"
              density="compact"
              hide-details
            />
          </v-col>
          <v-col cols="12" sm="5">
            <v-text-field
              v-model="newItem.url"
              label="链接地址"
              placeholder="例如：https://example.com"
              variant="outlined"
              density="compact"
              hide-details
            />
          </v-col>
          <v-col cols="12" sm="2">
            <v-btn color="primary" block height="40" @click="addSidebarItem">
              <v-icon start>mdi-plus</v-icon>
              添加
            </v-btn>
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>

    <!-- 侧边栏项列表 -->
    <v-card variant="flat" rounded="lg">
      <v-card-title class="pb-2">
        <v-icon start size="20">mdi-web</v-icon>
        侧边栏配置
      </v-card-title>
      <v-list lines="two" v-if="sidebarItems.length > 0">
        <v-list-item v-for="(item, index) in sidebarItems" :key="item.id || index" class="py-3">
          <template v-slot:prepend>
            <v-avatar size="48" color="primary" variant="tonal">
              <v-icon>mdi-link</v-icon>
            </v-avatar>
          </template>

          <v-list-item-title class="font-weight-medium mb-1">
            {{ item.title }}
          </v-list-item-title>

          <v-list-item-subtitle>
            {{ item.url }}
          </v-list-item-subtitle>

          <template v-slot:append>
            <v-btn size="small" color="error" variant="text" @click="removeSidebarItem(index)">
              <v-icon>mdi-delete</v-icon>
              <v-tooltip activator="parent">删除</v-tooltip>
            </v-btn>
          </template>
        </v-list-item>
      </v-list>

      <v-card-text v-else class="text-center py-8">
        <v-icon size="48" color="grey-lighten-1">mdi-inbox</v-icon>
        <div class="text-body-1 text-medium-emphasis mt-2">
          暂无侧边栏项
        </div>
      </v-card-text>

      <v-card-actions class="pa-4" v-if="sidebarItems.length > 0">
        <v-btn color="primary" variant="flat" @click="saveSidebarConfig" :loading="saving">
          <v-icon start>mdi-check</v-icon>
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