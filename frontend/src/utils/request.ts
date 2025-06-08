/**
 * HTTP Request utility for 45AI
 * Handles API communication with interceptors
 */

import type { RequestOptions, ResponseData } from '@/types/api'

// Environment configuration
const isDev = process.env.NODE_ENV === 'development'
const BASE_URL = isDev ? 'http://localhost:8080/api' : 'https://api.45ai.com/api'

// Request timeout (30 seconds)
const REQUEST_TIMEOUT = 30000

// Retry configuration
const MAX_RETRIES = 3
const RETRY_DELAY = 1000

// Response codes
const SUCCESS_CODE = 0
const UNAUTHORIZED_CODE = 401
const TOKEN_EXPIRED_CODE = 4001

// Request queue for token refresh
let isRefreshing = false
let refreshSubscribers: ((token: string) => void)[] = []

// Get token from storage
const getToken = (): string => {
  return uni.getStorageSync('token') || ''
}

// Set token to storage
const setToken = (token: string): void => {
  uni.setStorageSync('token', token)
}

// Clear auth data
const clearAuth = (): void => {
  uni.removeStorageSync('token')
  uni.removeStorageSync('userInfo')
}

// Subscribe to token refresh
const subscribeTokenRefresh = (cb: (token: string) => void): void => {
  refreshSubscribers.push(cb)
}

// Notify all subscribers when token is refreshed
const onTokenRefreshed = (token: string): void => {
  refreshSubscribers.forEach(cb => cb(token))
  refreshSubscribers = []
}

// Refresh token
const refreshToken = async (): Promise<string> => {
  if (isRefreshing) {
    return new Promise(resolve => {
      subscribeTokenRefresh(resolve)
    })
  }
  
  isRefreshing = true
  
  try {
    // TODO: Implement actual refresh token API call
    const refreshToken = uni.getStorageSync('refreshToken')
    if (!refreshToken) {
      throw new Error('No refresh token')
    }
    
    // Mock refresh token response
    const newToken = 'new-jwt-token'
    setToken(newToken)
    onTokenRefreshed(newToken)
    
    return newToken
  } catch (error) {
    clearAuth()
    uni.redirectTo({
      url: '/pages/login/login'
    })
    throw error
  } finally {
    isRefreshing = false
  }
}

// Show loading
let loadingCount = 0
const showLoading = (title = 'Loading...'): void => {
  if (loadingCount === 0) {
    uni.showLoading({
      title,
      mask: true
    })
  }
  loadingCount++
}

// Hide loading
const hideLoading = (): void => {
  loadingCount--
  if (loadingCount <= 0) {
    loadingCount = 0
    uni.hideLoading()
  }
}

// Request interceptor
const requestInterceptor = (options: RequestOptions): RequestOptions => {
  const token = getToken()
  
  // Add auth header
  if (token) {
    options.header = {
      ...options.header,
      'Authorization': `Bearer ${token}`
    }
  }
  
  // Add content type if not specified
  options.header = {
    'Content-Type': 'application/json',
    ...options.header
  }
  
  // Add timestamp to prevent caching
  if (options.method === 'GET') {
    options.data = {
      ...options.data,
      _t: Date.now()
    }
  }
  
  return options
}

// Response interceptor
const responseInterceptor = async (response: any, options: RequestOptions): Promise<any> => {
  const { statusCode, data } = response
  
  // Handle HTTP errors
  if (statusCode !== 200) {
    throw new Error(`HTTP Error ${statusCode}`)
  }
  
  // Handle business logic errors
  if (data.code !== SUCCESS_CODE) {
    // Token expired, try refresh
    if (data.code === TOKEN_EXPIRED_CODE || data.code === UNAUTHORIZED_CODE) {
      try {
        const newToken = await refreshToken()
        // Retry original request with new token
        options.header!['Authorization'] = `Bearer ${newToken}`
        return await request(options)
      } catch (error) {
        // Refresh failed, redirect to login
        throw error
      }
    }
    
    // Other business errors
    throw new Error(data.message || 'Request failed')
  }
  
  return data
}

// Main request function
export const request = async <T = any>(options: RequestOptions): Promise<ResponseData<T>> => {
  const { showLoading: shouldShowLoading = true, retries = 0 } = options
  
  // Show loading if needed
  if (shouldShowLoading) {
    showLoading(options.loadingText)
  }
  
  try {
    // Apply request interceptor
    const finalOptions = requestInterceptor({
      ...options,
      url: options.url.startsWith('http') ? options.url : `${BASE_URL}${options.url}`,
      timeout: options.timeout || REQUEST_TIMEOUT
    })
    
    // Make request
    const [error, response] = await uni.request(finalOptions)
    
    if (error) {
      throw error
    }
    
    // Apply response interceptor
    const result = await responseInterceptor(response, finalOptions)
    
    return result
  } catch (error: any) {
    // Retry logic
    if (retries < MAX_RETRIES) {
      await new Promise(resolve => setTimeout(resolve, RETRY_DELAY * (retries + 1)))
      return request({ ...options, retries: retries + 1, showLoading: false })
    }
    
    // Show error message
    const message = error.message || 'Network error'
    uni.showToast({
      title: message,
      icon: 'none',
      duration: 2000
    })
    
    throw error
  } finally {
    if (shouldShowLoading) {
      hideLoading()
    }
  }
}

// Convenience methods
export const get = <T = any>(url: string, data?: any, options?: Partial<RequestOptions>) => {
  return request<T>({
    method: 'GET',
    url,
    data,
    ...options
  })
}

export const post = <T = any>(url: string, data?: any, options?: Partial<RequestOptions>) => {
  return request<T>({
    method: 'POST',
    url,
    data,
    ...options
  })
}

export const put = <T = any>(url: string, data?: any, options?: Partial<RequestOptions>) => {
  return request<T>({
    method: 'PUT',
    url,
    data,
    ...options
  })
}

export const del = <T = any>(url: string, data?: any, options?: Partial<RequestOptions>) => {
  return request<T>({
    method: 'DELETE',
    url,
    data,
    ...options
  })
}

// Upload file
export const upload = async (options: {
  url: string
  filePath: string
  name?: string
  formData?: any
  header?: any
}): Promise<any> => {
  const token = getToken()
  
  showLoading('Uploading...')
  
  try {
    const [error, response] = await uni.uploadFile({
      url: `${BASE_URL}${options.url}`,
      filePath: options.filePath,
      name: options.name || 'file',
      formData: options.formData,
      header: {
        'Authorization': token ? `Bearer ${token}` : '',
        ...options.header
      }
    })
    
    if (error) {
      throw error
    }
    
    const data = JSON.parse(response.data)
    
    if (data.code !== SUCCESS_CODE) {
      throw new Error(data.message || 'Upload failed')
    }
    
    return data
  } catch (error: any) {
    uni.showToast({
      title: error.message || 'Upload failed',
      icon: 'none'
    })
    throw error
  } finally {
    hideLoading()
  }
}

// Export types
export * from '@/types/api' 