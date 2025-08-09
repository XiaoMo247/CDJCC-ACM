package main

import (
	"acm-site/config"
	"acm-site/database"
	"acm-site/middleware"
	"acm-site/router"

	"github.com/gin-gonic/gin"
)

func main() {
	config.InitConfig() // 加载配置
	database.InitDB()   // 初始化数据库连接

	r := gin.Default()

	r.Static("/uploads", "./uploads")

	r.Use(middleware.CorsMiddleware())

	router.RouterInit(r)

	r.Run(":8081")
}
