<template>
  <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0-beta3/css/all.min.css">
  <div class="admin-dashboard">
    <div class="sidebar">
      <h3 class="sidebar-title"><i class="fas fa-tools"></i> 控制面板</h3>
      <ul>
        <li :class="{ active: currentTab === 'announcement' }" @click="currentTab = 'announcement'">
          <i class="fas fa-bullhorn"></i> 公告管理
        </li>
        <li :class="{ active: currentTab === 'admin' }" @click="currentTab = 'admin'">
          <i class="fas fa-user-shield"></i> 管理员账户
        </li>
        <li :class="{ active: currentTab === 'join' }" @click="currentTab = 'join'">
          <i class="fas fa-clipboard-list"></i> 申请表管理
        </li>
        <li :class="{ active: currentTab === 'user' }" @click="currentTab = 'user'">
          <i class="fas fa-users"></i> 用户管理
        </li>
      </ul>
    </div>

    <div class="content">
      <AnnouncementManage v-if="currentTab === 'announcement'" />
      <AdminAccountManage v-if="currentTab === 'admin'" />
      <JoinFormManage v-if="currentTab === 'join'" />
      <UserManage v-if="currentTab === 'user'" />
    </div>
  </div>
</template>

<script>
import AnnouncementManage from '@/components/admin/AnnouncementManage.vue'
import AdminAccountManage from '@/components/admin/AdminAccountManage.vue'
import JoinFormManage from '@/components/admin/JoinFormManage.vue'
import UserManage from '@/components/admin/UserManage.vue'

export default {
  name: 'AdminDashboard',
  components: {
    AnnouncementManage,
    AdminAccountManage,
    JoinFormManage,
    UserManage
  },
  data() {
    return {
      currentTab: 'announcement'
    }
  }
}
</script>

<style scoped>
.admin-dashboard {
  display: flex;
  height: 100vh;
  background-color: #f8fafc;
}

/* 左侧导航栏 */
.sidebar {
  width: 240px;
  background: linear-gradient(180deg, #e2e8f0, #cbd5e1);
  padding: 20px;
  box-shadow: 2px 0 6px rgba(0, 0, 0, 0.05);
  display: flex;
  flex-direction: column;
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

/* 右侧内容区域 */
.content {
  flex: 1;
  padding: 24px;
  overflow-y: auto;
  background-color: #f8fafc;
}
</style>
