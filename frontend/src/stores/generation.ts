/**
 * Generation Store - AI Image Generation Management
 */

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { GenerateImageResponse, GetGenerationStatusResponse } from '@/types/api'
import { generation as generationApi } from '@/api'
import { useUserStore } from './user'
import { useTemplateStore } from './templates'

interface GenerationTask {
  id: string
  templateId: number
  templateName: string
  status: 'pending' | 'processing' | 'completed' | 'failed'
  resultUrl?: string
  errorMessage?: string
  progress: number
  startedAt: number
  completedAt?: number
}

export const useGenerationStore = defineStore('generation', () => {
  // State
  const currentTask = ref<GenerationTask | null>(null)
  const taskHistory = ref<GenerationTask[]>([])
  const userPhotoUrl = ref<string>('')
  const isGenerating = ref(false)
  const uploadProgress = ref(0)
  const pollingInterval = ref<number | null>(null)
  
  // Constants
  const POLLING_INTERVAL = 3000 // 3 seconds
  const MAX_HISTORY = 20

  // Getters
  const hasActiveTask = computed(() => 
    currentTask.value && ['pending', 'processing'].includes(currentTask.value.status)
  )
  
  const completedTasks = computed(() => 
    taskHistory.value.filter(t => t.status === 'completed')
  )
  
  const estimatedTimeRemaining = computed(() => {
    if (!currentTask.value || currentTask.value.status !== 'processing') return 0
    
    const elapsed = Date.now() - currentTask.value.startedAt
    const progress = currentTask.value.progress || 0
    
    if (progress === 0) return 30000 // Default 30s
    
    const totalEstimated = elapsed / (progress / 100)
    return Math.max(0, totalEstimated - elapsed)
  })

  // Actions
  const uploadPhoto = async (filePath: string) => {
    uploadProgress.value = 0
    
    try {
      // Show upload progress
      const uploadTask = uni.uploadFile({
        url: `${process.env.VUE_APP_API_URL}/generation/upload`,
        filePath,
        name: 'photo',
        header: {
          'Authorization': `Bearer ${useUserStore().token}`
        },
        success: (res) => {
          const data = JSON.parse(res.data)
          if (data.code === 0) {
            userPhotoUrl.value = data.data.url
          } else {
            throw new Error(data.message)
          }
        },
        fail: (err) => {
          throw new Error(err.errMsg)
        }
      })
      
      uploadTask.onProgressUpdate((res) => {
        uploadProgress.value = res.progress
      })
      
      return new Promise((resolve, reject) => {
        uploadTask.onHeadersReceived(() => resolve(userPhotoUrl.value))
        uploadTask.onProgressUpdate((res) => {
          if (res.progress === 100) {
            setTimeout(() => resolve(userPhotoUrl.value), 500)
          }
        })
      })
    } catch (e: any) {
      throw new Error(e.message || 'Failed to upload photo')
    }
  }

  const startGeneration = async (templateId: number) => {
    const userStore = useUserStore()
    const templateStore = useTemplateStore()
    
    // Check credits
    const template = await templateStore.selectTemplate(templateId)
    if (!template) {
      throw new Error('Template not found')
    }
    
    if (!userStore.hasCredits || userStore.credits < template.credit_cost) {
      throw new Error('Insufficient credits')
    }
    
    if (!userPhotoUrl.value) {
      throw new Error('Please upload a photo first')
    }
    
    isGenerating.value = true
    
    try {
      // Call generation API
      const response = await generationApi.generateImage({
        template_id: templateId,
        user_photo_url: userPhotoUrl.value
      })
      
      const { task_id, status, estimated_time } = response.data
      
      // Create task
      const task: GenerationTask = {
        id: task_id,
        templateId,
        templateName: template.name,
        status,
        progress: 0,
        startedAt: Date.now()
      }
      
      currentTask.value = task
      addToHistory(task)
      
      // Deduct credits immediately
      userStore.deductCredits(template.credit_cost)
      
      // Start polling for status
      startPolling(task_id)
      
      return task
    } catch (e: any) {
      isGenerating.value = false
      throw e
    }
  }

  const startPolling = (taskId: string) => {
    // Clear any existing polling
    stopPolling()
    
    // Poll for status updates
    pollingInterval.value = setInterval(async () => {
      try {
        const response = await generationApi.getGenerationStatus(taskId)
        const { status, result_url, error_message, progress } = response.data
        
        if (currentTask.value && currentTask.value.id === taskId) {
          currentTask.value.status = status
          currentTask.value.progress = progress || 0
          
          if (status === 'completed' && result_url) {
            currentTask.value.resultUrl = result_url
            currentTask.value.completedAt = Date.now()
            stopPolling()
            isGenerating.value = false
            
            // Update in history
            updateInHistory(currentTask.value)
            
            // Show success notification
            uni.showToast({
              title: 'Generation completed!',
              icon: 'success'
            })
          } else if (status === 'failed') {
            currentTask.value.errorMessage = error_message
            currentTask.value.completedAt = Date.now()
            stopPolling()
            isGenerating.value = false
            
            // Update in history
            updateInHistory(currentTask.value)
            
            // Refund credits on failure
            const userStore = useUserStore()
            const templateStore = useTemplateStore()
            const template = templateStore.templateById(currentTask.value.templateId)
            if (template) {
              userStore.updateCredits(template.credit_cost)
            }
            
            // Show error notification
            uni.showToast({
              title: error_message || 'Generation failed',
              icon: 'none'
            })
          }
        }
      } catch (e) {
        console.error('Polling error:', e)
      }
    }, POLLING_INTERVAL)
  }

  const stopPolling = () => {
    if (pollingInterval.value) {
      clearInterval(pollingInterval.value)
      pollingInterval.value = null
    }
  }

  const addToHistory = (task: GenerationTask) => {
    taskHistory.value.unshift({ ...task })
    
    // Limit history size
    if (taskHistory.value.length > MAX_HISTORY) {
      taskHistory.value = taskHistory.value.slice(0, MAX_HISTORY)
    }
    
    // Persist to storage
    saveHistoryToStorage()
  }

  const updateInHistory = (task: GenerationTask) => {
    const index = taskHistory.value.findIndex(t => t.id === task.id)
    if (index !== -1) {
      taskHistory.value[index] = { ...task }
      saveHistoryToStorage()
    }
  }

  const saveHistoryToStorage = () => {
    try {
      uni.setStorageSync('generation_history', JSON.stringify(taskHistory.value))
    } catch (e) {
      console.error('Failed to save generation history:', e)
    }
  }

  const loadHistoryFromStorage = () => {
    try {
      const saved = uni.getStorageSync('generation_history')
      if (saved) {
        taskHistory.value = JSON.parse(saved)
      }
    } catch (e) {
      console.error('Failed to load generation history:', e)
    }
  }

  const clearUserPhoto = () => {
    userPhotoUrl.value = ''
    uploadProgress.value = 0
  }

  const cancelGeneration = () => {
    stopPolling()
    isGenerating.value = false
    
    if (currentTask.value && ['pending', 'processing'].includes(currentTask.value.status)) {
      currentTask.value.status = 'failed'
      currentTask.value.errorMessage = 'Cancelled by user'
      currentTask.value.completedAt = Date.now()
      updateInHistory(currentTask.value)
    }
  }

  // Initialize
  loadHistoryFromStorage()

  return {
    // State
    currentTask,
    taskHistory,
    userPhotoUrl,
    isGenerating,
    uploadProgress,
    
    // Getters
    hasActiveTask,
    completedTasks,
    estimatedTimeRemaining,
    
    // Actions
    uploadPhoto,
    startGeneration,
    stopPolling,
    clearUserPhoto,
    cancelGeneration,
    loadHistoryFromStorage
  }
})

// Export type for easy access
export type GenerationStore = ReturnType<typeof useGenerationStore> 