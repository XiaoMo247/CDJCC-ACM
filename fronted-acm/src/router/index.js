import { createRouter, createWebHistory } from 'vue-router'

// 页面组件
import HomePage from '../pages/HomePage/HomePage.vue'
import KnowledgePage from '../pages/KnowledgePage.vue'
import CompetitionPage from '../pages/CompetitionPage.vue'
import RankingPage from '../pages/RankingPage.vue'
import FaqPage from '../pages/FaqPage.vue'
import LoginPage from '../pages/LoginPage.vue'
import AdminDashboard from '../pages/admin/Dashboard.vue'
import StudentDashboard from '../pages/student/Dashboard.vue'
import ForgotPasswordPage from '../pages/ForgotPassword.vue'
import MemberDashboard from '../pages/member/Dashboard.vue'

const routes = [
  { path: '/', component: HomePage },
  { path: '/knowledge', component: KnowledgePage },
  { path: '/competition', component: CompetitionPage },
  { path: '/ranking', component: RankingPage },
  { path: '/faq', component: FaqPage },
  { path: '/login', component: LoginPage },
  { path: '/admin/dashboard', component: AdminDashboard },
  { path: '/student/dashboard', component: StudentDashboard },
  { path: '/forgot-password', component: ForgotPasswordPage },
  { path: '/member/dashboard', component: MemberDashboard },
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 全局导航守卫（访问控制）
router.beforeEach((to, from, next) => {
  const isAdminPage = to.path.startsWith('/admin')
  const isUserPage = to.path.startsWith('/student')

  const adminToken = localStorage.getItem('admin_token')
  const userToken = localStorage.getItem('user_token')

  if (isAdminPage && !adminToken) {
    next('/login') // 管理员未登录
  } else if (isUserPage && !userToken) {
    next('/login') // 用户未登录
  } else {
    next() // 放行
  }
})

export default router