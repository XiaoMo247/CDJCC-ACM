import { createRouter, createWebHistory } from 'vue-router'

// 从 pages 文件夹中导入组件
import HomePage from '../pages/HomePage.vue'
import KnowledgePage from '../pages/KnowledgePage.vue'
import CompetitionPage from '../pages/CompetitionPage.vue'
import RankingPage from '../pages/RankingPage.vue'
import FaqPage from '../pages/FaqPage.vue'
import LoginPage from '../pages/LoginPage.vue'

const routes = [
  { path: '/', component: HomePage },
  { path: '/knowledge', component: KnowledgePage },
  { path: '/competition', component: CompetitionPage },
  { path: '/ranking', component: RankingPage },
  { path: '/faq', component: FaqPage },
  { path: '/login', component: LoginPage },
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router