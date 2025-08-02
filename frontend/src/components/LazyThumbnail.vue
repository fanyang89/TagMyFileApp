<template>
  <div class="lazy-thumbnail" ref="thumbnailRef">
    <div v-if="loading" class="thumbnail-loading">
      <n-spin :size="gridSize === 'small' ? 'small' : 'medium'" />
    </div>
    <div v-else-if="error" class="thumbnail-error">
      <n-icon :size="iconSize" color="#999">
        <Image />
      </n-icon>
    </div>
    <div v-else-if="imageUrl" class="thumbnail-content">
      <n-image
        :src="imageUrl"
        :alt="fileName"
        :img-props="{
          style: {
            objectFit: displayMode === 'crop' ? 'cover' : 'contain',
            width: '100%',
            height: gridSize === 'small' ? '48px' : '96px'
          }
        }"
        show-toolbar
      />
    </div>
    <div v-else class="thumbnail-placeholder">
      <n-icon :size="iconSize" color="#999">
        <Image />
      </n-icon>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { NIcon, NSpin, NImage } from 'naive-ui'
import { Image } from '@vicons/ionicons5'
import { GetImageData } from '../../wailsjs/go/main/App'

interface Props {
  filePath: string
  fileName: string
  displayMode: 'crop' | 'full'
  gridSize: 'small' | 'large'
  isVisible?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  isVisible: false
})

const emit = defineEmits<{
  loaded: [url: string]
  error: [error: string]
}>()

const thumbnailRef = ref<HTMLElement>()
const imageUrl = ref('')
const loading = ref(false)
const error = ref('')
const observer = ref<IntersectionObserver>()

const iconSize = computed(() => props.gridSize === 'small' ? 20 : 24)

const isImageFile = (fileName: string): boolean => {
  const imageExtensions = ['.jpg', '.jpeg', '.png', '.gif', '.bmp', '.webp']
  const ext = fileName.toLowerCase().substring(fileName.lastIndexOf('.'))
  return imageExtensions.includes(ext)
}

const loadImage = async () => {
  if (!isImageFile(props.fileName) || loading.value || imageUrl.value) {
    return
  }

  loading.value = true
  error.value = ''

  try {
    const imageData = await GetImageData(props.filePath)
    if (imageData) {
      imageUrl.value = imageData
      emit('loaded', imageData)
    } else {
      error.value = '无法加载图片'
      emit('error', '无法加载图片')
    }
  } catch (err) {
    console.error('Failed to load image:', err)
    error.value = '加载失败'
    emit('error', '加载失败')
  } finally {
    loading.value = false
  }
}

const setupIntersectionObserver = () => {
  if (!thumbnailRef.value) return

  observer.value = new IntersectionObserver(
    (entries) => {
      entries.forEach((entry) => {
        if (entry.isIntersecting && !imageUrl.value && !loading.value) {
          loadImage()
          // 加载一次后停止观察
          observer.value?.unobserve(entry.target)
        }
      })
    },
    {
      root: null,
      rootMargin: '50px',
      threshold: 0.1
    }
  )

  observer.value.observe(thumbnailRef.value)
}

onMounted(() => {
  // 如果组件已经可见，直接加载
  if (props.isVisible && !imageUrl.value && !loading.value) {
    loadImage()
  } else {
    setupIntersectionObserver()
  }
})

onUnmounted(() => {
  if (observer.value) {
    observer.value.disconnect()
  }
})
</script>

<style scoped>
.lazy-thumbnail {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
  overflow: hidden;
  background: #f5f5f5;
  position: relative;
}

.thumbnail-loading,
.thumbnail-error,
.thumbnail-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.thumbnail-content {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.thumbnail-error {
  color: #999;
}
</style>