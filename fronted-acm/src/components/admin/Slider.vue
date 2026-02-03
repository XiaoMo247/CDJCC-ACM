<template>
  <div class="admin-card">
    <h1><i class="fas fa-images"></i> 轮播图管理</h1>

    <div class="form-container">
      <el-button type="primary" @click="openAddSliderDialog" size="large" class="submit-btn">
        <i class="fas fa-plus"></i> 添加轮播图
      </el-button>
    </div>

    <el-table :data="sliderList" style="width: 100%; margin-top: 20px;" class="slider-table" stripe border>
      <el-table-column prop="title" label="标题" width="200" header-align="center" />
      <el-table-column prop="content" label="内容" header-align="center" show-overflow-tooltip />
      <el-table-column label="图片" width="200" header-align="center" align="center">
        <template #default="{ row }">
          <el-image
            :src="getImageUrl(row)"
            style="width: 160px; height: 90px"
            fit="cover"
            :preview-src-list="[getImageUrl(row)]"
          />
        </template>
      </el-table-column>
      <el-table-column prop="order" label="顺序" width="100" header-align="center" align="center" />
      <el-table-column label="状态" width="100" header-align="center" align="center">
        <template #default="{ row }">
          <el-tag :type="row.is_active ? 'success' : 'info'">
            {{ row.is_active ? '启用' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="120" header-align="center" align="center">
        <template #default="{ row }">
          <el-button size="small" type="danger" @click="deleteSlider(row.id)" class="delete-btn">
            <i class="fas fa-trash-alt"></i> 删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 添加轮播图对话框 -->
    <el-dialog v-model="addSliderDialogVisible" title="添加轮播图" width="600px">
      <el-form :model="newSlider" label-width="80px" ref="sliderFormRef">
        <el-form-item label="标题" prop="title" :rules="[{ required: true, message: '标题不能为空', trigger: 'blur' }]">
          <el-input v-model="newSlider.title" placeholder="请输入标题" size="large" clearable />
        </el-form-item>
        <el-form-item label="内容" prop="content">
          <el-input v-model="newSlider.content" type="textarea" :rows="3" placeholder="请输入内容" size="large" />
        </el-form-item>
        <el-form-item label="图片" prop="image_id" :rules="[{ required: true, message: '请选择图片', trigger: 'change' }]">
          <div class="image-select-area">
            <div v-if="selectedImage" class="selected-image-preview">
              <el-image
                :src="selectedImage.url"
                style="width: 200px; height: 120px"
                fit="cover"
              />
              <div class="image-info">
                <p>{{ selectedImage.file_name }}</p>
                <p class="image-size">{{ formatSize(selectedImage.size) }}</p>
              </div>
            </div>
            <el-button type="primary" @click="openImagePicker">
              <i class="fas fa-image"></i> {{ selectedImage ? '重新选择' : '选择图片' }}
            </el-button>
          </div>
        </el-form-item>
        <el-form-item label="顺序" prop="order">
          <el-input-number v-model="newSlider.order" :min="0" :max="999" size="large" />
          <span style="margin-left: 10px; color: #909399; font-size: 13px;">数字越小越靠前</span>
        </el-form-item>
        <el-form-item label="启用">
          <el-switch v-model="newSlider.is_active" />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="addSliderDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="addSlider">确定</el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 图片选择器 -->
    <ImagePicker
      v-model="imagePickerVisible"
      :multi-select="false"
      @confirm="handleImageSelected"
    />
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/request'
import ImagePicker from './ImagePicker.vue'

export default {
  name: 'SliderManagement',
  components: {
    ImagePicker
  },
  setup() {
    const sliderList = ref([])
    const loading = ref(false)
    const addSliderDialogVisible = ref(false)
    const imagePickerVisible = ref(false)
    const selectedImage = ref(null)
    const newSlider = ref({
      title: '',
      content: '',
      image_id: null,
      order: 0,
      is_active: true
    })
    const sliderFormRef = ref(null)

    // 获取轮播图列表
    const fetchSliderList = async () => {
      loading.value = true
      try {
        const response = await request.get('/slider/list')
        sliderList.value = response.data.data || []
      } catch (error) {
        ElMessage.error('获取轮播图列表失败')
      } finally {
        loading.value = false
      }
    }

    // 打开添加对话框
    const openAddSliderDialog = () => {
      newSlider.value = {
        title: '',
        content: '',
        image_id: null,
        order: 0,
        is_active: true
      }
      selectedImage.value = null
      addSliderDialogVisible.value = true
    }

    // 打开图片选择器
    const openImagePicker = () => {
      imagePickerVisible.value = true
    }

    // 处理图片选择
    const handleImageSelected = (image) => {
      selectedImage.value = image
      newSlider.value.image_id = image.id
    }

    // 添加轮播图
    const addSlider = async () => {
      if (!newSlider.value.title) {
        ElMessage.warning('请输入标题')
        return
      }
      if (!newSlider.value.image_id) {
        ElMessage.warning('请选择图片')
        return
      }

      try {
        const formData = new FormData()
        formData.append('title', newSlider.value.title)
        formData.append('content', newSlider.value.content || '')
        formData.append('image_id', newSlider.value.image_id)
        formData.append('order', newSlider.value.order)
        formData.append('is_active', newSlider.value.is_active ? 'true' : 'false')

        await request.post('/admin/slider/add', formData, {
          headers: { 'Content-Type': 'multipart/form-data' }
        })

        ElMessage.success('轮播图添加成功')
        addSliderDialogVisible.value = false
        fetchSliderList()
      } catch (error) {
        ElMessage.error(error.response?.data?.message || '添加失败')
      }
    }

    // 删除轮播图
    const deleteSlider = async (id) => {
      try {
        await ElMessageBox.confirm('确定要删除该轮播图吗？', '警告', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        })

        await request.delete(`/admin/slider/delete/${id}`)
        ElMessage.success('删除成功')
        fetchSliderList()
      } catch (error) {
        if (error !== 'cancel') {
          ElMessage.error('删除失败')
        }
      }
    }

    // 获取图片URL
    const getImageUrl = (row) => {
      // 优先使用ImageObj中的URL
      if (row.image_obj && row.image_obj.url) {
        const url = row.image_obj.url
        if (url.startsWith('http')) {
          return url
        }
        return `${import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'}${url}`.replace('/api', '')
      }
      // 回退到旧的image字段
      if (row.image) {
        if (row.image.startsWith('http')) {
          return row.image
        }
        return `${import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'}${row.image}`.replace('/api', '')
      }
      return ''
    }

    const formatSize = (bytes) => {
      if (bytes === 0) return '0 B'
      const k = 1024
      const sizes = ['B', 'KB', 'MB', 'GB']
      const i = Math.floor(Math.log(bytes) / Math.log(k))
      return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
    }

    onMounted(() => {
      fetchSliderList()
    })

    return {
      sliderList,
      loading,
      addSliderDialogVisible,
      imagePickerVisible,
      selectedImage,
      newSlider,
      sliderFormRef,
      openAddSliderDialog,
      openImagePicker,
      handleImageSelected,
      addSlider,
      deleteSlider,
      getImageUrl,
      formatSize
    }
  }
}
</script>

<style scoped>
.admin-card {
  padding: 24px;
  background: #fff;
  border-radius: 8px;
}

h1 {
  font-size: 24px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 24px;
}

h1 i {
  margin-right: 8px;
  color: #409eff;
}

.form-container {
  margin-bottom: 20px;
}

.slider-table {
  border-radius: 8px;
  overflow: hidden;
}

.slider-table :deep(.el-table__header) {
  background-color: #f5f7fa;
}

.image-select-area {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.selected-image-preview {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 12px;
  background: #f5f7fa;
  border-radius: 8px;
}

.image-info p {
  margin: 4px 0;
  font-size: 14px;
  color: #606266;
}

.image-size {
  font-size: 12px;
  color: #909399;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style>
