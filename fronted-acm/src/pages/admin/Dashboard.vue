<template>
  <div class="admin-dashboard">
    <!-- 移动端顶部导航栏 -->
    <div class="mobile-navbar" v-if="isMobile">
      <button class="hamburger" @click="toggleSidebar">
        <i class="fas fa-bars"></i>
      </button>
      <h3 class="mobile-title"><i class="fas fa-tools"></i> 控制面板</h3>
    </div>

    <!-- 侧边栏 -->
    <div class="sidebar" :class="{ 'mobile-active': sidebarActive }">
      <h3 class="sidebar-title"><i class="fas fa-tools"></i> 控制面板</h3>
      <ul>
        <li :class="{ active: currentTab === 'home' }" @click="changeTab('home')">
          <i class="fas fa-home"></i> 管理主页
        </li>
        <li :class="{ active: currentTab === 'announcement' }" @click="changeTab('announcement')">
          <i class="fas fa-bullhorn"></i> 公告管理
        </li>
        <li :class="{ active: currentTab === 'admin' }" @click="changeTab('admin')">
          <i class="fas fa-user-shield"></i> 管理员账户
        </li>
        <li :class="{ active: currentTab === 'join' }" @click="changeTab('join')">
          <i class="fas fa-clipboard-list"></i> 申请表管理
        </li>
        <li :class="{ active: currentTab === 'user' }" @click="changeTab('user')">
          <i class="fas fa-users"></i> 用户管理
        </li>
        <li :class="{ active: currentTab === 'folder' }" @click="changeTab('folder')">
          <i class="fas fa-copy"></i> 课件管理
        </li>
        <li :class="{ active: currentTab === 'student' }" @click="changeTab('student')">
          <i class="fas fa-user-graduate"></i> 学生管理
        </li>
        <li :class="{ active: currentTab === 'contest' }" @click="changeTab('contest')">
          <i class="fas fa-trophy"></i> 比赛管理
        </li>
        <li :class="{ active: currentTab === 'faq' }" @click="changeTab('faq')">
          <i class="fas fa-question-circle"></i> FAQ管理
        </li>
        <li :class="{ active: currentTab === 'slider' }" @click="changeTab('slider')">
          <i class="fas fa-images"></i> 轮播图管理
        </li>
        <li :class="{ active: currentTab === 'honorwallmanage' }" @click="changeTab('honorwallmanage')">
          <i class="fas fa-trophy"></i> 荣誉墙管理
        </li>
      </ul>
    </div>

    <!-- 内容区域 -->
    <div class="content" :class="{ 'mobile-padding': isMobile }" @touchstart="handleTouchStart"
      @touchmove="handleTouchMove" @touchend="handleTouchEnd">
      <DashboardHome v-if="currentTab === 'home'" @change-tab="changeTab" />
      <AnnouncementManage v-if="currentTab === 'announcement'" />
      <AdminAccountManage v-if="currentTab === 'admin'" />
      <JoinFormManage v-if="currentTab === 'join'" />
      <UserManage v-if="currentTab === 'user'" />
      <FolderManage v-if="currentTab === 'folder'" />
      <StudentManage v-if="currentTab === 'student'" />
      <ContestManage v-if="currentTab === 'contest'" />
      <FAQManage v-if="currentTab === 'faq'" />
      <Slider v-if="currentTab === 'slider'" />
      <HonorWallManage v-if="currentTab === 'honorwallmanage'" />
    </div>

    <!-- 移动端遮罩层 -->
    <div class="sidebar-overlay" v-if="isMobile && sidebarActive" @click="sidebarActive = false"></div>
  </div>
</template>


<script>
import DashboardHome from '@/components/admin/DashboardHome.vue'
import AnnouncementManage from '@/components/admin/AnnouncementManage.vue'
import AdminAccountManage from '@/components/admin/AdminAccountManage.vue'
import JoinFormManage from '@/components/admin/JoinFormManage.vue'
import UserManage from '@/components/admin/UserManage.vue'
import FolderManage from '@/components/admin/FolderManage.vue'
import StudentManage from '@/components/admin/StudentManage.vue'
import ContestManage from '@/components/admin/ContestManage.vue'
import FAQManage from '@/components/admin/FAQManage.vue'
import Slider from '@/components/admin/Slider.vue'
import HonorWallManage from '@/components/admin/HonorWallManage.vue'

export default {
  name: 'AdminDashboard',
  components: {
    DashboardHome,
    AnnouncementManage,
    AdminAccountManage,
    JoinFormManage,
    UserManage,
    FolderManage,
    StudentManage,
    ContestManage,
    FAQManage,
    Slider,
    HonorWallManage,
  },
  data() {
    return {
      currentTab: 'home',
      isMobile: false,
      sidebarActive: false,
      touchStartX: 0,
      touchStartY: 0,
      touchThreshold: 50
    }
  },
  methods: {
    changeTab(tab) {
      this.currentTab = tab
      if (this.isMobile) {
        this.sidebarActive = false
      }
    },
    toggleSidebar() {
      this.sidebarActive = !this.sidebarActive
    },
    checkMobile() {
      this.isMobile = window.innerWidth <= 768
      if (!this.isMobile) {
        this.sidebarActive = false
      }
    },
    handleTouchStart(e) {
      if (!this.isMobile || this.sidebarActive) return
      this.touchStartX = e.touches[0].clientX
      this.touchStartY = e.touches[0].clientY
    },
    handleTouchMove(e) {
      if (!this.touchStartX || this.sidebarActive) return
      const deltaX = e.touches[0].clientX - this.touchStartX
      const deltaY = e.touches[0].clientY - this.touchStartY
      if (Math.abs(deltaX) > Math.abs(deltaY) && deltaX > this.touchThreshold) {
        this.sidebarActive = true
      }
    },
    handleTouchEnd() {
      this.touchStartX = 0
      this.touchStartY = 0
    }
  },
  mounted() {
    this.checkMobile()
    window.addEventListener('resize', this.checkMobile)
  },
  beforeUnmount() {
    window.removeEventListener('resize', this.checkMobile)
  }
}
</script>


<style scoped>
.admin-dashboard {
  display: flex;
  height: 100vh;
  background-color: #f8fafc;
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
  transition: all 0.2s;
}

.sidebar li i {
  margin-right: 12px;
  width: 20px;
  height: 20px;
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
}

.mobile-navbar {
  display: none;
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 60px;
  background: linear-gradient(90deg, #e2e8f0, #cbd5e1);
  padding: 0 15px;
  align-items: center;
  z-index: 90;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
}

.mobile-title {
  font-size: 16px;
  font-weight: bold;
  color: #1e293b;
  margin: 0 auto;
  display: flex;
  align-items: center;
}

.mobile-title i {
  margin-right: 8px;
}

.hamburger {
  background: none;
  border: none;
  font-size: 20px;
  color: #1e293b;
  cursor: pointer;
  padding: 10px;
  z-index: 100;
}

.sidebar-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  z-index: 99;
}

.mobile-padding {
  padding-top: 80px;
}

/* 响应式 */
@media (max-width: 768px) {
  .mobile-navbar {
    display: flex;
  }

  .sidebar {
    position: fixed;
    top: 0;
    left: 0;
    bottom: 0;
    transform: translateX(-100%);
  }

  .sidebar.mobile-active {
    transform: translateX(0);
  }

  .content {
    width: 100%;
    padding: 15px;
    touch-action: pan-y;
  }
}

@media (max-width: 480px) {
  .sidebar {
    width: 220px;
    padding: 15px;
  }

  .sidebar li {
    padding: 10px 12px;
    font-size: 14px;
  }

  .content {
    padding: 10px;
  }
}
</style>
