import axios from 'axios'

const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'https://localhost:8080',
  timeout: 5000
})
console.log('Base URL:', request.defaults.baseURL);
// 请求拦截器：自动加 token（admin_token 和 user_token）
request.interceptors.request.use(config => {
  const adminToken = localStorage.getItem('admin_token')
  const userToken = localStorage.getItem('user_token')

  if (adminToken) {
    config.headers.Authorization = `Bearer ${adminToken}`
  }

  if (userToken) {
    config.headers.Authorization = `Bearer ${userToken}`
  }

  return config
}, err => Promise.reject(err))

export default request
