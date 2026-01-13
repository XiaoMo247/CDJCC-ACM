import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import { startTokenRefresh, isLoggedIn, migrateOldTokens } from './utils/tokenManager'
import { ElMessage } from 'element-plus'

// CSS 导入顺序：变量 -> reset -> 第三方库 -> 全局/组件/工具类
import './styles/abstracts/variables.css'
import './styles/abstracts/page-themes.css'
import './styles/base/reset.css'
import 'element-plus/dist/index.css'
import './styles/base/global.css'
import './styles/base/typography.css'
import './styles/components/cards.css'
import './styles/components/forms.css'
import './styles/components/buttons.css'
import './styles/components/tables.css'
import './styles/utilities/spacing.css'
import './styles/utilities/display.css'
import './styles/abstracts/animations.css'

// 应用启动时的初始化
// 1. 尝试迁移旧版本的 token（兼容性处理）
migrateOldTokens()

const app = createApp(App)

// 全局错误处理：避免意外错误导致白屏
app.config.errorHandler = (err, instance, info) => {
  console.error('全局错误捕获:', err)
  console.error('错误信息:', info)

  ElMessage.error({
    message: '抱歉，应用遇到了一些问题，请刷新页面重试',
    duration: 5000,
    showClose: true
  })
}

// 全局警告处理（仅开发环境）
if (import.meta.env.DEV) {
  app.config.warnHandler = (msg, instance, trace) => {
    console.warn('Vue 警告:', msg)
    console.warn('组件追踪:', trace)
  }
}

// Promise 未捕获错误
window.addEventListener('unhandledrejection', (event) => {
  console.error('未处理的 Promise 错误:', event.reason)
  ElMessage.error({
    message: '操作失败，请稍后重试',
    duration: 5000,
    showClose: true
  })
})

app.use(router)  // 使用 Vue Router
app.use(createPinia())  // 使用 Pinia 状态管理
app.use(ElementPlus)  // 使用 Element Plus UI 库

app.mount('#app')  // 挂载应用

// 2. 如果用户已登录，启动 token 自动刷新
if (isLoggedIn()) {
  startTokenRefresh()
}
