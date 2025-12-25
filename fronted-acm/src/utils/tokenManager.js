/**
 * Token 管理工具 - 统一版本
 * 使用单一 token 存储策略，token 中包含用户角色信息
 *
 * 建议后端实现：
 * 1. 统一的 token 签发接口（返回包含用户类型的 JWT）
 * 2. Token payload 包含: { user_id, username, role: 'admin'|'student'|'member', exp }
 * 3. 统一的 token 验证中间件
 */

import request from './request'

// ==================== 常量定义 ====================

// 统一的 token 存储键
const TOKEN_KEY = 'auth_token'
// 用户信息存储键
const USER_INFO_KEY = 'user_info'
// Token 刷新间隔（50分钟）
const REFRESH_INTERVAL = 50 * 60 * 1000

// 定时器ID
let refreshTimer = null

// ==================== Token 存储相关 ====================

/**
 * 保存 token 和用户信息
 * @param {string} token - JWT token
 * @param {object} userInfo - 用户信息 { id, username, role, ... }
 */
export function saveToken(token, userInfo) {
  localStorage.setItem(TOKEN_KEY, token)
  localStorage.setItem(USER_INFO_KEY, JSON.stringify(userInfo))
}

/**
 * 获取 token
 * @returns {string|null}
 */
export function getToken() {
  return localStorage.getItem(TOKEN_KEY)
}

/**
 * 获取用户信息
 * @returns {object|null}
 */
export function getUserInfo() {
  const userInfoStr = localStorage.getItem(USER_INFO_KEY)
  if (!userInfoStr) return null

  try {
    return JSON.parse(userInfoStr)
  } catch (error) {
    console.error('解析用户信息失败:', error)
    return null
  }
}

/**
 * 清除所有认证信息
 */
export function clearAuth() {
  localStorage.removeItem(TOKEN_KEY)
  localStorage.removeItem(USER_INFO_KEY)

  // 兼容旧版本：清除旧的多 token 存储
  localStorage.removeItem('admin_token')
  localStorage.removeItem('user_token')
  localStorage.removeItem('member_token')
  localStorage.removeItem('member_info')

  stopTokenRefresh()
}

/**
 * 检查是否已登录
 * @returns {boolean}
 */
export function isLoggedIn() {
  return !!getToken()
}

/**
 * 获取当前用户角色
 * @returns {'admin'|'student'|'member'|null}
 */
export function getCurrentRole() {
  const userInfo = getUserInfo()
  return userInfo?.role || null
}

/**
 * 检查用户是否有指定角色
 * @param {string|string[]} roles - 角色或角色数组
 * @returns {boolean}
 */
export function hasRole(roles) {
  const currentRole = getCurrentRole()
  if (!currentRole) return false

  if (Array.isArray(roles)) {
    return roles.includes(currentRole)
  }
  return currentRole === roles
}

// ==================== Token 刷新相关 ====================

/**
 * 启动 token 自动刷新
 * 在用户登录后调用此方法
 */
export function startTokenRefresh() {
  // 清除之前的定时器
  stopTokenRefresh()

  // 立即检查一次 token
  checkAndRefreshToken()

  // 设置定时刷新
  refreshTimer = setInterval(() => {
    checkAndRefreshToken()
  }, REFRESH_INTERVAL)
}

/**
 * 停止 token 自动刷新
 * 在用户登出时调用此方法
 */
export function stopTokenRefresh() {
  if (refreshTimer) {
    clearInterval(refreshTimer)
    refreshTimer = null
  }
}

/**
 * 检查并刷新 token
 */
async function checkAndRefreshToken() {
  const token = getToken()

  // 如果没有 token，停止刷新
  if (!token) {
    stopTokenRefresh()
    return
  }

  try {
    // 方式1: 如果后端提供了刷新接口
    // const res = await request.post('/auth/refresh-token')
    // const { token: newToken, user } = res.data
    // saveToken(newToken, user)
    // console.log('Token 已刷新')

    // 方式2: 使用静默验证（当前采用）
    const role = getCurrentRole()
    let endpoint = '/user/me' // 默认端点

    if (role === 'admin') {
      endpoint = '/admin/me'
    }

    await request.get(endpoint)
    console.log(`Token 验证成功 [${role}]`)
  } catch (error) {
    console.error('Token 刷新失败:', error)
    // 刷新失败时停止自动刷新（不清除token，由401拦截器处理）
    stopTokenRefresh()
  }
}

// ==================== 兼容性支持 ====================

/**
 * 从旧版本迁移 token
 * 如果发现旧的多 token 存储，尝试迁移到新格式
 */
export function migrateOldTokens() {
  // 如果已经有新格式的 token，不需要迁移
  if (getToken()) return

  // 检查旧版本的 token
  const adminToken = localStorage.getItem('admin_token')
  const userToken = localStorage.getItem('user_token')
  const memberToken = localStorage.getItem('member_token')

  if (adminToken) {
    // 迁移管理员 token（需要调用后端获取用户信息）
    console.warn('检测到旧版 admin_token，建议重新登录以更新到新格式')
    // 临时兼容：直接使用旧 token
    localStorage.setItem(TOKEN_KEY, adminToken)
    localStorage.setItem(USER_INFO_KEY, JSON.stringify({ role: 'admin' }))
  } else if (userToken) {
    console.warn('检测到旧版 user_token，建议重新登录以更新到新格式')
    localStorage.setItem(TOKEN_KEY, userToken)
    localStorage.setItem(USER_INFO_KEY, JSON.stringify({ role: 'student' }))
  } else if (memberToken) {
    console.warn('检测到旧版 member_token，建议重新登录以更新到新格式')
    const memberInfo = localStorage.getItem('member_info')
    localStorage.setItem(TOKEN_KEY, memberToken)
    localStorage.setItem(USER_INFO_KEY, memberInfo || JSON.stringify({ role: 'member' }))
  }
}
