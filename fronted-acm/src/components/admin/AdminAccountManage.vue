<template>
  <div class="admin-card">
    <h1><i class="fas fa-user-cog"></i>管理员账户管理</h1>

    <div class="form-container">
      <el-form :model="form" inline class="form-inline">
        <el-form-item label="用户名" class="form-item">
          <el-input v-model="form.username" placeholder="请输入管理员用户名" size="large" clearable />
        </el-form-item>
        <el-form-item label="密码" class="form-item">
          <el-input v-model="form.password" type="password" placeholder="请输入初始密码" size="large" show-password />
        </el-form-item>
        <el-button type="primary" @click="addAdmin" size="large" class="add-btn">
          <i class="fas fa-plus"></i> 新增管理员
        </el-button>
      </el-form>
    </div>

    <el-table :data="adminList" style="width: 100%; margin-top: 20px;" class="admin-table" stripe border>
      <el-table-column prop="username" label="用户名" header-align="center" align="center" />
      <el-table-column label="创建时间" header-align="center" align="center">
        <template #default="scope">
          <span class="create-time">
            {{ formatDate(scope.row.created_at) }}
          </span>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="150" header-align="center" align="center">
        <template #default="scope">
          <el-button size="large" type="danger" @click="deleteAdmin(scope.row.id)" class="delete-btn">
            <i class="fas fa-trash-alt"></i> 删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script>
import request from '@/utils/request'
import { ElMessage, ElMessageBox } from 'element-plus'
import dayjs from 'dayjs'

export default {
  name: 'AdminAccountManage',
  data() {
    return {
      adminList: [],
      form: {
        username: '',
        password: ''
      }
    }
  },
  mounted() {
    this.fetchAdmins()
  },
  methods: {
    async fetchAdmins() {
      try {
        const res = await request.get('/admin/list')
        this.adminList = res.data.admins || []
      } catch (err) {
        ElMessage.error('获取管理员列表失败')
      }
    },
    async addAdmin() {
      if (!this.form.username || !this.form.password) {
        ElMessage.warning('请填写用户名和密码')
        return
      }
      try {
        await request.post('/admin/add', {
          username: this.form.username,
          password: this.form.password
        })
        ElMessage.success('新增管理员成功')
        this.form.username = ''
        this.form.password = ''
        this.fetchAdmins()
      } catch (err) {
        ElMessage.error(err.response?.data?.message || '新增失败')
      }
    },
    async deleteAdmin(id) {
      try {
        await ElMessageBox.confirm('确定要删除该管理员吗？此操作不可恢复！', '警告', {
          confirmButtonText: '确定删除',
          cancelButtonText: '取消',
          type: 'warning',
          confirmButtonClass: 'confirm-delete-btn',
          cancelButtonClass: 'cancel-delete-btn'
        })
        await request.delete(`/admin/delete/${id}`)
        ElMessage.success('删除成功')
        this.fetchAdmins()
      } catch (err) {
        if (err !== 'cancel') {
          ElMessage.error(err.response?.data?.message || '删除失败')
        }
      }
    },
    formatDate(date) {
      return dayjs(date).format('YYYY-MM-DD HH:mm')
    }
  }
}
</script>

<style scoped>
.form-container {
  background-color: rgba(255, 255, 255, 0.7);
  padding: 20px;
  border-radius: 10px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.form-inline {
  display: flex;
  flex-wrap: wrap;
  gap: 15px;
  align-items: center;
}

.form-item {
  margin-bottom: 0;
}

.form-item :deep(.el-form-item__label) {
  font-size: 16px;
  font-weight: 500;
  color: #2d3748;
}

.add-btn {
  font-size: 16px;
  padding: 12px 20px;
  border-radius: 8px;
  transition: all 0.3s;
}

.add-btn i {
  margin-right: 8px;
}

.admin-table {
  font-size: 16px;
  border-radius: 10px;
  overflow: hidden;
}

.admin-table :deep(th) {
  background-color: #f8fafc;
  font-size: 16px;
  font-weight: 600;
  color: #2d3748;
}

.admin-table :deep(td) {
  padding: 12px 0;
}

.create-time {
  font-family: monospace;
  color: #4a5568;
}

.delete-btn {
  font-size: 15px;
  padding: 10px 15px;
  border-radius: 8px;
  transition: all 0.3s;
}

.delete-btn i {
  margin-right: 8px;
}

/* 对话框按钮样式 */
.confirm-delete-btn {
  background-color: #f56c6c;
  border-color: #f56c6c;
}

.confirm-delete-btn:hover {
  background-color: #e64c4c;
  border-color: #e64c4c;
}

.cancel-delete-btn {
  color: #909399;
}

.cancel-delete-btn:hover {
  color: #606266;
}

@media (max-width: 768px) {
  .form-inline {
    flex-direction: column;
    align-items: stretch;
  }

  .form-item {
    width: 100%;
  }

  .add-btn {
    width: 100%;
    font-size: 15px;
  }

  .admin-table {
    font-size: 14px;
    overflow-x: auto;
  }

  .delete-btn {
    font-size: 14px;
    width: 100%;
  }

  .admin-table :deep(td), 
  .admin-table :deep(th) {
    padding: 10px 5px;
  }
}

</style>
