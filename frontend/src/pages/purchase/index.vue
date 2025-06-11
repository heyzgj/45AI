<template>
  <view class="purchase-container page-wrapper animate-fade-in">
    <!-- Page Title -->
    <view class="page-header">
      <text class="page-title">购买胶卷</text>
    </view>

    <!-- Current Credits Display -->
    <view class="current-credits">
      <view class="credits-card">
        <text class="credits-label">当前胶卷</text>
        <text class="credits-amount">{{ userCredits }}</text>
        <text class="credits-unit">个</text>
      </view>
    </view>

    <!-- Credit Packs -->
    <view class="packs-section">
      <text class="section-title">选择套餐</text>
      <view class="packs-grid">
        <view
          v-for="pack in creditPacks"
          :key="pack.id"
          class="pack-card"
          :class="{ popular: pack.popular, selected: selectedPack?.id === pack.id }"
          @tap="selectPack(pack)"
        >
          <!-- Popular Badge -->
          <view v-if="pack.popular" class="popular-badge">
            <text class="badge-text">热门</text>
          </view>

          <!-- Pack Content -->
          <view class="pack-content">
            <text class="pack-title">{{ pack.title }}</text>
            <text class="pack-description">{{ pack.description }}</text>

            <view class="pack-credits">
              <text class="credits-number">{{ pack.credits }}</text>
              <text class="credits-text">胶卷</text>
            </view>

            <view class="pack-pricing">
              <text class="current-price">{{ formatPrice(pack.amount) }}</text>
              <text v-if="pack.originalPrice" class="original-price">
                {{ formatPrice(pack.originalPrice) }}
              </text>
            </view>

            <!-- Savings -->
            <view v-if="pack.originalPrice" class="savings">
              <text class="savings-text">
                省 {{ formatPrice(pack.originalPrice - pack.amount) }}
              </text>
            </view>
          </view>
        </view>
      </view>
    </view>

    <!-- Purchase Button -->
    <view class="purchase-section">
      <view
        class="purchase-btn"
        :class="{ disabled: !selectedPack || purchasing }"
        @tap="handlePurchase"
      >
        <text class="btn-text">
          {{
            purchasing ? '处理中...' : selectedPack ? `购买 ${selectedPack.title}` : '请选择套餐'
          }}
        </text>
      </view>

      <!-- Payment Methods -->
      <view class="payment-methods">
        <!-- #ifdef MP-WEIXIN -->
        <view class="payment-method">
          <text class="method-icon">💳</text>
          <text class="method-text">微信支付</text>
        </view>
        <!-- #endif -->

        <!-- #ifdef APP-IOS -->
        <view class="payment-method">
          <text class="method-icon">🍎</text>
          <text class="method-text">Apple Pay</text>
        </view>
        <!-- #endif -->
      </view>
    </view>

    <!-- Loading Overlay -->
    <view v-if="purchasing" class="loading-overlay">
      <view class="loading-content">
        <view class="loading-spinner"></view>
        <text class="loading-text">正在处理支付...</text>
      </view>
    </view>
  </view>
</template>

<script>
import {
  getCreditPackages,
  getAppleIAPProducts,
  completePayment,
  completeApplePurchase,
  formatPrice,
  formatUSDPrice,
} from '@/api/payment'
import { getUserInfo } from '@/api/user'

export default {
  name: 'PurchasePage',
  data() {
    return {
      userCredits: 0,
      creditPacks: [],
      selectedPack: null,
      purchasing: false,
      platform: 'wechat', // Default platform
    }
  },

  onLoad() {
    this.initPage()
  },

  methods: {
    async initPage() {
      try {
        // Determine platform
        // #ifdef MP-WEIXIN
        this.platform = 'wechat'
        this.creditPacks = getCreditPackages()
        // #endif

        // #ifdef APP-IOS
        this.platform = 'apple'
        this.creditPacks = getAppleIAPProducts()
        // #endif

        // #ifndef MP-WEIXIN
        // #ifndef APP-IOS
        // For development/testing
        this.platform = 'wechat'
        this.creditPacks = getCreditPackages()
        // #endif
        // #endif

        // Load user credits
        await this.loadUserCredits()
      } catch (error) {
        console.error('Failed to initialize purchase page:', error)
        uni.showToast({
          title: '加载失败',
          icon: 'error',
        })
      }
    },

    async loadUserCredits() {
      try {
        const response = await getUserInfo()
        if (response.success && response.data) {
          this.userCredits = response.data.credits || 0
        }
      } catch (error) {
        console.error('Failed to load user credits:', error)
      }
    },

    selectPack(pack) {
      // Radio button behavior - only one can be selected
      if (this.selectedPack && this.selectedPack.id === pack.id) {
        // If same pack clicked, deselect
        this.selectedPack = null
      } else {
        // Select new pack (deselects others)
        this.selectedPack = pack
      }
    },

    async handlePurchase() {
      if (!this.selectedPack || this.purchasing) {
        return
      }

      this.purchasing = true

      try {
        let result

        if (this.platform === 'wechat') {
          result = await completePayment({
            amount: this.selectedPack.amount,
            description: `购买${this.selectedPack.title}`,
          })
        } else if (this.platform === 'apple') {
          result = await completeApplePurchase({
            product_id: this.selectedPack.product_id,
            title: this.selectedPack.title,
            description: `购买${this.selectedPack.title}`,
          })
        }

        if (result.success) {
          // Success
          uni.showToast({
            title: `成功购买${result.credits}胶卷`,
            icon: 'success',
            duration: 2000,
          })

          // Update user credits
          this.userCredits += result.credits

          // Reset selection
          this.selectedPack = null

          // Navigate back after delay
          setTimeout(() => {
            this.goBack()
          }, 2000)
        } else {
          // Failed
          uni.showToast({
            title: result.message || '购买失败',
            icon: 'error',
          })
        }
      } catch (error) {
        console.error('Purchase failed:', error)
        uni.showToast({
          title: '购买失败，请重试',
          icon: 'error',
        })
      } finally {
        this.purchasing = false
      }
    },

    formatPrice(amountInCents) {
      if (this.platform === 'apple') {
        return formatUSDPrice(amountInCents)
      }
      return formatPrice(amountInCents)
    },

    goBack() {
      uni.navigateBack()
    },
  },
}
</script>

<style lang="scss" scoped>
.purchase-container {
  min-height: 100vh;
  background: var(--color-bg);
  padding-top: calc(env(safe-area-inset-top, 0px) + 20px);
  padding-bottom: 120rpx;
}

.page-header {
  padding: 0 20px 16px;
}

.page-title {
  font-size: 28px;
  font-weight: 700;
  color: var(--color-text);
  display: block;
}

.current-credits {
  padding: 48rpx 40rpx;

  .credits-card {
    background: linear-gradient(135deg, var(--color-secondary), var(--color-primary));
    border-radius: 32rpx;
    padding: 48rpx;
    display: flex;
    flex-direction: column;
    align-items: center;
    box-shadow: 0 8rpx 48rpx rgba(232, 155, 147, 0.15);

    .credits-label {
      font-size: 28rpx;
      color: rgba(255, 255, 255, 0.8);
      margin-bottom: 16rpx;
    }

    .credits-amount {
      font-size: 72rpx;
      font-weight: 600;
      color: white;
      line-height: 1;
    }

    .credits-unit {
      font-size: 24rpx;
      color: rgba(255, 255, 255, 0.8);
      margin-top: 8rpx;
    }
  }
}

.packs-section {
  padding: 0 40rpx;

  .section-title {
    font-size: 36rpx;
    font-weight: 600;
    color: var(--color-text);
    margin-bottom: 32rpx;
    display: block;
  }

  .packs-grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    grid-template-rows: repeat(2, 1fr);
    gap: 16rpx;
  }

  .pack-card {
    position: relative;
    background: var(--color-surface);
    border-radius: 32rpx;
    border: 2rpx solid transparent;
    box-shadow: 0 8rpx 48rpx rgba(74, 74, 74, 0.08);
    overflow: hidden;
    transition: all 0.3s cubic-bezier(0.6, 0.05, 0.4, 1);

    &.popular {
      border-color: var(--color-primary);
      box-shadow: 0 8rpx 48rpx rgba(232, 155, 147, 0.2);
    }

    &.selected {
      border-color: var(--color-primary);
      transform: scale(1.02) translateY(-8rpx);
      box-shadow: 0 16rpx 64rpx rgba(232, 155, 147, 0.25);
    }

    .popular-badge {
      position: absolute;
      top: 24rpx;
      right: 24rpx;
      background: var(--color-primary);
      border-radius: 24rpx;
      padding: 8rpx 16rpx;

      .badge-text {
        font-size: 20rpx;
        color: white;
        font-weight: 600;
      }
    }

    .pack-content {
      padding: 48rpx;

      .pack-title {
        font-size: 40rpx;
        font-weight: 600;
        color: var(--color-text);
        margin-bottom: 8rpx;
        display: block;
      }

      .pack-description {
        font-size: 24rpx;
        color: var(--color-text-subtle);
        margin-bottom: 32rpx;
        display: block;
      }

      .pack-credits {
        display: flex;
        align-items: baseline;
        margin-bottom: 24rpx;

        .credits-number {
          font-size: 64rpx;
          font-weight: 600;
          color: var(--color-primary);
          line-height: 1;
        }

        .credits-text {
          font-size: 28rpx;
          color: var(--color-text-subtle);
          margin-left: 8rpx;
        }
      }

      .pack-pricing {
        display: flex;
        align-items: center;
        gap: 16rpx;
        margin-bottom: 16rpx;

        .current-price {
          font-size: 36rpx;
          font-weight: 600;
          color: var(--color-text);
        }

        .original-price {
          font-size: 28rpx;
          color: var(--color-text-subtle);
          text-decoration: line-through;
        }
      }

      .savings {
        .savings-text {
          font-size: 24rpx;
          color: var(--color-primary);
          font-weight: 600;
        }
      }
    }
  }
}

.purchase-section {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  background: var(--color-surface);
  border-top: 1rpx solid #f0f0f0;
  padding: 32rpx 40rpx;
  padding-bottom: calc(32rpx + env(safe-area-inset-bottom));

  .purchase-btn {
    background: var(--color-primary);
    border-radius: 48rpx;
    height: 96rpx;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: 0 8rpx 24rpx rgba(232, 155, 147, 0.3);
    transition: all 0.3s cubic-bezier(0.6, 0.05, 0.4, 1);
    margin-bottom: 24rpx;

    &:active {
      transform: scale(0.97);
      box-shadow: 0 4rpx 16rpx rgba(232, 155, 147, 0.2);
    }

    &.disabled {
      background: #e0e0e0;
      box-shadow: none;
      transform: none;

      .btn-text {
        color: var(--color-text-subtle);
      }
    }

    .btn-text {
      font-size: 32rpx;
      font-weight: 600;
      color: white;
    }
  }

  .payment-methods {
    display: flex;
    justify-content: center;
    gap: 32rpx;

    .payment-method {
      display: flex;
      align-items: center;
      gap: 8rpx;

      .method-icon {
        font-size: 24rpx;
      }

      .method-text {
        font-size: 24rpx;
        color: var(--color-text-subtle);
      }
    }
  }
}

.loading-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;

  .loading-content {
    background: var(--color-surface);
    border-radius: 32rpx;
    padding: 64rpx;
    display: flex;
    flex-direction: column;
    align-items: center;

    .loading-spinner {
      width: 64rpx;
      height: 64rpx;
      border: 4rpx solid var(--color-secondary);
      border-top: 4rpx solid var(--color-primary);
      border-radius: 50%;
      animation: spin 1s linear infinite;
      margin-bottom: 32rpx;
    }

    .loading-text {
      font-size: 28rpx;
      color: var(--color-text);
    }
  }
}

@keyframes spin {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}

/* CSS Variables (should be defined globally) */
:root {
  --color-primary: #e89b93;
  --color-secondary: #f3d9d7;
  --color-bg: #fcfbf9;
  --color-surface: #ffffff;
  --color-text: #4a4a4a;
  --color-text-subtle: #9b9b9b;
}
</style>
