package database

import (
    "crypto/rand"
    "encoding/hex"
    "forum/models"
    "log"

    "golang.org/x/crypto/bcrypt"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(dsn string) {
    var err error
    DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("数据库连接失败:", err)
    }
    log.Println("数据库连接成功")
}

func AutoMigrate() {
    err := DB.AutoMigrate(
        &models.User{},
        &models.Article{},
        &models.Comment{},
        &models.Like{},
        &models.Category{},
        &models.DeletionRequest{},
        &models.SidebarConfig{},
        &models.Announcement{},
    )
    if err != nil {
        log.Fatal("数据库迁移失败:", err)
    }
    log.Println("数据库迁移完成")
}

func CheckAndInitAdmin() {
    var count int64
    DB.Model(&models.User{}).Count(&count)
    
    if count == 0 {
        log.Println("未检测到用户，需要初始化管理员")
        // 不自动创建，等待首次访问时通过API创建
        return
    }
}

func GenerateToken() string {
    bytes := make([]byte, 32)
    rand.Read(bytes)
    return hex.EncodeToString(bytes)
}

func HashPassword(password string) string {
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        log.Fatal("密码加密失败:", err)
    }
    return string(hash)
}

func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}