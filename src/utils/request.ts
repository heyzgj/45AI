/**
 * HTTP Request utility for 45AI
 * Handles API communication with interceptors
 */

import type { RequestOptions, ResponseData } from '@/types/api'
import { getEnvBaseUrl } from './index'

// Environment configuration
const BASE_URL = `${getEnvBaseUrl()}/api/v1`

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
  try {
    // Try to get from Pinia store first
    const token = uni.getStorageSync('token') || ''
    return token
  } catch (error) {
    console.error('Failed to get token:', error)
    return ''
  }
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
  refreshSubscribers.forEach((cb) => cb(token))
  refreshSubscribers = []
}

// Refresh token
const refreshToken = async (): Promise<string> => {
  if (isRefreshing) {
    return new Promise((resolve) => {
      subscribeTokenRefresh(resolve)
    })
  }

  isRefreshing = true

  try {
    // Import user store dynamically to avoid circular dependency
    const { useUserStore } = await import('@/store/user')
    const userStore = useUserStore()

    if (!userStore.refreshToken) {
      throw new Error('No refresh token available')
    }

    // Use real refresh token API
    const newToken = await userStore.refreshUserToken()
    onTokenRefreshed(newToken)

    return newToken
  } catch (error) {
    console.error('Token refresh failed:', error)
    clearAuth()

    // Show error message and redirect to login
    uni.showToast({
      title: '登录已过期，请重新登录',
      icon: 'none',
      duration: 2000,
    })

    setTimeout(() => {
      uni.redirectTo({
        url: '/pages/login/index',
      })
    }, 2000)

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
      mask: true,
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
      Authorization: `Bearer ${token}`,
    }
  }

  // Add content type if not specified
  options.header = {
    'Content-Type': 'application/json',
    ...options.header,
  }

  // Add timestamp to prevent caching
  if (options.method === 'GET') {
    options.data = {
      ...options.data,
      _t: Date.now(),
    }
  }

  return options
}

// Response interceptor
const responseInterceptor = async (response: any, options: RequestOptions): Promise<any> => {
  const { statusCode, data } = response

  console.log('responseInterceptor: Received response:', {
    statusCode,
    data,
    dataType: typeof data,
  })

  // Handle HTTP 401 errors specifically for token refresh
  if (statusCode === 401) {
    console.log('responseInterceptor: Received 401, attempting token refresh')
    try {
      const newToken = await refreshToken()
      // Retry original request with new token
      options.header!['Authorization'] = `Bearer ${newToken}`
      return await request(options)
    } catch (error) {
      // Refresh failed, redirect to login
      console.error('Token refresh failed:', error)
      throw new Error('Authentication failed - please login again')
    }
  }

  // Handle other HTTP errors
  if (statusCode < 200 || statusCode >= 300) {
    console.error(`HTTP Error ${statusCode}:`, data)

    // Extract error message from response
    let errorMessage = `HTTP Error ${statusCode}`
    if (data && typeof data === 'object') {
      errorMessage = data.message || data.error || errorMessage
    } else if (typeof data === 'string') {
      errorMessage = data
    }

    throw new Error(errorMessage)
  }

  // Check if data has the expected wrapper format with 'code' field
  if (data && typeof data === 'object' && 'code' in data) {
    console.log('responseInterceptor: Data has code field:', data.code)

    // Handle business logic errors
    if (data.code !== SUCCESS_CODE) {
      // Token expired, try refresh
      if (data.code === TOKEN_EXPIRED_CODE || data.code === UNAUTHORIZED_CODE) {
        console.log('responseInterceptor: Business logic 401, attempting token refresh')
        try {
          const newToken = await refreshToken()
          // Retry original request with new token
          options.header!['Authorization'] = `Bearer ${newToken}`
          return await request(options)
        } catch (error) {
          // Refresh failed, redirect to login
          console.error('Token refresh failed:', error)
          throw error
        }
      }

      // Other business errors
      throw new Error(data.message || 'Request failed')
    }

    return data
  } else {
    // Backend returns raw data without wrapper - this is fine for our current backend
    console.log('responseInterceptor: Backend returned raw data without wrapper, returning as-is')
    return { data: data, code: 0, message: 'success' }
  }
}

// Main request function
export const request = async <T = any>(options: RequestOptions): Promise<ResponseData<T>> => {
  const { showLoading: shouldShowLoading = true, retries = 0 } = options

  // Debug logging
  console.log('request: Initial options:', {
    url: options.url,
    method: options.method,
    data: options.data,
  })

  // Validate required arguments
  if (!options.url) {
    console.error('request: Missing required args: "url"', options)
    throw new Error('Missing required args: "url"')
  }

  // Show loading if needed
  if (shouldShowLoading) {
    showLoading(options.loadingText)
  }

  try {
    // Construct final URL with validation
    let finalUrl = options.url
    if (!finalUrl.startsWith('http')) {
      finalUrl = `${BASE_URL}${finalUrl}`
    }

    console.log('request: Final URL constructed:', finalUrl)

    // Additional validation to catch malformed URLs
    if (
      finalUrl.includes('/NaN') ||
      finalUrl.includes('/undefined') ||
      finalUrl.includes('/null')
    ) {
      console.error('Malformed URL detected:', finalUrl)
      throw new Error('Invalid API endpoint - malformed ID in URL')
    }

    // Apply request interceptor
    const finalOptions = requestInterceptor({
      ...options,
      url: finalUrl,
      timeout: options.timeout || REQUEST_TIMEOUT,
    })

    console.log('request: Final options for uni.request:', finalOptions)

    // Convert to UniApp request format
    const uniRequestOptions: UniApp.RequestOptions = {
      url: finalOptions.url!,
      method: (finalOptions.method as any) || 'GET',
      data: finalOptions.data,
      header: finalOptions.header,
      timeout: finalOptions.timeout,
      dataType: 'json',
      responseType: 'text',
      success: () => {}, // Will be overridden by Promise wrapper
      fail: () => {}, // Will be overridden by Promise wrapper
    }

    console.log('request: UniApp request options:', uniRequestOptions)

    // Final validation before making request
    if (!uniRequestOptions.url || uniRequestOptions.url.trim() === '') {
      console.error('request: URL is empty or invalid before uni.request call:', uniRequestOptions)
      throw new Error('Invalid URL: URL is empty')
    }

    console.log('request: About to call uni.request with URL:', uniRequestOptions.url)

    // Use native fetch() instead of broken uni.request in H5 mode
    const response = await new Promise<UniApp.RequestSuccessCallbackResult>((resolve, reject) => {
      const url = uniRequestOptions.url!
      const method = uniRequestOptions.method || 'GET'
      const headers = uniRequestOptions.header || {}

      // Prepare fetch options
      const fetchOptions: RequestInit = {
        method: method,
        headers: headers,
      }

      // Add body for non-GET requests
      if (method !== 'GET' && uniRequestOptions.data) {
        if (headers['Content-Type']?.includes('application/json')) {
          fetchOptions.body = JSON.stringify(uniRequestOptions.data)
        } else {
          fetchOptions.body = String(uniRequestOptions.data)
        }
      }

      console.log('request: Using native fetch with:', { url, fetchOptions })

      // Make the request with native fetch
      fetch(url, fetchOptions)
        .then(async (fetchResponse) => {
          console.log('request: Fetch response:', {
            status: fetchResponse.status,
            ok: fetchResponse.ok,
          })

          const responseText = await fetchResponse.text()
          let responseData

          try {
            responseData = JSON.parse(responseText)
          } catch (e) {
            responseData = responseText
          }

          const uniResponse: UniApp.RequestSuccessCallbackResult = {
            statusCode: fetchResponse.status,
            data: responseData,
            header: {} as any,
            cookies: [],
          }

          console.log('request: Success response converted:', uniResponse)
          resolve(uniResponse)
        })
        .catch((error) => {
          console.error('request: Fetch error:', error)
          reject(error)
        })
    })

    // Apply response interceptor
    const result = await responseInterceptor(response, finalOptions)

    return result
  } catch (error: any) {
    console.error('request: Error occurred:', error)
    // Retry logic
    if (retries < MAX_RETRIES) {
      await new Promise((resolve) => setTimeout(resolve, RETRY_DELAY * (retries + 1)))
      return request({ ...options, retries: retries + 1, showLoading: false })
    }

    // Show error message
    const message = error.message || 'Network error'
    uni.showToast({
      title: message,
      icon: 'none',
      duration: 2000,
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
    ...options,
  })
}

export const post = <T = any>(url: string, data?: any, options?: Partial<RequestOptions>) => {
  return request<T>({
    method: 'POST',
    url,
    data,
    ...options,
  })
}

export const put = <T = any>(url: string, data?: any, options?: Partial<RequestOptions>) => {
  return request<T>({
    method: 'PUT',
    url,
    data,
    ...options,
  })
}

export const del = <T = any>(url: string, data?: any, options?: Partial<RequestOptions>) => {
  return request<T>({
    method: 'DELETE',
    url,
    data,
    ...options,
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
    const response = await uni.uploadFile({
      url: `${BASE_URL}${options.url}`,
      filePath: options.filePath,
      name: options.name || 'file',
      formData: options.formData,
      header: {
        Authorization: token ? `Bearer ${token}` : '',
        ...options.header,
      },
    })

    const data = JSON.parse(response.data)

    if (data.code !== SUCCESS_CODE) {
      throw new Error(data.message || 'Upload failed')
    }

    return data
  } catch (error: any) {
    uni.showToast({
      title: error.message || 'Upload failed',
      icon: 'none',
    })
    throw error
  } finally {
    hideLoading()
  }
}

// Debug function for testing auth flow
export const testAuthFlow = async (): Promise<boolean> => {
  try {
    console.log('Testing auth flow...')

    // Test if token exists
    const token = getToken()
    console.log('Current token:', token ? 'exists' : 'missing')

    // Make a test request to authenticated endpoint
    const response = await get('/me')
    console.log('Auth test successful:', response)
    return true
  } catch (error) {
    console.error('Auth test failed:', error)
    return false
  }
}

// Export types
export * from '@/types/api'
