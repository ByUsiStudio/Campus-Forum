<template>
  <div class="content-panel">
    <div class="panel-header">
      <div class="header-left">
        <h2 class="panel-title">内容管理</h2>
        <p class="panel-subtitle">管理侧边栏、删除审核与公告</p>
      </div>
    </div>

    <v-row>
      <v-col cols="12" lg="4">
        <v-card class="section-card">
          <v-card-title class="section-title">
            <v-icon class="title-icon">mdi-sidebar</v-icon>
            侧边栏配置
          </v-card-title>
          <v-card-text>
            <div class="sidebar-items">
              <div v-for="(item, index) in localSidebarItems" :key="index" class="sidebar-item">
                <v-text-field
                  v-model="item.title"
                  label="标题"
                  variant="outlined"
                  density="compact"
                  hide-details
                  class="mb-2"
                ></v-text-field>
                <v-text-field
                  v-model="item.link"
                  label="链接"
                  variant="outlined"
                  density="compact"
                  hide-details
                  class="mb-2"
                ></v-text-field>
                <v-text-field
                  v-model="item.icon"
                  label="图标"
                  variant="outlined"
                  density="compact"
                  hide-details
                ></v-text-field>
                <v-btn
                  variant="text"
                  color="error"
                  size="small"
                  @click="$emit('remove-sidebar-item', index)"
                  class="mt-1"
                >
                  <v-icon start>mdi-close</v-icon>
                  移除
                </v-btn>
              </div>
            </div>

            <v-btn
              variant="outlined"
              color="primary"
              block
              class="mt-4"
              @click="$emit('add-sidebar-item')"
            >
              <v-icon start>mdi-plus</v-icon>
              添加链接
            </v-btn>

            <v-btn
              color="primary"
              block
              class="mt-3"
              @click="$emit('save-sidebar')"
            >
              <v-icon start>mdi-content-save</v-icon>
              保存配置
            </v-btn>
          </v-card-text>
        </v-card>
      </v-col>

      <v-col cols="12" lg="4">
        <v-card class="section-card">
          <v-card-title class="section-title">
            <v-icon class="title-icon text-error">mdi-file-question</v-icon>
            删除审核
            <v-chip v-if="deletionRequests.length" size="small" color="error" class="ml-2">
              {{ deletionRequests.length }} 待处理
            </v-chip>
          </v-card-title>
          <v-card-text>
            <div v-if="deletionRequests.length === 0" class="empty-state">
              <v-icon size="40" color="success">mdi-check-circle</v-icon>
              <div class="empty-text">暂无待审核申请</div>
            </div>

            <div v-else class="requests-list">
              <div v-for="req in deletionRequests" :key="req.id" class="request-card">
                <div class="request-header">
                  <span class="request-title">{{ req.article?.title || '文章已删除' }}</span>
                  <v-chip size="x-small" color="warning">待审核</v-chip>
                </div>
                <div class="request-meta">
                  <span>申请人：{{ req.user?.display_name }}</span>
                </div>
                <div class="request-meta">
                  <span>原因：{{ req.reason }}</span>
                </div>
                <div class="request-actions">
                  <v-btn
                    color="primary"
                    size="small"
                    variant="tonal"
                    @click="$emit('approve-deletion', req.id)"
                  >
                    批准
                  </v-btn>
                  <v-btn
                    color="error"
                    size="small"
                    variant="outlined"
                    @click="$emit('reject-deletion', req.id)"
                  >
                    拒绝
                  </v-btn>
                </div>
              </div>
            </div>
          </v-card-text>
        </v-card>
      </v-col>

      <v-col cols="12" lg="4">
        <v-card class="section-card">
          <v-card-title class="section-title">
            <v-icon class="title-icon text-warning">mdi-bullhorn</v-icon>
            公告管理
          </v-card-title>
          <v-card-text>
            <v-textarea
              v-model="localAnnouncement"
              label="公告内容（支持Markdown）"
              variant="outlined"
              rows="8"
              placeholder="输入公告内容..."
            ></v-textarea>
            <v-btn
              color="primary"
              block
              class="mt-4"
              @click="$emit('save-announcement')"
            >
              <v-icon start>mdi-content-save</v-icon>
              保存公告
            </v-btn>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </div>
</template>

<script>
import { ref, watch } from 'vue'

export default {
  name: 'ContentPanel',
  props: {
    sidebarItems: {
      type: Array,
      default: () => []
    },
    deletionRequests: {
      type: Array,
      default: () => []
    },
    announcement: {
      type: String,
      default: ''
    }
  },
  emits: [
    'add-sidebar-item',
    'remove-sidebar-item',
    'save-sidebar',
    'approve-deletion',
    'reject-deletion',
    'save-announcement',
    'refresh-sidebar',
    'refresh-deletions',
    'refresh-announcement'
  ],
  setup(props) {
    const localSidebarItems = ref([...props.sidebarItems])
    const localAnnouncement = ref(props.announcement)

    watch(() => props.sidebarItems, (val) => {
      localSidebarItems.value = [...val]
    }, { deep: true })

    watch(() => props.announcement, (val) => {
      localAnnouncement.value = val
    })

    return {
      localSidebarItems,
      localAnnouncement
    }
  }
}
</script>

<style scoped>
.content-panel {
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 24px;
  flex-wrap: wrap;
  gap: 16px;
}

.panel-title {
  font-size: 1.5rem;
  font-weight: 700;
  color: #1a1a2e;
  margin: 0 0 4px 0;
}

.panel-subtitle {
  font-size: 0.9rem;
  color: #6b7280;
  margin: 0;
}

.section-card {
  border-radius: 16px;
  height: 100%;
}

.section-title {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 1rem;
  font-weight: 600;
  padding: 16px 20px !important;
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
}

.title-icon {
  width: 32px;
  height: 32px;
  padding: 4px;
  border-radius: 8px;
  background: rgba(103, 80, 164, 0.1);
  color: #6750A4 !important;
}

.sidebar-items {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.sidebar-item {
  padding: 16px;
  background: #f8f9ff;
  border-radius: 10px;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 32px 16px;
  text-align: center;
}

.empty-text {
  margin-top: 12px;
  font-size: 0.9rem;
  color: #9ca3af;
}

.requests-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.request-card {
  padding: 14px;
  background: #fff8f8;
  border: 1px solid rgba(244, 67, 54, 0.1);
  border-radius: 10px;
}

.request-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 8px;
  margin-bottom: 8px;
}

.request-title {
  font-weight: 600;
  color: #1a1a2e;
  font-size: 0.95rem;
}

.request-meta {
  font-size: 0.85rem;
  color: #6b7280;
  margin-top: 4px;
}

.request-actions {
  display: flex;
  gap: 8px;
  margin-top: 12px;
}

@media (max-width: 960px) {
  .section-card {
    height: auto;
  }
}
</style>
