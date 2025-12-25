package router

import (
	"acm-site/api"
	"acm-site/middleware"

	"github.com/gin-gonic/gin"
)

// RouterUnifiedInit 使用统一认证中间件的路由配置示例
// 这是新版本的路由配置，使用统一的 token 管理
// 可以逐步替换 router.go 中的旧配置
func RouterUnifiedInit(r *gin.Engine) {
	apiGroup := r.Group("/api")
	{
		// ======================== 公开接口（无需认证） ========================
		// 登录接口
		apiGroup.POST("/admin/login", api.AdminLogin)
		apiGroup.POST("/user/login", api.UserLoginHandler)
		apiGroup.POST("/student/login", api.StudentLogin)

		// 公共查询接口
		apiGroup.GET("/contest/list", api.GetAllContests)
		apiGroup.GET("/faq/list", api.GetFAQList)
		apiGroup.GET("/slider/list", api.ListSliders)
		apiGroup.GET("/honor", api.GetAllHonors)
		apiGroup.GET("/admin/team-members", api.GetAllTeamMembers)
		apiGroup.GET("/announcement/list", api.ListAnnouncements)

		// 课件资源（无需登录）
		file := apiGroup.Group("/folder")
		{
			file.GET("/files", api.ListFiles)
			file.GET("/download/:id", api.DownloadFile)
			file.GET("/list", api.ListFolders)
		}

		// ======================== 需要认证的接口 ========================
		auth := apiGroup.Group("/")
		auth.Use(middleware.AuthMiddleware()) // 统一的认证中间件
		{
			// -------- 获取当前用户信息（任何登录用户都可以） --------
			auth.GET("/admin/me", api.GetAdminInfo) // 管理员信息（兼容性）
			auth.GET("/user/info", api.GetUserInfo) // 用户信息（兼容性）

			// ======================== 管理员专属接口 ========================
			admin := auth.Group("/admin")
			admin.Use(middleware.RequireAdmin()) // 要求管理员角色
			{
				// -------- 学生账号管理 --------
				admin.POST("/register-student", api.AdminRegisterStudent)
				admin.GET("/team-member", api.GetStudentByID)
				admin.DELETE("/delete-student", api.DeleteStudent)
				admin.POST("/reset-password", api.ResetPassword)
				admin.POST("/update-ratings", api.UpdateRatingsHandler)
				admin.POST("/update-all-ratings", api.UpdateAllRatingsHandler)

				// -------- 用户账号管理 --------
				admin.POST("/register", api.AdminBatchRegisterHandler)
				admin.GET("/get/all-user", api.AdminGetAllUserHandler)
				admin.GET("/get/user", api.AdminGetUsersHandler)
				admin.DELETE("/user/:id", api.AdminDeleteUserHandler)
				admin.PUT("/user/reset/:id", api.AdminResetPasswordHandler)

				// -------- 管理员账户管理 --------
				admin.POST("/add", api.AddAdmin)
				admin.DELETE("/delete/:id", api.DeleteAdmin)
				admin.GET("/list", api.ListAdmins)

				// -------- 比赛管理 --------
				admin.POST("/contest/upload", api.UploadContest)
				admin.DELETE("/contest/:id", api.DeleteContest)

				// -------- 文件夹管理 --------
				adminFolder := admin.Group("/folder")
				{
					adminFolder.POST("/upload", api.UploadFile)
					adminFolder.DELETE("/file/:id", api.DeleteFile)
					adminFolder.POST("/add", api.AddFolder)
					adminFolder.DELETE("/delete/:id", api.DeleteFolder)
				}

				// -------- FAQ 管理 --------
				adminFaq := admin.Group("/faq")
				{
					adminFaq.POST("/add", api.CreateFAQ)
					adminFaq.DELETE("/delete/:id", api.DeleteFAQ)
				}

				// -------- 轮播图管理 --------
				slider := admin.Group("/slider")
				{
					slider.POST("/add", api.AddSlider)
					slider.DELETE("/delete/:id", api.DeleteSlider)
				}

				// -------- 荣誉墙管理 --------
				honor := admin.Group("/honor")
				{
					honor.POST("", api.CreateHonor)
					honor.DELETE("/:id", api.DeleteHonor)
					honor.PUT("/:id", api.UpdateHonor)
				}

				// -------- 公告管理 --------
				announcement := admin.Group("/announcement")
				{
					announcement.POST("/create", api.CreateAnnouncement)
					announcement.PUT("/update/:id", api.UpdateAnnouncement)
					announcement.DELETE("/delete/:id", api.DeleteAnnouncement)
				}

				// -------- 申请审核 --------
				joinAdmin := admin.Group("/join")
				{
					joinAdmin.GET("/get", api.GetAllJoinAppliesHandler)
					joinAdmin.GET("/:id", api.GetJoinApplyByIDHandler)
					joinAdmin.PUT("/ac/:id", api.ApproveJoinApplyHandler)
					joinAdmin.PUT("/wa/:id", api.RejectJoinApply)
				}
			}

			// ======================== 学生专属接口 ========================
			user := auth.Group("/user")
			user.Use(middleware.RequireStudent()) // 要求学生角色
			{
				user.POST("/change-username", api.ChangeUsername)
				user.POST("/change-password", api.ChangePassword)
			}

			// ======================== 队员专属接口 ========================
			student := auth.Group("/student")
			student.Use(middleware.RequireMember()) // 要求队员角色
			{
				student.GET("/info", api.GetStudentInfo)
				student.POST("/update-password", api.UpdateStudentPassword)
				student.POST("/update-username", api.UpdateStudentUsername)
				student.POST("/update-info", api.UpdateStudentInfo)
			}

			// ======================== 申请模块（需要学生登录） ========================
			join := auth.Group("/join")
			join.Use(middleware.RequireStudent()) // 要求学生角色
			{
				join.POST("/send", api.SubmitJoinApplyHandler)
				join.GET("/my", api.GetMyJoinApplyHandler)
			}
		}
	}
}
