package router

import (
	"acm-site/api"
	"acm-site/middleware"

	"github.com/gin-gonic/gin"
)

func RouterInit(r *gin.Engine) {
	apiGroup := r.Group("/api")
	{
		// ======================== 管理员模块 ========================
		apiGroup.POST("/admin/login", api.AdminLogin) // 管理员登录

		admin := apiGroup.Group("/admin")
		admin.Use(middleware.AdminAuthMiddleware()) // 管理员身份验证中间件
		{
			// -------- 学生账号管理 --------
			admin.POST("/register-student", api.AdminRegisterStudent)      // 注册团队成员
			admin.GET("/team-member", api.GetStudentByID)                  // 获取指定成员信息（通过 ID）
			admin.DELETE("/delete-student", api.DeleteStudent)             // 删除成员
			admin.POST("/reset-password", api.ResetPassword)               // 重置成员密码
			admin.POST("/update-ratings", api.UpdateRatingsHandler)        // 更新指定成员的评分
			admin.POST("/update-all-ratings", api.UpdateAllRatingsHandler) // 更新所有成员评分

			// -------- 用户账号管理（登录系统的用户） --------
			admin.POST("/register", api.AdminBatchRegisterHandler)      // 区间批量注册用户
			admin.GET("/get/all-user", api.AdminGetAllUserHandler)      // 获取所有用户
			admin.GET("/get/user", api.AdminGetUsersHandler)            // 分页 & 搜索用户
			admin.DELETE("/user/:id", api.AdminDeleteUserHandler)       // 删除用户
			admin.PUT("/user/reset/:id", api.AdminResetPasswordHandler) // 重置用户密码

			// -------- 管理员账户管理 --------
			admin.GET("/me", api.GetAdminInfo)           // 获取当前管理员信息
			admin.POST("/add", api.AddAdmin)             // 添加管理员
			admin.DELETE("/delete/:id", api.DeleteAdmin) // 删除管理员
			admin.GET("/list", api.ListAdmins)           // 获取所有管理员列表

			// -------- 比赛上传管理 --------
			admin.POST("/contest/upload", api.UploadContest) // 上传比赛
			admin.DELETE("/contest/:id", api.DeleteContest)  // 删除比赛

			// -------- 课件管理（文件上传、删除） --------
			file := admin.Group("/folder")
			{
				file.POST("/upload", api.UploadFile)     // 上传文件
				file.DELETE("/file/:id", api.DeleteFile) // 删除文件
			}

			// -------- FAQ 管理 --------
			faq := admin.Group("/faq")
			{
				faq.POST("/add", api.CreateFAQ)          // 添加 FAQ
				faq.DELETE("/delete/:id", api.DeleteFAQ) // 删除 FAQ
			}

			// --------- 轮播图管理 --------
			slider := admin.Group("/slider")
			{
				slider.POST("/add", api.AddSlider)
				slider.DELETE("/delete/:id", api.DeleteSlider)
			}

			// -------- 荣誉墙管理 --------
			honorGroup := admin.Group("/honor")
			{
				honorGroup.POST("", api.CreateHonor)
				honorGroup.DELETE("/:id", api.DeleteHonor)
				honorGroup.PUT("/:id", api.UpdateHonor)
			}
		}

		// ======================== 课件资源下载模块（所有用户） ========================
		file := apiGroup.Group("/folder")
		{
			file.GET("/files", api.ListFiles)           // 获取文件列表
			file.GET("/download/:id", api.DownloadFile) // 下载指定文件
		}

		// -------- 文件夹管理（需要管理员权限） --------
		folder := apiGroup.Group("/folder")
		folder.GET("/list", api.ListFolders) // 获取所有文件夹
		folder.Use(middleware.AdminAuthMiddleware())
		{
			folder.POST("/add", api.AddFolder)             // 添加文件夹
			folder.DELETE("/delete/:id", api.DeleteFolder) // 删除文件夹
		}

		// ======================== 普通用户模块 ========================
		apiGroup.POST("/user/login", api.UserLoginHandler) // 用户登录

		userGroup := apiGroup.Group("/user")
		userGroup.Use(middleware.UserAuthMiddleware())
		{
			userGroup.POST("/change-username", api.ChangeUsername) // 修改用户名
			userGroup.POST("/change-password", api.ChangePassword) // 修改密码
			userGroup.GET("/info", api.GetUserInfo)                // 获取当前用户信息
		}

		// ======================== 公告模块 ========================
		announcement := apiGroup.Group("/announcement")
		{
			announcement.GET("/list", api.ListAnnouncements) // 公共获取公告列表（无需登录）

			// 管理员权限：增删改公告
			announcement.Use(middleware.AdminAuthMiddleware())
			{
				announcement.POST("/create", api.CreateAnnouncement)       // 创建公告
				announcement.PUT("/update/:id", api.UpdateAnnouncement)    // 更新公告
				announcement.DELETE("/delete/:id", api.DeleteAnnouncement) // 删除公告
			}
		}

		// ======================== 加入我们模块 ========================
		join := apiGroup.Group("/join")

		// 普通用户权限：提交申请
		join.Use(middleware.UserAuthMiddleware())
		{
			join.POST("/send", api.SubmitJoinApplyHandler) // 提交申请
			join.GET("/my", api.GetMyJoinApplyHandler)     // 查看自己的申请
		}

		// 管理员权限：审核申请（注意避免 UserMiddleware 拦住）
		joinAdmin := apiGroup.Group("/join")
		joinAdmin.Use(middleware.AdminAuthMiddleware())
		{
			joinAdmin.GET("/get", api.GetAllJoinAppliesHandler)   // 获取所有申请
			joinAdmin.GET("/:id", api.GetJoinApplyByIDHandler)    // 查看指定申请
			joinAdmin.PUT("/ac/:id", api.ApproveJoinApplyHandler) // 审核通过
			joinAdmin.PUT("/wa/:id", api.RejectJoinApply)         // 审核拒绝
		}

		// ======================== 队员模块 ========================
		studentGroup := apiGroup.Group("/student")
		studentGroup.POST("/login", api.StudentLogin) // 队员登录
		studentGroup.Use(middleware.StudentAuthMiddleware())
		{
			studentGroup.GET("/info", api.GetStudentInfo)                    // 获取队员信息
			studentGroup.POST("/update-password", api.UpdateStudentPassword) // 修改密码
			studentGroup.POST("/update-username", api.UpdateStudentUsername) // 修改用户名
			studentGroup.POST("/update-info", api.UpdateStudentInfo)         // 更新信息（新增）
		}

		// ======================== 比赛与FAQ公共接口 ========================
		apiGroup.GET("/contest/list", api.GetAllContests)          // 获取比赛列表
		apiGroup.GET("/faq/list", api.GetFAQList)                  // 获取 FAQ 列表
		apiGroup.GET("/admin/team-members", api.GetAllTeamMembers) // 获取所有团队成员信息

		// ========================轮播图公共接口 ========================
		apiGroup.GET("/slider/list", api.ListSliders)

		// ========================荣誉墙公共接口 ========================
		apiGroup.GET("/honor", api.GetAllHonors)
	}
}
