package router

import (
	"acm-site/api"
	"acm-site/middleware"

	"github.com/gin-gonic/gin"
)

func RouterInit(r *gin.Engine) {
	apiGroup := r.Group("/api")
	{
		// 管理员模块
		apiGroup.POST("/admin/login", api.AdminLogin)

		admin := apiGroup.Group("/admin")
		admin.Use(middleware.AdminAuthMiddleware()) // 管理员验证中间件
		{
			admin.GET("/me", api.GetAdminInfo)
			admin.POST("/add", api.AddAdmin)
			admin.DELETE("/delete/:id", api.DeleteAdmin)
			admin.GET("/list", api.ListAdmins)
		}

		// 公告模块
		announcement := apiGroup.Group("/announcement")
		{
			announcement.GET("/list", api.ListAnnouncements) // 公共获取公告

			// 以下需要管理员权限
			announcement.Use(middleware.AdminAuthMiddleware())
			{
				announcement.POST("/create", api.CreateAnnouncement)
				announcement.PUT("/update/:id", api.UpdateAnnouncement)
				announcement.DELETE("/delete/:id", api.DeleteAnnouncement)
			}
		}
	}
}
