<template>
  <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.5.0/css/all.min.css">
  <div class="knowledge-container">

    <!-- 博客弹窗 -->
    <el-dialog v-model="blogDialogVisible" title="团队成员博客" width="90%" class="blog-dialog">
      <ul class="blog-list">
        <li class="blog-item" v-for="(blog, index) in teamBlogs" :key="index">
          <img :src="blog.avatar" alt="Avatar" class="blog-avatar">
          <div class="blog-info">
            <span class="author">{{ blog.author }}</span>
            <a :href="blog.url" target="_blank" class="blog-link">{{ blog.title }}</a>
          </div>
        </li>
      </ul>
    </el-dialog>

    <!-- 知识库标题 -->
    <h1 class="knowledge-title">📖 知识库</h1>

    <!-- 垂直布局导航 -->
    <div class="vertical-links">
      <!-- 团队博客卡片，点击打开博客弹窗 -->
      <div class="link-card" @click="openBlogDialog">
        <div class="link-icon blog">
          <i class="fas fa-blog"></i>
        </div>
        <h3>团队博客</h3>
        <p>点击查看成员技术博客</p>
      </div>

      <!-- B站频道卡片，点击在新窗口打开B站链接 -->
      <a href="https://space.bilibili.com/3546651937475184?spm_id_from=333.337.search-card.all.click"
         target="_blank" class="link-card">
        <div class="link-icon bilibili">
          <i class="fab fa-bilibili"></i>
        </div>
        <h3>B站频道</h3>
        <p>观看团队视频教程</p>
      </a>

      <!-- OJ平台卡片，点击在新窗口打开OJ平台链接 -->
      <a href="https://hydro.ac/d/cdjcc_acm_2333/" target="_blank" class="link-card">
        <div class="link-icon oj">
          <i class="fas fa-laptop-code"></i>
        </div>
        <h3>OJ平台</h3>
        <p>在线编程训练</p>
      </a>
    </div>

    <!-- 课件资源部分 -->
    <div class="courseware-section">
      <!-- 课件资源标题 -->
      <h2>📁 课件资源</h2>

      <!-- 面包屑导航 -->
      <div class="breadcrumb">
        <el-breadcrumb separator="/">
          <!-- 返回根目录 -->
          <el-breadcrumb-item @click="goToRoot">
            <i class="fas fa-home"></i> 根目录
          </el-breadcrumb-item>
          <!-- 动态渲染面包屑路径 -->
          <el-breadcrumb-item v-for="(folder, index) in breadcrumbFolders" :key="index"
                              @click="goToBreadcrumb(index)">
            {{ folder.name }}
          </el-breadcrumb-item>
        </el-breadcrumb>
      </div>

      <!-- 搜索框 -->
      <div class="search-bar">
        <el-input v-model="searchQuery" placeholder="搜索课件..." clearable size="large" style="width: 100%"
                  @clear="resetSearch">
          <template #prefix>
            <i class="el-icon-search"></i>
          </template>
        </el-input>
      </div>

      <!-- 文件夹列表 -->
      <el-table :data="filteredFolders" style="width: 100%; margin-bottom: 20px;" stripe v-loading="loading"
                empty-text="暂无文件夹" @row-click="openFolder" class="clickable-table">
        <el-table-column prop="name" label="文件夹名称">
          <template #default="scope">
            <div class="folder-item">
              <!-- 文件夹图标 -->
              <i class="fas fa-folder" style="color: #FFD700; margin-right: 10px;"></i>
              <!-- 文件夹名称 -->
              <span>{{ scope.row.name }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="220" align="center">
          <template #default="scope">
                        <span class="create-time">
                            <!-- 日历图标 -->
                            <i class="far fa-calendar-alt" style="margin-right: 5px;"></i>
                          <!-- 格式化后的创建时间 -->
                            {{ formatDate(scope.row.created_at) }}
                        </span>
          </template>
        </el-table-column>
      </el-table>

      <!-- 文件列表 -->
      <el-table :data="filteredFiles" style="width: 100%" stripe v-loading="loading" empty-text="暂无课件资源">
        <el-table-column prop="name" label="课件名称">
          <template #default="scope">
            <div class="file-item">
              <!-- 根据文件类型显示不同图标 -->
              <i :class="getFileIcon(scope.row.name)" style="margin-right: 10px;"></i>
              <!-- 文件名称 -->
              <span>{{ scope.row.name }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="size" label="文件大小" width="150" align="center">
          <template #default="scope">
                        <span class="file-size">
                            <!-- 重量图标 -->
                            <i class="fas fa-weight-hanging" style="margin-right: 5px;"></i>
                          <!-- 格式化后的文件大小 -->
                            {{ formatFileSize(scope.row.size) }}
                        </span>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="上传时间" width="220" align="center">
          <template #default="scope">
                        <span class="upload-time">
                            <!-- 时钟图标 -->
                            <i class="far fa-clock" style="margin-right: 5px;"></i>
                          <!-- 格式化后的上传时间 -->
                            {{ formatDate(scope.row.created_at) }}
                        </span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" align="center" fixed="right">
          <template #default="scope">
            <!-- 下载按钮 -->
            <el-button type="primary" size="small" @click.stop="downloadFile(scope.row)" plain>
              <i class="fas fa-download"></i> 下载
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>
  </div>
</template>

<script>
import { ElMessage } from 'element-plus'
import dayjs from 'dayjs'
import request from '@/utils/request'

export default {
  name: 'KnowledgePage',
  data() {
    return {
      blogDialogVisible: false,
      teamBlogs: [
        {
          author: 'Martian148',
          title: '锦城学院ACM学习地图',
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
        }
      ],
      blogColors: ['#FF6B6B', '#4ECDC4', '#45B7D1', '#96CEB4'],
      searchQuery: '',
      loading: false,
      folderList: [],
      fileList: [],
      currentFolder: null,
      folderStack: []
    }
  },
  computed: {
    // 过滤文件夹列表
    filteredFolders() {
      if (!this.searchQuery) return this.folderList
      const query = this.searchQuery.toLowerCase()
      return this.folderList.filter(folder =>
          folder.name.toLowerCase().includes(query)
      )
    },
    // 过滤文件列表
    filteredFiles() {
      if (!this.searchQuery) return this.fileList
      const query = this.searchQuery.toLowerCase()
      return this.fileList.filter(file =>
          file.name.toLowerCase().includes(query)
      )
    },
    // 面包屑路径
    breadcrumbFolders() {
      return this.folderStack
    }
  },
  mounted() {
    // 页面加载时获取文件夹列表
    this.fetchFolders()
  },
  methods: {
    // 打开博客弹窗
    openBlogDialog() {
      this.blogDialogVisible = true
    },
    // 获取文件夹和文件列表
    async fetchFolders(folderId = null) {
      this.loading = true
      try {
        const url = folderId ? `/folder/list?parent_id=${folderId}` : '/folder/list'
        const res = await request.get(url)
        this.folderList = res.data.folders || []

        if (folderId) {
          const filesRes = await request.get(`/folder/files?folder_id=${folderId}`)
          this.fileList = filesRes.data.files || []
        } else {
          this.fileList = []
        }

        if (folderId) {
          this.currentFolder = {
            id: folderId,
            name: res.data.current_folder_name || '未知文件夹'
          }
        } else {
          this.currentFolder = null
        }
      } catch (error) {
        // 错误提示
        ElMessage.error('获取资源列表失败')
        console.error(error)
      } finally {
        this.loading = false
      }
    },
    // 打开文件夹
    openFolder(row) {
      this.folderStack = []
      this.folderStack.push({
        id: row.id,
        name: row.name
      })
      this.fetchFolders(row.id)
    },
    // 返回根目录
    goToRoot() {
      this.folderStack = []
      this.fetchFolders()
    },
    // 根据面包屑路径跳转
    goToBreadcrumb(index) {
      this.folderStack = this.folderStack.slice(0, index + 1)
      const folder = this.folderStack[index]
      this.fetchFolders(folder.id)
    },
    // 根据文件扩展名获取文件图标
    getFileIcon(filename) {
      const extension = filename.split('.').pop().toLowerCase()
      switch (extension) {
        case 'pdf': return 'fas fa-file-pdf text-red-500'
        case 'ppt': case 'pptx': return 'fas fa-file-powerpoint text-orange-500'
        case 'doc': case 'docx': return 'fas fa-file-word text-blue-500'
        case 'xls': case 'xlsx': return 'fas fa-file-excel text-green-500'
        case 'zip': case 'rar': return 'fas fa-file-archive text-purple-500'
        default: return 'fas fa-file text-gray-500'
      }
    },
    // 格式化文件大小
    formatFileSize(bytes) {
      if (bytes === 0) return '0 Bytes'
      const k = 1024
      const sizes = ['Bytes', 'KB', 'MB', 'GB']
      const i = Math.floor(Math.log(bytes) / Math.log(k))
      const size = (bytes / Math.pow(k, i)).toFixed(2)
      return `${size} ${sizes[i]}`
    },
    // 格式化日期
    formatDate(date) {
      return dayjs(date).format('YYYY-MM-DD HH:mm')
    },
    // 下载文件
    async downloadFile(file) {
      try {
        // 提示正在下载
        ElMessage.info(`正在下载: ${file.name}`)
        const res = await request.get(`/folder/download/${file.id}`)
        const { name, type, data } = res.data
        const byteCharacters = atob(data)
        const byteNumbers = new Array(byteCharacters.length).fill(0).map((_, i) => byteCharacters.charCodeAt(i))
        const byteArray = new Uint8Array(byteNumbers)
        const blob = new Blob([byteArray], { type })

        const link = document.createElement('a')
        link.href = URL.createObjectURL(blob)
        link.download = name
        document.body.appendChild(link)
        link.click()
        document.body.removeChild(link)

        // 提示下载成功
        ElMessage.success('文件下载成功')
      } catch (error) {
        // 提示下载失败
        ElMessage.error('文件下载失败')
        console.error(error)
      }
    },
    // 重置搜索框
    resetSearch() {
      this.searchQuery = ''
    }
  }
}
</script>
<style scoped>
.blog-avatar {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  margin-right: 1.5rem;
  object-fit: cover;
}
/* 基础样式 */
:root {
  --primary-color: #4361ee;
  --secondary-color: #3f37c9;
  --accent-color: #4895ef;
  --light-color: #f8f9fa;
  --dark-color: #212529;
  --success-color: #4cc9f0;
  --warning-color: #f72585;
  --info-color: #7209b7;
  --hover-bg-color: #e9f2ff;
}

body {
  margin: 0;
  padding: 0;
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
}

.knowledge-container {
  padding: 2rem;
  margin: 0 auto;
  min-height: 100vh;
  width: 90%;
  max-width: 1200px;
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(10px);
  box-shadow: 0 0 30px rgba(0, 0, 0, 0.1);
  border-radius: 20px;
  transition: all 0.3s ease;
}

.knowledge-container:hover {
  box-shadow: 0 0 40px rgba(0, 0, 0, 0.2);
}

.knowledge-title {
  font-size: 3rem;
  color: var(--dark-color);
  margin-bottom: 2rem;
  text-align: center;
  letter-spacing: 2px;
  font-weight: 700;
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
}

.knowledge-title:hover {
  transform: scale(1.05);
}

/* 链接卡片样式 */
.vertical-links {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 2rem;
  margin-bottom: 3rem;
}

.link-card {
  background: white;
  border-radius: 15px;
  padding: 2rem;
  text-align: center;
  transition: all 0.3s ease;
  cursor: pointer;
  text-decoration: none;
  color: inherit;
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.3);
  backdrop-filter: blur(5px);
  position: relative;
  overflow: hidden;
}

.link-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(255, 255, 255, 0.1);
  transform: scaleX(0);
  transform-origin: left;
  transition: transform 0.3s ease;
}

.link-card:hover::before {
  transform: scaleX(1);
}

.link-card:hover {
  transform: translateY(-10px);
  box-shadow: 0 15px 30px rgba(0, 0, 0, 0.2);
}

.link-icon {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  margin: 0 auto 1.5rem;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 2.5rem;
  color: white;
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.2);
  transition: all 0.3s ease;
}

.link-card:hover .link-icon {
  transform: scale(1.1);
}

.link-icon.blog {
  background: linear-gradient(135deg, #4ECDC4 0%, #45B7D1 100%);
}

.link-icon.bilibili {
  background: linear-gradient(135deg, #00A1D6 0%, #0066CC 100%);
}

.link-icon.oj {
  background: linear-gradient(135deg, #FF6B6B 0%, #FF8E8E 100%);
}

.link-card h3 {
  margin: 1rem 0;
  font-size: 1.8rem;
  color: var(--dark-color);
  font-weight: 600;
}

.link-card p {
  color: #6c757d;
  font-size: 1.1rem;
  margin: 0;
  line-height: 1.6;
}

/* 课件部分样式 */
.courseware-section {
  background: white;
  border-radius: 15px;
  padding: 2.5rem;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
  margin-top: 2rem;
  border: 1px solid rgba(255, 255, 255, 0.3);
  backdrop-filter: blur(5px);
  transition: all 0.3s ease;
}

.courseware-section:hover {
  box-shadow: 0 15px 40px rgba(0, 0, 0, 0.2);
}

.courseware-section h2 {
  font-size: 2.2rem;
  color: var(--dark-color);
  margin-bottom: 1.5rem;
  display: flex;
  align-items: center;
  gap: 10px;
  text-shadow: 1px 1px 2px rgba(0, 0, 0, 0.1);
}

.courseware-section h2 i {
  color: var(--accent-color);
}

.breadcrumb {
  margin-bottom: 25px;
  padding: 15px;
  background-color: #f8fafc;
  border-radius: 10px;
  font-size: 1.1rem;
  transition: all 0.3s ease;
}

.breadcrumb:hover {
  background-color: var(--hover-bg-color);
}

.search-bar {
  margin-bottom: 25px;
}

.search-bar .el-input__inner {
  border-radius: 10px;
  transition: all 0.3s ease;
}

.search-bar .el-input__inner:hover {
  border-color: var(--accent-color);
}

/* 博客弹窗样式 */
.blog-dialog {
  border-radius: 15px;
  overflow: hidden;
}

.blog-dialog .el-dialog__header {
  background: linear-gradient(135deg, #4361ee 0%, #3f37c9 100%);
  margin: 0;
  padding: 1.5rem;
}

.blog-dialog .el-dialog__title {
  color: white;
  font-size: 1.5rem;
  font-weight: 600;
}

.blog-item {
  display: flex;
  align-items: center;
  padding: 1.5rem;
  transition: all 0.3s ease;
  border-radius: 10px;
  margin: 5px 0;
}

.blog-item:hover {
  background: var(--hover-bg-color);
  transform: translateX(5px);
}

.blog-item i {
  font-size: 2rem;
  margin-right: 1.5rem;
  flex-shrink: 0;
}

.blog-info {
  display: flex;
  flex-direction: column;
}

.author {
  font-weight: 600;
  color: var(--dark-color);
  margin-bottom: 0.5rem;
  font-size: 1.2rem;
}

.blog-link {
  color: #6c757d;
  text-decoration: none;
  transition: color 0.3s ease;
  font-size: 1.1rem;
}

.blog-link:hover {
  color: var(--accent-color);
  text-decoration: underline;
}

/* 表格样式优化 */
.el-table {
  border-radius: 10px;
  overflow: hidden;
  font-size: 1.1rem;
  transition: all 0.3s ease;
}

.el-table:hover {
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.1);
}

.el-table th {
  font-size: 1.2rem;
  font-weight: 600;
  background-color: #f8fafc !important;
  transition: all 0.3s ease;
}

.el-table th:hover {
  background-color: var(--hover-bg-color) !important;
}

::v-deep(.clickable-table .el-table__row) {
  cursor: pointer;
  transition: all 0.3s ease;
}

::v-deep(.clickable-table .el-table__row:hover) {
  background-color: var(--hover-bg-color) !important;
  transform: scale(1.01);
}

.folder-item,
.file-item {
  display: flex;
  align-items: center;
  font-size: 1.1rem;
}

.create-time,
.upload-time,
.file-size {
  font-size: 1rem;
  color: #6c757d;
}

/* 下载按钮样式 */
.el-button {
  transition: all 0.3s ease;
}

.el-button:hover {
  transform: scale(1.05);
  background-color: var(--accent-color);
  color: white;
  border-color: var(--accent-color);
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .knowledge-container {
    padding: 1.5rem;
  }
}

@media (max-width: 992px) {
  .knowledge-title {
    font-size: 2.5rem;
  }

  .link-card h3 {
    font-size: 1.6rem;
  }

  .courseware-section h2 {
    font-size: 2rem;
  }
}

@media (max-width: 768px) {
  .vertical-links {
    grid-template-columns: 1fr;
  }

  .link-card {
    padding: 1.5rem;
  }

  .knowledge-title {
    font-size: 2rem;
  }

  .courseware-section {
    padding: 1.5rem;
  }

  .el-table-column--fixed-right {
    position: static !important;
  }

  .el-table th,
  .el-table td {
    padding: 12px 5px;
  }
}

@media (max-width: 576px) {
  .knowledge-container {
    padding: 1rem;
  }

  .knowledge-title {
    font-size: 1.8rem;
  }

  .link-icon {
    width: 60px;
    height: 60px;
    font-size: 2rem;
  }

  .link-card h3 {
    font-size: 1.4rem;
  }

  .courseware-section h2 {
    font-size: 1.6rem;
  }

  .folder-item,
  .file-item {
    font-size: 1rem;
  }

  .el-table th {
    font-size: 1rem;
  }

  .el-button {
    padding: 8px 12px;
    font-size: 0.9rem;
  }
}
</style>