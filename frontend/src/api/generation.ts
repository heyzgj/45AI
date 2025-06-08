/**
 * Image Generation API
 */

import { post, get, upload } from '@/utils/request'
import type { 
  GenerateImageRequest, 
  GenerateImageResponse,
  GetGenerationStatusRequest,
  GetGenerationStatusResponse 
} from '@/types/api'

// Generate image
export const generateImage = (data: GenerateImageRequest) => {
  return post<GenerateImageResponse>('/generation/create', data)
}

// Get generation status
export const getGenerationStatus = (taskId: string) => {
  return get<GetGenerationStatusResponse>(`/generation/status/${taskId}`)
}

// Upload user photo
export const uploadUserPhoto = (filePath: string) => {
  return upload({
    url: '/generation/upload',
    filePath,
    name: 'photo'
  })
}

// Get generation history
export const getGenerationHistory = (params?: { page?: number; pageSize?: number }) => {
  return get('/generation/history', params)
} 