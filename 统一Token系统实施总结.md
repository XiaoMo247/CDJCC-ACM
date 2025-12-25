# 统一 Token 系统实施总结

## 概述

已完成前后端统一 Token 认证系统的重构，将原先的多 Token 方案（admin_token, user_token, member_token）统一为单一 Token + role 字段的方案。

---

## 后端改动

### 1. 新增文件

#### `acm-site/utils/jwt/unified_token.go`
- 统一的 JWT Token 生成和解析
- `UnifiedClaims` 结构体包含 `UserID`, `Username`, `Role` 三个核心字段
- `GenerateUnifiedToken(userID, username, role)` - 生成 24 小时有效期的 token
- `ParseUnifiedToken(tokenString)` - 解析并验证 token

#### `acm-site/middleware/unified_auth.go`
- 统一认证中间件系统
- `AuthMiddleware()` - 验证 token 并将用户信息存入上下文（user_id, username, role）
- `RequireRole(roles...)` - 通用角色权限验证
- `RequireAdmin()` - 管理员权限（快捷方法）
- `RequireStudent()` - 学生权限（快捷方法）
- `RequireMember()` - 队员权限（快捷方法）

### 2. 修改文件

#### `acm-site/api/admin.go`
**变更点：**
- `AdminLogin()` 使用 `jwt.GenerateUnifiedToken(admin.ID, admin.Username, "admin")`
- 返回格式统一为：
  ```json
  {
    "token": "...",
    "user": {
      "id": 1,
      "username": "admin",
      "role": "admin"
    }
  }
  ```

#### `acm-site/api/user.go`
**变更点：**
- `UserLoginHandler()` 使用 `jwt.GenerateUnifiedToken(user.ID, user.StudentNumber, "student")`
- 返回格式统一为：
  ```json
  {
    "token": "...",
    "user": {
      "id": 1,
      "student_number": "2021001",
      "username": "张三",
      "role": "student"
    }
  }
  ```
- `ChangeUsername()` 和 `ChangePassword()` 使用统一的 `c.Get("userID")` 获取用户 ID

#### `acm-site/api/student.go`
**变更点：**
- `StudentLogin()` 使用 `jwt.GenerateUnifiedToken(student.ID, student.StudentID, "member")`
- 返回格式统一为：
  ```json
  {
    "code": 200,
    "msg": "登录成功",
    "token": "...",
    "user": {
      "id": 1,
      "username": "李四",
      "student_id": "ACM001",
      "role": "member"
    }
  }
  ```
- `GetStudentInfo()`, `UpdateStudentPassword()`, `UpdateStudentUsername()`, `UpdateStudentInfo()` 全部改用 `c.Get("user_id")` 从统一中间件获取用户 ID

#### `acm-site/router/router.go`
**重构要点：**
1. **公开接口区域**（无需认证）
   - 所有登录接口：`/api/admin/login`, `/api/user/login`, `/api/student/login`
   - 公共查询接口：比赛列表、FAQ、轮播图、荣誉墙、团队成员、公告列表
   - 课件资源下载

2. **认证接口区域**（需要登录）
   - 使用 `middleware.AuthMiddleware()` 作为基础认证
   - 子路由按角色分组：
     - `/api/admin/*` - 管理员专属（`middleware.RequireAdmin()`）
     - `/api/user/*` - 学生专属（`middleware.RequireStudent()`）
     - `/api/student/*` - 队员专属（`middleware.RequireMember()`）

---

## 前端改动

### 修改文件

#### `fronted-acm/src/utils/tokenManager.js`
**已完成的优化：**
- 统一存储：`auth_token` + `user_info`（包含 role 字段）
- `saveToken(token, userInfo)` - 保存 token 和用户信息
- `getToken()` - 获取 token
- `clearAuth()` - 清除认证信息
- `hasRole(roles)` - 基于 role 的权限判断
- `migrateOldTokens()` - 兼容旧版多 token 迁移
- `startTokenRefresh()` / `stopTokenRefresh()` - 自动刷新机制

#### `fronted-acm/src/utils/request.js`
**已完成的优化：**
- 请求拦截器：自动添加 `Authorization: Bearer ${token}` header
- 响应拦截器：捕获 401/403 错误，清理 token 并重定向到登录页
- 防抖机制：避免多个并发请求重复跳转

#### `fronted-acm/src/router/index.js`
**已完成的优化：**
- 使用 `isLoggedIn()` 和 `hasRole()` 进行路由守卫
- 保存原始目标路由到 `sessionStorage`，登录后自动跳转回原页面

#### 登录组件
- `AdminLogin.vue`
- `StudentLogin.vue`
- `MemberLogin.vue`

**已完成的优化：**
- 使用 `saveToken(token, userInfo)` 保存登录信息
- 调用 `startTokenRefresh()` 启动自动刷新
- 登录成功后根据 `sessionStorage` 中的 `redirect_after_login` 跳转

---

## 兼容性处理

### 后端兼容
在 `middleware/unified_auth.go` 中，为了兼容旧代码，AuthMiddleware 会同时设置新旧两套键名：
```go
c.Set("user_id", claims.UserID)     // 新键名
c.Set("username", claims.Username)  // 新键名
c.Set("role", claims.Role)          // 新键名

// 兼容旧键名
if claims.Role == "admin" {
    c.Set("adminID", claims.UserID)
    c.Set("adminUsername", claims.Username)
} else if claims.Role == "student" {
    c.Set("userID", claims.UserID)
}
```

### 前端兼容
`migrateOldTokens()` 函数会在应用启动时自动迁移旧版本的多 token 到新的统一格式。

---

## 测试建议

### 功能测试清单

1. **登录功能测试**
   - [ ] 管理员登录 - 验证返回格式包含 token 和 user.role="admin"
   - [ ] 学生登录 - 验证返回格式包含 token 和 user.role="student"
   - [ ] 队员登录 - 验证返回格式包含 token 和 user.role="member"

2. **权限验证测试**
   - [ ] 管理员访问管理后台 - 应成功
   - [ ] 学生访问管理后台 - 应被拦截（403）
   - [ ] 未登录用户访问需认证的接口 - 应被拦截（401）

3. **Token 过期测试**
   - [ ] Token 过期后访问接口 - 后端返回 401
   - [ ] 前端收到 401 - 自动清理 localStorage 并跳转登录页
   - [ ] 登录成功后 - 自动跳转回原页面

4. **跨角色功能测试**
   - [ ] 管理员修改密码
   - [ ] 学生修改用户名和密码
   - [ ] 队员更新个人信息（CF/AT/NC 账号）

---

## 迁移步骤（可选）

### 旧中间件清理

如果测试通过，可以考虑删除以下旧文件：
- `acm-site/middleware/admin_auth.go` → 已被 `unified_auth.go` 取代
- `acm-site/middleware/user_auth.go` → 已被 `unified_auth.go` 取代
- `acm-site/middleware/student_auth.go` → 已被 `unified_auth.go` 取代
- `acm-site/utils/jwt/token.go` → 已被 `unified_token.go` 取代
- `acm-site/utils/jwt/user_token.go` → 已被 `unified_token.go` 取代
- `acm-site/utils/jwt/student_jwt.go` → 已被 `unified_token.go` 取代

**注意：** 删除前请确保所有功能测试通过，并做好备份。

---

## 数据库影响

本次重构**不涉及**数据库结构变更，仅修改：
- JWT Token 的生成方式
- 前端 Token 存储方式
- 后端中间件验证逻辑

现有用户数据**无需迁移**。

---

## 安全改进

1. **统一 Token 密钥管理**
   - 当前使用 `"your_secret_key"` 作为密钥
   - **建议：** 改为从环境变量读取，避免硬编码

2. **Token 过期时间**
   - 当前设置为 24 小时
   - 前端实现了自动刷新机制

3. **HTTPS 传输**
   - Token 通过 `Authorization: Bearer` header 传输
   - **生产环境建议：** 启用 HTTPS 防止 token 被截获

---

## 总结

✅ **已完成：**
- 后端统一 Token 生成和验证系统
- 前端统一 Token 存储和管理
- 基于 role 的权限控制
- 旧系统兼容层
- Token 自动刷新机制
- 登录后原路返回功能

🔄 **建议后续优化：**
- 将 JWT 密钥改为环境变量
- 添加 Token 刷新接口（可选）
- 完善单元测试
- 清理旧中间件代码（测试通过后）

---

生成时间：2025-12-25
