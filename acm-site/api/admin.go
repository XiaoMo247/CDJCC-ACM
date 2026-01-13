package api

import (
	"acm-site/database"
	"acm-site/model"
	"acm-site/service"
	"acm-site/utils"
	"acm-site/utils/jwt"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func AdminLogin(c *gin.Context) {
	var loginDetails struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&loginDetails); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求数据不合法"})
		return
	}

	var admin model.Admin
	if err := database.DB.Where("username = ?", loginDetails.Username).First(&admin).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "用户名或密码错误"})
		return
	}

	if !utils.CheckPassword(loginDetails.Password, admin.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "用户名或密码错误"})
		return
	}

	// 使用统一的 token 生成函数
	token, err := jwt.GenerateUnifiedToken(admin.ID, admin.Username, "admin")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "生成 token 失败"})
		return
	}

	// 返回统一格式：token + user 对象
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":       admin.ID,
			"username": admin.Username,
			"role":     "admin",
		},
	})
}

// 添加新管理员接口
func AddAdmin(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "无效请求"})
		return
	}

	// 检查用户名是否已存在
	var existingAdmin model.Admin
	if err := database.DB.Where("username = ?", req.Username).First(&existingAdmin).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"msg": "用户名已存在"})
		return
	}

	// 创建新管理员
	newAdmin := model.Admin{
		Username: req.Username,
		Password: utils.HashPassword(req.Password),
	}

	if err := database.DB.Create(&newAdmin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "创建新管理员失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "新管理员已添加"})
}

func GetAdminInfo(c *gin.Context) {
	// 从统一中间件获取用户信息
	username, _ := c.Get("username")
	id, _ := c.Get("user_id")

	c.JSON(http.StatusOK, gin.H{
		"msg":      "你是合法管理员",
		"username": username,
		"id":       id,
	})
}

// 获取所有管理员列表（任意登录管理员可用）
func ListAdmins(c *gin.Context) {
	var admins []model.Admin

	if err := database.DB.Select("id", "username", "created_at").Find(&admins).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "获取管理员列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"admins": admins,
	})
}

// 删除管理员接口（仅初代管理员可用）
func DeleteAdmin(c *gin.Context) {
	// 从统一中间件获取当前管理员 ID
	adminIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "无权限"})
		return
	}

	adminID := fmt.Sprintf("%v", adminIDVal)
	if adminID != "1" {
		c.JSON(http.StatusForbidden, gin.H{"msg": "只有初代管理员可以删除其他管理员"})
		return
	}

	// 获取 URL 参数中的 ID（要删除的管理员 ID）
	targetID := c.Param("id")
	if targetID == "1" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "初代管理员不可删除"})
		return
	}

	// 检查是否存在该管理员
	var target model.Admin
	if err := database.DB.First(&target, targetID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": "管理员不存在"})
		return
	}

	// 删除管理员
	if err := database.DB.Delete(&target).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "管理员删除成功"})
}

type BatchRegisterRequest struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

// 批量注册
func AdminBatchRegisterHandler(c *gin.Context) {
	var req BatchRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数格式错误"})
		return
	}

	created, err := service.RegisterUserBatch(req.Start, req.End)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "批量注册成功",
		"createdList": created,
		"count":       len(created),
	})
}

// 获取所有用户
func AdminGetAllUserHandler(c *gin.Context) {
	users, err := service.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "获取成功",
		"data":    users,
		"count":   len(users),
	})
}

// 获取用户列表（分页 + 搜索）
func AdminGetUsersHandler(c *gin.Context) {
	keyword := c.Query("keyword")
	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	users, total, err := service.GetUsers(keyword, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  users,
		"total": total,
	})
}

// 删除用户
func AdminDeleteUserHandler(c *gin.Context) {
	id := c.Param("id")
	err := service.DeleteUserByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除用户失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// 重置密码（重置为默认密码）
func AdminResetPasswordHandler(c *gin.Context) {
	id := c.Param("id")
	err := service.ResetUserPassword(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "重置密码失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "密码已重置为默认密码"})
}

type RegisterStudentRequest struct {
	StudentID string `json:"student_id" binding:"required"`
}

func AdminRegisterStudent(c *gin.Context) {
	var req RegisterStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	// 学号长度检查
	if len(req.StudentID) != 9 {

		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "学号格式错误"})
		return
	}

	// 截取后六位
	suffix := req.StudentID[len(req.StudentID)-6:]
	initPw := "JCACM" + suffix
	hashedPw := utils.HashPassword(initPw)

	student := model.TeamMember{
		StudentID: req.StudentID,
		Password:  hashedPw,
		Username:  req.StudentID, // 初始用户名用学号
	}

	if err := database.DB.Create(&student).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "注册失败，可能该用户已存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"msg":     "注册成功",
		"data":    student,
		"raw_pwd": initPw, // 可选返回初始密码（前端提醒管理员记得通知）
	})
}

func GetAllTeamMembers(c *gin.Context) {
	var members []model.TeamMember

	// 查询所有成员，不包括密码和更新时间
	if err := database.DB.Select("id", "student_id", "username", "cf_name", "cf_rating", "at_name", "at_rating", "nc_id", "nc_rating", "rating", "count", "created_at").Find(&members).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "查询成功",
		"data": members,
	})
}

// 根据学号查询学生
func GetStudentByID(c *gin.Context) {
	studentID := c.DefaultQuery("student_id", "") // 获取查询参数中的学号
	if studentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "学号不能为空"})
		return
	}

	student, err := service.GetTeamMemberByStudentID(studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "查询学生信息失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": student})
}

// 删除学生
func DeleteStudent(c *gin.Context) {
	studentID := c.DefaultQuery("student_id", "")
	if studentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "学号不能为空"})
		return
	}

	err := service.DeleteStudent(studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "删除学生失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// 重置密码
func ResetPassword(c *gin.Context) {
	studentID := c.DefaultQuery("student_id", "")
	if studentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "学号不能为空"})
		return
	}

	// 根据学号生成新密码
	newPassword := "JCACM" + studentID[len(studentID)-6:] // 取学号后6位生成密码

	// 更新密码
	err := service.UpdatePassword(studentID, newPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "重置密码失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "密码重置成功", "new_password": newPassword})
}

func UpdateRatingsHandler(c *gin.Context) {
	// 从请求中获取学生 ID 和相应的用户名
	var request struct {
		StudentID  string `json:"student_id"` // 学生 ID
		Codeforces string `json:"codeforces"` // Codeforces 用户名
		Atcoder    string `json:"atcoder"`    // AtCoder 用户名
		Nowcoder   string `json:"nowcoder"`   // Nowcoder 用户名
	}

	// 解析请求的 JSON 数据
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "请求数据格式不正确", "error": err.Error()})
		return
	}

	// 调用外部 API 获取分数
	codeforcesRating, err := utils.FetchCodeforcesRating(request.Codeforces)
	if err != nil {
		codeforcesRating = ""
	}

	atcoderRating, err := utils.FetchAtCoderRating(request.Atcoder)
	if err != nil {
		atcoderRating = ""
	}

	nowcoderRating, err := utils.FetchNowcoderRating(request.Nowcoder)
	if err != nil {
		nowcoderRating = ""
	}

	// 更新数据库中的评分
	err = service.UpdateRatings(request.StudentID, codeforcesRating, atcoderRating, nowcoderRating)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "更新评分失败", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "评分更新成功"})
}

func getAllTeamMembers() ([]model.TeamMember, error) {
	var members []model.TeamMember
	if err := database.DB.Find(&members).Error; err != nil {
		return nil, err
	}
	return members, nil
}

func UpdateAllRatingsHandler(c *gin.Context) {
	members, err := getAllTeamMembers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "获取成员列表失败",
		})
		return
	}

	var wg sync.WaitGroup
	var mu sync.Mutex // 用于线程安全统计 success/fail
	successCount := 0
	failCount := 0

	for i := range members {
		wg.Add(1)
		go func(member *model.TeamMember) {
			defer wg.Done()
			err := service.UpdateRatingsByMember(member)
			mu.Lock()
			if err != nil {
				failCount++
			} else {
				successCount++
			}
			mu.Unlock()
		}(&members[i]) // 注意必须传指针 &members[i]
	}

	wg.Wait() // 等待所有协程完成

	c.JSON(http.StatusOK, gin.H{
		"message":      "所有成员分数更新完成",
		"successCount": successCount,
		"failCount":    failCount,
	})
}
