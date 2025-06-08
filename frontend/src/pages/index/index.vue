<template>
  <view class="home-container animate-fade-in">
    <!-- Hero Section -->
    <view class="hero-section">
      <text class="hero-title animate-slide-up">
        Transform Your Photos
        with AI Magic ✨
      </text>
      <text class="hero-subtitle animate-slide-up" style="animation-delay: 100ms">
        One tap, endless possibilities
      </text>
      
      <!-- CTA Button -->
      <button 
        class="btn-primary hero-cta animate-slide-up" 
        style="animation-delay: 200ms"
        @click="navigateToGallery"
      >
        Start Creating
      </button>
    </view>
    
    <!-- Features Section -->
    <view class="features-section">
      <view class="feature-card card animate-slide-up" style="animation-delay: 300ms">
        <text class="feature-icon">🎨</text>
        <text class="feature-title">Premium Templates</text>
        <text class="feature-desc">Curated styles by professional artists</text>
      </view>
      
      <view class="feature-card card animate-slide-up" style="animation-delay: 400ms">
        <text class="feature-icon">⚡</text>
        <text class="feature-title">Instant Results</text>
        <text class="feature-desc">Get your photos in under 30 seconds</text>
      </view>
      
      <view class="feature-card card animate-slide-up" style="animation-delay: 500ms">
        <text class="feature-icon">🛡️</text>
        <text class="feature-title">Safe & Private</text>
        <text class="feature-desc">Your photos are protected and secure</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const navigateToGallery = () => {
  uni.switchTab({
    url: '/pages/gallery/gallery'
  })
}
</script>

<style lang="scss" scoped>
.home-container {
  min-height: 100vh;
  background-color: #FCFBF9;
  padding: 20px;
}

.hero-section {
  text-align: center;
  padding: 60px 0 40px;
}

.hero-title {
  font-size: 32px;
  font-weight: 600;
  color: #4A4A4A;
  line-height: 1.3;
  display: block;
  margin-bottom: 12px;
}

.hero-subtitle {
  font-size: 16px;
  color: #9B9B9B;
  display: block;
  margin-bottom: 32px;
}

.hero-cta {
  background: linear-gradient(135deg, #E89B93, #F3D9D7);
  color: white;
  font-size: 16px;
  font-weight: 500;
  padding: 14px 32px;
  border-radius: 24px;
  border: none;
  box-shadow: 0 4px 12px rgba(232, 155, 147, 0.3);
  transition: all 0.3s cubic-bezier(0.6, 0.05, 0.4, 1);
  
  &:active {
    transform: scale(0.97);
    box-shadow: 0 2px 8px rgba(232, 155, 147, 0.3);
  }
}

.features-section {
  padding: 20px 0;
}

.feature-card {
  text-align: center;
  padding: 24px;
  margin-bottom: 16px;
  background: #FFFFFF;
  border: 1px solid #F0F0F0;
  border-radius: 16px;
  box-shadow: 0 4px 24px rgba(74, 74, 74, 0.08);
  transition: all 0.3s cubic-bezier(0.6, 0.05, 0.4, 1);
  
  &:active {
    transform: translateY(-4px) scale(1.02);
    box-shadow: 0 8px 32px rgba(74, 74, 74, 0.12);
  }
}

.feature-icon {
  font-size: 36px;
  display: block;
  margin-bottom: 12px;
}

.feature-title {
  font-size: 18px;
  font-weight: 600;
  color: #4A4A4A;
  display: block;
  margin-bottom: 8px;
}

.feature-desc {
  font-size: 14px;
  color: #9B9B9B;
  line-height: 1.5;
  display: block;
}

/* Animations */
@keyframes fade-in {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes slide-up {
  from {
    transform: translateY(16px);
    opacity: 0;
  }
  to {
    transform: translateY(0);
    opacity: 1;
  }
}

.animate-fade-in {
  animation: fade-in 0.4s cubic-bezier(0.6, 0.05, 0.4, 1) forwards;
}

.animate-slide-up {
  animation: slide-up 0.4s cubic-bezier(0.6, 0.05, 0.4, 1) forwards;
  opacity: 0;
}
</style> 