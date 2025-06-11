/**
 * Image Generation API
 */

import { post, get, upload } from '@/utils/request'
import type {
  GenerateImageRequest,
  GenerateImageResponse,
  GetGenerationStatusRequest,
  GetGenerationStatusResponse,
} from '@/types/api'

// Generate image (async - returns job_id)
export const generateImage = (data: GenerateImageRequest) => {
  return post<GenerateImageResponse>('/generate', data)
}

// Get generation status
export const getGenerationStatus = (jobId: string) => {
  return get<GetGenerationStatusResponse>(`/generate/${jobId}/status`)
}

// Get generation result
export const getGenerationResult = (jobId: string) => {
  return get(`/generate/${jobId}`)
}

// Upload user photo
export const uploadUserPhoto = (filePath: string) => {
  return upload({
    url: '/generation/upload',
    filePath,
    name: 'photo',
  })
}

// Get generation history
export const getGenerationHistory = (params?: { page?: number; pageSize?: number }) => {
  return get('/generation/history', params)
}
