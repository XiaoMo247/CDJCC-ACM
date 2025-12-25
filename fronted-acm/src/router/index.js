import { createRouter, createWebHistory } from 'vue-router'
import { isLoggedIn, hasRole } from '../utils/tokenManager'

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
import AnnouncementAndNews from '../pages/AnnouncementsAndNews.vue'
import HonorWall from '../pages/HonorWall.vue'

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
  { path: '/announcement-and-news', component: AnnouncementAndNews},
  { path: '/honor-wall', component: HonorWall},
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 全局导航守卫（访问控制）- 使用统一的 token 判断
router.beforeEach((to, from, next) => {
  const isAdminPage = to.path.startsWith('/admin')
  const isStudentPage = to.path.startsWith('/student')
  const isMemberPage = to.path.startsWith('/member')

  // 管理员页面：需要 admin 角色
  if (isAdminPage) {
    if (!hasRole('admin')) {
      sessionStorage.setItem('redirect_after_login', to.fullPath)
      next('/login')
      return
    }
  }

  // 学生页面：需要 student 角色
  if (isStudentPage) {
    if (!hasRole('student')) {
      sessionStorage.setItem('redirect_after_login', to.fullPath)
      next('/login')
      return
    }
  }

  // 成员页面：需要任意登录角色
  if (isMemberPage) {
    if (!isLoggedIn()) {
      sessionStorage.setItem('redirect_after_login', to.fullPath)
      next('/login')
      return
    }
  }

  // 如果已登录访问登录页，跳转到对应的 dashboard
  if (to.path === '/login' && isLoggedIn()) {
    if (hasRole('admin')) {
      next('/admin/dashboard')
    } else if (hasRole('student')) {
      next('/student/dashboard')
    } else if (hasRole('member')) {
      next('/member/dashboard')
    } else {
      next()
    }
    return
  }

  // 其他页面：放行
  next()
})

export default router