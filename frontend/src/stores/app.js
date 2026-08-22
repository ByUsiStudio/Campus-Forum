import { defineStore } from 'pinia'
import { ref } from 'vue'
import http from '@/api/http'

/**
 * 站点级配置 / 元信息状态
 */
export const useAppStore = defineStore('app', () => {
  const siteTitle = ref('校园论坛')
  const icpNumber = ref(null)
  const publicSecurityNumber = ref(null)
  const backendVersion = ref('')
  const frontendVersion = ref('')
  const announcement = ref({})
  const initialized = ref(false)

  async function initApp() {
    if (initialized.value) return
    await Promise.allSettled([loadSiteConfig(), loadVersion()])
    initialized.value = true
  }

  async function loadSiteConfig() {
    try {
      const { data } = await http.get('/site-config')
      siteTitle.value = data.site_title || '校园论坛'
      icpNumber.value = data.icp_number || null
      publicSecurityNumber.value = data.public_security_number || null
      if (typeof document !== 'undefined') document.title = siteTitle.value
    } catch (e) { /* http 层已提示 */ }
  }

  async function loadVersion() {
    try {
      const { data } = await http.get('/version')
      backendVersion.value = data.backend?.version || data.backend_version || data.version || ''
      frontendVersion.value = data.frontend?.version || ''
    } catch (e) { /* ignore */ }
  }

  async function loadAnnouncement() {
    try {
      const { data } = await http.get('/announcement')
      announcement.value = data || {}
    } catch (e) { /* ignore */ }
  }

  return {
    siteTitle, icpNumber, publicSecurityNumber,
    backendVersion, frontendVersion, announcement, initialized,
    initApp, loadSiteConfig, loadVersion, loadAnnouncement
  }
})
