<template>
  <n-modal 
    :show="show" 
    @update:show="$emit('update:show', $event)" 
    @after-leave="handleAfterLeave"
    :mask-closable="true"
    preset="card"
    :style="modalStyle"
    :closable="true"
    :close-on-esc="true"
  >
    <template #header>
      <div class="preview-header">
        <n-icon size="24"><Image /></n-icon>
        <span class="file-name">{{ file?.name }}</span>
      </div>
    </template>
    
    <div class="image-preview">
      <div v-if="loading" class="loading-container">
        <n-spin size="large" />
        <div class="loading-text">加载图片中...</div>
      </div>
      
      <div v-else-if="error" class="error-container">
        <n-icon size="48" color="#ea4335">
          <AlertCircle />
        </n-icon>
        <div class="error-text">{{ error }}</div>
      </div>
      
      <div v-else class="image-container">
        <n-image
          :src="imageUrl"
          :alt="file?.name"
          :img-props="{
            style: {
              maxWidth: '90%',
              maxHeight: '65vh',
              objectFit: 'contain'
            },
            onLoad: handleImageLoad,
            onError: handleImageError
          }"
          :previewed-img-props="{
            style: {
              maxWidth: '95vw',
              maxHeight: '95vh',
              objectFit: 'contain'
            }
          }"
          show-toolbar
        />
      </div>
    </div>

    <template #footer>
      <div class="preview-footer">
        <div class="file-info">
          <span>{{ formatFileSize(file?.size) }}</span>
          <span>{{ formatDate(file?.modified) }}</span>
        </div>
        <div class="actions">
          <n-button @click="openInSystem">
            <template #icon>
              <n-icon><OpenOutline /></n-icon>
            </template>
            在系统中打开
          </n-button>
          <n-button @click="close">
            关闭
          </n-button>
        </div>
      </div>
    </template>
  </n-modal>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { NModal, NIcon, NSpin, NButton, NImage, useMessage } from 'naive-ui'
import { Image, AlertCircle, OpenOutline } from '@vicons/ionicons5'
import { GetFileSystemItem, GetImageData } from '../../wailsjs/go/main/App'

interface FileItem {
  key: string
  name: string
  path: string
  isDirectory: boolean
  size?: number
  modified?: number
}

const props = defineProps<{
  show: boolean
  file?: FileItem
}>()

const emit = defineEmits<{
  'update:show': [value: boolean]
}>()

const message = useMessage()
const loading = ref(false)
const error = ref('')
const imageUrl = ref('')

const resetState = () => {
  // 安全地重置所有状态，避免触发错误
  try {
    loading.value = false
    error.value = ''
    imageUrl.value = ''
  } catch (err) {
    console.warn('Error resetting image preview state:', err)
  }
}

// 组件初始化时重置状态
resetState()

watch(() => props.show, (newShow) => {
  if (newShow && props.file) {
    loadImage()
  }
  // 当模态框关闭时，不在这里重置状态，而是等待动画完成后在 handleAfterLeave 中重置
})

const loadImage = async () => {
  if (!props.file) return
  
  // 检查模态框是否仍然显示
  if (!props.show) return
  
  loading.value = true
  error.value = ''
  imageUrl.value = ''
  
  try {
    // 获取图片的 base64 数据
    const imageData = await GetImageData(props.file.path)
    
    // 再次检查模态框是否仍然显示
    if (!props.show) return
    
    if (imageData) {
      imageUrl.value = imageData
    } else {
      error.value = '无法获取图片数据'
    }
  } catch (err) {
    // 只在模态框仍然显示时设置错误
    if (props.show) {
      console.error('Failed to load image:', err)
      error.value = '无法加载图片文件: ' + (err as Error).message
    }
  } finally {
    // 只在模态框仍然显示时更新加载状态
    if (props.show) {
      loading.value = false
    }
  }
}

const handleImageLoad = () => {
  // 图片加载成功
  console.log('Image loaded successfully')
}

const handleImageError = () => {
  console.error('Image load error')
  error.value = '图片加载失败，可能是文件格式不支持或文件已损坏'
}

const openInSystem = () => {
  if (props.file) {
    // 使用系统默认程序打开文件
    const { shell } = window as any
    if (shell && shell.openPath) {
      shell.openPath(props.file.path)
    } else {
      message.warning('无法在系统中打开文件')
    }
  }
}

const close = () => {
  // 只关闭模态框，状态重置会在动画完成后进行
  emit('update:show', false)
}

const handleAfterLeave = () => {
  // 关闭动画完成后重置状态
  resetState()
}

const formatFileSize = (bytes?: number): string => {
  if (!bytes) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB']
  let size = bytes
  let unitIndex = 0
  
  while (size >= 1024 && unitIndex < units.length - 1) {
    size /= 1024
    unitIndex++
  }
  
  return `${size.toFixed(1)} ${units[unitIndex]}`
}

const formatDate = (timestamp?: number): string => {
  if (!timestamp) return 'Unknown'
  return new Date(timestamp).toLocaleDateString()
}

// 动态计算模态框样式 - 全屏模式
const modalStyle = computed(() => {
  // 全屏预览模式
  return {
    width: 'calc(100vw - 80px)',
    maxWidth: 'none',
    height: '90vh',
    maxHeight: 'none',
    top: '20px',
    left: '40px',
    transform: 'none'
  }
})
</script>

<style scoped>
.preview-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 18px;
  font-weight: 500;
}

.file-name {
  font-weight: 500;
  color: #202124;
}

.image-preview {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 70vh;
  min-height: 400px;
  background: #f8f9fa;
  border-radius: 8px;
  overflow: hidden;
  position: relative;
  padding: 20px;
  box-sizing: border-box;
  width: 100%;
}

.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
  color: #5f6368;
}

.loading-text {
  font-size: 14px;
}

.error-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
  color: #ea4335;
}

.error-text {
  font-size: 14px;
  text-align: center;
  max-width: 300px;
}

.image-container {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100%;
  height: 100%;
  position: relative;
}

/* Naive UI Image 组件会自动处理图片样式 */

.preview-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 0;
}

.file-info {
  display: flex;
  gap: 16px;
  font-size: 12px;
  color: #5f6368;
}

.actions {
  display: flex;
  gap: 8px;
}
</style>