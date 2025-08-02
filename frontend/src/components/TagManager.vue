<template>
  <n-modal :show="show" @update:show="$emit('update:show', $event)" :mask-closable="true" preset="card" style="max-width: 600px;">
    <template #header>
      <div class="modal-header">
        <n-icon size="24"><Pricetag /></n-icon>
        <span>标签管理</span>
      </div>
    </template>
    
    <div class="tag-manager">
      <!-- 文件信息 -->
      <div class="file-info" v-if="file">
        <n-icon size="48">
          <Folder v-if="file.isDirectory" />
          <Document v-else />
        </n-icon>
        <div class="file-details">
          <div class="file-name">{{ file.name }}</div>
          <div class="file-meta">
            {{ formatFileSize(file.size) }} • {{ formatDate(file.modified) }}
          </div>
        </div>
      </div>

      <!-- 现有标签 -->
      <div class="section">
        <div class="section-title">现有标签</div>
        <div class="tags-container">
          <n-tag 
            v-for="tag in existingTags" 
            :key="tag.id"
            :type="selectedTags.includes(tag.id) ? 'success' : 'default'"
            :closable="selectedTags.includes(tag.id)"
            @close="removeTag(tag.id)"
            @click="toggleTag(tag.id)"
            style="cursor: pointer; margin: 4px;"
          >
            <template #icon>
              <n-icon :color="tag.color"><Bookmark /></n-icon>
            </template>
            {{ tag.name }}
          </n-tag>
        </div>
      </div>

      <!-- 创建新标签 -->
      <div class="section">
        <div class="section-title">创建新标签</div>
        <div class="create-tag">
          <n-input
            v-model:value="newTagName"
            placeholder="标签名称"
            style="flex: 1;"
            @keyup.enter="createTag"
          />
          <n-color-picker
            v-model:value="newTagColor"
            :show-alpha="false"
            style="margin: 0 8px;"
          />
          <n-button @click="createTag" :disabled="!newTagName.trim()">
            创建
          </n-button>
        </div>
      </div>

      <!-- 推荐标签 -->
      <div class="section">
        <div class="section-title">推荐标签</div>
        <div class="tags-container">
          <n-tag 
            v-for="tag in suggestedTags" 
            :key="tag.id"
            @click="addTag(tag)"
            style="cursor: pointer; margin: 4px;"
          >
            <template #icon>
              <n-icon :color="tag.color"><Bookmark /></n-icon>
            </template>
            {{ tag.name }}
          </n-tag>
        </div>
      </div>
    </div>

    <template #footer>
      <div class="modal-footer">
        <n-button @click="close">取消</n-button>
        <n-button type="primary" @click="saveTags" :disabled="!hasChanges">
          保存
        </n-button>
      </div>
    </template>
  </n-modal>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { 
  NModal, NIcon, NTag, NInput, NButton, NColorPicker, useMessage 
} from 'naive-ui'
import { Pricetag, Folder, Document, Bookmark } from '@vicons/ionicons5'

interface FileItem {
  key: string
  name: string
  isDirectory: boolean
  size?: number
  modified?: number
}

interface TagItem {
  id: string
  name: string
  color: string
}

const props = defineProps<{
  show: boolean
  file?: FileItem
}>()

const emit = defineEmits<{
  'update:show': [value: boolean]
  'save': [tags: string[]]
}>()

const message = useMessage()

const existingTags = ref<TagItem[]>([
  { id: 'work', name: '工作', color: '#4285f4' },
  { id: 'personal', name: '个人', color: '#34a853' },
  { id: 'important', name: '重要', color: '#ea4335' },
  { id: 'archive', name: '归档', color: '#fbbc04' },
  { id: 'images', name: '图片', color: '#9c27b0' },
  { id: 'documents', name: '文档', color: '#ff9800' }
])

const suggestedTags = ref<TagItem[]>([
  { id: 'music', name: '音乐', color: '#e91e63' },
  { id: 'videos', name: '视频', color: '#673ab7' },
  { id: 'downloads', name: '下载', color: '#607d8b' },
  { id: 'projects', name: '项目', color: '#795548' }
])

const selectedTags = ref<string[]>([])
const newTagName = ref('')
const newTagColor = ref('#4285f4')

const hasChanges = computed(() => {
  return selectedTags.value.length > 0
})

watch(() => props.show, (newShow) => {
  if (newShow && props.file) {
    loadFileTags()
  }
})

const loadFileTags = () => {
  // 模拟加载文件的标签
  selectedTags.value = []
}

const toggleTag = (tagId: string) => {
  const index = selectedTags.value.indexOf(tagId)
  if (index === -1) {
    selectedTags.value.push(tagId)
  } else {
    selectedTags.value.splice(index, 1)
  }
}

const removeTag = (tagId: string) => {
  const index = selectedTags.value.indexOf(tagId)
  if (index !== -1) {
    selectedTags.value.splice(index, 1)
  }
}

const addTag = (tag: TagItem) => {
  if (!selectedTags.value.includes(tag.id)) {
    selectedTags.value.push(tag.id)
  }
}

const createTag = () => {
  if (!newTagName.value.trim()) return

  const newTag: TagItem = {
    id: newTagName.value.toLowerCase().replace(/\s+/g, '-'),
    name: newTagName.value.trim(),
    color: newTagColor.value
  }

  existingTags.value.push(newTag)
  selectedTags.value.push(newTag.id)
  
  newTagName.value = ''
  newTagColor.value = '#4285f4'
  
  message.success('标签创建成功')
}

const saveTags = () => {
  emit('save', selectedTags.value)
  close()
}

const close = () => {
  emit('update:show', false)
  selectedTags.value = []
  newTagName.value = ''
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
</script>

<style scoped>
.modal-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 18px;
  font-weight: 500;
}

.tag-manager {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.file-info {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px;
  background: #f8f9fa;
  border-radius: 8px;
}

.file-details {
  flex: 1;
}

.file-name {
  font-size: 16px;
  font-weight: 500;
  color: #202124;
  margin-bottom: 4px;
}

.file-meta {
  font-size: 14px;
  color: #5f6368;
}

.section {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.section-title {
  font-size: 14px;
  font-weight: 500;
  color: #5f6368;
  margin-bottom: 8px;
}

.tags-container {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  min-height: 40px;
  padding: 8px;
  border: 1px solid #e0e0e0;
  border-radius: 4px;
  background: #fafafa;
}

.create-tag {
  display: flex;
  align-items: center;
  gap: 8px;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}
</style>