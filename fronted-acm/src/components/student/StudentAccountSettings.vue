<template>
    <div class="settings-card">
      <h1><i class="fas fa-user-cog"></i>账户设置</h1>
  
      <div class="form-container">
        <!-- 修改用户名 -->
        <el-form :model="usernameForm" label-width="100px" class="form-inline">
          <el-form-item label="新用户名">
            <el-input v-model="usernameForm.new_username" placeholder="请输入新的用户名" size="large" clearable />
          </el-form-item>
          <el-button type="primary" size="large" @click="submitUsername" class="submit-btn">
            <i class="fas fa-check-circle"></i> 修改用户名
          </el-button>
        </el-form>
  
        <el-divider />
  
        <!-- 修改密码 -->
        <el-form :model="passwordForm" label-width="100px" class="form-inline">
          <el-form-item label="原密码">
            <el-input v-model="passwordForm.old_password" type="password" placeholder="请输入原密码" size="large" show-password />
          </el-form-item>
          <el-form-item label="新密码">
            <el-input v-model="passwordForm.new_password" type="password" placeholder="请输入新密码" size="large" show-password />
          </el-form-item>
          <el-form-item label="确认新密码">
            <el-input v-model="passwordForm.confirm_password" type="password" placeholder="请再次输入新密码" size="large" show-password />
          </el-form-item>
          <el-button type="primary" size="large" @click="submitPassword" class="submit-btn">
            <i class="fas fa-key"></i> 修改密码
          </el-button>
        </el-form>
      </div>
    </div>
  </template>
  
  <script>
  import request from '@/utils/request'
  import { ElMessage } from 'element-plus'
  
  export default {
    name: 'StudentAccountSettings',
    data() {
      return {
        usernameForm: {
          new_username: ''
        },
        passwordForm: {
          old_password: '',
          new_password: '',
          confirm_password: ''
        }
      }
    },
    methods: {
      async submitUsername() {
        if (!this.usernameForm.new_username) {
          ElMessage.warning('请输入新用户名')
          return
        }
        try {
          await request.post('/user/change-username', this.usernameForm)
          ElMessage.success('用户名修改成功')
          // 刷新页面
          window.location.reload()
        } catch (err) {
          ElMessage.error(err.response?.data?.message || '修改失败')
        }
      },
      async submitPassword() {
        const { old_password, new_password, confirm_password } = this.passwordForm
        if (!old_password || !new_password || !confirm_password) {
          ElMessage.warning('请输入完整的密码信息')
          return
        }
        if (new_password !== confirm_password) {
          ElMessage.error('两次输入的新密码不一致')
          return
        }
        try {
          await request.post('/user/change-password', {
            old_password,
            new_password
          })
          ElMessage.success('密码修改成功，请重新登录')
          // 清空 token，跳转登录页
          localStorage.removeItem('token')
          setTimeout(() => {
            this.$router.push('/login')
          }, 1000)
        } catch (err) {
          ElMessage.error(err.response?.data?.message || '修改失败')
        }
      }
    }
  }
  </script>
  
  <style scoped>
  .settings-card {
    background: linear-gradient(135deg, #e0f2fe, #c7d2fe);
    padding: 25px;
    border-radius: 15px;
    box-shadow: 0 6px 15px rgba(0, 0, 0, 0.08);
    color: #1e293b;
    display: flex;
    flex-direction: column;
    gap: 25px;
  }
  
  .settings-card h1 {
    font-size: 24px;
    font-weight: bold;
    margin: 0 0 15px 0;
    display: flex;
    align-items: center;
    color: #1a365d;
  }
  
  .settings-card h1 i {
    color: #4a5568;
    margin-right: 20px;
    width: 20px;
    height: 20px;
  }
  
  .form-container {
    background-color: rgba(255, 255, 255, 0.85);
    padding: 20px;
    border-radius: 10px;
  }
  
  .form-inline {
    display: flex;
    flex-direction: column;
    gap: 15px;
  }
  
  .submit-btn {
    font-size: 16px;
    padding: 12px 20px;
    border-radius: 8px;
    width: fit-content;
  }
  
  .submit-btn i {
    margin-right: 8px;
  }
  </style>
  