<template>
  <view class="generate-container">
    <LoadingAnimation v-if="loading" />
    <view v-else class="result-grid">
      <view v-for="(image, index) in images" :key="index" class="image-container">
        <image
          :src="image"
          class="result-image"
          mode="aspectFill"
        />
        <button class="save-button" @tap="saveImage(image)">Save</button>
      </view>
    </view>
  </view>
</template>

<script>
import { defineComponent, ref, onMounted } from 'vue';
import LoadingAnimation from '../../components/LoadingAnimation/LoadingAnimation.vue';

export default defineComponent({
  components: {
    LoadingAnimation,
  },
  props: {
    result: {
      type: String,
      required: true,
    },
  },
  setup(props) {
    const images = ref([]);
    const loading = ref(true);

    onMounted(() => {
      const resultData = JSON.parse(decodeURIComponent(props.result));
      images.value = resultData.images;
      loading.value = false;
    });

    const saveImage = (imageUrl) => {
      uni.saveImageToPhotosAlbum({
        filePath: imageUrl,
        success: () => {
          uni.showToast({
            title: 'Saved successfully',
            icon: 'success',
          });
        },
        fail: (err) => {
          console.error('Failed to save image:', err);
          uni.showToast({
            title: 'Failed to save',
            icon: 'none',
          });
        },
      });
    };

    return {
      images,
      loading,
      saveImage,
    };
  },
});
</script>

<style scoped>
.generate-container {
  padding: 20px;
  background-color: #fcfbf9;
  min-height: 100vh;
}

.result-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 10px;
}

.result-image {
  width: 100%;
  height: 100%;
  border-radius: 8px;
}

.image-container {
  position: relative;
}

.save-button {
  position: absolute;
  bottom: 10px;
  right: 10px;
  background-color: rgba(0, 0, 0, 0.5);
  color: white;
  border: none;
  border-radius: 5px;
  padding: 5px 10px;
  cursor: pointer;
}
</style> 