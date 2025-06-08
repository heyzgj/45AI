<template>
  <view 
    class="template-card card interactive"
    :class="{ 'loading': loading, 'error': error }"
    @click="handleClick"
  >
    <!-- Loading State -->
    <view v-if="loading" class="template-loading">
      <view class="loading-blob"></view>
    </view>
    
    <!-- Error State -->
    <view v-else-if="error" class="template-error">
      <text class="error-icon">⚠️</text>
      <text class="error-text">Failed to load</text>
    </view>
    
    <!-- Content -->
    <template v-else>
      <!-- Preview Image -->
      <view class="template-image-wrapper">
        <image 
          :src="template.preview_image_url" 
          class="template-image"
          mode="aspectFill"
          :lazy-load="true"
          @error="handleImageError"
        />
        
        <!-- Hover Overlay -->
        <view class="template-overlay">
          <view class="overlay-content">
            <text class="overlay-text">Use This Style</text>
          </view>
        </view>
      </view>
      
      <!-- Template Info -->
      <view class="template-info">
        <text class="template-name">{{ template.name }}</text>
        <view class="template-footer">
          <text class="template-description">{{ template.description }}</text>
          <view class="template-cost">
            <text class="cost-icon">🎞️</text>
            <text class="cost-amount">{{ template.credit_cost }}</text>
          </view>
        </view>
      </view>
    </template>
  </view>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import type { Template } from '@/types/api'

// Props
interface Props {
  template: Template
  loading?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  loading: false
})

// State
const error = ref(false)

// Methods
const handleClick = () => {
  if (!props.loading && !error.value) {
    uni.navigateTo({
      url: `/pages/template-detail/index?id=${props.template.id}`
    });
  }
}

const handleImageError = () => {
  error.value = true
}
</script>

<style lang="scss" scoped>
@import '@/styles/variables.scss';

.template-card {
  position: relative;
  overflow: hidden;
  padding: 0;
  border-radius: $radius-md;
  background: $color-surface;
  transition: all $duration-fast $ease-custom;
  
  &.loading,
  &.error {
    min-height: 360px;
    display: flex;
    align-items: center;
    justify-content: center;
  }
  
  &:hover {
    .template-overlay {
      opacity: 1;
    }
    
    .template-image {
      transform: scale(1.05);
    }
  }
  
  &:active {
    transform: translateY(-4px) scale(1.02);
    box-shadow: $shadow-hover;
    
    .template-overlay {
      opacity: 0.9;
    }
  }
}

// Loading State
.template-loading {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: $spacing-xl;
  
  .loading-blob {
    width: 60px;
    height: 60px;
    background: linear-gradient(135deg, $color-secondary, $color-primary);
    border-radius: 50%;
    animation: pulseSoft 2s $ease-custom infinite;
  }
}

// Error State
.template-error {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: $spacing-xl;
  
  .error-icon {
    font-size: 32px;
    margin-bottom: $spacing-sm;
    opacity: 0.5;
  }
  
  .error-text {
    font-size: $font-size-caption;
    color: $color-text-subtle;
  }
}

// Image Wrapper
.template-image-wrapper {
  position: relative;
  overflow: hidden;
  height: 280px;
  background: $color-bg;
}

.template-image {
  width: 100%;
  height: 100%;
  display: block;
  transition: transform $duration-normal $ease-custom;
}

// Overlay
.template-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity $duration-fast $ease-custom;
  backdrop-filter: blur(4px);
}

.overlay-content {
  text-align: center;
  transform: translateY(10px);
  transition: transform $duration-fast $ease-custom;
  
  .template-overlay:hover & {
    transform: translateY(0);
  }
}

.overlay-text {
  color: white;
  font-size: 16px;
  font-weight: $font-weight-semibold;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
}

// Template Info
.template-info {
  padding: $spacing-md;
}

.template-name {
  display: block;
  font-size: 16px;
  font-weight: $font-weight-semibold;
  color: $color-text;
  margin-bottom: $spacing-xs;
  
  // Truncate long names
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.template-footer {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  gap: $spacing-sm;
}

.template-description {
  font-size: $font-size-caption;
  color: $color-text-subtle;
  flex: 1;
  
  // Limit to 2 lines
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.template-cost {
  display: flex;
  align-items: center;
  gap: 4px;
  color: $color-primary;
  font-weight: $font-weight-semibold;
  flex-shrink: 0;
}

.cost-icon {
  font-size: 16px;
}

.cost-amount {
  font-size: 14px;
}

// Responsive adjustments
@media (min-width: 768px) {
  .template-image-wrapper {
    height: 320px;
  }
}
</style> 