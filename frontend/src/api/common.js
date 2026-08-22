import api from './http'

// 閫氱敤API
export const commonApi = {
  // 鑾峰彇鐗堟湰淇℃伅
  getVersion: () => api.get('/version'),

  // 鑾峰彇鍒嗗尯鍒楄〃
  getCategories: () => api.get('/categories'),

  // 鑾峰彇鍏憡
  getAnnouncement: () => api.get('/announcement'),

  // 鑾峰彇渚ц竟鏍忛厤缃?  getSidebarConfig: () => api.get('/sidebar-config'),

  // 鑾峰彇缃戠珯閰嶇疆
  getSiteConfig: () => api.get('/site-config'),

  // 鑾峰彇鎵€鏈夊ご琛?  getTitles: () => api.get('/titles'),

  // 妫€鏌ョ鐞嗗憳鏉冮檺
  checkAdmin: () => api.get('/admin/check'),
}
