<template>
    <nav class="navbar">
      <!-- 品牌标志区域 -->
      <div class="navbar-brand" @click="$router.push('/')">
        <img src="@/assets/acm-logo.jpg" alt="ACM集训队Logo" class="brand-logo" />
        <div class="brand-text">
          <div class="school-name">成都锦城学院</div>
          <div class="team-name">ACM集训队</div>
        </div>
        <img src="@/assets/acm-icon.png" alt="ACM图标" class="acm-icon" />
      </div>
  
      <!-- 自适应间隔 -->
      <div class="spacer1"></div>
  
      <!-- 导航链接 -->
      <ul class="nav-links">
        <li><router-link to="/">首页</router-link></li>
        <li><router-link to="/knowledge">知识库</router-link></li>
        <li><router-link to="/competition">比赛页面</router-link></li>
        <li><router-link to="/ranking">排行榜</router-link></li>
        <li><router-link to="/faq">常见问题</router-link></li>
      </ul>
  
      <!-- 自适应间隔 -->
      <div class="spacer2"></div>
  
      <!-- 登录 / 控制面板 + 退出 -->
      <div v-if="isLoggedIn" class="auth-buttons">
        <router-link
          v-if="role === 'admin'"
          to="/admin/dashboard"
          class="login-btn"
        >管理员控制面板</router-link>
  
        <router-link
          v-else-if="role === 'member'"
          to="/member/dashboard"
          class="login-btn"
        >队员控制面板</router-link>
  
        <router-link
          v-else
          to="/student/dashboard"
          class="login-btn"
        >学生控制面板</router-link>
  
        <button class="logout-btn" @click="handleLogout">退出</button>
      </div>
  
      <router-link v-else to="/login" class="login-btn">登录</router-link>
    </nav>
  </template>
  
  <script>
  import emitter from '@/utils/eventBus'
  
  export default {
    name: 'NavBar',
    data() {
      return {
        role: null,
        isLoggedIn: false
      }
    },
    mounted() {
      this.syncLoginState()
  
      // 监听登录状态变化
      emitter.on('loginChange', this.syncLoginState)
    },
    beforeUnmount() {
      emitter.off('loginChange', this.syncLoginState)
    },
    methods: {
      syncLoginState() {
        const adminToken = localStorage.getItem('admin_token')
        const userToken = localStorage.getItem('user_token')
  
        if (adminToken) {
          this.role = 'admin'
          this.isLoggedIn = true
        } else if (userToken) {
          const userRole = localStorage.getItem('user_role') || 'user'
          this.role = userRole
          this.isLoggedIn = true
        } else {
          this.role = null
          this.isLoggedIn = false
        }
      },
      handleLogout() {
        localStorage.removeItem('admin_token')
        localStorage.removeItem('user_token')
        localStorage.removeItem('user_role')
        this.syncLoginState()
        emitter.emit('loginChange', { role: null })
        this.$router.push('/login')
      }
    }
  }
  </script>

<style scoped>
/* 基础导航栏样式 */
.navbar {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 80px;
  background: linear-gradient(90deg, #f8fafc, #e2e8f0, #cbd5e1);
  color: #1e293b;
  display: flex;
  align-items: center;
  padding: 0 2rem;
  z-index: 1000;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.04),
              0 2px 4px -1px rgba(0, 0, 0, 0.02);
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
}

/* 自适应间隔 */
.spacer1 {
    flex: 3;
}

.spacer2 {
    flex: 4;
}

/* 导航链接样式 */
.nav-links {
    list-style: none;
    display: flex;
    gap: 1.5rem;
    margin: 0;
    padding: 0;
}

.nav-links li {
    position: relative;
    padding: 0.5rem;
}

.nav-links a {
    color: #1e293b;
    text-decoration: none;
    font-weight: 500;
    font-size: 1rem;
    padding: 0.5rem 0;
    transition: all 0.3s ease;
    position: relative;
}

.nav-links a:not(.login-btn):hover {
    color: #2563eb;
}

.nav-links a:not(.login-btn)::after {
    content: '';
    position: absolute;
    bottom: 0;
    left: 0;
    width: 0;
    height: 2px;
    background-color: #2563eb;
    transition: width 0.3s ease;
}

.nav-links a:not(.login-btn):hover::after {
    width: 100%;
}

/* 登录按钮特殊样式 */
.login-btn {
    position: relative;
    background: linear-gradient(135deg, #3B82F6, #2563EB);
    color: white !important;
    padding: 0.75rem 1.75rem !important;
    border-radius: 12px;
    font-weight: 600 !important;
    font-size: 1rem;
    text-decoration: none;
    border: none;
    cursor: pointer;
    overflow: hidden;
    box-shadow: 
        0 4px 6px rgba(59, 130, 246, 0.25),
        0 1px 3px rgba(59, 130, 246, 0.1),
        inset 0 1px 0 rgba(255, 255, 255, 0.2);
    transition: all 0.4s cubic-bezier(0.175, 0.885, 0.32, 1.275);
    z-index: 1;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
}

.login-btn::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: linear-gradient(135deg, #2563EB, #1D4ED8);
    opacity: 0;
    transition: opacity 0.4s ease;
    z-index: -1;
}

.login-btn:hover {
    transform: translateY(-3px) scale(1.02);
    box-shadow: 
        0 10px 20px rgba(59, 130, 246, 0.3),
        0 6px 6px rgba(59, 130, 246, 0.2),
        inset 0 1px 0 rgba(255, 255, 255, 0.3);
}

.login-btn:hover::before {
    opacity: 1;
}

.login-btn:active {
    transform: translateY(1px) scale(0.98);
    box-shadow: 
        0 2px 5px rgba(59, 130, 246, 0.2),
        inset 0 1px 0 rgba(255, 255, 255, 0.1);
}

.login-btn i {
    transition: transform 0.3s ease;
}

.login-btn:hover i {
    transform: translateX(3px);
}


/* 品牌标志区域 */
.navbar-brand {
    display: flex;
    align-items: center;
    cursor: pointer;
    transition: all 0.3s ease;
    border-radius: 10px;
}


.navbar-brand:hover {
    background: rgba(0, 0, 0, 0.1);
}

.brand-logo {
    height: 50px;
    width: 50px;
    margin-right: 0.75rem;
    object-fit: cover;
    border-radius: 50%;
    filter: drop-shadow(0 0 6px rgba(100, 116, 139, 0.2));
    transition: transform 0.3s ease;
}


.navbar-brand:hover .brand-logo {
    transform: scale(1.05);
}

.brand-text {
    line-height: 1.2;
}

.school-name {
    font-size: 14px;
    color: #475569;
    font-weight: 500;
    display: flex;
    margin-left: 5px;
    justify-content: flex-start;
    text-align: left;
}

.team-name {
    font-size: 1.0625rem;
    font-weight: 700;
    color: #1e293b;
    letter-spacing: 0.025rem;
    text-shadow: 0 0 6px rgba(100, 116, 139, 0.2);
    transition: text-shadow 0.3s ease;
}

.navbar-brand:hover .team-name {
    text-shadow: 0 0 8px rgba(100, 116, 139, 0.4);
}

.acm-icon {
    width: 80px;
    height: 80px;
    margin-left: 8px;
    filter: drop-shadow(0 0 2px rgba(0, 0, 0, 0.2));
}

.auth-buttons {
  display: flex;
  align-items: center;
  gap: 10px;
}

.logout-btn {
    position: relative;
    background: linear-gradient(135deg, #EF4444, #DC2626);
    color: white;
    padding: 0.75rem 1.5rem;
    border-radius: 12px;
    font-weight: 600;
    font-size: 1rem;
    border: none;
    cursor: pointer;
    overflow: hidden;
    box-shadow: 
        0 4px 6px rgba(239, 68, 68, 0.2),
        0 1px 3px rgba(239, 68, 68, 0.1),
        inset 0 1px 0 rgba(255, 255, 255, 0.15);
    transition: all 0.4s cubic-bezier(0.175, 0.885, 0.32, 1.275);
    z-index: 1;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
}

.logout-btn::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: linear-gradient(135deg, #DC2626, #B91C1C);
    opacity: 0;
    transition: opacity 0.4s ease;
    z-index: -1;
}

.logout-btn::after {
    content: '';
    position: absolute;
    top: 50%;
    left: 50%;
    width: 5px;
    height: 5px;
    background: rgba(255, 255, 255, 0.5);
    opacity: 0;
    border-radius: 100%;
    transform: scale(1, 1) translate(-50%, -50%);
    transform-origin: 50% 50%;
}

.logout-btn:hover {
    transform: translateY(-3px) scale(1.02);
    box-shadow: 
        0 10px 20px rgba(239, 68, 68, 0.25),
        0 6px 6px rgba(239, 68, 68, 0.2),
        inset 0 1px 0 rgba(255, 255, 255, 0.25);
}

.logout-btn:hover::before {
    opacity: 1;
}

.logout-btn:focus:not(:active)::after {
    animation: ripple 0.6s ease-out;
}

.logout-btn:active {
    transform: translateY(1px) scale(0.98);
    box-shadow: 
        0 2px 5px rgba(239, 68, 68, 0.2),
        inset 0 1px 0 rgba(255, 255, 255, 0.1);
}

.logout-btn i {
    transition: transform 0.3s ease;
}

.logout-btn:hover i {
    transform: translateX(-3px);
}

@keyframes ripple {
    0% {
        transform: scale(0, 0);
        opacity: 0.5;
    }
    100% {
        transform: scale(20, 20);
        opacity: 0;
    }
}


/* 响应式设计 - 在较小屏幕上调整 */
@media (max-width: 768px) {
    .navbar {
        padding: 0 1rem;
        height: 60px;
    }

    .nav-links {
        gap: 0.75rem;
    }

    .brand-logo {
        height: 40px;
        width: 40px;
    }

    .school-name {
        font-size: 0.75rem;
    }

    .team-name {
        font-size: 0.9375rem;
    }
}
</style>