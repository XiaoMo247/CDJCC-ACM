<template>
  <div class="auth-form">
    <form @submit.prevent="handleLogin">
      <div class="form-group">
        <label for="student_number">
          <i class="fas fa-id-card"></i> 学号
        </label>
        <input type="text" id="student_number" v-model="student_number" placeholder="请输入学号" required :disabled="loading" />
      </div>
      <div class="form-group">
        <label for="password">
          <i class="fas fa-lock"></i> 密码
        </label>
        <input type="password" id="password" v-model="password" placeholder="请输入密码" required :disabled="loading" />
      </div>
      <button type="submit" class="submit-btn" :disabled="loading" :aria-busy="loading.toString()">
        <i class="fas fa-sign-in-alt"></i> {{ loading ? '登录中...' : '登录' }}
      </button>
    </form>
  </div>
</template>

<script>
import request from '@/utils/request'
import { ElMessage } from 'element-plus'
import emitter from '@/utils/eventBus'
import { saveToken, startTokenRefresh } from '@/utils/tokenManager'

export default {
  name: 'StudentLogin',
  data() {
    return {
      student_number: '',
      password: '',
      loading: false
    }
  },
  methods: {
    async handleLogin() {
      if (this.loading) return

      this.loading = true
      try {
        const res = await request.post('/user/login', {
          student_number: this.student_number,
          password: this.password
        })

        // 后端应该返回: { token, user: { id, student_number, role: 'student', ... } }
        const { token, user } = res.data

        // 如果后端没有返回 user 对象，临时构造一个
        const userInfo = user || {
          student_number: this.student_number,
          role: 'student'
        }

        // 使用统一的 token 存储
        saveToken(token, userInfo)

        emitter.emit('loginChange', { role: 'student', user: userInfo })

        // 启动 token 自动刷新
        startTokenRefresh()

        ElMessage.success('登录成功！')

        // 获取登录前访问的页面路径
        const redirectPath = sessionStorage.getItem('redirect_after_login')

        if (redirectPath) {
          // 清除已保存的路径
          sessionStorage.removeItem('redirect_after_login')
          // 跳转到之前的页面
          this.$router.push(redirectPath)
        } else {
          // 默认跳转到学生首页
          this.$router.push('/student/dashboard')
        }
      } catch (error) {
        console.error('登录失败:', error)
        ElMessage.error(error.response?.data?.message || '登录失败，请检查学号或密码')
      } finally {
        this.loading = false
      }
    }
  }
}
</script>
