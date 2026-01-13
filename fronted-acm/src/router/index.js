import { createRouter, createWebHistory } from 'vue-router'
import { isLoggedIn, hasRole } from '../utils/tokenManager'

// 页面组件（路由懒加载）
const HomePage = () => import('../pages/HomePage/HomePage.vue')
const KnowledgePage = () => import('../pages/KnowledgePage.vue')
const CompetitionPage = () => import('../pages/CompetitionPage.vue')
const RankingPage = () => import('../pages/RankingPage.vue')
const FaqPage = () => import('../pages/FaqPage.vue')
const LoginPage = () => import('../pages/LoginPage.vue')
const AdminDashboard = () => import('../pages/admin/Dashboard.vue')
const StudentDashboard = () => import('../pages/student/Dashboard.vue')
const ForgotPasswordPage = () => import('../pages/ForgotPassword.vue')
const MemberDashboard = () => import('../pages/member/Dashboard.vue')
const AnnouncementAndNews = () => import('../pages/AnnouncementsAndNews.vue')
const HonorWall = () => import('../pages/HonorWall.vue')
const NotFound = () => import('../pages/NotFound.vue')

const routes = [
  {
    path: '/',
    name: 'Home',
    component: HomePage,
    meta: { title: 'CDJC ACM 集训队', requiresAuth: false }
  },
  {
    path: '/announcement-and-news',
    name: 'AnnouncementAndNews',
    component: AnnouncementAndNews,
    meta: { title: '公告与新闻', requiresAuth: false }
  },
  {
    path: '/knowledge',
    name: 'Knowledge',
    component: KnowledgePage,
    meta: { title: '知识库', requiresAuth: false }
  },
  {
    path: '/competition',
    name: 'Competition',
    component: CompetitionPage,
    meta: { title: '比赛', requiresAuth: false }
  },
  {
    path: '/ranking',
    name: 'Ranking',
    component: RankingPage,
    meta: { title: '排行榜', requiresAuth: false }
  },
  {
    path: '/honor-wall',
    name: 'HonorWall',
    component: HonorWall,
    meta: { title: '荣誉墙', requiresAuth: false }
  },
  {
    path: '/faq',
    name: 'Faq',
    component: FaqPage,
    meta: { title: '常见问题', requiresAuth: false }
  },
  {
    path: '/login',
    name: 'Login',
    component: LoginPage,
    meta: { title: '登录', requiresAuth: false }
  },
  {
    path: '/forgot-password',
    name: 'ForgotPassword',
    component: ForgotPasswordPage,
    meta: { title: '找回密码', requiresAuth: false }
  },
  {
    path: '/admin/dashboard',
    name: 'AdminDashboard',
    component: AdminDashboard,
    meta: { title: '管理员控制面板', requiresAuth: true, roles: ['admin'] }
  },
  {
    path: '/student/dashboard',
    name: 'StudentDashboard',
    component: StudentDashboard,
    meta: { title: '学生控制面板', requiresAuth: true, roles: ['student'] }
  },
  {
    path: '/member/dashboard',
    name: 'MemberDashboard',
    component: MemberDashboard,
    meta: { title: '队员控制面板', requiresAuth: true, roles: ['member'] }
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: NotFound,
    meta: { title: '页面未找到', requiresAuth: false }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) return savedPosition
    if (to.hash) {
      return { el: to.hash, behavior: 'smooth' }
    }
    return { top: 0, behavior: 'smooth' }
  }
})

// 全局导航守卫（访问控制）- 使用统一的 token 判断
router.beforeEach((to, from, next) => {
  document.title = to.meta?.title || 'CDJC ACM 集训队'

  // 管理员页面：需要 admin 角色
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

  if (to.meta?.requiresAuth && !isLoggedIn()) {
    sessionStorage.setItem('redirect_after_login', to.fullPath)
    next('/login')
    return
  }

  if (to.meta?.roles && !hasRole(to.meta.roles)) {
    sessionStorage.setItem('redirect_after_login', to.fullPath)
    next('/login')
    return
  }

  // 其他页面：放行
  next()
})

export default router
