<script setup lang="ts">
import { onLaunch, onShow, onHide } from '@dcloudio/uni-app'
import { setupNavigationInterceptor } from '@/utils/navigation'
import { useUserStore } from '@/store/user'

onLaunch(() => {
  console.log('App Launch')

  // Setup navigation interceptor for auth guards
  setupNavigationInterceptor()

  // Initialize user store from local storage
  const userStore = useUserStore()
  userStore.initFromStorage()
})

onShow(() => {
  console.log('App Show')
})

onHide(() => {
  console.log('App Hide')
})
</script>

<style>
/* Global styles */
page {
  background-color: #fcfbf9; /* --color-bg */
  font-family: 'PingFang SC', Inter, sans-serif;
  color: #4a4a4a; /* --color-text */
  line-height: 1.6;
}

/* Remove default styles */
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

/* Ensure smooth animations */
view,
text,
image {
  transition: all 0.3s cubic-bezier(0.6, 0.05, 0.4, 1);
}
</style>
