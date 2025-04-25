package middleware

import (
	"acm-site/utils/jwt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func StudentAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "未提供合法 Token"})
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := jwt.ParseStudentToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "Token 无效或已过期"})
			c.Abort()
			return
		}

		c.Set("student_id", claims.StudentID)
		c.Set("student_number", claims.StudentNumber)

		c.Next()
	}
}
