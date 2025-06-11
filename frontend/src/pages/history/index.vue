<template>
  <view class="history-container page-wrapper animate-fade-in">
    <!-- Page Title -->
    <view class="page-header">
      <text class="page-title">交易记录</text>
    </view>

    <!-- Summary Card -->
    <view class="summary-section">
      <view class="summary-card">
        <view class="summary-item">
          <text class="summary-label">当前胶卷</text>
          <text class="summary-value">{{ userCredits }}</text>
        </view>
        <view class="summary-divider"></view>
        <view class="summary-item">
          <text class="summary-label">总交易</text>
          <text class="summary-value">{{ totalTransactions }}</text>
        </view>
      </view>
    </view>

    <!-- Filter Tabs -->
    <view class="filter-section">
      <view class="filter-tabs">
        <view
          v-for="filter in filters"
          :key="filter.key"
          class="filter-tab"
          :class="{ active: activeFilter === filter.key }"
          @tap="setFilter(filter.key)"
        >
          <text class="filter-text">{{ filter.label }}</text>
        </view>
      </view>
    </view>

    <!-- Transaction List -->
    <view class="transactions-section">
      <view v-if="loading && transactions.length === 0" class="loading-state">
        <view class="loading-spinner"></view>
        <text class="loading-text">加载中...</text>
      </view>

      <view v-else-if="transactions.length === 0" class="empty-state">
        <text class="empty-icon">📝</text>
        <text class="empty-title">暂无交易记录</text>
        <text class="empty-description">您还没有任何交易记录</text>
      </view>

      <view v-else class="transactions-list">
        <view
          v-for="transaction in filteredTransactions"
          :key="transaction.id"
          class="transaction-item"
          :class="{
            purchase: transaction.type === 'purchase',
            generation: transaction.type === 'generation',
          }"
        >
          <!-- Transaction Icon -->
          <view class="transaction-icon">
            <text class="icon-text">{{ getTransactionIcon(transaction.type) }}</text>
          </view>

          <!-- Transaction Details -->
          <view class="transaction-details">
            <text class="transaction-description">{{ transaction.description }}</text>
            <text class="transaction-date">{{ formatDate(transaction.created_at) }}</text>
            <view v-if="transaction.external_payment_id" class="transaction-id">
              <text class="id-label">订单号:</text>
              <text class="id-value">{{ transaction.external_payment_id.slice(-8) }}</text>
            </view>
          </view>

          <!-- Transaction Amount -->
          <view class="transaction-amount">
            <text
              class="amount-text"
              :class="{ positive: transaction.amount > 0, negative: transaction.amount < 0 }"
            >
              {{ transaction.amount > 0 ? '+' : '' }}{{ transaction.amount }}
            </text>
            <text class="amount-unit">胶卷</text>
          </view>
        </view>

        <!-- Load More -->
        <view v-if="hasMore" class="load-more" @tap="loadMore">
          <view v-if="loadingMore" class="loading-more">
            <view class="loading-spinner small"></view>
            <text class="loading-text">加载更多...</text>
          </view>
          <text v-else class="load-more-text">点击加载更多</text>
        </view>

        <view v-else-if="transactions.length > 0" class="end-message">
          <text class="end-text">已显示全部记录</text>
        </view>
      </view>
    </view>
  </view>
</template>

<script>
import { getUserTransactions } from '@/api/user'
import { getUserInfo } from '@/api/user'

export default {
  name: 'HistoryPage',
  data() {
    return {
      userCredits: 0,
      transactions: [],
      totalTransactions: 0,
      loading: false,
      loadingMore: false,
      hasMore: true,
      offset: 0,
      limit: 20,
      activeFilter: 'all',
      filters: [
        { key: 'all', label: '全部' },
        { key: 'purchase', label: '充值' },
        { key: 'generation', label: '消费' },
      ],
    }
  },

  computed: {
    filteredTransactions() {
      if (this.activeFilter === 'all') {
        return this.transactions
      }
      return this.transactions.filter((t) => t.type === this.activeFilter)
    },
  },

  onLoad() {
    this.initPage()
  },

  onReachBottom() {
    if (this.hasMore && !this.loadingMore) {
      this.loadMore()
    }
  },

  methods: {
    async initPage() {
      try {
        await Promise.all([this.loadUserInfo(), this.loadTransactions(true)])
      } catch (error) {
        console.error('Failed to initialize history page:', error)
        uni.showToast({
          title: '加载失败',
          icon: 'error',
        })
      }
    },

    async loadUserInfo() {
      try {
        const response = await getUserInfo()
        if (response.success && response.data) {
          this.userCredits = response.data.credits || 0
        }
      } catch (error) {
        console.error('Failed to load user info:', error)
      }
    },

    async loadTransactions(reset = false) {
      if (this.loading || this.loadingMore) return

      if (reset) {
        this.loading = true
        this.offset = 0
        this.transactions = []
      } else {
        this.loadingMore = true
      }

      try {
        const response = await getUserTransactions({
          limit: this.limit,
          offset: this.offset,
        })

        if (response.success && response.data) {
          const newTransactions = response.data.transactions || []
          const pagination = response.data.pagination || {}

          if (reset) {
            this.transactions = newTransactions
          } else {
            this.transactions.push(...newTransactions)
          }

          this.totalTransactions = pagination.total_count || 0
          this.hasMore = pagination.has_more || false
          this.offset += newTransactions.length
        }
      } catch (error) {
        console.error('Failed to load transactions:', error)
        uni.showToast({
          title: '加载失败',
          icon: 'error',
        })
      } finally {
        this.loading = false
        this.loadingMore = false
      }
    },

    async loadMore() {
      if (!this.hasMore || this.loadingMore) return
      await this.loadTransactions(false)
    },

    setFilter(filterKey) {
      this.activeFilter = filterKey
    },

    getTransactionIcon(type) {
      switch (type) {
        case 'purchase':
          return '💰'
        case 'generation':
          return '🎨'
        default:
          return '📝'
      }
    },

    formatDate(dateString) {
      const date = new Date(dateString)
      const now = new Date()
      const diffTime = Math.abs(now - date)
      const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))

      if (diffDays === 1) {
        return (
          '今天 ' +
          date.toLocaleTimeString('zh-CN', {
            hour: '2-digit',
            minute: '2-digit',
          })
        )
      } else if (diffDays === 2) {
        return (
          '昨天 ' +
          date.toLocaleTimeString('zh-CN', {
            hour: '2-digit',
            minute: '2-digit',
          })
        )
      } else if (diffDays <= 7) {
        return `${diffDays - 1}天前`
      } else {
        return date.toLocaleDateString('zh-CN', {
          month: '2-digit',
          day: '2-digit',
          hour: '2-digit',
          minute: '2-digit',
        })
      }
    },

    goBack() {
      uni.navigateBack()
    },
  },
}
</script>

<style lang="scss" scoped>
.history-container {
  min-height: 100vh;
  background: var(--color-bg);
  padding-top: calc(env(safe-area-inset-top, 0px) + 20px);
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

.summary-section {
  padding: 32rpx 40rpx;

  .summary-card {
    background: var(--color-surface);
    border-radius: 32rpx;
    padding: 48rpx;
    display: flex;
    align-items: center;
    box-shadow: 0 8rpx 48rpx rgba(74, 74, 74, 0.08);

    .summary-item {
      flex: 1;
      display: flex;
      flex-direction: column;
      align-items: center;

      .summary-label {
        font-size: 28rpx;
        color: var(--color-text-subtle);
        margin-bottom: 16rpx;
      }

      .summary-value {
        font-size: 48rpx;
        font-weight: 600;
        color: var(--color-text);
      }
    }

    .summary-divider {
      width: 1rpx;
      height: 80rpx;
      background: #f0f0f0;
      margin: 0 48rpx;
    }
  }
}

.filter-section {
  padding: 0 40rpx 32rpx;

  .filter-tabs {
    display: flex;
    background: var(--color-surface);
    border-radius: 24rpx;
    padding: 8rpx;
    box-shadow: 0 4rpx 24rpx rgba(74, 74, 74, 0.06);

    .filter-tab {
      flex: 1;
      height: 64rpx;
      display: flex;
      align-items: center;
      justify-content: center;
      border-radius: 16rpx;
      transition: all 0.3s cubic-bezier(0.6, 0.05, 0.4, 1);

      &.active {
        background: var(--color-primary);

        .filter-text {
          color: white;
          font-weight: 600;
        }
      }

      .filter-text {
        font-size: 28rpx;
        color: var(--color-text-subtle);
        transition: all 0.3s cubic-bezier(0.6, 0.05, 0.4, 1);
      }
    }
  }
}

.transactions-section {
  padding: 0 40rpx 40rpx;

  .loading-state,
  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 120rpx 40rpx;

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
      color: var(--color-text-subtle);
    }

    .empty-icon {
      font-size: 96rpx;
      margin-bottom: 32rpx;
    }

    .empty-title {
      font-size: 36rpx;
      font-weight: 600;
      color: var(--color-text);
      margin-bottom: 16rpx;
    }

    .empty-description {
      font-size: 28rpx;
      color: var(--color-text-subtle);
    }
  }

  .transactions-list {
    .transaction-item {
      background: var(--color-surface);
      border-radius: 24rpx;
      padding: 32rpx;
      margin-bottom: 16rpx;
      display: flex;
      align-items: center;
      box-shadow: 0 4rpx 24rpx rgba(74, 74, 74, 0.06);
      transition: all 0.3s cubic-bezier(0.6, 0.05, 0.4, 1);

      &:active {
        transform: scale(0.98);
      }

      .transaction-icon {
        width: 80rpx;
        height: 80rpx;
        border-radius: 20rpx;
        display: flex;
        align-items: center;
        justify-content: center;
        margin-right: 24rpx;

        .icon-text {
          font-size: 32rpx;
        }
      }

      &.purchase .transaction-icon {
        background: rgba(76, 175, 80, 0.1);
      }

      &.generation .transaction-icon {
        background: rgba(232, 155, 147, 0.1);
      }

      .transaction-details {
        flex: 1;

        .transaction-description {
          font-size: 32rpx;
          font-weight: 500;
          color: var(--color-text);
          margin-bottom: 8rpx;
          display: block;
        }

        .transaction-date {
          font-size: 24rpx;
          color: var(--color-text-subtle);
          display: block;
          margin-bottom: 4rpx;
        }

        .transaction-id {
          display: flex;
          align-items: center;
          gap: 8rpx;

          .id-label {
            font-size: 20rpx;
            color: var(--color-text-subtle);
          }

          .id-value {
            font-size: 20rpx;
            color: var(--color-text-subtle);
            font-family: monospace;
          }
        }
      }

      .transaction-amount {
        display: flex;
        flex-direction: column;
        align-items: flex-end;

        .amount-text {
          font-size: 32rpx;
          font-weight: 600;
          line-height: 1;

          &.positive {
            color: #4caf50;
          }

          &.negative {
            color: var(--color-primary);
          }
        }

        .amount-unit {
          font-size: 20rpx;
          color: var(--color-text-subtle);
          margin-top: 4rpx;
        }
      }
    }

    .load-more {
      display: flex;
      align-items: center;
      justify-content: center;
      padding: 32rpx;
      margin-top: 16rpx;

      .loading-more {
        display: flex;
        align-items: center;
        gap: 16rpx;

        .loading-spinner.small {
          width: 32rpx;
          height: 32rpx;
          border-width: 2rpx;
        }

        .loading-text {
          font-size: 28rpx;
          color: var(--color-text-subtle);
        }
      }

      .load-more-text {
        font-size: 28rpx;
        color: var(--color-primary);
        font-weight: 500;
      }
    }

    .end-message {
      display: flex;
      align-items: center;
      justify-content: center;
      padding: 32rpx;

      .end-text {
        font-size: 24rpx;
        color: var(--color-text-subtle);
      }
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
