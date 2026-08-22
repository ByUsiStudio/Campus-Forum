<template>
  <div class="page-container">
    <el-card shadow="never" class="card-surface">
      <template #header>
        <div class="d-flex align-center justify-space-between flex-wrap gap-2">
          <div class="d-flex align-center">
            <el-button text circle @click="router.back()" aria-label="返回">
              <el-icon><ArrowLeft /></el-icon>
            </el-button>
            <span class="title">我的通知</span>
          </div>
          <el-button
            v-if="unreadCount > 0"
            type="primary"
            plain
            size="small"
            @click="markAllRead"
          >
            <el-icon class="mr-1"><Select /></el-icon>
            全部标记已读
          </el-button>
        </div>
      </template>

      <template v-if="notifications.length > 0">
        <el-list class="notification-list">
          <el-list-item
            v-for="notification in notifications"
            :key="notification.id"
            class="notification-item"
            :class="{ 'notification-item-unread': !notification.is_read }"
          >
            <template #default>
              <div class="notification-row">
                <el-avatar
                  :size="48"
                  :style="{ backgroundColor: getTypeColor(notification.type) }"
                  class="notification-avatar"
                >
                  <el-icon :size="22" color="white">
                    <component :is="getTypeIcon(notification.type)" />
                  </el-icon>
                </el-avatar>

                <div class="notification-body">
                  <div class="notification-title-row">
                    <el-tag
                      size="small"
                      :style="tagStyle(getTypeColor(notification.type))"
                      class="notification-tag"
                    >
                      {{ getTypeText(notification.type) }}
                    </el-tag>
                    <span class="notification-title">{{ notification.title }}</span>
                  </div>
                  <div class="notification-content">{{ notification.content }}</div>
                  <div class="notification-time">
                    <el-icon class="mr-1"><Clock /></el-icon>
                    {{ formatDate(notification.created_at) }}
                  </div>
                </div>

                <el-button
                  v-if="!notification.is_read"
                  type="primary"
                  plain
                  size="small"
                  @click="markRead(notification)"
                >
                  标记已读
                </el-button>
              </div>
            </template>
          </el-list-item>
        </el-list>
      </template>

      <div v-else class="notification-empty text-center">
        <el-icon :size="80" color="#e4e7ed"><Bell /></el-icon>
        <div class="empty-title">暂无通知</div>
        <div class="text-secondary empty-desc">暂无新的通知消息</div>
      </div>

      <div v-if="totalPages > 1" class="pagination-wrap">
        <el-pagination
          v-model:current-page="page"
          :page-count="totalPages"
          :page-size="pageSize"
          layout="prev, pager, next"
          background
        />
      </div>
    </el-card>
  </div>
</template>

<script>
import { ref, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { ArrowLeft, Bell, Clock, InfoFilled, Calendar, Refresh, Warning, Select } from '@element-plus/icons-vue'
import { notificationApi } from '../api'

export default {
  name: 'Notifications',
  setup() {
    const router = useRouter()
    const notifications = ref([])
    const unreadCount = ref(0)
    const page = ref(1)
    const totalPages = ref(1)
    const pageSize = 20

    const loadNotifications = async () => {
      try {
        const response = await notificationApi.getNotifications()
        const allNotifications = response.data.notifications || []

        const total = allNotifications.length
        totalPages.value = Math.ceil(total / pageSize)

        const start = (page.value - 1) * pageSize
        notifications.value = allNotifications.slice(start, start + pageSize)
      } catch (error) {
        console.error('加载通知失败', error)
      }
    }

    const loadUnreadCount = async () => {
      try {
        const response = await notificationApi.getUnreadCount()
        unreadCount.value = response.data.unread_count || 0
      } catch (error) {
        console.error('加载未读数量失败', error)
      }
    }

    const markRead = async (notification) => {
      try {
        await notificationApi.markNotificationRead(notification.id)
        notification.is_read = true
        unreadCount.value = Math.max(0, unreadCount.value - 1)
      } catch (error) {
        console.error('标记已读失败', error)
      }
    }

    const markAllRead = async () => {
      try {
        await notificationApi.markAllNotificationsRead()
        notifications.value.forEach((n) => { n.is_read = true })
        unreadCount.value = 0
      } catch (error) {
        console.error('标记全部已读失败', error)
      }
    }

    const getTypeColor = (type) => {
      const colors = {
        system: '#6750a4',
        activity: '#67c23a',
        update: '#409eff',
        warning: '#e6a23c'
      }
      return colors[type] || '#6b7280'
    }

    const tagStyle = (color) => {
      return {
        color,
        backgroundColor: `${color}1a`,
        borderColor: `${color}33`
      }
    }

    const getTypeIcon = (type) => {
      const icons = {
        system: InfoFilled,
        activity: Calendar,
        update: Refresh,
        warning: Warning
      }
      return icons[type] || Bell
    }

    const getTypeText = (type) => {
      const texts = {
        system: '系统通知',
        activity: '活动公告',
        update: '更新通知',
        warning: '警告通知'
      }
      return texts[type] || type
    }

    const formatDate = (date) => {
      return new Date(date).toLocaleString('zh-CN')
    }

    watch(page, () => {
      loadNotifications()
    })

    onMounted(() => {
      loadNotifications()
      loadUnreadCount()
    })

    return {
      router,
      notifications,
      unreadCount,
      page,
      totalPages,
      pageSize,
      markRead,
      markAllRead,
      getTypeColor,
      getTypeIcon,
      getTypeText,
      tagStyle,
      formatDate
    }
  }
}
</script>

<style scoped>
.d-flex {
  display: flex;
}
.align-center {
  align-items: center;
}
.justify-space-between {
  justify-content: space-between;
}
.flex-wrap {
  flex-wrap: wrap;
}
.gap-2 {
  gap: 8px;
}
.mr-1 {
  margin-right: 4px;
}
.text-center {
  text-align: center;
}
.title {
  font-size: 18px;
  font-weight: 600;
  color: var(--campus-text);
}
.notification-list {
  padding: 0;
}
.notification-item {
  border-bottom: 1px solid var(--campus-border);
  padding: 12px 8px;
}
.notification-item:last-child {
  border-bottom: none;
}
.notification-item-unread {
  background: #f5f5fa;
}
.notification-row {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  width: 100%;
}
.notification-avatar {
  flex-shrink: 0;
}
.notification-body {
  flex: 1;
  min-width: 0;
}
.notification-title-row {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 6px;
  margin-bottom: 4px;
}
.notification-tag {
  flex-shrink: 0;
  border: 1px solid transparent;
}
.notification-title {
  font-weight: 600;
  font-size: 14px;
  color: var(--campus-text);
}
.notification-content {
  font-size: 13px;
  color: var(--campus-text);
  margin-top: 4px;
  word-break: break-word;
}
.notification-time {
  display: flex;
  align-items: center;
  font-size: 12px;
  color: var(--campus-text-secondary);
  margin-top: 6px;
}
.notification-empty {
  padding: 48px 0;
}
.empty-title {
  font-size: 16px;
  font-weight: 500;
  color: var(--campus-text-secondary);
  margin-top: 16px;
}
.empty-desc {
  margin-top: 8px;
  font-size: 13px;
}
.pagination-wrap {
  display: flex;
  justify-content: center;
  padding: 16px 0 8px;
  border-top: 1px solid var(--campus-border);
  margin-top: 8px;
}
</style>
