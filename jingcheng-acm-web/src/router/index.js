import { createRouter, createWebHistory } from 'vue-router'

// 引入页面组件
const Home = () => import('@/views/Home.vue')
const About = () => import('@/views/About.vue')
const Contests = () => import('@/views/Contests.vue')
const Rank = () => import('@/views/Rank.vue')
const Join = () => import('@/views/Join.vue')
const AdminLogin = () => import('@/views/AdminLogin.vue')  // 管理员登录页
const AdminDashboard = () => import('@/views/AdminDashboard.vue')  // 管理员面板

// 路由配置
const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home
  },
  {
    path: '/about',
    name: 'About',
    component: About
  },
  {
    path: '/contests',
    name: 'Contests',
    component: Contests
  },
  {
    path: '/rank',
    name: 'Rank',
    component: Rank
  },
  {
    path: '/join',
    name: 'Join',
    component: Join
  },
  {
    path: '/admin/login',
    name: 'AdminLogin',
    component: AdminLogin  // 管理员登录页
  },
  {
    path: '/admin/dashboard',
    name: 'AdminDashboard',
    component: AdminDashboard,  // 管理员面板
    meta: { requiresAuth: true } // 需要认证
  },
  {
    path: '/:pathMatch(.*)*',  // 404页面
    name: 'NotFound',
    component: () => import('@/views/NotFound.vue')  // 创建 NotFound.vue 页面
  }
]

// 创建 router 实例
const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫：确保只有登录用户能访问管理员面板
router.beforeEach((to, from, next) => {
  if (to.meta.requiresAuth && !localStorage.getItem('adminToken')) {
    next('/admin/login')  // 未登录则跳转到登录页面
  } else {
    next()  // 允许访问
  }
})

export default router
