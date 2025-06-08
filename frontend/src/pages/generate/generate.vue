<template>
  <view class="generate-container page-wrapper animate-fade-in">
    <!-- Header -->
    <view class="generate-header">
      <text class="generate-title">Create Your Magic</text>
      <CreditDisplay size="small" :showLabel="true" />
    </view>
    
    <!-- Steps -->
    <view class="steps-container">
      <!-- Step 1: Upload Photo -->
      <view class="step-card card" :class="{ 'active': currentStep === 1 }">
        <view class="step-header">
          <text class="step-number">1</text>
          <text class="step-title">Upload Your Photo</text>
        </view>
        
        <view class="step-content">
          <view v-if="!generationStore.userPhotoUrl" class="upload-area" @click="choosePhoto">
            <text class="upload-icon">📸</text>
            <text class="upload-text">Tap to select a photo</text>
            <text class="upload-hint">Best results with clear face photos</text>
          </view>
          
          <view v-else class="photo-preview">
            <image 
              :src="generationStore.userPhotoUrl" 
              mode="aspectFill"
              class="preview-image"
            />
            <button class="btn-ghost small change-photo" @click="choosePhoto">
              Change Photo
            </button>
          </view>
          
          <!-- Upload Progress -->
          <view v-if="generationStore.uploadProgress > 0 && generationStore.uploadProgress < 100" class="upload-progress">
            <LoadingAnimation 
              variant="progress" 
              :progress="generationStore.uploadProgress"
              size="small"
              text="Uploading..."
            />
          </view>
        </view>
      </view>
      
      <!-- Step 2: Confirm Template -->
      <view class="step-card card" :class="{ 'active': currentStep === 2, 'disabled': !generationStore.userPhotoUrl }">
        <view class="step-header">
          <text class="step-number">2</text>
          <text class="step-title">Confirm Style</text>
        </view>
        
        <view v-if="selectedTemplate" class="step-content">
          <view class="template-preview">
            <TemplateCard 
              :template="selectedTemplate"
              @click="goToGallery"
            />
          </view>
          
          <view class="template-confirm">
            <text class="confirm-text">
              This will use <text class="highlight">{{ selectedTemplate.credit_cost }} credits</text>
            </text>
            <button 
              class="btn-primary gradient"
              :disabled="!canGenerate"
              @click="startGeneration"
            >
              Generate Now
            </button>
          </view>
        </view>
      </view>
    </view>
    
    <!-- Generation Progress -->
    <view v-if="generationStore.isGenerating" class="generation-overlay">
      <view class="generation-modal card">
        <LoadingAnimation 
          variant="blob" 
          size="large"
          text="Creating your masterpiece..."
        />
        
        <view v-if="generationStore.currentTask" class="generation-info">
          <view class="progress-bar">
            <view 
              class="progress-fill"
              :style="{ width: `${generationStore.currentTask.progress}%` }"
            ></view>
          </view>
          
          <text class="time-estimate">
            Estimated time: {{ formatTime(generationStore.estimatedTimeRemaining) }}
          </text>
        </view>
        
        <button class="btn-ghost small" @click="cancelGeneration">
          Cancel
        </button>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useTemplateStore } from '@/stores/templates'
import { useGenerationStore } from '@/stores/generation'
import { useUserStore } from '@/stores/user'
import { customTransitions } from '@/utils/navigation'
import TemplateCard from '@/components/TemplateCard/TemplateCard.vue'
import CreditDisplay from '@/components/CreditDisplay/CreditDisplay.vue'
import LoadingAnimation from '@/components/LoadingAnimation/LoadingAnimation.vue'

// Stores
const templateStore = useTemplateStore()
const generationStore = useGenerationStore()
const userStore = useUserStore()

// State
const currentStep = ref(1)
const selectedTemplate = ref<any>(null)

// Get template ID from route params
const pages = getCurrentPages()
const currentPage = pages[pages.length - 1]
const templateId = currentPage.options?.templateId

// Computed
const canGenerate = computed(() => {
  return generationStore.userPhotoUrl && 
         selectedTemplate.value && 
         userStore.credits >= selectedTemplate.value.credit_cost &&
         !generationStore.isGenerating
})

// Methods
const choosePhoto = () => {
  uni.chooseImage({
    count: 1,
    sizeType: ['compressed'],
    sourceType: ['album', 'camera'],
    success: async (res) => {
      const tempFilePath = res.tempFilePaths[0]
      
      try {
        await generationStore.uploadPhoto(tempFilePath)
        currentStep.value = 2
      } catch (error: any) {
        uni.showToast({
          title: error.message || 'Upload failed',
          icon: 'none'
        })
      }
    }
  })
}

const goToGallery = () => {
  customTransitions.toGallery()
}

const startGeneration = async () => {
  if (!canGenerate.value) return
  
  try {
    const task = await generationStore.startGeneration(selectedTemplate.value.id)
    
    // Poll for completion
    const checkInterval = setInterval(() => {
      if (generationStore.currentTask?.status === 'completed') {
        clearInterval(checkInterval)
        // Navigate to result page
        uni.navigateTo({
          url: `/pages/result/result?taskId=${task.id}`
        })
      } else if (generationStore.currentTask?.status === 'failed') {
        clearInterval(checkInterval)
      }
    }, 1000)
  } catch (error: any) {
    uni.showToast({
      title: error.message || 'Generation failed',
      icon: 'none'
    })
  }
}

const cancelGeneration = () => {
  uni.showModal({
    title: 'Cancel Generation?',
    content: 'Your credits will be refunded if you cancel now.',
    confirmText: 'Yes, Cancel',
    confirmColor: '#E89B93',
    success: (res) => {
      if (res.confirm) {
        generationStore.cancelGeneration()
      }
    }
  })
}

const formatTime = (ms: number) => {
  const seconds = Math.round(ms / 1000)
  if (seconds < 60) return `${seconds}s`
  const minutes = Math.floor(seconds / 60)
  return `${minutes}m ${seconds % 60}s`
}

// Load template on mount
onMounted(async () => {
  if (templateId) {
    try {
      selectedTemplate.value = await templateStore.selectTemplate(parseInt(templateId))
    } catch (error) {
      uni.showToast({
        title: 'Failed to load template',
        icon: 'none'
      })
      setTimeout(() => {
        uni.navigateBack()
      }, 1500)
    }
  }
})
</script>

<style lang="scss" scoped>
@import '@/styles/variables.scss';

.generate-container {
  background-color: $color-bg;
  min-height: 100vh;
  padding-bottom: $spacing-xl;
}

.generate-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: $spacing-lg $page-padding;
}

.generate-title {
  font-size: $font-size-h1;
  font-weight: $font-weight-semibold;
  color: $color-text;
}

// Steps
.steps-container {
  padding: 0 $page-padding;
}

.step-card {
  margin-bottom: $spacing-lg;
  transition: all $duration-fast $ease-custom;
  
  &.active {
    border-color: $color-primary;
    box-shadow: 0 4px 16px rgba($color-primary, 0.15);
  }
  
  &.disabled {
    opacity: 0.6;
    pointer-events: none;
  }
}

.step-header {
  display: flex;
  align-items: center;
  gap: $spacing-md;
  margin-bottom: $spacing-lg;
}

.step-number {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: $color-secondary;
  color: white;
  border-radius: $radius-full;
  font-weight: $font-weight-semibold;
}

.step-title {
  font-size: 18px;
  font-weight: $font-weight-semibold;
  color: $color-text;
}

.step-content {
  position: relative;
}

// Upload Area
.upload-area {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: $spacing-xl * 2;
  background: $color-bg;
  border: 2px dashed $color-secondary;
  border-radius: $radius-md;
  cursor: pointer;
  transition: all $duration-fast $ease-custom;
  
  &:hover {
    border-color: $color-primary;
    background: rgba($color-secondary, 0.05);
  }
  
  &:active {
    transform: scale(0.98);
  }
}

.upload-icon {
  font-size: 48px;
  margin-bottom: $spacing-md;
}

.upload-text {
  font-size: $font-size-body;
  color: $color-text;
  margin-bottom: $spacing-xs;
}

.upload-hint {
  font-size: $font-size-caption;
  color: $color-text-subtle;
}

// Photo Preview
.photo-preview {
  position: relative;
}

.preview-image {
  width: 100%;
  height: 300px;
  border-radius: $radius-md;
}

.change-photo {
  position: absolute;
  bottom: $spacing-md;
  right: $spacing-md;
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(10px);
}

.upload-progress {
  margin-top: $spacing-lg;
}

// Template Confirm
.template-confirm {
  margin-top: $spacing-lg;
  text-align: center;
}

.confirm-text {
  display: block;
  font-size: $font-size-body;
  color: $color-text;
  margin-bottom: $spacing-md;
  
  .highlight {
    color: $color-primary;
    font-weight: $font-weight-semibold;
  }
}

// Generation Overlay
.generation-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.8);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: $z-modal;
  animation: fadeIn $duration-fast $ease-custom;
}

.generation-modal {
  width: 90%;
  max-width: 400px;
  padding: $spacing-xl;
  text-align: center;
}

.generation-info {
  margin-top: $spacing-lg;
  width: 100%;
}

.progress-bar {
  width: 100%;
  height: 4px;
  background: $color-bg;
  border-radius: 2px;
  overflow: hidden;
  margin-bottom: $spacing-md;
}

.progress-fill {
  height: 100%;
  background: $color-gradient;
  transition: width $duration-fast $ease-custom;
}

.time-estimate {
  font-size: $font-size-caption;
  color: $color-text-subtle;
}
</style> 