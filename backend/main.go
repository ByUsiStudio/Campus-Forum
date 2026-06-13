package main

import (
	"encoding/json"
	"forum/controllers"
	"forum/database"
	"forum/im"
	initpkg "forum/init"
	"forum/middleware"
	"forum/service"
	"forum/utils"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type Config struct {
	Port      string         `json:"port"`
	JWTSecret string         `json:"jwt_secret"`
	WebDAV    WebDAVConfig   `json:"webdav"`
	Database  DatabaseConfig `json:"database"`
	IM        IMConfig       `json:"im"`
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

type IMConfig struct {
	Enabled        bool   `json:"enabled"`
	APIURL         string `json:"api_url"`    // 野火IM API地址
	AppID          string `json:"app_id"`     // 应用ID
	AppSecret      string `json:"app_secret"` // 应用密钥
	IMHost         string `json:"im_host"`    // IM服务器地址
	IMPort         int    `json:"im_port"`    // IM服务器端口
	ApiGatewayPort int    `json:"api_gateway_port"`
	WsPort         int    `json:"ws_port"`
	AdminPort      int    `json:"admin_port"`
}

var config Config

func main() {
	// 加载配置
	loadConfig()

	// 初始化JWT
	controllers.InitJWT(config.JWTSecret)

	// 初始化WebDAV
	utils.InitWebDAV(config.WebDAV.URL, config.WebDAV.Username, config.WebDAV.Password)

	// 初始化数据库
	dsn := config.Database.Username + ":" + config.Database.Password + "@tcp(" + config.Database.Host + ":" + config.Database.Port + ")/" + config.Database.DBName + "?charset=utf8mb4&parseTime=True&loc=Local"
	database.InitDB(dsn)

	// 初始化服务
	service.InitServices()

	// 自动迁移和初始化
	database.AutoMigrate()
	database.CheckAndInitAdmin()

	// 系统初始化（权限组、头衔等）
	initpkg.SystemInit()

	// 设置路由
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())

	// 全局速率限制（每秒100次请求）
	r.Use(middleware.RateLimit(100, time.Second))

	// 处理404请求
	r.NoRoute(func(c *gin.Context) {
		// 对 OPTIONS 请求返回 204 让 CORS 中间件处理
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.JSON(404, gin.H{
			"error":   "404 Not Found",
			"message": "请求的接口不存在",
		})
	})

	// 设置最大请求体大小为 20G（适用于视频上传）
	r.MaxMultipartMemory = 20480 << 20 // 20G

	// WebDAV代理路由
	r.Any("/proxy/webdav/*path", utils.ProxyWebDAVHandler)

	// WebSocket路由（IM服务）
	imCtrl := im.NewController()
	r.GET("/ws", imCtrl.WebSocketHandle)

	// API路由
	api := r.Group("/api")
	{
		// 公开路由
		auth := api.Group("/auth")
		{
			auth.POST("/register", middleware.RateLimit(5, time.Minute), controllers.Register)
			auth.POST("/login", middleware.RateLimit(10, time.Minute), controllers.Login)
			auth.POST("/init-admin", controllers.InitAdmin)
			auth.GET("/check-init", controllers.CheckInit)
		}

		// 需要认证的路由
		protected := api.Group("/")
		protected.Use(middleware.Auth(config.JWTSecret))
		{
			protected.GET("/profile", controllers.GetProfile)
			protected.PUT("/profile", controllers.UpdateProfile)
			protected.POST("/upload/avatar", controllers.UploadAvatar)
			protected.POST("/upload/image", controllers.UploadImage)
			protected.POST("/upload/video", controllers.UploadVideo)
			protected.POST("/upload/voice", controllers.UploadVoice)

			// 文章相关
			protected.POST("/articles", controllers.CreateArticle)
			protected.PUT("/articles/:id", controllers.UpdateArticle)
			protected.DELETE("/articles/:id", controllers.DeleteArticle)
			protected.POST("/articles/:id/restore", middleware.AdminOnly(), controllers.RestoreArticle)
			protected.POST("/articles/:id/like", controllers.LikeArticle)
			protected.DELETE("/articles/:id/like", controllers.UnlikeArticle)
			protected.GET("/my/articles", controllers.GetMyArticles)
			protected.POST("/articles/:id/share", controllers.ShareArticle)
			protected.GET("/my/drafts", controllers.GetDraftArticles)
			protected.POST("/articles/:id/publish", controllers.PublishDraft)
			protected.POST("/articles/:id/pin", middleware.AdminOnly(), controllers.PinArticle)
			protected.DELETE("/articles/:id/pin", middleware.AdminOnly(), controllers.UnpinArticle)

			// 举报相关
			protected.POST("/reports", controllers.CreateReport)
			protected.GET("/reports", middleware.AdminOnly(), controllers.GetReports)
			protected.GET("/reports/:id", middleware.AdminOnly(), controllers.GetReport)
			protected.PUT("/reports/:id/handle", middleware.AdminOnly(), controllers.HandleReport)

			// 评论相关
			protected.POST("/articles/:id/comments", controllers.CreateComment)
			protected.DELETE("/comments/:id", controllers.DeleteComment)
			protected.POST("/comments/:id/like", controllers.LikeComment)
			protected.DELETE("/comments/:id/like", controllers.UnlikeComment)

			// 分区相关
			protected.POST("/categories", middleware.AdminOnly(), controllers.CreateCategory)
			protected.PUT("/categories/:id", middleware.AdminOnly(), controllers.UpdateCategory)
			protected.DELETE("/categories/:id", middleware.AdminOnly(), controllers.DeleteCategory)

			// 侧边栏配置
			protected.PUT("/sidebar-config", middleware.AdminOnly(), controllers.UpdateSidebarConfig)

			// 删除审核
			protected.GET("/deletion-requests", middleware.AdminOnly(), controllers.GetDeletionRequests)
			protected.POST("/deletion-requests/:id/approve", middleware.AdminOnly(), controllers.ApproveDeletion)
			protected.POST("/deletion-requests/:id/reject", middleware.AdminOnly(), controllers.RejectDeletion)

			// 公告
			protected.PUT("/announcement", middleware.AdminOnly(), controllers.UpdateAnnouncement)

			// 网站配置
			protected.PUT("/site-config", middleware.AdminOnly(), controllers.UpdateSiteConfig)
			protected.POST("/site-config/test-smtp", middleware.AdminOnly(), controllers.TestSMTPConfig)

			// 通知相关
			protected.GET("/notifications", controllers.GetNotifications)
			protected.GET("/notifications/unread-count", controllers.GetUnreadCount)
			protected.POST("/notifications/:id/read", controllers.MarkNotificationRead)
			protected.POST("/notifications/read-all", controllers.MarkAllNotificationsRead)
			protected.POST("/notifications", middleware.AdminOnly(), controllers.CreateNotification)
			protected.GET("/notifications/admin", middleware.AdminOnly(), controllers.GetAdminNotifications)
			protected.DELETE("/notifications/:id", middleware.AdminOnly(), controllers.DeleteNotification)

			// 评论回复通知
			protected.GET("/comment-reply-notifications", controllers.GetCommentReplyNotifications)
			protected.POST("/comment-reply-notifications/:id/read", controllers.MarkCommentReplyNotificationRead)
			protected.POST("/comment-reply-notifications/read-all", controllers.MarkAllCommentReplyNotificationsRead)

			// 用户个人通知（单独通知）
			protected.POST("/user-notifications/send", middleware.AdminOnly(), controllers.SendUserNotification)
			protected.POST("/user-notifications/send-batch", middleware.AdminOnly(), controllers.SendBatchNotifications)
			protected.GET("/user-notifications", controllers.GetUserNotifications)
			protected.GET("/user-notifications/:id", controllers.GetNotification)
			protected.POST("/user-notifications/:id/read", controllers.MarkNotificationAsRead)
			protected.POST("/user-notifications/read-all", controllers.MarkAllNotificationsAsRead)
			protected.DELETE("/user-notifications/:id", controllers.DeletePersonalNotification)
			protected.DELETE("/user-notifications/clear", controllers.ClearAllNotifications)
			protected.GET("/admin/user-notifications/:user_id", middleware.AdminOnly(), controllers.AdminGetUserNotifications)

			// IM（goim）相关
			protected.GET("/im/online-status", imCtrl.GetOnlineStatus)
			protected.GET("/im/online-users", imCtrl.GetOnlineUsers)
			protected.POST("/im/send-message", imCtrl.SendMessage)

			// 权限组管理
			protected.GET("/permission-groups", controllers.GetPermissionGroups)
			protected.GET("/permission-groups/:id", controllers.GetPermissionGroup)
			protected.POST("/permission-groups", middleware.AdminOnly(), controllers.CreatePermissionGroup)
			protected.PUT("/permission-groups/:id", middleware.AdminOnly(), controllers.UpdatePermissionGroup)
			protected.DELETE("/permission-groups/:id", middleware.AdminOnly(), controllers.DeletePermissionGroup)
			protected.POST("/permission-groups/grant", middleware.AdminOnly(), controllers.GrantPermissionGroup)
			protected.DELETE("/permission-groups/:id/revoke-user/:user_id", middleware.AdminOnly(), controllers.RevokePermissionGroup)
			protected.GET("/users/:id/permission-groups", controllers.GetUserPermissionGroups)
			protected.GET("/permissions/check", controllers.CheckUserPermissions)
			protected.POST("/permission-groups/init", middleware.AdminOnly(), controllers.InitializeDefaultPermissionGroups)

			// 好友相关（原关注系统改为好友系统）
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

			// 收藏相关
			protected.POST("/articles/:id/favorite", controllers.AddFavorite)
			protected.DELETE("/articles/:id/favorite", controllers.RemoveFavorite)
			protected.GET("/favorites", controllers.GetFavorites)
			protected.GET("/articles/:id/favorite/check", controllers.CheckFavorite)

			// 签到相关
			protected.POST("/signin", controllers.SignIn)
			protected.GET("/signin/status", controllers.GetSignInStatus)
			protected.GET("/signin/history", controllers.GetSignInHistory)
			protected.GET("/signin/rankings", controllers.GetSignInRankings)
			protected.GET("/signin/config", controllers.GetSignInConfig)
			protected.PUT("/signin/config", middleware.AdminOnly(), controllers.UpdateSignInConfig)

			// 后台管理
			protected.GET("/admin/check", controllers.CheckAdmin)
			protected.GET("/admin/statistics", middleware.AdminOnly(), controllers.GetStatistics)
			protected.GET("/admin/users", middleware.AdminOnly(), controllers.GetAllUsers)
			protected.PUT("/admin/users/:id", middleware.AdminOnly(), controllers.UpdateUser)
			protected.PUT("/admin/users/:id/role", middleware.AdminOnly(), controllers.UpdateUserRole)
			protected.POST("/admin/users/:id/ban", middleware.AdminOnly(), controllers.BanUser)
			protected.POST("/admin/users/:id/unban", middleware.AdminOnly(), controllers.UnbanUser)
			protected.DELETE("/admin/users/:id", middleware.AdminOnly(), controllers.DeleteUser)
			protected.GET("/admin/articles", middleware.AdminOnly(), controllers.GetAllArticles)
			protected.PUT("/admin/articles/:id/status", middleware.AdminOnly(), controllers.UpdateArticleStatus)
			protected.GET("/admin/comments", middleware.AdminOnly(), controllers.GetAllComments)
			protected.DELETE("/admin/comments/:id", middleware.AdminOnly(), controllers.DeleteCommentAdmin)

			// 用户在线状态管理
			protected.POST("/user/status", controllers.UpdateUserStatus)
			protected.GET("/user/status/:id", controllers.GetUserStatus)
			protected.GET("/users/status", middleware.AdminOnly(), controllers.GetAllUserStatuses)
			protected.GET("/users/online", middleware.AdminOnly(), controllers.GetOnlineUsers)
			protected.POST("/users/status/cleanup", middleware.AdminOnly(), controllers.CleanupUserStatuses)

			// 头衔管理
			protected.GET("/titles", controllers.GetAllTitles)
			protected.POST("/titles", middleware.AdminOnly(), controllers.CreateTitle)
			protected.PUT("/titles/:id", middleware.AdminOnly(), controllers.UpdateTitle)
			protected.DELETE("/titles/:id", middleware.AdminOnly(), controllers.DeleteTitle)
			protected.POST("/titles/grant", middleware.AdminOnly(), controllers.GrantTitle)
			protected.POST("/titles/revoke", middleware.AdminOnly(), controllers.RevokeTitle)
			protected.GET("/users/:id/titles", controllers.GetUserTitles)

			// 系统日志
			protected.GET("/system-logs", middleware.AdminOnly(), controllers.GetSystemLogs)
			protected.GET("/system-logs/modules", middleware.AdminOnly(), controllers.GetLogModules)
			protected.DELETE("/system-logs/old", middleware.AdminOnly(), controllers.DeleteOldLogs)
			protected.GET("/my-logs", controllers.GetMyLogs)
		}

		// 公开路由
		api.GET("/articles", controllers.GetArticles)
		api.GET("/articles/:id", controllers.GetArticle)
		api.GET("/articles/search", controllers.SearchArticles)
		api.GET("/categories", controllers.GetCategories)
		api.GET("/announcement", controllers.GetAnnouncement)
		api.GET("/sidebar-config", controllers.GetSidebarConfig)
		api.GET("/site-config", controllers.GetSiteConfig)
		api.GET("/version", controllers.GetVersion)

		// 用户公开信息
		api.GET("/users/:id", controllers.GetUserByID)
		api.GET("/users/:id/articles", controllers.GetUserArticles)
	}

	// 密码重置（公开接口，但有限制）
	passwordReset := api.Group("/password")
	{
		passwordReset.POST("/reset-code", middleware.RateLimit(3, time.Minute), controllers.SendResetCode)
		passwordReset.POST("/reset", middleware.RateLimit(5, time.Minute), controllers.ResetPassword)
	}

	// 启动服务
	utils.Section("服务启动")
	utils.Info("论坛服务器地址: http://0.0.0.0:%s", config.Port)
	utils.Info("WebSocket地址: ws://0.0.0.0:%s/ws", config.Port)
	utils.Section("")
	r.Run(":" + config.Port)

	// 等待退出信号
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
