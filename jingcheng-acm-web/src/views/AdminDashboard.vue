<template>
  <div class="admin-dashboard">
    <h2 class="title">管理员面板</h2>
    <p class="description">欢迎来到管理员后台管理系统！</p>

    <!-- 管理员管理和公告管理按钮 -->
    <div class="button-group">
      <el-button @click="showAdminManagement = true" type="primary" class="nav-button" :class="{ active: showAdminManagement }">管理员管理</el-button>
      <el-button @click="showAdminManagement = false" type="primary" class="nav-button" :class="{ active: !showAdminManagement }">公告管理</el-button>
    </div>

    <!-- 根据按钮选择渲染组件 -->
    <div v-if="showAdminManagement">
      <AdminManagement />
    </div>
    <div v-else>
      <AnnouncementManagement />
    </div>

    <el-button @click="logout" type="danger" class="logout-button">退出登录</el-button>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import AdminManagement from '@/components/AdminManagement.vue'  // 引入管理员管理组件
import AnnouncementManagement from '@/components/AnnouncementManagement.vue'  // 引入公告管理组件

// 控制显示的组件
const showAdminManagement = ref(true)  // 默认为显示管理员管理

// 获取路由对象
const router = useRouter()

// 退出登录
const logout = () => {
  localStorage.removeItem('adminToken')  // 移除存储的 Token
  router.push('/admin/login')  // 跳转到管理员登录页面
}
</script>

<style scoped>
.admin-dashboard {
  padding: 40px 20px;
  background-color: #f7f9fc;
  border-radius: 10px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  text-align: center;
  animation: fadeInUp 0.6s ease-in-out forwards;
}

.title {
  font-size: 2rem;
  font-weight: bold;
  color: #3498db;
  margin-bottom: 10px;
}

.description {
  font-size: 1.1rem;
  color: #7f8c8d;
  margin-bottom: 30px;
}

.button-group {
  margin-bottom: 30px;
}

.nav-button {
  width: 180px;
  margin-right: 15px;
  transition: all 0.3s ease;
}

.nav-button.active {
  background-color: #3498db;
  color: #fff;
}

.logout-button {
  width: 200px;
  margin-top: 30px;
  transition: transform 0.3s ease;
}

.logout-button:hover {
  transform: scale(1.05);
  box-shadow: 0px 4px 15px rgba(52, 152, 219, 0.3);
}

/* 动画效果 */
@keyframes fadeInUp {
  0% {
    opacity: 0;
    transform: translateY(20px);
  }
  100% {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 移动端适配 */
@media (max-width: 500px) {
  .title {
    font-size: 1.5rem;
  }

  .description {
    font-size: 1rem;
  }

  .nav-button {
    width: 100%;
    margin-right: 0;
    margin-bottom: 10px;
  }

  .logout-button {
    width: 100%;
  }
}
</style>
