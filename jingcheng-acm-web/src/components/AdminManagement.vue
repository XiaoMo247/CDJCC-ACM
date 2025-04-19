<template>
  <div class="admin-management-container">
    <!-- 顶部标题区 -->
    <div class="header-section">
      <h1 class="page-title">
        <el-icon><UserFilled /></el-icon>
        <span>管理员管理系统</span>
      </h1>
      <div class="page-subtitle">管理后台管理员账号及权限</div>
    </div>

    <!-- 添加管理员卡片 -->
    <el-card class="add-admin-card" shadow="hover">
  <template #header>
    <div class="card-header">
      <el-icon><CirclePlus /></el-icon>
      <span>添加新管理员</span>
    </div>
  </template>

  <el-form 
    :model="form" 
    label-position="left" 
    class="admin-form"
    @submit.prevent="addAdmin"
  >
    <el-form-item label="用户名" required>
      <el-input
        v-model="form.username"
        placeholder="输入新管理员用户名"
        size="large"
        clearable
      >
        <template #prefix>
          <el-icon><User /></el-icon>
        </template>
      </el-input>
    </el-form-item>

    <el-form-item label="密码" required>
      <el-input
        v-model="form.password"
        type="password"
        placeholder="输入6位以上密码"
        size="large"
        show-password
      >
        <template #prefix>
          <el-icon><Lock /></el-icon>
        </template>
      </el-input>
    </el-form-item>

    <el-button
      type="primary"
      class="submit-btn"
      @click="addAdmin"
      :loading="isAdding"
    >
      <el-icon><Plus /></el-icon>
      <span>添加管理员</span>
    </el-button>
  </el-form>
</el-card>

    <!-- 管理员列表卡片 -->
    <el-card class="admin-list-card" shadow="hover">
      <template #header>
        <div class="card-header">
          <el-icon><List /></el-icon>
          <span>管理员列表</span>
          <el-tag type="info" class="count-tag">
            共 {{ admins.length }} 位管理员
          </el-tag>
        </div>
      </template>

      <el-table
        :data="admins"
        style="width: 100%"
        v-loading="isLoading"
        empty-text="暂无管理员数据"
        class="admin-table"
      >
        <el-table-column prop="username" label="用户名" width="220">
          <template #default="{ row }">
            <div class="user-cell">
              <el-avatar :size="36" :src="getAvatar(row.username)" />
              <span class="username">{{ row.username }}</span>
            </div>
          </template>
        </el-table-column>

        <el-table-column prop="created_at" label="创建时间" width="220">
          <template #default="{ row }">
            <div class="time-cell">
              <el-icon><Clock /></el-icon>
              <span>{{ formatDate(row.created_at) }}</span>
            </div>
          </template>
        </el-table-column>

        <el-table-column label="操作" width="180" align="right">
          <template #default="{ row }">
            <el-popconfirm
              title="确定要删除此管理员吗？"
              confirm-button-text="确认"
              cancel-button-text="取消"
              @confirm="deleteAdmin(row.id)"
            >
              <template #reference>
                <el-button
                  type="danger"
                  size="small"
                  :icon="Delete"
                  circle
                  v-if="row.username !== currentAdmin"
                />
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import axios from 'axios'
import { 
  UserFilled, CirclePlus, User, Lock, Plus, 
  List, Clock, Delete 
} from '@element-plus/icons-vue'

const form = ref({
  username: '',
  password: ''
})
const admins = ref([])
const isAdding = ref(false)
const isLoading = ref(false)
const currentAdmin = ref('admin') // 假设当前登录管理员

// 获取随机头像
const getAvatar = (username) => {
  const colors = ['#3b82f6', '#6366f1', '#8b5cf6', '#ec4899', '#ef4444']
  const color = colors[username.length % colors.length]
  return `https://ui-avatars.com/api/?name=${username}&background=${color.slice(1)}&color=fff&size=128`
}

// 格式化时间
const formatDate = (timestamp) => {
  const date = new Date(timestamp)
  return date.toLocaleString('zh-CN', { 
    year: 'numeric', 
    month: '2-digit', 
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 获取管理员列表
const getAdmins = async () => {
  try {
    isLoading.value = true
    const response = await axios.get('http://localhost:8081/api/admin/list', {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('adminToken')}`
      }
    })
    admins.value = response.data.admins || []
  } catch (error) {
    ElMessage.error('获取管理员列表失败: ' + (error.response?.data?.msg || error.message))
  } finally {
    isLoading.value = false
  }
}

// 添加管理员
const addAdmin = async () => {
  if (!form.value.username || !form.value.password) {
    return ElMessage.warning('请填写用户名和密码')
  }
  if (form.value.password.length < 6) {
    return ElMessage.warning('密码长度不能少于6位')
  }

  try {
    isAdding.value = true
    const response = await axios.post('http://localhost:8081/api/admin/add', {
      username: form.value.username,
      password: form.value.password
    }, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('adminToken')}`
      }
    })

    ElMessage.success(response.data.msg || '管理员添加成功')
    form.value.username = ''
    form.value.password = ''
    await getAdmins()
  } catch (error) {
    ElMessage.error('添加失败: ' + (error.response?.data?.msg || error.message))
  } finally {
    isAdding.value = false
  }
}

// 删除管理员
const deleteAdmin = async (adminId) => {
  try {
    const response = await axios.delete(`http://localhost:8081/api/admin/delete/${adminId}`, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('adminToken')}`
      }
    })
    ElMessage.success(response.data.msg || '管理员删除成功')
    await getAdmins()
  } catch (error) {
    ElMessage.error('删除失败: ' + (error.response?.data?.msg || error.message))
  }
}

onMounted(() => {
  getAdmins()
})
</script>

<style scoped>
.admin-management-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 30px;
}

.header-section {
  margin-bottom: 30px;
  text-align: center;
}

.page-title {
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 28px;
  font-weight: 600;
  color: #1e293b;
  margin-bottom: 8px;
}

.page-title .el-icon {
  margin-right: 10px;
  font-size: 32px;
  color: #3b82f6;
}

.page-subtitle {
  font-size: 14px;
  color: #64748b;
}

.add-admin-card {
  margin-bottom: 30px;
  border-radius: 12px;
  border: 1px solid #e2e8f0;
}

.admin-list-card {
  border-radius: 12px;
  border: 1px solid #e2e8f0;
}

.card-header {
  display: flex;
  align-items: center;
  font-size: 18px;
  font-weight: 500;
}

.card-header .el-icon {
  margin-right: 8px;
  color: #3b82f6;
}

.count-tag {
  margin-left: 12px;
}

.admin-form {
  padding: 0 20px;
}

.submit-btn {
  width: 100%;
  height: 48px;
  margin-top: 10px;
  font-size: 16px;
  letter-spacing: 1px;
}

.submit-btn .el-icon {
  margin-right: 8px;
}

.admin-table {
  margin-top: -10px;
}

.user-cell {
  display: flex;
  align-items: center;
}

.user-cell .username {
  margin-left: 12px;
  font-weight: 500;
}

.time-cell {
  display: flex;
  align-items: center;
  color: #64748b;
}

.time-cell .el-icon {
  margin-right: 8px;
  font-size: 16px;
}

/* 动画效果 */
.el-table__row {
  transition: all 0.3s ease;
}

.el-table__row:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .admin-management-container {
    padding: 20px 15px;
  }
  
  .page-title {
    font-size: 22px;
  }
  
  .card-header {
    font-size: 16px;
  }
}
</style>