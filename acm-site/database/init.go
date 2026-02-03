package database

import (
	"acm-site/config"
	"acm-site/model"
	"acm-site/utils"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func generateSecurePassword(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		panic("failed to generate password: " + err.Error())
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length]
}

func InitDB() {
	c := config.GlobalConfig.MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.DBName)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("database connection failed: " + err.Error())
	}

	// 预迁移：在AutoMigrate之前清理数据，避免外键约束失败
	prepareSliderMigration()

	if err := DB.AutoMigrate(
		&model.Admin{},
		&model.Announcement{},
		&model.User{},
		&model.JoinApply{},
		&model.Folder{},
		&model.File{},
		&model.TeamMember{},
		&model.Contest{},
		&model.FAQ{},
		&model.Slider{},
		&model.Honor{},
		&model.Image{}, // 新增Image模型
	); err != nil {
		panic("database migration failed: " + err.Error())
	}

	// Best-effort schema upgrade for legacy folders table.
	ensureFolderSchema()

	// Image management system migrations
	ensureImageSchema()

	// Seed default admin account (admin/admin or generated).
	var count int64
	DB.Model(&model.Admin{}).Where("username = ?", "admin").Count(&count)
	if count == 0 {
		initialPassword := os.Getenv("INITIAL_ADMIN_PASSWORD")
		if initialPassword == "" {
			initialPassword = generateSecurePassword(16)
			fmt.Println("=========================================")
			fmt.Println("Default admin created")
			fmt.Println("  username: admin")
			fmt.Println("  password:", initialPassword)
			fmt.Println("  Please log in and change it ASAP.")
			fmt.Println("=========================================")
		} else {
			fmt.Println("Created default admin from INITIAL_ADMIN_PASSWORD")
		}

		admin := model.Admin{
			Username: "admin",
			Password: utils.HashPassword(initialPassword),
		}
		DB.Create(&admin)
	}
}

// prepareSliderMigration 预迁移处理，清理sliders表中的无效image_id
func prepareSliderMigration() {
	// 检查sliders表是否存在
	var tableExists int64
	DB.Raw("SELECT COUNT(*) FROM information_schema.TABLES WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sliders'").Scan(&tableExists)

	if tableExists == 0 {
		return // 表不存在，跳过
	}

	// 检查image_id列是否存在
	var columnExists int64
	DB.Raw("SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sliders' AND COLUMN_NAME = 'image_id'").Scan(&columnExists)

	if columnExists > 0 {
		// 将无效的image_id设置为NULL（为外键约束做准备）
		DB.Exec("UPDATE sliders SET image_id = NULL WHERE image_id = 0 OR image_id IS NOT NULL AND image_id NOT IN (SELECT id FROM images)")
		fmt.Println("已清理sliders表中的无效image_id")
	}
}

