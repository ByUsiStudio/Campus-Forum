import api from './http'

// 璁よ瘉鐩稿叧API
export const authApi = {
  // 鐢ㄦ埛娉ㄥ唽
register: (data) => api.post('/auth/register', {
  username: data.username,
  qq_number: data.qq_number,
  display_name: data.display_name,
  password: data.password
}),

  // 鐢ㄦ埛鐧诲綍
  login: (data) => api.post('/auth/login', data),

  // 鍒濆鍖栫鐞嗗憳
  initAdmin: (data) => api.post('/auth/init-admin', data),

  // 妫€鏌ョ郴缁熸槸鍚﹀凡鍒濆鍖?  checkInit: () => api.get('/auth/check-init'),

  // 鍙戦€佸瘑鐮侀噸缃獙璇佺爜
  sendResetCode: (data) => api.post('/password/reset-code', data),

  // 閲嶇疆瀵嗙爜
  resetPassword: (data) => api.post('/password/reset', data),

  // 鑾峰彇褰撳墠鐢ㄦ埛淇℃伅
  getProfile: () => api.get('/profile'),

  // 鏇存柊褰撳墠鐢ㄦ埛淇℃伅
  updateProfile: (data) => api.put('/profile', data),
}
