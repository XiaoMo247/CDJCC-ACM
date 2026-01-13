<template>
  <div class="knowledge-page page-knowledge">
    <el-dialog v-model="blogDialogVisible" title="团队成员博客" width="820px" class="blog-dialog">
      <div class="blog-list">
        <a v-for="(blog, index) in teamBlogs" :key="index" class="blog-item" :href="blog.url" target="_blank"
          rel="noreferrer">
          <img :src="blog.avatar" alt="Avatar" class="blog-avatar" />
          <div class="blog-info">
            <div class="blog-author">{{ blog.author }}</div>
            <div class="blog-title">{{ blog.title }}</div>
          </div>
        </a>
      </div>
    </el-dialog>

    <header class="page-header page-hero">
      <h1 class="page-title">
        <span class="page-title-icon"><i class="fas fa-book"></i></span>
        知识库
      </h1>
      <p class="page-subtitle">课程资料、课件资源与常用链接</p>
    </header>

    <section class="quick-links">
      <div class="link-card link-card--blog" role="button" tabindex="0" @click="openBlogDialog" @keydown.enter="openBlogDialog">
        <div class="link-icon"><i class="fas fa-blog"></i></div>
        <div class="link-main">
          <div class="link-title">团队博客</div>
          <div class="link-desc">查看队员的学习记录与题解总结</div>
        </div>
      </div>

      <a class="link-card link-card--bilibili" :href="biliBiliUrl"
        target="_blank" rel="noreferrer">
        <div class="link-icon"><i class="fab fa-bilibili"></i></div>
        <div class="link-main">
          <div class="link-title">B 站频道</div>
          <div class="link-desc">团队公开视频与训练分享</div>
        </div>
      </a>

      <a class="link-card link-card--oj" :href="ojUrl" target="_blank" rel="noreferrer">
        <div class="link-icon"><i class="fas fa-laptop-code"></i></div>
        <div class="link-main">
          <div class="link-title">OJ 平台</div>
          <div class="link-desc">在线训练与题目集合</div>
        </div>
      </a>
    </section>

    <section class="courseware">
      <div class="section-title">
        <i class="fas fa-folder-open"></i>
        <span>课件资源</span>
      </div>

      <div class="explorer">
        <aside class="sidebar">
          <div class="sidebar-header">
            <span class="sidebar-title">目录</span>
            <el-button size="small" text :loading="treeLoading" @click="fetchTree">刷新</el-button>
          </div>

          <el-tree ref="treeRef" class="tree" :data="treeData" node-key="id" :props="treeProps"
            :highlight-current="true" :expand-on-click-node="false" default-expand-all @node-click="handleTreeClick">
            <template #default="{ data }">
              <span class="tree-node">
                <i class="fas fa-folder tree-icon"></i>
                <span class="tree-label">{{ data.name }}</span>
              </span>
            </template>
          </el-tree>
        </aside>

        <main class="main">
          <div class="main-header">
            <el-breadcrumb separator="/" class="breadcrumb">
              <el-breadcrumb-item v-for="item in breadcrumb" :key="item.id">
                <a href="#" class="breadcrumb-link" @click.prevent="openFolder(item.id)">{{ item.name }}</a>
              </el-breadcrumb-item>
            </el-breadcrumb>

            <div class="toolbar">
              <el-input v-model="filterQuery" placeholder="筛选当前目录..." clearable style="max-width: 260px" />

              <el-input v-model="globalSearchQuery" placeholder="全局搜索..." clearable style="max-width: 260px"
                @keyup.enter="runGlobalSearch" />
              <el-button :loading="searchLoading" @click="runGlobalSearch">搜索</el-button>
              <el-button v-if="searchMode" @click="exitSearch">返回目录</el-button>
            </div>
          </div>

          <div class="content">
            <el-table :data="displayRows" stripe class="table" height="100%"
              v-loading="contentLoading || treeLoading || searchLoading" @row-dblclick="handleRowDblClick">
              <el-table-column label="名称" min-width="320">
                <template #default="{ row }">
                  <div class="name-cell">
                    <span :class="['icon-badge', iconBadgeClass(row)]">
                      <i :class="rowIconClass(row)"></i>
                    </span>
                    <span class="name-text">{{ row.name }}</span>
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

              <el-table-column label="大小" width="170" align="right">
                <template #default="{ row }">
                  <span v-if="row.type === 'file'">{{ formatFileSize(row.size) }}</span>
                  <span v-else class="muted">
                    {{ formatFileSize(row.size || 0) }} · {{ typeof row.fileCount === 'number' ? `${row.fileCount} 文件` :
                      '-' }}
                  </span>
                </template>
              </el-table-column>

              <el-table-column label="路径" min-width="220">
                <template #default="{ row }">
                  <span class="muted">{{ searchMode ? row.path : currentPath }}</span>
                </template>
              </el-table-column>

              <el-table-column label="上传时间" width="180">
                <template #default="{ row }">
                  <span v-if="row.type === 'file'">{{ row.uploadedAt || '-' }}</span>
                  <span v-else class="muted">-</span>
                </template>
              </el-table-column>

              <el-table-column label="操作" width="260" align="center">
                <template #default="{ row }">
                  <el-button v-if="row.type === 'folder'" size="small" @click="openFolder(row.id)">打开</el-button>
                  <el-button v-if="row.type === 'file'" size="small" type="primary" @click="downloadFile(row.id)">
                    下载
                  </el-button>
                  <el-button v-if="searchMode && row.type === 'file' && row.folderId" size="small"
                    @click="openFolder(row.folderId)">
                    打开所在文件夹
                  </el-button>
                </template>
              </el-table-column>
            </el-table>

            <el-empty v-if="!contentLoading && !searchLoading && displayRows.length === 0" description="暂无内容" />
          </div>
        </main>
      </div>
    </section>
  </div>
</template>

<script>
import { ElMessage } from 'element-plus'
import request from '@/utils/request'

export default {
  name: 'KnowledgePage',
  data() {
    return {
      blogDialogVisible: false,
      // 外部链接配置
      biliBiliUrl: import.meta.env.VITE_BILIBILI_URL || 'https://space.bilibili.com/3546651937475184?spm_id_from=333.337.search-card.all.click',
      ojUrl: import.meta.env.VITE_OJ_URL || 'https://hydro.ac/d/cdjcc_acm_2333/',

      teamBlogs: [
        {
          author: 'Martian148',
          title: '锦城学院 ACM 学习地图',
          url: 'https://www.cnblogs.com/martian148/p/18221024',
          avatar: 'https://cdn.luogu.com.cn/upload/image_hosting/npue5pns.png'
        },
        {
          author: 'Martian148',
          title: '个人博客',
          url: 'https://martian148.xyz/',
          avatar: 'https://cdn.luogu.com.cn/upload/image_hosting/npue5pns.png'
        },
        {
          author: 'XiaoMo247',
          title: '个人博客',
          url: 'https://www.cnblogs.com/XiaoMo247',
          avatar: 'https://cdn.luogu.com.cn/upload/image_hosting/5of8epxu.png'
        },
        {
          author: 'Martian148',
          title: 'XCPC 自学指北',
          url: 'https://chengyiwei.github.io/mkdocs/',
          avatar: 'https://cdn.luogu.com.cn/upload/image_hosting/npue5pns.png'
        }
      ],

      // Explorer state
      treeProps: { label: 'name', children: 'children' },
      treeData: [],
      treeLoading: false,

      currentFolderId: 0,
      currentPath: '/',
      breadcrumb: [{ id: 0, name: '根目录' }],

      subFolders: [],
      files: [],
      contentLoading: false,

      filterQuery: '',
      globalSearchQuery: '',
      searchMode: false,
      searchResults: [],
      searchLoading: false,

      contentCache: {}
    }
  },
  computed: {
    displayRows() {
      let rows = []

      if (this.searchMode) {
        rows = (this.searchResults || []).map(r => ({
          id: r.id,
          name: r.name,
          type: r.type,
          path: r.path,
          size: r.size || 0,
          folderId: r.folderId || null
        }))
      } else {
        const folderRows = (this.subFolders || []).map(f => ({ ...f, type: 'folder' }))
        const fileRows = (this.files || []).map(f => ({
          ...f,
          type: 'file',
          path: `${this.currentPath.replace(/\/$/, '')}/${f.name}`
        }))
        rows = [...folderRows, ...fileRows]
      }

      const q = (this.filterQuery || '').trim().toLowerCase()
      if (!q) return rows
      return rows.filter(r => (r.name || '').toLowerCase().includes(q))
    }
  },
  mounted() {
    this.initExplorer()
  },
  methods: {
    openBlogDialog() {
      this.blogDialogVisible = true
    },
    async initExplorer() {
      await this.fetchTree()
      await this.openFolder(0)
    },
    async fetchTree() {
      this.treeLoading = true
      try {
        const res = await request.get('/folder/tree')
        const tree = res.data.tree || []
        this.treeData = [{ id: 0, name: '根目录', children: tree }]
      } catch (error) {
        ElMessage.error(error.response?.data?.message || '获取文件夹树失败')
        console.error(error)
      } finally {
        this.treeLoading = false
      }
    },
    handleTreeClick(node) {
      this.openFolder(node.id)
    },
    async openFolder(folderId) {
      this.searchMode = false
      this.searchResults = []
      this.currentFolderId = folderId

      await Promise.all([this.fetchBreadcrumb(folderId), this.fetchContent(folderId)])

      if (this.$refs.treeRef) {
        this.$refs.treeRef.setCurrentKey(folderId)
      }
    },
    async fetchBreadcrumb(folderId) {
      try {
        const res = await request.get(`/folder/${folderId}/breadcrumb`)
        this.breadcrumb = res.data.breadcrumb || [{ id: 0, name: '根目录' }]
      } catch (error) {
        this.breadcrumb = [{ id: 0, name: '根目录' }]
      }
    },
    async fetchContent(folderId) {
      const cacheKey = String(folderId)
      if (this.contentCache[cacheKey]) {
        const cached = this.contentCache[cacheKey]
        this.currentPath = cached.currentPath
        this.subFolders = cached.subFolders
        this.files = cached.files
        return
      }

      this.contentLoading = true
      try {
        const res = await request.get(`/folder/${folderId}/content`)
        const data = res.data

        this.currentPath = data.folder?.path || '/'
        this.subFolders = data.subFolders || []
        this.files = data.files || []

        this.contentCache[cacheKey] = {
          currentPath: this.currentPath,
          subFolders: this.subFolders,
          files: this.files
        }
      } catch (error) {
        ElMessage.error(error.response?.data?.message || '获取目录内容失败')
        console.error(error)
        this.currentPath = '/'
        this.subFolders = []
        this.files = []
      } finally {
        this.contentLoading = false
      }
    },
    async runGlobalSearch() {
      const q = (this.globalSearchQuery || '').trim()
      if (!q) {
        ElMessage.warning('请输入搜索关键字')
        return
      }

      this.searchLoading = true
      try {
        const res = await request.get('/folder/search', { params: { q } })
        this.searchResults = res.data.results || []
        this.searchMode = true
      } catch (error) {
        ElMessage.error(error.response?.data?.message || '搜索失败')
        console.error(error)
      } finally {
        this.searchLoading = false
      }
    },
    exitSearch() {
      this.searchMode = false
      this.searchResults = []
    },
    handleRowDblClick(row) {
      if (row.type === 'folder') {
        this.openFolder(row.id)
      } else {
        this.downloadFile(row.id)
      }
    },
    async downloadFile(fileId) {
      try {
        const res = await request.get(`/folder/download/${fileId}`)
        const { name, type, data } = res.data

        const byteCharacters = atob(data)
        const bytes = new Uint8Array(byteCharacters.length)
        for (let i = 0; i < byteCharacters.length; i++) bytes[i] = byteCharacters.charCodeAt(i)

        const blob = new Blob([bytes], { type: type || 'application/octet-stream' })
        const url = URL.createObjectURL(blob)

        const link = document.createElement('a')
        link.href = url
        link.download = name || 'download'
        document.body.appendChild(link)
        link.click()
        link.remove()

        URL.revokeObjectURL(url)
      } catch (error) {
        ElMessage.error(error.response?.data?.message || '下载失败')
        console.error(error)
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
.knowledge-page {
  max-width: 1200px;
  margin: 0 auto;
  padding: var(--spacing-6) var(--spacing-4);
  display: flex;
  flex-direction: column;
  gap: var(--spacing-6);
}

.page-header {
  text-align: center;
  padding: 60px 20px;
  background: linear-gradient(135deg, #c7d2fe 0%, #a5b4fc 100%);
  color: #312e81;
  border-radius: 12px;
  margin-bottom: 60px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.page-title {
  margin: 0 0 15px 0;
  font-size: clamp(2rem, 3vw, 2.8rem);
  font-weight: 700;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 15px;
}

.page-title i {
  font-size: 2.5rem;
}

.page-subtitle {
  margin: 10px 0 0 0;
  color: rgba(255, 255, 255, 0.95);
  font-size: 1.125rem;
}

.quick-links {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: var(--spacing-4);
}

.link-card {
  display: flex;
  gap: var(--spacing-3);
  align-items: center;
  padding: var(--spacing-4);
  border-radius: var(--radius-lg);
  background: rgba(255, 255, 255, 0.75);
  border: 1px solid rgba(255, 255, 255, 0.5);
  box-shadow: var(--shadow-sm);
  text-decoration: none;
  color: inherit;
  transition: transform var(--duration-base) var(--ease-out), box-shadow var(--duration-base) var(--ease-out);
}

.link-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}

.link-icon {
  width: 44px;
  height: 44px;
  border-radius: var(--radius-md);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--color-white);
  flex: 0 0 auto;
}

.link-card--blog .link-icon {
  background: linear-gradient(135deg, #4e54c8, #8f94fb);
  box-shadow: 0 10px 25px rgba(79, 70, 229, 0.25);
}

.link-card--bilibili .link-icon {
  background: linear-gradient(135deg, #00a1d6, #0066cc);
  box-shadow: 0 10px 25px rgba(0, 161, 214, 0.22);
}

.link-card--oj .link-icon {
  background: linear-gradient(135deg, #10b981, #3b82f6);
  box-shadow: 0 10px 25px rgba(59, 130, 246, 0.18);
}

.link-title {
  font-weight: var(--font-weight-semibold);
  color: var(--color-gray-900);
}

.link-desc {
  color: var(--color-gray-600);
  font-size: var(--font-sm);
  margin-top: 2px;
}

.courseware {
  background: rgba(255, 255, 255, 0.75);
  border: 1px solid rgba(255, 255, 255, 0.5);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-sm);
  padding: var(--spacing-5);
}

.section-title {
  display: flex;
  align-items: center;
  gap: var(--spacing-3);
  font-size: var(--font-xl);
  font-weight: var(--font-weight-semibold);
  color: var(--color-gray-900);
  margin-bottom: var(--spacing-4);
}

.explorer {
  display: flex;
  gap: var(--spacing-5);
  min-height: 520px;
}

.sidebar {
  width: clamp(240px, 24vw, 320px);
  flex: 0 0 auto;
  background: rgba(255, 255, 255, 0.65);
  border: 1px solid rgba(255, 255, 255, 0.5);
  border-radius: var(--radius-lg);
  padding: var(--spacing-4);
  display: flex;
  flex-direction: column;
  min-width: 0;
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

.tree {
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

.tree-label {
  color: var(--color-gray-800);
}

.main {
  flex: 1 1 auto;
  min-width: 0;
  background: rgba(255, 255, 255, 0.65);
  border: 1px solid rgba(255, 255, 255, 0.5);
  border-radius: var(--radius-lg);
  padding: var(--spacing-4);
  display: flex;
  flex-direction: column;
  gap: var(--spacing-4);
}

.main-header {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-3);
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
}

.content {
  flex: 1;
  min-height: 0;
}

.table {
  border-radius: var(--radius-md);
  overflow: hidden;
}

.name-cell {
  display: inline-flex;
  align-items: center;
  gap: var(--spacing-2);
}

.name-icon {
  width: 18px;
  text-align: center;
  color: var(--color-gray-600);
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

.muted {
  color: var(--color-gray-500);
}

.blog-dialog :deep(.el-dialog) {
  border-radius: var(--radius-lg);
}

.blog-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-3);
}

.blog-item {
  display: flex;
  gap: var(--spacing-3);
  align-items: center;
  padding: var(--spacing-3);
  border-radius: var(--radius-md);
  text-decoration: none;
  color: inherit;
  border: 1px solid var(--color-gray-200);
  background: var(--color-white);
  transition: background-color var(--duration-base) var(--ease-out);
}

.blog-item:hover {
  background: var(--color-gray-50);
}

.blog-avatar {
  width: 44px;
  height: 44px;
  border-radius: var(--radius-full);
  object-fit: cover;
}

.blog-author {
  font-weight: var(--font-weight-semibold);
  color: var(--color-gray-900);
}

.blog-title {
  color: var(--color-gray-600);
  font-size: var(--font-sm);
}

@media (max-width: 1024px) {
  .quick-links {
    grid-template-columns: 1fr;
  }

  .explorer {
    flex-direction: column;
  }

  .sidebar {
    width: 100%;
  }
}
</style>
