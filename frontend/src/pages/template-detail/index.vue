<template>
  <view class="detail-container" v-if="template">
    <image class="preview-image" :src="template.preview_image_url" mode="widthFix" />
    <view class="info-card">
      <text class="title">{{ template.name }}</text>
      <text class="description">{{ template.description }}</text>
      <view class="cost">
        <text>Credit Cost: {{ template.credit_cost }}</text>
      </view>
      <button class="generate-button" @tap="handleGenerate" :disabled="loading">
        {{ loading ? 'Generating...' : 'Generate with this Template' }}
      </button>
    </view>
  </view>
  <LoadingAnimation v-else />
</template>

<script>
import { defineComponent, ref, onMounted } from 'vue';
import { getTemplateByID } from '../../api/template';
import { generateImage } from '../../api/generation';
import LoadingAnimation from '../../components/LoadingAnimation/LoadingAnimation.vue';

export default defineComponent({
  components: {
    LoadingAnimation,
  },
  props: {
    id: {
      type: String,
      required: true,
    },
  },
  setup(props) {
    const template = ref(null);
    const loading = ref(false);

    onMounted(async () => {
      try {
        const response = await getTemplateByID(props.id);
        template.value = response.data;
      } catch (error) {
        console.error('Failed to fetch template details:', error);
      }
    });

    const handleGenerate = () => {
      uni.chooseImage({
        count: 1,
        success: async (res) => {
          loading.value = true;
          try {
            const response = await generateImage(props.id, res.tempFiles[0]);
            uni.navigateTo({
              url: `/pages/generate/index?result=${JSON.stringify(response.data)}`,
            });
          } catch (error) {
            console.error('Failed to generate image:', error);
            uni.showToast({
              title: 'Generation failed',
              icon: 'none',
            });
          } finally {
            loading.value = false;
          }
        },
      });
    };

    return {
      template,
      loading,
      handleGenerate,
    };
  },
});
</script>

<style scoped>
.detail-container {
  background-color: #fcfbf9;
  min-height: 100vh;
}

.preview-image {
  width: 100%;
}

.info-card {
  background-color: #fff;
  padding: 20px;
  border-radius: 16px;
  margin: -30px 20px 0;
  box-shadow: 0 4px 24px rgba(74, 74, 74, 0.08);
}

.title {
  font-size: 24px;
  font-weight: 600;
  color: #4a4a4a;
  display: block;
  margin-bottom: 10px;
}

.description {
  font-size: 15px;
  color: #9b9b9b;
  line-height: 1.6;
  display: block;
  margin-bottom: 20px;
}

.cost {
  font-size: 16px;
  color: #4a4a4a;
  margin-bottom: 20px;
}

.generate-button {
  background-color: #e89b93;
  color: #fff;
  border: none;
  padding: 12px 24px;
  border-radius: 24px;
  font-size: 16px;
  cursor: pointer;
  box-shadow: 0 4px 12px rgba(232, 155, 147, 0.3);
  transition: all 0.2s ease-in-out;
  width: 100%;
}

.generate-button:active {
  transform: scale(0.97);
  box-shadow: 0 2px 6px rgba(232, 155, 147, 0.3);
}
</style> 