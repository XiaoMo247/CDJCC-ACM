<template>
  <div class="admin-card">
    <h1><i class="fas fa-users-cog"></i> 批量注册用户</h1>

    <div class="card-content">
      <!-- 批量注册表单 -->
      <div class="form-section">
        <el-form :model="form" :rules="rules" ref="registerForm" label-position="top">
          <h2>批量注册表单</h2>
          <div class="form-grid">
            <el-form-item label="起始学号" prop="start">
              <el-input v-model="form.start" placeholder="请输入起始学号" size="large" clearable @keyup.enter="handleBatchRegister" />
            </el-form-item>

            <el-form-item label="终止学号" prop="end">
              <el-input v-model="form.end" placeholder="请输入终止学号" size="large" clearable @keyup.enter="handleBatchRegister" />
            </el-form-item>
          </div>

          <el-button type="primary" size="large" @click="handleBatchRegister" :loading="registerLoading" class="action-btn">
            <i class="fas fa-user-plus"></i> 批量注册
          </el-button>

          <!-- 进度条和消息提示 -->
          <div class="progress-section" v-if="registerLoading || message">
            <el-progress :percentage="progress" :status="progressStatus" :stroke-width="12" class="progress-bar" v-if="registerLoading" />
            <div class="message-box" :class="messageType" v-if="message">
              <i :class="messageIcon"></i>
              <span>{{ message }}</span>
            </div>
          </div>
        </el-form>
      </div>

      <!-- 用户列表 -->
      <div class="user-list-section">
        <div class="section-header">
          <h2>用户列表</h2>
          <div class="search-box">
            <el-input v-model="searchKeyword" placeholder="搜索学号或姓名" clearable size="large" @clear="handleSearchClear" @keyup.enter="fetchUsers">
              <template #append>
                <el-button @click="fetchUsers" :icon="SearchIcon" />
              </template>
            </el-input>
          </div>
        </div>

        <el-table :data="userList" v-loading="tableLoading" stripe border :default-sort="{ prop: 'created_at', order: 'descending' }" class="user-table">
          <el-table-column prop="id" label="ID" width="80" align="center" />
          <el-table-column prop="student_number" label="学号" width="150" align="center" sortable />
          <el-table-column prop="username" label="姓名" align="center" sortable/>
          <el-table-column prop="created_at" label="注册时间" width="180" align="center" sortable>
            <template #default="{ row }">
              <span class="time-text">{{ formatDateTime(row.created_at) }}</span>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="200" align="center" fixed="right">
            <template #default="{ row }">
              <el-button size="default" type="warning" @click="showResetDialog(row)" :loading="row.resetting">
                <i class="fas fa-key"></i> 重置
              </el-button>
              <el-button size="default" type="danger" @click="showDeleteDialog(row)" :loading="row.deleting">
                <i class="fas fa-trash-alt"></i> 删除
              </el-button>
            </template>
          </el-table-column>
        </el-table>

        <el-pagination v-model:current-page="pagination.currentPage" v-model:page-size="pagination.pageSize" :total="pagination.total" :page-sizes="[10, 20, 50, 100]" layout="total, sizes, prev, pager, next, jumper" background @size-change="fetchUsers" @current-change="fetchUsers" class="pagination" />
      </div>
    </div>

    <!-- 重置密码确认对话框 -->
    <el-dialog v-model="resetDialogVisible" title="重置密码确认" width="400px" :close-on-click-modal="false" @closed="handleResetDialogClose">
      <div class="dialog-content">
        <i class="fas fa-exclamation-circle warning-icon"></i>
        <span>确定要将用户 <strong>{{ currentUser.username }}</strong> (学号: {{ currentUser.student_number }}) 的密码重置为默认密码吗？</span>
      </div>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="resetDialogVisible = false" size="large">取消</el-button>
          <el-button type="warning" @click="handleResetPassword" :loading="resetLoading" size="large">
            确认重置
          </el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 删除用户确认对话框 -->
    <el-dialog v-model="deleteDialogVisible" title="删除用户确认" width="400px" :close-on-click-modal="false" @closed="handleDeleteDialogClose">
      <div class="dialog-content">
        <i class="fas fa-exclamation-triangle danger-icon"></i>
        <span>确定要删除用户 <strong>{{ currentUser.username }}</strong> (学号: {{ currentUser.student_number }}) 吗？此操作不可恢复！</span>
      </div>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="deleteDialogVisible = false" size="large">取消</el-button>
          <el-button type="danger" @click="handleDeleteUser" :loading="deleteLoading" size="large">
            确认删除
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { ref, reactive, computed } from 'vue'
import { ElMessage } from 'element-plus'
import request from '@/utils/request'
import dayjs from 'dayjs'
import { Search as SearchIcon } from '@element-plus/icons-vue'

export default {
  name: 'BatchRegister',
  setup() {
    // 批量注册相关
    const registerForm = ref(null)
    const form = reactive({
      start: '',
      end: ''
    })

    const rules = reactive({
      start: [
        { required: true, message: '请输入起始学号', trigger: 'blur' },
        { pattern: /^\d+$/, message: '学号必须为数字', trigger: 'blur' }
      ],
      end: [
        { required: true, message: '请输入终止学号', trigger: 'blur' },
        { pattern: /^\d+$/, message: '学号必须为数字', trigger: 'blur' }
      ]
    })

    const registerLoading = ref(false)
    const progress = ref(0)
    const progressStatus = ref('')
    const message = ref('')
    const messageType = computed(() => {
      return progressStatus.value === 'success' ? 'success' :
        progressStatus.value === 'exception' ? 'error' : 'info'
    })
    const messageIcon = computed(() => {
      return messageType.value === 'success' ? 'fas fa-check-circle' :
        messageType.value === 'error' ? 'fas fa-times-circle' : 'fas fa-info-circle'
    })

    // 用户列表相关
    const userList = ref([])
    const tableLoading = ref(false)
    const searchKeyword = ref('')
    const pagination = reactive({
      currentPage: 1,
      pageSize: 10,
      total: 0
    })

    // 对话框相关
    const resetDialogVisible = ref(false)
    const deleteDialogVisible = ref(false)
    const resetLoading = ref(false)
    const deleteLoading = ref(false)
    const currentUser = reactive({
      id: null,
      username: '',
      student_number: ''
    })

    // 方法
    const formatDateTime = (date) => {
      return dayjs(date).format('YYYY-MM-DD HH:mm:ss')
    }

    const handleBatchRegister = async () => {
      try {
        await registerForm.value.validate()
        registerLoading.value = true
        progress.value = 0
        progressStatus.value = ''
        message.value = ''

        // 模拟进度更新
        const timer = setInterval(() => {
          if (progress.value < 90) {
            progress.value += Math.floor(Math.random() * 5) + 1
          }
        }, 300)

        // 实际API调用
        const response = await request.post('/admin/register', {
          start: form.start,
          end: form.end
        })

        clearInterval(timer)
        progress.value = 100
        progressStatus.value = 'success'
        message.value = `批量注册成功，共注册 ${response.data.count} 个用户`

        // 刷新用户列表
        fetchUsers()
      } catch (error) {
        progressStatus.value = 'exception'
        message.value = error.response?.data?.message || '批量注册失败'
        console.error('批量注册失败:', error)
      } finally {
        registerLoading.value = false
      }
    }

    const fetchUsers = async () => {
      try {
        tableLoading.value = true
        const params = {
          page: pagination.currentPage,
          limit: pagination.pageSize,
          keyword: searchKeyword.value
        }
        const response = await request.get('/admin/get/user', { params })
        userList.value = response.data.data.map(user => ({
          ...user,
          resetting: false,
          deleting: false
        }))
        pagination.total = response.data.total
      } catch (error) {
        ElMessage.error('获取用户列表失败')
        console.error('获取用户列表失败:', error)
      } finally {
        tableLoading.value = false
      }
    }

    const handleSearchClear = () => {
      searchKeyword.value = ''
      pagination.currentPage = 1
      fetchUsers()
    }

    const showResetDialog = (user) => {
      currentUser.id = user.id
      currentUser.username = user.username
      currentUser.student_number = user.student_number
      resetDialogVisible.value = true
    }

    const handleResetPassword = async () => {
      try {
        resetLoading.value = true
        await request.put(`/admin/user/reset/${currentUser.id}`)

        ElMessage.success('密码重置成功')
        resetDialogVisible.value = false
        fetchUsers()
      } catch (error) {
        ElMessage.error(error.response?.data?.message || '密码重置失败')
      } finally {
        resetLoading.value = false
      }
    }

    const handleResetDialogClose = () => {
      currentUser.id = null
      currentUser.username = ''
      currentUser.student_number = ''
    }

    const showDeleteDialog = (user) => {
      currentUser.id = user.id
      currentUser.username = user.username
      currentUser.student_number = user.student_number
      deleteDialogVisible.value = true
    }

    const handleDeleteUser = async () => {
      try {
        deleteLoading.value = true
        await request.delete(`/admin/user/${currentUser.id}`)

        ElMessage.success('用户删除成功')
        deleteDialogVisible.value = false

        // 如果当前页只剩一条数据且不是第一页，则返回上一页
        if (userList.value.length === 1 && pagination.currentPage > 1) {
          pagination.currentPage -= 1
        }

        fetchUsers()
      } catch (error) {
        ElMessage.error(error.response?.data?.message || '用户删除失败')
      } finally {
        deleteLoading.value = false
      }
    }

    const handleDeleteDialogClose = () => {
      currentUser.id = null
      currentUser.username = ''
      currentUser.student_number = ''
    }

    // 初始化加载用户列表
    fetchUsers()

    return {
      // 表单相关
      registerForm,
      form,
      rules,
      registerLoading,
      progress,
      progressStatus,
      message,
      messageType,
      messageIcon,

      // 用户列表相关
      userList,
      tableLoading,
      searchKeyword,
      pagination,
      SearchIcon,

      // 对话框相关
      resetDialogVisible,
      deleteDialogVisible,
      resetLoading,
      deleteLoading,
      currentUser,

      // 方法
      formatDateTime,
      handleBatchRegister,
      fetchUsers,
      handleSearchClear,
      showResetDialog,
      handleResetPassword,
      handleResetDialogClose,
      showDeleteDialog,
      handleDeleteUser,
      handleDeleteDialogClose
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

.card-content {
  display: flex;
  flex-direction: column;
  gap: 30px;
}

.form-section {
  background: rgba(255, 255, 255, 0.7);
  padding: 20px;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.form-section h2 {
  font-size: 18px;
  color: #2d3748;
  margin-bottom: 20px;
  padding-bottom: 8px;
  border-bottom: 1px solid #e2e8f0;
}

.form-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 20px;
  margin-bottom: 20px;
}

.action-btn {
  width: 100%;
  margin-top: 10px;
}

.progress-section {
  margin-top: 20px;
}

.message-box {
  padding: 12px 16px;
  border-radius: 6px;
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 15px;
  font-size: 14px;
}

.message-box.success {
  background-color: #f0fdf4;
  color: #166534;
  border: 1px solid #bbf7d0;
}

.message-box.error {
  background-color: #fef2f2;
  color: #991b1b;
  border: 1px solid #fecaca;
}

.message-box.info {
  background-color: #eff6ff;
  color: #1e40af;
  border: 1px solid #bfdbfe;
}

.user-list-section {
  background: rgba(255, 255, 255, 0.7);
  padding: 20px;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  flex-wrap: wrap;
  gap: 15px;
}

.section-header h2 {
  font-size: 18px;
  color: #2d3748;
  margin: 0;
  padding-bottom: 8px;
  border-bottom: 1px solid #e2e8f0;
}

.search-box {
  flex: 1;
  min-width: 250px;
  max-width: 400px;
}

.user-table {
  width: 100%;
  margin: 15px 0;
}

.time-text {
  font-family: monospace;
  font-size: 16px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}

.dialog-content {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  line-height: 1.6;
}

.dialog-content strong {
  color: #1a365d;
}

.warning-icon {
  color: #d97706;
  font-size: 20px;
  margin-top: 2px;
}

.danger-icon {
  color: #dc2626;
  font-size: 20px;
  margin-top: 2px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

@media (max-width: 768px) {
  .form-grid {
    grid-template-columns: 1fr;
  }
  
  .section-header {
    flex-direction: column;
    align-items: stretch;
  }
  
  .search-box {
    width: 100%;
    max-width: none;
  }
}
</style>
