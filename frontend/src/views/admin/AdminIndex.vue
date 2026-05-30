<template>
  <OverviewPanel 
    :statistics="statistics" 
    :loading="loading" 
    @refresh="loadStatistics" 
  />
</template>

<script setup>
import { ref, onMounted } from 'vue'
import OverviewPanel from './OverviewPanel.vue'
import api from '../../api'

const statistics = ref({
  user_count: 0,
  article_count: 0,
  comment_count: 0,
  view_count: 0
})
const loading = ref(true)

const loadStatistics = async () => {
  loading.value = true
  try {
    const response = await api.get('/admin/statistics')
    statistics.value = response.data
  } catch (error) {
    console.error('加载统计数据失败', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadStatistics()
})
</script>