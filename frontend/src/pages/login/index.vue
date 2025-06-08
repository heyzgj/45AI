<template>
  <view class="login-container">
    <button class="login-button" @tap="handleLogin">Login with WeChat</button>
  </view>
</template>

<script>
import { defineComponent } from 'vue';
import { useAuthStore } from '../../stores/auth';

export default defineComponent({
  setup() {
    const authStore = useAuthStore();

    const handleLogin = () => {
      uni.login({
        provider: 'weixin',
        success: (res) => {
          authStore.login(res.code);
        },
        fail: (err) => {
          console.error('Login failed:', err);
        },
      });
    };

    return {
      handleLogin,
    };
  },
});
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background-color: #fcfbf9; /* --color-bg */
}

.login-button {
  background-color: #e89b93; /* --color-primary */
  color: #fff;
  border: none;
  padding: 12px 24px;
  border-radius: 24px;
  font-size: 16px;
  cursor: pointer;
  box-shadow: 0 4px 12px rgba(232, 155, 147, 0.3);
  transition: all 0.2s ease-in-out;
}

.login-button:active {
  transform: scale(0.97);
  box-shadow: 0 2px 6px rgba(232, 155, 147, 0.3);
}
</style> 