<template>
  <div class="google-drive-layout">
    <!-- 顶部工具栏 -->
    <div class="toolbar">
      <div class="toolbar-left">
        <div class="app-title">
          <n-icon size="32" color="#4285f4"><Folder /></n-icon>
          <span class="app-name">TagMyFile</span>
        </div>
        <n-button-group>
          <n-button @click="handleNewFile" type="primary">
            <template #icon>
              <n-icon><Add /></n-icon>
            </template>
            新建
          </n-button>
          <n-button @click="handleUpload">
            <template #icon>
              <n-icon><CloudUpload /></n-icon>
            </template>
            上传
          </n-button>
        </n-button-group>
      </div>
      
      <div class="toolbar-right">
        <n-input
          v-model:value="searchQuery"
          placeholder="搜索文件和文件夹"
          clearable
          style="width: 300px;"
        >
          <template #prefix>
            <n-icon><Search /></n-icon>
          </template>
        </n-input>
        
        <n-button-group>
          <n-button 
            :type="viewMode === 'grid' ? 'primary' : 'default'"
            @click="viewMode = 'grid'"
          >
            <template #icon>
              <n-icon><Grid /></n-icon>
            </template>
          </n-button>
          <n-button 
            :type="viewMode === 'list' ? 'primary' : 'default'"
            @click="viewMode = 'list'"
          >
            <template #icon>
              <n-icon><List /></n-icon>
            </template>
          </n-button>
        </n-button-group>
        
        <n-button-group v-if="viewMode === 'grid'">
          <n-button 
            :type="imageDisplayMode === 'crop' ? 'primary' : 'default'"
            @click="() => { imageDisplayMode = 'crop'; console.log('Image display mode changed to:', imageDisplayMode) }"
            title="裁切模式"
          >
            <template #icon>
              <n-icon><Crop /></n-icon>
            </template>
          </n-button>
          <n-button 
            :type="imageDisplayMode === 'full' ? 'primary' : 'default'"
            @click="() => { imageDisplayMode = 'full'; console.log('Image display mode changed to:', imageDisplayMode) }"
            title="完整显示模式"
          >
            <template #icon>
              <n-icon><Expand /></n-icon>
            </template>
          </n-button>
        </n-button-group>
        
        <n-dropdown :options="accountOptions" @select="handleAccountAction">
          <n-avatar round size="small" style="cursor: pointer;">
            <n-icon><Person /></n-icon>
          </n-avatar>
        </n-dropdown>
      </div>
    </div>

    <div class="main-content">
      <!-- 左侧导航栏 -->
      <div class="sidebar">
        <div class="nav-content">
          <div class="nav-section">
            <div class="nav-item" :class="{ active: currentView === 'my-drive' }" @click="currentView = 'my-drive'">
              <n-icon><Cloud /></n-icon>
              <span>我的云端硬盘</span>
            </div>
            <div class="nav-item" :class="{ active: currentView === 'shared' }" @click="currentView = 'shared'">
              <n-icon><People /></n-icon>
              <span>共享</span>
            </div>
            <div class="nav-item" :class="{ active: currentView === 'recent' }" @click="currentView = 'recent'">
              <n-icon><Time /></n-icon>
              <span>最近</span>
            </div>
            <div class="nav-item" :class="{ active: currentView === 'starred' }" @click="currentView = 'starred'">
              <n-icon><Star /></n-icon>
              <span>已加星标</span>
            </div>
            <div class="nav-item" :class="{ active: currentView === 'trash' }" @click="currentView = 'trash'">
              <n-icon><Trash /></n-icon>
              <span>回收站</span>
            </div>
          </div>

          <div class="nav-section">
            <div class="nav-title">标签</div>
            <div class="nav-item" :class="{ active: currentView === 'all-files' }" @click="currentView = 'all-files'">
              <n-icon><FolderOpen /></n-icon>
              <span>所有文件</span>
            </div>
            <div class="nav-item" :class="{ active: currentView === 'images' }" @click="currentView = 'images'">
              <n-icon><Image /></n-icon>
              <span>图片</span>
            </div>
            <div class="nav-item" :class="{ active: currentView === 'documents' }" @click="currentView = 'documents'">
              <n-icon><DocumentText /></n-icon>
              <span>文档</span>
            </div>
            <div class="nav-item" :class="{ active: currentView === 'videos' }" @click="currentView = 'videos'">
              <n-icon><Videocam /></n-icon>
              <span>视频</span>
            </div>
          </div>
        </div>

        <!-- 存储空间信息 -->
        <div class="storage-info">
          <div class="storage-bar">
            <div class="storage-used" :style="{ width: diskSpaceInfo.usagePercent + '%' }"></div>
          </div>
          <div class="storage-text">
            <span>{{ diskSpaceInfo.usedSpaceGB.toFixed(1) }} GB 已使用</span>
            <span>共 {{ diskSpaceInfo.totalSpaceGB.toFixed(1) }} GB</span>
          </div>
        </div>
      </div>

      <!-- 右侧内容区域 -->
      <div class="content-area">
        <!-- 路径导航 -->
        <div class="path-navigation">
          <div class="path-controls">
            <n-button text @click="goToParent" :disabled="!currentPath.trim()" title="返回上级文件夹">
              <n-icon><ArrowBack /></n-icon>
            </n-button>
            <n-button text @click="goToRoot" title="根目录">
              <n-icon><Home /></n-icon>
            </n-button>
            <n-button text @click="resetToRoot" title="重置到根目录（清除记忆）">
              <n-icon><Refresh /></n-icon>
            </n-button>
          </div>
          <div class="path-input-wrapper">
            <n-input
              v-model:value="currentPath"
              placeholder="输入路径 (如: C:\Users\YourName\Documents) 或清空查看根目录"
              clearable
              @keyup.enter="handlePathEnter"
              @clear="handlePathClear"
            >
              <template #suffix>
                <n-button text @click="handlePathEnter" :disabled="!currentPath.trim()">
                  <n-icon><ArrowForward /></n-icon>
                </n-button>
              </template>
            </n-input>
          </div>
        </div>

        <!-- 文件列表 -->
        <div class="file-container">
          <div v-if="viewMode === 'grid'" class="grid-view">
            <div class="grid-container">
              <div 
                v-for="file in currentFiles" 
                :key="file.key"
                class="grid-item"
                @click="handleFileClick(file)"
                @contextmenu.prevent="handleRightClick($event, file)"
              >
                <div class="file-icon">
                  <n-icon size="48" v-if="file.isDirectory">
                    <Folder />
                  </n-icon>
                  <lazy-thumbnail
                    v-else-if="isImageFile(file.name)"
                    :file-path="file.path"
                    :file-name="file.name"
                    :display-mode="imageDisplayMode"
                    grid-size="large"
                    @loaded="(url) => imageUrls.set(file.path, url)"
                  />
                  <n-icon size="48" v-else>
                    <Document />
                  </n-icon>
                </div>
                <div class="file-name">{{ file.name }}</div>
                <div class="file-info" v-if="!file.isDirectory">
                  {{ formatFileSize(file.size) }}
                </div>
                <div class="file-info" v-else-if="isDriveLetter(file.name)">
                  本地磁盘
                </div>
                <div class="file-info" v-else>
                  文件夹
                </div>
                <div class="file-info" v-if="file.modified">
                  {{ formatDate(file.modified) }}
                </div>
              </div>
            </div>
          </div>

          <div v-else class="list-view">
            <n-data-table
              :columns="columns"
              :data="currentFiles"
              :row-key="row => row.key"
              :scroll-x="800"
              :scroll-y="600"
              :max-height="600"
              :virtual-scroll="true"
              :item-size="54"
              @row-click="handleRowClick"
              @contextmenu="handleRowRightClick"
            />
          </div>
        </div>
      </div>
    </div>

    <!-- 右键菜单 -->
    <n-dropdown
      :show="contextMenu.show"
      :options="contextMenuOptions"
      :x="contextMenu.x"
      :y="contextMenu.y"
      @clickoutside="closeContextMenu"
      @select="handleContextMenuAction"
    />

    <!-- 标签管理器 -->
    <TagManager
      v-model:show="tagManager.show"
      :file="tagManager.file || undefined"
      @save="handleTagSave"
    />

    <!-- 图片预览 -->
    <ImagePreview
      v-model:show="imagePreview.show"
      :file="imagePreview.file || undefined"
    />

    </div>
</template>

<script setup lang="ts">
import { ref, computed, h, onMounted } from 'vue'
import { 
  NButton, NButtonGroup, NIcon, NInput, NAvatar, NDropdown, 
  NDataTable, NVirtualList, useMessage
} from 'naive-ui'
import TagManager from './TagManager.vue'
import LazyThumbnail from './LazyThumbnail.vue'
import ImagePreview from './ImagePreview.vue'
import {
  Add, CloudUpload, Search, Grid, List, Person, Cloud, People, 
  Time, Star, Trash, FolderOpen, Image, DocumentText, Videocam, 
  Home, Folder, Document, ArrowForward, ArrowBack, Crop, Expand, Refresh
} from '@vicons/ionicons5'
import { GetRootDirectories, GetDirectoryContents, GetImageData, GetDiskSpaceInfo, OpenFile } from '../../wailsjs/go/main/App'

interface FileItem {
  key: string
  label: string
  path: string
  name: string
  isDirectory: boolean
  size?: number
  modified?: number
  children?: FileItem[]
}

const message = useMessage()
const viewMode = ref<'grid' | 'list'>('grid')
const imageDisplayMode = ref<'crop' | 'full'>('crop')
const currentView = ref('my-drive')
const searchQuery = ref('')
const currentPath = ref('')
const files = ref<FileItem[]>([])
const selectedFile = ref<FileItem | null>(null)
const imageUrls = ref<Map<string, string>>(new Map())
const diskSpaceInfo = ref({
  totalSpaceGB: 0,
  usedSpaceGB: 0,
  freeSpaceGB: 0,
  usagePercent: 0
})

const contextMenu = ref({
  show: false,
  x: 0,
  y: 0,
  file: null as FileItem | null
})

const tagManager = ref({
  show: false,
  file: null as FileItem | null
})

const imagePreview = ref({
  show: false,
  file: null as FileItem | null
})

const openingFiles = ref(new Set<string>()) // Track files being opened


const accountOptions = [
  { label: '账户设置', key: 'account' },
  { label: '帮助', key: 'help' },
  { label: '退出', key: 'logout' }
]


const currentFiles = computed(() => {
  if (!searchQuery.value) {
    return files.value
  }
  const query = searchQuery.value.toLowerCase()
  return files.value.filter(file => 
    file.name.toLowerCase().includes(query)
  )
})

const columns = [
  {
    title: '名称',
    key: 'name',
    render: (row: FileItem) => {
      let iconElement: any
      
      if (row.isDirectory) {
        iconElement = h(Folder)
      } else if (isImageFile(row.name)) {
        iconElement = h(LazyThumbnail, {
          filePath: row.path,
          fileName: row.name,
          displayMode: imageDisplayMode.value,
          gridSize: 'small',
          onLoaded: (url: string) => imageUrls.value.set(row.path, url)
        })
      } else {
        iconElement = h(Document)
      }
      
      return h('div', { class: 'file-row-name' }, [
        h(NIcon, { size: 20 }, { default: () => iconElement }),
        h('span', { class: 'file-name-text' }, row.name)
      ])
    }
  },
  {
    title: '所有者',
    key: 'owner',
    render: () => '我'
  },
  {
    title: '最后修改',
    key: 'modified',
    render: (row: FileItem) => formatDate(row.modified)
  },
  {
    title: '文件大小',
    key: 'size',
    render: (row: FileItem) => row.isDirectory ? '-' : formatFileSize(row.size)
  }
]

const contextMenuOptions = [
  { label: '打开', key: 'open' },
  { label: '重命名', key: 'rename' },
  { label: '添加标签', key: 'tag' },
  { type: 'divider' },
  { label: '下载', key: 'download' },
  { label: '移至回收站', key: 'delete' }
]

onMounted(async () => {
  await loadDiskSpaceInfo()
  await restoreLastPath()
})

const loadRootDirectories = async () => {
  try {
    const rootDirs = await GetRootDirectories()
    files.value = rootDirs.map(item => ({
      ...item,
      modified: item.modified * 1000
    }))
  } catch (error) {
    message.error('Failed to load root directories: ' + error)
  }
}

const loadDiskSpaceInfo = async () => {
  try {
    const info = await GetDiskSpaceInfo()
    if (info) {
      diskSpaceInfo.value = {
        totalSpaceGB: info.totalSpaceGB,
        usedSpaceGB: info.usedSpaceGB,
        freeSpaceGB: info.freeSpaceGB,
        usagePercent: info.usagePercent
      }
    }
  } catch (error) {
    console.error('Failed to load disk space info:', error)
    // 使用默认值保持UI正常显示
    diskSpaceInfo.value = {
      totalSpaceGB: 100,
      usedSpaceGB: 65,
      freeSpaceGB: 35,
      usagePercent: 65
    }
  }
}

const loadDirectoryContents = async (path: string) => {
  try {
    const contents = await GetDirectoryContents(path)
    files.value = contents.map(item => ({
      ...item,
      modified: item.modified * 1000
    }))
    // 确保路径格式正确
    currentPath.value = path.replace(/\//g, '\\')
    
    // 保存当前路径到本地存储
    saveCurrentPath(currentPath.value)
    
    // 清空图片缓存
    imageUrls.value.clear()
  } catch (error) {
    message.error(`无法访问路径: ${path}`)
    console.error('Directory load error:', error)
  }
}

const handleFileClick = (file: FileItem) => {
  if (file.isDirectory) {
    loadDirectoryContents(file.path)
  } else {
    // For image files, show preview
    if (isImageFile(file.name)) {
      showImagePreview(file)
    } else {
      // For non-image files, open with default program
      openFile(file)
    }
  }
}

const handleRowClick = (row: FileItem) => {
  handleFileClick(row)
}

const handleRightClick = (event: MouseEvent, file: FileItem) => {
  contextMenu.value = {
    show: true,
    x: event.clientX,
    y: event.clientY,
    file
  }
}

const handleRowRightClick = (row: FileItem, event: MouseEvent) => {
  handleRightClick(event, row)
}

const closeContextMenu = () => {
  contextMenu.value.show = false
}

const handleContextMenuAction = (key: string) => {
  const file = contextMenu.value.file
  if (!file) return

  switch (key) {
    case 'open':
      if (file.isDirectory) {
        handleFileClick(file)
      } else {
        openFile(file)
      }
      break
    case 'tag':
      openTagManager(file)
      break
    case 'delete':
      message.info(`Move to trash: ${file.name}`)
      break
  }
  closeContextMenu()
}

const showImagePreview = (file: FileItem) => {
  imagePreview.value = {
    show: true,
    file
  }
}

const openFile = async (file: FileItem) => {
  // Check if file is already being opened
  if (openingFiles.value.has(file.path)) {
    message.info(`正在打开 ${file.name}，请稍候...`)
    return
  }
  
  // Add to opening files set
  openingFiles.value.add(file.path)
  
  // Show loading message
  const loadingMessage = message.loading(`正在打开 ${file.name}...`, { duration: 0 })
  
  try {
    await OpenFile(file.path)
    // Success - file opened with default program
    loadingMessage.destroy()
    message.success(`已打开 ${file.name}`)
  } catch (error) {
    loadingMessage.destroy()
    message.error(`无法打开文件: ${file.name}`)
    console.error('Failed to open file:', error)
  } finally {
    // Remove from opening files set
    openingFiles.value.delete(file.path)
  }
}

const openTagManager = (file: FileItem) => {
  tagManager.value = {
    show: true,
    file
  }
}

const handleTagSave = (tags: string[]) => {
  const file = tagManager.value.file
  if (file) {
    message.success(`已为 ${file.name} 添加 ${tags.length} 个标签`)
  }
}

const handlePathEnter = () => {
  const path = currentPath.value.trim()
  if (!path) {
    loadRootDirectories()
    return
  }
  
  // 标准化路径分隔符，确保使用反斜杠
  let normalizedPath = path.replace(/\//g, '\\')
  
  // 确保盘符路径格式正确（如 D:\ 而不是 D:）
  if (/^[A-Za-z]:$/.test(normalizedPath)) {
    normalizedPath += '\\'
  }
  
  loadDirectoryContents(normalizedPath)
}

const handlePathClear = () => {
  currentPath.value = ''
  // 清空保存的路径
  saveCurrentPath('')
  loadRootDirectories()
}

const saveCurrentPath = (path: string) => {
  try {
    localStorage.setItem('tagMyFile_currentPath', path)
  } catch (error) {
    console.warn('Failed to save current path:', error)
  }
}

const restoreLastPath = async () => {
  try {
    const savedPath = localStorage.getItem('tagMyFile_currentPath')
    if (savedPath && savedPath.trim()) {
      // 尝试恢复保存的路径
      await loadDirectoryContents(savedPath)
    } else {
      // 如果没有保存的路径，显示根目录
      await loadRootDirectories()
    }
  } catch (error) {
    console.warn('Failed to restore last path:', error)
    // 恢复失败时显示根目录
    await loadRootDirectories()
  }
}

const goToRoot = () => {
  currentPath.value = ''
  // 清空保存的路径
  saveCurrentPath('')
  loadRootDirectories()
}

const resetToRoot = () => {
  currentPath.value = ''
  // 清空保存的路径
  saveCurrentPath('')
  loadRootDirectories()
  message.success('已重置到根目录')
}

const goToParent = () => {
  const path = currentPath.value.trim()
  if (!path) return
  
  // 使用路径分隔符获取父目录
  let parentPath = path.substring(0, path.lastIndexOf('\\'))
  
  // 处理特殊情况：如果当前是盘符根目录（如 D:\），则返回根目录列表
  if (/^[A-Za-z]:\\$/.test(path)) {
    goToRoot()
    return
  }
  
  if (parentPath) {
    // 如果父路径是盘符（如 D:），确保添加反斜杠
    if (/^[A-Za-z]:$/.test(parentPath)) {
      parentPath += '\\'
    }
    loadDirectoryContents(parentPath)
  } else {
    // 如果没有父路径，返回根目录
    goToRoot()
  }
}

const navigateTo = (path: string) => {
  if (path === '') {
    loadRootDirectories()
    currentPath.value = ''
  }
}

const handleNewFile = () => {
  message.info('Create new file functionality')
}

const handleUpload = () => {
  message.info('Upload file functionality')
}

const handleAccountAction = (key: string) => {
  message.info(`Account action: ${key}`)
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

const isImageFile = (fileName: string): boolean => {
  const imageExtensions = ['.jpg', '.jpeg', '.png', '.gif', '.bmp', '.webp']
  const ext = fileName.toLowerCase().substring(fileName.lastIndexOf('.'))
  return imageExtensions.includes(ext)
}

const isDriveLetter = (fileName: string): boolean => {
  // 盘符格式通常是 "C:" 或 "D:" 等
  return /^[A-Za-z]:$/.test(fileName)
}

const formatDate = (timestamp?: number): string => {
  if (!timestamp) return 'Unknown'
  return new Date(timestamp).toLocaleString()
}


</script>

<style scoped>
.google-drive-layout {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: #f8f9fa;
}

.toolbar {
  height: 64px;
  background: white;
  border-bottom: 1px solid #e0e0e0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.app-title {
  display: flex;
  align-items: center;
  gap: 8px;
}

.app-logo {
  width: 32px;
  height: 32px;
}

.app-name {
  font-size: 22px;
  font-weight: 400;
  color: #5f6368;
}

.toolbar-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.main-content {
  flex: 1;
  display: flex;
  overflow: hidden;
}

.sidebar {
  width: 256px;
  background: white;
  border-right: 1px solid #e0e0e0;
  display: flex;
  flex-direction: column;
  height: 100%;
}

.nav-section {
  padding: 8px 0;
}

.nav-content {
  flex: 1;
  overflow-y: auto;
  min-height: 0;
}

.nav-title {
  padding: 8px 24px;
  font-size: 12px;
  font-weight: 500;
  color: #5f6368;
  text-transform: uppercase;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 20px;
  padding: 8px 24px;
  cursor: pointer;
  transition: background-color 0.2s;
  color: #5f6368;
}

.nav-item:hover {
  background: #f1f3f4;
}

.nav-item.active {
  background: #e8f0fe;
  color: #1967d2;
}

.nav-item .n-icon {
  font-size: 20px;
  color: #5f6368;
}

.nav-item.active .n-icon {
  color: #1967d2;
}

.storage-info {
  margin-top: auto;
  padding: 16px 24px;
  border-top: 1px solid #e0e0e0;
}

.storage-bar {
  height: 4px;
  background: #e0e0e0;
  border-radius: 2px;
  overflow: hidden;
  margin-bottom: 8px;
}

.storage-used {
  height: 100%;
  background: #34a853;
}

.storage-text {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: #5f6368;
}

.content-area {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.path-navigation {
  padding: 16px 24px;
  background: white;
  border-bottom: 1px solid #e0e0e0;
  display: flex;
  align-items: center;
  gap: 8px;
}

.path-controls {
  display: flex;
  gap: 4px;
}

.path-input-wrapper {
  flex: 1;
}

.path-navigation .n-input {
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  text-align: left;
}

.path-navigation .n-input .n-icon {
  transition: color 0.2s ease;
}

.path-navigation .n-input .n-icon:hover {
  color: #4285f4;
}

.file-container {
  flex: 1;
  overflow-y: auto;
  padding: 24px;
}

.grid-view {
  width: 100%;
  height: 100%;
  overflow-y: auto;
}

.grid-container {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 16px;
  padding: 16px;
  width: 100%;
  box-sizing: border-box;
}

.grid-item {
  background: white;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  padding: 16px;
  text-align: center;
  cursor: pointer;
  transition: all 0.2s;
  min-width: 0;
  width: 100%;
  box-sizing: border-box;
}

.grid-item:hover {
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  transform: translateY(-2px);
}

.file-icon {
  margin-bottom: 8px;
  color: #5f6368;
  width: 100%;
  height: 96px;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
}

.image-container {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.file-name {
  font-size: 14px;
  font-weight: 500;
  color: #202124;
  margin-bottom: 4px;
  word-break: break-word;
}

.file-info {
  font-size: 12px;
  color: #5f6368;
}

.list-view {
  background: white;
  border-radius: 8px;
  overflow: hidden;
  height: 600px;
}

.list-view .n-data-table {
  height: 100%;
}

.file-row-name {
  display: flex;
  align-items: center;
  gap: 8px;
}

.file-name-text {
  font-weight: 500;
}

</style>