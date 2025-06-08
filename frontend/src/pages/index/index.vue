<template>
  <view class="gallery-container">
    <LoadingAnimation v-if="loading" />
    <view v-else class="grid">
      <TemplateCard
        v-for="(template, index) in templates"
        :key="template.id"
        :template="template"
        class="animate-stagger"
        :style="{ animationDelay: `${index * 50}ms` }"
      />
    </view>
  </view>
</template>

<script>
import { defineComponent, ref, onMounted } from 'vue';
import { getTemplates } from '../../api/template';
import TemplateCard from '../../components/TemplateCard/TemplateCard.vue';
import LoadingAnimation from '../../components/LoadingAnimation/LoadingAnimation.vue';

export default defineComponent({
  components: {
    TemplateCard,
    LoadingAnimation,
  },
  setup() {
    const templates = ref([]);
    const loading = ref(true);

    onMounted(async () => {
      try {
        const response = await getTemplates();
        templates.value = response.data.templates;
      } catch (error) {
        console.error('Failed to fetch templates:', error);
      } finally {
        loading.value = false;
      }
    });

    return {
      templates,
      loading,
    };
  },
});
</script>

<style lang="scss" scoped>
.gallery-container {
  padding: 20px;
  background-color: #fcfbf9; /* --color-bg */
}

.grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
  gap: 20px;
}

@keyframes stagger-in {
  from {
    opacity: 0;
    transform: translateY(16px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.animate-stagger {
  animation: stagger-in 0.4s cubic-bezier(0.6, 0.05, 0.4, 1) forwards;
  opacity: 0;
}
</style> 