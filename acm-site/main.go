package main

import (
	"acm-site/config"
	"acm-site/database"
	"acm-site/middleware"
	"acm-site/router"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// 检查必需的环境变量
	checkRequiredEnvVars()

	config.InitConfig() // 加载配置
	database.InitDB()   // 初始化数据库连接

	r := gin.Default()

	r.Static("/uploads", "./uploads")

	r.Use(middleware.CorsMiddleware())

	router.RouterInit(r)

	r.Run(":8081")
}

// checkRequiredEnvVars 检查必需的环境变量
func checkRequiredEnvVars() {
	jwtKey := os.Getenv("JWT_SECRET_KEY")

	// 生产环境警告
	if gin.Mode() == gin.ReleaseMode {
		if jwtKey == "" || jwtKey == "your_super_secret_jwt_key_here_change_this_in_production" || jwtKey == "default_dev_secret_key_please_change_in_production" {
			fmt.Println("⚠️  [安全警告] 生产环境必须设置 JWT_SECRET_KEY 环境变量！")
			fmt.Println("   生成密钥: openssl rand -base64 32")
			fmt.Println("   然后在 .env 文件中设置: JWT_SECRET_KEY=<生成的密钥>")
			os.Exit(1)
		}

		allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
		if allowedOrigins == "" {
			fmt.Println("⚠️  [安全警告] 生产环境建议设置 ALLOWED_ORIGINS 环境变量")
			fmt.Println("   当前将允许所有源访问，请在 .env 文件中配置")
		}

		// 检查数据库配置
		dbHost := os.Getenv("DB_HOST")
		dbPassword := os.Getenv("DB_PASSWORD")
		if dbHost == "127.0.0.1" || dbPassword == "123456" {
			fmt.Println("⚠️  [安全警告] 检测到默认数据库配置")
			fmt.Println("   请在 .env 文件中更新数据库凭据")
		}
	} else {
		// 开发环境提示
		if jwtKey == "" {
			fmt.Println("ℹ️  [提示] 未设置 JWT_SECRET_KEY，使用默认开发密钥")
		}
		fmt.Println("ℹ️  [提示] 当前运行在开发模式，从 .env 文件读取配置")
	}
}
