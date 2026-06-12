<script setup>
import { ref, onMounted } from 'vue'
import { adminApi, siteApi } from '../../api'

const loading = ref(false)
const saving = ref(false)
const sidebarItems = ref([])
const newItem = ref({ title: '', url: '' })

const loadSidebarConfig = async () => {
  loading.value = true
  try {
    const response = await adminSidebarApi.getConfig()
    sidebarItems.value = response.data.items || []
  } catch (error) {
    console.error('加载侧边栏配置失败:', error)
  } finally {
    loading.value = false
  }
}

const addSidebarItem = () => {
  if (!newItem.value.title || !newItem.value.url) return
  sidebarItems.value.push({
    id: Date.now(),
    title: newItem.value.title,
    url: newItem.value.url
  })
  newItem.value = { title: '', url: '' }
}

const removeSidebarItem = (index) => {
  sidebarItems.value.splice(index, 1)
}

const saveSidebarConfig = async () => {
  saving.value = true
  try {
    await adminApi.updateSidebarConfig({ items: sidebarItems.value })
  } catch (error) {
    console.error('保存侧边栏配置失败:', error)
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  loadSidebarConfig()
})
</script>

<template>
  <v-container fluid class="pa-0">
    <!-- 页面标题 -->
    <div class="mb-6">
      <h1 class="text-h5 font-weight-bold">侧边栏配置</h1>
      <p class="text-body-2 text-grey">配置网站侧边栏链接</p>
    </div>

    <!-- 添加新项 -->
    <v-card class="mb-4">
      <v-card-title>添加新侧边栏项</v-card-title>
      <v-card-text>
        <v-row dense>
          <v-col cols="12" md="5">
            <v-text-field v-model="newItem.title" label="标题" hide-details />
          </v-col>
          <v-col cols="12" md="5">
            <v-text-field v-model="newItem.url" label="链接地址" hide-details />
          </v-col>
          <v-col cols="12" md="2">
            <v-btn color="primary" block @click="addSidebarItem">
              <v-icon start>mdi-plus</v-icon>
              添加
            </v-btn>
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>

    <!-- 侧边栏项列表 -->
    <v-card>
      <v-card-title>侧边栏配置</v-card-title>
      <v-list>
        <v-list-item v-for="(item, index) in sidebarItems" :key="item.id || index">
          <v-list-item-title>{{ item.title }}</v-list-item-title>
          <v-list-item-subtitle>{{ item.url }}</v-list-item-subtitle>
          <v-list-item-actions>
            <v-btn icon variant="text" size="small" color="error" @click="removeSidebarItem(index)">
              <v-icon>mdi-delete</v-icon>
            </v-btn>
          </v-list-item-actions>
        </v-list-item>
      </v-list>
      <v-card-text v-if="sidebarItems.length === 0" class="text-center text-grey py-8">
        暂无侧边栏项
      </v-card-text>
      <v-card-actions v-if="sidebarItems.length > 0">
        <v-spacer></v-spacer>
        <v-btn color="primary" @click="saveSidebarConfig" :loading="saving">
          <v-icon start>mdi-check</v-icon>
          保存配置
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-container>
</template>