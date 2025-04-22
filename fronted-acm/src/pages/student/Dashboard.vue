<template>
  <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0-beta3/css/all.min.css">
  <div class="student-dashboard">
    <div class="sidebar">
      <h3 class="sidebar-title">
        <i class="fas fa-user-graduate"></i> 学生中心 : {{ Username }}
      </h3>
      <ul>
        <li :class="{ active: currentTab === 'join' }" @click="currentTab = 'join'">
          <i class="fas fa-clipboard-list"></i> 申请表管理
        </li>
        <li :class="{ active: currentTab === 'account' }" @click="currentTab = 'account'">
          <i class="fas fa-user-cog"></i> 账户设置
        </li>
      </ul>
    </div>

    <div class="content">
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
      Username: '加载中...'
    }
  },
  mounted() {
    this.fetchUsername()
  },
  methods: {
    async fetchUsername() {
      try {
        const res = await request.get('/user/info') // 假设你后端提供这个接口返回当前登录用户
        this.Username = res.data.data.username || '未知用户'
      } catch (err) {
        ElMessage.error('获取用户名失败')
        this.Username = '未登录'
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
  }
  
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
  
  .content {
    flex: 1;
    padding: 24px;
    overflow-y: auto;
    background-color: #f8fafc;
  }
  </style>
  