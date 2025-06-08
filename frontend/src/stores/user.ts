/**
 * User Store - Authentication and Profile Management
 */

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User, LoginResponse } from '@/types/api'
import { auth } from '@/api'
import { onAuthStateChange } from '@/utils/navigation'

export const useUserStore = defineStore('user', () => {
  // State
  const token = ref<string>('')
  const refreshToken = ref<string>('')
  const userInfo = ref<User | null>(null)
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  // Getters
  const isAuthenticated = computed(() => !!token.value && !!userInfo.value)
  const credits = computed(() => userInfo.value?.credits || 0)
  const hasCredits = computed(() => credits.value > 0)
  const userId = computed(() => userInfo.value?.id || '')

  // Actions
  const setAuth = (data: LoginResponse) => {
    token.value = data.token
    refreshToken.value = data.refreshToken
    userInfo.value = data.user
    
    // Persist to storage
    uni.setStorageSync('token', data.token)
    uni.setStorageSync('refreshToken', data.refreshToken)
    uni.setStorageSync('userInfo', JSON.stringify(data.user))
    
    // Notify navigation system
    onAuthStateChange(true)
  }

  const clearAuth = () => {
    token.value = ''
    refreshToken.value = ''
    userInfo.value = null
    
    // Clear storage
    uni.removeStorageSync('token')
    uni.removeStorageSync('refreshToken')
    uni.removeStorageSync('userInfo')
    
    // Notify navigation system
    onAuthStateChange(false)
  }

  const initFromStorage = () => {
    try {
      const storedToken = uni.getStorageSync('token')
      const storedRefreshToken = uni.getStorageSync('refreshToken')
      const storedUserInfo = uni.getStorageSync('userInfo')
      
      if (storedToken && storedRefreshToken && storedUserInfo) {
        token.value = storedToken
        refreshToken.value = storedRefreshToken
        userInfo.value = JSON.parse(storedUserInfo)
      }
    } catch (e) {
      console.error('Failed to init user from storage:', e)
      clearAuth()
    }
  }

  const login = async (code: string) => {
    isLoading.value = true
    error.value = null
    
    try {
      const response = await auth.login({ code })
      setAuth(response.data)
      return response.data
    } catch (e: any) {
      error.value = e.message || 'Login failed'
      throw e
    } finally {
      isLoading.value = false
    }
  }

  const logout = async () => {
    try {
      await auth.logout()
    } catch (e) {
      // Ignore logout errors
    } finally {
      clearAuth()
    }
  }

  const fetchUserInfo = async () => {
    if (!token.value) return
    
    isLoading.value = true
    error.value = null
    
    try {
      const response = await auth.getUserInfo()
      userInfo.value = response.data
      uni.setStorageSync('userInfo', JSON.stringify(response.data))
    } catch (e: any) {
      error.value = e.message || 'Failed to fetch user info'
      throw e
    } finally {
      isLoading.value = false
    }
  }

  const updateCredits = (amount: number) => {
    if (userInfo.value) {
      userInfo.value.credits += amount
      uni.setStorageSync('userInfo', JSON.stringify(userInfo.value))
    }
  }

  const deductCredits = (amount: number) => {
    if (userInfo.value && userInfo.value.credits >= amount) {
      userInfo.value.credits -= amount
      uni.setStorageSync('userInfo', JSON.stringify(userInfo.value))
      return true
    }
    return false
  }

  const refreshUserToken = async () => {
    if (!refreshToken.value) {
      throw new Error('No refresh token available')
    }
    
    try {
      const response = await auth.refreshToken(refreshToken.value)
      setAuth(response.data)
      return response.data.token
    } catch (e) {
      clearAuth()
      throw e
    }
  }

  // WeChat specific methods
  const loginWithWeChat = async () => {
    return new Promise<LoginResponse>((resolve, reject) => {
      uni.login({
        provider: 'weixin',
        success: async (loginRes) => {
          if (loginRes.code) {
            try {
              const response = await login(loginRes.code)
              resolve(response)
            } catch (e) {
              reject(e)
            }
          } else {
            reject(new Error('WeChat login failed: no code'))
          }
        },
        fail: (err) => {
          reject(new Error(err.errMsg || 'WeChat login failed'))
        }
      })
    })
  }

  return {
    // State
    token,
    refreshToken,
    userInfo,
    isLoading,
    error,
    
    // Getters
    isAuthenticated,
    credits,
    hasCredits,
    userId,
    
    // Actions
    setAuth,
    clearAuth,
    initFromStorage,
    login,
    logout,
    fetchUserInfo,
    updateCredits,
    deductCredits,
    refreshUserToken,
    loginWithWeChat
  }
})

// Export type for easy access
export type UserStore = ReturnType<typeof useUserStore> 