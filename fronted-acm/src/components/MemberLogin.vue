<template>
  <div class="auth-form">
    <form @submit.prevent="handleLogin">
      <div class="form-group">
        <label for="memberId">
          <i class="fas fa-user-tie"></i> 队员ID
        </label>
        <input 
          type="text" 
          id="memberId" 
          v-model="memberId" 
          placeholder="请输入队员ID"
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
export default {
  name: 'MemberLogin',
  data() {
    return {
      memberId: '',
      password: ''
    }
  },
  methods: {
    handleLogin() {
      // 处理队员登录逻辑
      console.log('队员登录', this.memberId, this.password)
      // 这里可以添加API调用
    }
  }
}
</script>

<style scoped>
/* 样式与StudentLogin相同，保持一致性 */
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