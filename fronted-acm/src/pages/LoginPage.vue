<template>
  <div class="login-page">
    <div class="login-container">
      <div class="login-card">
        <!-- 顶部装饰元素 -->
        <div class="login-header">
          <div class="logo">
            <img src="@/assets/acm-logo.jpg" alt="ACM Logo">
          </div>
          <h1 class="title">ACM 集训队登录</h1>
          <p class="subtitle">成都锦城学院 · ACM集训队</p>
        </div>
  
        <!-- 登录方式选择 -->
        <div class="login-tabs">
          <button 
            class="tab-btn" 
            :class="{ active: loginType === 'student' }"
            @click="loginType = 'student'"
          >
            <i class="fas fa-user-graduate"></i> 学号登录
          </button>
          <button 
            class="tab-btn" 
            :class="{ active: loginType === 'member' }"
            @click="loginType = 'member'"
          >
            <i class="fas fa-users"></i> 队员登录
          </button>
          <button 
            class="tab-btn" 
            :class="{ active: loginType === 'admin' }"
            @click="loginType = 'admin'"
          >
            <i class="fas fa-user-shield"></i> 管理员
          </button>
        </div>
  
        <!-- 动态加载登录组件 -->
        <transition name="fade" mode="out-in">
          <div class="login-content">
            <StudentLogin v-if="loginType === 'student'" />
            <MemberLogin v-if="loginType === 'member'" />
            <AdminLogin v-if="loginType === 'admin'" />
          </div>
        </transition>
  
        <!-- 底部链接 -->
        <div class="login-footer">
          <router-link to="/forgot-password" class="footer-link">
            <i class="fas fa-question-circle"></i> 忘记密码?
          </router-link>
        </div>
      </div>
    </div>
  </div>
  </template>

<script>
import StudentLogin from '../components/StudentLogin.vue';
import MemberLogin from '../components/MemberLogin.vue';
import AdminLogin from '../components/AdminLogin.vue';

export default {
  name: 'LoginPage',
  components: {
    StudentLogin,
    MemberLogin,
    AdminLogin
  },
  data() {
    return {
      loginType: 'student'
    };
  }
};
</script>

<style scoped>
.login-page {
  min-height: calc(100vh - 64px); /* 如果 NavBar 高度为 64px，可根据实际值调整 */
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  padding: 40px 20px 20px; /* 上 padding 稍微大一点，底部小一些，更居中 */
  box-sizing: border-box;
}


.login-container {
  width: 100%;
  max-width: 480px;
  margin-top: -100px; /* 向上移动减少空白 */
}

.login-card {
  background: white;
  border-radius: 12px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
  overflow: hidden;
  padding: 30px;
  transition: all 0.3s ease;
}

.login-header {
  text-align: center;
  margin-bottom: 24px;
}

.logo img {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  object-fit: cover;
  border: 3px solid #eaeaea;
  margin-bottom: 16px;
}

.title {
  font-size: 24px;
  color: #2c3e50;
  margin-bottom: 8px;
  font-weight: 600;
}

.subtitle {
  font-size: 14px;
  color: #7f8c8d;
  margin-bottom: 0;
}

.login-tabs {
  display: flex;
  justify-content: space-between;
  margin-bottom: 24px;
  border-bottom: 1px solid #eee;
}

.tab-btn {
  flex: 1;
  padding: 12px 0;
  background: none;
  border: none;
  cursor: pointer;
  font-size: 14px;
  color: #7f8c8d;
  position: relative;
  transition: all 0.3s ease;
}

.tab-btn i {
  margin-right: 8px;
}

.tab-btn.active {
  color: #3498db;
  font-weight: 500;
}

.tab-btn.active::after {
  content: '';
  position: absolute;
  bottom: -1px;
  left: 0;
  right: 0;
  height: 2px;
  background: #3498db;
}

.login-content {
  min-height: 200px;
}

.login-footer {
  text-align: center;
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px solid #eee;
}

.footer-link {
  color: #7f8c8d;
  text-decoration: none;
  font-size: 13px;
  transition: color 0.3s;
}

.footer-link:hover {
  color: #3498db;
}

.footer-link i {
  margin-right: 5px;
}

/* 过渡动画 */
.fade-enter-active, .fade-leave-active {
  transition: opacity 0.3s;
}
.fade-enter, .fade-leave-to {
  opacity: 0;
}

@media (max-width: 576px) {
  .login-card {
    padding: 20px;
  }
  
  .login-tabs {
    flex-direction: column;
  }
  
  .tab-btn {
    padding: 10px;
    text-align: center;
  }
  
  .tab-btn.active::after {
    bottom: 0;
    left: 25%;
    right: 25%;
  }
}
</style>