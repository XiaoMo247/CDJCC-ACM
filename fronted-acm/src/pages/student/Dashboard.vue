<template>
  <div
    class="student-dashboard"
    @touchstart="handleTouchStart"
    @touchmove="handleTouchMove"
    @touchend="handleTouchEnd"
  >
    <div class="sidebar" :class="{ hidden: !sidebarVisible, 'mobile-view': isMobile }">
      <h3 class="sidebar-title">
        <i class="fas fa-user-graduate"></i> 学生中心 : {{ Username }}
      </h3>
      <ul>
        <li :class="{ active: currentTab === 'join' }" @click="switchTab('join')">
          <i class="fas fa-clipboard-list"></i> 申请表管理
        </li>
        <li :class="{ active: currentTab === 'account' }" @click="switchTab('account')">
          <i class="fas fa-user-cog"></i> 账户设置
        </li>
      </ul>
      <button class="sidebar-toggle" @click="toggleSidebar" v-if="isMobile">
        <i :class="sidebarVisible ? 'fas fa-times' : 'fas fa-bars'"></i>
      </button>
    </div>

    <div class="content" :class="{ 'sidebar-hidden': !sidebarVisible && isMobile }">
      <div class="mobile-header" v-if="isMobile">
        <button @click="toggleSidebar">
          <i class="fas fa-bars"></i>
        </button>
        <h3>{{ currentTab === 'join' ? '申请表管理' : '账户设置' }}</h3>
      </div>
      <StudentJoinForm v-if="currentTab === 'join'" />
      <StudentAccountSettings v-if="currentTab === 'account'" @updated="fetchUsername" />
    </div>
  </div>
</template>

<script>
import StudentJoinForm from '@/components/student/StudentJoinForm.vue'
import StudentAccountSettings from '@/components/student/StudentAccountSettings.vue'
import request from '@/utils/request'
import { ElMessage } from 'element-plus'

export default {
  name: 'StudentDashboard',
  components: {
    StudentJoinForm,
    StudentAccountSettings
  },
  data() {
    return {
      currentTab: 'join',
      Username: '加载中...',
      sidebarVisible: true,
      touchStartX: 0,
      touchEndX: 0,
      isMobile: false
    }
  },
  mounted() {
    this.fetchUsername()
    this.checkMobile()
    window.addEventListener('resize', this.checkMobile)
  },
  beforeUnmount() {
    window.removeEventListener('resize', this.checkMobile)
  },
  methods: {
    async fetchUsername() {
      try {
        const res = await request.get('/user/info')
        this.Username = res.data.data.username || '未知用户'
      } catch (err) {
        ElMessage.error('获取用户名失败')
        this.Username = '未登录'
      }
    },
    checkMobile() {
      this.isMobile = window.innerWidth <= 768
      if (this.isMobile) {
        this.sidebarVisible = false
      } else {
        this.sidebarVisible = true
      }
    },
    toggleSidebar() {
      this.sidebarVisible = !this.sidebarVisible
    },
    switchTab(tab) {
      this.currentTab = tab
      if (this.isMobile) {
        this.sidebarVisible = false
      }
    },
    handleTouchStart(e) {
      this.touchStartX = e.changedTouches[0].clientX
    },
    handleTouchMove(e) {
      this.touchEndX = e.changedTouches[0].clientX
    },
    handleTouchEnd() {
      if (!this.isMobile) return
      
      const deltaX = this.touchEndX - this.touchStartX
      const threshold = 30 // 降低阈值，更适合移动设备

      if (deltaX > threshold) {
        this.sidebarVisible = true
      } else if (deltaX < -threshold) {
        this.sidebarVisible = false
      }
    }
  }
}
</script>

<style scoped>
.student-dashboard {
  display: flex;
  height: 100vh;
  background-color: #f8fafc;
  overflow: hidden;
  position: relative;
}

.sidebar {
  width: 240px;
  background: linear-gradient(180deg, #e2e8f0, #cbd5e1);
  padding: 20px;
  box-shadow: 2px 0 6px rgba(0, 0, 0, 0.05);
  display: flex;
  flex-direction: column;
  transition: transform 0.3s ease;
  z-index: 100;
  position: relative;
}

.sidebar.hidden {
  transform: translateX(-100%);
}

.sidebar.mobile-view {
  position: fixed;
  height: 100%;
  width: 260px;
  left: 0;
  top: 0;
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
  cursor: pointer;
  color: #1e293b;
  font-size: 15px;
  display: flex;
  align-items: center;
  transition: background-color 0.2s;
  line-height: 20px;
}

.sidebar li i {
  margin-right: 12px;
  width: 20px;
  height: 20px;
  display: inline-block;
  text-align: center;
}

.sidebar li:hover {
  background-color: #e0f2fe;
}

.sidebar li.active {
  background-color: #38bdf8;
  color: white;
}

.content {
  flex: 1;
  padding: 24px;
  overflow-y: auto;
  background-color: #f8fafc;
  transition: margin-left 0.3s ease;
}

.content.sidebar-hidden {
  margin-left: 0;
  width: 100%;
}

.sidebar-toggle {
  position: absolute;
  right: -40px;
  top: 20px;
  width: 40px;
  height: 40px;
  background: #e2e8f0;
  border: none;
  border-radius: 0 4px 4px 0;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 2px 0 6px rgba(0, 0, 0, 0.1);
}

.mobile-header {
  display: none;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 15px;
  border-bottom: 1px solid #e2e8f0;
}

.mobile-header button {
  background: none;
  border: none;
  font-size: 20px;
  margin-right: 15px;
  cursor: pointer;
}

.mobile-header h3 {
  margin: 0;
  font-size: 18px;
}

@media (max-width: 768px) {
  .content {
    padding: 15px;
    width: 100%;
  }
  
  .mobile-header {
    display: flex;
  }
  
  .sidebar li {
    padding: 15px 16px;
    font-size: 16px;
  }
  
  .sidebar-title {
    font-size: 16px;
  }
}

@media (max-width: 480px) {
  .sidebar {
    width: 220px;
  }
  
  .content {
    padding: 10px;
  }
}
</style>
