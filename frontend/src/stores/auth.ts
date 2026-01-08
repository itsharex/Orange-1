import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import { authApi, type User, type LoginRequest } from '@/api/auth'

export const useAuthStore = defineStore('auth', () => {
  // 从 localStorage 读取登录状态
  const token = ref<string | null>(localStorage.getItem('token'))
  const user = ref<User | null>(
    JSON.parse(localStorage.getItem('user') || 'null')
  )
  const loading = ref(false)
  const error = ref<string | null>(null)

  // 计算属性
  const isLoggedIn = computed(() => !!token.value)
  const isAuthenticated = computed(() => !!token.value)

  // 登录
  async function login(credentials: LoginRequest) {
    loading.value = true
    error.value = null

    try {
      const response = await authApi.login(credentials)
      const { token: newToken, user: userData } = response.data.data

      // 保存到 state
      token.value = newToken
      user.value = userData

      // 保存到 localStorage
      localStorage.setItem('token', newToken)
      localStorage.setItem('user', JSON.stringify(userData))
      localStorage.setItem('isAuthenticated', 'true')

      return true
    } catch (err: unknown) {
      error.value = err instanceof Error ? err.message : '登录失败'
      return false
    } finally {
      loading.value = false
    }
  }

  // 注册
  async function register(data: { username: string; name: string; email?: string; phone?: string; password: string }) {
    loading.value = true
    error.value = null

    try {
      await authApi.register(data)
      return true
    } catch (err: unknown) {
      error.value = err instanceof Error ? err.message : '注册失败'
      return false
    } finally {
      loading.value = false
    }
  }

  // 退出登录
  async function logout() {
    try {
      await authApi.logout()
    } catch {
      // 忽略错误，仍然清除本地状态
    }

    // 清除状态
    token.value = null
    user.value = null

    // 清除 localStorage
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    localStorage.removeItem('isAuthenticated')
  }

  // 刷新用户信息
  async function refreshUser() {
    if (!token.value) return

    try {
      const response = await authApi.getCurrentUser()
      user.value = response.data.data
      localStorage.setItem('user', JSON.stringify(response.data.data))
    } catch {
      // Token 可能已过期
      await logout()
    }
  }

  // 更新个人信息
  async function updateProfile(data: { name?: string; phone?: string; department?: string; position?: string }) {
    loading.value = true
    error.value = null

    try {
      const response = await authApi.updateProfile(data)
      user.value = response.data.data
      localStorage.setItem('user', JSON.stringify(response.data.data))
      return true
    } catch (err: unknown) {
      error.value = err instanceof Error ? err.message : '更新失败'
      return false
    } finally {
      loading.value = false
    }
  }

  // 修改密码
  async function changePassword(oldPassword: string, newPassword: string) {
    loading.value = true
    error.value = null

    try {
      await authApi.changePassword({ old_password: oldPassword, new_password: newPassword })
      return true
    } catch (err: unknown) {
      error.value = err instanceof Error ? err.message : '修改失败'
      return false
    } finally {
      loading.value = false
    }
  }

  return {
    // State
    token,
    user,
    loading,
    error,
    // Computed
    isLoggedIn,
    isAuthenticated,
    // Actions
    login,
    register,
    logout,
    refreshUser,
    updateProfile,
    changePassword,
  }
})
