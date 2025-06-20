/**
 * Templates API
 */

import { get } from '@/utils/request'
import type { GetTemplatesRequest, GetTemplatesResponse, Template } from '@/types/api'

// Get templates list
export const getTemplates = (params?: GetTemplatesRequest) => {
  return get<GetTemplatesResponse>('/templates', params)
}

// Get template by ID
export const getTemplateById = (id: number) => {
  return get<Template>(`/templates/${id}`)
}

// Get featured templates
export const getFeaturedTemplates = () => {
  return get<Template[]>('/templates/featured')
} 