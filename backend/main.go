package main

import (
	"embed"
	"encoding/json"
	"forum/controllers"
	"forum/database"
	initpkg "forum/init"
	"forum/middleware"
	"forum/service"
	"forum/utils"
	"io/fs"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type Config struct {
	Port      string         `json:"port"`
	JWTSecret string         `json:"jwt_secret"`
	WebDAV    WebDAVConfig   `json:"webdav"`
	Database  DatabaseConfig `json:"database"`
}

type WebDAVConfig struct {
	URL      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
}

var config Config

//go:embed dist
var distFS embed.FS

func main() {
	rand.Seed(time.Now().UnixNano())

	loadConfig()

	controllers.InitJWT(config.JWTSecret)

	utils.InitWebDAV(config.WebDAV.URL, config.WebDAV.Username, config.WebDAV.Password)

	dsn := config.Database.Username + ":" + config.Database.Password + "@tcp(" + config.Database.Host + ":" + config.Database.Port + ")/" + config.Database.DBName + "?charset=utf8mb4&parseTime=True&loc=Local"
	database.InitDB(dsn)

	service.InitServices()

	database.AutoMigrate()
	database.CheckAndInitAdmin()

	initpkg.SystemInit()

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())

	r.Use(middleware.RateLimit(100, time.Second))

	fs, err := fs.Sub(distFS, "dist")
	if err != nil {
		utils.Error("无法创建子文件系统: %v", err)
		os.Exit(1)
	}

	r.NoRoute(func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.JSON(404, gin.H{
				"error":   "404 Not Found",
				"message": "请求的接口不存在",
			})
			return
		}

		if strings.HasPrefix(c.Request.URL.Path, "/assets/") ||
			strings.HasPrefix(c.Request.URL.Path, "/icons/") ||
			strings.HasSuffix(c.Request.URL.Path, ".js") ||
			strings.HasSuffix(c.Request.URL.Path, ".css") ||
			strings.HasSuffix(c.Request.URL.Path, ".svg") ||
			strings.HasSuffix(c.Request.URL.Path, ".png") ||
			strings.HasSuffix(c.Request.URL.Path, ".ico") {

			filePath := c.Request.URL.Path
			if filePath[0] == '/' {
				filePath = filePath[1:]
			}

			file, err := fs.Open(filePath)
			if err != nil {
				c.Status(http.StatusNotFound)
				return
			}
			defer file.Close()

			stat, err := file.Stat()
			if err != nil {
				c.Status(http.StatusInternalServerError)
				return
			}

			c.DataFromReader(http.StatusOK, stat.Size(), getContentType(filePath), file, nil)
			return
		}

		indexFile, err := fs.Open("index.html")
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		defer indexFile.Close()

		stat, err := indexFile.Stat()
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		c.DataFromReader(http.StatusOK, stat.Size(), "text/html; charset=utf-8", indexFile, nil)
	})

	r.MaxMultipartMemory = 20480 << 20

	r.Any("/proxy/webdav/*path", utils.ProxyWebDAVHandler)

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", middleware.RateLimit(5, time.Minute), controllers.Register)
			auth.POST("/login", middleware.RateLimit(10, time.Minute), controllers.Login)
			auth.POST("/refresh-token", controllers.RefreshToken)
			auth.POST("/init-admin", controllers.InitAdmin)
			auth.GET("/check-init", controllers.CheckInit)
		}

		protected := api.Group("/")
		protected.Use(middleware.Auth(config.JWTSecret))
		{
			protected.GET("/profile", controllers.GetProfile)
			protected.PUT("/profile", controllers.UpdateProfile)
			protected.POST("/upload/avatar", controllers.UploadAvatar)
			protected.POST("/upload/image", controllers.UploadImage)
			protected.POST("/upload/video", controllers.UploadVideo)
			protected.POST("/upload/voice", controllers.UploadVoice)

			protected.POST("/articles", controllers.CreateArticle)
			protected.PUT("/articles/:id", controllers.UpdateArticle)
			protected.DELETE("/articles/:id", controllers.DeleteArticle)
			protected.POST("/articles/:id/restore", middleware.RequireMinLevel(80), controllers.RestoreArticle)
			protected.POST("/articles/:id/like", controllers.LikeArticle)
			protected.DELETE("/articles/:id/like", controllers.UnlikeArticle)
			protected.POST("/articles/:id/coin", controllers.CoinArticle)
			protected.GET("/my/articles", controllers.GetMyArticles)
			protected.POST("/articles/:id/share", controllers.ShareArticle)
			protected.GET("/my/drafts", controllers.GetDraftArticles)
			protected.POST("/articles/:id/publish", controllers.PublishDraft)
			protected.POST("/articles/:id/pin", middleware.RequireMinLevel(80), controllers.PinArticle)
			protected.DELETE("/articles/:id/pin", middleware.RequireMinLevel(80), controllers.UnpinArticle)

			protected.POST("/reports", controllers.CreateReport)
			protected.GET("/reports", middleware.RequireMinLevel(80), controllers.GetReports)
			protected.GET("/reports/:id", middleware.RequireMinLevel(80), controllers.GetReport)
			protected.PUT("/reports/:id/handle", middleware.RequireMinLevel(80), controllers.HandleReport)

			protected.POST("/articles/:id/comments", controllers.CreateComment)
			protected.DELETE("/comments/:id", controllers.DeleteComment)
			protected.POST("/comments/:id/like", controllers.LikeComment)
			protected.DELETE("/comments/:id/like", controllers.UnlikeComment)

			protected.POST("/categories", middleware.RequireMinLevel(80), controllers.CreateCategory)
			protected.PUT("/categories/:id", middleware.RequireMinLevel(80), controllers.UpdateCategory)
			protected.DELETE("/categories/:id", middleware.RequireMinLevel(80), controllers.DeleteCategory)

			protected.PUT("/sidebar-config", middleware.RequireMinLevel(80), controllers.UpdateSidebarConfig)

			protected.GET("/deletion-requests", middleware.RequireMinLevel(80), controllers.GetDeletionRequests)
			protected.POST("/deletion-requests/:id/approve", middleware.RequireMinLevel(80), controllers.ApproveDeletion)
			protected.POST("/deletion-requests/:id/reject", middleware.RequireMinLevel(80), controllers.RejectDeletion)

			protected.PUT("/announcement", middleware.RequireMinLevel(80), controllers.UpdateAnnouncement)

			protected.PUT("/site-config", middleware.RequireSystemAdmin(), controllers.UpdateSiteConfig)
			protected.POST("/site-config/test-smtp", middleware.RequireSystemAdmin(), controllers.TestSMTPConfig)

			protected.GET("/notifications", controllers.GetNotifications)
			protected.GET("/notifications/unread-count", controllers.GetUnreadCount)
			protected.POST("/notifications/:id/read", controllers.MarkNotificationRead)
			protected.POST("/notifications/read-all", controllers.MarkAllNotificationsRead)
			protected.POST("/notifications", middleware.RequireMinLevel(80), controllers.CreateNotification)
			protected.GET("/notifications/admin", middleware.RequireMinLevel(80), controllers.GetAdminNotifications)
			protected.DELETE("/notifications/:id", middleware.RequireMinLevel(80), controllers.DeleteNotification)

			protected.GET("/comment-reply-notifications", controllers.GetCommentReplyNotifications)
			protected.POST("/comment-reply-notifications/:id/read", controllers.MarkCommentReplyNotificationRead)
			protected.POST("/comment-reply-notifications/read-all", controllers.MarkAllCommentReplyNotificationsRead)

			protected.POST("/user-notifications/send", middleware.RequireMinLevel(80), controllers.SendUserNotification)
			protected.POST("/user-notifications/send-batch", middleware.RequireMinLevel(80), controllers.SendBatchNotifications)
			protected.GET("/user-notifications", controllers.GetUserNotifications)
			protected.GET("/user-notifications/:id", controllers.GetNotification)
			protected.POST("/user-notifications/:id/read", controllers.MarkNotificationAsRead)
			protected.POST("/user-notifications/read-all", controllers.MarkAllNotificationsAsRead)
			protected.DELETE("/user-notifications/:id", controllers.DeletePersonalNotification)
			protected.DELETE("/user-notifications/clear", controllers.ClearAllNotifications)
			protected.GET("/admin/user-notifications/:user_id", middleware.RequireMinLevel(80), controllers.AdminGetUserNotifications)

			protected.GET("/permission-groups", controllers.GetPermissionGroups)
			protected.GET("/permission-groups/:id", controllers.GetPermissionGroup)
			protected.POST("/permission-groups", middleware.RequireMinLevel(80), controllers.CreatePermissionGroup)
			protected.PUT("/permission-groups/:id", middleware.RequireMinLevel(80), controllers.UpdatePermissionGroup)
			protected.DELETE("/permission-groups/:id", middleware.RequireMinLevel(80), controllers.DeletePermissionGroup)
			protected.POST("/permission-groups/grant", middleware.RequireMinLevel(80), controllers.GrantPermissionGroup)
			protected.DELETE("/permission-groups/:id/revoke-user/:user_id", middleware.RequireMinLevel(80), controllers.RevokePermissionGroup)
			protected.GET("/users/:id/permission-groups", controllers.GetUserPermissionGroups)
			protected.GET("/permissions/check", controllers.CheckUserPermissions)
			protected.POST("/permission-groups/init", middleware.RequireSystemAdmin(), controllers.InitializeDefaultPermissionGroups)

			protected.POST("/friends/request", controllers.SendFriendRequest)
			protected.POST("/friends/request/:id/accept", controllers.AcceptFriendRequest)
			protected.POST("/friends/request/:id/reject", controllers.RejectFriendRequest)
			protected.DELETE("/friends/:id", controllers.DeleteFriend)
			protected.GET("/friends", controllers.GetFriendList)
			protected.GET("/friends/requests", controllers.GetFriendRequests)
			protected.GET("/friends/sent-requests", controllers.GetSentFriendRequests)
			protected.PUT("/friends/:id/display-name", controllers.UpdateFriendDisplayName)
			protected.GET("/friends/status/:id", controllers.CheckFriendStatus)
			protected.GET("/friends/mutual/:id", controllers.GetMutualFriends)

			protected.POST("/articles/:id/favorite", controllers.AddFavorite)
			protected.DELETE("/articles/:id/favorite", controllers.RemoveFavorite)
			protected.GET("/favorites", controllers.GetFavorites)
			protected.GET("/articles/:id/favorite/check", controllers.CheckFavorite)

			protected.POST("/signin", controllers.SignIn)
			protected.GET("/signin/status", controllers.GetSignInStatus)
			protected.GET("/signin/history", controllers.GetSignInHistory)
			protected.GET("/signin/rankings", controllers.GetSignInRankings)
			protected.GET("/signin/config", controllers.GetSignInConfig)
			protected.PUT("/signin/config", middleware.RequireMinLevel(80), controllers.UpdateSignInConfig)

			protected.GET("/admin/check", controllers.CheckAdmin)
			protected.GET("/admin/statistics", middleware.RequireMinLevel(80), controllers.GetStatistics)
			protected.GET("/admin/users", middleware.RequireMinLevel(80), controllers.GetAllUsers)
			protected.PUT("/admin/users/:id", middleware.RequireMinLevel(80), controllers.UpdateUser)
			protected.PUT("/admin/users/:id/role", middleware.RequireMinLevel(80), controllers.UpdateUserRole)
			protected.POST("/admin/users/:id/ban", middleware.RequireMinLevel(80), controllers.BanUser)
			protected.POST("/admin/users/:id/unban", middleware.RequireMinLevel(80), controllers.UnbanUser)
			protected.DELETE("/admin/users/:id", middleware.RequireMinLevel(80), controllers.DeleteUser)
			protected.GET("/admin/articles", middleware.RequireMinLevel(80), controllers.GetAllArticles)
			protected.PUT("/admin/articles/:id/status", middleware.RequireMinLevel(80), controllers.UpdateArticleStatus)
			protected.GET("/admin/comments", middleware.RequireMinLevel(80), controllers.GetAllComments)
			protected.DELETE("/admin/comments/:id", middleware.RequireMinLevel(80), controllers.DeleteCommentAdmin)

			protected.POST("/user/status", controllers.UpdateUserStatus)
			protected.GET("/user/status/:id", controllers.GetUserStatus)
			protected.GET("/users/status", middleware.RequireMinLevel(80), controllers.GetAllUserStatuses)
			protected.GET("/users/online", middleware.RequireMinLevel(80), controllers.GetOnlineUsers)
			protected.POST("/users/status/cleanup", middleware.RequireMinLevel(80), controllers.CleanupUserStatuses)

			protected.GET("/titles", controllers.GetAllTitles)
			protected.POST("/titles", middleware.RequireMinLevel(80), controllers.CreateTitle)
			protected.PUT("/titles/:id", middleware.RequireMinLevel(80), controllers.UpdateTitle)
			protected.DELETE("/titles/:id", middleware.RequireMinLevel(80), controllers.DeleteTitle)
			protected.POST("/titles/grant", middleware.RequireMinLevel(80), controllers.GrantTitle)
			protected.POST("/titles/revoke", middleware.RequireMinLevel(80), controllers.RevokeTitle)
			protected.GET("/users/:id/titles", controllers.GetUserTitles)

			protected.GET("/system-logs", middleware.RequireSystemAdmin(), controllers.GetSystemLogs)
			protected.GET("/system-logs/modules", middleware.RequireSystemAdmin(), controllers.GetLogModules)
			protected.DELETE("/system-logs/old", middleware.RequireSystemAdmin(), controllers.DeleteOldLogs)
			protected.GET("/my-logs", controllers.GetMyLogs)

			protected.GET("/level", controllers.GetUserLevel)
			protected.GET("/level/experience-records", controllers.GetUserExperienceRecords)
			protected.GET("/achievements", controllers.GetUserAchievements)
			protected.POST("/level/config", middleware.RequireMinLevel(80), controllers.CreateLevelConfig)
			protected.PUT("/level/config/:id", middleware.RequireMinLevel(80), controllers.UpdateLevelConfig)
			protected.POST("/achievements", middleware.RequireMinLevel(80), controllers.CreateAchievement)
			protected.PUT("/achievements/:id", middleware.RequireMinLevel(80), controllers.UpdateAchievement)
			protected.DELETE("/achievements/:id", middleware.RequireMinLevel(80), controllers.DeleteAchievement)

			protected.GET("/statistics", controllers.GetUserStatistics)
			protected.GET("/statistics/daily", controllers.GetDailyStatistics)
			protected.GET("/statistics/overview", controllers.GetSystemOverview)
			protected.GET("/statistics/activity", controllers.GetUserActivity)
			protected.GET("/statistics/dashboard", middleware.RequireMinLevel(80), controllers.GetStatisticsDashboard)
			protected.GET("/articles/:id/statistics", controllers.GetArticleStatistics)

			protected.GET("/collections", controllers.GetCollections)
			protected.GET("/collections/:id", controllers.GetCollection)
			protected.POST("/collections", controllers.CreateCollection)
			protected.PUT("/collections/:id", controllers.UpdateCollection)
			protected.DELETE("/collections/:id", controllers.DeleteCollection)
			protected.POST("/collections/:id/articles", controllers.AddArticleToCollection)
			protected.DELETE("/collections/:id/articles/:article_id", controllers.RemoveArticleFromCollection)
			protected.GET("/articles/:id/versions", controllers.GetArticleVersions)
			protected.GET("/articles/:id/versions/:version", controllers.GetArticleVersion)
			protected.POST("/articles/:id/versions/:version/restore", controllers.RestoreArticleVersion)

			protected.POST("/topics", middleware.RequireMinLevel(80), controllers.CreateTopic)
			protected.PUT("/topics/:id", middleware.RequireMinLevel(80), controllers.UpdateTopic)
			protected.DELETE("/topics/:id", middleware.RequireMinLevel(80), controllers.DeleteTopic)
			protected.POST("/topics/:id/follow", controllers.FollowTopic)
			protected.DELETE("/topics/:id/follow", controllers.UnfollowTopic)
			protected.GET("/topics/followed", controllers.GetFollowedTopics)
			protected.POST("/articles/:id/topics", controllers.AddTopicToArticle)
			protected.DELETE("/articles/:id/topics/:topic_id", controllers.RemoveTopicFromArticle)

			protected.GET("/leaderboard/rank", controllers.GetUserRank)
			protected.GET("/badges", controllers.GetUserBadges)
			protected.PUT("/badges/:id/display", controllers.UpdateBadgeDisplay)
			protected.POST("/badges/grant", middleware.RequireMinLevel(80), controllers.GrantBadge)
			protected.DELETE("/badges/:id", middleware.RequireMinLevel(80), controllers.RevokeBadge)
		}

		api.GET("/articles", controllers.GetArticles)
		api.GET("/articles/:id", controllers.GetArticle)
		api.GET("/articles/search", controllers.SearchArticles)
		api.GET("/categories", controllers.GetCategories)
		api.GET("/announcement", controllers.GetAnnouncement)
		api.GET("/sidebar-config", controllers.GetSidebarConfig)
		api.GET("/site-config", controllers.GetSiteConfig)
		api.GET("/version", controllers.GetVersion)

		api.GET("/users/:id", controllers.GetUserByID)
		api.GET("/users/:id/articles", controllers.GetUserArticles)

		api.GET("/topics", controllers.GetTopics)
		api.GET("/topics/:id", controllers.GetTopic)
		api.GET("/topics/hot", controllers.GetHotTopics)
		api.GET("/leaderboard", controllers.GetLeaderboard)
		api.GET("/achievements/all", controllers.GetAllAchievements)
		api.GET("/level/config", controllers.GetLevelConfig)
	}

	passwordReset := api.Group("/password")
	{
		passwordReset.POST("/reset-code", middleware.RateLimit(3, time.Minute), controllers.SendResetCode)
		passwordReset.POST("/reset", middleware.RateLimit(5, time.Minute), controllers.ResetPassword)
	}

	utils.Section("服务启动")
	utils.Info("论坛服务器地址: http://0.0.0.0:%s", config.Port)
	utils.Info("WebSocket地址: ws://0.0.0.0:%s/ws", config.Port)
	utils.Section("")
	r.Run(":" + config.Port)

	closeChan := make(chan struct{})
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sigChan
		signal.Stop(sigChan)
		close(closeChan)
	}()

	<-closeChan
}

func loadConfig() {
	_, err := os.Stat("config.json")
	if os.IsNotExist(err) {
		utils.Info("配置文件不存在，正在生成默认配置...")
		err = generateDefaultConfig()
		if err != nil {
			utils.Error("生成默认配置文件失败: %v", err)
			os.Exit(1)
		}
		utils.Success("默认配置文件已生成: config.json")
	}

	file, err := os.Open("config.json")
	if err != nil {
		utils.Error("无法加载配置文件: %v", err)
		os.Exit(1)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		utils.Error("解析配置文件失败: %v", err)
		os.Exit(1)
	}
	utils.Success("配置文件加载成功")
}

func generateDefaultConfig() error {
	defaultConfig := Config{
		Port:      "3620",
		JWTSecret: generateRandomString(32),
		WebDAV: WebDAVConfig{
			URL:      "",
			Username: "",
			Password: "",
		},
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     "3306",
			Username: "root",
			Password: "",
			DBName:   "forum",
		},
	}

	file, err := os.Create("config.json")
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	return encoder.Encode(defaultConfig)
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

func getContentType(filePath string) string {
	switch {
	case strings.HasSuffix(filePath, ".js"):
		return "application/javascript; charset=utf-8"
	case strings.HasSuffix(filePath, ".css"):
		return "text/css; charset=utf-8"
	case strings.HasSuffix(filePath, ".html"):
		return "text/html; charset=utf-8"
	case strings.HasSuffix(filePath, ".svg"):
		return "image/svg+xml"
	case strings.HasSuffix(filePath, ".png"):
		return "image/png"
	case strings.HasSuffix(filePath, ".jpg") || strings.HasSuffix(filePath, ".jpeg"):
		return "image/jpeg"
	case strings.HasSuffix(filePath, ".gif"):
		return "image/gif"
	case strings.HasSuffix(filePath, ".ico"):
		return "image/x-icon"
	default:
		return "application/octet-stream"
	}
}
