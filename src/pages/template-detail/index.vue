<route lang="json5">
{
  style: {
    navigationBarTitleText: 'Template Detail',
  },
}
</route>

<template>
  <view class="template-detail-container page-wrapper animate-fade-in">
    <!-- Loading State -->
    <LoadingAnimation v-if="loading" variant="blob" size="large" text="Loading template..." />

    <!-- Error State -->
    <view v-else-if="error" class="error-container">
      <text class="error-text">Failed to load template details</text>
      <button class="retry-button" @click="loadTemplate">Retry</button>
    </view>

    <!-- Template Content -->
    <template v-else-if="template">
      <view class="template-content">
        <view class="image-container">
          <!-- Back Button Overlay -->
          <view class="back-button-overlay" @click="goBack">
            <view class="back-button">
              <text class="back-icon">←</text>
            </view>
          </view>
          <image class="preview-image" :src="template.preview_image_url" mode="aspectFill" />
        </view>

        <view class="info-section">
          <view class="template-info">
            <text class="title">{{ template.name }}</text>
            <text class="description">{{ template.description }}</text>
          </view>

          <view class="action-section">
            <view class="cost-badge-bottom">
              <text class="cost-icon">🎞️</text>
              <text class="cost-text">{{ template.credit_cost }} 积分</text>
            </view>

            <button class="generate-button" @click="handleGenerate" :disabled="generating">
              <text class="button-text">{{ generating ? '生成中...' : '立即生成' }}</text>
            </button>
          </view>
        </view>
      </view>
    </template>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import LoadingAnimation from '@/components/LoadingAnimation/LoadingAnimation.vue'
import type { Template } from '@/types/api'

// State
const template = ref<Template | null>(null)
const loading = ref(true)
const generating = ref(false)
const error = ref(false)
const templateId = ref<string>('')

// Load template data
const loadTemplate = async () => {
  if (!templateId.value) return

  loading.value = true
  error.value = false

  try {
    // For now, use mock data since we don't have backend ready
    // In real implementation, this would call the API
    // const response = await templatesApi.getById(templateId.value)
    // template.value = response.data

    // Mock template data
    const mockTemplate: Template = {
      id: parseInt(templateId.value),
      name: 'Dusty Rose Dream',
      description: 'Soft, romantic vibes with dreamy pink tones and ethereal lighting effects',
      preview_image_url: `https://picsum.photos/400/600?random=${templateId.value}`,
      credit_cost: 15,
      is_active: true,
      created_at: '2024-01-01T00:00:00Z',
    }

    template.value = mockTemplate
  } catch (err) {
    console.error('Failed to load template:', err)
    error.value = true
  } finally {
    loading.value = false
  }
}

// Handle generate button click
const handleGenerate = () => {
  if (!template.value) return

  // Navigate to generate page with template
  uni.navigateTo({
    url: `/pages/generate/index?templateId=${template.value.id}`,
  })
}

// Handle back button click
const goBack = () => {
  uni.navigateBack({
    fail() {
      uni.switchTab({
        url: '/pages/gallery/index',
      })
    },
  })
}

// Page lifecycle
onLoad((options) => {
  templateId.value = options.id || ''
  loadTemplate()
})

onMounted(() => {
  if (templateId.value) {
    loadTemplate()
  }
})
</script>

<style lang="scss" scoped>
@import '@/style/variables.scss';

.detail-container {
  background-color: $color-bg;
  min-height: 100vh;
}

.error-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: $spacing-xl;
  min-height: 300px;
}

.error-text {
  color: $color-text-subtle;
  font-size: $font-size-body;
  margin-bottom: $spacing-md;
}

.retry-button {
  background-color: $color-primary;
  color: white;
  border: none;
  padding: $spacing-sm $spacing-lg;
  border-radius: $radius-md;
  font-size: $font-size-body;
}

.template-content {
  background-color: $color-bg;
  min-height: 100vh;
}

.image-container {
  position: relative;
  height: 70vh;
  overflow: hidden;
}

.back-button-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  z-index: 10;
  padding: calc(env(safe-area-inset-top, 0px) + 16px) 20px 16px;
}

.back-button {
  width: 40px;
  height: 40px;
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(10px);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all $duration-fast $ease-custom;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);

  &:active {
    transform: scale(0.95);
    background: rgba(0, 0, 0, 0.7);
  }
}

.back-icon {
  color: white;
  font-size: 18px;
  font-weight: bold;
}

.preview-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.image-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(transparent 40%, rgba(0, 0, 0, 0.7));
  display: flex;
  align-items: flex-end;
  padding: $spacing-xl;
}

.cost-badge {
  display: flex;
  align-items: center;
  gap: $spacing-xs;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  padding: $spacing-sm $spacing-md;
  border-radius: $radius-full;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
}

.cost-icon {
  font-size: 18px;
}

.cost-text {
  font-size: $font-size-body;
  font-weight: $font-weight-medium;
  color: $color-primary;
}

.info-section {
  background-color: $color-surface;
  border-radius: $radius-xl $radius-xl 0 0;
  margin-top: -$spacing-xl;
  padding: $spacing-xl;
  position: relative;
  z-index: 2;
}

.template-info {
  margin-bottom: $spacing-xl;
}

.title {
  font-size: 24px;
  font-weight: $font-weight-semibold;
  color: $color-text;
  display: block;
  margin-bottom: $spacing-md;
  line-height: 1.3;
}

.description {
  font-size: $font-size-body;
  color: $color-text-subtle;
  line-height: 1.6;
  display: block;
}

.action-section {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.cost-badge-bottom {
  display: flex;
  align-items: center;
  gap: $spacing-xs;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  padding: $spacing-sm $spacing-md;
  border-radius: $radius-full;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
}

.generate-button {
  background: linear-gradient(135deg, $color-primary, lighten($color-primary, 10%));
  color: white;
  border: none;
  padding: $spacing-md $spacing-xl;
  border-radius: $radius-xl;
  font-size: $font-size-body;
  font-weight: $font-weight-semibold;
  cursor: pointer;
  transition: all $duration-fast $ease-custom;
  box-shadow: 0 8px 24px rgba($color-primary, 0.3);
  flex: 1;
  margin-left: $spacing-md;

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 12px 32px rgba($color-primary, 0.4);
  }

  &:active {
    transform: translateY(0);
  }

  &:disabled {
    opacity: 0.6;
    cursor: not-allowed;
    transform: none !important;
  }
}

.button-text {
  font-size: $font-size-body;
  font-weight: $font-weight-semibold;
  color: white;
}
</style>
