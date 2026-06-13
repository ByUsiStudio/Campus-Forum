import api from './index'

// 上传相关API
export const uploadApi = {
  // 上传头像
  uploadAvatar: (formData) => api.post('/upload/avatar', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  }),

  // 上传图片
  uploadImage: (formData) => api.post('/upload/image', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  }),

  // 上传视频
  uploadVideo: (formData) => api.post('/upload/video', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  }),

  // 上传语音
  uploadVoice: (formData) => api.post('/upload/voice', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  }),
}
