package middleware

import (
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CorsMiddleware() gin.HandlerFunc {
	// 从环境变量读取允许的来源（逗号分隔），方便生产环境配置：
	// ALLOWED_ORIGINS=http://xxx,http://yyy
	allowedOriginsEnv := strings.TrimSpace(os.Getenv("ALLOWED_ORIGINS"))
	allowedOrigins := []string{}

	if allowedOriginsEnv != "" {
		parts := strings.Split(allowedOriginsEnv, ",")
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p != "" {
				allowedOrigins = append(allowedOrigins, p)
			}
		}
	} else {
		// 默认允许本地 Vite 开发服务器（5173）
		allowedOrigins = []string{
			"http://localhost:5173",
			"http://127.0.0.1:5173",
		}
	}

	return cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
