import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'

const app = createApp(App)

app.use(router)  // 使用 Vue Router
app.use(createPinia())  // 使用 Pinia 状态管理
app.use(ElementPlus)  // 使用 Element Plus UI 库

app.mount('#app')  // 挂载应用
