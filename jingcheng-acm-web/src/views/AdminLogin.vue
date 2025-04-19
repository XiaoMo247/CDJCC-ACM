<template>
  <div class="admin-login-container">
    <div class="login-card">
      <!-- Logo区域 -->
      <div class="logo-section">
        <img 
          src="@/assets/acm-logo.jpg" 
          alt="ACM集训队Logo"
          class="logo-img"
        />
        <h1 class="login-title">管理员登录</h1>
        <div class="login-subtitle">成都锦城学院ACM集训队后台管理系统</div>
      </div>

      <!-- 登录表单 -->
      <el-form 
        :model="form" 
        label-position="top" 
        @submit.native.prevent="login"
        class="login-form"
      >
        <el-form-item label="用户名" class="form-item">
          <el-input
            v-model="form.username"
            placeholder="请输入管理员账号"
            size="large"
          >
            <template #prefix>
              <el-icon><User /></el-icon>
            </template>
          </el-input>
        </el-form-item>

        <el-form-item label="密码" class="form-item">
          <el-input
            v-model="form.password"
            type="password"
            placeholder="请输入登录密码"
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
          @click="login" 
          class="login-btn"
          size="large"
        >
          登 录
        </el-button>
      </el-form>

      <!-- 底部版权信息 -->
      <div class="footer">
        <div class="copyright">© 2025 成都锦城学院ACM集训队 版权所有</div>
      </div>
    </div>

    <!-- 背景装饰元素 -->
    <div class="decoration-circle circle-1"></div>
    <div class="decoration-circle circle-2"></div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { User, Lock } from '@element-plus/icons-vue'
import { useAdminStore } from '@/stores/admin'

const form = ref({
  username: '',
  password: ''
})

const router = useRouter()
const adminStore = useAdminStore()

const login = async () => {
  if (!form.value.username || !form.value.password) {
    ElMessage.error('请填写用户名和密码')
    return
  }

  try {
    const response = await fetch('http://localhost:8081/api/admin/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(form.value),
    })

    const data = await response.json()

    if (data.token) {
      adminStore.login(data.token)
      ElMessage.success('登录成功')
      router.push('/admin/dashboard')
    } else {
      ElMessage.error('登录失败: ' + (data.msg || '未知错误'))
    }
  } catch (error) {
    ElMessage.error('登录失败: ' + error.message)
  }
}
</script>

<style scoped>
.admin-login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #f0f9ff 0%, #e0f2fe 100%);
  position: relative;
  overflow: hidden;
}

.login-card {
  width: 420px;
  padding: 40px;
  background: white;
  border-radius: 16px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.08);
  z-index: 1;
  position: relative;
}

.logo-section {
  text-align: center;
  margin-bottom: 36px;
}

.logo-img {
  height: 80px;
  width: 80px;
  object-fit: contain;
  margin-bottom: 16px;
}

.login-title {
  font-size: 24px;
  font-weight: 600;
  color: #1e293b;
  margin-bottom: 8px;
}

.login-subtitle {
  font-size: 14px;
  color: #64748b;
}

.login-form {
  margin-top: 30px;
}

.form-item {
  margin-bottom: 24px;
}

.form-item :deep(.el-form-item__label) {
  font-weight: 500;
  color: #475569;
  padding-bottom: 8px;
}

.login-btn {
  width: 100%;
  margin-top: 10px;
  height: 48px;
  font-size: 16px;
  font-weight: 500;
  letter-spacing: 2px;
  background: linear-gradient(90deg, #3b82f6 0%, #6366f1 100%);
  border: none;
  transition: all 0.3s ease;
}

.login-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.3);
}

.footer {
  margin-top: 40px;
  text-align: center;
}

.copyright {
  font-size: 12px;
  color: #94a3b8;
}

/* 背景装饰元素 */
.decoration-circle {
  position: absolute;
  border-radius: 50%;
  background: linear-gradient(135deg, rgba(56, 189, 248, 0.1) 0%, rgba(99, 102, 241, 0.1) 100%);
}

.circle-1 {
  width: 300px;
  height: 300px;
  top: -100px;
  left: -100px;
}

.circle-2 {
  width: 400px;
  height: 400px;
  bottom: -150px;
  right: -150px;
}
</style>