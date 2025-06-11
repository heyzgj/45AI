<template>
  <view
    class="credit-display"
    :class="[
      `size-${size}`,
      {
        clickable: clickable,
        loading: loading,
      },
    ]"
    @click="handleClick"
  >
    <!-- Loading State -->
    <view v-if="loading" class="credit-loading">
      <view class="loading-dots">
        <view class="dot"></view>
        <view class="dot"></view>
        <view class="dot"></view>
      </view>
    </view>

    <!-- Credit Display -->
    <view v-else class="credit-content">
      <text class="credit-icon">🎞️</text>
      <text class="credit-amount" :class="{ changing: isChanging }">
        {{ animatedCredits }}
      </text>
      <text v-if="showLabel" class="credit-label">credits</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useUserStore } from '@/store/user'
import { customTransitions } from '@/utils/navigation'

// Props
interface Props {
  size?: 'small' | 'medium' | 'large'
  clickable?: boolean
  showLabel?: boolean
  showWarning?: boolean
  loading?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  size: 'medium',
  clickable: true,
  showLabel: true,
  showWarning: true,
  loading: false,
})

// Store
const userStore = useUserStore()

// State
const animatedCredits = ref(userStore.credits)
const isChanging = ref(false)

// No low balance logic needed

// Watch for credit changes and animate
watch(
  () => userStore.credits,
  (newValue, oldValue) => {
    if (newValue !== oldValue) {
      animateCredits(oldValue, newValue)
    }
  },
)

// Animate credit number changes
const animateCredits = (from: number, to: number) => {
  isChanging.value = true
  const duration = 600
  const startTime = Date.now()
  const diff = to - from

  const animate = () => {
    const elapsed = Date.now() - startTime
    const progress = Math.min(elapsed / duration, 1)

    // Easing function
    const easeOutCubic = (t: number) => 1 - Math.pow(1 - t, 3)
    const easedProgress = easeOutCubic(progress)

    animatedCredits.value = Math.round(from + diff * easedProgress)

    if (progress < 1) {
      requestAnimationFrame(animate)
    } else {
      animatedCredits.value = to
      setTimeout(() => {
        isChanging.value = false
      }, 200)
    }
  }

  requestAnimationFrame(animate)
}

// Methods
const handleClick = () => {
  if (props.clickable) {
    customTransitions.toLogin() // Navigate to purchase page
    // TODO: Change to purchase page when implemented
  }
}
</script>

<style lang="scss" scoped>
@import '@/style/variables.scss';

.credit-display {
  display: inline-flex;
  align-items: center;
  gap: $spacing-xs;
  transition: all $duration-fast $ease-custom;

  &.clickable {
    cursor: pointer;

    &:hover {
      transform: translateY(-1px);

      .credit-icon {
        transform: scale(1.1);
      }
    }

    &:active {
      transform: scale(0.95);
    }
  }

  &.loading {
    min-width: 80px;
  }
}

// Credit Content
.credit-content {
  display: flex;
  align-items: center;
  gap: $spacing-xs;
}

.credit-icon {
  transition: transform $duration-fast $ease-custom;

  &.pulse {
    animation: pulseSoft 2s $ease-custom infinite;
  }
}

.credit-amount {
  font-weight: $font-weight-semibold;
  color: $color-primary;
  transition: all $duration-fast $ease-custom;

  &.changing {
    transform: scale(1.1);
    color: darken($color-primary, 10%);
  }
}

.credit-label {
  color: $color-text-subtle;
}

// Size Variants
.size-small {
  .credit-icon {
    font-size: 16px;
  }

  .credit-amount {
    font-size: $font-size-body;
  }

  .credit-label {
    font-size: $font-size-caption;
  }
}

.size-medium {
  .credit-icon {
    font-size: 20px;
  }

  .credit-amount {
    font-size: 18px;
  }

  .credit-label {
    font-size: $font-size-body;
  }
}

.size-large {
  gap: $spacing-sm;

  .credit-icon {
    font-size: 28px;
  }

  .credit-amount {
    font-size: 24px;
  }

  .credit-label {
    font-size: 16px;
  }
}

// Loading State
.credit-loading {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 24px;
}

.loading-dots {
  display: flex;
  gap: 4px;

  .dot {
    width: 6px;
    height: 6px;
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

// Low balance warnings removed

// Loading animation
@keyframes loadingDot {
  0%,
  80%,
  100% {
    transform: scale(0);
    opacity: 0;
  }
  40% {
    transform: scale(1);
    opacity: 1;
  }
}
</style>
