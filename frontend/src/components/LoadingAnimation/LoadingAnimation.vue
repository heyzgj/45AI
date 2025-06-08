<template>
  <view 
    class="loading-animation"
    :class="[`variant-${variant}`, `size-${size}`]"
  >
    <!-- Dots Variant -->
    <view v-if="variant === 'dots'" class="loading-dots">
      <view class="dot"></view>
      <view class="dot"></view>
      <view class="dot"></view>
    </view>
    
    <!-- Blob Variant -->
    <view v-else-if="variant === 'blob'" class="loading-blob-container">
      <view class="loading-blob"></view>
    </view>
    
    <!-- Progress Variant -->
    <view v-else-if="variant === 'progress'" class="loading-progress">
      <view class="progress-circle">
        <svg viewBox="0 0 100 100" class="progress-svg">
          <!-- Background circle -->
          <circle
            cx="50"
            cy="50"
            r="45"
            fill="none"
            stroke="#F3D9D7"
            stroke-width="6"
          />
          <!-- Progress circle -->
          <circle
            cx="50"
            cy="50"
            r="45"
            fill="none"
            stroke="#E89B93"
            stroke-width="6"
            :stroke-dasharray="`${283 * (progress / 100)} 283`"
            stroke-dashoffset="0"
            stroke-linecap="round"
            transform="rotate(-90 50 50)"
            class="progress-bar"
          />
        </svg>
        <text class="progress-text">{{ Math.round(progress) }}%</text>
      </view>
    </view>
    
    <!-- Spinner Variant -->
    <view v-else-if="variant === 'spinner'" class="loading-spinner">
      <view class="spinner"></view>
    </view>
    
    <!-- Loading Text -->
    <text v-if="text" class="loading-text" :class="{ 'with-margin': variant !== 'dots' }">
      {{ text }}
    </text>
  </view>
</template>

<script setup lang="ts">
// Props
interface Props {
  variant?: 'dots' | 'blob' | 'progress' | 'spinner'
  size?: 'small' | 'medium' | 'large'
  text?: string
  progress?: number
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'dots',
  size: 'medium',
  text: '',
  progress: 0
})
</script>

<style lang="scss" scoped>
@import '@/styles/variables.scss';

.loading-animation {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: $spacing-md;
}

// Dots Variant
.loading-dots {
  display: flex;
  gap: $spacing-xs;
  
  .dot {
    border-radius: 50%;
    background: $color-secondary;
    animation: loadingDot 1.4s ease-in-out infinite;
    
    &:nth-child(1) {
      animation-delay: -0.32s;
    }
    
    &:nth-child(2) {
      animation-delay: -0.16s;
    }
  }
}

// Blob Variant
.loading-blob-container {
  position: relative;
}

.loading-blob {
  background: linear-gradient(135deg, $color-secondary, $color-primary);
  border-radius: 50%;
  animation: pulseSoft 2s $ease-custom infinite;
  filter: blur(1px);
  
  // Add glow effect
  &::after {
    content: '';
    position: absolute;
    top: -10%;
    left: -10%;
    right: -10%;
    bottom: -10%;
    background: inherit;
    border-radius: inherit;
    filter: blur(20px);
    opacity: 0.4;
    animation: pulseSoft 2s $ease-custom infinite;
  }
}

// Progress Variant
.loading-progress {
  position: relative;
}

.progress-circle {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
}

.progress-svg {
  width: 100%;
  height: 100%;
  transform: rotate(-90deg);
}

.progress-bar {
  transition: stroke-dasharray 0.3s $ease-custom;
}

.progress-text {
  position: absolute;
  font-weight: $font-weight-semibold;
  color: $color-primary;
}

// Spinner Variant
.loading-spinner {
  position: relative;
}

.spinner {
  border-radius: 50%;
  border: 3px solid $color-secondary;
  border-top-color: $color-primary;
  animation: spin 1s linear infinite;
}

// Loading Text
.loading-text {
  color: $color-text-subtle;
  text-align: center;
  animation: fadeIn 0.6s $ease-custom;
  
  &.with-margin {
    margin-top: $spacing-sm;
  }
}

// Size Variants
.size-small {
  gap: $spacing-sm;
  
  .dot {
    width: 8px;
    height: 8px;
  }
  
  .loading-blob {
    width: 40px;
    height: 40px;
  }
  
  .progress-circle {
    width: 60px;
    height: 60px;
  }
  
  .spinner {
    width: 24px;
    height: 24px;
  }
  
  .progress-text {
    font-size: $font-size-caption;
  }
  
  .loading-text {
    font-size: $font-size-caption;
  }
}

.size-medium {
  .dot {
    width: 12px;
    height: 12px;
  }
  
  .loading-blob {
    width: 80px;
    height: 80px;
  }
  
  .progress-circle {
    width: 100px;
    height: 100px;
  }
  
  .spinner {
    width: 40px;
    height: 40px;
  }
  
  .progress-text {
    font-size: $font-size-body;
  }
  
  .loading-text {
    font-size: $font-size-body;
  }
}

.size-large {
  gap: $spacing-lg;
  
  .dot {
    width: 16px;
    height: 16px;
  }
  
  .loading-blob {
    width: 120px;
    height: 120px;
  }
  
  .progress-circle {
    width: 140px;
    height: 140px;
  }
  
  .spinner {
    width: 60px;
    height: 60px;
    border-width: 4px;
  }
  
  .progress-text {
    font-size: 20px;
  }
  
  .loading-text {
    font-size: 18px;
  }
}

// Animations
@keyframes loadingDot {
  0%, 80%, 100% {
    transform: scale(0);
    opacity: 0;
  }
  40% {
    transform: scale(1);
    opacity: 1;
  }
}

@keyframes pulseSoft {
  0%, 100% {
    opacity: 1;
    transform: scale(1);
  }
  50% {
    opacity: 0.8;
    transform: scale(1.05);
  }
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

// Ensure smooth 60fps animations
* {
  will-change: transform, opacity;
  transform: translateZ(0);
  backface-visibility: hidden;
}
</style> 