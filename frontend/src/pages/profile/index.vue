<template>
  <view class="profile-container" v-if="user">
    <view class="profile-header">
      <image class="avatar" :src="user.avatar_url" mode="aspectFill" />
      <text class="nickname">{{ user.nickname }}</text>
    </view>
    <CreditDisplay :credits="user.credits" />
    <view class="transactions">
      <text class="section-title">Recent Transactions</text>
      <view v-for="tx in transactions" :key="tx.id" class="transaction-item">
        <text>{{ tx.description }}</text>
        <text :class="tx.amount > 0 ? 'amount-positive' : 'amount-negative'">
          {{ tx.amount > 0 ? '+' : '' }}{{ tx.amount }}
        </text>
      </view>
    </view>
  </view>
  <LoadingAnimation v-else />
</template>

<script>
import { defineComponent, ref, onMounted } from 'vue';
import { getProfile, getTransactions } from '../../api/user';
import CreditDisplay from '../../components/CreditDisplay/CreditDisplay.vue';
import LoadingAnimation from '../../components/LoadingAnimation/LoadingAnimation.vue';

export default defineComponent({
  components: {
    CreditDisplay,
    LoadingAnimation,
  },
  setup() {
    const user = ref(null);
    const transactions = ref([]);

    onMounted(async () => {
      try {
        const profileRes = await getProfile();
        user.value = profileRes.data;

        const transactionsRes = await getTransactions();
        transactions.value = transactionsRes.data;
      } catch (error) {
        console.error('Failed to fetch profile data:', error);
      }
    });

    return {
      user,
      transactions,
    };
  },
});
</script>

<style scoped>
.profile-container {
  padding: 20px;
  background-color: #fcfbf9;
}

.profile-header {
  display: flex;
  align-items: center;
  margin-bottom: 20px;
}

.avatar {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  margin-right: 20px;
}

.nickname {
  font-size: 24px;
  font-weight: 600;
  color: #4a4a4a;
}

.transactions {
  margin-top: 30px;
}

.section-title {
  font-size: 18px;
  font-weight: 600;
  color: #4a4a4a;
  margin-bottom: 15px;
  display: block;
}

.transaction-item {
  display: flex;
  justify-content: space-between;
  padding: 15px;
  background-color: #fff;
  border-radius: 8px;
  margin-bottom: 10px;
  box-shadow: 0 2px 8px rgba(74, 74, 74, 0.05);
}

.amount-positive {
  color: #2ecc71;
}

.amount-negative {
  color: #e74c3c;
}
</style> 