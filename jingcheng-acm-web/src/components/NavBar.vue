<template>
  <el-menu 
    mode="horizontal" 
    router 
    class="acm-navbar"
    :style="navbarStyle"
    text-color="#f8fafc"
    active-text-color="#38bdf8"
  >
    <!-- 左侧Logo区域 -->
    <div class="navbar-brand" @click="$router.push('/')">
      <img 
        src="@/assets/acm-logo.jpg" 
        alt="ACM集训队Logo"
        class="brand-logo"
      >
      <div class="brand-text">
        <div class="school-name">成都锦城学院</div>
        <div class="team-name">ACM集训队</div>
      </div>
    </div>

    <!-- 主导航菜单 -->
    <div class="nav-items">
      <el-menu-item index="/" class="nav-item">
        <el-icon><House /></el-icon>
        <span>首页</span>
      </el-menu-item>
      <el-menu-item index="/about" class="nav-item">
        <el-icon><User /></el-icon>
        <span>关于我们</span>
      </el-menu-item>
      <el-menu-item index="/contests" class="nav-item">
        <el-icon><Trophy /></el-icon>
        <span>比赛</span>
      </el-menu-item>
      <el-menu-item index="/rank" class="nav-item">
        <el-icon><Medal /></el-icon>
        <span>排行榜</span>
      </el-menu-item>
      <el-menu-item index="/join" class="nav-item">
        <el-icon><Connection /></el-icon>
        <span>加入我们</span>
      </el-menu-item>
    </div>

    <!-- 右侧管理员区域 -->
    <div class="admin-area">
      <el-button 
        v-if="!isLoggedIn"
        type="primary"
        class="admin-login-btn"
        @click="$router.push('/admin/login')"
      >
        <el-icon><Lock /></el-icon>
        <span>管理员登录</span>
      </el-button>
      
      <el-dropdown v-else>
        <div class="admin-profile">
          <el-avatar 
            :size="36" 
            src="@/assets/admin-avatar.png"
            class="admin-avatar"
          />
          <div class="admin-info">
            <div class="admin-name">管理员</div>
            <div class="admin-role">超级管理员</div>
          </div>
          <el-icon class="dropdown-icon"><ArrowDown /></el-icon>
        </div>
        <template #dropdown>
          <el-dropdown-menu class="admin-dropdown-menu">
            <el-dropdown-item 
              class="dropdown-item"
              @click="$router.push('/admin/dashboard')"
            >
              <el-icon><Monitor /></el-icon>
              <span>控制面板</span>
            </el-dropdown-item>
            <el-dropdown-item 
              class="dropdown-item"
              @click="$router.push('/admin/settings')"
            >
              <el-icon><Setting /></el-icon>
              <span>系统设置</span>
            </el-dropdown-item>
            <el-dropdown-divider />
            <el-dropdown-item 
              class="dropdown-item logout-item"
              @click="logout"
            >
              <el-icon><SwitchButton /></el-icon>
              <span>退出登录</span>
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </el-menu>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAdminStore } from '@/stores/admin' // 引入 pinia 状态

import {
  House,
  User,
  Trophy,
  Medal,
  Connection,
  Lock,
  Monitor,
  Setting,
  SwitchButton,
  ArrowDown
} from '@element-plus/icons-vue'

const navbarStyle = {
  backgroundColor: '#0f172a', // 深蓝色背景
  color: '#f8fafc',           // 字体颜色
}

const router = useRouter()
const adminStore = useAdminStore()

const isLoggedIn = computed(() => adminStore.isLoggedIn) // 响应式绑定

const logout = () => {
  adminStore.logout()
  router.push('/')
}
</script>



<style>
/* 基础导航栏样式 */
.acm-navbar {
  height: 70px;
  padding: 0 2.5rem;
  border: none !important;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.15);
  display: flex;
  align-items: center;
  justify-content: space-between;
  transition: all 0.3s ease;
}

/* Logo品牌区域 */
.navbar-brand {
  display: flex;
  align-items: center;
  cursor: pointer;
  transition: all 0.3s ease;
  padding: 8px 12px;
  border-radius: 8px;
}

.navbar-brand:hover {
  background: rgba(56, 189, 248, 0.1);
}

.brand-logo {
  height: 42px;
  width: 42px;
  margin-right: 12px;
  object-fit: contain;
  filter: drop-shadow(0 0 6px rgba(56, 189, 248, 0.3));
}

.brand-text {
  line-height: 1.2;
}

.school-name {
  font-size: 13px;
  color: #94a3b8;
  font-weight: 500;
}

.team-name {
  font-size: 17px;
  font-weight: 700;
  color: #f8fafc;
  letter-spacing: 0.5px;
  text-shadow: 0 0 8px rgba(56, 189, 248, 0.2);
}

/* 导航菜单项 */
.nav-items {
  display: flex;
  gap: 6px;
}

.nav-item {
  font-size: 15px;
  font-weight: 500;
  margin: 0 4px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border-radius: 8px;
  padding: 0 16px !important;
}

.nav-item:hover {
  background: rgba(56, 189, 248, 0.1) !important;
  transform: translateY(-2px);
}

.nav-item .el-icon {
  margin-right: 8px;
  font-size: 18px;
  color: #38bdf8;
}

.nav-item.is-active {
  background: rgba(56, 189, 248, 0.15) !important;
  color: #38bdf8 !important;
}

/* 管理员区域 */
.admin-area {
  display: flex;
  align-items: center;
  gap: 16px;
}

.admin-login-btn {
  background: rgba(56, 189, 248, 0.15) !important;
  border: 1px solid rgba(56, 189, 248, 0.3) !important;
  color: #38bdf8 !important;
  font-weight: 500;
  padding: 8px 16px;
  border-radius: 8px;
  transition: all 0.3s ease;
}

.admin-login-btn:hover {
  background: rgba(56, 189, 248, 0.25) !important;
  transform: translateY(-1px);
  box-shadow: 0 2px 12px rgba(56, 189, 248, 0.15);
}

.admin-profile {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 6px 12px;
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.admin-profile:hover {
  background: rgba(56, 189, 248, 0.1);
}

.admin-avatar {
  border: 2px solid rgba(56, 189, 248, 0.3);
  box-shadow: 0 0 10px rgba(56, 189, 248, 0.2);
}

.admin-info {
  display: flex;
  flex-direction: column;
}

.admin-name {
  font-size: 14px;
  font-weight: 600;
  color: #f8fafc;
}

.admin-role {
  font-size: 12px;
  color: #94a3b8;
  margin-top: 2px;
}

.dropdown-icon {
  color: #94a3b8;
  font-size: 14px;
  margin-left: 4px;
  transition: transform 0.3s;
}

.el-dropdown:hover .dropdown-icon {
  transform: rotate(180deg);
}

/* 下拉菜单样式 */
.admin-dropdown-menu {
  border: none !important;
  background: #1e293b !important;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.25);
  border-radius: 8px !important;
  padding: 8px !important;
  min-width: 180px;
}

.dropdown-item {
  color: #e2e8f0 !important;
  font-size: 14px !important;
  padding: 8px 12px !important;
  border-radius: 6px !important;
  margin: 4px 0 !important;
  transition: all 0.2s ease !important;
}

.dropdown-item:hover {
  background: rgba(56, 189, 248, 0.1) !important;
  color: #38bdf8 !important;
}

.dropdown-item .el-icon {
  margin-right: 10px;
  font-size: 16px;
}

.logout-item {
  color: #f87171 !important;
}

.logout-item:hover {
  background: rgba(248, 113, 113, 0.1) !important;
}

.el-dropdown-menu__item:not(.is-disabled):focus {
  background: transparent !important;
}
</style>