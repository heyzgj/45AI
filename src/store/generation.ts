/**
 * Generation Store - AI Image Generation Management
 */

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { GenerateImageResponse, GetGenerationStatusResponse } from '@/types/api'
import * as generationApi from '@/api/generation'
import { useUserStore } from './user'
import { useTemplateStore } from './templates'
import { getEnvBaseUploadUrl } from '@/utils/index'

interface GenerationTask {
  id: string
  templateId: number
  templateName: string
  status: 'pending' | 'processing' | 'completed' | 'failed' | 'succeeded'
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
  const pollingInterval = ref<any>(null)

  // Constants
  const POLLING_INTERVAL = 3000 // 3 seconds
  const MAX_HISTORY = 20

  // Getters
  const hasActiveTask = computed(
    () => currentTask.value && ['pending', 'processing'].includes(currentTask.value.status),
  )

  const completedTasks = computed(() => taskHistory.value.filter((t) => t.status === 'completed'))

  const estimatedTimeRemaining = computed(() => {
    if (!currentTask.value || currentTask.value.status !== 'processing') return 0

    const elapsed = Date.now() - currentTask.value.startedAt
    const progress = currentTask.value.progress || 0

    if (progress === 0) return 30000 // Default 30s

    const totalEstimated = elapsed / (progress / 100)
    return Math.max(0, totalEstimated - elapsed)
  })

  // Actions
  const startGeneration = async (templateId: number, filePath: string) => {
    const userStore = useUserStore()
    const templateStore = useTemplateStore()

    if (!filePath) {
      throw new Error('Please select a photo first')
    }

    // Check credits
    const template = await templateStore.selectTemplate(templateId)
    if (!template) {
      throw new Error('Template not found')
    }

    if (!userStore.hasCredits || userStore.credits < template.credit_cost) {
      throw new Error('Insufficient credits')
    }

    isGenerating.value = true
    uploadProgress.value = 0

    try {
      // Wrap uni.uploadFile in a Promise for better error handling
      console.log('=== DEBUGGING ENVIRONMENT VARIABLES ===')
      console.log('import.meta.env.VITE_UPLOAD_BASEURL:', import.meta.env.VITE_UPLOAD_BASEURL)
      console.log(
        'typeof import.meta.env.VITE_UPLOAD_BASEURL:',
        typeof import.meta.env.VITE_UPLOAD_BASEURL,
      )
      console.log('import.meta.env:', import.meta.env)

      const uploadUrl = getEnvBaseUploadUrl()
      console.log('getEnvBaseUploadUrl() returned:', uploadUrl)
      console.log('typeof uploadUrl:', typeof uploadUrl)

      // Failsafe: if uploadUrl is undefined or 'undefined' string, use hardcoded value
      let finalUploadUrl = uploadUrl
      if (!finalUploadUrl || finalUploadUrl === 'undefined' || finalUploadUrl === undefined) {
        console.warn('WARNING: uploadUrl is invalid, using fallback')
        finalUploadUrl = 'http://localhost:8080/api/v1/generate'
      }

      console.log('Final URL to be used:', finalUploadUrl)
      console.log('Platform check - isH5:', true) // We're in H5 mode

      if (!finalUploadUrl || finalUploadUrl === 'undefined') {
        throw new Error('Upload URL is undefined or invalid after failsafe')
      }

      // H5-specific file upload handling
      // uni.uploadFile has issues in H5 mode, use XMLHttpRequest instead
      console.log('Using H5-specific file upload method instead of uni.uploadFile')

      await new Promise<void>((resolve, reject) => {
        // Create FormData for H5 upload
        const formData = new FormData()

        // Convert filePath to File object (H5 specific)
        fetch(filePath)
          .then((response) => response.blob())
          .then((blob) => {
            const file = new File([blob], 'upload.jpg', { type: 'image/jpeg' })
            formData.append('image', file)
            formData.append('template_id', templateId.toString())

            // Use XMLHttpRequest for H5 compatibility
            const xhr = new XMLHttpRequest()

            xhr.upload.onprogress = (event) => {
              if (event.lengthComputable) {
                uploadProgress.value = Math.round((event.loaded / event.total) * 100)
              }
            }

            xhr.onload = () => {
              try {
                console.log('H5 Upload success response:', xhr.responseText)

                if (xhr.status === 202) {
                  // Accepted - this is the expected response for async generation
                  const response = JSON.parse(xhr.responseText)
                  if (!response.job_id) {
                    reject(new Error(response.message || 'No job_id returned'))
                    return
                  }

                  const jobId = response.job_id

                  const task: GenerationTask = {
                    id: jobId,
                    templateId,
                    templateName: template.name,
                    status: 'pending',
                    progress: 0,
                    startedAt: Date.now(),
                  }

                  currentTask.value = task
                  addToHistory(task)
                  userStore.deductCredits(template.credit_cost)
                  startPolling(jobId)
                  resolve()
                } else {
                  // Handle other status codes
                  const response = JSON.parse(xhr.responseText)
                  reject(
                    new Error(
                      response.error || response.message || `Unexpected status: ${xhr.status}`,
                    ),
                  )
                }
              } catch (parseError: any) {
                console.error('Failed to parse H5 upload response:', parseError)
                reject(new Error(parseError.message || 'Failed to parse response'))
              }
            }

            xhr.onerror = () => {
              console.error('H5 Upload failed with network error')
              reject(new Error('Network error during upload'))
            }

            xhr.open('POST', finalUploadUrl)
            xhr.setRequestHeader('Authorization', `Bearer ${userStore.token}`)

            console.log('H5: Starting XMLHttpRequest upload to:', finalUploadUrl)
            xhr.send(formData)
          })
          .catch((error) => {
            console.error('Failed to convert filePath to blob:', error)
            reject(new Error('Failed to process file for upload'))
          })
      })
    } catch (e: any) {
      isGenerating.value = false
      throw e
    }
  }

  const startPolling = async (jobId: string) => {
    // Clear any existing polling
    stopPolling()

    // Validate jobId
    if (!jobId) {
      console.error('startPolling: jobId is required')
      return
    }

    // Poll for status updates
    pollingInterval.value = setInterval(async () => {
      try {
        const response = await generationApi.getGenerationStatus(jobId)

        // Handle potential non-200 responses gracefully
        if (!response || !response.data) {
          console.warn(`Invalid response for job ${jobId}, skipping update.`)
          return
        }

        const { status, progress, error } = response.data

        if (currentTask.value && currentTask.value.id === jobId) {
          currentTask.value.status = status
          currentTask.value.progress = progress || currentTask.value.progress || 0

          if (status === 'completed') {
            // Get the final result with image URL
            try {
              if (jobId) {
                const resultResponse = await generationApi.getGenerationResult(jobId)
                if (resultResponse.data && resultResponse.data.image_url) {
                  let imageUrl = resultResponse.data.image_url

                  // Convert relative URL to full URL for H5 mode
                  if (imageUrl.startsWith('/')) {
                    const baseUrl = import.meta.env.VITE_SERVER_BASEURL || 'http://localhost:8080'
                    // Remove /api/v1 suffix if present for image URLs
                    const serverBase = baseUrl.replace('/api/v1', '')
                    imageUrl = `${serverBase}${imageUrl}`
                  }

                  currentTask.value.resultUrl = imageUrl
                  console.log('Set result URL:', imageUrl)
                }
              }
            } catch (resultError) {
              console.error('Failed to get generation result:', resultError)
            }

            currentTask.value.completedAt = Date.now()
            stopPolling()
            isGenerating.value = false

            // Update in history
            updateInHistory(currentTask.value)

            // Show success notification
            uni.showToast({
              title: '生成完成!',
              icon: 'success',
            })
          } else if (status === 'failed') {
            currentTask.value.errorMessage = error
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
              title: error || '生成失败',
              icon: 'none',
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
      taskHistory.value.pop()
    }

    saveHistoryToStorage()
  }

  const updateInHistory = (task: GenerationTask) => {
    const index = taskHistory.value.findIndex((t) => t.id === task.id)
    if (index !== -1) {
      taskHistory.value[index] = { ...task }
      saveHistoryToStorage()
    }
  }

  const saveHistoryToStorage = () => {
    uni.setStorageSync('generation_history', JSON.stringify(taskHistory.value))
  }

  const loadHistoryFromStorage = () => {
    const history = uni.getStorageSync('generation_history')
    if (history) {
      taskHistory.value = JSON.parse(history)
    }
  }

  const clearUserPhoto = () => {
    userPhotoUrl.value = ''
    uploadProgress.value = 0
  }

  const cancelGeneration = () => {
    if (!currentTask.value) return

    // Stop polling
    stopPolling()

    // Call API to cancel if possible (not implemented in this example)
    // generationApi.cancelGeneration(currentTask.value.id)

    // Visually update
    isGenerating.value = false

    // Refund credits
    const userStore = useUserStore()
    const templateStore = useTemplateStore()
    const template = templateStore.templateById(currentTask.value.templateId)
    if (template) {
      userStore.updateCredits(template.credit_cost)
    }

    // Update task status
    currentTask.value.status = 'failed' // or 'cancelled'
    currentTask.value.errorMessage = 'User cancelled'
    updateInHistory(currentTask.value)

    currentTask.value = null

    uni.showToast({
      title: 'Generation cancelled',
      icon: 'none',
    })
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
    startGeneration,
    stopPolling,
    clearUserPhoto,
    cancelGeneration,
    loadHistoryFromStorage,
  }
})

// Export type for easy access
export type GenerationStore = ReturnType<typeof useGenerationStore>
