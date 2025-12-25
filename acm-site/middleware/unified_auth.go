package middleware

import (
	"acm-site/utils/jwt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware 统一的认证中间件
// 验证 token 并将用户信息存入上下文
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

		tokenString := parts[1]

		// 验证 token
		claims, err := jwt.ParseUnifiedToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "认证令牌无效或已过期"})
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		// 兼容旧代码：也保存旧的键名
		if claims.Role == "admin" {
			c.Set("adminID", claims.UserID)
			c.Set("adminUsername", claims.Username)
		} else if claims.Role == "student" {
			c.Set("userID", claims.UserID)
		}

		c.Next()
	}
}

// RequireRole 角色权限中间件
// 限制只有指定角色才能访问
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

// RequireAdmin 要求管理员权限（快捷方法）
func RequireAdmin() gin.HandlerFunc {
	return RequireRole("admin")
}

// RequireStudent 要求学生权限（快捷方法）
func RequireStudent() gin.HandlerFunc {
	return RequireRole("student")
}

// RequireMember 要求队员权限（快捷方法）
func RequireMember() gin.HandlerFunc {
	return RequireRole("member")
}

// RequireAnyRole 要求任意登录角色（快捷方法）
func RequireAnyRole() gin.HandlerFunc {
	return RequireRole("admin", "student", "member")
}
