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

        <!-- 存储空间信息 -->
        <div class="storage-info">
          <div class="storage-bar">
            <div class="storage-used" :style="{ width: '65%' }"></div>
          </div>
          <div class="storage-text">
            <span>6.5 GB 已使用</span>
            <span>共 10 GB</span>
          </div>
        </div>
      </div>

      <!-- 右侧内容区域 -->
      <div class="content-area">
        <!-- 面包屑导航 -->
        <div class="breadcrumb" v-if="currentPath">
          <n-breadcrumb>
            <n-breadcrumb-item @click="navigateTo('')">
              <n-icon><Home /></n-icon>
            </n-breadcrumb-item>
            <n-breadcrumb-item 
              v-for="(segment, index) in pathSegments" 
              :key="index"
              @click="navigateToPath(index)"
            >
              {{ segment }}
            </n-breadcrumb-item>
          </n-breadcrumb>
        </div>

        <!-- 文件列表 -->
        <div class="file-container">
          <div v-if="viewMode === 'grid'" class="grid-view">
            <div 
              v-for="file in currentFiles" 
              :key="file.key"
              class="grid-item"
              @click="handleFileClick(file)"
              @contextmenu.prevent="handleRightClick($event, file)"
            >
              <div class="file-icon">
                <n-icon size="48">
                  <Folder v-if="file.isDirectory" />
                  <Document v-else />
                </n-icon>
              </div>
              <div class="file-name">{{ file.name }}</div>
              <div class="file-info">
                {{ formatFileSize(file.size) }}
              </div>
            </div>
          </div>

          <div v-else class="list-view">
            <n-data-table
              :columns="columns"
              :data="currentFiles"
              :row-key="row => row.key"
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
  </div>
</template>

<script setup lang="ts">
import { ref, computed, h, onMounted } from 'vue'
import { 
  NButton, NButtonGroup, NIcon, NInput, NAvatar, NDropdown, 
  NBreadcrumb, NBreadcrumbItem, NDataTable, useMessage
} from 'naive-ui'
import TagManager from './TagManager.vue'
import {
  Add, CloudUpload, Search, Grid, List, Person, Cloud, People, 
  Time, Star, Trash, FolderOpen, Image, DocumentText, Videocam, 
  Home, Folder, Document
} from '@vicons/ionicons5'
import { GetRootDirectories, GetDirectoryContents } from '../../wailsjs/go/main/App'

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
const currentView = ref('my-drive')
const searchQuery = ref('')
const currentPath = ref('')
const files = ref<FileItem[]>([])
const selectedFile = ref<FileItem | null>(null)

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

const accountOptions = [
  { label: '账户设置', key: 'account' },
  { label: '帮助', key: 'help' },
  { label: '退出', key: 'logout' }
]

const pathSegments = computed(() => {
  return currentPath.value ? currentPath.value.split('\\').filter(Boolean) : []
})

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
      return h('div', { class: 'file-row-name' }, [
        h(NIcon, { size: 20 }, { default: () => h(row.isDirectory ? Folder : Document) }),
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
    render: (row: FileItem) => formatFileSize(row.size)
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
  await loadRootDirectories()
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

const loadDirectoryContents = async (path: string) => {
  try {
    const contents = await GetDirectoryContents(path)
    files.value = contents.map(item => ({
      ...item,
      modified: item.modified * 1000
    }))
    currentPath.value = path
  } catch (error) {
    message.error('Failed to load directory contents: ' + error)
  }
}

const handleFileClick = (file: FileItem) => {
  if (file.isDirectory) {
    loadDirectoryContents(file.path)
  } else {
    selectedFile.value = file
    message.info(`Selected file: ${file.name}`)
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
      handleFileClick(file)
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

const navigateTo = (path: string) => {
  if (path === '') {
    loadRootDirectories()
    currentPath.value = ''
  }
}

const navigateToPath = (index: number) => {
  const targetPath = pathSegments.value.slice(0, index + 1).join('\\')
  if (targetPath) {
    loadDirectoryContents(targetPath)
  } else {
    loadRootDirectories()
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

const formatDate = (timestamp?: number): string => {
  if (!timestamp) return 'Unknown'
  return new Date(timestamp).toLocaleDateString()
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
  overflow-y: auto;
}

.nav-section {
  padding: 8px 0;
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

.breadcrumb {
  padding: 16px 24px;
  background: white;
  border-bottom: 1px solid #e0e0e0;
}

.file-container {
  flex: 1;
  overflow-y: auto;
  padding: 24px;
}

.grid-view {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 16px;
}

.grid-item {
  background: white;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  padding: 16px;
  text-align: center;
  cursor: pointer;
  transition: all 0.2s;
}

.grid-item:hover {
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  transform: translateY(-2px);
}

.file-icon {
  margin-bottom: 8px;
  color: #5f6368;
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