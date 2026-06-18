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
	DB.Exec("SET FOREIGN_KEY_CHECKS = 0")
	err := DB.AutoMigrate(
		&models.User{},
		&models.Article{},
		&models.Comment{},
		&models.CommentLike{},
		&models.CommentReplyNotification{},
		&models.ViewHistory{},
		&models.Like{},
		&models.Favorite{},
		&models.Title{},
		&models.UserTitle{},
		&models.Category{},
		&models.DeletionRequest{},
		&models.SidebarConfig{},
		&models.Announcement{},
		&models.Notification{},
		&models.UserNotification{},
		&models.Friend{},
		&models.FriendRequest{},
		&models.SiteConfig{},
		&models.Report{},
		&models.PersonalNotification{},
		&models.PermissionGroup{},
		&models.UserPermissionGroup{},
		&models.SystemLog{},
		&models.SignInRecord{},
		&models.SignInConfig{},
		&models.CoinRecord{},
	)
	DB.Exec("SET FOREIGN_KEY_CHECKS = 1")
	if err != nil {
		utils.Error("数据库迁移失败: %v", err)
		os.Exit(1)
	}
	utils.Success("数据库迁移完成")

	// 删除旧关注系统相关表
	DropOldFollowTables()
}

// DropOldFollowTables 删除旧的关注系统表
func DropOldFollowTables() {
	oldTables := []string{"follows", "follow_notifications"}
	for _, table := range oldTables {
		if DB.Migrator().HasTable(table) {
			DB.Exec("SET FOREIGN_KEY_CHECKS = 0")
			DB.Migrator().DropTable(table)
			utils.Info("已删除旧表: %s", table)
			DB.Exec("SET FOREIGN_KEY_CHECKS = 1")
		}
	}
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

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		utils.Error("密码加密失败: %v", err)
		return "", err
	}
	return string(hash), nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
