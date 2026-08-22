import { defineStore } from 'pinia'
import { ref } from 'vue'
import http from '@/api/http'

/**
 * 用户通知状态（顶部铃铛 + 通知列表）
 */
export const useNotificationStore = defineStore('notification', () => {
  const unreadCount = ref(0)
  const notifications = ref([])
  let pollTimer = null

  async function loadUnreadCount() {
    try {
      const { data } = await http.get('/notifications/unread-count')
      unreadCount.value = data.unread_count || 0
    } catch (e) { /* ignore */ }
  }

  async function loadNotifications(limit = 5) {
    try {
      const { data } = await http.get('/notifications')
      notifications.value = (data.notifications || []).slice(0, limit)
    } catch (e) { /* ignore */ }
  }

  async function markRead(id) {
    try {
      await http.post(`/notifications/${id}/read`)
      const n = notifications.value.find(x => x.id === id)
      if (n) n.is_read = true
      unreadCount.value = Math.max(0, unreadCount.value - 1)
    } catch (e) { /* ignore */ }
  }

  async function markAllRead() {
    try {
      await http.post('/notifications/read-all')
      notifications.value.forEach(n => (n.is_read = true))
      unreadCount.value = 0
    } catch (e) { /* ignore */ }
  }

  function startPolling() {
    stopPolling()
    loadUnreadCount()
    loadNotifications()
    pollTimer = setInterval(() => loadUnreadCount(), 60000)
  }

  function stopPolling() {
    if (pollTimer) {
      clearInterval(pollTimer)
      pollTimer = null
    }
  }

  return {
    unreadCount, notifications,
    loadUnreadCount, loadNotifications, markRead, markAllRead, startPolling, stopPolling
  }
})
