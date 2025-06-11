<template>
  <view class="generate-container page-wrapper animate-fade-in">
    <view class="header">
      <view class="title">AI 写真馆</view>
      <view class="subtitle">创建你的专属形象</view>
    </view>

    <view class="content-card">
      <!-- State 1: Upload Photo -->
      <view v-if="!localPhotoPath" class="upload-state">
        <view class="upload-placeholder" @click="handleChoosePhoto">
          <image class="upload-icon" src="/static/icons/camera-add.svg" />
          <view class="upload-text">上传一张清晰的正面照</view>
          <view class="upload-subtext">让AI为你创造惊喜</view>
        </view>
      </view>

      <!-- State 2: Preview & Generate -->
      <view
        v-else-if="
          !generationStore.currentTask || generationStore.currentTask.status !== 'completed'
        "
        class="preview-state"
      >
        <image class="preview-image" :src="localPhotoPath" mode="aspectFill" />

        <view class="template-info">
          <image
            class="template-preview-thumb"
            :src="selectedTemplate?.preview_image_url"
            mode="aspectFill"
          />
          <view class="template-details">
            <view class="template-name">{{ selectedTemplate?.name || '加载中...' }}</view>
            <view class="template-cost">
              <image class="credit-icon" src="/static/icons/credit-coin.svg" />
              <text>{{ selectedTemplate?.credit_cost || '-' }} 积分</text>
            </view>
          </view>
        </view>

        <button
          class="generate-button"
          :class="{ disabled: !canGenerate }"
          :disabled="!canGenerate"
          @click="handleGenerate"
        >
          <text>立即生成</text>
        </button>
        <view class="re-upload-text" @click="handleChoosePhoto">重新上传照片</view>
      </view>

      <!-- State 3: Result Display -->
      <view v-else-if="generationStore.currentTask?.status === 'completed'" class="result-state">
        <view class="result-header">
          <text class="result-title">🎉 生成完成！</text>
          <text class="result-subtitle">你的专属AI作品</text>
        </view>

        <!-- Generated Image -->
        <view class="result-image-container">
          <image
            v-if="generationStore.currentTask.resultUrl"
            class="result-image"
            :src="generationStore.currentTask.resultUrl"
            mode="aspectFit"
            @error="handleImageError"
          />
          <view v-else class="image-placeholder">
            <text class="placeholder-text">图片加载中...</text>
          </view>
        </view>

        <!-- Action Buttons -->
        <view class="result-actions">
          <button class="action-button secondary" @click="handleSaveImage">
            <text>保存到相册</text>
          </button>
          <button class="action-button primary" @click="handleGenerateAgain">
            <text>再生成一张</text>
          </button>
        </view>

        <view class="result-info">
          <text class="template-used">
            使用模板：{{ generationStore.currentTask.templateName }}
          </text>
          <text class="generation-time">生成时间：{{ formatGenerationTime() }}</text>
          <!-- Debug info for development -->
          <text v-if="isDev" class="debug-info">
            图片URL: {{ generationStore.currentTask.resultUrl || '未设置' }}
          </text>
        </view>
      </view>
    </view>

    <!-- Generating Overlay -->
    <view v-if="generationStore.isGenerating" class="generating-overlay">
      <view class="generating-modal">
        <LoadingAnimation variant="spinner" size="large" />
        <view class="generating-text">AI正在为你绘制作品</view>
        <view class="generating-subtext">
          {{
            generationStore.currentTask?.status === 'pending'
              ? '正在排队中...'
              : generationStore.currentTask?.status === 'processing'
                ? '正在生成中...'
                : '预计需要 30 秒，请稍候'
          }}
        </view>
        <view class="progress-bar-container">
          <view
            class="progress-bar"
            :style="{
              width: `${generationStore.currentTask?.progress || generationStore.uploadProgress}%`,
            }"
          ></view>
        </view>
        <view class="progress-text">
          {{ Math.round(generationStore.currentTask?.progress || generationStore.uploadProgress) }}%
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useTemplateStore } from '@/store/templates'
import { useGenerationStore } from '@/store/generation'
import { useUserStore } from '@/store/user'
import LoadingAnimation from '@/components/LoadingAnimation/LoadingAnimation.vue'
import { navigateBack } from '@/utils/navigation'

// Stores
const templateStore = useTemplateStore()
const generationStore = useGenerationStore()
const userStore = useUserStore()

// State
const localPhotoPath = ref<string>('')
const selectedTemplate = ref<any>(null) // Using any for simplicity with template structure

// Get template ID from route
const pages = getCurrentPages()
const currentPage = pages[pages.length - 1]
const routeTemplateId = (currentPage as any)?.options?.templateId

// Safe template ID parsing with comprehensive validation
let templateId: number | null = null
if (routeTemplateId && routeTemplateId !== 'undefined' && routeTemplateId !== 'null') {
  const parsed = parseInt(routeTemplateId.toString(), 10)
  if (!isNaN(parsed) && parsed > 0) {
    templateId = parsed
  }
}

// Debug logging for templateId issues
console.log('Template ID processing:', {
  routeTemplateId,
  type: typeof routeTemplateId,
  parsed: templateId,
  options: (currentPage as any)?.options,
})

// Computed
const canGenerate = computed(() => {
  const hasPhoto = !!localPhotoPath.value
  const hasTemplate = !!selectedTemplate.value
  const notGenerating = !generationStore.isGenerating
  const hasCredits = userStore.credits >= (selectedTemplate.value?.credit_cost || 999)

  console.log('canGenerate check:', {
    hasPhoto,
    hasTemplate,
    notGenerating,
    hasCredits,
    userCredits: userStore.credits,
    userInfo: userStore.userInfo,
    isAuthenticated: userStore.isAuthenticated,
    templateCost: selectedTemplate.value?.credit_cost,
    photoPath: localPhotoPath.value,
    template: selectedTemplate.value?.name,
  })

  return hasPhoto && hasTemplate && notGenerating && hasCredits
})

const isDev = computed(() => {
  return import.meta.env.DEV || import.meta.env.VITE_USER_NODE_ENV === 'development'
})

// Methods
const handleChoosePhoto = () => {
  // Check platform and use appropriate method
  // #ifdef H5
  // Use native HTML file input for H5 mode (web browsers)
  const input = document.createElement('input')
  input.type = 'file'
  input.accept = 'image/*'
  input.style.display = 'none'

  input.onchange = (event: any) => {
    const file = event.target.files?.[0]
    if (file) {
      // Create object URL for preview
      const tempFilePath = URL.createObjectURL(file)
      localPhotoPath.value = tempFilePath

      // Also update store for any components that might use it for preview
      generationStore.userPhotoUrl = tempFilePath

      console.log('Photo selected (H5):', { fileName: file.name, size: file.size, type: file.type })
    }

    // Clean up the input element
    document.body.removeChild(input)
  }

  input.onerror = () => {
    uni.showToast({ title: '选择照片失败', icon: 'none' })
    document.body.removeChild(input)
  }

  // Add to DOM and trigger click
  document.body.appendChild(input)
  input.click()
  // #endif

  // #ifndef H5
  // Use uni.chooseImage for mini-program and mobile apps
  uni.chooseImage({
    count: 1,
    sizeType: ['compressed'],
    sourceType: ['album', 'camera'],
    success: (res) => {
      const tempFilePath = (res.tempFilePaths as string[])[0]
      localPhotoPath.value = tempFilePath
      // Also update store for any components that might use it for preview
      generationStore.userPhotoUrl = tempFilePath
      console.log('Photo selected (Mini-program/App):', { path: tempFilePath })
    },
    fail: () => {
      uni.showToast({ title: '选择照片失败', icon: 'none' })
    },
  })
  // #endif
}

const handleGenerate = async () => {
  if (!canGenerate.value) {
    if (userStore.credits < (selectedTemplate.value?.credit_cost || 999)) {
      uni.showToast({ title: '积分不足', icon: 'none' })
    }
    return
  }

  try {
    await generationStore.startGeneration(selectedTemplate.value.id, localPhotoPath.value)

    // Watch for generation completion
    const unwatch = generationStore.$subscribe((mutation, state) => {
      const task = state.currentTask
      if (task?.status === 'completed') {
        // Show success message
        uni.showToast({
          title: '🎉 生成完成！',
          icon: 'success',
          duration: 2000,
        })

        unwatch() // Stop watching after completion
      } else if (task?.status === 'failed') {
        uni.showToast({ title: task.errorMessage || '生成失败', icon: 'none' })
        unwatch() // Stop watching
      }
    })
  } catch (error: any) {
    uni.showToast({ title: error.message || '生成请求失败', icon: 'none' })
  }
}

const handleImageError = () => {
  const resultUrl = generationStore.currentTask?.resultUrl
  console.error('Generated image failed to load:', {
    resultUrl,
    taskStatus: generationStore.currentTask?.status,
    taskId: generationStore.currentTask?.id,
  })

  // Show more helpful error message
  if (!resultUrl) {
    uni.showToast({ title: '图片URL缺失', icon: 'none' })
  } else {
    uni.showToast({ title: '图片加载失败，请重试', icon: 'none' })
  }
}

const handleSaveImage = () => {
  const resultUrl = generationStore.currentTask?.resultUrl
  if (!resultUrl) {
    uni.showToast({ title: '没有可保存的图片', icon: 'none' })
    return
  }

  // #ifdef H5
  // For H5, try to trigger download
  const link = document.createElement('a')
  link.href = resultUrl
  link.download = `ai-generated-${Date.now()}.png`
  link.click()
  uni.showToast({ title: '开始下载', icon: 'success' })
  // #endif

  // #ifndef H5
  // For mobile platforms, save to album
  uni.saveImageToPhotosAlbum({
    filePath: resultUrl,
    success: () => {
      uni.showToast({ title: '保存成功', icon: 'success' })
    },
    fail: () => {
      uni.showToast({ title: '保存失败', icon: 'none' })
    },
  })
  // #endif
}

const handleGenerateAgain = () => {
  // Reset generation state
  generationStore.currentTask = null
  generationStore.isGenerating = false

  // Keep the same photo and template, just regenerate
  if (canGenerate.value) {
    handleGenerate()
  } else {
    uni.showToast({ title: '请检查积分余额', icon: 'none' })
  }
}

const formatGenerationTime = () => {
  const task = generationStore.currentTask
  if (!task?.completedAt) return ''

  const completedTime = new Date(task.completedAt)
  const now = new Date()
  const diffMinutes = Math.floor((now.getTime() - completedTime.getTime()) / (1000 * 60))

  if (diffMinutes < 1) return '刚刚'
  if (diffMinutes < 60) return `${diffMinutes}分钟前`

  const diffHours = Math.floor(diffMinutes / 60)
  if (diffHours < 24) return `${diffHours}小时前`

  return completedTime.toLocaleDateString('zh-CN')
}

// Lifecycle
onMounted(async () => {
  // Clear any previous generation state
  generationStore.clearUserPhoto()
  localPhotoPath.value = ''

  // Load user data to ensure credits are available
  try {
    await userStore.fetchUserInfo()
    console.log('User info loaded:', { credits: userStore.credits, userInfo: userStore.userInfo })
  } catch (error) {
    console.error('Failed to load user info:', error)
  }

  if (templateId !== null) {
    try {
      const template = await templateStore.selectTemplate(templateId)
      if (template) {
        selectedTemplate.value = template
        console.log('Template loaded:', { name: template.name, cost: template.credit_cost })
      } else {
        throw new Error('Template not found.')
      }
    } catch (error: any) {
      console.error('Failed to load template:', error)
      uni.showToast({ title: error.message || '加载模板失败', icon: 'none' })
      setTimeout(() => navigateBack(), 1500)
    }
  } else {
    console.error('Invalid or missing templateId:', { routeTemplateId, templateId })
    uni.showToast({ title: '无效的模板ID', icon: 'none' })
    setTimeout(() => navigateBack(), 1500)
  }
})
</script>

<style lang="scss" scoped>
@import '@/style/variables.scss';

.generate-container {
  background: linear-gradient(180deg, #fdf4f3 0%, #f5f7fa 100%);
  min-height: 100vh;
  padding-top: calc(env(safe-area-inset-top, 0px) + 20px);
}

.header {
  text-align: center;
  padding: 40rpx 0;
  .title {
    font-size: 48rpx;
    font-weight: 600;
    color: $color-text;
  }
  .subtitle {
    font-size: 28rpx;
    color: $color-text-subtle;
    margin-top: 10rpx;
  }
}

.content-card {
  background-color: #ffffff;
  border-radius: 24rpx;
  padding: 40rpx;
  box-shadow: 0 8rpx 32rpx rgba(0, 0, 0, 0.05);
  min-height: 700rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

// Upload State
.upload-state {
  width: 100%;
}
.upload-placeholder {
  border: 2rpx dashed $color-primary;
  border-radius: 24rpx;
  padding: 100rpx 40rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
  cursor: pointer;
  transition: background-color 0.2s;
  &:hover {
    background-color: #fffbfb;
  }
}
.upload-icon {
  width: 100rpx;
  height: 100rpx;
  opacity: 0.8;
}
.upload-text {
  font-size: 32rpx;
  font-weight: 500;
  color: $color-text;
  margin-top: 30rpx;
}
.upload-subtext {
  font-size: 26rpx;
  color: $color-text-subtle;
  margin-top: 10rpx;
}

// Preview State
.preview-state {
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
}
.preview-image {
  width: 100%;
  height: 500rpx;
  border-radius: 16rpx;
  object-fit: cover;
  margin-bottom: 30rpx;
}
.template-info {
  display: flex;
  align-items: center;
  width: 100%;
  background-color: #f9fafc;
  padding: 20rpx;
  border-radius: 16rpx;
  margin-bottom: 40rpx;
}
.template-preview-thumb {
  width: 100rpx;
  height: 100rpx;
  border-radius: 12rpx;
  margin-right: 20rpx;
}
.template-details {
  flex: 1;
}
.template-name {
  font-size: 30rpx;
  font-weight: 500;
  color: $color-text;
}
.template-cost {
  display: flex;
  align-items: center;
  font-size: 26rpx;
  color: $color-text-subtle;
  margin-top: 8rpx;
}
.credit-icon {
  width: 32rpx;
  height: 32rpx;
  margin-right: 10rpx;
}

.generate-button {
  width: 100%;
  padding: 24rpx;
  font-size: 32rpx;
  font-weight: 500;
  color: #ffffff;
  background: $color-primary;
  border: none;
  border-radius: 16rpx;
  cursor: pointer;
  transition: opacity 0.2s;
  &.disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
}
.re-upload-text {
  font-size: 26rpx;
  color: $color-text-subtle;
  margin-top: 30rpx;
  cursor: pointer;
  text-decoration: underline;
}

// Result State
.result-state {
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.result-header {
  text-align: center;
  margin-bottom: 40rpx;
}

.result-title {
  font-size: 36rpx;
  font-weight: 600;
  color: $color-text;
  display: block;
  margin-bottom: 10rpx;
}

.result-subtitle {
  font-size: 26rpx;
  color: $color-text-subtle;
  display: block;
}

.result-image-container {
  width: 100%;
  height: 500rpx;
  border-radius: 16rpx;
  overflow: hidden;
  margin-bottom: 30rpx;
  background-color: #f5f5f5;
  display: flex;
  align-items: center;
  justify-content: center;
}

.result-image {
  width: 100%;
  height: 100%;
}

.image-placeholder {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
}

.placeholder-text {
  color: $color-text-subtle;
  font-size: 28rpx;
}

.result-actions {
  display: flex;
  gap: 20rpx;
  width: 100%;
  margin-bottom: 30rpx;
}

.action-button {
  flex: 1;
  padding: 24rpx;
  font-size: 28rpx;
  font-weight: 500;
  border: none;
  border-radius: 16rpx;
  cursor: pointer;
  transition: opacity 0.2s;

  &.primary {
    color: #ffffff;
    background: $color-primary;
  }

  &.secondary {
    color: $color-primary;
    background: rgba($color-primary, 0.1);
    border: 1rpx solid $color-primary;
  }

  &:active {
    opacity: 0.8;
  }
}

.result-info {
  width: 100%;
  text-align: center;
}

.template-used,
.generation-time,
.debug-info {
  display: block;
  font-size: 24rpx;
  color: $color-text-subtle;
  margin-bottom: 8rpx;
}

.debug-info {
  font-family: monospace;
  font-size: 20rpx;
  color: #999;
  word-break: break-all;
}

// Generating Overlay
.generating-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}
.generating-modal {
  background-color: #ffffff;
  border-radius: 24rpx;
  padding: 60rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 80%;
  max-width: 600rpx;
}
.generating-text {
  font-size: 36rpx;
  font-weight: 500;
  margin-top: 40rpx;
}
.generating-subtext {
  font-size: 28rpx;
  color: $color-text-subtle;
  margin-top: 10rpx;
  margin-bottom: 40rpx;
}
.progress-bar-container {
  width: 100%;
  height: 8rpx;
  background-color: #f0f0f0;
  border-radius: 4rpx;
  overflow: hidden;
}
.progress-bar {
  height: 100%;
  background-color: $color-primary;
  transition: width 0.3s ease;
}
.progress-text {
  font-size: 24rpx;
  color: $color-text-subtle;
  margin-top: 20rpx;
  text-align: center;
}
</style>
