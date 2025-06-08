import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import { api } from '../api';

export const useAuthStore = defineStore('auth', () => {
  const token = ref(uni.getStorageSync('token') || null);
  const user = ref(null);

  const isAuthenticated = computed(() => !!token.value);

  const login = async (code) => {
    try {
      const response = await api.post('/auth/login', { code });
      token.value = response.data.token;
      user.value = response.data.user;
      uni.setStorageSync('token', token.value);
      uni.switchTab({ url: '/pages/index/index' });
    } catch (error) {
      console.error('Login failed:', error);
      uni.showToast({
        title: 'Login failed',
        icon: 'none',
      });
    }
  };

  const logout = () => {
    token.value = null;
    user.value = null;
    uni.removeStorageSync('token');
    uni.reLaunch({ url: '/pages/login/index' });
  };

  const fetchUser = async () => {
    if (token.value && !user.value) {
      try {
        const response = await api.get('/me');
        user.value = response.data;
      } catch (error) {
        console.error('Failed to fetch user:', error);
        // If the token is invalid, log the user out
        if (error.response && error.response.status === 401) {
          logout();
        }
      }
    }
  };

  // Fetch user on store initialization
  fetchUser();

  return {
    token,
    user,
    isAuthenticated,
    login,
    logout,
    fetchUser,
  };
}); 