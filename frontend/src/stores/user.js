import { ref, computed } from 'vue'

const user = ref(null)

export function useUserStore() {
  const isLoggedIn = computed(() => !!user.value)

  const setUser = (userData) => {
    user.value = userData
  }

  const clearUser = () => {
    user.value = null
  }

  const initUser = () => {
    const userStr = localStorage.getItem('user')
    if (userStr) {
      try {
        user.value = JSON.parse(userStr)
      } catch (e) {
        console.error('解析用户信息失败', e)
      }
    }
  }

  return {
    user,
    isLoggedIn,
    setUser,
    clearUser,
    initUser
  }
}
