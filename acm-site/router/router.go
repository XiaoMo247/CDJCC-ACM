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
			admin.POST("/register", api.AdminBatchRegisterHandler)
			admin.GET("/me", api.GetAdminInfo)
			admin.POST("/add", api.AddAdmin)
			admin.DELETE("/delete/:id", api.DeleteAdmin)
			admin.GET("/list", api.ListAdmins)
		}

		// 用户模块：普通用户登录
		apiGroup.POST("/user/login", api.UserLoginHandler)

		user := apiGroup.Group("/user")
		user.Use(middleware.UserAuthMiddleware()) // 用户验证中间件
		{
			// 如果你之后有用户信息接口，比如 `/user/me`，可以放这里
			// user.GET("/me", api.GetUserInfo)
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

		// 加入我们模块
		join := apiGroup.Group("/join")

		// 用户权限：必须登录
		join.Use(middleware.UserAuthMiddleware())
		{
			join.POST("/send", api.SubmitJoinApplyHandler) // 提交申请表
			join.GET("/my", api.GetMyJoinApplyHandler)     // 获取自己申请
		}

		// 管理员权限：单独分组，避免被 UserAuthMiddleware 拦住
		joinAdmin := apiGroup.Group("/join")
		joinAdmin.Use(middleware.AdminAuthMiddleware())
		{
			joinAdmin.GET("/get", api.GetAllJoinAppliesHandler)   // 所有申请表
			joinAdmin.GET("/:id", api.GetJoinApplyByIDHandler)    // 查看某一条详情
			joinAdmin.PUT("/ac/:id", api.ApproveJoinApplyHandler) // 通过某条请求
			joinAdmin.PUT("/wa/:id", api.RejectJoinApply)         //拒绝某条请求
		}
	}
}
