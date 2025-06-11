<template>
  <view class="profile-container page-wrapper animate-fade-in">
    <!-- User Info Section -->
    <view class="user-section card animate-slide-up">
      <view class="user-header">
        <image
          :src="userInfo.avatar_url || '/static/default-avatar.png'"
          class="user-avatar"
          mode="aspectFill"
        />
        <view class="user-details">
          <text class="user-name">{{ userInfo.nickname || 'User' }}</text>
          <text class="user-id">ID: {{ userInfo.id }}</text>
        </view>
      </view>

      <!-- Credit Balance -->
      <view class="credit-section">
        <view class="credit-balance">
          <text class="credit-label">Your Balance</text>
          <CreditDisplay size="large" :showLabel="false" :clickable="false" />
        </view>
        <button class="btn-primary gradient" @click="navigateToPurchase">Add Credits</button>
      </view>
    </view>

    <!-- Quick Actions -->
    <view class="actions-section">
      <view class="action-grid">
        <view
          class="action-card card interactive animate-slide-up delay-100"
          @click="navigateToHistory"
        >
          <text class="action-icon">📜</text>
          <text class="action-label">Transaction History</text>
        </view>
        <view
          class="action-card card interactive animate-slide-up delay-200"
          @click="navigateToSettings"
        >
          <text class="action-icon">⚙️</text>
          <text class="action-label">Settings</text>
        </view>
      </view>
    </view>

    <!-- Recent Activity -->
    <view class="activity-section animate-slide-up delay-300">
      <text class="section-title">Recent Activity</text>
      <view class="activity-list">
        <view
          v-for="(activity, index) in recentActivities"
          :key="activity.id"
          class="activity-item"
          :class="{ 'animate-slide-up': true }"
          :style="{ animationDelay: `${300 + index * 50}ms` }"
        >
          <view class="activity-info">
            <text class="activity-type">{{ getActivityTypeLabel(activity.type) }}</text>
            <text class="activity-desc">{{ activity.description }}</text>
            <text class="activity-time">{{ formatTime(activity.created_at) }}</text>
          </view>
          <text
            class="activity-amount"
            :class="{ positive: activity.amount > 0, negative: activity.amount < 0 }"
          >
            {{ activity.amount > 0 ? '+' : '' }}{{ activity.amount }}
          </text>
        </view>
      </view>

      <!-- Empty State -->
      <view v-if="recentActivities.length === 0" class="empty-activity">
        <text class="empty-text">No recent activity</text>
      </view>
    </view>

    <!-- Logout Button -->
    <view class="logout-section">
      <button class="btn-secondary" @click="handleLogout">Log Out</button>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useUserStore } from '@/store/user'
import CreditDisplay from '@/components/CreditDisplay/CreditDisplay.vue'
import * as transactionsApi from '@/api/transactions'

// Use actual user store instead of mock data
const userStore = useUserStore()

// Computed properties to get real user data
const userInfo = computed(
  () =>
    userStore.userInfo || {
      id: '',
      nickname: 'Guest User',
      avatar_url: '',
      credits: 0,
    },
)

// Real user activities (fetched from API)
const recentActivities = ref([])
const loadingActivities = ref(false)

// Fetch recent activities
const fetchRecentActivities = async () => {
  if (!userStore.isAuthenticated) return

  loadingActivities.value = true
  try {
    const response = await transactionsApi.getTransactions({
      limit: 5,
      offset: 0,
    })

    // Handle response structure
    const transactions = response.data?.transactions || []
    recentActivities.value = transactions.map((transaction: any) => ({
      id: transaction.id,
      type: transaction.type,
      description: transaction.description || `${transaction.type} transaction`,
      amount: transaction.amount,
      created_at: new Date(transaction.created_at),
    }))
  } catch (error) {
    console.error('Failed to fetch recent activities:', error)
    // Don't show error to user, just keep empty state
  } finally {
    loadingActivities.value = false
  }
}

const navigateToPurchase = () => {
  uni.navigateTo({
    url: '/pages/purchase/index',
  })
}

const navigateToHistory = () => {
  uni.navigateTo({
    url: '/pages/history/index',
  })
}

const navigateToSettings = () => {
  // TODO: Implement settings page
  uni.showToast({
    title: 'Coming soon',
    icon: 'none',
  })
}

const handleLogout = () => {
  uni.showModal({
    title: 'Confirm Logout',
    content: 'Are you sure you want to log out?',
    confirmText: 'Log Out',
    confirmColor: '#E89B93',
    success: (res) => {
      if (res.confirm) {
        // Use user store logout method
        userStore.logout()
        uni.redirectTo({
          url: '/pages/login/index',
        })
      }
    },
  })
}

const getActivityTypeLabel = (type: string) => {
  return type === 'purchase' ? '💰 Purchase' : '🎨 Generation'
}

const formatTime = (date: Date) => {
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))

  if (days === 0) return 'Today'
  if (days === 1) return 'Yesterday'
  return `${days} days ago`
}

onMounted(async () => {
  // Initialize user store and fetch latest user data
  if (!userStore.isInitialized) {
    userStore.initFromStorage()
  }

  // Fetch fresh user data if authenticated
  if (userStore.isAuthenticated) {
    try {
      await userStore.fetchUserInfo()
      // Fetch recent activities after user info is loaded
      await fetchRecentActivities()
    } catch (error) {
      console.error('Failed to fetch user info:', error)
    }
  }

  console.log('Profile page - user info:', userStore.userInfo)
  console.log('Profile page - user credits:', userStore.credits)
})
</script>

<style lang="scss" scoped>
@import '@/style/variables.scss';

.profile-container {
  background-color: $color-bg;
  min-height: 100vh;
  padding: calc(env(safe-area-inset-top, 0px) + 20px) 0 $spacing-lg;
}

.user-section {
  margin: $spacing-md $page-padding;
}

.user-header {
  display: flex;
  align-items: center;
  gap: $spacing-md;
  margin-bottom: $spacing-lg;
}

.user-avatar {
  width: 64px;
  height: 64px;
  border-radius: $radius-full;
  border: 3px solid $color-secondary;
}

.user-details {
  flex: 1;
}

.user-name {
  display: block;
  font-size: $font-size-h2;
  font-weight: $font-weight-semibold;
  color: $color-text;
  margin-bottom: 4px;
  max-width: 200px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.user-id {
  display: block;
  font-size: $font-size-caption;
  color: $color-text-subtle;
}

.credit-section {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-top: $spacing-lg;
  border-top: 1px solid $border-color;
}

.credit-balance {
  flex: 1;
}

.credit-label {
  display: block;
  font-size: $font-size-caption;
  color: $color-text-subtle;
  margin-bottom: 4px;
}

.actions-section {
  margin: $spacing-md $page-padding $spacing-sm;
}

.section-title {
  display: block;
  font-size: $font-size-h2;
  font-weight: $font-weight-semibold;
  color: $color-text;
  margin-bottom: $spacing-sm;
}

.action-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: $spacing-sm;
}

.action-card {
  padding: $spacing-md;
  text-align: center;

  &:active {
    transform: translateY(-2px) scale(1.01);
  }
}

.action-icon {
  display: block;
  font-size: 28px;
  margin-bottom: $spacing-xs;
}

.action-label {
  display: block;
  font-size: $font-size-body;
  color: $color-text;
}

.activity-section {
  margin: $spacing-sm $page-padding $spacing-lg;
}

.activity-list {
  margin-top: $spacing-md;
}

.activity-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: $spacing-md 0;
  border-bottom: 1px solid $border-color;

  &:last-child {
    border-bottom: none;
  }
}

.activity-info {
  flex: 1;
}

.activity-type {
  display: block;
  font-size: $font-size-body;
  color: $color-text;
  margin-bottom: 2px;
}

.activity-desc {
  display: block;
  font-size: $font-size-caption;
  color: $color-text-subtle;
  margin-bottom: 2px;
}

.activity-time {
  display: block;
  font-size: $font-size-caption;
  color: $color-text-subtle;
}

.activity-amount {
  font-size: $font-size-body;
  font-weight: $font-weight-semibold;

  &.positive {
    color: $color-success;
  }

  &.negative {
    color: $color-primary;
  }
}

.empty-activity {
  text-align: center;
  padding: $spacing-xl 0;
}

.empty-text {
  color: $color-text-subtle;
  font-size: $font-size-body;
}

.logout-section {
  margin: $spacing-xl $page-padding 0;

  .btn-secondary {
    width: 100%;
  }
}
</style>
