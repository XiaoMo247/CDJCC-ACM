<template>
  <div class="auth-form">
    <form @submit.prevent="handleLogin">
      <div class="form-group">
        <label for="student_number">
          <i class="fas fa-id-card"></i> 学号
        </label>
        <input 
          type="text" 
          id="student_number" 
          v-model="student_number" 
          placeholder="请输入学号"
          required
        />
      </div>
      <div class="form-group">
        <label for="password">
          <i class="fas fa-lock"></i> 密码
        </label>
        <input 
          type="password" 
          id="password" 
          v-model="password" 
          placeholder="请输入密码"
          required
        />
      </div>
      <button type="submit" class="submit-btn">
        <i class="fas fa-sign-in-alt"></i> 登录
      </button>
    </form>
  </div>
</template>

<script>
import request from '@/utils/request'
import { ElMessage } from 'element-plus'
import emitter from '@/utils/eventBus'

export default {
  name: 'StudentLogin',
  data() {
    return {
      student_number: '',
      password: ''
    }
  },
  methods: {
    async handleLogin() {
      try {
        const res = await request.post('/user/login', {
          student_number: this.student_number,
          password: this.password
        })

        const token = res.data.token
        localStorage.setItem('user_token', token)
        emitter.emit('loginChange', { role: 'user' })

        ElMessage.success('登录成功！')
        this.$router.push('/student/dashboard')
      } catch (error) {
        console.error('登录失败:', error)
        ElMessage.error(error.response?.data?.message || '登录失败，请检查学号或密码')
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
  width: 100%;
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

.submit-btn:hover {
  background-color: #2980b9;
}

.submit-btn i {
  margin-right: 8px;
}
</style>