/**
 * Navigation utilities and guards for 45AI
 */

import { useUserStore } from '@/stores/user'

// Routes that require authentication
const AUTH_REQUIRED_ROUTES = [
  '/pages/generate/generate',
  '/pages/profile/profile',
  '/pages/purchase/purchase',
  '/pages/history/history',
  '/pages/result/result'
]

// Routes that should skip auth (accessible when logged in)
const PUBLIC_ROUTES = [
  '/pages/index/index',
  '/pages/gallery/gallery',
  '/pages/login/login'
]

// Navigation options type
interface NavigateOptions {
  url: string
  animationType?: 'slide-in-right' | 'slide-in-left' | 'slide-in-top' | 'slide-in-bottom' | 'fade-in' | 'zoom-out' | 'zoom-fade-out' | 'none'
  animationDuration?: number
  events?: Record<string, Function>
  success?: (result: any) => void
  fail?: (result: any) => void
  complete?: (result: any) => void
}

// Check if route requires authentication
const requiresAuth = (url: string): boolean => {
  const path = url.split('?')[0] // Remove query params
  return AUTH_REQUIRED_ROUTES.some(route => path.includes(route))
}

// Check if user is authenticated
const isAuthenticated = (): boolean => {
  const userStore = useUserStore()
  return !!userStore.token && !!userStore.userInfo
}

// Enhanced navigation with auth check and transitions
export const navigateTo = (options: NavigateOptions | string) => {
  const opts: NavigateOptions = typeof options === 'string' 
    ? { url: options } 
    : options

  // Default animation settings
  opts.animationType = opts.animationType || 'slide-in-right'
  opts.animationDuration = opts.animationDuration || 400

  // Check authentication
  if (requiresAuth(opts.url) && !isAuthenticated()) {
    // Redirect to login with return URL
    const currentPage = getCurrentPages().pop()
    const returnUrl = currentPage ? currentPage.route : ''
    
    uni.showToast({
      title: 'Please login first',
      icon: 'none',
      duration: 1500
    })

    setTimeout(() => {
      uni.navigateTo({
        url: `/pages/login/login?returnUrl=${encodeURIComponent(opts.url)}`,
        animationType: 'slide-in-bottom',
        animationDuration: 400
      })
    }, 500)
    
    return
  }

  // Navigate with animation
  uni.navigateTo({
    ...opts,
    fail: (err) => {
      console.error('Navigation failed:', err)
      if (opts.fail) opts.fail(err)
    }
  })
}

// Navigate back with animation
export const navigateBack = (delta: number = 1, animationType = 'slide-out-right') => {
  uni.navigateBack({
    delta,
    animationType,
    animationDuration: 400
  })
}

// Redirect with fade animation
export const redirectTo = (options: NavigateOptions | string) => {
  const opts: NavigateOptions = typeof options === 'string' 
    ? { url: options } 
    : options

  // Check authentication
  if (requiresAuth(opts.url) && !isAuthenticated()) {
    uni.redirectTo({
      url: '/pages/login/login'
    })
    return
  }

  uni.redirectTo(opts)
}

// Relaunch app
export const reLaunch = (url: string) => {
  uni.reLaunch({ url })
}

// Switch tab with custom animation
export const switchTab = (url: string) => {
  // Add fade effect for tab switching
  const pages = getCurrentPages()
  const currentPage = pages[pages.length - 1]
  
  // Apply fade out to current page
  if (currentPage && currentPage.$vm) {
    currentPage.$vm.$el?.classList.add('animate-fade-out')
  }
  
  setTimeout(() => {
    uni.switchTab({ url })
  }, 200)
}

// Get current page path
export const getCurrentPath = (): string => {
  const pages = getCurrentPages()
  if (pages.length === 0) return ''
  
  const currentPage = pages[pages.length - 1]
  return `/${currentPage.route}`
}

// Navigation interceptor setup
export const setupNavigationInterceptor = () => {
  // Intercept uni.navigateTo
  const originalNavigateTo = uni.navigateTo
  uni.navigateTo = function(options: any) {
    // Apply our navigation logic
    navigateTo(options)
  }

  // Listen for page show to apply animations
  uni.addInterceptor('navigateTo', {
    success() {
      // Page successfully navigated
      const pages = getCurrentPages()
      const currentPage = pages[pages.length - 1]
      
      // Apply enter animation
      if (currentPage && currentPage.$vm) {
        currentPage.$vm.$el?.classList.add('animate-fade-in')
      }
    }
  })
}

// Auth state change handler
export const onAuthStateChange = (isLoggedIn: boolean) => {
  const currentPath = getCurrentPath()
  
  if (!isLoggedIn && requiresAuth(currentPath)) {
    // User logged out on protected page
    redirectTo('/pages/login/login')
  } else if (isLoggedIn && currentPath === '/pages/login/login') {
    // User logged in on login page
    const pages = getCurrentPages()
    const loginPage = pages[pages.length - 1]
    const returnUrl = loginPage.options?.returnUrl
    
    if (returnUrl) {
      redirectTo(decodeURIComponent(returnUrl))
    } else {
      switchTab('/pages/index/index')
    }
  }
}

// Page transition mixins for Vue components
export const pageTransitionMixin = {
  onShow() {
    // Apply page show animation
    if (this.$el) {
      this.$el.classList.add('animate-fade-in')
    }
  },
  onHide() {
    // Prepare for hide animation
    if (this.$el) {
      this.$el.classList.add('animate-fade-out')
    }
  },
  onUnload() {
    // Clean up animations
    if (this.$el) {
      this.$el.classList.remove('animate-fade-in', 'animate-fade-out')
    }
  }
}

// Custom page transitions for specific routes
export const customTransitions = {
  toGallery: () => navigateTo({
    url: '/pages/gallery/gallery',
    animationType: 'slide-in-right',
    animationDuration: 400
  }),
  
  toGenerate: (templateId: number) => navigateTo({
    url: `/pages/generate/generate?templateId=${templateId}`,
    animationType: 'slide-in-bottom',
    animationDuration: 500
  }),
  
  toResult: (taskId: string) => navigateTo({
    url: `/pages/result/result?taskId=${taskId}`,
    animationType: 'zoom-fade-out',
    animationDuration: 600
  }),
  
  toLogin: (returnUrl?: string) => navigateTo({
    url: returnUrl ? `/pages/login/login?returnUrl=${encodeURIComponent(returnUrl)}` : '/pages/login/login',
    animationType: 'slide-in-bottom',
    animationDuration: 400
  })
} 