<template>
  <div class="image-thumbnail">
    <img 
      v-if="thumbnailUrl"
      :src="thumbnailUrl"
      :alt="fileName"
      :class="['thumbnail-image', { 'full': displayMode === 'full' }]"
      @load="handleImageLoad"
      @error="handleImageError"
    />
    <div v-else-if="loading" class="thumbnail-loading">
      <n-spin :size="spinSize" />
    </div>
    <div v-else class="thumbnail-placeholder">
      <n-icon :size="iconSize" color="#999">
        <Image />
      </n-icon>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, computed, onUnmounted } from 'vue'
import { NIcon, NSpin } from 'naive-ui'
import { Image } from '@vicons/ionicons5'
import { GetImageData } from '../../wailsjs/go/main/App'

interface Props {
  filePath: string
  fileName: string
  size?: number
  displayMode?: 'crop' | 'full'
}

const props = withDefaults(defineProps<Props>(), {
  size: 48,
  displayMode: 'crop'
})

const thumbnailUrl = ref('')
const loading = ref(false)
const error = ref(false)
const loadId = ref<number | string | null>(null)

const iconSize = computed(() => Math.max(12, props.size * 0.5))
const spinSize = computed<'small' | 'medium' | 'large' | number>(() => props.size > 30 ? 'small' : 20)

// 缩略图加载队列管理
class ThumbnailLoader {
  private static instance: ThumbnailLoader
  private queue: Array<{
    id: number
    filePath: string
    fileName: string
    resolve: (url: string) => void
    reject: (error: Error) => void
  }> = []
  private loading = new Set<number | string>()
  private maxConcurrent = 4 // 最大并发数
  private nextId = 1
  private cache = new Map<string, string>() // 缓存已加载的图片
  
  // 设置最大并发数
  setMaxConcurrent(max: number) {
    this.maxConcurrent = Math.max(1, Math.min(max, 10)) // 限制在 1-10 之间
  }
  
  // 获取当前队列状态
  getStatus() {
    return {
      queueLength: this.queue.length,
      loadingCount: this.loading.size,
      maxConcurrent: this.maxConcurrent,
      cacheSize: this.cache.size
    }
  }
  
  // 清除缓存
  clearCache() {
    this.cache.clear()
  }
  
  // 清除指定文件的缓存
  clearFileCache(filePath: string, fileName: string) {
    const cacheKey = `${filePath}_${fileName}`
    this.cache.delete(cacheKey)
  }
  
  static getInstance(): ThumbnailLoader {
    if (!ThumbnailLoader.instance) {
      ThumbnailLoader.instance = new ThumbnailLoader()
    }
    return ThumbnailLoader.instance
  }
  
  async loadThumbnail(filePath: string, fileName: string): Promise<string> {
    // 检查缓存
    const cacheKey = `${filePath}_${fileName}`
    if (this.cache.has(cacheKey)) {
      return this.cache.get(cacheKey)!
    }
    
    const id = this.nextId++
    
    return new Promise((resolve, reject) => {
      this.queue.push({ id, filePath, fileName, resolve, reject })
      this.processQueue()
    })
  }
  
  private async processQueue() {
    if (this.loading.size >= this.maxConcurrent || this.queue.length === 0) {
      return
    }
    
    const task = this.queue.shift()
    if (!task) return
    
    this.loading.add(task.id)
    
    try {
      const imageData = await GetImageData(task.filePath)
      if (imageData) {
        // 缓存结果
        const cacheKey = `${task.filePath}_${task.fileName}`
        this.cache.set(cacheKey, imageData)
        
        // 限制缓存大小
        if (this.cache.size > 100) {
          // 删除最旧的缓存项
          const firstKey = this.cache.keys().next().value
          this.cache.delete(firstKey)
        }
        
        task.resolve(imageData)
      } else {
        task.reject(new Error('Failed to load image data'))
      }
    } catch (err) {
      task.reject(err as Error)
    } finally {
      this.loading.delete(task.id)
      // 处理队列中的下一个任务
      this.processQueue()
    }
  }
  
  cancel(id: number | string) {
    // 从队列中移除未开始的任务
    this.queue = this.queue.filter(task => task.id !== id)
    // 如果正在加载，我们无法取消，但可以标记为已取消
    this.loading.delete(id)
  }
}

const thumbnailLoader = ThumbnailLoader.getInstance()

const isImageFile = (fileName: string): boolean => {
  const imageExtensions = ['.jpg', '.jpeg', '.png', '.gif', '.bmp', '.webp']
  const ext = fileName.toLowerCase().substring(fileName.lastIndexOf('.'))
  return imageExtensions.includes(ext)
}

const loadThumbnail = async () => {
  if (!isImageFile(props.fileName)) {
    return
  }
  
  loading.value = true
  error.value = false
  thumbnailUrl.value = ''
  
  try {
    // 获取一个唯一的加载ID
    const currentLoadId = Date.now() + Math.random()
    loadId.value = currentLoadId
    
    const imageData = await thumbnailLoader.loadThumbnail(props.filePath, props.fileName)
    
    // 检查是否已经被取消或组件已卸载
    if (loadId.value === currentLoadId) {
      thumbnailUrl.value = imageData
    }
  } catch (err) {
    console.warn('Failed to load thumbnail:', err)
    if (loadId.value !== null) {
      error.value = true
    }
  } finally {
    if (loadId.value !== null) {
      loading.value = false
    }
  }
}

const handleImageLoad = () => {
  // 图片加载成功
}

const handleImageError = () => {
  error.value = true
  thumbnailUrl.value = ''
}


// 监听文件路径变化
watch(() => props.filePath, () => {
  thumbnailUrl.value = ''
  error.value = false
  loadThumbnail()
})

// 组件挂载时加载缩略图
onMounted(() => {
  // 延迟加载，避免大量组件同时初始化造成的性能问题
  setTimeout(() => {
    loadThumbnail()
  }, Math.random() * 1000) // 随机延迟 0-1 秒
})

// 组件卸载时取消加载
onUnmounted(() => {
  if (loadId.value !== null) {
    thumbnailLoader.cancel(loadId.value)
    loadId.value = null
  }
})

// 暴露调试方法到全局作用域（开发环境）
if (import.meta.env.DEV) {
  ;(window as any).__thumbnailLoader = thumbnailLoader
}
</script>

<style scoped>
.image-thumbnail {
  width: 100%;
  height: v-bind(size + 'px');
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
  overflow: hidden;
  background: #f5f5f5;
  position: relative;
}

.thumbnail-image {
  width: 100%;
  height: 100%;
  border-radius: 4px;
  object-fit: cover;
}

.thumbnail-image.full {
  object-fit: contain;
}

.thumbnail-loading {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
}

.thumbnail-placeholder {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
  color: #999;
}

</style>