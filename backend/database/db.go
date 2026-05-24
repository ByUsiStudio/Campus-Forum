package database

import (
	"crypto/rand"
	"encoding/hex"
	"forum/models"
	"forum/utils"
	"os"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(dsn string) {
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		utils.Error("数据库连接失败: %v", err)
		os.Exit(1)
	}
	utils.Success("数据库连接成功")
}

func AutoMigrate() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Article{},
		&models.Comment{},
		&models.CommentLike{},
		&models.ViewHistory{},
		&models.Like{},
		&models.Category{},
		&models.DeletionRequest{},
		&models.SidebarConfig{},
		&models.Announcement{},
		&models.Notification{},
		&models.UserNotification{},
		&models.Follow{},
		&models.ChatMessage{},
		&models.ChatSession{},
		&models.SiteConfig{},
	)
	if err != nil {
		utils.Error("数据库迁移失败: %v", err)
		os.Exit(1)
	}
	utils.Success("数据库迁移完成")
}

func CheckAndInitAdmin() {
	var count int64
	DB.Model(&models.User{}).Count(&count)

	if count == 0 {
		utils.Warn("未检测到用户，需要初始化管理员")
		return
	}
	utils.Info("检测到 %d 个用户", count)
}

func GenerateToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func HashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		utils.Error("密码加密失败: %v", err)
		os.Exit(1)
	}
	return string(hash)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
