<template>
  <view class="login-container">
    <button class="login-button" @tap="handleLogin">Login with WeChat</button>
  </view>
</template>

<script>
import { defineComponent } from 'vue'
import { useUserStore } from '../../store/user'

export default defineComponent({
  setup() {
    const userStore = useUserStore()

    const handleLogin = async () => {
      try {
        // Check if we're in development mode (H5)
        // #ifdef H5
        // For development, create a test user without using mock authentication
        // This simulates a proper WeChat login flow
        const testCode = `test_${Date.now()}_${Math.random().toString(36).substring(7)}`
        await userStore.login(testCode)
        // Navigate to gallery page after successful login
        uni.switchTab({ url: '/pages/gallery/index' })
        // #endif

        // #ifndef H5
        // Use real WeChat login for production
        uni.login({
          provider: 'weixin',
          success: async (res) => {
            try {
              await userStore.login(res.code)
              uni.switchTab({ url: '/pages/gallery/index' })
            } catch (error) {
              console.error('Login failed:', error)
              uni.showToast({
                title: 'Login failed',
                icon: 'none',
              })
            }
          },
          fail: (err) => {
            console.error('Login failed:', err)
            uni.showToast({
              title: 'Login failed',
              icon: 'none',
            })
          },
        })
        // #endif
      } catch (error) {
        console.error('Login failed:', error)
        uni.showToast({
          title: 'Login failed',
          icon: 'none',
        })
      }
    }

    return {
      handleLogin,
    }
  },
})
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
