<template>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.5.0/css/all.min.css">
    <div class="knowledge-container">
        <!-- 博客弹窗 -->
        <el-dialog v-model="blogDialogVisible" title="团队成员博客" width="90%" class="blog-dialog" :fullscreen="isMobile">
            <div class="dialog-content">
                <ul class="blog-list">
                    <li class="blog-item" v-for="(blog, index) in teamBlogs" :key="index" @click="openBlog(blog.url)">
                        <div class="blog-icon-container">
                            <i class="fas fa-blog" :style="{ color: blogColors[index % 4] }"></i>
                        </div>
                        <div class="blog-info">
                            <span class="author">{{ blog.author }}</span>
                            <a :href="blog.url" target="_blank" class="blog-link">{{ blog.title }}</a>
                            <div class="blog-url">{{ formatUrl(blog.url) }}</div>
                        </div>
                        <i class="fas fa-external-link-alt link-arrow"></i>
                    </li>
                </ul>
            </div>
        </el-dialog>

        <!-- 页面标题区域 -->
        <div class="header-section">
            <h1 class="knowledge-title">
                <span class="title-icon">📚</span>
                <span class="title-text">知识库</span>
                <span class="title-decoration"></span>
            </h1>
            <p class="knowledge-subtitle">汇聚优质资源，助力学习成长</p>
        </div>

        <!-- 快速链接卡片 -->
        <div class="quick-links-section">
            <div class="link-card blog-card" @click="openBlogDialog">
                <div class="link-icon-container">
                    <i class="fas fa-users"></i>
                </div>
                <div class="link-content">
                    <h3>团队博客</h3>
                    <p>探索成员技术分享</p>
                    <div class="link-hover-effect">
                        <span>立即查看</span>
                        <i class="fas fa-arrow-right"></i>
                    </div>
                </div>
            </div>

            <a href="https://space.bilibili.com/3546651937475184?spm_id_from=333.337.search-card.all.click"
                target="_blank" class="link-card bilibili-card">
                <div class="link-icon-container">
                    <i class="fab fa-bilibili"></i>
                </div>
                <div class="link-content">
                    <h3>B站频道</h3>
                    <p>观看视频教程</p>
                    <div class="link-hover-effect">
                        <span>前往观看</span>
                        <i class="fas fa-arrow-right"></i>
                    </div>
                </div>
            </a>

            <a href="https://hydro.ac/d/cdjcc_acm_2333/" target="_blank" class="link-card oj-card">
                <div class="link-icon-container">
                    <i class="fas fa-code"></i>
                </div>
                <div class="link-content">
                    <h3>OJ平台</h3>
                    <p>在线编程训练</p>
                    <div class="link-hover-effect">
                        <span>开始练习</span>
                        <i class="fas fa-arrow-right"></i>
                    </div>
                </div>
            </a>
        </div>

        <!-- 课件资源部分 -->
        <div class="courseware-section">
            <div class="section-header">
                <h2>
                    <i class="fas fa-book-open"></i>
                    <span>课件资源</span>
                </h2>
                <div class="section-actions">
                    <el-tooltip content="刷新数据" placement="top">
                        <el-button circle @click="refreshData">
                            <i class="fas fa-sync-alt"></i>
                        </el-button>
                    </el-tooltip>
                </div>
            </div>

            <div class="navigation-controls">
                <div class="breadcrumb-container">
                    <el-breadcrumb separator="/">
                        <el-breadcrumb-item @click="goToRoot" class="breadcrumb-item">
                            <i class="fas fa-home"></i>
                            <span>根目录</span>
                        </el-breadcrumb-item>
                        <el-breadcrumb-item v-for="(folder, index) in breadcrumbFolders" :key="index"
                            @click="goToBreadcrumb(index)" class="breadcrumb-item">
                            <i class="fas fa-folder"></i>
                            <span>{{ folder.name }}</span>
                        </el-breadcrumb-item>
                    </el-breadcrumb>
                </div>

                <div class="search-container">
                    <el-input v-model="searchQuery" placeholder="搜索课件..." clearable size="large" class="search-input"
                        @clear="resetSearch">
                        <template #prefix>
                            <i class="el-icon-search"></i>
                        </template>
                        <template #append>
                            <el-button @click="searchFiles">
                                <i class="fas fa-search"></i>
                            </el-button>
                        </template>
                    </el-input>
                </div>
            </div>

            <!-- 文件夹列表 -->
            <el-table :data="filteredFolders" style="width: 100%; margin-bottom: 20px;" stripe v-loading="loading"
                empty-text="暂无文件夹" @row-click="openFolder" class="clickable-table folder-table">
                <el-table-column prop="name" label="文件夹名称" width="400">
                    <template #default="scope">
                        <div class="folder-item">
                            <div class="folder-icon">
                                <i class="fas fa-folder"></i>
                            </div>
                            <div class="folder-info">
                                <span class="folder-name">{{ scope.row.name }}</span>
                                <span class="folder-items" v-if="scope.row.item_count">
                                    {{ scope.row.item_count }}个项目
                                </span>
                            </div>
                        </div>
                    </template>
                </el-table-column>
                <el-table-column prop="created_at" label="创建时间" width="220" align="center">
                    <template #default="scope">
                        <div class="create-time">
                            <i class="far fa-calendar-alt"></i>
                            <span>{{ formatDate(scope.row.created_at) }}</span>
                        </div>
                    </template>
                </el-table-column>
                <el-table-column label="操作" width="120" align="center">
                    <template #default="scope">
                        <el-button type="text" @click.stop="previewFolder(scope.row)">
                            <i class="fas fa-eye"></i>
                        </el-button>
                    </template>
                </el-table-column>
            </el-table>

            <!-- 文件列表 -->
            <el-table :data="filteredFiles" style="width: 100%" stripe v-loading="loading" empty-text="暂无课件资源"
                class="file-table">
                <el-table-column prop="name" label="课件名称" width="400">
                    <template #default="scope">
                        <div class="file-item">
                            <div class="file-icon" :class="getFileType(scope.row.name)">
                                <i :class="getFileIcon(scope.row.name)"></i>
                            </div>
                            <div class="file-info">
                                <span class="file-name">{{ scope.row.name }}</span>
                                <span class="file-type">{{ getFileExtension(scope.row.name) }}</span>
                            </div>
                        </div>
                    </template>
                </el-table-column>
                <el-table-column prop="size" label="文件大小" width="150" align="center">
                    <template #default="scope">
                        <div class="file-size">
                            <i class="fas fa-weight-hanging"></i>
                            <span>{{ formatFileSize(scope.row.size) }}</span>
                        </div>
                    </template>
                </el-table-column>
                <el-table-column prop="created_at" label="上传时间" width="220" align="center">
                    <template #default="scope">
                        <div class="upload-time">
                            <i class="far fa-clock"></i>
                            <span>{{ formatDate(scope.row.created_at) }}</span>
                        </div>
                    </template>
                </el-table-column>
                <el-table-column label="操作" width="180" align="center" fixed="right">
                    <template #default="scope">
                        <el-button-group>
                            <el-tooltip content="预览" placement="top">
                                <el-button type="primary" size="small" @click.stop="previewFile(scope.row)" circle>
                                    <i class="fas fa-eye"></i>
                                </el-button>
                            </el-tooltip>
                            <el-tooltip content="下载" placement="top">
                                <el-button type="success" size="small" @click.stop="downloadFile(scope.row)" circle>
                                    <i class="fas fa-download"></i>
                                </el-button>
                            </el-tooltip>
                            <el-tooltip content="分享" placement="top">
                                <el-button type="info" size="small" @click.stop="shareFile(scope.row)" circle>
                                    <i class="fas fa-share-alt"></i>
                                </el-button>
                            </el-tooltip>
                        </el-button-group>
                    </template>
                </el-table-column>
            </el-table>

            <!-- 分页控件 -->
            <div class="pagination-container" v-if="showPagination">
                <el-pagination background layout="prev, pager, next" :total="totalItems" :page-size="pageSize"
                    :current-page="currentPage" @current-change="handlePageChange" />
            </div>
        </div>
    </div>
</template>

<script>
import { ElMessage, ElMessageBox } from 'element-plus'
import dayjs from 'dayjs'
import request from '@/utils/request'

export default {
    name: 'KnowledgePage',
    data() {
        return {
            isMobile: window.innerWidth < 768,
            blogDialogVisible: false,
            teamBlogs: [
                { author: 'Martian148', title: '锦城学院ACM学习地图', url: 'https://www.cnblogs.com/martian148/p/18221024' },
                { author: 'Martian148', title: '个人博客', url: 'https://martian148.xyz/' },
                { author: 'XiaoMo247', title: '个人博客', url: 'https://www.cnblogs.com/XiaoMo247' }
            ],
            blogColors: ['#FF6B6B', '#4ECDC4', '#45B7D1', '#96CEB4'],
            searchQuery: '',
            loading: false,
            folderList: [],
            fileList: [],
            currentFolder: null,
            folderStack: [],
            currentPage: 1,
            pageSize: 10,
            totalItems: 0
        }
    },
    computed: {
        filteredFolders() {
            if (!this.searchQuery) return this.folderList
            const query = this.searchQuery.toLowerCase()
            return this.folderList.filter(folder =>
                folder.name.toLowerCase().includes(query)
            )
        },
        filteredFiles() {
            if (!this.searchQuery) return this.fileList
            const query = this.searchQuery.toLowerCase()
            return this.fileList.filter(file =>
                file.name.toLowerCase().includes(query)
            )
        },
        breadcrumbFolders() {
            return this.folderStack
        },
        showPagination() {
            return this.totalItems > this.pageSize
        }
    },
    mounted() {
        this.fetchFolders()
        window.addEventListener('resize', this.handleResize)
    },
    beforeUnmount() {
        window.removeEventListener('resize', this.handleResize)
    },
    methods: {
        openBlogDialog() {
            this.blogDialogVisible = true;
        },
        handleResize() {
            this.isMobile = window.innerWidth < 768
        },
        openBlog(url) {
            window.open(url, '_blank')
        },
        formatUrl(url) {
            try {
                const urlObj = new URL(url)
                return urlObj.hostname.replace('www.', '')
            } catch {
                return url.length > 30 ? url.substring(0, 30) + '...' : url
            }
        },
        async refreshData() {
            if (this.currentFolder) {
                await this.fetchFolders(this.currentFolder.id)
            } else {
                await this.fetchFolders()
            }
            ElMessage.success('数据已刷新')
        },
        async fetchFolders(folderId = null) {
            this.loading = true
            try {
                const url = folderId ? `/folder/list?parent_id=${folderId}` : '/folder/list'
                const res = await request.get(url)
                this.folderList = res.data.folders || []

                if (folderId) {
                    const filesRes = await request.get(`/folder/files?folder_id=${folderId}`)
                    this.fileList = filesRes.data.files || []
                    this.totalItems = filesRes.data.total || 0
                } else {
                    this.fileList = []
                    this.totalItems = 0
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
                ElMessage.error('获取资源列表失败')
                console.error(error)
            } finally {
                this.loading = false
            }
        },
        openFolder(row) {
            this.currentPage = 1
            this.folderStack.pop()
            this.folderStack.push({
                id: row.id,
                name: row.name
            })
            this.fetchFolders(row.id)
        },
        goToRoot() {
            this.currentPage = 1
            this.folderStack = []
            this.fetchFolders()
        },
        goToBreadcrumb(index) {
            this.currentPage = 1
            this.folderStack = this.folderStack.slice(0, index + 1)
            const folder = this.folderStack[index]
            this.fetchFolders(folder.id)
        },
        getFileIcon(filename) {
            const extension = this.getFileExtension(filename).toLowerCase()
            switch (extension) {
                case 'pdf': return 'fas fa-file-pdf'
                case 'ppt': case 'pptx': return 'fas fa-file-powerpoint'
                case 'doc': case 'docx': return 'fas fa-file-word'
                case 'xls': case 'xlsx': return 'fas fa-file-excel'
                case 'zip': case 'rar': return 'fas fa-file-archive'
                case 'jpg': case 'jpeg': case 'png': case 'gif': return 'fas fa-file-image'
                case 'mp4': case 'avi': case 'mov': return 'fas fa-file-video'
                case 'mp3': case 'wav': return 'fas fa-file-audio'
                case 'txt': return 'fas fa-file-alt'
                case 'html': case 'htm': return 'fas fa-file-code'
                default: return 'fas fa-file'
            }
        },
        getFileType(filename) {
            const extension = this.getFileExtension(filename).toLowerCase()
            return `file-type-${extension}`
        },
        getFileExtension(filename) {
            return filename.split('.').pop()
        },
        formatFileSize(bytes) {
            if (bytes === 0) return '0 Bytes'
            const k = 1024
            const sizes = ['Bytes', 'KB', 'MB', 'GB']
            const i = Math.floor(Math.log(bytes) / Math.log(k))
            const size = (bytes / Math.pow(k, i)).toFixed(2)
            return `${size} ${sizes[i]}`
        },
        formatDate(date) {
            return dayjs(date).format('YYYY-MM-DD HH:mm')
        },
        async downloadFile(file) {
            try {
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

                ElMessage.success('文件下载成功')
            } catch (error) {
                ElMessage.error('文件下载失败')
                console.error(error)
            }
        },
        previewFile(file) {
            ElMessageBox.alert('文件预览功能正在开发中，敬请期待！', '预览提示', {
                confirmButtonText: '确定'
            })
        },
        previewFolder(folder) {
            ElMessageBox.alert(`即将查看文件夹: ${folder.name}`, '文件夹预览', {
                confirmButtonText: '确定'
            })
        },
        shareFile(file) {
            const shareUrl = `${window.location.origin}/share/file/${file.id}`
            navigator.clipboard.writeText(shareUrl).then(() => {
                ElMessage.success('分享链接已复制到剪贴板')
            }).catch(() => {
                ElMessage.warning('无法复制到剪贴板，请手动复制链接')
            })
        },
        searchFiles() {
            if (this.searchQuery.trim()) {
                ElMessage.info(`正在搜索: ${this.searchQuery}`)
            }
        },
        resetSearch() {
            this.searchQuery = ''
        },
        handlePageChange(page) {
            this.currentPage = page
            // 这里可以添加分页加载数据的逻辑
        }
    }
}
</script>

<style scoped>
/* 基础变量 */
:root {
    --primary-color: #4361ee;
    --primary-light: #eef2ff;
    --secondary-color: #3f37c9;
    --accent-color: #4895ef;
    --success-color: #4cc9f0;
    --warning-color: #f72585;
    --info-color: #7209b7;
    --light-color: #f8f9fa;
    --dark-color: #212529;
    --gray-color: #6c757d;
    --light-gray: #e9ecef;
    --border-radius: 12px;
    --box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
    --transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1);
}

/* 基础样式 */
.knowledge-container {
    max-width: 1400px;
    margin: 0 auto;
    padding: 2rem 1.5rem;
    min-height: 100vh;
}

/* 头部区域 */
.header-section {
    text-align: center;
    margin-bottom: 3rem;
    position: relative;
}

.knowledge-title {
    font-size: 3rem;
    font-weight: 800;
    color: var(--dark-color);
    margin-bottom: 1rem;
    position: relative;
    display: inline-block;
}

.knowledge-title .title-icon {
    display: inline-block;
    margin-right: 15px;
    transform: rotate(-10deg);
}

.knowledge-title .title-text {
    position: relative;
    z-index: 2;
}

.knowledge-title .title-decoration {
    position: absolute;
    bottom: 5px;
    left: 0;
    width: 100%;
    height: 15px;
    background: linear-gradient(90deg, rgba(67, 97, 238, 0.2), transparent);
    z-index: 1;
    border-radius: 10px;
}

.knowledge-subtitle {
    font-size: 1.2rem;
    color: var(--gray-color);
    max-width: 600px;
    margin: 0 auto;
    line-height: 1.6;
}

/* 快速链接区域 */
.quick-links-section {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
    gap: 2rem;
    margin-bottom: 4rem;
}

.link-card {
    background: white;
    border-radius: var(--border-radius);
    overflow: hidden;
    transition: var(--transition);
    box-shadow: var(--box-shadow);
    position: relative;
    z-index: 1;
    border: none;
    text-decoration: none;
    color: inherit;
    cursor: pointer;
}

.link-card::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: linear-gradient(135deg, rgba(255, 255, 255, 0.9), rgba(255, 255, 255, 0.7));
    backdrop-filter: blur(5px);
    z-index: -1;
}

.link-card:hover {
    transform: translateY(-10px);
    box-shadow: 0 15px 40px rgba(0, 0, 0, 0.15);
}

.link-icon-container {
    height: 120px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 3.5rem;
    color: white;
    position: relative;
    overflow: hidden;
}

.link-icon-container::after {
    content: '';
    position: absolute;
    top: -50%;
    left: -50%;
    width: 200%;
    height: 200%;
    background: linear-gradient(135deg, rgba(255, 255, 255, 0.2), transparent);
    transform: rotate(30deg);
    transition: var(--transition);
}

.link-card:hover .link-icon-container::after {
    transform: rotate(30deg) translateX(20%);
}

.link-content {
    padding: 1.5rem;
    position: relative;
}

.link-content h3 {
    font-size: 1.5rem;
    font-weight: 700;
    margin-bottom: 0.5rem;
    color: var(--dark-color);
}

.link-content p {
    color: var(--gray-color);
    margin-bottom: 1.5rem;
    font-size: 1rem;
}

.link-hover-effect {
    display: flex;
    align-items: center;
    color: var(--primary-color);
    font-weight: 600;
    opacity: 0;
    transform: translateX(-10px);
    transition: var(--transition);
}

.link-card:hover .link-hover-effect {
    opacity: 1;
    transform: translateX(0);
}

.link-hover-effect i {
    margin-left: 8px;
    transition: var(--transition);
}

.link-card:hover .link-hover-effect i {
    transform: translateX(5px);
}

/* 卡片特定样式 */
.blog-card .link-icon-container {
    background: linear-gradient(135deg, #4ECDC4, #45B7D1);
}

.bilibili-card .link-icon-container {
    background: linear-gradient(135deg, #00A1D6, #0066CC);
}

.oj-card .link-icon-container {
    background: linear-gradient(135deg, #FF6B6B, #FF8E8E);
}

/* 课件区域 */
.courseware-section {
    background: white;
    border-radius: var(--border-radius);
    padding: 2.5rem;
    box-shadow: var(--box-shadow);
    margin-top: 2rem;
    position: relative;
    overflow: hidden;
}

.courseware-section::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    width: 5px;
    height: 100%;
    background: linear-gradient(to bottom, var(--primary-color), var(--accent-color));
}

.section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 2rem;
}

.section-header h2 {
    font-size: 1.8rem;
    font-weight: 700;
    color: var(--dark-color);
    display: flex;
    align-items: center;
    margin: 0;
}

.section-header h2 i {
    color: var(--primary-color);
    margin-right: 12px;
    font-size: 1.5rem;
}

.section-actions .el-button {
    background: var(--primary-light);
    color: var(--primary-color);
    border: none;
    font-size: 1rem;
    width: 40px;
    height: 40px;
}

.section-actions .el-button:hover {
    background: var(--primary-color);
    color: white;
}

/* 导航控制 */
.navigation-controls {
    display: flex;
    flex-wrap: wrap;
    gap: 1.5rem;
    margin-bottom: 2rem;
}

.breadcrumb-container {
    flex: 1;
    min-width: 300px;
}

.breadcrumb-item {
    display: inline-flex;
    align-items: center;
    cursor: pointer;
    transition: var(--transition);
}

.breadcrumb-item:hover {
    color: var(--primary-color);
}

.breadcrumb-item i {
    margin-right: 8px;
    font-size: 0.9rem;
}

.search-container {
    min-width: 300px;
    flex: 1;
}

.search-input {
    border-radius: 30px;
    overflow: hidden;
}

.search-input :deep(.el-input-group__append) {
    background: var(--primary-color);
    border: none;
    color: white;
}

.search-input :deep(.el-input-group__append:hover) {
    background: var(--secondary-color);
}

/* 表格样式 */
.clickable-table {
    --el-table-border-color: var(--light-gray);
    --el-table-header-bg-color: var(--light-color);
}

.clickable-table :deep(.el-table__row) {
    cursor: pointer;
    transition: var(--transition);
}

.clickable-table :deep(.el-table__row:hover) {
    background-color: var(--primary-light) !important;
    transform: scale(1.005);
}

.folder-table :deep(.el-table__cell) {
    padding: 15px 0;
}

.file-table :deep(.el-table__cell) {
    padding: 12px 0;
}

.folder-item,
.file-item {
    display: flex;
    align-items: center;
    padding: 0 15px;
}

.folder-icon {
    width: 40px;
    height: 40px;
    background: rgba(255, 215, 0, 0.1);
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    margin-right: 15px;
    color: #FFD700;
    font-size: 1.2rem;
}

.file-icon {
    width: 40px;
    height: 40px;
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    margin-right: 15px;
    color: white;
    font-size: 1.2rem;
}

.file-type-pdf {
    background: #FF6B6B;
}

.file-type-ppt,
.file-type-pptx {
    background: #FF9E3F;
}

.file-type-doc,
.file-type-docx {
    background: #4A89DC;
}

.file-type-xls,
.file-type-xlsx {
    background: #48CFAD;
}

.file-type-zip,
.file-type-rar {
    background: #AC92EC;
}

.file-type-jpg,
.file-type-jpeg,
.file-type-png,
.file-type-gif {
    background: #ED5565;
}

.file-type-mp4,
.file-type-avi,
.file-type-mov {
    background: #5D9CEC;
}

.file-type-mp3,
.file-type-wav {
    background: #A0D468;
}

.file-type-txt {
    background: #656D78;
}

.file-type-html,
.file-type-htm {
    background: #E9573F;
}

.file-type-default {
    background: #AAB2BD;
}

.folder-info,
.file-info {
    display: flex;
    flex-direction: column;
}

.folder-name,
.file-name {
    font-weight: 500;
    margin-bottom: 3px;
}

.folder-items,
.file-type {
    font-size: 0.8rem;
    color: var(--gray-color);
}

.create-time,
.upload-time,
.file-size {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
    font-size: 0.9rem;
}

.create-time i,
.upload-time i,
.file-size i {
    color: var(--gray-color);
}

.el-button-group {
    display: flex;
    gap: 5px;
}

.el-button-group .el-button {
    width: 32px;
    height: 32px;
    padding: 0;
    display: flex;
    align-items: center;
    justify-content: center;
}

/* 分页样式 */
.pagination-container {
    margin-top: 2rem;
    display: flex;
    justify-content: center;
}

.pagination-container :deep(.el-pagination.is-background .el-pager li:not(.is-disabled).is-active) {
    background: var(--primary-color);
}

.pagination-container :deep(.el-pagination.is-background .el-pager li:not(.is-disabled):hover) {
    color: var(--primary-color);
}

/* 博客弹窗样式 */
.blog-dialog :deep(.el-dialog) {
    border-radius: var(--border-radius);
    overflow: hidden;
    max-width: 800px;
}

.blog-dialog :deep(.el-dialog__header) {
    background: linear-gradient(135deg, var(--primary-color), var(--secondary-color));
    margin: 0;
    padding: 1.5rem;
}

.blog-dialog :deep(.el-dialog__title) {
    color: white;
    font-size: 1.5rem;
    font-weight: 600;
}

.blog-dialog :deep(.el-dialog__headerbtn .el-dialog__close) {
    color: white;
    font-size: 1.5rem;
}

.dialog-content {
    max-height: 60vh;
    overflow-y: auto;
    padding: 0 20px;
}

.blog-list {
    list-style: none;
    padding: 0;
    margin: 0;
}

.blog-item {
    display: flex;
    align-items: center;
    padding: 1.5rem;
    transition: var(--transition);
    border-radius: 8px;
    margin: 5px 0;
    border: 1px solid var(--light-gray);
    cursor: pointer;
}

.blog-item:hover {
    background: var(--light-color);
    border-color: var(--primary-light);
    transform: translateX(5px);
}

.blog-icon-container {
    width: 50px;
    height: 50px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 1.5rem;
    margin-right: 1.5rem;
    flex-shrink: 0;
    background: rgba(0, 0, 0, 0.05);
}

.blog-info {
    flex: 1;
    min-width: 0;
}

.author {
    font-weight: 600;
    color: var(--dark-color);
    margin-bottom: 0.5rem;
    font-size: 1.1rem;
}

.blog-link {
    color: var(--dark-color);
    text-decoration: none;
    transition: var(--transition);
    font-size: 1rem;
    display: block;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}

.blog-item:hover .blog-link {
    color: var(--primary-color);
}

.blog-url {
    font-size: 0.8rem;
    color: var(--gray-color);
    margin-top: 5px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}

.link-arrow {
    color: var(--gray-color);
    margin-left: 1rem;
    opacity: 0;
    transition: var(--transition);
}

.blog-item:hover .link-arrow {
    opacity: 1;
    color: var(--primary-color);
}

/* 响应式设计 */
@media (max-width: 1200px) {
    .knowledge-container {
        padding: 1.5rem;
    }
}

@media (max-width: 992px) {
    .knowledge-title {
        font-size: 2.8rem;
    }

    .section-header h2 {
        font-size: 1.6rem;
    }

    .quick-links-section {
        grid-template-columns: 1fr 1fr;
    }
}

@media (max-width: 768px) {
    .knowledge-title {
        font-size: 2.2rem;
    }

    .knowledge-subtitle {
        font-size: 1rem;
    }

    .quick-links-section {
        grid-template-columns: 1fr;
    }

    .navigation-controls {
        flex-direction: column;
        gap: 1rem;
    }

    .breadcrumb-container,
    .search-container {
        min-width: 100%;
    }

    .courseware-section {
        padding: 1.5rem;
    }

    .folder-item,
    .file-item {
        padding: 0 10px;
    }

    .folder-icon,
    .file-icon {
        margin-right: 10px;
        width: 36px;
        height: 36px;
        font-size: 1rem;
    }
}

@media (max-width: 576px) {
    .knowledge-title {
        font-size: 1.8rem;
    }

    .header-section {
        margin-bottom: 2rem;
    }

    .link-content h3 {
        font-size: 1.3rem;
    }

    .section-header {
        flex-direction: column;
        align-items: flex-start;
        gap: 1rem;
    }

    .section-actions {
        width: 100%;
        display: flex;
        justify-content: flex-end;
    }

    .blog-item {
        padding: 1rem;
    }

    .blog-icon-container {
        width: 40px;
        height: 40px;
        font-size: 1.2rem;
        margin-right: 1rem;
    }
}
</style>