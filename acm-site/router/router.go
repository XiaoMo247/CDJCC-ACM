package router

import (
	"acm-site/api"
	"acm-site/middleware"

	"github.com/gin-gonic/gin"
)

func RouterInit(r *gin.Engine) {
	apiGroup := r.Group("/api")
	{
		// Public (no auth)
		apiGroup.POST("/admin/login", api.AdminLogin)
		apiGroup.POST("/user/login", api.UserLoginHandler)
		apiGroup.POST("/student/login", api.StudentLogin)

		apiGroup.GET("/contest/list", api.GetAllContests)
		apiGroup.GET("/faq/list", api.GetFAQList)
		apiGroup.GET("/slider/list", api.ListSliders)
		apiGroup.GET("/honor", api.GetAllHonors)
		apiGroup.GET("/admin/team-members", api.GetAllTeamMembers)
		apiGroup.GET("/announcement/list", api.ListAnnouncements)

		// Folder resources (public)
		folder := apiGroup.Group("/folder")
		{
			folder.GET("/files", api.ListFiles)
			folder.GET("/download/:id", api.DownloadFile)
			folder.GET("/list", api.ListFolders)
			folder.GET("/tree", api.GetFolderTree)
			folder.GET("/:id/content", api.GetFolderContent)
			folder.GET("/:id/breadcrumb", api.GetFolderBreadcrumb)
			folder.GET("/search", api.SearchPublic)
		}

		// Authenticated
		auth := apiGroup.Group("/")
		auth.Use(middleware.AuthMiddleware())
		{
			auth.GET("/admin/me", api.GetAdminInfo)
			auth.GET("/user/info", api.GetUserInfo)

			// Admin-only
			admin := auth.Group("/admin")
			admin.Use(middleware.RequireAdmin())
			{
				// Dashboard统计
				admin.GET("/dashboard/stats", api.GetDashboardStats)

				// Students
				admin.POST("/register-student", api.AdminRegisterStudent)
				admin.GET("/team-member", api.GetStudentByID)
				admin.DELETE("/delete-student", api.DeleteStudent)
				admin.POST("/reset-password", api.ResetPassword)
				admin.POST("/update-ratings", api.UpdateRatingsHandler)
				admin.POST("/update-all-ratings", api.UpdateAllRatingsHandler)

				// Users
				admin.POST("/register", api.AdminBatchRegisterHandler)
				admin.GET("/get/all-user", api.AdminGetAllUserHandler)
				admin.GET("/get/user", api.AdminGetUsersHandler)
				admin.DELETE("/user/:id", api.AdminDeleteUserHandler)
				admin.PUT("/user/reset/:id", api.AdminResetPasswordHandler)

				// Admin accounts
				admin.POST("/add", api.AddAdmin)
				admin.DELETE("/delete/:id", api.DeleteAdmin)
				admin.GET("/list", api.ListAdmins)

				// Contests
				admin.POST("/contest/upload", api.UploadContest)
				admin.DELETE("/contest/:id", api.DeleteContest)

				// Folder manager (admin)
				adminFolder := admin.Group("/folder")
				{
					// New nested folder APIs
					adminFolder.GET("/tree", api.AdminGetFolderTree)
					adminFolder.GET("/:id/content", api.AdminGetFolderContent)
					adminFolder.GET("/:id/breadcrumb", api.AdminGetFolderBreadcrumb)
					adminFolder.POST("/create", api.AdminCreateFolder)
					adminFolder.POST("/:id/upload", api.AdminUploadFileToFolder)

					// Legacy APIs (kept for backward compatibility)
					adminFolder.POST("/upload", api.UploadFile)
					adminFolder.DELETE("/file/:id", api.DeleteFile)
					adminFolder.POST("/add", api.AddFolder)
					adminFolder.DELETE("/delete/:id", api.DeleteFolder)
				}

				admin.GET("/search", api.AdminSearch)
				admin.POST("/item/move", api.AdminMoveItem)

				// FAQ
				adminFaq := admin.Group("/faq")
				{
					adminFaq.POST("/add", api.CreateFAQ)
					adminFaq.DELETE("/delete/:id", api.DeleteFAQ)
				}

				// Slider
				slider := admin.Group("/slider")
				{
					slider.POST("/add", api.AddSlider)
					slider.DELETE("/delete/:id", api.DeleteSlider)
				}

				// Honor
				honor := admin.Group("/honor")
				{
					honor.POST("", api.CreateHonor)
					honor.DELETE("/:id", api.DeleteHonor)
					honor.PUT("/:id", api.UpdateHonor)
				}

				// Announcement
				announcement := admin.Group("/announcement")
				{
					announcement.POST("/create", api.CreateAnnouncement)
					announcement.PUT("/update/:id", api.UpdateAnnouncement)
					announcement.DELETE("/delete/:id", api.DeleteAnnouncement)
				}

				// Join review
				joinAdmin := admin.Group("/join")
				{
					joinAdmin.GET("/get", api.GetAllJoinAppliesHandler)
					joinAdmin.GET("/:id", api.GetJoinApplyByIDHandler)
					joinAdmin.PUT("/ac/:id", api.ApproveJoinApplyHandler)
					joinAdmin.PUT("/wa/:id", api.RejectJoinApply)
				}
			}

			// Student-only
			user := auth.Group("/user")
			user.Use(middleware.RequireStudent())
			{
				user.POST("/change-username", api.ChangeUsername)
				user.POST("/change-password", api.ChangePassword)
			}

			// Member-only
			student := auth.Group("/student")
			student.Use(middleware.RequireMember())
			{
				student.GET("/info", api.GetStudentInfo)
				student.POST("/update-password", api.UpdateStudentPassword)
				student.POST("/update-username", api.UpdateStudentUsername)
				student.POST("/update-info", api.UpdateStudentInfo)
			}

			// Join module (student)
			join := auth.Group("/join")
			join.Use(middleware.RequireStudent())
			{
				join.POST("/send", api.SubmitJoinApplyHandler)
				join.GET("/my", api.GetMyJoinApplyHandler)
			}
		}
	}
}
