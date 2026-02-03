<template>
    <div class="member-dashboard">
      <!-- 移动端菜单按钮 -->
      <div class="mobile-menu-btn" @click="toggleSidebar">
        <i class="fas fa-bars"></i>
      </div>
      
      <!-- 遮罩层 -->
      <div class="sidebar-overlay" :class="{ 'active': sidebarActive }" @click="toggleSidebar"></div>
      
      <!-- 侧边栏 -->
      <div class="sidebar" :class="{ 'active': sidebarActive }" @touchstart="handleTouchStart" 
           @touchmove="handleTouchMove" @touchend="handleTouchEnd">
        <h3 class="sidebar-title">
          <i class="fas fa-user-graduate"></i> 队员中心 : {{ Username }}
        </h3>
        <ul>
          <li class="active">
            <i class="fas fa-user-cog"></i> 账户设置
          </li>
        </ul>
      </div>
  
      <div class="content" :class="{ 'sidebar-active': sidebarActive }">
        <MemberAccountSettings @updated="fetchUsername" />
      </div>
    </div>
  </template>
  
  <script>
  import MemberAccountSettings from '@/components/member/MemberAccountSettings.vue'
  import request from '@/utils/request'
  import { ElMessage } from 'element-plus'
  import { isLoggedIn } from '@/utils/tokenManager'
  
  export default {
    name: 'MemberDashboard',
    components: {
      MemberAccountSettings
    },
    data() {
      return {
        Username: '加载中...',
        sidebarActive: false,
        touchStartX: 0,
        touchEndX: 0
      }
    },
    mounted() {
      this.fetchUsername()
    },
    methods: {
      async fetchUsername() {
        try {
          if (!isLoggedIn()) {
            ElMessage.error('未登录，请重新登录')
            this.Username = '未登录'
            return
          }

          const res = await request.get('/student/info')
          const data = res.data?.data || {}
          this.Username = data.Username || data.username || data.name || '未知队员'
        } catch (err) {
          ElMessage.error('获取用户名失败')
          this.Username = '未登录'
        }
      },
      toggleSidebar() {
        this.sidebarActive = !this.sidebarActive
      },
      handleTouchStart(e) {
        this.touchStartX = e.changedTouches[0].screenX
      },
      handleTouchMove(e) {
        e.preventDefault()
      },
      handleTouchEnd(e) {
        this.touchEndX = e.changedTouches[0].screenX
        this.handleSwipe()
      },
      handleSwipe() {
        const diff = this.touchEndX - this.touchStartX
        // 从左向右滑动，且滑动距离大于50px
        if (diff > 50 && !this.sidebarActive) {
          this.sidebarActive = true
        }
        // 从右向左滑动，且滑动距离大于50px
        else if (diff < -50 && this.sidebarActive) {
          this.sidebarActive = false
        }
      }
    }
  }
  </script>
  
  <style scoped>
  .member-dashboard {
    display: flex;
    height: 100%; /* 使用100%适配.main-content的内容区域 */
    max-height: 100%; /* 确保不会超出 */
    background-color: #f8fafc;
    position: relative;
    overflow: hidden; /* 禁止Dashboard层级滚动 */
  }
  
  /* 侧边栏样式 - 保留原配色 */
  .sidebar {
    width: 240px;
    background: linear-gradient(180deg, #e2e8f0, #cbd5e1);
    padding: 20px;
    box-shadow: 2px 0 6px rgba(0, 0, 0, 0.05);
    display: flex;
    flex-direction: column;
    transition: transform 0.3s ease;
    z-index: 100;
    position: fixed;
    height: 100%;
    color: #1e293b;
    overflow-y: auto;
    scrollbar-width: none; /* Firefox */
    -ms-overflow-style: none; /* IE 10+ */
  }

  /* 隐藏侧边栏滚动条 */
  .sidebar::-webkit-scrollbar {
    width: 0;
    height: 0;
  }
  
  .sidebar-title {
    font-size: 18px;
    font-weight: bold;
    color: #1e293b;
    margin-bottom: 20px;
    display: flex;
    align-items: center;
  }
  
  .sidebar-title i {
    margin-right: 8px;
  }
  
  .sidebar ul {
    list-style: none;
    padding: 0;
    margin: 0;
    flex: 1;
  }
  
  .sidebar li {
    padding: 12px 16px;
    margin-bottom: 10px;
    border-radius: 6px;
    cursor: default;
    color: white;
    font-size: 15px;
    display: flex;
    align-items: center;
    background-color: #38bdf8;
  }
  
  .sidebar li i {
    margin-right: 12px;
    width: 20px;
    height: 20px;
    display: inline-block;
    text-align: center;
  }
  
  /* 内容区域 */
  .content {
    flex: 1;
    padding: 24px;
    overflow-y: auto;
    background-color: #f8fafc;
    transition: margin-left 0.3s ease;
    margin-left: 240px;
    height: 100%; /* 使用100%而不是min-height */
    max-height: 100%; /* 确保高度不超过父容器 */
    scrollbar-width: none; /* Firefox */
    -ms-overflow-style: none; /* IE 10+ */
  }

  /* 隐藏滚动条但保留滚动功能 */
  .content::-webkit-scrollbar {
    width: 0;
    height: 0;
  }
  
  /* 移动端菜单按钮 */
  .mobile-menu-btn {
    display: none;
    position: fixed;
    top: 15px;
    left: 15px;
    z-index: 90;
    background: #38bdf8;
    color: white;
    width: 40px;
    height: 40px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: 0 2px 5px rgba(0, 0, 0, 0.2);
    cursor: pointer;
  }
  
  /* 遮罩层 */
  .sidebar-overlay {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: rgba(0, 0, 0, 0.5);
    z-index: 99;
    opacity: 0;
    visibility: hidden;
    transition: all 0.3s ease;
  }
  
  .sidebar-overlay.active {
    opacity: 1;
    visibility: visible;
  }
  
  /* 响应式设计 */
  @media (max-width: 768px) {
    .sidebar {
      transform: translateX(-100%);
    }
    
    .sidebar.active {
      transform: translateX(0);
    }
    
    .content {
      margin-left: 0;
    }
    
    .content.sidebar-active {
      margin-left: 0;
    }
    
    .mobile-menu-btn {
      display: flex;
    }
  }
  
  /* 美化滚动条 */
  ::-webkit-scrollbar {
    width: 6px;
  }
  
  ::-webkit-scrollbar-track {
    background: #f1f1f1;
  }
  
  ::-webkit-scrollbar-thumb {
    background: #c1c1c1;
    border-radius: 3px;
  }
  
  ::-webkit-scrollbar-thumb:hover {
    background: #a8a8a8;
  }
  </style>
