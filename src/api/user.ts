import { get, put } from '@/utils/request'
import type { ResponseData } from '@/types/api'

export interface UserInfo {
  id: number
  openid?: string
  avatar_url?: string
  name?: string
  credits: number
  created_at: string
  updated_at?: string
}

export interface Transaction {
  id: number
  user_id: number
  type: 'purchase' | 'consume'
  credits: number
  description: string
  created_at: string
}

// Get user information
export const getUserInfo = (userId?: string): Promise<ResponseData<UserInfo>> => {
  const url = userId ? `/users/${userId}` : '/me'
  return get(url)
}

// Update user information
export const updateUserInfo = (data: Partial<UserInfo>): Promise<ResponseData<UserInfo>> => {
  return put('/me', data)
}

// Get user transaction history
export const getUserTransactions = (page = 1, limit = 20): Promise<ResponseData<Transaction[]>> => {
  return get('/me/transactions', { page, limit })
}

// Get user credit balance
export const getUserCredits = (): Promise<ResponseData<{ credits: number }>> => {
  return get('/me/credits')
}
