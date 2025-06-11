<template>
  <div class="image-upload">
    <!-- Upload Area -->
    <div
      class="upload-area"
      :class="{
        'upload-area--dragging': isDragging,
        'upload-area--has-image': selectedImage,
      }"
      @click="selectImage"
    >
      <!-- No Image State -->
      <div v-if="!selectedImage" class="upload-placeholder">
        <div class="upload-icon">
          <svg width="48" height="48" viewBox="0 0 24 24" fill="none">
            <path
              d="M12 2L13.09 8.26L22 9L13.09 9.74L12 16L10.91 9.74L2 9L10.91 8.26L12 2Z"
              fill="currentColor"
            />
          </svg>
        </div>
        <h3 class="upload-title">选择照片开始创作</h3>
        <p class="upload-subtitle">支持 JPG、PNG 格式，最大 10MB</p>
        <div class="upload-actions">
          <button class="btn btn--primary" @click.stop="selectFromGallery">
            <i class="icon-gallery"></i>
            从相册选择
          </button>
          <button class="btn btn--secondary" @click.stop="takePhoto">
            <i class="icon-camera"></i>
            拍照
          </button>
        </div>
      </div>

      <!-- Image Preview -->
      <div v-else class="image-preview">
        <img :src="selectedImage.url" :alt="selectedImage.name" class="preview-image" />
        <div class="image-overlay">
          <div class="image-info">
            <span class="image-name">{{ selectedImage.name }}</span>
            <span class="image-size">{{ formatFileSize(selectedImage.size) }}</span>
          </div>
          <div class="image-actions">
            <button class="btn-icon" @click.stop="cropImage" title="裁剪">
              <i class="icon-crop"></i>
            </button>
            <button class="btn-icon" @click.stop="removeImage" title="删除">
              <i class="icon-delete"></i>
            </button>
          </div>
        </div>
      </div>

      <!-- Upload Progress -->
      <div v-if="uploadProgress > 0 && uploadProgress < 100" class="upload-progress">
        <div class="progress-bar">
          <div class="progress-fill" :style="{ width: uploadProgress + '%' }"></div>
        </div>
        <span class="progress-text">上传中... {{ uploadProgress }}%</span>
      </div>
    </div>

    <!-- Error Message -->
    <div v-if="errorMessage" class="error-message">
      <i class="icon-warning"></i>
      {{ errorMessage }}
    </div>

    <!-- Image Crop Modal -->
    <div v-if="showCropModal" class="crop-modal" @click="closeCropModal">
      <div class="crop-container" @click.stop>
        <div class="crop-header">
          <h3>裁剪图片</h3>
          <button class="btn-close" @click="closeCropModal">×</button>
        </div>
        <div class="crop-area">
          <canvas ref="cropCanvas" class="crop-canvas"></canvas>
        </div>
        <div class="crop-actions">
          <button class="btn btn--secondary" @click="closeCropModal">取消</button>
          <button class="btn btn--primary" @click="applyCrop">确定</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, nextTick, watch } from 'vue'

interface ImageFile {
  file: File | string | Blob
  url: string
  name: string
  size: number
}

// Props
interface Props {
  maxSize?: number // MB
  acceptedFormats?: string[]
  multiple?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  maxSize: 10,
  acceptedFormats: () => ['image/jpeg', 'image/png', 'image/webp'],
  multiple: false,
})

// Emits
const emit = defineEmits<{
  upload: [file: ImageFile]
  error: [message: string]
  progress: [progress: number]
}>()

// Platform detection
const isWeChat = ref(false)
const isIOS = ref(false)
const isMobile = ref(false)
const uploadProgress = ref(0)

// State
const selectedImage = ref<ImageFile | null>(null)
const isDragging = ref(false)
const errorMessage = ref('')
const showCropModal = ref(false)
const cropCanvas = ref<HTMLCanvasElement>()

// Methods
const selectImage = () => {
  if (selectedImage.value) return
  selectFromGallery()
}

const selectFromGallery = async () => {
  try {
    errorMessage.value = ''

    if (isWeChat.value) {
      // WeChat Mini Program
      const res = await uni.chooseImage({
        count: props.multiple ? 9 : 1,
        sizeType: ['original', 'compressed'],
        sourceType: ['album'],
      })

      if (res.tempFilePaths && res.tempFilePaths.length > 0) {
        await processWeChatImage(res.tempFilePaths[0])
      }
    } else if (isIOS.value) {
      // iOS Capacitor - simplified for demo
      showError('iOS相机功能需要在真机环境中使用')
    } else {
      // Web fallback
      const input = document.createElement('input')
      input.type = 'file'
      input.accept = props.acceptedFormats.join(',')
      input.multiple = props.multiple

      input.onchange = (e) => {
        const files = (e.target as HTMLInputElement).files
        if (files && files.length > 0) {
          processWebImage(files[0])
        }
      }

      input.click()
    }
  } catch (error) {
    console.error('选择图片失败:', error)
    showError('选择图片失败，请重试')
  }
}

const takePhoto = async () => {
  try {
    errorMessage.value = ''

    if (isWeChat.value) {
      // WeChat Mini Program
      const res = await uni.chooseImage({
        count: 1,
        sizeType: ['original'],
        sourceType: ['camera'],
      })

      if (res.tempFilePaths && res.tempFilePaths.length > 0) {
        await processWeChatImage(res.tempFilePaths[0])
      }
    } else if (isIOS.value) {
      // iOS Capacitor - simplified for demo
      showError('iOS拍照功能需要在真机环境中使用')
    } else {
      showError('当前环境不支持拍照功能')
    }
  } catch (error) {
    console.error('拍照失败:', error)
    showError('拍照失败，请重试')
  }
}

const processWeChatImage = async (tempFilePath: string) => {
  try {
    // Get file info
    const fileInfo = await uni.getFileInfo({
      filePath: tempFilePath,
    })

    // Validate file size
    if (fileInfo.size > props.maxSize * 1024 * 1024) {
      showError(`图片大小不能超过 ${props.maxSize}MB`)
      return
    }

    selectedImage.value = {
      file: tempFilePath,
      url: tempFilePath,
      name: `image_${Date.now()}.jpg`,
      size: fileInfo.size,
    }

    await uploadSelectedImage()
  } catch (error) {
    console.error('处理微信图片失败:', error)
    showError('处理图片失败')
  }
}

const processCapacitorImage = async (webPath: string) => {
  try {
    // Convert to blob for size validation
    const response = await fetch(webPath)
    const blob = await response.blob()

    // Validate file size
    if (blob.size > props.maxSize * 1024 * 1024) {
      showError(`图片大小不能超过 ${props.maxSize}MB`)
      return
    }

    selectedImage.value = {
      file: blob,
      url: webPath,
      name: `image_${Date.now()}.jpg`,
      size: blob.size,
    }

    await uploadSelectedImage()
  } catch (error) {
    console.error('处理Capacitor图片失败:', error)
    showError('处理图片失败')
  }
}

const processWebImage = async (file: File) => {
  try {
    // Validate file type
    if (!props.acceptedFormats.includes(file.type)) {
      showError('不支持的文件格式')
      return
    }

    // Validate file size
    if (file.size > props.maxSize * 1024 * 1024) {
      showError(`图片大小不能超过 ${props.maxSize}MB`)
      return
    }

    const url = URL.createObjectURL(file)

    selectedImage.value = {
      file,
      url,
      name: file.name,
      size: file.size,
    }

    await uploadSelectedImage()
  } catch (error) {
    console.error('处理Web图片失败:', error)
    showError('处理图片失败')
  }
}

const uploadSelectedImage = async () => {
  if (!selectedImage.value) return

  try {
    // Mock upload for demo - replace with actual upload logic
    uploadProgress.value = 0
    const interval = setInterval(() => {
      uploadProgress.value += 10
      if (uploadProgress.value >= 100) {
        clearInterval(interval)
        emit('upload', selectedImage.value!)
      }
    }, 200)
  } catch (error) {
    console.error('上传失败:', error)
    showError('上传失败，请重试')
  }
}

const removeImage = () => {
  if (selectedImage.value?.url.startsWith('blob:')) {
    URL.revokeObjectURL(selectedImage.value.url)
  }
  selectedImage.value = null
  errorMessage.value = ''
}

const cropImage = () => {
  showCropModal.value = true
  // Initialize crop canvas in next tick
  nextTick(() => {
    initializeCrop()
  })
}

const initializeCrop = () => {
  if (!cropCanvas.value || !selectedImage.value) return

  const canvas = cropCanvas.value
  const ctx = canvas.getContext('2d')
  const img = new Image()

  img.onload = () => {
    canvas.width = img.width
    canvas.height = img.height
    ctx?.drawImage(img, 0, 0)
  }

  img.src = selectedImage.value.url
}

const applyCrop = () => {
  // Apply crop logic here
  showCropModal.value = false
}

const closeCropModal = () => {
  showCropModal.value = false
}

const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const showError = (message: string) => {
  errorMessage.value = message
  emit('error', message)
  setTimeout(() => {
    errorMessage.value = ''
  }, 5000)
}

// Watch upload progress
watch(uploadProgress, (progress) => {
  emit('progress', progress)
})
</script>

<style lang="scss" scoped>
.image-upload {
  width: 100%;
  max-width: 400px;
  margin: 0 auto;
}

.upload-area {
  position: relative;
  border: 2px dashed var(--color-border);
  border-radius: 16px;
  background: var(--color-background-soft);
  transition: all 0.3s ease;
  overflow: hidden;

  &:hover {
    border-color: var(--color-primary);
    background: var(--color-background-mute);
  }

  &--dragging {
    border-color: var(--color-primary);
    background: var(--color-primary-light);
    transform: scale(1.02);
  }

  &--has-image {
    border: none;
    background: transparent;
  }
}

.upload-placeholder {
  padding: 48px 24px;
  text-align: center;
  cursor: pointer;
}

.upload-icon {
  width: 64px;
  height: 64px;
  margin: 0 auto 16px;
  color: var(--color-primary);
  animation: pulse 2s infinite;
}

.upload-title {
  font-size: 18px;
  font-weight: 600;
  color: var(--color-text);
  margin: 0 0 8px;
}

.upload-subtitle {
  font-size: 14px;
  color: var(--color-text-soft);
  margin: 0 0 24px;
}

.upload-actions {
  display: flex;
  gap: 12px;
  justify-content: center;
  flex-wrap: wrap;
}

.image-preview {
  position: relative;
  aspect-ratio: 1;
  border-radius: 16px;
  overflow: hidden;
}

.preview-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.image-overlay {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  background: linear-gradient(transparent, rgba(0, 0, 0, 0.7));
  padding: 16px;
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
}

.image-info {
  flex: 1;
  color: white;
}

.image-name {
  display: block;
  font-weight: 500;
  margin-bottom: 4px;
}

.image-size {
  font-size: 12px;
  opacity: 0.8;
}

.image-actions {
  display: flex;
  gap: 8px;
}

.btn-icon {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.2);
  border: none;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s ease;

  &:hover {
    background: rgba(255, 255, 255, 0.3);
    transform: scale(1.1);
  }
}

.upload-progress {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  background: rgba(0, 0, 0, 0.8);
  padding: 12px 16px;
  color: white;
}

.progress-bar {
  height: 4px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 2px;
  overflow: hidden;
  margin-bottom: 8px;
}

.progress-fill {
  height: 100%;
  background: var(--color-primary);
  transition: width 0.3s ease;
}

.progress-text {
  font-size: 12px;
}

.error-message {
  margin-top: 12px;
  padding: 12px 16px;
  background: var(--color-error-light);
  color: var(--color-error);
  border-radius: 8px;
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
}

.crop-modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.8);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.crop-container {
  background: white;
  border-radius: 16px;
  max-width: 90vw;
  max-height: 90vh;
  overflow: hidden;
}

.crop-header {
  padding: 16px 20px;
  border-bottom: 1px solid var(--color-border);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.btn-close {
  width: 32px;
  height: 32px;
  border: none;
  background: none;
  font-size: 24px;
  cursor: pointer;
  color: var(--color-text-soft);
}

.crop-area {
  padding: 20px;
  max-height: 60vh;
  overflow: auto;
}

.crop-canvas {
  max-width: 100%;
  height: auto;
}

.crop-actions {
  padding: 16px 20px;
  border-top: 1px solid var(--color-border);
  display: flex;
  gap: 12px;
  justify-content: flex-end;
}

@keyframes pulse {
  0%,
  100% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
}

// Mobile optimizations
@media (max-width: 768px) {
  .upload-placeholder {
    padding: 32px 16px;
  }

  .upload-actions {
    flex-direction: column;
    align-items: center;
  }

  .crop-container {
    margin: 20px;
    max-width: calc(100vw - 40px);
  }
}
</style>
