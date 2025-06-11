<template>
  <div class="loading-screen" :class="{ 'loading-screen--visible': visible }">
    <div class="loading-container">
      <!-- Main Loading Animation -->
      <div class="loading-animation">
        <div class="pulse-rings">
          <div class="pulse-ring pulse-ring--1"></div>
          <div class="pulse-ring pulse-ring--2"></div>
          <div class="pulse-ring pulse-ring--3"></div>
        </div>
        <div class="loading-icon">
          <svg width="64" height="64" viewBox="0 0 24 24" fill="none">
            <path
              d="M12 2L13.09 8.26L22 9L13.09 9.74L12 16L10.91 9.74L2 9L10.91 8.26L12 2Z"
              fill="currentColor"
            />
          </svg>
        </div>
      </div>

      <!-- Progress Information -->
      <div class="progress-info">
        <h2 class="loading-title">{{ currentStage.title }}</h2>
        <p class="loading-subtitle">{{ currentStage.description }}</p>

        <!-- Progress Bar -->
        <div class="progress-container">
          <div class="progress-bar">
            <div class="progress-fill" :style="{ width: progress + '%' }"></div>
          </div>
          <div class="progress-text">
            <span class="progress-percentage">{{ Math.round(progress) }}%</span>
            <span class="progress-time">预计剩余 {{ estimatedTime }}s</span>
          </div>
        </div>

        <!-- Stage Indicators -->
        <div class="stage-indicators">
          <div
            v-for="(stage, index) in stages"
            :key="index"
            class="stage-indicator"
            :class="{
              'stage-indicator--active': index === currentStageIndex,
              'stage-indicator--completed': index < currentStageIndex,
            }"
          >
            <div class="stage-dot">
              <i v-if="index < currentStageIndex" class="icon-check"></i>
              <span v-else>{{ index + 1 }}</span>
            </div>
            <span class="stage-label">{{ stage.name }}</span>
          </div>
        </div>
      </div>

      <!-- Fun Facts / Tips -->
      <div class="loading-tips">
        <div class="tip-container">
          <transition name="fade" mode="out-in">
            <div :key="currentTipIndex" class="tip-content">
              <i class="tip-icon icon-lightbulb"></i>
              <p class="tip-text">{{ currentTip }}</p>
            </div>
          </transition>
        </div>
      </div>

      <!-- Cancel Button -->
      <button v-if="showCancel" class="cancel-button" @click="handleCancel">取消生成</button>
    </div>

    <!-- Background Effects -->
    <div class="background-effects">
      <div class="gradient-orb gradient-orb--1"></div>
      <div class="gradient-orb gradient-orb--2"></div>
      <div class="gradient-orb gradient-orb--3"></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'

interface LoadingStage {
  name: string
  title: string
  description: string
  duration: number // seconds
}

interface Props {
  visible?: boolean
  showCancel?: boolean
  totalDuration?: number // seconds
}

const props = withDefaults(defineProps<Props>(), {
  visible: false,
  showCancel: true,
  totalDuration: 30,
})

const emit = defineEmits<{
  cancel: []
  complete: []
}>()

// Loading stages
const stages: LoadingStage[] = [
  {
    name: '准备',
    title: '正在准备您的图片',
    description: '分析图片内容和风格要求',
    duration: 5,
  },
  {
    name: '处理',
    title: 'AI 正在创作中',
    description: '运用先进算法生成独特图像',
    duration: 20,
  },
  {
    name: '优化',
    title: '精细调整画面',
    description: '优化细节和色彩表现',
    duration: 5,
  },
]

// Fun tips to show during loading
const tips = [
  '💡 AI 图像生成基于深度学习神经网络',
  '🎨 每张图片都是独一无二的艺术创作',
  '⚡ 我们的 GPU 集群每秒处理数百万次计算',
  '🌟 尝试不同的提示词可以获得意想不到的效果',
  '🎯 更详细的描述通常能产生更好的结果',
  '🔥 我们的模型训练了数百万张高质量图片',
  '✨ 创意无限，每次生成都是新的探索',
]

// State
const progress = ref(0)
const currentStageIndex = ref(0)
const currentTipIndex = ref(0)
const startTime = ref(0)
const progressInterval = ref<NodeJS.Timeout>()
const tipInterval = ref<NodeJS.Timeout>()

// Computed
const currentStage = computed(() => stages[currentStageIndex.value])
const currentTip = computed(() => tips[currentTipIndex.value])

const estimatedTime = computed(() => {
  const elapsed = (Date.now() - startTime.value) / 1000
  const remaining = Math.max(0, props.totalDuration - elapsed)
  return Math.ceil(remaining)
})

// Methods
const startProgress = () => {
  startTime.value = Date.now()
  progress.value = 0
  currentStageIndex.value = 0

  // Progress simulation
  progressInterval.value = setInterval(() => {
    const elapsed = (Date.now() - startTime.value) / 1000
    const newProgress = Math.min(100, (elapsed / props.totalDuration) * 100)
    progress.value = newProgress

    // Update stage based on progress
    let cumulativeDuration = 0
    for (let i = 0; i < stages.length; i++) {
      cumulativeDuration += stages[i].duration
      if (elapsed <= cumulativeDuration) {
        currentStageIndex.value = i
        break
      }
    }

    // Complete when progress reaches 100%
    if (newProgress >= 100) {
      clearInterval(progressInterval.value)
      setTimeout(() => {
        emit('complete')
      }, 500)
    }
  }, 100)
}

const startTipRotation = () => {
  tipInterval.value = setInterval(() => {
    currentTipIndex.value = (currentTipIndex.value + 1) % tips.length
  }, 4000)
}

const stopProgress = () => {
  if (progressInterval.value) {
    clearInterval(progressInterval.value)
    progressInterval.value = undefined
  }
  if (tipInterval.value) {
    clearInterval(tipInterval.value)
    tipInterval.value = undefined
  }
}

const handleCancel = () => {
  stopProgress()
  emit('cancel')
}

// Watchers
watch(
  () => props.visible,
  (visible) => {
    if (visible) {
      startProgress()
      startTipRotation()
    } else {
      stopProgress()
    }
  },
)

// Lifecycle
onMounted(() => {
  if (props.visible) {
    startProgress()
    startTipRotation()
  }
})

onUnmounted(() => {
  stopProgress()
})
</script>

<style lang="scss" scoped>
.loading-screen {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(
    135deg,
    rgba(99, 102, 241, 0.1) 0%,
    rgba(168, 85, 247, 0.1) 50%,
    rgba(236, 72, 153, 0.1) 100%
  );
  backdrop-filter: blur(20px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
  opacity: 0;
  visibility: hidden;
  transition: all 0.3s ease;

  &--visible {
    opacity: 1;
    visibility: visible;
  }
}

.loading-container {
  text-align: center;
  max-width: 400px;
  padding: 40px 20px;
  position: relative;
  z-index: 2;
}

.loading-animation {
  position: relative;
  margin-bottom: 40px;
}

.pulse-rings {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
}

.pulse-ring {
  position: absolute;
  border: 2px solid var(--color-primary);
  border-radius: 50%;
  opacity: 0;
  animation: pulse 2s infinite;

  &--1 {
    width: 80px;
    height: 80px;
    margin: -40px 0 0 -40px;
    animation-delay: 0s;
  }

  &--2 {
    width: 120px;
    height: 120px;
    margin: -60px 0 0 -60px;
    animation-delay: 0.5s;
  }

  &--3 {
    width: 160px;
    height: 160px;
    margin: -80px 0 0 -80px;
    animation-delay: 1s;
  }
}

.loading-icon {
  width: 64px;
  height: 64px;
  margin: 0 auto;
  color: var(--color-primary);
  animation: rotate 3s linear infinite;
  position: relative;
  z-index: 1;
}

.progress-info {
  margin-bottom: 32px;
}

.loading-title {
  font-size: 24px;
  font-weight: 600;
  color: var(--color-text);
  margin: 0 0 8px;
}

.loading-subtitle {
  font-size: 16px;
  color: var(--color-text-soft);
  margin: 0 0 24px;
}

.progress-container {
  margin-bottom: 24px;
}

.progress-bar {
  height: 6px;
  background: var(--color-background-mute);
  border-radius: 3px;
  overflow: hidden;
  margin-bottom: 8px;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, var(--color-primary), var(--color-primary-light));
  border-radius: 3px;
  transition: width 0.3s ease;
  position: relative;

  &::after {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.3), transparent);
    animation: shimmer 2s infinite;
  }
}

.progress-text {
  display: flex;
  justify-content: space-between;
  font-size: 14px;
}

.progress-percentage {
  font-weight: 600;
  color: var(--color-primary);
}

.progress-time {
  color: var(--color-text-soft);
}

.stage-indicators {
  display: flex;
  justify-content: center;
  gap: 24px;
  margin-bottom: 32px;
}

.stage-indicator {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  opacity: 0.5;
  transition: all 0.3s ease;

  &--active {
    opacity: 1;
    transform: scale(1.1);
  }

  &--completed {
    opacity: 0.8;
  }
}

.stage-dot {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: var(--color-background-mute);
  border: 2px solid var(--color-border);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 600;
  transition: all 0.3s ease;

  .stage-indicator--active & {
    background: var(--color-primary);
    border-color: var(--color-primary);
    color: white;
    animation: pulse-dot 2s infinite;
  }

  .stage-indicator--completed & {
    background: var(--color-success);
    border-color: var(--color-success);
    color: white;
  }
}

.stage-label {
  font-size: 12px;
  color: var(--color-text-soft);

  .stage-indicator--active & {
    color: var(--color-text);
    font-weight: 500;
  }
}

.loading-tips {
  margin-bottom: 32px;
}

.tip-container {
  min-height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.tip-content {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px 20px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  backdrop-filter: blur(10px);
}

.tip-icon {
  font-size: 20px;
  color: var(--color-warning);
}

.tip-text {
  font-size: 14px;
  color: var(--color-text);
  margin: 0;
}

.cancel-button {
  padding: 12px 24px;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 8px;
  color: var(--color-text);
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s ease;
  backdrop-filter: blur(10px);

  &:hover {
    background: rgba(255, 255, 255, 0.2);
    border-color: rgba(255, 255, 255, 0.3);
  }
}

.background-effects {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  overflow: hidden;
  z-index: 1;
}

.gradient-orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(60px);
  opacity: 0.3;
  animation: float 6s ease-in-out infinite;

  &--1 {
    width: 200px;
    height: 200px;
    background: var(--color-primary);
    top: 20%;
    left: 10%;
    animation-delay: 0s;
  }

  &--2 {
    width: 150px;
    height: 150px;
    background: var(--color-secondary);
    top: 60%;
    right: 20%;
    animation-delay: 2s;
  }

  &--3 {
    width: 100px;
    height: 100px;
    background: var(--color-accent);
    bottom: 20%;
    left: 30%;
    animation-delay: 4s;
  }
}

// Animations
@keyframes pulse {
  0% {
    opacity: 0;
    transform: scale(0.8);
  }
  50% {
    opacity: 1;
  }
  100% {
    opacity: 0;
    transform: scale(1.2);
  }
}

@keyframes rotate {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

@keyframes shimmer {
  0% {
    transform: translateX(-100%);
  }
  100% {
    transform: translateX(100%);
  }
}

@keyframes pulse-dot {
  0%,
  100% {
    transform: scale(1);
  }
  50% {
    transform: scale(1.1);
  }
}

@keyframes float {
  0%,
  100% {
    transform: translateY(0px);
  }
  50% {
    transform: translateY(-20px);
  }
}

// Transitions
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.5s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

// Mobile optimizations
@media (max-width: 768px) {
  .loading-container {
    padding: 20px;
    max-width: 320px;
  }

  .loading-title {
    font-size: 20px;
  }

  .loading-subtitle {
    font-size: 14px;
  }

  .stage-indicators {
    gap: 16px;
  }

  .stage-dot {
    width: 28px;
    height: 28px;
  }

  .tip-content {
    padding: 12px 16px;
  }

  .tip-text {
    font-size: 13px;
  }
}
</style>
