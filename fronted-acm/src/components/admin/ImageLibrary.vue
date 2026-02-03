<template>
  <div class="admin-card">
    <h1><i class="fas fa-images"></i> 图片管理</h1>

    <!-- 上传区域 -->
    <div class="upload-section">
      <el-upload
        class="upload-container"
        drag
        :action="uploadUrl"
        :headers="uploadHeaders"
        :on-success="handleUploadSuccess"
        :on-error="handleUploadError"
        :before-upload="beforeUpload"
        :show-file-list="false"
        accept="image/*"
        name="image"
        multiple
      >
        <el-icon class="el-icon--upload"><upload-filled /></el-icon>
        <div class="el-upload__text">
          拖拽图片到此处或<em>点击上传</em>
        </div>
        <template #tip>
          <div class="el-upload__tip">
            支持 jpg/png/gif/webp 格式，单个文件不超过 5MB
          </div>
        </template>
      </el-upload>
    </div>

    <!-- 过滤器 -->
    <div class="filter-section">
      <el-radio-group v-model="filterType" @change="fetchImages">
        <el-radio-button label="all">全部图片</el-radio-button>
        <el-radio-button label="slider">轮播图片</el-radio-button>
        <el-radio-button label="normal">普通图片</el-radio-button>
      </el-radio-group>
      <div class="stats">
        共 {{ total }} 张图片
      </div>
    </div>

    <!-- 图片网格 -->
    <div v-loading="loading" class="image-grid">
      <div v-for="img in images" :key="img.id" class="image-card">
        <el-image
          :src="getImageUrl(img.url)"
          fit="cover"
          :preview-src-list="[getImageUrl(img.url)]"
          class="image-preview"
        >
          <template #error>
            <div class="image-error">
              <el-icon><picture-filled /></el-icon>
            </div>
          </template>
        </el-image>

        <div class="image-overlay">
          <div class="image-actions">
            <el-button
              type="primary"
              size="small"
              circle
              @click="togglePreview(img)"
            >
              <el-icon><zoom-in /></el-icon>
            </el-button>
            <el-button
              type="danger"
              size="small"
              circle
              @click="deleteImage(img.id)"
            >
              <el-icon><delete /></el-icon>
            </el-button>
          </div>
        </div>

        <div class="image-info">
          <div class="image-name" :title="img.file_name">
            {{ truncateText(img.file_name, 20) }}
          </div>
          <div class="image-meta">
            <span class="image-size">{{ formatSize(img.size) }}</span>
            <span class="image-dimensions">{{ img.width }} × {{ img.height }}</span>
          </div>
          <div class="image-controls">
            <el-checkbox
              v-model="img.is_slider"
              @change="updateSliderFlag(img)"
            >
              轮播图
            </el-checkbox>
          </div>
        </div>
      </div>

      <!-- 空状态 -->
      <div v-if="!loading && images.length === 0" class="empty-state">
        <el-icon class="empty-icon"><picture /></el-icon>
        <p>暂无图片</p>
        <p class="empty-tip">点击上方上传区域添加图片</p>
      </div>
    </div>

    <!-- 分页 -->
    <div v-if="total > pageSize" class="pagination">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[12, 24, 48, 96]"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="fetchImages"
        @current-change="fetchImages"
      />
    </div>
  </div>
</template>

<script>
import request from '@/utils/request'
import { ElMessage, ElMessageBox } from 'element-plus'
import { UploadFilled, ZoomIn, Delete, PictureFilled, Picture } from '@element-plus/icons-vue'
import { getToken } from '@/utils/tokenManager'

export default {
  name: 'ImageLibrary',
  components: {
    UploadFilled,
    ZoomIn,
    Delete,
    PictureFilled,
    Picture
  },
  data() {
    return {
      images: [],
      loading: false,
      filterType: 'all',
      currentPage: 1,
      pageSize: 24,
      total: 0,
      uploadUrl: `${import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api'}/admin/image/upload`,
      uploadHeaders: {
        Authorization: `Bearer ${getToken()}`
      }
    }
  },
  mounted() {
    this.fetchImages()
  },
  methods: {
    async fetchImages() {
      this.loading = true
      try {
        const params = {
          page: this.currentPage,
          limit: this.pageSize
        }

        if (this.filterType === 'slider') {
          params.is_slider = 'true'
        } else if (this.filterType === 'normal') {
          params.is_slider = 'false'
        }

        const res = await request.get('/admin/image/list', { params })
        this.images = res.data.data.images || []
        this.total = res.data.data.total || 0
      } catch (error) {
        ElMessage.error('获取图片列表失败')
        console.error(error)
      } finally {
        this.loading = false
      }
    },

    beforeUpload(file) {
      const isImage = file.type.startsWith('image/')
      const isLt5M = file.size / 1024 / 1024 < 5

      if (!isImage) {
        ElMessage.error('只能上传图片文件')
        return false
      }
      if (!isLt5M) {
        ElMessage.error('图片大小不能超过 5MB')
        return false
      }
      return true
    },

    handleUploadSuccess(response) {
      ElMessage.success('上传成功')
      this.fetchImages()
    },

    handleUploadError(error) {
      console.error('Upload error:', error)
      ElMessage.error('上传失败')
    },

    async updateSliderFlag(img) {
      try {
        await request.put(`/admin/image/update/${img.id}`, {
          is_slider: img.is_slider
        })
        ElMessage.success(img.is_slider ? '已标记为轮播图' : '已取消轮播图标记')
      } catch (error) {
        ElMessage.error('更新失败')
        img.is_slider = !img.is_slider // 回滚
      }
    },

    async deleteImage(id) {
      try {
        await ElMessageBox.confirm(
          '确定要删除这张图片吗？此操作不可恢复！',
          '警告',
          {
            confirmButtonText: '确定删除',
            cancelButtonText: '取消',
            type: 'warning'
          }
        )

        await request.delete(`/admin/image/delete/${id}`)
        ElMessage.success('删除成功')
        this.fetchImages()
      } catch (error) {
        if (error !== 'cancel') {
          ElMessage.error('删除失败')
        }
      }
    },

    togglePreview(img) {
      // Element Plus的el-image组件已经支持预览，点击图片即可
    },

    getImageUrl(url) {
      if (url.startsWith('http')) {
        return url
      }
      return `${import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'}${url}`.replace('/api', '')
    },

    formatSize(bytes) {
      if (bytes === 0) return '0 B'
      const k = 1024
      const sizes = ['B', 'KB', 'MB', 'GB']
      const i = Math.floor(Math.log(bytes) / Math.log(k))
      return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
    },

    truncateText(text, maxLength) {
      if (text.length <= maxLength) return text
      return text.substring(0, maxLength) + '...'
    }
  }
}
</script>

<style scoped>
.admin-card {
  padding: 20px;
}

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

/* 上传区域 */
.upload-section {
  margin-bottom: 30px;
}

.upload-container {
  width: 100%;
}

.upload-container :deep(.el-upload-dragger) {
  width: 100%;
  height: 180px;
}

.el-icon--upload {
  font-size: 67px;
  color: #c0c4cc;
  margin-bottom: 16px;
}

/* 过滤器 */
.filter-section {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding: 15px;
  background: #f5f7fa;
  border-radius: 8px;
}

.stats {
  color: #606266;
  font-size: 14px;
}

/* 图片网格 */
.image-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 20px;
  min-height: 300px;
}

.image-card {
  position: relative;
  border-radius: 8px;
  overflow: hidden;
  background: #fff;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  transition: transform 0.3s, box-shadow 0.3s;
}

.image-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
}

.image-card:hover .image-overlay {
  opacity: 1;
}

.image-preview {
  width: 100%;
  height: 200px;
  display: block;
}

.image-preview :deep(.el-image__inner) {
  object-fit: cover;
}

.image-error {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  background: #f5f7fa;
  color: #c0c4cc;
  font-size: 48px;
}

.image-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.3s;
  pointer-events: none;
}

.image-actions {
  display: flex;
  gap: 10px;
  pointer-events: all;
}

.image-info {
  padding: 12px;
}

.image-name {
  font-size: 14px;
  font-weight: 500;
  color: #303133;
  margin-bottom: 8px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.image-meta {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: #909399;
  margin-bottom: 10px;
}

.image-controls {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

/* 空状态 */
.empty-state {
  grid-column: 1 / -1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  color: #909399;
}

.empty-icon {
  font-size: 80px;
  margin-bottom: 20px;
}

.empty-state p {
  margin: 5px 0;
  font-size: 16px;
}

.empty-tip {
  font-size: 14px;
  color: #c0c4cc;
}

/* 分页 */
.pagination {
  margin-top: 30px;
  display: flex;
  justify-content: center;
}

/* 响应式 */
@media (max-width: 768px) {
  .image-grid {
    grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
    gap: 15px;
  }

  .image-preview {
    height: 150px;
  }

  .filter-section {
    flex-direction: column;
    gap: 15px;
  }
}
</style>
