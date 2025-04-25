<template>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0-beta3/css/all.min.css">
    <div class="member-dashboard">
        <div class="sidebar">
            <h3 class="sidebar-title">
                <i class="fas fa-user-graduate"></i> 队员中心 : {{ Username }}
            </h3>
            <ul>
                <li class="active">
                    <i class="fas fa-user-cog"></i> 账户设置
                </li>
            </ul>
        </div>

        <div class="content">
            <MemberAccountSettings @updated="fetchUsername" />
        </div>
    </div>
</template>

<script>
import MemberAccountSettings from '@/components/member/MemberAccountSettings.vue'
import request from '@/utils/request'
import { ElMessage } from 'element-plus'

export default {
    name: 'MemberDashboard',
    components: {
        MemberAccountSettings
    },
    data() {
        return {
            Username: '加载中...'
        }
    },
    mounted() {
        this.fetchUsername()
    },
    methods: {
        async fetchUsername() {
            try {
                // 从 localStorage 或 Vuex 中获取保存的 Token
                const token = localStorage.getItem('member_token') // 假设 token 存储在 localStorage
                if (!token) {
                    ElMessage.error('未登录，请重新登录')
                    this.Username = '未登录'
                    return
                }
                // 请求学生信息，带上 Authorization Token
                const res = await request.get('/student/info', {
                    headers: {
                        Authorization: `Bearer ${token}` // 加上 Bearer token
                    }
                })
                this.Username = res.data.data.Username || '未知队员'
            } catch (err) {
                ElMessage.error('获取用户名失败')
                this.Username = '未登录'
            }
        }
    }
}
</script>

<style scoped>
.member-dashboard {
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
    cursor: default;
    color: #1e293b;
    font-size: 15px;
    display: flex;
    align-items: center;
    background-color: #38bdf8;
    color: white;
}

.sidebar li i {
    margin-right: 12px;
    width: 20px;
    height: 20px;
    display: inline-block;
    text-align: center;
}

.content {
    flex: 1;
    padding: 24px;
    overflow-y: auto;
    background-color: #f8fafc;
}
</style>
