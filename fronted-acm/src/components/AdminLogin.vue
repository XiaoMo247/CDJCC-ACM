<template>
  <div class="auth-form">
    <form @submit.prevent="handleLogin">
      <div class="form-group">
        <label for="username">
          <i class="fas fa-user-shield"></i> 管理员ID
        </label>
        <input type="text" id="username" v-model="username" placeholder="请输入管理员ID" required />
      </div>
      <div class="form-group">
        <label for="password">
          <i class="fas fa-lock"></i> 密码
        </label>
        <input type="password" id="password" v-model="password" placeholder="请输入密码" required />
      </div>
      <button type="submit" class="submit-btn admin-btn">
        <i class="fas fa-sign-in-alt"></i> 管理员登录
      </button>
    </form>
  </div>
</template>

<script>
import request from '@/utils/request'
import { ElMessage } from 'element-plus'
import emitter from '@/utils/eventBus'

export default {
  name: 'AdminLogin',
  data() {
    return {
      username: '',
      password: ''
    }
  },
  methods: {
    async handleLogin() {
      try {
        const res = await request.post('/admin/login', {
          username: this.username,
          password: this.password
        })

        const { token } = res.data
        localStorage.setItem('admin_token', token)
        emitter.emit('loginChange', { role: 'admin' })

        ElMessage.success('登录成功')
        this.$router.push('/admin/dashboard')
      } catch (error) {
        console.error('登录失败:', error)
        ElMessage.error(error.response?.data?.message || '登录失败')
      }
    }
  }
}
</script>


<style scoped>
.auth-form {
  padding: 0 10px;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  color: #555;
  font-size: 14px;
}

.form-group input {
  width: 95%;
  padding: 12px 15px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 15px;
  transition: border-color 0.3s;
}

.form-group input:focus {
  border-color: #3498db;
  outline: none;
  box-shadow: 0 0 0 2px rgba(52, 152, 219, 0.2);
}

.submit-btn {
  width: 100%;
  padding: 12px;
  background-color: #3498db;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 15px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.submit-btn i {
  margin-right: 8px;
}
</style>
