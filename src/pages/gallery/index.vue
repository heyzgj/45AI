<route lang="json5" type="home">
{
  style: {
    navigationBarTitleText: '45AI Gallery',
  },
}
</route>

<template>
  <view class="gallery-container page-wrapper">
    <!-- Hero Carousel Section -->
    <view class="hero-section">
      <swiper
        class="hero-carousel"
        :indicator-dots="true"
        :autoplay="true"
        :circular="true"
        :interval="4000"
        indicator-color="rgba(255, 255, 255, 0.3)"
        :indicator-active-color="colorPrimary"
      >
        <swiper-item v-for="(slide, index) in heroSlides" :key="index" class="hero-slide">
          <view class="slide-content">
            <image :src="slide.image" mode="aspectFill" class="slide-image" />
          </view>
        </swiper-item>
      </swiper>
    </view>

    <!-- Main Content Container with Rounded Corners -->
    <view class="main-content-container">
      <!-- AI Function Tiles Section -->
      <view class="ai-functions-section">
        <view class="function-tiles">
          <view
            v-for="(func, index) in aiFunctions"
            :key="index"
            class="function-tile"
            :style="{ animationDelay: `${index * 100}ms` }"
            @click="handleFunctionClick(func)"
          >
            <view class="tile-icon">{{ func.icon }}</view>
            <text class="tile-label">{{ func.name }}</text>
          </view>
        </view>
      </view>

      <!-- Coming Soon Popup -->
      <wd-popup
        v-model="showComingSoonDialog"
        position="center"
        :close-on-click-modal="true"
        custom-style="border-radius: 16px; padding: 0;"
      >
        <view class="coming-soon-popup">
          <view class="coming-soon-header">
            <text class="coming-soon-popup-title">功能预告</text>
          </view>
          <view class="coming-soon-content">
            <view class="coming-soon-icon">{{ comingSoonFeature.icon }}</view>
            <text class="coming-soon-title">{{ comingSoonFeature.name }}</text>
            <text class="coming-soon-message">{{ comingSoonMessage }}</text>
            <text class="coming-soon-subtitle">我们正在努力开发中，敬请期待！</text>
          </view>
          <view class="coming-soon-footer">
            <wd-button @click="showComingSoonDialog = false" type="primary" size="large" block>
              好的
            </wd-button>
          </view>
        </view>
      </wd-popup>

      <!-- New Search & Filter Section -->
      <view class="search-filter-section">
        <!-- Row 1: Tab Switcher + Search -->
        <view class="search-filter-row-1">
          <view class="tab-switcher">
            <view
              v-for="tab in mainTabs"
              :key="tab.id"
              class="main-tab"
              :class="{ active: activeMainTab === tab.id }"
              @click="handleMainTabSelect(tab.id)"
            >
              <text class="tab-text">{{ tab.name }}</text>
            </view>
          </view>

          <view class="search-pill">
            <text class="search-icon">🔍</text>
            <input
              class="search-input"
              type="text"
              placeholder="搜索模板"
              v-model="searchQuery"
              @input="handleSearch"
            />
          </view>
        </view>

        <!-- Row 2: Category Chips -->
        <view class="search-filter-row-2">
          <scroll-view class="category-scroll" scroll-x enable-flex>
            <view class="category-chips">
              <view
                v-for="chip in categoryChips"
                :key="chip.id"
                class="category-chip"
                :class="{ active: activeCategory === chip.id }"
                @click="handleCategorySelect(chip.id)"
              >
                <text class="chip-text">{{ chip.name }}</text>
              </view>
            </view>
          </scroll-view>
        </view>
      </view>

      <!-- Templates Section -->
      <view class="templates-section">
        <view class="template-grid">
          <TemplateCard
            v-for="(template, index) in filteredTemplates"
            :key="template.id"
            :template="template"
            class="template-item animate-slide-up"
            :style="{ animationDelay: `${index * 50}ms` }"
          />
        </view>
      </view>

      <!-- Loading State -->
      <view v-if="loading" class="loading-overlay">
        <LoadingAnimation variant="blob" size="large" text="Loading beautiful styles..." />
      </view>

      <!-- Empty State -->
      <view v-if="!loading && filteredTemplates.length === 0" class="empty-state">
        <text class="empty-icon">🎨</text>
        <text class="empty-title">No Templates Found</text>
        <text class="empty-message">Try adjusting your search or filters</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useTemplateStore } from '@/store/templates'
import TemplateCard from '@/components/TemplateCard/TemplateCard.vue'
import LoadingAnimation from '@/components/LoadingAnimation/LoadingAnimation.vue'

// Color constants for carousel
const colorPrimary = '#E89B93'

// Hero carousel data
const heroSlides = ref([
  {
    title: '玩转日常穿搭',
    subtitle: '一键换装，穿越童年动画世界',
    image: 'https://picsum.photos/800/400?random=hero1',
  },
  {
    title: '夏日海边拍贴',
    subtitle: 'AI魔法让你的照片更出彩',
    image: 'https://picsum.photos/800/400?random=hero2',
  },
  {
    title: '创意无限',
    subtitle: '探索更多AI创作可能',
    image: 'https://picsum.photos/800/400?random=hero3',
  },
])

// AI Features coming soon
const aiFunctions = ref([
  {
    name: '生成',
    icon: '📸',
    comingSoon: true,
  },
  {
    name: '模板',
    icon: '👗',
    comingSoon: true,
  },
  {
    name: '我的',
    icon: '🌟',
    comingSoon: true,
  },
])

// Main tabs (row 1)
const mainTabs = ref([
  { id: 'latest', name: '最新', highlight: false },
  { id: 'recommended', name: '推荐', highlight: true },
])

// Category chips (row 2)
const categoryChips = ref([
  { id: 'all', name: '全部' },
  { id: 'christmas', name: '圣诞写真' },
  { id: 'newyear', name: '新年' },
  { id: 'male', name: '男生' },
  { id: 'retro', name: '时代映像' },
])

// State
const searchQuery = ref('')
const activeMainTab = ref('recommended')
const activeCategory = ref('all')
const loading = ref(false)
const showComingSoonDialog = ref(false)
const comingSoonMessage = ref('')
const comingSoonFeature = ref({ name: '', icon: '' })

// Mock template data
const templates = ref([
  {
    id: 1,
    name: 'Dusty Rose Dream',
    description: 'Soft, romantic vibes',
    preview_image_url: 'https://picsum.photos/400/600?random=1',
    credit_cost: 15,
    is_active: true,
    category: 'aesthetic',
    created_at: '2024-01-01T00:00:00Z',
  },
  {
    id: 2,
    name: 'Fairy Tale Magic',
    description: 'Enchanted forest aesthetic',
    preview_image_url: 'https://picsum.photos/400/600?random=2',
    credit_cost: 20,
    is_active: true,
    category: 'dreamy',
    created_at: '2024-01-01T00:00:00Z',
  },
  {
    id: 3,
    name: 'Vintage Film',
    description: 'Classic cinema look',
    preview_image_url: 'https://picsum.photos/400/600?random=3',
    credit_cost: 10,
    is_active: true,
    category: 'vintage',
    created_at: '2024-01-01T00:00:00Z',
  },
  {
    id: 4,
    name: 'Modern Minimalist',
    description: 'Clean and sophisticated',
    preview_image_url: 'https://picsum.photos/400/600?random=4',
    credit_cost: 15,
    is_active: true,
    category: 'modern',
    created_at: '2024-01-01T00:00:00Z',
  },
  {
    id: 5,
    name: 'Golden Hour',
    description: 'Warm sunset glow',
    preview_image_url: 'https://picsum.photos/400/600?random=5',
    credit_cost: 12,
    is_active: true,
    category: 'portrait',
    created_at: '2024-01-01T00:00:00Z',
  },
  {
    id: 6,
    name: 'Artistic Portrait',
    description: 'Creative expression',
    preview_image_url: 'https://picsum.photos/400/600?random=6',
    credit_cost: 25,
    is_active: true,
    category: 'artistic',
    created_at: '2024-01-01T00:00:00Z',
  },
])

// Computed
const filteredTemplates = computed(() => {
  let filtered = templates.value

  // Apply category filter
  if (activeCategory.value !== 'all') {
    filtered = filtered.filter((t) => t.category === activeCategory.value)
  }

  // Apply search filter
  if (searchQuery.value.trim()) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(
      (t) => t.name.toLowerCase().includes(query) || t.description.toLowerCase().includes(query),
    )
  }

  return filtered
})

// Methods
const handleSearch = () => {
  // Search is reactive through computed property
}

const handleFilterSelect = (filterId: string) => {
  activeCategory.value = filterId
}

const handleMainTabSelect = (tabId: string) => {
  activeMainTab.value = tabId
}

const handleCategorySelect = (chipId: string) => {
  activeCategory.value = chipId
}

const handleFunctionClick = (func: any) => {
  if (func.comingSoon) {
    comingSoonFeature.value = func
    comingSoonMessage.value = `${func.name} 即将上线`
    showComingSoonDialog.value = true
  }
}

onMounted(() => {
  // TODO: Load templates from API
})
</script>

<style lang="scss" scoped>
@import '@/style/variables.scss';

.gallery-container {
  background: linear-gradient(180deg, #fafafa 0%, #ffffff 100%);
  min-height: 100vh;
}

// Main Content Container with Rounded Top Corners
.main-content-container {
  background-color: $color-surface;
  border-radius: $radius-xl $radius-xl 0 0;
  margin-top: -$spacing-lg;
  padding-top: $spacing-xl;
  box-shadow: 0 -4px 24px rgba(0, 0, 0, 0.08);
  position: relative;
  z-index: 2;

  // Add subtle gradient border like detail page
  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 2px;
    background: linear-gradient(
      90deg,
      transparent,
      rgba($color-primary, 0.3),
      rgba($color-secondary, 0.4),
      rgba($color-primary, 0.3),
      transparent
    );
    border-radius: $radius-xl $radius-xl 0 0;
  }
}

// Hero Carousel
.hero-section {
  width: 100%;
  aspect-ratio: 16/9;
  margin-bottom: $spacing-md;
  position: relative;
  z-index: 1;
}

.hero-carousel {
  width: 100%;
  height: 240px;
  border-radius: 0;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}

.carousel-item {
  position: relative;
  width: 100%;
  height: 100%;
  aspect-ratio: 16/9;
  border-radius: 16px;
  overflow: hidden;
  margin: 0 8px;
}

.hero-slide {
  position: relative;
  width: 100%;
  height: 100%;
}

.slide-content {
  position: relative;
  width: 100%;
  height: 100%;
}

.slide-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

// AI Function Tiles Section
.ai-functions-section {
  margin-bottom: $spacing-lg;
  padding: 0 $page-padding;
}

.section-title {
  display: block;
  font-size: $font-size-h2;
  font-weight: $font-weight-semibold;
  color: $color-text;
  padding: 0 $page-padding $spacing-md;
}

.function-tiles {
  display: flex;
  gap: $spacing-md;
  padding-right: $page-padding;
}

.function-tile {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  aspect-ratio: 1;
  background: linear-gradient(135deg, #ffffff 0%, #f8f9fa 100%);
  border-radius: 16px;
  border: 1px solid rgba(0, 0, 0, 0.06);
  transition: all $duration-fast $ease-custom;

  &:active {
    transform: scale(0.98);
    background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.12);
  }
}

.tile-icon {
  position: relative;
  margin-bottom: $spacing-sm;
}

.tile-label {
  display: block;
  font-size: $font-size-body;
  font-weight: $font-weight-medium;
  color: $color-text;
  margin-bottom: 2px;
}

// Search & Filter Section
.search-filter-section {
  padding: 0 $page-padding $spacing-lg;
}

.search-filter-row-1 {
  display: flex;
  gap: $spacing-sm;
  align-items: center;
  margin-bottom: $spacing-sm;
  justify-content: space-between;
}

.tab-switcher {
  display: flex;
  gap: $spacing-xs;
}

.main-tab {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: $spacing-xs $spacing-sm;
  border-radius: $radius-full;
  background: $color-surface;
  border: 1px solid $border-color;
  transition: all $duration-fast $ease-custom;

  &.active {
    background: $color-primary;
    border-color: $color-primary;

    .tab-text {
      color: white;
    }
  }

  &:active {
    transform: scale(0.96);
  }
}

.tab-text {
  font-size: $font-size-caption;
  font-weight: $font-weight-medium;
  color: $color-text;
}

.tab-highlight {
  font-size: 12px;
}

.search-pill {
  display: flex;
  align-items: center;
  gap: $spacing-xs;
  background: $color-surface;
  border-radius: $radius-full;
  padding: $spacing-xs $spacing-sm;
  border: 1px solid $border-color;
  min-width: 120px;
}

.search-icon {
  font-size: 14px;
  opacity: 0.6;
}

.search-input {
  flex: 1;
  border: none;
  outline: none;
  background: transparent;
  font-size: $font-size-caption;
  color: $color-text;

  &::placeholder {
    color: $color-text-subtle;
  }
}

.search-filter-row-2 {
  overflow: hidden;
}

.category-scroll {
  scrollbar-width: none;
  scrollbar-color: transparent transparent;
  white-space: nowrap;

  &::-webkit-scrollbar {
    display: none;
  }
}

.category-chips {
  display: flex;
  gap: $spacing-xs;
  padding-right: $spacing-lg;
}

.category-chip {
  flex-shrink: 0;
  padding: $spacing-xs $spacing-sm;
  border-radius: $radius-full;
  background: $color-surface;
  border: 1px solid $border-color;
  transition: all $duration-fast $ease-custom;

  &.active {
    background: $color-secondary;
    border-color: $color-secondary;

    .chip-text {
      color: $color-text;
    }
  }

  &:active {
    transform: scale(0.96);
  }
}

.chip-text {
  font-size: $font-size-caption;
  font-weight: $font-weight-medium;
  color: $color-text-subtle;
}

.filter-scroll {
  flex: 1;

  /* Hide scrollbar for all browsers and screen sizes */
  &::-webkit-scrollbar {
    display: none;
    width: 0;
    height: 0;
    background: transparent;
  }

  /* For Firefox */
  scrollbar-width: none;
  scrollbar-color: transparent transparent;
}

.filter-tags {
  scrollbar-width: none;
  scrollbar-color: transparent transparent;
  display: flex;
  gap: $spacing-xs;
  padding-right: $spacing-md;

  /* Additional scrollbar hiding */
  &::-webkit-scrollbar {
    display: none;
    width: 0;
    height: 0;
    background: transparent;
  }
}

.search-input-wrapper {
  display: flex !important;
  align-items: center;
  background: $color-surface;
  border-radius: $radius-full;
  padding: $spacing-xs $spacing-sm;
  box-shadow: $shadow-card;
  min-width: 120px;
  width: 120px;
  flex-shrink: 0;
  position: relative;
  z-index: 10;
  visibility: visible !important;
  opacity: 1 !important;

  // Responsive: Force visibility on all screen sizes
  @media (max-width: 767px) {
    display: flex !important;
    min-width: 110px;
    width: 110px;
    flex-shrink: 0;
    visibility: visible !important;
    opacity: 1 !important;
  }

  @media (max-width: 480px) {
    min-width: 100px;
    width: 100px;
  }
}

.search-icon {
  font-size: 14px;
  color: $color-text-subtle;
  margin-right: $spacing-xs;
}

.search-input {
  flex: 1;
  border: none;
  outline: none;
  background: transparent;
  font-size: $font-size-caption;
  color: $color-text;
  width: 100%;

  &::placeholder {
    color: $color-text-subtle;
  }
}

.filter-tag {
  flex-shrink: 0;
  padding: 6px 12px;
  background: $color-surface;
  border-radius: $radius-full;
  border: 1px solid rgba($color-primary, 0.2);
  transition: all $duration-fast $ease-custom;

  &.active {
    background: $color-primary;
    border-color: $color-primary;

    .tag-text {
      color: white;
      font-weight: $font-weight-medium;
    }
  }

  &:active {
    transform: scale(0.95);
  }
}

.tag-text {
  font-size: 11px;
  color: $color-text;
  transition: color $duration-fast $ease-custom;
  white-space: nowrap;
}

// Templates Section
.templates-section {
  padding: 0 $page-padding $spacing-xl;
}

.section-header {
  margin-bottom: $spacing-lg;
}

.section-subtitle {
  display: block;
  font-size: $font-size-body;
  color: $color-text-subtle;
  margin-top: 4px;
}

.template-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: $spacing-lg;
}

.template-item {
  animation-fill-mode: both;
}

// Loading & Empty States
.loading-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba($color-bg, 0.9);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 999;
}

.empty-state {
  text-align: center;
  padding: $spacing-xl;

  .empty-icon {
    font-size: 48px;
    display: block;
    margin-bottom: $spacing-md;
  }

  .empty-title {
    display: block;
    font-size: $font-size-h2;
    font-weight: $font-weight-semibold;
    color: $color-text;
    margin-bottom: $spacing-sm;
  }

  .empty-message {
    display: block;
    font-size: $font-size-body;
    color: $color-text-subtle;
  }
}

// Responsive adjustments
@media (min-width: 768px) {
  .template-grid {
    grid-template-columns: repeat(3, 1fr);
    gap: $spacing-lg;
  }

  .function-tile {
    width: 140px;
  }
}

// Coming Soon Popup
.coming-soon-popup {
  background: white;
  border-radius: 16px;
  overflow: hidden;
  max-width: 300px;
  width: 90vw;
  margin: 0 auto;
}

.coming-soon-header {
  background: linear-gradient(135deg, $color-primary, $color-secondary);
  padding: $spacing-md $spacing-lg;
  text-align: center;

  .coming-soon-popup-title {
    font-size: $font-size-h2;
    font-weight: $font-weight-semibold;
    color: white;
    display: block;
  }
}

.coming-soon-content {
  text-align: center;
  padding: $spacing-lg;

  .coming-soon-icon {
    font-size: 48px;
    margin-bottom: $spacing-md;
    display: block;
  }

  .coming-soon-title {
    font-size: $font-size-h2;
    font-weight: $font-weight-semibold;
    color: $color-text;
    margin-bottom: $spacing-sm;
    display: block;
  }

  .coming-soon-message {
    font-size: $font-size-body;
    color: $color-primary;
    font-weight: $font-weight-medium;
    margin-bottom: $spacing-sm;
    display: block;
  }

  .coming-soon-subtitle {
    font-size: $font-size-caption;
    color: $color-text-subtle;
    display: block;
    line-height: 1.5;
  }
}

.coming-soon-footer {
  padding: 0 $spacing-lg $spacing-lg;
}

// Animations
@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.animate-slide-up {
  animation: slideUp 0.6s $ease-custom forwards;
  opacity: 0;
}
</style>
