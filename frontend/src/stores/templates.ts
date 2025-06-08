/**
 * Templates Store - Template Gallery Management
 */

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Template, GetTemplatesResponse } from '@/types/api'
import { templates as templatesApi } from '@/api'

export const useTemplateStore = defineStore('templates', () => {
  // State
  const templates = ref<Template[]>([])
  const featuredTemplates = ref<Template[]>([])
  const currentTemplate = ref<Template | null>(null)
  const isLoading = ref(false)
  const error = ref<string | null>(null)
  
  // Pagination
  const currentPage = ref(1)
  const pageSize = ref(20)
  const totalCount = ref(0)
  const hasMore = ref(true)
  
  // Cache
  const lastFetchTime = ref<number>(0)
  const CACHE_DURATION = 5 * 60 * 1000 // 5 minutes

  // Getters
  const activeTemplates = computed(() => 
    templates.value.filter(t => t.is_active)
  )
  
  const templateById = computed(() => (id: number) => 
    templates.value.find(t => t.id === id)
  )
  
  const templatesSortedByCost = computed(() => 
    [...activeTemplates.value].sort((a, b) => a.credit_cost - b.credit_cost)
  )
  
  const isCacheValid = computed(() => 
    Date.now() - lastFetchTime.value < CACHE_DURATION
  )

  // Actions
  const fetchTemplates = async (force = false) => {
    // Use cache if valid and not forced
    if (!force && isCacheValid.value && templates.value.length > 0) {
      return templates.value
    }
    
    isLoading.value = true
    error.value = null
    
    try {
      const response = await templatesApi.getTemplates({
        page: currentPage.value,
        pageSize: pageSize.value,
        isActive: true
      })
      
      const data = response.data as GetTemplatesResponse
      templates.value = data.list
      totalCount.value = data.total
      hasMore.value = data.hasMore
      lastFetchTime.value = Date.now()
      
      // Cache in storage
      uni.setStorageSync('templates_cache', JSON.stringify({
        templates: templates.value,
        timestamp: lastFetchTime.value
      }))
      
      return templates.value
    } catch (e: any) {
      error.value = e.message || 'Failed to fetch templates'
      throw e
    } finally {
      isLoading.value = false
    }
  }

  const fetchMoreTemplates = async () => {
    if (!hasMore.value || isLoading.value) return
    
    currentPage.value++
    isLoading.value = true
    error.value = null
    
    try {
      const response = await templatesApi.getTemplates({
        page: currentPage.value,
        pageSize: pageSize.value,
        isActive: true
      })
      
      const data = response.data as GetTemplatesResponse
      templates.value = [...templates.value, ...data.list]
      hasMore.value = data.hasMore
      
      return data.list
    } catch (e: any) {
      error.value = e.message || 'Failed to fetch more templates'
      currentPage.value-- // Revert page on error
      throw e
    } finally {
      isLoading.value = false
    }
  }

  const fetchFeaturedTemplates = async () => {
    try {
      const response = await templatesApi.getFeaturedTemplates()
      featuredTemplates.value = response.data
      return featuredTemplates.value
    } catch (e: any) {
      console.error('Failed to fetch featured templates:', e)
      // Don't throw, just return empty array
      return []
    }
  }

  const selectTemplate = async (id: number) => {
    // Check if we already have it
    let template = templateById.value(id)
    
    if (!template) {
      // Fetch individual template
      isLoading.value = true
      error.value = null
      
      try {
        const response = await templatesApi.getTemplateById(id)
        template = response.data
        
        // Add to templates array if not exists
        const index = templates.value.findIndex(t => t.id === id)
        if (index === -1) {
          templates.value.push(template)
        } else {
          templates.value[index] = template
        }
      } catch (e: any) {
        error.value = e.message || 'Failed to fetch template'
        throw e
      } finally {
        isLoading.value = false
      }
    }
    
    currentTemplate.value = template
    return template
  }

  const clearCurrentTemplate = () => {
    currentTemplate.value = null
  }

  const initFromCache = () => {
    try {
      const cached = uni.getStorageSync('templates_cache')
      if (cached) {
        const { templates: cachedTemplates, timestamp } = JSON.parse(cached)
        
        // Check if cache is still valid
        if (Date.now() - timestamp < CACHE_DURATION) {
          templates.value = cachedTemplates
          lastFetchTime.value = timestamp
        }
      }
    } catch (e) {
      console.error('Failed to init templates from cache:', e)
    }
  }

  const resetPagination = () => {
    currentPage.value = 1
    hasMore.value = true
    templates.value = []
  }

  // Search templates
  const searchTemplates = (query: string) => {
    const searchQuery = query.toLowerCase()
    return templates.value.filter(t => 
      t.name.toLowerCase().includes(searchQuery) ||
      t.description?.toLowerCase().includes(searchQuery)
    )
  }

  return {
    // State
    templates,
    featuredTemplates,
    currentTemplate,
    isLoading,
    error,
    currentPage,
    pageSize,
    totalCount,
    hasMore,
    
    // Getters
    activeTemplates,
    templateById,
    templatesSortedByCost,
    isCacheValid,
    
    // Actions
    fetchTemplates,
    fetchMoreTemplates,
    fetchFeaturedTemplates,
    selectTemplate,
    clearCurrentTemplate,
    initFromCache,
    resetPagination,
    searchTemplates
  }
})

// Export type for easy access
export type TemplateStore = ReturnType<typeof useTemplateStore> 