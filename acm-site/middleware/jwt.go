package middleware

import (
	"acm-site/utils/jwt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "未提供合法 Token"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := jwt.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "Token 无效或已过期"})
			c.Abort()
			return
		}

		// 将管理员信息保存到上下文中，后续接口可以通过 c.Get() 获取
		c.Set("adminID", claims.AdminID)
		c.Set("adminUsername", claims.Username)

		c.Next()
	}
}
