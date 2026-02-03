<template>
  <el-dialog
    v-model="visible"
    title="选择图片"
    width="80%"
    :before-close="handleClose"
  >
    <div class="image-picker">
      <!-- 搜索和过滤 -->
      <div class="picker-header">
        <el-radio-group v-model="filterType" size="small" @change="fetchImages">
          <el-radio-button label="all">全部</el-radio-button>
          <el-radio-button label="slider">轮播图</el-radio-button>
          <el-radio-button label="normal">普通图片</el-radio-button>
        </el-radio-group>
        <div class="selected-info" v-if="multiSelect">
          已选择 {{ selectedImages.length }} 张图片
        </div>
      </div>

      <!-- 图片网格 -->
      <div v-loading="loading" class="picker-grid">
        <div
          v-for="img in images"
          :key="img.id"
          class="picker-item"
          :class="{ selected: isSelected(img.id) }"
          @click="toggleSelect(img)"
        >
          <el-image
            :src="getImageUrl(img.url)"
            fit="cover"
            class="picker-image"
          >
            <template #error>
              <div class="image-error">
                <el-icon><picture-filled /></el-icon>
              </div>
            </template>
          </el-image>

          <div v-if="isSelected(img.id)" class="selected-mark">
            <el-icon class="check-icon"><check /></el-icon>
          </div>

          <div class="picker-info">
            <div class="image-name" :title="img.file_name">
              {{ truncateText(img.file_name, 15) }}
            </div>
            <div class="image-size">{{ formatSize(img.size) }}</div>
          </div>
        </div>

        <!-- 空状态 -->
        <div v-if="!loading && images.length === 0" class="empty-state">
          <el-icon class="empty-icon"><picture /></el-icon>
          <p>暂无图片</p>
        </div>
      </div>

      <!-- 分页 -->
      <div v-if="total > pageSize" class="picker-pagination">
        <el-pagination
          v-model:current-page="currentPage"
          :page-size="pageSize"
          :total="total"
          layout="prev, pager, next"
          small
          @current-change="fetchImages"
        />
      </div>
    </div>

    <template #footer>
      <span class="dialog-footer">
        <el-button @click="handleClose">取消</el-button>
        <el-button type="primary" @click="confirmSelection" :disabled="selectedImages.length === 0">
          确认选择 {{ selectedImages.length > 0 ? `(${selectedImages.length})` : '' }}
        </el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script>
import request from '@/utils/request'
import { ElMessage } from 'element-plus'
import { PictureFilled, Check, Picture } from '@element-plus/icons-vue'

export default {
  name: 'ImagePicker',
  components: {
    PictureFilled,
    Check,
    Picture
  },
  props: {
    modelValue: {
      type: Boolean,
      default: false
    },
    multiSelect: {
      type: Boolean,
      default: false
    },
    maxSelect: {
      type: Number,
      default: 0 // 0表示不限制
    }
  },
  emits: ['update:modelValue', 'confirm'],
  data() {
    return {
      images: [],
      loading: false,
      filterType: 'all',
      currentPage: 1,
      pageSize: 12,
      total: 0,
      selectedImages: []
    }
  },
  computed: {
    visible: {
      get() {
        return this.modelValue
      },
      set(value) {
        this.$emit('update:modelValue', value)
      }
    }
  },
  watch: {
    visible(newVal) {
      if (newVal) {
        this.selectedImages = []
        this.currentPage = 1
        this.fetchImages()
      }
    }
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

    isSelected(imageId) {
      return this.selectedImages.some(img => img.id === imageId)
    },

    toggleSelect(image) {
      const index = this.selectedImages.findIndex(img => img.id === image.id)

      if (index > -1) {
        // 取消选择
        this.selectedImages.splice(index, 1)
      } else {
        // 选择图片
        if (this.multiSelect) {
          // 多选模式
          if (this.maxSelect > 0 && this.selectedImages.length >= this.maxSelect) {
            ElMessage.warning(`最多只能选择 ${this.maxSelect} 张图片`)
            return
          }
          this.selectedImages.push(image)
        } else {
          // 单选模式
          this.selectedImages = [image]
        }
      }
    },

    confirmSelection() {
      if (this.selectedImages.length === 0) {
        ElMessage.warning('请选择至少一张图片')
        return
      }

      this.$emit('confirm', this.multiSelect ? this.selectedImages : this.selectedImages[0])
      this.visible = false
    },

    handleClose() {
      this.visible = false
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
.image-picker {
  min-height: 400px;
}

/* 头部 */
.picker-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 15px;
  border-bottom: 1px solid #ebeef5;
}

.selected-info {
  color: #409eff;
  font-size: 14px;
  font-weight: 500;
}

/* 图片网格 */
.picker-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
  gap: 15px;
  min-height: 300px;
  margin-bottom: 20px;
}

.picker-item {
  position: relative;
  border-radius: 8px;
  overflow: hidden;
  cursor: pointer;
  border: 2px solid transparent;
  transition: all 0.3s;
}

.picker-item:hover {
  border-color: #409eff;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.3);
}

.picker-item.selected {
  border-color: #409eff;
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.3);
}

.picker-image {
  width: 100%;
  height: 120px;
  display: block;
}

.picker-image :deep(.el-image__inner) {
  object-fit: cover;
}

.image-error {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  background: #f5f7fa;
  color: #c0c4cc;
  font-size: 36px;
}

.selected-mark {
  position: absolute;
  top: 5px;
  right: 5px;
  width: 24px;
  height: 24px;
  background: #409eff;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 14px;
}

.check-icon {
  font-weight: bold;
}

.picker-info {
  padding: 8px;
  background: #fff;
}

.image-name {
  font-size: 12px;
  color: #303133;
  margin-bottom: 4px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.image-size {
  font-size: 11px;
  color: #909399;
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
  font-size: 60px;
  margin-bottom: 15px;
}

.empty-state p {
  margin: 0;
  font-size: 14px;
}

/* 分页 */
.picker-pagination {
  display: flex;
  justify-content: center;
}

/* 对话框底部 */
.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style>
