<template>
  <div class="file-explorer">
    <n-card title="File Explorer" class="explorer-card">
      <n-tree
        :data="fileTreeData"
        :selected-keys="selectedKeys"
        :expanded-keys="expandedKeys"
        :on-update:selected-keys="handleSelect"
        :on-update:expanded-keys="handleExpand"
        :render-label="renderLabel"
        :render-prefix="renderPrefix"
        :load="handleLoad"
        :default-expanded-keys="defaultExpandedKeys"
        cascade
        selectable
      />
    </n-card>
    
    <n-drawer
      v-model:show="showDetails"
      :width="400"
      placement="right"
    >
      <n-drawer-content title="File Details">
        <div v-if="selectedFile">
          <n-space vertical>
            <n-descriptions :column="1" label-placement="left" bordered>
              <n-descriptions-item label="Name">
                {{ selectedFile.name }}
              </n-descriptions-item>
              <n-descriptions-item label="Path">
                {{ selectedFile.path }}
              </n-descriptions-item>
              <n-descriptions-item label="Type">
                {{ selectedFile.isDirectory ? 'Directory' : 'File' }}
              </n-descriptions-item>
              <n-descriptions-item label="Size">
                {{ formatFileSize(selectedFile.size) }}
              </n-descriptions-item>
              <n-descriptions-item label="Modified">
                {{ formatDate(selectedFile.modified) }}
              </n-descriptions-item>
            </n-descriptions>
            
            <n-space v-if="!selectedFile.isDirectory">
              <n-button type="primary" @click="handleTagFile">
                Tag File
              </n-button>
            </n-space>
          </n-space>
        </div>
      </n-drawer-content>
    </n-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, h } from 'vue'
import { 
  NTree, NCard, NDrawer, NDrawerContent, NDescriptions, 
  NDescriptionsItem, NSpace, NButton, NIcon, useMessage 
} from 'naive-ui'
import { FolderOpen, Folder, Document } from '@vicons/ionicons5'
import { GetRootDirectories, GetDirectoryContents, GetFileSystemItem } from '../../wailsjs/go/main/App'

interface FileItem {
  key: string
  label: string
  path: string
  name: string
  isDirectory: boolean
  children?: FileItem[]
  size?: number
  modified?: number
  isLeaf?: boolean
}

const fileTreeData = ref<FileItem[]>([])
const selectedKeys = ref<string[]>([])
const expandedKeys = ref<string[]>([])
const defaultExpandedKeys = ref<string[]>([])
const selectedFile = ref<FileItem | null>(null)
const showDetails = ref(false)
const message = useMessage()

const loadingMap = ref<Record<string, boolean>>({})

onMounted(async () => {
  try {
    await loadRootDirectories()
  } catch (error) {
    message.error('Failed to load file system: ' + error)
  }
})

const loadRootDirectories = async () => {
  try {
    console.log('Loading root directories...')
    const rootDirs = await GetRootDirectories()
    console.log('Root directories loaded:', rootDirs)
    fileTreeData.value = rootDirs.map(item => ({
      ...item,
      children: undefined, // Will be loaded lazily
      modified: item.modified * 1000 // Convert to milliseconds
    }))
    
    if (rootDirs.length > 0) {
      defaultExpandedKeys.value = [rootDirs[0].key]
      expandedKeys.value = [rootDirs[0].key]
    }
  } catch (error) {
    console.error('Failed to load root directories:', error)
    message.error('Failed to load root directories: ' + error)
    throw error
  }
}

const handleSelect = (keys: string[]) => {
  selectedKeys.value = keys
  if (keys.length > 0) {
    const fileItem = findFileByKey(fileTreeData.value, keys[0])
    if (fileItem) {
      selectedFile.value = fileItem
      showDetails.value = true
    }
  }
}

const handleExpand = (keys: string[]) => {
  expandedKeys.value = keys
}

const handleLoad = async (node: FileItem) => {
  console.log('Loading node:', node)
  if (!node.isDirectory) {
    return []
  }
  
  const nodeKey = node.key
  loadingMap.value[nodeKey] = true
  
  try {
    console.log('Getting directory contents for:', nodeKey)
    const contents = await GetDirectoryContents(nodeKey)
    console.log('Directory contents:', contents)
    return contents.map(item => ({
      ...item,
      children: undefined, // Will be loaded lazily for directories
      modified: item.modified * 1000 // Convert to milliseconds
    }))
  } catch (error) {
    console.error('Failed to load directory contents:', error)
    message.error('Failed to load directory contents: ' + error)
    return []
  } finally {
    loadingMap.value[nodeKey] = false
  }
}

const findFileByKey = (items: FileItem[], key: string): FileItem | null => {
  for (const item of items) {
    if (item.key === key) {
      return item
    }
    if (item.children) {
      const found = findFileByKey(item.children, key)
      if (found) return found
    }
  }
  return null
}

const renderLabel = ({ option }: { option: any }) => {
  return h('span', { class: 'file-label' }, option.label)
}

const renderPrefix = ({ option }: { option: any }) => {
  if (option.isDirectory) {
    return h(NIcon, null, h(FolderOpen))
  } else {
    return h(NIcon, null, h(Document))
  }
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

const handleTagFile = () => {
  if (selectedFile.value) {
    // This will be implemented later for tagging functionality
    console.log('Tag file:', selectedFile.value.path)
  }
}
</script>

<style scoped>
.file-explorer {
  height: 100vh;
  padding: 16px;
}

.explorer-card {
  height: calc(100vh - 32px);
}

.file-label {
  user-select: none;
}

:deep(.n-tree-node-content) {
  padding: 4px 0;
}

:deep(.n-tree-node-content:hover) {
  background-color: rgba(24, 160, 88, 0.1);
}

:deep(.n-tree-node--selected .n-tree-node-content) {
  background-color: rgba(24, 160, 88, 0.2);
}
</style>