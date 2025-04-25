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
			// student
			admin.POST("/register-student", api.AdminRegisterStudent) //注册团队成员
			admin.GET("/team-members", api.GetAllTeamMembers)         // 获取所有成员信息
			admin.GET("/team-member", api.GetStudentByID)             // 新增查询接口
			admin.DELETE("/delete-student", api.DeleteStudent)
			admin.POST("/reset-password", api.ResetPassword)
			admin.POST("/update-ratings", api.UpdateRatingsHandler)
			admin.POST("/update-all-ratings", api.UpdateAllRatingsHandler)

			// user
			admin.POST("/register", api.AdminBatchRegisterHandler)      // 区间批量注册用户
			admin.GET("/get/all-user", api.AdminGetAllUserHandler)      // 得到所有用户
			admin.GET("/get/user", api.AdminGetUsersHandler)            // 查询用户(分页查询 + 模糊搜索)
			admin.DELETE("/user/:id", api.AdminDeleteUserHandler)       // 删除用户
			admin.PUT("/user/reset/:id", api.AdminResetPasswordHandler) // 重置用户

			admin.GET("/me", api.GetAdminInfo)
			admin.POST("/add", api.AddAdmin)
			admin.DELETE("/delete/:id", api.DeleteAdmin)
			admin.GET("/list", api.ListAdmins)

			file := admin.Group("/folder")
			{
				file.POST("/upload", api.UploadFile)
				file.DELETE("/file/:id", api.DeleteFile)
			}
		}
		file := apiGroup.Group("/folder")
		{
			file.GET("/files", api.ListFiles)
			file.GET("/download/:id", api.DownloadFile)
		}
		// 课件管理模块
		folder := apiGroup.Group("/folder")
		folder.GET("/list", api.ListFolders)
		folder.Use(middleware.AdminAuthMiddleware()) // 管理员权限
		{
			folder.POST("/add", api.AddFolder)
			folder.DELETE("/delete/:id", api.DeleteFolder)
		}

		// 用户模块：普通用户登录
		apiGroup.POST("/user/login", api.UserLoginHandler)

		userGroup := apiGroup.Group("/user")
		userGroup.Use(middleware.UserAuthMiddleware()) // 用户验证中间件
		{
			userGroup.POST("/change-username", api.ChangeUsername) // 修改用户名
			userGroup.POST("/change-password", api.ChangePassword) // 修改密码
			userGroup.GET("/info", api.GetUserInfo)
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

		studentGroup := apiGroup.Group("/student")
		studentGroup.POST("/login", api.StudentLogin)        //学生登录
		studentGroup.Use(middleware.StudentAuthMiddleware()) // 下面这些都需要登录后才可操作
		{
			studentGroup.GET("/info", api.GetStudentInfo)
			studentGroup.POST("/update-password", api.UpdateStudentPassword)
			studentGroup.POST("/update-username", api.UpdateStudentUsername)
			studentGroup.POST("/update-info", api.UpdateStudentInfo) // 新增的修改信息接口
		}
	}
}
