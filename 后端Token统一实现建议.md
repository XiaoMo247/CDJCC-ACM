# 后端 Token 统一实现建议

## 🎯 目标

统一所有用户类型（admin、student、member）的 token 签发和验证机制，使用单一的 JWT token 包含用户角色信息。

---

## 📋 当前问题

目前的实现：
- ❌ 三个独立的登录接口（`/admin/login`、`/user/login`、`/student/login`）
- ❌ 三种不同的 token 存储（`admin_token`、`user_token`、`member_token`）
- ❌ 需要维护多套认证逻辑

---

## ✅ 推荐方案

### 1. 统一的 JWT Token 结构

```go
type TokenClaims struct {
    UserID   uint   `json:"user_id"`
    Username string `json:"username"`
    Role     string `json:"role"` // "admin" | "student" | "member"
    jwt.RegisteredClaims
}
```

### 2. 统一的 Token 签发

```go
// utils/jwt.go
package utils

import (
    "time"
    "github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte("your-secret-key") // 从配置文件读取

// GenerateToken 生成统一的 JWT token
func GenerateToken(userID uint, username string, role string) (string, error) {
    claims := TokenClaims{
        UserID:   userID,
        Username: username,
        Role:     role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 24小时过期
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "acm-site",
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}

// ParseToken 解析并验证 token
func ParseToken(tokenString string) (*TokenClaims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
        return jwtSecret, nil
    })

    if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
        return claims, nil
    }
    return nil, err
}
```

### 3. 统一的认证中间件

```go
// middleware/auth.go
package middleware

import (
    "net/http"
    "strings"
    "acm-site/utils"
    "github.com/gin-gonic/gin"
)

// AuthMiddleware 统一的认证中间件
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 获取 Authorization header
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"message": "未提供认证令牌"})
            c.Abort()
            return
        }

        // 解析 Bearer token
        parts := strings.SplitN(authHeader, " ", 2)
        if !(len(parts) == 2 && parts[0] == "Bearer") {
            c.JSON(http.StatusUnauthorized, gin.H{"message": "认证令牌格式错误"})
            c.Abort()
            return
        }

        // 验证 token
        claims, err := utils.ParseToken(parts[1])
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"message": "认证令牌无效或已过期"})
            c.Abort()
            return
        }

        // 将用户信息存入上下文
        c.Set("user_id", claims.UserID)
        c.Set("username", claims.Username)
        c.Set("role", claims.Role)

        c.Next()
    }
}

// RequireRole 角色权限中间件
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        role, exists := c.Get("role")
        if !exists {
            c.JSON(http.StatusForbidden, gin.H{"message": "无权限访问"})
            c.Abort()
            return
        }

        userRole := role.(string)
        for _, allowedRole := range allowedRoles {
            if userRole == allowedRole {
                c.Next()
                return
            }
        }

        c.JSON(http.StatusForbidden, gin.H{"message": "角色权限不足"})
        c.Abort()
    }
}
```

### 4. 修改登录接口返回格式

```go
// api/admin.go
package api

import (
    "net/http"
    "acm-site/model"
    "acm-site/utils"
    "github.com/gin-gonic/gin"
    "golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
    Token string      `json:"token"`
    User  interface{} `json:"user"`
}

// AdminLogin 管理员登录
func AdminLogin(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "请求参数错误"})
        return
    }

    // 查询管理员
    var admin model.Admin
    if err := database.DB.Where("username = ?", req.Username).First(&admin).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"message": "用户名或密码错误"})
        return
    }

    // 验证密码
    if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(req.Password)); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"message": "用户名或密码错误"})
        return
    }

    // 生成统一的 token，包含角色信息
    token, err := utils.GenerateToken(admin.ID, admin.Username, "admin")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "生成令牌失败"})
        return
    }

    // 返回 token 和用户信息
    c.JSON(http.StatusOK, LoginResponse{
        Token: token,
        User: gin.H{
            "id":       admin.ID,
            "username": admin.Username,
            "role":     "admin",
        },
    })
}
```

### 5. 用户登录接口（类似结构）

```go
// api/user.go
func UserLogin(c *gin.Context) {
    // ... 验证逻辑 ...

    // 生成 token，role 设为 "student"
    token, err := utils.GenerateToken(user.ID, user.StudentNumber, "student")

    c.JSON(http.StatusOK, LoginResponse{
        Token: token,
        User: gin.H{
            "id":             user.ID,
            "student_number": user.StudentNumber,
            "role":           "student",
        },
    })
}
```

### 6. 路由配置示例

```go
// router/router.go
package router

import (
    "acm-site/api"
    "acm-site/middleware"
    "github.com/gin-gonic/gin"
)

func RouterInit(r *gin.Engine) {
    // 公开接口
    r.POST("/admin/login", api.AdminLogin)
    r.POST("/user/login", api.UserLogin)
    r.POST("/student/login", api.StudentLogin)

    // 需要认证的接口
    auth := r.Group("/", middleware.AuthMiddleware())
    {
        // 任意登录用户可访问
        auth.GET("/user/me", api.GetCurrentUser)

        // 仅管理员可访问
        admin := auth.Group("/admin", middleware.RequireRole("admin"))
        {
            admin.GET("/dashboard", api.AdminDashboard)
            admin.POST("/announcement/create", api.CreateAnnouncement)
        }

        // 仅学生可访问
        student := auth.Group("/student", middleware.RequireRole("student"))
        {
            student.GET("/dashboard", api.StudentDashboard)
        }

        // 管理员和队员都可访问
        member := auth.Group("/member", middleware.RequireRole("admin", "member"))
        {
            member.GET("/dashboard", api.MemberDashboard)
        }
    }
}
```

---

## 🔄 迁移步骤

### 第一阶段：准备工作
1. 创建 `utils/jwt.go` 实现统一的 token 生成和解析
2. 创建 `middleware/auth.go` 实现统一的认证中间件
3. 在 `model` 中确保所有用户模型都有明确的角色标识

### 第二阶段：修改登录接口
1. 修改 `/admin/login` 返回格式包含 `user` 对象
2. 修改 `/user/login` 返回格式包含 `user` 对象
3. 修改 `/student/login` 返回格式包含 `user` 对象
4. 确保所有登录接口都使用 `utils.GenerateToken`

### 第三阶段：替换认证中间件
1. 将所有路由的旧中间件替换为 `AuthMiddleware()`
2. 使用 `RequireRole()` 添加角色权限控制
3. 删除旧的多套认证中间件代码

### 第四阶段：测试
1. 测试所有登录接口返回格式正确
2. 测试 token 能否正确携带和解析
3. 测试角色权限控制是否生效
4. 测试 token 过期处理

---

## 📝 前后端对接说明

### 前端期望的响应格式

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "admin",
    "role": "admin"
  }
}
```

### Token Payload 示例

```json
{
  "user_id": 1,
  "username": "admin",
  "role": "admin",
  "exp": 1735123200,
  "iat": 1735036800,
  "iss": "acm-site"
}
```

---

## ⚡ 优势

1. **代码简化**: 只需要维护一套 token 逻辑
2. **安全性提升**: 统一的验证机制，减少安全漏洞
3. **扩展性好**: 新增用户类型只需添加角色标识
4. **前后端一致**: 统一的认证流程，减少沟通成本
5. **便于调试**: token 中包含所有必要信息，易于追踪

---

## 🚀 可选增强功能

### 1. Token 刷新机制

```go
// RefreshToken 刷新 token
func RefreshToken(c *gin.Context) {
    // 从上下文获取用户信息（已通过认证中间件）
    userID := c.GetUint("user_id")
    username := c.GetString("username")
    role := c.GetString("role")

    // 生成新 token
    newToken, err := utils.GenerateToken(userID, username, role)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "刷新令牌失败"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"token": newToken})
}
```

### 2. Token 黑名单（用于登出）

```go
// 使用 Redis 存储已登出的 token
func Logout(c *gin.Context) {
    token := c.GetHeader("Authorization")
    // 将 token 加入黑名单（存储到 Redis）
    // redis.Set(token, "blacklist", tokenExpTime)
    c.JSON(http.StatusOK, gin.H{"message": "登出成功"})
}
```

---

## 📞 需要前端配合的地方

1. ✅ 前端已经实现了统一的 token 存储（`auth_token`）
2. ✅ 前端已经实现了统一的请求拦截器
3. ✅ 前端已经实现了角色判断逻辑
4. ⏳ 等待后端修改登录接口返回格式

---

## 🎯 总结

通过这次重构，整个系统的认证流程将更加**统一、安全、易维护**。前端已经完成了相应的适配工作，后端只需要按照上述建议进行修改即可实现无缝对接。

如有疑问，请随时沟通！
