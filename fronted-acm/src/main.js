import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import { startTokenRefresh, isLoggedIn, migrateOldTokens } from './utils/tokenManager'

const app = createApp(App)

app.use(router)  // 使用 Vue Router
app.use(createPinia())  // 使用 Pinia 状态管理
app.use(ElementPlus)  // 使用 Element Plus UI 库

app.mount('#app')  // 挂载应用

// 应用启动时的初始化
// 1. 尝试迁移旧版本的 token（兼容性处理）
migrateOldTokens()

// 2. 如果用户已登录，启动 token 自动刷新
if (isLoggedIn()) {
  startTokenRefresh()
}
