<template>
  <div class="admin-card folder-manage">
    <h1><i class="fas fa-folder-open"></i> 文件资源管理</h1>

    <div class="folder-layout">
      <aside class="folder-sidebar">
        <div class="sidebar-header">
          <span class="sidebar-title">文件夹</span>
          <el-button size="small" text :loading="treeLoading" @click="refreshTree">
            刷新
          </el-button>
        </div>

        <el-tree
          ref="treeRef"
          class="folder-tree"
          :data="treeData"
          node-key="id"
          :props="treeProps"
          :highlight-current="true"
          :expand-on-click-node="false"
          @node-click="handleTreeNodeClick"
        >
          <template #default="{ data }">
            <span class="tree-node">
              <i class="fas fa-folder tree-icon"></i>
              <span>{{ data.name }}</span>
            </span>
          </template>
        </el-tree>
      </aside>

      <section class="folder-main">
        <div class="main-header">
          <el-breadcrumb separator="/" class="breadcrumb">
            <el-breadcrumb-item v-for="item in breadcrumb" :key="item.id">
              <a href="#" class="breadcrumb-link" @click.prevent="navigateTo(item.id)">
                {{ item.name }}
              </a>
            </el-breadcrumb-item>
          </el-breadcrumb>

          <div class="toolbar">
            <el-input
              v-model="searchQuery"
              placeholder="搜索文件/文件夹（全局）"
              clearable
              style="max-width: 280px"
              @keyup.enter="runSearch"
            />
            <el-button :loading="searchLoading" @click="runSearch">搜索</el-button>
            <el-button v-if="searchMode" @click="exitSearch">返回当前目录</el-button>

            <el-button type="primary" :disabled="searchMode" @click="promptCreateFolder">
              新建文件夹
            </el-button>

            <el-button
              v-if="selectedRows.length > 0"
              :disabled="searchMode"
              @click="openMoveDialog()"
            >
              批量移动（{{ selectedRows.length }}）
            </el-button>
            <el-button
              v-if="selectedRows.length > 0"
              :disabled="searchMode"
              type="danger"
              @click="handleBatchDelete"
            >
              批量删除（{{ selectedRows.length }}）
            </el-button>

            <el-upload
              :action="uploadUrl"
              :headers="uploadHeaders"
              :show-file-list="false"
              :before-upload="beforeUpload"
              :on-success="handleUploadSuccess"
              :on-error="handleUploadError"
              :disabled="currentFolderId === 0 || searchMode"
            >
              <el-button type="primary" plain :disabled="currentFolderId === 0 || searchMode">
                上传文件
              </el-button>
            </el-upload>
          </div>
        </div>

        <div class="table-wrap">
          <el-table
            :data="tableRows"
            class="content-table"
            height="100%"
            v-loading="contentLoading || searchLoading"
            @row-dblclick="handleRowDblClick"
            @selection-change="handleSelectionChange"
          >
            <template #empty>
              <el-empty description="暂无内容" />
            </template>

            <el-table-column type="selection" width="46" />

          <el-table-column label="名称" min-width="320">
            <template #default="{ row }">
              <div class="name-cell">
                <span :class="['icon-badge', iconBadgeClass(row)]">
                  <i :class="rowIconClass(row)"></i>
                </span>
                <div class="name-main">
                  <div class="name-text">{{ row.name }}</div>
                  <div v-if="searchMode && row.path" class="name-sub muted">{{ row.path }}</div>
                  <div v-else-if="row.type === 'file' && row.uploadedAt" class="name-sub muted">
                    上传：{{ row.uploadedAt }}
                  </div>
                </div>
              </div>
            </template>
          </el-table-column>

          <el-table-column label="类型" width="110" align="center">
            <template #default="{ row }">
              <el-tag size="small" :type="row.type === 'folder' ? 'info' : 'success'">
                {{ row.type === 'folder' ? '文件夹' : '文件' }}
              </el-tag>
            </template>
          </el-table-column>

          <el-table-column label="大小" width="150" align="right">
            <template #default="{ row }">
              <span v-if="row.type === 'file'">{{ formatFileSize(row.size) }}</span>
              <span v-else class="muted">
                {{ formatFileSize(row.size || 0) }} /
                {{ typeof row.fileCount === 'number' ? `${row.fileCount} 文件` : '-' }}
              </span>
            </template>
          </el-table-column>

          <el-table-column label="操作" width="220" align="center" fixed="right">
            <template #default="{ row }">
              <el-tooltip v-if="row.type === 'folder'" content="打开" placement="top">
                <el-button size="small" circle @click="openFolder(row.id)">
                  <i class="fas fa-folder-open"></i>
                </el-button>
              </el-tooltip>
              <el-tooltip v-if="row.type === 'file'" content="下载" placement="top">
                <el-button size="small" circle type="success" @click="downloadFile(row.id)">
                  <i class="fas fa-download"></i>
                </el-button>
              </el-tooltip>
              <el-tooltip content="移动" placement="top">
                <el-button size="small" circle @click="openMoveDialog(row)">
                  <i class="fas fa-arrows-alt"></i>
                </el-button>
              </el-tooltip>
              <el-tooltip content="删除" placement="top">
                <el-button size="small" circle type="danger" @click="handleDelete(row)">
                  <i class="fas fa-trash-alt"></i>
                </el-button>
              </el-tooltip>
            </template>
          </el-table-column>
          </el-table>
        </div>
      </section>
    </div>

    <el-dialog v-model="moveDialogVisible" title="移动到..." width="560px">
      <div class="move-dialog">
        <el-alert
          type="info"
          :closable="false"
          show-icon
          title="选择目标文件夹（文件不能移动到根目录）"
        />

        <el-tree
          ref="moveTreeRef"
          class="move-tree"
          :data="moveTreeData"
          node-key="id"
          :props="treeProps"
          :highlight-current="true"
          default-expand-all
          :expand-on-click-node="false"
          @node-click="handleMoveTargetClick"
        >
          <template #default="{ data }">
            <span class="tree-node">
              <i class="fas fa-folder tree-icon"></i>
              <span>{{ data.name }}</span>
            </span>
          </template>
        </el-tree>
      </div>

      <template #footer>
        <el-button @click="moveDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="moveSubmitting" @click="confirmMove">
          确定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/request'
import { getToken } from '@/utils/tokenManager'

export default {
  name: 'FolderManage',
  data() {
    return {
      treeProps: { label: 'name', children: 'children' },
      treeData: [],
      treeLoading: false,

      currentFolderId: 0,
      currentPath: '/',
      breadcrumb: [],
      subFolders: [],
      files: [],
      contentLoading: false,

      searchQuery: '',
      searchMode: false,
      searchResults: [],
      searchLoading: false,

      selectedRows: [],

      moveDialogVisible: false,
      moveSubmitting: false,
      moveTargetId: null,
      moveSources: []
    }
  },
  computed: {
    uploadUrl() {
      return `${request.defaults.baseURL}/admin/folder/${this.currentFolderId}/upload`
    },
    uploadHeaders() {
      const token = getToken()
      return token ? { Authorization: `Bearer ${token}` } : {}
    },
    moveTreeData() {
      return [{ id: 0, name: '根目录', children: this.treeData }]
    },
    tableRows() {
      if (this.searchMode) {
        return this.searchResults.map(r => ({
          ...r,
          type: r.type,
          uploadedAt: r.uploadedAt,
          fileCount: null
        }))
      }

      const folderRows = this.subFolders.map(f => ({
        ...f,
        type: 'folder'
      }))

      const fileRows = this.files.map(f => ({
        ...f,
        type: 'file',
        path: `${this.currentPath.replace(/\/$/, '')}/${f.name}`
      }))

      return [...folderRows, ...fileRows]
    }
  },
  async mounted() {
    await this.refreshTree()
    await this.setCurrentFolder(0)
  },
  methods: {
    async refreshTree() {
      this.treeLoading = true
      try {
        const res = await request.get('/admin/folder/tree')
        this.treeData = res.data.tree || []
      } catch (err) {
        ElMessage.error(err.response?.data?.message || '获取文件夹树失败')
      } finally {
        this.treeLoading = false
      }
    },
    async setCurrentFolder(folderId) {
      this.searchMode = false
      this.currentFolderId = folderId
      this.selectedRows = []
      await Promise.all([this.fetchBreadcrumb(folderId), this.fetchContent(folderId)])

      if (this.$refs.treeRef && folderId !== 0) {
        this.$refs.treeRef.setCurrentKey(folderId)
      }
    },
    async fetchBreadcrumb(folderId) {
      try {
        const res = await request.get(`/admin/folder/${folderId}/breadcrumb`)
        this.breadcrumb = res.data.breadcrumb || [{ id: 0, name: '根目录' }]
      } catch (err) {
        this.breadcrumb = [{ id: 0, name: '根目录' }]
      }
    },
    async fetchContent(folderId) {
      this.contentLoading = true
      try {
        const res = await request.get(`/admin/folder/${folderId}/content`)
        const data = res.data

        this.currentPath = data.folder?.path || '/'
        this.subFolders = data.subFolders || []
        this.files = data.files || []
      } catch (err) {
        ElMessage.error(err.response?.data?.message || '获取目录内容失败')
        this.currentPath = '/'
        this.subFolders = []
        this.files = []
      } finally {
        this.contentLoading = false
      }
    },
    handleTreeNodeClick(node) {
      this.setCurrentFolder(node.id)
    },
    navigateTo(folderId) {
      this.setCurrentFolder(folderId)
    },
    openFolder(folderId) {
      this.setCurrentFolder(folderId)
    },
    handleRowDblClick(row) {
      if (row.type === 'folder') {
        this.openFolder(row.id)
      } else {
        this.downloadFile(row.id)
      }
    },
    async promptCreateFolder() {
      try {
        const { value } = await ElMessageBox.prompt('请输入新文件夹名称', '新建文件夹', {
          confirmButtonText: '创建',
          cancelButtonText: '取消',
          inputPattern: /^.{1,50}$/,
          inputErrorMessage: '名称不能为空，且不能超过 50 个字符'
        })

        const parentId = this.currentFolderId === 0 ? null : this.currentFolderId
        await request.post('/admin/folder/create', { name: value, parentId })
        ElMessage.success('创建成功')
        await this.refreshTree()
        await this.fetchContent(this.currentFolderId)
      } catch (err) {
        if (err !== 'cancel') {
          ElMessage.error(err.response?.data?.message || '创建失败')
        }
      }
    },
    beforeUpload(file) {
      const validTypes = [
        'application/pdf',
        'application/msword',
        'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
        'application/vnd.ms-powerpoint',
        'application/vnd.openxmlformats-officedocument.presentationml.presentation',
        'application/vnd.ms-excel',
        'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
        'image/jpeg',
        'image/png',
        'image/gif',
        'application/zip',
        'application/x-rar-compressed',
        'application/x-zip-compressed',
        'application/octet-stream',
        'text/markdown',
        'text/x-markdown'
      ]

      const validExtensions = [
        '.pdf',
        '.doc',
        '.docx',
        '.ppt',
        '.pptx',
        '.xls',
        '.xlsx',
        '.jpg',
        '.jpeg',
        '.png',
        '.gif',
        '.zip',
        '.rar',
        '.md'
      ]

      const fileExtension = file.name.slice(file.name.lastIndexOf('.')).toLowerCase()
      const isTypeValid = validTypes.includes(file.type) || validExtensions.includes(fileExtension)
      const isLt50M = file.size / 1024 / 1024 < 50

      if (this.currentFolderId === 0) {
        ElMessage.warning('请先选择一个文件夹再上传文件')
        return false
      }
      if (!isTypeValid) {
        ElMessage.error('仅支持 PDF/Office/图片/压缩包/Markdown 文件')
      }
      if (!isLt50M) {
        ElMessage.error('文件大小不能超过 50MB')
      }
      return isTypeValid && isLt50M
    },
    async handleUploadSuccess() {
      ElMessage.success('上传成功')
      await this.fetchContent(this.currentFolderId)
      await this.refreshTree()
    },
    handleUploadError(error) {
      ElMessage.error('上传失败: ' + (error.message || '未知错误'))
      console.error(error)
    },
    async downloadFile(fileId) {
      try {
        const res = await request.get(`/folder/download/${fileId}`)
        const { name, type, data } = res.data

        const binary = atob(data)
        const bytes = new Uint8Array(binary.length)
        for (let i = 0; i < binary.length; i++) bytes[i] = binary.charCodeAt(i)

        const blob = new Blob([bytes], { type: type || 'application/octet-stream' })
        const url = URL.createObjectURL(blob)

        const link = document.createElement('a')
        link.href = url
        link.download = name || 'download'
        document.body.appendChild(link)
        link.click()
        link.remove()

        URL.revokeObjectURL(url)
      } catch (err) {
        ElMessage.error(err.response?.data?.message || '下载失败')
      }
    },
    async handleDelete(row) {
      const isFolder = row.type === 'folder'
      const title = isFolder ? '删除文件夹' : '删除文件'
      const message = isFolder
        ? '确定删除该文件夹及其子文件夹/文件吗？'
        : '确定删除该文件吗？'

      try {
        await ElMessageBox.confirm(message, title, {
          confirmButtonText: '确定删除',
          cancelButtonText: '取消',
          type: 'warning'
        })

        if (isFolder) {
          await request.delete(`/admin/folder/delete/${row.id}`)
        } else {
          await request.delete(`/admin/folder/file/${row.id}`)
        }

        ElMessage.success('删除成功')
        await this.refreshTree()
        if (!this.searchMode) {
          await this.fetchContent(this.currentFolderId)
        } else {
          await this.runSearch()
        }
      } catch (err) {
        if (err !== 'cancel') {
          ElMessage.error(err.response?.data?.message || '删除失败')
        }
      }
    },
    async runSearch() {
      const q = (this.searchQuery || '').trim()
      if (!q) {
        ElMessage.warning('请输入搜索关键字')
        return
      }

      this.searchLoading = true
      try {
        const res = await request.get('/admin/search', { params: { q } })
        this.searchResults = res.data.results || []
        this.searchMode = true
      } catch (err) {
        ElMessage.error(err.response?.data?.message || '搜索失败')
      } finally {
        this.searchLoading = false
      }
    },
    exitSearch() {
      this.searchMode = false
      this.searchResults = []
      this.selectedRows = []
    },
    handleSelectionChange(rows) {
      this.selectedRows = rows || []
    },
    openMoveDialog(row) {
      const sources = row ? [row] : this.selectedRows
      if (!sources || sources.length === 0) {
        ElMessage.warning('请选择要移动的项目')
        return
      }
      this.moveSources = sources
      this.moveTargetId = null
      this.moveDialogVisible = true
    },
    handleMoveTargetClick(node) {
      this.moveTargetId = node.id
    },
    async confirmMove() {
      if (!this.moveSources || this.moveSources.length === 0) return
      if (this.moveTargetId == null) {
        ElMessage.warning('请选择目标文件夹')
        return
      }

      const hasFile = this.moveSources.some(s => s.type === 'file')
      if (hasFile && this.moveTargetId === 0) {
        ElMessage.warning('文件不能移动到根目录，请选择一个文件夹')
        return
      }

      this.moveSubmitting = true
      try {
        for (const source of this.moveSources) {
          await request.post('/admin/item/move', {
            type: source.type,
            id: source.id,
            targetFolderId: this.moveTargetId === 0 ? null : this.moveTargetId
          })
        }

        ElMessage.success('移动成功')
        this.moveDialogVisible = false
        this.selectedRows = []
        await this.refreshTree()
        if (!this.searchMode) {
          await this.fetchContent(this.currentFolderId)
        } else {
          await this.runSearch()
        }
      } catch (err) {
        ElMessage.error(err.response?.data?.message || '移动失败')
      } finally {
        this.moveSubmitting = false
      }
    },
    async handleBatchDelete() {
      if (!this.selectedRows || this.selectedRows.length === 0) {
        ElMessage.warning('请选择要删除的项目')
        return
      }

      try {
        await ElMessageBox.confirm(
          `确定删除选中的 ${this.selectedRows.length} 项吗？（文件夹会递归删除子内容）`,
          '批量删除',
          {
            confirmButtonText: '确定删除',
            cancelButtonText: '取消',
            type: 'warning'
          }
        )

        for (const row of this.selectedRows) {
          if (row.type === 'folder') {
            await request.delete(`/admin/folder/delete/${row.id}`)
          } else {
            await request.delete(`/admin/folder/file/${row.id}`)
          }
        }

        ElMessage.success('删除成功')
        this.selectedRows = []
        await this.refreshTree()
        await this.fetchContent(this.currentFolderId)
      } catch (err) {
        if (err !== 'cancel') {
          ElMessage.error(err.response?.data?.message || '批量删除失败')
        }
      }
    },
    rowIconClass(row) {
      if (row.type === 'folder') return 'fas fa-folder'
      return this.getFileIcon(row.name)
    },
    iconBadgeClass(row) {
      if (row.type === 'folder') return 'icon-badge--folder'
      const name = (row.name || '').toLowerCase()
      if (name.endsWith('.pdf')) return 'icon-badge--pdf'
      if (name.endsWith('.doc') || name.endsWith('.docx')) return 'icon-badge--word'
      if (name.endsWith('.ppt') || name.endsWith('.pptx')) return 'icon-badge--ppt'
      if (name.endsWith('.xls') || name.endsWith('.xlsx')) return 'icon-badge--excel'
      if (name.endsWith('.zip') || name.endsWith('.rar')) return 'icon-badge--archive'
      if (name.endsWith('.png') || name.endsWith('.jpg') || name.endsWith('.jpeg') || name.endsWith('.gif')) return 'icon-badge--image'
      if (name.endsWith('.md')) return 'icon-badge--md'
      return 'icon-badge--file'
    },
    getFileIcon(filename) {
      const name = (filename || '').toLowerCase()
      if (name.endsWith('.pdf')) return 'fas fa-file-pdf'
      if (name.endsWith('.doc') || name.endsWith('.docx')) return 'fas fa-file-word'
      if (name.endsWith('.ppt') || name.endsWith('.pptx')) return 'fas fa-file-powerpoint'
      if (name.endsWith('.xls') || name.endsWith('.xlsx')) return 'fas fa-file-excel'
      if (name.endsWith('.zip') || name.endsWith('.rar')) return 'fas fa-file-archive'
      if (name.endsWith('.png') || name.endsWith('.jpg') || name.endsWith('.jpeg') || name.endsWith('.gif')) return 'fas fa-file-image'
      if (name.endsWith('.md')) return 'fas fa-file-alt'
      return 'fas fa-file'
    },
    formatFileSize(bytes) {
      if (!bytes) return '0 Bytes'
      const k = 1024
      const sizes = ['Bytes', 'KB', 'MB', 'GB']
      const i = Math.floor(Math.log(bytes) / Math.log(k))
      const size = (bytes / Math.pow(k, i)).toFixed(2)
      return `${size} ${sizes[i]}`
    }
  }
}
</script>

<style scoped>
h1 {
  font-size: 24px;
  color: #303133;
  margin-bottom: 24px;
  font-weight: 600;
}

h1 i {
  margin-right: 10px;
  color: #409eff;
}

.folder-layout {
  display: flex;
  gap: var(--spacing-5);
  min-height: calc(100vh - 260px);
}

.folder-sidebar,
.folder-main {
  background-color: rgba(255, 255, 255, 0.75);
  border: 1px solid rgba(255, 255, 255, 0.5);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-sm);
  min-width: 0;
}

.folder-sidebar {
  flex: 0 0 clamp(240px, 24vw, 360px);
  padding: var(--spacing-4);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.sidebar-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--spacing-3);
}

.sidebar-title {
  font-weight: var(--font-weight-semibold);
  color: var(--color-gray-800);
}

.folder-tree {
  flex: 1;
  min-height: 0;
  overflow: auto;
}

.tree-node {
  display: inline-flex;
  align-items: center;
  gap: var(--spacing-2);
}

.tree-icon {
  color: var(--color-warning);
}

.folder-main {
  flex: 1 1 auto;
  padding: var(--spacing-4);
  display: flex;
  flex-direction: column;
  gap: var(--spacing-4);
  overflow: hidden;
}

.main-header {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-3);
}

.breadcrumb {
  font-size: var(--font-sm);
}

.breadcrumb-link {
  color: var(--color-primary);
  text-decoration: none;
}

.breadcrumb-link:hover {
  text-decoration: underline;
}

.toolbar {
  display: flex;
  flex-wrap: wrap;
  gap: var(--spacing-3);
  align-items: center;
  justify-content: flex-start;
}

.content-table {
  border-radius: var(--radius-md);
  overflow: hidden;
  height: 100%;
}

.table-wrap {
  flex: 1;
  min-height: 0;
}

.name-cell {
  display: inline-flex;
  align-items: center;
  gap: var(--spacing-2);
}

.name-main {
  display: flex;
  flex-direction: column;
  min-width: 0;
  gap: 2px;
}

.name-sub {
  font-size: var(--font-xs);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 52ch;
}

.muted {
  color: var(--color-gray-500);
}

.icon-badge {
  width: 34px;
  height: 34px;
  border-radius: var(--radius-md);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  flex: 0 0 auto;
  color: var(--color-white);
}

.icon-badge--folder {
  background: linear-gradient(135deg, #f59e0b, #fbbf24);
}
.icon-badge--pdf {
  background: linear-gradient(135deg, #ef4444, #f87171);
}
.icon-badge--word {
  background: linear-gradient(135deg, #2563eb, #60a5fa);
}
.icon-badge--ppt {
  background: linear-gradient(135deg, #f97316, #fb923c);
}
.icon-badge--excel {
  background: linear-gradient(135deg, #16a34a, #4ade80);
}
.icon-badge--archive {
  background: linear-gradient(135deg, #7c3aed, #a78bfa);
}
.icon-badge--image {
  background: linear-gradient(135deg, #06b6d4, #22c55e);
}
.icon-badge--md {
  background: linear-gradient(135deg, #334155, #64748b);
}
.icon-badge--file {
  background: linear-gradient(135deg, #475569, #94a3b8);
}

.move-dialog {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-3);
}

.move-tree {
  max-height: 360px;
  overflow: auto;
  border: 1px solid var(--color-gray-200);
  border-radius: var(--radius-md);
  padding: var(--spacing-2);
}

@media (max-width: 1024px) {
  .folder-layout {
    flex-direction: column;
    min-height: auto;
  }

  .folder-tree {
    max-height: 320px;
  }
}
</style>
