import axios from 'axios'
import router from '../router'
import { ElMessage } from 'element-plus'
import { getToken, clearAuth, stopTokenRefresh } from './tokenManager'

const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'https://localhost:8080',
  timeout: 5000
})

// 防抖标志：避免多个401请求同时触发跳转
let isRedirecting = false

// 请求拦截器：自动添加统一的 token
request.interceptors.request.use(config => {
  const token = getToken()

  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }

  return config
}, err => Promise.reject(err))

// 响应拦截器：处理 token 失效
request.interceptors.response.use(
  response => {
    // 请求成功，直接返回数据
    return response
  },
  error => {
    // 处理错误响应
    if (error.response) {
      const { status, data } = error.response

      // token 失效或未授权
      if (status === 401 || status === 403) {
        // 防抖：避免多次触发
        if (!isRedirecting) {
          isRedirecting = true

          // 停止 token 自动刷新并清除所有认证信息
          clearAuth()

          // 保存当前路由，用于登录后返回
          const currentPath = router.currentRoute.value.fullPath
          if (currentPath !== '/login') {
            sessionStorage.setItem('redirect_after_login', currentPath)
          }

          // 提示用户
          ElMessage.error({
            message: data?.message || '登录已过期，请重新登录',
            duration: 3000,
            onClose: () => {
              isRedirecting = false
            }
          })

          // 重定向到登录页
          router.push('/login').then(() => {
            // 延迟重置标志，确保跳转完成
            setTimeout(() => {
              isRedirecting = false
            }, 1000)
          })
        }
      } else if (status === 500) {
        ElMessage.error({
          message: data?.message || '服务器错误，请稍后重试',
          duration: 3000
        })
      } else if (status === 404) {
        ElMessage.error({
          message: '请求的资源不存在',
          duration: 3000
        })
      } else if (status === 400) {
        ElMessage.warning({
          message: data?.message || '请求参数错误',
          duration: 3000
        })
      }
    } else if (error.request) {
      // 请求已发出但没有收到响应
      ElMessage.error({
        message: '网络连接失败，请检查网络',
        duration: 3000
      })
    } else {
      // 其他错误
      ElMessage.error({
        message: '请求失败: ' + error.message,
        duration: 3000
      })
    }

    return Promise.reject(error)
  }
)

export default request
