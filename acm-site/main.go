package main

import (
	"acm-site/config"
	"acm-site/database"
	"acm-site/middleware"
	"acm-site/router"
	"bufio"
	"bytes"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	// 尝试加载 .env（优先 acm-site 目录，其次仓库根目录）。
	// 这样仓库只需要维护一份根目录 .env。
	loadDotEnv([]string{".env", filepath.Join("..", ".env")})
	checkRequiredEnvVars()

	// 加载配置 & 初始化数据库连接
	config.InitConfig()
	database.InitDB()

	r := gin.Default()
	// 只信任本机代理，避免 Gin "trust all proxies" 的安全警告刷屏
	_ = r.SetTrustedProxies([]string{"127.0.0.1"})
	// 设置文件上传大小限制为10MB
	r.MaxMultipartMemory = 10 << 20 // 10MB
	r.Static("/uploads", "./uploads")
	r.Use(middleware.CorsMiddleware())
	router.RouterInit(r)

	addr := getListenAddr()
	if err := r.Run(addr); err != nil {
		panic(err)
	}
}

// getListenAddr 获取服务监听地址。
// 默认只监听 IPv4（127.0.0.1），端口默认 8080，可通过环境变量覆盖：
// - SERVER_HOST（默认 127.0.0.1）
// - SERVER_PORT（默认 8080）
func getListenAddr() string {
	host := strings.TrimSpace(os.Getenv("SERVER_HOST"))
	if host == "" {
		host = "127.0.0.1"
	}

	port := strings.TrimSpace(os.Getenv("SERVER_PORT"))
	if port == "" {
		port = "8080"
	}

	return net.JoinHostPort(host, port)
}

// loadDotEnv 是一个轻量的 .env 读取器（不引入额外依赖）。
// 规则：
// - 只设置当前进程中尚未存在的环境变量（不覆盖外部传入的 env）
// - 支持：KEY=VALUE / export KEY=VALUE / # 注释行
func loadDotEnv(candidatePaths []string) {
	for _, p := range candidatePaths {
		data, err := os.ReadFile(p)
		if err != nil {
			continue
		}

		scanner := bufio.NewScanner(bytes.NewReader(data))
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			if strings.HasPrefix(line, "export ") {
				line = strings.TrimSpace(strings.TrimPrefix(line, "export "))
			}

			key, value, ok := strings.Cut(line, "=")
			if !ok {
				continue
			}
			key = strings.TrimSpace(key)
			value = strings.TrimSpace(value)
			if key == "" {
				continue
			}

			value = strings.Trim(value, "\"'")

			if os.Getenv(key) == "" {
				_ = os.Setenv(key, value)
			}
		}

		break
	}
}

// checkRequiredEnvVars 检查必需的环境变量，并在生产环境给出安全提示。
func checkRequiredEnvVars() {
	jwtKey := os.Getenv("JWT_SECRET_KEY")

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
			fmt.Println("   当前将使用默认允许源列表，请在 .env 文件中配置")
		}

		dbHost := os.Getenv("DB_HOST")
		dbPassword := os.Getenv("DB_PASSWORD")
		if dbHost == "127.0.0.1" || dbPassword == "123456" {
			fmt.Println("⚠️  [安全警告] 检测到默认数据库配置")
			fmt.Println("   请在 .env 文件中更新数据库凭据")
		}
		return
	}

	if jwtKey == "" {
		fmt.Println("ℹ️  [提示] 未设置 JWT_SECRET_KEY，将使用默认开发密钥（仅建议本地开发）")
	}
	fmt.Println("ℹ️  [提示] 当前运行在开发模式，优先从仓库根目录 .env 读取配置（如存在）")
}
