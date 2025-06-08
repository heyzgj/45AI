<template>
  <view class="gallery-container page-wrapper animate-fade-in">
    <!-- Header -->
    <view class="gallery-header">
      <text class="gallery-title">Choose Your Style</text>
      <text class="gallery-subtitle">Select a template to transform your photo</text>
    </view>
    
    <!-- Template Grid -->
    <view class="template-grid">
      <TemplateCard
        v-for="(template, index) in templates" 
        :key="template.id"
        :template="template"
        class="animate-slide-up"
        :style="{ animationDelay: `${index * 50}ms` }"
        @click="selectTemplate"
      />
    </view>
    
    <!-- Loading State -->
    <view v-if="loading" class="loading-container">
      <LoadingAnimation 
        variant="blob" 
        size="large" 
        text="Loading beautiful styles..."
      />
    </view>
    
    <!-- Empty State -->
    <view v-if="!loading && templates.length === 0" class="empty-state">
      <text class="empty-icon">🎨</text>
      <text class="empty-title">No Templates Available</text>
      <text class="empty-message">Check back soon for new styles!</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useTemplateStore } from '@/stores/templates'
import { customTransitions } from '@/utils/navigation'
import TemplateCard from '@/components/TemplateCard/TemplateCard.vue'
import LoadingAnimation from '@/components/LoadingAnimation/LoadingAnimation.vue'

// Mock template data for now
const templates = ref([
  {
    id: 1,
    name: 'Dusty Rose Dream',
    description: 'Soft, romantic vibes',
    preview_image_url: 'https://placeholder.com/400x600',
    credit_cost: 15
  },
  {
    id: 2,
    name: 'Fairy Tale Magic',
    description: 'Enchanted forest aesthetic',
    preview_image_url: 'https://placeholder.com/400x600',
    credit_cost: 20
  },
  {
    id: 3,
    name: 'Vintage Film',
    description: 'Classic cinema look',
    preview_image_url: 'https://placeholder.com/400x600',
    credit_cost: 10
  },
  {
    id: 4,
    name: 'Modern Minimalist',
    description: 'Clean and sophisticated',
    preview_image_url: 'https://placeholder.com/400x600',
    credit_cost: 15
  },
  {
    id: 5,
    name: 'Golden Hour',
    description: 'Warm sunset glow',
    preview_image_url: 'https://placeholder.com/400x600',
    credit_cost: 12
  },
  {
    id: 6,
    name: 'Neon Dreams',
    description: 'Vibrant cyberpunk style',
    preview_image_url: 'https://placeholder.com/400x600',
    credit_cost: 25
  }
])

const loading = ref(false)

const selectTemplate = (template: any) => {
  // Navigate to generation page with selected template
  customTransitions.toGenerate(template.id)
}

onMounted(() => {
  // TODO: Load templates from API
})
</script>

<style lang="scss" scoped>
@import '@/styles/variables.scss';

.gallery-container {
  background-color: $color-bg;
  min-height: 100vh;
}

.gallery-header {
  text-align: center;
  padding: $spacing-xl $page-padding;
}

.gallery-title {
  display: block;
  font-size: $font-size-h1;
  font-weight: $font-weight-semibold;
  color: $color-text;
  margin-bottom: $spacing-xs;
}

.gallery-subtitle {
  display: block;
  font-size: $font-size-body;
  color: $color-text-subtle;
}

.template-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: $spacing-md;
  padding: 0 $page-padding $spacing-xl;
}

// Responsive adjustments
@media (min-width: 768px) {
  .template-grid {
    grid-template-columns: repeat(3, 1fr);
    gap: $spacing-lg;
  }
}
</style> 