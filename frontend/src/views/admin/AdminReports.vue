<script setup>
import { ref, onMounted } from 'vue'
import { reportApi } from '../../api'

const loading = ref(false)
const reports = ref([])

const loadReports = async () => {
  loading.value = true
  try {
    const response = await reportApi.getReports()
    reports.value = response.data.reports || []
  } catch (error) {
    console.error('加载举报失败:', error)
  } finally {
    loading.value = false
  }
}

const handleReport = async (report, action) => {
  if (!confirm(`确定要${action === 'approve' ? '批准' : '拒绝'}此举报吗？`)) return
  
  loading.value = true
  try {
    await reportApi.handleReport(report.id, { action })
    loadReports()
  } catch (error) {
    console.error('处理举报失败:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadReports()
})
</script>

<template>
  <v-container fluid class="pa-0">
    <!-- 页面标题 -->
    <div class="mb-6">
      <h1 class="text-h5 font-weight-bold">举报管理</h1>
      <p class="text-body-2 text-grey">处理用户举报内容</p>
    </div>

    <!-- 刷新按钮 -->
    <v-card class="mb-4">
      <v-card-text>
        <v-btn color="primary" @click="loadReports" :loading="loading">
          <v-icon start>mdi-refresh</v-icon>
          刷新
        </v-btn>
      </v-card-text>
    </v-card>

    <!-- 举报列表 -->
    <v-card>
      <v-card-title>举报列表</v-card-title>
      <v-list>
        <v-list-item v-for="report in reports" :key="report.id">
          <v-list-item-title>{{ report.type }}</v-list-item-title>
          <v-list-item-subtitle>{{ report.reason }}</v-list-item-subtitle>
          <v-list-item-actions>
            <v-btn text color="success" @click="handleReport(report, 'approve')">批准</v-btn>
            <v-btn text color="error" @click="handleReport(report, 'reject')">拒绝</v-btn>
          </v-list-item-actions>
        </v-list-item>
      </v-list>
      <v-card-text v-if="reports.length === 0" class="text-center text-grey py-8">
        暂无举报
      </v-card-text>
    </v-card>
  </v-container>
</template>