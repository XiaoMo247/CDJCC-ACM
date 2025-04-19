<template>
  <div class="announcement-management">
    <!-- 头部标题 -->
    <div class="header">
      <h3 class="title">
        <el-icon :size="24" class="icon"><Message /></el-icon>
        公告管理
      </h3>
      <p class="subtitle">发布和管理系统公告信息</p>
    </div>

    <!-- 发布公告表单 -->
    <el-card class="form-container" shadow="never">
      <template #header>
        <div class="card-header">
          <span class="card-title">发布新公告</span>
          <el-tag type="info" size="small">最多500字</el-tag>
        </div>
      </template>
      
      <el-form 
        :model="form" 
        label-width="120px" 
        label-position="top"
        @submit.prevent="submitAnnouncement"
        class="announcement-form"
      >
        <el-form-item 
          label="公告标题" 
          prop="title"
          :rules="formRules.title"
        >
          <el-input 
            v-model="form.title" 
            placeholder="请输入公告标题"
            clearable
            class="input-with-count"
          />
          <span class="input-count">{{ form.title.length }}/50</span>
        </el-form-item>

        <el-form-item 
          label="公告内容" 
          prop="content"
          :rules="formRules.content"
        >
          <el-input 
            type="textarea" 
            v-model="form.content" 
            placeholder="请输入公告内容"
            :rows="5"
            resize="none"
            show-word-limit
            maxlength="500"
            class="content-textarea"
          />
        </el-form-item>

        <div class="form-actions">
          <el-button 
            type="primary" 
            native-type="submit"
            :loading="submitting"
            class="submit-btn"
          >
            <template #icon>
              <el-icon><Upload /></el-icon>
            </template>
            立即发布
          </el-button>
          <el-button @click="resetForm">重置</el-button>
        </div>
      </el-form>
    </el-card>

    <!-- 当前公告列表 -->
    <el-card class="list-container" shadow="never">
      <template #header>
        <div class="list-header">
          <span class="header-title">
            <el-icon><List /></el-icon>
            当前公告（{{ announcements.length }}条）
          </span>
          <el-button 
            type="text" 
            @click="fetchAnnouncements"
            :loading="loading"
            class="refresh-btn"
          >
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </div>
      </template>

      <el-table 
        :data="announcements" 
        v-loading="loading"
        empty-text="暂无公告"
        stripe
        style="width: 100%"
        class="announcement-table"
      >
      <el-table-column label="公告标题" width="220" class-name="title-column">
    <template #default="{ row }">
      <div v-if="row.isEditing" class="edit-container">
        <el-input 
          v-model="row.title" 
          placeholder="请输入标题"
          class="full-width-input"
        />
      </div>
      <div v-else class="display-content">
        <el-tag effect="plain" type="info">
          {{ row.title || '--' }}
        </el-tag>
      </div>
    </template>
  </el-table-column>

  <el-table-column label="公告内容" class-name="content-column">
    <template #default="{ row }">
      <div v-if="row.isEditing" class="edit-container">
        <el-input 
          v-model="row.content" 
          placeholder="请输入内容"
          type="textarea"
          :rows="4"
          resize="none"
          class="full-width-textarea"
        />
      </div>
      <div v-else class="display-content">
        {{ row.content || '--' }}
      </div>
    </template>
  </el-table-column>

        <el-table-column label="发布时间" width="180" class-name="time-column">
          <template #default="{ row }">
            <div class="time-wrapper">
              <el-icon class="time-icon"><Clock /></el-icon>
              <span class="time">{{ formatDate(row.created_at) }}</span>
            </div>
          </template>
        </el-table-column>

        <el-table-column label="操作" width="200" fixed="right" class-name="action-column">
          <template #default="{ row }">
            <div class="action-buttons">
              <template v-if="row.isEditing">
                <el-button 
                  type="success" 
                  size="small"
                  @click="updateAnnouncement(row)"
                  :loading="row.updating"
                  class="save-btn"
                >
                  <el-icon><Check /></el-icon>
                  保存
                </el-button>
                <el-button 
                  size="small"
                  @click="row.isEditing = false"
                  class="cancel-btn"
                >
                  <el-icon><Close /></el-icon>
                  取消
                </el-button>
              </template>
              <template v-else>
                <el-button 
                  type="primary" 
                  size="small"
                  @click="editAnnouncement(row)"
                  class="edit-btn"
                >
                  <el-icon><Edit /></el-icon>
                  编辑
                </el-button>
                <el-popconfirm 
                  title="确定要删除该公告吗？"
                  confirm-button-text="确定"
                  cancel-button-text="取消"
                  @confirm="deleteAnnouncement(row.id)"
                >
                  <template #reference>
                    <el-button 
                      type="danger" 
                      size="small"
                      class="delete-btn"
                    >
                      <el-icon><Delete /></el-icon>
                      删除
                    </el-button>
                  </template>
                </el-popconfirm>
              </template>
            </div>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { 
  Message, 
  Upload, 
  List, 
  Refresh, 
  Edit, 
  Delete, 
  Check, 
  Close,
  Clock
} from '@element-plus/icons-vue'
import axios from 'axios'

// 表单数据
const form = ref({
  title: '',
  content: ''
})

// 验证规则
const formRules = {
  title: [
    { required: true, message: '请输入公告标题', trigger: 'blur' },
    { max: 50, message: '标题长度不能超过50个字符', trigger: 'blur' }
  ],
  content: [
    { required: true, message: '请输入公告内容', trigger: 'blur' },
    { max: 500, message: '内容长度不能超过500个字符', trigger: 'blur' }
  ]
}

// 状态管理
const announcements = ref([])
const submitting = ref(false)
const loading = ref(false)

// API请求相关
const fetchAnnouncements = async () => {
  try {
    loading.value = true
    const response = await axios.get('http://localhost:8081/api/announcement/list', {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('adminToken')}`
      }
    })
    announcements.value = (response.data.announcements || []).map(item => ({
      ...item,
      isEditing: false,
      updating: false
    }))
  } catch (error) {
    ElMessage.error('获取公告失败')
    console.error('获取公告失败:', error)
  } finally {
    loading.value = false
  }
}

// 提交公告
const submitAnnouncement = async () => {
  if (!form.value.title || !form.value.content) {
    ElMessage.error('请填写公告标题和内容')
    return
  }
  
  submitting.value = true
  try {
    const response = await axios.post('http://localhost:8081/api/announcement/create', form.value, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('adminToken')}`
      }
    })
    
    if (response.data.msg === '创建成功') {
      ElMessage.success('公告发布成功')
      form.value.title = ''
      form.value.content = ''
      fetchAnnouncements()
    } else {
      ElMessage.error(response.data.msg || '发布公告失败')
    }
  } catch (error) {
    ElMessage.error('发布公告失败')
    console.error('发布公告失败:', error)
  } finally {
    submitting.value = false
  }
}

// 编辑公告
const editAnnouncement = (announcement) => {
  // 保存原始值以便取消时恢复
  announcement.originalTitle = announcement.title
  announcement.originalContent = announcement.content
  announcements.value.forEach(item => item.isEditing = false)
  announcement.isEditing = true
}

const updateAnnouncement = async (announcement) => {
  if (!announcement.title || !announcement.content) {
    ElMessage.error('请填写公告标题和内容')
    return
  }
  
  try {
    announcement.updating = true
    const response = await axios.put(
      `http://localhost:8081/api/announcement/update/${announcement.id}`, 
      {
        title: announcement.title,
        content: announcement.content
      }, 
      {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('adminToken')}`
        }
      }
    )

    if (response.data.msg === '修改成功') {
      ElMessage.success('公告更新成功')
      announcement.isEditing = false
      // 清除原始值
      delete announcement.originalTitle
      delete announcement.originalContent
    } else {
      ElMessage.error(response.data.msg || '更新公告失败')
    }
  } catch (error) {
    ElMessage.error('更新公告失败')
    console.error('更新公告失败:', error)
  } finally {
    announcement.updating = false
  }
}

// 取消编辑时恢复原始值
const cancelEdit = (row) => {
  if (row.originalTitle) {
    row.title = row.originalTitle
  }
  if (row.originalContent) {
    row.content = row.originalContent
  }
  row.isEditing = false
}

// 删除公告
const deleteAnnouncement = async (id) => {
  try {
    const response = await axios.delete(
      `http://localhost:8081/api/announcement/delete/${id}`, 
      {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('adminToken')}`
        }
      }
    )
    
    if (response.data.msg === '删除成功') {
      ElMessage.success('公告删除成功')
      fetchAnnouncements()
    } else {
      ElMessage.error(response.data.msg || '删除公告失败')
    }
  } catch (error) {
    ElMessage.error('删除公告失败')
    console.error('删除公告失败:', error)
  }
}

// 重置表单
const resetForm = () => {
  form.value = {
    title: '',
    content: ''
  }
}

// 格式化时间
const formatDate = (timestamp) => {
  const date = new Date(timestamp)
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  
  return `${year}-${month}-${day} ${hours}:${minutes}`
}

// 初始化
onMounted(fetchAnnouncements)
</script>

<style lang="scss" scoped>
.announcement-management {
  padding: 20px;
  background-color: #f8fafc;
  min-height: calc(100vh - 40px);
  
  .header {
    margin-bottom: 24px;
    
    .title {
      display: flex;
      align-items: center;
      font-size: 22px;
      color: #303133;
      margin: 0 0 8px 0;
      font-weight: 600;
      
      .icon {
        margin-right: 10px;
        color: var(--el-color-primary);
      }
    }
    
    .subtitle {
      margin: 0;
      font-size: 14px;
      color: #909399;
    }
  }
  
  .form-container {
    margin-bottom: 24px;
    border-radius: 8px;
    border: 1px solid #ebeef5;
    background-color: #fff;
    
    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      
      .card-title {
        font-size: 16px;
        font-weight: 500;
      }
    }
    
    .announcement-form {
      padding: 0 10px;
      
      :deep(.el-form-item__label) {
        font-weight: 500;
        padding-bottom: 6px;
      }
      
      .input-with-count {
        position: relative;
      }
      
      .input-count {
        position: absolute;
        right: 10px;
        bottom: -20px;
        font-size: 12px;
        color: #909399;
      }
      
      .content-textarea {
        :deep(.el-textarea__inner) {
          line-height: 1.5;
          font-family: inherit;
        }
      }
      
      .form-actions {
        display: flex;
        justify-content: flex-end;
        margin-top: 24px;
        
        .submit-btn {
          padding: 10px 20px;
        }
      }
    }
  }
  
  .list-container {
    border-radius: 8px;
    border: 1px solid #ebeef5;
    background-color: #fff;
    
    .list-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      
      .header-title {
        font-size: 16px;
        font-weight: 500;
        display: flex;
        align-items: center;
        
        .el-icon {
          margin-right: 8px;
          color: var(--el-color-primary);
        }
      }
      
      .refresh-btn {
        padding: 4px 8px;
      }
    }
    
    .announcement-table {
      :deep(.el-table__header) {
        th {
          background-color: #f5f7fa;
          font-weight: 600;
        }
      }
      
      :deep(.el-table__row) {
        td {
          padding: 12px 0;
        }
      }
      
      .title-column {
        .display-content {
          font-weight: 500;
        }
      }
      
      .content-column {
        .display-content {
          color: #606266;
          line-height: 1.5;
          white-space: pre-wrap;
          word-break: break-word;
        }
      }
      
      .time-wrapper {
        display: flex;
        align-items: center;
        
        .time-icon {
          margin-right: 6px;
          color: #909399;
          font-size: 14px;
        }
        
        .time {
          color: #909399;
          font-size: 13px;
        }
      }
      
      .action-buttons {
        display: flex;
        gap: 8px;
        
        .el-button {
          display: flex;
          align-items: center;
          
          .el-icon {
            margin-right: 4px;
          }
        }
      }
      
      .editing-input {
        :deep(.el-input__inner),
        :deep(.el-textarea__inner) {
          padding: 0 8px;
          font-size: inherit;
        }
      }
    }
  }
}

.edit-container {
  padding: 8px 0;
  
  .full-width-input {
    width: 100%;
    
    :deep(.el-input__inner) {
      padding: 8px 12px;
      font-size: inherit;
    }
  }
  
  .full-width-textarea {
    width: 100%;
    
    :deep(.el-textarea__inner) {
      padding: 12px;
      line-height: 1.5;
      font-size: inherit;
      min-height: 100px;
    }
  }
}

@media (max-width: 768px) {
  .announcement-management {
    padding: 12px;
    
    .form-container,
    .list-container {
      :deep(.el-card__body) {
        padding: 16px;
      }
    }
    
    .action-buttons {
      flex-direction: column;
      gap: 4px !important;
      
      .el-button {
        width: 100%;
        margin: 0 !important;
      }
    }
  }
  .action-buttons {
    flex-direction: column;
    gap: 4px;
    
    .el-button {
      width: 100%;
    }
  }
}
</style>