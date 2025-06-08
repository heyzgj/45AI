/**
 * Authentication API
 */

import { post } from '@/utils/request'
import type { LoginRequest, LoginResponse, User } from '@/types/api'

// WeChat login
export const login = (data: LoginRequest) => {
  return post<LoginResponse>('/auth/wechat/login', data)
}

// Get current user info
export const getUserInfo = () => {
  return post<User>('/auth/me')
}

// Refresh token
export const refreshToken = (refreshToken: string) => {
  return post<LoginResponse>('/auth/refresh', { refresh_token: refreshToken })
}

// Logout
export const logout = () => {
  return post('/auth/logout')
} 