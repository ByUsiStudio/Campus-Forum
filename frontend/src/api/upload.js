import api from './http'

// 涓婁紶鐩稿叧API
export const uploadApi = {
  // 涓婁紶澶村儚
  uploadAvatar: (formData) => api.post('/upload/avatar', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  }),

  // 涓婁紶鍥剧墖
  uploadImage: (formData) => api.post('/upload/image', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  }),

  // 涓婁紶瑙嗛
  uploadVideo: (formData) => api.post('/upload/video', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  }),

  // 涓婁紶璇煶
  uploadVoice: (formData) => api.post('/upload/voice', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  }),
}
