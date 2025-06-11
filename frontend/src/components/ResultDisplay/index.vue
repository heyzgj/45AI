<template>
  <div class="result-display" :class="{ 'result-display--visible': visible }">
    <div class="result-container">
      <!-- Header -->
      <div class="result-header">
        <h2 class="result-title">✨ 生成完成！</h2>
        <p class="result-subtitle">您的专属 AI 艺术作品已经准备好了</p>
      </div>

      <!-- Image Grid -->
      <div class="image-grid" v-if="images.length > 0">
        <div
          v-for="(image, index) in images"
          :key="image.id"
          class="image-item"
          :style="{ animationDelay: `${index * 0.2}s` }"
          @click="openPreview(image, index)"
        >
          <div class="image-wrapper">
            <img :src="image.url" :alt="`Generated image ${index + 1}`" class="result-image" />
            <div class="image-overlay">
              <div class="image-actions">
                <button class="action-btn" @click.stop="downloadImage(image)" title="下载">
                  <i class="icon-download"></i>
                </button>
                <button class="action-btn" @click.stop="shareImage(image)" title="分享">
                  <i class="icon-share"></i>
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Action Buttons -->
      <div class="result-actions">
        <button class="btn btn--secondary" @click="regenerate">
          <i class="icon-refresh"></i>
          重新生成
        </button>
        <button class="btn btn--primary" @click="generateMore">
          <i class="icon-plus"></i>
          生成更多
        </button>
      </div>
    </div>

    <!-- Image Preview Modal -->
    <div v-if="previewImage" class="preview-modal" @click="closePreview">
      <div class="preview-container" @click.stop>
        <div class="preview-header">
          <span class="image-counter">{{ currentImageIndex + 1 }} / {{ images.length }}</span>
          <button class="close-btn" @click="closePreview">×</button>
        </div>
        <div class="preview-content">
          <img :src="previewImage.url" class="preview-image" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

interface GeneratedImage {
  id: string
  url: string
  width: number
  height: number
  size: number
}

interface Props {
  visible?: boolean
  images?: GeneratedImage[]
}

const props = withDefaults(defineProps<Props>(), {
  visible: false,
  images: () => [],
})

const emit = defineEmits<{
  regenerate: []
  generateMore: []
}>()

const previewImage = ref<GeneratedImage | null>(null)
const currentImageIndex = ref(0)

const openPreview = (image: GeneratedImage, index: number) => {
  previewImage.value = image
  currentImageIndex.value = index
}

const closePreview = () => {
  previewImage.value = null
}

const downloadImage = async (image: GeneratedImage) => {
  try {
    const response = await fetch(image.url)
    const blob = await response.blob()
    const url = window.URL.createObjectURL(blob)

    const link = document.createElement('a')
    link.href = url
    link.download = `ai-generated-${image.id}.jpg`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
  } catch (error) {
    console.error('下载失败:', error)
  }
}

const shareImage = async (image: GeneratedImage) => {
  try {
    if (navigator.share) {
      await navigator.share({
        title: '我的 AI 艺术作品',
        text: '看看我用 AI 生成的这张图片！',
        url: image.url,
      })
    } else {
      await navigator.clipboard.writeText(image.url)
    }
  } catch (error) {
    console.error('分享失败:', error)
  }
}

const regenerate = () => emit('regenerate')
const generateMore = () => emit('generateMore')
</script>

<style lang="scss" scoped>
.result-display {
  opacity: 0;
  transform: translateY(20px);
  transition: all 0.5s ease;

  &--visible {
    opacity: 1;
    transform: translateY(0);
  }
}

.result-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

.result-header {
  text-align: center;
  margin-bottom: 32px;
}

.result-title {
  font-size: 28px;
  font-weight: 700;
  color: var(--color-text);
  margin: 0 0 8px;
}

.result-subtitle {
  font-size: 16px;
  color: var(--color-text-soft);
  margin: 0;
}

.image-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
  margin-bottom: 32px;
}

.image-item {
  animation: fadeInUp 0.6s ease forwards;
  opacity: 0;
  transform: translateY(30px);
}

.image-wrapper {
  position: relative;
  border-radius: 16px;
  overflow: hidden;
  cursor: pointer;
  transition: transform 0.3s ease;

  &:hover {
    transform: translateY(-4px);

    .image-overlay {
      opacity: 1;
    }
  }
}

.result-image {
  width: 100%;
  height: auto;
  aspect-ratio: 1;
  object-fit: cover;
  display: block;
}

.image-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(to bottom, transparent 60%, rgba(0, 0, 0, 0.8) 100%);
  opacity: 0;
  transition: opacity 0.3s ease;
  display: flex;
  align-items: flex-end;
  padding: 16px;
}

.image-actions {
  display: flex;
  gap: 8px;
  margin-left: auto;
}

.action-btn {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.9);
  border: none;
  color: var(--color-text);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s ease;

  &:hover {
    background: white;
    transform: scale(1.1);
  }
}

.result-actions {
  display: flex;
  gap: 16px;
  justify-content: center;
  flex-wrap: wrap;
}

.preview-modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.9);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 10000;
  padding: 20px;
}

.preview-container {
  max-width: 90vw;
  max-height: 90vh;
  background: white;
  border-radius: 16px;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.preview-header {
  padding: 16px 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid var(--color-border);
}

.image-counter {
  font-size: 14px;
  color: var(--color-text-soft);
  font-weight: 500;
}

.close-btn {
  width: 32px;
  height: 32px;
  border: none;
  background: none;
  font-size: 24px;
  cursor: pointer;
  color: var(--color-text-soft);
}

.preview-content {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
  min-height: 400px;
}

.preview-image {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
  border-radius: 8px;
}

@keyframes fadeInUp {
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@media (max-width: 768px) {
  .result-container {
    padding: 16px;
  }

  .image-grid {
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 16px;
  }

  .result-actions {
    flex-direction: column;
    align-items: center;
  }
}
</style>
